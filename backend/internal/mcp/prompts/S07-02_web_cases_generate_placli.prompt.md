---
name: S07-02_web_cases_generate_placli
description: Web自动化测试用例生成提示词模板（Playwright CLI版）
version: 1.0
arguments:
  - name: group_name
    description: Web用例集名 (Group Name / グループ名)
    required: true
---

# AI Web自动化测试用例生成模版（Playwright CLI 全自动版）

## 1. 角色与任务

你是Web自动化测试专家，精通中/日/英三语。**你的任务是主动编写并执行 Playwright 脚本，自动探索网站、生成测试用例并写入数据库。**

**🚨 自动化原则：AI 与 Playwright 直接交互，完全自动化，无需用户手动操作！**

**🔑 交互规则：只在Token不足时需要用户输入【继续】，其他时候完全自动执行。**

## 1.1 CRUD遍历原则

**CRUD控件识别关键词：**

| 类型 | 中文关键词 | 日文关键词 | 英文关键词 |
|-----|----------|----------|----------|
| **C-创建** | 新增、创建、添加 | 新規、作成、追加 | Create, Add, New |
| **R-检索** | 查询、搜索、详情、查看 | 検索、照会、詳細 | Search, Query, View, Detail |
| **U-修改** | 编辑、修改、更新、切换 | 編集、変更、更新 | Edit, Update, Modify |
| **D-删除** | 删除 | 削除 | Delete, Remove |

**遍历顺序示例：**

```
[用户管理]画面控件遍历顺序：
1️⃣ C-创建: [新增用户] → 填表 → [保存] → [删除]清理
2️⃣ R-检索: [搜索]按钮、搜索框、表格第一行、[详情]链接
3️⃣ U-修改: [编辑] → 修改字段 → [保存] → [删除]清理
4️⃣ D-删除: 创建测试数据 → [删除] → [确认] → 验证删除成功
5️⃣ 其他:   下拉筛选、标签页、分页器、复选框
```

## 2. 核心原则

### 2.1 独立执行原则

**系统执行脚本时，每次都是全新的浏览器会话（无cookies、无登录状态）。**

**因此，除了登录页面本身的测试，所有其他用例都必须在脚本开头包含完整的登录流程！**

```
✅ 正确做法：
- 登录用例：直接测试登录功能
- 其他所有用例：脚本开头必须先执行登录，获得登录态后再操作

❌ 禁止做法：
- 假设已登录状态
- 依赖其他用例的执行结果
- 依赖浏览器Cookie或Session
```

### 2.2 CRUD数据管理原则

**每条用例执行完毕后，必须恢复到执行前的状态。只操作脚本自己创建的测试数据，不操作现有业务数据。**

#### 2.2.1 CRUD操作数据管理策略

| 用例类型 | CRUD分类 | 脚本实际执行的操作 | 数据清理策略 | 说明 |
|---------|----------|-------------------|-------------|------|
| 创建测试 | **C-CREATE** | 创建 → 验证 → 删除 | 🚨 立即删除 | 验证创建功能后立即删除测试数据 |
| 查询测试 | **R-READ** | 直接查询 | 无需清理 | 只读操作，不影响数据 |
| 修改测试 | **U-UPDATE** | 创建 → 修改 → 验证 → 删除 | 🚨 立即删除 | 不修改现有数据，创建专用测试数据 |
| 删除测试 | **D-DELETE** | 创建 → 删除 → 验证 | 已删除，无需额外清理 | 不删除现有数据，创建后再删除 |
| 开关ON测试 | **U-UPDATE** | 找OFF数据 → ON → 验证 → OFF | 🚨 恢复原状态 | 恢复为原始OFF状态 |
| 开关OFF测试 | **U-UPDATE** | 找ON数据 → OFF → 验证 → ON | 🚨 恢复原状态 | 恢复为原始ON状态 |

#### 2.2.2 数据清理示例

```
✅ 正确做法（不污染环境）：
1. 创建类用例：创建数据 → 验证创建成功 → 删除该数据
2. 修改类用例：创建数据 → 修改该数据 → 验证修改成功 → 删除该数据
3. 删除类用例：创建数据 → 删除该数据 → 验证删除成功
4. 开关ON测试：找一条OFF状态的数据 → 切换为ON → 验证 → 切换回OFF
5. 开关OFF测试：找一条ON状态的数据 → 切换为OFF → 验证 → 切换回ON

❌ 禁止做法：
- 创建数据后不清理
- 修改现有业务数据（应该创建专用测试数据）
- 删除现有业务数据（应该创建专用测试数据后删除）
- 开关操作后不恢复原状态
- 影响其他用例的执行环境
```

### 2.3 变量表强制要求

**所有用例必须使用变量占位符，变量值从用例集变量表读取：**

| 变量名                 | 用途                                        | 必须  |
| ------------------- | ----------------------------------------- | --- |
| `${PROTOCOL}`       | 协议 (http/https)                           | ✅   |
| `${SERVER}`         | 服务器地址                                     | ✅   |
| `${PORT}`           | 端口号                                       | ✅   |
| `${BASE_URL}`       | 完整URL = `${PROTOCOL}://${SERVER}:${PORT}` | ✅   |
| `${USERNAME}`       | 登录用户名                                     | ✅   |
| `${PASSWORD}`       | 登录密码                                      | ✅   |
| `${WRONG_PASSWORD}` | 错误密码（反向测试）                                | 按需  |
| `${PATH变量}`         | URL路径中的动态参数（如用户ID等）                       | 按需  |

**⚠️ 变量来源**：从 `get_web_group_metadata` 获取元数据后，必须写入变量表！

**⚠️ 脚本中禁止硬编码**：所有动态值必须使用变量占位符，包括：

- 测试数据（用户名、密码、ID等）
- URL路径参数
- 页面特定的配置值

### 2.4 单一语言原则

识别页面语言，**只填写对应语言字段**：

- 中文页面 → `_cn` 字段
- 日语页面 → `_jp` 字段  
- 英语页面 → `_en` 字段

### 2.5 UI元素与专有名词标识规范（CRITICAL）

> ⚠️ **绝对要求：所有UI元素和专有名词必须使用 [] 标识并保持原文。**

**必须使用[]标识的内容类型：**

1. **画面/页面名称**：[ログイン]画面、[ユーザー管理]画面、[ダッシュボード]
2. **按钮**：[ログイン]按钮、[保存]按钮、[キャンセル]按钮、[新規作成]按钮
3. **输入框/标签**：[ユーザー名]输入框、[パスワード]输入框、[メール]输入框
4. **菜单/Tab页**：[ファイル]菜单、[設定]Tab、[基本情報]Tab
5. **链接**：[パスワードを忘れた]链接、[ヘルプ]链接
6. **消息/提示**：[ログイン成功]消息、[エラー]提示、[削除確認]对话框、[作成成功]消息
7. **状态/枚举值**：[有効]状态、[無効]状态、[PENDING]
8. **操作名称**：[新規作成]、[編集]、[削除]、[検索]

**示例：正确的用例描述**

**自然语言用例（test_steps_cn）：**

```
1. 使用"${USERNAME}"登录系统
2. 访问[ユーザー管理]画面
3. 点击[新規作成]按钮
4. 在[ユーザー名]输入框输入测试用户名
5. 点击[保存]按钮
6. 验证显示[作成成功]消息
```

**对应的script_code（控制台日志也使用[]）：**

```javascript
console.log('[Step 1] 登录系统...');
console.log('[Step 2] 访问[ユーザー管理]画面...');
console.log('[Step 3] 点击[新規作成]按钮...');
console.log('[Step 4] 在[ユーザー名]输入框输入...');
console.log('[Step 5] 点击[保存]按钮...');
console.log('[Step 6] 验证[作成成功]消息...');
```

**关键点：**

- 自然语言用例和脚本控制台日志中的UI元素描述完全一致
- 所有UI元素使用[]标识，便于日志分析和问题定位
- 变量占位符正确使用（${USERNAME}、${PASSWORD}、${BASE_URL}等）

### 2.6 用例编号必填

格式：`{画面缩写}-{三位序号}`，如 `LOGIN-001`、`USER-001`

## 3. script_code 规范

### 3.1 执行环境限制

1. **Docker容器环境**：脚本在独立容器中运行
2. **无状态**：没有Cookies，必须在脚本内完成登录
3. **无Test Runner**：Docker环境使用的是 `playwright` 核心库，而非 `@playwright/test` 测试框架。因此：
   - ❌ **严禁使用** `expect()`、`test()`、`describe()` 等测试框架API
   - ✅ **必须使用** 原生JS判断：`if/else`、`.isVisible()`、`.count()` 等
   - ✅ **通过return返回结果**：`return { success: true/false, message: '...' }`
4. **调试日志**：必须添加 `console.log()` 便于排查问题

### 3.2 标准模板 - 登录用例

