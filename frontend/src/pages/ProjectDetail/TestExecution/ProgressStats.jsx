import React from 'react';
import { Card, Statistic, Row, Col } from 'antd';
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  MinusCircleOutlined,
  StopOutlined,
  PieChartOutlined
} from '@ant-design/icons';

/**
 * 进度统计组件
 * @param {Object} props
 * @param {Object} props.statistics - 统计对象 {total, nr_count, ok_count, ng_count, block_count}
 */
const ProgressStats = ({ statistics }) => {
  const { total = 0, nr_count = 0, ok_count = 0, ng_count = 0, block_count = 0 } = statistics || {};

  // 计算完成率
  const completedCount = ok_count + ng_count + block_count;
  const completionRate = total > 0 ? ((completedCount / total) * 100).toFixed(1) : 0;

  return (
    <Card title="执行进度统计" size="small">
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} md={8} lg={4}>
          <Statistic
            title="总用例数"
            value={total}
            prefix={<PieChartOutlined />}
            valueStyle={{ color: '#1890ff' }}
          />
        </Col>

        <Col xs={24} sm={12} md={8} lg={4}>
          <Statistic
            title="未执行(NR)"
            value={nr_count}
            prefix={<MinusCircleOutlined />}
            valueStyle={{ color: '#8c8c8c' }}
          />
        </Col>

        <Col xs={24} sm={12} md={8} lg={4}>
          <Statistic
            title="通过(OK)"
            value={ok_count}
            prefix={<CheckCircleOutlined />}
            valueStyle={{ color: '#52c41a' }}
          />
        </Col>

        <Col xs={24} sm={12} md={8} lg={4}>
          <Statistic
            title="失败(NG)"
            value={ng_count}
            prefix={<CloseCircleOutlined />}
            valueStyle={{ color: '#ff4d4f' }}
          />
        </Col>

        <Col xs={24} sm={12} md={8} lg={4}>
          <Statistic
            title="阻塞(Block)"
            value={block_count}
            prefix={<StopOutlined />}
            valueStyle={{ color: '#faad14' }}
          />
        </Col>

        <Col xs={24} sm={12} md={8} lg={4}>
          <Statistic
            title="完成率"
            value={completionRate}
            suffix="%"
            prefix={<PieChartOutlined />}
            valueStyle={{ color: completionRate >= 100 ? '#52c41a' : '#1890ff' }}
          />
        </Col>
      </Row>
    </Card>
  );
};

export default ProgressStats;
