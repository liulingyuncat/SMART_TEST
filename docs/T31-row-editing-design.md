# T31 行内编辑功能设计文档

## 1. 设计概述

### 1.1 设计目标

实现高性能、用户友好的测试用例行内编辑功能，支持多种用例类型和多语言字段编辑。

### 1.2 架构设计

```
┌─────────────────────────────────────────────────────┐
│                  EditableTable Component             │
├─────────────────────────────────────────────────────┤
│  State Management:                                   │
│  - editingKey: string (当前编辑行的case_id)          │
│  - cases: Array (用例数据)                           │
│  - hasEditChanges: boolean (是否有未保存的插入)     │
│  - form: Form instance (表单实例)                    │
├─────────────────────────────────────────────────────┤
│  Core Functions:                                     │
│  - startEdit(record) → 进入编辑状态                  │
│  - saveEdit(record) → 保存编辑                       │
│  - cancelEdit() → 取消编辑                           │
│  - createDefaultEmptyRow() → 创建默认空行           │
├─────────────────────────────────────────────────────┤
│  Components:                                         │
│  - EditableCell → 可编辑单元格                       │
│  - MultiLangEditModal → 多语言编辑对话框            │
└─────────────────────────────────────────────────────┘
```

## 2. 状态管理设计

### 2.1 核心状态

```javascript
// 编辑状态
const [editingKey, setEditingKey] = useState('');  // 当前编辑行的case_id

// 表单实例
const [form] = Form.useForm();

// 用例数据
const [cases, setCases] = useState([]);

// 插入操作标记
const [hasEditChanges, setHasEditChanges] = useState(false);
```

### 2.2 状态转换图

```
┌─────────────┐
│   初始状态   │ (editingKey = '')
└──────┬──────┘
       │ 点击Edit
       ↓
┌─────────────┐
│   编辑状态   │ (editingKey = record.case_id)
└──────┬──────┘
       │ 点击Save/Cancel
       ↓
┌─────────────┐
│   初始状态   │ (editingKey = '')
└─────────────┘
```

### 2.3 状态优先级

1. **编辑状态 (editingKey)**：最高优先级
   - 控制当前行的编辑模式
   - 影响其他操作按钮的可用性

2. **插入状态 (hasEditChanges)**：高优先级
   - 阻止编辑和删除操作
   - 控制分页显示

3. **数据加载 (loading)**：中优先级
   - 控制加载指示器显示
   - 阻止用户操作

## 3. 核心功能设计

### 3.1 startEdit - 进入编辑状态

#### 3.1.1 函数签名

```javascript
const startEdit = useCallback((record) => {
  // 实现逻辑
}, [form, caseType, language]);
```

#### 3.1.2 处理流程

```
1. 数据规范化
   ├─ getDisplayValue(value)
   │  ├─ null/undefined → ''
   │  └─ 其他 → String(value)
   │
2. 根据用例类型构造表单值
   ├─ AI用例 (caseType='ai')
   │  └─ 单语言字段
   │
   ├─ Role用例 (caseType.startsWith('role'))
   │  └─ 多语言字段 + 语言后缀
   │
   └─ 整体/变更/受入用例
      └─ 多语言字段 + 语言后缀
      
3. 设置表单值
   └─ form.setFieldsValue(formValues)
   
4. 设置编辑键
   └─ setEditingKey(record.case_id)
```

#### 3.1.3 字段映射规则

**AI用例字段映射**：
```javascript
{
  case_number: record.case_number,
  major_function: record.major_function,
  middle_function: record.middle_function,
  minor_function: record.minor_function,
  precondition: record.precondition,
  test_steps: record.test_steps,
  expected_result: record.expected_result,
  remark: record.remark
}
```

