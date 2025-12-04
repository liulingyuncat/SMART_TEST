package handlers

import (
	"log"
	"strconv"
	"webtest/internal/models"
	"webtest/internal/services"
	"webtest/internal/utils"

	"github.com/gin-gonic/gin"
)

// DefectConfigHandler 缺陷配置处理器接口
type DefectConfigHandler interface {
	// Subject管理
	GetSubjects(c *gin.Context)
	CreateSubject(c *gin.Context)
	UpdateSubject(c *gin.Context)
	DeleteSubject(c *gin.Context)
	// Phase管理
	GetPhases(c *gin.Context)
	CreatePhase(c *gin.Context)
	UpdatePhase(c *gin.Context)
	DeletePhase(c *gin.Context)
}

type defectConfigHandler struct {
	configService services.DefectConfigService
}

// NewDefectConfigHandler 创建缺陷配置处理器实例
func NewDefectConfigHandler(configService services.DefectConfigService) DefectConfigHandler {
	return &defectConfigHandler{
		configService: configService,
	}
}

// ========== Subject管理 ==========

// GetSubjects 获取Subject列表
// GET /api/v1/projects/:id/defect-subjects
func (h *defectConfigHandler) GetSubjects(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid project id")
		return
	}

	subjects, err := h.configService.ListSubjects(uint(projectID))
	if err != nil {
		log.Printf("[Subject List Failed] project_id=%d, error=%v", projectID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, subjects)
}

// CreateSubject 创建Subject
// POST /api/v1/projects/:id/defect-subjects
func (h *defectConfigHandler) CreateSubject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid project id")
		return
	}

	var req models.DefectSubjectCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "validation failed: "+err.Error())
		return
	}

	subject, err := h.configService.CreateSubject(uint(projectID), &req)
	if err != nil {
		if err.Error() == "subject name already exists" {
			utils.ResponseError(c, 409, err.Error())
			return
		}
		log.Printf("[Subject Create Failed] project_id=%d, error=%v", projectID, err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccessWithCode(c, 201, subject)
}

// UpdateSubject 更新Subject
// PUT /api/v1/projects/:id/defect-subjects/:subjectId
func (h *defectConfigHandler) UpdateSubject(c *gin.Context) {
	subjectIDStr := c.Param("subjectId")
	subjectID, err := strconv.ParseUint(subjectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid subject id")
		return
	}

	var req models.DefectSubjectUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "validation failed: "+err.Error())
		return
	}

	if err := h.configService.UpdateSubject(uint(subjectID), &req); err != nil {
		if err.Error() == "subject not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		if err.Error() == "subject name already exists" {
			utils.ResponseError(c, 409, err.Error())
			return
		}
		log.Printf("[Subject Update Failed] subject_id=%d, error=%v", subjectID, err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "subject updated successfully"})
}

// DeleteSubject 删除Subject
// DELETE /api/v1/projects/:id/defect-subjects/:subjectId
func (h *defectConfigHandler) DeleteSubject(c *gin.Context) {
	subjectIDStr := c.Param("subjectId")
	subjectID, err := strconv.ParseUint(subjectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid subject id")
		return
	}

	if err := h.configService.DeleteSubject(uint(subjectID)); err != nil {
		if err.Error() == "subject not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		log.Printf("[Subject Delete Failed] subject_id=%d, error=%v", subjectID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "subject deleted successfully"})
}

// ========== Phase管理 ==========

// GetPhases 获取Phase列表
// GET /api/v1/projects/:id/defect-phases
func (h *defectConfigHandler) GetPhases(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid project id")
		return
	}

	phases, err := h.configService.ListPhases(uint(projectID))
	if err != nil {
		log.Printf("[Phase List Failed] project_id=%d, error=%v", projectID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, phases)
}

// CreatePhase 创建Phase
// POST /api/v1/projects/:id/defect-phases
func (h *defectConfigHandler) CreatePhase(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid project id")
		return
	}

	var req models.DefectPhaseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "validation failed: "+err.Error())
		return
	}

	phase, err := h.configService.CreatePhase(uint(projectID), &req)
	if err != nil {
		if err.Error() == "phase name already exists" {
			utils.ResponseError(c, 409, err.Error())
			return
		}
		log.Printf("[Phase Create Failed] project_id=%d, error=%v", projectID, err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccessWithCode(c, 201, phase)
}

// UpdatePhase 更新Phase
// PUT /api/v1/projects/:id/defect-phases/:phaseId
func (h *defectConfigHandler) UpdatePhase(c *gin.Context) {
	phaseIDStr := c.Param("phaseId")
	phaseID, err := strconv.ParseUint(phaseIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid phase id")
		return
	}

	var req models.DefectPhaseUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "validation failed: "+err.Error())
		return
	}

	if err := h.configService.UpdatePhase(uint(phaseID), &req); err != nil {
		if err.Error() == "phase not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		if err.Error() == "phase name already exists" {
			utils.ResponseError(c, 409, err.Error())
			return
		}
		log.Printf("[Phase Update Failed] phase_id=%d, error=%v", phaseID, err)
		utils.ResponseError(c, 400, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "phase updated successfully"})
}

// DeletePhase 删除Phase
// DELETE /api/v1/projects/:id/defect-phases/:phaseId
func (h *defectConfigHandler) DeletePhase(c *gin.Context) {
	phaseIDStr := c.Param("phaseId")
	phaseID, err := strconv.ParseUint(phaseIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, 400, "invalid phase id")
		return
	}

	if err := h.configService.DeletePhase(uint(phaseID)); err != nil {
		if err.Error() == "phase not found" {
			utils.ResponseError(c, 404, err.Error())
			return
		}
		log.Printf("[Phase Delete Failed] phase_id=%d, error=%v", phaseID, err)
		utils.ResponseError(c, 500, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "phase deleted successfully"})
}
