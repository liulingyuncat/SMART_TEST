import React, { useState, useEffect, useRef } from 'react';
import { Button, message, Spin, Modal, Input, Radio } from 'antd';
import { PlusOutlined, SaveOutlined, RightOutlined, LeftOutlined, UnorderedListOutlined, DownloadOutlined, UploadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import {
  saveApiVersion,
  exportApiTemplate,
  importApiCases
} from '../../../../api/apiCase';
import {
  getApiCaseGroupsFromTable,
  createApiCaseGroupInTable,
  updateApiCaseGroupInTable,
  deleteApiCaseGroupInTable
} from '../../../../api/autoCase';
import CaseListItem from '../../../ProjectDetail/ManualTestTabs/components/CaseListItem';
import WebVersionList from '../../AutoTestTabs/components/WebVersionList';
import './ApiLeftSidePanel.css';

/**
 * API用例库左侧操作面板
 * 功能区：1.创建用例集(蓝色按钮) 2.版本保存(白色按钮) 3.版本一览(白色按钮+弹窗) 
 *        4.模版下载 5.用例导入 6.用例集一览(标题右侧收束控件)
 * 与WebLeftSidePanel的差异：不包含语言筛选相关逻辑
 */
const ApiLeftSidePanel = ({
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
  const [caseGroupNameError, setCaseGroupNameError] = useState('');
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
      console.warn('[ApiLeftSidePanel] projectId为空');
      return;
    }

    console.log('[ApiLeftSidePanel] 🔍 开始加载API用例集列表 - projectId:', projectId);
    setCasesLoading(true);
    try {
      const response = await getApiCaseGroupsFromTable(projectId);
      console.log('[ApiLeftSidePanel] 📦 用例集API完整返回:', response);

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
        console.log('[ApiLeftSidePanel] 🔄 转换后的用例集数据:', formattedGroups);
        
        setCaseGroups(formattedGroups);

        // 默认选中第一个用例集
        if (formattedGroups.length > 0 && onCaseSwitch && !selectedCaseGroup) {
          console.log('[ApiLeftSidePanel] 🎯 默认选中第一个用例集:', formattedGroups[0].case_group);
          setTimeout(() => {
            onCaseSwitch(formattedGroups[0]);
          }, 0);
        }
      } else {
        console.log('[ApiLeftSidePanel] ⚠️ 用例集列表为空');
        setCaseGroups([]);
      }
    } catch (error) {
      console.error('[ApiLeftSidePanel] ❌ 加载用例集列表失败:', error);
      message.error(t('api_case.loadCaseGroupsFailed', { defaultValue: '加载用例集列表失败' }));
    } finally {
      setCasesLoading(false);
    }
  };

  // 打开创建用例集对话框
  const handleCreateCaseGroup = () => {
    setNewCaseGroupName('');
    setCaseGroupNameError('');
    setCreateModalVisible(true);
  };

  // 保存新用例集
  const handleSaveNewCaseGroup = async () => {
    const trimmedName = newCaseGroupName.trim();
    console.log('[ApiLeftSidePanel] 🆕 开始创建用例集 - 名称:', trimmedName);
    console.log('[ApiLeftSidePanel] 🆕 当前projectId:', projectId);
    console.log('[ApiLeftSidePanel] 🆕 当前用例集列表:', caseGroups);
    
    if (!trimmedName) {
      console.warn('[ApiLeftSidePanel] ⚠️ 用例集名称为空');
      setCaseGroupNameError(t('api_case.caseGroupNameRequired', { defaultValue: '请输入用例集名称' }));
      return;
    }

    // 检查重复
    const isDuplicate = caseGroups.some(group => group.case_group === trimmedName);
    if (isDuplicate) {
      console.warn('[ApiLeftSidePanel] ⚠️ 用例集名称重复:', trimmedName);
      setCaseGroupNameError(t('api_case.caseGroupNameDuplicate', { defaultValue: '用例集名称已存在' }));
      return;
    }

    try {
      console.log('[ApiLeftSidePanel] 📤 调用createApiCaseGroupInTable API...');
      const createResponse = await createApiCaseGroupInTable(projectId, { groupName: trimmedName });
      console.log('[ApiLeftSidePanel] ✅ 创建API返回:', createResponse);
      
      message.success(t('api_case.createCaseGroupSuccess', { defaultValue: '用例集创建成功' }));
      setCreateModalVisible(false);
      setNewCaseGroupName('');
      setCaseGroupNameError('');
      
      console.log('[ApiLeftSidePanel] 🔄 关闭Modal，开始重新加载用例集列表...');
      
      // 重新加载用例集列表
      await loadCaseGroups();
      console.log('[ApiLeftSidePanel] ✅ loadCaseGroups执行完成');
      
      if (onCaseGroupsUpdated) {
        console.log('[ApiLeftSidePanel] 🔔 调用onCaseGroupsUpdated回调');
        onCaseGroupsUpdated();
      }
      
      console.log('[ApiLeftSidePanel] 🎉 创建用例集流程全部完成');
    } catch (error) {
      console.error('[ApiLeftSidePanel] ❌ 创建用例集失败:', error);
      console.error('[ApiLeftSidePanel] ❌ 错误状态码:', error.response?.status);
      console.error('[ApiLeftSidePanel] ❌ 错误数据:', error.response?.data);
      
      const errorMessage = error.response?.data?.message || t('api_case.createCaseGroupFailed', { defaultValue: '创建用例集失败' });
      
      // 如果是重复名称错误(409状态码)，显示在输入框下方的红字提示
      if (error.response?.status === 409) {
        console.log('[ApiLeftSidePanel] 🔴 显示红字提示:', errorMessage);
        setCaseGroupNameError(errorMessage);
      } else {
        // 其他错误显示全局提示
        console.log('[ApiLeftSidePanel] 🔴 显示全局错误提示:', errorMessage);
        message.error(errorMessage);
      }
    }
  };

  // 保存版本
  const handleVersionSave = async () => {
    setSavingVersion(true);
    try {
      await saveApiVersion(projectId);
      message.success(t('api_version.saveSuccess', { defaultValue: '版本保存成功' }));
      setVersionListKey(prev => prev + 1); // 刷新版本列表
    } catch (error) {
      console.error('[ApiLeftSidePanel] 保存版本失败:', error);
      message.error(t('api_version.saveFailed', { defaultValue: '保存版本失败' }));
    } finally {
      setSavingVersion(false);
    }
  };

  // 模版下载
  const handleTemplateDownload = async () => {
    try {
      await exportApiTemplate(projectId);
      message.success(t('api_case.templateDownloadSuccess', { defaultValue: '模版下载成功' }));
    } catch (error) {
      console.error('[ApiLeftSidePanel] 模版下载失败:', error);
      message.error(t('api_case.templateDownloadFailed', { defaultValue: '模版下载失败' }));
    }
  };

  // 打开导入用例对话框
  const handleOpenImportModal = () => {
    setImportFile(null);
    setNewImportGroupName('');
    // 默认选中当前选中的用例集，如果没有则选第一个
    setSelectedImportGroup(selectedCaseGroup || (caseGroups.length > 0 ? caseGroups[0].case_group : null));
    setImportModalVisible(true);
  };

  // 选择文件
  const handleFileSelect = (e) => {
    const file = e.target.files[0];
    if (file) {
      setImportFile(file);
    }
    e.target.value = null; // 清空input，允许重复选择同一文件
  };

  // 确认导入
  const handleConfirmImport = async () => {
    if (!importFile) {
      message.warning(t('api_case.selectFile', { defaultValue: '请选择要导入的文件' }));
      return;
    }

    if (!selectedImportGroup) {
      message.warning(t('api_case.selectCaseGroup', { defaultValue: '请选择用例集' }));
      return;
    }

    // 如果选择新建用例集，检查名称是否输入
    let targetGroup = selectedImportGroup;
    if (selectedImportGroup === '__new__') {
      if (!newImportGroupName || !newImportGroupName.trim()) {
        message.warning(t('api_case.enterCaseGroupName', { defaultValue: '请输入用例集名称' }));
        return;
      }
      targetGroup = newImportGroupName.trim();
    }

    setImportLoading(true);
    try {
      const response = await importApiCases(projectId, importFile, targetGroup);
      const { insert_count = 0, update_count = 0 } = response.data || {};
      message.success(
        t('api_case.importSuccess', { 
          defaultValue: `导入成功：新增${insert_count}条，更新${update_count}条`,
          insert_count, 
          update_count 
        })
      );
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
      console.error('[ApiLeftSidePanel] 导入用例失败:', error);
      message.error(error.response?.data?.message || t('api_case.importFailed', { defaultValue: '导入失败' }));
    } finally {
      setImportLoading(false);
    }
  };

  // 切换用例集
  const handleCaseSwitch = (caseGroup) => {
    console.log('[ApiLeftSidePanel] 切换用例集:', caseGroup);
    if (onCaseSwitch) {
      onCaseSwitch(caseGroup);
    }
  };

  // 编辑用例集名称
  const handleCaseEdit = async (groupId, newName) => {
    if (!newName || !newName.trim()) {
      message.error(t('api_case.caseGroupNameRequired', { defaultValue: '请输入用例集名称' }));
      return;
    }

    const trimmedName = newName.trim();

    // 检查重复（排除自己）
    const isDuplicate = caseGroups.some(
      group => group.case_group === trimmedName && group._groupId !== groupId
    );
    if (isDuplicate) {
      message.error(t('api_case.caseGroupNameDuplicate', { defaultValue: '用例集名称已存在' }));
      return;
    }

    try {
      await updateApiCaseGroupInTable(groupId, { group_name: trimmedName });
      message.success(t('api_case.updateCaseGroupSuccess', { defaultValue: '用例集名称更新成功' }));
      
      // 重新加载用例集列表
      await loadCaseGroups();
      
      // 如果修改的是当前选中的用例集，需要更新选中状态
      if (selectedCaseGroup === groupId) {
        onCaseSwitch(trimmedName);
      }
      
      if (onCaseGroupsUpdated) {
        onCaseGroupsUpdated();
      }
    } catch (error) {
      console.error('[ApiLeftSidePanel] 更新用例集失败:', error);
      message.error(error.response?.data?.message || t('api_case.updateCaseGroupFailed', { defaultValue: '更新用例集失败' }));
    }
  };

  // 删除用例集
  const handleCaseDelete = async (groupId) => {
    try {
      await deleteApiCaseGroupInTable(groupId);
      message.success(t('api_case.deleteCaseGroupSuccess', { defaultValue: '用例集删除成功' }));
      
      // 如果删除的是当前选中的用例集，清空选中状态
      if (selectedCaseGroup === groupId) {
        onCaseSwitch(null);
      }
      
      // 重新加载用例集列表
      await loadCaseGroups();
      if (onCaseGroupsUpdated) {
        onCaseGroupsUpdated();
      }
    } catch (error) {
      console.error('[ApiLeftSidePanel] 删除用例集失败:', error);
      message.error(error.response?.data?.message || t('api_case.deleteCaseGroupFailed', { defaultValue: '删除用例集失败' }));
    }
  };

  // 切换收起/展开
  const handleToggleCollapse = () => {
    if (onCollapse) {
      onCollapse();
    }
  };

  if (collapsed) {
    return (
      <div className="api-left-side-panel collapsed">
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
    <div className="api-left-side-panel">
      {/* 功能区1: 创建API用例集 - 蓝色按钮 */}
      <div className="function-area create-case-group">
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={handleCreateCaseGroup}
          block
        >
          {t('api_case.createCaseGroup', { defaultValue: '创建API用例集' })}
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
          {t('api_version.save', { defaultValue: '版本保存' })}
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
          {t('api_version.versionList', { defaultValue: '版本一览' })}
        </Button>
      </div>

      {/* 功能区4: 模版下载 - 白色按钮 */}
      <div className="function-area template-download">
        <Button
          icon={<DownloadOutlined />}
          onClick={handleTemplateDownload}
          block
        >
          {t('api_case.templateDownload', { defaultValue: '模版下载' })}
        </Button>
      </div>

      {/* 功能区5: 用例导入 - 白色按钮 */}
      <div className="function-area case-import">
        <Button
          icon={<UploadOutlined />}
          onClick={handleOpenImportModal}
          block
        >
          {t('api_case.importCases', { defaultValue: '用例导入' })}
        </Button>
      </div>

      {/* 功能区6: 用例集一览 */}
      <div className="function-area case-group-list">
        <div className="case-group-list-header">
          <span>{t('api_case.caseGroupList', { defaultValue: '用例集一览' })}</span>
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
              <div className="empty-tip">{t('api_case.noCaseGroups', { defaultValue: '暂无用例集' })}</div>
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
        title={t('api_case.createCaseGroup', { defaultValue: '创建API用例集' })}
        open={createModalVisible}
        onOk={handleSaveNewCaseGroup}
        onCancel={() => {
          setCreateModalVisible(false);
          setCaseGroupNameError('');
        }}
        okText={t('common.confirm', { defaultValue: '确定' })}
        cancelText={t('common.cancel', { defaultValue: '取消' })}
      >
        <Input
          placeholder={t('api_case.enterCaseGroupName', { defaultValue: '请输入用例集名称' })}
          value={newCaseGroupName}
          onChange={(e) => {
            setNewCaseGroupName(e.target.value);
            setCaseGroupNameError('');
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
        title={t('api_version.versionList', { defaultValue: '版本一览' })}
        open={versionModalVisible}
        onCancel={() => setVersionModalVisible(false)}
        footer={null}
        width={1000}
        bodyStyle={{ padding: '16px' }}
      >
        <WebVersionList
          key={versionListKey}
          projectId={projectId}
          apiModule="api-cases"
          onVersionDeleted={() => {
            console.log('[ApiLeftSidePanel] 版本已删除');
            setVersionListKey(prev => prev + 1);
          }}
        />
      </Modal>

      {/* 导入用例Modal */}
      <Modal
        title={t('api_case.importCases', { defaultValue: '用例导入' })}
        open={importModalVisible}
        onOk={handleConfirmImport}
        onCancel={() => setImportModalVisible(false)}
        okText={t('common.confirm', { defaultValue: '确定' })}
        cancelText={t('common.cancel', { defaultValue: '取消' })}
        confirmLoading={importLoading}
        width={600}
      >
        <div style={{ marginBottom: 16 }}>
          <div style={{ marginBottom: 8 }}>{t('api_case.selectExcelFile', { defaultValue: '选择Excel文件' })}:</div>
          <Button 
            icon={<UploadOutlined />} 
            onClick={() => fileInputRef.current?.click()}
          >
            {t('api_case.selectFile', { defaultValue: '选择文件' })}
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
          <div style={{ marginBottom: 8 }}>{t('api_case.selectTargetCaseGroup', { defaultValue: '选择目标用例集' })}:</div>
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
              + {t('api_case.createNewCaseGroup', { defaultValue: '新建用例集' })}
            </Radio>
          </Radio.Group>
          {selectedImportGroup === '__new__' && (
            <Input
              placeholder={t('api_case.enterCaseGroupName', { defaultValue: '请输入用例集名称' })}
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
          <div style={{ marginBottom: 4 }}>💡 {t('common.tips', { defaultValue: '提示' })}:</div>
          <div>• {t('api_case.importTip1', { defaultValue: 'UUID为空的行将创建新用例' })}</div>
          <div>• {t('api_case.importTip2', { defaultValue: 'UUID存在的行将更新现有用例' })}</div>
          <div>• {t('api_case.importTip3', { defaultValue: '请确保Excel格式与模版一致' })}</div>
        </div>
      </Modal>
    </div>
  );
};

export default ApiLeftSidePanel;
