import React, { useState, useEffect, useCallback } from 'react';
import { Row, Col, Spin, Empty, message, Typography } from 'antd';
import { useTranslation } from 'react-i18next';
import { getProjects } from '../../api/project';
import CreateProjectModal from './CreateProjectModal';
import ProjectCard from './ProjectCard';
import './index.css';

const { Title } = Typography;

const ProjectList = () => {
  const { t } = useTranslation();
  const [projects, setProjects] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);

  const fetchProjects = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getProjects();
      // apiClient已经返回response.data，直接使用
      setProjects(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('[ProjectList] Failed to load projects:', error);
      message.error(t('project.loadFailed'));
    } finally {
      setLoading(false);
    }
  }, [t]);

  useEffect(() => {
    fetchProjects();
    
    // 监听来自侧边栏的创建项目事件
    const handleOpenModal = () => {
      setModalVisible(true);
    };
    
    window.addEventListener('openCreateProjectModal', handleOpenModal);
    
    // 检查URL参数,如果有action=create则打开modal
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get('action') === 'create') {
      setModalVisible(true);
      // 清除URL参数
      window.history.replaceState({}, '', '/projects');
    }
    
    return () => {
      window.removeEventListener('openCreateProjectModal', handleOpenModal);
    };
  }, [fetchProjects]);

  const handleCreateSuccess = () => {
    setModalVisible(false);
    fetchProjects();
  };

  const handleUpdateProject = (updatedProject) => {
    console.log('Updating project in list:', updatedProject);
    setProjects((prevProjects) =>
      prevProjects.map((p) => (p.id === updatedProject.id ? updatedProject : p))
    );
  };

  const handleDeleteProject = (deletedProjectId) => {
    console.log('Deleting project from list, id:', deletedProjectId, 'type:', typeof deletedProjectId);
    console.log('Current projects:', projects);
    setProjects((prevProjects) => {
      const filtered = prevProjects.filter((p) => {
        console.log('Comparing p.id:', p.id, 'type:', typeof p.id, 'with deletedProjectId:', deletedProjectId, 'equal?', p.id !== deletedProjectId);
        return p.id !== deletedProjectId;
      });
      console.log('Filtered projects:', filtered);
      return filtered;
    });
  };

  return (
    <div className="project-list-container">
      <Spin spinning={loading}>
        {projects.length === 0 && !loading ? (
          <Empty description={t('project.empty')} />
        ) : (
          <Row gutter={[16, 16]}>
            {projects.map((project) => (
              <Col xs={24} sm={12} md={8} lg={6} key={project.id}>
                <ProjectCard
                  project={project}
                  onUpdate={handleUpdateProject}
                  onDelete={handleDeleteProject}
                />
              </Col>
            ))}
          </Row>
        )}
      </Spin>

      <CreateProjectModal
        visible={modalVisible}
        onCancel={() => setModalVisible(false)}
        onSuccess={handleCreateSuccess}
      />
    </div>
  );
};

export default ProjectList;
