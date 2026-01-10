import { useState, useEffect } from 'react';
import { Button, Spin, Typography, message, Modal } from 'antd';
import { PlusOutlined, ProjectOutlined, DeleteOutlined, MenuFoldOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { getProjects, deleteProject } from '../api/project';
import CreateProjectModal from '../pages/ProjectList/CreateProjectModal';

const { Text } = Typography;

const ProjectListSidebar = ({ onProjectSelect, selectedProjectId: externalSelectedProjectId, onCollapse }) => {
  const navigate = useNavigate();
  const { projectId } = useParams();
  const { t } = useTranslation();
  const { user } = useSelector(state => state.auth);
  
  const [projects, setProjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedProjectId, setSelectedProjectId] = useState(externalSelectedProjectId || projectId);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [deleteModalVisible, setDeleteModalVisible] = useState(false);
  const [projectToDelete, setProjectToDelete] = useState(null);

  // 加载项目列表
  const fetchProjects = async () => {
    try {
      setLoading(true);
      console.log('[ProjectListSidebar] Fetching projects...');
      const data = await getProjects();
      console.log('[ProjectListSidebar] Received data:', data);
      console.log('[ProjectListSidebar] Is array?', Array.isArray(data));
      console.log('[ProjectListSidebar] Data length:', Array.isArray(data) ? data.length : 'N/A');
      // apiClient已经返回response.data，直接使用
      const projectList = Array.isArray(data) ? data : [];
      console.log('[ProjectListSidebar] Setting projects:', projectList);
      setProjects(projectList);
    } catch (error) {
      console.error('[ProjectListSidebar] Failed to load projects:', error);
      message.error(t('project.loadFailed'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();

    // 监听项目创建成功事件
    const handleProjectCreated = () => {
      console.log('[ProjectListSidebar] Project created, refreshing list...');
      fetchProjects();
    };

    window.addEventListener('projectCreated', handleProjectCreated);

    return () => {
      window.removeEventListener('projectCreated', handleProjectCreated);
    };
  }, [t]);

  // 同步选中状态
  useEffect(() => {
    if (externalSelectedProjectId !== undefined) {
      setSelectedProjectId(externalSelectedProjectId?.toString());
    } else if (projectId) {
      setSelectedProjectId(projectId);
    }
  }, [externalSelectedProjectId, projectId]);

  // 处理项目点击
  const handleProjectClick = (project) => {
    setSelectedProjectId(project.id.toString());
    if (onProjectSelect) {
      // 嵌入式模式：通过回调通知父组件
      onProjectSelect(project.id);
    } else {
      // 路由模式：跳转到详情页
      navigate(`/projects/${project.id}?tab=project-info`);
    }
  };

  // 处理删除项目
  const handleDeleteProject = (project, e) => {
    console.log('[ProjectListSidebar] Delete button clicked for project:', project);
    e.stopPropagation();
    setProjectToDelete(project);
    setDeleteModalVisible(true);
    console.log('[ProjectListSidebar] Showing delete confirmation modal');
  };

  // 确认删除
  const handleConfirmDelete = async () => {
    if (!projectToDelete) return;
    
    console.log('[ProjectListSidebar] User confirmed deletion');
    try {
      await deleteProject(projectToDelete.id);
      console.log('[ProjectListSidebar] Delete API success');
      message.success(t('project.deleteSuccess'));
      setDeleteModalVisible(false);
      setProjectToDelete(null);
      
      // 重新加载项目列表
      const data = await getProjects();
      const newProjectList = Array.isArray(data) ? data : [];
      setProjects(newProjectList);
      
      // 如果删除的是当前项目
      if (selectedProjectId === projectToDelete.id.toString()) {
        if (newProjectList.length > 0) {
          // 选择第一个项目
          if (onProjectSelect) {
            // 嵌入模式：使用回调
            onProjectSelect(newProjectList[0].id);
          } else {
            // 路由模式：导航
            navigate(`/projects/${newProjectList[0].id}?tab=project-info`);
          }
        } else {
          // 没有项目了
          if (!onProjectSelect) {
            navigate('/projects');
          }
        }
      }
    } catch (error) {
      console.error('[ProjectListSidebar] Failed to delete project:', error);
      message.error(t('project.deleteFailed'));
    }
  };

  // 取消删除
  const handleCancelDelete = () => {
    setDeleteModalVisible(false);
    setProjectToDelete(null);
  };

  // 处理创建项目
  const handleCreateProject = () => {
    console.log('[ProjectListSidebar] Opening create project modal');
    setCreateModalVisible(true);
  };

  // 处理创建成功
  const handleCreateSuccess = () => {
    console.log('[ProjectListSidebar] Project created successfully');
    setCreateModalVisible(false);
    fetchProjects();
  };

  if (loading) {
    return (
      <div style={{ 
        width: 130, 
        background: '#fafafa', 
        borderRight: '1px solid #f0f0f0',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100%'
      }}>
        <Spin tip={t('common.loading')} />
      </div>
    );
  }

  return (
    <div style={{ 
      width: 130, 
      background: '#fafafa', 
      borderRight: '1px solid #f0f0f0',
      display: 'flex',
      flexDirection: 'column',
      height: '100%'
    }}>
      {/* 顶部按钮区域 */}
      <div style={{ padding: 8, borderBottom: '1px solid #f0f0f0' }}>
        {/* 创建项目按钮 (仅项目管理员) */}
        {user?.role === 'project_manager' && (
          <Button
            type="primary"
            icon={<PlusOutlined />}
            size="small"
            block
            onClick={handleCreateProject}
            style={{ marginBottom: 8, fontSize: '11px', padding: '4px 8px', height: 'auto', lineHeight: '1.4' }}
          >
            {t('project.createProject')}
          </Button>
        )}
        {/* 收起按钮 */}
        {onCollapse && (
          <Button
            type="text"
            icon={<MenuFoldOutlined />}
            size="small"
            onClick={onCollapse}
            style={{ width: '100%' }}
          />
        )}
      </div>

      {/* 项目列表 */}
      <div style={{ 
        flex: 1, 
        overflowY: 'auto',
        padding: '8px 0'
      }}>
        {projects.length === 0 ? (
          <div style={{ 
            padding: '24px 8px', 
            textAlign: 'center', 
            color: 'rgba(0, 0, 0, 0.45)' 
          }}>
            <Text type="secondary" style={{ fontSize: '12px' }}>{t('project.empty')}</Text>
          </div>
        ) : (
          <div style={{ padding: '4px 0' }}>
            {projects.map(project => (
              <div
                key={project.id.toString()}
                onClick={() => handleProjectClick(project)}
                style={{
                  padding: '8px',
                  margin: '2px 4px',
                  borderRadius: 4,
                  cursor: 'pointer',
                  backgroundColor: selectedProjectId === project.id.toString() ? '#e6f7ff' : 'transparent',
                  border: selectedProjectId === project.id.toString() ? '1px solid #91d5ff' : '1px solid transparent',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  transition: 'all 0.3s',
                }}
                onMouseEnter={(e) => {
                  if (selectedProjectId !== project.id.toString()) {
                    e.currentTarget.style.backgroundColor = '#f5f5f5';
                  }
                }}
                onMouseLeave={(e) => {
                  if (selectedProjectId !== project.id.toString()) {
                    e.currentTarget.style.backgroundColor = 'transparent';
                  }
                }}
                title={project.name}
              >
                <div style={{
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                  fontSize: '12px',
                  flex: 1,
                }}>
                  <ProjectOutlined style={{ marginRight: 4, fontSize: '12px' }} />
                  {project.name}
                </div>
                {user?.role === 'project_manager' && (
                  <DeleteOutlined
                    onClick={(e) => handleDeleteProject(project, e)}
                    style={{
                      fontSize: '12px',
                      color: '#ff4d4f',
                      padding: '2px',
                      opacity: 0.6,
                    }}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.opacity = 1;
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.opacity = 0.6;
                    }}
                  />
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* 创建项目Modal */}
      <CreateProjectModal
        visible={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onSuccess={handleCreateSuccess}
      />

      {/* 删除确认Modal */}
      <Modal
        title={t('project.deleteProject')}
        open={deleteModalVisible}
        onOk={handleConfirmDelete}
        onCancel={handleCancelDelete}
        okText={t('common.confirm')}
        cancelText={t('common.cancel')}
        okType="danger"
      >
        {projectToDelete && (
          <p>{t('project.confirmDelete', { name: projectToDelete.name })}</p>
        )}
      </Modal>
    </div>
  );
};

export default ProjectListSidebar;
