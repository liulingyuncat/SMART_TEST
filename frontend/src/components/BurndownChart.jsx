import React, { useState, useEffect, useMemo } from 'react';
import { Card, DatePicker, Spin, message, Empty, Row, Col, Statistic } from 'antd';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';
import isBetween from 'dayjs/plugin/isBetween';
import { fetchDefects } from '../api/defect';
import { generateSmoothPath, chartConfig, getOptimalGridLineCount } from '../utils/chartUtils';

dayjs.extend(isBetween);

/**
 * 燃尽图组件
 * 功能：
 * 1. 展示 Bug 总数和关闭数量趋势
 * 2. 支持日期范围选择（拖拽选择）
 * 3. 显示收敛趋势
 */

const BurndownChart = ({ projectId }) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [chartData, setChartData] = useState([]);
  const [dateRange, setDateRange] = useState([
    dayjs().subtract(30, 'day'),
    dayjs()
  ]);
  const [defectStats, setDefectStats] = useState(null);

  // 获取真实缺陷数据并生成图表数据
  useEffect(() => {
    const loadDefects = async () => {
      if (!projectId) return;
      
      try {
        setLoading(true);
        // 获取所有缺陷数据（不分页）
        const response = await fetchDefects(projectId, { page: 1, page_size: 1000 });
        const defects = response.defects || [];
        
        console.log('[BurndownChart] Loaded defects:', defects.length);
        
        // 生成燃尽图数据
        const data = generateBurndownData(defects);
        setChartData(data);
      } catch (error) {
        console.error('Failed to load defects:', error);
        message.error(t('burndownChart.loadFailed'));
      } finally {
        setLoading(false);
      }
    };

    loadDefects();
  }, [projectId, t]);

  // 根据真实缺陷数据生成燃尽图数据
  const generateBurndownData = (defects) => {
    if (!defects || defects.length === 0) return [];
    
    // 找到最早的创建日期
    const dates = defects.map(d => dayjs(d.created_at));
    // dayjs没有min方法，需要手动找最小日期
    let minDate = dates[0];
    dates.forEach(date => {
      if (date.isBefore(minDate)) {
        minDate = date;
      }
    });
    const maxDate = dayjs();
    
    // 生成日期序列（从最早的Bug创建日期到今天）
    const data = [];
    let currentDate = minDate.startOf('day');
    const endDate = maxDate.startOf('day');
    
    while (currentDate.isBefore(endDate) || currentDate.isSame(endDate, 'day')) {
      // 统计截至当前日期的Bug总数和已关闭数
      let total = 0;
      let closed = 0;
      
      defects.forEach(defect => {
        const createdAt = dayjs(defect.created_at);
        const updatedAt = dayjs(defect.updated_at);
        
        // 如果在当前日期之前或当天创建，计入总数
        if (createdAt.isBefore(currentDate, 'day') || createdAt.isSame(currentDate, 'day')) {
          total++;
          
          // 如果状态是Closed，并且updated_at在当前日期之前或当天，说明在这天或之前已关闭
          if (defect.status === 'Closed' && 
              (updatedAt.isBefore(currentDate, 'day') || updatedAt.isSame(currentDate, 'day'))) {
            closed++;
          }
        }
      });
      
      data.push({
        date: currentDate.format('YYYY-MM-DD'),
        dateObj: currentDate,
        total,
        closed,
        open: total - closed
      });
      
      currentDate = currentDate.add(1, 'day');
    }
    
    return data;
  };

  // 过滤日期范围内的数据
  const filteredData = useMemo(() => {
    if (!dateRange || dateRange.length !== 2) return [];
    
    const [start, end] = dateRange;
    
    return chartData.filter(item => {
      const itemDate = dayjs(item.date);
      return itemDate.isBetween(start, end, 'day', '[]');
    });
  }, [chartData, dateRange]);

  // 计算统计数据
  useEffect(() => {
    if (filteredData.length === 0) {
      setDefectStats(null);
      return;
    }

    const lastDay = filteredData[filteredData.length - 1];
    const stats = {
      currentTotal: lastDay?.total || 0,
      totalClosed: lastDay?.closed || 0,
      openCount: (lastDay?.total || 0) - (lastDay?.closed || 0)
    };

    setDefectStats(stats);
  }, [filteredData]);

  // 如果没有 @ant-design/charts，使用优化的 SVG 图表
  const renderChart = () => {
    if (filteredData.length === 0) {
      return <Empty description={t('burndownChart.noData')} />;
    }

    const width = 1200;
    const height = 400;
    const padding = { top: 40, right: 40, bottom: 60, left: 60 };
    const chartWidth = width - padding.left - padding.right;
    const chartHeight = height - padding.top - padding.bottom;

    // 计算坐标
    const maxValue = Math.max(...filteredData.map(d => d.total)) || 1;
    const xStep = chartWidth / (filteredData.length - 1 || 1);
    const yScale = chartHeight / maxValue;

    // 生成数据点
    const totalPoints = filteredData.map((d, i) => ({
      x: padding.left + i * xStep,
      y: padding.top + chartHeight - d.total * yScale
    }));
    
    const closedPoints = filteredData.map((d, i) => ({
      x: padding.left + i * xStep,
      y: padding.top + chartHeight - d.closed * yScale
    }));

    // 使用平滑曲线路径
    const totalPath = generateSmoothPath(totalPoints);
    const closedPath = generateSmoothPath(closedPoints);

    // 计算最优网格线数量
    const gridLines = getOptimalGridLineCount(maxValue);

    return (
      <svg width={width} height={height} style={{ 
        border: '1px solid #f0f0f0', 
        borderRadius: 4,
        backgroundColor: '#fafafa'
      }}>
        <defs>
          {/* 渐变定义 */}
          <linearGradient id="totalGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor="#ff4d4f" stopOpacity="0.3" />
            <stop offset="100%" stopColor="#ff4d4f" stopOpacity="0.01" />
          </linearGradient>
          <linearGradient id="closedGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor="#52c41a" stopOpacity="0.3" />
            <stop offset="100%" stopColor="#52c41a" stopOpacity="0.01" />
          </linearGradient>
          {/* 阴影定义 */}
          <filter id="shadow" x="-50%" y="-50%" width="200%" height="200%">
            <feDropShadow dx="0" dy="2" stdDeviation="3" floodOpacity="0.15" />
          </filter>
        </defs>

        {/* Y 轴网格线 */}
        {Array.from({ length: gridLines + 1 }).map((_, i) => {
          const value = Math.floor((maxValue / gridLines) * i);
          const y = padding.top + chartHeight - (value * yScale);
          return (
            <line
              key={`grid-${i}`}
              x1={padding.left}
              y1={y}
              x2={width - padding.right}
              y2={y}
              stroke="#e8e8e8"
              strokeDasharray="4"
              strokeWidth="0.5"
            />
          );
        })}

        {/* Y 轴 */}
        <line
          x1={padding.left}
          y1={padding.top}
          x2={padding.left}
          y2={padding.top + chartHeight}
          stroke={chartConfig.defaults.axisColor}
          strokeWidth="1.5"
        />
        
        {/* X 轴 */}
        <line
          x1={padding.left}
          y1={padding.top + chartHeight}
          x2={width - padding.right}
          y2={padding.top + chartHeight}
          stroke={chartConfig.defaults.axisColor}
          strokeWidth="1.5"
        />

        {/* Y 轴标签和刻度 */}
        {Array.from({ length: gridLines + 1 }).map((_, i) => {
          const value = Math.floor((maxValue / gridLines) * i);
          const y = padding.top + chartHeight - (value * yScale);
          return (
            <g key={`y-${i}`}>
              <line
                x1={padding.left - 5}
                y1={y}
                x2={padding.left}
                y2={y}
                stroke={chartConfig.defaults.axisColor}
                strokeWidth="1"
              />
              <text
                x={padding.left - 10}
                y={y + 4}
                textAnchor="end"
                fontSize="12"
                fontFamily="system-ui"
                fill={chartConfig.defaults.labelColor}
              >
                {value}
              </text>
            </g>
          );
        })}

        {/* X 轴标签和刻度 */}
        {filteredData.map((d, i) => {
          if (filteredData.length > 15 && i % Math.ceil(filteredData.length / 8) !== 0) return null;
          const x = padding.left + i * xStep;
          return (
            <g key={`x-${i}`}>
              <line
                x1={x}
                y1={padding.top + chartHeight}
                x2={x}
                y2={padding.top + chartHeight + 5}
                stroke={chartConfig.defaults.axisColor}
                strokeWidth="1"
              />
              <text
                x={x}
                y={padding.top + chartHeight + 20}
                textAnchor="middle"
                fontSize="12"
                fontFamily="system-ui"
                fill={chartConfig.defaults.labelColor}
              >
                {d.date}
              </text>
            </g>
          );
        })}

        {/* Bug 总数曲线 - 带渐变填充 */}
        <path
          d={`${totalPath} L ${padding.left + (filteredData.length - 1) * xStep} ${padding.top + chartHeight} L ${padding.left} ${padding.top + chartHeight} Z`}
          fill="url(#totalGradient)"
          opacity="0.6"
        />
        
        {/* Bug 总数曲线 */}
        <path
          d={totalPath}
          fill="none"
          stroke={chartConfig.colors.red}
          strokeWidth={chartConfig.defaults.strokeWidth}
          strokeLinecap="round"
          strokeLinejoin="round"
          filter="url(#shadow)"
        />

        {/* Bug 已关闭曲线 - 带渐变填充 */}
        <path
          d={`${closedPath} L ${padding.left + (filteredData.length - 1) * xStep} ${padding.top + chartHeight} L ${padding.left} ${padding.top + chartHeight} Z`}
          fill="url(#closedGradient)"
          opacity="0.6"
        />
        
        {/* Bug 已关闭曲线 */}
        <path
          d={closedPath}
          fill="none"
          stroke={chartConfig.colors.green}
          strokeWidth={chartConfig.defaults.strokeWidth}
          strokeLinecap="round"
          strokeLinejoin="round"
          filter="url(#shadow)"
        />

        {/* 数据点 */}
        {totalPoints.map((point, i) => (
          <g key={`points-${i}`}>
            {/* 总数点 */}
            <circle
              cx={point.x}
              cy={point.y}
              r={chartConfig.defaults.pointRadius}
              fill={chartConfig.colors.red}
              opacity="0.8"
              style={{ cursor: 'pointer', transition: 'all 0.2s ease' }}
            />
            {/* 已关闭点 */}
            <circle
              cx={closedPoints[i].x}
              cy={closedPoints[i].y}
              r={chartConfig.defaults.pointRadius}
              fill={chartConfig.colors.green}
              opacity="0.8"
              style={{ cursor: 'pointer', transition: 'all 0.2s ease' }}
            />
          </g>
        ))}

        {/* 图例 */}
        <g>
          <rect x={40} y={10} width={160} height={28} fill="white" stroke="#e0e0e0" strokeWidth="0.5" rx="3" opacity="0.95"/>
          
          <line x1={50} y1={20} x2={80} y2={20} stroke={chartConfig.colors.red} strokeWidth={chartConfig.defaults.strokeWidth} />
          <text x={90} y={25} fontSize="12" fontFamily="system-ui" fontWeight="500" fill={chartConfig.defaults.labelColor}>
            {t('burndownChart.totalDefects')}
          </text>

          <line x1={50} y1={32} x2={80} y2={32} stroke={chartConfig.colors.green} strokeWidth={chartConfig.defaults.strokeWidth} />
          <text x={90} y={37} fontSize="12" fontFamily="system-ui" fontWeight="500" fill={chartConfig.defaults.labelColor}>
            {t('burndownChart.closed')}
          </text>
        </g>

        {/* Y 轴标题 */}
        <text
          x={20}
          y={padding.top + chartHeight / 2}
          textAnchor="middle"
          fontSize="13"
          fontFamily="system-ui"
          fontWeight="500"
          fill={chartConfig.defaults.labelColor}
          transform={`rotate(-90, 20, ${padding.top + chartHeight / 2})`}
        >
          {t('burndownChart.defectCount')}
        </text>

        {/* X 轴标题 */}
        <text
          x={width / 2}
          y={height - 10}
          textAnchor="middle"
          fontSize="13"
          fontFamily="system-ui"
          fontWeight="500"
          fill={chartConfig.defaults.labelColor}
        >
          {t('burndownChart.date')}
        </text>
      </svg>
    );
  };

  return (
    <Card
      title={t('burndownChart.title')}
      extra={
        <DatePicker.RangePicker
          value={dateRange}
          onChange={setDateRange}
          style={{ width: 280 }}
        />
      }
      style={{ marginTop: 16 }}
    >
      {loading ? (
        <div style={{ display: 'flex', justifyContent: 'center', padding: 48 }}>
          <Spin />
        </div>
      ) : (
        <>
          {/* 统计指标 */}
          {defectStats && (
            <Row gutter={16} style={{ marginBottom: 24 }}>
              <Col span={8}>
                <Statistic
                  title={t('burndownChart.currentDefects')}
                  value={defectStats.currentTotal}
                  valueStyle={{ color: '#ff4d4f' }}
                />
              </Col>
              <Col span={8}>
                <Statistic
                  title={t('burndownChart.closedDefects')}
                  value={defectStats.totalClosed}
                  valueStyle={{ color: '#52c41a' }}
                />
              </Col>
              <Col span={8}>
                <Statistic
                  title={t('burndownChart.openDefects')}
                  value={defectStats.openCount}
                  valueStyle={{ color: '#faad14' }}
                />
              </Col>
            </Row>
          )}

          {/* 图表 */}
          <div style={{ overflowX: 'auto', marginTop: 16 }}>
            {renderChart()}
          </div>

          {/* 图表说明 */}
          <div style={{ marginTop: 16, padding: 12, background: '#fafafa', borderRadius: 4 }}>
            <p style={{ margin: '8px 0', fontSize: 12, color: '#666' }}>
              {t('burndownChart.description')}
            </p>
          </div>
        </>
      )}
    </Card>
  );
};

export default BurndownChart;
