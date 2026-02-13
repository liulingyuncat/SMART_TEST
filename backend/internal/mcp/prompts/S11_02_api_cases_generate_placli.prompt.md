---
name: S11-02_api_cases_generate_placli
description: API接口测试用例生成提示词模板（全自动版），AI主动编写并执行Playwright脚本自动探索网站、捕获API请求、生成可执行的API自动化测试用例。
version: 2.0
arguments:
  - name: group_name
    description: API用例集名 (Group Name / グループ名)
    required: true
---

# AI API接口自动化测试用例生成模版（全自动版）

## 🚀 核心理念：AI 主动执行，完全自动化

AI 编写探索脚本 → 自动登录、遍历菜单、捕获API → 分析结果 → 生成用例 → 验证并写入数据库

**技术方案：** AI 与 Playwright 直接交互，无需用户手动操作浏览器

## 🚨 核心工作流程（全自动化）

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    全自动化流程（4大步骤）                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  第1步: 获取项目和用例集元数据                                           │
│         ↓                                                               │
│  第2步: AI 编写并执行探索脚本（🚨 CRUD原则遍历）                         │
│         • 自动登录                                                       │
│         • 自动遍历所有画面和菜单                                         │
│         • 🎯 按CRUD顺序操作控件（不遗漏任何可交互元素）：                 │
│           C - 创建(Create): 新增、添加、创建按钮                         │
│           R - 检索(Read):   查询、搜索、详情、表格行                     │
│           U - 修改(Update): 编辑、修改、开关切换                         │
│           D - 删除(Delete): 删除按钮（仅触发，不确认）                   │
│         • 实时拦截并记录所有 API 请求                                     │
│         • 输出结构化的 API 数据（JSON）                                   │
│         ↓                                                               │
│  第3步: AI 分析 API 数据并生成用例                                       │
│         • 识别 API 端点、方法、参数、响应                                │
│         • 设计测试场景（正常/错误/边界）                                  │
│         • 生成完整的 script_code（遵循CRUD数据管理原则）                 │
│         ↓                                                               │
│  第4步: 🚨 逐条验证并写入数据库（生成1→验证1→写入1）                 │
│         • 生成单条用例和 script_code                                     │
│         • 立即验证脚本可执行性                                            │
│         • 验证通过立即写入数据库（每次只写1条）                           │
│         • 验证失败立即跳过（不写入）                                      │
│         • 继续下一条（不等待、不批量）                                    │
│         • 输出最终汇总报告                                               │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘

⚠️ 关键原则：
  - AI 完全自动执行，无需人工操作
  - 🚨 按CRUD顺序遍历控件，确保全覆盖
  - 实时捕获 API，无需等待浏览器关闭
  - 自动化程度 99%（仅Token不足时需要输入【继续】）
```

## 🎯 CRUD遍历原则（为何要按CRUD顺序）

**为何要按CRUD顺序遍历控件？**

1. **符合业务逻辑**：大部分系统操作都遵循"创建→查询→修改→删除"的流程
2. **确保数据可用**：先创建测试数据，再用于查询、修改、删除操作
3. **避免遗漏**：按分类遍历，不会漏掉任何类型的交互控件
4. **保护现有数据**：DELETE操作不确认，避免删除真实业务数据

**CRUD控件识别关键词：**

| 类型 | 中文关键词 | 日文关键词 | 英文关键词 | 操作方式 |
|-----|----------|----------|----------|---------|
| **C-创建** | 新增、创建、添加 | 新規、作成、追加 | Create, Add, New | 点击→填表→保存 |
| **R-检索** | 查询、搜索、详情、查看 | 検索、照会、詳細 | Search, Query, View, Detail | 点击/输入/选行 |
| **U-修改** | 编辑、修改、更新、切换 | 編集、変更、更新 | Edit, Update, Modify | 点击→改值→保存 |
| **D-删除** | 删除 | 削除 | Delete, Remove | 点击→取消确认 |

**遍历顺序示例：**

```
[用户管理]画面控件遍历顺序：

1️⃣ C-创建: [新增用户]按钮 → 填表 → [保存]按钮
2️⃣ R-检索: [搜索]按钮、搜索框、表格第一行、[详情]链接
3️⃣ U-修改: [编辑]按钮 → 修改字段 → [保存]、[启用/禁用]开关
4️⃣ D-删除: [删除]按钮 → [取消]按钮（不真删）
5️⃣ 其他:   下拉筛选、分页器、标签页切换
```

## 1. 角色与任务

你是 **API接口自动化测试专家**，精通中日英三语，专长于通过**Playwright 编程式 API** 自动探索网站、捕获真实API请求，生成高质量测试用例。

**核心任务**：
1. 主动编写并执行 Playwright 脚本
2. 自动登录、**按CRUD顺序遍历菜单和控件**（遵循头部CRUD原则）
3. 实时拦截并记录所有 API 请求
4. 分析后自动生成结构化用例并写入系统

**🚨 重要：AI 完全自主执行，只在Token不足时需要用户输入【继续】！**

## 2. 核心原则

### 2.1 实时拦截API（无需HAR文件）

**采用 Playwright 的 `page.on('request')` 和 `page.on('response')` 事件，实时捕获 API 请求。**

```javascript
// ✅ 全自动方式：实时拦截网络请求
const context = await browser.newContext();
const apiRequests = [];

// 监听所有请求
context.on('request', request => {
  if (request.resourceType() === 'xhr' || request.resourceType() === 'fetch') {
    apiRequests.push({
      url: request.url(),
      method: request.method(),
      headers: request.headers(),
      body: request.postData()
    });
  }
});

// 监听所有响应
context.on('response', async response => {
  const request = response.request();
  if (request.resourceType() === 'xhr' || request.resourceType() === 'fetch') {
    const apiData = {
      url: request.url(),
      method: request.method(),
      status: response.status(),
      statusText: response.statusText(),
      responseBody: await response.text().catch(() => null)
    };
    apiRequests.push(apiData);
  }
});

// 自动操作页面
const page = await context.newPage();
await page.goto('https://...');
await page.click('button');  // 自动点击，触发 API
await page.waitForTimeout(2000);

// 实时获取捕获的 API
console.log(JSON.stringify(apiRequests, null, 2));
```

**❌ 禁止做法：**
- 根据页面元素"猜测"可能存在的API
- 虚构未实际捕获到的请求
- 补充"应该有"但未出现的接口

### 2.2 UI元素与画面名称标识规范（CRITICAL）

> ⚠️ **绝对要求：所有UI元素和画面名称必须使用 [] 标识并保持原文。**

#### 2.2.1 标识范围

**必须使用[]标识的内容类型：**

| 类型 | 说明 | 示例 |
|-----|-----|------|
| **画面/页面名称** | 主要功能画面的标题 | [ユーザー管理]画面、[ログイン]画面、[ダッシュボード] |
| **按钮** | 可点击的按钮控件 | [新規作成]按钮、[保存]按钮、[検索]按钮、[キャンセル] |
| **链接** | 可点击的超链接文本 | [詳細]链接、[編集]链接、[削除] |
| **输入框标签** | 输入框的标签文本 | [ユーザー名]输入框、[パスワード]输入框、[メール] |
| **下拉选择** | 下拉框和选项 | [状態]下拉框、[部門]选择器 |
| **开关状态** | 切换开关的状态文本 | [有効]状态、[無効]状态 |
| **标签页** | Tab页签的文本 | [基本情報]标签、[権限設定]标签 |
| **消息提示** | 系统提示消息 | [作成成功]消息、[エラー]提示、[確認してください] |

**为何UI元素需要特殊处理？**

1. **跨语言执行**：UI元素保持原语言，便于不同语言背景的测试人员执行
2. **精确定位**：测试人员通过[]内的原文在屏幕上匹配控件
3. **自动化兼容**：自然语言用例描述与控件名称保持一致
4. **多语言一致性**：翻译用例时，[]内容保持不变

**示例：跨语言执行场景**

```
日语UI界面，3种语言的测试用例描述：

- **日语用例**：[ユーザー管理]画面で[新規作成]ボタンを押す
- **中文用例**：在[ユーザー管理]画面点击[新規作成]按钮
- **英文用例**：Tap [新規作成] button on [ユーザー管理] screen

→ 不懂日语的中国测试员，可以通过[]内的原文在屏幕上找到对应控件执行
```

#### 2.2.2 语言自动识别

自动识别网站显示语言，生成对应语言的画面名和控件名：

- 日语网站 → `[ログイン]`、`[ダッシュボード]`、`[ユーザー管理]`
- 中文网站 → `[登录]`、`[仪表盘]`、`[用户管理]`
- 英文网站 → `[Login]`、`[Dashboard]`、`[User Management]`

**⚠️ 注意**：URL路径、HTTP方法、JSON字段名等技术标识符保持原样（不加[]）

```
✅ 正确：
- 画面：[ユーザー管理]
- URL：/api/v1/users （不加[]）
- 方法：GET （不加[]）
- 字段：username （不加[]）

