package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"webtest/internal/constants"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// ImportHandler 导入处理器
type ImportHandler struct {
	excelService services.ExcelService
}

// NewImportHandler 创建导入处理器实例
func NewImportHandler(excelService services.ExcelService) *ImportHandler {
	return &ImportHandler{
		excelService: excelService,
	}
}

// ImportCases 导入用例(支持UUID匹配更新)
// @Summary 导入用例
// @Tags Import
// @Accept multipart/form-data
// @Param id path int true "项目ID"
// @Param caseType formData string true "用例类型(overall/change)"
// @Param file formData file true "Excel文件"
// @Success 200 {object} map[string]interface{} "导入结果"
// @Router /api/manual-cases/:id/import [post]
func (h *ImportHandler) ImportCases(c *gin.Context) {
	fmt.Println("========================================")
	fmt.Println("ImportCases handler called!")
	fmt.Println("========================================")

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	caseType := c.PostForm("caseType")
	fmt.Printf("ProjectID: %d, CaseType: %s\n", projectID, caseType)

	if caseType != "overall" && caseType != "change" && caseType != "acceptance" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("Error getting file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrFileRequired)})
		return
	}
	fmt.Printf("File received: %s, Size: %d bytes\n", file.Filename, file.Size)

	// 文件验证
	if err := utils.ValidateFileExtension(file.Filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := utils.ValidateFileSize(file.Size); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 读取文件数据
	fileHandle, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.GetErrorMessage(constants.ErrImportFailed)})
		return
	}
	defer fileHandle.Close()

	fileData, err := io.ReadAll(fileHandle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.GetErrorMessage(constants.ErrImportFailed)})
		return
	}

	// 执行导入
	fmt.Println("Calling excelService.ImportCases...")
	updateCount, insertCount, err := h.excelService.ImportCases(uint(projectID), caseType, fileData)
	if err != nil {
		fmt.Printf("Import failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.GetErrorMessage(constants.ErrImportFailed), "details": err.Error()})
		return
	}

	fmt.Printf("Import successful: %d updated, %d inserted\n", updateCount, insertCount)
	c.JSON(http.StatusOK, gin.H{
		"message":     "导入成功",
		"updateCount": updateCount,
		"insertCount": insertCount,
	})
}
