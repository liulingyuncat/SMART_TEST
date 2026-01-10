package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// ExcelService Excel导入导出服务接口
type ExcelService interface {
	ExportAICases(projectID uint) ([]byte, string, error)
	ExportTemplate(projectID uint, caseType string) ([]byte, string, error)
	// T44: 扩展ExportCases支持language和caseGroup参数
	ExportCases(projectID uint, caseType string, taskUUID string, language string, caseGroup string) ([]byte, string, error)
	// T44: 扩展ImportCases支持language参数
	ImportCases(projectID uint, caseType string, fileData []byte, language string, caseGroup string) (updateCount, insertCount int, err error)
	// T45: Web用例多语言导出
	ExportWebCasesByLanguage(projectName string, caseGroups []models.CaseGroup, cases []models.AutoTestCase, language string) ([]byte, string, error)
	GenerateWebCasesZip(projectID uint, projectName string, cases []models.AutoTestCase) (zipPath string, fileSize int64, err error)
}

type excelService struct {
	caseRepo    repositories.ManualTestCaseRepository
	projectRepo repositories.ProjectRepository
	ecrRepo     repositories.ExecutionCaseResultRepository
}

// NewExcelService 创建Excel服务实例
func NewExcelService(
	caseRepo repositories.ManualTestCaseRepository,
	projectRepo repositories.ProjectRepository,
	ecrRepo repositories.ExecutionCaseResultRepository,
) ExcelService {
	return &excelService{
		caseRepo:    caseRepo,
		projectRepo: projectRepo,
		ecrRepo:     ecrRepo,
	}
}

// containsString 检查字符串切片中是否包含指定字符串
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ExportAICases 导出AI用例(9列单Sheet)
func (s *excelService) ExportAICases(projectID uint) ([]byte, string, error) {
	// 1. 查询用例数据
	cases, err := s.caseRepo.GetByProjectAndTypeOrdered(projectID, "ai")
	if err != nil {
		return nil, "", fmt.Errorf("get cases: %w", err)
	}

	// 2. 创建Excel文件
	f := excelize.NewFile()
	sheetName := "用例数据"
	index, _ := f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	// 3. 写入标题行
	headers := []string{"No.", "CaseID", "Maj.Category", "Mid.Category", "Min.Category",
		"Precondition", "Test Step", "Expect", "Remark"}
	for i, h := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, h)
	}

	// 4. 写入数据行
	for i, c := range cases {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), c.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), c.CaseNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), c.MajorFunction)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), c.MiddleFunction)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), c.MinorFunction)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), c.Precondition)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), c.TestSteps)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), c.ExpectedResult)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), c.Remark)
	}

	// 5. 生成文件名和字节流
	project, _ := s.projectRepo.GetByID(projectID)
	projectName := "project"
	if project != nil {
		projectName = project.Name
	}
	filename := fmt.Sprintf("%s_AI_Cases_%s.xlsx", projectName, time.Now().Format("2006-01-02_150405"))

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", fmt.Errorf("write buffer: %w", err)
	}

	return buffer.Bytes(), filename, nil
}

// ExportTemplate 导出模板(23列空模板+示例行)
func (s *excelService) ExportTemplate(projectID uint, caseType string) ([]byte, string, error) {
	f := excelize.NewFile()
	sheetName := "用例数据"
	index, _ := f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	// 写入23列标题
	headers := []string{"No.", "CaseID",
		"Maj.CategoryCN", "Maj.CategoryJP", "Maj.CategoryEN",
		"Mid.CategoryCN", "Mid.CategoryJP", "Mid.CategoryEN",
		"Min.CategoryCN", "Min.CategoryJP", "Min.CategoryEN",
		"PreconditionCN", "PreconditionJP", "PreconditionEN",
		"Test StepCN", "Test StepJP", "Test StepEN",
		"ExpectCN", "ExpectJP", "ExpectEN",
		"TestResult", "Remark", "UUID"}

	for i, h := range headers {
		cell := fmt.Sprintf("%s1", columnName(i))
		f.SetCellValue(sheetName, cell, h)
	}

	// 写入示例数据行（完整示例）
	f.SetCellValue(sheetName, "A2", "1")
	f.SetCellValue(sheetName, "B2", "TC001")
	f.SetCellValue(sheetName, "C2", "登录功能")
	f.SetCellValue(sheetName, "D2", "ログイン機能")
	f.SetCellValue(sheetName, "E2", "Login Function")
	f.SetCellValue(sheetName, "F2", "用户登录")
	f.SetCellValue(sheetName, "G2", "ユーザーログイン")
	f.SetCellValue(sheetName, "H2", "User Login")
	f.SetCellValue(sheetName, "I2", "登录界面")
	f.SetCellValue(sheetName, "J2", "ログイン画面")
	f.SetCellValue(sheetName, "K2", "Login Page")
	f.SetCellValue(sheetName, "L2", "用户已注册")
	f.SetCellValue(sheetName, "M2", "ユーザー登録済み")
	f.SetCellValue(sheetName, "N2", "User registered")
	f.SetCellValue(sheetName, "O2", "1. 打开登录页面\n2. 输入用户名和密码\n3. 点击登录按钮")
	f.SetCellValue(sheetName, "P2", "1. ログインページを開く\n2. ユーザー名とパスワードを入力\n3. ログインボタンをクリック")
	f.SetCellValue(sheetName, "Q2", "1. Open login page\n2. Enter username and password\n3. Click login button")
	f.SetCellValue(sheetName, "R2", "成功登录并跳转到主页")
	f.SetCellValue(sheetName, "S2", "正常にログインしてホームページに移動")
	f.SetCellValue(sheetName, "T2", "Login successfully and redirect to homepage")
	f.SetCellValue(sheetName, "U2", "NR")
	f.SetCellValue(sheetName, "V2", "示例备注")

	project, _ := s.projectRepo.GetByID(projectID)
	projectName := "project"
	if project != nil {
		projectName = project.Name
	}

	typeMap := map[string]string{"overall": "Overall", "change": "Change", "acceptance": "Acceptance"}
	filename := fmt.Sprintf("%s_%s_Template.xlsx", projectName, typeMap[caseType])

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", fmt.Errorf("write buffer: %w", err)
	}

	return buffer.Bytes(), filename, nil
}

