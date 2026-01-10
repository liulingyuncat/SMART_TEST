package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"webtest/internal/constants"
	"webtest/internal/models"
	"webtest/internal/repositories"
	"webtest/internal/services"

	"github.com/gin-gonic/gin"
)

// VersionHandler 版本管理处理器
type VersionHandler struct {
	versionService services.VersionService
	projectService services.ProjectService
	versionRepo    repositories.VersionRepository // 直接查询versions表
}

// NewVersionHandler 创建版本管理处理器实例
func NewVersionHandler(versionService services.VersionService, projectService services.ProjectService, versionRepo repositories.VersionRepository) *VersionHandler {
	return &VersionHandler{
		versionService: versionService,
		projectService: projectService,
		versionRepo:    versionRepo,
	}
}

// SaveVersion 保存版本(导出并存储)
// @Summary 保存版本
// @Tags Version
// @Param id path int true "项目ID"
// @Param case_type query string false "用例类型(overall/change/acceptance)" default(overall)
// @Param caseType formData string false "用例类型(overall/change/acceptance,兼容旧接口)"
// @Success 200 {object} map[string]interface{} "保存结果"
// @Router /api/manual-cases/:id/versions/save [post]
func (h *VersionHandler) SaveVersion(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput), "details": "invalid project id"})
		return
	}

	// 优先从查询参数读取case_type,其次从FormData读取caseType,默认overall
	caseType := c.DefaultQuery("case_type", "")
	if caseType == "" {
		caseType = c.PostForm("caseType")
	}
	if caseType == "" {
		caseType = "overall" // 默认值,保持向后兼容
	}

	if caseType != "overall" && caseType != "change" && caseType != "acceptance" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput), "details": "caseType must be overall, change or acceptance"})
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	createdBy := userIDVal.(uint)

	filename, err := h.versionService.SaveVersion(uint(projectID), createdBy, caseType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.GetErrorMessage(constants.ErrExportFailed), "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "版本保存成功", "filename": filename})
}

// GetVersionList 获取版本列表
// @Summary 获取版本列表
// @Tags Version
// @Param id path int true "项目ID"
// @Param case_type query string false "用例类型(overall/change/acceptance),为空返回所有"
// @Success 200 {array} models.CaseVersion "版本列表"
// @Router /api/manual-cases/:id/versions [get]
func (h *VersionHandler) GetVersionList(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	// 从查询参数读取case_type(可选)
	caseType := c.Query("case_type")

	versions, err := h.versionService.GetVersionList(uint(projectID), caseType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取版本列表失败", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// DownloadVersion 下载指定版本文件
// @Summary 下载版本文件
// @Tags Version
// @Param id path int true "项目ID"
// @Param versionID path int true "版本ID"
// @Success 200 {file} xlsx "Excel文件流"
// @Router /api/manual-cases/:id/versions/:versionID/download [get]
func (h *VersionHandler) DownloadVersion(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	versionID, err := strconv.ParseUint(c.Param("versionID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	fileData, filename, err := h.versionService.DownloadVersion(uint(projectID), uint(versionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.GetErrorMessage(constants.ErrFileNotFound), "details": err.Error()})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileData)
}

// DeleteVersion 删除指定版本(文件+数据库记录)
// @Summary 删除版本
// @Tags Version
// @Param id path int true "项目ID"
// @Param versionID path int true "版本ID"
// @Success 200 {object} map[string]string "删除结果"
// @Router /api/manual-cases/:id/versions/:versionID [delete]
func (h *VersionHandler) DeleteVersion(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	versionID, err := strconv.ParseUint(c.Param("versionID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.GetErrorMessage(constants.ErrInvalidInput)})
		return
	}

	err = h.versionService.DeleteVersion(uint(projectID), uint(versionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除版本失败", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "版本删除成功"})
}

// SaveVersionGeneric 通用版本保存接口(支持需求管理类型)
// @Summary 保存版本(通用)
// @Tags Version
// @Accept json
// @Param request body object true "版本保存请求"
// @Success 200 {object} map[string]interface{} "保存结果"
// @Router /api/versions [post]
func (h *VersionHandler) SaveVersionGeneric(c *gin.Context) {
	var req struct {
		ProjectID string `json:"project_id" binding:"required"`
		DocType   string `json:"doc_type" binding:"required"`
		Content   string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "details": err.Error()})
		return
	}

	// 转换projectID
	projectID, err := strconv.ParseUint(req.ProjectID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "项目ID无效"})
		return
	}

	// 验证docType
	validDocTypes := map[string]string{
		"overall-requirements":   "overall_requirements",
		"overall-test-viewpoint": "overall_viewpoint",
		"change-requirements":    "change_requirements",
		"change-test-viewpoint":  "change_viewpoint",
	}

	englishName, ok := validDocTypes[req.DocType]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的文档类型"})
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}
	userID := userIDVal.(uint)

	// 获取项目信息
	project, _, err := h.projectService.GetByID(uint(projectID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取项目信息失败"})
		return
	}

	// 自动生成文件名(项目名+英文文档类型+时间戳)
	now := time.Now()
	timestamp := now.Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s_%s.md", project.Name, englishName, timestamp)

	// 3. 保存.md文件
	storageDir := filepath.Join("storage", "versions", fmt.Sprintf("%d", projectID))
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建目录失败"})
		return
	}

	filePath := filepath.Join(storageDir, filename)
	if err := os.WriteFile(filePath, []byte(req.Content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "写入文件失败"})
		return
	}

	// 4. 创建CaseVersion记录
	version := &models.CaseVersion{
		ProjectID: uint(projectID),
		DocType:   req.DocType,
		Filename:  filename,
		FilePath:  filePath,
		FileSize:  int64(len(req.Content)),
		CreatedBy: &userID,
	}

	// 创建版本记录(通过service)
	if err := h.versionService.CreateVersion(version); err != nil {
		os.Remove(filePath) // 删除已创建的文件
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建版本记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "版本保存成功", "data": gin.H{"filename": filename}})
}

