package prompts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"webtest/internal/mcp/client"
	"webtest/internal/models"
)

// PromptsListParams prompts/list请求参数
type PromptsListParams struct {
	UserID     uint   `json:"user_id,omitempty"`
	UserRole   string `json:"user_role,omitempty"`
	ProjectID  uint   `json:"project_id,omitempty"`
	Scope      string `json:"scope,omitempty"`       // 可选: system/project/user/all
	IncludeAll bool   `json:"include_all,omitempty"` // 是否返回所有可见提示词
	APIToken   string `json:"api_token,omitempty"`   // 用于认证后端API的Token
}

// HandlePromptsList 处理prompts/list请求
func HandlePromptsList(
	ctx context.Context,
	registry *PromptsRegistry,
	backendClient *client.BackendClient,
	params json.RawMessage,
) (interface{}, error) {
	// 解析参数
	var listParams PromptsListParams
	if err := json.Unmarshal(params, &listParams); err != nil {
		// 参数解析失败，返回系统提示词
		systemPrompts := registry.List()
		return map[string]interface{}{
			"prompts": systemPrompts,
		}, nil
	}

	// 如果params中没有api_token，尝试从context中提取
	if listParams.APIToken == "" {
		if token, ok := ctx.Value(client.TokenContextKey).(string); ok && token != "" {
			listParams.APIToken = token
			log.Printf("[MCP] Extracted api_token from context")
		}
	}

	log.Printf("[MCP] prompts/list called with: userID=%d, projectID=%d, apiToken=%s (len=%d)",
		listParams.UserID, listParams.ProjectID,
		func() string {
			if listParams.APIToken != "" {
				return "***"
			} else {
				return "empty"
			}
		}(),
		len(listParams.APIToken))

	// 从PromptsRegistry获取系统Prompt元数据
	systemPrompts := registry.List()
	log.Printf("[MCP] Found %d system prompts from registry", len(systemPrompts))

	// 准备返回的prompts列表
	allPrompts := make([]PromptMetadata, 0, len(systemPrompts)*3)

	// 添加系统Prompt（全员可见，无需认证）
	allPrompts = append(allPrompts, systemPrompts...)
	log.Printf("[MCP] Added %d system prompts to result", len(systemPrompts))

	// 从后端HTTP API获取全员和个人提示词
	// 即使没有userID也应该获取全员提示词（scope='project'）
	if backendClient != nil {
		// 优先使用提供的projectID，默认为1
		projectID := listParams.ProjectID
		if projectID == 0 {
			projectID = 1
		}

		// 如果MCP请求中提供了Token，添加到context中
		reqCtx := ctx
		userIDForPersonal := listParams.UserID // 用户明确指定的ID

		if listParams.APIToken != "" {
			reqCtx = context.WithValue(ctx, client.TokenContextKey, listParams.APIToken)
			log.Printf("[MCP] APIToken is present, length: %d", len(listParams.APIToken))

			// 如果没有明确指定userID，从Token中自动提取
			if userIDForPersonal == 0 {
				log.Printf("[MCP] userIDForPersonal is 0, attempting to extract from token...")
				userID, err := getUserIDFromToken(reqCtx, backendClient, listParams.APIToken)
				if err != nil {
					log.Printf("[MCP] Failed to extract userID from token: %v", err)
					// 错误不阻塞，继续返回系统+全员提示词
				} else if userID > 0 {
					userIDForPersonal = userID
					log.Printf("[MCP] Extracted userID=%d from token", userID)
				} else {
					log.Printf("[MCP] getUserIDFromToken returned 0 without error")
				}
			} else {
				log.Printf("[MCP] userIDForPersonal already set to %d, skipping token extraction", userIDForPersonal)
			}
		} else {
			log.Printf("[MCP] No APIToken in listParams")
		}

		// 获取全员提示词（scope='project'，无需认证）
		// 即使没有Token也可以获取全员提示词
		// 优先尝试使用公开API（无需认证），如果失败则使用需认证的API
		log.Printf("[MCP] Fetching project prompts (scope=project)...")
		projectPrompts, err := getCustomPromptsFromBackend(
			ctx, // 不需要Token，使用原始context
			backendClient,
			0, // userID=0，不查询个人提示词
			listParams.UserRole,
			projectID,
			"project", // 明确指定获取全员提示词
		)
		if err != nil {
			log.Printf("[MCP] Error fetching project prompts: %v", err)
		} else {
			log.Printf("[MCP] Fetched %d project prompts", len(projectPrompts))
		}
		if err == nil && len(projectPrompts) > 0 {
			allPrompts = append(allPrompts, projectPrompts...)
			log.Printf("[MCP] Added %d project prompts, total now: %d", len(projectPrompts), len(allPrompts))
		}
		// 错误不阻塞，继续

		// 如果有有效的userID（无论是明确指定还是从Token提取），获取个人提示词
		if userIDForPersonal > 0 {
			log.Printf("[MCP] Attempting to fetch user prompts: userID=%d, projectID=%d, hasToken=%v",
				userIDForPersonal, projectID, ctx.Value(client.TokenContextKey) != nil)
			log.Printf("[MCP] Context token value: %v", reqCtx.Value(client.TokenContextKey))

			userPrompts, err := getCustomPromptsFromBackend(
				reqCtx,
				backendClient,
				userIDForPersonal,
				listParams.UserRole,
				projectID,
				"user", // 获取用户个人提示词
			)
			if err != nil {
				log.Printf("[MCP] Error fetching user prompts (userID=%d): %v", userIDForPersonal, err)
			} else {
				log.Printf("[MCP] Fetched %d user prompts for userID=%d", len(userPrompts), userIDForPersonal)
			}
			if err == nil && len(userPrompts) > 0 {
				allPrompts = append(allPrompts, userPrompts...)
				log.Printf("[MCP] Added %d user prompts, total now: %d", len(userPrompts), len(allPrompts))
			} else if err == nil {
				log.Printf("[MCP] No user prompts returned (empty list)")
			}
			// 错误不阻塞，继续返回系统提示词
		} else {
			log.Printf("[MCP] No valid userID found, skipping personal prompts (userIDForPersonal=%d)", userIDForPersonal)
		}
	}

	// 返回符合MCP规范的响应
	log.Printf("[MCP] Returning total of %d prompts", len(allPrompts))
	return map[string]interface{}{
		"prompts": allPrompts,
	}, nil
}

