import { useState, useEffect } from 'react';
import { Layout, Spin, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import ProjectDetail from '../ProjectDetail';
import './index.css';

const { Content } = Layout;

const ProjectManagement = () => {
  const { t } = useTranslation();
  const { currentProjectId } = useSelector(state => state.project);

  return (
    <div style={{ height: '100%', width: '100%', background: '#f0f2f5', overflow: 'auto', padding: '24px' }}>
      {currentProjectId ? (
        <ProjectDetail projectId={currentProjectId} />
      ) : (
        <div style={{ 
          display: 'flex', 
          justifyContent: 'center', 
          alignItems: 'center', 
          height: '100%',
          color: 'rgba(0, 0, 0, 0.45)'
        }}>
          {t('project.selectProjectHint') || t('project.empty')}
        </div>
      )}
    </div>
  );
};

export default ProjectManagement;
