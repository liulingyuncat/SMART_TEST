package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"webtest/internal/models"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// ApiTestCaseHandler 接口测试用例处理器
type ApiTestCaseHandler struct {
	service services.ApiTestCaseService
}

// NewApiTestCaseHandler 创建处理器实例
func NewApiTestCaseHandler(service services.ApiTestCaseService) *ApiTestCaseHandler {
	return &ApiTestCaseHandler{service: service}
}

// GetCases 获取用例列表
// GET /api/v1/projects/:id/api-cases?case_type=role1&case_group=xxx&page=1&size=50
func (h *ApiTestCaseHandler) GetCases(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	caseType := c.DefaultQuery("case_type", "api")
	caseGroup := c.Query("case_group") // 用例集筛选参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))

	cases, total, err := h.service.GetCasesByGroup(uint(projectID), userID, caseType, caseGroup, page, size)
	if err != nil {
		log.Printf("[API Cases Get Failed] user=%d, project=%d, case_group=%s, error=%v", userID, projectID, caseGroup, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取用例列表失败")
		return
	}

	log.Printf("[API Cases Get] user=%d, project=%d, type=%s, case_group=%s, total=%d", userID, projectID, caseType, caseGroup, total)
	utils.SuccessResponse(c, gin.H{
		"cases": cases,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// CreateCase 创建用例
// POST /api/v1/projects/:id/api-cases
func (h *ApiTestCaseHandler) CreateCase(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	var req struct {
		CaseType   string `json:"case_type"`
		CaseGroup  string `json:"case_group"`
		GroupID    int    `json:"group_id"` // 支持通过group_id传入
		CaseNumber string `json:"case_number"`
		Screen     string `json:"screen"`
		URL        string `json:"url"`
		Method     string `json:"method"`
		Header     string `json:"header"`
		Body       string `json:"body"`
		Response   string `json:"response"`
		TestResult string `json:"test_result"`
		Remark     string `json:"remark"`
		ScriptCode string `json:"script_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 创建用例数据
	newCase := &models.ApiTestCase{
		ProjectID:    uint(projectID),
		CaseType:     req.CaseType,
		CaseGroup:    req.CaseGroup,
		CaseNumber:   req.CaseNumber,
		Screen:       req.Screen,
		URL:          req.URL,
		Method:       req.Method,
		Header:       req.Header,
		Body:         req.Body,
		Response:     req.Response,
		TestResult:   req.TestResult,
		Remark:       req.Remark,
		ScriptCode:   req.ScriptCode,
		DisplayOrder: 1, // 默认为1,后续会重新分配
	}

	// 调用Service创建
	createdCase, err := h.service.CreateCase(uint(projectID), userID, newCase)
	if err != nil {
		log.Printf("[API Case Create Failed] user=%d, project=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建用例失败")
		return
	}

	log.Printf("[API Case Create] user=%d, project=%d, case_id=%s", userID, projectID, createdCase.ID)
	utils.SuccessResponse(c, createdCase)
}

// InsertCase 插入用例
// POST /api/v1/projects/:id/api-cases/insert
func (h *ApiTestCaseHandler) InsertCase(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	var req struct {
		CaseType     string                 `json:"case_type"`
		Position     string                 `json:"position"`
		TargetCaseID string                 `json:"target_case_id"`
		CaseGroup    string                 `json:"case_group"`
		CaseData     map[string]interface{} `json:"case_data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 如果请求中包含case_group，添加到case_data中
	if req.CaseGroup != "" {
		if req.CaseData == nil {
			req.CaseData = make(map[string]interface{})
		}
		req.CaseData["case_group"] = req.CaseGroup
	}

	newCase, err := h.service.InsertCase(uint(projectID), userID, req.CaseType, req.Position, req.TargetCaseID, req.CaseData)
	if err != nil {
		log.Printf("[API Case Insert Failed] user=%d, project=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[API Case Insert] user=%d, project=%d, case_id=%s", userID, projectID, newCase.ID)
	utils.SuccessResponse(c, newCase)
}

// DeleteCase 删除单个用例
// DELETE /api/v1/projects/:id/api-cases/:caseId
func (h *ApiTestCaseHandler) DeleteCase(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	caseID := c.Param("caseId")
	if caseID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "用例ID不能为空")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 调用service层删除
	err = h.service.DeleteCase(uint(projectID), userID, caseID)
	if err != nil {
		log.Printf("[API Case Delete Failed] user=%d, project=%d, case_id=%s, error=%v", userID, projectID, caseID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "用例不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除失败")
		return
	}

	log.Printf("[API Case Delete] user=%d, project=%d, case_id=%s", userID, projectID, caseID)
	utils.SuccessResponse(c, gin.H{"message": "删除成功"})
}

// BatchDeleteCases 批量删除用例
// POST /api/v1/projects/:id/api-cases/batch-delete
func (h *ApiTestCaseHandler) BatchDeleteCases(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	var req struct {
		CaseType string   `json:"case_type"`
		CaseIDs  []string `json:"case_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	deletedCount, failedIDs, err := h.service.BatchDeleteCases(uint(projectID), userID, req.CaseType, req.CaseIDs)
	if err != nil {
		log.Printf("[API Case BatchDelete Failed] user=%d, project=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量删除失败")
		return
	}

	log.Printf("[API Case BatchDelete] user=%d, project=%d, deleted=%d, failed=%d", userID, projectID, deletedCount, len(failedIDs))
	utils.SuccessResponse(c, gin.H{
		"deleted_count": deletedCount,
		"failed_ids":    failedIDs,
	})
}

// UpdateCase 更新用例
// PATCH /api/v1/projects/:id/api-cases/:caseId
func (h *ApiTestCaseHandler) UpdateCase(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	caseID := c.Param("caseId")
	if caseID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用例ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	err = h.service.UpdateCase(uint(projectID), userID, caseID, updates)
	if err != nil {
		log.Printf("[API Case Update Failed] user=%d, project=%d, case_id=%s, error=%v", userID, projectID, caseID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新用例失败")
		return
	}

	log.Printf("[API Case Update] user=%d, project=%d, case_id=%s", userID, projectID, caseID)
	utils.MessageResponse(c, http.StatusOK, "用例更新成功")
}

// ========== 版本管理 ==========

// SaveVersion 保存版本
// POST /api/v1/projects/:id/api-cases/versions
func (h *ApiTestCaseHandler) SaveVersion(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req) // 可选参数

	version, err := h.service.SaveVersion(uint(projectID), userID, req.Remark)
	if err != nil {
		log.Printf("[API Version Save Failed] user=%d, project=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[API Version Save] user=%d, project=%d, version_id=%s", userID, projectID, version.ID)
	utils.SuccessResponse(c, gin.H{
		"version_id":    version.ID,
		"xlsx_filename": version.XlsxFilename,
		"created_at":    version.CreatedAt,
	})
}

// GetVersions 获取版本列表
// GET /api/v1/projects/:id/api-cases/versions?page=1&size=10
func (h *ApiTestCaseHandler) GetVersions(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	versions, total, err := h.service.GetVersions(uint(projectID), userID, page, size)
	if err != nil {
		log.Printf("[API Versions Get Failed] user=%d, project=%d, error=%v", userID, projectID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取版本列表失败")
		return
	}

	log.Printf("[API Versions Get] user=%d, project=%d, total=%d", userID, projectID, total)
	utils.SuccessResponse(c, gin.H{
		"versions": versions,
		"total":    total,
		"page":     page,
		"size":     size,
	})
}

// DownloadVersion 下载版本ZIP
// GET /api/v1/projects/:id/api-cases/versions/:versionId/export
func (h *ApiTestCaseHandler) DownloadVersion(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	versionID := c.Param("versionId")
	if versionID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的版本ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 查询版本记录
	versions, _, err := h.service.GetVersions(uint(projectID), userID, 1, 1000)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "查询版本失败")
		return
	}

	var targetVersion *models.ApiTestCaseVersion
	for _, v := range versions {
		if v.ID == versionID {
			targetVersion = v
			break
		}
	}

	if targetVersion == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "版本不存在")
		return
	}

	storageDir := filepath.Join("storage", "versions", "api-cases")

	// 新版本：直接下载XLSX文件
	if targetVersion.XlsxFilename != "" {
		filePath := filepath.Join(storageDir, targetVersion.XlsxFilename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("[API Version Download] 读取XLSX文件失败: %s, error=%v", targetVersion.XlsxFilename, err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "读取文件失败")
			return
		}

		// 设置响应头
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", targetVersion.XlsxFilename))
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", content)
	} else {
		// 兼容旧版本：生成ZIP文件包含CSV文件
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		filenames := []string{
			targetVersion.FilenameRole1,
			targetVersion.FilenameRole2,
			targetVersion.FilenameRole3,
			targetVersion.FilenameRole4,
		}

		for _, filename := range filenames {
			if filename == "" {
				continue
			}
			filePath := filepath.Join(storageDir, filename)
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("[API Version Download] 读取文件失败: %s, error=%v", filename, err)
				continue
			}

			fw, err := zipWriter.Create(filename)
			if err != nil {
				continue
			}
			fw.Write(content)
		}

		zipWriter.Close()

		// 设置响应头
		zipFilename := fmt.Sprintf("%s.zip", versionID)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFilename))
		c.Header("Content-Type", "application/zip")
		c.Data(http.StatusOK, "application/zip", buf.Bytes())
	}

	log.Printf("[API Version Download] user=%d, project=%d, version_id=%s", userID, projectID, versionID)
}

// DeleteVersion 删除版本
// DELETE /api/v1/projects/:id/api-cases/versions/:versionId
func (h *ApiTestCaseHandler) DeleteVersion(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	versionID := c.Param("versionId")
	if versionID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的版本ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	err = h.service.DeleteVersion(uint(projectID), userID, versionID)
	if err != nil {
		log.Printf("[API Version Delete Failed] user=%d, project=%d, version_id=%s, error=%v", userID, projectID, versionID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "版本不存在" || err.Error() == "版本不属于该项目" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除版本失败")
		return
	}

	log.Printf("[API Version Delete] user=%d, project=%d, version_id=%s", userID, projectID, versionID)
	utils.MessageResponse(c, http.StatusOK, "版本删除成功")
}

// UpdateVersionRemark 更新版本备注
// PUT /api/v1/projects/:id/api-cases/versions/:versionId/remark
func (h *ApiTestCaseHandler) UpdateVersionRemark(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	versionID := c.Param("versionId")
	if versionID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的版本ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	var req struct {
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	err = h.service.UpdateVersionRemark(uint(projectID), userID, versionID, req.Remark)
	if err != nil {
		log.Printf("[API Version Remark Update Failed] user=%d, project=%d, version_id=%s, error=%v", userID, projectID, versionID, err)
		if err.Error() == "无项目访问权限" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "版本不存在" || err.Error() == "版本不属于该项目" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新备注失败")
		return
	}

	log.Printf("[API Version Remark Update] user=%d, project=%d, version_id=%s", userID, projectID, versionID)
	utils.MessageResponse(c, http.StatusOK, "备注更新成功")
}

// ExportTemplate 导出API用例模版
// GET /api/v1/projects/:id/api-cases/template
func (h *ApiTestCaseHandler) ExportTemplate(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 导出空模版（无数据）
	data, filename, err := h.service.ExportApiTemplate(uint(projectID))
	if err != nil {
		log.Printf("[API Template Export Failed] project=%d, error=%v", projectID, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("导出模版失败: %v", err))
		return
	}

	log.Printf("[API Template Export] user=%d, project=%d, filename=%s", userIDVal, projectID, filename)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// ImportCases 导入API用例
// POST /api/v1/projects/:id/api-cases/import
func (h *ApiTestCaseHandler) ImportCases(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取case_group参数
	caseGroup := c.PostForm("case_group")
	if caseGroup == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "请选择用例集")
		return
	}

	// 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请选择要导入的文件")
		return
	}

	// 读取文件内容
	fileReader, err := file.Open()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "读取文件失败")
		return
	}
	defer fileReader.Close()

	fileData := make([]byte, file.Size)
	_, err = fileReader.Read(fileData)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "读取文件内容失败")
		return
	}

	// 调用service导入
	insertCount, updateCount, err := h.service.ImportApiCases(uint(projectID), userID, caseGroup, fileData)
	if err != nil {
		log.Printf("[API Cases Import Failed] user=%d, project=%d, case_group=%s, error=%v", userID, projectID, caseGroup, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("导入失败: %v", err))
		return
	}

	log.Printf("[API Cases Import] user=%d, project=%d, case_group=%s, insert=%d, update=%d", userID, projectID, caseGroup, insertCount, updateCount)
	utils.SuccessResponse(c, gin.H{
		"message":      "导入成功",
		"insert_count": insertCount,
		"update_count": updateCount,
	})
}