// getUserIDFromToken 从Token中提取userID
func getUserIDFromToken(ctx context.Context, backendClient *client.BackendClient, token string) (uint, error) {
	if token == "" {
		log.Printf("[getUserIDFromToken] Empty token provided")
		return 0, fmt.Errorf("empty token")
	}

	log.Printf("[getUserIDFromToken] Calling /api/v1/auth/me with token (len=%d)", len(token))

	// 将Token添加到context中
	reqCtx := context.WithValue(ctx, client.TokenContextKey, token)

	// 调用后端API的 /auth/me 端点获取当前用户信息
	result, err := backendClient.Get(reqCtx, "/api/v1/auth/me", nil)
	if err != nil {
		log.Printf("[getUserIDFromToken] API call failed: %v", err)
		return 0, fmt.Errorf("failed to fetch user info from backend: %w", err)
	}

	log.Printf("[getUserIDFromToken] API call succeeded, parsing response...")

	// 解析返回的数据
	var response struct {
		Code int `json:"code"`
		Data struct {
			UserID   uint   `json:"user_id"`  // 后端返回user_id
			Username string `json:"username"` // 后端返回username
		} `json:"data"`
	}

	if err := json.Unmarshal(result, &response); err != nil {
		log.Printf("[getUserIDFromToken] Failed to parse response: %v", err)
		return 0, fmt.Errorf("failed to parse user response: %w", err)
	}

	if response.Code != 0 {
		log.Printf("[getUserIDFromToken] Backend returned error code: %d", response.Code)
		return 0, fmt.Errorf("backend returned error code: %d", response.Code)
	}

	log.Printf("[getUserIDFromToken] Successfully extracted userID=%d (username=%s)", response.Data.UserID, response.Data.Username)
	return response.Data.UserID, nil
}

