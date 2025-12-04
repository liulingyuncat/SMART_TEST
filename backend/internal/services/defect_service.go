package services

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"gorm.io/gorm"
)

// DefectService 缺陷服务接口
type DefectService interface {
	// CRUD
	Create(projectID uint, userID uint, req *models.DefectCreateRequest) (*models.Defect, error)
	GetByID(id string) (*models.Defect, error)
	GetByDefectID(defectID string) (*models.Defect, error)
	Update(id string, userID uint, req *models.DefectUpdateRequest) error
	Delete(id string) error

	// 列表
	List(projectID uint, status, keyword string, page, size int) (*models.DefectListResponse, error)

	// 导入导出
	GenerateTemplate() ([]byte, error)
	Import(projectID uint, userID uint, reader io.Reader) (*models.ImportResult, error)
	Export(projectID uint) ([]byte, error)
}

type defectService struct {
	repo repositories.DefectRepository
}

// NewDefectService 创建缺陷服务实例
func NewDefectService(repo repositories.DefectRepository) DefectService {
	return &defectService{repo: repo}
}

// generateDefectID 生成缺陷显示ID
func (s *defectService) generateDefectID(projectID uint) (string, error) {
	maxSeq, err := s.repo.GetMaxDefectSeq(projectID)
	if err != nil {
		return "", fmt.Errorf("get max defect seq: %w", err)
	}
	return fmt.Sprintf("DEF-%06d", maxSeq+1), nil
}

// Create 创建缺陷
func (s *defectService) Create(projectID uint, userID uint, req *models.DefectCreateRequest) (*models.Defect, error) {
	// 生成DefectID
	defectID, err := s.generateDefectID(projectID)
	if err != nil {
		return nil, err
	}

	// 验证优先级
	if req.Priority != "" && !models.IsValidDefectPriority(req.Priority) {
		return nil, errors.New("invalid priority value")
	}

	// 验证严重程度
	if req.Severity != "" && !models.IsValidDefectSeverity(req.Severity) {
		return nil, errors.New("invalid severity value")
	}

	// 处理Subject：如果提供了SubjectID，查找名称
	subject := req.Subject
	if req.SubjectID != nil && *req.SubjectID > 0 {
		var subjectModel models.DefectSubject
		if err := s.repo.GetDB().Where("project_id = ?", projectID).First(&subjectModel, *req.SubjectID).Error; err == nil {
			subject = subjectModel.Name
			log.Printf("[Defect Create] SubjectID=%d, SubjectName=%s", *req.SubjectID, subject)
		} else {
			log.Printf("[Defect Create] Failed to find subject: subject_id=%d, error=%v", *req.SubjectID, err)
		}
	}

	// 处理Phase：如果提供了PhaseID，查找名称
	phase := req.Phase
	if req.PhaseID != nil && *req.PhaseID > 0 {
		var phaseModel models.DefectPhase
		if err := s.repo.GetDB().Where("project_id = ?", projectID).First(&phaseModel, *req.PhaseID).Error; err == nil {
			phase = phaseModel.Name
			log.Printf("[Defect Create] PhaseID=%d, PhaseName=%s", *req.PhaseID, phase)
		} else {
			log.Printf("[Defect Create] Failed to find phase: phase_id=%d, error=%v", *req.PhaseID, err)
		}
	}

	defect := &models.Defect{
		DefectID:          defectID,
		ProjectID:         projectID,
		Title:             req.Title,
		Subject:           subject,
		Description:       req.Description,
		RecoveryMethod:    req.RecoveryMethod,
		Priority:          req.Priority,
		Severity:          req.Severity,
		Frequency:         req.Frequency,
		DetectedInRelease: req.DetectedInRelease,
		Phase:             phase,
		Status:            string(models.DefectStatusNew),
		CreatedBy:         userID,
		UpdatedBy:         userID,
	}

	if err := s.repo.Create(defect); err != nil {
		return nil, fmt.Errorf("create defect: %w", err)
	}

	log.Printf("[Defect Create] user_id=%d, project_id=%d, defect_id=%s", userID, projectID, defectID)
	return defect, nil
}

// GetByID 根据UUID获取缺陷
func (s *defectService) GetByID(id string) (*models.Defect, error) {
	defect, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("defect not found")
		}
		return nil, fmt.Errorf("get defect: %w", err)
	}
	return defect, nil
}

// GetByDefectID 根据显示ID获取缺陷
func (s *defectService) GetByDefectID(defectID string) (*models.Defect, error) {
	defect, err := s.repo.GetByDefectID(defectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("defect not found")
		}
		return nil, fmt.Errorf("get defect: %w", err)
	}
	return defect, nil
}

