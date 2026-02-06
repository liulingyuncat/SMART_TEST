package services

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"strings"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"github.com/xuri/excelize/v2"
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
	GenerateTemplate(format string) ([]byte, error)
	ImportWithFormat(projectID uint, userID uint, reader io.Reader, isXLSX bool) (*models.ImportResult, error)
	Import(projectID uint, userID uint, reader io.Reader) (*models.ImportResult, error)
	Export(projectID uint) ([]byte, error)
	ExportWithFormat(projectID uint, format string) ([]byte, error)
}

type defectService struct {
	repo     repositories.DefectRepository
	userRepo repositories.UserRepository
}

// NewDefectService 创建缺陷服务实例
func NewDefectService(repo repositories.DefectRepository, userRepo repositories.UserRepository) DefectService {
	return &defectService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// generateDefectID 生成缺陷显示ID
func (s *defectService) generateDefectID(projectID uint) (string, error) {
	// 使用原子性方法生成缺陷ID，确保在并发场景下不会重复
	return s.repo.GenerateNextDefectID(projectID)
}

// decodeHTMLEntities 解码HTML实体（处理&quot;、&#39;等）
func decodeHTMLEntities(text string) string {
	if text == "" {
		return text
	}
	return html.UnescapeString(text)
}

// Create 创建缺陷
func (s *defectService) Create(projectID uint, userID uint, req *models.DefectCreateRequest) (*models.Defect, error) {
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

	// 生成DefectID（原子性）
	defectID, err := s.generateDefectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("generate defect id: %w", err)
	}

	// 处理状态：如果提供了有效状态则使用，否则默认New
	status := string(models.DefectStatusNew)
	if req.Status != "" && models.IsValidDefectStatus(req.Status) {
		status = req.Status
	}

	// 处理DetectedBy：如果请求中有值则使用，否则使用创建人的用户名
	detectedBy := req.DetectedBy
	if detectedBy == "" {
		// 查询创建人的用户名
		user, err := s.userRepo.FindByID(userID)
		if err == nil && user != nil {
			detectedBy = user.Username
		}
	}

	defect := &models.Defect{
		DefectID:        defectID,
		ProjectID:       projectID,
		Title:           req.Title,
		Subject:         subject,
		Description:     req.Description,
		RecoveryMethod:  req.RecoveryMethod,
		Priority:        req.Priority,
		Severity:        req.Severity,
		Type:            req.Type,
		Frequency:       req.Frequency,
		DetectedVersion: req.DetectedVersion,
		Phase:           phase,
		CaseID:          req.CaseID,
		RecoveryRank:    req.RecoveryRank,
		DetectionTeam:   req.DetectionTeam,
		Location:        req.Location,
		FixVersion:      req.FixVersion,
		SQAMemo:         req.SQAMemo,
		Component:       req.Component,
		Resolution:      req.Resolution,
		Models:          req.Models,
		DetectedBy:      detectedBy,
		Status:          status,
		CreatedBy:       userID,
		UpdatedBy:       userID,
	}

	// 处理CreatedAt：如果提供了有效日期则使用
	if req.CreatedAt != "" {
		log.Printf("[Defect Create] Parsing CreatedAt: '%s'", req.CreatedAt)
		// 尝试多种日期格式
		formats := []string{
			"2006-01-02",
			"2006-01-02 15:04:05",
			"2006/01/02", // Excel常用格式
			"2006/1/2",   // Excel可能的短格式
			"01-02-06",   // Excel短格式：月-日-年(YY)
			"1-2-06",     // Excel短格式：月-日-年(YY)，单数字
			time.RFC3339,
		}
		parsed := false
		for _, format := range formats {
			if parsedTime, err := time.Parse(format, req.CreatedAt); err == nil {
				defect.CreatedAt = parsedTime
				log.Printf("[Defect Create] Successfully parsed CreatedAt with format '%s': %v", format, parsedTime)
				parsed = true
				break
			}
		}
		if !parsed {
			log.Printf("[Defect Create] WARNING: Failed to parse CreatedAt: '%s'", req.CreatedAt)
		}
	}

	if err = s.repo.Create(defect); err != nil {
		return nil, fmt.Errorf("create defect: %w", err)
	}

	log.Printf("[Defect Create] user_id=%d, project_id=%d, defect_id=%s, created_at=%v", userID, projectID, defect.DefectID, defect.CreatedAt)
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
	if req.Type != nil {
		if *req.Type != "" && !models.IsValidDefectType(*req.Type) {
			return errors.New("invalid type value")
		}
		updates["type"] = *req.Type
	}
	if req.Frequency != nil {
		updates["frequency"] = *req.Frequency
	}
	if req.DetectedVersion != nil {
		updates["detected_version"] = *req.DetectedVersion
	}
	if req.RecoveryRank != nil {
		updates["recovery_rank"] = *req.RecoveryRank
	}
	if req.DetectionTeam != nil {
		updates["detection_team"] = *req.DetectionTeam
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.FixVersion != nil {
		updates["fix_version"] = *req.FixVersion
	}
	if req.SQAMemo != nil {
		updates["sqa_memo"] = *req.SQAMemo
	}
	if req.Component != nil {
		updates["component"] = *req.Component
	}
	if req.Resolution != nil {
		updates["resolution"] = *req.Resolution
	}
	if req.Models != nil {
		updates["models"] = *req.Models
	}
	if req.DetectedBy != nil {
		updates["detected_by"] = *req.DetectedBy
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
	if req.CaseID != nil {
		updates["case_id"] = *req.CaseID
	}
	if req.Assignee != nil {
		updates["assignee"] = *req.Assignee
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
	if size < 1 {
		size = 50
	} else if size > 100000 {
		size = 100000
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

// getExcelColumn 将列索引转换为Excel列名（1-based: 1->A, 2->B, ..., 27->AA）
func getExcelColumn(col int) string {
	result := ""
	for col > 0 {
		col--
		result = string(rune('A'+col%26)) + result
		col /= 26
	}
	return result
}

// CSV列定义（除ID外的全部字段）
var csvHeaders = []string{
	"Title", "Module", "Description", "Recovery Method",
	"Priority", "Severity", "Type", "Frequency", "Detected Version", "Phase",
	"Detection Team", "Location", "Fix Version", "SQA MEMO", "Component",
	"Resolution", "Models",
	"Status", "Created At", "Detected By",
}

// GenerateTemplate 生成导入模板（支持CSV和XLSX格式）
func (s *defectService) GenerateTemplate(format string) ([]byte, error) {
	if format == "xlsx" {
		return s.generateXLSXTemplate()
	}
	// 默认返回CSV
	return s.generateCSVTemplate()
}

// generateCSVTemplate 生成CSV模板
func (s *defectService) generateCSVTemplate() ([]byte, error) {
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
		"(Required)", // Title - 必填
		"(Optional)", // Module - 可选，模块名称
		"(Optional)", // Description - 可选
		"(Optional)", // Recovery Method - 可选
		"(A/B/C/D)",  // Priority - 必须是A/B/C/D之一
		"(Critical/Major/Minor/Trivial or A/B/C/D)", // Severity - 新值或旧值
		"(Optional: Functional/UI/UIInteraction/Compatibility/BrowserSpecific/Performance/Security/Environment/UserError)", // Type - 可选
		"(Optional, e.g., 100%)",   // Frequency - 可选
		"(Optional, e.g., v1.0.0)", // Detected Version - 可选
		"(Optional)",               // Phase - 可选，阶段名称
		"(Optional)",               // Detection Team - 可选
		"(Optional)",               // Location - 可选
		"(Optional, e.g., v1.1.0)", // Fix Version - 可选
		"(Optional)",               // SQA MEMO - 可选
		"(Optional)",               // Component - 可选
		"(Optional)",               // Resolution - 可选，解决方案
		"(Optional)",               // Models - 可选，机型
		"(Optional: New/InProgress/Confirmed/Resolved/Reopened/Rejected/Closed, default: New)", // Status - 可选
		"(Auto-generated or YYYY-MM-DD)", // Created At - 自动生成或指定日期
		"(Optional)",                     // Detected By - 可选，提出人
	}
	if err := writer.Write(instructions); err != nil {
		return nil, fmt.Errorf("write csv instructions: %w", err)
	}

	// 写入示例数据
	example := []string{
		"登录页面无法输入用户名",           // Title
		"登录模块",                  // Module
		"用户无法在用户名输入框中输入内容",      // Description
		"刷新页面后重试",               // Recovery Method
		"A",                     // Priority
		"Major",                 // Severity
		"Functional",            // Type
		"100%",                  // Frequency
		"v1.0.0",                // Detected Version
		"系统测试",                  // Phase
		"QA Team A",             // Detection Team
		"Login Page",            // Location
		"v1.1.0",                // Fix Version
		"Needs code review",     // SQA MEMO
		"Auth Module",           // Component
		"Clear cache and retry", // Resolution
		"Chrome, Firefox",       // Models
		"New",                   // Status
		"2025-12-03",            // Created At
		"test_user",             // Detected By
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

// generateXLSXTemplate 生成XLSX模板（极简版）
func (s *defectService) generateXLSXTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// 获取默认工作表
	sheetName := f.GetSheetName(0)

	// 写入表头
	for col, header := range csvHeaders {
		cell := fmt.Sprintf("%s1", getExcelColumn(col+1))
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入说明行
	instructions := []string{
		"(Required)", // Title - 必填
		"(Optional)", // Module - 可选
		"(Optional)", // Description - 可选
		"(Optional)", // Recovery Method - 可选
		"(A/B/C/D)",  // Priority
		"(Critical/Major/Minor/Trivial or A/B/C/D)", // Severity
		"(Optional: Functional/UI/UIInteraction/Compatibility/BrowserSpecific/Performance/Security/Environment/UserError)", // Type
		"(Optional, e.g., 100%)",   // Frequency
		"(Optional, e.g., v1.0.0)", // Detected Version
		"(Optional)",               // Phase
		"(Optional)",               // Detection Team
		"(Optional)",               // Location
		"(Optional, e.g., v1.1.0)", // Fix Version
		"(Optional)",               // SQA MEMO
		"(Optional)",               // Component
		"(Optional)",               // Resolution
		"(Optional)",               // Models
		"(Optional: New/InProgress/Confirmed/Resolved/Reopened/Rejected/Closed, default: New)", // Status
		"(Auto-generated or YYYY-MM-DD)", // Created At
		"(Optional)",                     // Detected By
	}
	for col, instr := range instructions {
		cell := fmt.Sprintf("%s2", getExcelColumn(col+1))
		f.SetCellValue(sheetName, cell, instr)
	}

	// 写入示例数据
	example := []string{
		"登录页面无法输入用户名",           // Title
		"登录模块",                  // Module
		"用户无法在用户名输入框中输入内容",      // Description
		"刷新页面后重试",               // Recovery Method
		"A",                     // Priority
		"Major",                 // Severity
		"Functional",            // Type
		"100%",                  // Frequency
		"v1.0.0",                // Detected Version
		"系统测试",                  // Phase
		"QA Team",               // Detection Team
		"Login Page",            // Location
		"v1.1.0",                // Fix Version
		"需要优先修复",                // SQA MEMO
		"Login Component",       // Component
		"Clear cache and retry", // Resolution
		"Chrome, Firefox",       // Models
		"New",                   // Status
		"2025-12-03",            // Created At
		"test_user",             // Detected By
	}
	for col, value := range example {
		cell := fmt.Sprintf("%s3", getExcelColumn(col+1))
		f.SetCellValue(sheetName, cell, value)
	}

	// 保存到缓冲区
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("write xlsx: %w", err)
	}

	return buf.Bytes(), nil
}

// ImportWithFormat 根据格式导入缺陷
func (s *defectService) ImportWithFormat(projectID uint, userID uint, reader io.Reader, isXLSX bool) (*models.ImportResult, error) {
	if isXLSX {
		return s.importXLSX(projectID, userID, reader)
	}
	return s.Import(projectID, userID, reader)
}

// importXLSX 直接导入 XLSX 文件
func (s *defectService) importXLSX(projectID uint, userID uint, reader io.Reader) (*models.ImportResult, error) {
	// 读取 XLSX 文件
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read xlsx file: %w", err)
	}

	// 打开 XLSX 文件
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("open xlsx file: %w", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)

	// 获取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("get rows from xlsx: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("xlsx file is empty")
	}

	// 清理重复的 defect_id（防止主键冲突）
	s.repo.CleanupDuplicateDefects(projectID, "000001")
	log.Printf("[Defect Import XLSX] Cleaned up duplicate defect IDs before importing")

	// 直接处理 XLSX 数据（跳过表头和说明行）
	headers := rows[0]
	log.Printf("[Defect Import XLSX DEBUG] XLSX file loaded. Total rows: %d. Headers: %v", len(rows), headers)

	result := &models.ImportResult{
		SuccessCount: 0,
		FailCount:    0,
		Errors:       []models.ImportError{},
	}

	// 检测说明行：如果第2行包含 "(Required)" 或 "(Optional)" 则认为存在说明行
	startRowIdx := 1 // 默认从第2行(索引1)开始处理数据
	if len(rows) > 1 {
		secondRow := rows[1]
		isDescriptionRow := false
		for _, cell := range secondRow {
			cellVal := strings.TrimSpace(cell)
			if strings.Contains(cellVal, "(Required)") || strings.Contains(cellVal, "(Optional)") ||
				strings.Contains(cellVal, "(Auto-generated)") {
				isDescriptionRow = true
				break
			}
		}
		if isDescriptionRow {
			startRowIdx = 2 // 如果有说明行，从第3行(索引2)开始
			log.Printf("[Defect Import XLSX DEBUG] Description row detected at row 2. Starting data processing from row 3")
		} else {
			log.Printf("[Defect Import XLSX DEBUG] No description row detected. Starting data processing from row 2")
		}
	}

	// 遍历数据行
	for rowIdx := startRowIdx; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		excelRowNum := rowIdx + 1
		dataRowNum := rowIdx - startRowIdx + 1

		// 检查是否为空行
		isEmpty := true
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			log.Printf("[Defect Import XLSX DEBUG] Row %d (Excel row %d) SKIPPED: empty row", dataRowNum, excelRowNum)
			continue
		}

		// 将行数据映射到缺陷请求
		req := &models.DefectCreateRequest{}
		var defectID string

		// 映射列
		for colIdx, header := range headers {
			var value string
			if colIdx < len(row) {
				value = strings.TrimSpace(row[colIdx])
				// 解码HTML实体
				value = decodeHTMLEntities(value)
			}

			switch strings.TrimSpace(header) {
			case "Defect ID":
				defectID = value
			case "Title":
				req.Title = value
			case "Module":
				req.Subject = value
			case "Description":
				req.Description = value
			case "Recovery Method":
				req.RecoveryMethod = value
			case "Priority":
				if value != "" && models.IsValidDefectPriority(value) {
					req.Priority = value
				}
			case "Severity":
				if value != "" && models.IsValidDefectSeverity(value) {
					req.Severity = value
				}
			case "Type":
				if value != "" && models.IsValidDefectType(value) {
					req.Type = value
				}
			case "Frequency":
				req.Frequency = value
			case "Detected Version":
				req.DetectedVersion = value
			case "Phase":
				req.Phase = value
			case "Detection Team":
				req.DetectionTeam = value
			case "Location":
				req.Location = value
			case "Fix Version":
				req.FixVersion = value
			case "SQA MEMO":
				req.SQAMemo = value
			case "Component":
				req.Component = value
			case "Resolution":
				req.Resolution = value
			case "Models":
				req.Models = value
			case "Status":
				req.Status = value
			case "Created At":
				req.CreatedAt = value
			case "Detected By":
				req.DetectedBy = value
			}
		}

		// 验证必填字段
		if strings.TrimSpace(req.Title) == "" {
			log.Printf("[Defect Import XLSX DEBUG] Row %d (Excel row %d) SKIPPED: empty Title", dataRowNum, excelRowNum)
			continue
		}

		// 如果有Defect ID，尝试更新已有缺陷
		if defectID != "" {
			existingDefect, err := s.repo.GetByDefectID(defectID)
			if err == nil && existingDefect != nil {
				// 找到已有缺陷，进行更新
				updateReq := &models.DefectUpdateRequest{
					Title:           &req.Title,
					Subject:         &req.Subject,
					Description:     &req.Description,
					RecoveryMethod:  &req.RecoveryMethod,
					Priority:        &req.Priority,
					Severity:        &req.Severity,
					Type:            &req.Type,
					Frequency:       &req.Frequency,
					DetectedVersion: &req.DetectedVersion,
					Phase:           &req.Phase,
					DetectionTeam:   &req.DetectionTeam,
					Location:        &req.Location,
					FixVersion:      &req.FixVersion,
					SQAMemo:         &req.SQAMemo,
					Component:       &req.Component,
					Resolution:      &req.Resolution,
					Models:          &req.Models,
					DetectedBy:      &req.DetectedBy,
					Status:          &req.Status,
				}
				err := s.Update(existingDefect.ID, userID, updateReq)
				if err != nil {
					log.Printf("[Defect Import XLSX DEBUG] Row %d (Excel row %d) UPDATE FAILED: Defect ID=%s, error=%v", dataRowNum, excelRowNum, defectID, err)
					continue
				}
				log.Printf("[Defect Import XLSX DEBUG] Row %d (Excel row %d) UPDATED: Defect ID=%s", dataRowNum, excelRowNum, defectID)
				result.SuccessCount++
				continue
			}
		}

		// 创建新缺陷
		_, err := s.Create(projectID, userID, req)
		if err != nil {
			log.Printf("[Defect Import XLSX DEBUG] Row %d (Excel row %d) CREATE FAILED: %v", dataRowNum, excelRowNum, err)
			continue
		}
		log.Printf("[Defect Import XLSX DEBUG] Row %d (Excel row %d) CREATED", dataRowNum, excelRowNum)
		result.SuccessCount++
	}

	log.Printf("[Defect Import XLSX] ===== COMPLETE ===== success=%d", result.SuccessCount)
	return result, nil
}

// Import 导入缺陷（CSV格式 - 支持有无说明行）
func (s *defectService) Import(projectID uint, userID uint, reader io.Reader) (*models.ImportResult, error) {
	// 读取所有内容以检测和移除BOM
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// 检测并移除UTF-8 BOM (EF BB BF)
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
		log.Printf("[Defect Import CSV] UTF-8 BOM removed")
	}

	// 创建CSV读取器并读取所有行
	csvReader := csv.NewReader(bytes.NewReader(data))
	allRows, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv file: %w", err)
	}

	if len(allRows) == 0 {
		return nil, errors.New("csv file is empty")
	}

	result := &models.ImportResult{
		Errors: make([]models.ImportError, 0),
	}

	// 第1行是表头
	headers := allRows[0]
	if len(headers) > 0 {
		headers[0] = strings.TrimPrefix(headers[0], "\ufeff")
	}

	// 验证表头
	if len(headers) < len(csvHeaders) {
		return nil, errors.New("csv format error: insufficient columns")
	}
	log.Printf("[Defect Import CSV DEBUG] Headers validated, starting to process data rows")

	// 清理重复的 defect_id（防止主键冲突）
	s.repo.CleanupDuplicateDefects(projectID, "000001")
	log.Printf("[Defect Import CSV] Cleaned up duplicate defect IDs before importing")

	// 确定数据行起始位置 - 智能检测是否有说明行
	dataStartIdx := 1 // 默认从第2行(索引1)开始
	if len(allRows) > 1 {
		isDescriptionRow := false
		for _, cell := range allRows[1] {
			cellVal := strings.TrimSpace(cell)
			if strings.Contains(cellVal, "(Required)") || strings.Contains(cellVal, "(Optional)") ||
				strings.Contains(cellVal, "(Auto-generated)") {
				isDescriptionRow = true
				break
			}
		}
		if isDescriptionRow {
			dataStartIdx = 2
			log.Printf("[Defect Import CSV DEBUG] Description row detected - starting from row 3")
		} else {
			log.Printf("[Defect Import CSV DEBUG] No description row - starting from row 2")
		}
	}

	// 检查是否有Defect ID列（导出文件会有）
	hasDefectIDColumn := false
	if len(headers) > 0 && strings.TrimSpace(headers[0]) == "Defect ID" {
		hasDefectIDColumn = true
		log.Printf("[Defect Import CSV DEBUG] Defect ID column detected - will support update")
	}

	// 处理数据行
	for rowIdx := dataStartIdx; rowIdx < len(allRows); rowIdx++ {
		record := allRows[rowIdx]
		excelRowNum := rowIdx + 1
		dataRowNum := rowIdx - dataStartIdx + 1

		// 解析Defect ID（如果存在）
		var defectID string
		titleIdx := 0
		if hasDefectIDColumn {
			if len(record) > 0 {
				defectID = strings.TrimSpace(record[0])
			}
			titleIdx = 1
		}

		// 验证必填字段
		if len(record) <= titleIdx || strings.TrimSpace(record[titleIdx]) == "" {
			log.Printf("[Defect Import CSV DEBUG] Row %d (Excel row %d) SKIPPED: empty Title", dataRowNum, excelRowNum)
			continue
		}

		// 构建创建请求
		req := &models.DefectCreateRequest{
			Title: decodeHTMLEntities(strings.TrimSpace(record[titleIdx])),
		}

		if len(record) > titleIdx+1 {
			req.Subject = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+1]))
		}
		if len(record) > titleIdx+2 {
			req.Description = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+2]))
		}
		if len(record) > titleIdx+3 {
			req.RecoveryMethod = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+3]))
		}
		if len(record) > titleIdx+4 {
			priority := strings.TrimSpace(record[titleIdx+4])
			if priority != "" && models.IsValidDefectPriority(priority) {
				req.Priority = priority
			}
		}
		if len(record) > titleIdx+5 {
			severity := strings.TrimSpace(record[titleIdx+5])
			if severity != "" && models.IsValidDefectSeverity(severity) {
				req.Severity = severity
			}
		}
		if len(record) > titleIdx+6 {
			defectType := strings.TrimSpace(record[titleIdx+6])
			if defectType != "" && models.IsValidDefectType(defectType) {
				req.Type = defectType
			}
		}
		if len(record) > titleIdx+7 {
			req.Frequency = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+7]))
		}
		if len(record) > titleIdx+8 {
			req.DetectedVersion = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+8]))
		}
		if len(record) > titleIdx+9 {
			req.Phase = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+9]))
		}
		if len(record) > titleIdx+10 {
			req.DetectionTeam = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+10]))
		}
		if len(record) > titleIdx+11 {
			req.Location = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+11]))
		}
		if len(record) > titleIdx+12 {
			req.FixVersion = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+12]))
		}
		if len(record) > titleIdx+13 {
			req.SQAMemo = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+13]))
		}
		if len(record) > titleIdx+14 {
			req.Component = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+14]))
		}
		if len(record) > titleIdx+15 {
			req.Resolution = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+15]))
		}
		if len(record) > titleIdx+16 {
			req.Models = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+16]))
		}
		if len(record) > titleIdx+17 {
			req.Status = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+17]))
		}
		if len(record) > titleIdx+18 {
			req.CreatedAt = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+18]))
		}
		if len(record) > titleIdx+19 {
			req.DetectedBy = decodeHTMLEntities(strings.TrimSpace(record[titleIdx+19]))
		}

		// 如果有Defect ID，尝试更新已有缺陷
		if defectID != "" {
			existingDefect, err := s.repo.GetByDefectID(defectID)
			if err == nil && existingDefect != nil {
				// 找到已有缺陷，进行更新
				updateReq := &models.DefectUpdateRequest{
					Title:           &req.Title,
					Subject:         &req.Subject,
					Description:     &req.Description,
					RecoveryMethod:  &req.RecoveryMethod,
					Priority:        &req.Priority,
					Severity:        &req.Severity,
					Type:            &req.Type,
					Frequency:       &req.Frequency,
					DetectedVersion: &req.DetectedVersion,
					Phase:           &req.Phase,
					RecoveryRank:    &req.RecoveryRank,
					DetectionTeam:   &req.DetectionTeam,
					Location:        &req.Location,
					FixVersion:      &req.FixVersion,
					SQAMemo:         &req.SQAMemo,
					Component:       &req.Component,
					Resolution:      &req.Resolution,
					Models:          &req.Models,
					DetectedBy:      &req.DetectedBy,
					Status:          &req.Status,
				}
				err := s.Update(existingDefect.ID, userID, updateReq)
				if err != nil {
					log.Printf("[Defect Import CSV DEBUG] Row %d (Excel row %d) UPDATE FAILED: Defect ID=%s, error=%v", dataRowNum, excelRowNum, defectID, err)
					continue
				}
				log.Printf("[Defect Import CSV DEBUG] Row %d (Excel row %d) UPDATED: Defect ID=%s", dataRowNum, excelRowNum, defectID)
				result.SuccessCount++
				continue
			}
		}

		// 创建新缺陷
		_, err := s.Create(projectID, userID, req)
		if err != nil {
			log.Printf("[Defect Import CSV DEBUG] Row %d (Excel row %d) CREATE FAILED: %v", dataRowNum, excelRowNum, err)
			continue
		}

		log.Printf("[Defect Import CSV DEBUG] Row %d (Excel row %d) CREATED", dataRowNum, excelRowNum)
		result.SuccessCount++
	}

	log.Printf("[Defect Import CSV] ===== COMPLETE ===== success=%d", result.SuccessCount)
	return result, nil
}