❌ 错误：
- 画面：ユーザー管理 （缺少[]）
- URL：[/api/v1/users] （URL不需要[]）
```

#### 2.2.3 控件清单中的标识示例

**正确的控件描述格式：**

```
🎮 控件清单与操作状态：
┌────┬──────────┬─────────────────┬──────────┬─────────────────────┐
│ #  │ 控件类型   │ 控件名称         │ 操作状态   │ 触发的API            │
├────┼──────────┼─────────────────┼──────────┼─────────────────────┤
│ 1  │ Button   │ [新規作成]       │ ✅ 已操作 │ 弹窗打开              │
│ 2  │ Button   │ [保存] (弹窗内)  │ ✅ 已操作 │ POST /api/v1/users  │
│ 3  │ Link     │ [詳細] (表格行)  │ ✅ 已操作 │ GET /api/v1/users/1 │
│ 4  │ Icon     │ [編集]图标       │ ✅ 已操作 │ GET /api/v1/users/1 │
│ 5  │ Input    │ [検索]输入框     │ ✅ 已操作 │ 无API（需配合搜索按钮）│
└────┴──────────┴─────────────────┴──────────┴─────────────────────┘
```

#### 2.2.4 用例描述中的标识示例

**生成的用例screen字段格式：**

```json
{
  "screen": "[ユーザー管理]画面",
  "url": "/api/v1/users",
  "method": "POST",
  "function": "创建用户 - 通过[新規作成]按钮触发",
  "precondition": "已登录系统，进入[ユーザー管理]画面",
  "test_steps": "1. 点击[新規作成]按钮\n2. 在弹窗中填写用户信息\n3. 点击[保存]按钮\n4. 验证API响应为201",
  "expected_result": "返回状态码201，显示[作成成功]消息"
}
```

### 2.3 语言与多语言处理

自动识别网站显示语言，生成对应语言的描述性文本：

- 日语网站 → 用例描述使用日语
- 中文网站 → 用例描述使用中文
- 英文网站 → 用例描述使用英语
- **注意**：URL、HTTP方法、JSON字段名等技术标识符保持原样

### 2.4 CRUD完整覆盖原则

- **一接口多用例**：同一接口不同响应码场景拆分为独立用例
- **响应码覆盖**：200、201、400、401、403、404、500
- **🚨 CRUD完整覆盖原则**：
  - **GET(查询)**: 列表查询、详情查询、条件查询、分页查询
  - **POST(创建)**: 创建→验证→删除（避免垃圾数据）
  - **PUT/PATCH(更新)**: 创建测试数据→修改→验证→删除（不修改真实数据）
  - **DELETE(删除)**: 创建测试数据→删除→验证（不删除真实数据）
  - **开关切换**: 找OFF数据→ON→验证→恢复OFF（保持原状态）

### 2.5 测试数据管理规则

#### 2.5.1 数据管理核心原则

**⚠️ 绝对要求：不操作现有业务数据，只操作脚本自己创建的测试数据**

| 用例类型        | 脚本实际执行的操作              | 数据清理策略         | 说明               |
| ----------- | ---------------------- | -------------- | ---------------- |
| GET 查询      | 直接查询                   | 无需清理           | 只读操作，不影响数据       |
| POST 创建     | 创建 → 验证 → 删除           | 🚨 立即删除         | 验证创建功能后立即删除测试数据  |
| PUT 修改      | 创建 → 修改 → 验证 → 删除      | 🚨 立即删除         | 不修改现有数据，创建专用测试数据 |
| DELETE 删除   | 创建 → 删除 → 验证           | 已删除，无需额外清理     | 不删除现有数据，创建后再删除   |
| PATCH 开关ON  | 找OFF数据 → ON → 验证 → OFF | 🚨 恢复原状态        | 恢复为原始OFF状态       |
| PATCH 开关OFF | 找ON数据 → OFF → 验证 → ON  | 🚨 恢复原状态        | 恢复为原始ON状态        |

**为何要如此严格？**

1. **保护生产数据**：避免测试过程中误删、误改真实业务数据
2. **可重复执行**：每次执行都创建新的测试数据，不依赖环境状态
3. **无副作用**：测试执行前后，系统数据状态保持一致
4. **隔离性**：不同测试用例之间互不干扰

#### 2.5.2 script_code必须使用真实可执行数据

**script_code中的路径参数、请求体数据必须来自实际捕获的请求，确保脚本可直接执行成功：**

```
✅ 正确做法：
- **URL中的ID**：使用探索脚本捕获到的真实ID
- 请求体：使用捕获的请求中的真实数据结构和值
- Token：使用实际登录后获取的有效Token

❌ 禁止做法：
- 使用虚构的ID（如 /api/user/99999）
- 编造请求体字段（未在实际请求中出现的字段）
- 使用过期或无效的Token
```

**数据来源优先级**：

1. **探索脚本捕获**：从Playwright脚本实时拦截的API请求中提取真实数据
2. **页面观察**：从脚本输出的JSON数据中提取列表第一行的真实ID
3. **元数据凭证**：登录接口使用 `get_api_group_metadata` 返回的 user/password

#### 2.5.3 CRUD用例script_code生成规范

下面是各类CRUD操作的标准脚本模板，生成用例时必须遵循这些模板：

```javascript
// ✅ GET查询用例 - 直接查询（无需清理）
async (page) => {
    const loginRes = await page.request.post('${base_url}/api/v1/auth/login', {
        data: { username: '${username}', password: '${password}' },
        ignoreHTTPSErrors: true
    });
    const loginData = await loginRes.json();
    const token = loginData.data?.token || loginData.token;
    if (!token) return { passed: false, error: 'Login failed' };

    const res = await page.request.get('${base_url}/api/users', {
        headers: { 'Authorization': 'Bearer ' + token },
        ignoreHTTPSErrors: true
    });

    return { passed: res.status() === 200, status: res.status() };
}
```
// ✅ POST创建用例 - 创建→验证→删除
async (page) => {
    // 1. 🚨 先登录获取Token
    const loginRes = await page.request.post('${base_url}/api/v1/auth/login', {
        data: { username: '${username}', password: '${password}' },
        ignoreHTTPSErrors: true
    });
    const loginData = await loginRes.json();
    const token = loginData.data?.token || loginData.token;
    if (!token) return { passed: false, error: 'Login failed' };

    // 2. 执行创建
    const createRes = await page.request.post('${base_url}/api/users', {
        headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
        data: { username: 'test_auto_' + Date.now(), password: 'Test123!' },
        ignoreHTTPSErrors: true
    });
    const created = await createRes.json();
    const testPassed = createRes.status() === 201;

    // 3. 🚨 删除刚创建的数据
    if (created.data?.id || created.id) {
        await page.request.delete('${base_url}/api/users/' + (created.data?.id || created.id), {
            headers: { 'Authorization': 'Bearer ' + token },
            ignoreHTTPSErrors: true
        });
    }

    return { passed: testPassed, status: createRes.status(), cleaned: true };
}
```

```javascript
// ✅ PUT修改用例 - 创建→修改→验证→删除（不修改现有数据）
async (page) => {
    // 1. 🚨 先登录获取Token
    const loginRes = await page.request.post('${base_url}/api/v1/auth/login', {
        data: { username: '${username}', password: '${password}' },
        ignoreHTTPSErrors: true
    });
    const loginData = await loginRes.json();
    const token = loginData.data?.token || loginData.token;
    if (!token) return { passed: false, error: 'Login failed' };

    // 2. 先创建专用测试数据
    const createRes = await page.request.post('${base_url}/api/users', {
        headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
        data: { username: 'test_auto_' + Date.now(), nickname: 'before' },
        ignoreHTTPSErrors: true
    });
    const created = await createRes.json();
    const testId = created.data?.id || created.id;
    if (!testId) return { passed: false, error: 'Create test data failed' };

    // 3. 修改刚创建的数据
    const updateRes = await page.request.put('${base_url}/api/users/' + testId, {
        headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
        data: { nickname: 'after_modified' },
        ignoreHTTPSErrors: true
    });
    const testPassed = updateRes.status() === 200;

    // 4. 🚨 删除测试数据
    await page.request.delete('${base_url}/api/users/' + testId, {
        headers: { 'Authorization': 'Bearer ' + token },
        ignoreHTTPSErrors: true
    });

    return { passed: testPassed, status: updateRes.status(), cleaned: true };
}
```

```javascript
// ✅ DELETE删除用例 - 创建→删除→验证（不删除现有数据）
async (page) => {
    // 1. 🚨 先登录获取Token
    const loginRes = await page.request.post('${base_url}/api/v1/auth/login', {
        data: { username: '${username}', password: '${password}' },
        ignoreHTTPSErrors: true
    });
    const loginData = await loginRes.json();
    const token = loginData.data?.token || loginData.token;
    if (!token) return { passed: false, error: 'Login failed' };

    // 2. 先创建专用测试数据（专门用于删除测试）
    const createRes = await page.request.post('${base_url}/api/users', {
        headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
        data: { username: 'test_auto_delete_' + Date.now() },
        ignoreHTTPSErrors: true
    });
    const created = await createRes.json();
    const testId = created.data?.id || created.id;
    if (!testId) return { passed: false, error: 'Create test data failed' };

    // 3. 删除刚创建的数据
    const deleteRes = await page.request.delete('${base_url}/api/users/' + testId, {
        headers: { 'Authorization': 'Bearer ' + token },
        ignoreHTTPSErrors: true
    });

    return { passed: deleteRes.status() === 200 || deleteRes.status() === 204, status: deleteRes.status() };
}
```

```javascript
// ✅ 开关ON测试 - 找OFF数据→ON→验证→OFF（恢复原状态）
async (page) => {
    // 1. 🚨 先登录获取Token
    const loginRes = await page.request.post('${base_url}/api/v1/auth/login', {
        data: { username: '${username}', password: '${password}' },
        ignoreHTTPSErrors: true
    });
    const loginData = await loginRes.json();
    const token = loginData.data?.token || loginData.token;
    if (!token) return { passed: false, error: 'Login failed' };

    // 2. OFF → ON（测试开启功能）
    const enableRes = await page.request.patch('${base_url}/api/items/${test_off_item_id}/enable', {
        headers: { 'Authorization': 'Bearer ' + token },
        ignoreHTTPSErrors: true
    });
    const testPassed = enableRes.status() === 200;

    // 3. 🚨 ON → OFF（恢复原状态）
    await page.request.patch('${base_url}/api/items/${test_off_item_id}/disable', {
        headers: { 'Authorization': 'Bearer ' + token },
        ignoreHTTPSErrors: true
    });

    return { passed: testPassed, status: enableRes.status(), restored: true };
}
```

