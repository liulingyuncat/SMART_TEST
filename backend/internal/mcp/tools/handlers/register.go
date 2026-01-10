package handlers

import (
	"webtest/internal/mcp/client"
	"webtest/internal/mcp/tools"
)

// RegisterAllTools registers all MCP tool handlers to the registry.
// Total: 37 tools across 11 handler files.
func RegisterAllTools(registry *tools.ToolRegistry, c *client.BackendClient) {
	// ==================== 用户与项目信息相关 (1 tool) ====================
	// registry.Register(NewGetCurrentUserInfoHandler(c)) // 已禁用：提示词中不再使用
	registry.Register(NewGetCurrentProjectNameHandler(c))

	// ==================== 原始文档相关 (2 tools) ====================
	registry.Register(NewListRawDocumentsHandler(c))
	registry.Register(NewGetConvertedDocumentHandler(c))

	// ==================== 需求条目相关 (4 tools) ====================
	registry.Register(NewListRequirementItemsHandler(c))
	registry.Register(NewGetRequirementItemHandler(c))
	registry.Register(NewCreateRequirementItemHandler(c))
	registry.Register(NewUpdateRequirementItemHandler(c))

	// ==================== 测试观点相关 (4 tools) ====================
	registry.Register(NewListViewpointItemsHandler(c))
	registry.Register(NewGetViewpointItemHandler(c))
	registry.Register(NewCreateViewpointItemHandler(c))
	registry.Register(NewUpdateViewpointItemHandler(c))

	// ==================== 用例集与手工用例相关 (6 tools) ====================
	registry.Register(NewListManualGroupsHandler(c))
	registry.Register(NewListManualCasesHandler(c))
	registry.Register(NewCreateCaseGroupHandler(c))
	registry.Register(NewCreateManualCasesHandler(c))
	registry.Register(NewUpdateManualCaseHandler(c))
	registry.Register(NewUpdateManualCasesHandler(c))

	// ==================== Web自动化用例相关 (5 tools) ====================
	registry.Register(NewListWebGroupsHandler(c))
	registry.Register(NewGetWebGroupMetadataHandler(c))
	registry.Register(NewListWebCasesHandler(c))
	// registry.Register(NewCreateWebGroupHandler(c)) // 已禁用：不再使用
	registry.Register(NewUpdateWebCasesHandler(c))
	registry.Register(NewCreateWebCasesHandler(c))

	// ==================== API接口用例相关 (5 tools) ====================
	registry.Register(NewListApiGroupsHandler(c))
	registry.Register(NewGetApiGroupMetadataHandler(c))
	registry.Register(NewListApiCasesHandler(c))
	// registry.Register(NewCreateApiGroupHandler(c)) // 已禁用：不再使用
	registry.Register(NewCreateApiCaseHandler(c))
	registry.Register(NewUpdateApiCaseHandler(c))

	// ==================== 用例评审相关 (1 tool) ====================
	registry.Register(NewCreateReviewItemHandler(c))

	// ==================== 执行任务相关 (4 tools) ====================
	registry.Register(NewListExecutionTasksHandler(c))
	registry.Register(NewGetExecutionTaskMetadataHandler(c))
	registry.Register(NewGetExecutionTaskCasesHandler(c))
	registry.Register(NewUpdateExecutionCaseResultHandler(c))

	// ==================== 缺陷管理相关 (2 tools) ====================
	registry.Register(NewListDefectsHandler(c))
	registry.Register(NewUpdateDefectHandler(c))

	// ==================== AI报告相关 (2 tools) ====================
	registry.Register(NewCreateAIReportHandler(c))
	registry.Register(NewUpdateAIReportHandler(c))
}
