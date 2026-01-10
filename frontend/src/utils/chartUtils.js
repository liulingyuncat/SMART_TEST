/**
 * 图表工具函数
 * 用于优化SVG图表的视觉效果
 */

/**
 * 使用贝塞尔曲线生成平滑的路径
 * @param {Array} points - 坐标点数组 [{x, y}, ...]
 * @returns {string} SVG路径字符串
 */
export function generateSmoothPath(points) {
  if (!points || points.length === 0) return '';
  if (points.length === 1) return `M ${points[0].x} ${points[0].y}`;

  let path = `M ${points[0].x} ${points[0].y}`;
  
  for (let i = 0; i < points.length - 1; i++) {
    const current = points[i];
    const next = points[i + 1];
    
    // 计算控制点
    const prev = i > 0 ? points[i - 1] : current;
    const nextNext = i + 2 < points.length ? points[i + 2] : next;
    
    // 使用Catmull-Rom样条曲线的控制点公式
    const cp1x = current.x + (next.x - prev.x) / 6;
    const cp1y = current.y + (next.y - prev.y) / 6;
    const cp2x = next.x - (nextNext.x - current.x) / 6;
    const cp2y = next.y - (nextNext.y - current.y) / 6;
    
    path += ` C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${next.x} ${next.y}`;
  }
  
  return path;
}

/**
 * 生成直线路径（不使用平滑曲线）
 * @param {Array} points - 坐标点数组 [{x, y}, ...]
 * @returns {string} SVG路径字符串
 */
export function generateLinePath(points) {
  if (!points || points.length === 0) return '';
  if (points.length === 1) return `M ${points[0].x} ${points[0].y}`;

  let path = `M ${points[0].x} ${points[0].y}`;
  
  for (let i = 1; i < points.length; i++) {
    path += ` L ${points[i].x} ${points[i].y}`;
  }
  
  return path;
}

/**
 * 生成SVG渐变ID和定义
 * @param {string} id - 渐变ID
 * @param {string} color1 - 起始颜色
 * @param {string} color2 - 结束颜色
 * @returns {JSX} SVG defs元素
 */
export function createGradientDef(id, color1, color2, opacity1 = 1, opacity2 = 0.1) {
  return (
    <defs>
      <linearGradient id={id} x1="0%" y1="0%" x2="0%" y2="100%">
        <stop offset="0%" stopColor={color1} stopOpacity={opacity1} />
        <stop offset="100%" stopColor={color2} stopOpacity={opacity2} />
      </linearGradient>
    </defs>
  );
}

/**
 * 获取图表配置
 */
export const chartConfig = {
  // 颜色主题
  colors: {
    red: '#ff4d4f',
    green: '#52c41a',
    blue: '#1890ff',
    orange: '#faad14',
    purple: '#722ed1',
    cyan: '#13c2c2'
  },
  
  // 默认样式
  defaults: {
    strokeWidth: 2,
    pointRadius: 3,
    pointRadiusHover: 5,
    gridColor: '#f0f0f0',
    gridDasharray: '3,3',
    axisColor: '#999',
    labelColor: '#666',
    backgroundColor: '#fafafa'
  },
  
  // 动画配置
  animation: {
    duration: 150,
    easing: 'ease'
  }
};

/**
 * 生成网格线
 * @param {number} count - 网格线数量
 * @param {number} top - 顶部间距
 * @param {number} bottom - 底部间距
 * @param {number} left - 左边间距
 * @param {number} right - 右边间距
 * @param {number} svgWidth - SVG宽度
 * @param {number} svgHeight - SVG高度
 * @returns {Array} 网格线信息数组
 */
export function generateGridLines(count, top, bottom, left, right, svgWidth, svgHeight) {
  const chartHeight = svgHeight - top - bottom;
  const lines = [];
  
  for (let i = 0; i <= count; i++) {
    const y = top + chartHeight - (i * (chartHeight / count));
    lines.push({
      y,
      value: i,
      x1: left,
      x2: svgWidth - right
    });
  }
  
  return lines;
}

/**
 * 根据数据范围计算最佳网格线数量
 * @param {number} maxValue - 最大值
 * @returns {number} 推荐的网格线数量
 */
export function getOptimalGridLineCount(maxValue) {
  if (maxValue <= 5) return maxValue;
  if (maxValue <= 10) return 5;
  if (maxValue <= 20) return 4;
  if (maxValue <= 50) return 5;
  return Math.ceil(maxValue / 10);
}