**Role用例字段映射**：
```javascript
const langSuffix = language === '中文' ? '_cn' : language === 'English' ? '_en' : '_jp';
{
  case_number: record.case_number,
  [`screen${langSuffix}`]: record[`screen${langSuffix}`],
  [`function${langSuffix}`]: record[`function${langSuffix}`],
  [`precondition${langSuffix}`]: record[`precondition${langSuffix}`],
  [`test_steps${langSuffix}`]: record[`test_steps${langSuffix}`],
  [`expected_result${langSuffix}`]: record[`expected_result${langSuffix}`],
  test_result: record.test_result,
  remark: record.remark
}
```

**整体/变更/受入用例字段映射**：
```javascript
const langSuffix = language === '中文' ? '_cn' : language === 'English' ? '_en' : '_jp';
{
  case_number: record.case_number,
  [`major_function${langSuffix}`]: record[`major_function${langSuffix}`],
  [`middle_function${langSuffix}`]: record[`middle_function${langSuffix}`],
  [`minor_function${langSuffix}`]: record[`minor_function${langSuffix}`],
  [`precondition${langSuffix}`]: record[`precondition${langSuffix}`],
  [`test_steps${langSuffix}`]: record[`test_steps${langSuffix}`],
  [`expected_result${langSuffix}`]: record[`expected_result${langSuffix}`],
  test_result: record.test_result,
  remark: record.remark
}
```

### 3.2 saveEdit - 保存编辑

#### 3.2.1 函数签名

```javascript
const saveEdit = useCallback(async (recordParam) => {
  // 实现逻辑
}, [form, projectId, pagination, fetchCases, caseType, language, cases]);
```

#### 3.2.2 处理流程

```
1. 表单验证
   └─ form.validateFields()
      ├─ 成功 → 继续
      └─ 失败 → 抛出异常

2. 获取最新记录
   └─ cases.find(c => c.case_id === recordParam.case_id)
      ├─ 找到 → 继续
      └─ 未找到 → 错误提示

3. 收集变更字段
   ├─ normalizeValue(value) → trim后的字符串
   ├─ 对比 formValue !== recordValue
   └─ 收集到 updates 对象

4. 发送更新请求
   ├─ 有变更 → 调用updateCaseAPI
   └─ 无变更 → 跳过

5. 更新本地数据
   └─ setCases(prevCases => prevCases.map(...))

6. 退出编辑状态
   └─ setEditingKey('')

7. 显示提示
   ├─ 成功 → message.success()
   └─ 失败 → message.error()
```

#### 3.2.3 值规范化逻辑

```javascript
const normalizeValue = (value) => {
  if (value === null || value === undefined) return '';
  return String(value).trim();
};

// 构造更新值
const updateValue = row[key] === undefined || row[key] === null 
  ? '' 
  : String(row[key]);
```

#### 3.2.4 差异检测算法

```javascript
Object.keys(row).forEach(key => {
  const formValue = normalizeValue(row[key]);
  const recordValue = normalizeValue(currentRecord[key]);
  
  if (formValue !== recordValue) {
    updates[key] = row[key] === undefined || row[key] === null 
      ? '' 
      : String(row[key]);
  }
});
```

#### 3.2.5 本地数据更新

```javascript
setCases(prevCases => 
  prevCases.map(c => 
    c.case_id === currentRecord.case_id 
      ? { ...c, ...updates }  // 合并更新
      : c                      // 保持不变
  )
);
```

**优势**：
- 不触发fetchCases，避免整页刷新
- 保留通过Above/Below插入的新行
- 保持hasEditChanges状态
- 更新UI立即响应

### 3.3 cancelEdit - 取消编辑

#### 3.3.1 函数签名

```javascript
const cancelEdit = useCallback(() => {
  setEditingKey('');
  form.resetFields();
}, [form]);
```

#### 3.3.2 处理流程

```
1. 清空编辑键
   └─ setEditingKey('')

2. 重置表单
   └─ form.resetFields()
      └─ 恢复初始值
```

### 3.4 createDefaultEmptyRow - 创建默认空行

#### 3.4.1 函数签名

```javascript
const createDefaultEmptyRow = useCallback(async () => {
  // 实现逻辑
}, [projectId, caseType, language]);
```