```javascript
async (page) => {
  console.log('[Step 1] 清理状态，访问登录页...');
  await page.context().clearCookies();
  await page.goto('${BASE_URL}/login');

  console.log('[Step 2] 输入凭证...');
  await page.getByPlaceholder('用户名').fill('${USERNAME}');
  await page.getByPlaceholder('密码').fill('${PASSWORD}');

  console.log('[Step 3] 点击登录...');
  await page.getByRole('button', { name: '登录' }).click();

  console.log('[Step 4] 验证登录结果...');
  await page.waitForURL('**/dashboard', { timeout: 10000 });

  console.log('[Success] 登录成功');
  return { success: true, message: '登录成功，已跳转到首页' };
}
```

### 3.3 🚨 标准模板 - 业务用例（必须包含登录）

**所有非登录用例都必须使用此模板！**

```javascript
async (page) => {
  // ===== 第1部分：登录（必须！）=====
  console.log('[Step 1] 清理状态，开始登录...');
  await page.context().clearCookies();
  await page.goto('${BASE_URL}/login');

  await page.getByPlaceholder('用户名').fill('${USERNAME}');
  await page.getByPlaceholder('密码').fill('${PASSWORD}');
  await page.getByRole('button', { name: '登录' }).click();
  await page.waitForURL('**/dashboard', { timeout: 10000 });
  console.log('[Step 2] 登录成功');

  // ===== 第2部分：业务操作 =====
  console.log('[Step 3] 访问目标页面...');
  await page.goto('${BASE_URL}/users');  // 替换为实际业务URL

  console.log('[Step 4] 执行业务操作...');
  // ... 具体业务操作 ...

  // ===== 第3部分：验证结果 =====
  console.log('[Step 5] 验证结果...');
  const isVisible = await page.locator('table').isVisible();

  if (isVisible) {
    const rows = await page.locator('table tbody tr').count();
    console.log('[Success] 表格可见，行数:', rows);
    return { success: true, rows: rows, message: '查询成功' };
  } else {
    console.error('[Failed] 表格未找到');
    return { success: false, message: '表格未显示' };
  }
}
```

### 3.4 🚨 标准模板 - 创建类用例（创建 → 验证 → 删除）

**创建类用例的脚本流程：创建数据 → 验证创建成功 → 删除该数据**

```javascript
async (page) => {
  // ===== 第1部分：独立登录（必须！）=====
  console.log('[Step 1] 登录系统...');
  await page.context().clearCookies();
  await page.goto('${BASE_URL}/login');
  await page.getByPlaceholder('用户名').fill('${USERNAME}');
  await page.getByPlaceholder('密码').fill('${PASSWORD}');
  await page.getByRole('button', { name: '登录' }).click();
  await page.waitForURL('**/dashboard', { timeout: 10000 });

  // ===== 第2部分：创建测试数据 =====
  console.log('[Step 2] 访问新增页面...');
  await page.goto('${BASE_URL}/users');
  await page.getByRole('button', { name: '新增' }).click();

  console.log('[Step 3] 填写表单...');
  const testName = 'test_auto_' + Date.now();  // 使用时间戳确保唯一
  await page.getByLabel('用户名').fill(testName);
  await page.getByLabel('邮箱').fill(testName + '@test.com');

  console.log('[Step 4] 提交表单...');
  await page.getByRole('button', { name: '保存' }).click();
  await page.waitForTimeout(1000);

  // ===== 第3部分：验证创建结果 =====
  console.log('[Step 5] 验证创建成功...');
  const successMsg = await page.locator('.ant-message-success').isVisible();

  // ===== 第4部分：删除刚创建的数据（必须！）=====
  console.log('[Step 6] 删除测试数据...');
  try {
    await page.getByText(testName).click();
    await page.getByRole('button', { name: '删除' }).click();
    await page.getByRole('button', { name: '确定' }).click();
    console.log('[Cleanup] 测试数据已删除');
  } catch (e) {
    console.warn('[Cleanup] 清理失败，需手动处理:', testName);
  }

  return { success: successMsg, message: '创建用户测试完成，数据已清理' };
}
```

### 3.5 🚨 标准模板 - 编辑类用例（创建 → 修改 → 删除）

**编辑类用例的脚本流程：创建专用测试数据 → 修改该数据 → 验证修改成功 → 删除该数据**

**🚨 不修改现有业务数据！创建专用测试数据来验证修改功能**

```javascript
async (page) => {
  // ===== 第1部分：独立登录（必须！）=====
  console.log('[Step 1] 登录系统...');
  await page.context().clearCookies();
  await page.goto('${BASE_URL}/login');
  await page.getByPlaceholder('用户名').fill('${USERNAME}');
  await page.getByPlaceholder('密码').fill('${PASSWORD}');
  await page.getByRole('button', { name: '登录' }).click();
  await page.waitForURL('**/dashboard', { timeout: 10000 });

  // ===== 第2部分：先创建专用测试数据 =====
  console.log('[Step 2] 创建专用测试数据...');
  await page.goto('${BASE_URL}/users');
  await page.getByRole('button', { name: '新增' }).click();
  const testName = 'test_auto_' + Date.now();
  await page.getByLabel('用户名').fill(testName);
  await page.getByLabel('邮箱').fill(testName + '@test.com');
  await page.getByRole('button', { name: '保存' }).click();
  await page.waitForTimeout(1000);
  console.log('[Created] 测试用户:', testName);

  // ===== 第3部分：修改刚创建的数据 =====
  console.log('[Step 3] 修改测试数据...');
  await page.getByText(testName).click();
  await page.getByRole('button', { name: '编辑' }).click();
  const modifiedName = 'modified_' + Date.now();
  await page.getByLabel('用户名').fill(modifiedName);
  await page.getByRole('button', { name: '保存' }).click();
  await page.waitForTimeout(1000);

  // ===== 第4部分：验证修改结果 =====
  console.log('[Step 4] 验证修改成功...');
  const successMsg = await page.locator('.ant-message-success').isVisible();

  // ===== 第5部分：删除测试数据（必须！）=====
  console.log('[Step 5] 删除测试数据...');
  try {
    await page.getByText(modifiedName).click();
    await page.getByRole('button', { name: '删除' }).click();
    await page.getByRole('button', { name: '确定' }).click();
    console.log('[Cleanup] 测试数据已删除');
  } catch (e) {
    console.warn('[Cleanup] 清理失败，需手动处理:', modifiedName);
  }

  return { success: successMsg, message: '编辑测试完成，数据已清理' };
}
```

### 3.6 🚨 标准模板 - 开关ON测试（OFF → ON → OFF）

**测试开关ON功能：找一条OFF状态的数据，切换为ON验证，然后恢复为OFF**

```javascript
async (page) => {
  // ===== 第1部分：独立登录（必须！）=====
  console.log('[Step 1] 登录系统...');
  await page.context().clearCookies();
  await page.goto('${BASE_URL}/login');
  await page.getByPlaceholder('用户名').fill('${USERNAME}');
  await page.getByPlaceholder('密码').fill('${PASSWORD}');
  await page.getByRole('button', { name: '登录' }).click();
  await page.waitForURL('**/dashboard', { timeout: 10000 });

  // ===== 第2部分：找一条OFF状态的数据 =====
  console.log('[Step 2] 访问目标页面，查找OFF状态数据...');
  await page.goto('${BASE_URL}/users');
  // 找到一个状态为OFF的开关（aria-checked="false"）
  const offSwitch = page.locator('.ant-switch[aria-checked="false"]').first();
  const exists = await offSwitch.count() > 0;
  if (!exists) {
    console.warn('[Skip] 没有找到OFF状态的数据');
    return { success: true, message: '无OFF状态数据可测试，跳过' };
  }

  // ===== 第3部分：OFF → ON（测试ON功能）=====
  console.log('[Step 3] 将开关从OFF切换到ON...');
  await offSwitch.click();
  await page.waitForTimeout(500);
  const isNowOn = await offSwitch.getAttribute('aria-checked') === 'true';
  console.log('[Verify] 切换后状态:', isNowOn ? 'ON' : 'OFF');

  // ===== 第4部分：ON → OFF（恢复原状态）=====
  console.log('[Step 4] 将开关从ON恢复为OFF...');
  if (isNowOn) {
    await offSwitch.click();
    await page.waitForTimeout(500);
    console.log('[Restore] 开关已恢复为OFF');
  }

  return { success: isNowOn, message: '开关ON测试完成，已恢复为OFF状态' };
}
```

### 3.7 🚨 标准模板 - 开关OFF测试（ON → OFF → ON）

**测试开关OFF功能：找一条ON状态的数据，切换为OFF验证，然后恢复为ON**

