import React, { useState } from 'react';
import { Button, Modal, Input, List, Empty, message, Space, Collapse } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ReportListItem from './ReportListItem';
import './AIReportLeftPanel.css';

const AIReportLeftPanel = ({ 
  reports, 
  selectedReportId, 
  loading, 
  onSelect, 
  onCreate, 
  onDelete 
}) => {
  const { t } = useTranslation();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [newReportName, setNewReportName] = useState('');
  const [isCreating, setIsCreating] = useState(false);
  const [nameError, setNameError] = useState('');

  // 打开创建对话框
  const handleCreateClick = () => {
    setNewReportName('');
    setNameError('');
    setIsCreateModalOpen(true);
  };

  // 保存新报告
  const handleCreateSave = async () => {
    if (!newReportName.trim()) {
      setNameError(t('aiReport.reportNameRequired'));
      return;
    }

    setIsCreating(true);
    try {
      await onCreate(newReportName);
      setIsCreateModalOpen(false);
      setNewReportName('');
      setNameError('');
      message.success(t('aiReport.createSuccess'));
    } catch (error) {
      if (error.response?.status === 409) {
        setNameError(t('aiReport.reportNameDuplicate'));
      } else {
        message.error(error.message || t('aiReport.createFailed'));
      }
    } finally {
      setIsCreating(false);
    }
  };

  return (
    <div className="ai-report-left-panel">
      {/* 创建报告按钮 */}
      <Button
        type="primary"
        icon={<PlusOutlined />}
        block
        onClick={handleCreateClick}
        className="create-button"
      >
        {t('aiReport.createReport')}
      </Button>

      {/* 报告一览 */}
      <div className="reports-container">
        <div className="panel-title">{t('aiReport.reportList')}</div>
        
        {loading ? (
          <Empty description={t('common.loading')} />
        ) : reports.length === 0 ? (
          <Empty
            description={t('aiReport.noReports')}
            style={{ marginTop: '20px' }}
          />
        ) : (
          <List
            dataSource={reports}
            renderItem={(report) => (
              <ReportListItem
                key={report.id}
                reportData={report}
                isSelected={selectedReportId === report.id}
                onSelect={onSelect}
                onDelete={onDelete}
              />
            )}
          />
        )}
      </div>

      {/* 创建报告对话框 */}
      <Modal
        title={t('aiReport.createReport')}
        open={isCreateModalOpen}
        onOk={handleCreateSave}
        onCancel={() => setIsCreateModalOpen(false)}
        okText={t('common.create')}
        cancelText={t('common.cancel')}
        confirmLoading={isCreating}
      >
        <Input
          value={newReportName}
          onChange={(e) => {
            setNewReportName(e.target.value);
            if (e.target.value.trim()) {
              setNameError('');
            }
          }}
          placeholder={t('aiReport.enterReportName')}
          onPressEnter={handleCreateSave}
          autoFocus
          status={nameError ? 'error' : ''}
        />
        {nameError && (
          <div className="error-message">{nameError}</div>
        )}
      </Modal>
    </div>
  );
};

export default AIReportLeftPanel;
