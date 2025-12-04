package handlers

import (
	"log"
	"strconv"
	"webtest/internal/models"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// DefectHandler 缺陷处理器接口
type DefectHandler interface {
	GetDefects(c *gin.Context)
	CreateDefect(c *gin.Context)
	GetDefect(c *gin.Context)
	UpdateDefect(c *gin.Context)
	DeleteDefect(c *gin.Context)
	ExportTemplate(c *gin.Context)
	ImportDefects(c *gin.Context)
	ExportDefects(c *gin.Context)
}

type defectHandler struct {
	defectService services.DefectService
}

// NewDefectHandler 创建缺陷处理器实例
func NewDefectHandler(defectService services.DefectService) DefectHandler {
	return &defectHandler{
		defectService: defectService,
	}
}

// GetDefects 获取缺陷列表
// GET /api/v1/projects/:id/defects
func (h *defectHandler) GetDefects(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid project id")
		return
	}

	// 获取查询参数
	status := c.Query("status")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))

	result, err := h.defectService.List(uint(projectID), status, keyword, page, size)
	if err != nil {
		log.Printf("[Defect List Failed] project_id=%d, error=%v", projectID, err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, result)
}

// CreateDefect 创建缺陷
// POST /api/v1/projects/:id/defects
func (h *defectHandler) CreateDefect(c *gin.Context) {
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

	var req models.DefectCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "validation failed: "+err.Error())
		return
	}

	defect, err := h.defectService.Create(uint(projectID), userID, &req)
	if err != nil {
		log.Printf("[Defect Create Failed] project_id=%d, user_id=%d, error=%v", projectID, userID, err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccessWithCode(c, 201, gin.H{
		"id":        defect.ID,
		"defect_id": defect.DefectID,
	})
}

// GetDefect 获取缺陷详情
// GET /api/v1/projects/:id/defects/:defectId
func (h *defectHandler) GetDefect(c *gin.Context) {
	defectID := c.Param("defectId")

	defect, err := h.defectService.GetByDefectID(defectID)
	if err != nil {
		if err.Error() == "defect not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		log.Printf("[Defect Get Failed] defect_id=%s, error=%v", defectID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, defect)
}

// UpdateDefect 更新缺陷
// PUT /api/v1/projects/:id/defects/:defectId
func (h *defectHandler) UpdateDefect(c *gin.Context) {
	defectID := c.Param("defectId")

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ResponseError(c, 401, "unauthorized")
		return
	}
	userID := userIDVal.(uint)

	// 先通过defectID获取UUID
	defect, err := h.defectService.GetByDefectID(defectID)
	if err != nil {
		if err.Error() == "defect not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		utils.ResponseError(c, 500, err.Error())
		return
	}

	var req models.DefectUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "validation failed: "+err.Error())
		return
	}

	if err := h.defectService.Update(defect.ID, userID, &req); err != nil {
		if err.Error() == "defect not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		log.Printf("[Defect Update Failed] defect_id=%s, user_id=%d, error=%v", defectID, userID, err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "defect updated successfully"})
}

// DeleteDefect 删除缺陷
// DELETE /api/v1/projects/:id/defects/:defectId
func (h *defectHandler) DeleteDefect(c *gin.Context) {
	defectID := c.Param("defectId")

	// 先通过defectID获取UUID
	defect, err := h.defectService.GetByDefectID(defectID)
	if err != nil {
		if err.Error() == "defect not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		utils.ResponseError(c, 500, err.Error())
		return
	}

	if err := h.defectService.Delete(defect.ID); err != nil {
		if err.Error() == "defect not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		log.Printf("[Defect Delete Failed] defect_id=%s, error=%v", defectID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "defect deleted successfully"})
}

// ExportTemplate 下载CSV导入模板
// GET /api/v1/projects/:id/defects/template
func (h *defectHandler) ExportTemplate(c *gin.Context) {
	data, err := h.defectService.GenerateTemplate()
	if err != nil {
		log.Printf("[Defect Template Failed] error=%v", err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=defect_template.csv")
	c.Data(200, "text/csv", data)
}

// ImportDefects 导入缺陷
// POST /api/v1/projects/:id/defects/import
func (h *defectHandler) ImportDefects(c *gin.Context) {
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

	src, err := file.Open()
	if err != nil {
		utils.ResponseError(c, 500, "failed to open file")
		return
	}
	defer src.Close()

	// 读取文件内容并检测BOM
	// 注意：这里直接传递src，在service层会自动处理UTF-8 BOM
	result, err := h.defectService.Import(uint(projectID), userID, src)
	if err != nil {
		log.Printf("[Defect Import Failed] project_id=%d, user_id=%d, error=%v", projectID, userID, err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, result)
}

// ExportDefects 导出缺陷
// GET /api/v1/projects/:id/defects/export
func (h *defectHandler) ExportDefects(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid project id")
		return
	}

	data, err := h.defectService.Export(uint(projectID))
	if err != nil {
		log.Printf("[Defect Export Failed] project_id=%d, error=%v", projectID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=defects_export.csv")
	c.Data(200, "text/csv", data)
}
