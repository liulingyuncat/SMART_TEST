package handlers

import (
	"net/http"
	"strconv"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"github.com/gin-gonic/gin"
)

// CaseGroupHandler 用例集处理器
type CaseGroupHandler struct {
	repo       *repositories.CaseGroupRepository
	autoRepo   repositories.AutoTestCaseRepository
	manualRepo repositories.ManualTestCaseRepository
	apiRepo    repositories.ApiTestCaseRepository
}

// NewCaseGroupHandler 创建用例集处理器
func NewCaseGroupHandler(repo *repositories.CaseGroupRepository, autoRepo repositories.AutoTestCaseRepository, manualRepo repositories.ManualTestCaseRepository, apiRepo repositories.ApiTestCaseRepository) *CaseGroupHandler {
	return &CaseGroupHandler{
		repo:       repo,
		autoRepo:   autoRepo,
		manualRepo: manualRepo,
		apiRepo:    apiRepo,
	}
}

// GetCaseGroups 获取用例集列表
// GET /api/v1/projects/:id/case-groups?case_type=overall
func (h *CaseGroupHandler) GetCaseGroups(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	caseType := c.DefaultQuery("case_type", "overall")

	groups, err := h.repo.GetByProjectAndType(uint(projectID), caseType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetCaseGroup 获取单个用例集详情
// GET /api/v1/case-groups/:id
func (h *CaseGroupHandler) GetCaseGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case group ID"})
		return
	}

	group, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// CreateCaseGroup 创建用例集
// POST /api/v1/projects/:id/case-groups
func (h *CaseGroupHandler) CreateCaseGroup(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		CaseType     string `json:"case_type" binding:"required"`
		GroupName    string `json:"group_name" binding:"required"`
		Description  string `json:"description"`
		DisplayOrder int    `json:"display_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否已存在
	existing, err := h.repo.GetByName(uint(projectID), req.CaseType, req.GroupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Case group already exists"})
		return
	}

	group := &models.CaseGroup{
		ProjectID:    uint(projectID),
		CaseType:     req.CaseType,
		GroupName:    req.GroupName,
		Description:  req.Description,
		DisplayOrder: req.DisplayOrder,
	}

	if err := h.repo.Create(group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// UpdateCaseGroup 更新用例集
// PUT /api/v1/case-groups/:id
func (h *CaseGroupHandler) UpdateCaseGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case group ID"})
		return
	}

	group, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case group not found"})
		return
	}

	var req struct {
		GroupName    string `json:"group_name"`
		Description  string `json:"description"`
		DisplayOrder int    `json:"display_order"`
		// 元数据字段
		MetaProtocol *string `json:"meta_protocol"`
		MetaServer   *string `json:"meta_server"`
		MetaPort     *string `json:"meta_port"`
		MetaUser     *string `json:"meta_user"`
		MetaPassword *string `json:"meta_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果修改了组名，检查是否重复
	if req.GroupName != "" && req.GroupName != group.GroupName {
		existing, err := h.repo.GetByName(group.ProjectID, group.CaseType, req.GroupName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if existing != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Case group name already exists"})
			return
		}
		group.GroupName = req.GroupName
	}

	if req.Description != "" {
		group.Description = req.Description
	}
	group.DisplayOrder = req.DisplayOrder

	// 更新元数据字段
	if req.MetaProtocol != nil {
		group.MetaProtocol = *req.MetaProtocol
	}
	if req.MetaServer != nil {
		group.MetaServer = *req.MetaServer
	}
	if req.MetaPort != nil {
		group.MetaPort = *req.MetaPort
	}
	if req.MetaUser != nil {
		group.MetaUser = *req.MetaUser
	}
	if req.MetaPassword != nil {
		group.MetaPassword = *req.MetaPassword
	}

	if err := h.repo.Update(group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteCaseGroup 删除用例集（级联删除用例）
// DELETE /api/v1/case-groups/:id
func (h *CaseGroupHandler) DeleteCaseGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case group ID"})
		return
	}

	group, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case group not found"})
		return
	}

	// 级联删除：先硬删除该用例集下的所有用例
	switch group.CaseType {
	case "web", "role1", "role2", "role3", "role4":
		// Web/自动化用例
		err = h.autoRepo.DeleteByCaseGroup(group.ProjectID, group.CaseType, group.GroupName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated auto test cases: " + err.Error()})
			return
		}
	case "overall", "change":
		// Manual用例
		err = h.manualRepo.DeleteByCaseGroup(group.ProjectID, group.CaseType, group.GroupName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated manual test cases: " + err.Error()})
			return
		}
	case "api":
		// API接口用例
		err = h.apiRepo.DeleteByCaseGroup(group.ProjectID, group.GroupName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated API test cases: " + err.Error()})
			return
		}
	}

	// 硬删除用例集本身
	if err := h.repo.HardDelete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Case group and associated test cases deleted successfully"})
}
