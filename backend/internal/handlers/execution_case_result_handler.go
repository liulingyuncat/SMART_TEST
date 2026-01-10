package handlers

import (
	"log"
	"net/http"
	"strconv"
	"webtest/internal/models"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// ExecutionCaseResultHandler 测试执行用例结果处理器
type ExecutionCaseResultHandler struct {
	service     services.ExecutionCaseResultService
	taskService services.ExecutionTaskService
}

// NewExecutionCaseResultHandler 创建处理器实例
func NewExecutionCaseResultHandler(
	service services.ExecutionCaseResultService,
	taskService services.ExecutionTaskService,
) *ExecutionCaseResultHandler {
	return &ExecutionCaseResultHandler{
		service:     service,
		taskService: taskService,
	}
}

// GetExecutionCaseResults 获取任务的执行结果列表
// GET /api/v1/execution-tasks/:taskUuid/case-results
func (h *ExecutionCaseResultHandler) GetExecutionCaseResults(c *gin.Context) {
	// 获取任务UUID
	taskUUID := c.Param("taskUuid")
	log.Printf("[ExecutionCaseResult Get] Starting - taskUUID=%s", taskUUID)

	if taskUUID == "" {
		log.Printf("[ExecutionCaseResult Get] ERROR: taskUUID is empty")
		utils.ErrorResponse(c, http.StatusBadRequest, "任务UUID不能为空")
		return
	}

	// 获取用户ID(中间件已验证)
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[ExecutionCaseResult Get] ERROR: userID not found")
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)
	log.Printf("[ExecutionCaseResult Get] userID=%d, taskUUID=%s", userID, taskUUID)

	// 调用服务
	results, err := h.service.GetCaseResults(taskUUID)
	if err != nil {
		log.Printf("[ExecutionCaseResult Get] Service error: user_id=%d, task_uuid=%s, error=%v", userID, taskUUID, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取执行结果失败: "+err.Error())
		return
	}

	// 如果没有结果，返回空数组而不是null
	if results == nil {
		results = []*models.ExecutionCaseResult{}
	}

	log.Printf("[ExecutionCaseResult Get] SUCCESS: user_id=%d, task_uuid=%s, count=%d", userID, taskUUID, len(results))
	utils.SuccessResponse(c, results)
}

// SaveExecutionCaseResults 保存或更新执行结果
// PATCH /api/v1/execution-tasks/:taskUuid/case-results
func (h *ExecutionCaseResultHandler) SaveExecutionCaseResults(c *gin.Context) {
	// 获取任务UUID
	taskUUID := c.Param("taskUuid")
	log.Printf("[ExecutionCaseResult Save] Starting - taskUUID=%s", taskUUID)

	if taskUUID == "" {
		log.Printf("[ExecutionCaseResult Save] ERROR: taskUUID is empty")
		utils.ErrorResponse(c, http.StatusBadRequest, "任务UUID不能为空")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[ExecutionCaseResult Save] ERROR: userID not found in context")
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)
	log.Printf("[ExecutionCaseResult Save] userID=%d, taskUUID=%s", userID, taskUUID)

	// 解析请求体
	var requests []services.SaveCaseResultRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		log.Printf("[ExecutionCaseResult Save] Bind Failed: user_id=%d, task_uuid=%s, error=%v", userID, taskUUID, err)
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	log.Printf("[ExecutionCaseResult Save] Parsed %d requests", len(requests))
	if len(requests) > 0 {
		log.Printf("[ExecutionCaseResult Save] First request: case_id=%s, case_num=%s, display_id=%d",
			requests[0].CaseID, requests[0].CaseNum, requests[0].DisplayID)
	}

	// 调用服务
	err := h.service.SaveCaseResults(taskUUID, userID, requests)
	if err != nil {
		log.Printf("[ExecutionCaseResult Save] Service Failed: user_id=%d, task_uuid=%s, count=%d, error=%v", userID, taskUUID, len(requests), err)
		if err.Error() == "requests array is empty" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "保存执行结果失败: "+err.Error())
		return
	}

	log.Printf("[ExecutionCaseResult Save] SUCCESS: user_id=%d, task_uuid=%s, count=%d", userID, taskUUID, len(requests))
	utils.MessageResponse(c, http.StatusOK, "保存成功")
}

// GetExecutionStatistics 获取任务的统计信息
// GET /api/v1/execution-tasks/:taskUuid/statistics
func (h *ExecutionCaseResultHandler) GetExecutionStatistics(c *gin.Context) {
	// 获取任务UUID
	taskUUID := c.Param("taskUuid")
	if taskUUID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "任务UUID不能为空")
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
	stats, err := h.service.GetStatistics(taskUUID)
	if err != nil {
		log.Printf("[ExecutionStatistics Get Failed] user_id=%d, task_uuid=%s, error=%v", userID, taskUUID, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取统计信息失败")
		return
	}

	log.Printf("[ExecutionStatistics Get] user_id=%d, task_uuid=%s, total=%d", userID, taskUUID, stats["total"])
	utils.SuccessResponse(c, stats)
}

