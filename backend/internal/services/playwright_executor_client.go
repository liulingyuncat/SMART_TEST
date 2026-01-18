package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// PlaywrightExecutorConfig 执行器配置
type PlaywrightExecutorConfig struct {
	ExecutorURL    string        // 执行器服务地址，默认 http://playwright-executor:3001
	ExecuteTimeout time.Duration // 执行超时，默认 60s
	MaxRetries     int           // 最大重试次数，默认 3
}

// DefaultExecutorConfig 返回默认配置
func DefaultExecutorConfig() PlaywrightExecutorConfig {
	executorURL := os.Getenv("PLAYWRIGHT_EXECUTOR_URL")
	if executorURL == "" {
		executorURL = "http://playwright-executor:53730"
	}

	return PlaywrightExecutorConfig{
		ExecutorURL:    executorURL,
		ExecuteTimeout: 60 * time.Second,
		MaxRetries:     3,
	}
}

// PlaywrightExecutorClient 基于 HTTP 的 Playwright 执行器客户端
type PlaywrightExecutorClient struct {
	config     PlaywrightExecutorConfig
	httpClient *http.Client
}

// NewPlaywrightExecutorClient 创建客户端实例
func NewPlaywrightExecutorClient(config PlaywrightExecutorConfig) *PlaywrightExecutorClient {
	return &PlaywrightExecutorClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.ExecuteTimeout + 10*time.Second, // 额外 10 秒作为 HTTP 超时
		},
	}
}

// ExecuteRequest 执行请求
type ExecuteRequest struct {
	ScriptCode string `json:"scriptCode"`
	Timeout    int    `json:"timeout,omitempty"` // 毫秒
}

// ExecuteResponse 执行响应
type ExecuteResponse struct {
	Success      bool   `json:"success"`
	Output       string `json:"output,omitempty"`
	Error        string `json:"error,omitempty"`
	Stack        string `json:"stack,omitempty"`
	ResponseTime int    `json:"responseTime"`
}

// Execute 执行 Playwright 脚本
func (c *PlaywrightExecutorClient) Execute(ctx context.Context, scriptCode string) (*DockerExecResult, error) {
	fmt.Printf("[ExecutorClient] 开始执行脚本，长度: %d 字节\n", len(scriptCode))

	var lastErr error
	startTime := time.Now()

	for i := 0; i < c.config.MaxRetries; i++ {
		if i > 0 {
			fmt.Printf("[ExecutorClient] 重试 %d/%d...\n", i, c.config.MaxRetries)
		}

		result, err := c.doExecute(ctx, scriptCode)
		if err == nil {
			responseTime := int(time.Since(startTime).Milliseconds())
			result.ResponseTime = responseTime
			fmt.Printf("[ExecutorClient] 执行成功，耗时: %dms\n", result.ResponseTime)
			return result, nil
		}
		lastErr = err
		fmt.Printf("[ExecutorClient] 执行失败: %v\n", err)
	}

	responseTime := int(time.Since(startTime).Milliseconds())
	return &DockerExecResult{
		Success:      false,
		Output:       lastErr.Error(),
		ResponseTime: responseTime,
	}, fmt.Errorf("execute failed after %d retries: %w", c.config.MaxRetries, lastErr)
}

// doExecute 执行脚本的实际逻辑
func (c *PlaywrightExecutorClient) doExecute(ctx context.Context, scriptCode string) (*DockerExecResult, error) {
	// 构造请求
	reqBody := ExecuteRequest{
		ScriptCode: scriptCode,
		Timeout:    int(c.config.ExecuteTimeout.Milliseconds()),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// 发送 HTTP 请求
	url := c.config.ExecutorURL + "/execute"
	fmt.Printf("[ExecutorClient] 发送请求到: %s\n", url)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// 解析响应
	var execResp ExecuteResponse
	if err := json.Unmarshal(body, &execResp); err != nil {
		return nil, fmt.Errorf("parse response: %w (body: %s)", err, string(body))
	}

	// 检查执行结果
	if !execResp.Success {
		errorMsg := execResp.Error
		if errorMsg == "" {
			errorMsg = "Unknown error"
		}
		return nil, fmt.Errorf("script execution failed: %s", errorMsg)
	}

	return &DockerExecResult{
		Success:      true,
		Output:       execResp.Output,
		ResponseTime: execResp.ResponseTime,
	}, nil
}

// Close 关闭客户端
func (c *PlaywrightExecutorClient) Close() error {
	// HTTP client 不需要显式关闭
	return nil
}