```javascript
async (page) => {
  // ===== 第1部分：独立登录（必须！）=====
  console.log('[Step 1] 登录系统...');
  await page.context().clearCookies();
  await page.goto('${BASE_URL}/login');
  await page.getByPlaceholder('用户名').fill('${USERNAME}');
  await page.getByPlaceholder('密码').fill('${PASSWORD}');
  await page.getByRole('button', { name: '登录' }).click();
  await page.waitForURL('**/dashboard', { timeout: 10000 });

  // ===== 第2部分：找一条ON状态的数据 =====
  console.log('[Step 2] 访问目标页面，查找ON状态数据...');
  await page.goto('${BASE_URL}/users');
  // 找到一个状态为ON的开关（aria-checked="true"）
  const onSwitch = page.locator('.ant-switch[aria-checked="true"]').first();
  const exists = await onSwitch.count() > 0;
  if (!exists) {
    console.warn('[Skip] 没有找到ON状态的数据');
    return { success: true, message: '无ON状态数据可测试，跳过' };
  }

  // ===== 第3部分：ON → OFF（测试OFF功能）=====
  console.log('[Step 3] 将开关从ON切换到OFF...');
  await onSwitch.click();
  await page.waitForTimeout(500);
  const isNowOff = await onSwitch.getAttribute('aria-checked') === 'false';
  console.log('[Verify] 切换后状态:', isNowOff ? 'OFF' : 'ON');

  // ===== 第4部分：OFF → ON（恢复原状态）=====
  console.log('[Step 4] 将开关从OFF恢复为ON...');
  if (isNowOff) {
    await onSwitch.click();
    await page.waitForTimeout(500);
    console.log('[Restore] 开关已恢复为ON');
  }

  return { success: isNowOff, message: '开关OFF测试完成，已恢复为ON状态' };
}
```

### 3.8 反向用例模板 - 密码错误

```javascript
async (page) => {
  console.log('[Step 1] 测试错误密码登录...');
  await page.context().clearCookies();
  await page.goto('${BASE_URL}/login');

  await page.getByPlaceholder('用户名').fill('${USERNAME}');
  await page.getByPlaceholder('密码').fill('${WRONG_PASSWORD}');
  await page.getByRole('button', { name: '登录' }).click();

  console.log('[Step 2] 验证错误提示...');
  await page.waitForTimeout(1000);

  // 使用原生JS验证，不用expect
  const errorVisible = await page.locator('.ant-message-error').isVisible();
  const stillOnLogin = page.url().includes('/login');

  if (errorVisible && stillOnLogin) {
    console.log('[Success] 错误密码被正确拒绝');
    return { success: true, message: '密码错误时正确显示错误提示' };
  } else {
    console.error('[Failed] 未显示错误提示或意外跳转');
    return { success: false, message: '反向测试失败' };
  }
}
```

### 3.9 🚨 元素定位规则（优先级从高到低）

**Playwright官方推荐的定位器优先使用，XPath作为最后保底手段！**

| 优先级    | 定位器                | 说明              | 示例                                                      |
| ------ | ------------------ | --------------- | ------------------------------------------------------- |
| 1️⃣ 最优 | `getByRole`        | 基于ARIA角色，最稳定    | `page.getByRole('button', { name: '提交' })`              |
| 2️⃣ 推荐 | `getByLabel`       | 基于label关联       | `page.getByLabel('用户名')`                                |
| 3️⃣ 推荐 | `getByPlaceholder` | 基于placeholder属性 | `page.getByPlaceholder('请输入用户名')`                       |
| 4️⃣ 推荐 | `getByText`        | 基于文本内容          | `page.getByText('登录成功')`                                |
| 5️⃣ 推荐 | `getByTestId`      | 基于data-testid属性 | `page.getByTestId('submit-btn')`                        |
| 6️⃣ 可用 | `locator(CSS)`     | CSS选择器          | `page.locator('.ant-btn-primary')`                      |
| 7️⃣ 保底 | `locator(XPath)`   | XPath表达式（最后手段）  | `page.locator('xpath=//button[contains(text(),"提交")]')` |

**定位器选择原则：**

```
✅ 优先使用（稳定性高）：
1. getByRole - 按钮、链接、输入框等有明确角色的元素
2. getByLabel - 有label标签关联的表单元素
3. getByPlaceholder - 有placeholder的输入框
4. getByText - 唯一文本内容的元素

⚠️ 谨慎使用（可能不稳定）：
5. getByTestId - 需要开发配合添加data-testid
6. CSS选择器 - 类名可能变化

🔧 保底手段（仅当上述都不可用时）：
7. XPath - 复杂结构、动态元素的最后选择
```

**XPath使用场景（仅在以下情况使用）：**

```javascript
// 场景1：需要基于父子关系定位
page.locator('xpath=//div[@class="form-item"]//input')

// 场景2：需要基于兄弟元素定位
page.locator('xpath=//label[text()="用户名"]/following-sibling::input')

// 场景3：需要复杂条件组合
page.locator('xpath=//tr[contains(@class,"ant-table-row") and .//td[text()="admin"]]//button')

// 场景4：动态索引定位
page.locator('xpath=(//button[@type="submit"])[1]')
```

**常用定位器示例：**

```javascript
// ✅ 推荐写法
await page.getByRole('button', { name: '登录' }).click();
await page.getByRole('textbox', { name: '用户名' }).fill('admin');
await page.getByRole('link', { name: '首页' }).click();
await page.getByRole('checkbox', { name: '记住我' }).check();
await page.getByRole('combobox', { name: '选择部门' }).selectOption('IT');
await page.getByRole('tab', { name: '基本信息' }).click();
await page.getByRole('row', { name: /admin/ }).getByRole('button', { name: '编辑' }).click();

// ⚠️ 备选写法
await page.locator('.ant-btn-primary').click();
await page.locator('#username').fill('admin');

// 🔧 XPath保底
await page.locator('xpath=//button[normalize-space()="登录"]').click();
```

### 3.10 常见问题处理

| 问题     | 解决方案                     |
| ------ | ------------------------ |
| 多个相同元素 | `.first()` 或 `.nth(0)`   |
| 元素被遮挡  | `click({ force: true })` |
| 超时     | 增加 `timeout` 或检查定位器      |

## 4. 数据结构

### 4.1 字段定义

| 字段  | CN                 | JP                 | EN                 | 必填  |
| --- | ------------------ | ------------------ | ------------------ | --- |
| 画面  | screen_cn          | screen_jp          | screen_en          | ✅   |
| 功能  | function_cn        | function_jp        | function_en        | ✅   |
| 前置  | precondition_cn    | precondition_jp    | precondition_en    | ❌   |
| 步骤  | test_steps_cn      | test_steps_jp      | test_steps_en      | ✅   |
| 期望  | expected_result_cn | expected_result_jp | expected_result_en | ✅   |

公共字段：`case_number`(必填), `script_code`(必填)

### 4.2 🚨 自然语言用例与脚本一致性原则

**自然语言用例（test_steps、expected_result）是给人看的，必须与script_code内容完全一致：**

```
✅ 正确示例：
test_steps_cn: "1. 使用\"${USERNAME}\"登录系统\n2. 点击[用户管理]菜单\n3. 等待用户列表加载"
script_code: 对应的脚本确实执行了登录→点击菜单→等待加载

❌ 错误示例：
test_steps_cn: "1. 查看用户列表"  
script_code: 实际包含登录、导航等步骤（步骤不一致）
```

**作用**：

- 自然语言用例：供测试人员阅读理解
- 脚本用例：供自动化执行
- 两者必须描述相同的操作步骤

### 4.3 用例示例

```json
{
  "case_number": "USER-001",
  "screen_cn": "[用户管理]页面",
  "function_cn": "用户列表 - 查看用户列表",
  "precondition_cn": "1. 系统正常运行\n2. 存在有效的用户账号",
  "test_steps_cn": "1. 使用\"${USERNAME}\"登录系统\n2. 访问[用户管理]页面\n3. 等待用户列表加载",
  "expected_result_cn": "1. 页面正常显示\n2. 用户列表表格可见\n3. 显示用户数据",
  "script_code": "async (page) => { await page.context().clearCookies(); await page.goto('${BASE_URL}/login'); await page.getByPlaceholder('用户名').fill('${USERNAME}'); await page.getByPlaceholder('密码').fill('${PASSWORD}'); await page.getByRole('button', { name: '登录' }).click(); await page.waitForURL('**/users', { timeout: 10000 }); const rows = await page.locator('table tbody tr').count(); console.log('用户数量:', rows); return { success: rows > 0, rows: rows }; }"
}
```

## 5. 工作流程（全自动化）

**执行步骤：**当你收到Web自动化测试用例生成任务时，严格按照以下12步顺序执行。

**技术方案：**
- 使用 Playwright 编程 API（非 codegen 录制模式）
- AI 编写探索脚本 → 执行脚本 → 分析结果 → 生成用例
- 全流程自动化，无需用户手动操作浏览器

**🔐 HTTPS证书处理：**
```javascript
const browser = await chromium.launch({ 
  headless: true,
  ignoreHTTPSErrors: true  // 自动跳过自签名证书
});
```

---

**第一步：获取项目信息**

* 调用 `get_current_project_name()` 获取 `project_id` 和项目名称
* 失败则终止流程并报告错误

**第二步：列出Web用例集**

