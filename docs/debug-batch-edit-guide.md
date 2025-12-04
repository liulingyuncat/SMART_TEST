# 批量修改功能调试指南

## 调试代码已添加

已在以下位置添加详细的调试日志：

### 1. EditableTable.jsx - handleMultiLangSave 函数
- ✅ 开始和结束标记
- ✅ 输入数据记录
- ✅ 字段类型检查
- ✅ 获取所有用例
- ✅ 匹配用例查找
- ✅ 批量修改确认对话框
- ✅ API 调用追踪
- ✅ 数据刷新追踪

### 2. MultiLangEditModal.jsx - handleOk 函数
- ✅ 开始和结束标记
- ✅ Loading 状态追踪
- ✅ 表单验证追踪
- ✅ onSave 回调追踪

## 如何测试

### 步骤1：打开浏览器开发者工具
1. 按 F12 打开开发者工具
2. 切换到 "Console" (控制台) 标签
3. 清空现有日志

### 步骤2：执行测试操作
1. 进入"手工测试用例"页面
2. 切换到"整体用例"Tab
3. 确保语言选择为"中文"
4. 点击某条用例的"大功能"字段
5. 在弹出的多语言编辑对话框中修改内容
6. 点击"保存"按钮

### 步骤3：观察控制台输出

#### 期望的日志流程：

```
🔵 MultiLangEditModal handleOk 开始
设置 loading = true
表单验证成功，值: {cn: "...", jp: "...", en: "..."}
调用 onSave 回调, fieldName: major_function

=== handleMultiLangSave 开始 ===
data: {fieldName: "major_function", cn: "...", jp: "...", en: "..."}
multiLangData: {record: {...}, fieldName: "major_function", ...}
record: {case_id: "...", ...}
updates: {major_function_cn: "...", major_function_jp: "...", major_function_en: "..."}
isFunctionField: true fieldName: major_function
检查批量更新...
使用的 listAPI: getCasesList
allCasesResponse: {...}
获取到的所有用例数量: X
原始值 CN: xxx JP: xxx EN: xxx
匹配的用例数量: X
匹配的用例: [...]

# 如果有匹配的用例：
🚀 发现匹配用例，弹出批量修改确认对话框
创建 Modal.confirm
Modal.confirm 已创建

# 用户点击"批量修改"后：
用户点击了"批量修改"
使用的 updateAPI: updateCase
开始更新当前用例: xxx
当前用例更新完成
开始批量更新其他用例: X
更新用例: xxx
...
批量更新完成
关闭多语言对话框
开始刷新用例列表
刷新完成

onSave 回调完成
表单已重置
🔵 MultiLangEditModal handleOk 成功结束
设置 loading = false

# 或用户点击"仅修改当前"后：
用户点击了"仅修改当前"
使用的 updateAPI: updateCase
开始更新当前用例: xxx
当前用例更新完成
关闭多语言对话框
开始刷新用例列表
刷新完成

onSave 回调完成
表单已重置
🔵 MultiLangEditModal handleOk 成功结束
设置 loading = false

# 如果没有匹配的用例：
没有匹配的用例，直接保存
执行直接保存逻辑
保存完成
关闭多语言对话框
开始刷新用例列表
刷新完成
=== handleMultiLangSave 结束 ===

onSave 回调完成
表单已重置
🔵 MultiLangEditModal handleOk 成功结束
设置 loading = false
```

## 常见问题诊断

### 问题1：没有弹出批量修改确认对话框

**检查日志**：
- 查找 `isFunctionField: true` - 如果是 false，说明不是大功能/中功能字段
- 查找 `匹配的用例数量: X` - 如果是 0，说明没有匹配的用例
- 查找 `🚀 发现匹配用例` - 如果没有这行，说明没有进入批量修改逻辑

**可能原因**：
1. 修改的不是大功能或中功能字段
2. 没有其他用例与当前值完全相同（CN/JP/EN 三个值都要相同）
3. 数据库中的值与当前显示的值不一致

### 问题2：对话框卡死

**检查日志**：
- 查找 `onSave 回调完成` - 如果没有这行，说明 onSave 没有正常返回
- 查找 `❌` 标记的错误日志
- 查找 `刷新完成` - 如果没有这行，说明数据刷新被卡住

**可能原因**：
1. API 调用失败但没有正确处理错误
2. fetchCases 函数执行时间过长或卡住
3. Promise 没有正确 resolve/reject

### 问题3：保存后数据没有刷新

**检查日志**：
- 查找 `开始刷新用例列表` 和 `刷新完成`
- 如果这两行都出现了，说明刷新逻辑执行了
- 检查网络请求，看是否真的发送了更新请求

## 调试技巧

### 1. 设置断点
在以下位置设置断点：
- `handleMultiLangSave` 函数开始处
- `Modal.confirm` 创建处
- `onOk` 和 `onCancel` 回调函数内部

### 2. 检查网络请求
1. 切换到 "Network" (网络) 标签
2. 筛选 XHR/Fetch 请求
3. 观察更新用例的 API 请求
4. 检查请求参数和响应

### 3. 查看 React DevTools
1. 安装 React DevTools 扩展
2. 检查组件状态：
   - `multiLangModalVisible`
   - `multiLangData`
   - `cases`
   - `loading`

## 下一步操作

根据控制台输出的日志，确定问题出现在哪个环节：

1. **如果日志显示正常流程**：问题可能在 UI 渲染层面
2. **如果日志在某个环节停止**：问题在该环节的逻辑
3. **如果出现错误日志**：根据错误信息修复具体问题

## 联系支持

如果问题依然存在，请提供：
1. 完整的控制台日志
2. 网络请求详情
3. 具体的操作步骤
4. 截图或录屏

---

**最后更新**: 2025-11-25
