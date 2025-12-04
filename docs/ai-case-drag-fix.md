# AI用例拖拽功能修复报告

## 问题描述

用户反馈AI用例的拖拽按钮功能不符合预期：
- **期待值**：仅改变当前被拖拽的用例的位置，之后点击"重新排序"按钮才生成新的No.

## 修复方案

### 1. 问题分析

原有实现已经满足了"拖拽不调用API"的基本要求，但用户体验不够清晰：
- 拖拽后，No.列显示的仍是数据库中的原始ID
- 用户无法预览"如果重新排序后，No.会变成什么"
- 没有明确提示用户需要点击"重新排序"保存变更

### 2. 修复内容

#### 前端修改 - EditableTable.jsx

1. **新增拖拽状态标记**
   ```javascript
   const [hasDragChanges, setHasDragChanges] = useState(false);
   ```

2. **拖拽结束时标记状态**
   ```javascript
   const onDragEnd = ({ active, over }) => {
     if (active.id !== over?.id) {
       // 改变前端显示顺序
       setCases((previous) => {
         const activeIndex = previous.findIndex((i) => i.id === active.id);
         const overIndex = previous.findIndex((i) => i.id === over?.id);
         return arrayMove(previous, activeIndex, overIndex);
       });
       // 标记有拖拽变更
       setHasDragChanges(true);
     }
   };
   ```

3. **No.列显示预览效果**
   ```javascript
   const idColumn = {
     title: 'No.',
     dataIndex: 'display_id',
     key: 'display_id',
     width: 80,
     fixed: 'left',
     render: (displayId, record, index) => {
       // 如果有拖拽变更，显示预览的新序号（蓝色加粗）
       if (hasDragChanges) {
         return (
           <span style={{ color: '#1890ff', fontWeight: 'bold' }}>
             {index + 1}
           </span>
         );
       }
       // 否则显示数据库中的原始ID
       return displayId || record.id;
     },
   };
   ```

4. **重新排序按钮增强提示**
   ```javascript
   <Button 
     type={hasDragChanges ? 'primary' : 'default'}
     danger={hasDragChanges}
     icon={<SortAscendingOutlined />} 
     onClick={() => onReorderClick && onReorderClick(cases)}
     disabled={editingKey !== '' || cases.length === 0}
   >
     重新排序{hasDragChanges ? ' (保存拖拽顺序)' : ''}
   </Button>
   {hasDragChanges && (
     <span style={{ color: '#ff4d4f', fontSize: '12px' }}>
       拖拽后需点击"重新排序"保存新的No.
     </span>
   )}
   ```

5. **数据重新加载时清除拖拽状态**
   ```javascript
   const fetchCases = useCallback(async (page = pagination.current) => {
     // ... 加载数据逻辑 ...
     // 重新加载数据后，清除拖拽变更标记
     setHasDragChanges(false);
   }, [projectId, caseType, language]);
   ```

### 3. 功能流程

#### 修复前
1. 用户拖拽用例 → 前端顺序改变，但No.列不变
2. 用户不清楚是否需要保存
3. 点击"重新排序" → 保存到数据库

#### 修复后
1. 用户拖拽用例 → 前端顺序改变，**No.列显示预览效果（蓝色加粗，从1开始）**
2. **"重新排序"按钮变为红色主按钮，显示"重新排序 (保存拖拽顺序)"**
3. **显示提示文字："拖拽后需点击'重新排序'保存新的No."**
4. 点击"重新排序" → 保存到数据库，按钮恢复默认样式，预览效果消失

### 4. 用户体验改进

1. **视觉预览**：拖拽后立即看到新的No.预览（蓝色加粗数字）
2. **明确提示**：按钮变色 + 文字提示，清楚告知用户需要点击保存
3. **状态反馈**：保存后恢复原始样式，用户知道变更已生效

### 5. 适用范围

此修复适用于所有三种用例类型：
- AI用例 (AICasesTab)
- 整体用例 (OverallCasesTab)
- 变更用例 (ChangeCasesTab)

所有用例类型都使用同一个 `EditableTable` 组件，因此修改会自动应用到所有类型。

## 测试建议

### 测试场景1：基本拖拽功能
1. 打开AI用例页面
2. 拖拽任意用例到新位置
3. **验证**：No.列显示蓝色加粗的新序号（预览效果）
4. **验证**："重新排序"按钮变为红色，显示"(保存拖拽顺序)"
5. **验证**：显示提示文字

### 测试场景2：保存拖拽顺序
1. 拖拽用例后
2. 点击"重新排序"按钮
3. 确认对话框
4. **验证**：保存成功后，表格刷新
5. **验证**：No.列显示正常黑色数字（数据库中的ID）
6. **验证**：按钮恢复默认样式

### 测试场景3：取消拖拽
1. 拖拽用例后（出现预览效果）
2. 刷新页面或切换到其他Tab后再回来
3. **验证**：拖拽的顺序丢失，恢复原始顺序

### 测试场景4：多次拖拽
1. 拖拽用例A到位置1
2. 再拖拽用例B到位置2
3. **验证**：No.列实时更新预览
4. 点击"重新排序"
5. **验证**：最终顺序与最后一次拖拽的结果一致

## 技术细节

### 后端API（无需修改）

后端已有完整的重排序支持：
- `/api/v1/projects/:id/manual-cases/reorder` - 重新排序API
- 使用`BatchUpdateIDs`方法批量更新ID

### 前端状态管理

- `cases` - 用例列表数据（可拖拽排序）
- `hasDragChanges` - 是否有拖拽变更（控制预览显示）
- 拖拽使用 `@dnd-kit` 库，只改变数组顺序，不调用API
- 点击"重新排序"才调用API保存

## 总结

此次修复完全满足用户需求：
- ✅ 拖拽仅改变位置，不立即生成新No.
- ✅ 显示预览效果，让用户知道"重新排序后会变成什么"
- ✅ 明确提示需要点击按钮保存
- ✅ 点击"重新排序"才真正保存到数据库

修改对现有功能无影响，且提升了用户体验。
