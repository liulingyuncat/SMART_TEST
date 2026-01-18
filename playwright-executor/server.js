const express = require('express');
const { chromium } = require('playwright');

const app = express();
const PORT = process.env.PORT || 53730;
const PLAYWRIGHT_WS = process.env.PLAYWRIGHT_WS || 'ws://playwright-runner:53729/';

app.use(express.json({ limit: '10mb' }));

// 健康检查
app.get('/health', (req, res) => {
  res.json({ 
    status: 'ok', 
    playwright_ws: PLAYWRIGHT_WS
  });
});

// 执行脚本
app.post('/execute', async (req, res) => {
  const { scriptCode, timeout = 60000 } = req.body;
  
  if (!scriptCode) {
    return res.status(400).json({ 
      success: false, 
      error: 'scriptCode is required' 
    });
  }

  console.log('[Executor] 接收到执行请求');
  console.log('[Executor] 脚本长度:', scriptCode.length);
  console.log('[Executor] Playwright Server:', PLAYWRIGHT_WS);
  
  const startTime = Date.now();
  let browser = null;
  let context = null;
  let page = null;

  try {
    // 连接到远程 Playwright Server
    console.log('[Executor] 连接到 Playwright Server...');
    browser = await chromium.connect(PLAYWRIGHT_WS, { timeout });
    
    // 创建浏览器上下文和页面（忽略 HTTPS 证书错误）
    context = await browser.newContext({
      ignoreHTTPSErrors: true  // 跳过自签名证书验证
    });
    page = await context.newPage();
    
    console.log('[Executor] 开始执行脚本...');
    
    // 解析用户脚本（格式: async (page) => { ... }）
    let userFunction;
    try {
      // 使用 eval 解析函数
      userFunction = eval(`(${scriptCode})`);
      
      if (typeof userFunction !== 'function') {
        throw new Error('scriptCode must be a function');
      }
    } catch (error) {
      throw new Error(`Failed to parse script: ${error.message}`);
    }
    
    // 执行用户脚本
    const result = await userFunction(page);
    
    const responseTime = Date.now() - startTime;
    console.log('[Executor] 执行成功，耗时:', responseTime, 'ms');
    
    res.json({
      success: true,
      output: result !== undefined ? JSON.stringify(result) : 'Script executed successfully',
      responseTime
    });
    
  } catch (error) {
    const responseTime = Date.now() - startTime;
    console.error('[Executor] 执行失败:', error.message);
    console.error('[Executor] 错误堆栈:', error.stack);
    
    res.json({
      success: false,
      error: error.message,
      stack: error.stack,
      responseTime
    });
    
  } finally {
    // 清理资源
    try {
      if (page) await page.close();
      if (context) await context.close();
      if (browser) await browser.close();
    } catch (cleanupError) {
      console.error('[Executor] 清理资源失败:', cleanupError.message);
    }
  }
});

// 启动服务器
app.listen(PORT, () => {
  console.log(`[Executor] 服务启动在端口 ${PORT}`);
  console.log(`[Executor] Playwright Server: ${PLAYWRIGHT_WS}`);
});

// 优雅关闭
process.on('SIGTERM', () => {
  console.log('[Executor] 收到 SIGTERM 信号，正在关闭...');
  process.exit(0);
});

process.on('SIGINT', () => {
  console.log('[Executor] 收到 SIGINT 信号，正在关闭...');
  process.exit(0);
});