// getCustomPromptsFromBackend 从后端API获取自定义提示词
func getCustomPromptsFromBackend(
	ctx context.Context,
	backendClient *client.BackendClient,
	userID uint,
	userRole string,
	projectID uint,
	scope string,
) ([]PromptMetadata, error) {
	// 构造请求参数（转换为字符串类型）
	// 提示词与project_id完全无关，不传递project_id
	params := map[string]string{
		"page":      "1",
		"page_size": "1000", // 获取更大的限制以支持/获取所有提示词
	}

	// 对于个人提示词（scope=user），不传递user_id参数，后端会从token中自动提取userID
	// 对于其他scope，可以传递user_id和user_role
	if scope != "user" && userID > 0 {
		params["user_id"] = fmt.Sprintf("%d", userID)
	}
	if userRole != "" {
		params["user_role"] = userRole
	}

	// 设置scope参数（如果指定了scope）
	if scope != "" {
		params["scope"] = scope
	}

	log.Printf("[MCP] Fetching prompts: scope=%s, userID=%d, hasToken=%v, tokenInCtx=%v",
		scope, userID, ctx.Value(client.TokenContextKey) != nil,
		fmt.Sprintf("%v", ctx.Value(client.TokenContextKey)))

	// 调用后端API
	// 对于全员提示词（scope=project），优先使用公开API（不需要认证）
	apiPath := "/api/v1/prompts"
	if scope == "project" {
		apiPath = "/api/v1/prompts/public"
	}

	log.Printf("[MCP] Calling backend API: %s with params: %+v", apiPath, params)

	result, err := backendClient.Get(ctx, apiPath, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch custom prompts from backend (scope=%s, path=%s): %w", scope, apiPath, err)
	}

	// 解析返回的数据
	var response struct {
		Code int `json:"code"`
		Data struct {
			Items []models.PromptDTO `json:"items"` // 后端返回的是PromptDTO
			Total int                `json:"total"`
		} `json:"data"`
	}

	if err := json.Unmarshal(result, &response); err != nil {
		return nil, fmt.Errorf("failed to parse backend response (scope=%s): %w", scope, err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("backend returned error code: %d (scope=%s)", response.Code, scope)
	}

	// 转换PromptDTO为PromptMetadata
	metadata := make([]PromptMetadata, len(response.Data.Items))
	for i, item := range response.Data.Items {
		// 解析UpdatedAt时间戳
		var updatedAt int64 = 0
		if item.UpdatedAt != "" {
			if t, err := time.Parse(time.RFC3339, item.UpdatedAt); err == nil {
				updatedAt = t.Unix()
			} else {
				log.Printf("[MCP] Warning: could not parse updated_at '%s': %v", item.UpdatedAt, err)
			}
		}

		// 转换Arguments
		args := make([]PromptArgument, len(item.Arguments))
		for j, arg := range item.Arguments {
			args[j] = PromptArgument{
				Name:        arg.Name,
				Description: arg.Description,
				Required:    arg.Required,
			}
		}

		metadata[i] = PromptMetadata{
			Name:        item.Name,
			Description: item.Description,
			Version:     item.Version,
			Arguments:   args,
			UpdatedAt:   updatedAt,
		}
	}

	return metadata, nil
}

// PromptsGetParams prompts/get请求参数
type PromptsGetParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	UserID    uint                   `json:"user_id,omitempty"`    // 用于查询个人/全员提示词
	UserRole  string                 `json:"user_role,omitempty"`  // 用于权限检查
	ProjectID uint                   `json:"project_id,omitempty"` // 用于查询项目全员提示词
	APIToken  string                 `json:"api_token,omitempty"`  // API Token
}

