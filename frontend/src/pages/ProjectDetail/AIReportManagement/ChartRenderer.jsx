import React, { useMemo } from 'react';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
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
import { Empty, Alert } from 'antd';

/**
 * Markdown中的图表渲染组件
 * 支持在Markdown中嵌入Recharts图表
 * 
 * 使用方式：
 * ```chart:line
 * {
 *   "title": "图表标题",
 *   "data": [...],
 *   "dataKey": "value"
 * }
 * ```
 */

const ChartRenderer = ({ type, config }) => {
  // 调试日志
  console.log('[ChartRenderer] Rendering chart:', { type, config });

  // 验证配置
  const isValid = useMemo(() => {
    if (!config || typeof config !== 'object') {
      console.warn('[ChartRenderer] Invalid config - not an object:', config);
      return false;
    }
    if (!config.data || !Array.isArray(config.data)) {
      console.warn('[ChartRenderer] Invalid config - no data array:', config);
      return false;
    }
    if (config.data.length === 0) {
      console.warn('[ChartRenderer] Invalid config - empty data array');
      return false;
    }
    console.log('[ChartRenderer] Config is valid, data length:', config.data.length);
    return true;
  }, [config]);

  if (!isValid) {
    return (
      <Alert
        type="error"
        message="图表配置错误"
        description="请确保配置包含有效的data数组"
        style={{ margin: '16px 0' }}
      />
    );
  }

  const { data, title, dataKey = 'value', colors = ['#8884d8', '#82ca9d', '#ffc658', '#ff7c7c'] } = config;

  try {
    switch (type) {
      case 'line':
        return (
          <div style={{ width: '100%', margin: '16px 0' }}>
            {title && <h4 style={{ marginBottom: '16px' }}>{title}</h4>}
            <ResponsiveContainer width="100%" height={360}>
              <LineChart data={data} margin={{ top: 5, right: 30, left: 0, bottom: 90 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey={Object.keys(data[0])[0]} />
                <YAxis />
                <Tooltip 
                  contentStyle={{
                    backgroundColor: '#fff',
                    border: '1px solid #ccc',
                    borderRadius: '4px',
                    padding: '8px'
                  }}
                />
                <Legend 
                  verticalAlign="bottom" 
                  height={60}
                  wrapperStyle={{ 
                    paddingTop: '40px',
                    paddingLeft: '20px',
                    paddingRight: '20px',
                    width: '100%',
                    boxSizing: 'border-box'
                  }}
                  formatter={(value) => (
                    <span style={{ color: '#333', fontSize: '14px', fontWeight: '500' }}>
                      {value}
                    </span>
                  )}
                />
                {Array.isArray(dataKey) ? (
                  dataKey.map((key, idx) => (
                    <Line
                      key={key}
                      type="monotone"
                      name={key}
                      dataKey={key}
                      stroke={colors[idx % colors.length]}
                      dot={{ r: 4 }}
                      isAnimationActive={true}
                    />
                  ))
                ) : (
                  <Line type="monotone" name={dataKey} dataKey={dataKey} stroke={colors[0]} dot={{ r: 4 }} />
                )}
              </LineChart>
            </ResponsiveContainer>
          </div>
        );

      case 'bar':
        return (
          <div style={{ width: '100%', margin: '16px 0' }}>
            {title && <h4 style={{ marginBottom: '16px' }}>{title}</h4>}
            <ResponsiveContainer width="100%" height={360}>
              <BarChart data={data} margin={{ top: 5, right: 30, left: 0, bottom: 90 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey={Object.keys(data[0])[0]} />
                <YAxis />
                <Tooltip 
                  contentStyle={{
                    backgroundColor: '#fff',
                    border: '1px solid #ccc',
                    borderRadius: '4px',
                    padding: '8px'
                  }}
                />
                <Legend 
                  verticalAlign="bottom" 
                  height={70}
                  wrapperStyle={{ 
                    paddingTop: '50px',
                    paddingLeft: '20px',
                    paddingRight: '20px',
                    width: '100%',
                    boxSizing: 'border-box',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    gap: '20px'
                  }}
                  formatter={(value, entry) => {
                    return (
                      <span 
                        style={{ 
                          color: '#333', 
                          fontSize: '14px', 
                          fontWeight: '500',
                          display: 'flex',
                          alignItems: 'center',
                          gap: '6px',
                          whiteSpace: 'nowrap'
                        }}
                      >
                        <span 
                          style={{
                            width: '10px',
                            height: '10px',
                            backgroundColor: entry?.payload?.fill || '#333',
                            borderRadius: '2px',
                            display: 'inline-block',
                            flexShrink: 0
                          }}
                        />
                        {value}
                      </span>
                    );
                  }}
                />
                {Array.isArray(dataKey) ? (
                  dataKey.map((key, idx) => (
                    <Bar key={key} name={key} dataKey={key} fill={colors[idx % colors.length]} />
                  ))
                ) : (
                  <Bar name={dataKey} dataKey={dataKey} fill={colors[0]} />
                )}
              </BarChart>
            </ResponsiveContainer>
          </div>
        );

      case 'area':
        return (
          <div style={{ width: '100%', margin: '16px 0' }}>
            {title && <h4 style={{ marginBottom: '16px' }}>{title}</h4>}
            <ResponsiveContainer width="100%" height={360}>
              <AreaChart data={data} margin={{ top: 5, right: 30, left: 0, bottom: 90 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey={Object.keys(data[0])[0]} />
                <YAxis />
                <Tooltip 
                  contentStyle={{
                    backgroundColor: '#fff',
                    border: '1px solid #ccc',
                    borderRadius: '4px',
                    padding: '8px'
                  }}
                />
                <Legend 
                  verticalAlign="bottom" 
                  height={60}
                  wrapperStyle={{ 
                    paddingTop: '40px',
                    paddingLeft: '20px',
                    paddingRight: '20px',
                    width: '100%',
                    boxSizing: 'border-box'
                  }}
                  formatter={(value) => (
                    <span style={{ color: '#333', fontSize: '14px', fontWeight: '500' }}>
                      {value}
                    </span>
                  )}
                />
                {Array.isArray(dataKey) ? (
                  dataKey.map((key, idx) => (
                    <Area
                      key={key}
                      type="monotone"
                      name={key}
                      dataKey={key}
                      stackId="1"
                      stroke={colors[idx % colors.length]}
                      fill={colors[idx % colors.length]}
                    />
                  ))
                ) : (
                  <Area type="monotone" name={dataKey} dataKey={dataKey} stroke={colors[0]} fill={colors[0]} />
                )}
              </AreaChart>
            </ResponsiveContainer>
          </div>
        );

      case 'pie':
        return (
          <div style={{ width: '100%', margin: '16px 0' }}>
            {title && <h4 style={{ marginBottom: '16px', textAlign: 'center' }}>{title}</h4>}
            <ResponsiveContainer width="100%" height={480}>
              <PieChart>
                <Pie
                  data={data}
                  cx="50%"
                  cy="42%"
                  labelLine={true}
                  label={({ name, value, percent }) => {
                    // 计算总值用于百分比计算
                    const total = data.reduce((sum, item) => sum + (item[dataKey] || item.value || 0), 0);
                    const percentage = ((value / total) * 100).toFixed(1);
                    return `${name}: ${value} (${percentage}%)`;
                  }}
                  outerRadius={100}
                  innerRadius={0}
                  fill="#8884d8"
                  dataKey={dataKey || 'value'}
                  startAngle={90}
                  endAngle={-270}
                >
                  {data.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={colors[index % colors.length]} />
                  ))}
                </Pie>
                <Tooltip 
                  formatter={(value, name, props) => {
                    const total = data.reduce((sum, item) => sum + (item[dataKey] || item.value || 0), 0);
                    const percentage = ((value / total) * 100).toFixed(1);
                    return `${value} (${percentage}%)`;
                  }}
                  contentStyle={{
                    backgroundColor: '#fff',
                    border: '1px solid #ccc',
                    borderRadius: '4px',
                    padding: '8px',
                    fontSize: '12px'
                  }}
                />
                <Legend 
                  verticalAlign="bottom" 
                  height={50}
                  wrapperStyle={{ 
                    paddingTop: '30px',
                    paddingLeft: '20px',
                    paddingRight: '20px',
                    width: '100%',
                    boxSizing: 'border-box'
                  }}
                  formatter={(value, entry) => {
                    try {
                      const item = data[entry.index];
                      if (!item) {
                        console.warn('[ChartRenderer] Legend formatter - item is undefined at index:', entry.index);
                        return <span style={{ color: '#333', fontSize: '14px' }}>{value}</span>;
                      }
                      const total = data.reduce((sum, d) => sum + (d[dataKey] || d.value || 0), 0);
                      const val = item[dataKey] || item.value || 0;
                      const percentage = total > 0 ? ((val / total) * 100).toFixed(1) : 0;
                      return <span style={{ color: '#333', fontSize: '14px', fontWeight: '500' }}>{item.name || value}: {val} ({percentage}%)</span>;
                    } catch (err) {
                      console.error('[ChartRenderer] Legend formatter error:', err, 'entry:', entry);
                      return <span style={{ color: '#333', fontSize: '14px' }}>{value}</span>;
                    }
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
          </div>
        );

      case 'radar':
        return (
          <div style={{ width: '100%', margin: '16px 0' }}>
            {title && <h4 style={{ marginBottom: '16px' }}>{title}</h4>}
            <ResponsiveContainer width="100%" height={400}>
              <RadarChart data={data}>
                <PolarGrid />
                <PolarAngleAxis dataKey="subject" />
                <PolarRadiusAxis angle={90} domain={[0, 150]} />
                {Array.isArray(dataKey) ? (
                  dataKey.map((key, idx) => (
                    <Radar
                      key={key}
                      name={key}
                      dataKey={key}
                      stroke={colors[idx % colors.length]}
                      fill={colors[idx % colors.length]}
                      fillOpacity={0.6}
                    />
                  ))
                ) : (
                  <Radar
                    name={dataKey}
                    dataKey={dataKey}
                    stroke={colors[0]}
                    fill={colors[0]}
                    fillOpacity={0.6}
                  />
                )}
                <Legend />
                <Tooltip />
              </RadarChart>
            </ResponsiveContainer>
          </div>
        );

      default:
        return (
          <Alert
            type="warning"
            message="未知的图表类型"
            description={`不支持的图表类型: ${type}`}
            style={{ margin: '16px 0' }}
          />
        );
    }
  } catch (error) {
    return (
      <Alert
        type="error"
        message="图表渲染错误"
        description={error.message}
        style={{ margin: '16px 0' }}
      />
    );
  }
};

export default ChartRenderer;
