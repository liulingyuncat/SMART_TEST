package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// ManualCasesHandler 手工测试用例处理器
type ManualCasesHandler struct {
	service        services.ManualTestCaseService
	versionService services.VersionService
}

// NewManualCasesHandler 创建处理器实例
func NewManualCasesHandler(service services.ManualTestCaseService, versionService services.VersionService) *ManualCasesHandler {
	return &ManualCasesHandler{
		service:        service,
		versionService: versionService,
	}
}

// GetMetadata 获取元数据
// GET /api/v1/projects/:id/manual-cases/metadata?type=overall|change
func (h *ManualCasesHandler) GetMetadata(c *gin.Context) {
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

	// 获取类型参数(默认为 overall)
	caseType := c.DefaultQuery("type", "overall")

	// 调用服务
	metadata, err := h.service.GetMetadata(uint(projectID), userID, caseType)
	if err != nil {
		log.Printf("[Metadata Get Failed] user_id=%d, project_id=%d, type=%s, error=%v", userID, projectID, caseType, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取元数据失败")
		return
	}

	log.Printf("[Metadata Get] user_id=%d, project_id=%d, type=%s", userID, projectID, caseType)
	utils.SuccessResponse(c, metadata)
}

// UpdateMetadata 更新元数据
// PUT /api/v1/projects/:id/manual-cases/metadata?type=overall|change
func (h *ManualCasesHandler) UpdateMetadata(c *gin.Context) {
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

	// 获取类型参数(默认为 overall)
	caseType := c.DefaultQuery("type", "overall")

	// 解析请求体
	var req services.UpdateMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	err = h.service.UpdateMetadata(uint(projectID), userID, caseType, req)
	if err != nil {
		log.Printf("[Metadata Update Failed] user_id=%d, project_id=%d, type=%s, error=%v", userID, projectID, caseType, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新元数据失败")
		return
	}

	log.Printf("[Metadata Update] user_id=%d, project_id=%d, type=%s", userID, projectID, caseType)
	utils.MessageResponse(c, http.StatusOK, "元数据更新成功")
}

// GetCases 获取用例列表
// GET /api/v1/projects/:id/manual-cases?language=中文&page=1&size=50
func (h *ManualCasesHandler) GetCases(c *gin.Context) {
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
	caseType := c.DefaultQuery("case_type", "overall")
	language := c.DefaultQuery("language", "中文")
	caseGroup := c.Query("case_group") // 获取用例集过滤参数（可选）
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "50")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	// 调用服务
	caseList, err := h.service.GetCases(uint(projectID), userID, caseType, language, page, size, caseGroup)
	if err != nil {
		log.Printf("[Cases Get Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取用例列表失败")
		return
	}

	log.Printf("[Cases Get] user_id=%d, project_id=%d, case_type=%s, language=%s, total=%d", userID, projectID, caseType, language, caseList.Total)
	utils.SuccessResponse(c, caseList)
}

// CreateCase 创建新用例
// POST /api/v1/projects/:id/manual-cases
func (h *ManualCasesHandler) CreateCase(c *gin.Context) {
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
	var req services.CreateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	caseDTO, err := h.service.CreateCase(uint(projectID), userID, req)
	if err != nil {
		log.Printf("[Case Create Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建用例失败")
		return
	}

	log.Printf("[Case Create] user_id=%d, project_id=%d, case_id=%d", userID, projectID, caseDTO.ID)
	utils.SuccessResponse(c, caseDTO)
}

// UpdateCase 更新用例
// PATCH /api/v1/projects/:id/manual-cases/:caseId
func (h *ManualCasesHandler) UpdateCase(c *gin.Context) {
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
	var req services.UpdateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务（caseID现在是UUID字符串）
	err = h.service.UpdateCase(uint(projectID), userID, caseID, req)
	if err != nil {
		log.Printf("[Case Update Failed] user_id=%d, project_id=%d, case_id=%s, error=%v", userID, projectID, caseID, err)
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

	log.Printf("[Case Update] user_id=%d, project_id=%d, case_id=%s", userID, projectID, caseID)
	utils.MessageResponse(c, http.StatusOK, "用例更新成功")
}

// DeleteCase 删除用例
// DELETE /api/v1/projects/:id/manual-cases/:caseId
func (h *ManualCasesHandler) DeleteCase(c *gin.Context) {
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

	// 调用服务（caseID现在是UUID字符串）
	err = h.service.DeleteCase(uint(projectID), userID, caseID)
	if err != nil {
		log.Printf("[Case Delete Failed] user_id=%d, project_id=%d, case_id=%s, error=%v", userID, projectID, caseID, err)
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

	log.Printf("[Case Delete] user_id=%d, project_id=%d, case_id=%s", userID, projectID, caseID)
	utils.MessageResponse(c, http.StatusOK, "用例删除成功")
}

// ReorderCases 重新排序用例
// POST /api/v1/projects/:id/manual-cases/reorder
func (h *ManualCasesHandler) ReorderCases(c *gin.Context) {
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
		CaseType string `json:"case_type" binding:"required,oneof=overall change acceptance ai"`
		CaseIDs  []uint `json:"case_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	newIDs, err := h.service.ReorderCases(uint(projectID), userID, req.CaseType, req.CaseIDs)
	if err != nil {
		log.Printf("[Cases Reorder Failed] user_id=%d, project_id=%d, type=%s, error=%v", userID, projectID, req.CaseType, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "用例重排失败")
		return
	}

	log.Printf("[Cases Reorder] user_id=%d, project_id=%d, type=%s, count=%d", userID, projectID, req.CaseType, len(newIDs))
	utils.SuccessResponse(c, gin.H{"new_ids": newIDs})
}

// ReorderCasesByDrag 拖拽重排处理器
func (h *ManualCasesHandler) ReorderCasesByDrag(c *gin.Context) {
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
		CaseType    string   `json:"case_type" binding:"required,oneof=overall change acceptance ai"`
		CaseIDOrder []string `json:"case_id_order" binding:"required"` // case_id数组（按拖拽后的顺序）
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	err = h.service.ReorderCasesByDrag(uint(projectID), userID, req.CaseType, req.CaseIDOrder)
	if err != nil {
		log.Printf("[Drag Reorder Failed] user_id=%d, project_id=%d, type=%s, error=%v", userID, projectID, req.CaseType, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "拖拽排序失败")
		return
	}

	log.Printf("[Drag Reorder Success] user_id=%d, project_id=%d, type=%s, count=%d", userID, projectID, req.CaseType, len(req.CaseIDOrder))
	utils.SuccessResponse(c, gin.H{"message": "拖拽排序成功"})
}

// ReorderAllCasesByID 按现有ID顺序重新编号所有用例
func (h *ManualCasesHandler) ReorderAllCasesByID(c *gin.Context) {
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
		CaseType string `json:"case_type" binding:"required,oneof=overall change acceptance ai"`
		Language string `json:"language" binding:"required,oneof=中文 English 日本語"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	count, err := h.service.ReorderAllCasesByID(uint(projectID), userID, req.CaseType, req.Language)
	if err != nil {
		log.Printf("[Reorder All Failed] user_id=%d, project_id=%d, type=%s, error=%v", userID, projectID, req.CaseType, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "重新编号失败")
		return
	}

	log.Printf("[Reorder All Success] user_id=%d, project_id=%d, type=%s, count=%d", userID, projectID, req.CaseType, count)
	utils.SuccessResponse(c, gin.H{
		"message": "重新编号成功",
		"count":   count,
	})
}

// ClearAICases 清空AI用例
func (h *ManualCasesHandler) ClearAICases(c *gin.Context) {
	// 获取用户ID (注意: 中间件设置的键名是 "userID" 不是 "user_id")
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

	// 调用服务
	log.Printf("[Clear AI Cases Start] user_id=%d, project_id=%d", userID, projectID)
	deletedCount, err := h.service.ClearAICases(uint(projectID), userID)
	if err != nil {
		log.Printf("[Clear AI Cases Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "清空AI用例失败")
		return
	}

	log.Printf("[Clear AI Cases Success] user_id=%d, project_id=%d, deleted_count=%d", userID, projectID, deletedCount)
	utils.SuccessResponse(c, gin.H{
		"message":       "清空成功",
		"deleted_count": deletedCount,
	})
}

// InsertCase 在指定位置插入用例
// POST /api/v1/projects/:id/manual-cases/insert
func (h *ManualCasesHandler) InsertCase(c *gin.Context) {
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
		CaseType     string `json:"case_type" binding:"required,oneof=overall change acceptance ai"`
		Position     string `json:"position" binding:"required,oneof=before after"`
		TargetCaseID string `json:"target_case_id" binding:"required"`
		Language     string `json:"language" binding:"omitempty,oneof=中文 English 日本語"`
		CaseGroup    string `json:"case_group"` // 用例集名称（可选）
	}
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 对于AI用例,固定使用中文;对于其他用例类型,语言参数必需
	if req.CaseType == "ai" {
		if req.Language == "" {
			req.Language = "中文"
		}
	} else {
		if req.Language == "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "该用例类型必须指定语言参数")
			return
		}
	}

	// 调用服务
	log.Printf("[Insert Case Start] user_id=%d, project_id=%d, type=%s, position=%s, target=%s",
		userID, projectID, req.CaseType, req.Position, req.TargetCaseID)
	newCase, err := h.service.InsertCase(uint(projectID), userID, req.CaseType, req.Position, req.TargetCaseID, req.Language)
	if err != nil {
		log.Printf("[Insert Case Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "插入用例失败")
		return
	}

	log.Printf("[Insert Case Success] user_id=%d, project_id=%d, new_case_id=%s", userID, projectID, newCase.CaseID)
	utils.SuccessResponse(c, newCase)
}

// BatchDeleteCases 批量删除用例
// POST /api/v1/projects/:id/manual-cases/batch-delete
func (h *ManualCasesHandler) BatchDeleteCases(c *gin.Context) {
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
		CaseType string   `json:"case_type" binding:"required,oneof=overall change acceptance ai"`
		CaseIDs  []string `json:"case_ids" binding:"required,min=1"`
	}
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	log.Printf("[Batch Delete Start] user_id=%d, project_id=%d, type=%s, count=%d",
		userID, projectID, req.CaseType, len(req.CaseIDs))
	deletedCount, failedCaseIDs, err := h.service.BatchDeleteCases(uint(projectID), userID, req.CaseType, req.CaseIDs)
	if err != nil {
		log.Printf("[Batch Delete Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量删除失败")
		return
	}

	log.Printf("[Batch Delete Success] user_id=%d, project_id=%d, deleted=%d, failed=%d",
		userID, projectID, deletedCount, len(failedCaseIDs))
	utils.SuccessResponse(c, gin.H{
		"message":         "批量删除完成",
		"deleted_count":   deletedCount,
		"failed_case_ids": failedCaseIDs,
	})
}

// ReassignIDs 重新分配用例ID
// POST /api/v1/projects/:id/manual-cases/reassign-ids
func (h *ManualCasesHandler) ReassignIDs(c *gin.Context) {
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
		CaseType string `json:"caseType" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务重新分配ID
	log.Printf("[Reassign IDs Start] user_id=%d, project_id=%d, type=%s", userID, projectID, req.CaseType)
	if err := h.service.ReassignAllIDs(uint(projectID), userID, req.CaseType); err != nil {
		log.Printf("[Reassign IDs Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "重新分配ID失败")
		return
	}

	log.Printf("[Reassign IDs Success] user_id=%d, project_id=%d, type=%s", userID, projectID, req.CaseType)
	utils.SuccessResponse(c, gin.H{
		"message": "重新分配ID成功",
	})
}

// SaveMultiLangVersion 保存多语言版本
// POST /api/v1/projects/:id/manual-cases/save-version
func (h *ManualCasesHandler) SaveMultiLangVersion(c *gin.Context) {
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

	// TODO: SaveMultiLangVersion方法未实现，使用SaveVersion代替
	// 调用版本服务保存多语言版本
	log.Printf("[SaveVersion Start] user_id=%d, project_id=%d", userID, projectID)
	filename, err := h.versionService.SaveVersion(uint(projectID), userID, "overall") // 使用overall类型
	if err != nil {
		log.Printf("[SaveVersion Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "版本保存失败: "+err.Error())
		return
	}

	log.Printf("[SaveVersion Success] user_id=%d, project_id=%d, filename=%s", userID, projectID, filename)
	utils.SuccessResponse(c, gin.H{
		"filename": filename,
		"message":  "版本保存成功",
	})
}

// CreateCaseForGroup 为指定的用例集创建用例
// POST /api/v1/projects/:id/case-groups/:groupId/manual-cases
func (h *ManualCasesHandler) CreateCaseForGroup(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用例集ID
	groupIDStr := c.Param("groupId")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用例集ID")
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
	var req services.CreateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 如果没有指定case_type，默认为overall
	if req.CaseType == "" {
		req.CaseType = "overall"
	}

	// TODO: GetCaseGroupName方法未实现，直接使用groupID作为caseGroup参数
	// 根据groupID获取用例集的group_name
	groupName := "" // 暂时使用空值，后续实现GetCaseGroupName方法

	// 设置请求中的case_group为组名称
	req.CaseGroup = groupName

	// 调用服务创建用例
	caseDTO, err := h.service.CreateCase(uint(projectID), userID, req)
	if err != nil {
		log.Printf("[Case Create For Group Failed] user_id=%d, project_id=%d, group_id=%d, error=%v", userID, projectID, groupID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建用例失败")
		return
	}

	log.Printf("[Case Create For Group] user_id=%d, project_id=%d, group_id=%d, group_name=%s, case_id=%s", userID, projectID, groupID, groupName, caseDTO.CaseID)
	utils.SuccessResponse(c, caseDTO)
}

// UpdateCaseForGroup 为指定的用例集更新用例
// PUT /api/v1/projects/:id/case-groups/:groupId/manual-cases/:caseId
func (h *ManualCasesHandler) UpdateCaseForGroup(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用例集ID
	groupIDStr := c.Param("groupId")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用例集ID")
		return
	}

	// 获取用例ID（整数）
	caseIDStr := c.Param("caseId")
	caseID, err := strconv.ParseUint(caseIDStr, 10, 32)
	if err != nil {
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
	var req services.UpdateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// TODO: UpdateCaseByID方法未实现，需要先查询case_id然后使用UpdateCase
	// 查询用例获取case_id
	// 暂时跳过此操作
	err = fmt.Errorf("UpdateCaseByID not implemented, please use UpdateCase with case_id")
	if err != nil {
		log.Printf("[Case Update For Group Failed] user_id=%d, project_id=%d, group_id=%d, case_id=%d, error=%v", userID, projectID, groupID, caseID, err)
		if err.Error() == "无项目访问权限" {
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

	log.Printf("[Case Update For Group] user_id=%d, project_id=%d, group_id=%d, case_id=%d", userID, projectID, groupID, caseID)
	utils.SuccessResponse(c, map[string]interface{}{
		"code":    0,
		"message": "更新成功",
		"data": map[string]interface{}{
			"case_id": caseID,
		},
	})
}
