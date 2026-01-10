package services

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"webtest/internal/models"
	"webtest/internal/repositories"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ApiTestCaseService 接口测试用例服务接口
type ApiTestCaseService interface {
	// 用例管理
	GetCases(projectID uint, userID uint, caseType string, page int, size int) ([]*models.ApiTestCase, int64, error)
	GetCasesByGroup(projectID uint, userID uint, caseType string, caseGroup string, page int, size int) ([]*models.ApiTestCase, int64, error)
	CreateCase(projectID uint, userID uint, testCase *models.ApiTestCase) (*models.ApiTestCase, error)
	InsertCase(projectID uint, userID uint, caseType string, position string, targetCaseID string, caseData map[string]interface{}) (*models.ApiTestCase, error)
	DeleteCase(projectID uint, userID uint, caseID string) error
	BatchDeleteCases(projectID uint, userID uint, caseType string, caseIDs []string) (int, []string, error)
	UpdateCase(projectID uint, userID uint, caseID string, updates map[string]interface{}) error

	// 版本管理
	SaveVersion(projectID uint, userID uint, remark string) (*models.ApiTestCaseVersion, error)
	GetVersions(projectID uint, userID uint, page int, size int) ([]*models.ApiTestCaseVersion, int64, error)
	DeleteVersion(projectID uint, userID uint, versionID string) error
	UpdateVersionRemark(projectID uint, userID uint, versionID string, remark string) error

	// 模版与导入导出
	ExportApiTemplate(projectID uint) ([]byte, string, error)
	ImportApiCases(projectID uint, userID uint, caseGroup string, fileData []byte) (insertCount int, updateCount int, err error)
}

type apiTestCaseService struct {
	repo           repositories.ApiTestCaseRepository
	projectService ProjectService
}

// NewApiTestCaseService 创建服务实例
func NewApiTestCaseService(repo repositories.ApiTestCaseRepository, projectService ProjectService) ApiTestCaseService {
	return &apiTestCaseService{
		repo:           repo,
		projectService: projectService,
	}
}

// GetCases 分页查询用例列表
func (s *apiTestCaseService) GetCases(projectID uint, userID uint, caseType string, page int, size int) ([]*models.ApiTestCase, int64, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, 0, errors.New("无项目访问权限")
	}

	// 计算offset
	offset := (page - 1) * size

	// 查询数据(按display_order排序)
	cases, total, err := s.repo.List(projectID, caseType, offset, size)
	if err != nil {
		return nil, 0, fmt.Errorf("list api cases: %w", err)
	}

	return cases, total, nil
}

// GetCasesByGroup 按用例集分页查询用例列表
func (s *apiTestCaseService) GetCasesByGroup(projectID uint, userID uint, caseType string, caseGroup string, page int, size int) ([]*models.ApiTestCase, int64, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, 0, errors.New("无项目访问权限")
	}

	// 计算offset
	offset := (page - 1) * size

	// 如果没有指定用例集，使用原始List方法
	if caseGroup == "" {
		cases, total, err := s.repo.List(projectID, caseType, offset, size)
		if err != nil {
			return nil, 0, fmt.Errorf("list api cases: %w", err)
		}
		return cases, total, nil
	}

	// 按用例集筛选查询数据(按display_order排序)
	cases, total, err := s.repo.ListByGroup(projectID, caseType, caseGroup, offset, size)
	if err != nil {
		return nil, 0, fmt.Errorf("list api cases by group: %w", err)
	}

	return cases, total, nil
}

// CreateCase 创建用例
func (s *apiTestCaseService) CreateCase(projectID uint, userID uint, testCase *models.ApiTestCase) (*models.ApiTestCase, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 字段验证
	if err := testCase.Validate(); err != nil {
		return nil, fmt.Errorf("validate case data: %w", err)
	}

	// 获取当前类型的最大display_order
	cases, _, err := s.repo.List(projectID, testCase.CaseType, 0, 1)
	if err == nil && len(cases) > 0 {
		// 如果有数据,获取最大display_order并+1
		allCases, _ := s.repo.GetByProjectAndType(projectID, testCase.CaseType)
		maxOrder := 0
		for _, c := range allCases {
			if c.DisplayOrder > maxOrder {
				maxOrder = c.DisplayOrder
			}
		}
		testCase.DisplayOrder = maxOrder + 1
	} else {
		// 如果没有数据,设置为1
		testCase.DisplayOrder = 1
	}

	// 保存新用例
	if err := s.repo.Create(testCase); err != nil {
		return nil, fmt.Errorf("create api test case: %w", err)
	}

	return testCase, nil
}