// Update 更新缺陷
func (s *defectService) Update(id string, userID uint, req *models.DefectUpdateRequest) error {
	// 先获取缺陷以得到projectID
	defect, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("defect not found")
		}
		return fmt.Errorf("get defect: %w", err)
	}

	// 构建更新字段
	updates := make(map[string]interface{})

	if req.Title != nil {
		if *req.Title == "" {
			return errors.New("title cannot be empty")
		}
		updates["title"] = *req.Title
	}

	// 处理Subject：如果提供了SubjectID，查找名称
	if req.SubjectID != nil {
		if *req.SubjectID > 0 {
			var subjectModel models.DefectSubject
			if err := s.repo.GetDB().Where("project_id = ?", defect.ProjectID).First(&subjectModel, *req.SubjectID).Error; err == nil {
				updates["subject"] = subjectModel.Name
			}
		} else {
			updates["subject"] = ""
		}
	} else if req.Subject != nil {
		updates["subject"] = *req.Subject
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.RecoveryMethod != nil {
		updates["recovery_method"] = *req.RecoveryMethod
	}
	if req.Priority != nil {
		if !models.IsValidDefectPriority(*req.Priority) {
			return errors.New("invalid priority value")
		}
		updates["priority"] = *req.Priority
	}
	if req.Severity != nil {
		if !models.IsValidDefectSeverity(*req.Severity) {
			return errors.New("invalid severity value")
		}
		updates["severity"] = *req.Severity
	}
	if req.Frequency != nil {
		updates["frequency"] = *req.Frequency
	}
	if req.DetectedInRelease != nil {
		updates["detected_in_release"] = *req.DetectedInRelease
	}

	// 处理Phase：如果提供了PhaseID，查找名称
	if req.PhaseID != nil {
		if *req.PhaseID > 0 {
			var phaseModel models.DefectPhase
			if err := s.repo.GetDB().Where("project_id = ?", defect.ProjectID).First(&phaseModel, *req.PhaseID).Error; err == nil {
				updates["phase"] = phaseModel.Name
			}
		} else {
			updates["phase"] = ""
		}
	} else if req.Phase != nil {
		updates["phase"] = *req.Phase
	}
	if req.Status != nil {
		if !models.IsValidDefectStatus(*req.Status) {
			return errors.New("invalid status value")
		}
		oldStatus := defect.Status
		updates["status"] = *req.Status
		log.Printf("[Defect Status Change] defect_id=%s, from=%s, to=%s, user_id=%d", defect.DefectID, oldStatus, *req.Status, userID)
	}

	updates["updated_by"] = userID

	if len(updates) == 1 { // 只有updated_by
		return nil
	}

	if err := s.repo.Update(id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("defect not found")
		}
		return fmt.Errorf("update defect: %w", err)
	}

	log.Printf("[Defect Update] user_id=%d, defect_id=%s, fields=%v", userID, id, updates)
	return nil
}

// Delete 删除缺陷
func (s *defectService) Delete(id string) error {
	// 先获取缺陷信息用于日志
	defect, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("defect not found")
		}
		return fmt.Errorf("get defect: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete defect: %w", err)
	}

	log.Printf("[Defect Delete] defect_id=%s", defect.DefectID)
	return nil
}

// List 分页查询缺陷列表
func (s *defectService) List(projectID uint, status, keyword string, page, size int) (*models.DefectListResponse, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 50
	}

	// 验证状态
	if status != "" && !models.IsValidDefectStatus(status) {
		return nil, errors.New("invalid status value")
	}

	defects, total, err := s.repo.List(projectID, status, keyword, page, size)
	if err != nil {
		return nil, fmt.Errorf("list defects: %w", err)
	}

	statusCounts, err := s.repo.GetStatusCounts(projectID)
	if err != nil {
		return nil, fmt.Errorf("get status counts: %w", err)
	}

	return &models.DefectListResponse{
		Defects:      convertToDefectSlice(defects),
		Total:        total,
		Page:         page,
		Size:         size,
		StatusCounts: statusCounts,
	}, nil
}

// convertToDefectSlice 转换指针切片为值切片
func convertToDefectSlice(defects []*models.Defect) []models.Defect {
	result := make([]models.Defect, len(defects))
	for i, d := range defects {
		result[i] = *d
	}
	return result
}

// CSV列定义
var csvHeaders = []string{
	"Title", "Subject", "Description", "Recovery Method",
	"Priority", "Severity", "Frequency", "Detected In Release", "Phase",
	"Status", "Created At",
}

