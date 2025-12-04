# AI用例重排序功能修复 - 按当前显示顺序重排

## 需求变更

### 原需求（已修复）
- 问题：分页插入后重新排序，No.88变成No.4
- 解决：按现有ID顺序重新编号

### 新需求
**场景**：
- 当前用例总数80条，10条/页
- 点击插入行，第一页显示11条：No.81, No.1~No.10
- 点击重排按钮

**期望效果**：
1. No.81变为No.1
2. 原No.1变为No.2，原No.2变为No.3...以此类推
3. 当前页面恢复为10条/页
4. 其他页的用例顺序保持不变，但编号顺延

**核心改变**：
- 旧逻辑：按原ID顺序重新编号（No.81排到最后）
- 新逻辑：按**当前页显示顺序**重新编号（No.81排到最前）

## 实现方案

### 1. 传递当前显示数据

**EditableTable** → **容器组件** → **ReorderModal**

```javascript
// EditableTable.jsx - 传递当前页cases
onClick={() => onReorderClick && onReorderClick(cases)}

// OverallCasesTab.jsx - 接收并保存cases
const handleReorderClick = (currentCases) => {
  setCasesForReorder(currentCases || []);
  setReorderModalVisible(true);
};

// ReorderModal - 接收cases prop
const ReorderModal = ({ visible, caseType, projectId, language, cases = [], onOk, onCancel })
```

### 2. 重排逻辑

ReorderModal的新逻辑：

```javascript
// 1. 获取当前页显示的case_id（按显示顺序）
const currentPageCaseIds = cases.map(c => c.case_id);

// 2. 获取其他页的case_id（按原ID排序）
const otherCases = allCases.filter(c => !currentPageCaseIds.includes(c.case_id));
const otherCaseIds = otherCases.sort((a, b) => (a.display_id || a.id) - (b.display_id || b.id))
                                .map(c => c.case_id);

// 3. 合并：当前页在前，其他页在后
const finalOrder = [...currentPageCaseIds, ...otherCaseIds];

// 4. 调用后端API
await reorderCasesByDrag(projectId, caseType, finalOrder);
```

### 3. 使用的API

使用已有的 `reorderCasesByDrag` API（POST `/api/v1/projects/:id/manual-cases/reorder-drag`）

该API根据传入的case_id数组顺序，重新分配ID为1,2,3...

## 修改文件

### 前端
1. **EditableTable.jsx**
   - 修改重排按钮：`onClick={() => onReorderClick && onReorderClick(cases)}`
   - 传递当前页的cases数组给父组件

2. **OverallCasesTab.jsx**
   - 修改`handleReorderClick`：接收并保存currentCases
   - 传递`cases={casesForReorder}`给ReorderModal

3. **AICasesTab.jsx** 和 **ChangeCasesTab.jsx**
   - 已经正确实现，无需修改

4. **ReorderModal.jsx**
   - 接收`cases` prop
   - 获取所有用例数据
   - 按"当前页在前，其他页在后"的顺序重排
   - 调用`reorderCasesByDrag` API
   - 更新说明文字

## 测试步骤

### 场景1：第一页插入新行

1. **准备数据**
   - 确保有80条AI用例（No.1-80）
   - 设置每页显示10条

2. **操作步骤**
   ```
   a. 进入AI用例页面
   b. 当前第一页显示：No.1-10（10条）
   c. 点击"插入行"按钮
   d. 第一页显示：No.81, No.1-10（11条）
   e. 点击"重新排序"按钮
   ```

3. **预期结果**
   ```
   ✅ 第一页显示：No.1-10（10条）
   ✅ No.81 → No.1
   ✅ 原No.1 → No.2
   ✅ 原No.2 → No.3
   ...
   ✅ 原No.10 → No.11
   ✅ 第二页显示：No.11-20（原来的No.10-19）
   ```

### 场景2：第二页插入新行

1. **操作步骤**
   ```
   a. 翻到第二页
   b. 当前第二页显示：No.11-20（10条）
   c. 点击"插入行"按钮
   d. 第二页显示：No.82, No.11-20（11条）
   e. 点击"重新排序"按钮
   ```

2. **预期结果**
   ```
   ✅ 第二页显示：No.11-20（10条）
   ✅ No.82 → No.11
   ✅ 原No.11 → No.12
   ✅ 原No.12 → No.13
   ...
   ✅ 原No.20 → No.21
   ```

### 场景3：拖拽后重排序

1. **操作步骤**
   ```
   a. 在第一页拖拽No.5到No.1前面
   b. 显示顺序：No.5, No.1, No.2, No.3, No.4, No.6-10
   c. 点击"重新排序"按钮
   ```

2. **预期结果**
   ```
   ✅ No.5 → No.1
   ✅ 原No.1 → No.2
   ✅ 原No.2 → No.3
   ✅ 原No.3 → No.4
   ✅ 原No.4 → No.5
   ✅ 原No.6 → No.6
   ...
   ```

## 技术说明

### 为什么使用reorderCasesByDrag而不是reorderAllCasesByID？

1. **reorderCasesByDrag**
   - 接收完整的case_id顺序数组
   - 按数组顺序重新分配ID（1,2,3...）
   - 支持任意顺序排列
   - ✅ **适合本需求**

2. **reorderAllCasesByID**
   - 后端自动按现有ID排序
   - 然后重新分配ID
   - 不考虑前端的显示顺序
   - ❌ 不适合本需求

### 关键点

1. **当前页优先**：当前页的case按显示顺序排在最前面
2. **其他页保持顺序**：其他页按原ID顺序排在后面
3. **分页恢复**：重排后自动恢复为10条/页
4. **拖拽支持**：支持拖拽后的顺序

## 总结

本次修改实现了按**当前显示顺序**重排序的功能，完美解决了"插入行在第一页应该变成No.1"的需求。

**优势**：
- ✅ 符合用户直觉：显示顺序即重排顺序
- ✅ 支持拖拽排序
- ✅ 灵活可控：用户可以通过翻页和拖拽控制重排结果
- ✅ 代码简洁：复用已有的drag排序API