// InsertCase 插入用例(指定位置)
func (s *apiTestCaseService) InsertCase(projectID uint, userID uint, caseType string, position string, targetCaseID string, caseData map[string]interface{}) (*models.ApiTestCase, error) {
	// 如果caseType为空，使用目标用例的caseType或默认值"api"
	if caseType == "" {
		caseType = "api" // 默认值
	}

	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 查询目标用例
	targetCase, err := s.repo.GetByID(targetCaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目标用例不存在")
		}
		return nil, fmt.Errorf("get target case: %w", err)
	}

	// 如果caseType是默认值，优先使用目标用例的caseType
	if caseType == "api" && targetCase.CaseType != "" {
		caseType = targetCase.CaseType
	}

	// 提取case_group（用于筛选同一用例集内的用例）
	caseGroup := ""
	if cg, ok := caseData["case_group"].(string); ok {
		caseGroup = cg
	}
	// 如果没有传入case_group，使用目标用例的case_group
	if caseGroup == "" {
		caseGroup = targetCase.CaseGroup
	}

	log.Printf("[InsertCase] caseType=%s, targetCase.ID=%s, targetCase.DisplayOrder=%d, position=%s, caseGroup=%s",
		caseType, targetCase.ID, targetCase.DisplayOrder, position, caseGroup)

	// 计算新用例的display_order
	var newOrder int
	if position == "before" {
		newOrder = targetCase.DisplayOrder
	} else { // after
		newOrder = targetCase.DisplayOrder + 1
	}

	log.Printf("[InsertCase] newOrder=%d, will IncrementOrderAfter afterOrder=%d", newOrder, targetCase.DisplayOrder-1)

	// 调整现有用例的display_order（只影响同一用例集）
	if position == "before" {
		if err := s.repo.IncrementOrderAfter(projectID, caseType, caseGroup, targetCase.DisplayOrder-1); err != nil {
			return nil, fmt.Errorf("increment order: %w", err)
		}
	} else {
		if err := s.repo.IncrementOrderAfter(projectID, caseType, caseGroup, targetCase.DisplayOrder); err != nil {
			return nil, fmt.Errorf("increment order: %w", err)
		}
	}

	// 创建新用例(UUID自动生成)
	newCase := &models.ApiTestCase{
		ProjectID:    projectID,
		CaseType:     caseType,
		CaseGroup:    caseGroup, // 设置用例集
		DisplayOrder: newOrder,
		TestResult:   "NR",
		Method:       "GET",
	}

	// 填充caseData字段（case_group已在上面设置）
	if caseNumber, ok := caseData["case_number"].(string); ok {
		newCase.CaseNumber = caseNumber
	}
	if screen, ok := caseData["screen"].(string); ok {
		newCase.Screen = screen
	}
	if url, ok := caseData["url"].(string); ok {
		newCase.URL = url
	}
	if header, ok := caseData["header"].(string); ok {
		newCase.Header = header
	}
	if method, ok := caseData["method"].(string); ok {
		newCase.Method = method
	}
	if body, ok := caseData["body"].(string); ok {
		newCase.Body = body
	}
	if response, ok := caseData["response"].(string); ok {
		newCase.Response = response
	}
	if scriptCode, ok := caseData["script_code"].(string); ok {
		newCase.ScriptCode = scriptCode
	}
	if testResult, ok := caseData["test_result"].(string); ok {
		newCase.TestResult = testResult
	}
	if remark, ok := caseData["remark"].(string); ok {
		newCase.Remark = remark
	}

	// 字段验证
	if err := newCase.Validate(); err != nil {
		return nil, fmt.Errorf("validate case data: %w", err)
	}

	// 保存新用例
	if err := s.repo.Create(newCase); err != nil {
		return nil, fmt.Errorf("create new case: %w", err)
	}

	// 重新分配display_order（只影响同一用例集）
	if err := s.repo.ReassignDisplayOrders(projectID, caseType, caseGroup); err != nil {
		return nil, fmt.Errorf("reassign display orders: %w", err)
	}

	return newCase, nil
}