```javascript
// ✅ 开关OFF测试 - 找ON数据→OFF→验证→ON（恢复原状态）
async (page) => {
    // 1. 🚨 先登录获取Token
    const loginRes = await page.request.post('${base_url}/api/v1/auth/login', {
        data: { username: '${username}', password: '${password}' },
        ignoreHTTPSErrors: true
    });
    const loginData = await loginRes.json();
    const token = loginData.data?.token || loginData.token;
    if (!token) return { passed: false, error: 'Login failed' };

    // 2. ON → OFF（测试关闭功能）
    const disableRes = await page.request.patch('${base_url}/api/items/${test_on_item_id}/disable', {
        headers: { 'Authorization': 'Bearer ' + token },
        ignoreHTTPSErrors: true
    });
    const testPassed = disableRes.status() === 200;

    // 3. 🚨 OFF → ON（恢复原状态）
    await page.request.patch('${base_url}/api/items/${test_on_item_id}/enable', {
        headers: { 'Authorization': 'Bearer ' + token },
        ignoreHTTPSErrors: true
    });

    return { passed: testPassed, status: disableRes.status(), restored: true };
}
```

**数据管理规则汇总**：
| 操作类型 | 脚本流程 | 说明 |
|---------|------------|--------|
| GET 查询 | 直接查询 | 无需清理 |
| POST 创建 | 创建 → 验证 → DELETE | 创建后必须删除 |
| PUT 修改 | POST → PUT → 验证 → DELETE | 创建测试数据后修改，最后删除 |
| DELETE 删除 | POST → DELETE → 验证 | 创建测试数据后删除 |
| PATCH 开关 | 找相反状态 → 切换 → 验证 → 恢复 | 必须恢复原状态 |

### 2.6 完整输出规则

- **画面完整遍历**：必须遍历网站的**所有主要画面**，不得只做部分画面就结束。典型网站应覆盖：登录、Dashboard、各功能模块列表页、详情页、设置页等

- **API全量覆盖**：每个画面中**实时捕获的所有API接口**都必须生成测试用例，不得遗漏

- **用例数量参考基准**：
  
  | 网站规模 | 画面数   | 预期用例数    |
  | ---- | ----- | -------- |
  | 小型   | 5-10  | 50-100条  |
  | 中型   | 10-20 | 100-200条 |
  | 大型   | 20+   | 200+条    |
  
  **如果生成的用例数量明显偏少，必须检查是否遗漏了画面或接口**

- **🚨 强制继续机制（最重要）**：
  
  **触发条件（满足任一即触发）**：
  
  1. 还有画面未遍历完成
  2. 当前画面的API未全部生成用例
  3. 单次输出即将达到token限制
  4. 已生成用例数量未达到预期基准
  
  **仅在以下情况时暂停并提示用户输入【继续】**：
  
  1. 输出token接近限制（约80%使用量）
  2. 任务执行时间超过30分钟（避免超时）
  3. 已生成用例数量超过100条（建议分批写入）
  4. AI判断需要中断优化（如错误率过高）
  
  **⚠️ 除上述情况外，必须持续自动执行，不得中断！**
  
  ```
  ⏸️ API用例生成进度报告
  
  ✅ 已完成画面：
  - [ログイン] - 8条用例 ✓
  - [ダッシュボード] - 12条用例 ✓
  
  ⏳ 待处理画面：
  - [ライセンス一覧] - 预计15条
  - [ファイル管理] - 预计10条
  - [設定] - 预计8条
  
  📊 当前进度：20/60条（33%）
  
  👉 请输入【继续】生成剩余画面的用例
  ```
  
  **⚠️ 严禁行为**：
  
  - ❌ 在未遍历完所有画面时输出"完成"报告
  - ❌ 跳过画面直接结束
  - ❌ 只捕获部分API就认为画面完成
  - ❌ 在输出token不足时直接截断而不提示继续
  - ❌ 不遍历控件就认为画面API采集完成

- **完成确认**：**只有当所有画面都遍历完成后**，才输出最终汇总报告：
  
  ```
  ✅ API用例生成完成！
  
  📊 生成统计：
  - 总画面数：12个
  - 总控件数：86个（已操作82个，跳过4个）
  - 总API数：45个
  - 总用例数：156条（正向98/反向58）
  
  📋 各画面用例分布：
  - [ログイン]: 8条 ✓ (控件: 5/5)
  - [ダッシュボード]: 12条 ✓ (控件: 8/8)
  - [ライセンス一覧]: 15条 ✓ (控件: 12/12)
  ...
  
  🎉 全部画面、全部控件遍历完成，任务结束！
  ```

- **🚨 画面控件全覆盖（强制要求 + CRUD原则）**：
  
  每个画面必须按**CRUD顺序**识别并操作**所有可交互控件**：
  
  **🎯 遍历顺序（严格执行）：**
  
  ```
  1️⃣ C-创建(CREATE)
     → [新增]、[创建]、[添加]等按钮 → 填表 → [保存]
  
  2️⃣ R-检索(READ)  
     → [查询]、[搜索]按钮、搜索框、表格行、[详情]链接
  
  3️⃣ U-修改(UPDATE)
     → [编辑]、[修改]按钮 → 改值 → [保存]、开关切换
  
  4️⃣ D-删除(DELETE)
     → [删除]按钮 → [取消]确认（不真删）
  
  5️⃣ 其他控件
     → 下拉筛选、标签页、分页器、复选框
  ```
  
  **控件类型详解：**
  
  | 控件类型        | 所属CRUD | 常见形态 | 容易遗漏的场景                |
  | ----------- | ------- | ---- | ---------------------- |
  | Button(新增)  | **C**   | 按钮   | [新規作成]、[追加]、对话框内的[保存] |
  | Button(查询)  | **R**   | 按钮   | [検索]、[照会]、[リセット]      |
  | Link        | **R**   | 文字链接 | [詳細]、[查看]、面包屑导航        |
  | Button(编辑)  | **U**   | 按钮   | 表格行内的[編集]图标🖊️          |
  | Switch      | **U**   | 开关   | [有効/無効]切换              |
  | Button(删除)  | **D**   | 按钮   | [削除]图标🗑️（点击后取消确认）     |
  | Input       | C/R    | 输入框  | 搜索框（输入后回车）、表单字段        |
  | Select      | R      | 下拉框  | 状态筛选、类型筛选              |
  | Tab         | -      | 标签页  | 切换后继续遍历该标签内控件          |
  | Pagination  | R      | 分页   | 下一页触发列表查询API           |
  | Table Row   | R      | 表格行  | 点击行查看详情                |

## 3. 数据结构定义

### 3.1 API用例7字段结构

| 字段名         | 说明          | 示例                                  |
| ----------- | ----------- | ----------------------------------- |
| screen      | 画面名称（用[]标识） | [ログイン]                              |
| url         | 接口路径（不含域名）  | /api/auth/login                     |
| method      | HTTP方法      | GET, POST, PUT, DELETE              |
| header      | 请求头JSON     | {"Authorization": "Bearer {token}"} |
| body        | 请求体JSON     | {"username": "admin"}               |
| response    | 期望响应        | {"code": 200}                       |
| script_code | 可执行的JS测试脚本  | 见下方模板                               |

### 3.2 script_code 字段生成规则

**script_code 用于后续自动执行测试和性能测试，必须为每个用例生成：**

#### 脚本格式说明

**API脚本使用 `async (page) => {}` 格式配合 Playwright 的 `page.request` API。**

**为什么使用 page.request？**

1. **原生API支持**：Playwright提供的原生HTTP请求API，无需浏览器上下文
2. **HTTPS证书跳过**：支持 `ignoreHTTPSErrors: true` 参数，可直接跳过自签名证书验证
3. **更简洁高效**：无需page.evaluate包装，代码更直观
4. **Node执行兼容**：在 node + playwright 环境中运行，自动处理证书

**脚本结构：**

```javascript
// ✅ 推荐：使用 page.request API（更简洁）
async (page) => {
  // 1. 登录获取Token
  const loginRes = await page.request.post('${base_url}/api/v1/auth/login', {
    data: { username: '${username}', password: '${password}' },
    ignoreHTTPSErrors: true  // 🔐 跳过HTTPS证书验证
  });
  const token = (await loginRes.json()).data?.token;
  
  // 2. 执行API请求
  const res = await page.request.get('${base_url}/api/users', {
    headers: { 'Authorization': 'Bearer ' + token },
    ignoreHTTPSErrors: true  // 🔐 跳过HTTPS证书验证
  });
  
  return { passed: res.status() === 200, status: res.status() };
}
```

**⚠️ 注意事项：**
- 每个请求都需要添加 `ignoreHTTPSErrors: true` 参数（当目标系统使用自签名证书时）
- 使用 `res.status()` 获取状态码（注意是方法调用，不是属性）
- 使用 `await res.json()` 解析响应体

#### 🚨 script_code 脚本独立原则

**每个script_code必须完全独立可执行，包含登录获取Token的完整流程：**

