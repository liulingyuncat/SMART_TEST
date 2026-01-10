import React, { useState, useEffect, useMemo, useRef } from 'react';
import { Card, DatePicker, Spin, message, Empty } from 'antd';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';
import isBetween from 'dayjs/plugin/isBetween';
import { fetchDefects } from '../api/defect';
import { generateLinePath, chartConfig, getOptimalGridLineCount } from '../utils/chartUtils';

dayjs.extend(isBetween);

/**
 * 缺陷趋势图组件
 * 功能：
 * 1. 展示缺陷总数和已关闭数量趋势
 * 2. 支持日期范围选择
 * 3. 横坐标自适应间隔,无横向滚动条
 */

const DefectTrendChart = ({ projectId }) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [chartData, setChartData] = useState([]);
  const [dateRange, setDateRange] = useState([
    dayjs().subtract(30, 'day'),
    dayjs()
  ]);
  const [defectStats, setDefectStats] = useState(null);
  const [allDefects, setAllDefects] = useState([]);
  const [hoveredIndex, setHoveredIndex] = useState(null);
  const chartContainerRef = useRef(null);
  const [containerHeight, setContainerHeight] = useState(500); // 保存原始缺陷数据用于统计

  // 监听容器大小变化
  useEffect(() => {
    const updateHeight = () => {
      if (chartContainerRef.current) {
        // 获取容器高度，动态计算为屏幕可用高度的1/2
        const rect = chartContainerRef.current.getBoundingClientRect();
        const availableHeight = window.innerHeight - rect.top - 60;
        // 高度为1/2，最小250px
        const finalHeight = Math.max(Math.floor(availableHeight / 2), 250);
        setContainerHeight(finalHeight);
      }
    };

    updateHeight();
    window.addEventListener('resize', updateHeight);
    return () => window.removeEventListener('resize', updateHeight);
  }, []);

  // 获取真实缺陷数据并生成图表数据
  useEffect(() => {
    const loadDefects = async () => {
      if (!projectId) return;
      
      try {
        setLoading(true);
        // 获取所有缺陷数据（不分页）
        const response = await fetchDefects(projectId, { page: 1, page_size: 1000 });
        const defects = response.defects || [];
        
        console.log('[DefectTrendChart] Loaded defects:', defects.length);
        
        // 保存原始缺陷数据
        setAllDefects(defects);
        
        // 生成趋势图数据
        const data = generateTrendData(defects);
        setChartData(data);
      } catch (error) {
        console.error('Failed to load defects:', error);
        message.error(t('trendChart.loadFailed'));
      } finally {
        setLoading(false);
      }
    };

    loadDefects();
  }, [projectId, t]);

  // 根据真实缺陷数据生成趋势图数据
  const generateTrendData = (defects) => {
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

  // 计算统计数据（从原始缺陷数据按状态统计）
  useEffect(() => {
    if (allDefects.length === 0) {
      setDefectStats(null);
      return;
    }

    // 按状态统计
    const stats = {
      total: allDefects.length,
      resolved: allDefects.filter(d => d.status === 'Resolved').length,
      active: allDefects.filter(d => d.status === 'Active').length,
      new: allDefects.filter(d => d.status === 'New').length,
      closed: allDefects.filter(d => d.status === 'Closed').length
    };

    setDefectStats(stats);
  }, [allDefects]);

  // 使用 SVG 原生渲染图表
  const renderChart = () => {
    if (filteredData.length === 0) {
      return <Empty description={t('trendChart.noData')} />;
    }

    // 计算SVG尺寸和缩放 - 宽度自适应，高度固定
    // 根据容器获取宽度
    let containerWidth = 1000;
    if (chartContainerRef.current) {
      containerWidth = chartContainerRef.current.offsetWidth;
    }
    
    const svgWidth = Math.max(containerWidth - 40, 800); // 留出左右边距
    const svgHeight = 280; // 固定高度，不使用动态高度
    const padding = { top: 35, right: 50, bottom: 70, left: 70 };
    const chartWidth = svgWidth - padding.left - padding.right;
    const chartHeight = svgHeight - padding.top - padding.bottom;
    
    // 找到数据的最大值用于缩放
    const maxValue = Math.max(...filteredData.map(d => d.total || 0), 1);
    const yScale = chartHeight / maxValue;
    const xStep = chartWidth / (filteredData.length > 1 ? filteredData.length - 1 : 1);
    
    // 计算坐标
    const totalLine = filteredData.map((d, i) => ({
      x: padding.left + i * xStep,
      y: padding.top + chartHeight - (d.total * yScale),
      data: d
    }));
    
    const closedLine = filteredData.map((d, i) => ({
      x: padding.left + i * xStep,
      y: padding.top + chartHeight - (d.closed * yScale),
      data: d
    }));

    // 动态计算显示日期标签的间隔（保证标签之间至少间隔50px）
    const minLabelSpacing = 50;
    const labelInterval = Math.max(1, Math.ceil((filteredData.length * minLabelSpacing) / chartWidth));
    
    // 使用 getOptimalGridLineCount 计算最优网格线数量
    const gridLines = getOptimalGridLineCount(maxValue);

    return (
      <div style={{ position: 'relative', width: '100%', height: '100%' }}>
        <svg 
          viewBox={`0 0 ${svgWidth} ${svgHeight}`} 
          style={{ width: '100%', height: '100%', display: 'block' }}
          preserveAspectRatio="xMidYMid meet"
        >
          {/* 背景 */}
          <rect width={svgWidth} height={svgHeight} fill="#fafafa"/>

          <defs>
            {/* 渐变定义 */}
            <linearGradient id="totalGradient2" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stopColor="#ff4d4f" stopOpacity="0.3" />
              <stop offset="100%" stopColor="#ff4d4f" stopOpacity="0.01" />
            </linearGradient>
            <linearGradient id="closedGradient2" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stopColor="#52c41a" stopOpacity="0.3" />
              <stop offset="100%" stopColor="#52c41a" stopOpacity="0.01" />
            </linearGradient>
            {/* 阴影定义 */}
            <filter id="shadow2" x="-50%" y="-50%" width="200%" height="200%">
              <feDropShadow dx="0" dy="2" stdDeviation="2.5" floodOpacity="0.12" />
            </filter>
          </defs>
          
          {/* Y轴网格线 */}
          {Array.from({ length: gridLines + 1 }, (_, i) => {
            const y = padding.top + chartHeight - (i * (chartHeight / gridLines));
            const value = Math.round((i * maxValue) / gridLines);
            return (
              <g key={`grid-${i}`}>
                <line 
                  x1={padding.left} 
                  y1={y} 
                  x2={svgWidth - padding.right} 
                  y2={y} 
                  stroke="#efefef" 
                  strokeWidth="0.5"
                  strokeDasharray="3,3"
                />
                <text 
                  x={padding.left - 10} 
                  y={y + 4} 
                  fontSize="11" 
                  textAnchor="end" 
                  fill="#999"
                  fontFamily="system-ui"
                  fontWeight="400"
                >
                  {value}
                </text>
              </g>
            );
          })}
          
          {/* X轴和Y轴 */}
          <line 
            x1={padding.left} 
            y1={padding.top} 
            x2={padding.left} 
            y2={svgHeight - padding.bottom} 
            stroke={chartConfig.defaults.axisColor} 
            strokeWidth="1"
          />
          <line 
            x1={padding.left} 
            y1={svgHeight - padding.bottom} 
            x2={svgWidth - padding.right} 
            y2={svgHeight - padding.bottom} 
            stroke={chartConfig.defaults.axisColor} 
            strokeWidth="1"
          />
          
          {/* 总缺陷数曲线填充 */}
          <path
            d={`${generateLinePath(totalLine)} L ${padding.left + (filteredData.length - 1) * xStep} ${svgHeight - padding.bottom} L ${padding.left} ${svgHeight - padding.bottom} Z`}
            fill="url(#totalGradient2)"
            opacity="0.5"
          />
          
          {/* 总缺陷数曲线（红色） */}
          <path
            d={generateLinePath(totalLine)}
            fill="none"
            stroke={chartConfig.colors.red}
            strokeWidth={chartConfig.defaults.strokeWidth}
            strokeLinecap="round"
            strokeLinejoin="round"
            filter="url(#shadow2)"
          />
          
          {/* 已关闭缺陷数曲线填充 */}
          <path
            d={`${generateLinePath(closedLine)} L ${padding.left + (filteredData.length - 1) * xStep} ${svgHeight - padding.bottom} L ${padding.left} ${svgHeight - padding.bottom} Z`}
            fill="url(#closedGradient2)"
            opacity="0.5"
          />
          
          {/* 已关闭缺陷数曲线（绿色） */}
          <path
            d={generateLinePath(closedLine)}
            fill="none"
            stroke={chartConfig.colors.green}
            strokeWidth={chartConfig.defaults.strokeWidth}
            strokeLinecap="round"
            strokeLinejoin="round"
            filter="url(#shadow2)"
          />
          
          {/* 总缺陷数数据点 */}
          {totalLine.map((point, i) => (
            <circle
              key={`total-point-${i}`}
              cx={point.x}
              cy={point.y}
              r={hoveredIndex === i ? 4 : 2.5}
              fill={chartConfig.colors.red}
              cursor="pointer"
              onMouseEnter={() => setHoveredIndex(i)}
              onMouseLeave={() => setHoveredIndex(null)}
              style={{ transition: 'r 0.15s ease', pointerEvents: 'none', opacity: 0.85 }}
            />
          ))}
          
          {/* 已关闭缺陷数数据点 */}
          {closedLine.map((point, i) => (
            <circle
              key={`closed-point-${i}`}
              cx={point.x}
              cy={point.y}
              r={hoveredIndex === i ? 4 : 2.5}
              fill={chartConfig.colors.green}
              cursor="pointer"
              onMouseEnter={() => setHoveredIndex(i)}
              onMouseLeave={() => setHoveredIndex(null)}
              style={{ transition: 'r 0.15s ease', pointerEvents: 'none', opacity: 0.85 }}
            />
          ))}
          {/* X轴标签 */}
          {filteredData.map((d, i) => {
            const showLabel = i % labelInterval === 0 || i === filteredData.length - 1;
            if (!showLabel) return null;
            
            const dateStr = dayjs(d.date).format('MM-DD');
            const xPos = padding.left + i * xStep;
            
            return (
              <g key={`x-label-${i}`}>
                <line 
                  x1={xPos} 
                  y1={svgHeight - padding.bottom} 
                  x2={xPos} 
                  y2={svgHeight - padding.bottom + 3} 
                  stroke="#999" 
                  strokeWidth="0.5"
                />
                <text 
                  x={xPos} 
                  y={svgHeight - padding.bottom + 18} 
                  fontSize="11" 
                  textAnchor="middle" 
                  fill="#999"
                  fontFamily="system-ui"
                >
                  {dateStr}
                </text>
              </g>
            );
          })}
          
          {/* Y轴标签 */}
          <text
            x={20}
            y={svgHeight / 2}
            fontSize="12"
            fill={chartConfig.defaults.labelColor}
            textAnchor="middle"
            fontFamily="system-ui"
            fontWeight="500"
            transform={`rotate(-90, 20, ${svgHeight / 2})`}
          >
            {t('trendChart.defectCount')}
          </text>
          
          {/* 图例 */}
          <g>
            <rect x={svgWidth - 140} y={10} width={120} height={42} fill="white" stroke="#e8e8e8" strokeWidth="0.5" rx="4" opacity="0.97"/>
            
            <line x1={svgWidth - 130} y1={20} x2={svgWidth - 110} y2={20} stroke={chartConfig.colors.red} strokeWidth={chartConfig.defaults.strokeWidth}/>
            <text x={svgWidth - 105} y={25} fontSize="11" fill={chartConfig.defaults.labelColor} fontFamily="system-ui" fontWeight="500">
              {t('trendChart.totalLine')}
            </text>
            
            <line x1={svgWidth - 130} y1={37} x2={svgWidth - 110} y2={37} stroke={chartConfig.colors.green} strokeWidth={chartConfig.defaults.strokeWidth}/>
            <text x={svgWidth - 105} y={42} fontSize="11" fill={chartConfig.defaults.labelColor} fontFamily="system-ui" fontWeight="500">
              {t('trendChart.closedLine')}
            </text>
          </g>
          
          {/* 隐形悬停区域 */}
          {filteredData.map((d, i) => {
            const xPos = padding.left + i * xStep;
            return (
              <rect
                key={`hover-area-${i}`}
                x={xPos - xStep / 2}
                y={padding.top}
                width={xStep}
                height={chartHeight}
                fill="transparent"
                cursor="pointer"
                pointerEvents="auto"
                onMouseEnter={() => setHoveredIndex(i)}
                onMouseLeave={() => setHoveredIndex(null)}
                style={{ transition: 'none' }}
              />
            );
          })}
        </svg>

        {/* Tooltip - 显示数据 */}
        {hoveredIndex !== null && filteredData[hoveredIndex] && (
          <div
            style={{
              position: 'absolute',
              left: `${((padding.left + hoveredIndex * xStep) / svgWidth) * 100}%`,
              top: '10px',
              backgroundColor: '#fff',
              border: '1px solid #d9d9d9',
              borderRadius: '4px',
              padding: '10px 14px',
              boxShadow: '0 3px 10px rgba(0, 0, 0, 0.12)',
              fontSize: '12px',
              fontFamily: 'system-ui',
              zIndex: 10,
              transform: 'translateX(-50%)',
              whiteSpace: 'nowrap',
              pointerEvents: 'none'
            }}
          >
            <div style={{ color: '#999', fontSize: '11px', marginBottom: '6px', fontWeight: 500 }}>
              {dayjs(filteredData[hoveredIndex].date).format('YYYY-MM-DD')}
            </div>
            <div style={{ color: chartConfig.colors.red, fontWeight: '600', fontSize: '12px', marginBottom: '3px' }}>
              {t('trendChart.totalLine')}: {filteredData[hoveredIndex].total}
            </div>
            <div style={{ color: chartConfig.colors.green, fontWeight: '600', fontSize: '12px' }}>
              {t('trendChart.closedLine')}: {filteredData[hoveredIndex].closed}
            </div>
          </div>
        )}
      </div>
    );
  };

  return (
    <Card
      title={
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '20px' }}>
            <span style={{ fontSize: '16px', fontWeight: 600 }}>{t('trendChart.title')}</span>
            {defectStats && (
              <div style={{ display: 'flex', alignItems: 'center', gap: '16px', fontSize: '13px' }}>
                <span>
                  <span style={{ color: 'rgba(0,0,0,0.45)' }}>{t('defect.totalCount')}:</span>
                  <span style={{ color: '#1890ff', fontWeight: 600, marginLeft: '6px' }}>{defectStats.total}</span>
                </span>
                <span>
                  <span style={{ color: 'rgba(0,0,0,0.45)' }}>{t('defect.resolvedCount')}:</span>
                  <span style={{ color: '#52c41a', fontWeight: 600, marginLeft: '6px' }}>{defectStats.resolved}</span>
                </span>
                <span>
                  <span style={{ color: 'rgba(0,0,0,0.45)' }}>{t('defect.activeCount')}:</span>
                  <span style={{ color: '#faad14', fontWeight: 600, marginLeft: '6px' }}>{defectStats.active}</span>
                </span>
                <span>
                  <span style={{ color: 'rgba(0,0,0,0.45)' }}>{t('defect.newCount')}:</span>
                  <span style={{ color: '#ff4d4f', fontWeight: 600, marginLeft: '6px' }}>{defectStats.new}</span>
                </span>
                <span>
                  <span style={{ color: 'rgba(0,0,0,0.45)' }}>{t('defect.closedCount')}:</span>
                  <span style={{ color: '#8c8c8c', fontWeight: 600, marginLeft: '6px' }}>{defectStats.closed}</span>
                </span>
              </div>
            )}
          </div>
          <DatePicker.RangePicker
            value={dateRange}
            onChange={setDateRange}
            style={{ width: 280 }}
          />
        </div>
      }
      style={{ marginBottom: 16 }}
      bodyStyle={{ padding: '20px 24px' }}
    >
      {loading ? (
        <div style={{ display: 'flex', justifyContent: 'center', padding: 48 }}>
          <Spin />
        </div>
      ) : (
        <>
          {/* 图表 */}
          <div 
            ref={chartContainerRef}
            style={{ width: '100%', height: '320px', marginTop: 16 }}
          >
            {renderChart()}
          </div>

          {/* 图表说明 */}
          <div style={{ marginTop: 16, padding: 12, background: '#fafafa', borderRadius: 4 }}>
            <p style={{ margin: '8px 0', fontSize: 12, color: '#666' }}>
              {t('trendChart.description')}
            </p>
          </div>
        </>
      )}
    </Card>
  );
};

export default DefectTrendChart;
