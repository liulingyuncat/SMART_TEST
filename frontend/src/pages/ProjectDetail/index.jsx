import React, { useState, useEffect } from 'react';
import { useParams, useSearchParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Tabs, Card, Spin, message, Button } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import { getProjectById } from '../../api/project';
import ProjectInfoTab from '../../components/ProjectInfoTab';
import AIRequirementTab from './RequirementManagement/AIRequirementTab';
import AIViewpointTab from './RequirementManagement/AIViewpointTab';
import RawDocumentTab from './RawDocumentTab';
import ManualTestTabs from './ManualTestTabs';
import AutoTestTab from './AutoTestTabs/containers/AutoTestTab';
import ApiTestTabs from './ApiTestTabs';
import TestExecution from './TestExecution';
import DefectManagement from './DefectManagement';
import AIReportManagementTab from './AIReportManagement';
import ComingSoonPlaceholder from '../../components/ComingSoonPlaceholder';
import './index.css';

// Tab配置数组
const TAB_CONFIG = [
  { key: 'project-info', labelKey: 'projectDetail.projectInfo', component: ProjectInfoTab },
  { key: 'raw-document', labelKey: 'projectDetail.rawDocument', component: RawDocumentTab },
  { key: 'ai-requirement', labelKey: 'projectDetail.aiRequirement', component: AIRequirementTab },
  { key: 'ai-viewpoint', labelKey: 'projectDetail.aiViewpoint', component: AIViewpointTab },
  { key: 'manual-test', labelKey: 'projectDetail.manualTest', component: ManualTestTabs },
  { key: 'auto-test', labelKey: 'projectDetail.autoTest', component: AutoTestTab },
  { key: 'api-test', labelKey: 'projectDetail.apiTest', component: ApiTestTabs },
  { key: 'execution', labelKey: 'projectDetail.execution', component: TestExecution },
  { key: 'bug', labelKey: 'projectDetail.bug', component: DefectManagement },
  { key: 'ai-report', labelKey: 'projectDetail.aiReport', component: AIReportManagementTab },
  // { key: 'statistics', labelKey: 'projectDetail.statistics', component: ComingSoonPlaceholder }, // 暂时隐藏统计分析
];

const ProjectDetail = ({ projectId: externalProjectId }) => {
  const { t } = useTranslation();
  const { id: paramId } = useParams();
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  
  const id = externalProjectId || paramId;
  const [loading, setLoading] = useState(true);
  const [project, setProject] = useState(null);
  const [activeTab, setActiveTab] = useState('project-info');

  // 初始化时从URL读取tab参数
  useEffect(() => {
    const tabParam = searchParams.get('tab');
    if (tabParam && TAB_CONFIG.find(tab => tab.key === tabParam)) {
      setActiveTab(tabParam);
    }
  }, [searchParams]);

  // 加载项目数据
  useEffect(() => {
    const fetchProject = async () => {
      setLoading(true);
      try {
        const data = await getProjectById(id);
        setProject(data);
      } catch (error) {
        if (error.response?.status === 403) {
          message.error(t('projectDetail.permissionDenied'));
          navigate('/projects');
        } else if (error.response?.status === 404) {
          message.error(t('projectDetail.notFound'));
          navigate('/projects');
        } else {
          message.error(t('projectDetail.loadFailed'));
        }
      } finally {
        setLoading(false);
      }
    };

    fetchProject();
  }, [id, navigate, t]);

  // Tab切换处理
  const handleTabChange = (key) => {
    setActiveTab(key);
    setSearchParams({ tab: key });
  };

  // 渲染Tab内容
  const renderTabContent = () => {
    const currentTab = TAB_CONFIG.find(tab => tab.key === activeTab);
    if (currentTab) {
      const Component = currentTab.component;
      // 为特定Tab传递props
      if (activeTab === 'project-info') {
        return <Component projectId={id} />;
      }
      if (['ai-requirement', 'ai-viewpoint'].includes(activeTab)) {
        return <Component projectId={id} projectName={project.name} />;
      }
      if (['raw-document', 'manual-test', 'auto-test', 'api-test'].includes(activeTab)) {
        return <Component projectId={id} />;
      }
      if (activeTab === 'execution') {
        return <Component projectId={id} projectName={project.name} />;
      }
      if (activeTab === 'bug') {
        return <Component projectId={id} projectName={project.name} />;
      }
      if (activeTab === 'ai-report') {
        return <Component projectId={id} projectName={project.name} />;
      }
      if (activeTab === 'statistics') {
        return <Component feature={t(currentTab.labelKey)} />;
      }
      return <Component />;
    }
    return null;
  };

  if (loading) {
    return (
      <div className="project-detail-loading">
        <Spin size="large" />
      </div>
    );
  }

  if (!project) {
    return null;
  }

  return (
    <div className="project-detail-container">
      <Tabs activeKey={activeTab} onChange={handleTabChange}>
        {TAB_CONFIG.map(tab => (
          <Tabs.TabPane tab={t(tab.labelKey)} key={tab.key} />
        ))}
      </Tabs>
      
      <Card>
        {renderTabContent()}
      </Card>
    </div>
  );
};

export default ProjectDetail;
