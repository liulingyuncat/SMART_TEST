package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"webtest/internal/models"
)

// ScriptTestRequest 脚本测试请求
type ScriptTestRequest struct {
	ScriptCode string `json:"script_code" binding:"required"`
	GroupID    uint   `json:"group_id"`   // 用例集ID，用于获取变量
	GroupType  string `json:"group_type"` // 用例集类型：web 或 api
	ProjectID  uint   `json:"project_id"` // 项目ID
}

// ScriptTestResult 脚本测试结果
type ScriptTestResult struct {
	Success      bool      `json:"success"`
	Output       string    `json:"output"`
	ErrorMessage string    `json:"error_message,omitempty"`
	ResponseTime int       `json:"response_time"` // 毫秒
	ExecutedAt   time.Time `json:"executed_at"`
}

// ScriptTestService 脚本测试服务接口
type ScriptTestService interface {
	// TestScript 测试脚本（直接执行，不保存结果）
	TestScript(projectID uint, userID uint, req ScriptTestRequest) (*ScriptTestResult, error)
}

type scriptTestService struct {
	pwClient        *PlaywrightExecutorClient
	variableService UserDefinedVariableService
}

// NewScriptTestService 创建脚本测试服务实例
func NewScriptTestService(
	variableService UserDefinedVariableService,
) ScriptTestService {
	pwClient := NewPlaywrightExecutorClient(DefaultExecutorConfig())
	return &scriptTestService{
		pwClient:        pwClient,
		variableService: variableService,
	}
}

// TestScript 测试脚本
func (s *scriptTestService) TestScript(projectID uint, userID uint, req ScriptTestRequest) (*ScriptTestResult, error) {
	fmt.Printf("[ScriptTest] 开始测试脚本: projectID=%d, userID=%d, groupID=%d, groupType=%s\n",
		projectID, userID, req.GroupID, req.GroupType)

	// 1. 检查脚本是否为空
	if req.ScriptCode == "" {
		return nil, errors.New("脚本代码不能为空")
	}

	// 2. 获取用例集变量
	var variables []*models.UserDefinedVariable
	var err error
	if req.GroupID > 0 && req.GroupType != "" {
		variables, err = s.variableService.GetVariablesByGroup(req.GroupID, req.GroupType)
		if err != nil {
			fmt.Printf("[ScriptTest] 警告: 获取变量失败: %v\n", err)
			variables = []*models.UserDefinedVariable{} // 继续执行，不中断
		}
		fmt.Printf("[ScriptTest] 获取到 %d 个变量\n", len(variables))
		// 🔍 打印变量详情用于调试
		for i, v := range variables {
			fmt.Printf("[ScriptTest]   变量 %d: var_key=%s, var_value=%s (长度:%d)\n",
				i+1, v.VarKey, maskSensitive(v.VarKey, v.VarValue), len(v.VarValue))
		}
	} else {
		fmt.Printf("[ScriptTest] 跳过变量获取 (groupID=%d, groupType=%s)\n", req.GroupID, req.GroupType)
	}

	// 3. 替换脚本中的变量
	replacedScript := s.replaceVariables(req.ScriptCode, variables)

	// 🔍 调试日志：检查变量替换效果
	fmt.Printf("[ScriptTest] 脚本替换前长度: %d bytes\n", len(req.ScriptCode))
	fmt.Printf("[ScriptTest] 脚本替换后长度: %d bytes\n", len(replacedScript))
	if strings.Contains(req.ScriptCode, "${") {
		fmt.Printf("[ScriptTest] ⚠️ 原始脚本包含变量占位符\n")
	}
	if strings.Contains(replacedScript, "${") {
		fmt.Printf("[ScriptTest] ❌ 替换后脚本仍包含 '${' ，变量替换可能失败！\n")
		// 打印前200个字符用于调试
		preview := replacedScript
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		fmt.Printf("[ScriptTest] 脚本预览: %s\n", preview)
	} else {
		fmt.Printf("[ScriptTest] ✅ 变量替换成功\n")
	}

	// 4. 执行脚本
	fmt.Printf("[ScriptTest] 开始执行 Playwright 脚本...\n")
	ctx := context.Background()
	execResult, execErr := s.pwClient.Execute(ctx, replacedScript)

	now := time.Now()

	if execErr != nil {
		fmt.Printf("[ScriptTest] 执行失败: %v\n", execErr)
		responseTime := 0
		errorMessage := execErr.Error()

		// 🚨 错误分析：检查是否缺少变量
		if strings.Contains(replacedScript, "${BASE_URL}") {
			errorMessage += "\n\n⚠️ 错误分析: 脚本中仍包含 '${BASE_URL}'，说明该变量未被替换。\n请检查变量表中是否存在 'base_url' 变量。"
		} else if strings.Contains(replacedScript, "${") {
			errorMessage += "\n\n⚠️ 错误分析: 脚本中可能存在未替换的变量 (检测到 '${' 符号)。\n请检查变量表配置。"
		}

		if execResult != nil {
			responseTime = execResult.ResponseTime
		}
		return &ScriptTestResult{
			Success:      false,
			Output:       "",
			ErrorMessage: errorMessage,
			ResponseTime: responseTime,
			ExecutedAt:   now,
		}, nil
	}

	fmt.Printf("[ScriptTest] 执行成功: response_time=%dms\n", execResult.ResponseTime)
	return &ScriptTestResult{
		Success:      true,
		Output:       execResult.Output,
		ErrorMessage: "",
		ResponseTime: execResult.ResponseTime,
		ExecutedAt:   now,
	}, nil
}

