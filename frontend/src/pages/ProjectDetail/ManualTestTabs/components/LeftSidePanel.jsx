import React, { useState, useEffect } from 'react';
import { Button, message, Spin, Modal, Input } from 'antd';
import { PlusOutlined, SaveOutlined, DownloadOutlined, LeftOutlined, RightOutlined, UploadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { createCase, getCasesList, saveMultiLangVersion, exportMultiLangTemplate, importCasesByLanguage, getCaseGroups, createCaseGroup } from '../../../../api/manualCase';
import CaseListItem from './CaseListItem';
import ImportCaseModal from './ImportCaseModal';
import './LeftSidePanel.css';

/**
 * 左侧操作面板组件
 * 包含4个功能区：用例创建、版本保存、模版导出、用例一览
 */
const LeftSidePanel = ({ 
  projectId, 
  language, 
  collapsed = false,
  selectedCaseGroup,
  onCaseCreated, 
  onCaseUpdated,
  onCaseDeleted,
  onCaseSwitch,
  onCollapse 
}) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [cases, setCases] = useState([]);
  const [casesLoading, setCasesLoading] = useState(false);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [newCaseName, setNewCaseName] = useState('');
  const [importModalVisible, setImportModalVisible] = useState(false);
  const [importLoading, setImportLoading] = useState(false);

  // 获取语言字段后缀
  const getLanguageSuffix = (lang) => {
    const map = {
      '中文': 'cn',
      'English': 'en',
      '日本語': 'jp'
    };
    return map[lang] || 'cn';
  };

  // 加载用例一览列表（从 case_groups 表获取）
  const loadCases = async () => {
    if (!projectId) {
      console.warn('[LeftSidePanel] loadCases: projectId为空');
      return;
    }
    
    console.log('[LeftSidePanel] 开始加载用例集列表:', { projectId, language });
    setCasesLoading(true);
    try {
      // 先从case_groups表获取用例集列表
      const groupsResponse = await getCaseGroups(projectId, 'overall');
      console.log('[LeftSidePanel] 用例集API返回:', groupsResponse);
      
      if (groupsResponse && groupsResponse.length > 0) {
        // 有独立的用例集记录，直接使用
        console.log('[LeftSidePanel] 从case_groups表加载到', groupsResponse.length, '个用例集');
        // 转换数据格式以兼容现有组件
        const formattedGroups = groupsResponse.map(group => ({
          case_group: group.group_name,
          id: group.id,
          // 保留其他可能需要的字段
          _groupId: group.id
        }));
        setCases(formattedGroups);
        
        // 默认选中第一个用例集
        if (formattedGroups.length > 0 && onCaseSwitch && !selectedCaseGroup) {
          console.log('[LeftSidePanel] 默认选中第一个用例集:', formattedGroups[0].case_group);
          setTimeout(() => {
            onCaseSwitch(formattedGroups[0].case_group);
          }, 0);
        }
      } else {
        // case_groups表为空，尝试从manual_test_cases中提取并创建
        console.log('[LeftSidePanel] case_groups表为空，从用例中提取');
        const casesResponse = await getCasesList(projectId, { 
          caseType: 'overall',
          language: language,
          page: 1,
          size: 1000
        });
        
        if (casesResponse && casesResponse.cases) {
          // 提取不重复的case_group
          const caseGroupNames = new Set();
          casesResponse.cases.forEach(caseItem => {
            if (caseItem.case_group) {
              caseGroupNames.add(caseItem.case_group);
            }
          });
          
          console.log('[LeftSidePanel] 从用例中提取到', caseGroupNames.size, '个不重复的用例集');
          
          // 为每个用例集创建记录
          const createdGroups = [];
          for (const groupName of caseGroupNames) {
            try {
              const newGroup = await createCaseGroup(projectId, {
                caseType: 'overall',
                groupName: groupName
              });
              console.log('[LeftSidePanel] 创建用例集记录:', groupName);
              createdGroups.push({
                case_group: groupName,
                id: newGroup.id,
                _groupId: newGroup.id
              });
            } catch (err) {
              console.error('[LeftSidePanel] 创建用例集失败:', groupName, err);
            }
          }
          
          setCases(createdGroups);
          
          if (createdGroups.length > 0 && onCaseSwitch && !selectedCaseGroup) {
            console.log('[LeftSidePanel] 默认选中第一个用例集:', createdGroups[0].case_group);
            setTimeout(() => {
              onCaseSwitch(createdGroups[0].case_group);
            }, 0);
          }
        } else {
          console.log('[LeftSidePanel] 没有返回用例数据');
          setCases([]);
          // 清空选中状态
          if (onCaseSwitch) {
            onCaseSwitch(null);
          }
        }
      }
    } catch (error) {
      console.error('[LeftSidePanel] 加载用例列表失败:', error);
      message.error(t('manualTest.loadCasesFailed'));
      setCases([]);
    } finally {
      setCasesLoading(false);
    }
  };

  // 初始加载和语言变更时重新加载
  useEffect(() => {
    loadCases();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [projectId, language]);

  // 显示创建用例对话框
  const handleShowCreateModal = () => {
    setNewCaseName('');
    setCreateModalVisible(true);
  };

  // 创建用例集
  const handleCreateCase = async () => {
    console.log('[LeftSidePanel] 开始创建用例集:', { projectId, newCaseName });
    
    if (!projectId) {
      console.error('[LeftSidePanel] projectId为空');
      message.error(t('manualTest.projectIdRequired'));
      return;
    }

    if (!newCaseName.trim()) {
      console.warn('[LeftSidePanel] 用例集名称为空');
      message.error(t('manualTest.caseNameRequired'));
      return;
    }

    // 检查用例集名称是否重复
    const trimmedName = newCaseName.trim();
    const isDuplicate = cases.some(c => c.case_group === trimmedName);
    if (isDuplicate) {
      message.error('用例集名称已存在，请使用不同的名称');
      return;
    }

    setLoading(true);
    try {
      // 直接创建用例集记录
      const groupData = {
        caseType: 'overall',
        groupName: newCaseName.trim()
      };

      console.log('[LeftSidePanel] 创建用例集数据:', groupData);
      console.log('[LeftSidePanel] 开始调用createCaseGroup API');
      
      const result = await createCaseGroup(projectId, groupData);
      console.log('[LeftSidePanel] 创建用例集API返回:', result);

      message.success(t('manualTest.createCaseSuccess'));
      setCreateModalVisible(false);
      setNewCaseName('');
      
      console.log('[LeftSidePanel] 开始刷新用例集列表');
      // 短暂延迟后刷新，确保后端数据已保存
      setTimeout(async () => {
        await loadCases();
        if (onCaseCreated) {
          onCaseCreated();
        }
      }, 300);
    } catch (error) {
      console.error('[LeftSidePanel] 创建用例集失败:', error);
      message.error(t('manualTest.createCaseFailed'));
    } finally {
      setLoading(false);
    }
  };

  // 版本保存（生成CN/JP/EN三语言版本）
  const [saveVersionLoading, setSaveVersionLoading] = useState(false);
  const handleSaveVersion = async () => {
    if (!projectId) {
      message.error(t('manualTest.projectIdRequired'));
      return;
    }

    setSaveVersionLoading(true);
    try {
      const response = await saveMultiLangVersion(projectId);
      if (response && response.filename) {
        message.success(`${t('manualTest.saveVersionSuccess')}: ${response.filename}`);
      } else {
        message.success(t('manualTest.saveVersionSuccess'));
      }
    } catch (error) {
      console.error('[LeftSidePanel] 版本保存失败:', error);
      message.error(t('manualTest.saveVersionFailed'));
    } finally {
      setSaveVersionLoading(false);
    }
  };

  // 模版导出（CN/JP/EN空白xlsx打包成zip）
  const [exportTemplateLoading, setExportTemplateLoading] = useState(false);
  const handleExportTemplate = async () => {
    setExportTemplateLoading(true);
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
      console.error('[LeftSidePanel] 模版导出失败:', error);
      message.error(t('manualTest.exportTemplateFailed'));
    } finally {
      setExportTemplateLoading(false);
    }
  };

  // 导入用例
  const handleImport = async (file, targetCaseGroup) => {
    console.log('[LeftSidePanel] 开始导入:', { file: file.name, targetCaseGroup, language, projectId });
    
    setImportLoading(true);
    try {
      const result = await importCasesByLanguage(projectId, 'overall', file, language, targetCaseGroup);
      console.log('[LeftSidePanel] 导入成功:', result);
      message.success(`导入成功: 更新${result.updateCount}条, 新增${result.insertCount}条`);
      
      // 关闭对话框
      setImportModalVisible(false);
      
      // 刷新用例列表
      setTimeout(async () => {
        await loadCases();
        if (onCaseCreated) {
          onCaseCreated();
        }
      }, 300);
    } catch (error) {
      console.error('[LeftSidePanel] 导入失败:', error);
      message.error(`导入失败: ${error.message || '未知错误'}`);
    } finally {
      setImportLoading(false);
    }
  };

  // 用例更新回调
  const handleCaseUpdate = (caseId, newName) => {
    console.log('[LeftSidePanel] 用例集名称已更新:', { caseId, newName, currentSelected: selectedCaseGroup });
    
    // 如果修改的是当前选中的用例集，需要通知父组件更新selectedCaseGroup
    const oldCaseGroup = cases.find(c => (c._groupId || c.id) === caseId)?.case_group;
    if (oldCaseGroup === selectedCaseGroup) {
      console.log('[LeftSidePanel] 修改的是当前选中的用例集，通知父组件切换到新名称');
      if (onCaseSwitch) {
        onCaseSwitch(newName);
      }
    }
    
    // 刷新列表
    loadCases();
    if (onCaseUpdated) {
      onCaseUpdated(caseId, newName);
    }
  };

  // 用例删除回调
  const handleCaseDelete = (caseId) => {
    console.log('[LeftSidePanel] 收到删除回调:', caseId);
    // 刷新列表
    loadCases();
    if (onCaseDeleted) {
      onCaseDeleted(caseId);
    }
  };

  // 切换用例回调
  const handleCaseSwitch = (caseData) => {
    console.log('[LeftSidePanel] 切换用例集:', caseData.case_group);
    message.info(`${t('manualTest.switchToCase')}: ${caseData.case_group || t('manualTest.untitledCase')}`);
    if (onCaseSwitch) {
      // 传递用例集名称，用于表格过滤
      onCaseSwitch(caseData.case_group);
    }
  };

  // 收束/展开切换
  const handleToggleCollapse = () => {
    const newCollapsed = !collapsed;
    
    // 持久化到 sessionStorage
    try {
      sessionStorage.setItem(`left_panel_collapsed_${projectId}`, JSON.stringify(newCollapsed));
    } catch (e) {
      console.warn('Failed to save collapse state:', e);
    }
    
    if (onCollapse) {
      onCollapse(newCollapsed);
    }
  };

  // 从 sessionStorage 恢复收束状态
  useEffect(() => {
    if (!projectId) return;
    
    try {
      const saved = sessionStorage.getItem(`left_panel_collapsed_${projectId}`);
      if (saved !== null) {
        const savedCollapsed = JSON.parse(saved);
        if (savedCollapsed !== collapsed && onCollapse) {
          onCollapse(savedCollapsed);
        }
      }
    } catch (e) {
      console.warn('Failed to load collapse state:', e);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [projectId]);

  if (collapsed) {
    return (
      <div className="left-side-panel collapsed">
        <Button 
          type="text" 
          icon={<RightOutlined />} 
          onClick={handleToggleCollapse}
          className="collapse-toggle-collapsed"
          title={t('manualTest.expand')}
        />
      </div>
    );
  }

  return (
    <div className="left-side-panel">
      {/* 收束按钮 */}
      <Button 
        type="text" 
        icon={<LeftOutlined />} 
        onClick={handleToggleCollapse}
        className="collapse-toggle"
      />

      {/* 功能区1: 用例创建 */}
      <div className="panel-section">
        <h3 className="section-title">{t('manualTest.caseOperations')}</h3>
        <Button 
          type="primary" 
          icon={<PlusOutlined />}
          onClick={handleShowCreateModal}
          block
        >
          {t('manualTest.createCase')}
        </Button>
      </div>

      {/* 创建用例对话框 */}
      <Modal
        title={t('manualTest.createCase')}
        open={createModalVisible}
        onOk={handleCreateCase}
        onCancel={() => {
          setCreateModalVisible(false);
          setNewCaseName('');
        }}
        confirmLoading={loading}
        okText={t('common.confirm')}
        cancelText={t('common.cancel')}
      >
        <Input
          placeholder={t('manualTest.enterCaseName')}
          value={newCaseName}
          onChange={(e) => setNewCaseName(e.target.value)}
          onPressEnter={handleCreateCase}
          maxLength={100}
          status={cases.some(c => c.case_group === newCaseName.trim()) ? 'error' : ''}
        />
        {cases.some(c => c.case_group === newCaseName.trim()) && (
          <div style={{ color: '#ff4d4f', marginTop: '8px', fontSize: '14px' }}>
            用例集名称已存在，请使用不同的名称
          </div>
        )}
      </Modal>

      {/* 功能区2: 版本保存 */}
      <div className="panel-section">
        <Button 
          icon={<SaveOutlined />}
          onClick={handleSaveVersion}
          loading={saveVersionLoading}
          block
        >
          {t('manualTest.saveVersion')}
        </Button>
      </div>

      {/* 功能区3: 模版导出 */}
      <div className="panel-section">
        <Button 
          icon={<DownloadOutlined />}
          onClick={handleExportTemplate}
          loading={exportTemplateLoading}
          disabled={exportTemplateLoading}
          block
        >
          {t('manualTest.exportTemplate')}
        </Button>
      </div>

      {/* 功能区3.5: 导入用例 */}
      <div className="panel-section">
        <Button 
          icon={<UploadOutlined />}
          onClick={() => setImportModalVisible(true)}
          block
        >
          {t('manualTest.importCases')}
        </Button>
      </div>

      {/* 导入用例对话框 */}
      <ImportCaseModal
        visible={importModalVisible}
        onCancel={() => setImportModalVisible(false)}
        onImport={handleImport}
        caseGroups={cases.map(c => c.case_group)}
        loading={importLoading}
      />

      {/* 功能区4: 用例一览 */}
      <div className="panel-section case-list-section">
        <h3 className="section-title">{t('manualTest.caseOverview')}</h3>
        <Spin spinning={casesLoading}>
          <div className="case-list">
            {cases.length === 0 ? (
              <div className="empty-hint">{t('manualTest.noCases')}</div>
            ) : (
              cases.map(caseItem => (
                <CaseListItem
                  key={caseItem.id || caseItem.case_id}
                  projectId={projectId}
                  caseData={caseItem}
                  language={language}
                  isSelected={caseItem.case_group === selectedCaseGroup}
                  onUpdate={handleCaseUpdate}
                  onDelete={handleCaseDelete}
                  onSwitch={handleCaseSwitch}
                />
              ))
            )}
          </div>
        </Spin>
      </div>
    </div>
  );
};

export default LeftSidePanel;