* 调用 `list_web_groups(project_id)` 获取Web用例集列表
* 使用用户指定的用例集名称，或第一个可用的用例集

**第三步：获取用例集元数据**

* 调用 `get_web_group_metadata(group_name, project_id)` 获取元数据
* 元数据包含：协议、服务器、端口、用户名、密码等
* 失败则报告错误并终止流程

**第四步：AI编写并执行登录探索脚本**

1. **AI 自动编写登录探索脚本**（使用 Playwright）：

```javascript
// 登录探索脚本示例
const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext({ ignoreHTTPSErrors: true });
  const page = await context.newPage();
  
  // 访问登录页面
  await page.goto('${BASE_URL}');
  
  // 尝试识别登录表单元素
  const usernameSelectors = [
    'input[name="username"]',
    'input[type="text"]',
    'input[placeholder*="用户"]',
    'input[placeholder*="ユーザー"]',
    '#username', '#user', '.username'
  ];
  
  const passwordSelectors = [
    'input[name="password"]',
    'input[type="password"]',
    '#password', '#pass', '.password'
  ];
  
  const loginButtonSelectors = [
    'button[type="submit"]',
    'button:has-text("登录")',
    'button:has-text("ログイン")',
    'button:has-text("Login")',
    'input[type="submit"]'
  ];
  
  // 自动填写并登录
  for (const selector of usernameSelectors) {
    try {
      await page.locator(selector).fill('${USERNAME}');
      break;
    } catch (e) { continue; }
  }
  
  for (const selector of passwordSelectors) {
    try {
      await page.locator(selector).fill('${PASSWORD}');
      break;
    } catch (e) { continue; }
  }
  
  for (const selector of loginButtonSelectors) {
    try {
      await page.locator(selector).click();
      break;
    } catch (e) { continue; }
  }
  
  // 等待登录完成
  await page.waitForTimeout(2000);
  
  // 输出登录后的 URL 和页面信息
  console.log('Login successful! URL:', page.url());
  console.log('Page title:', await page.title());
  
  await browser.close();
})();
```

2. **AI 使用 powershell 工具执行脚本**，自动验证登录是否成功
3. **AI 提取登录逻辑**，作为后续所有用例的公共登录模块

**第五步：AI自动编写并执行画面识别脚本 (AI Identifies All Screens Automatically)**

🚨 **在生成任何用例之前，必须先完成画面识别！**

**AI 自动编写画面识别脚本**（识别所有菜单和画面）：

```javascript
// 画面识别脚本示例
const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext({ ignoreHTTPSErrors: true });
  const page = await context.newPage();
  
  // 登录
  await page.goto('${BASE_URL}');
  // ... 执行登录逻辑 ...
  
  // 识别所有菜单项
  const menuSelectors = [
    'nav a', 'nav button', '.menu a', '.menu button',
    '.sidebar a', '.sidebar button', '[role="navigation"] a'
  ];
  
  let allMenus = [];
  for (const selector of menuSelectors) {
    try {
      const items = await page.locator(selector).all();
      for (const item of items) {
        const text = await item.textContent();
        const href = await item.getAttribute('href');
        allMenus.push({ text: text.trim(), href });
      }
    } catch (e) { continue; }
  }
  
  // 输出识别到的画面清单
  console.log('识别到的画面：', JSON.stringify(allMenus, null, 2));
  
  await browser.close();
})();
```

**AI 自动输出画面清单**：

```
📋 自动识别到的主要画面/菜单：
┌────┬─────────────────┬──────────────────┐
│ #  │ 画面/菜单名称    │ 状态             │
├────┼─────────────────┼──────────────────┤
│ 1  │ [ログイン]画面   │ ⏳ 待处理        │
│ 2  │ [ダッシュボード] │ ⏳ 待处理        │
│ 3  │ [ユーザー管理]  │ ⏳ 待処理        │
│ 4  │ [プロジェクト管理]│ ⏳ 待処理        │
│ 5  │ [設定]          │ ⏳ 待処理        │
└────┴─────────────────┴──────────────────┘

✅ 画面识别完成，开始自动生成用例...
```

**第六步：任务规模评估 (Task Scale Assessment)**

**在识别完所有画面后，进行规模评估：**

#### 6.1 统计画面和功能

**收集数据：**

- 识别到的画面总数
- 每个画面的功能点数量（基于可交互元素）
- 预计生成的用例总数（功能点 × 1.2，考虑反向用例和边界值用例）

#### 6.2 规模评估表

| 预计用例数   | 画面数    | 功能点数    | 预计耗时     | 中断可能性 | 建议处理方式   |
| ------- | ------ | ------- | -------- | ----- | -------- |
| ≤ 30条   | ≤ 5个   | ≤ 25个   | 15-30分钟  | 低     | 一次性完成    |
| 31-80条  | 6-10个  | 26-70个  | 30-60分钟  | 中等    | 一次性完成或分段 |
| 81-150条 | 11-20个 | 71-130个 | 60-120分钟 | 高     | 建议分段处理   |
| > 150条  | > 20个  | > 130个  | > 120分钟  | 很高    | 强烈建议分段   |

#### 6.3 输出任务评估

```
## 📊 Web用例生成任务评估

### 基本信息
- 目标网站: {base_url}
- 用例集: {group_name}
- 登录状态: ✅ 已登录成功

### 画面识别结果
- 识别到画面: {screen_count}个
- 功能点总数: {function_count}个
- 预计用例数: {estimated_cases}条（含反向用例）

### 画面清单
| 画面名称 | 功能点数 | 预计用例数 |
|---------|---------|----------|
| [ログイン]画面 | 3 | 5条（正常/错误密码/空用户名等） |
| [ダッシュボード]画面 | 2 | 3条 |
| [ユーザー管理]画面 | 8 | 12条（查询/新增/编辑/削除/开关等） |
| [設定]画面 | 5 | 7条 |
| ... | ... | ... |
| **合计** | **{function_count}** | **{estimated_cases}** |

### 规模预估
- 预计耗时: {estimated_time}分钟
- 中断可能性: {risk}
- Token预估: 约{tokens}

{如果用例数 > 80}
⚠️  **检测到大规模任务（{estimated_cases}条用例）**

AI将自动处理：
- 自动连续处理全部{screen_count}个画面
- 预计耗时{estimated_time}分钟
- 只在Token不足时暂停（需输入【继续】）
- 自动优化批次大小以避免超时

✅ AI将自动开始处理，无需用户选择方案
{/如果}

{如果用例数 ≤ 80}
✅ 任务规模适中，开始生成用例...
{/如果}
```

**用户决策处理：**

- 如果选择A：直接进入第七步，一次性处理所有画面
- 如果选择B：记录分段策略，按阶段执行（每阶段后暂停）

**第七步：按画面逐个录制操作并生成用例 (Record Actions and Generate Cases by Screen)**

#### 7.1 进度汇报策略（5阶段输出）

**阶段1：启动阶段**（任务开始时输出）

```
## 🚀 Web用例生成任务启动

### 任务配置
- 项目: {project_name}
- 用例集: {group_name}
- 目标网站: {base_url}
- 处理模式: {一次性完成 / 分阶段处理}
- 工具: Playwright CLI (npx playwright codegen)

### 画面处理计划
| 序号 | 画面名称 | 功能点数 | 预计用例数 | 状态 |
|-----|---------|---------|----------|------|
| 1 | [ログイン]画面 | 3 | 5 | ⏳ 待处理 |
| 2 | [ダッシュボード]画面 | 2 | 3 | ⏳ 待処理 |
| 3 | [ユーザー管理]画面 | 8 | 12 | ⏳ 待処理 |
| ... | ... | ... | ... | ... |
| **合计** | **{screen_count}个画面** | **{function_count}** | **{total_cases}** | |

### 开始处理...
```

**阶段2：执行阶段**（每处理完一个画面输出）

```
✅ [1/{screen_count}] [ログイン]画面 - 完成
   - 生成用例: 5条
   - 脚本验证: 5/5 通过
   - 成功回写: 5条
   - 累计进度: 5/{total_cases} (6.25%)

📝 [2/{screen_count}] [ダッシュボード]画面 - 処理中...
```

**简洁汇报规则：**

- 小批量（≤5个画面）：每完成1个画面汇报
- 中批量（6-15个画面）：每完成2个画面汇报
- 大批量（>15个画面）：每完成5个画面汇报详细信息

**阶段3：中断阶段**（Token不足时输出）

```
⏸️  进度暂停 - Token限制

### 当前进度
- **最后完成的画面**: [ユーザー管理]画面 (第3/{screen_count}个)
- **已生成用例**: {generated_cases}/{total_cases} ({percentage}%)
- **脚本验证通过**: {verified_cases}条
- **成功回写**: {success_count}条
- **失败**: {failed_count}条

### 剩余任务
| 序号 | 画面名称 | 预计用例数 | 状态 |
|-----|---------|----------|------|
| 4 | [設定]画面 | 7 | ⏳ 待処理 |
| 5 | [レポート]画面 | 10 | ⏳ 待処理 |
| ... | ... | ... | ... |
| **小计** | **{remaining_screens}个画面** | **約{remaining_cases}条** | |

### 已生成用例編号範囲
- 首条: LOGIN-001
- 末条: USER-012

**阶段4：自动继续阶段**（自动处理剩余任务）

```
▶️  自动继续生成用例...

