package handlers

import (
	"log"
	"net/http"
	"strconv"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// AutoCasesHandler 自动化测试用例处理器
type AutoCasesHandler struct {
	service services.AutoTestCaseService
}

// NewAutoCasesHandler 创建处理器实例
func NewAutoCasesHandler(service services.AutoTestCaseService) *AutoCasesHandler {
	return &AutoCasesHandler{service: service}
}

// GetMetadata 获取元数据
// GET /api/v1/projects/:id/auto-cases/metadata?type=role1|role2|role3|role4
func (h *AutoCasesHandler) GetMetadata(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取类型参数(默认为 role1)
	caseType := c.DefaultQuery("type", "role1")

	// 调用服务
	metadata, err := h.service.GetMetadata(uint(projectID), userID, caseType)
	if err != nil {
		log.Printf("[Auto Metadata Get Failed] user_id=%d, project_id=%d, type=%s, error=%v", userID, projectID, caseType, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取元数据失败")
		return
	}

	log.Printf("[Auto Metadata Get] user_id=%d, project_id=%d, type=%s", userID, projectID, caseType)
	utils.SuccessResponse(c, metadata)
}

// UpdateMetadata 更新元数据
// PUT /api/v1/projects/:id/auto-cases/metadata?type=role1|role2|role3|role4
func (h *AutoCasesHandler) UpdateMetadata(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取类型参数(默认为 role1)
	caseType := c.DefaultQuery("type", "role1")

	// 解析请求体
	var req services.UpdateAutoMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	err = h.service.UpdateMetadata(uint(projectID), userID, caseType, req)
	if err != nil {
		log.Printf("[Auto Metadata Update Failed] user_id=%d, project_id=%d, type=%s, error=%v", userID, projectID, caseType, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新元数据失败")
		return
	}

	log.Printf("[Auto Metadata Update] user_id=%d, project_id=%d, type=%s", userID, projectID, caseType)
	utils.MessageResponse(c, http.StatusOK, "元数据更新成功")
}

// GetCases 获取用例列表
// GET /api/v1/projects/:id/auto-cases?case_type=role1&page=1&size=50
func (h *AutoCasesHandler) GetCases(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取查询参数
	caseType := c.DefaultQuery("case_type", "role1")
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "50")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	// 调用服务
	caseList, err := h.service.GetCases(uint(projectID), userID, caseType, page, size)
	if err != nil {
		log.Printf("[Auto Cases Get Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取用例列表失败")
		return
	}

	log.Printf("[Auto Cases Get] user_id=%d, project_id=%d, case_type=%s, total=%d", userID, projectID, caseType, caseList.Total)
	utils.SuccessResponse(c, caseList)
}

// CreateCase 创建新用例
// POST /api/v1/projects/:id/auto-cases
func (h *AutoCasesHandler) CreateCase(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 解析请求体
	var req services.CreateAutoCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	caseDTO, err := h.service.CreateCase(uint(projectID), userID, req)
	if err != nil {
		log.Printf("[Auto Case Create Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建用例失败")
		return
	}

	log.Printf("[Auto Case Create] user_id=%d, project_id=%d, case_id=%s, id=%d", userID, projectID, caseDTO.CaseID, caseDTO.ID)
	log.Printf("[Auto Case Create] Returning DTO: %+v", caseDTO)
	utils.SuccessResponse(c, caseDTO)
}

// UpdateCase 更新用例
// PATCH /api/v1/projects/:id/auto-cases/:caseId
func (h *AutoCasesHandler) UpdateCase(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用例ID（UUID字符串）
	caseID := c.Param("caseId")
	if caseID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用例ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 解析请求体
	var req services.UpdateAutoCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	err = h.service.UpdateCase(uint(projectID), userID, caseID, req)
	if err != nil {
		log.Printf("[Auto Case Update Failed] user_id=%d, project_id=%d, case_id=%s, error=%v", userID, projectID, caseID, err)
		if err.Error() == "无项目访问权限" || err.Error() == "用例不属于当前项目" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "用例不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新用例失败")
		return
	}

	log.Printf("[Auto Case Update] user_id=%d, project_id=%d, case_id=%s", userID, projectID, caseID)
	utils.MessageResponse(c, http.StatusOK, "用例更新成功")
}

// DeleteCase 删除用例
// DELETE /api/v1/projects/:id/auto-cases/:caseId
func (h *AutoCasesHandler) DeleteCase(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用例ID（UUID字符串）
	caseID := c.Param("caseId")
	if caseID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用例ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 调用服务
	err = h.service.DeleteCase(uint(projectID), userID, caseID)
	if err != nil {
		log.Printf("[Auto Case Delete Failed] user_id=%d, project_id=%d, case_id=%s, error=%v", userID, projectID, caseID, err)
		if err.Error() == "无项目访问权限" || err.Error() == "用例不属于当前项目" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "用例不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除用例失败")
		return
	}

	log.Printf("[Auto Case Delete] user_id=%d, project_id=%d, case_id=%s", userID, projectID, caseID)
	utils.MessageResponse(c, http.StatusOK, "用例删除成功")
}

