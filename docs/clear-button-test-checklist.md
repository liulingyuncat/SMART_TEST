# 清空AI用例按钮 - 快速测试清单

## ✅ 测试清单

### 1️⃣ 打开页面
- [ ] 访问项目详情页
- [ ] 点击"手工用例"菜单
- [ ] 切换到"AI用例"Tab

### 2️⃣ 检查控制台
打开浏览器开发者工具(F12) → Console面板

**预期看到**:
```
[AICasesTab] Component mounted, projectId: <数字>
```

**❌ 如果没有看到**:
→ 组件未渲染,检查Tab配置
→ 跳转到[问题排查](#组件未渲染)

### 3️⃣ 测试基本点击
- [ ] 点击"测试按钮"

**预期看到**:
- 控制台: `Test button clicked!`
- 页面消息: "测试按钮点击成功!"

**❌ 如果测试按钮也不能点击**:
→ 按钮被遮挡或事件被阻止
→ 跳转到[问题排查](#按钮不可点击)

### 4️⃣ 测试清空按钮
- [ ] 点击"清空AI用例"按钮

**预期看到**:
- 控制台: `[AICasesTab] Clear button clicked, projectId: <数字>`
- 弹出确认对话框

**❌ 如果没有控制台输出**:
→ 按钮事件没有绑定
→ 跳转到[问题排查](#事件未绑定)

**❌ 如果有输出但没有对话框**:
→ Modal组件问题
→ 跳转到[问题排查](#Modal不显示)

### 5️⃣ 测试确认操作
在确认对话框中:

- [ ] 点击"取消"
  - 预期: 控制台显示 `[AICasesTab] User cancelled clear operation`

- [ ] 再次点击清空按钮,然后点击"确认清空"
  - 预期: 控制台显示 `[AICasesTab] User confirmed clear operation`

### 6️⃣ 检查API调用
点击确认后,检查:

**控制台 - 预期输出**:
```
[AICasesTab] User confirmed clear operation
[API] clearAICases called, projectId: <数字>
[API] Request URL: /projects/<数字>/manual-cases/clear-ai
[API] clearAICases response: {...}
[AICasesTab] Clear API response: {...}
```

**Network面板 - 预期请求**:
- Method: DELETE
- URL: `http://localhost:8080/api/v1/projects/<id>/manual-cases/clear-ai`
- Status: 200
- Response: `{"code": 200, "data": {...}, "message": "success"}`

**❌ 如果没有API调用**:
→ API函数问题
→ 跳转到[问题排查](#API未调用)

### 7️⃣ 验证结果
- [ ] 页面显示成功消息
- [ ] 表格自动刷新

---

## 🔍 问题排查

### 组件未渲染
**现象**: 没有"Component mounted"日志

**排查步骤**:
1. 在Console中运行: `document.querySelector('.ai-cases-tab')`
2. 如果返回`null`,组件确实未渲染

**可能原因**:
- Tab配置错误
- 路由参数丢失
- 组件导入失败

**解决方案**:
```javascript
// 检查 src/pages/ProjectDetail/ManualTestTabs/index.jsx
const MANUAL_TEST_TABS = [
  { key: 'ai-cases', labelKey: 'manualTest.aiCases', component: AICasesTab },
];

// 确认导入
import AICasesTab from './containers/AICasesTab';
```

### 按钮不可点击
**现象**: 测试按钮也无法点击

**排查步骤**:
1. 在Elements面板选中按钮
2. 查看Computed样式
3. 检查 `pointer-events`, `z-index`, `opacity`

**解决方案**:
- 已添加 `style={{ zIndex: 1000, position: 'relative' }}`
- 检查是否有其他元素覆盖

### 事件未绑定
**现象**: 点击无任何控制台输出

**排查步骤**:
```javascript
// 在Console中测试
const btn = document.querySelectorAll('.ai-cases-tab button')[0];
console.log('Button element:', btn);
console.log('onClick handler:', btn.onclick);
```

**解决方案**:
- 确认 `onClick={handleClearAICases}` 正确书写
- 重新启动开发服务器

### Modal不显示
**现象**: 有点击日志但无对话框

**排查步骤**:
1. 检查antd版本: `npm list antd`
2. 临时替换为简单alert测试

```javascript
// 临时测试代码
const handleClearAICases = () => {
  console.log('[AICasesTab] Clear button clicked');
  alert('Button clicked!'); // 如果这个能显示,说明Modal有问题
};
```

**解决方案**:
- 升级antd到最新版本
- 检查Modal的z-index
- 检查是否有CSS冲突

### API未调用
**现象**: 确认后没有Network请求

**排查步骤**:
1. 检查clearAICases导入: `import { clearAICases } from '../../../../api/manualCase';`
2. 在API文件中添加调试日志(已添加)

**解决方案**:
- 确认API函数导出: `export const clearAICases`
- 检查client配置
- 测试简化版本:
```javascript
onOk: async () => {
  console.log('About to call API...');
  fetch(`http://localhost:8080/api/v1/projects/${projectId}/manual-cases/clear-ai`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
    }
  }).then(r => console.log('Direct fetch result:', r));
};
```

---

## 📋 调试代码位置

所有调试代码已添加到:

1. **组件文件**: `frontend/src/pages/ProjectDetail/ManualTestTabs/containers/AICasesTab.jsx`
   - 组件挂载日志: 第23-28行
   - 按钮点击日志: 第44行
   - 确认/取消日志: 第56、69行
   - 测试按钮: 第88-95行

2. **API文件**: `frontend/src/api/manualCase.js`
   - API调用日志: 第119-121行

---

## 🧹 清理调试代码

测试完成后,请移除以下内容:

### AICasesTab.jsx
```javascript
// 删除 useEffect
React.useEffect(() => {
  console.log('[AICasesTab] Component mounted, projectId:', projectId);
  return () => {
    console.log('[AICasesTab] Component unmounted');
  };
}, [projectId]);

// 删除所有 console.log
console.log('[AICasesTab] Clear button clicked, projectId:', projectId);
console.log('[AICasesTab] User confirmed clear operation');
console.log('[AICasesTab] Clear API response:', result);
console.log('[AICasesTab] User cancelled clear operation');

// 删除测试按钮
<Button onClick={...}>测试按钮</Button>

// 删除 style 属性(如果不需要)
style={{ zIndex: 1000, position: 'relative' }}
```

### manualCase.js
```javascript
// 删除所有 console.log
console.log('[API] clearAICases called, projectId:', projectId);
console.log('[API] Request URL:', ...);
console.log('[API] clearAICases response:', response);
```

---

## ✨ 成功标准

所有以下测试都通过:
- ✅ 组件正常挂载
- ✅ 测试按钮可点击
- ✅ 清空按钮可点击
- ✅ 确认对话框正常显示
- ✅ API请求正常发送
- ✅ 响应正常接收
- ✅ 页面正常更新

如果所有测试通过,说明问题已解决! 🎉
