# 清空AI用例按钮无响应 - 调试指南

## 问题描述
点击"清空AI用例"按钮时,控制台和Network都没有任何输出。

## 添加的调试代码

### 1. 组件挂载日志
```javascript
React.useEffect(() => {
  console.log('[AICasesTab] Component mounted, projectId:', projectId);
  return () => {
    console.log('[AICasesTab] Component unmounted');
  };
}, [projectId]);
```

### 2. 按钮点击日志
```javascript
const handleClearAICases = () => {
  console.log('[AICasesTab] Clear button clicked, projectId:', projectId);
  // ... Modal.confirm
};
```

### 3. 测试按钮
添加了一个简单的测试按钮,用于验证按钮点击事件是否正常工作。

## 调试步骤

### Step 1: 检查组件是否正确渲染
1. 打开浏览器开发者工具(F12)
2. 切换到Console标签
3. 导航到项目详情页 -> 手工用例
4. 点击"AI用例"Tab

**期望输出**:
```
[AICasesTab] Component mounted, projectId: <项目ID>
```

**如果没有输出**: 
- 组件没有被渲染
- 检查Tab配置和路由

### Step 2: 测试基本点击事件
1. 点击"测试按钮"

**期望结果**:
- 控制台输出: `Test button clicked!`
- 页面显示消息: "测试按钮点击成功!"

**如果测试按钮能点击,清空按钮不能**:
- 说明是Modal.confirm的问题
- 或者是handleClearAICases函数的问题

**如果测试按钮也不能点击**:
- 可能是CSS遮挡问题
- 可能是事件冒泡被阻止
- 检查EditableTable组件是否覆盖了按钮区域

### Step 3: 检查清空按钮点击
1. 点击"清空AI用例"按钮

**期望输出**:
```
[AICasesTab] Clear button clicked, projectId: <项目ID>
```

**如果有输出但没有弹出确认框**:
- Modal组件问题
- 检查antd版本
- 检查Modal导入

**如果点击后有确认框弹出**:
1. 点击"取消"
   - 期望输出: `[AICasesTab] User cancelled clear operation`

2. 点击"确认清空"
   - 期望输出: `[AICasesTab] User confirmed clear operation`
   - 期望输出: `[AICasesTab] Clear API response: {...}`
   - 期望Network面板有DELETE请求到 `/api/v1/projects/{id}/manual-cases/clear-ai`

## 可能的问题和解决方案

### 问题1: 组件没有挂载
**症状**: 控制台没有"Component mounted"日志

**原因**:
- Tab切换逻辑问题
- 组件导入错误

**解决**:
```javascript
// 检查 ManualTestTabs/index.jsx
const MANUAL_TEST_TABS = [
  { key: 'ai-cases', labelKey: 'manualTest.aiCases', component: AICasesTab },
  // ...
];
```

### 问题2: 按钮被遮挡
**症状**: 测试按钮也不能点击

**检查**:
```javascript
// 在浏览器开发者工具的Elements面板
// 选中按钮元素,检查Computed样式
// 查看 z-index, position, pointer-events
```

**解决**: 已添加 `style={{ zIndex: 1000, position: 'relative' }}`

### 问题3: projectId为undefined
**症状**: 点击后projectId显示undefined

**原因**: 
- useParams()没有正确获取路由参数
- 路由配置错误

**检查**:
```javascript
// 在组件中添加
console.log('URL:', window.location.pathname);
console.log('projectId:', projectId);
```

### 问题4: Modal不显示
**症状**: 有点击日志,但没有确认框

**原因**:
- antd Modal组件问题
- Modal被其他元素遮挡

**解决**:
```javascript
// 测试直接调用message
onClick={() => {
  message.info('按钮被点击了!');
}}
```

### 问题5: API调用失败但无错误
**症状**: 确认后没有Network请求

**检查**:
```javascript
// 在clearAICases函数中
export const clearAICases = async (projectId) => {
  console.log('[API] clearAICases called, projectId:', projectId);
  const response = await client.delete(`/projects/${projectId}/manual-cases/clear-ai`);
  console.log('[API] clearAICases response:', response);
  return response.data;
};
```

## 快速测试命令

### 测试1: 检查组件是否存在
在浏览器Console中运行:
```javascript
document.querySelector('.ai-cases-tab')
// 应该返回DOM元素,如果返回null则组件未渲染
```

### 测试2: 检查按钮是否存在
```javascript
document.querySelector('.ai-cases-tab button')
// 应该返回按钮元素
```

### 测试3: 手动触发点击
```javascript
const btn = document.querySelector('.ai-cases-tab button');
if (btn) {
  btn.click();
  console.log('Button clicked programmatically');
}
```

## 下一步操作

根据上述调试步骤的结果:

1. **如果组件没有挂载** → 检查Tab配置和路由
2. **如果按钮不可点击** → 检查CSS和z-index
3. **如果点击无日志** → 检查事件绑定
4. **如果有日志但无Modal** → 检查Modal组件
5. **如果Modal显示但无API请求** → 检查API函数

## 当前修改的文件

- `frontend/src/pages/ProjectDetail/ManualTestTabs/containers/AICasesTab.jsx`
  - 添加了组件挂载日志
  - 添加了按钮点击日志
  - 添加了测试按钮
  - 添加了z-index样式

记得在调试完成后移除测试按钮和多余的console.log!
