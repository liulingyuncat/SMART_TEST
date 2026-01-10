import React, { useState } from 'react';
import { Card, Row, Col, Tabs } from 'antd';
import { useTranslation } from 'react-i18next';
import {
  BarChart,
  Bar,
  LineChart,
  Line,
  AreaChart,
  Area,
  PieChart,
  Pie,
  Cell,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Radar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';

/**
 * Recharts 图表示例组件
 * 展示各种图表类型的使用方法
 */
const ChartExamples = () => {
  const { t } = useTranslation();

  // 柱状图数据
  const barData = [
    { name: 'Jan', value: 400, uv: 240 },
    { name: 'Feb', value: 300, uv: 139 },
    { name: 'Mar', value: 200, uv: 221 },
    { name: 'Apr', value: 278, uv: 229 },
    { name: 'May', value: 189, uv: 200 },
    { name: 'Jun', value: 239, uv: 200 }
  ];

  // 线图数据
  const lineData = [
    { name: 'Week 1', value: 400 },
    { name: 'Week 2', value: 300 },
    { name: 'Week 3', value: 200 },
    { name: 'Week 4', value: 278 },
    { name: 'Week 5', value: 189 }
  ];

  // 面积图数据
  const areaData = [
    { name: 'Jan', sales: 400, profit: 240 },
    { name: 'Feb', sales: 300, profit: 139 },
    { name: 'Mar', sales: 200, profit: 221 },
    { name: 'Apr', sales: 278, profit: 229 },
    { name: 'May', sales: 189, profit: 200 }
  ];

  // 饼图数据
  const pieData = [
    { name: 'Group A', value: 400 },
    { name: 'Group B', value: 300 },
    { name: 'Group C', value: 200 },
    { name: 'Group D', value: 200 }
  ];

  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

  // 雷达图数据
  const radarData = [
    { subject: 'Math', A: 120, B: 110, fullMark: 150 },
    { subject: 'Chinese', A: 98, B: 130, fullMark: 150 },
    { subject: 'English', A: 86, B: 130, fullMark: 150 },
    { subject: 'Geography', A: 99, B: 100, fullMark: 150 },
    { subject: 'Physics', A: 85, B: 90, fullMark: 150 },
    { subject: 'History', A: 65, B: 85, fullMark: 150 }
  ];

  const tabItems = [
    {
      key: 'bar',
      label: 'BarChart (柱状图)',
      children: (
        <Card>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={barData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="value" fill="#8884d8" />
              <Bar dataKey="uv" fill="#82ca9d" />
            </BarChart>
          </ResponsiveContainer>
        </Card>
      )
    },
    {
      key: 'line',
      label: 'LineChart (线图)',
      children: (
        <Card>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={lineData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Line type="monotone" dataKey="value" stroke="#8884d8" />
            </LineChart>
          </ResponsiveContainer>
        </Card>
      )
    },
    {
      key: 'area',
      label: 'AreaChart (面积图)',
      children: (
        <Card>
          <ResponsiveContainer width="100%" height={300}>
            <AreaChart data={areaData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Area type="monotone" dataKey="sales" stackId="1" stroke="#8884d8" fill="#8884d8" />
              <Area type="monotone" dataKey="profit" stackId="1" stroke="#82ca9d" fill="#82ca9d" />
            </AreaChart>
          </ResponsiveContainer>
        </Card>
      )
    },
    {
      key: 'pie',
      label: 'PieChart (饼图)',
      children: (
        <Card>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={pieData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, value }) => `${name}: ${value}`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {pieData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </Card>
      )
    },
    {
      key: 'radar',
      label: 'RadarChart (雷达图)',
      children: (
        <Card>
          <ResponsiveContainer width="100%" height={400}>
            <RadarChart data={radarData}>
              <PolarGrid />
              <PolarAngleAxis dataKey="subject" />
              <PolarRadiusAxis angle={90} domain={[0, 150]} />
              <Radar name="Mike" dataKey="A" stroke="#8884d8" fill="#8884d8" fillOpacity={0.6} />
              <Radar name="Lily" dataKey="B" stroke="#82ca9d" fill="#82ca9d" fillOpacity={0.6} />
              <Legend />
              <Tooltip />
            </RadarChart>
          </ResponsiveContainer>
        </Card>
      )
    }
  ];

  return (
    <Card title="Recharts 图表示例集合" style={{ margin: '20px' }}>
      <Tabs items={tabItems} />
    </Card>
  );
};

export default ChartExamples;
