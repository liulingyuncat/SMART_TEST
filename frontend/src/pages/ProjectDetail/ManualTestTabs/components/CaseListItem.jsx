import React, { useState } from 'react';
import { Input, message, Modal, App } from 'antd';
import { EditOutlined, DeleteOutlined, CheckOutlined, ExclamationCircleOutlined, CopyOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { updateCase, deleteCase, updateCaseGroup, deleteCaseGroup } from '../../../../api/manualCase';
import './CaseListItem.css';

/**
 * 用例一览列表项组件
 * 支持内联编辑和删除用例集
 */
const CaseListItem = ({ projectId, caseData, language, isSelected = false, onUpdate, onEdit, onDelete, onSwitch }) => {
  const { t } = useTranslation();
  const { modal } = App.useApp();
  const [modalVisible, setModalVisible] = useState(false);
  const [editValue, setEditValue] = useState('');
  const [saving, setSaving] = useState(false);

  // 调试日志
  console.log('[CaseListItem] Render:', { 
    projectId, 
    caseId: caseData?.case_id, 
    caseGroup: caseData?.case_group,
    hasOnDelete: !!onDelete 
  });

  // 获取语言字段后缀
  const getLanguageSuffix = (lang) => {
    const map = {
      '中文': 'cn',
      'English': 'en',
      '日本語': 'jp'
    };
    return map[lang] || 'cn';
  };

  // 获取用例集名称（使用 case_group 字段）
  const getCaseName = () => {
    return caseData.case_group || t('manualTest.untitledCase');
  };

  // 打开编辑Modal
  const handleOpenModal = () => {
    console.log('[CaseListItem] 打开编辑Modal:', getCaseName());
    setEditValue(getCaseName());
    setModalVisible(true);
  };

  // 保存编辑
  const handleSave = async () => {
    const groupId = caseData._groupId || caseData.id;
    console.log('[CaseListItem] 开始保存:', { 
      projectId, 
      groupId, 
      oldName: getCaseName(), 
      newName: editValue.trim(),
      hasOnEdit: !!onEdit
    });

    if (!editValue.trim()) {
      message.error(t('manualTest.caseNameRequired'));
      return;
    }

    setSaving(true);
    try {
      // 如果有外部onEdit回调，使用它（用于API用例集等特殊场景）
      if (onEdit) {
        await onEdit(editValue.trim());
        console.log('[CaseListItem] 通过onEdit回调更新成功');
        message.success(t('manualTest.updateCaseSuccess'));
        caseData.case_group = editValue.trim();
        setModalVisible(false);
      } else {
        // 否则使用默认的updateCaseGroup API
        const updateData = {
          groupName: editValue.trim()
        };

        console.log('[CaseListItem] 调用updateCaseGroup API:', { groupId, updateData });
        await updateCaseGroup(groupId, updateData);
        console.log('[CaseListItem] 更新成功');
        message.success(t('manualTest.updateCaseSuccess'));
        
        // 更新本地显示
        caseData.case_group = editValue.trim();
        setModalVisible(false);
        if (onUpdate) {
          onUpdate(groupId, editValue.trim());
        }
      }
    } catch (error) {
      console.error('[CaseListItem] 更新用例集失败:', error);
      message.error(t('manualTest.updateCaseFailed'));
    } finally {
      setSaving(false);
    }
  };

  // 删除用例集
  const handleDelete = async () => {
    const groupId = caseData._groupId || caseData.id;
    console.log('[CaseListItem] 点击删除按钮:', { 
      projectId, 
      groupId,
      caseGroup: caseData.case_group,
      caseName: getCaseName(),
      hasOnDelete: !!onDelete
    });
    
    // 如果有外部onDelete回调且没有projectId（API用例集场景），直接使用回调
    if (onDelete && !projectId) {
      modal.confirm({
        title: t('manualTest.deleteCaseGroup'),
        icon: <ExclamationCircleOutlined />,
        content: `确认删除用例集 "${getCaseName()}"？`,
        okText: t('common.confirm'),
        cancelText: t('common.cancel'),
        okType: 'danger',
        onOk: async () => {
          console.log('[CaseListItem] 通过onDelete回调删除');
          onDelete();
        }
      });
      return;
    }
    
    // 否则使用默认逻辑（手工用例集场景）
    // 先获取该用例集的所有用例数量
    try {
      const { getCasesList } = await import('../../../../api/manualCase');
      const response = await getCasesList(projectId, { caseType: 'overall' });
      const casesToUpdate = response.cases.filter(c => c.case_group === caseData.case_group);
      
      console.log('[CaseListItem] 用例集包含用例数量:', casesToUpdate.length);
      
      modal.confirm({
        title: t('manualTest.deleteCaseGroup'),
        icon: <ExclamationCircleOutlined />,
        content: `确认删除用例集 "${getCaseName()}"？该用例集内的 ${casesToUpdate.length} 条用例将被移至"未分组"`,
        okText: t('common.confirm'),
        cancelText: t('common.cancel'),
        okType: 'danger',
        onOk: async () => {
          console.log('[CaseListItem] 确认删除用例集');
          try {
            // 1. 将该用例集的所有用例的case_group清空
            if (casesToUpdate.length > 0) {
              const { updateCase } = await import('../../../../api/manualCase');
              const updatePromises = casesToUpdate.map(c => 
                updateCase(projectId, c.case_id, { case_group: '' })
              );
              await Promise.all(updatePromises);
            }
            
            // 2. 删除case_groups表中的记录
            await deleteCaseGroup(groupId);
            
            console.log('[CaseListItem] 用例集删除成功，用例已移至未分组');
            message.success('用例集已删除，用例已移至未分组');
            if (onDelete) {
              console.log('[CaseListItem] 调用onDelete回调');
              onDelete(groupId);
            }
          } catch (error) {
            console.error('[CaseListItem] 删除用例集失败:', error);
            message.error(t('message.deleteFailed'));
          }
        }
      });
    } catch (error) {
      console.error('[CaseListItem] 获取用例列表失败:', error);
      message.error('获取用例列表失败');
    }
  };

  // 点击用例名，切换到该用例集
  const handleCaseClick = () => {
    console.log('[CaseListItem] 点击用例集:', caseData.case_group);
    if (onSwitch) {
      // 传递完整的caseData给父组件
      onSwitch(caseData);
    }
  };

  return (
    <>
      <div className={`case-list-item ${isSelected ? 'selected' : ''}`}>
        <span className="case-name" onClick={handleCaseClick}>{getCaseName()}</span>
        <div className="action-icons">
          <CopyOutlined
            className="copy-icon"
            onClick={(e) => {
              e.stopPropagation();
              e.preventDefault();
              navigator.clipboard.writeText(getCaseName());
              message.success('复制成功');
            }}
            style={{ cursor: 'pointer' }}
            title="复制用例集名称"
          />
          <DeleteOutlined 
            className="delete-icon" 
            onClick={(e) => {
              console.log('[CaseListItem] DeleteOutlined clicked');
              e.stopPropagation();
              e.preventDefault();
              handleDelete();
            }}
            style={{ cursor: 'pointer' }}
          />
        </div>
      </div>

      <Modal
        title={t('manualTest.editCaseGroup')}
        open={modalVisible}
        onOk={handleSave}
        onCancel={() => setModalVisible(false)}
        confirmLoading={saving}
        okText={t('common.save')}
        cancelText={t('common.cancel')}
      >
        <div style={{ marginBottom: 16 }}>
          <label>{t('manualTest.caseGroupName')}:</label>
          <Input
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
            placeholder={t('manualTest.enterCaseName')}
            disabled={saving}
            onPressEnter={handleSave}
            style={{ marginTop: 8 }}
          />
        </div>
      </Modal>
    </>
  );
};

export default CaseListItem;