// GenerateTemplate 生成CSV导入模板
func (s *defectService) GenerateTemplate() ([]byte, error) {
	var buf bytes.Buffer

	// 添加UTF-8 BOM头，确保Excel正确识别编码
	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(&buf)

	// 写入表头
	if err := writer.Write(csvHeaders); err != nil {
		return nil, fmt.Errorf("write csv header: %w", err)
	}

	// 写入说明行（字段要求）
	instructions := []string{
		"(Required)",               // Title - 必填
		"(Optional)",               // Subject - 可选，模块名称
		"(Optional)",               // Description - 可选
		"(Optional)",               // Recovery Method - 可选
		"(A/B/C/D)",                // Priority - 必须是A/B/C/D之一
		"(A/B/C/D)",                // Severity - 必须是A/B/C/D之一
		"(Optional, e.g., 100%)",   // Frequency - 可选
		"(Optional, e.g., v1.0.0)", // Detected In Release - 可选
		"(Optional)",               // Phase - 可选，阶段名称
		"(Optional: New/Active/Resolved/Closed, default: New)", // Status - 可选
		"(Auto-generated or YYYY-MM-DD)",                       // Created At - 自动生成或指定日期
	}
	if err := writer.Write(instructions); err != nil {
		return nil, fmt.Errorf("write csv instructions: %w", err)
	}

	// 写入示例数据
	example := []string{
		"登录页面无法输入用户名",      // Title
		"登录模块",             // Subject
		"用户无法在用户名输入框中输入内容", // Description
		"刷新页面后重试",          // Recovery Method
		"A",                // Priority (必须是A/B/C/D)
		"B",                // Severity (必须是A/B/C/D)
		"100%",             // Frequency
		"v1.0.0",           // Detected In Release
		"系统测试",             // Phase
		"New",              // Status (可选：New/Active/Resolved/Closed)
		"2025-12-03",       // Created At
	}
	if err := writer.Write(example); err != nil {
		return nil, fmt.Errorf("write csv example: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("flush csv writer: %w", err)
	}

	return buf.Bytes(), nil
}

// Import 导入缺陷
func (s *defectService) Import(projectID uint, userID uint, reader io.Reader) (*models.ImportResult, error) {
	// 读取所有内容以检测和移除BOM
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// 检测并移除UTF-8 BOM (EF BB BF)
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
		log.Printf("[Defect Import] UTF-8 BOM detected and removed")
	}

	// 创建CSV读取器
	csvReader := csv.NewReader(bytes.NewReader(data))

	// 读取表头
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("read csv header: %w", err)
	}

	// 清理表头（移除可能的BOM残留）
	if len(headers) > 0 {
		headers[0] = strings.TrimPrefix(headers[0], "\ufeff") // 移除UTF-8 BOM字符
	}

	// 验证表头
	if len(headers) < len(csvHeaders) {
		return nil, errors.New("csv format error: insufficient columns")
	}

	result := &models.ImportResult{
		Errors: make([]models.ImportError, 0),
	}

	rowNum := 1 // 表头为第1行，数据从第2行开始

	for {
		rowNum++
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, models.ImportError{
				Row:    rowNum,
				Reason: fmt.Sprintf("read row error: %v", err),
			})
			continue
		}

		// 验证必填字段
		if len(record) < 1 || strings.TrimSpace(record[0]) == "" {
			result.FailCount++
			result.Errors = append(result.Errors, models.ImportError{
				Row:    rowNum,
				Reason: "Title is required",
			})
			continue
		}

		// 构建创建请求
		req := &models.DefectCreateRequest{
			Title: strings.TrimSpace(record[0]),
		}

		if len(record) > 1 {
			req.Subject = strings.TrimSpace(record[1])
		}
		if len(record) > 2 {
			req.Description = strings.TrimSpace(record[2])
		}
		if len(record) > 3 {
			req.RecoveryMethod = strings.TrimSpace(record[3])
		}
		if len(record) > 4 {
			priority := strings.TrimSpace(record[4])
			if priority != "" && !models.IsValidDefectPriority(priority) {
				result.FailCount++
				result.Errors = append(result.Errors, models.ImportError{
					Row:    rowNum,
					Reason: "Invalid Priority value",
				})
				continue
			}
			req.Priority = priority
		}
		if len(record) > 5 {
			severity := strings.TrimSpace(record[5])
			if severity != "" && !models.IsValidDefectSeverity(severity) {
				result.FailCount++
				result.Errors = append(result.Errors, models.ImportError{
					Row:    rowNum,
					Reason: "Invalid Severity value",
				})
				continue
			}
			req.Severity = severity
		}
		if len(record) > 6 {
			req.Frequency = strings.TrimSpace(record[6])
		}
		if len(record) > 7 {
			req.DetectedInRelease = strings.TrimSpace(record[7])
		}
		if len(record) > 8 {
			req.Phase = strings.TrimSpace(record[8])
		}

		// 解析 Status（可选，默认为 New）
		statusStr := ""
		if len(record) > 9 {
			statusStr = strings.TrimSpace(record[9])
			if statusStr != "" && statusStr != "(Optional)" {
				if !models.IsValidDefectStatus(statusStr) {
					result.FailCount++
					result.Errors = append(result.Errors, models.ImportError{
						Row:    rowNum,
						Reason: fmt.Sprintf("Invalid Status value: %s (expected: New/Active/Resolved/Closed)", statusStr),
					})
					continue
				}
			}
		}

		// 解析创建日期（支持多种格式）
		var createdAt time.Time
		var importCreatedBy uint = userID // 默认使用导入用户

		if len(record) > 10 {
			createdAtStr := strings.TrimSpace(record[10])
			if createdAtStr != "" && createdAtStr != "(Auto-generated)" {
				// 尝试多种日期格式
				formats := []string{
					"2006-01-02 15:04:05",
					"2006-01-02 15:04",
					"2006-01-02",
					"2006/01/02 15:04:05",
					"2006/01/02 15:04",
					"2006/01/02",
				}

				parsed := false
				for _, format := range formats {
					if t, err := time.Parse(format, createdAtStr); err == nil {
						createdAt = t
						parsed = true
						break
					}
				}

				if !parsed {
					result.FailCount++
					result.Errors = append(result.Errors, models.ImportError{
						Row:    rowNum,
						Reason: fmt.Sprintf("Invalid date format: %s (expected: YYYY-MM-DD HH:MM:SS or YYYY-MM-DD)", createdAtStr),
					})
					continue
				}
			}
		}

		// 创建缺陷（使用导入用户ID作为创建人）
		defect, err := s.Create(projectID, importCreatedBy, req)
		if err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, models.ImportError{
				Row:    rowNum,
				Reason: err.Error(),
			})
			continue
		}

		// 如果CSV提供了 Status，更新数据库中的 status
		if statusStr != "" && statusStr != "(Optional)" {
			if err := s.repo.GetDB().Model(&models.Defect{}).
				Where("id = ?", defect.ID).
				Update("status", statusStr).Error; err != nil {
				log.Printf("[Defect Import] Failed to update status for defect %s: %v", defect.DefectID, err)
			}
		}

		// 如果CSV提供了创建日期，更新数据库中的created_at
		if !createdAt.IsZero() {
			if err := s.repo.GetDB().Model(&models.Defect{}).
				Where("id = ?", defect.ID).
				Update("created_at", createdAt).Error; err != nil {
				log.Printf("[Defect Import] Failed to update created_at for defect %s: %v", defect.DefectID, err)
			}
		}

		result.SuccessCount++
	}

	log.Printf("[Defect Import] user_id=%d, project_id=%d, success=%d, fail=%d",
		userID, projectID, result.SuccessCount, result.FailCount)

	return result, nil
}

