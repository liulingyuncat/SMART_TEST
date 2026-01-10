package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"
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
	projectRepo   repositories.ProjectRepository
}

// NewDefectHandler 创建缺陷处理器实例
func NewDefectHandler(defectService services.DefectService, projectRepo repositories.ProjectRepository) DefectHandler {
	return &defectHandler{
		defectService: defectService,
		projectRepo:   projectRepo,
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

// ExportTemplate 下载导入模板（支持CSV和XLSX格式）
// GET /api/v1/projects/:id/defects/template?format=csv|xlsx
func (h *defectHandler) ExportTemplate(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	if format != "csv" && format != "xlsx" {
		format = "csv"
	}

	data, err := h.defectService.GenerateTemplate(format)
	if err != nil {
		log.Printf("[Defect Template Failed] format=%s, error=%v", format, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	if format == "xlsx" {
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=defect_template.xlsx")
		c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
	} else {
		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=defect_template.csv")
		c.Data(200, "text/csv", data)
	}
}

// ImportDefects 导入缺陷（支持CSV和XLSX格式）
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

	// 根据文件扩展名检测格式
	filename := strings.ToLower(file.Filename)
	isXLSX := strings.HasSuffix(filename, ".xlsx")

	// 导入缺陷
	result, err := h.defectService.ImportWithFormat(uint(projectID), userID, src, isXLSX)
	if err != nil {
		log.Printf("[Defect Import Failed] project_id=%d, user_id=%d, error=%v", projectID, userID, err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, result)
}

// ExportDefects 导出缺陷（支持CSV和XLSX格式）
// GET /api/v1/projects/:id/defects/export?format=csv|xlsx
func (h *defectHandler) ExportDefects(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid project id")
		return
	}

	// 获取项目信息
	project, err := h.projectRepo.GetByID(uint(projectID))
	if err != nil {
		log.Printf("[Defect Export] Failed to get project: project_id=%d, error=%v", projectID, err)
		utils.ResponseError(c, 500, "project not found")
		return
	}

	// 获取导出格式（默认csv）
	format := c.DefaultQuery("format", "csv")
	if format != "csv" && format != "xlsx" {
		format = "csv"
	}

	data, err := h.defectService.ExportWithFormat(uint(projectID), format)
	if err != nil {
		log.Printf("[Defect Export Failed] project_id=%d, format=%s, error=%v", projectID, format, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	// 生成带项目名和时间戳的文件名
	timestamp := time.Now().Unix()
	var filename string
	if format == "xlsx" {
		filename = fmt.Sprintf("%s_defects_export_%d.xlsx", project.Name, timestamp)
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
	} else {
		filename = fmt.Sprintf("%s_defects_export_%d.csv", project.Name, timestamp)
		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(200, "text/csv", data)
	}
}
