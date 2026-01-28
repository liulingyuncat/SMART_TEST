import React, { useState, useCallback } from 'react';
import { message, Button, Space, Popconfirm } from 'antd';
import { DownloadOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ManualLeftSidePanel from '../components/ManualLeftSidePanel';
import LanguageFilter from '../components/LanguageFilter';
import EditableTable from '../components/EditableTable';
import ReorderModal from '../components/ReorderModal';
import { exportCasesByLanguage } from '../../../../api/manualCase';
import './ManualCaseManagementTab.css';

/**
 * 手工用例管理Tab容器组件
 * 采用左右分栏布局：左栏200px固定宽度，右栏自适应
 */
const ManualCaseManagementTab = ({ projectId }) => {
  const { t } = useTranslation();
  const [language, setLanguage] = useState('中文');
  const [collapsed, setCollapsed] = useState(false); // 左栏收束状态
  const [reorderModalVisible, setReorderModalVisible] = useState(false);
  const [casesForReorder, setCasesForReorder] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [refreshKey, setRefreshKey] = useState(0); // 用于刷新表格
  const [selectedCaseGroup, setSelectedCaseGroup] = useState(null); // 当前选中的用例集
  const [batchDeleteInfo, setBatchDeleteInfo] = useState(null); // 批量删除信息

  // 语言筛选变更
  const handleLanguageChange = (newLanguage) => {
    setLanguage(newLanguage);
  };

  // 打开重排对话框（接收当前显示的cases数组和页码）
  const handleReorderClick = (currentCases, pageNumber) => {
    setCasesForReorder(currentCases || []);
    setCurrentPage(pageNumber || 1);
    setReorderModalVisible(true);
  };

  // 重排成功回调
  const handleReorderSuccess = () => {
    setReorderModalVisible(false);
    setRefreshKey(prev => prev + 1); // 刷新表格
  };

  // 用例创建成功回调
  const handleCaseCreated = () => {
    setRefreshKey(prev => prev + 1); // 刷新表格和用例一览
  };

  // 用例更新回调
  const handleCaseUpdated = () => {
    setRefreshKey(prev => prev + 1); // 刷新表格
  };

  // 左栏收束状态变更回调
  const handleCollapseChange = (isCollapsed) => {
    setCollapsed(isCollapsed);
  };

  // 用例集切换回调
  const handleCaseSwitch = (caseGroup) => {
    console.log('[ManualCaseManagementTab] 切换用例集:', caseGroup);
    setSelectedCaseGroup(caseGroup);
    setRefreshKey(prev => prev + 1); // 刷新表格
  };

  // 用例删除回调
  const handleCaseDeleted = () => {
    console.log('[ManualCaseManagementTab] 用例删除');
    setRefreshKey(prev => prev + 1); // 刷新表格
  };

  // 用例集更新回调（创建/编辑/删除用例集后触发）
  const handleCaseGroupsUpdated = () => {
    setRefreshKey(prev => prev + 1); // 刷新表格
  };

  // T44: 按语言导出用例
  const handleExport = async () => {
    if (!selectedCaseGroup) {
      message.warning(t('manualTest.selectCaseGroupFirst'));
      return;
    }

    try {
      await exportCasesByLanguage(projectId, 'overall', language, selectedCaseGroup.case_group);
      message.success(t('message.exportSuccess'));
    } catch (error) {
      console.error('导出用例失败:', error);
      message.error(t('message.exportFailed'));
    }
  };

  // 批量删除 - 调用EditableTable暴露的删除函数
  const handleBatchDelete = () => {
    if (!batchDeleteInfo || !batchDeleteInfo.executeDelete) {
      message.warning(t('manualTest.selectCasesToDelete'));
      return;
    }
    // 调用EditableTable暴露的批量删除函数
    batchDeleteInfo.executeDelete();
  };

  // 接收EditableTable的批量删除请求
  const handleBatchDeleteRequest = useCallback((info) => {
    setBatchDeleteInfo(info);
  }, []);

  return (
    <div className="manual-case-management-tab">
      {/* 左栏操作面板 */}
      <ManualLeftSidePanel
        projectId={projectId}
        language={language}
        collapsed={collapsed}
        selectedCaseGroup={selectedCaseGroup}
        onCaseSwitch={handleCaseSwitch}
        onCollapse={handleCollapseChange}
        onCaseGroupsUpdated={handleCaseGroupsUpdated}
      />

      {/* 右栏内容区 */}
      <div className={`right-content-panel ${collapsed ? 'full-width' : ''}`}>
        {/* 顶部工具栏：语言切换 + 导出/删除按钮 */}
        <div style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '8px',
          padding: '0 8px'
        }}>
          <LanguageFilter
            value={language}
            onChange={handleLanguageChange}
          />

          {/* 右侧操作按钮 */}
          <Space size={8}>
            <Button
              icon={<DownloadOutlined />}
              onClick={handleExport}
              disabled={!selectedCaseGroup}
            >
              {t('manualTest.exportCases')}
            </Button>
            <Popconfirm
              title={t('manualTest.batchDeleteConfirm', { count: batchDeleteInfo?.selectedCount || 0 })}
              onConfirm={handleBatchDelete}
              okText={t('common.ok')}
              cancelText={t('common.cancel')}
              disabled={!selectedCaseGroup || !batchDeleteInfo || batchDeleteInfo.selectedCount === 0}
            >
              <Button
                danger
                icon={<DeleteOutlined />}
                disabled={!selectedCaseGroup || !batchDeleteInfo || batchDeleteInfo.selectedCount === 0}
              >
                {t('manualTest.batchDelete')}
              </Button>
            </Popconfirm>
          </Space>
        </div>

        {selectedCaseGroup === null ? (
          <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            height: 'calc(100vh - 200px)',
            fontSize: '16px',
            color: '#999'
          }}>
            {t('manualTest.clickCreateCaseGroup')}
          </div>
        ) : (
          <EditableTable
            key={refreshKey}
            projectId={projectId}
            caseType="overall"
            language={language}
            caseGroupFilter={selectedCaseGroup.case_group}
            onReorderClick={handleReorderClick}
            onBatchDeleteRequest={handleBatchDeleteRequest}
            hiddenButtons={['saveVersion', 'exportTemplate', 'aiSupplement', 'exportCases', 'importCases']}
          />
        )}

        {/* 重排对话框 */}
        <ReorderModal
          visible={reorderModalVisible}
          cases={casesForReorder}
          currentPage={currentPage}
          projectId={projectId}
          caseType="overall"
          language={language}
          onSuccess={handleReorderSuccess}
          onCancel={() => setReorderModalVisible(false)}
        />
      </div>
    </div>
  );
};

export default ManualCaseManagementTab;
