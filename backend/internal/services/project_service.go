package services

import (
	"errors"
	"fmt"
	"log"
	"webtest/internal/constants"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"gorm.io/gorm"
)

// 错误定义
var (
	ErrProjectNameExists    = errors.New("项目名已存在")
	ErrProjectNotFound      = errors.New("项目不存在")
	ErrPermissionDenied     = errors.New("无权限访问此项目")
	ErrSelfRemovalForbidden = errors.New("当前登录管理员不可移出项目")
	ErrInvalidUserIDs       = errors.New("存在无效的用户ID")
	ErrRoleMismatch         = errors.New("用户角色不匹配")
)

// MemberInfo 成员信息
type MemberInfo struct {
	UserID        uint   `json:"user_id"`
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
	IsCurrentUser bool   `json:"is_current_user"`
}

// ProjectMembersResponse 项目成员响应
type ProjectMembersResponse struct {
	ProjectID   uint         `json:"project_id"`
	ProjectName string       `json:"project_name"`
	Managers    []MemberInfo `json:"managers"`
	Members     []MemberInfo `json:"members"`
}

// ProjectService 项目服务接口
type ProjectService interface {
	GetUserProjects(userID uint, role string) ([]models.Project, error)
	IsProjectMember(projectID uint, userID uint) (bool, error)
	CreateProject(name string, description string, creatorID uint) (*models.Project, error)
	UpdateProject(projectID uint, newName string, userID uint, role string) (*models.Project, error)
	UpdateProjectMetadata(projectID uint, updates map[string]interface{}, userID uint, role string) (*models.Project, error)
	DeleteProject(projectID uint, userID uint, role string) error
	GetByID(projectID uint, userID uint) (*models.Project, string, error)
	GetProjectMembers(projectID uint, currentUserID uint) (*ProjectMembersResponse, error)
	UpdateProjectMembers(projectID uint, managers []uint, members []uint, currentUserID uint) (*ProjectMembersResponse, error)
}

// projectService 项目服务实现
type projectService struct {
	projectRepo repositories.ProjectRepository
	memberRepo  repositories.ProjectMemberRepository
	userRepo    repositories.UserRepository
	db          *gorm.DB
}

// NewProjectService 创建项目服务实例
func NewProjectService(
	projectRepo repositories.ProjectRepository,
	memberRepo repositories.ProjectMemberRepository,
	userRepo repositories.UserRepository,
	db *gorm.DB,
) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		memberRepo:  memberRepo,
		userRepo:    userRepo,
		db:          db,
	}
}

// GetUserProjects 获取用户可访问的项目列表
func (s *projectService) GetUserProjects(userID uint, role string) ([]models.Project, error) {
	// 系统管理员不参与项目,返回空
	if role == constants.RoleSystemAdmin {
		return []models.Project{}, nil
	}

	// 项目管理员和项目成员通过 project_members 表过滤
	return s.projectRepo.FindProjectsByUserID(userID)
}

// IsProjectMember 判断用户是否是项目成员
func (s *projectService) IsProjectMember(projectID uint, userID uint) (bool, error) {
	return s.memberRepo.IsMember(projectID, userID)
}

// CreateProject 创建项目并自动将创建者加入成员
func (s *projectService) CreateProject(name string, description string, creatorID uint) (*models.Project, error) {
	// 1. 检查项目名是否已存在
	exists, err := s.projectRepo.ExistsByName(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrProjectNameExists
	}

	// 2. 构造项目对象，自动设置owner_id为创建人
	ownerID := int(creatorID)
	project := &models.Project{
		Name:        name,
		Description: description,
		OwnerID:     &ownerID,
	}

	// 3. 开启事务执行创建项目和添加成员
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 创建项目 - 使用事务tx
		if err := tx.Create(project).Error; err != nil {
			return err
		}

		// 将创建者添加为项目管理员 - 使用事务tx
		member := &models.ProjectMember{
			ProjectID: project.ID,
			UserID:    creatorID,
			Role:      constants.RoleProjectManager,
		}
		if err := tx.Create(member).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 4. 重新查询项目以获取owner_name
	projectWithOwner, err := s.projectRepo.GetByID(project.ID)
	if err != nil {
		return nil, err
	}

	return projectWithOwner, nil
}