// ExportCases 导出整体/变更用例(双Sheet: 元数据+23列用例)
// taskUUID为空时导出用例库数据,非空时合并执行结果(test_result/bug_id/remark)
// T44: 扩展支持language和caseGroup参数
// language: CN/JP/EN (为空则导出全部语言)
// caseGroup: 用例集名称 (为空则导出全部)
func (s *excelService) ExportCases(projectID uint, caseType string, taskUUID string, language string, caseGroup string) ([]byte, string, error) {
	// 1. 查询元数据和用例数据
	metadata, _ := s.caseRepo.GetMetadataByProjectID(projectID, caseType)
	cases, err := s.caseRepo.GetByProjectAndTypeOrdered(projectID, caseType)
	if err != nil {
		return nil, "", fmt.Errorf("get cases: %w", err)
	}

	// T44: 根据caseGroup参数过滤用例
	if caseGroup != "" {
		var filteredCases []*models.ManualTestCase
		for _, c := range cases {
			// 根据case_group字段匹配用例集
			if c.CaseGroup == caseGroup {
				filteredCases = append(filteredCases, c)
			}
		}
		cases = filteredCases
	}

	// 2. 如果taskUUID非空,合并执行结果
	var executionResults map[string]*models.ExecutionCaseResult
	if taskUUID != "" {
		results, err := s.ecrRepo.GetByTaskUUID(taskUUID)
		if err == nil && len(results) > 0 {
			executionResults = make(map[string]*models.ExecutionCaseResult)
			for _, r := range results {
				executionResults[r.CaseID] = r
			}
		}
	}

	// 3. 创建Excel文件
	f := excelize.NewFile()

	// Sheet1: 元数据
	metaSheet := "元数据"
	f.SetSheetName("Sheet1", metaSheet)
	if metadata != nil {
		f.SetCellValue(metaSheet, "A1", "Test Version")
		f.SetCellValue(metaSheet, "B1", metadata.TestVersion)
		f.SetCellValue(metaSheet, "A2", "Test Environment")
		f.SetCellValue(metaSheet, "B2", metadata.TestEnv)
		f.SetCellValue(metaSheet, "A3", "Test Date")
		f.SetCellValue(metaSheet, "B3", metadata.TestDate)
		f.SetCellValue(metaSheet, "A4", "Tester")
		f.SetCellValue(metaSheet, "B4", metadata.Executor)
	}

	// 4. Sheet2: 用例数据
	// T44: 如果指定caseGroup,Sheet名使用用例集名称
	dataSheet := "用例数据"
	if caseGroup != "" {
		dataSheet = caseGroup
	}
	index, _ := f.NewSheet(dataSheet)
	f.SetActiveSheet(index)

	// T44: 根据language参数动态生成列头（与前端表头保持一致）
	var headers []string
	if language != "" {
		// 按语言导出(8列格式)
		switch language {
		case "CN":
			headers = []string{"UUID", "CaseID", "Maj.CategoryCN", "Mid.CategoryCN", "Min.CategoryCN", "PreconditionCN", "Test StepCN", "ExpectCN"}
		case "JP":
			headers = []string{"UUID", "CaseID", "Maj.CategoryJP", "Mid.CategoryJP", "Min.CategoryJP", "PreconditionJP", "Test StepJP", "ExpectJP"}
		case "EN":
			headers = []string{"UUID", "CaseID", "Maj.CategoryEN", "Mid.CategoryEN", "Min.CategoryEN", "PreconditionEN", "Test StepEN", "ExpectEN"}
		default:
			headers = []string{"UUID", "CaseID", "Maj.CategoryEN", "Mid.CategoryEN", "Min.CategoryEN", "PreconditionEN", "Test StepEN", "ExpectEN"}
		}
	} else {
		// 全语言导出(原23/25列格式)
		if executionResults != nil {
			headers = []string{"No.", "CaseID",
				"Maj.CategoryCN", "Maj.CategoryJP", "Maj.CategoryEN",
				"Mid.CategoryCN", "Mid.CategoryJP", "Mid.CategoryEN",
				"Min.CategoryCN", "Min.CategoryJP", "Min.CategoryEN",
				"PreconditionCN", "PreconditionJP", "PreconditionEN",
				"Test StepCN", "Test StepJP", "Test StepEN",
				"ExpectCN", "ExpectJP", "ExpectEN",
				"TestResult", "BugID", "ExecutionRemark", "Remark", "UUID"}
		} else {
			headers = []string{"No.", "CaseID",
				"Maj.CategoryCN", "Maj.CategoryJP", "Maj.CategoryEN",
				"Mid.CategoryCN", "Mid.CategoryJP", "Mid.CategoryEN",
				"Min.CategoryCN", "Min.CategoryJP", "Min.CategoryEN",
				"PreconditionCN", "PreconditionJP", "PreconditionEN",
				"Test StepCN", "Test StepJP", "Test StepEN",
				"ExpectCN", "ExpectJP", "ExpectEN",
				"TestResult", "Remark", "UUID"}
		}
	}

	for i, h := range headers {
		cell := fmt.Sprintf("%s1", columnName(i))
		f.SetCellValue(dataSheet, cell, h)
	}

	// 5. 写入数据行
	for i, c := range cases {
		row := i + 2

		// T44: 根据language参数选择导出字段
		if language != "" {
			// 按单语言导出(8列格式: UUID, 用例编号, 一级功能, 二级功能, 三级功能, 前置条件, 测试步骤, 期望结果)
			f.SetCellValue(dataSheet, fmt.Sprintf("A%d", row), c.CaseID)
			f.SetCellValue(dataSheet, fmt.Sprintf("B%d", row), c.CaseNumber)
			switch language {
			case "CN":
				f.SetCellValue(dataSheet, fmt.Sprintf("C%d", row), c.MajorFunctionCN)
				f.SetCellValue(dataSheet, fmt.Sprintf("D%d", row), c.MiddleFunctionCN)
				f.SetCellValue(dataSheet, fmt.Sprintf("E%d", row), c.MinorFunctionCN)
				f.SetCellValue(dataSheet, fmt.Sprintf("F%d", row), c.PreconditionCN)
				f.SetCellValue(dataSheet, fmt.Sprintf("G%d", row), c.TestStepsCN)
				f.SetCellValue(dataSheet, fmt.Sprintf("H%d", row), c.ExpectedResultCN)
			case "JP":
				f.SetCellValue(dataSheet, fmt.Sprintf("C%d", row), c.MajorFunctionJP)
				f.SetCellValue(dataSheet, fmt.Sprintf("D%d", row), c.MiddleFunctionJP)
				f.SetCellValue(dataSheet, fmt.Sprintf("E%d", row), c.MinorFunctionJP)
				f.SetCellValue(dataSheet, fmt.Sprintf("F%d", row), c.PreconditionJP)
				f.SetCellValue(dataSheet, fmt.Sprintf("G%d", row), c.TestStepsJP)
				f.SetCellValue(dataSheet, fmt.Sprintf("H%d", row), c.ExpectedResultJP)
			case "EN":
				f.SetCellValue(dataSheet, fmt.Sprintf("C%d", row), c.MajorFunctionEN)
				f.SetCellValue(dataSheet, fmt.Sprintf("D%d", row), c.MiddleFunctionEN)
				f.SetCellValue(dataSheet, fmt.Sprintf("E%d", row), c.MinorFunctionEN)
				f.SetCellValue(dataSheet, fmt.Sprintf("F%d", row), c.PreconditionEN)
				f.SetCellValue(dataSheet, fmt.Sprintf("G%d", row), c.TestStepsEN)
				f.SetCellValue(dataSheet, fmt.Sprintf("H%d", row), c.ExpectedResultEN)
			}
			continue
		}

		// 全语言导出(原格式)
		f.SetCellValue(dataSheet, fmt.Sprintf("A%d", row), c.ID)
		f.SetCellValue(dataSheet, fmt.Sprintf("B%d", row), c.CaseNumber)
		f.SetCellValue(dataSheet, fmt.Sprintf("C%d", row), c.MajorFunctionCN)
		f.SetCellValue(dataSheet, fmt.Sprintf("D%d", row), c.MajorFunctionJP)
		f.SetCellValue(dataSheet, fmt.Sprintf("E%d", row), c.MajorFunctionEN)
		f.SetCellValue(dataSheet, fmt.Sprintf("F%d", row), c.MiddleFunctionCN)
		f.SetCellValue(dataSheet, fmt.Sprintf("G%d", row), c.MiddleFunctionJP)
		f.SetCellValue(dataSheet, fmt.Sprintf("H%d", row), c.MiddleFunctionEN)
		f.SetCellValue(dataSheet, fmt.Sprintf("I%d", row), c.MinorFunctionCN)
		f.SetCellValue(dataSheet, fmt.Sprintf("J%d", row), c.MinorFunctionJP)
		f.SetCellValue(dataSheet, fmt.Sprintf("K%d", row), c.MinorFunctionEN)
		f.SetCellValue(dataSheet, fmt.Sprintf("L%d", row), c.PreconditionCN)
		f.SetCellValue(dataSheet, fmt.Sprintf("M%d", row), c.PreconditionJP)
		f.SetCellValue(dataSheet, fmt.Sprintf("N%d", row), c.PreconditionEN)
		f.SetCellValue(dataSheet, fmt.Sprintf("O%d", row), c.TestStepsCN)
		f.SetCellValue(dataSheet, fmt.Sprintf("P%d", row), c.TestStepsJP)
		f.SetCellValue(dataSheet, fmt.Sprintf("Q%d", row), c.TestStepsEN)
		f.SetCellValue(dataSheet, fmt.Sprintf("R%d", row), c.ExpectedResultCN)
		f.SetCellValue(dataSheet, fmt.Sprintf("S%d", row), c.ExpectedResultJP)
		f.SetCellValue(dataSheet, fmt.Sprintf("T%d", row), c.ExpectedResultEN)

		// 合并执行结果(如果有)
		if executionResults != nil {
			if execResult, ok := executionResults[c.CaseID]; ok {
				f.SetCellValue(dataSheet, fmt.Sprintf("U%d", row), execResult.TestResult)
				f.SetCellValue(dataSheet, fmt.Sprintf("V%d", row), execResult.BugID)
				f.SetCellValue(dataSheet, fmt.Sprintf("W%d", row), execResult.Remark)
			} else {
				f.SetCellValue(dataSheet, fmt.Sprintf("U%d", row), c.TestResult)
				f.SetCellValue(dataSheet, fmt.Sprintf("V%d", row), "")
				f.SetCellValue(dataSheet, fmt.Sprintf("W%d", row), "")
			}
			f.SetCellValue(dataSheet, fmt.Sprintf("X%d", row), c.Remark)
			f.SetCellValue(dataSheet, fmt.Sprintf("Y%d", row), c.CaseID)
		} else {
			f.SetCellValue(dataSheet, fmt.Sprintf("U%d", row), c.TestResult)
			f.SetCellValue(dataSheet, fmt.Sprintf("V%d", row), c.Remark)
			f.SetCellValue(dataSheet, fmt.Sprintf("W%d", row), c.CaseID)
		}
	}

	// 生成文件名
	project, _ := s.projectRepo.GetByID(projectID)
	projectName := "project"
	if project != nil {
		projectName = project.Name
	}

	// T44: 文件名格式调整为: 项目名_Manual_用例集名_语言_时间戳.xlsx
	typeMap := map[string]string{"overall": "Overall", "change": "Change", "acceptance": "Acceptance"}
	var filename string
	if language != "" && caseGroup != "" {
		// 按语言导出时使用新格式
		filename = fmt.Sprintf("%s_Manual_%s_%s_%s.xlsx", projectName, caseGroup, language, time.Now().Format("20060102_150405"))
	} else {
		// 兼容旧格式
		filename = fmt.Sprintf("%s_%s_Cases_%s.xlsx", projectName, typeMap[caseType], time.Now().Format("2006-01-02_150405"))
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", fmt.Errorf("write buffer: %w", err)
	}

	return buffer.Bytes(), filename, nil
}

