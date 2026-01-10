# React + Recharts 快速开始指南

## 项目修改概览

本次已将前端项目升级为支持 Recharts 图表库，实现美观、交互式的数据可视化。

## 🚀 快速开始

### 1. 安装依赖
```bash
cd frontend
npm install
# 或
yarn install
```

### 2. 启动项目
```bash
npm start
# 或
yarn start
```

## 📝 修改清单

### ✅ 已完成的修改

#### 1. **package.json**
- 添加 `recharts` 依赖 `^2.10.3`

```json
"recharts": "^2.10.3"
```

#### 2. **DefectTrendChart.jsx** (改进)
- 将手写的SVG图表替换为Recharts LineChart
- 移除了复杂的SVG计算逻辑
- 添加了自动响应式设计
- 支持交互式Tooltip和Legend
- 改进的视觉效果和动画

**主要改进：**
```jsx
// 之前：手写SVG图表
<svg width={containerWidth} height={height}>
  {/* 复杂的SVG路径绘制 */}
</svg>

// 之后：使用Recharts
<ResponsiveContainer width="100%" height={400}>
  <LineChart data={filteredData}>
    <CartesianGrid strokeDasharray="3 3" />
    <XAxis dataKey="date" />
    <YAxis />
    <Tooltip />
    <Legend />
    <Line dataKey="total" stroke="#ff4d4f" name={t('trendChart.totalLine')} />
    <Line dataKey="closed" stroke="#52c41a" name={t('trendChart.closedLine')} />
  </LineChart>
</ResponsiveContainer>
```

#### 3. **ChartExamples.jsx** (新增)
新增示例组件，展示Recharts支持的所有主要图表类型：

- 📊 **BarChart** - 柱状图
- 📈 **LineChart** - 线图/趋势图
- 📉 **AreaChart** - 面积图/累积图
- 🥧 **PieChart** - 饼图
- 🎯 **RadarChart** - 雷达图

位置：`frontend/src/components/ChartExamples.jsx`

可在路由中引入：
```jsx
import ChartExamples from './components/ChartExamples';

// 在路由配置中添加
<Route path="/chart-examples" element={<ChartExamples />} />
```

#### 4. **RECHARTS_INTEGRATION_GUIDE.md** (文档)
位置：`frontend/RECHARTS_INTEGRATION_GUIDE.md`

详细的集成指南，包含：
- 安装说明
- 所有图表类型的使用方法
- 常用属性说明
- 最佳实践
- 常见注意事项

## 🎨 图表特性

### DefectTrendChart 的新特性

| 特性 | 说明 |
|------|------|
| 响应式设计 | 自动适应容器宽度 |
| Tooltip交互 | 鼠标悬停显示详细信息 |
| Legend图例 | 自动生成和交互 |
| 平滑动画 | 数据更新时有动画效果 |
| 日期过滤 | 支持日期范围选择 |
| 统计信息 | 显示缺陷总数、已解决、激活等统计 |

## 📚 使用示例

### 基础线图
```jsx
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const MyChart = ({ data }) => (
  <ResponsiveContainer width="100%" height={300}>
    <LineChart data={data}>
      <CartesianGrid strokeDasharray="3 3" />
      <XAxis dataKey="name" />
      <YAxis />
      <Tooltip />
      <Legend />
      <Line type="monotone" dataKey="value" stroke="#8884d8" />
    </LineChart>
  </ResponsiveContainer>
);
```

### 柱状图
```jsx
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const MyBarChart = ({ data }) => (
  <ResponsiveContainer width="100%" height={300}>
    <BarChart data={data}>
      <CartesianGrid strokeDasharray="3 3" />
      <XAxis dataKey="name" />
      <YAxis />
      <Tooltip />
      <Legend />
      <Bar dataKey="value" fill="#8884d8" />
    </BarChart>
  </ResponsiveContainer>
);
```

### 饼图
```jsx
import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const MyPieChart = ({ data }) => {
  const COLORS = ['#0088FE', '#00C49F', '#FFBB28'];
  
  return (
    <ResponsiveContainer width="100%" height={300}>
      <PieChart>
        <Pie data={data} dataKey="value" cx="50%" cy="50%" outerRadius={100}>
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip />
        <Legend />
      </PieChart>
    </ResponsiveContainer>
  );
};
```

## 🔧 开发建议

### 1. 性能优化
```jsx
import { useMemo } from 'react';

// 缓存图表数据，避免不必要的重新计算
const chartData = useMemo(() => {
  return processData(rawData);
}, [rawData]);
```

### 2. 导出图表
```bash
npm install html2canvas
```

```jsx
import html2canvas from 'html2canvas';

const exportChart = async (ref) => {
  const canvas = await html2canvas(ref.current);
  const image = canvas.toDataURL('image/png');
  // 下载或分享
};
```

### 3. 自定义样式
创建主题配置文件：
```jsx
// theme/chartTheme.js
export const chartTheme = {
  colors: {
    primary: '#8884d8',
    success: '#52c41a',
    warning: '#faad14',
    error: '#ff4d4f'
  },
  fonts: {
    size: 12,
    family: 'Arial'
  }
};
```

## ⚠️ 注意事项

1. **数据格式**：必须是数组，元素为对象
   ```jsx
   // ✅ 正确
   const data = [
     { date: '2025-01-01', value: 100 },
     { date: '2025-01-02', value: 200 }
   ];
   
   // ❌ 错误
   const data = [[2025, 1, 1, 100], [2025, 1, 2, 200]];
   ```

2. **ResponsiveContainer高度**：必须设置固定值
   ```jsx
   // ✅ 正确
   <ResponsiveContainer width="100%" height={300}>
   
   // ❌ 错误（会导致不显示）
   <ResponsiveContainer width="100%" height="100%">
   ```

3. **dataKey匹配**：确保dataKey与数据对象字段一致
   ```jsx
   // 数据
   const data = [{ date: '2025-01-01', total: 10 }];
   
   // 图表
   <Line dataKey="total" />  // ✅ 正确
   <Line dataKey="count" />  // ❌ 错误，无法显示
   ```

## 📖 更多资源

- [Recharts官方文档](https://recharts.org/)
- [示例库](https://recharts.org/en-US/examples)
- [API参考](https://recharts.org/en-US/api)

## 🤝 反馈和支持

如遇到问题，请参考：
1. 完整的集成指南：`frontend/RECHARTS_INTEGRATION_GUIDE.md`
2. 示例组件：`frontend/src/components/ChartExamples.jsx`
3. 改进的缺陷图表：`frontend/src/components/DefectTrendChart.jsx`