从 第4个画面: [設定]画面 開始
当前进度: {processed_screens}/{total_screens} ({percentage}%)
```

**阶段5：完成阶段**（全部処理完成后）

```
✅ Web用例生成任務完成

### 📊 生成統計
- 処理画面数: {processed_screens}/{total_screens} ✅
- 識別控件数: {total_controls}个
- 総生成用例: {total_cases}条
- 脚本驗証通過: {verified_cases}条 ✅
- 🚨 成功回寫数据库: {success_count}条 ✅
- ⚠️ 验证失败未写入: {failed_count}条 ❌
- 成功率: {success_rate}%

### 🎯 CRUD操作覆盖情况
| CRUD类型 | 控件数 | 生成用例数 | 通过数 | 覆盖率 |
|---------|--------|-----------|--------|--------|
| **C-创建(CREATE)** | {create_controls} | {create_cases} | {create_passed} | {create_coverage}% |
| **R-检索(READ)** | {read_controls} | {read_cases} | {read_passed} | {read_coverage}% |
| **U-修改(UPDATE)** | {update_controls} | {update_cases} | {update_passed} | {update_coverage}% |
| **D-删除(DELETE)** | {delete_controls} | {delete_cases} | {delete_passed} | {delete_coverage}% |
| **其他控件** | {other_controls} | {other_cases} | {other_passed} | {other_coverage}% |
| **总计** | **{total_controls}** | **{total_cases}** | **{verified_cases}** | **{total_coverage}%** |

### 📋 用例編号範囲
- 首条: LOGIN-001
- 末条: REPORT-025

### 🎯 画面覆盖情況
| 画面名称 | 用例数 | CRUD覆盖 | 通過数 | 失敗数 |
|---------|--------|---------|--------|--------|
| [ログイン]画面 | 5 | R:5 | 5 | 0 |
| [ダッシュボード]画面 | 3 | R:3 | 3 | 0 |
| [ユーザー管理]画面 | 12 | C:2,R:4,U:3,D:3 | 11 | 1 |
| ... | ... | ... | ... | ... |
| **合計** | **{total_cases}** | **CRUD全覆盖** | **{verified_cases}** | **{failed_count}** |

{如果有失敗用例}
### ⚠️  验证失败用例详情（未写入数据库）
| 用例編号 | 画面 | CRUD类型 | 失敗原因 | 重试次数 | 状态 |
|---------|------|---------|---------|---------|------|
| USER-008 | [ユーザー管理] | C-CREATE | 元素定位失敗 | 3次 | ❌ 跳过 |

**说明：**
- 以上用例验证失败（result.success !== true），已跳过，未写入数据库
- 建議人工检查后手动添加或修正脚本后重新生成
{/如果}

### 🔧 変数表已更新
以下変数已写入用例集：
- ${BASE_URL} = {base_url}
- ${USERNAME} = {username}
- ${PASSWORD} = {password}
- ${WRONG_PASSWORD} = {wrong_password}

---
用例已保存到用例集：**{group_name}**
```

#### 7.2 生成与回写流程

```
已处理画面数 = 0
总画面数 = total_screens
已生成用例数 = 0
预计总用例数 = total_estimated_cases
失败用例列表 = []
变量列表 = []