// ImportCases 导入用例(UUID匹配逻辑)
// T44: 扩展支持language和caseGroup参数
// language: CN/JP/EN (为空则更新全部语言字段)
// caseGroup: 用例集名称 (为空则不过滤)
func (s *excelService) ImportCases(projectID uint, caseType string, fileData []byte, language string, caseGroup string) (int, int, error) {
	fmt.Printf("\n🔍 [ImportCases] 开始导入:\n")
	fmt.Printf("  ProjectID: %d\n", projectID)
	fmt.Printf("  CaseType: %q\n", caseType)
	fmt.Printf("  Language: %q\n", language)
	fmt.Printf("  CaseGroup: %q (长度: %d)\n", caseGroup, len(caseGroup))
	fmt.Printf("  文件大小: %d bytes\n", len(fileData))

	if caseGroup == "" {
		fmt.Println("❌ [ImportCases] 严重错误: caseGroup参数为空！")
	} else {
		fmt.Printf("✅ [ImportCases] caseGroup已接收: %q\n", caseGroup)
	}

	// 1. 打开Excel文件
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return 0, 0, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	// 2. 读取Sheet数据
	// T44: 智能选择Sheet - 优先用例数据，否则选第一个非空Sheet
	sheetList := f.GetSheetList()
	fmt.Printf("📊 [ImportCases] Excel包含的Sheets: %v\n", sheetList)

	var sheetName string
	// 优先尝试 "用例数据"
	if containsString(sheetList, "用例数据") {
		sheetName = "用例数据"
	} else if len(sheetList) > 0 {
		// 使用第一个Sheet
		sheetName = sheetList[0]
		fmt.Printf("⚠️ [ImportCases] 未找到'用例数据'工作表，使用第一个Sheet: %s\n", sheetName)
	} else {
		return 0, 0, fmt.Errorf("Excel文件中没有可用的工作表")
	}

	fmt.Printf("✅ [ImportCases] 读取Sheet: %s\n", sheetName)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0, 0, fmt.Errorf("读取Sheet '%s' 失败: %w", sheetName, err)
	}

	if len(rows) < 2 {
		return 0, 0, fmt.Errorf("no data rows")
	}

	// 检测Excel格式和语言：根据表头判断
	headerRow := rows[0]
	is9ColumnFormat := false
	detectedLanguage := language // 默认使用传入的language参数

	if len(headerRow) >= 8 {
		// 检查是否是9列格式
		if len(headerRow) < 12 { // 少于12列，肯定不是23列格式
			is9ColumnFormat = true
			fmt.Printf("📝 [ImportCases] 检测到9列单语言格式\n")

			// 自动检测语言：根据表头第3列的语言后缀
			if len(headerRow) > 2 {
				header := headerRow[2] // Maj.CategoryXX
				if strings.HasSuffix(header, "CN") {
					detectedLanguage = "CN"
				} else if strings.HasSuffix(header, "JP") {
					detectedLanguage = "JP"
				} else if strings.HasSuffix(header, "EN") {
					detectedLanguage = "EN"
				}
				if detectedLanguage != language {
					fmt.Printf("⚠️ [ImportCases] 语言覆盖: 传入=%s, 检测到=%s (使用检测值)\n", language, detectedLanguage)
					language = detectedLanguage // 使用检测到的语言
				}
			}
		}
	}
	fmt.Printf("📋 [ImportCases] 表头: %v\n", headerRow)
	fmt.Printf("🔍 [ImportCases] 格式判断: is9ColumnFormat=%v, 列数=%d, 最终语言=%s\n", is9ColumnFormat, len(headerRow), language)

	updateCount := 0
	insertCount := 0

	// 获取当前最大ID（在循环外获取一次，避免并发问题）
	currentMaxID, err := s.caseRepo.GetMaxID(projectID, caseType)
	if err != nil {
		return 0, 0, fmt.Errorf("get max id: %w", err)
	}

	// 3. 遍历数据行(跳过标题行)
	for _, row := range rows[1:] {
		// 安全读取列数据的辅助函数
		getCol := func(index int) string {
			if index < len(row) {
				return strings.TrimSpace(row[index])
			}
			return ""
		}

		// T44: 根据language参数读取不同列格式
		var caseNumber, uuidStr string
		var majorFunctionCN, majorFunctionJP, majorFunctionEN string
		var middleFunctionCN, middleFunctionJP, middleFunctionEN string
		var minorFunctionCN, minorFunctionJP, minorFunctionEN string
		var preconditionCN, preconditionJP, preconditionEN string
		var testStepsCN, testStepsJP, testStepsEN string
		var expectedResultCN, expectedResultJP, expectedResultEN string
		var testResult, remark string

		// 根据检测到的格式读取数据
		if is9ColumnFormat && language != "" {
			// 9列单语言格式: No., CaseID, Maj.Category, Mid.Category, Min.Category, Precondition, TestStep, Expect, UUID
			caseNumber = getCol(1)
			uuidStr = getCol(8)
			switch language {
			case "CN":
				majorFunctionCN = getCol(2)
				middleFunctionCN = getCol(3)
				minorFunctionCN = getCol(4)
				preconditionCN = getCol(5)
				testStepsCN = getCol(6)
				expectedResultCN = getCol(7)
			case "JP":
				majorFunctionJP = getCol(2)
				middleFunctionJP = getCol(3)
				minorFunctionJP = getCol(4)
				preconditionJP = getCol(5)
				testStepsJP = getCol(6)
				expectedResultJP = getCol(7)
			case "EN":
				majorFunctionEN = getCol(2)
				middleFunctionEN = getCol(3)
				minorFunctionEN = getCol(4)
				preconditionEN = getCol(5)
				testStepsEN = getCol(6)
				expectedResultEN = getCol(7)
			}
		} else {
			// 23列格式(原格式)
			caseNumber = getCol(1)
			uuidStr = getCol(22)
			majorFunctionCN = getCol(2)
			majorFunctionJP = getCol(3)
			majorFunctionEN = getCol(4)
			middleFunctionCN = getCol(5)
			middleFunctionJP = getCol(6)
			middleFunctionEN = getCol(7)
			minorFunctionCN = getCol(8)
			minorFunctionJP = getCol(9)
			minorFunctionEN = getCol(10)
			preconditionCN = getCol(11)
			preconditionJP = getCol(12)
			preconditionEN = getCol(13)
			testStepsCN = getCol(14)
			testStepsJP = getCol(15)
			testStepsEN = getCol(16)
			expectedResultCN = getCol(17)
			expectedResultJP = getCol(18)
			expectedResultEN = getCol(19)
			testResult = getCol(20)
			remark = getCol(21)
		}

		// 检查是否为完全空行：所有字段都为空才跳过
		hasData := caseNumber != "" ||
			majorFunctionCN != "" || majorFunctionJP != "" || majorFunctionEN != "" ||
			middleFunctionCN != "" || middleFunctionJP != "" || middleFunctionEN != "" ||
			minorFunctionCN != "" || minorFunctionJP != "" || minorFunctionEN != "" ||
			preconditionCN != "" || preconditionJP != "" || preconditionEN != "" ||
			testStepsCN != "" || testStepsJP != "" || testStepsEN != "" ||
			expectedResultCN != "" || expectedResultJP != "" || expectedResultEN != "" ||
			testResult != "" || remark != ""

		if !hasData {
			continue
		}

		// 解析数据
		fmt.Printf("\n📝 [ImportCases] 创建用例对象:\n")
		fmt.Printf("  将设置 CaseGroup = %q\n", caseGroup)

		testCase := &models.ManualTestCase{
			ProjectID:        projectID,
			CaseType:         caseType,
			CaseNumber:       caseNumber,
			CaseGroup:        caseGroup, // T44: 设置用例集字段
			MajorFunctionCN:  majorFunctionCN,
			MajorFunctionJP:  majorFunctionJP,
			MajorFunctionEN:  majorFunctionEN,
			MiddleFunctionCN: middleFunctionCN,
			MiddleFunctionJP: middleFunctionJP,
			MiddleFunctionEN: middleFunctionEN,
			MinorFunctionCN:  minorFunctionCN,
			MinorFunctionJP:  minorFunctionJP,
			MinorFunctionEN:  minorFunctionEN,
			PreconditionCN:   preconditionCN,
			PreconditionJP:   preconditionJP,
			PreconditionEN:   preconditionEN,
			TestStepsCN:      testStepsCN,
			TestStepsJP:      testStepsJP,
			TestStepsEN:      testStepsEN,
			ExpectedResultCN: expectedResultCN,
			ExpectedResultJP: expectedResultJP,
			ExpectedResultEN: expectedResultEN,
			TestResult:       testResult,
			Remark:           remark,
		}

		// 调试：打印导入的数据
		fmt.Printf("=== Importing Row ===\n")
		fmt.Printf("CaseType: %s, CaseNumber: %q\n", caseType, caseNumber)
		fmt.Printf("MajorCN=%q, MajorJP=%q, MajorEN=%q\n", majorFunctionCN, majorFunctionJP, majorFunctionEN)
		fmt.Printf("MiddleCN=%q, MiddleJP=%q, MiddleEN=%q\n", middleFunctionCN, middleFunctionJP, middleFunctionEN)
		fmt.Printf("MinorCN=%q, MinorJP=%q, MinorEN=%q\n", minorFunctionCN, minorFunctionJP, minorFunctionEN)
		fmt.Printf("PrecondCN=%q, PrecondJP=%q, PrecondEN=%q\n", preconditionCN, preconditionJP, preconditionEN)
		fmt.Printf("TestStepsCN=%q\nTestStepsJP=%q\nTestStepsEN=%q\n", testStepsCN, testStepsJP, testStepsEN)
		fmt.Printf("ExpectCN=%q\nExpectJP=%q\nExpectEN=%q\n", expectedResultCN, expectedResultJP, expectedResultEN)
		fmt.Printf("TestResult=%q, Remark=%q, UUID=%q\n", testResult, remark, uuidStr)
		fmt.Printf("==================\n")

		// UUID匹配逻辑
		if uuidStr != "" && uuidStr != " " {
			// 非空UUID: 尝试查找
			existing, _ := s.caseRepo.GetByCaseID(uuidStr)
			if existing != nil {
				// UUID存在于数据库: 覆盖更新(保留created_at和ID)
				// T44: 根据language参数仅更新对应语言字段
				fmt.Printf("\n🔄 [ImportCases] 准备更新已有用例 (UUID: %s):\n", uuidStr)
				fmt.Printf("  旧CaseGroup: %q\n", existing.CaseGroup)
				fmt.Printf("  新CaseGroup: %q\n", testCase.CaseGroup)

				updates := map[string]interface{}{
					"case_number": testCase.CaseNumber,
					"case_group":  testCase.CaseGroup, // T44: 更新用例集字段
					"test_result": testCase.TestResult,
					"remark":      testCase.Remark,
				}

				fmt.Printf("  Updates map: %+v\n", updates)

				if language == "" {
					// 全语言更新
					updates["major_function_cn"] = testCase.MajorFunctionCN
					updates["major_function_jp"] = testCase.MajorFunctionJP
					updates["major_function_en"] = testCase.MajorFunctionEN
					updates["middle_function_cn"] = testCase.MiddleFunctionCN
					updates["middle_function_jp"] = testCase.MiddleFunctionJP
					updates["middle_function_en"] = testCase.MiddleFunctionEN
					updates["minor_function_cn"] = testCase.MinorFunctionCN
					updates["minor_function_jp"] = testCase.MinorFunctionJP
					updates["minor_function_en"] = testCase.MinorFunctionEN
					updates["precondition_cn"] = testCase.PreconditionCN
					updates["precondition_jp"] = testCase.PreconditionJP
					updates["precondition_en"] = testCase.PreconditionEN
					updates["test_steps_cn"] = testCase.TestStepsCN
					updates["test_steps_jp"] = testCase.TestStepsJP
					updates["test_steps_en"] = testCase.TestStepsEN
					updates["expected_result_cn"] = testCase.ExpectedResultCN
					updates["expected_result_jp"] = testCase.ExpectedResultJP
					updates["expected_result_en"] = testCase.ExpectedResultEN
				} else if language == "CN" {
					updates["major_function_cn"] = testCase.MajorFunctionCN
					updates["middle_function_cn"] = testCase.MiddleFunctionCN
					updates["minor_function_cn"] = testCase.MinorFunctionCN
					updates["precondition_cn"] = testCase.PreconditionCN
					updates["test_steps_cn"] = testCase.TestStepsCN
					updates["expected_result_cn"] = testCase.ExpectedResultCN
				} else if language == "JP" {
					updates["major_function_jp"] = testCase.MajorFunctionJP
					updates["middle_function_jp"] = testCase.MiddleFunctionJP
					updates["minor_function_jp"] = testCase.MinorFunctionJP
					updates["precondition_jp"] = testCase.PreconditionJP
					updates["test_steps_jp"] = testCase.TestStepsJP
					updates["expected_result_jp"] = testCase.ExpectedResultJP
				} else if language == "EN" {
					updates["major_function_en"] = testCase.MajorFunctionEN
					updates["middle_function_en"] = testCase.MiddleFunctionEN
					updates["minor_function_en"] = testCase.MinorFunctionEN
					updates["precondition_en"] = testCase.PreconditionEN
					updates["test_steps_en"] = testCase.TestStepsEN
					updates["expected_result_en"] = testCase.ExpectedResultEN
				}

				if err := s.caseRepo.UpdateByCaseID(uuidStr, updates); err != nil {
					return 0, 0, fmt.Errorf("update case: %w", err)
				}
				updateCount++
			} else {
				// UUID不存在于数据库: 视为新用例，生成新UUID插入
				currentMaxID++
				testCase.ID = currentMaxID
				testCase.CaseID = uuid.New().String() // 生成新UUID
				if err := s.caseRepo.Create(testCase); err != nil {
					return 0, 0, fmt.Errorf("create case: %w", err)
				}
				insertCount++
			}
		} else {
			// 空UUID: 视为新用例，生成新UUID插入
			currentMaxID++
			testCase.ID = currentMaxID
			testCase.CaseID = uuid.New().String()

			fmt.Printf("\n➕ [ImportCases] 准备插入新用例:\n")
			fmt.Printf("  ID: %d\n", testCase.ID)
			fmt.Printf("  CaseID: %s\n", testCase.CaseID)
			fmt.Printf("  CaseGroup: %q\n", testCase.CaseGroup)
			fmt.Printf("  MajorFunctionCN: %q\n", testCase.MajorFunctionCN)

			if err := s.caseRepo.Create(testCase); err != nil {
				fmt.Printf("❌ [ImportCases] 插入失败: %v\n", err)
				return 0, 0, fmt.Errorf("create case: %w", err)
			}
			fmt.Printf("✅ [ImportCases] 插入成功\n")
			insertCount++
		}
	}

	return updateCount, insertCount, nil
}

