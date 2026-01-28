package handlers

import (
	"strconv"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// AIReportHandler AI报告处理器接口
type AIReportHandler interface {
	ListReports(c *gin.Context)
	CreateReport(c *gin.Context)
	GetReport(c *gin.Context)
	UpdateReport(c *gin.Context)
	DeleteReport(c *gin.Context)
}

// aiReportHandler AI报告处理器实现
type aiReportHandler struct {
	service services.AIReportService
}

// NewAIReportHandler 创建AI报告处理器实例
func NewAIReportHandler(service services.AIReportService) AIReportHandler {
	return &aiReportHandler{service: service}
}

// ListReports 获取报告列表
// GET /api/projects/:id/ai-reports?type=R|A|T|O
func (h *aiReportHandler) ListReports(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的项目ID")
		return
	}

	// 获取可选的type参数
	reportType := c.Query("type")

	reports, err := h.service.ListReports(uint(projectID), reportType)
	if err != nil {
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, reports)
}

// CreateReport 创建报告
// POST /api/projects/:id/ai-reports
func (h *aiReportHandler) CreateReport(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "无效的项目ID")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type"` // R=用例审阅, A=品质分析, T=测试结果, O=其他(默认)
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	report, err := h.service.CreateReport(uint(projectID), req.Type, req.Name)
	if err != nil {
		// 名称已存在返回409
		if err.Error() == "报告名称已存在" {
			utils.ResponseError(c, 409, err.Error())
		} else {
			utils.ResponseError(c, 500, err.Error())
		}
		return
	}

	utils.ResponseSuccess(c, report)
}

// GetReport 获取报告详情
// GET /api/projects/:id/ai-reports/:reportId
func (h *aiReportHandler) GetReport(c *gin.Context) {
	reportID := c.Param("reportId")

	report, err := h.service.GetReport(reportID)
	if err != nil {
		utils.ResponseError(c, 404, "报告不存在")
		return
	}

	utils.ResponseSuccess(c, report)
}

// UpdateReport 更新报告
// PUT /api/projects/:id/ai-reports/:reportId
func (h *aiReportHandler) UpdateReport(c *gin.Context) {
	reportID := c.Param("reportId")

	var req struct {
		Name    *string `json:"name"`
		Content *string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	report, err := h.service.UpdateReport(reportID, req.Name, req.Content)
	if err != nil {
		if err.Error() == "报告不存在" {
			utils.ResponseError(c, 404, err.Error())
		} else if err.Error() == "报告名称已存在" {
			utils.ResponseError(c, 409, err.Error())
		} else {
			utils.ResponseError(c, 500, err.Error())
		}
		return
	}

	utils.ResponseSuccess(c, report)
}

// DeleteReport 删除报告
// DELETE /api/projects/:id/ai-reports/:reportId
func (h *aiReportHandler) DeleteReport(c *gin.Context) {
	reportID := c.Param("reportId")

	err := h.service.DeleteReport(reportID)
	if err != nil {
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "删除成功"})
}
