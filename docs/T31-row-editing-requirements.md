# T31 行内编辑功能需求文档

## 1. 功能概述

Web智能测试平台的测试用例管理模块需要支持行内编辑功能，允许用户直接在表格中编辑测试用例的各个字段，提供流畅的用户体验。

## 2. 用户角色

- 测试工程师：创建、编辑、查看测试用例
- 项目成员：查看和编辑所属项目的测试用例

## 3. 功能需求

### 3.1 基本行编辑功能

#### 3.1.1 进入编辑状态

**需求描述**：
- 用户点击用例行的"Edit"按钮，该行进入编辑状态
- 编辑状态下，可编辑字段显示为输入框或文本域
- 一次只能编辑一行

**触发条件**：
- 点击用例行的"Edit"按钮
- 当前没有其他行在编辑状态
- 没有未保存的插入操作（hasEditChanges=false）

**前置条件**：
- 用户已登录
- 用户是项目成员
- 没有其他编辑操作进行中

#### 3.1.2 字段编辑

**需求描述**：
根据用例类型和当前语言，显示对应的可编辑字段

**AI用例（case_type='ai'）**：
- case_number：用例编号（文本输入）
- major_function：大功能分类（文本输入）
- middle_function：中功能分类（文本输入）
- minor_function：小功能分类（文本输入）
- precondition：前置条件（文本域）
- test_steps：测试步骤（文本域）
- expected_result：期望结果（文本域）
- remark：备注（文本输入）

**整体用例/变更用例/受入用例（case_type='overall'/'change'/'acceptance'）**：
- case_number：用例编号（文本输入）
- major_function_{cn/en/jp}：大功能分类（根据当前语言）
- middle_function_{cn/en/jp}：中功能分类（根据当前语言）
- minor_function_{cn/en/jp}：小功能分类（根据当前语言）
- precondition_{cn/en/jp}：前置条件（根据当前语言）
- test_steps_{cn/en/jp}：测试步骤（根据当前语言）
- expected_result_{cn/en/jp}：期望结果（根据当前语言）
- test_result：测试结果（下拉选择：NR/OK/NG/Block）
- remark：备注（文本输入）

**Role类型自动化用例（case_type='role1'/'role2'/'role3'/'role4'）**：
- case_num：用例编号（文本输入）
- screen_{cn/en/jp}：画面（根据当前语言）
- function_{cn/en/jp}：功能（根据当前语言）
- precondition_{cn/en/jp}：前置条件（根据当前语言）
- test_steps_{cn/en/jp}：测试步骤（根据当前语言）
- expected_result_{cn/en/jp}：期望结果（根据当前语言）
- test_result：测试结果（下拉选择：NR/OK/NG）
- remark：备注（文本输入）

**字段验证规则**：
- 所有字段允许为空
- 文本输入字段最大长度：100字符
- 文本域字段无长度限制
- test_result必须从预定义选项中选择

#### 3.1.3 保存编辑

**需求描述**：
- 用户点击"Save"按钮保存编辑内容
- 只更新有变化的字段
- 保存成功后退出编辑状态
- 更新本地显示数据，不刷新整页

**保存逻辑**：
1. 验证表单字段
2. 对比表单值与原始值，收集变更字段
3. 如果有变更，调用更新API
4. 更新成功后，更新本地cases数组
5. 退出编辑状态（setEditingKey('')）
6. 显示成功提示

**空值处理**：
- null/undefined值转换为空字符串''
- 表单提交时保持空字符串
- 对比时统一trim()去除首尾空格

**失败处理**：
- 保存失败后仍退出编辑状态
- 显示错误提示
- 不更新本地数据

#### 3.1.4 取消编辑

**需求描述**：
- 用户点击"Cancel"按钮取消编辑
- 取消后恢复原始数据
- 退出编辑状态

**取消逻辑**：
1. 清空editingKey状态
2. 重置表单字段（form.resetFields()）

### 3.2 特殊字段处理

#### 3.2.1 TestResult字段即时保存

**需求描述**：
- TestResult字段在编辑状态下显示为下拉选择框
- 用户选择新值后，立即保存到后端
- 不需要点击Save按钮
- 保存成功后更新本地数据，保持编辑状态

**实现要点**：
- 使用Select组件的onChange事件
- 直接调用updateCaseAPI
- 更新本地cases数组中的对应记录
- 不调用fetchCases（避免触发整页刷新）
- 保持hasEditChanges状态不变

#### 3.2.2 多语言字段处理

**整体/变更/受入用例的多语言编辑**：

**中文界面（language='中文'）**：
- 显示蓝色链接文本，点击打开多语言编辑对话框
- 对话框中可同时编辑CN/EN/JP三种语言
- 保存后更新所有三种语言的字段

**英文/日文界面（language='English'/'日本語'）**：
- 字段进入编辑状态时显示对应语言的输入框
- 只编辑当前语言的字段
- 保存时只更新当前语言字段

### 3.3 编辑状态限制

#### 3.3.1 单行编辑限制

**需求描述**：
- 同一时间只能编辑一行
- 编辑过程中，其他行的Edit按钮禁用

