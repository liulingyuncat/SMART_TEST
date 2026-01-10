import React from 'react';
import { Button, Popconfirm, message, Space, Tooltip } from 'antd';
import { CopyOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import './ReportListItem.css';

const ReportListItem = ({ reportData, isSelected, onSelect, onDelete }) => {
  const { t } = useTranslation();

  // 格式化日期
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    const year = date.getFullYear().toString().slice(-2);
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    return `${year}/${month}/${day}`;
  };

  // 复制报告名称
  const handleCopy = (e) => {
    e.stopPropagation();
    navigator.clipboard.writeText(reportData.name);
    message.success(t('common.copySuccess', { defaultValue: '复制成功' }));
  };

  // 删除
  const handleDelete = (e) => {
    e.stopPropagation();
  };

  return (
    <>
      <div
        className={`report-list-item ${isSelected ? 'selected' : ''}`}
        onClick={() => onSelect(reportData.id)}
      >
        <div className="report-info">
          <div className="report-name">{reportData.name}</div>
          <div className="report-date">{formatDate(reportData.created_at)}</div>
        </div>
        <div className="report-actions">
          <Space size="small">
            <Tooltip title={t('common.copy', { defaultValue: '复制' })}>
              <Button
                type="text"
                size="small"
                icon={<CopyOutlined />}
                onClick={handleCopy}
                className="copy-button"
              />
            </Tooltip>
            <Tooltip title={t('aiReport.delete')}>
              <Popconfirm
                title={t('aiReport.deleteConfirm')}
                onConfirm={() => onDelete(reportData.id)}
                okText={t('common.confirm')}
                cancelText={t('common.cancel')}
              >
                <Button
                  type="text"
                  size="small"
                  icon={<DeleteOutlined />}
                  onClick={handleDelete}
                  danger
                />
              </Popconfirm>
            </Tooltip>
          </Space>
        </div>
      </div>
    </>
  );
};

export default ReportListItem;