```javascript
// ✅ 正确：脚本自行登录获取Token，使用 page.request API
async (page) => {
  // 1. 先登录获取Token（每个脚本独立获取）
  const loginRes = await page.request.post('${base_url}/api/v1/auth/login', {
    data: { username: '${username}', password: '${password}' },
    ignoreHTTPSErrors: true  // 🔐 跳过HTTPS证书验证（自签名证书）
  });
  const loginData = await loginRes.json();
  const token = loginData.data?.token || loginData.token;
  if (!token) return { passed: false, error: 'Login failed' };

  // 2. 使用获取的token执行实际测试
  const res = await page.request.get('${base_url}/api/users', {
    headers: { 'Authorization': 'Bearer ' + token },
    ignoreHTTPSErrors: true  // 🔐 跳过HTTPS证书验证
  });
  return { passed: res.status() === 200, status: res.status() };
}

// ❌ 错误：依赖外部token变量（不独立）
async (page) => {
  const res = await page.request.get('${base_url}/api/users', {
    headers: { 'Authorization': 'Bearer ${token}' },  // 依赖变量表中的token，token会过期
    ignoreHTTPSErrors: true
  });
  return { passed: res.status() === 200, status: res.status() };
}

// ❌ 错误：硬编码具体值
async (page) => {
  const res = await page.request.get('https://example.com:443/api/users', {  // 硬编码URL
    headers: { 'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIs...' },  // 硬编码Token
    ignoreHTTPSErrors: true
  });
  return { passed: res.status() === 200, status: res.status() };
}
```

**可用的变量占位符：**

| 占位符           | 来源      | 说明      |
| ------------- | ------- | ------- |
| `${base_url}` | 元数据自动生成 | 完整基础URL |
| `${username}` | 元数据     | 登录用户名   |
| `${password}` | 元数据     | 登录密码    |
| `${user_id}`  | 实时捕获   | 当前用户ID  |
| `${test_id}`  | 实时捕获   | 测试数据ID  |
| `${自定义变量}`    | 动态写入    | 运行时动态变量 |

> 🚨 **重要**：`token` 不写入变量表！每个脚本必须自行调用登录接口获取Token，确保脚本完全独立可执行。

#### 正向用例模板（需要认证）

```javascript
// {screen} - {method} {url} - 正常场景
async (page) => {
  // 1. 先登录获取Token
  const loginRes = await page.request.post('${base_url}/api/v1/auth/login', {
    data: { username: '${username}', password: '${password}' },
    ignoreHTTPSErrors: true
  });
  const loginData = await loginRes.json();
  const token = loginData.data?.token || loginData.token;
  if (!token) return { passed: false, error: 'Login failed' };

  // 2. 执行实际测试
  const res = await page.request.{method}('${base_url}{url}', {
    headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
    data: {body_or_null},
    ignoreHTTPSErrors: true
  });
  return { passed: res.status() === {expected_status}, status: res.status() };
}
```

#### 反向用例模板（无Token场景）

```javascript
// {screen} - {method} {url} - 无Token访问被拒绝
async (page) => {
  const res = await page.request.{method}('${base_url}{url}', {
    headers: { 'Content-Type': 'application/json' },  // 🚨 无Authorization头
    ignoreHTTPSErrors: true
  });
  return { passed: res.status() === 401, status: res.status() };
}
```

#### 反向用例模板（无效Token场景）

```javascript
// {screen} - {method} {url} - 无效Token被拒绝
async (page) => {
  const res = await page.request.{method}('${base_url}{url}', {
    headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer invalid_token_12345' },
    ignoreHTTPSErrors: true
  });
  return { passed: res.status() === 401, status: res.status() };
}
```

**生成规则：**

- 将用例的 url/method/header/body/response 信息嵌入脚本

- `{expected_status}` 从 response 字段中提取状态码

- GET/DELETE 请求不需要 body 参数

- 脚本必须可独立执行，便于后续批量运行和性能测试

- **🚨 Token使用规则**：
  
  | 用例场景    | Authorization头          | 期望状态码   |
  | ------- | ----------------------- | ------- |
  | 正常访问    | `Bearer ` + 脚本内获取的token | 200/201 |
  | 无Token  | 不传                      | 401     |
  | 无效Token | `Bearer invalid_token`  | 401     |
  | 权限不足    | 使用低权限用户登录获取的token       | 403     |
  
  > 🚨 **脚本独立原则**：每个脚本必须自行调用登录API获取Token，不依赖变量表中的token值

### 3.3 🚨 路径参数处理规则（重要）

**URL字段与script_code字段的参数处理方式不同：**

| 字段              | 处理方式         | 目的        |
| --------------- | ------------ | --------- |
| **url**         | 保留占位符 `{id}` | 便于理解API结构 |
| **script_code** | 替换为真实值       | 便于直接执行测试  |

**示例：**

```
捕获到的请求：GET /api/softsim/440070700060217 - 200

生成的用例：
{
  "url": "/api/softsim/{imsi}",                    ← 保留占位符，便于理解
  "script_code": "...page.request.get('${base_url}/api/softsim/440070700060217'..."  ← 使用真实值，便于执行
}
```

**占位符命名规范：**

- 数字ID → `{id}`
- 用户ID → `{userId}`
- IMSI号 → `{imsi}`
- 项目ID → `{projectId}`
- UUID → `{uuid}`

**真实值获取方法（CLI版）：**

1. 从**实时捕获的请求URL**中提取
2. 如果无法获取，使用页面上显示的数据（如列表第一行的ID）
3. 确保 script_code 中的值在目标系统中真实存在

**完整示例：**

```json
{
  "screen": "[端末情報]",
  "url": "/api/softsim/{imsi}",
  "method": "GET",
  "header": "{\"Authorization\": \"Bearer {token}\"}",
  "body": "",
  "response": "{\"code\": 200}",
  "script_code": "async (page) => { const loginRes = await page.request.post('${base_url}/api/v1/auth/login', { data: { username: '${username}', password: '${password}' }, ignoreHTTPSErrors: true }); const loginData = await loginRes.json(); const token = loginData.data?.token || loginData.token; if (!token) return { passed: false, error: 'Login failed' }; const res = await page.request.get('${base_url}/api/softsim/440070700060217', { headers: { 'Authorization': 'Bearer ' + token }, ignoreHTTPSErrors: true }); return { passed: res.status() === 200, status: res.status() }; }"
}
```

### 3.4 字段填写规范

- **remark字段必须留空**（由执行阶段填写）
- **URL字段**：只填Path部分，如 `/api/version`
- **Header字段**：无需认证填 `{}`，需Token填 `{"Authorization": "Bearer ${token}"}`

## 4. 全自动化工作流程

### 第一步：获取项目和用例集信息

#### 1.1 获取当前项目

```
get_current_project_name()
```

**执行后输出**：

```
✅ 当前项目：
- 项目ID：{project_id}
- 项目名称：{project_name}
```

#### 1.2 获取API用例集列表

```
list_api_groups(project_id={上一步获取的project_id})
```

**执行后输出**：

```
✅ API用例集列表：
| ID | 用例集名称 | 目标服务器 |
|----|----------|----------|
| 45 | apitest  | 192.168.50.32:8443 |

🎯 目标用例集：{{group_name}}
```

#### 1.3 获取用例集元数据

```
get_api_group_metadata(group_name="{{group_name}}")
```

**执行后输出**：

```
✅ 用例集元数据：
- 用例集ID：{group_id}
- 协议：{meta_protocol}
- 服务器：{meta_server}
- 端口：{meta_port}
- 用户名：{meta_user}
- 密码：{meta_password}
- BASE_URL：{meta_protocol}://{meta_server}:{meta_port}
```

---

### 第二步：AI 编写并执行探索脚本

**🚨 核心：AI 自动编写 Playwright 脚本，完全自动化探索！**

#### 2.1 创建探索脚本

**AI 主动创建以下文件：**

**文件 1: `package.json`**

```json
{
  "name": "api-explorer",
  "version": "1.0.0",
  "dependencies": {
    "playwright": "^1.40.0"
  }
}
```

**文件 2: `explore-api.js`** (核心探索脚本)