FOR 每个画面 IN 识别的画面列表（严格按顺序）:
    当前批次用例 = []
    当前画面名称 = screen.name

    # 第八步：AI自动编写画面功能探索脚本 (AI Writes Screen Exploration Script)
    
    **🚨 AI自动编写脚本，自动探索当前画面的所有功能点！**
    
    1. **AI 编写探索脚本**，自动识别当前画面的所有可交互元素：
    
    ```javascript
    // 画面功能探索脚本示例（按CRUD顺序遍历）
    const { chromium } = require('playwright');
    
    (async () => {
      const browser = await chromium.launch({ headless: false });
      const context = await browser.newContext({ ignoreHTTPSErrors: true });
      const page = await context.newPage();
      
      // 执行登录
      await page.goto('${BASE_URL}');
      // ... 登录逻辑 ...
      
      // 导航到目标画面（如果需要）
      // await page.click('text=用户管理');
      
      console.log('🎯 开始按CRUD顺序识别控件...');
      
      let features = {
        create: [],   // C-创建类控件
        read: [],     // R-检索类控件
        update: [],   // U-修改类控件
        delete: [],   // D-删除类控件
        other: []     // 其他控件
      };
      
      // ===== 第1阶段：识别 C-创建 类控件 =====
      console.log('\\n📝 [C-CREATE] 识别创建类控件...');
      const createKeywords = ['新增', '创建', '添加', '新規', '作成', '追加', 'Create', 'Add', 'New'];
      const allButtons = await page.locator('button:visible, a.btn:visible, input[type="submit"]:visible').all();
      
      for (const btn of allButtons) {
        const text = await btn.textContent().catch(() => '');
        const trimmedText = text.trim();
        
        if (createKeywords.some(kw => trimmedText.includes(kw))) {
          features.create.push({
            type: 'button',
            text: trimmedText,
            crud: 'CREATE',
            action: 'click → fill form → save → delete'
          });
          console.log(`  ✓ [C] ${trimmedText}`);
        }
      }
      
      // ===== 第2阶段：识别 R-检索 类控件 =====
      console.log('\\n🔍 [R-READ] 识别检索类控件...');
      const readKeywords = ['查询', '搜索', '检索', '查看', '详情', '検索', '照会', '詳細', 'Search', 'Query', 'View', 'Detail'];
      
      // 2.1 检索按钮
      for (const btn of allButtons) {
        const text = await btn.textContent().catch(() => '');
        const trimmedText = text.trim();
        
        if (readKeywords.some(kw => trimmedText.includes(kw))) {
          features.read.push({
            type: 'button',
            text: trimmedText,
            crud: 'READ',
            action: 'click'
          });
          console.log(`  ✓ [R-Button] ${trimmedText}`);
        }
      }
      
      // 2.2 搜索框
      const searchInputs = await page.locator('input[type="search"]:visible, input[type="text"]:visible').all();
      for (const input of searchInputs) {
        const placeholder = await input.getAttribute('placeholder').catch(() => '');
        if (readKeywords.some(kw => placeholder.includes(kw)) || placeholder.toLowerCase().includes('search')) {
          features.read.push({
            type: 'search-input',
            placeholder: placeholder,
            crud: 'READ',
            action: 'fill → enter/click search button'
          });
          console.log(`  ✓ [R-Input] ${placeholder}`);
        }
      }
      
      // 2.3 表格行（查看详情）
      const tableRows = await page.locator('tbody tr:visible').count();
      if (tableRows > 0) {
        features.read.push({
          type: 'table-row',
          text: '表格第一行',
          crud: 'READ',
          action: 'click row to view details'
        });
        console.log(`  ✓ [R-Table] 表格（${tableRows}行）`);
      }
      
      // 2.4 详情链接
      const detailLinks = await page.locator('a:visible').all();
      for (const link of detailLinks.slice(0, 5)) {
        const text = await link.textContent().catch(() => '');
        if (readKeywords.some(kw => text.includes(kw))) {
          features.read.push({
            type: 'link',
            text: text.trim(),
            crud: 'READ',
            action: 'click → view details'
          });
          console.log(`  ✓ [R-Link] ${text.trim()}`);
        }
      }
      
      // ===== 第3阶段：识别 U-修改 类控件 =====
      console.log('\\n✏️ [U-UPDATE] 识别修改类控件...');
      const updateKeywords = ['编辑', '修改', '更新', '変更', '編集', '更新', 'Edit', 'Update', 'Modify'];
      
      // 3.1 编辑按钮
      for (const btn of allButtons) {
        const text = await btn.textContent().catch(() => '');
        const trimmedText = text.trim();
        
        if (updateKeywords.some(kw => trimmedText.includes(kw))) {
          features.update.push({
            type: 'button',
            text: trimmedText,
            crud: 'UPDATE',
            action: 'click → modify → save → delete'
          });
          console.log(`  ✓ [U-Button] ${trimmedText}`);
        }
      }
      
      // 3.2 开关切换
      const switches = await page.locator('input[type="checkbox"]:visible, .switch:visible, .toggle:visible').count();
      if (switches > 0) {
        features.update.push({
          type: 'switch',
          text: '开关切换',
          crud: 'UPDATE',
          action: 'toggle → verify → restore'
        });
        console.log(`  ✓ [U-Switch] 开关（${switches}个）`);
      }
      
      // ===== 第4阶段：识别 D-删除 类控件 =====
      console.log('\\n🗑️  [D-DELETE] 识别删除类控件...');
      const deleteKeywords = ['删除', '刪除', '削除', 'Delete', 'Remove'];
      
      for (const btn of allButtons) {
        const text = await btn.textContent().catch(() => '');
        const trimmedText = text.trim();
        
        if (deleteKeywords.some(kw => trimmedText.includes(kw))) {
          features.delete.push({
            type: 'button',
            text: trimmedText,
            crud: 'DELETE',
            action: 'create test data → click → confirm → verify'
          });
          console.log(`  ✓ [D] ${trimmedText}`);
        }
      }
      
      // ===== 第5阶段：识别其他控件 =====
      console.log('\\n🎛️  [OTHER] 识别其他控件...');
      
      // 5.1 下拉框
      const selects = await page.locator('select:visible').all();
      for (const select of selects) {
        const name = await select.getAttribute('name').catch(() => '');
        features.other.push({
          type: 'select',
          name: name,
          crud: 'OTHER',
          action: 'select option'
        });
        console.log(`  ✓ [Select] ${name || '下拉框'}`);
      }
      
      // 5.2 标签页
      const tabs = await page.locator('[role="tab"]:visible, .tab:visible').count();
      if (tabs > 0) {
        features.other.push({
          type: 'tabs',
          text: '标签页',
          crud: 'OTHER',
          action: 'click tab → explore tab content'
        });
        console.log(`  ✓ [Tabs] ${tabs}个标签页`);
      }
      
      // 5.3 分页器
      const pagination = await page.locator('.pagination:visible, .ant-pagination:visible').count();
      if (pagination > 0) {
        features.other.push({
          type: 'pagination',
          text: '分页器',
          crud: 'OTHER',
          action: 'click next page'
        });
        console.log(`  ✓ [Pagination] 分页器`);
      }
      
      // 输出CRUD统计
      console.log('\\n📊 CRUD控件统计：');
      console.log(`  C-CREATE: ${features.create.length}个`);
      console.log(`  R-READ:   ${features.read.length}个`);
      console.log(`  U-UPDATE: ${features.update.length}个`);
      console.log(`  D-DELETE: ${features.delete.length}个`);
      console.log(`  OTHER:    ${features.other.length}个`);
      console.log(`  总计:     ${features.create.length + features.read.length + features.update.length + features.delete.length + features.other.length}个`);
      
      // 输出完整的功能点列表
      console.log('\\n📋 画面功能点（按CRUD分类）：', JSON.stringify(features, null, 2));
      
      await browser.close();
    })();
    ```
    
    2. **AI 执行脚本**，自动获取当前画面的功能点列表（按CRUD分类）
    3. **AI 自动分析**，识别页面语言（中文/日文/英文）

    # 第九步：AI基于CRUD分类自动设计测试用例 (AI Designs Test Cases by CRUD)
    
    **🚨 重要：按CRUD顺序生成用例，确保数据依赖正确！**
    
    # 9.1 先处理 C-CREATE 类用例
    FOR 每个创建类控件 IN features.create:
        AI 自动生成创建用例：
        - 用例类型：C-CREATE
        - 脚本流程：登录 → 点击创建 → 填表 → 保存 → 验证 → 删除测试数据
        - 数据清理：🚨 必须删除创建的测试数据
        将用例加入当前批次
    END FOR
    
    # 9.2 再处理 R-READ 类用例
    FOR 每个检索类控件 IN features.read:
        AI 自动生成检索用例：
        - 用例类型：R-READ
        - 脚本流程：登录 → 执行检索操作 → 验证结果
        - 数据清理：无需清理（只读操作）
        将用例加入当前批次
    END FOR
    
    # 9.3 接着处理 U-UPDATE 类用例
    FOR 每个修改类控件 IN features.update:
        AI 自动生成修改用例：
        - 用例类型：U-UPDATE
        - 脚本流程：登录 → 创建测试数据 → 执行修改 → 验证 → 删除测试数据
        - 数据清理：🚨 必须删除修改测试用的数据
        将用例加入当前批次
    END FOR
    
    # 9.4 最后处理 D-DELETE 类用例
    FOR 每个删除类控件 IN features.delete:
        AI 自动生成删除用例：
        - 用例类型：D-DELETE
        - 脚本流程：登录 → 创建测试数据 → 执行删除 → 验证删除成功
        - 数据清理：已删除，无需额外清理
        将用例加入当前批次
    END FOR
    
    # 9.5 处理其他控件
    FOR 每个其他控件 IN features.other:
        AI 自动生成对应用例
        将用例加入当前批次
    END FOR
    
    FOR 每条用例:
        AI 基于探索脚本的识别结果，自动生成测试操作代码
        基于功能点类型和CRUD分类设计正向用例（正常流程）
        IF 需要反向用例:
            设计反向用例（错误输入、边界值等）
        END IF

        AI 自动生成用例数据（**注意UI元素标识**）：
        - case_number = {画面缩写}-{序号}
        - screen_{语言} = 画面名称（用[]标识）
        - function_{语言} = 功能描述（标注CRUD类型）
        - precondition_{语言} = 前置条件（包含登录状态）
        - test_steps_{语言} = 测试步骤（**确保UI元素用[]标识，包含数据清理步骤**）
        - expected_result_{语言} = 期望结果（**确保UI元素用[]标识**）
        - script_code = 完整脚本（包含登录逻辑、业务操作、数据恢复）

        ⚠️ 脚本生成注意：
        - AI 自动编写 async (page) => { ... } 格式的脚本
        - 添加 console.log() 调试日志
        - 添加 return { success: true/false, message: '...' } 返回值
        - 将硬编码值替换为变量占位符（${BASE_URL}、${USERNAME}等）
        - 确保脚本可以独立执行（包含完整登录流程）
        - 🚨 根据CRUD类型添加相应的数据清理逻辑

        确保：
        - 每条用例能独立执行（除登录画面外，都包含登录流程）
        - C-CREATE用例包含：创建→验证→删除
        - U-UPDATE用例包含：创建→修改→验证→删除
        - D-DELETE用例包含：创建→删除→验证
        - 开关类用例包含：切换→验证→恢复
        - 所有变量使用占位符（${BASE_URL}、${USERNAME}等）

        将用例加入当前批次
    END FOR

    # 第十步：主动验证脚本执行 (Validate Script Execution)
    # 🚨 关键原则：只有验证通过(result.success === true)的用例才能写入数据库
    验证通过的用例列表 = []
    验证失败的用例列表 = []
    
    FOR 每条用例 IN 当前批次用例:
        TRY:
            # 使用 powershell 工具执行 node 验证脚本
            # 指数退避重试（1s, 2s, 4s）
            FOR retry IN [1, 2, 3]:
                主动执行：使用 node 验证脚本（通过 powershell 工具）
                IF result.success === true:
                    标记为"验证通过" ✅
                    将用例加入验证通过列表
                    BREAK
                ELSE:
                    IF retry < 3:
                        等待 2^(retry-1) 秒
                        继续重试
                    ELSE:
                        标记为"验证失败" ❌
                        记录失败: {
                            case_number: 用例编号,
                            error: result.error,
                            suggestion: "检查元素定位器或页面结构",
                            action: "跳过，不写入数据库"
                        }
                        将用例加入验证失败列表（不写入数据库）
                    END IF
                END IF
            END FOR
        CATCH SCRIPT_ERROR:
            标记为"验证失败" ❌
            记录失败详情到验证失败列表
        END TRY
    END FOR
    
    # 🚨 关键检查：只处理验证通过的用例
    IF 验证通过的用例列表为空:
        输出: "⚠️  本批次所有用例验证失败，跳过写入数据库"
        记录失败详情
        继续处理下一个画面
    END IF

    # 质量检查点（每个画面完成后）
    检查本批次用例：
    - UI元素是否都使用[]标识
    - script_code中的控制台日志是否使用[]标识
    - 变量占位符是否正确（${...}格式）
    - 数据恢复逻辑是否完整
    - 登录流程是否包含（非登录画面）
    - 录制代码是否已正确转换（移除test/expect等）

    # 第十一步：主动批量创建用例并写入变量
    # 🚨 关键原则：只写入验证通过的用例（result.success === true）
    TRY:
        # 首次回写时，同时写入变量表
        IF 已处理画面数 == 0:
            主动构建变量列表 = [
                {var_key: "base_url", var_value: base_url, var_desc: "系统基础URL"},
                {var_key: "username", var_value: username, var_desc: "登录用户名"},
                {var_key: "password", var_value: password, var_desc: "登录密码"},
                {var_key: "wrong_password", var_value: "WrongPass@123", var_desc: "错误密码"}
            ]
        END IF

        # 🚨 只写入验证通过的用例！
        主动调用 create_web_cases(
            project_id=project_id,
            group_name=group_name,
            cases=验证通过的用例列表,  # 🚨 关键：只包含 result.success === true 的用例
            variables=变量列表,
            continue_on_error=true
        )

        已处理画面数 += 1
        已生成用例数 += 验证通过的用例列表长度
        
        # 记录失败用例到全局失败列表
        IF 验证失败的用例列表不为空:
            失败用例列表.extend(验证失败的用例列表)
        END IF

        # 根据批次大小决定汇报详细程度
        IF 总画面数 <= 5 OR 当前画面数 % 汇报频率 == 0:
            主动输出详细进度（阶段2格式）
        ELSE:
            主动输出简化进度
        END IF

    CATCH API_ERROR:
        # 指数退避重试
        FOR retry IN [1, 2, 3]:
            等待 2^(retry-1) 秒
            TRY:
                重新调用 create_web_cases()
                成功则 BREAK
            CATCH:
                IF retry == 3:
                    记录失败: 失败用例列表.append({
                        screen: 当前画面名称,
                        cases_count: 本批次用例数,
                        error: error.message
                    })
                END IF
            END TRY
        END FOR
    END TRY

    # Token检查
    IF 即将达到Token上限:
        主动输出中断信息（阶段3格式）
        STOP
    END IF
