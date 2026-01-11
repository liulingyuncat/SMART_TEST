import React, { useState, useCallback, useEffect } from 'react';
import { Button, Popconfirm, message, Input, Row, Col, Space } from 'antd';
import { DeleteOutlined, EditOutlined, SaveOutlined, CloseOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ApiLeftSidePanel from '../components/ApiLeftSidePanel';
import EditableTable from '../../ManualTestTabs/components/EditableTable';
import { updateCaseGroup } from '../../../../api/manualCase';
import './ApiCaseManagementTab.css';

/**
 * API用例管理Tab容器组件
 * 采用左右分栏布局：左栏200px固定宽度，右栏自适应
 * 与WebCaseManagementTab的差异：不包含语言筛选和LanguageFilter组件
 */
const ApiCaseManagementTab = ({ projectId }) => {
  const { t } = useTranslation();
  const [collapsed, setCollapsed] = useState(false); // 左栏收束状态
  const [selectedCaseGroup, setSelectedCaseGroup] = useState(null); // 当前选中的用例集
  const [refreshKey, setRefreshKey] = useState(0); // 用于刷新表格
  const [batchDeleteInfo, setBatchDeleteInfo] = useState(null); // 批量删除信息
  
  // 元数据状态
  const [metadata, setMetadata] = useState({
    protocol: 'https',
    server: '',
    port: '',
    user: '',
    password: ''
  });
  const [metadataEditing, setMetadataEditing] = useState(false); // 元数据编辑状态
  const [metadataBackup, setMetadataBackup] = useState(null); // 编辑前的备份
  const [metadataSaving, setMetadataSaving] = useState(false); // 保存中状态

  // 当选中用例集变化时，加载该用例集的元数据
  useEffect(() => {
    if (selectedCaseGroup) {
      setMetadata({
        protocol: selectedCaseGroup.meta_protocol || 'https',
        server: selectedCaseGroup.meta_server || '',
        port: selectedCaseGroup.meta_port || '',
        user: selectedCaseGroup.meta_user || '',
        password: selectedCaseGroup.meta_password || ''
      });
    } else {
      setMetadata({ protocol: 'https', server: '', port: '', user: '', password: '' });
    }
  }, [selectedCaseGroup]);

  // 元数据变更处理
  const handleMetadataChange = (field, value) => {
    setMetadata(prev => ({
      ...prev,
      [field]: value
    }));
  };

  // 开始编辑元数据
  const handleEditMetadata = () => {
    setMetadataBackup({ ...metadata });
    setMetadataEditing(true);
  };

  // 取消编辑
  const handleCancelEditMetadata = () => {
    if (metadataBackup) {
      setMetadata(metadataBackup);
    }
    setMetadataEditing(false);
    setMetadataBackup(null);
  };

  // 保存元数据到后端
  const handleSaveMetadata = async () => {
    if (!selectedCaseGroup) return;
    
    setMetadataSaving(true);
    try {
      await updateCaseGroup(selectedCaseGroup.id, {
        metaProtocol: metadata.protocol,
        metaServer: metadata.server,
        metaPort: metadata.port,
        metaUser: metadata.user,
        metaPassword: metadata.password
      });
      // 更新本地缓存的用例集数据
      selectedCaseGroup.meta_protocol = metadata.protocol;
      selectedCaseGroup.meta_server = metadata.server;
      selectedCaseGroup.meta_port = metadata.port;
      selectedCaseGroup.meta_user = metadata.user;
      selectedCaseGroup.meta_password = metadata.password;
      
      message.success(t('message.saveSuccess'));
      setMetadataEditing(false);
      setMetadataBackup(null);
    } catch (error) {
      console.error('[ApiCaseManagementTab] Failed to save metadata:', error);
      message.error(t('message.saveFailed'));
    } finally {
      setMetadataSaving(false);
    }
  };

  // 左栏收束状态变更回调
  const handleCollapseChange = () => {
    setCollapsed(!collapsed);
  };

  // 用例集切换回调
  const handleCaseSwitch = (caseGroup) => {
    console.log('[ApiCaseManagementTab] 切换用例集:', caseGroup);
    setSelectedCaseGroup(caseGroup);
    setMetadataEditing(false); // 切换用例集时退出编辑状态
    setMetadataBackup(null);
    setRefreshKey(prev => prev + 1); // 刷新表格
  };

  // 用例集更新回调（创建/编辑/删除用例集后触发）
  const handleCaseGroupsUpdated = () => {
    setRefreshKey(prev => prev + 1); // 刷新表格
  };

  // 批量删除 - 调用EditableTable暴露的删除函数
  const handleBatchDelete = () => {
    if (!batchDeleteInfo || !batchDeleteInfo.executeDelete) {
      message.warning('请先在表格中选择要删除的用例');
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
    <div className="api-case-management-tab">
      {/* 左栏操作面板 */}
      <ApiLeftSidePanel
        projectId={projectId}
        collapsed={collapsed}
        selectedCaseGroup={selectedCaseGroup}
        onCaseSwitch={handleCaseSwitch}
        onCollapse={handleCollapseChange}
        onCaseGroupsUpdated={handleCaseGroupsUpdated}
      />

      {/* 右栏内容区 */}
      <div className={`right-content-panel ${collapsed ? 'full-width' : ''}`}>
        {/* 顶部工具栏：仅批量删除按钮，无语言筛选 */}
        <div style={{ 
          display: 'flex', 
          justifyContent: 'flex-end', 
          alignItems: 'center',
          marginBottom: '4px',
          padding: '0 4px'
        }}>
          <Popconfirm
            title={t('project.batchDeleteConfirm', { count: batchDeleteInfo?.selectedCount || 0 })}
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
              {t('project.batchDelete')}
            </Button>
          </Popconfirm>
        </div>

        {/* 元数据输入区 */}
        <div style={{
          padding: '12px 8px',
          background: '#fafafa',
          borderRadius: '4px',
          marginBottom: '8px'
        }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
            <div style={{ fontSize: '13px', fontWeight: 500, color: 'rgba(0,0,0,0.85)' }}>Web Server</div>
            <Space size={4}>
              {!metadataEditing ? (
                <Button 
                  size="small" 
                  icon={<EditOutlined />}
                  onClick={handleEditMetadata}
                  disabled={!selectedCaseGroup}
                >
                  {t('common.edit')}
                </Button>
              ) : (
                <>
                  <Button 
                    size="small" 
                    icon={<CloseOutlined />}
                    onClick={handleCancelEditMetadata}
                  >
                    {t('common.cancel')}
                  </Button>
                  <Button 
                    size="small" 
                    type="primary"
                    icon={<SaveOutlined />}
                    onClick={handleSaveMetadata}
                    loading={metadataSaving}
                  >
                    {t('common.save')}
                  </Button>
                </>
              )}
            </Space>
          </div>
          {/* 第一行: Protocol / Server Name or IP / Port Number */}
          <Row gutter={[12, 8]}>
            <Col>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '60px', textAlign: 'right' }}>Protocol：</span>
                <Input
                  size="small"
                  style={{ width: '180px' }}
                  value={metadata.protocol}
                  onChange={(e) => handleMetadataChange('protocol', e.target.value)}
                  disabled={!metadataEditing}
                />
              </div>
            </Col>
            <Col>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '120px', textAlign: 'right' }}>Server Name or IP：</span>
                <Input
                  size="small"
                  style={{ width: '360px' }}
                  value={metadata.server}
                  onChange={(e) => handleMetadataChange('server', e.target.value)}
                  disabled={!metadataEditing}
                />
              </div>
            </Col>
            <Col>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '80px', textAlign: 'right' }}>Port Number：</span>
                <Input
                  size="small"
                  style={{ width: '180px' }}
                  value={metadata.port}
                  onChange={(e) => handleMetadataChange('port', e.target.value)}
                  disabled={!metadataEditing}
                />
              </div>
            </Col>
          </Row>
          {/* 第二行: User 和 Password */}
          <Row gutter={[12, 8]} style={{ marginTop: '8px' }}>
            <Col>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '60px', textAlign: 'right' }}>User：</span>
                <Input
                  size="small"
                  style={{ width: '180px' }}
                  value={metadata.user}
                  onChange={(e) => handleMetadataChange('user', e.target.value)}
                  disabled={!metadataEditing}
                />
              </div>
            </Col>
            <Col>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '120px', textAlign: 'right' }}>Password：</span>
                <Input.Password
                  size="small"
                  style={{ width: '180px' }}
                  value={metadata.password}
                  onChange={(e) => handleMetadataChange('password', e.target.value)}
                  disabled={!metadataEditing}
                />
              </div>
            </Col>
          </Row>
        </div>

        {/* 表格内容区 */}
        <div className="table-container">
          {selectedCaseGroup === null ? (
            <div className="empty-state">
              <div className="empty-state-icon">📋</div>
              <div>请点击左侧"创建API用例集"按钮添加第一个用例集</div>
            </div>
          ) : (
            <EditableTable
              key={refreshKey}
              projectId={projectId}
              apiModule="api-cases"
              caseType="api"
              caseGroupFilter={selectedCaseGroup.case_group}
              onBatchDeleteRequest={handleBatchDeleteRequest}
              hiddenButtons={['saveVersion', 'exportTemplate', 'exportCases', 'importCases']}
              knownPasswords={[metadata.password].filter(Boolean)}
            />
          )}
        </div>
      </div>
    </div>
  );
};

export default ApiCaseManagementTab;