// ExportAutoCasesAllLanguages 导出自动化用例为Excel(包含所有三种语言,19列)
func (s *excelService) ExportAutoCasesAllLanguages(cases []*models.AutoTestCase, filePath string) error {
	f := excelize.NewFile()
	sheetName := "用例数据"
	f.SetSheetName("Sheet1", sheetName)

	// 设置标题行(19列,使用英文标题参考手工测试用例格式)
	headers := []string{
		"No.", "CaseID",
		"ScreenCN", "ScreenJP", "ScreenEN",
		"FunctionCN", "FunctionJP", "FunctionEN",
		"PreconditionCN", "PreconditionJP", "PreconditionEN",
		"Test StepCN", "Test StepJP", "Test StepEN",
		"ExpectCN", "ExpectJP", "ExpectEN",
		"TestResult", "Remark",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%s1", columnLetter(i))
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return fmt.Errorf("set header: %w", err)
		}
	}

	// 设置数据行(包含所有语言字段)
	for i, tc := range cases {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), tc.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), tc.CaseNumber)

		// 画面(三语言)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), tc.ScreenCN)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), tc.ScreenJP)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), tc.ScreenEN)

		// 功能(三语言)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), tc.FunctionCN)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), tc.FunctionJP)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), tc.FunctionEN)

		// 前置条件(三语言)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), tc.PreconditionCN)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), tc.PreconditionJP)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), tc.PreconditionEN)

		// 测试步骤(三语言)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), tc.TestStepsCN)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), tc.TestStepsJP)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), tc.TestStepsEN)

		// 期待值(三语言)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), tc.ExpectedResultCN)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), tc.ExpectedResultJP)
		f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), tc.ExpectedResultEN)

		// 测试结果和备注
		f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), tc.TestResult)
		f.SetCellValue(sheetName, fmt.Sprintf("S%d", row), tc.Remark)
	}

	// 应用样式(设置列宽和自动换行)
	s.applyAutoExcelStyles(f, sheetName, len(cases))

	// 保存文件
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("save excel: %w", err)
	}
	return nil
}

