package services

import (
	"bytes"
	"fmt"
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
	ExportCases(projectID uint, caseType string, taskUUID string) ([]byte, string, error)
	ImportCases(projectID uint, caseType string, fileData []byte) (updateCount, insertCount int, err error)
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
func (s *excelService) ExportCases(projectID uint, caseType string, taskUUID string) ([]byte, string, error) {
	// 1. 查询元数据和用例数据
	metadata, _ := s.caseRepo.GetMetadataByProjectID(projectID, caseType)
	cases, err := s.caseRepo.GetByProjectAndTypeOrdered(projectID, caseType)
	if err != nil {
		return nil, "", fmt.Errorf("get cases: %w", err)
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
	dataSheet := "用例数据"
	index, _ := f.NewSheet(dataSheet)
	f.SetActiveSheet(index)

	// 根据是否有执行结果决定列数(23列或25列)
	var headers []string
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

	for i, h := range headers {
		cell := fmt.Sprintf("%s1", columnName(i))
		f.SetCellValue(dataSheet, cell, h)
	}

	// 5. 写入数据行
	for i, c := range cases {
		row := i + 2
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

	typeMap := map[string]string{"overall": "Overall", "change": "Change", "acceptance": "Acceptance"}
	filename := fmt.Sprintf("%s_%s_Cases_%s.xlsx", projectName, typeMap[caseType], time.Now().Format("2006-01-02_150405"))

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", fmt.Errorf("write buffer: %w", err)
	}

	return buffer.Bytes(), filename, nil
}

// ImportCases 导入用例(UUID匹配逻辑)
func (s *excelService) ImportCases(projectID uint, caseType string, fileData []byte) (int, int, error) {
	// 1. 打开Excel文件
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return 0, 0, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	// 2. 读取Sheet2数据
	rows, err := f.GetRows("用例数据")
	if err != nil {
		return 0, 0, fmt.Errorf("get rows: %w", err)
	}

	if len(rows) < 2 {
		return 0, 0, fmt.Errorf("no data rows")
	}

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

		// 读取关键字段
		caseNumber := getCol(1)
		uuidStr := getCol(22)

		// 读取所有字段
		majorFunctionCN := getCol(2)
		majorFunctionJP := getCol(3)
		majorFunctionEN := getCol(4)
		middleFunctionCN := getCol(5)
		middleFunctionJP := getCol(6)
		middleFunctionEN := getCol(7)
		minorFunctionCN := getCol(8)
		minorFunctionJP := getCol(9)
		minorFunctionEN := getCol(10)
		preconditionCN := getCol(11)
		preconditionJP := getCol(12)
		preconditionEN := getCol(13)
		testStepsCN := getCol(14)
		testStepsJP := getCol(15)
		testStepsEN := getCol(16)
		expectedResultCN := getCol(17)
		expectedResultJP := getCol(18)
		expectedResultEN := getCol(19)
		testResult := getCol(20)
		remark := getCol(21)

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
		testCase := &models.ManualTestCase{
			ProjectID:        projectID,
			CaseType:         caseType,
			CaseNumber:       caseNumber,
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
				updates := map[string]interface{}{
					"case_number":        testCase.CaseNumber,
					"major_function_cn":  testCase.MajorFunctionCN,
					"major_function_jp":  testCase.MajorFunctionJP,
					"major_function_en":  testCase.MajorFunctionEN,
					"middle_function_cn": testCase.MiddleFunctionCN,
					"middle_function_jp": testCase.MiddleFunctionJP,
					"middle_function_en": testCase.MiddleFunctionEN,
					"minor_function_cn":  testCase.MinorFunctionCN,
					"minor_function_jp":  testCase.MinorFunctionJP,
					"minor_function_en":  testCase.MinorFunctionEN,
					"precondition_cn":    testCase.PreconditionCN,
					"precondition_jp":    testCase.PreconditionJP,
					"precondition_en":    testCase.PreconditionEN,
					"test_steps_cn":      testCase.TestStepsCN,
					"test_steps_jp":      testCase.TestStepsJP,
					"test_steps_en":      testCase.TestStepsEN,
					"expected_result_cn": testCase.ExpectedResultCN,
					"expected_result_jp": testCase.ExpectedResultJP,
					"expected_result_en": testCase.ExpectedResultEN,
					"test_result":        testCase.TestResult,
					"remark":             testCase.Remark,
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
			if err := s.caseRepo.Create(testCase); err != nil {
				return 0, 0, fmt.Errorf("create case: %w", err)
			}
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
