package handlers

import (
	"log"
	"net/http"
	"strconv"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// ExecutionTaskHandler 测试执行任务处理器
type ExecutionTaskHandler struct {
	service services.ExecutionTaskService
}

// NewExecutionTaskHandler 创建处理器实例
func NewExecutionTaskHandler(service services.ExecutionTaskService) *ExecutionTaskHandler {
	return &ExecutionTaskHandler{service: service}
}

// GetTasks 获取任务列表
// GET /api/v1/projects/:id/execution-tasks
func (h *ExecutionTaskHandler) GetTasks(c *gin.Context) {
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

	// 调用服务
	tasks, err := h.service.GetTasksByProject(uint(projectID), userID)
	if err != nil {
		log.Printf("[ExecutionTask Get Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取任务列表失败")
		return
	}

	log.Printf("[ExecutionTask Get] user_id=%d, project_id=%d, count=%d", userID, projectID, len(tasks))
	utils.SuccessResponse(c, tasks)
}

// CreateTask 创建新任务
// POST /api/v1/projects/:id/execution-tasks
func (h *ExecutionTaskHandler) CreateTask(c *gin.Context) {
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
	var req services.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ExecutionTask Create Bind Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	task, err := h.service.CreateTask(uint(projectID), userID, req)
	if err != nil {
		log.Printf("[ExecutionTask Create Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		if err.Error() == "任务名已存在" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建任务失败")
		return
	}

	log.Printf("[ExecutionTask Create] user_id=%d, project_id=%d, task_uuid=%s, task_name=%s", userID, projectID, task.TaskUUID, task.TaskName)
	utils.SuccessResponse(c, task)
}

// UpdateTask 更新任务
// PUT /api/v1/projects/:id/execution-tasks/:task_uuid
func (h *ExecutionTaskHandler) UpdateTask(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取任务UUID
	taskUUID := c.Param("task_uuid")
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

	// 解析请求体
	var req services.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ExecutionTask Update Bind Failed] user_id=%d, project_id=%d, task_uuid=%s, error=%v", userID, projectID, taskUUID, err)
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务
	task, err := h.service.UpdateTask(uint(projectID), userID, taskUUID, req)
	if err != nil {
		log.Printf("[ExecutionTask Update Failed] user_id=%d, project_id=%d, task_uuid=%s, error=%v", userID, projectID, taskUUID, err)
		if err.Error() == "任务不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "任务不属于该项目" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "任务名已存在" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		if err.Error() == "结束日期不能早于开始日期" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新任务失败")
		return
	}

	log.Printf("[ExecutionTask Update] user_id=%d, project_id=%d, task_uuid=%s", userID, projectID, taskUUID)
	utils.SuccessResponse(c, task)
}

// DeleteTask 删除任务
// DELETE /api/v1/projects/:id/execution-tasks/:task_uuid
func (h *ExecutionTaskHandler) DeleteTask(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取任务UUID
	taskUUID := c.Param("task_uuid")
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
	err = h.service.DeleteTask(uint(projectID), userID, taskUUID)
	if err != nil {
		log.Printf("[ExecutionTask Delete Failed] user_id=%d, project_id=%d, task_uuid=%s, error=%v", userID, projectID, taskUUID, err)
		if err.Error() == "任务不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "任务不属于该项目" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除任务失败")
		return
	}

	log.Printf("[ExecutionTask Delete] user_id=%d, project_id=%d, task_uuid=%s", userID, projectID, taskUUID)
	utils.MessageResponse(c, http.StatusOK, "任务已删除")
}

// ExecuteTask 执行测试任务
// POST /api/v1/projects/:id/execution-tasks/:task_uuid/execute
func (h *ExecutionTaskHandler) ExecuteTask(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取任务UUID
	taskUUID := c.Param("task_uuid")
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

	// 调用服务执行
	result, err := h.service.ExecuteTask(uint(projectID), userID, taskUUID)
	if err != nil {
		log.Printf("[ExecutionTask Execute Failed] user_id=%d, project_id=%d, task_uuid=%s, error=%v",
			userID, projectID, taskUUID, err)

		if err.Error() == "任务不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "手工测试类型不支持自动执行" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "没有可执行的用例" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "执行任务失败")
		return
	}

	log.Printf("[ExecutionTask Execute] user_id=%d, project_id=%d, task_uuid=%s, total=%d, ok=%d, ng=%d",
		userID, projectID, taskUUID, result.Total, result.OKCount, result.NGCount)
	utils.SuccessResponse(c, result)
}

// ExecuteSingleCase 执行单条测试用例
// POST /api/v1/projects/:id/execution-tasks/:task_uuid/cases/:case_result_id/execute
func (h *ExecutionTaskHandler) ExecuteSingleCase(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取任务UUID
	taskUUID := c.Param("task_uuid")
	if taskUUID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "任务UUID不能为空")
		return
	}

	// 获取用例结果ID
	caseResultIDStr := c.Param("case_result_id")
	caseResultID, err := strconv.ParseUint(caseResultIDStr, 10, 32)
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

	// 调用服务执行
	result, err := h.service.ExecuteSingleCase(uint(projectID), userID, taskUUID, uint(caseResultID))
	if err != nil {
		log.Printf("[ExecutionTask ExecuteSingleCase Failed] user_id=%d, project_id=%d, task_uuid=%s, case_result_id=%d, error=%v",
			userID, projectID, taskUUID, caseResultID, err)

		if err.Error() == "任务不存在" || err.Error() == "用例不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "手工测试类型不支持自动执行" || err.Error() == "用例没有脚本代码，无法执行" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "执行用例失败")
		return
	}

	log.Printf("[ExecutionTask ExecuteSingleCase] user_id=%d, project_id=%d, task_uuid=%s, case_result_id=%d, ok_count=%d",
		userID, projectID, taskUUID, caseResultID, result.OKCount)
	utils.SuccessResponse(c, result)
}