// DeleteCase 删除单个用例
func (s *apiTestCaseService) DeleteCase(projectID uint, userID uint, caseID string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 先检查用例是否存在并属于该项目
	existingCase, err := s.repo.GetByID(caseID)
	if err != nil {
		return errors.New("用例不存在")
	}
	if existingCase.ProjectID != projectID {
		return errors.New("无权限删除该用例")
	}

	// 获取用例类型和用例集用于后续重排序
	caseType := existingCase.CaseType
	caseGroup := existingCase.CaseGroup

	// 删除用例
	if err := s.repo.Delete(caseID); err != nil {
		return fmt.Errorf("delete case: %w", err)
	}

	// 重新分配display_order（只影响同一用例集）
	if err := s.repo.ReassignDisplayOrders(projectID, caseType, caseGroup); err != nil {
		return fmt.Errorf("reassign display orders: %w", err)
	}

	return nil
}

// BatchDeleteCases 批量删除用例
func (s *apiTestCaseService) BatchDeleteCases(projectID uint, userID uint, caseType string, caseIDs []string) (int, []string, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return 0, nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return 0, nil, errors.New("无项目访问权限")
	}

	deletedCount := 0
	var failedCaseIDs []string

	// 循环删除每个用例
	for _, caseID := range caseIDs {
		if err := s.repo.Delete(caseID); err != nil {
			failedCaseIDs = append(failedCaseIDs, caseID)
		} else {
			deletedCount++
		}
	}

	// 重新分配display_order（批量删除时传空字符串重排所有case_type的用例）
	if err := s.repo.ReassignDisplayOrders(projectID, caseType, ""); err != nil {
		return deletedCount, failedCaseIDs, fmt.Errorf("reassign display orders: %w", err)
	}

	return deletedCount, failedCaseIDs, nil
}

// UpdateCase 更新用例
func (s *apiTestCaseService) UpdateCase(projectID uint, userID uint, caseID string, updates map[string]interface{}) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 更新用例
	if err := s.repo.Update(caseID, updates); err != nil {
		return fmt.Errorf("update case: %w", err)
	}

	return nil
}

// ========== 版本管理 ==========

// SaveVersion 保存版本(生成XLSX文件，每个用例集一个sheet)
func (s *apiTestCaseService) SaveVersion(projectID uint, userID uint, remark string) (*models.ApiTestCaseVersion, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, errors.New("无项目访问权限")
	}

	// 获取项目信息(用于文件名)
	project, _, err := s.projectService.GetByID(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	// 生成时间戳
	timestamp := time.Now().Format("20060102_150405")

	// 创建存储目录
	storageDir := filepath.Join("storage", "versions", "api-cases")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}

	// 生成XLSX文件名：项目名_AIAPI_TestCase_时间戳.xlsx
	xlsxFilename := fmt.Sprintf("%s_AIAPI_TestCase_%s.xlsx", project.Name, timestamp)
	xlsxPath := filepath.Join(storageDir, xlsxFilename)

	// 创建Excel工作簿
	f := excelize.NewFile()

	// 创建Cover sheet并设置内容
	coverSheet := "Cover"
	coverIndex, _ := f.NewSheet(coverSheet)
	f.SetCellValue(coverSheet, "A1", "API Test Case Version")
	f.SetCellValue(coverSheet, "A2", fmt.Sprintf("Project: %s", project.Name))
	f.SetCellValue(coverSheet, "A3", fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	if remark != "" {
		f.SetCellValue(coverSheet, "A4", fmt.Sprintf("Remark: %s", remark))
	}

	// 删除默认的Sheet1
	f.DeleteSheet("Sheet1")

	// 设置Cover为默认激活sheet
	f.SetActiveSheet(coverIndex)

	// 获取所有用例集
	caseGroups, err := s.repo.GetCaseGroups(projectID)
	if err != nil {
		return nil, fmt.Errorf("get case groups: %w", err)
	}

	// 为每个用例集创建一个sheet
	for _, groupName := range caseGroups {
		// 获取该用例集的所有用例
		cases, err := s.repo.GetByProjectAndGroup(projectID, groupName)
		if err != nil {
			return nil, fmt.Errorf("get cases for group %s: %w", groupName, err)
		}

		// 跳过空用例集
		if len(cases) == 0 {
			continue
		}

		// 创建sheet（使用用例集名称）
		sheetName := groupName
		f.NewSheet(sheetName)

		// 写入表头
		headers := []string{"No.", "Screen", "URL", "Header", "Method", "Body", "Response", "ScriptCode"}
		for col, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}

		// 写入数据行
		for i, c := range cases {
			row := i + 2                                                     // 从第2行开始（第1行是表头）
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)          // No.
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), c.Screen)     // Screen
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), c.URL)        // URL
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), c.Header)     // Header
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), c.Method)     // Method
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), c.Body)       // Body
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), c.Response)   // Response
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), c.ScriptCode) // ScriptCode
		}

		// 设置列宽
		f.SetColWidth(sheetName, "A", "A", 6)  // No.
		f.SetColWidth(sheetName, "B", "B", 20) // Screen
		f.SetColWidth(sheetName, "C", "C", 40) // URL
		f.SetColWidth(sheetName, "D", "D", 30) // Header
		f.SetColWidth(sheetName, "E", "E", 10) // Method
		f.SetColWidth(sheetName, "F", "F", 40) // Body
		f.SetColWidth(sheetName, "G", "G", 40) // Response
		f.SetColWidth(sheetName, "H", "H", 60) // ScriptCode
	}

	// 保存Excel文件
	if err := f.SaveAs(xlsxPath); err != nil {
		return nil, fmt.Errorf("save xlsx file: %w", err)
	}

	// 创建版本记录(UUID自动生成)
	version := &models.ApiTestCaseVersion{
		ProjectID:    projectID,
		XlsxFilename: xlsxFilename,
		Remark:       remark,
		CreatedBy:    userID,
	}

	if err := s.repo.CreateVersion(version); err != nil {
		return nil, fmt.Errorf("create version: %w", err)
	}

	return version, nil
}