#### 3.4.2 处理流程

```
1. 判断用例类型
   └─ isRoleType = caseType.startsWith('role')

2. 选择创建API
   ├─ Role类型 → createAutoCase
   └─ 其他类型 → createCase

3. 构造请求参数
   └─ {
        case_type: caseType,
        language: !isRoleType ? language : undefined
      }

4. 调用API创建
   └─ await createAPI(projectId, createData)

5. 返回新记录
   └─ 包含真实的case_id (UUID)

6. 错误处理
   └─ 捕获异常，记录日志，重新抛出
```

#### 3.4.3 参数构造

**手工用例（ai/overall/change/acceptance）**：
```javascript
{
  case_type: caseType,
  language: language  // 必填：'中文' | 'English' | '日本語'
}
```

**Role类型用例（role1/role2/role3/role4）**：
```javascript
{
  case_type: caseType  // 不需要language参数
}
```

#### 3.4.4 调用场景

1. **fetchCases检测到空数据**：
```javascript
if (casesData.length === 0) {
  try {
    const newCase = await createDefaultEmptyRow();
    casesData = [newCase];
  } catch (error) {
    casesData = [];
  }
}
```

2. **清空AI用例后**：
```javascript
await clearAICases(projectId);
try {
  const newCase = await createDefaultEmptyRow();
  setCases([newCase]);
  setPagination(prev => ({ ...prev, current: 1, total: 1 }));
} catch (error) {
  setCases([]);
  setPagination(prev => ({ ...prev, current: 1, total: 0 }));
}
```

3. **删除到空表格**：
```javascript
if (cases.length === 1) {  // 删除后为空
  try {
    const newCase = await createDefaultEmptyRow();
    setCases([newCase]);
    setPagination(prev => ({ ...prev, current: 1, total: 1 }));
  } catch (error) {
    setCases([]);
    setPagination(prev => ({ ...prev, current: 1, total: 0 }));
  }
}
```

4. **批量删除到空表格**：
```javascript
if (cases.length === selectedRowKeys.length) {  // 全部删除
  try {
    const newCase = await createDefaultEmptyRow();
    setCases([newCase]);
    setPagination(prev => ({ ...prev, current: 1, total: 1 }));
  } catch (error) {
    setCases([]);
    setPagination(prev => ({ ...prev, current: 1, total: 0 }));
  }
}
```

## 4. 组件设计

### 4.1 EditableCell - 可编辑单元格

#### 4.1.1 组件职责

- 根据编辑状态显示输入框或普通文本
- 管理单元格的Form.Item
- 支持不同输入类型（text/textarea/select）

#### 4.1.2 组件结构

```jsx
const EditableCell = ({
  editing,          // 是否处于编辑状态
  dataIndex,        // 字段名
  title,            // 列标题
  inputType,        // 输入类型：'text' | 'textarea' | 'select'
  record,           // 当前行数据
  children,         // 子元素（非编辑时显示）
  ...restProps
}) => {
  // 编辑状态判断
  const isEditable = editing && !hasEditChanges;
  
  return (
    <td {...restProps}>
      {isEditable ? (
        <Form.Item name={dataIndex} style={{ margin: 0 }}>
          {inputType === 'textarea' ? (
            <TextArea rows={2} />
          ) : (
            <Input />
          )}
        </Form.Item>
      ) : (
        children
      )}
    </td>
  );
};
```

#### 4.1.3 编辑状态控制

**关键逻辑**：
```javascript
editing && !hasEditChanges
```

**原因**：
- `editing`：当前单元格是否在编辑状态
- `!hasEditChanges`：没有未保存的插入操作
- 两个条件都满足才显示输入框

### 4.2 MultiLangEditModal - 多语言编辑对话框

#### 4.2.1 触发条件

- 仅在中文界面（language='中文'）
- 点击蓝色链接文本
- 当前行不在编辑状态

#### 4.2.2 编辑字段

