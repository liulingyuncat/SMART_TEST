package handlers

import (
	"net/http"
	"strconv"
	"webtest/internal/constants"
	"webtest/internal/services"

	"github.com/gin-gonic/gin"
)

// ReviewHandler 评审管理处理器
type ReviewHandler struct {
	reviewService services.ReviewService
}

// NewReviewHandler 创建评审管理处理器实例
func NewReviewHandler(reviewService services.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// GetCaseReview 获取评审内容
// @Summary 获取评审内容
// @Tags Review
// @Param projectID path int true "项目ID"
// @Param caseType query string true "用例类型(overall/change/ai)"
// @Success 200 {object} map[string]string "评审内容"
// @Router /api/manual-cases/:projectID/review [get]
func (h *ReviewHandler) GetCaseReview(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	caseType := c.Query("caseType")
	if caseType != "overall" && caseType != "change" && caseType != "ai" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	content, err := h.reviewService.GetReview(uint(projectID), caseType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评审内容失败", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": content})
}

// SaveCaseReview 保存评审内容(UPSERT)
// @Summary 保存评审内容
// @Tags Review
// @Accept json
// @Param projectID path int true "项目ID"
// @Param body body object true "请求体 {caseType, content}"
// @Success 200 {object} map[string]string "保存结果"
// @Router /api/manual-cases/:projectID/review [post]
func (h *ReviewHandler) SaveCaseReview(c *gin.Context) {
	projectIDStr := c.Param("id")
	println("[DEBUG] SaveCaseReview - Raw projectID param:", projectIDStr)

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		println("[ERROR] SaveCaseReview - ParseUint error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}
	println("[DEBUG] SaveCaseReview - Parsed projectID:", projectID)

	var req struct {
		CaseType string `json:"caseType" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		println("[ERROR] SaveCaseReview - ShouldBindJSON error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}
	println("[DEBUG] SaveCaseReview - CaseType:", req.CaseType, "ContentLength:", len(req.Content))

	if req.CaseType != "overall" && req.CaseType != "change" && req.CaseType != "ai" {
		println("[ERROR] SaveCaseReview - Invalid caseType:", req.CaseType)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	err = h.reviewService.SaveReview(uint(projectID), req.CaseType, req.Content)
	if err != nil {
		println("[ERROR] SaveCaseReview - SaveReview error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存评审内容失败", "details": err.Error()})
		return
	}

	println("[DEBUG] SaveCaseReview - Success")
	c.JSON(http.StatusOK, gin.H{"message": "评审内容保存成功"})
}