**实现方式**：
- 使用editingKey状态记录当前编辑行的case_id
- Edit按钮的disabled属性：`disabled={editingKey !== '' && editingKey !== record.case_id}`

#### 3.3.2 与插入操作的互斥

**需求描述**：
- 有未保存的插入操作时，禁止编辑
- 编辑过程中，禁止Above/Below插入操作（通过强制取消编辑实现）

**实现方式**：
- Edit按钮disabled条件：`disabled={editingKey !== '' || hasEditChanges}`
- Above/Below按钮点击时强制取消编辑（setEditingKey(''), form.resetFields()）

#### 3.3.3 与删除操作的互斥

**需求描述**：
- 有未保存的插入操作时，禁止删除
- 编辑状态下的行，Delete按钮禁用

**实现方式**：
- Delete按钮disabled条件：`disabled={editingKey === record.case_id || hasEditChanges}`
- 批量删除按钮disabled条件：`disabled={selectedRowKeys.length === 0 || editingKey !== '' || hasEditChanges}`

### 3.4 默认空行功能

#### 3.4.1 空表格默认行

**需求描述**：
- 当用例列表为空时，自动创建一条真实的空白记录
- 空白记录拥有真实的UUID，可以正常编辑和保存
- 用户可以直接编辑这条记录

**触发场景**：
1. 首次加载时表格为空
2. 清空AI用例后
3. 删除最后一条用例后
4. 批量删除所有用例后

**实现逻辑**：
1. 调用createDefaultEmptyRow函数
2. 构造正确的请求参数（case_type, language）
3. 调用后端createCase/createAutoCase API
4. 获取返回的真实记录（包含case_id UUID）
5. 设置到cases数组和pagination

**参数构造**：
```javascript
const createData = {
  case_type: caseType,  // 注意使用下划线格式
};

// 手工用例需要language参数
if (!isRoleType) {
  createData.language = language;
}
```

**失败处理**：
- 如果创建失败，显示空表格（cases=[], total=0）
- 不影响其他功能正常使用

## 4. 非功能需求

### 4.1 性能要求

- 编辑状态切换响应时间 < 100ms
- 保存操作响应时间 < 500ms
- 本地数据更新不触发整页刷新

### 4.2 用户体验

- 编辑框失去焦点时不自动保存，需用户点击Save
- 提供明确的保存/取消按钮
- 操作成功/失败有明确的提示信息
- 编辑状态下的视觉反馈清晰

### 4.3 数据一致性

- 保存前验证表单数据
- 只更新有变化的字段
- 失败时不更新本地数据
- 空值统一处理为空字符串

## 5. 验收标准

### 5.1 基本编辑功能

- [ ] 点击Edit按钮能进入编辑状态
- [ ] 编辑状态下字段显示为输入框
- [ ] 只能同时编辑一行
- [ ] 点击Save能保存修改
- [ ] 点击Cancel能取消修改
- [ ] 保存成功后显示提示信息

### 5.2 字段处理

- [ ] AI用例显示单语言字段
- [ ] 整体/变更用例显示多语言字段
- [ ] 当前语言字段正确映射
- [ ] TestResult字段选择后即时保存
- [ ] 空值正确处理为空字符串

### 5.3 状态控制

- [ ] 编辑时其他行Edit按钮禁用
- [ ] 有插入操作时Edit按钮禁用
- [ ] 编辑时Delete按钮禁用
- [ ] Above/Below点击时取消当前编辑

### 5.4 默认空行

- [ ] 空表格自动创建默认行
- [ ] 清空AI用例后创建默认行
- [ ] 删除最后一条用例后创建默认行
- [ ] 默认行拥有真实UUID
- [ ] 默认行可正常编辑保存

## 6. 测试场景

### 6.1 基本编辑测试

1. 进入编辑状态
2. 修改多个字段
3. 保存成功
4. 验证数据已更新

### 6.2 取消编辑测试

1. 进入编辑状态
2. 修改字段
3. 点击Cancel
4. 验证数据未改变

### 6.3 TestResult即时保存测试

1. 进入编辑状态
2. 修改TestResult
3. 验证立即保存成功
4. 保持编辑状态

### 6.4 多语言编辑测试

1. 切换到中文界面
2. 点击蓝色链接
3. 编辑多语言字段
4. 保存并验证三种语言都已更新

### 6.5 默认空行测试

1. 删除所有用例
2. 验证自动创建默认行
3. 编辑默认行
4. 保存成功
5. 验证数据已保存

### 6.6 边界测试

1. 编辑时切换语言（应保持编辑状态）
2. 编辑时点击Above/Below（应取消编辑）
3. 空值输入测试
4. 最大长度输入测试

## 7. 技术约束

- 使用Ant Design Form组件管理表单
- 使用React Hooks实现状态管理
- 使用useCallback优化性能
- 后端API使用PATCH方法部分更新

## 8. 依赖项

- Ant Design 5.x Table, Form, Input, Select组件
- React 18.x Hooks
- 后端API：PATCH /projects/:id/manual-cases/:caseId
- 后端API：PATCH /projects/:id/auto-cases/:caseId
- 后端API：POST /projects/:id/manual-cases
- 后端API：POST /projects/:id/auto-cases