// checkProjectManagerPermission 检查用户是否有项目管理员权限
func (s *projectService) checkProjectManagerPermission(projectID uint, userID uint, role string) error {
	// 系统管理员无权访问项目
	if role == constants.RoleSystemAdmin {
		return ErrPermissionDenied
	}

	// 检查是否为项目成员且角色为project_manager
	isMember, err := s.memberRepo.IsMemberWithRole(projectID, userID, constants.RoleProjectManager)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrPermissionDenied
	}
	return nil
}

// UpdateProject 更新项目名称
func (s *projectService) UpdateProject(projectID uint, newName string, userID uint, role string) (*models.Project, error) {
	// 1. 权限验证
	if err := s.checkProjectManagerPermission(projectID, userID, role); err != nil {
		return nil, err
	}

	// 2. 检查新名称是否已被其他项目使用
	exists, err := s.projectRepo.ExistsByName(newName)
	if err != nil {
		return nil, err
	}
	if exists {
		// 验证是否是同一个项目(允许不修改名称的更新)
		currentProject, err := s.projectRepo.GetByID(projectID)
		if err != nil {
			return nil, ErrProjectNotFound
		}
		if currentProject.Name != newName {
			return nil, ErrProjectNameExists
		}
	}

	// 3. 更新项目名称
	return s.projectRepo.UpdateName(projectID, newName)
}

// DeleteProject 删除项目及其关联数据
func (s *projectService) DeleteProject(projectID uint, userID uint, role string) error {
	// 1. 权限验证
	if err := s.checkProjectManagerPermission(projectID, userID, role); err != nil {
		return err
	}

	// 2. 验证项目存在
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return ErrProjectNotFound
	}

	// 3. 级联删除项目
	return s.projectRepo.DeleteWithCascade(projectID)
}

// GetByID 获取项目详情并返回用户角色
func (s *projectService) GetByID(projectID uint, userID uint) (*models.Project, string, error) {
	// 1. 查询项目
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[Project GetByID] Project not found: project_id=%d", projectID)
			return nil, "", ErrProjectNotFound
		}
		log.Printf("[Project GetByID] Query failed: project_id=%d, error=%v", projectID, err)
		return nil, "", err
	}

	// 2. 查询用户角色
	userRole, err := s.memberRepo.GetMemberRole(projectID, userID)
	if err != nil {
		// 如果不是项目成员,检查是否是系统管理员
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 非项目成员,返回无权限错误
			log.Printf("[Project GetByID] Permission denied: project_id=%d, user_id=%d", projectID, userID)
			return nil, "", ErrPermissionDenied
		}
		log.Printf("[Project GetByID] GetMemberRole failed: project_id=%d, user_id=%d, error=%v", projectID, userID, err)
		return nil, "", err
	}

	log.Printf("[Project GetByID] Success: project_id=%d, user_id=%d, role=%s", projectID, userID, userRole)
	return project, userRole, nil
}

// UpdateProjectMetadata 更新项目元数据(name/description/status/owner_id)
func (s *projectService) UpdateProjectMetadata(projectID uint, updates map[string]interface{}, userID uint, role string) (*models.Project, error) {
	// 1. 权限验证 - 仅项目管理员可更新
	if err := s.checkProjectManagerPermission(projectID, userID, role); err != nil {
		return nil, err
	}

	// 2. 验证项目存在
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	// 3. 如果包含 name 字段,检查名称唯一性
	if newName, ok := updates["name"].(string); ok && newName != "" {
		exists, err := s.projectRepo.ExistsByName(newName)
		if err != nil {
			return nil, err
		}
		if exists {
			// 检查是否是同一个项目
			currentProject, err := s.projectRepo.GetByID(projectID)
			if err != nil {
				return nil, err
			}
			if currentProject.Name != newName {
				return nil, ErrProjectNameExists
			}
		}
	}

	// 4. 如果包含 status 字段,校验枚举值
	if status, ok := updates["status"].(string); ok {
		validStatuses := map[string]bool{
			constants.ProjectStatusPending:    true,
			constants.ProjectStatusInProgress: true,
			constants.ProjectStatusCompleted:  true,
		}
		if !validStatuses[status] {
			return nil, errors.New("invalid status value")
		}
	}

	// 5. 如果包含 owner_id 字段,验证用户存在性
	if ownerID, ok := updates["owner_id"].(int); ok && ownerID > 0 {
		var userExists int64
		if err := s.db.Model(&models.User{}).Where("id = ?", ownerID).Count(&userExists).Error; err != nil {
			return nil, err
		}
		if userExists == 0 {
			return nil, errors.New("owner user not found")
		}
	}

	// 6. 调用 Repository 层更新
	return s.projectRepo.Update(projectID, updates)
}

