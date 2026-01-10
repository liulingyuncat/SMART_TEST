# React + Recharts 集成指南

本指南说明如何在React项目中使用Recharts库创建美观的图表。

## 一、安装依赖

在项目根目录执行：

```bash
npm install recharts
```

或者如果使用yarn：

```bash
yarn add recharts
```

## 二、基础导入

在组件中导入Recharts的相关组件：

```jsx
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';
```

## 三、项目中的修改

### 1. DefectTrendChart 组件
已将原来的SVG图表替换为Recharts的LineChart组件：

**特点：**
- 使用`ResponsiveContainer`实现自适应宽度和高度
- 支持鼠标悬停显示Tooltip
- 自动生成图例
- 动画效果

**使用方式：**
```jsx
<ResponsiveContainer width="100%" height={400}>
  <LineChart data={filteredData}>
    <CartesianGrid strokeDasharray="3 3" />
    <XAxis dataKey="date" />
    <YAxis />
    <Tooltip />
    <Legend />
    <Line 
      type="monotone" 
      dataKey="total" 
      stroke="#ff4d4f" 
      strokeWidth={2}
      dot={{ fill: '#ff4d4f', r: 4 }}
    />
    <Line 
      type="monotone" 
      dataKey="closed" 
      stroke="#52c41a" 
      strokeWidth={2}
      dot={{ fill: '#52c41a', r: 4 }}
    />
  </LineChart>
</ResponsiveContainer>
```

### 2. ChartExamples 组件（新增）
位于 `frontend/src/components/ChartExamples.jsx`

提供了以下图表示例：
- **BarChart**：柱状图
- **LineChart**：线图
- **AreaChart**：面积图
- **PieChart**：饼图
- **RadarChart**：雷达图

可在路由中引入此组件进行展示。

## 四、常用图表类型详解

### 1. LineChart（线图/趋势图）
适用于：时间序列数据、趋势分析

```jsx
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
```

### 2. BarChart（柱状图）
适用于：分类数据对比、多维度比较

```jsx
<ResponsiveContainer width="100%" height={300}>
  <BarChart data={data}>
    <CartesianGrid strokeDasharray="3 3" />
    <XAxis dataKey="name" />
    <YAxis />
    <Tooltip />
    <Legend />
    <Bar dataKey="value" fill="#8884d8" />
    <Bar dataKey="uv" fill="#82ca9d" />
  </BarChart>
</ResponsiveContainer>
```

### 3. AreaChart（面积图）
适用于：累积数据、份额变化

```jsx
<ResponsiveContainer width="100%" height={300}>
  <AreaChart data={data}>
    <CartesianGrid strokeDasharray="3 3" />
    <XAxis dataKey="name" />
    <YAxis />
    <Tooltip />
    <Legend />
    <Area type="monotone" dataKey="value" fill="#8884d8" />
  </AreaChart>
</ResponsiveContainer>
```

### 4. PieChart（饼图）
适用于：比例分布、占比展示

```jsx
<ResponsiveContainer width="100%" height={300}>
  <PieChart>
    <Pie
      data={data}
      cx="50%"
      cy="50%"
      outerRadius={80}
      fill="#8884d8"
      dataKey="value"
    >
      {data.map((entry, index) => (
        <Cell key={`cell-${index}`} fill={COLORS[index]} />
      ))}
    </Pie>
    <Tooltip />
    <Legend />
  </PieChart>
</ResponsiveContainer>
```

### 5. RadarChart（雷达图）
适用于：多维度综合对比、性能评估

```jsx
<ResponsiveContainer width="100%" height={400}>
  <RadarChart data={data}>
    <PolarGrid />
    <PolarAngleAxis dataKey="subject" />
    <PolarRadiusAxis angle={90} domain={[0, 150]} />
    <Radar name="A" dataKey="value" stroke="#8884d8" fill="#8884d8" fillOpacity={0.6} />
    <Legend />
    <Tooltip />
  </RadarChart>
</ResponsiveContainer>
```

