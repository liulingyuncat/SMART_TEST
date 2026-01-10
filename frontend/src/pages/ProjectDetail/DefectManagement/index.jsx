import React, { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { message } from 'antd';
import DefectListPage from './DefectListPage';
import DefectCreatePage from './DefectCreatePage';
import DefectDetailPage from './DefectDetailPage';
import { fetchDefectSubjects, fetchDefectPhases } from '../../../api/defect';
import './index.css';

/**
 * 缺陷管理模块入口组件
 * 管理 list/create/detail 三种视图状态切换
 */
const DefectManagement = ({ projectId, projectName }) => {
  const { t } = useTranslation();
  // 视图状态: 'list' | 'create' | 'detail'
  const [view, setView] = useState('list');
  const [currentDefectId, setCurrentDefectId] = useState(null);
  
  // 配置数据
  const [subjects, setSubjects] = useState([]);
  const [phases, setPhases] = useState([]);
  const [configLoading, setConfigLoading] = useState(false);

  // 加载配置数据 (Subject/Phase)
  const loadConfig = useCallback(async () => {
    if (!projectId) return;
    setConfigLoading(true);
    try {
      const [subjectsData, phasesData] = await Promise.all([
        fetchDefectSubjects(projectId),
        fetchDefectPhases(projectId),
      ]);
      setSubjects(Array.isArray(subjectsData) ? subjectsData : []);
      setPhases(Array.isArray(phasesData) ? phasesData : []);
    } catch (error) {
      console.error('Failed to load defect config:', error);
    } finally {
      setConfigLoading(false);
    }
  }, [projectId]);

  // 初始加载配置
  useEffect(() => {
    loadConfig();
  }, [loadConfig]);

  // 切换到创建页面
  const handleCreate = useCallback(() => {
    setView('create');
    setCurrentDefectId(null);
  }, []);

  // 切换到详情页面
  const handleDetail = useCallback((defect) => {
    // 使用显示ID（defect_id）而不是内部UUID
    setCurrentDefectId(defect.defect_id || defect.id);
    setView('detail');
  }, []);

  // 返回列表页面
  const handleBackToList = useCallback(() => {
    setView('list');
    setCurrentDefectId(null);
  }, []);

  // 创建成功后切换到详情页
  const handleCreateSuccess = useCallback((defect) => {
    console.log('[DEBUG] handleCreateSuccess: received defect', defect);
    message.success(t('defect.createSuccess'));
    // 使用 defect_id (显示ID如000001)，而不是 id (UUID)
    const displayId = defect?.defect_id || defect?.id;
    console.log('[DEBUG] handleCreateSuccess: setting currentDefectId to', displayId);
    setCurrentDefectId(displayId);
    setView('detail');
  }, [t]);

  // 配置更新回调
  const handleConfigUpdate = useCallback(() => {
    loadConfig();
  }, [loadConfig]);

  // 根据视图状态渲染对应组件
  const renderContent = () => {
    switch (view) {
      case 'create':
        return (
          <DefectCreatePage
            projectId={projectId}
            subjects={subjects}
            phases={phases}
            onSuccess={handleCreateSuccess}
            onCancel={handleBackToList}
          />
        );
      case 'detail':
        return (
          <DefectDetailPage
            projectId={projectId}
            defectId={currentDefectId}
            subjects={subjects}
            phases={phases}
            onBack={handleBackToList}
            onCreate={handleCreate}
            onConfigUpdate={handleConfigUpdate}
          />
        );
      case 'list':
      default:
        return (
          <DefectListPage
            projectId={projectId}
            projectName={projectName}
            subjects={subjects}
            phases={phases}
            configLoading={configLoading}
            onCreate={handleCreate}
            onDetail={handleDetail}
            onConfigUpdate={handleConfigUpdate}
          />
        );
    }
  };

  return (
    <div className="defect-management-container">
      {renderContent()}
    </div>
  );
};

export default DefectManagement;
