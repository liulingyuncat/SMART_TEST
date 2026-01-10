import React, { useState, useEffect, useRef } from 'react';
import { Button, message, Spin, Modal, Input, Radio } from 'antd';
import { PlusOutlined, SaveOutlined, RightOutlined, LeftOutlined, UnorderedListOutlined, DownloadOutlined, UploadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import {
  getWebCaseGroups,
  createWebCaseGroup,
  updateWebCaseGroup,
  deleteWebCaseGroup,
  saveWebVersion,
  exportWebTemplate,
  importWebCases
} from '../../../../api/autoCase';
import CaseListItem from '../../../ProjectDetail/ManualTestTabs/components/CaseListItem';
import WebVersionList from './WebVersionList';
import './WebLeftSidePanel.css';

/**
 * Web用例库左侧操作面板
 * 功能区：1.创建用例集(蓝色按钮) 2.版本保存(白色按钮) 3.版本一览(白色按钮+弹窗) 4.用例集一览(标题右侧收束控件)
 */
const WebLeftSidePanel = ({
  projectId,
  collapsed = false,
  selectedCaseGroup,
  onCaseSwitch,
  onCollapse,
  onCaseGroupsUpdated
}) => {
  const { t } = useTranslation();
  const [caseGroups, setCaseGroups] = useState([]);
  const [casesLoading, setCasesLoading] = useState(false);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [newCaseGroupName, setNewCaseGroupName] = useState('');
  const [caseGroupNameError, setCaseGroupNameError] = useState(''); // 用例集名称错误提示
  const [savingVersion, setSavingVersion] = useState(false);
  const [versionModalVisible, setVersionModalVisible] = useState(false);
  const [versionListKey, setVersionListKey] = useState(0);
  const versionListRef = useRef(null);
  const [importModalVisible, setImportModalVisible] = useState(false);
  const [selectedImportGroup, setSelectedImportGroup] = useState(null);
  const [newImportGroupName, setNewImportGroupName] = useState('');
  const [importFile, setImportFile] = useState(null);
  const [importLoading, setImportLoading] = useState(false);
  const fileInputRef = useRef(null);

  useEffect(() => {
    if (projectId) {
      loadCaseGroups();
    }
  }, [projectId]);

  // 加载用例集列表
  const loadCaseGroups = async () => {
    if (!projectId) {
      console.warn('[WebLeftSidePanel] projectId为空');
      return;
    }

    console.log('[WebLeftSidePanel] 开始加载Web用例集列表:', projectId);
    setCasesLoading(true);
    try {
      const response = await getWebCaseGroups(projectId);
      console.log('[WebLeftSidePanel] 用例集API返回:', response);

      if (response && response.length > 0) {
        // 转换数据格式以兼容CaseListItem组件，同时保留元数据
        const formattedGroups = response.map(group => ({
          case_group: group.group_name,
          id: group.id,
          _groupId: group.id,
          meta_protocol: group.meta_protocol,
          meta_server: group.meta_server,
          meta_port: group.meta_port,
          meta_user: group.meta_user,
          meta_password: group.meta_password
        }));
        setCaseGroups(formattedGroups);

        // 默认选中第一个用例集
        if (formattedGroups.length > 0 && onCaseSwitch && !selectedCaseGroup) {
          console.log('[WebLeftSidePanel] 默认选中第一个用例集:', formattedGroups[0].case_group);
          setTimeout(() => {
            onCaseSwitch(formattedGroups[0]);
          }, 0);
        }
      } else {
        console.log('[WebLeftSidePanel] 用例集列表为空');
        setCaseGroups([]);
      }
    } catch (error) {
      console.error('[WebLeftSidePanel] 加载用例集列表失败:', error);
      message.error(t('web_case.loadCaseGroupsFailed'));
    } finally {
      setCasesLoading(false);
    }
  };

  // 打开创建用例集对话框
  const handleCreateCaseGroup = () => {
    setNewCaseGroupName('');
    setCaseGroupNameError(''); // 清除错误提示
    setCreateModalVisible(true);
  };

  // 保存新用例集
  const handleSaveNewCaseGroup = async () => {
    if (!projectId) {
      console.error('[WebLeftSidePanel] projectId为空');
      message.error(t('common.error'));
      return;
    }

    if (!newCaseGroupName || newCaseGroupName.trim() === '') {
      setCaseGroupNameError(t('web_case.caseGroupNameRequired'));
      return;
    }

    const trimmedName = newCaseGroupName.trim();

    // 检查重名
    const isDuplicate = caseGroups.some(group => group.case_group === trimmedName);
    if (isDuplicate) {
      setCaseGroupNameError(t('web_case.caseGroupNameDuplicate'));
      return;
    }
    
    setCaseGroupNameError(''); // 清除错误提示

    try {
      const groupData = {
        groupName: trimmedName,
        description: '',
        displayOrder: caseGroups.length
      };

      console.log('[WebLeftSidePanel] 创建用例集:', groupData);
      await createWebCaseGroup(projectId, groupData);

      message.success(t('web_case.createCaseGroupSuccess'));
      setCreateModalVisible(false);
      setNewCaseGroupName('');

      // 刷新列表
      await loadCaseGroups();

      // 触发父组件刷新
      if (onCaseGroupsUpdated) {
        onCaseGroupsUpdated();
      }
    } catch (error) {
      console.error('[WebLeftSidePanel] 创建用例集失败:', error);
      // 检查是否是重名错误（UNIQUE constraint failed）
      const errorMsg = error.response?.data?.error || error.message || '';
      if (error.response?.status === 409 || errorMsg.includes('UNIQUE') || errorMsg.includes('constraint')) {
        setCaseGroupNameError(t('web_case.caseGroupNameDuplicate'));
      } else {
        message.error(t('web_case.createCaseGroupFailed'));
      }
    }
  };

  // 版本保存
  const handleVersionSave = async () => {
    if (!projectId) {
      console.error('[WebLeftSidePanel] projectId为空');
      message.error(t('common.error'));
      return;
    }

    setSavingVersion(true);
    try {
      console.log('[WebLeftSidePanel] 开始保存Web用例版本:', projectId);
      const result = await saveWebVersion(projectId);
      console.log('[WebLeftSidePanel] 版本保存成功:', result);

      message.success(t('web_version.saveSuccess'));

      // 刷新版本列表
      if (versionListRef.current && versionListRef.current.loadVersions) {
        versionListRef.current.loadVersions();
      }
    } catch (error) {
      console.error('[WebLeftSidePanel] 版本保存失败:', error);
      message.error(t('web_version.saveFailed'));
    } finally {
      setSavingVersion(false);
    }
  };

  // 模版下载
  const handleTemplateDownload = async () => {
    try {
      await exportWebTemplate(projectId);
      message.success(t('web_case.templateDownloadSuccess'));
    } catch (error) {
      console.error('[WebLeftSidePanel] 模版下载失败:', error);
      message.error(t('web_case.templateDownloadFailed'));
    }
  };

  // 打开导入对话框
  const handleOpenImportModal = () => {
    setSelectedImportGroup(null);
    setNewImportGroupName('');
    setImportFile(null);
    setImportModalVisible(true);
  };

  // 选择文件
  const handleFileSelect = (e) => {
    const file = e.target.files[0];
    if (file) {
      setImportFile(file);
    }
  };

  // 执行导入
  const handleConfirmImport = async () => {
    if (!importFile) {
      message.warning(t('web_case.selectFileFirst'));
      return;
    }
    if (!selectedImportGroup) {
      message.warning(t('web_case.selectCaseGroupFirst'));
      return;
    }

    // 如果选择新建用例集，检查名称是否输入
    let targetGroup = selectedImportGroup;
    if (selectedImportGroup === '__new__') {
      if (!newImportGroupName || !newImportGroupName.trim()) {
        message.warning(t('web_case.enterCaseGroupName'));
        return;
      }
      targetGroup = newImportGroupName.trim();
    }

    setImportLoading(true);
    try {
      const result = await importWebCases(projectId, targetGroup, importFile);
      message.success(t('web_case.importSuccess', { 
        insertCount: result.insertCount, 
        updateCount: result.updateCount 
      }));
      setImportModalVisible(false);
      setImportFile(null);
      // 刷新用例集列表
      await loadCaseGroups();
      // 延迟切换到导入的用例集，确保列表已刷新
      setTimeout(() => {
        handleCaseSwitch(targetGroup);
        // 触发父组件刷新
        if (onCaseGroupsUpdated) {
          onCaseGroupsUpdated();
        }
      }, 100);
    } catch (error) {
      console.error('[WebLeftSidePanel] 导入失败:', error);
      const errorMsg = error.response?.data?.message || error.message || t('web_case.importFailed');
      message.error(errorMsg);
    } finally {
      setImportLoading(false);
    }
  };

  // 用例集切换
  const handleCaseSwitch = (caseGroup) => {
    console.log('[WebLeftSidePanel] 切换用例集:', caseGroup);
    if (onCaseSwitch) {
      onCaseSwitch(caseGroup);
    }
  };

  // 用例集编辑
  const handleCaseEdit = async (groupId, newName) => {
    if (!newName || newName.trim() === '') {
      message.warning(t('web_case.caseGroupNameRequired'));
      return;
    }

    const trimmedName = newName.trim();

    // 检查重名（排除自己）
    const isDuplicate = caseGroups.some(
      group => group.case_group === trimmedName && group._groupId !== groupId
    );
    if (isDuplicate) {
      message.error(t('web_case.caseGroupNameDuplicate'));
      return;
    }

    try {
      await updateWebCaseGroup(groupId, { group_name: trimmedName });
      message.success(t('web_case.updateCaseGroupSuccess'));

      // 刷新列表
      await loadCaseGroups();

      // 如果当前选中的用例集被重命名，更新选中状态
      const oldName = caseGroups.find(g => g._groupId === groupId)?.case_group;
      if (oldName === selectedCaseGroup) {
        handleCaseSwitch(trimmedName);
      }

      // 触发父组件刷新
      if (onCaseGroupsUpdated) {
        onCaseGroupsUpdated();
      }
    } catch (error) {
      console.error('[WebLeftSidePanel] 更新用例集失败:', error);
      message.error(t('web_case.updateCaseGroupFailed'));
    }
  };

  // 用例集删除
  const handleCaseDelete = async (groupId) => {
    try {
      await deleteWebCaseGroup(groupId);
      message.success(t('web_case.deleteCaseGroupSuccess'));

      // 刷新列表
      await loadCaseGroups();

      // 如果删除的是当前选中的用例集，选中第一个
      const deletedGroup = caseGroups.find(g => g._groupId === groupId);
      if (deletedGroup && deletedGroup.case_group === selectedCaseGroup) {
        const remainingGroups = caseGroups.filter(g => g._groupId !== groupId);
        if (remainingGroups.length > 0) {
          handleCaseSwitch(remainingGroups[0].case_group);
        } else {
          handleCaseSwitch(null);
        }
      }

      // 触发父组件刷新
      if (onCaseGroupsUpdated) {
        onCaseGroupsUpdated();
      }
    } catch (error) {
      console.error('[WebLeftSidePanel] 删除用例集失败:', error);
      message.error(t('web_case.deleteCaseGroupFailed'));
    }
  };

  // 收束/展开切换
  const handleToggleCollapse = () => {
    if (onCollapse) {
      onCollapse(!collapsed);
    }
  };

  if (collapsed) {
    return (
      <div className="web-left-side-panel collapsed">
        <Button
          type="text"
          icon={<RightOutlined />}
          onClick={handleToggleCollapse}
          className="collapse-toggle"
        />
      </div>
    );
  }

  return (
    <div className="web-left-side-panel">
      {/* 功能区1: 创建用例集 - 蓝色按钮 */}
      <div className="function-area create-case-group">
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={handleCreateCaseGroup}
          block
        >
          {t('web_case.createCaseGroup')}
        </Button>
      </div>

      {/* 功能区2: 版本保存 - 白色按钮 */}
      <div className="function-area version-save">
        <Button
          icon={<SaveOutlined />}
          onClick={handleVersionSave}
          loading={savingVersion}
          block
        >
          {t('web_version.save')}
        </Button>
      </div>

      {/* 功能区3: 版本一览 - 白色按钮 */}
      <div className="function-area version-list">
        <Button
          icon={<UnorderedListOutlined />}
          onClick={() => {
            setVersionModalVisible(true);
            // 立即刷新版本列表
            setVersionListKey(prev => prev + 1);
          }}
          block
        >
          {t('web_version.versionList')}
        </Button>
      </div>

      {/* 功能区4: 模版下载 - 白色按钮 */}
      <div className="function-area template-download">
        <Button
          icon={<DownloadOutlined />}
          onClick={handleTemplateDownload}
          block
        >
          {t('web_case.templateDownload')}
        </Button>
      </div>

      {/* 功能区5: 用例导入 - 白色按钮 */}
      <div className="function-area case-import">
        <Button
          icon={<UploadOutlined />}
          onClick={handleOpenImportModal}
          block
        >
          {t('web_case.importCases')}
        </Button>
      </div>

      {/* 功能区6: 用例集一览 */}
      <div className="function-area case-group-list">
        <div className="case-group-list-header">
          <span>{t('web_case.caseGroupList')}</span>
          <Button
            type="text"
            size="small"
            icon={collapsed ? <RightOutlined /> : <LeftOutlined />}
            onClick={handleToggleCollapse}
            style={{ padding: '0 4px' }}
          />
        </div>
        <Spin spinning={casesLoading}>
          <div className="case-group-list-content">
            {caseGroups.length === 0 ? (
              <div className="empty-tip">{t('web_case.noCaseGroups')}</div>
            ) : (
              caseGroups.map((caseGroup) => (
                <CaseListItem
                  key={caseGroup._groupId}
                  caseData={caseGroup}
                  isSelected={selectedCaseGroup && caseGroup._groupId === selectedCaseGroup.id}
                  onSwitch={() => handleCaseSwitch(caseGroup)}
                  onEdit={(newName) => handleCaseEdit(caseGroup._groupId, newName)}
                  onDelete={() => handleCaseDelete(caseGroup._groupId)}
                />
              ))
            )}
          </div>
        </Spin>
      </div>

      {/* 创建用例集Modal */}
      <Modal
        title={t('web_case.createCaseGroup')}
        open={createModalVisible}
        onOk={handleSaveNewCaseGroup}
        onCancel={() => {
          setCreateModalVisible(false);
          setCaseGroupNameError(''); // 关闭时清除错误
        }}
        okText={t('common.confirm')}
        cancelText={t('common.cancel')}
      >
        <Input
          placeholder={t('web_case.enterCaseGroupName')}
          value={newCaseGroupName}
          onChange={(e) => {
            setNewCaseGroupName(e.target.value);
            setCaseGroupNameError(''); // 输入时清除错误
          }}
          onPressEnter={handleSaveNewCaseGroup}
          maxLength={100}
          status={caseGroupNameError ? 'error' : undefined}
        />
        {caseGroupNameError && (
          <div style={{ color: '#ff4d4f', marginTop: 8, fontSize: 12 }}>
            {caseGroupNameError}
          </div>
        )}
      </Modal>

      {/* 版本一览Modal */}
      <Modal
        title={t('web_version.versionList')}
        open={versionModalVisible}
        onCancel={() => setVersionModalVisible(false)}
        footer={null}
        width={1000}
        bodyStyle={{ padding: '16px' }}
      >
        <WebVersionList
          key={versionListKey}
          projectId={projectId}
          onVersionDeleted={() => {
            console.log('[WebLeftSidePanel] 版本已删除');
            setVersionListKey(prev => prev + 1);
          }}
        />
      </Modal>

      {/* 导入用例Modal */}
      <Modal
        title={t('web_case.importCases')}
        open={importModalVisible}
        onOk={handleConfirmImport}
        onCancel={() => setImportModalVisible(false)}
        okText={t('common.confirm')}
        cancelText={t('common.cancel')}
        confirmLoading={importLoading}
        width={600}
      >
        <div style={{ marginBottom: 16 }}>
          <div style={{ marginBottom: 8 }}>{t('web_case.selectExcelFile')}:</div>
          <Button 
            icon={<UploadOutlined />} 
            onClick={() => fileInputRef.current?.click()}
          >
            {t('web_case.selectFile')}
          </Button>
          {importFile && (
            <span style={{ marginLeft: 12, color: '#52c41a' }}>
              {importFile.name}
            </span>
          )}
          <input
            ref={fileInputRef}
            type="file"
            accept=".xlsx,.xls"
            style={{ display: 'none' }}
            onChange={handleFileSelect}
          />
        </div>

        <div style={{ marginBottom: 16 }}>
          <div style={{ marginBottom: 8 }}>{t('web_case.selectTargetCaseGroup')}:</div>
          <Radio.Group 
            value={selectedImportGroup} 
            onChange={(e) => setSelectedImportGroup(e.target.value)}
            style={{ width: '100%' }}
          >
            {caseGroups.map((group) => (
              <Radio 
                key={group._groupId} 
                value={group.case_group}
                style={{ display: 'block', marginBottom: 8 }}
              >
                {group.case_group}
              </Radio>
            ))}
            <Radio value="__new__" style={{ display: 'block', marginTop: 12 }}>
              + {t('web_case.createNewCaseGroup')}
            </Radio>
          </Radio.Group>
          {selectedImportGroup === '__new__' && (
            <Input
              placeholder={t('web_case.enterCaseGroupName')}
              style={{ marginTop: 8 }}
              value={newImportGroupName}
              onChange={(e) => setNewImportGroupName(e.target.value)}
            />
          )}
        </div>

        <div style={{ 
          padding: 12, 
          background: '#e6f7ff', 
          borderLeft: '3px solid #1890ff',
          marginTop: 16 
        }}>
          <div style={{ marginBottom: 4 }}>💡 {t('common.tips')}:</div>
          <div>• {t('web_case.importTip1')}</div>
          <div>• {t('web_case.importTip2')}</div>
          <div>• {t('web_case.importTip3')}</div>
        </div>
      </Modal>
    </div>
  );
};

export default WebLeftSidePanel;