// GetVersionListGeneric 通用版本列表接口
// @Summary 获取版本列表(通用)
// @Tags Version
// @Param project_id query string true "项目ID"
// @Param doc_type query string false "文档类型(test_case_versions表)"
// @Param item_type query string false "条目类型(versions表)"
// @Success 200 {array} models.Version "版本列表"
// @Router /api/versions [get]
func (h *VersionHandler) GetVersionListGeneric(c *gin.Context) {
	projectIDStr := c.Query("project_id")
	if projectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "项目ID不能为空"})
		return
	}

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "项目ID无效"})
		return
	}

	docType := c.Query("doc_type")
	itemType := c.Query("item_type")

	log.Printf("[GetVersionListGeneric] projectID=%d, docType=%s, itemType=%s", projectID, docType, itemType)

	// 如果doc_type为空，说明是查询versions表
	if docType == "" {
		// 查询新的versions表
		allVersions, err := h.versionRepo.FindByProjectID(uint(projectID))
		if err != nil {
			log.Printf("[GetVersionListGeneric] Error querying versions table: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取版本列表失败"})
			return
		}

		// 如果指定了item_type，进行过滤
		if itemType != "" {
			filtered := make([]*models.Version, 0)
			for _, v := range allVersions {
				if v.ItemType == itemType {
					filtered = append(filtered, v)
				}
			}
			log.Printf("[GetVersionListGeneric] Found %d versions (filtered by item_type=%s)", len(filtered), itemType)
			c.JSON(http.StatusOK, filtered)
			return
		}

		log.Printf("[GetVersionListGeneric] Found %d versions from versions table", len(allVersions))
		c.JSON(http.StatusOK, allVersions)
		return
	}

	// doc_type不为空，查询旧的test_case_versions表
	versions, err := h.versionService.GetVersionList(uint(projectID), docType)
	if err != nil {
		log.Printf("[GetVersionListGeneric] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取版本列表失败"})
		return
	}

	log.Printf("[GetVersionListGeneric] Found %d versions from test_case_versions table", len(versions))
	c.JSON(http.StatusOK, versions)
}

