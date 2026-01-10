import React, { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Table, Tag, Button, Space, Popconfirm, Tooltip } from 'antd';
import { EyeOutlined, DeleteOutlined } from '@ant-design/icons';
import {
  DEFECT_STATUS,
  DEFECT_STATUS_COLORS,
  DEFECT_PRIORITY,
  DEFECT_PRIORITY_COLORS,
  DEFECT_SEVERITY,
  DEFECT_SEVERITY_COLORS,
} from '../../../../constants/defect';

/**
 * 缺陷列表表格
 * 显示缺陷列表，支持分页、查看详情、删除操作
 */
const DefectTable = ({
  defects,
  total,
  page,
  pageSize,
  onPageChange,
  onDetail,
  onDelete,
}) => {
  const { t, i18n } = useTranslation();

  // 使用 useMemo 缓存翻译标签，只在语言变化时重新计算
  const labels = useMemo(() => ({
    defectId: t('defect.defectId', '缺陷ID'),
    title: t('defect.title', '标题'),
    status: t('defect.status', '状态'),
    priority: t('defect.priority', '优先级'),
    severity: t('defect.severity', '严重程度'),
    subject: t('defect.subject', '主题'),
    phase: t('defect.phase', '阶段'),
    assignee: t('defect.assignee', '负责人'),
    createdAt: t('common.createdAt', '创建时间'),
    actions: t('common.actions', '操作'),
    confirmDelete: t('common.confirmDelete', '确认删除吗？'),
    ok: t('common.ok', '确定'),
    cancel: t('common.cancel', '取消'),
  }), [t, i18n.language]);

  // 状态标签映射
  const statusLabels = useMemo(() => ({
    [DEFECT_STATUS.NEW]: t('defect.statusNew', '新建'),
    [DEFECT_STATUS.ACTIVE]: t('defect.statusActive', '处理中'),
    [DEFECT_STATUS.RESOLVED]: t('defect.statusResolved', '已解决'),
    [DEFECT_STATUS.CLOSED]: t('defect.statusClosed', '已关闭'),
  }), [t, i18n.language]);

  // 优先级标签映射（A/B/C/D）
  const priorityLabels = useMemo(() => ({
    [DEFECT_PRIORITY.A]: t('defect.priorityA', 'A'),
    [DEFECT_PRIORITY.B]: t('defect.priorityB', 'B'),
    [DEFECT_PRIORITY.C]: t('defect.priorityC', 'C'),
    [DEFECT_PRIORITY.D]: t('defect.priorityD', 'D'),
  }), [t, i18n.language]);

  // 严重程度标签映射（A/B/C/D）
  const severityLabels = useMemo(() => ({
    [DEFECT_SEVERITY.A]: t('defect.severityA', 'A'),
    [DEFECT_SEVERITY.B]: t('defect.severityB', 'B'),
    [DEFECT_SEVERITY.C]: t('defect.severityC', 'C'),
    [DEFECT_SEVERITY.D]: t('defect.severityD', 'D'),
  }), [t, i18n.language]);

  // 获取状态标签
  const getStatusTag = (status) => {
    const color = DEFECT_STATUS_COLORS[status];
    const label = statusLabels[status];
    if (!label) return status;
    return <Tag color={color}>{label}</Tag>;
  };

  // 获取优先级标签
  const getPriorityTag = (priority) => {
    const color = DEFECT_PRIORITY_COLORS[priority];
    const label = priorityLabels[priority];
    if (!label) return priority;
    return <Tag color={color}>{label}</Tag>;
  };

  // 获取严重程度标签
  const getSeverityTag = (severity) => {
    const color = DEFECT_SEVERITY_COLORS[severity];
    const label = severityLabels[severity];
    if (!label) return severity;
    return <Tag color={color}>{label}</Tag>;
  };

  // 表格列定义
  const columns = useMemo(() => [
    {
      title: labels.defectId,
      dataIndex: 'defect_id',
      key: 'defect_id',
      width: 120,
      fixed: 'left',
      render: (text, record) => (
        <a onClick={() => onDetail?.(record)}>{text}</a>
      ),
    },
    {
      title: labels.title,
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
      render: (text, record) => (
        <Tooltip title={text}>
          <a onClick={() => onDetail?.(record)}>{text}</a>
        </Tooltip>
      ),
    },
    {
      title: labels.status,
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => getStatusTag(status),
    },
    {
      title: labels.priority,
      dataIndex: 'priority',
      key: 'priority',
      width: 100,
      render: (priority) => getPriorityTag(priority),
    },
    {
      title: labels.severity,
      dataIndex: 'severity',
      key: 'severity',
      width: 100,
      render: (severity) => getSeverityTag(severity),
    },
    {
      title: labels.subject,
      dataIndex: ['subject', 'name'],
      key: 'subject',
      width: 120,
      ellipsis: true,
    },
    {
      title: labels.phase,
      dataIndex: ['phase', 'name'],
      key: 'phase',
      width: 120,
      ellipsis: true,
    },
    {
      title: labels.createdAt,
      dataIndex: 'created_at',
      key: 'created_at',
      width: 160,
    },
    {
      title: labels.actions,
      key: 'actions',
      width: 120,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => onDetail?.(record)}
          />
        </Space>
      ),
    },
  ], [labels, statusLabels, priorityLabels, severityLabels, onDetail, onDelete]);

  // 翻译总数文本
  const totalLabel = useMemo(() => {
    return (total) => t('common.totalItems', { total, defaultValue: `共 ${total} 条` });
  }, [t, i18n.language]);

  return (
    <Table
      dataSource={defects}
      columns={columns}
      rowKey="id"
      pagination={{
        current: page,
        pageSize: pageSize,
        total: total,
        showSizeChanger: true,
        showQuickJumper: true,
        showTotal: totalLabel,
        pageSizeOptions: ['20', '50', '100'],
        onChange: onPageChange,
      }}
      scroll={{ x: 1200 }}
      size="middle"
    />
  );
};

export default DefectTable;
