package handlers

import (
	"log"
	"os"
	"strconv"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// RawDocumentHandler 原始文档处理器接口
type RawDocumentHandler interface {
	Upload(c *gin.Context)
	List(c *gin.Context)
	Convert(c *gin.Context)
	GetConvertStatus(c *gin.Context)
	DownloadOriginal(c *gin.Context)
	DownloadConverted(c *gin.Context)
	PreviewConverted(c *gin.Context)
	DeleteOriginal(c *gin.Context)
	DeleteConverted(c *gin.Context)
}

type rawDocumentHandler struct {
	documentService services.RawDocumentService
}

// NewRawDocumentHandler 创建原始文档处理器实例
func NewRawDocumentHandler(documentService services.RawDocumentService) RawDocumentHandler {
	return &rawDocumentHandler{
		documentService: documentService,
	}
}

// Upload 上传原始文档
// POST /api/v1/projects/:id/raw-documents
func (h *rawDocumentHandler) Upload(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid project id")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}
	userID := userIDVal.(uint)

	file, err := c.FormFile("file")
	if err != nil {
		utils.ResponseError(c, 400, "file is required")
		return
	}

	result, err := h.documentService.Upload(uint(projectID), userID, file)
	if err != nil {
		log.Printf("[RawDocument Upload Failed] project_id=%d, user_id=%d, error=%v", projectID, userID, err)
		if err.Error() == "file type not allowed" {
			utils.ResponseError(c, 400, "file type not supported")
			return
		}
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, result)
}

// List 获取原始文档列表
// GET /api/v1/projects/:id/raw-documents
func (h *rawDocumentHandler) List(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid project id")
		return
	}

	documents, err := h.documentService.List(uint(projectID))
	if err != nil {
		log.Printf("[RawDocument List Failed] project_id=%d, error=%v", projectID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"documents": documents,
		"total":     len(documents),
	})
}

// Convert 启动文档转换
// POST /api/v1/raw-documents/:id/convert
func (h *rawDocumentHandler) Convert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid document id")
		return
	}

	result, err := h.documentService.StartConvert(uint(id))
	if err != nil {
		log.Printf("[RawDocument Convert Failed] document_id=%d, error=%v", id, err)
		if err.Error() == "document not found" {
			utils.ResponseError(c, 404, "document not found")
			return
		}
		if err.Error() == "document conversion already in progress" {
			utils.ResponseError(c, 409, "document conversion already in progress")
			return
		}
		utils.ResponseError(c, 500, "internal server error")
		return
	}

	log.Printf("[Convert Triggered] document_id=%d, task_id=%s", id, result.TaskID)
	utils.ResponseSuccess(c, result)
}

// GetConvertStatus 查询转换状态
// GET /api/v1/raw-documents/:id/convert-status
func (h *rawDocumentHandler) GetConvertStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid document id")
		return
	}

	status, err := h.documentService.GetConvertStatus(uint(id))
	if err != nil {
		log.Printf("[RawDocument GetConvertStatus Failed] document_id=%d, error=%v", id, err)
		if err.Error() == "document not found" {
			utils.ResponseError(c, 404, "document not found")
			return
		}
		utils.ResponseError(c, 500, "internal server error")
		return
	}

	utils.ResponseSuccess(c, status)
}

// DownloadOriginal 下载原始文档
// GET /api/v1/raw-documents/:id/download
func (h *rawDocumentHandler) DownloadOriginal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid document id")
		return
	}

	doc, file, err := h.documentService.DownloadOriginal(uint(id))
	if err != nil {
		if err.Error() == "document not found" || err.Error() == "document file not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		log.Printf("[RawDocument Download Failed] id=%d, error=%v", id, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}
	defer file.Close()

	c.Header("Content-Type", doc.MimeType)
	c.Header("Content-Disposition", "attachment; filename="+doc.OriginalFilename)
	c.Header("Content-Length", strconv.FormatInt(doc.FileSize, 10))

	c.DataFromReader(200, doc.FileSize, doc.MimeType, file, nil)
}

// DownloadConverted 下载转换后的文档
// GET /api/v1/raw-documents/:id/converted/download
func (h *rawDocumentHandler) DownloadConverted(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid document id")
		return
	}

	doc, file, err := h.documentService.DownloadConverted(uint(id))
	if err != nil {
		if err.Error() == "document not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		if err.Error() == "document not converted yet" || err.Error() == "converted file not found" {
			utils.ResponseError(c, 404, "converted file not available")
			return
		}
		log.Printf("[RawDocument DownloadConverted Failed] id=%d, error=%v", id, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}
	defer file.Close()

	// 获取转换后文件的大小
	fileInfo, _ := file.(*os.File).Stat()
	fileSize := fileInfo.Size()

	c.Header("Content-Type", "text/markdown")
	c.Header("Content-Disposition", "attachment; filename="+doc.ConvertedFilename)
	c.Header("Content-Length", strconv.FormatInt(fileSize, 10))

	c.DataFromReader(200, fileSize, "text/markdown", file, nil)
}

// PreviewConverted 预览转换后的Markdown文档
// GET /api/v1/raw-documents/:id/converted/preview
func (h *rawDocumentHandler) PreviewConverted(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid document id")
		return
	}

	doc, content, err := h.documentService.PreviewConverted(uint(id))
	if err != nil {
		if err.Error() == "document not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		if err.Error() == "document not converted yet" || err.Error() == "converted file not found" {
			utils.ResponseError(c, 404, "converted file not available")
			return
		}
		log.Printf("[RawDocument PreviewConverted Failed] id=%d, error=%v", id, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"filename": doc.ConvertedFilename,
		"content":  content,
	})
}

// DeleteOriginal 删除原始文档（包括转换文件）
// DELETE /api/v1/raw-documents/:id
func (h *rawDocumentHandler) DeleteOriginal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid document id")
		return
	}

	err = h.documentService.DeleteOriginal(uint(id))
	if err != nil {
		log.Printf("[RawDocument Delete Failed] id=%d, error=%v", id, err)
		if err.Error() == "document not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"message": "document deleted successfully",
	})
}

// DeleteConverted 仅删除转换后的文档
// DELETE /api/v1/raw-documents/:id/converted
func (h *rawDocumentHandler) DeleteConverted(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid document id")
		return
	}

	err = h.documentService.DeleteConverted(uint(id))
	if err != nil {
		log.Printf("[RawDocument DeleteConverted Failed] id=%d, error=%v", id, err)
		if err.Error() == "document not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"message": "converted file deleted successfully",
	})
}
