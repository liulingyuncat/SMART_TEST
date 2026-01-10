/**
 * AI报告中的图表示例模板
 * 
 * 在AIReport中，可以使用以下格式嵌入Recharts图表
 */

// 线图示例
export const LINE_CHART_EXAMPLE = `
## 缺陷趋势分析

以下展示了项目中缺陷的趋势变化：

\`\`\`chart:line
{
  "title": "缺陷总数趋势",
  "data": [
    { "date": "2025-01-01", "total": 10, "closed": 2 },
    { "date": "2025-01-02", "total": 12, "closed": 3 },
    { "date": "2025-01-03", "total": 15, "closed": 5 },
    { "date": "2025-01-04", "total": 14, "closed": 6 },
    { "date": "2025-01-05", "total": 16, "closed": 8 }
  ],
  "dataKey": ["total", "closed"]
}
\`\`\`
`;

// 柱状图示例
export const BAR_CHART_EXAMPLE = `
## 按优先级统计

缺陷按优先级的分布情况：

\`\`\`chart:bar
{
  "title": "优先级分布",
  "data": [
    { "priority": "Critical", "count": 5 },
    { "priority": "High", "count": 12 },
    { "priority": "Medium", "count": 8 },
    { "priority": "Low", "count": 3 }
  ],
  "dataKey": "count"
}
\`\`\`
`;

// 饼图示例
export const PIE_CHART_EXAMPLE = `
## 状态分布

缺陷按状态的分布：

\`\`\`chart:pie
{
  "title": "缺陷状态分布",
  "data": [
    { "name": "New", "value": 8 },
    { "name": "Active", "value": 15 },
    { "name": "Resolved", "value": 12 },
    { "name": "Closed", "value": 20 }
  ],
  "dataKey": "value"
}
\`\`\`
`;

// 面积图示例
export const AREA_CHART_EXAMPLE = `
## 累积缺陷数据

项目中累积的缺陷数据：

\`\`\`chart:area
{
  "title": "累积缺陷统计",
  "data": [
    { "week": "Week 1", "resolved": 5, "active": 10 },
    { "week": "Week 2", "resolved": 8, "active": 14 },
    { "week": "Week 3", "resolved": 12, "active": 16 },
    { "week": "Week 4", "resolved": 18, "active": 18 }
  ],
  "dataKey": ["resolved", "active"]
}
\`\`\`
`;

// 雷达图示例
export const RADAR_CHART_EXAMPLE = `
## 质量维度评估

项目质量的多维度评估：

\`\`\`chart:radar
{
  "title": "质量维度评分",
  "data": [
    { "subject": "功能完整性", "current": 80, "target": 95 },
    { "subject": "代码质量", "current": 72, "target": 90 },
    { "subject": "测试覆盖", "current": 65, "target": 85 },
    { "subject": "性能", "current": 88, "target": 90 },
    { "subject": "安全性", "current": 85, "target": 95 }
  ],
  "dataKey": ["current", "target"]
}
\`\`\`
`;

/**
 * 合并所有示例
 */
export const COMPLETE_REPORT_EXAMPLE = `# AI测试质量报告

## 项目概览

本报告展示了项目的测试质量情况，包括缺陷趋势、分布统计和质量评估。

---

${LINE_CHART_EXAMPLE}

---

${BAR_CHART_EXAMPLE}

---

${PIE_CHART_EXAMPLE}

---

${AREA_CHART_EXAMPLE}

---

${RADAR_CHART_EXAMPLE}

---

## 总体评价

基于以上数据分析，项目整体质量处于良好状态，建议继续关注高优先级缺陷的解决。
`;

/**
 * 获取图表示例
 */
export function getChartExample(type) {
  const examples = {
    line: LINE_CHART_EXAMPLE,
    bar: BAR_CHART_EXAMPLE,
    pie: PIE_CHART_EXAMPLE,
    area: AREA_CHART_EXAMPLE,
    radar: RADAR_CHART_EXAMPLE,
    complete: COMPLETE_REPORT_EXAMPLE
  };
  
  return examples[type] || '';
}

/**
 * 验证图表配置
 */
export function validateChartConfig(type, config) {
  if (!config || typeof config !== 'object') {
    return { valid: false, error: '配置必须是一个对象' };
  }

  if (!Array.isArray(config.data) || config.data.length === 0) {
    return { valid: false, error: '配置必须包含非空的data数组' };
  }

  // 验证dataKey
  if (config.dataKey) {
    if (Array.isArray(config.dataKey)) {
      // 验证所有key都存在于数据中
      const firstData = config.data[0];
      for (const key of config.dataKey) {
        if (!(key in firstData)) {
          return { valid: false, error: `dataKey "${key}" 在数据中不存在` };
        }
      }
    } else if (typeof config.dataKey === 'string') {
      if (!(config.dataKey in config.data[0])) {
        return { valid: false, error: `dataKey "${config.dataKey}" 在数据中不存在` };
      }
    }
  }

  return { valid: true };
}