```javascript
const { chromium } = require('playwright');
const fs = require('fs');

(async () => {
  const browser = await chromium.launch({
    headless: true,
    ignoreHTTPSErrors: true  // 跳过HTTPS证书
  });

  const context = await browser.newContext();
  const page = await context.newPage();

  const BASE_URL = process.env.BASE_URL;
  const USERNAME = process.env.USERNAME;
  const PASSWORD = process.env.PASSWORD;

  // ===== 数据收集容器 =====
  const apiRequests = [];
  const result = {
    baseUrl: BASE_URL,
    timestamp: new Date().toISOString(),
    apis: [],
    screens: []
  };

  // ===== 核心：实时拦截所有 API 请求 =====
  const requestMap = new Map();  // 存储请求信息

  context.on('request', request => {
    const resourceType = request.resourceType();
    
    // 只捕获 XHR 和 Fetch 请求（API 请求）
    if (resourceType === 'xhr' || resourceType === 'fetch') {
      const requestData = {
        id: Date.now() + Math.random(),
        url: request.url(),
        method: request.method(),
        headers: request.headers(),
        body: request.postData() || null,
        timestamp: new Date().toISOString()
      };
      
      requestMap.set(request.url() + request.method(), requestData);
      console.log(`[Request] ${request.method()} ${request.url()}`);
    }
  });

  context.on('response', async response => {
    const request = response.request();
    const resourceType = request.resourceType();
    
    if (resourceType === 'xhr' || resourceType === 'fetch') {
      const key = request.url() + request.method();
      const requestData = requestMap.get(key);
      
      if (requestData) {
        try {
          const responseBody = await response.text().catch(() => null);
          
          const apiData = {
            ...requestData,
            status: response.status(),
            statusText: response.statusText(),
            responseHeaders: response.headers(),
            responseBody: responseBody
          };
          
          apiRequests.push(apiData);
          console.log(`[Response] ${response.status()} ${request.method()} ${request.url()}`);
        } catch (e) {
          console.warn(`  ⚠️  Failed to capture response: ${e.message}`);
        }
      }
    }
  });

  try {
    console.log('========== API EXPLORATION START ==========');
    console.log('🔐 [Step 1] 自动登录...');
    
    // 访问登录页面
    await page.goto(BASE_URL, { waitUntil: 'domcontentloaded', timeout: 15000 });
    await page.waitForTimeout(2000);

    // 智能识别登录表单
    const usernameSelectors = [
      'input[type="text"]',
      'input[name*="user"]',
      'input[name*="login"]',
      'input[placeholder*="用户"]',
      'input[placeholder*="ユーザー"]',
      'input[placeholder*="User"]'
    ];

    for (const selector of usernameSelectors) {
      try {
        const element = page.locator(selector).first();
        if (await element.isVisible({ timeout: 500 })) {
          await element.fill(USERNAME);
          console.log(`  ✅ 用户名字段: ${selector}`);
          break;
        }
      } catch (e) { /* 继续 */ }
    }

    await page.locator('input[type="password"]').first().fill(PASSWORD);
    console.log('  ✅ 密码已填写');

    // 点击登录按钮
    const loginButtonSelectors = [
      'button[type="submit"]',
      'input[type="submit"]',
      'button:has-text("登录")',
      'button:has-text("ログイン")',
      'button:has-text("Login")'
    ];

    for (const selector of loginButtonSelectors) {
      try {
        const button = page.locator(selector).first();
        if (await button.isVisible({ timeout: 500 })) {
          await button.click();
          console.log(`  ✅ 登录按钮已点击`);
          break;
        }
      } catch (e) { /* 继续 */ }
    }

    await page.waitForTimeout(3000);
    console.log('✅ 登录成功');

    console.log('\\n📊 [Step 2] 扫描菜单结构...');

    // 获取所有导航链接
    const menuSelectors = ['nav a', '.menu a', '.sidebar a', 'header a', '[role="menuitem"]'];
    let allMenuItems = [];

    for (const selector of menuSelectors) {
      try {
        const items = await page.locator(selector).evaluateAll(links => {
          return links
            .filter(link => link.offsetParent !== null && link.textContent.trim())
            .map(link => ({ text: link.textContent.trim(), href: link.href }));
        });
        allMenuItems = allMenuItems.concat(items);
      } catch (e) { /* 继续 */ }
    }

    // 去重
    const uniqueMenus = Array.from(new Map(allMenuItems.map(item => [item.href, item])).values());
    console.log(`找到 ${uniqueMenus.length} 个唯一菜单项`);

    console.log('\\n🔍 [Step 3] 遍历画面并自动操作（遵循CRUD原则）...');
    console.log('  📝 CRUD操作顺序: C创建 → R检索 → U修改 → D删除');

    // 遍历每个菜单（最多15个）
    const maxScreens = Math.min(uniqueMenus.length, 15);
    
    for (let i = 0; i < maxScreens; i++) {
      const menu = uniqueMenus[i];

      if (!menu.href || 
          menu.href.includes('logout') || 
          menu.href.includes('javascript:') || 
          menu.href === BASE_URL) {
        continue;
      }

      try {
        console.log(`\\n[${i+1}/${maxScreens}] 访问画面: "${menu.text}"`);
        await page.goto(menu.href, { waitUntil: 'domcontentloaded', timeout: 10000 });
        await page.waitForTimeout(2000);

        const screenInfo = {
          index: i + 1,
          name: menu.text,
          url: page.url(),
          capturedApis: 0,
          controlActions: []
        };

        // === 按CRUD顺序自动识别并操作所有可交互控件 ===
        console.log(`  🎯 自动操作控件（CRUD顺序）...`);

        // === 阶段1：C - 创建操作（CREATE） ===
        console.log(`  \\n  📝 [C-CREATE] 查找并触发创建相关控件...`);
        const createKeywords = ['新增', '创建', '添加', '新規', '作成', '追加', 'Create', 'Add', 'New'];
        const buttons = await page.locator('button, input[type="submit"], a.btn').all();
        
        for (let btnIdx = 0; btnIdx < buttons.length && btnIdx < 10; btnIdx++) {
          try {
            const button = buttons[btnIdx];
            const buttonText = await button.textContent().catch(() => '');
            const trimmedText = buttonText.trim();
            
            // 优先处理创建按钮
            if (createKeywords.some(kw => trimmedText.includes(kw))) {
              if (await button.isVisible({ timeout: 500 })) {
                console.log(`    ✓ [C] 创建按钮: [${trimmedText}]`);
                screenInfo.controlActions.push({ type: 'CREATE', control: trimmedText });
                await button.click({ timeout: 3000 });
                await page.waitForTimeout(1500);
                
                // 如果弹出表单，填写并提交
                const dialogInputs = await page.locator('input:visible:not([type="hidden"])').count();
                if (dialogInputs > 0) {
                  console.log(`      → 检测到表单，填写 ${dialogInputs} 个字段...`);
                  const inputs = await page.locator('input:visible:not([type="hidden"])').all();
                  for (let inp of inputs.slice(0, 8)) {
                    try {
                      const inputType = await inp.getAttribute('type').catch(() => 'text');
                      let testValue = 'test_' + Date.now();
                      if (inputType === 'email') testValue = `test${Date.now()}@example.com`;
                      if (inputType === 'number') testValue = '123';
                      await inp.fill(testValue);
                      await page.waitForTimeout(300);
                    } catch (e) { /* 继续 */ }
                  }
                  
                  // 查找保存/确认按钮
                  const saveKeywords = ['保存', '确认', '提交', '保存', '確認', '送信', 'Save', 'Confirm', 'Submit'];
                  const dialogButtons = await page.locator('button:visible, input[type="submit"]:visible').all();
                  for (let saveBtn of dialogButtons) {
                    const saveBtnText = await saveBtn.textContent().catch(() => '');
                    if (saveKeywords.some(kw => saveBtnText.includes(kw))) {
                      console.log(`      → 点击保存按钮: [${saveBtnText.trim()}]`);
                      await saveBtn.click({ timeout: 3000 });
                      await page.waitForTimeout(2000);
                      break;
                    }
                  }
                }
              }
            }
          } catch (e) {
            console.warn(`    ⚠️  创建操作失败: ${e.message.slice(0, 50)}`);
          }
        }

        // === 阶段2：R - 检索操作（READ/RETRIEVE） ===
        console.log(`  \\n  🔍 [R-READ] 查找并触发检索相关控件...`);
        const readKeywords = ['查询', '搜索', '检索', '查看', '详情', '検索', '照会', '詳細', 'Search', 'Query', 'View', 'Detail'];
        
        // 2.1 检索按钮
        for (let btnIdx = 0; btnIdx < buttons.length && btnIdx < 10; btnIdx++) {
          try {
            const button = buttons[btnIdx];
            const buttonText = await button.textContent().catch(() => '');
            const trimmedText = buttonText.trim();
            
            if (readKeywords.some(kw => trimmedText.includes(kw))) {
              if (await button.isVisible({ timeout: 500 })) {
                console.log(`    ✓ [R] 检索按钮: [${trimmedText}]`);
                screenInfo.controlActions.push({ type: 'READ', control: trimmedText });
                await button.click({ timeout: 3000 });
                await page.waitForTimeout(1500);
              }
            }
          } catch (e) { /* 继续 */ }
        }

        // 2.2 填写搜索框
        const searchInputs = await page.locator('input[type="text"]:visible, input[type="search"]:visible').all();
        for (let input of searchInputs.slice(0, 3)) {
          try {
            const placeholder = await input.getAttribute('placeholder').catch(() => '');
            if (readKeywords.some(kw => placeholder.includes(kw)) || placeholder.includes('search')) {
              console.log(`    ✓ [R] 搜索框: [${placeholder}]`);
              screenInfo.controlActions.push({ type: 'READ', control: '搜索框' });
              await input.fill('test');
              await page.waitForTimeout(500);
              await page.keyboard.press('Enter');
              await page.waitForTimeout(1500);
            }
          } catch (e) { /* 继续 */ }
        }

        // 2.3 点击表格行/列表项（查看详情）
        const tableRows = await page.locator('tbody tr, .list-item, .table-row').count();
        if (tableRows > 0) {
          try {
            const firstRow = page.locator('tbody tr, .list-item, .table-row').first();
            console.log(`    ✓ [R] 点击第一行数据（查看详情）`);
            screenInfo.controlActions.push({ type: 'READ', control: '表格行详情' });
            await firstRow.click({ timeout: 3000 });
            await page.waitForTimeout(1500);
            await page.goBack({ timeout: 5000 }).catch(() => {});
          } catch (e) { /* 继续 */ }
        }

        // === 阶段3：U - 修改操作（UPDATE） ===
        console.log(`  \\n  ✏️ [U-UPDATE] 查找并触发修改相关控件...`);
        const updateKeywords = ['编辑', '修改', '更新', '変更', '編集', '更新', 'Edit', 'Update', 'Modify'];
        
        for (let btnIdx = 0; btnIdx < buttons.length && btnIdx < 10; btnIdx++) {
          try {
            const button = buttons[btnIdx];
            const buttonText = await button.textContent().catch(() => '');
            const trimmedText = buttonText.trim();
            
            if (updateKeywords.some(kw => trimmedText.includes(kw))) {
              if (await button.isVisible({ timeout: 500 })) {
                console.log(`    ✓ [U] 修改按钮: [${trimmedText}]`);
                screenInfo.controlActions.push({ type: 'UPDATE', control: trimmedText });
                await button.click({ timeout: 3000 });
                await page.waitForTimeout(1500);
                
                // 如果弹出编辑表单，修改字段
                const editInputs = await page.locator('input:visible:not([type="hidden"])').count();
                if (editInputs > 0) {
                  console.log(`      → 修改表单字段...`);
                  const inputs = await page.locator('input:visible:not([type="hidden"])').all();
                  for (let inp of inputs.slice(0, 3)) {
                    try {
                      await inp.fill('modified_' + Date.now());
                      await page.waitForTimeout(300);
                    } catch (e) { /* 继续 */ }
                  }
                  
                  // 点击保存
                  const saveKeywords = ['保存', '确认', '保存', '確認', 'Save', 'Confirm'];
                  const saveButtons = await page.locator('button:visible').all();
                  for (let saveBtn of saveButtons) {
                    const saveBtnText = await saveBtn.textContent().catch(() => '');
                    if (saveKeywords.some(kw => saveBtnText.includes(kw))) {
                      console.log(`      → 保存修改: [${saveBtnText.trim()}]`);
                      await saveBtn.click({ timeout: 3000 });
                      await page.waitForTimeout(2000);
                      break;
                    }
                  }
                }
              }
            }
          } catch (e) {
            console.warn(`    ⚠️  修改操作失败: ${e.message.slice(0, 50)}`);
          }
        }

        // 2.4 开关状态切换（特殊的UPDATE）
        const switches = await page.locator('input[type="checkbox"]:visible, .switch, .toggle').all();
        for (let sw of switches.slice(0, 3)) {
          try {
            console.log(`    ✓ [U] 切换开关`);
            screenInfo.controlActions.push({ type: 'UPDATE', control: '开关切换' });
            await sw.click({ timeout: 2000 });
            await page.waitForTimeout(1500);
          } catch (e) { /* 继续 */ }
        }

        // === 阶段4：D - 删除操作（DELETE） ===
        console.log(`  \\n  🗑️  [D-DELETE] 查找并触发删除相关控件...`);
        const deleteKeywords = ['删除', '刪除', 'Delete', 'Remove'];
        
        for (let btnIdx = 0; btnIdx < buttons.length && btnIdx < 10; btnIdx++) {
          try {
            const button = buttons[btnIdx];
            const buttonText = await button.textContent().catch(() => '');
            const trimmedText = buttonText.trim();
            
            if (deleteKeywords.some(kw => trimmedText.includes(kw))) {
              if (await button.isVisible({ timeout: 500 })) {
                console.log(`    ✓ [D] 删除按钮: [${trimmedText}] (仅触发，不确认)`);
                screenInfo.controlActions.push({ type: 'DELETE', control: trimmedText });
                await button.click({ timeout: 3000 });
                await page.waitForTimeout(1500);
                
                // 检测确认对话框，点击取消
                const cancelKeywords = ['取消', 'キャンセル', 'Cancel', 'No'];
                const dialogButtons = await page.locator('button:visible').all();
                for (let cancelBtn of dialogButtons) {
                  const cancelBtnText = await cancelBtn.textContent().catch(() => '');
                  if (cancelKeywords.some(kw => cancelBtnText.includes(kw))) {
                    console.log(`      → 取消删除确认（避免删除真实数据）`);
                    await cancelBtn.click({ timeout: 2000 });
                    await page.waitForTimeout(1000);
                    break;
                  }
                }
              }
            }
          } catch (e) {
            console.warn(`    ⚠️  删除操作失败: ${e.message.slice(0, 50)}`);
          }
        }

        // === 阶段5：其他控件（下拉、标签页、分页等） ===
        console.log(`  \\n  🎛️  [OTHER] 操作其他交互控件...`);
        
        // 5.1 下拉选择
        const selects = await page.locator('select:visible').all();
        for (let select of selects.slice(0, 3)) {
          try {
            const options = await select.locator('option').count();
            if (options > 1) {
              console.log(`    ✓ 下拉框选择（${options}个选项）`);
              await select.selectOption({ index: 1 });
              await page.waitForTimeout(1000);
            }
          } catch (e) { /* 继续 */ }
        }

        // 5.2 标签页切换
        const tabs = await page.locator('[role="tab"], .tab, .nav-tabs a').all();
        for (let tab of tabs.slice(0, 3)) {
          try {
            const tabText = await tab.textContent().catch(() => '');
            console.log(`    ✓ 切换标签页: [${tabText.trim()}]`);
            await tab.click({ timeout: 2000 });
            await page.waitForTimeout(1000);
          } catch (e) { /* 继续 */ }
        }

        // 5.3 分页器
        const paginationNext = await page.locator('.pagination .next, .pagination button:has-text("下一页"), .pagination button:has-text("Next")').count();
        if (paginationNext > 0) {
          try {
            console.log(`    ✓ 点击下一页`);
            await page.locator('.pagination .next, .pagination button:has-text("下一页"), .pagination button:has-text("Next")').first().click({ timeout: 2000 });
            await page.waitForTimeout(1500);
          } catch (e) { /* 继续 */ }
        }

        screenInfo.capturedApis = apiRequests.length;
        result.screens.push(screenInfo);
        
        console.log(`  ✅ 完成，当前已捕获 ${apiRequests.length} 个 API`);

      } catch (err) {
        console.warn(`  ⚠️  访问失败: ${err.message}`);
      }
    }

    console.log('\\n========== EXPLORATION COMPLETE ==========');
    console.log(`📊 总计识别: ${result.screens.length} 个画面`);
    console.log(`📊 总计捕获: ${apiRequests.length} 个 API 请求`);

    // === 处理 API 数据 ===
    const apiMap = new Map();

    apiRequests.forEach(api => {
      // 过滤掉非业务 API
      const url = api.url;
      if (url.match(/\\.(js|css|png|jpg|jpeg|gif|svg|woff|woff2|ttf|ico)$/)) return;
      if (url.includes('/static/')) return;
      if (url.includes('/_next/')) return;
      
      // 提取路径
      const urlObj = new URL(url);
      const path = urlObj.pathname + urlObj.search;
      
      const key = `${api.method} ${path}`;
      
      if (!apiMap.has(key)) {
        apiMap.set(key, {
          method: api.method,
          path: path,
          fullUrl: url,
          status: api.status,
          requestHeaders: api.headers,
          requestBody: api.body,
          responseStatus: api.status,
          responseBody: api.responseBody
        });
      }
    });

    result.apis = Array.from(apiMap.values());

    // 保存结果
    const jsonOutput = JSON.stringify(result, null, 2);
    fs.writeFileSync('api-data.json', jsonOutput);

    console.log('\\n========== JSON OUTPUT START ==========');
    console.log(jsonOutput);
    console.log('========== JSON OUTPUT END ==========');

  } catch (error) {
    console.error('\\n❌ 探索失败:', error.message);
    console.error(error.stack);
    process.exit(1);
  } finally {
    await browser.close();
  }
})();
```