// columnLetter 将列索引转换为Excel列字母(0->A, 1->B, ..., 18->S)
func columnLetter(index int) string {
	if index < 26 {
		return string(rune('A' + index))
	}
	// 处理超过Z的列(AA, AB, ...)
	return string(rune('A'+index/26-1)) + string(rune('A'+index%26))
}

// applyAutoExcelStyles 应用自动化测试用例Excel样式
func (s *excelService) applyAutoExcelStyles(f *excelize.File, sheetName string, rowCount int) {
	// 设置标题行样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	f.SetCellStyle(sheetName, "A1", "S1", headerStyle)

	// 设置数据行样式(自动换行)
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", WrapText: true},
	})
	if rowCount > 0 {
		f.SetCellStyle(sheetName, "A2", fmt.Sprintf("S%d", rowCount+1), dataStyle)
	}

	// 设置列宽(根据内容调整)
	columnWidths := map[string]float64{
		"A": 8, "B": 12, // ID, 用例编号
		"C": 15, "D": 15, "E": 15, // 画面
		"F": 20, "G": 20, "H": 20, // 功能
		"I": 25, "J": 25, "K": 25, // 前置条件
		"L": 30, "M": 30, "N": 30, // 测试步骤
		"O": 25, "P": 25, "Q": 25, // 期待值
		"R": 10, "S": 20, // 测试结果, 备注
	}
	for col, width := range columnWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	// 设置行高(标题行)
	f.SetRowHeight(sheetName, 1, 25)
}