// Export 导出缺陷
func (s *defectService) Export(projectID uint) ([]byte, error) {
	defects, err := s.repo.GetByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("get defects: %w", err)
	}

	// 获取所有创建人ID
	creatorIDs := make(map[uint]bool)
	for _, defect := range defects {
		creatorIDs[defect.CreatedBy] = true
	}

	// 批量查询创建人信息
	var users []models.User
	userIDList := make([]uint, 0, len(creatorIDs))
	for id := range creatorIDs {
		userIDList = append(userIDList, id)
	}
	if len(userIDList) > 0 {
		if err := s.repo.GetDB().Where("id IN ?", userIDList).Find(&users).Error; err != nil {
			return nil, fmt.Errorf("get creators: %w", err)
		}
	}

	// 构建用户ID到昵称的映射
	userMap := make(map[uint]string)
	for _, user := range users {
		userMap[user.ID] = user.Nickname
	}

	var buf bytes.Buffer
	// 添加UTF-8 BOM头，解决Excel乱码问题
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(&buf)

	// 写入表头（包含 Defect ID 和 Status，其他字段按 csvHeaders 顺序）
	exportHeaders := []string{"Defect ID", "Title", "Subject", "Description", "Recovery Method",
		"Priority", "Severity", "Frequency", "Detected In Release", "Phase",
		"Status", "Created At", "Created By"}
	if err := writer.Write(exportHeaders); err != nil {
		return nil, fmt.Errorf("write csv header: %w", err)
	}

	// 写入数据
	for _, defect := range defects {
		// 格式化创建时间
		createdAt := defect.CreatedAt.Format("2006-01-02 15:04:05")

		// 获取创建人昵称
		createdBy := userMap[defect.CreatedBy]
		if createdBy == "" {
			createdBy = fmt.Sprintf("User#%d", defect.CreatedBy)
		}

		// Subject和Phase在创建时已经存储为名称字符串
		row := []string{
			defect.DefectID,
			defect.Title,
			defect.Subject,
			defect.Description,
			defect.RecoveryMethod,
			defect.Priority,
			defect.Severity,
			defect.Frequency,
			defect.DetectedInRelease,
			defect.Phase,
			defect.Status,
			createdAt,
			createdBy,
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("write csv row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("flush csv writer: %w", err)
	}

	return buf.Bytes(), nil
}