#### 2.2 执行探索脚本

**AI 主动执行命令：**

```powershell
# Step 1: 进入工作目录
cd api-explorer

# Step 2: 安装依赖
npm install

# Step 3: 设置环境变量并执行
$env:BASE_URL="{base_url}"; $env:USERNAME="{username}"; $env:PASSWORD="{password}"; node explore-api.js
```

**等待时间：** `initial_wait: 180` (API探索需要2-3分钟)

#### 2.3 解析探索结果

**AI 主动解析输出：**

1. 从控制台输出中提取 `JSON OUTPUT START` 到 `END` 之间的 JSON
2. 或直接读取 `api-data.json` 文件
3. 解析 JSON 结构：

```json
{
  "baseUrl": "https://192.168.11.104:8443",
  "screens": [
    {
      "name": "ダッシュボード",
      "url": "https://192.168.11.104:8443/dashboard",
      "capturedApis": 5
    }
  ],
  "apis": [
    {
      "method": "GET",
      "path": "/api/users",
      "status": 200,
      "requestHeaders": {...},
      "responseBody": "{\"users\": [...]}"
    },
    {
      "method": "POST",
      "path": "/api/users",
      "status": 201,
      "requestBody": "{\"name\":\"test\"}",
      "responseBody": "{\"id\": 123}"
    }
  ]
}
```

#### 2.4 AI 主动输出 API 清单

```
📋 自动探索完成！

🌐 访问了 {screens.length} 个画面
📊 捕获了 {apis.length} 个 API 端点

┌─────┬──────────────┬──────────────────────┬────────┬──────────┐
│ #   │ 方法         │ 路径                  │ 状态码 │ 来源画面 │
├─────┼──────────────┼──────────────────────┼────────┼──────────┤
│ 1   │ GET          │ /api/users            │ 200    │ 用户管理 │
│ 2   │ POST         │ /api/users            │ 201    │ 用户管理 │
│ 3   │ PUT          │ /api/users/{id}       │ 200    │ 用户管理 │
│ 4   │ DELETE       │ /api/users/{id}       │ 204    │ 用户管理 │
│ 5   │ GET          │ /api/projects         │ 200    │ 项目管理 │
└─────┴──────────────┴──────────────────────┴────────┴──────────┘

预计生成用例: 约 {estimated_cases} 条

---
请输入【继续】开始生成用例
```

---

### 第三步：AI 自动生成测试用例

**现在 AI 拥有完整的 API 数据，可以自动生成测试用例！**

#### 3.1 用例生成策略

**对于每个 API 端点：**

1. **正常场景**（200/201）
   - 使用捕获的真实请求参数
   - 验证响应状态码

2. **错误场景**（401/403/404）
   - 无Token → 401
   - 无效Token → 401
   - 不存在的ID → 404

3. **边界场景**（400）
   - 缺少必填参数
   - 参数类型错误

#### 3.2 自动生成 script_code

**AI 根据探索结果，自动编写 Playwright API 测试脚本：**

**示例：GET /api/users（正常场景）**

```javascript
async (page) => {
  // 登录获取 Token
  const loginRes = await page.request.post('${BASE_URL}/api/v1/auth/login', {
    data: { username: '${USERNAME}', password: '${PASSWORD}' },
    ignoreHTTPSErrors: true
  });
  const loginData = await loginRes.json();
  const token = loginData.data?.token || loginData.token;
  
  if (!token) {
    return { passed: false, error: 'Login failed' };
  }
  
  // 调用目标 API
  const res = await page.request.get('${BASE_URL}/api/users', {
    headers: { 'Authorization': 'Bearer ' + token },
    ignoreHTTPSErrors: true
  });
  
  return {
    passed: res.status() === 200,
    status: res.status(),
    message: `Status: ${res.status()}`
  };
}
```

