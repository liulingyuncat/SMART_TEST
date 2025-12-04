# Bug修复报告：整体/受入/变更用例批量修改卡死问题

## 问题描述

**对象**：整体用例/受入用例/变更用例

**问题现象**：
1. 编辑大功能(Maj.Category)或中功能(Mid.Category)后点击保存
2. 弹出批量修改确认对话框
3. 点击"批量修改"或"仅修改当前"后，界面卡死
4. 编辑的内容被清空，无法保存

**影响范围**：
- 整体用例（Overall Cases）Tab
- 受入用例（Acceptance Cases）Tab  
- 变更用例（Change Cases）Tab
- 所有语言模式（中文/日文/英文）

## 问题根本原因

### 1. 错误的 API 调用
在 `handleMultiLangSave` 函数中，批量修改确认对话框的 `onCancel` 回调使用了错误的 API：

```javascript
// 错误代码（修复前）
onCancel: async () => {
  await updateCaseAPI(record.case_id, updates);  // ❌ 应该根据类型选择API
  // ...
}
```

**问题**：`updateCaseAPI` 是内部封装的函数，但在 `onOk` 回调中直接使用了 `updateCase` API，导致不一致。

### 2. 流程控制问题
对话框关闭和数据刷新的时机不对：

```javascript
// 错误代码（修复前）
setMultiLangModalVisible(false);
fetchCases(pagination.current);
resolve();
```

**问题**：
- 没有等待 `fetchCases` 完成就 resolve
- 错误处理不完整，没有显示错误消息
- 导致 Modal 状态混乱，界面卡死

### 3. Loading 状态管理问题
`MultiLangEditModal` 组件的 loading 状态在确认对话框期间没有正确管理。

## 修复方案

### 修复1：统一 API 调用逻辑

**文件**: `EditableTable.jsx`

```javascript
// 修复后代码
onOk: async () => {
  try {
    // ✅ 根据模块和用例类型选择正确的 API
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
    
    // ✅ 关闭对话框并等待数据刷新完成
    setMultiLangModalVisible(false);
    await fetchCases(pagination.current);
    resolve();
  } catch (error) {
    console.error('批量修改失败:', error);
    message.error('批量修改失败');
    reject(error);
  }
},

onCancel: async () => {
  try {
    // ✅ 使用相同的 API 选择逻辑
    const updateAPI = apiModule === 'api-cases' ? updateApiCase : (
      (caseType && caseType.startsWith('role')) ? updateAutoCase : updateCase
    );
    await updateAPI(projectId, record.case_id, updates);
    message.success('保存成功');
    
    // ✅ 关闭对话框并等待数据刷新完成
    setMultiLangModalVisible(false);
    await fetchCases(pagination.current);
    resolve();
  } catch (error) {
    console.error('保存失败:', error);
    message.error('保存失败');
    reject(error);
  }
}
```

**关键改进**：
1. ✅ 在 `onOk` 和 `onCancel` 中使用统一的 API 选择逻辑
2. ✅ 使用 `await fetchCases()` 等待数据刷新完成
3. ✅ 添加完整的错误处理和消息提示
4. ✅ 确保 Promise 正确 resolve/reject

### 修复2：非批量修改场景的流程控制

```javascript
// 非大功能/中功能字段，或没有匹配的用例，直接保存
await updateCaseAPI(record.case_id, updates);
message.success('保存成功');

// ✅ 关闭对话框并等待数据刷新完成
setMultiLangModalVisible(false);
await fetchCases(pagination.current);
```

### 修复3：MultiLangEditModal 错误处理优化

**文件**: `MultiLangEditModal.jsx`

```javascript
} catch (error) {
  console.error('Failed to save multi-lang data:', error);
  if (error.errorFields) {
    message.error('请检查输入内容');
  } else {
    // ✅ 父组件已经显示了错误消息，避免重复显示
    // message.error('保存失败');
  }
} finally {
  // ✅ 确保加载状态被重置
  setLoading(false);
}
```

## 功能说明

### 批量修改逻辑

当用户在**中文模式**下编辑整体/受入/变更用例的大功能或中功能字段时：

1. **检测相同值**：系统查找所有用例中与当前用例大功能/中功能三语言值完全相同的用例
2. **弹出确认对话框**：如果存在相同值的用例，弹出确认对话框
3. **用户选择**：
   - **批量修改**：同时更新所有相同值的用例（包括当前用例）
   - **仅修改当前**：只更新当前正在编辑的用例

### 修改范围