## 五、关键属性说明

### ResponsiveContainer
- `width`：宽度（"100%"表示占满父容器）
- `height`：固定高度（像素值）

### LineChart / BarChart / AreaChart
- `data`：图表数据（数组）
- `margin`：边距配置
- `layout`："vertical"设置为竖向

### Line / Bar / Area
- `type`：线条类型（"monotone"、"linear"等）
- `dataKey`：对应数据对象的字段名
- `stroke`：颜色
- `strokeWidth`：线宽
- `dot`：数据点样式
- `fill`：填充颜色

### XAxis / YAxis
- `dataKey`：对应数据对象的字段名（仅XAxis）
- `angle`：标签旋转角度
- `tick`：标签样式配置

### Tooltip
- `contentStyle`：悬停框样式
- `formatter`：值格式化函数
- `labelFormatter`：标签格式化函数

### Legend
- `wrapperStyle`：图例包装样式
- `iconType`：图标类型（"line"、"circle"等）

## 六、常见注意事项

1. **数据格式**：确保传入的数据是数组格式，每个元素都是对象
   ```jsx
   const data = [
     { name: 'Jan', value: 400 },
     { name: 'Feb', value: 300 }
   ];
   ```

2. **响应式设计**：始终使用`ResponsiveContainer`包裹图表，并设置固定高度
   ```jsx
   <ResponsiveContainer width="100%" height={300}>
     {/* 图表组件 */}
   </ResponsiveContainer>
   ```

3. **性能优化**：
   - 避免频繁重新渲染，使用`useMemo`缓存数据
   - 大数据集可考虑分页或虚拟化
   - 使用`isAnimationActive={false}`禁用动画以提升性能

4. **颜色配置**：定义颜色常量便于维护
   ```jsx
   const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];
   ```

5. **国际化**：图表标签应使用i18n翻译
   ```jsx
   <XAxis label={{ value: t('chart.xAxisLabel') }} />
   ```

## 七、实际使用示例

### 示例1：缺陷趋势统计
```jsx
const filteredData = [
  { date: '2025-01-01', total: 10, closed: 2, open: 8 },
  { date: '2025-01-02', total: 12, closed: 3, open: 9 }
];

<ResponsiveContainer width="100%" height={400}>
  <LineChart data={filteredData}>
    <CartesianGrid strokeDasharray="3 3" />
    <XAxis dataKey="date" />
    <YAxis />
    <Tooltip />
    <Legend />
    <Line dataKey="total" stroke="#ff4d4f" name="总数" />
    <Line dataKey="closed" stroke="#52c41a" name="已关闭" />
    <Line dataKey="open" stroke="#faad14" name="未关闭" />
  </LineChart>
</ResponsiveContainer>
```

### 示例2：按优先级分布
```jsx
const priorityData = [
  { name: 'Critical', value: 20 },
  { name: 'High', value: 35 },
  { name: 'Medium', value: 30 },
  { name: 'Low', value: 15 }
];

<ResponsiveContainer width="100%" height={300}>
  <PieChart>
    <Pie
      data={priorityData}
      cx="50%"
      cy="50%"
      outerRadius={80}
      dataKey="value"
    >
      {priorityData.map((entry, index) => (
        <Cell key={`cell-${index}`} fill={COLORS[index]} />
      ))}
    </Pie>
    <Legend />
    <Tooltip />
  </PieChart>
</ResponsiveContainer>
```

## 八、更多资源

- 官方文档：https://recharts.org/
- API参考：https://recharts.org/en-US/api
- 示例库：https://recharts.org/en-US/examples

## 九、后续优化建议

1. **自定义样式**：创建一个统一的主题配置文件
2. **性能监测**：使用React DevTools Profiler检查渲染性能
3. **导出功能**：集成html2canvas库实现图表导出为图片
4. **高级交互**：使用Recharts的回调函数实现自定义交互
5. **3D图表**：如需3D效果，可结合Three.js使用，但会增加复杂度

