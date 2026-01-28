package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// GetCurrentTimestamp 返回当前时间的格式化字符串（格式：20060102_150405）
func GetCurrentTimestamp() string {
	return time.Now().Format("20060102_150405")
}

// CreateAIReportHandler handles creating an AI report.
type CreateAIReportHandler struct {
	*BaseHandler
}

func NewCreateAIReportHandler(c *client.BackendClient) *CreateAIReportHandler {
	return &CreateAIReportHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *CreateAIReportHandler) Name() string {
	return "create_ai_report"
}

func (h *CreateAIReportHandler) Description() string {
	return "创建AI测试报告，根据报告类型自动生成报告名称"
}

func (h *CreateAIReportHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"report_type": map[string]interface{}{
				"type":        "string",
				"description": "报告类型：R-用例审阅, A-品质分析, T-测试结果, O-其他",
				"enum":        []interface{}{"R", "A", "T", "O"},
			},
			"case_group_name": map[string]interface{}{
				"type":        "string",
				"description": "用例集名称（仅当report_type为R时需要）",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "报告内容（Markdown格式）",
			},
		},
		"required": []interface{}{"project_id", "report_type", "content"},
	}
}

func (h *CreateAIReportHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	reportType, err := GetString(args, "report_type")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 验证report_type
	if reportType != "R" && reportType != "A" && reportType != "T" && reportType != "O" {
		return tools.NewErrorResult("report_type必须是R/A/T/O其中之一"), nil
	}

	content, err := GetString(args, "content")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 生成报告标题
	var title string
	timestamp := GetCurrentTimestamp()

	switch reportType {
	case "R": // 用例审阅
		caseGroupName := GetOptionalString(args, "case_group_name", "")
		if caseGroupName == "" {
			return tools.NewErrorResult("report_type为R时，case_group_name参数为必填项"), nil
		}
		title = fmt.Sprintf("%s_Review_%s", caseGroupName, timestamp)
	case "A": // 品质分析
		title = fmt.Sprintf("Quality_Analyse_%s", timestamp)
	case "T": // 测试结果
		title = fmt.Sprintf("TestResult_Analyse_%s", timestamp)
	case "O": // 其他
		title = fmt.Sprintf("Others_%s", timestamp)
	}

	// 首先创建报告（只设置名称）
	createBody := map[string]interface{}{
		"name": title, // 使用自动生成的标题
	}

	path := fmt.Sprintf("/api/v1/projects/%d/ai-reports", projectID)
	data, err := h.client.Post(ctx, path, createBody)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 解析创建响应，获取报告ID
	var createResponse struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &createResponse); err != nil {
		return tools.NewErrorResult("解析创建响应失败: " + err.Error()), nil
	}

	reportID := createResponse.Data.ID
	if reportID == "" {
		return tools.NewErrorResult("创建报告失败：未获取到报告ID"), nil
	}

	// 如果有内容，则更新报告内容
	if content != "" {
		updateBody := map[string]interface{}{
			"content": content,
		}

		updatePath := fmt.Sprintf("/api/v1/projects/%d/ai-reports/%s", projectID, reportID)
		updateData, err := h.client.Put(ctx, updatePath, updateBody)
		if err != nil {
			return tools.NewErrorResult("更新报告内容失败: " + err.Error()), nil
		}

		return tools.NewJSONResult(string(updateData)), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// UpdateAIReportHandler handles updating an AI report.
type UpdateAIReportHandler struct {
	*BaseHandler
}

func NewUpdateAIReportHandler(c *client.BackendClient) *UpdateAIReportHandler {
	return &UpdateAIReportHandler{BaseHandler: NewBaseHandler(c)}
}

func (h *UpdateAIReportHandler) Name() string {
	return "update_ai_report"
}

func (h *UpdateAIReportHandler) Description() string {
	return "更新AI测试报告，支持通过报告ID或报告名称（精确或模糊匹配）更新"
}

func (h *UpdateAIReportHandler) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "integer",
				"description": "项目ID",
			},
			"id": map[string]interface{}{
				"type":        "string",
				"description": "报告ID（字符串格式，如 report_xxx，可选）",
			},
			"title": map[string]interface{}{
				"type":        "string",
				"description": "报告名称/标题（可选，支持精确匹配或模糊匹配。如果提供则通过名称查找报告，用于更新其他字段）",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "报告内容（Markdown格式，可选）",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "报告新名称（可选，用于重命名报告）",
			},
		},
		"required": []interface{}{"project_id"},
	}
}

