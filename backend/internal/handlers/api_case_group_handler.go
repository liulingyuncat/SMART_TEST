package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"webtest/internal/models"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApiCaseGroupHandler 接口用例集处理器
type ApiCaseGroupHandler struct {
	db *gorm.DB
}

// NewApiCaseGroupHandler 创建处理器实例
func NewApiCaseGroupHandler(db *gorm.DB) *ApiCaseGroupHandler {
	return &ApiCaseGroupHandler{db: db}
}

// GetCaseGroups 获取用例集列表
// GET /api/v1/projects/:id/api-case-groups
func (h *ApiCaseGroupHandler) GetCaseGroups(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 验证用户权限
	_, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 从api_test_cases表按case_group字段去重查询
	var caseGroups []string
	fmt.Printf("[GetCaseGroups] 🔍 开始查询项目 %d 的用例集\n", projectID)

	err = h.db.Model(&models.ApiTestCase{}).
		Where("project_id = ? AND case_group != ''", uint(projectID)).
		Distinct("case_group").
		Pluck("case_group", &caseGroups).Error

	if err != nil {
		fmt.Printf("[GetCaseGroups] ❌ 查询失败: %v\n", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("查询用例集失败: %v", err))
		return
	}

	fmt.Printf("[GetCaseGroups] ✅ 查询成功，找到 %d 个用例集\n", len(caseGroups))
	fmt.Printf("[GetCaseGroups] 📋 用例集列表: %v\n", caseGroups)

	utils.SuccessResponse(c, gin.H{
		"case_groups": caseGroups,
	})
}

// CreateCaseGroup 创建用例集
// POST /api/v1/projects/:id/api-case-groups
// Body: { "group_name": "xxx" }
func (h *ApiCaseGroupHandler) CreateCaseGroup(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 验证用户权限
	_, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req struct {
		GroupName string `json:"group_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请输入用例集名称")
		return
	}

	// 验证group_name非空
	groupName := strings.TrimSpace(req.GroupName)
	if groupName == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "用例集名称不能为空")
		return
	}

	// 验证group_name不重复
	var count int64
	err = h.db.Model(&models.ApiTestCase{}).
		Where("project_id = ? AND case_group = ?", uint(projectID), groupName).
		Count(&count).Error

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("验证用例集失败: %v", err))
		return
	}

	if count > 0 {
		utils.ErrorResponse(c, http.StatusConflict, "用例集名称已存在")
		return
	}

	fmt.Printf("[CreateCaseGroup] 🆕 开始创建用例集 - 项目: %d, 名称: %s\n", projectID, groupName)

	// 创建一个空的占位用例来标识用例集的存在
	// 这样GetCaseGroups就能查询到这个用例集
	placeholderCase := models.ApiTestCase{
		ProjectID:  uint(projectID),
		CaseGroup:  groupName,
		CaseNumber: "", // 空用例编号，作为占位符
		Method:     "GET",
		URL:        "",
		Screen:     "",
		Remark:     "",
	}

	fmt.Printf("[CreateCaseGroup] 📤 准备插入占位记录: %+v\n", placeholderCase)

	if err := h.db.Create(&placeholderCase).Error; err != nil {
		fmt.Printf("[CreateCaseGroup] ❌ 插入失败: %v\n", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("创建用例集失败: %v", err))
		return
	}

	fmt.Printf("[CreateCaseGroup] ✅ 插入成功，ID: %s\n", placeholderCase.ID)

	// 验证插入后是否能查询到
	var verifyCount int64
	h.db.Model(&models.ApiTestCase{}).
		Where("project_id = ? AND case_group = ?", uint(projectID), groupName).
		Count(&verifyCount)
	fmt.Printf("[CreateCaseGroup] 🔍 验证查询: 项目 %d 中名为 '%s' 的记录数: %d\n", projectID, groupName, verifyCount)

	utils.SuccessResponse(c, gin.H{
		"message":  "用例集创建成功",
		"group_id": placeholderCase.ID,
	})
}

// UpdateCaseGroup 更新用例集名称
// PUT /api/v1/api-case-groups/:groupId
// Body: { "group_name": "new_name" }
func (h *ApiCaseGroupHandler) UpdateCaseGroup(c *gin.Context) {
	oldGroupName := c.Param("groupId")

	// 验证用户权限
	_, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req struct {
		GroupName string `json:"group_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请输入新的用例集名称")
		return
	}

	// 验证新group_name非空
	newGroupName := strings.TrimSpace(req.GroupName)
	if newGroupName == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "用例集名称不能为空")
		return
	}

	// 更新所有匹配记录的case_group字段
	result := h.db.Model(&models.ApiTestCase{}).
		Where("case_group = ?", oldGroupName).
		Update("case_group", newGroupName)

	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("更新用例集失败: %v", result.Error))
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":        "用例集更新成功",
		"updated_count":  result.RowsAffected,
		"new_group_name": newGroupName,
	})
}

// DeleteCaseGroup 硬删除用例集（级联删除用例集内所有用例）
// DELETE /api/v1/api-case-groups/:groupId
func (h *ApiCaseGroupHandler) DeleteCaseGroup(c *gin.Context) {
	groupName := c.Param("groupId")
	projectIDStr := c.Param("id")

	// 验证用户权限
	_, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析项目ID
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	fmt.Printf("[DeleteCaseGroup] 🗑️ 开始硬删除API用例集: %s (项目ID: %d)\n", groupName, projectID)

	// 1. 硬删除指定case_group的所有API用例记录
	result := h.db.Unscoped().
		Where("project_id = ? AND case_group = ?", uint(projectID), groupName).
		Delete(&models.ApiTestCase{})

	if result.Error != nil {
		fmt.Printf("[DeleteCaseGroup] ❌ 删除API用例失败: %v\n", result.Error)
		utils.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("删除API用例失败: %v", result.Error))
		return
	}

	deletedCaseCount := result.RowsAffected
	fmt.Printf("[DeleteCaseGroup] ✅ 已硬删除 %d 条API用例\n", deletedCaseCount)

	// 2. 硬删除case_groups表中的用例集记录
	result = h.db.Unscoped().
		Where("project_id = ? AND case_type = 'api' AND group_name = ?", uint(projectID), groupName).
		Delete(&models.CaseGroup{})

	if result.Error != nil {
		fmt.Printf("[DeleteCaseGroup] ❌ 删除用例集记录失败: %v\n", result.Error)
		utils.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("删除用例集记录失败: %v", result.Error))
		return
	}

	deletedGroupCount := result.RowsAffected
	fmt.Printf("[DeleteCaseGroup] ✅ 已硬删除用例集记录 (删除数: %d)\n", deletedGroupCount)

	utils.SuccessResponse(c, gin.H{
		"message":            "用例集及其所有用例已硬删除成功",
		"deleted_cases":      deletedCaseCount,
		"deleted_case_group": deletedGroupCount,
	})
}