// GetVersions 获取版本列表
func (s *apiTestCaseService) GetVersions(projectID uint, userID uint, page int, size int) ([]*models.ApiTestCaseVersion, int64, error) {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return nil, 0, errors.New("无项目访问权限")
	}

	// 计算offset
	offset := (page - 1) * size

	// 查询版本列表
	versions, total, err := s.repo.ListVersions(projectID, offset, size)
	if err != nil {
		return nil, 0, fmt.Errorf("list versions: %w", err)
	}

	return versions, total, nil
}

// DeleteVersion 删除版本
func (s *apiTestCaseService) DeleteVersion(projectID uint, userID uint, versionID string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 查询版本记录
	version, err := s.repo.GetVersionByID(versionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("版本不存在")
		}
		return fmt.Errorf("get version: %w", err)
	}

	// 验证版本归属
	if version.ProjectID != projectID {
		return errors.New("版本不属于该项目")
	}

	// 删除XLSX文件（优先使用新字段）
	storageDir := filepath.Join("storage", "versions", "api-cases")
	if version.XlsxFilename != "" {
		filePath := filepath.Join(storageDir, version.XlsxFilename)
		os.Remove(filePath) // 忽略删除错误
	} else {
		// 兼容旧版本的CSV文件
		filenames := []string{
			version.FilenameRole1,
			version.FilenameRole2,
			version.FilenameRole3,
			version.FilenameRole4,
		}
		for _, filename := range filenames {
			if filename != "" {
				filePath := filepath.Join(storageDir, filename)
				os.Remove(filePath) // 忽略删除错误
			}
		}
	}

	// 删除版本记录
	if err := s.repo.DeleteVersion(versionID); err != nil {
		return fmt.Errorf("delete version: %w", err)
	}

	return nil
}

// UpdateVersionRemark 更新版本备注
func (s *apiTestCaseService) UpdateVersionRemark(projectID uint, userID uint, versionID string, remark string) error {
	// 权限校验
	isMember, err := s.projectService.IsProjectMember(projectID, userID)
	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}
	if !isMember {
		return errors.New("无项目访问权限")
	}

	// 查询版本记录
	version, err := s.repo.GetVersionByID(versionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("版本不存在")
		}
		return fmt.Errorf("get version: %w", err)
	}

	// 验证版本归属
	if version.ProjectID != projectID {
		return errors.New("版本不属于该项目")
	}

	// 更新备注(限制500字符由前端校验)
	if err := s.repo.UpdateVersionRemark(versionID, remark); err != nil {
		return fmt.Errorf("update version remark: %w", err)
	}

	return nil
}

// ExportApiTemplate 导出API用例模版
func (s *apiTestCaseService) ExportApiTemplate(projectID uint) ([]byte, string, error) {
	// 使用excelize创建XLSX工作簿
	f := excelize.NewFile()
	sheetName := "API Cases"
	index, _ := f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	// 设置表头（9列）
	headers := []string{"No.", "Screen", "URL", "Header", "Method", "Body", "Response", "ScriptCode", "UUID"}
	for i, h := range headers {
		cell := fmt.Sprintf("%s1", columnName(i))
		f.SetCellValue(sheetName, cell, h)
	}

	// 生成文件名（ISO时间戳格式）
	timestamp := time.Now().Format("20060102T150405")
	filename := fmt.Sprintf("API_Case_Template_%s.xlsx", timestamp)

	// 保存到buffer
	var buffer bytes.Buffer
	if err := f.Write(&buffer); err != nil {
		return nil, "", fmt.Errorf("write xlsx: %w", err)
	}

	return buffer.Bytes(), filename, nil
}