func (h *UpdateAIReportHandler) Execute(ctx context.Context, args map[string]interface{}) (tools.ToolResult, error) {
	projectID, err := GetInt(args, "project_id")
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	// 获取报告ID：优先通过 title（报告名称）查询，其次使用 id 参数
	var reportID string

	// 优先使用 title 参数通过名称查找报告（更用户友好）
	if titleToFind := GetOptionalString(args, "title", ""); titleToFind != "" {
		foundReportID, err := h.findReportByName(ctx, projectID, titleToFind)
		if err != nil {
			return tools.NewErrorResult(fmt.Sprintf("通过报告名称 '%s' 查找报告失败: %s", titleToFind, err.Error())), nil
		}
		reportID = foundReportID
	} else if id := GetOptionalString(args, "id", ""); id != "" {
		// 验证 id 格式（应该以 report_ 开头）
		if !strings.HasPrefix(id, "report_") {
			return tools.NewErrorResult(fmt.Sprintf("无效的报告ID格式 '%s'，报告ID应该以 'report_' 开头（如 report_xxx）。建议使用 'title' 参数通过报告名称来查找报告", id)), nil
		}
		reportID = id
	} else {
		return tools.NewErrorResult("必须提供 'id' 或 'title' 参数中的至少一个来标识要更新的报告"), nil
	}

	if reportID == "" {
		return tools.NewErrorResult("无法确定要更新的报告ID"), nil
	}

	// 构建更新数据
	body := make(map[string]interface{})

	// 处理名称更新（使用 name 参数表示新名称）
	if newName := GetOptionalString(args, "name", ""); newName != "" {
		body["name"] = newName
	}

	// 处理内容更新
	if content := GetOptionalString(args, "content", ""); content != "" {
		body["content"] = content
	}

	if len(body) == 0 {
		return tools.NewErrorResult("至少需要提供 'content' 或 'name' 其中一个字段进行更新"), nil
	}

	path := fmt.Sprintf("/api/v1/projects/%d/ai-reports/%s", projectID, reportID)
	data, err := h.client.Put(ctx, path, body)
	if err != nil {
		return tools.NewErrorResult(err.Error()), nil
	}

	return tools.NewJSONResult(string(data)), nil
}

// findReportByName 通过报告名称查找报告ID（支持精确匹配和模糊匹配）
func (h *UpdateAIReportHandler) findReportByName(ctx context.Context, projectID int, reportName string) (string, error) {
	// 调用 ListReports API 获取所有报告
	path := fmt.Sprintf("/api/v1/projects/%d/ai-reports", projectID)
	data, err := h.client.Get(ctx, path, nil)
	if err != nil {
		return "", fmt.Errorf("获取报告列表失败: %v", err)
	}

	// 解析响应
	var listResponse struct {
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &listResponse); err != nil {
		return "", fmt.Errorf("解析报告列表响应失败: %v", err)
	}

	// 收集所有报告名称用于错误提示
	var allNames []string
	var exactMatch string
	var fuzzyMatches []struct {
		ID   string
		Name string
	}

	// 使用小写进行比较以支持大小写不敏感的匹配
	lowerSearchName := strings.ToLower(reportName)

	for _, report := range listResponse.Data {
		allNames = append(allNames, report.Name)

		// 精确匹配（优先）
		if report.Name == reportName {
			exactMatch = report.ID
		}

		// 模糊匹配：名称包含搜索词（大小写不敏感）
		if strings.Contains(strings.ToLower(report.Name), lowerSearchName) {
			fuzzyMatches = append(fuzzyMatches, struct {
				ID   string
				Name string
			}{ID: report.ID, Name: report.Name})
		}
	}

	// 优先返回精确匹配
	if exactMatch != "" {
		return exactMatch, nil
	}

	// 如果只有一个模糊匹配，返回它
	if len(fuzzyMatches) == 1 {
		return fuzzyMatches[0].ID, nil
	}

	// 如果有多个模糊匹配，返回错误提示用户精确指定
	if len(fuzzyMatches) > 1 {
		var matchedNames []string
		for _, m := range fuzzyMatches {
			matchedNames = append(matchedNames, m.Name)
		}
		return "", fmt.Errorf("找到多个匹配的报告 '%s'，请精确指定报告名称。匹配的报告: %v", reportName, matchedNames)
	}

	// 没有任何匹配
	return "", fmt.Errorf("未找到名称包含 '%s' 的报告。当前项目中的报告列表: %v", reportName, allNames)
}