// columnName 根据索引生成Excel列名(A-Z, AA-AZ, ...)
func columnName(index int) string {
	name := ""
	for index >= 0 {
		name = string(rune('A'+index%26)) + name
		index = index/26 - 1
	}
	return name
}

// ExportWebCasesByLanguage 导出Web用例到Excel（按语言）
// language: "All"(所有语言) / "CN"(仅中文) / "JP"(仅日文) / "EN"(仅英文)
func (s *excelService) ExportWebCasesByLanguage(projectName string, caseGroups []models.CaseGroup, cases []models.AutoTestCase, language string) ([]byte, string, error) {
	f := excelize.NewFile()

	// 1. 创建Cover页
	coverSheet := "Cover"
	coverIndex, _ := f.NewSheet(coverSheet)
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(coverIndex)

	// 写入Cover页内容
	f.SetCellValue(coverSheet, "A1", fmt.Sprintf("项目名称 / Project Name: %s", projectName))
	f.SetCellValue(coverSheet, "A2", fmt.Sprintf("导出时间 / Export Time: %s", time.Now().Format(time.RFC3339)))
	f.SetCellValue(coverSheet, "A3", fmt.Sprintf("用例总数 / Total Cases: %d", len(cases)))
	f.SetCellValue(coverSheet, "A4", "用例集列表 / Case Group List:")
	for i, cg := range caseGroups {
		f.SetCellValue(coverSheet, fmt.Sprintf("A%d", 5+i), fmt.Sprintf("%d. %s", i+1, cg.GroupName))
	}

	// 2. 按用例集分组并创建Sheet页
	casesByGroup := make(map[string][]models.AutoTestCase)
	for _, c := range cases {
		if c.CaseGroup != "" {
			casesByGroup[c.CaseGroup] = append(casesByGroup[c.CaseGroup], c)
		}
	}

	// 3. 为每个用例集创建Sheet页
	for _, cg := range caseGroups {
		groupCases, exists := casesByGroup[cg.GroupName]
		if !exists || len(groupCases) == 0 {
			continue
		}

		sheetName := cg.GroupName
		sheetIndex, err := f.NewSheet(sheetName)
		if err != nil {
			return nil, "", fmt.Errorf("create sheet %s: %w", sheetName, err)
		}
		f.SetActiveSheet(sheetIndex)

		// 根据语言设置表头
		var headers []string
		switch language {
		case "All":
			headers = []string{"No.", "CaseID",
				"ScreenCN", "FunctionCN", "PreconditionCN", "Test StepCN", "ExpectCN",
				"ScreenJP", "FunctionJP", "PreconditionJP", "Test StepJP", "ExpectJP",
				"ScreenEN", "FunctionEN", "PreconditionEN", "Test StepEN", "ExpectEN",
				"ScriptCode", "UUID"}
		case "CN":
			headers = []string{"No.", "CaseID", "ScreenCN", "FunctionCN", "PreconditionCN", "Test StepCN", "ExpectCN", "ScriptCode", "UUID"}
		case "JP":
			headers = []string{"No.", "CaseID", "ScreenJP", "FunctionJP", "PreconditionJP", "Test StepJP", "ExpectJP", "ScriptCode", "UUID"}
		case "EN":
			headers = []string{"No.", "CaseID", "ScreenEN", "FunctionEN", "PreconditionEN", "Test StepEN", "ExpectEN", "ScriptCode", "UUID"}
		default:
			return nil, "", fmt.Errorf("unsupported language: %s", language)
		}

		// 写入表头
		for i, h := range headers {
			cell := fmt.Sprintf("%s1", columnName(i))
			f.SetCellValue(sheetName, cell, h)
		}

		// 写入数据行
		for i, c := range groupCases {
			row := i + 2
			colIndex := 0

			// No.
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), i+1)
			colIndex++

			// CaseID
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.CaseNumber)
			colIndex++

			// 根据语言写入字段
			switch language {
			case "All":
				// CN字段
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ScreenCN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.FunctionCN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.PreconditionCN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.TestStepsCN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ExpectedResultCN)
				colIndex++

				// JP字段
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ScreenJP)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.FunctionJP)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.PreconditionJP)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.TestStepsJP)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ExpectedResultJP)
				colIndex++

				// EN字段
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ScreenEN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.FunctionEN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.PreconditionEN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.TestStepsEN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ExpectedResultEN)
				colIndex++

				// ScriptCode字段
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ScriptCode)
				colIndex++

			case "CN":
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ScreenCN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.FunctionCN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.PreconditionCN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.TestStepsCN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ExpectedResultCN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ScriptCode)
				colIndex++

			case "JP":
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ScreenJP)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.FunctionJP)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.PreconditionJP)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.TestStepsJP)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ExpectedResultJP)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ScriptCode)
				colIndex++

			case "EN":
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ScreenEN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.FunctionEN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.PreconditionEN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.TestStepsEN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ExpectedResultEN)
				colIndex++
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.ScriptCode)
				colIndex++
			}

			// UUID
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnName(colIndex), row), c.CaseID)
		}
	}

	// 4. 生成文件名和字节流
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_AIWeb_%s_TestCase_%s.xlsx", projectName, language, timestamp)

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", fmt.Errorf("write buffer: %w", err)
	}

	return buffer.Bytes(), filename, nil
}

