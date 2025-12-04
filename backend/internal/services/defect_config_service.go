package services

import (
	"errors"
	"fmt"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"gorm.io/gorm"
)

// DefectConfigService 缺陷配置服务接口
type DefectConfigService interface {
	// Subject管理
	CreateSubject(projectID uint, req *models.DefectSubjectCreateRequest) (*models.DefectSubject, error)
	GetSubject(id uint) (*models.DefectSubject, error)
	UpdateSubject(id uint, req *models.DefectSubjectUpdateRequest) error
	DeleteSubject(id uint) error
	ListSubjects(projectID uint) ([]*models.DefectSubject, error)

	// Phase管理
	CreatePhase(projectID uint, req *models.DefectPhaseCreateRequest) (*models.DefectPhase, error)
	GetPhase(id uint) (*models.DefectPhase, error)
	UpdatePhase(id uint, req *models.DefectPhaseUpdateRequest) error
	DeletePhase(id uint) error
	ListPhases(projectID uint) ([]*models.DefectPhase, error)
}

type defectConfigService struct {
	subjectRepo repositories.DefectSubjectRepository
	phaseRepo   repositories.DefectPhaseRepository
}

// NewDefectConfigService 创建缺陷配置服务实例
func NewDefectConfigService(
	subjectRepo repositories.DefectSubjectRepository,
	phaseRepo repositories.DefectPhaseRepository,
) DefectConfigService {
	return &defectConfigService{
		subjectRepo: subjectRepo,
		phaseRepo:   phaseRepo,
	}
}

// ========== Subject管理 ==========

// CreateSubject 创建Subject
func (s *defectConfigService) CreateSubject(projectID uint, req *models.DefectSubjectCreateRequest) (*models.DefectSubject, error) {
	// 检查名称是否重复
	exists, err := s.subjectRepo.ExistsByName(projectID, req.Name, 0)
	if err != nil {
		return nil, fmt.Errorf("check subject exists: %w", err)
	}
	if exists {
		return nil, errors.New("subject name already exists")
	}

	subject := &models.DefectSubject{
		ProjectID: projectID,
		Name:      req.Name,
		SortOrder: req.SortOrder,
	}

	if err := s.subjectRepo.Create(subject); err != nil {
		return nil, fmt.Errorf("create subject: %w", err)
	}

	return subject, nil
}

// GetSubject 获取Subject
func (s *defectConfigService) GetSubject(id uint) (*models.DefectSubject, error) {
	subject, err := s.subjectRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, fmt.Errorf("get subject: %w", err)
	}
	return subject, nil
}

// UpdateSubject 更新Subject
func (s *defectConfigService) UpdateSubject(id uint, req *models.DefectSubjectUpdateRequest) error {
	// 获取现有记录
	subject, err := s.subjectRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject not found")
		}
		return fmt.Errorf("get subject: %w", err)
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		// 检查名称是否与其他记录重复
		exists, err := s.subjectRepo.ExistsByName(subject.ProjectID, *req.Name, id)
		if err != nil {
			return fmt.Errorf("check subject exists: %w", err)
		}
		if exists {
			return errors.New("subject name already exists")
		}
		updates["name"] = *req.Name
	}

	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) == 0 {
		return nil
	}

	if err := s.subjectRepo.Update(id, updates); err != nil {
		return fmt.Errorf("update subject: %w", err)
	}

	return nil
}

// DeleteSubject 删除Subject
func (s *defectConfigService) DeleteSubject(id uint) error {
	if err := s.subjectRepo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject not found")
		}
		return fmt.Errorf("delete subject: %w", err)
	}
	return nil
}

// ListSubjects 获取Subject列表
func (s *defectConfigService) ListSubjects(projectID uint) ([]*models.DefectSubject, error) {
	subjects, err := s.subjectRepo.ListByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("list subjects: %w", err)
	}
	return subjects, nil
}

// ========== Phase管理 ==========

// CreatePhase 创建Phase
func (s *defectConfigService) CreatePhase(projectID uint, req *models.DefectPhaseCreateRequest) (*models.DefectPhase, error) {
	// 检查名称是否重复
	exists, err := s.phaseRepo.ExistsByName(projectID, req.Name, 0)
	if err != nil {
		return nil, fmt.Errorf("check phase exists: %w", err)
	}
	if exists {
		return nil, errors.New("phase name already exists")
	}

	phase := &models.DefectPhase{
		ProjectID: projectID,
		Name:      req.Name,
		SortOrder: req.SortOrder,
	}

	if err := s.phaseRepo.Create(phase); err != nil {
		return nil, fmt.Errorf("create phase: %w", err)
	}

	return phase, nil
}

// GetPhase 获取Phase
func (s *defectConfigService) GetPhase(id uint) (*models.DefectPhase, error) {
	phase, err := s.phaseRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("phase not found")
		}
		return nil, fmt.Errorf("get phase: %w", err)
	}
	return phase, nil
}

// UpdatePhase 更新Phase
func (s *defectConfigService) UpdatePhase(id uint, req *models.DefectPhaseUpdateRequest) error {
	// 获取现有记录
	phase, err := s.phaseRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("phase not found")
		}
		return fmt.Errorf("get phase: %w", err)
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		// 检查名称是否与其他记录重复
		exists, err := s.phaseRepo.ExistsByName(phase.ProjectID, *req.Name, id)
		if err != nil {
			return fmt.Errorf("check phase exists: %w", err)
		}
		if exists {
			return errors.New("phase name already exists")
		}
		updates["name"] = *req.Name
	}

	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) == 0 {
		return nil
	}

	if err := s.phaseRepo.Update(id, updates); err != nil {
		return fmt.Errorf("update phase: %w", err)
	}

	return nil
}

// DeletePhase 删除Phase
func (s *defectConfigService) DeletePhase(id uint) error {
	if err := s.phaseRepo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("phase not found")
		}
		return fmt.Errorf("delete phase: %w", err)
	}
	return nil
}

// ListPhases 获取Phase列表
func (s *defectConfigService) ListPhases(projectID uint) ([]*models.DefectPhase, error) {
	phases, err := s.phaseRepo.ListByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("list phases: %w", err)
	}
	return phases, nil
}