同时编辑三种语言：
- {field}_cn：中文
- {field}_en：英文
- {field}_jp：日文

#### 4.2.3 保存逻辑

```javascript
const handleSave = async () => {
  const updates = {
    [`${fieldName}_cn`]: values.cn,
    [`${fieldName}_en`]: values.en,
    [`${fieldName}_jp`]: values.jp
  };
  
  await updateCaseAPI(record.case_id, updates);
  
  // 更新本地数据
  setCases(prevCases =>
    prevCases.map(c =>
      c.case_id === record.case_id
        ? { ...c, ...updates }
        : c
    )
  );
};
```

## 5. 特殊字段处理

### 5.1 TestResult字段即时保存

#### 5.1.1 实现方式

```jsx
<Select 
  defaultValue={record.test_result || 'NR'}
  onChange={(value) => {
    if (value !== record.test_result) {
      updateCaseAPI(record.case_id, { test_result: value })
        .then(() => {
          message.success('保存成功');
          setCases(prevCases =>
            prevCases.map(c =>
              c.case_id === record.case_id
                ? { ...c, test_result: value }
                : c
            )
          );
        })
        .catch(error => {
          message.error('保存失败');
        });
    }
  }}
>
  <Select.Option value="NR">NR</Select.Option>
  <Select.Option value="OK">OK</Select.Option>
  <Select.Option value="NG">NG</Select.Option>
  <Select.Option value="Block">Block</Select.Option>
</Select>
```

#### 5.1.2 关键点

- 使用defaultValue而非value（避免受控问题）
- onChange时立即调用API
- 只在值改变时保存
- 成功后更新本地数据
- 不调用fetchCases（保持编辑状态）

### 5.2 语言字段映射

#### 5.2.1 后缀计算

```javascript
const langFieldSuffix = language === '中文' 
  ? '_cn' 
  : language === 'English' 
  ? '_en' 
  : '_jp';
```

#### 5.2.2 动态字段名

```javascript
const fieldName = `major_function${langFieldSuffix}`;
// language='中文' → 'major_function_cn'
// language='English' → 'major_function_en'
// language='日本語' → 'major_function_jp'
```

## 6. 按钮状态控制

### 6.1 Edit按钮

```jsx
<Button
  disabled={editingKey !== '' || hasEditChanges}
  onClick={() => startEdit(record)}
>
  Edit
</Button>
```

**禁用条件**：
- `editingKey !== ''`：有行在编辑
- `hasEditChanges`：有未保存的插入操作

### 6.2 Delete按钮

```jsx
<Button
  disabled={editingKey === record.case_id || hasEditChanges}
  onClick={() => handleDelete(record)}
>
  Delete
</Button>
```

**禁用条件**：
- `editingKey === record.case_id`：当前行在编辑
- `hasEditChanges`：有未保存的插入操作

### 6.3 Save/Cancel按钮

```jsx
{editingKey === record.case_id ? (
  <>
    <Button onClick={() => saveEdit(record)}>Save</Button>
    <Button onClick={cancelEdit}>Cancel</Button>
  </>
) : (
  <Button onClick={() => startEdit(record)}>Edit</Button>
)}
```

**显示条件**：
- 仅当前行在编辑时显示Save/Cancel
- 否则显示Edit按钮

### 6.4 Above/Below按钮

```jsx
<Button
  disabled={editingKey !== '' || loading}
  onClick={() => handleInsertAbove(record.case_id)}
>
  Above
</Button>
```

**行为**：
- 点击时强制取消当前编辑
- 不受hasEditChanges限制（可多次插入）

## 7. 性能优化

### 7.1 useCallback优化

所有事件处理函数使用useCallback：
```javascript
const startEdit = useCallback((record) => {
  // ...
}, [form, caseType, language]);

const saveEdit = useCallback(async (recordParam) => {
  // ...
}, [form, projectId, pagination, fetchCases, caseType, language, cases]);

const cancelEdit = useCallback(() => {
  // ...
}, [form]);
```