// maskSensitive masks sensitive variable values like passwords
func maskSensitive(key, value string) string {
	lowerKey := strings.ToLower(key)
	if strings.Contains(lowerKey, "password") || strings.Contains(lowerKey, "secret") || strings.Contains(lowerKey, "token") {
		if len(value) <= 3 {
			return "***"
		}
		return value[:2] + "***"
	}
	return value
}

// replaceVariables 替换脚本中的变量占位符
// 支持 ${VAR_NAME} 格式（脚本标准格式）和 {{VAR_KEY}} 格式（兼容旧格式）
func (s *scriptTestService) replaceVariables(script string, variables []*models.UserDefinedVariable) string {
	if len(variables) == 0 {
		fmt.Printf("[ScriptTest] 变量列表为空，跳过替换\n")
		return script
	}

	result := script
	replacedCount := 0
	for _, v := range variables {
		// 1. 使用 VarKey 替换大写格式 "${BASE_URL}"
		if v.VarKey != "" {
			upperKey := strings.ToUpper(v.VarKey)
			placeholder := fmt.Sprintf("${%s}", upperKey)
			if strings.Contains(result, placeholder) {
				result = strings.ReplaceAll(result, placeholder, v.VarValue)
				replacedCount++
				fmt.Printf("[ScriptTest] 替换变量: ${%s} -> %s\n", upperKey, maskSensitive(v.VarKey, v.VarValue))
			}
		}

		// 2. 同时支持小写格式 "${base_url}" (增强兼容性)
		if v.VarKey != "" {
			lowerKey := strings.ToLower(v.VarKey)
			placeholder := fmt.Sprintf("${%s}", lowerKey)
			if strings.Contains(result, placeholder) {
				result = strings.ReplaceAll(result, placeholder, v.VarValue)
				replacedCount++
				fmt.Printf("[ScriptTest] 替换变量: ${%s} -> %s\n", lowerKey, maskSensitive(v.VarKey, v.VarValue))
			}
		}

		// 3. 兼容：使用 VarName 字段
		if v.VarName != "" {
			result = strings.ReplaceAll(result, v.VarName, v.VarValue)
		}

		// 4. 兼容：{{key}} 格式
		if v.VarKey != "" {
			placeholder := fmt.Sprintf("{{%s}}", v.VarKey)
			if strings.Contains(result, placeholder) {
				result = strings.ReplaceAll(result, placeholder, v.VarValue)
				replacedCount++
				fmt.Printf("[ScriptTest] 替换变量: {{%s}} -> %s\n", v.VarKey, maskSensitive(v.VarKey, v.VarValue))
			}
		}
	}

	fmt.Printf("[ScriptTest] 变量替换完成: 共替换 %d 个变量\n", replacedCount)
	return result
}
