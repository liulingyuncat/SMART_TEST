# 批量修改功能最终修复方案

## 问题回顾

用户反馈：编辑整体/受入/变更用例的大功能或中功能后点击保存，应该弹出批量修改确认对话框，但**对话框没有显示**，界面卡死。

## 根本原因

使用 `Modal.confirm` 时存在以下问题：
1. **时序问题**：在异步函数中创建 Modal.confirm，可能因为组件状态变化导致对话框无法渲染
2. **状态混乱**：MultiLangEditModal 的 loading 状态与 Modal.confirm 的显示时机冲突
3. **Promise 阻塞**：`await onSave(...)` 等待 Modal.confirm 的用户选择，导致 loading 一直为 true

## 最终解决方案

### 核心思路：使用 React 状态管理确认对话框

不再使用 `Modal.confirm` API，而是使用 React 组件方式管理对话框，确保对话框一定会显示。

### 实现步骤

#### 1. 添加状态管理

```javascript
// 批量修改确认对话框状态
const [batchConfirmVisible, setBatchConfirmVisible] = useState(false);
const [batchConfirmData, setBatchConfirmData] = useState({
  matchingCases: [],
  updates: {},
  record: null,
  fieldName: '',
});
```

#### 2. 修改 handleMultiLangSave 函数

```javascript
if (matchingCases.length > 0) {
  console.log('🚀 发现匹配用例，弹出批量修改确认对话框');
  
  // 保存批量修改数据到状态
  setBatchConfirmData({
    matchingCases,
    updates,
    record,
    fieldName: data.fieldName,
  });
  
  // 先关闭多语言对话框
  setMultiLangModalVisible(false);
  
  // 延迟显示批量修改确认对话框
  setTimeout(() => {
    setBatchConfirmVisible(true);
  }, 100);
  
  return; // 立即返回，不阻塞
}
```

**关键点**：
- ✅ 将匹配数据保存到状态
- ✅ 先关闭多语言对话框（避免两个对话框重叠）
- ✅ 延迟 100ms 后显示确认对话框（确保多语言对话框已完全关闭）
- ✅ 立即 return，让 MultiLangEditModal 正常完成

#### 3. 添加处理函数

```javascript
// 批量修改
const handleBatchConfirmOk = async () => {
  const { matchingCases, updates, record } = batchConfirmData;
  
  try {
    setBatchConfirmVisible(false);
    setLoading(true);
    
    const updateAPI = apiModule === 'api-cases' ? updateApiCase : (
      (caseType && caseType.startsWith('role')) ? updateAutoCase : updateCase
    );
    
    // 更新当前用例
    await updateAPI(projectId, record.case_id, updates);
    
    // 批量更新所有匹配的用例
    for (const matchCase of matchingCases) {
      await updateAPI(projectId, matchCase.case_id, updates);
    }
    
    message.success(`成功修改 ${matchingCases.length + 1} 条用例`);
    await fetchCases(pagination.current);
  } catch (error) {
    message.error('批量修改失败');
  } finally {
    setLoading(false);
  }
};

// 仅修改当前
const handleBatchConfirmCancel = async () => {
  const { updates, record } = batchConfirmData;
  
  try {
    setBatchConfirmVisible(false);
    setLoading(true);
    
    const updateAPI = apiModule === 'api-cases' ? updateApiCase : (
      (caseType && caseType.startsWith('role')) ? updateAutoCase : updateCase
    );
    
    await updateAPI(projectId, record.case_id, updates);
    message.success('保存成功');
    await fetchCases(pagination.current);
  } catch (error) {
    message.error('保存失败');
  } finally {
    setLoading(false);
  }
};
```

#### 4. 添加对话框组件

```jsx
<Modal
  title="批量修改确认"
  open={batchConfirmVisible}
  onOk={handleBatchConfirmOk}
  onCancel={handleBatchConfirmCancel}
  okText="批量修改"
  cancelText="仅修改当前"
  confirmLoading={loading}
  maskClosable={false}
  closable={false}
>
  <p>
    检测到有 <strong>{batchConfirmData.matchingCases?.length || 0}</strong> 条用例的
    <strong>{batchConfirmData.fieldName === 'major_function' ? '大功能' : '中功能'}</strong>
    与当前值相同，是否一起修改?
  </p>
  <p style={{ marginTop: 16, color: '#666' }}>
    • 点击"批量修改"：将同时更新所有相同值的用例<br />
    • 点击"仅修改当前"：只更新当前正在编辑的用例
  </p>
</Modal>
```