**优势**：
- 避免不必要的重新渲染
- 稳定的函数引用
- 减少子组件更新

### 7.2 本地数据更新

使用函数式setState：
```javascript
setCases(prevCases => 
  prevCases.map(c => 
    c.case_id === currentRecord.case_id 
      ? { ...c, ...updates }
      : c
  )
);
```

**优势**：
- 不依赖外部cases状态
- 避免闭包陷阱
- 确保使用最新状态

### 7.3 避免整页刷新

保存成功后不调用fetchCases：
```javascript
// ❌ 旧方式 - 触发整页刷新
await fetchCases(pagination.current);

// ✅ 新方式 - 本地更新
setCases(prevCases => prevCases.map(...));
```

**优势**：
- 保留Above/Below插入的新行
- 保持hasEditChanges状态
- 更快的响应速度
- 更好的用户体验

## 8. 错误处理

### 8.1 表单验证失败

```javascript
try {
  const row = await form.validateFields();
  // 继续处理
} catch (error) {
  // 验证失败，自动显示错误提示
  // 不需要额外处理
}
```

### 8.2 API调用失败

```javascript
try {
  await updateCaseAPI(currentRecord.case_id, updates);
  message.success('保存成功');
} catch (error) {
  console.error('Failed to save edit:', error);
  message.error('保存失败');
} finally {
  setEditingKey('');  // 无论成功失败都退出编辑
}
```

### 8.3 创建默认空行失败

```javascript
try {
  const newCase = await createDefaultEmptyRow();
  setCases([newCase]);
  setPagination(prev => ({ ...prev, current: 1, total: 1 }));
} catch (error) {
  console.error('[createDefaultEmptyRow] Failed:', error);
  // 降级处理：显示空表格
  setCases([]);
  setPagination(prev => ({ ...prev, current: 1, total: 0 }));
}
```

## 9. API接口设计

### 9.1 更新用例

**手工用例**：
```
PATCH /api/v1/projects/:projectId/manual-cases/:caseId
Content-Type: application/json

{
  "case_number": "TC001",
  "major_function_cn": "用户管理",
  "test_steps_cn": "1. 点击登录\n2. 输入用户名密码",
  ...
}
```

**自动化用例**：
```
PATCH /api/v1/projects/:projectId/auto-cases/:caseId
Content-Type: application/json

{
  "case_num": "AUTO001",
  "screen_cn": "登录画面",
  ...
}
```

### 9.2 创建空用例

**手工用例**：
```
POST /api/v1/projects/:projectId/manual-cases
Content-Type: application/json

{
  "case_type": "overall",
  "language": "中文"
}
```

**自动化用例**：
```
POST /api/v1/projects/:projectId/auto-cases
Content-Type: application/json

{
  "case_type": "role1"
}
```

**返回**：
```json
{
  "case_id": "uuid-string",
  "id": 1,
  "display_id": 1,
  "project_id": 1,
  "case_type": "overall",
  ...所有字段初始化为空值
}
```

## 10. 测试策略

### 10.1 单元测试

- startEdit函数测试
- saveEdit函数测试
- cancelEdit函数测试
- createDefaultEmptyRow函数测试
- 值规范化函数测试
- 字段映射逻辑测试

### 10.2 集成测试

- 完整的编辑-保存流程
- 完整的编辑-取消流程
- TestResult即时保存
- 多语言编辑对话框
- 默认空行创建和编辑

### 10.3 E2E测试

- 用户完整操作流程
- 跨语言切换测试
- 并发操作测试
- 错误场景测试

## 11. 未来改进

### 11.1 性能优化

- 虚拟滚动支持大量数据
- 批量更新优化
- 防抖/节流优化

### 11.2 功能增强

- 支持拖拽排序
- 支持复制粘贴
- 支持快捷键操作
- 支持批量编辑

### 11.3 用户体验

- 自动保存草稿
- 撤销/重做功能
- 更丰富的验证提示
- 字段联动提示