// GenerateWebCasesZip 生成包含4个语言版本的zip包
func (s *excelService) GenerateWebCasesZip(projectID uint, projectName string, cases []models.AutoTestCase) (zipPath string, fileSize int64, err error) {
	// 1. 从用例中提取所有用例集信息（去重）
	caseGroupMap := make(map[string]bool)
	for _, c := range cases {
		if c.CaseGroup != "" {
			caseGroupMap[c.CaseGroup] = true
		}
	}

	// 构造用例集列表
	caseGroups := make([]models.CaseGroup, 0, len(caseGroupMap))
	for groupName := range caseGroupMap {
		caseGroups = append(caseGroups, models.CaseGroup{
			GroupName: groupName,
			CaseType:  "web",
			ProjectID: projectID,
		})
	}

	// 2. 使用临时目录存储Excel文件
	tmpDir := fmt.Sprintf("storage/tmp/web-cases-%d-%d", projectID, time.Now().Unix())
	err = s.createDir(tmpDir)
	if err != nil {
		return "", 0, fmt.Errorf("create tmp dir: %w", err)
	}
	defer s.removeDir(tmpDir)

	// 3. 并发生成4个Excel文件
	languages := []string{"All", "CN", "JP", "EN"}
	type excelResult struct {
		language string
		data     []byte
		filename string
		err      error
	}

	resultChan := make(chan excelResult, 4)

	for _, lang := range languages {
		go func(language string) {
			data, filename, err := s.ExportWebCasesByLanguage(projectName, caseGroups, cases, language)
			resultChan <- excelResult{language: language, data: data, filename: filename, err: err}
		}(lang)
	}

	// 4. 收集结果并保存文件
	excelFiles := make(map[string]string) // language -> filepath
	for i := 0; i < 4; i++ {
		result := <-resultChan
		if result.err != nil {
			return "", 0, fmt.Errorf("generate %s excel: %w", result.language, result.err)
		}

		filePath := fmt.Sprintf("%s/%s", tmpDir, result.filename)
		err := s.writeFile(filePath, result.data)
		if err != nil {
			return "", 0, fmt.Errorf("write %s file: %w", result.language, err)
		}

		excelFiles[result.language] = filePath
	}
	close(resultChan)

	// 5. 打包为zip
	timestamp := time.Now().Format("20060102_150405")
	zipFilename := fmt.Sprintf("%s_AIWeb_TestCase_%s.zip", projectName, timestamp)
	zipDir := fmt.Sprintf("storage/versions/web-cases/%d", projectID)
	err = s.createDir(zipDir)
	if err != nil {
		return "", 0, fmt.Errorf("create zip dir: %w", err)
	}

	zipPath = fmt.Sprintf("%s/%s", zipDir, zipFilename)
	err = s.createZipArchive(zipPath, excelFiles)
	if err != nil {
		return "", 0, fmt.Errorf("create zip archive: %w", err)
	}

	// 6. 获取文件大小
	fileSize, err = s.getFileSize(zipPath)
	if err != nil {
		return "", 0, fmt.Errorf("get file size: %w", err)
	}

	return zipPath, fileSize, nil
}

// 辅助方法：创建目录
func (s *excelService) createDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// 辅助方法：删除目录
func (s *excelService) removeDir(dirPath string) error {
	return os.RemoveAll(dirPath)
}

// 辅助方法：写入文件
func (s *excelService) writeFile(filePath string, data []byte) error {
	return os.WriteFile(filePath, data, 0644)
}

// 辅助方法：创建zip归档
func (s *excelService) createZipArchive(zipPath string, files map[string]string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, filePath := range files {
		// 读取Excel文件
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read file %s: %w", filePath, err)
		}

		// 获取文件名（不包含路径）
		_, filename := filepath.Split(filePath)

		// 在zip中创建文件
		writer, err := zipWriter.Create(filename)
		if err != nil {
			return fmt.Errorf("create zip entry %s: %w", filename, err)
		}

		// 写入数据
		_, err = writer.Write(fileData)
		if err != nil {
			return fmt.Errorf("write zip entry %s: %w", filename, err)
		}
	}

	return nil
}

// 辅助方法：获取文件大小
func (s *excelService) getFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("stat file: %w", err)
	}
	return fileInfo.Size(), nil
}
