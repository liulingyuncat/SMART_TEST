# AI用例重排序功能修复报告

## 问题描述

当前AI用例在分页情况下，重新排序功能存在bug：

- 当前设置：10条/页
- 操作步骤：
  1. 翻到第二页（显示No.11-20）
  2. 点击"插入行"，新行编号为No.88
  3. 此时第二页显示：No.88, No.10, No.11
  4. 点击"重新排序"按钮

**错误结果**：No.88变成了No.4

**期望结果**：No.88应该变成No.11（因为它在第二页第一条的位置）

## 问题根因

前端`ReorderModal`组件的重排序逻辑有误：

1. 使用`getCasesList`获取所有用例（传入size: 10000）
2. 将获取到的用例按ID排序：`sort((a, b) => idA - idB)`
3. 将排序后的case_id数组传给后端

这导致插入在第二页的No.88（ID=88）被排序到第88位，而不是保持在第二页的位置。

## 解决方案

### 方案说明

后端已有`ReorderAllCasesByID`方法，可以按现有ID顺序重新编号所有用例。但缺少对应的API路由。

**修复方案**：
1. 添加后端API路由：`POST /api/v1/projects/:id/manual-cases/reorder-all`
2. 修复后端service层方法调用
3. 添加前端API接口函数
4. 修改`ReorderModal`直接调用新API，不再前端排序

### 实现细节

#### 1. 后端路由添加

**文件**：`backend/cmd/server/main.go`

```go
projects.POST("/:id/manual-cases/reorder-all",
    middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
    manualCasesHandler.ReorderAllCasesByID)
```

#### 2. 修复Service层方法

**文件**：`backend/internal/services/manual_test_case_service.go`

修改前：
```go
allCases, err := s.repo.GetAllByProjectTypeAndLanguage(projectID, caseType, language)
```

修改后：
```go
allCases, err := s.repo.GetByProjectAndTypeOrdered(projectID, caseType)
```

使用已有的`GetByProjectAndTypeOrdered`方法，该方法会按ID升序返回所有用例。

#### 3. 前端API接口

**文件**：`frontend/src/api/manualCase.js`

```javascript
/**
 * 按现有ID顺序重新编号所有用例（用于重新排序按钮）
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change'|'ai')
 * @param {string} language - 语言 ('中文'|'English'|'日本語')
 * @returns {Promise<{count: number}>}
 */
export const reorderAllCasesByID = async (projectId, caseType, language) => {
  const response = await client.post(`/projects/${projectId}/manual-cases/reorder-all`, {
    case_type: caseType,
    language: language
  });
  return response.data;
};
```

#### 4. 修改ReorderModal组件

**文件**：`frontend/src/pages/ProjectDetail/ManualTestTabs/components/ReorderModal.jsx`

修改前（有问题的逻辑）：
```javascript
// 1. 获取所有用例
const data = await getCasesList(projectId, { 
  caseType, 
  language, 
  page: 1, 
  size: 10000
});

// 2. 按现有ID排序
const sortedCases = [...data.cases].sort((a, b) => {
  const idA = a.display_id || a.id;
  const idB = b.display_id || b.id;
  return idA - idB;
});

// 3. 提取case_id数组
const caseIds = sortedCases.map(c => c.case_id);

// 4. 调用后端API
await reorderCasesByDrag(projectId, caseType, caseIds);
```

修改后（正确逻辑）：
```javascript
// 直接调用后端的重新编号API
// 后端会获取所有用例，按现有ID排序，然后重新分配连续的ID（1,2,3...）
const result = await reorderAllCasesByID(projectId, caseType, language);

const count = result.count || 0;
if (count === 0) {
  message.warning('没有可重排的用例');
  return;
}

message.success(`成功重排 ${count} 条用例，ID已重新编号为 1-${count}`);
```

## 修改文件清单

### 后端
1. `backend/cmd/server/main.go` - 添加路由
2. `backend/internal/services/manual_test_case_service.go` - 修复方法调用

### 前端
1. `frontend/src/api/manualCase.js` - 添加API接口函数
2. `frontend/src/pages/ProjectDetail/ManualTestTabs/components/ReorderModal.jsx` - 使用新API

## 测试步骤

### 手动测试

1. **启动后端服务**
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. **启动前端服务**
   ```bash
   cd frontend
   npm start
   ```

3. **测试流程**
   - 登录系统，进入项目的AI用例页
   - 设置每页显示10条
   - 翻到第二页（显示No.11-20）
   - 在第二页点击"插入行"按钮
   - 新行的No应该是一个较大的数字（如No.88）
   - 点击"重新排序"按钮
   
4. **验证结果**
   - ✅ **期望**：No.88应该变成No.11（第二页第一条）
   - ❌ **错误**（修复前）：No.88变成No.4

### API测试（可选）

使用curl测试新API：

```bash
curl -X POST http://localhost:8080/api/v1/projects/1/manual-cases/reorder-all \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "case_type": "ai",
    "language": "中文"
  }'
```

期望返回：
```json
{
  "message": "重新编号成功",
  "count": 25
}
```

## 技术说明

### 为什么需要后端专用排序API？

1. **数据一致性**：后端直接从数据库按ID排序获取所有数据，确保顺序正确
2. **性能优化**：避免前端获取大量数据（10000条），减少网络传输
3. **简化逻辑**：前端不需要关心排序逻辑，直接调用后端API
4. **事务保证**：后端在事务中批量更新ID，保证原子性

### 后端排序逻辑

后端`ReorderAllCasesByID`方法：
1. 调用`GetByProjectAndTypeOrdered`获取所有用例（按ID升序）
2. 提取case_id数组
3. 调用`BatchUpdateIDsByCaseID`批量更新ID（按数组顺序分配1,2,3...）

这样确保了按现有ID顺序重新编号，插入的新行会保持在正确的位置。

## 总结

本次修复通过添加后端专用排序API，解决了前端分页情况下重排序逻辑错误的问题。修改简洁高效，符合前后端分离的设计原则。
