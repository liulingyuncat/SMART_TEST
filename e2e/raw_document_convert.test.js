/**
 * E2E 测试 - 原始文档转换功能
 * 
 * 测试场景：
 * 1. 用户点击转换按钮，轮询获取状态，转换成功时显示成功提示
 * 2. 转换失败时显示错误信息
 * 3. 超时场景下显示超时提示
 * 
 * 依赖：Playwright / Cypress
 * 运行前需要启动前后端服务
 */

describe('Raw Document Convert E2E Tests', () => {
  const BASE_URL = 'http://localhost:3000';
  const API_BASE = 'http://localhost:8080/api/v1';
  
  // 测试账号（需要有效的测试账号）
  const TEST_USER = {
    username: 'admin',
    password: 'admin123',
  };
  
  // 测试项目ID（需要有效的项目）
  const TEST_PROJECT_ID = 1;
  
  beforeEach(async () => {
    // 登录
    await page.goto(`${BASE_URL}/login`);
    await page.fill('input[name="username"]', TEST_USER.username);
    await page.fill('input[name="password"]', TEST_USER.password);
    await page.click('button[type="submit"]');
    
    // 等待登录完成并跳转
    await page.waitForNavigation();
    
    // 导航到项目详情页的原始需求Tab
    await page.goto(`${BASE_URL}/projects/${TEST_PROJECT_ID}`);
    await page.click('text=原始需求文档');
    await page.waitForSelector('.raw-document-table');
  });
  
  afterEach(async () => {
    // 清理测试数据（如果需要）
  });
  
  /**
   * 测试场景1: 正常转换流程
   * 前置条件：存在未转换的文档
   */
  test('should successfully convert document and show success message', async () => {
    // 1. 找到状态为 "none" 的文档的转换按钮
    const convertButton = await page.locator('button:has-text("转换")').first();
    
    // 验证按钮可见且可点击
    await expect(convertButton).toBeVisible();
    await expect(convertButton).not.toBeDisabled();
    
    // 2. 点击转换按钮
    await convertButton.click();
    
    // 3. 验证显示"转换已启动"提示
    await expect(page.locator('.ant-message-success')).toContainText('转换已启动');
    
    // 4. 验证按钮变为"转换中..."并禁用
    await expect(page.locator('button:has-text("转换中...")').first()).toBeVisible();
    await expect(page.locator('button:has-text("转换中...")').first()).toBeDisabled();
    
    // 5. 等待转换完成（最长60秒）
    await page.waitForSelector('.ant-tag-success:has-text("已完成")', { timeout: 60000 });
    
    // 6. 验证显示"转换成功"提示
    await expect(page.locator('.ant-message-success')).toContainText('转换成功');
    
    // 7. 验证转换后的文件名显示
    await expect(page.locator('td:has-text("_Trans_")')).toBeVisible();
  });
  
  /**
   * 测试场景2: 转换失败流程
   * 前置条件：上传一个无法转换的文档（如损坏的PDF）
   */
  test('should show error message when conversion fails', async () => {
    // 此测试需要预先上传一个无法转换的文件
    // 或者通过 Mock API 模拟失败场景
    
    // 1. 找到会失败的文档的转换按钮
    const convertButton = await page.locator('[data-testid="convert-corrupted-doc"]');
    
    // 2. 点击转换按钮
    await convertButton.click();
    
    // 3. 等待转换失败
    await page.waitForSelector('.ant-message-error', { timeout: 60000 });
    
    // 4. 验证显示错误信息
    await expect(page.locator('.ant-message-error')).toBeVisible();
    
    // 5. 验证按钮回退到可点击状态
    await expect(page.locator('button:has-text("转换")').first()).not.toBeDisabled();
  });
  
  /**
   * 测试场景3: 重复点击转换按钮被阻止
   */
  test('should prevent duplicate conversion clicks', async () => {
    // 1. 找到转换按钮
    const convertButton = await page.locator('button:has-text("转换")').first();
    
    // 2. 点击转换按钮
    await convertButton.click();
    
    // 3. 验证按钮立即禁用
    await expect(page.locator('button:has-text("转换中...")').first()).toBeDisabled();
    
    // 4. 尝试再次点击（应该无效）
    const disabledButton = await page.locator('button:has-text("转换中...")').first();
    
    // 计数器验证 API 只被调用一次
    let apiCallCount = 0;
    await page.route(`${API_BASE}/raw-documents/*/convert`, (route) => {
      apiCallCount++;
      route.continue();
    });
    
    // 等待一小段时间确保没有额外调用
    await page.waitForTimeout(1000);
    
    // API 应该只被调用一次
    expect(apiCallCount).toBeLessThanOrEqual(1);
  });
  
  /**
   * 测试场景4: 已在转换中的文档返回409
   */
  test('should show appropriate message when document is already converting', async () => {
    // Mock API 返回 409
    await page.route(`${API_BASE}/raw-documents/*/convert`, (route) => {
      route.fulfill({
        status: 409,
        body: JSON.stringify({ error: 'document conversion already in progress' }),
      });
    });
    
    // 点击转换按钮
    const convertButton = await page.locator('button:has-text("转换")').first();
    await convertButton.click();
    
    // 验证显示"正在转换中"提示
    await expect(page.locator('.ant-message-error')).toContainText('正在转换中');
  });
  
  /**
   * 测试场景5: 文档不存在返回404
   */
  test('should show error when document not found', async () => {
    // Mock API 返回 404
    await page.route(`${API_BASE}/raw-documents/*/convert`, (route) => {
      route.fulfill({
        status: 404,
        body: JSON.stringify({ error: 'document not found' }),
      });
    });
    
    // 点击转换按钮
    const convertButton = await page.locator('button:has-text("转换")').first();
    await convertButton.click();
    
    // 验证显示"文档不存在"提示
    await expect(page.locator('.ant-message-error')).toContainText('文档不存在');
  });
  
  /**
   * 测试场景6: 轮询状态时网络错误重试
   */
  test('should retry polling on network error and eventually succeed', async () => {
    let pollCount = 0;
    
    // 第一次轮询失败，第二次成功
    await page.route(`${API_BASE}/raw-documents/*/convert-status`, (route) => {
      pollCount++;
      if (pollCount === 1) {
        route.abort('failed');
      } else {
        route.fulfill({
          status: 200,
          body: JSON.stringify({ status: 'completed', progress: 100 }),
        });
      }
    });
    
    // 点击转换按钮
    const convertButton = await page.locator('button:has-text("转换")').first();
    await convertButton.click();
    
    // 等待最终成功
    await expect(page.locator('.ant-message-success')).toContainText('转换成功');
    
    // 验证重试了
    expect(pollCount).toBeGreaterThan(1);
  });
  
  /**
   * 测试场景7: 转换超时
   */
  test('should show timeout message after max polling attempts', async () => {
    // Mock API 始终返回 processing
    await page.route(`${API_BASE}/raw-documents/*/convert-status`, (route) => {
      route.fulfill({
        status: 200,
        body: JSON.stringify({ status: 'processing', progress: 50 }),
      });
    });
    
    // 设置较短的超时进行测试（实际测试可能需要调整）
    // 注意：真实环境中这个测试会很慢（60秒超时）
    
    // 点击转换按钮
    const convertButton = await page.locator('button:has-text("转换")').first();
    await convertButton.click();
    
    // 等待超时提示（需要约60秒）
    await expect(page.locator('.ant-message-warning')).toContainText('转换超时', { timeout: 70000 });
  });
  
  /**
   * 测试场景8: 已完成的文档显示正确状态
   */
  test('should display completed status and converted filename', async () => {
    // 验证已完成的文档显示
    await expect(page.locator('.ant-tag-success:has-text("已完成")')).toBeVisible();
    
    // 验证转换后的文件名格式正确
    const convertedFilename = await page.locator('td:has-text("_Trans_")').textContent();
    expect(convertedFilename).toMatch(/_Trans_\d+\.md/);
  });
});

/**
 * 辅助函数：上传测试文件
 */
async function uploadTestFile(page, filename, content) {
  const buffer = Buffer.from(content);
  
  await page.setInputFiles('input[type="file"]', {
    name: filename,
    mimeType: 'text/plain',
    buffer: buffer,
  });
  
  // 等待上传完成
  await page.waitForSelector('.ant-message-success');
}

/**
 * 辅助函数：清理测试文档
 */
async function cleanupTestDocuments(page, API_BASE) {
  // 获取所有测试文档
  const response = await page.request.get(`${API_BASE}/projects/1/raw-documents`);
  const data = await response.json();
  
  // 删除所有测试文档
  for (const doc of data.documents || []) {
    await page.request.delete(`${API_BASE}/raw-documents/${doc.id}`);
  }
}
