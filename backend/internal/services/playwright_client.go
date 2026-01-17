package services

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

// PlaywrightClientConfig 客户端配置
type PlaywrightClientConfig struct {
	WSEndpoint     string        // WebSocket 地址，默认 ws://playwright-runner:3000/
	ConnectTimeout time.Duration // 连接超时，默认 10s
	ExecuteTimeout time.Duration // 执行超时，默认 60s
	MaxRetries     int           // 最大重试次数，默认 3
}

// DefaultPlaywrightConfig 返回默认配置
func DefaultPlaywrightConfig() PlaywrightClientConfig {
	wsEndpoint := os.Getenv("PLAYWRIGHT_WS_ENDPOINT")
	if wsEndpoint == "" {
		wsEndpoint = "ws://playwright-runner:3000/"
	}

	return PlaywrightClientConfig{
		WSEndpoint:     wsEndpoint,
		ConnectTimeout: 10 * time.Second,
		ExecuteTimeout: 60 * time.Second,
		MaxRetries:     3,
	}
}

// PlaywrightClient Playwright WebSocket 客户端
type PlaywrightClient struct {
	config PlaywrightClientConfig
	pw     *playwright.Playwright
	mu     sync.Mutex
}

// NewPlaywrightClient 创建客户端实例
func NewPlaywrightClient(config PlaywrightClientConfig) *PlaywrightClient {
	return &PlaywrightClient{
		config: config,
	}
}

// ensurePlaywright 确保 Playwright 实例已启动
// 需要本地 SDK 驱动来建立 WebSocket 协议通信
func (c *PlaywrightClient) ensurePlaywright() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pw != nil {
		return nil
	}

	// 启动本地 Playwright SDK
	// 注意：playwright-go 会自动下载所需的驱动文件到 ~/.cache/ms-playwright-go/
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("start playwright: %w", err)
	}

	c.pw = pw
	fmt.Printf("[PlaywrightClient] Playwright SDK 已启动，将连接到远程服务器: %s\n", c.config.WSEndpoint)
	return nil
}

// ExecuteScript 执行脚本
// 连接到远程 Playwright Server，创建浏览器上下文，执行脚本并返回结果
func (c *PlaywrightClient) ExecuteScript(ctx context.Context, scriptCode string) (*DockerExecResult, error) {
	fmt.Printf("[PlaywrightClient] 开始执行脚本，长度: %d bytes\n", len(scriptCode))
	fmt.Printf("[PlaywrightClient] WebSocket 端点: %s\n", c.config.WSEndpoint)

	startTime := time.Now()

	// 确保 Playwright 实例已启动
	if err := c.ensurePlaywright(); err != nil {
		return nil, fmt.Errorf("ensure playwright: %w", err)
	}

	// 带超时的上下文
	execCtx, cancel := context.WithTimeout(ctx, c.config.ExecuteTimeout)
	defer cancel()

	// 重试机制
	var lastErr error
	for i := 0; i < c.config.MaxRetries; i++ {
		if i > 0 {
			fmt.Printf("[PlaywrightClient] 重试第 %d 次...\n", i)
			time.Sleep(time.Second)
		}

		result, err := c.doExecute(execCtx, scriptCode)
		if err == nil {
			result.ResponseTime = int(time.Since(startTime).Milliseconds())
			fmt.Printf("[PlaywrightClient] 执行成功，耗时: %dms\n", result.ResponseTime)
			return result, nil
		}
		lastErr = err
		fmt.Printf("[PlaywrightClient] 执行失败: %v\n", err)
	}

	responseTime := int(time.Since(startTime).Milliseconds())
	return &DockerExecResult{
		Success:      false,
		Output:       lastErr.Error(),
		ResponseTime: responseTime,
	}, fmt.Errorf("execute failed after %d retries: %w", c.config.MaxRetries, lastErr)
}

// doExecute 执行脚本的实际逻辑
func (c *PlaywrightClient) doExecute(ctx context.Context, scriptCode string) (*DockerExecResult, error) {
	// 连接到远程 Playwright Server
	fmt.Printf("[PlaywrightClient] 连接到 Playwright Server...\n")
	browser, err := c.pw.Chromium.Connect(c.config.WSEndpoint)
	if err != nil {
		return nil, fmt.Errorf("connect to playwright server: %w", err)
	}
	defer browser.Close()

	// 创建浏览器上下文
	browserContext, err := browser.NewContext()
	if err != nil {
		return nil, fmt.Errorf("create browser context: %w", err)
	}
	defer browserContext.Close()

	// 创建页面
	page, err := browserContext.NewPage()
	if err != nil {
		return nil, fmt.Errorf("create page: %w", err)
	}
	defer page.Close()

	// 包装脚本为可执行函数
	// 用户脚本通常是完整的 Playwright 脚本，我们需要在页面上下文中执行
	wrappedScript := fmt.Sprintf(`
		(async () => {
			try {
				%s
				return { success: true, output: 'Script executed successfully' };
			} catch (error) {
				return { success: false, output: error.message || error.toString() };
			}
		})()
	`, scriptCode)

	// 执行脚本
	fmt.Printf("[PlaywrightClient] 执行脚本...\n")
	result, err := page.Evaluate(wrappedScript)
	if err != nil {
		return nil, fmt.Errorf("evaluate script: %w", err)
	}

	// 解析结果
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return &DockerExecResult{
			Success: true,
			Output:  fmt.Sprintf("%v", result),
		}, nil
	}

	success, _ := resultMap["success"].(bool)
	output, _ := resultMap["output"].(string)

	if !success {
		return nil, fmt.Errorf("script error: %s", output)
	}

	return &DockerExecResult{
		Success: true,
		Output:  output,
	}, nil
}

// Close 关闭 Playwright 实例
func (c *PlaywrightClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pw != nil {
		if err := c.pw.Stop(); err != nil {
			return fmt.Errorf("stop playwright: %w", err)
		}
		c.pw = nil
	}
	return nil
}

// IsReady 检查客户端是否就绪
func (c *PlaywrightClient) IsReady() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pw != nil
}
