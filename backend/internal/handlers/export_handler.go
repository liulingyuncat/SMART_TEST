package handlers

import (
	"net/http"
	"strconv"
	"webtest/internal/constants"
	"webtest/internal/services"

	"github.com/gin-gonic/gin"
)

// ExportHandler 导出处理器
type ExportHandler struct {
	excelService services.ExcelService
}

// NewExportHandler 创建导出处理器实例
func NewExportHandler(excelService services.ExcelService) *ExportHandler {
	return &ExportHandler{
		excelService: excelService,
	}
}

// ExportAICases 导出AI用例(9列单Sheet)
// @Summary 导出AI用例
// @Tags Export
// @Param id path int true "项目ID"
// @Success 200 {file} xlsx "Excel文件流"
// @Router /api/manual-cases/:id/export/ai [get]
func (h *ExportHandler) ExportAICases(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	fileData, filename, err := h.excelService.ExportAICases(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.GetErrorMessage(constants.ErrExportFailed)})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileData)
}

// ExportTemplate 导出用例模板(23列空模板+示例行)
// @Summary 导出用例模板
// @Tags Export
// @Param id path int true "项目ID"
// @Param caseType query string true "用例类型(overall/change)"
// @Success 200 {file} xlsx "Excel模板文件"
// @Router /api/manual-cases/:id/export/template [get]
func (h *ExportHandler) ExportTemplate(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	caseType := c.Query("caseType")
	if caseType != "overall" && caseType != "change" && caseType != "acceptance" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	fileData, filename, err := h.excelService.ExportTemplate(uint(projectID), caseType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.GetErrorMessage(constants.ErrExportFailed)})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileData)
}

// ExportCases 导出整体/变更用例(双Sheet: 元数据+23列数据)
// @Summary 导出用例数据
// @Tags Export
// @Param id path int true "项目ID"
// @Param caseType query string true "用例类型(overall/change)"
// @Success 200 {file} xlsx "Excel数据文件"
// @Router /api/manual-cases/:id/export/cases [get]
func (h *ExportHandler) ExportCases(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	caseType := c.Query("caseType")
	if caseType != "overall" && caseType != "change" && caseType != "acceptance" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	// 获取可选的task_uuid参数(用于导出执行结果)
	taskUUID := c.Query("task_uuid")

	// T44: 新增language和case_group参数支持按语言导出
	language := c.Query("language")    // CN/JP/EN
	caseGroup := c.Query("case_group") // 用例集名称

	fileData, filename, err := h.excelService.ExportCases(uint(projectID), caseType, taskUUID, language, caseGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.GetErrorMessage(constants.ErrExportFailed)})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileData)
}