// Export 导出缺陷（默认CSV格式）
func (s *defectService) Export(projectID uint) ([]byte, error) {
	return s.ExportWithFormat(projectID, "csv")
}

// ExportWithFormat 导出缺陷（支持CSV和XLSX格式）
func (s *defectService) ExportWithFormat(projectID uint, format string) ([]byte, error) {
	if format == "xlsx" {
		return s.exportXLSX(projectID)
	}
	return s.exportCSV(projectID)
}

// exportCSV 导出CSV格式
func (s *defectService) exportCSV(projectID uint) ([]byte, error) {
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

	// 写入表头（除了ID之外的全部字段）
	exportHeaders := []string{"Defect ID", "Title", "Module", "Description", "Recovery Method",
		"Priority", "Severity", "Type", "Frequency", "Detected Version", "Phase",
		"Detection Team", "Location", "Fix Version", "SQA MEMO", "Component",
		"Resolution", "Models",
		"Status", "Created At", "Detected By"}
	if err := writer.Write(exportHeaders); err != nil {
		return nil, fmt.Errorf("write csv header: %w", err)
	}

	// 写入数据
	for _, defect := range defects {
		// 格式化创建时间
		createdAt := defect.CreatedAt.Format("2006-01-02 15:04:05")

		// 获取提出人：优先使用DetectedBy字段，否则使用创建人昵称
		detectedBy := defect.DetectedBy
		if detectedBy == "" {
			detectedBy = userMap[defect.CreatedBy]
			if detectedBy == "" {
				detectedBy = fmt.Sprintf("User#%d", defect.CreatedBy)
			}
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
			defect.Type,
			defect.Frequency,
			defect.DetectedVersion,
			defect.Phase,
			defect.DetectionTeam,
			defect.Location,
			defect.FixVersion,
			defect.SQAMemo,
			defect.Component,
			defect.Resolution,
			defect.Models,
			defect.Status,
			createdAt,
			detectedBy,
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

// exportXLSX 导出XLSX格式
func (s *defectService) exportXLSX(projectID uint) ([]byte, error) {
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

	// 创建XLSX文件
	f := excelize.NewFile()
	sheetName := "Defects"
	f.SetSheetName("Sheet1", sheetName)

	// 写入表头
	headers := []string{"Defect ID", "Title", "Module", "Description", "Recovery Method",
		"Priority", "Severity", "Type", "Frequency", "Detected Version", "Phase",
		"Detection Team", "Location", "Fix Version", "SQA MEMO", "Component",
		"Resolution", "Models",
		"Status", "Created At", "Detected By"}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for rowIndex, defect := range defects {
		row := rowIndex + 2 // 从第2行开始（第1行是表头）

		// 格式化创建时间
		createdAt := defect.CreatedAt.Format("2006-01-02 15:04:05")

		// 获取提出人：优先使用DetectedBy字段，如果为空则使用创建人昵称
		detectedBy := defect.DetectedBy
		if detectedBy == "" {
			detectedBy = userMap[defect.CreatedBy]
			if detectedBy == "" {
				detectedBy = fmt.Sprintf("User#%d", defect.CreatedBy)
			}
		}

		data := []interface{}{
			defect.DefectID,
			defect.Title,
			defect.Subject,
			defect.Description,
			defect.RecoveryMethod,
			defect.Priority,
			defect.Severity,
			defect.Type,
			defect.Frequency,
			defect.DetectedVersion,
			defect.Phase,
			defect.DetectionTeam,
			defect.Location,
			defect.FixVersion,
			defect.SQAMemo,
			defect.Component,
			defect.Resolution,
			defect.Models,
			defect.Status,
			createdAt,
			detectedBy,
		}

		for colIndex, value := range data {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 将文件写入内存
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("write xlsx: %w", err)
	}

	return buf.Bytes(), nil
}