// InitExecutionResults 初始化任务的执行结果
// POST /api/v1/execution-tasks/:taskUuid/case-results/init
func (h *ExecutionCaseResultHandler) InitExecutionResults(c *gin.Context) {
	// 获取任务UUID
	taskUUID := c.Param("taskUuid")
	if taskUUID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "任务UUID不能为空")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 解析请求体(包含projectID、executionType和可选的caseGroupID)
	type InitRequest struct {
		ProjectID     uint   `json:"project_id" binding:"required"`
		ExecutionType string `json:"execution_type" binding:"required,oneof=manual automation api"`
		CaseGroupID   uint   `json:"case_group_id"`   // 可选：指定用例集ID，仅导入该用例集的用例
		CaseGroupName string `json:"case_group_name"` // 可选：指定用例集名称
	}
	var req InitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ExecutionResult Init Bind Failed] user_id=%d, task_uuid=%s, error=%v", userID, taskUUID, err)
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务（传入用例集ID和名称）
	err := h.service.InitTaskResults(taskUUID, req.ProjectID, req.ExecutionType, userID, req.CaseGroupID, req.CaseGroupName)
	if err != nil {
		log.Printf("[ExecutionResult Init Failed] user_id=%d, task_uuid=%s, project_id=%d, type=%s, case_group_id=%d, error=%v",
			userID, taskUUID, req.ProjectID, req.ExecutionType, req.CaseGroupID, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "初始化执行结果失败: "+err.Error())
		return
	}

	log.Printf("[ExecutionResult Init] user_id=%d, task_uuid=%s, project_id=%d, type=%s, case_group_id=%d",
		userID, taskUUID, req.ProjectID, req.ExecutionType, req.CaseGroupID)
	utils.MessageResponse(c, http.StatusOK, "初始化成功")
}

// UpdateSingleCaseResult 更新单个用例的执行结果
// PUT /api/v1/execution-task-cases/:id
func (h *ExecutionCaseResultHandler) UpdateSingleCaseResult(c *gin.Context) {
	// 获取记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("[ExecutionCaseResult UpdateSingle] ERROR: invalid id=%s", idStr)
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[ExecutionCaseResult UpdateSingle] ERROR: userID not found")
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 解析请求体
	type UpdateRequest struct {
		Result       string `json:"result" binding:"required,oneof=NR OK NG Block"`
		Comment      string `json:"comment"`
		BugID        string `json:"bug_id"`
		ResponseTime string `json:"response_time"`
	}
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ExecutionCaseResult UpdateSingle] Bind Failed: id=%d, error=%v", id, err)
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	log.Printf("[ExecutionCaseResult UpdateSingle] id=%d, result=%s, comment=%s, bug_id=%s, response_time=%s, userID=%d", id, req.Result, req.Comment, req.BugID, req.ResponseTime, userID)

	// 调用服务（传入bug_id和response_time）
	err = h.service.UpdateSingleResultWithBugIDAndResponseTime(uint(id), req.Result, req.Comment, req.BugID, req.ResponseTime, userID)
	if err != nil {
		log.Printf("[ExecutionCaseResult UpdateSingle] Service Failed: id=%d, error=%v", id, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新执行结果失败: "+err.Error())
		return
	}

	log.Printf("[ExecutionCaseResult UpdateSingle] SUCCESS: id=%d", id)
	utils.SuccessResponse(c, map[string]interface{}{
		"id":      id,
		"result":  req.Result,
		"message": "更新成功",
	})
}

// BatchUpdateCaseResults 批量更新用例的执行结果
// PUT /api/v1/execution-task-cases/batch
func (h *ExecutionCaseResultHandler) BatchUpdateCaseResults(c *gin.Context) {
	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[ExecutionCaseResult BatchUpdate] ERROR: userID not found")
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 解析请求体
	var requests []services.UpdateResultRequestWithBugID
	if err := c.ShouldBindJSON(&requests); err != nil {
		log.Printf("[ExecutionCaseResult BatchUpdate] Bind Failed: error=%v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	if len(requests) == 0 {
		log.Printf("[ExecutionCaseResult BatchUpdate] ERROR: empty requests")
		utils.ErrorResponse(c, http.StatusBadRequest, "请求数组不能为空")
		return
	}

	log.Printf("[ExecutionCaseResult BatchUpdate] userID=%d, count=%d", userID, len(requests))

	// 调用服务
	responses, err := h.service.BatchUpdateResultsWithBugID(requests, userID)
	if err != nil {
		log.Printf("[ExecutionCaseResult BatchUpdate] Service Failed: error=%v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量更新失败: "+err.Error())
		return
	}

	// 统计成功和失败数量
	successCount := 0
	failCount := 0
	for _, resp := range responses {
		if resp.Success {
			successCount++
		} else {
			failCount++
		}
	}

	log.Printf("[ExecutionCaseResult BatchUpdate] SUCCESS: total=%d, success=%d, fail=%d", len(requests), successCount, failCount)
	utils.SuccessResponse(c, map[string]interface{}{
		"total":   len(requests),
		"success": successCount,
		"fail":    failCount,
		"results": responses,
	})
}