// GetProjectMembers 获取项目成员列表
func (s *projectService) GetProjectMembers(projectID uint, currentUserID uint) (*ProjectMembersResponse, error) {
	// 1. 验证项目存在
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("get project: %w", err)
	}

	// 2. 验证当前用户是项目成员
	isMember, err := s.memberRepo.IsMember(projectID, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		return nil, ErrPermissionDenied
	}

	// 3. 查询项目成员及用户信息
	membersWithUser, err := s.memberRepo.FindMembersWithUser(projectID)
	if err != nil {
		return nil, fmt.Errorf("find members: %w", err)
	}

	// 4. 分类并标记当前用户
	response := &ProjectMembersResponse{
		ProjectID:   projectID,
		ProjectName: project.Name,
		Managers:    []MemberInfo{},
		Members:     []MemberInfo{},
	}

	for _, member := range membersWithUser {
		memberInfo := MemberInfo{
			UserID:        member.UserID,
			Username:      member.Username,
			Nickname:      member.Nickname,
			IsCurrentUser: member.UserID == currentUserID,
		}

		if member.Role == constants.RoleProjectManager {
			response.Managers = append(response.Managers, memberInfo)
		} else if member.Role == constants.RoleProjectMember {
			response.Members = append(response.Members, memberInfo)
		}
	}

	return response, nil
}

// UpdateProjectMembers 批量更新项目成员
func (s *projectService) UpdateProjectMembers(projectID uint, managers []uint, members []uint, currentUserID uint) (*ProjectMembersResponse, error) {
	// 1. 验证项目存在
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("get project: %w", err)
	}

	// 2. 验证当前用户是项目管理员
	isManager, err := s.memberRepo.IsMemberWithRole(projectID, currentUserID, constants.RoleProjectManager)
	if err != nil {
		return nil, fmt.Errorf("check manager role: %w", err)
	}
	if !isManager {
		return nil, ErrPermissionDenied
	}

	// 3. 验证当前用户ID必须在managers列表中（自我保护）
	currentUserInManagers := false
	for _, uid := range managers {
		if uid == currentUserID {
			currentUserInManagers = true
			break
		}
	}
	if !currentUserInManagers {
		return nil, ErrSelfRemovalForbidden
	}

	// 4. 验证所有用户ID有效且角色匹配
	allUserIDs := append([]uint{}, managers...)
	allUserIDs = append(allUserIDs, members...)
	if len(allUserIDs) > 0 {
		users, err := s.userRepo.FindByIDs(allUserIDs)
		if err != nil {
			return nil, fmt.Errorf("find users by IDs: %w", err)
		}

		// 验证数量匹配
		if len(users) != len(allUserIDs) {
			return nil, ErrInvalidUserIDs
		}

		// 验证角色匹配
		userRoleMap := make(map[uint]string)
		for _, user := range users {
			userRoleMap[user.ID] = user.Role
		}

		for _, uid := range managers {
			if role, ok := userRoleMap[uid]; !ok || role != constants.RoleProjectManager {
				return nil, ErrRoleMismatch
			}
		}

		for _, uid := range members {
			if role, ok := userRoleMap[uid]; !ok || role != constants.RoleProjectMember {
				return nil, ErrRoleMismatch
			}
		}
	}

	// 5. 调用Repository执行批量更新
	if err := s.memberRepo.BatchUpdateMembers(projectID, managers, members); err != nil {
		return nil, fmt.Errorf("batch update members: %w", err)
	}

	// 6. 查询并返回更新后的成员列表
	return s.GetProjectMembers(projectID, currentUserID)
}
