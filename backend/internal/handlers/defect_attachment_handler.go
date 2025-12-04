package handlers

import (
	"log"
	"strconv"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// DefectAttachmentHandler 缺陷附件处理器接口
type DefectAttachmentHandler interface {
	List(c *gin.Context)
	Upload(c *gin.Context)
	Download(c *gin.Context)
	Delete(c *gin.Context)
}

type defectAttachmentHandler struct {
	attachmentService services.DefectAttachmentService
}

// NewDefectAttachmentHandler 创建缺陷附件处理器实例
func NewDefectAttachmentHandler(attachmentService services.DefectAttachmentService) DefectAttachmentHandler {
	return &defectAttachmentHandler{
		attachmentService: attachmentService,
	}
}

// List 获取附件列表
// GET /api/v1/projects/:id/defects/:defectId/attachments
func (h *defectAttachmentHandler) List(c *gin.Context) {
	defectID := c.Param("defectId")

	attachments, err := h.attachmentService.ListByDefectID(defectID)
	if err != nil {
		log.Printf("[Attachment List Failed] defect_id=%s, error=%v", defectID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"attachments": attachments,
		"total":       len(attachments),
	})
}

// Upload 上传附件
// POST /api/v1/defects/:defectId/attachments
func (h *defectAttachmentHandler) Upload(c *gin.Context) {
	defectID := c.Param("defectId")

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}
	userID := userIDVal.(uint)

	// 获取项目ID（从上下文或请求参数）
	projectIDVal, exists := c.Get("projectID")
	var projectID uint
	if exists {
		projectID = projectIDVal.(uint)
	} else {
		// 尝试从查询参数获取
		projectIDStr := c.Query("project_id")
		if projectIDStr != "" {
			pid, err := strconv.ParseUint(projectIDStr, 10, 32)
			if err == nil {
				projectID = uint(pid)
			}
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.ResponseError(c, 400, "file is required")
		return
	}

	result, err := h.attachmentService.Upload(defectID, projectID, userID, file)
	if err != nil {
		log.Printf("[Attachment Upload Failed] defect_id=%s, user_id=%d, error=%v", defectID, userID, err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, result)
}

// Download 下载附件
// GET /api/v1/defects/:defectId/attachments/:attId
func (h *defectAttachmentHandler) Download(c *gin.Context) {
	attIDStr := c.Param("attId")
	attID, err := strconv.ParseUint(attIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid attachment id")
		return
	}

	attachment, file, err := h.attachmentService.Download(uint(attID))
	if err != nil {
		if err.Error() == "attachment not found" || err.Error() == "attachment file not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		log.Printf("[Attachment Download Failed] att_id=%d, error=%v", attID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}
	defer file.Close()

	c.Header("Content-Type", attachment.MimeType)
	c.Header("Content-Disposition", "attachment; filename="+attachment.FileName)
	c.Header("Content-Length", strconv.FormatInt(attachment.FileSize, 10))

	c.DataFromReader(200, attachment.FileSize, attachment.MimeType, file, nil)
}

// Delete 删除附件
// DELETE /api/v1/defects/:defectId/attachments/:attId
func (h *defectAttachmentHandler) Delete(c *gin.Context) {
	attIDStr := c.Param("attId")
	attID, err := strconv.ParseUint(attIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid attachment id")
		return
	}

	if err := h.attachmentService.Delete(uint(attID)); err != nil {
		if err.Error() == "attachment not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		log.Printf("[Attachment Delete Failed] att_id=%d, error=%v", attID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "attachment deleted successfully"})
}