**示例：POST /api/users（创建用户）**

```javascript
async (page) => {
  // 登录
  const loginRes = await page.request.post('${BASE_URL}/api/v1/auth/login', {
    data: { username: '${USERNAME}', password: '${PASSWORD}' },
    ignoreHTTPSErrors: true
  });
  const token = (await loginRes.json()).data.token;
  
  // 创建用户
  const testUser = { name: 'test_' + Date.now(), email: 'test@example.com' };
  const createRes = await page.request.post('${BASE_URL}/api/users', {
    data: testUser,
    headers: { 'Authorization': 'Bearer ' + token },
    ignoreHTTPSErrors: true
  });
  
  // 验证并删除（数据恢复原则）
  if (createRes.status() === 201) {
    const userId = (await createRes.json()).id;
    await page.request.delete(\`\${BASE_URL}/api/users/\${userId}\`, {
      headers: { 'Authorization': 'Bearer ' + token },
      ignoreHTTPSErrors: true
    });
  }
  
  return {
    passed: createRes.status() === 201,
    status: createRes.status()
  };
}
```

#### 3.3 批量生成用例

**AI 主动遍历所有 API：**

```
FOR 每个 API IN result.apis:
    识别 API 类型 (GET/POST/PUT/DELETE)
    
    生成用例数据:
        - screen = 来源画面名称
        - url = API路径（保留占位符）
        - method = HTTP方法
        - header = "{\"Authorization\": \"Bearer ${token}\"}"
        - body = 请求体（如有）
        - response = "{\"code\": {status}}"
        - script_code = 完整 Playwright API 脚本
    
    将用例加入批次
END FOR

// 批量写入
调用 create_api_cases(...)
```

---

### 第四步：验证与写入

**🚨 核心：逐条验证并写入（生成1→验证1→写入1）**

#### 4.1 单条用例处理流程

**严格遵循以下循环，直到所有用例处理完成：**

```
FOR 每个 API IN api-data.json:
    步骤A - 生成单条用例：
    → 根据API数据生成1条用例（包含script_code）
    
    步骤B - 立即验证：
    → 创建临时验证脚本 validate.js
    → 执行验证：node validate.js
    → 检查结果：IF result.passed === true THEN 继续，ELSE 跳过
    
    步骤C - 立即写入数据库（仅验证通过的）：
    → IF 是第1条用例 THEN
        调用 create_api_cases(project_id, group_name, [case], variables)
      ELSE
        调用 create_api_cases(project_id, group_name, [case])  // 无variables参数
    → 写入成功：已写入用例数 += 1
    → 写入失败：记录到失败列表
    
    步骤D - 实时输出：
    ✅ [N/total] 已写入: [画面名] 方法 路径 - 描述(状态码)
    
    步骤E - 判断是否继续：
    → IF 未达到token限制 THEN 继续下一条
    → ELSE 输出进度报告，等待用户输入【继续】
END FOR
```

#### 4.2 验证脚本模板

**AI 为每条用例创建独立验证脚本：**

```javascript
// validate.js
const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ ignoreHTTPSErrors: true });
  
  // 将 ${变量} 替换为实际值
  const BASE_URL = '实际URL';
  const USERNAME = '实际用户名';
  const PASSWORD = '实际密码';
  
  try {
    // 执行用例的 script_code
    const scriptCode = `用例的script_code内容`;
    const result = await eval(`(${scriptCode})(page)`);
    
    console.log(JSON.stringify({
      passed: result.passed === true && result.status === 200,  // 严格检查
      status: result.status,
      message: result.message || 'OK'
    }));
  } catch (error) {
    console.log(JSON.stringify({
      passed: false,
      status: 'ERROR',
      message: error.message
    }));
  } finally {
    await browser.close();
  }
})();
```

#### 4.3 实时进度输出

**每写入1条成功，立即输出：**

```
✅ [1/47] 已写入: [用户管理] GET /api/v1/users - 正常访问(200)
✅ [2/47] 已写入: [用户管理] GET /api/v1/users - 无Token(401)
❌ [3/47] 验证失败，跳过: [用户管理] POST /api/v1/users - 重复创建(409) - 原因：脚本执行超时
✅ [4/47] 已写入: [用户管理] POST /api/v1/users - 创建用户(201)
...
```

**进度频率：**
- 总数 ≤ 20：每条输出
- 总数 21-50：每5条输出汇总
- 总数 > 50：每10条输出汇总

---

### 第五步：进度检查与继续（Token管理）

> **⚠️ 每个画面的用例写入完成后，执行进度检查，决定继续或暂停！**

#### 5.1 自动继续（Token充足）

**AI 自动判断逻辑：**

```javascript
const remainingTokens = getCurrentTokens();
const estimatedNeededTokens = calculateNeeded(remainingCases);

if (remainingTokens > estimatedNeededTokens * 2) {
  // Token充足，自动继续
  continueProcessing();
} else {
  // Token不足，请求用户确认
  askUserToContinue();
}
```

**Token充足时的输出：**

```
📊 进度检查：
- 已处理用例：23/47 (48.9%)
- 预估Token剩余：> 40%

✅ Token充足，继续处理...
```

→ **AI 不询问用户，直接继续执行**

#### 5.2 等待用户确认（Token不足）

**触发条件（满足任一）：**
1. 预估剩余Token < 15%
2. 单次输出即将超过限制
3. 已处理用例数 > 总数的60%且仍有大量剩余

**必须输出以下提示并等待：**

```
⏸️ API用例生成进度报告

📊 本批次统计：
- 已写入用例：23 条
- 验证失败跳过：2 条
- 成功率：92.0%

📋 已完成画面：
- [用户管理]: 15条 ✓
- [提示词管理]: 8条 ✓

⏳ 待处理画面：
- [个人中心]: 预计12条
- [项目管理]: 预计10条

⚠️ Token剩余不足，请输入【继续】以完成剩余用例
   或输入【停止】以结束生成
```

→ **AI 暂停，等待用户输入**

#### 5.3 最终汇总报告

**所有用例处理完成后输出：**

```
🎉 API用例生成完成！

📊 最终统计：
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ 成功写入：    45 条
❌ 验证失败跳过： 2 条
📈 成功率：      95.7%
⏱️  总耗时：      约 5 分钟

📋 各画面明细：
┌──────────────┬─────────┬─────────┬──────────┐
│ 画面名称      │ 写入成功 │ 失败跳过 │ 成功率   │
├──────────────┼─────────┼─────────┼──────────┤
│ [用户管理]    │ 15条    │ 0条     │ 100.0%   │
│ [提示词管理]  │ 8条     │ 1条     │ 88.9%    │
│ [个人中心]    │ 12条    │ 1条     │ 92.3%    │
│ [项目管理]    │ 10条    │ 0条     │ 100.0%   │
└──────────────┴─────────┴─────────┴──────────┘

❌ 失败用例详情：
1. [提示词管理] POST /api/v1/prompts - 重复创建(409)
   原因：无法模拟409冲突场景
2. [个人中心] PUT /api/v1/profile/avatar - 上传头像(200)
   原因：文件上传脚本执行失败

✅ 所有用例已写入数据库，可以在系统中查看和执行！
```

---

## 6. 错误处理与异常场景

> **重要：S11-02的错误处理必须确保画面遍历和API捕获的完整性，任何错误都不应阻断整体流程。**

### 6.1 项目/用例集获取失败

**场景：** `get_current_project_name()` 或 `list_api_groups()` 失败

```
⚠️  项目信息获取失败

错误: {error_message}

处理方式：
1. 检查当前选择的项目
2. AI自动检查MCP工具连接状态并重试
3. 如果持续失败，输出错误并终止

→ 严重错误时自动终止任务
```

### 5.2 目标网站无法访问

**场景：** 浏览器无法打开目标网站

```
⚠️  目标网站访问异常

目标: {base_url}
错误: 连接超时 / 网络错误 / DNS解析失败

处理方式：
1. AI自动检查URL格式
2. AI自动重试3次（每次间隔5秒）
3. 尝试使用备用协议（http/https切换）
4. 如果是HTTPS证书问题，AI自动添加 ignoreHTTPSErrors 参数

→ 连续失败3次后自动终止任务
```

### 5.3 登录失败

**场景：** 无法成功登录目标网站

```
⚠️  登录失败

用户名: {username}
错误: 登录按钮未找到 / 凭证错误 / 验证码拦截

处理方式：
1. AI自动尝试多种常见登录表单选择器
2. AI自动重试3次
3. 如果是验证码问题，输出错误并终止（需手动处理）
4. 检查元数据中的凭证是否正确

→ 自动尝试失败后终止任务，输出详细错误信息
```

### 5.4 画面识别不完整

**场景：** 自动脚本未能识别所有画面

```
⚠️  画面识别可能不完整

识别到: {screen_count}个画面
风险: 可能存在下拉菜单、标签页、权限限制的隐藏画面

处理方式：
1. AI自动尝试展开下拉菜单
2. AI自动切换标签页
3. AI自动遍历所有可见菜单项
4. 继续处理已识别的画面

→ AI自动继续处理，不等待用户确认
```

### 5.5 控件遍历失败

**场景：** 某些控件无法正常操作

```
⚠️  控件操作异常

画面: {screen_name}
控件: {widget_name} ({widget_type})
错误: 元素不可见 / 点击无响应 / 操作超时

处理方式：
1. AI自动滚动到元素位置
2. AI自动等待元素可见
3. AI自动重试3次（每次间隔2秒）
4. 持续失败则标记为"跳过"，记录原因

→ 自动跳过失败的控件，继续处理下一个
```

