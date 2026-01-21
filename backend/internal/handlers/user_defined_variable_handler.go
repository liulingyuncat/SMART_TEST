package handlers

import (
	"net/http"
	"strconv"

	"webtest/internal/models"
	"webtest/internal/services"

	"github.com/gin-gonic/gin"
)

// UserDefinedVariableHandler 用户自定义变量处理器
type UserDefinedVariableHandler struct {
	service services.UserDefinedVariableService
}

// NewUserDefinedVariableHandler 创建处理器实例
func NewUserDefinedVariableHandler(service services.UserDefinedVariableService) *UserDefinedVariableHandler {
	return &UserDefinedVariableHandler{service: service}
}

// GetVariables 获取用例集的变量列表
// GET /api/v1/case-groups/:groupId/variables?group_type=web
func (h *UserDefinedVariableHandler) GetVariables(c *gin.Context) {
	groupIDStr := c.Param("groupId")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group_id"})
		return
	}

	groupType := c.Query("group_type")
	if groupType == "" {
		groupType = "web" // 默认web类型
	}

	variables, err := h.service.GetVariablesByGroup(uint(groupID), groupType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group_id":   groupID,
		"group_type": groupType,
		"variables":  variables,
	})
}

// SaveVariablesRequest 批量保存变量请求
type SaveVariablesRequest struct {
	ProjectID uint                          `json:"project_id"`
	GroupType string                        `json:"group_type"`
	Variables []*models.UserDefinedVariable `json:"variables"`
}

// SaveVariables 批量保存变量（替换模式）
// PUT /api/v1/case-groups/:groupId/variables
func (h *UserDefinedVariableHandler) SaveVariables(c *gin.Context) {
	groupIDStr := c.Param("groupId")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group_id"})
		return
	}

	var req SaveVariablesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.GroupType == "" {
		req.GroupType = "web"
	}

	err = h.service.SaveVariables(uint(groupID), req.GroupType, req.ProjectID, req.Variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "variables saved successfully",
	})
}

// AddVariableRequest 添加变量请求
type AddVariableRequest struct {
	ProjectID uint   `json:"project_id"`
	GroupType string `json:"group_type"`
	VarKey    string `json:"var_key" binding:"required"`
	VarDesc   string `json:"var_desc"`
	VarValue  string `json:"var_value"`
	VarType   string `json:"var_type"`
}

// AddVariable 添加单个变量
// POST /api/v1/case-groups/:groupId/variables
func (h *UserDefinedVariableHandler) AddVariable(c *gin.Context) {
	groupIDStr := c.Param("groupId")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group_id"})
		return
	}

	var req AddVariableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.GroupType == "" {
		req.GroupType = "web"
	}

	variable := &models.UserDefinedVariable{
		ProjectID: req.ProjectID,
		GroupID:   uint(groupID),
		GroupType: req.GroupType,
		VarKey:    req.VarKey,
		VarDesc:   req.VarDesc,
		VarValue:  req.VarValue,
		VarType:   req.VarType,
	}

	err = h.service.AddVariable(variable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, variable)
}

// UpdateVariableRequest 更新变量请求
type UpdateVariableRequest struct {
	VarKey   string `json:"var_key"`
	VarDesc  string `json:"var_desc"`
	VarValue string `json:"var_value"`
	VarType  string `json:"var_type"`
}

// UpdateVariable 更新单个变量
// PUT /api/v1/case-groups/:groupId/variables/:varId
func (h *UserDefinedVariableHandler) UpdateVariable(c *gin.Context) {
	varIDStr := c.Param("varId")
	varID, err := strconv.ParseUint(varIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid var_id"})
		return
	}

	var req UpdateVariableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variable := &models.UserDefinedVariable{
		ID:       uint(varID),
		VarKey:   req.VarKey,
		VarDesc:  req.VarDesc,
		VarValue: req.VarValue,
		VarType:  req.VarType,
	}

	err = h.service.UpdateVariable(variable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "variable updated successfully",
	})
}

// DeleteVariable 删除单个变量
// DELETE /api/v1/case-groups/:groupId/variables/:varId
func (h *UserDefinedVariableHandler) DeleteVariable(c *gin.Context) {
	varIDStr := c.Param("varId")
	varID, err := strconv.ParseUint(varIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid var_id"})
		return
	}

	err = h.service.DeleteVariable(uint(varID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "variable deleted successfully",
	})
}

// GetTaskVariables 获取执行任务的变量列表
// GET /api/v1/execution-tasks/:taskUuid/variables?group_id=1&group_type=web
func (h *UserDefinedVariableHandler) GetTaskVariables(c *gin.Context) {
	taskUUID := c.Param("task_uuid")
	if taskUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_uuid is required"})
		return
	}

	groupIDStr := c.Query("group_id")
	groupID, _ := strconv.ParseUint(groupIDStr, 10, 64)

	groupType := c.Query("group_type")
	if groupType == "" {
		groupType = "web"
	}

	variables, err := h.service.GetVariablesByTask(taskUUID, uint(groupID), groupType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task_uuid":  taskUUID,
		"group_id":   groupID,
		"group_type": groupType,
		"variables":  variables,
	})
}

// SaveTaskVariablesRequest 保存任务变量请求
type SaveTaskVariablesRequest struct {
	ProjectID uint                          `json:"project_id"`
	GroupID   uint                          `json:"group_id"`
	GroupType string                        `json:"group_type"`
	Variables []*models.UserDefinedVariable `json:"variables"`
}

// SaveTaskVariables 批量保存任务变量（替换模式）
// PUT /api/v1/execution-tasks/:taskUuid/variables
func (h *UserDefinedVariableHandler) SaveTaskVariables(c *gin.Context) {
	taskUUID := c.Param("task_uuid")
	if taskUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_uuid is required"})
		return
	}

	var req SaveTaskVariablesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.GroupType == "" {
		req.GroupType = "web"
	}

	err := h.service.SaveTaskVariables(taskUUID, req.GroupID, req.GroupType, req.ProjectID, req.Variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task variables saved successfully",
	})
}