END FOR

# 全部完成后主动输出阶段5报告
```

**第十一步：主动输出进度并等待继续 (Output Progress and Wait)**

* **每完成一个画面后，主动输出进度**：

```
✅ [登录]画面 - 已完成（生成 8 条用例）
┌────┬─────────────────┬──────────────────┐
│ #  │ 画面/菜单名称    │ 状态             │
├────┼─────────────────┼──────────────────┤
│ 1  │ [登录]画面       │ ✅ 已完成 (8条)  │
│ 2  │ [首页/Dashboard] │ ⏳ 待处理        │
│ 3  │ [用户管理]      │ ⏳ 待处理        │
│ 4  │ [项目管理]      │ ⏳ 待处理        │
│ 5  │ [设置]          │ ⏳ 待处理        │
└────┴─────────────────┴──────────────────┘

📊 当前进度：1/5 画面已完成
✅ AI自动继续处理下一个画面（仅Token不足时暂停）
```

* 自动处理下一个画面，返回 **第七步**。
* 如果所有画面已处理完成，自动进入下一步。

**第十二步：主动输出汇总报告 (Output Summary Report)**

* 主动汇总所有画面的用例生成情况，自动输出最终报告：

```
📊 用例生成完成汇总：
┌────┬─────────────────┬──────────┬──────────┐
│ #  │ 画面名称         │ 正向用例 │ 反向用例 │
├────┼─────────────────┼──────────┼──────────┤
│ 1  │ [登录]画面       │ 4        │ 4        │
│ 2  │ [首页/Dashboard] │ 6        │ 2        │
│ 3  │ [用户管理]      │ 12       │ 8        │
│ ...│ ...             │ ...      │ ...      │
├────┼─────────────────┼──────────┼──────────┤
│    │ 合计            │ 45       │ 30       │
└────┴─────────────────┴──────────┴──────────┘
总计：75 条用例
```

## 6. 错误处理与异常场景 (Error Handling)

### 6.1 项目/用例集获取失败

```
❌ 错误：无法获取项目信息 / 用例集不存在

处理方式：
1. 确认已选择项目
2. 检查用例集名称是否正确
3. 列出可用的用例集供用户选择
```

### 6.2 目标网站无法访问

**场景1：网络连接失败**

```
❌ 错误：无法连接到 {base_url}

可能原因：
- 网络不可达
- 服务器未启动
- 防火墙阻止
- 端口配置错误

处理方式：
1. 检查网络连接
2. 确认服务器状态
3. 验证URL格式：{protocol}://{server}:{port}
4. 尝试ping服务器地址
```

**场景2：HTTPS证书错误**

```
⚠️  检测到HTTPS证书错误（ERR_CERT_AUTHORITY_INVALID）

自动处理：
✓ 使用 --ignore-https-errors 参数重新启动 Playwright CLI
✓ 继续访问网站...

命令：
npx playwright codegen --ignore-https-errors {protocol}://{server}:{port}

说明：
- 录制阶段：使用 --ignore-https-errors 参数启动 playwright codegen
- script_code字段：无需特殊处理（Docker环境已配置）
```

### 6.3 登录失败

**场景：凭证错误或登录逻辑变化**

**处理流程：**

```python
TRY:
    AI 自动执行登录探索脚本
    脚本自动检测：登录是否成功（检查URL跳转或特定元素）

    IF 登录成功:
        继续画面识别
    ELSE:
        抛出 LOGIN_FAILED
    END IF

CATCH LOGIN_FAILED:
    输出：
    """
    ❌ 登录失败，无法继续生成用例

    可能原因：
    - 用户名/密码错误（meta_user / meta_password）
    - 登录页面结构变化
    - 登录按钮定位失败
    - 登录后跳转URL不符合预期

    自动尝试的解决方案：
    1. ✅ 已尝试多种常见的登录表单选择器
    2. ✅ 已尝试等待登录后页面跳转
    3. ❌ 所有尝试均失败

    建议：
    - 检查用例集元数据中的凭证是否正确
    - 检查登录页面是否可访问
    - 建议排查登录流程配置
    """
    终止流程
END TRY
```

### 6.4 画面识别不完整

**场景：某些画面无法通过自动脚本识别**

```
⚠️  画面识别可能不完整

已自动识别画面：{identified_count}个
可能原因：
- 某些画面需要特定权限才能访问
- 某些菜单项是动态加载的
- 某些子菜单被隐藏

处理方式：
✅ AI将继续处理已识别的画面，自动生成测试用例
```

### 6.5 脚本验证失败（带重试机制）

**场景：脚本执行失败**

**🚨 关键原则：验证失败的用例禁止写入数据库！**

**处理流程（指数退避）：**

```python
FOR retry_count IN [1, 2, 3]:
    TRY:
        主动使用 powershell 工具执行脚本验证

        IF result.success === true:
            标记为"验证通过" ✅
            → 可以写入数据库
            BREAK
        ELSE:
            抛出 SCRIPT_ERROR(result.message)
        END IF

    CATCH SCRIPT_ERROR AS error:
        IF retry_count < 3:
            等待时间 = 1 * (2 ^ (retry_count - 1))  # 1s, 2s, 4s
            主动输出: "⚠️  脚本验证失败，{等待时间}秒后重试（第{retry_count}/3次）..."
            等待(等待时间)
        ELSE:
            主动输出: "❌ 脚本验证失败3次，跳过此用例（不写入数据库）"
            标记为"验证失败" ❌
            记录错误详情到最终报告: {
                case_number: 用例编号,
                screen: 画面名称,
                error: error.message,
                retry_count: 3,
                suggestion: "检查元素定位器或页面结构，可能需要调整定位策略",
                action: "跳过，未写入数据库"
            }
            # 🚨 继续处理下一个用例（不阻断整体流程，不写入此条失败用例）
        END IF
    END TRY
END FOR
```

**常见失败原因及解决方案：**

| 失败原因       | 解决方案                        |
| ---------- | --------------------------- |
| 元素定位超时     | 增加timeout或检查定位器             |
| 元素不可见      | 使用 `click({ force: true })` |
| 页面未加载完成    | 增加 `waitForTimeout()`       |
| 登录状态丢失     | 检查cookies清理逻辑               |
| XPath表达式错误 | 使用更稳定的getByRole定位器          |
| 录制代码未转换    | 移除test()/expect()等框架API     |

### 6.6 元素定位失败自动降级

**场景：首选定位器失败**

**自动降级策略（在脚本中实现）：**

```javascript
// 策略1：getByRole失败 → 尝试getByText
try {
    await page.getByRole('button', { name: 'ログイン' }).click();
} catch (e) {
    console.warn('[Fallback] getByRole失败，尝试getByText');
    await page.getByText('ログイン').click();
}

// 策略2：语义定位失败 → 尝试CSS选择器
try {
    await page.getByLabel('ユーザー名').fill(username);
} catch (e) {
    console.warn('[Fallback] getByLabel失败，尝试CSS');
    await page.locator('input[name="username"]').fill(username);
}

// 策略3：最后使用XPath
try {
    await page.locator('.ant-btn-primary').click();
} catch (e) {
    console.warn('[Fallback] CSS失败，使用XPath');
    await page.locator('xpath=//button[@type="submit"]').click();
}
```

### 6.7 用例回写失败

**处理流程（指数退避）：**

```python
TRY:
    调用 create_web_cases(
        project_id=project_id,
        group_name=group_name,
        cases=current_batch,
        variables=variables_list,
        continue_on_error=true
    )
    成功 ✅

CATCH API_ERROR AS error:
    # 指数退避重试
    FOR retry IN [1, 2, 3]:
        等待 2^(retry-1) 秒  # 1s, 2s, 4s
        TRY:
            重新调用 create_web_cases()
            IF 成功:
                BREAK
            END IF
        CATCH:
            IF retry == 3:
                记录失败: {
                    screen: 当前画面,
                    cases_count: 本批次用例数,
                    error: error.message,
                    suggestion: "检查网络连接或稍后重试"
                }
                # 继续处理下一个画面（continue_on_error=true）
            END IF
        END TRY
    END FOR
