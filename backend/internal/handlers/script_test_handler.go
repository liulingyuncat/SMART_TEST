package handlers

import (
	"log"
	"net/http"
	"strconv"
	"webtest/internal/services"

	"github.com/gin-gonic/gin"
)

// ScriptTestHandler 脚本测试处理器
type ScriptTestHandler struct {
	service services.ScriptTestService
}

// NewScriptTestHandler 创建处理器实例
func NewScriptTestHandler(service services.ScriptTestService) *ScriptTestHandler {
	return &ScriptTestHandler{service: service}
}

// TestScriptRequest 脚本测试请求体
type TestScriptRequest struct {
	ScriptCode string `json:"script_code" binding:"required"`
	GroupID    uint   `json:"group_id"`   // 用例集ID
	GroupType  string `json:"group_type"` // 用例集类型：web 或 api
}

// TestScript 测试脚本
// @Summary 测试脚本
// @Description 在Docker环境中测试脚本，支持变量替换
// @Tags ScriptTest
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Param body body TestScriptRequest true "脚本测试请求"
// @Success 200 {object} services.ScriptTestResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id}/script-test [post]
func (h *ScriptTestHandler) TestScript(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析请求体
	var req TestScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 调用服务
	serviceReq := services.ScriptTestRequest{
		ScriptCode: req.ScriptCode,
		GroupID:    req.GroupID,
		GroupType:  req.GroupType,
		ProjectID:  uint(projectID),
	}

	result, err := h.service.TestScript(uint(projectID), userID.(uint), serviceReq)
	if err != nil {
		log.Printf("[ScriptTest] 执行失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ScriptTest] 执行完成: success=%v, response_time=%dms", result.Success, result.ResponseTime)
	c.JSON(http.StatusOK, result)
}