// HandlePromptsGet 处理prompts/get请求
func HandlePromptsGet(
	ctx context.Context,
	registry *PromptsRegistry,
	backendClient *client.BackendClient,
	params json.RawMessage,
) (interface{}, error) {
	// 尝试从context中提取api_token
	if apiToken, ok := ctx.Value(client.TokenContextKey).(string); ok && apiToken != "" {
		log.Printf("[MCP] prompts/get: Extracted api_token from context")
	}

	// 解析参数
	var getParams PromptsGetParams
	if err := json.Unmarshal(params, &getParams); err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	if getParams.Name == "" {
		return nil, fmt.Errorf("missing required parameter: name")
	}

	log.Printf("[MCP] prompts/get called for: %s", getParams.Name)

	log.Printf("[MCP] prompts/get called for: %s", getParams.Name)

	// 优先查找系统Prompt
	prompt, found := registry.Get(getParams.Name)
	if found {
		log.Printf("[MCP] Found system prompt: %s", getParams.Name)
		// 获取Prompt内容
		content, err := prompt.GetContent()
		if err != nil {
			return nil, fmt.Errorf("failed to load prompt content: %w", err)
		}

		// 如果提供了arguments，执行模板替换
		if len(getParams.Arguments) > 0 {
			content = applyArguments(content, getParams.Arguments)
		}

		// 构造符合MCP规范的响应
		response := map[string]interface{}{
			"description": prompt.Description,
			"messages": []map[string]interface{}{
				{
					"role": "user",
					"content": map[string]interface{}{
						"type": "text",
						"text": content,
					},
				},
			},
		}

		return response, nil
	}

	log.Printf("[MCP] Not a system prompt, checking custom prompts...")

	// 如果不是系统Prompt，尝试从token中提取userID
	userIDForQuery := getParams.UserID
	hasToken := false
	if userIDForQuery == 0 {
		// 尝试从context中获取token
		if apiToken, ok := ctx.Value(client.TokenContextKey).(string); ok && apiToken != "" {
			hasToken = true
			log.Printf("[MCP] Attempting to extract userID from token for prompt: %s", getParams.Name)
			userID, err := getUserIDFromToken(ctx, backendClient, apiToken)
			if err != nil {
				log.Printf("[MCP] Failed to extract userID: %v", err)
			} else if userID > 0 {
				userIDForQuery = userID
				log.Printf("[MCP] Extracted userID=%d from token", userID)
			}
		}
	}

	// 通过BackendClient查询数据库
	// 全员提示词（project scope）不需要userID也可以查询
	// 个人提示词（user scope）需要userID
	if backendClient != nil {
		log.Printf("[MCP] Querying backend for custom prompt: %s (userID=%d, hasToken=%v)", getParams.Name, userIDForQuery, hasToken)

		reqCtx := ctx
		if apiToken, ok := ctx.Value(client.TokenContextKey).(string); ok && apiToken != "" {
			reqCtx = context.WithValue(ctx, client.TokenContextKey, apiToken)
		}

		customPromptContent, err := getCustomPromptFromBackend(
			reqCtx,
			backendClient,
			getParams.Name,
			userIDForQuery,
			getParams.UserRole,
			getParams.ProjectID,
		)
		if err == nil && customPromptContent != "" {
			log.Printf("[MCP] Found custom prompt: %s", getParams.Name)
			// 如果提供了arguments，执行模板替换
			if len(getParams.Arguments) > 0 {
				customPromptContent = applyArguments(customPromptContent, getParams.Arguments)
			}

			// 构造符合MCP规范的响应
			response := map[string]interface{}{
				"description": fmt.Sprintf("Custom prompt: %s", getParams.Name),
				"messages": []map[string]interface{}{
					{
						"role": "user",
						"content": map[string]interface{}{
							"type": "text",
							"text": customPromptContent,
						},
					},
				},
			}

			return response, nil
		} else if err != nil {
			log.Printf("[MCP] Error fetching custom prompt: %v", err)
		}
	} else {
		log.Printf("[MCP] Cannot query custom prompt: backendClient is nil")
	}

	log.Printf("[MCP] Prompt not found: %s", getParams.Name)
	return nil, fmt.Errorf("prompt not found: %s", getParams.Name)
}

// getCustomPromptFromBackend 从后端API获取自定义提示词内容
func getCustomPromptFromBackend(
	ctx context.Context,
	backendClient *client.BackendClient,
	promptName string,
	userID uint,
	userRole string,
	projectID uint,
) (string, error) {
	// 调用后端API获取提示词详情
	params := map[string]string{
		"name": promptName,
	}

	log.Printf("[MCP] Calling backend API: /api/v1/prompts/by-name with name=%s, userID=%d", promptName, userID)

	// 调用后端API获取提示词详情
	result, err := backendClient.Get(ctx, "/api/v1/prompts/by-name", params)
	if err != nil {
		return "", fmt.Errorf("failed to fetch custom prompt from backend: %w", err)
	}

	// 解析返回的数据
	var response struct {
		Code int `json:"code"`
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}

	if err := json.Unmarshal(result, &response); err != nil {
		return "", fmt.Errorf("failed to parse backend response: %w", err)
	}

	if response.Code != 0 {
		return "", fmt.Errorf("backend returned error code: %d", response.Code)
	}

	return response.Data.Content, nil
}

// applyArguments 执行简单的模板参数替换
func applyArguments(content string, args map[string]interface{}) string {
	result := content
	for key, value := range args {
		placeholder := fmt.Sprintf("{{%s}}", key)
		valueStr := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, valueStr)
	}
	return result
}