// ReorderAllCases 按现有ID顺序重新编号所有用例
// POST /api/v1/projects/:id/auto-cases/reorder
func (h *AutoCasesHandler) ReorderAllCases(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 解析请求体
	var req struct {
		CaseType string   `json:"case_type" binding:"required,oneof=role1 role2 role3 role4"`
		CaseIDs  []string `json:"case_ids"` // 可选：指定重排顺序
	}
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	var count int
	if len(req.CaseIDs) > 0 {
		// 按指定顺序重排
		count, err = h.service.ReorderByIDs(uint(projectID), userID, req.CaseType, req.CaseIDs)
	} else {
		// 按现有ID顺序重排
		count, err = h.service.ReorderAllCases(uint(projectID), userID, req.CaseType)
	}
	if err != nil {
		log.Printf("[Auto Cases Reorder Failed] user_id=%d, project_id=%d, type=%s, error=%v", userID, projectID, req.CaseType, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "用例重排失败")
		return
	}

	log.Printf("[Auto Cases Reorder] user_id=%d, project_id=%d, type=%s, count=%d", userID, projectID, req.CaseType, count)
	utils.SuccessResponse(c, gin.H{"message": "用例重排成功", "count": count})
}

// InsertCase 在指定位置插入用例
// POST /api/v1/projects/:id/auto-cases/insert
func (h *AutoCasesHandler) InsertCase(c *gin.Context) {
	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取项目ID
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 解析请求体
	var req struct {
		CaseType     string `json:"case_type" binding:"required,oneof=role1 role2 role3 role4"`
		Position     string `json:"position" binding:"required,oneof=before after"`
		TargetCaseID string `json:"target_case_id" binding:"required"`
	}
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	log.Printf("[Auto Insert Case Start] user_id=%d, project_id=%d, type=%s, position=%s, target=%s",
		userID, projectID, req.CaseType, req.Position, req.TargetCaseID)
	newCase, err := h.service.InsertCase(uint(projectID), userID, req.CaseType, req.Position, req.TargetCaseID)
	if err != nil {
		log.Printf("[Auto Insert Case Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "插入用例失败")
		return
	}

	log.Printf("[Auto Insert Case Success] user_id=%d, project_id=%d, new_case_id=%s", userID, projectID, newCase.CaseID)
	utils.SuccessResponse(c, newCase)
}

// BatchDeleteCases 批量删除用例
// POST /api/v1/projects/:id/auto-cases/batch-delete
func (h *AutoCasesHandler) BatchDeleteCases(c *gin.Context) {
	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取项目ID
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 解析请求体
	var req struct {
		CaseType string   `json:"case_type" binding:"required,oneof=role1 role2 role3 role4"`
		CaseIDs  []string `json:"case_ids" binding:"required,min=1"`
	}
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	log.Printf("[Auto Batch Delete Start] user_id=%d, project_id=%d, type=%s, count=%d",
		userID, projectID, req.CaseType, len(req.CaseIDs))
	deletedCount, failedCaseIDs, err := h.service.BatchDeleteCases(uint(projectID), userID, req.CaseType, req.CaseIDs)
	if err != nil {
		log.Printf("[Auto Batch Delete Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量删除失败")
		return
	}

	log.Printf("[Auto Batch Delete Success] user_id=%d, project_id=%d, deleted=%d, failed=%d",
		userID, projectID, deletedCount, len(failedCaseIDs))
	utils.SuccessResponse(c, gin.H{
		"message":         "批量删除完成",
		"deleted_count":   deletedCount,
		"failed_case_ids": failedCaseIDs,
	})
}

// ReassignIDs 重新分配用例ID
// POST /api/v1/projects/:id/auto-cases/reassign-ids
func (h *AutoCasesHandler) ReassignIDs(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 绑定请求参数
	var req struct {
		CaseType string `json:"caseType" binding:"required,oneof=role1 role2 role3 role4"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务重新分配ID
	log.Printf("[Auto Reassign IDs Start] user_id=%d, project_id=%d, type=%s", userID, projectID, req.CaseType)
	if err := h.service.ReassignAllIDs(uint(projectID), userID, req.CaseType); err != nil {
		log.Printf("[Auto Reassign IDs Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "重新分配ID失败")
		return
	}

	log.Printf("[Auto Reassign IDs Success] user_id=%d, project_id=%d, type=%s", userID, projectID, req.CaseType)
	utils.SuccessResponse(c, gin.H{
		"message": "重新分配ID成功",
	})
}

// BatchSaveVersion 批量保存版本(ROLE1-4所有用例)
// POST /api/v1/projects/:id/auto-cases/versions
func (h *AutoCasesHandler) BatchSaveVersion(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	log.Printf("[Auto Batch Save Version Start] user_id=%d, project_id=%d", userID, projectID)
	versionInfo, err := h.service.BatchSaveVersion(uint(projectID), userID)
	if err != nil {
		log.Printf("[Auto Batch Save Version Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "版本保存失败: "+err.Error())
		return
	}

	log.Printf("[Auto Batch Save Version Success] user_id=%d, project_id=%d, version_id=%d",
		userID, projectID, versionInfo.VersionID)
	utils.SuccessResponse(c, versionInfo)
}

// GetAutoVersions 获取版本列表
// GET /api/v1/projects/:id/auto-cases/versions?page=1&size=20
func (h *AutoCasesHandler) GetAutoVersions(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	log.Printf("[Auto Get Versions] user_id=%d, project_id=%d, page=%d, size=%d", userID, projectID, page, size)
	result, err := h.service.GetVersionList(uint(projectID), userID, page, size)
	if err != nil {
		log.Printf("[Auto Get Versions Failed] error=%v", err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取版本列表失败")
		return
	}

	utils.SuccessResponse(c, result)
}

// DownloadAutoVersion 下载版本(zip打包4个Excel)
// GET /api/v1/projects/:id/auto-cases/versions/:versionId/export
func (h *AutoCasesHandler) DownloadAutoVersion(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	versionIDStr := c.Param("versionId")
	versionID := versionIDStr // versionID 现在是 string 类型

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	log.Printf("[Auto Download Version] user_id=%d, project_id=%d, version_id=%s", userID, projectID, versionID)
	zipData, filename, err := h.service.DownloadVersion(uint(projectID), userID, versionID)
	if err != nil {
		log.Printf("[Auto Download Version Failed] error=%v", err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "版本不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "下载失败: "+err.Error())
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/zip", zipData)
}

// DeleteAutoVersion 删除版本
// DELETE /api/v1/projects/:id/auto-cases/versions/:versionId
func (h *AutoCasesHandler) DeleteAutoVersion(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	versionIDStr := c.Param("versionId")
	versionID := versionIDStr // versionID 现在是 string 类型

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	log.Printf("[Auto Delete Version] user_id=%d, project_id=%d, version_id=%s", userID, projectID, versionID)
	if err := h.service.DeleteVersion(uint(projectID), userID, versionID); err != nil {
		log.Printf("[Auto Delete Version Failed] error=%v", err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "版本不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	log.Printf("[Auto Delete Version Success] user_id=%d, project_id=%d, version_id=%s", userID, projectID, versionID)
	utils.SuccessResponse(c, gin.H{"message": "版本已删除"})
}

// UpdateAutoVersionRemark 更新版本备注
// PUT /api/v1/projects/:id/auto-cases/versions/:versionId/remark
func (h *AutoCasesHandler) UpdateAutoVersionRemark(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	versionIDStr := c.Param("versionId")
	versionID := versionIDStr // versionID 现在是 string 类型

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	var req struct {
		Remark string `json:"remark" binding:"max=200"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "备注长度不能超过200字符")
		return
	}

	log.Printf("[Auto Update Version Remark] user_id=%d, project_id=%d, version_id=%s", userID, projectID, versionID)
	err = h.service.UpdateVersionRemark(uint(projectID), userID, versionID, req.Remark)
	if err != nil {
		log.Printf("[Auto Update Version Remark Failed] error=%v", err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "版本不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	log.Printf("[Auto Update Version Remark Success] user_id=%d, project_id=%d, version_id=%s", userID, projectID, versionID)
	utils.SuccessResponse(c, gin.H{"message": "备注已更新"})
}