END TRY
```

### 6.8 变量表写入失败

**场景：用例集变量表更新失败**

```
⚠️  警告：变量表写入失败

影响：
- 用例已生成并保存 ✅
- 但变量值未写入用例集 ❌
- AI将自动重试写入变量

已生成的变量列表：
- ${BASE_URL} = {base_url}
- ${USERNAME} = {username}
- ${PASSWORD} = {password}
- ${WRONG_PASSWORD} = {wrong_password}

处理方式：
1. AI自动重试调用 create_web_cases 带上 variables 参数
2. 如持续失败，记录错误并继续执行
```

### 6.9 Token超限处理（CRITICAL）

**🚨 唯一需要用户交互的情况：Token不足！**

✅ **必须的行为：**

1. 检测到Token即将用完时立即暂停（当前画面处理完成后）
2. 按照"阶段3：中断阶段"格式输出详细进度
3. 明确显示已处理和待处理的画面清单
4. 记录最后一条用例编号和画面名称
5. **提示用户输入【继续】以恢复执行**
6. 用户输入【继续】后，自动从断点恢复

❌ **禁止的行为（即使Token不足也不能做）：**

- 简化用例内容或脚本
- 跳过某些画面
- 省略数据恢复逻辑（创建→验证→删除）
- 合并多个功能为一条用例
- 省略脚本验证环节

**中断时的状态保存：**

```
⏸️  进度暂停 - Token限制

当前状态：
- 已处理画面: [ログイン]、[ダッシュボード]、[ユーザー管理] (3/{screen_count})
- 已生成用例: LOGIN-001 ~ USER-012 (20条)
- 脚本验证通过: 18条
- 成功回写: 18条
- 失败: 2条
- 断点位置: 第4个画面 [設定]画面

👉 请输入【继续】恢复处理剩余画面
```

### 6.10 Playwright CLI 启动失败

**场景：CLI工具无法启动或异常退出**

```
❌ 错误：Playwright CLI 启动失败

常见原因：
1. Playwright 未安装（需运行 npm install -D @playwright/test）
2. 浏览器未安装（需运行 npx playwright install）
3. 系统依赖缺失（需运行 npx playwright install-deps）
4. 端口被占用或网络问题

处理方式：
1. 检查 Playwright 安装状态：
   npx playwright --version

2. 安装浏览器：
   npx playwright install chromium

3. 安装系统依赖（Linux）：
   npx playwright install-deps

4. 重新启动录制：
   npx playwright codegen --ignore-https-errors <URL>

5. 如果持续失败，检查 Node.js 版本（需要 >= 16）
```

### 6.11 录制代码转换常见问题

**场景：Playwright CLI 录制的代码需要转换为 script_code 格式**

**Playwright CLI 录制格式 → script_code 格式 转换规则：**

```javascript
// ❌ CLI 录制的原始代码（@playwright/test 格式）
import { test, expect } from '@playwright/test';

test('test', async ({ page }) => {
  await page.goto('https://192.168.11.104:8443/login');
  await page.getByPlaceholder('用户名').fill('admin');
  await page.getByPlaceholder('密码').fill('password123');
  await page.getByRole('button', { name: 'ログイン' }).click();
  await expect(page).toHaveURL(/.*dashboard/);
});

// ✅ 转换后的 script_code 格式
async (page) => {
  console.log('[Step 1] 清理状态，访问登录页...');
  await page.context().clearCookies();
  await page.goto('${BASE_URL}/login');

  console.log('[Step 2] 输入凭证...');
  await page.getByPlaceholder('用户名').fill('${USERNAME}');
  await page.getByPlaceholder('密码').fill('${PASSWORD}');

  console.log('[Step 3] 点击登录...');
  await page.getByRole('button', { name: 'ログイン' }).click();

  console.log('[Step 4] 验证登录结果...');
  await page.waitForURL('**/dashboard', { timeout: 10000 });

  console.log('[Success] 登录成功');
  return { success: true, message: '登录成功，已跳转到首页' };
}
```

**转换检查清单：**

| 检查项        | 录制代码                                  | 转换后代码                                            |
| ---------- | ------------------------------------- | ------------------------------------------------ |
| 函数签名       | `test('...', async ({ page }) => {`   | `async (page) => {`                              |
| import语句   | `import { test, expect }`             | 删除                                               |
| expect断言   | `await expect(page).toHaveURL(...)`   | `await page.waitForURL(...)`                     |
| expect元素断言 | `await expect(locator).toBeVisible()` | `const visible = await locator.isVisible()`      |
| 硬编码URL     | `'https://192.168.11.104:8443'`       | `'${BASE_URL}'`                                  |
| 硬编码用户名     | `'admin'`                             | `'${USERNAME}'`                                  |
| 硬编码密码      | `'password123'`                       | `'${PASSWORD}'`                                  |
| 返回值        | 无                                     | `return { success: true/false, message: '...' }` |
| 调试日志       | 无                                     | `console.log('[Step N] ...')`                    |
| Cookies清理  | 无                                     | `await page.context().clearCookies()`            |

### 6.12 错误处理总原则

1. **画面遍历优先**：任何单个用例的错误不应阻断画面遍历
2. **记录继续**：失败的用例记录原因后继续处理下一个
3. **🚨 验证失败不写入**：只有 result.success === true 的用例才能写入数据库
4. **指数退避**：重试使用1s→2s→4s的延迟
5. **完整报告**：所有失败在最终报告中详细列出（包括未写入的用例）
6. **自动化优先**：AI自动尝试解决问题，严重错误才终止任务
7. **Token限制**：只在Token不足时暂停，等待用户输入【继续】

---

## 7. 画面穷尽原则

### 7.1 必须穷尽所有画面

```
❌ 禁止行为：
- 只生成部分画面的用例就结束
- 跳过"不重要"的画面
- 未遍历所有菜单项

✅ 正确做法：
- 逐个遍历所有主菜单和子菜单
- 每个画面都必须生成用例
- 自动处理下一个画面，无需等待用户输入
```

### 7.2 必须穷尽画面上的所有控件

**每个画面必须识别并为以下所有控件生成用例：**

| 控件类型 | 正向用例     | 反向用例             |
| ---- | -------- | ---------------- |
| 按钮   | 点击执行正常功能 | 禁用状态、无权限         |
| 输入框  | 正常输入     | 空值、超长、特殊字符、SQL注入 |
| 下拉框  | 选择有效选项   | 无选项、默认值          |
| 复选框  | 勾选/取消    | 必选未勾选            |
| 表格   | 查看、排序、分页 | 空数据、大数据量         |
| 链接   | 正常跳转     | 无权限页面            |
| 文件上传 | 正常上传     | 超大文件、错误格式        |

### 7.3 自动继续机制

```
🚨 严禁提前结束！

在所有画面处理完毕之前：
1. 每完成一个画面必须输出进度
2. 自动继续处理下一个画面
3. 只有所有画面处理完毕才能结束
```

## 8. 工具速查

### 8.1 测试管理工具

| 工具                                                         | 用途         |
| ---------------------------------------------------------- | ---------- |
| `get_current_project_name()`                               | 获取当前项目     |
| `list_web_groups(project_id)`                              | 获取Web用例集列表 |
| `get_web_group_metadata(group_name)`                       | 获取元数据      |
| `list_web_cases(project_id, group_id)`                     | 获取现有用例     |
| `create_web_cases(project_id, group_id, cases, variables)` | 创建用例+变量    |
| `update_web_cases(project_id, group_id, cases)`            | 批量更新用例     |

### 8.2 Playwright CLI 命令速查

| 命令                                                      | 用途                       |
| ------------------------------------------------------- | ------------------------ |
| `npx playwright codegen <url>`                          | 启动录制浏览器，自动生成代码           |
| `npx playwright codegen --ignore-https-errors <url>`    | 跳过HTTPS证书验证启动录制（自签名证书必用） |
| `npx playwright codegen --output=<file> <url>`          | 保存录制脚本到文件                |
| `npx playwright codegen --browser=chromium <url>`       | 指定浏览器启动录制                |
| `npx playwright codegen --viewport-size=1280,720 <url>` | 指定视窗大小                   |
| `npx playwright codegen --device="iPhone 13" <url>`     | 模拟移动设备                   |
| `npx playwright codegen --save-storage=auth.json <url>` | 保存认证状态                   |
| `npx playwright codegen --load-storage=auth.json <url>` | 加载认证状态（跳过登录）             |
| `npx playwright install`                                | 安装所有浏览器                  |
| `npx playwright install chromium`                       | 仅安装Chromium              |
| `npx playwright --version`                              | 查看Playwright版本           |

> 🚨 **重要提醒**：HTTPS自签名证书站点必须使用 `--ignore-https-errors` 参数！

---

生成Web自动化测试用例，目标用例集：**{{group_name}}**

```

```
