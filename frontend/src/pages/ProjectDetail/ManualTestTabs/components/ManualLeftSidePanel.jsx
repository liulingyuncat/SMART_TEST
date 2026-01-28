import React, { useState, useEffect, useRef } from 'react';
import { Button, message, Spin, Modal, Input, Radio } from 'antd';
import { PlusOutlined, SaveOutlined, RightOutlined, LeftOutlined, UnorderedListOutlined, DownloadOutlined, UploadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import {
    getCaseGroups,
    createCaseGroup,
    updateCaseGroup,
    deleteCaseGroup,
    saveMultiLangVersion,
    exportMultiLangTemplate,
    importCasesByLanguage
} from '../../../../api/manualCase';
import { getVersionList, downloadVersion, deleteVersion, updateVersionRemark } from '../../../../api/manualCase';
import CaseListItem from './CaseListItem';
import ManualVersionList from './ManualVersionList';
import './ManualLeftSidePanel.css';

/**
 * 手工用例库左侧操作面板
 * 功能区：1.创建手工用例集(蓝色按钮) 2.版本保存(白色按钮) 3.版本一览(白色按钮+弹窗) 4.模版下载 5.导入用例 6.用例集一览(标题右侧收束控件)
 */
const ManualLeftSidePanel = ({
    projectId,
    language = '中文',
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
            console.warn('[ManualLeftSidePanel] projectId为空');
            return;
        }

        console.log('[ManualLeftSidePanel] 开始加载手工用例集列表:', projectId);
        setCasesLoading(true);
        try {
            const response = await getCaseGroups(projectId, 'overall');
            console.log('[ManualLeftSidePanel] 用例集API返回:', response);

            if (response && response.length > 0) {
                // 转换数据格式以兼容CaseListItem组件
                const formattedGroups = response.map(group => ({
                    case_group: group.group_name,
                    id: group.id,
                    _groupId: group.id,
                }));
                setCaseGroups(formattedGroups);

                // 默认选中第一个用例集
                if (formattedGroups.length > 0 && onCaseSwitch && !selectedCaseGroup) {
                    console.log('[ManualLeftSidePanel] 默认选中第一个用例集:', formattedGroups[0].case_group);
                    setTimeout(() => {
                        onCaseSwitch(formattedGroups[0]);
                    }, 0);
                }
            } else {
                console.log('[ManualLeftSidePanel] 用例集列表为空');
                setCaseGroups([]);
            }
        } catch (error) {
            console.error('[ManualLeftSidePanel] 加载用例集列表失败:', error);
            message.error(t('manualTest.loadCasesFailed'));
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
        if (!projectId) {
            console.error('[ManualLeftSidePanel] projectId为空');
            message.error(t('common.error'));
            return;
        }

        if (!newCaseGroupName || newCaseGroupName.trim() === '') {
            setCaseGroupNameError(t('manualTest.caseNameRequired'));
            return;
        }

        const trimmedName = newCaseGroupName.trim();

        // 检查重名
        const isDuplicate = caseGroups.some(group => group.case_group === trimmedName);
        if (isDuplicate) {
            setCaseGroupNameError(t('manualTest.caseGroupNameDuplicate'));
            return;
        }

        setCaseGroupNameError('');

        try {
            const groupData = {
                caseType: 'overall',
                groupName: trimmedName,
                description: '',
                displayOrder: caseGroups.length
            };

            console.log('[ManualLeftSidePanel] 创建用例集:', groupData);
            await createCaseGroup(projectId, groupData);

            message.success(t('manualTest.createCaseSuccess'));
            setCreateModalVisible(false);
            setNewCaseGroupName('');

            // 刷新列表
            await loadCaseGroups();

            // 触发父组件刷新
            if (onCaseGroupsUpdated) {
                onCaseGroupsUpdated();
            }
        } catch (error) {
            console.error('[ManualLeftSidePanel] 创建用例集失败:', error);
            const errorMsg = error.response?.data?.error || error.message || '';
            if (error.response?.status === 409 || errorMsg.includes('UNIQUE') || errorMsg.includes('constraint')) {
                setCaseGroupNameError(t('manualTest.caseGroupNameDuplicate'));
            } else {
                message.error(t('manualTest.createCaseFailed'));
            }
        }
    };

    // 版本保存
    const handleVersionSave = async () => {
        if (!projectId) {
            console.error('[ManualLeftSidePanel] projectId为空');
            message.error(t('common.error'));
            return;
        }

        setSavingVersion(true);
        try {
            console.log('[ManualLeftSidePanel] 开始保存手工用例版本:', projectId);
            const result = await saveMultiLangVersion(projectId);
            console.log('[ManualLeftSidePanel] 版本保存成功:', result);

            message.success(t('manualTest.saveVersionSuccess'));

            // 刷新版本列表
            setVersionListKey(prev => prev + 1);
        } catch (error) {
            console.error('[ManualLeftSidePanel] 版本保存失败:', error);
            message.error(t('manualTest.saveVersionFailed'));
        } finally {
            setSavingVersion(false);
        }
    };

    // 模版下载
    const handleTemplateDownload = async () => {
        try {
            const blob = await exportMultiLangTemplate();

            // 创建下载链接
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `manual_case_template_${new Date().getTime()}.zip`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);

            message.success(t('manualTest.exportTemplateSuccess'));
        } catch (error) {
            console.error('[ManualLeftSidePanel] 模版下载失败:', error);
            message.error(t('manualTest.exportTemplateFailed'));
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
            message.warning(t('manualTest.selectFileFirst'));
            return;
        }
        if (!selectedImportGroup) {
            message.warning(t('manualTest.selectCaseGroupFirst'));
            return;
        }

        // 如果选择新建用例集，检查名称是否输入
        let targetGroup = selectedImportGroup;
        if (selectedImportGroup === '__new__') {
            if (!newImportGroupName || !newImportGroupName.trim()) {
                message.warning(t('manualTest.enterCaseName'));
                return;
            }
            targetGroup = newImportGroupName.trim();
        }

        setImportLoading(true);
        try {
            const result = await importCasesByLanguage(projectId, 'overall', importFile, language, targetGroup);
            message.success(t('manualTest.importSuccess', {
                insertCount: result.insertCount,
                updateCount: result.updateCount
            }));
            setImportModalVisible(false);
            setImportFile(null);
            // 刷新用例集列表
            await loadCaseGroups();
            // 延迟切换到导入的用例集，确保列表已刷新
            setTimeout(() => {
                const targetCaseGroup = caseGroups.find(g => g.case_group === targetGroup);
                if (targetCaseGroup) {
                    handleCaseSwitch(targetCaseGroup);
                }
                // 触发父组件刷新
                if (onCaseGroupsUpdated) {
                    onCaseGroupsUpdated();
                }
            }, 100);
        } catch (error) {
            console.error('[ManualLeftSidePanel] 导入失败:', error);
            const errorMsg = error.response?.data?.message || error.message || t('manualTest.importFailed');
            message.error(errorMsg);
        } finally {
            setImportLoading(false);
        }
    };

    // 用例集切换
    const handleCaseSwitch = (caseGroup) => {
        console.log('[ManualLeftSidePanel] 切换用例集:', caseGroup);
        if (onCaseSwitch) {
            onCaseSwitch(caseGroup);
        }
    };

    // 用例集编辑
    const handleCaseEdit = async (groupId, newName) => {
        if (!newName || newName.trim() === '') {
            message.warning(t('manualTest.caseNameRequired'));
            return;
        }

        const trimmedName = newName.trim();

        // 检查重名（排除自己）
        const isDuplicate = caseGroups.some(
            group => group.case_group === trimmedName && group._groupId !== groupId
        );
        if (isDuplicate) {
            message.error(t('manualTest.caseGroupNameDuplicate'));
            return;
        }

        try {
            await updateCaseGroup(groupId, { groupName: trimmedName });
            message.success(t('manualTest.updateCaseSuccess'));

            // 刷新列表
            await loadCaseGroups();

            // 如果当前选中的用例集被重命名，更新选中状态
            const oldGroup = caseGroups.find(g => g._groupId === groupId);
            if (oldGroup && selectedCaseGroup && oldGroup._groupId === selectedCaseGroup.id) {
                // 找到更新后的用例集并切换
                setTimeout(() => {
                    const updatedGroup = { ...selectedCaseGroup, case_group: trimmedName };
                    handleCaseSwitch(updatedGroup);
                }, 0);
            }

            // 触发父组件刷新
            if (onCaseGroupsUpdated) {
                onCaseGroupsUpdated();
            }
        } catch (error) {
            console.error('[ManualLeftSidePanel] 更新用例集失败:', error);
            message.error(t('manualTest.updateCaseFailed'));
        }
    };

    // 用例集删除
    const handleCaseDelete = async (groupId) => {
        try {
            await deleteCaseGroup(groupId);
            message.success(t('message.deleteSuccess'));

            // 刷新列表
            await loadCaseGroups();

            // 如果删除的是当前选中的用例集，选中第一个
            if (selectedCaseGroup && selectedCaseGroup.id === groupId) {
                const remainingGroups = caseGroups.filter(g => g._groupId !== groupId);
                if (remainingGroups.length > 0) {
                    handleCaseSwitch(remainingGroups[0]);
                } else {
                    handleCaseSwitch(null);
                }
            }

            // 触发父组件刷新
            if (onCaseGroupsUpdated) {
                onCaseGroupsUpdated();
            }
        } catch (error) {
            console.error('[ManualLeftSidePanel] 删除用例集失败:', error);
            message.error(t('message.deleteFailed'));
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
            <div className="manual-left-side-panel collapsed">
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
        <div className="manual-left-side-panel">
            {/* 功能区1: 创建手工用例集 - 蓝色按钮 */}
            <div className="function-area create-case-group">
                <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={handleCreateCaseGroup}
                    block
                >
                    {t('manualTest.createManualCaseGroup')}
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
                    {t('manualTest.saveVersion')}
                </Button>
            </div>

            {/* 功能区3: 版本一览 - 白色按钮 */}
            <div className="function-area version-list">
                <Button
                    icon={<UnorderedListOutlined />}
                    onClick={() => {
                        setVersionModalVisible(true);
                        setVersionListKey(prev => prev + 1);
                    }}
                    block
                >
                    {t('manualTest.versionList')}
                </Button>
            </div>

            {/* 功能区4: 模版下载 - 白色按钮 */}
            <div className="function-area template-download">
                <Button
                    icon={<DownloadOutlined />}
                    onClick={handleTemplateDownload}
                    block
                >
                    {t('manualTest.templateDownload')}
                </Button>
            </div>

            {/* 功能区5: 用例导入 - 白色按钮 */}
            <div className="function-area case-import">
                <Button
                    icon={<UploadOutlined />}
                    onClick={handleOpenImportModal}
                    block
                >
                    {t('manualTest.importCases')}
                </Button>
            </div>

            {/* 功能区6: 用例集一览 */}
            <div className="function-area case-group-list">
                <div className="case-group-list-header">
                    <span>{t('manualTest.caseOverview')}</span>
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
                            <div className="empty-tip">{t('manualTest.noCases')}</div>
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
                title={t('manualTest.createManualCaseGroup')}
                open={createModalVisible}
                onOk={handleSaveNewCaseGroup}
                onCancel={() => {
                    setCreateModalVisible(false);
                    setCaseGroupNameError('');
                }}
                okText={t('common.confirm')}
                cancelText={t('common.cancel')}
            >
                <Input
                    placeholder={t('manualTest.enterCaseName')}
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
                title={t('manualTest.versionList')}
                open={versionModalVisible}
                onCancel={() => setVersionModalVisible(false)}
                footer={null}
                width={1000}
                bodyStyle={{ padding: '16px' }}
            >
                <ManualVersionList
                    key={versionListKey}
                    projectId={projectId}
                    onVersionDeleted={() => {
                        console.log('[ManualLeftSidePanel] 版本已删除');
                        setVersionListKey(prev => prev + 1);
                    }}
                />
            </Modal>

            {/* 导入用例Modal */}
            <Modal
                title={t('manualTest.importCases')}
                open={importModalVisible}
                onOk={handleConfirmImport}
                onCancel={() => setImportModalVisible(false)}
                okText={t('common.confirm')}
                cancelText={t('common.cancel')}
                confirmLoading={importLoading}
                width={600}
            >
                <div style={{ marginBottom: 16 }}>
                    <div style={{ marginBottom: 8 }}>{t('manualTest.selectExcelFile')}:</div>
                    <Button
                        icon={<UploadOutlined />}
                        onClick={() => fileInputRef.current?.click()}
                    >
                        {t('manualTest.selectFile')}
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
                    <div style={{ marginBottom: 8 }}>{t('manualTest.selectTargetCaseGroup')}:</div>
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
                            + {t('manualTest.createNewCaseGroup')}
                        </Radio>
                    </Radio.Group>
                    {selectedImportGroup === '__new__' && (
                        <Input
                            placeholder={t('manualTest.enterCaseName')}
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
                    <div>• {t('manualTest.importTip1')}</div>
                    <div>• {t('manualTest.importTip2')}</div>
                    <div>• {t('manualTest.importTip3')}</div>
                </div>
            </Modal>
        </div>
    );
};

export default ManualLeftSidePanel;