// ImportApiCases 导入API用例
func (s *apiTestCaseService) ImportApiCases(projectID uint, userID uint, caseGroup string, fileData []byte) (insertCount int, updateCount int, err error) {
	// 使用excelize解析XLSX
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return 0, 0, fmt.Errorf("open xlsx: %w", err)
	}
	defer f.Close()

	// 获取第一个sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return 0, 0, errors.New("文件中没有工作表")
	}
	sheetName := sheets[0]

	// 读取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0, 0, fmt.Errorf("get rows: %w", err)
	}

	if len(rows) <= 1 {
		return 0, 0, errors.New("文件中没有数据")
	}

	// 解析表头，识别字段对应的列索引
	headerRow := rows[0]
	colIndex := make(map[string]int)

	// 字段名称映射（支持中英文和多种写法）
	fieldMappings := map[string][]string{
		"no":       {"no", "no.", "序号", "编号"},
		"screen":   {"screen", "画面", "模块", "功能"},
		"url":      {"url", "接口", "地址", "接口地址", "api"},
		"header":   {"header", "headers", "请求头", "头部"},
		"method":   {"method", "方法", "请求方法", "http方法"},
		"body":     {"body", "请求体", "请求内容", "参数"},
		"response": {"response", "响应", "预期响应", "期望结果", "返回"},
		"uuid":     {"uuid", "id", "用例id", "唯一标识"},
	}

	// 遍历表头，匹配字段
	for colIdx, header := range headerRow {
		headerLower := strings.ToLower(strings.TrimSpace(header))
		for field, aliases := range fieldMappings {
			for _, alias := range aliases {
				if headerLower == alias || strings.Contains(headerLower, alias) {
					if _, exists := colIndex[field]; !exists {
						colIndex[field] = colIdx
					}
					break
				}
			}
		}
	}

	// 辅助函数：安全获取列值
	getColValue := func(row []string, field string) string {
		if idx, exists := colIndex[field]; exists && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// 跳过表头，处理数据行
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue // 跳过空行
		}

		// 根据表头识别的列索引获取数据
		uuid := getColValue(row, "uuid")
		screen := getColValue(row, "screen")
		url := getColValue(row, "url")
		header := getColValue(row, "header")
		method := getColValue(row, "method")
		body := getColValue(row, "body")
		response := getColValue(row, "response")
		remark := "" // Remark列已删除，默认为空

		// 如果所有关键字段都为空，跳过该行
		if screen == "" && url == "" && method == "" && body == "" && response == "" {
			continue
		}

		// 根据UUID判断是UPDATE还是INSERT
		if uuid != "" {
			// UPDATE：根据UUID查找用例
			existingCase, err := s.repo.GetByID(uuid)
			if err == nil && existingCase != nil && existingCase.ProjectID == projectID {
				// 更新用例
				updates := map[string]interface{}{
					"screen":     screen,
					"url":        url,
					"header":     header,
					"method":     method,
					"body":       body,
					"response":   response,
					"remark":     remark,
					"case_group": caseGroup,
				}
				if err := s.repo.Update(uuid, updates); err != nil {
					return insertCount, updateCount, fmt.Errorf("update case: %w", err)
				}
				updateCount++
			} else {
				// UUID不存在或不属于该项目，视为INSERT
				newCase := &models.ApiTestCase{
					ProjectID: projectID,
					CaseType:  "api", // API用例默认类型
					CaseGroup: caseGroup,
					Screen:    screen,
					URL:       url,
					Header:    header,
					Method:    method,
					Body:      body,
					Response:  response,
					Remark:    remark,
				}
				if err := s.repo.Create(newCase); err != nil {
					return insertCount, updateCount, fmt.Errorf("create case: %w", err)
				}
				insertCount++
			}
		} else {
			// INSERT：没有UUID，创建新用例
			newCase := &models.ApiTestCase{
				ProjectID: projectID,
				CaseType:  "api", // API用例默认类型
				CaseGroup: caseGroup,
				Screen:    screen,
				URL:       url,
				Header:    header,
				Method:    method,
				Body:      body,
				Response:  response,
				Remark:    remark,
			}
			if err := s.repo.Create(newCase); err != nil {
				return insertCount, updateCount, fmt.Errorf("create case: %w", err)
			}
			insertCount++
		}
	}

	return insertCount, updateCount, nil
}
