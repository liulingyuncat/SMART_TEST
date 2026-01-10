import React from 'react';
import { Menu, Dropdown, Button } from 'antd';
import { InsertRowAboveOutlined } from '@ant-design/icons';
import {
  getChartExample,
  LINE_CHART_EXAMPLE,
  BAR_CHART_EXAMPLE,
  PIE_CHART_EXAMPLE,
  AREA_CHART_EXAMPLE,
  RADAR_CHART_EXAMPLE
} from './chartTemplates';

/**
 * 图表插入菜单组件
 * 提供快速插入各种Recharts图表的菜单
 */
const ChartInsertMenu = ({ onInsert }) => {
  const items = [
    {
      key: 'line',
      label: '📈 线图 (Line Chart)',
      onClick: () => onInsert(getChartExample('line'))
    },
    {
      key: 'bar',
      label: '📊 柱状图 (Bar Chart)',
      onClick: () => onInsert(getChartExample('bar'))
    },
    {
      key: 'pie',
      label: '🥧 饼图 (Pie Chart)',
      onClick: () => onInsert(getChartExample('pie'))
    },
    {
      key: 'area',
      label: '📉 面积图 (Area Chart)',
      onClick: () => onInsert(getChartExample('area'))
    },
    {
      key: 'radar',
      label: '🎯 雷达图 (Radar Chart)',
      onClick: () => onInsert(getChartExample('radar'))
    },
    {
      type: 'divider'
    },
    {
      key: 'help',
      label: '❓ 图表使用帮助',
      onClick: () => {
        alert(`图表使用说明：

1. 使用特殊的代码块语法：
   \`\`\`chart:类型
   {JSON配置}
   \`\`\`

2. 支持的图表类型：
   - chart:line    - 线图/趋势图
   - chart:bar     - 柱状图
   - chart:pie     - 饼图
   - chart:area    - 面积图
   - chart:radar   - 雷达图

3. 配置示例：
   {
     "title": "图表标题",
     "data": [{...}, {...}],
     "dataKey": "字段名或字段名数组",
     "colors": ["#8884d8", "#82ca9d"]
   }

4. 数据要求：
   - data 必须是数组
   - dataKey 必须匹配数据中的字段名
   - 推荐至少3条数据记录`);
      }
    }
  ];

  return (
    <Dropdown
      menu={{ items }}
      placement="bottomLeft"
      trigger={['click']}
    >
      <Button icon={<InsertRowAboveOutlined />} type="dashed">
        插入图表
      </Button>
    </Dropdown>
  );
};

export default ChartInsertMenu;