### 5.6 API捕获为空

**场景：** 实时拦截未能捕获任何API请求

```
⚠️  API捕获异常

预期: 应捕获到操作触发的API请求
实际: 拦截器未捕获任何xhr/fetch请求

可能原因：
- 页面未发起任何API请求（纯静态页面）
- API请求类型不是xhr/fetch（如WebSocket）
- 拦截器配置错误

处理方式：
1. AI自动检查拦截器配置
2. AI自动增加等待时间并重试
3. 如确认该画面无API，标记为"无API"

→ 自动继续处理下一个画面
```

### 5.7 脚本验证失败（带重试机制）

**场景：** `node validate.js` 执行失败或返回非预期状态码

**🚨 关键原则：验证失败的用例禁止写入数据库！**

**处理流程（指数退避）：**

```python
FOR retry_count IN [1, 2, 3]:
    TRY:
        result = 运行 node validate.js
        
        IF result.passed == true && result.status == expected_status:
            标记为"验证通过" ✅
            → 可以写入数据库
            BREAK
        ELSE:
            抛出 SCRIPT_ERROR(result.error)
        END IF
        
    CATCH SCRIPT_ERROR AS error:
        IF retry_count < 3:
            等待时间 = 1 * (2 ^ (retry_count - 1))  # 1s, 2s, 4s
            输出: "⚠️  脚本验证失败，{等待时间}秒后重试（第{retry_count}/3次）..."
            
            # 尝试修正脚本
            IF error包含"Token获取失败":
                检查登录逻辑是否正确
            ELSE IF error包含"资源不存在":
                检查URL中的ID是否为真实捕获的ID
            ELSE IF error包含"参数错误":
                检查请求体字段是否完整
            END IF
            
            等待(等待时间)
        ELSE:
            输出: "❌ 脚本验证失败3次，跳过此用例（不写入数据库）"
            标记为"验证失败" ❌
            记录错误详情到最终报告: {
                api: "{method} {url}",
                scenario: "{场景}",
                error: "{error_message}",
                retries: 3,
                action: "跳过，未写入数据库"
            }
            # 🚨 继续处理下一个场景（不阻断整体流程，不写入此条失败用例）
        END IF
    END TRY
END FOR
```

**常见验证失败原因与修正方法：**

| 失败原因 | 修正方法 |
|---------|---------| 
| Token获取失败 | 检查登录API路径和响应结构 |
| 资源不存在(404) | 从捕获的API中提取真实ID |
| 参数错误(400) | 补充必填字段或修正字段格式 |
| 权限不足(403) | 确认当前用户权限，或跳过该场景 |
| 网络超时 | 增加等待时间，检查网络状态 |

### 5.8 用例回写失败

**场景：** `create_api_cases` 调用失败

**处理流程（指数退避重试）：**

```python
FOR retry_count IN [1, 2, 3]:
    TRY:
        result = 调用 create_api_cases(
            project_id=project_id,
            group_name=group_name,
            cases=[case_data],
            variables=variables  # 仅第一条用例携带
        )
        
        IF result.success:
            输出: "✅ [{current}/{total}] 已写入: {screen} {method} {url} - {scenario}"
            BREAK
        ELSE:
            抛出 API_ERROR(result.error)
        END IF
        
    CATCH API_ERROR AS error:
        IF retry_count < 3:
            等待时间 = 1 * (2 ^ (retry_count - 1))  # 1s, 2s, 4s
            输出: "⚠️  用例回写失败，{等待时间}秒后重试（第{retry_count}/3次）..."
            等待(等待时间)
        ELSE:
            输出: "❌ 用例回写失败3次，跳过该用例"
            记录失败详情: {
                case_number: "{case_number}",
                api: "{method} {url}",
                error: "{error_message}"
            }
            # 继续处理下一个用例
        END IF
    END TRY
END FOR
```

### 5.9 Token超限处理（CRITICAL）

**场景：** 单次输出即将达到Token限制

```
⚠️  检测到输出即将达到Token限制

当前状态：
- 已处理画面: {completed_screens}个
- 已生成用例: {generated_cases}条
- 剩余画面: {remaining_screens}个

必须执行的操作：
1. 立即停止继续生成新用例
2. 确保当前画面的API用例全部写入完成
3. 输出详细的进度报告（阶段3格式）
4. 提示用户输入【继续】以继续

⚠️ Token不足时暂停，等待用户输入【继续】后恢复
```

### 5.10 Playwright CLI 命令执行失败

**场景：** Playwright CLI 未安装或命令执行失败

```
⚠️  Playwright CLI 执行失败

命令: {command}
错误: npx playwright not found / chromium not installed

处理方式：
1. AI自动检查 Playwright 版本: npx playwright --version
2. AI自动安装浏览器: npx playwright install chromium
3. AI自动重试命令

→ 自动尝试修复，失败则终止
```
4. 确认npx命令可用

→ 重新执行安装步骤，确保CLI可用后继续
```

### 5.11 数据恢复失败

**场景：** 测试数据未能正确清理

```
⚠️  数据恢复异常

操作: {operation}（如：DELETE测试数据）
错误: {error_message}

影响: 可能残留测试数据

处理方式：
1. 记录未清理的数据ID
2. 在完成报告中列出残留数据清单
3. 提供手动清理的SQL/API命令
4. 继续执行（不阻断流程）

残留数据记录：
- 资源ID: {test_id}
- 类型: {resource_type}
- 清理建议: DELETE /api/{resource}/{test_id}
```

### 5.12 错误处理总原则

1. **画面遍历优先**：任何单个控件/API的错误不应阻断画面遍历
2. **记录继续**：失败的用例记录原因后继续处理下一个
3. **🚨 验证失败不写入**：只有 passed === true 且状态码匹配的用例才能写入数据库
4. **指数退避**：重试使用1s→2s→4s的延迟
5. **完整报告**：所有失败在最终报告中详细列出（包括未写入的用例）
6. **自动化优先**：AI自动尝试解决问题，严重错误才终止任务
7. **Token限制**：只在Token不足时暂停，等待用户输入【继续】

## 7. 工具速查

### 7.1 AIGO 测试管理工具

| 工具                                                          | 用途                       |
| ----------------------------------------------------------- | ------------------------ |
| `get_current_project_name()`                                | 1.1 获取当前项目               |
| `list_api_groups(project_id)`                               | 1.2 获取API用例集列表           |
| `get_api_group_metadata(group_name)`                        | 1.3 获取用例集元数据（用名称查询）      |
| `create_api_cases(project_id, group_name, cases, variables)` | 创建用例+写入变量（variables自动检重） |

### 6.2 变量表管理说明

**`create_api_cases` 的 `variables` 参数：**

```javascript
variables: [
  { var_key: 'base_url', var_value: 'https://example.com', var_desc: '基础URL' },
  { var_key: 'token', var_value: 'xxx', var_desc: '认证Token' }
]
```

**检重规则**：

| 情况    | 处理方式          | 示例                                                |
| ----- | ------------- | ------------------------------------------------- |
| 同名同值  | **跳过**，不重复创建  | 已有 `token=abc`，再写入 `token=abc` → 跳过               |
| 同名不同值 | **新建**带序号的变量名 | 已有 `token=abc`，再写入 `token=xyz` → 创建 `token_2=xyz` |
| 新变量   | 直接创建          | 写入 `user_id=123` → 创建                             |

> 🚨 **注意**：元数据变量（`base_url`、`username`、`password`）除外，这些变量会直接覆盖更新

**变量命名示例**：

```
第1次写入 token=abc     → token=abc
第2次写入 token=xyz     → token_2=xyz  (值不同，新建)
第3次写入 token=abc     → 跳过 (与token值相同)
第4次写入 token=123     → token_3=123  (值不同，继续新建)
```

### 6.3 Playwright CLI 命令速查

| 命令                                                                  | 用途          |
| ------------------------------------------------------------------- | ----------- |
| `npx playwright --version`                                          | 验证环境        |
| `npx playwright install chromium`                                   | 安装浏览器       |

> ⚠️ **重要**：所有命令必须包含 `--ignore-https-errors` 参数（自签名证书场景）

## 8. 用例场景模板

### 8.1 成功响应码

| 场景    | 方法     | 响应码 | 说明         |
| ----- | ------ | --- | ---------- |
| 正常查询  | GET    | 200 | OK         |
| 正常创建  | POST   | 201 | Created    |
| 无返回内容 | DELETE | 204 | No Content |
| 正常更新  | PUT    | 200 | OK         |
| 正常删除  | DELETE | 200 | OK         |

### 7.2 客户端错误码 (4xx)

| 场景          | 方法   | 响应码 | 说明                   |
| ----------- | ---- | --- | -------------------- |
| 参数缺失/格式错误   | POST | 400 | Bad Request          |
| 未登录/Token无效 | GET  | 401 | Unauthorized         |
| 无权限访问       | GET  | 403 | Forbidden            |
| 资源不存在       | GET  | 404 | Not Found            |
| 方法不允许       | POST | 405 | Method Not Allowed   |
| 资源冲突(如重复创建) | POST | 409 | Conflict             |
| 数据验证失败      | POST | 422 | Unprocessable Entity |
| 请求过于频繁      | GET  | 429 | Too Many Requests    |

### 7.3 服务端错误码 (5xx)

| 场景      | 方法  | 响应码 | 说明                    |
| ------- | --- | --- | --------------------- |
| 服务器内部错误 | ANY | 500 | Internal Server Error |
| 网关错误    | ANY | 502 | Bad Gateway           |
| 服务暂不可用  | ANY | 503 | Service Unavailable   |
| 网关超时    | ANY | 504 | Gateway Timeout       |

---

## 开始生成

生成API接口测试用例，目标用例集：**{{group_name}}**
