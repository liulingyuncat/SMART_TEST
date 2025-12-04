package handlers

import (
	"log"
	"strconv"
	"webtest/internal/models"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// DefectCommentHandler 缺陷说明处理器接口
type DefectCommentHandler interface {
	GetComments(c *gin.Context)
	CreateComment(c *gin.Context)
	UpdateComment(c *gin.Context)
	DeleteComment(c *gin.Context)
}

type defectCommentHandler struct {
	commentService services.DefectCommentService
}

// NewDefectCommentHandler 创建缺陷说明处理器实例
func NewDefectCommentHandler(commentService services.DefectCommentService) DefectCommentHandler {
	return &defectCommentHandler{
		commentService: commentService,
	}
}

// GetComments 获取缺陷说明列表
// GET /api/v1/projects/:id/defects/:defectId/comments
func (h *defectCommentHandler) GetComments(c *gin.Context) {
	defectID := c.Param("defectId")
	if defectID == "" {
		utils.ResponseError(c, 400, "defect id is required")
		return
	}

	result, err := h.commentService.List(defectID)
	if err != nil {
		log.Printf("[DefectComment List Failed] defect_id=%s, error=%v", defectID, err)
		utils.ResponseError(c, 404, "defect not found or failed to list comments")
		return
	}

	utils.ResponseSuccess(c, result)
}

// CreateComment 创建缺陷说明
// POST /api/v1/projects/:id/defects/:defectId/comments
func (h *defectCommentHandler) CreateComment(c *gin.Context) {
	defectID := c.Param("defectId")
	if defectID == "" {
		utils.ResponseError(c, 400, "defect id is required")
		return
	}

	// 获取当前用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}
	userID := userIDVal.(uint)

	// 解析请求体
	var req models.DefectCommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "validation failed: "+err.Error())
		return
	}

	// 创建说明
	comment, err := h.commentService.Create(defectID, userID, &req)
	if err != nil {
		log.Printf("[DefectComment Create Failed] defect_id=%s, user_id=%d, error=%v", defectID, userID, err)
		if err.Error() == "defect not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		utils.ResponseError(c, 400, err.Error())
		return
	}

	c.JSON(201, gin.H{
		"code": 201,
		"data": gin.H{
			"comment": comment,
		},
	})
}

// UpdateComment 更新缺陷说明
// PUT /api/v1/projects/:id/defects/:defectId/comments/:commentId
func (h *defectCommentHandler) UpdateComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid comment id")
		return
	}

	// 获取当前用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}
	userID := userIDVal.(uint)

	// 解析请求体
	var req models.DefectCommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "validation failed: "+err.Error())
		return
	}

	// 更新说明
	if err := h.commentService.Update(uint(commentID), userID, &req); err != nil {
		log.Printf("[DefectComment Update Failed] comment_id=%d, user_id=%d, error=%v", commentID, userID, err)
		if err.Error() == "comment not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		if err.Error() == "permission denied: only creator can edit" {
			utils.ResponseError(c, 403, err.Error())
			return
		}
		utils.ResponseError(c, 400, err.Error())
		return
	}

	// 获取更新后的说明
	comment, err := h.commentService.GetByID(uint(commentID))
	if err != nil {
		utils.ResponseError(c, 400, "failed to get updated comment")
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"comment": comment,
	})
}

// DeleteComment 删除缺陷说明
// DELETE /api/v1/projects/:id/defects/:defectId/comments/:commentId
func (h *defectCommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid comment id")
		return
	}

	// 获取当前用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}
	userID := userIDVal.(uint)

	// 删除说明
	if err := h.commentService.Delete(uint(commentID), userID); err != nil {
		log.Printf("[DefectComment Delete Failed] comment_id=%d, user_id=%d, error=%v", commentID, userID, err)
		if err.Error() == "comment not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		if err.Error() == "permission denied: only creator can delete" {
			utils.ResponseError(c, 403, err.Error())
			return
		}
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"message": "success",
	})
}