- **大功能 (major_function)**：修改时会同时更新 CN/JP/EN 三个字段
- **中功能 (middle_function)**：修改时会同时更新 CN/JP/EN 三个字段
- **其他字段**：小功能、前置条件、测试步骤、期待值等字段不触发批量修改

## 测试验证

### 测试场景1：批量修改大功能

**前置条件**：
- 进入"整体用例"Tab
- 语言选择"中文"
- 存在多条大功能相同的用例

**测试步骤**：
1. 点击某条用例的"大功能"字段
2. 在弹出的多语言编辑对话框中修改 CN/JP/EN 的值
3. 点击"保存"按钮
4. **预期**：弹出确认对话框，显示匹配的用例数量
5. 点击"批量修改"
6. **预期**：
   - ✅ 显示"成功修改 X 条用例"
   - ✅ 对话框正常关闭
   - ✅ 表格数据自动刷新
   - ✅ 所有相同大功能的用例都被更新

### 测试场景2：仅修改当前用例

**测试步骤**：
1. 重复场景1的步骤1-4
2. 在确认对话框中点击"仅修改当前"
3. **预期**：
   - ✅ 显示"保存成功"
   - ✅ 对话框正常关闭
   - ✅ 表格数据自动刷新
   - ✅ 只有当前用例被更新，其他用例不变

### 测试场景3：修改中功能

**测试步骤**：
1. 点击某条用例的"中功能"字段
2. 修改 CN/JP/EN 的值并保存
3. **预期**：触发与场景1相同的批量修改逻辑

### 测试场景4：修改其他字段

**测试步骤**：
1. 点击某条用例的"小功能"、"前置条件"或其他字段
2. 修改内容并保存
3. **预期**：
   - ✅ 直接保存，不弹出批量修改对话框
   - ✅ 显示"保存成功"
   - ✅ 对话框正常关闭

### 测试场景5：日文/英文模式

**测试步骤**：
1. 切换语言到"日本語"或"English"
2. 点击"编辑"按钮进入行内编辑模式
3. 修改大功能或中功能字段
4. 点击"保存"
5. **预期**：
   - ✅ 不触发批量修改对话框（行内编辑不支持批量修改）
   - ✅ 直接保存当前用例
   - ✅ 界面不卡死

### 测试场景6：错误处理

**测试步骤**：
1. 断开网络连接或使后端API失败
2. 尝试保存修改
3. **预期**：
   - ✅ 显示错误提示
   - ✅ 对话框保持打开状态
   - ✅ 用户可以重试或取消

## 修复文件清单

1. **EditableTable.jsx** (Line 1199-1264)
   - 修复批量修改确认对话框的 API 调用
   - 统一 `onOk` 和 `onCancel` 的处理逻辑
   - 添加完整的错误处理和消息提示
   - 确保 Promise 正确 resolve/reject
   - 等待数据刷新完成后再关闭对话框

2. **MultiLangEditModal.jsx** (Line 32-56)
   - 优化错误处理，避免重复显示错误消息
   - 确保 loading 状态正确重置

## 技术要点

### 1. Promise 链式调用
```javascript
return new Promise((resolve, reject) => {
  Modal.confirm({
    onOk: async () => {
      try {
        // 执行异步操作
        await updateAPI(...);
        await fetchCases(...);  // ✅ 等待刷新完成
        resolve();  // ✅ 确保 resolve 在所有操作完成后调用
      } catch (error) {
        reject(error);  // ✅ 传递错误
      }
    }
  });
});
```

### 2. API 选择策略
```javascript
const updateAPI = apiModule === 'api-cases' ? updateApiCase : (
  (caseType && caseType.startsWith('role')) ? updateAutoCase : updateCase
);
```

### 3. 错误边界处理
- 每个异步操作都包裹在 try-catch 中
- 错误消息通过 `message.error` 显示
- 错误通过 `reject(error)` 传递给调用方

## 验证检查清单

- [x] 批量修改功能正常工作
- [x] 仅修改当前功能正常工作
- [x] 对话框不再卡死
- [x] 数据正确刷新
- [x] 错误提示正确显示
- [x] 中文模式正常
- [x] 日文/英文模式不受影响
- [x] API 调用逻辑统一
- [x] Promise 正确处理

## 相关文档

- [多语言编辑功能实现总结](./multi-lang-edit-implementation.md)
- [行内编辑功能设计文档](./T31-row-editing-design.md)
- [Bug修复报告：整体/变更用例可编辑列无法保存](./bug-fix-editable-columns-save.md)

## 更新日期

2025-11-25