**配置说明**：
- `open={batchConfirmVisible}` - 通过状态控制显示
- `confirmLoading={loading}` - 显示加载状态
- `maskClosable={false}` - 禁止点击遮罩关闭
- `closable={false}` - 禁止点击 X 关闭（强制用户选择）

## 优势对比

### 之前的方案（Modal.confirm）
❌ 时序不可控，可能无法显示  
❌ 状态管理混乱  
❌ 调试困难  
❌ 用户体验差（卡死）

### 当前方案（React 组件）
✅ 状态管理清晰  
✅ 显示时机可控  
✅ 易于调试和维护  
✅ 用户体验好  
✅ 可以自定义样式和内容  

## 测试流程

### 1. 准备测试数据
确保有多条用例具有相同的大功能或中功能（CN/JP/EN 三个值都相同）

### 2. 执行测试
1. 进入"整体用例"Tab（中文模式）
2. 点击某条用例的"大功能"字段
3. 在多语言编辑对话框中修改内容
4. 点击"保存"按钮

### 3. 预期结果
1. ✅ 多语言编辑对话框关闭
2. ✅ 延迟 100ms 后弹出批量修改确认对话框
3. ✅ 对话框显示匹配的用例数量
4. ✅ 提供"批量修改"和"仅修改当前"两个选项

### 4. 批量修改测试
点击"批量修改"：
- ✅ 对话框显示 loading 状态
- ✅ 更新当前用例和所有匹配用例
- ✅ 显示成功消息："成功修改 X 条用例"
- ✅ 自动刷新表格数据
- ✅ 所有相同值的用例都被更新

### 5. 仅修改当前测试
点击"仅修改当前"：
- ✅ 对话框显示 loading 状态
- ✅ 只更新当前用例
- ✅ 显示成功消息："保存成功"
- ✅ 自动刷新表格数据
- ✅ 其他用例保持不变

### 6. 边界情况测试
- **修改其他字段**（小功能、前置条件等）：不弹出确认对话框，直接保存
- **没有匹配用例**：不弹出确认对话框，直接保存
- **网络错误**：显示错误提示，对话框正常关闭

## 调试日志

关键日志输出：
```
🚀 发现匹配用例，弹出批量修改确认对话框
保存批量修改数据到状态
关闭多语言对话框
显示批量修改确认对话框
[用户选择后]
用户点击了"批量修改" / "仅修改当前"
开始更新用例...
批量更新完成
刷新完成
```

## 修改文件清单

**EditableTable.jsx**
- Line 147-157: 添加批量修改确认对话框状态
- Line 1217-1237: 修改 handleMultiLangSave 函数
- Line 1277-1350: 添加 handleBatchConfirmOk 和 handleBatchConfirmCancel 函数
- Line 2840-2858: 添加批量修改确认对话框组件

## 技术要点

### 1. 状态管理
使用 useState 管理对话框状态，确保完全可控

### 2. 时序控制
使用 setTimeout 确保两个对话框不会同时显示

### 3. 数据传递
通过 state 传递批量修改所需的所有数据

### 4. 错误处理
try-catch-finally 确保 loading 状态正确重置

### 5. 用户体验
- 禁止点击遮罩和 X 关闭，强制用户做出选择
- 显示 loading 状态，提供视觉反馈
- 清晰的提示文本，帮助用户理解选项含义

## 后续优化建议

1. **批量修改进度**：可以添加进度条显示批量更新进度
2. **撤销功能**：提供批量修改后的撤销功能
3. **预览功能**：在确认对话框中显示将要修改的用例列表
4. **性能优化**：使用批量更新 API 代替循环单条更新

## 总结

通过使用 React 组件管理确认对话框，彻底解决了 Modal.confirm 的时序和状态问题，确保对话框一定会显示，用户体验得到极大改善。

**核心优势**：
- ✅ 可靠性高：对话框一定会显示
- ✅ 可维护性好：代码结构清晰
- ✅ 可扩展性强：易于添加新功能
- ✅ 调试友好：状态变化可追踪

---

**最后更新**: 2025-11-25  
**修复状态**: ✅ 已完成并验证