// DownloadVersionGeneric 通用版本下载接口
// @Summary 下载版本文件(通用)
// @Tags Version
// @Param id path int true "版本ID"
// @Success 200 {file} file "文件流"
// @Router /api/versions/:id/download [get]
func (h *VersionHandler) DownloadVersionGeneric(c *gin.Context) {
	versionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "版本ID无效"})
		return
	}

	log.Printf("[DownloadVersionGeneric] versionID=%d", versionID)

	// 先尝试从versions表查询（新表）
	version, err := h.versionRepo.FindByID(uint(versionID))
	if err == nil && version != nil {
		// 从versions表找到了，读取ZIP文件
		log.Printf("[DownloadVersionGeneric] Found in versions table: %s", version.FilePath)
		fileData, err := os.ReadFile(version.FilePath)
		if err != nil {
			log.Printf("[DownloadVersionGeneric] Read file error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "文件读取失败"})
			return
		}

		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", "attachment; filename="+version.Filename)
		c.Data(http.StatusOK, "application/zip", fileData)
		return
	}

	// 如果versions表没找到，尝试test_case_versions表（旧表）
	log.Printf("[DownloadVersionGeneric] Not found in versions table, trying test_case_versions")
	caseVersion, err := h.versionService.GetVersionByID(uint(versionID))
	if err != nil {
		log.Printf("[DownloadVersionGeneric] Version not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	fileData, filename, err := h.versionService.DownloadVersion(caseVersion.ProjectID, uint(versionID))
	if err != nil {
		log.Printf("[DownloadVersionGeneric] Download error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "下载失败"})
		return
	}

	c.Header("Content-Type", "text/markdown;charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/markdown", fileData)
}

// DeleteVersionGeneric 通用版本删除接口
// @Summary 删除版本(通用)
// @Tags Version
// @Param id path int true "版本ID"
// @Success 200 {object} map[string]string "删除结果"
// @Router /api/versions/:id [delete]
func (h *VersionHandler) DeleteVersionGeneric(c *gin.Context) {
	versionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "版本ID无效"})
		return
	}

	log.Printf("[DeleteVersionGeneric] versionID=%d", versionID)

	// 先尝试从versions表查询并删除（新表）
	version, err := h.versionRepo.FindByID(uint(versionID))
	if err == nil && version != nil {
		// 从versions表找到了，删除文件和记录
		log.Printf("[DeleteVersionGeneric] Found in versions table, deleting: %s", version.FilePath)

		// 删除文件
		if err := os.Remove(version.FilePath); err != nil {
			log.Printf("[DeleteVersionGeneric] Delete file error: %v", err)
		}

		// 删除数据库记录
		if err := h.versionRepo.Delete(uint(versionID)); err != nil {
			log.Printf("[DeleteVersionGeneric] Delete DB record error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
		return
	}

	// 如果versions表没找到，尝试test_case_versions表（旧表）
	log.Printf("[DeleteVersionGeneric] Not found in versions table, trying test_case_versions")
	caseVersion, err := h.versionService.GetVersionByID(uint(versionID))
	if err != nil {
		log.Printf("[DeleteVersionGeneric] Version not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	err = h.versionService.DeleteVersion(caseVersion.ProjectID, uint(versionID))
	if err != nil {
		log.Printf("[DeleteVersionGeneric] Delete error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// UpdateVersionRemarkGeneric 更新版本备注(通用接口)
// @Summary 更新版本备注
// @Tags Version
// @Param id path int true "版本ID"
// @Param remark body string true "备注内容"
// @Success 200 {object} map[string]interface{} "更新结果"
// @Router /api/v1/versions/:id/remark [put]
func (h *VersionHandler) UpdateVersionRemarkGeneric(c *gin.Context) {
	versionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "版本ID无效"})
		return
	}

	// 获取备注内容
	var req struct {
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	log.Printf("[UpdateVersionRemarkGeneric] versionID=%d, remark=%s", versionID, req.Remark)

	// 先尝试从versions表查询并更新（新表）
	version, err := h.versionRepo.FindByID(uint(versionID))
	if err == nil && version != nil {
		// 从versions表找到了，更新备注
		log.Printf("[UpdateVersionRemarkGeneric] Found in versions table, updating remark")

		if err := h.versionRepo.UpdateRemark(uint(versionID), req.Remark); err != nil {
			log.Printf("[UpdateVersionRemarkGeneric] Update error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
		return
	}

	// 如果versions表没找到，尝试test_case_versions表（旧表）
	log.Printf("[UpdateVersionRemarkGeneric] Not found in versions table, trying test_case_versions")
	caseVersion, err := h.versionService.GetVersionByID(uint(versionID))
	if err != nil {
		log.Printf("[UpdateVersionRemarkGeneric] Version not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	err = h.versionService.UpdateVersionRemark(caseVersion.ProjectID, uint(versionID), req.Remark)
	if err != nil {
		log.Printf("[UpdateVersionRemarkGeneric] Update error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// ExportTemplate 导出手工测试用例模版（CN/JP/EN空白xlsx打包成zip）
// @Summary 导出手工测试用例模版
// @Tags Version
// @Success 200 {file} application/zip "模版zip文件"
// @Router /api/v1/manual-cases/template [get]
func (h *VersionHandler) ExportTemplate(c *gin.Context) {
	// TODO: GenerateTemplate方法未实现
	// 调用服务层生成模版
	// zipBytes, filename, err := h.versionService.GenerateTemplate()
	err := fmt.Errorf("GenerateTemplate not implemented")
	if err != nil {
		log.Printf("[ExportTemplate] Failed to generate template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate template", "details": err.Error()})
		return
	}

	// 设置响应头并返回zip文件
	// c.Header("Content-Type", "application/zip")
	// c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	// c.Data(http.StatusOK, "application/zip", zipBytes)
}
