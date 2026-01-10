package handlers

import (
	"log"
	"net/http"
	"strconv"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// WebVersionHandler Web用例版本处理器
type WebVersionHandler struct {
	service services.WebVersionService
}

// NewWebVersionHandler 创建Web版本处理器实例
func NewWebVersionHandler(service services.WebVersionService) *WebVersionHandler {
	return &WebVersionHandler{service: service}
}

// SaveWebVersion 保存Web用例版本
// POST /api/v1/projects/:id/web-cases/versions
func (h *WebVersionHandler) SaveWebVersion(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 调用服务保存版本
	versionInfo, err := h.service.SaveVersion(uint(projectID), userID)
	if err != nil {
		log.Printf("[Web Version Save Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)

		// 根据错误类型返回不同状态码
		if err.Error() == "项目不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "没有可导出的Web用例" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "版本保存失败")
		return
	}

	log.Printf("[Web Version Save] user_id=%d, project_id=%d, version_id=%s", userID, projectID, versionInfo.VersionID)
	utils.SuccessResponse(c, versionInfo)
}

// GetWebVersionList 获取Web用例版本列表
// GET /api/v1/projects/:id/web-cases/versions
func (h *WebVersionHandler) GetWebVersionList(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	} else if size > 100000 {
		size = 100000
	}

	// 调用服务获取版本列表
	versionList, err := h.service.GetVersionList(uint(projectID), userID, page, size)
	if err != nil {
		log.Printf("[Web Version List Get Failed] user_id=%d, project_id=%d, error=%v", userID, projectID, err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取版本列表失败")
		return
	}

	log.Printf("[Web Version List Get] user_id=%d, project_id=%d, page=%d, size=%d, total=%d", userID, projectID, page, size, versionList.Total)
	utils.SuccessResponse(c, versionList)
}

// DownloadWebVersion 下载Web用例版本
// GET /api/v1/projects/:id/web-cases/versions/:versionId/export
func (h *WebVersionHandler) DownloadWebVersion(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取版本ID
	versionID := c.Param("versionId")
	if versionID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "版本ID不能为空")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 调用服务下载版本文件
	fileData, filename, err := h.service.DownloadVersion(uint(projectID), userID, versionID)
	if err != nil {
		log.Printf("[Web Version Download Failed] user_id=%d, project_id=%d, version_id=%s, error=%v", userID, projectID, versionID, err)

		if err.Error() == "版本不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "文件不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, "版本文件不存在")
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "下载版本失败")
		return
	}

	log.Printf("[Web Version Download] user_id=%d, project_id=%d, version_id=%s, filename=%s", userID, projectID, versionID, filename)

	// 设置响应头
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(len(fileData)))

	// 返回文件字节流
	c.Data(http.StatusOK, "application/zip", fileData)
}

// DeleteWebVersion 删除Web用例版本
// DELETE /api/v1/projects/:id/web-cases/versions/:versionId
func (h *WebVersionHandler) DeleteWebVersion(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取版本ID
	versionID := c.Param("versionId")
	if versionID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "版本ID不能为空")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 调用服务删除版本
	err = h.service.DeleteVersion(uint(projectID), userID, versionID)
	if err != nil {
		log.Printf("[Web Version Delete Failed] user_id=%d, project_id=%d, version_id=%s, error=%v", userID, projectID, versionID, err)

		if err.Error() == "版本不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "删除版本失败")
		return
	}

	log.Printf("[Web Version Delete] user_id=%d, project_id=%d, version_id=%s", userID, projectID, versionID)
	utils.SuccessResponse(c, gin.H{"message": "版本删除成功"})
}

// UpdateWebVersionRemark 更新Web用例版本备注
// PUT /api/v1/projects/:id/web-cases/versions/:versionId/remark
func (h *WebVersionHandler) UpdateWebVersionRemark(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取版本ID
	versionID := c.Param("versionId")
	if versionID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "版本ID不能为空")
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}
	userID := userIDVal.(uint)

	// 解析请求体
	var req struct {
		Remark string `json:"remark" binding:"max=200"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数验证失败")
		return
	}

	// 调用服务更新备注
	err = h.service.UpdateVersionRemark(uint(projectID), userID, versionID, req.Remark)
	if err != nil {
		log.Printf("[Web Version Remark Update Failed] user_id=%d, project_id=%d, version_id=%s, error=%v", userID, projectID, versionID, err)

		if err.Error() == "版本不存在" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "更新备注失败")
		return
	}

	log.Printf("[Web Version Remark Update] user_id=%d, project_id=%d, version_id=%s", userID, projectID, versionID)
	utils.SuccessResponse(c, gin.H{"message": "备注更新成功"})
}
