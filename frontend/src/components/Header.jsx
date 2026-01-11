import { Layout, Button, Typography, message, Select, Space, Tooltip, Modal, Tag } from 'antd';
import { LogoutOutlined, BulbOutlined, UserOutlined, TeamOutlined, IdcardOutlined, ProjectOutlined, CheckOutlined, CloseOutlined, PlusOutlined, DeleteOutlined, TagOutlined } from '@ant-design/icons';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useMemo, useState, useEffect } from 'react';
import { logout } from '../store/authSlice';
import { setCurrentProject, setProjects, clearProjects } from '../store/projectSlice';
import LanguageSwitch from './LanguageSwitch';
import { getProjects, deleteProject } from '../api/project';
import { setCurrentProject as setCurrentProjectAPI } from '../api/profile';
import CreateProjectModal from '../pages/ProjectList/CreateProjectModal';import { VERSION } from '../version';
const { Header: AntHeader } = Layout;
const { Title } = Typography;

const Header = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useSelector(state => state.auth);
  const { currentProjectId, projects } = useSelector(state => state.project);
  const { t } = useTranslation();
  const displayName = user?.nickname || user?.username;
  
  // 选中的导航项状态
  const [activeNav, setActiveNav] = useState('');
  
  // 临时选中的项目（未确认）
  const [pendingProjectId, setPendingProjectId] = useState(null);
  const [isSyncing, setIsSyncing] = useState(false);
  
  // 创建和删除项目相关状态
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [deleteModalVisible, setDeleteModalVisible] = useState(false);

  // 根据当前路由和角色确定默认选中项
  useEffect(() => {
    const pathname = location.pathname;
    const role = user?.role;
    
    // 根据路由确定选中项
    if (pathname.startsWith('/projects')) {
      setActiveNav('projects');
    } else if (pathname.startsWith('/users')) {
      setActiveNav('users');
    } else if (pathname.startsWith('/prompt')) {
      setActiveNav('prompt');
    } else if (pathname.startsWith('/assign')) {
      setActiveNav('assign');
    } else if (pathname.startsWith('/profile')) {
      setActiveNav('profile');
    } else {
      // 默认选中第一个导航项
      if (role === 'project_manager' || role === 'project_member') {
        setActiveNav('projects');
      } else if (role === 'system_admin') {
        setActiveNav('users');
      }
    }
  }, [location.pathname, user?.role]);

  // 获取项目列表并恢复用户的项目选择
  // 【关键流程】
  // 1. 初始加载：项目选择为空
  // 2. 用户选择项目 + 保存：保存到localStorage和后端
  // 3. 创建新项目：项目选择不变（已保存到localStorage）
  // 4. 登出系统：认证信息清除，但currentProjectId在localStorage中保留
  // 5. 重新登录：此useEffect重新执行，从localStorage恢复项目选择
  useEffect(() => {
    const fetchProjects = async () => {
      if (user?.role === 'project_manager' || user?.role === 'project_member') {
        try {
          console.log('[Header] Fetching projects for user:', user?.username);
          const data = await getProjects();
          const projectList = Array.isArray(data) ? data : [];
          dispatch(setProjects(projectList));
          console.log('[Header] Projects fetched:', projectList.length, 'projects');
          
          // 【关键修复】从localStorage恢复项目选择（包括登出后重新登录的场景）
          const savedProjectId = localStorage.getItem('currentProjectId');
          if (savedProjectId) {
            const projectId = parseInt(savedProjectId, 10);
            // 验证保存的项目ID是否仍在项目列表中
            const projectExists = projectList.some(p => p.id === projectId);
            if (projectExists) {
              dispatch(setCurrentProject(projectId));
              console.log('[Header] ✓ Restored project selection from localStorage:', projectId);
            } else {
              // 项目不存在（可能被删除），清除localStorage中的缓存
              localStorage.removeItem('currentProjectId');
              console.log('[Header] ⚠ Saved project no longer exists, cleared from localStorage');
            }
          } else {
            console.log('[Header] No saved project selection in localStorage');
          }
        } catch (error) {
          console.error('[Header] Failed to fetch projects:', error);
        }
      }
    };
    fetchProjects();
  }, [user?.role, dispatch]);

  // 处理项目选择变化
  const handleProjectChange = (value) => {
    // 只是设置临时值，不立即同步
    setPendingProjectId(value);
  };

  // 确认项目选择
  const handleConfirmProject = async () => {
    if (pendingProjectId === null) return;
    
    setIsSyncing(true);
    try {
      // 调用API同步到后端
      await setCurrentProjectAPI(pendingProjectId);
      // 更新本地Redux状态
      dispatch(setCurrentProject(pendingProjectId));
      // 【关键修复】持久化项目选择到localStorage
      localStorage.setItem('currentProjectId', pendingProjectId.toString());
      // 清除临时状态
      setPendingProjectId(null);
      message.success(t('common.updateSuccess') || 'Project updated successfully');
    } catch (error) {
      console.error('[Header] Failed to set current project:', error);
      message.error(t('common.updateFailed') || 'Failed to update project');
      // 同步失败时，恢复Select的值
      setPendingProjectId(null);
    } finally {
      setIsSyncing(false);
    }
  };

  // 取消项目选择
  const handleCancelProject = () => {
    setPendingProjectId(null);
  };

  // 刷新项目列表
  const refreshProjects = async () => {
    try {
      const data = await getProjects();
      const projectList = Array.isArray(data) ? data : [];
      dispatch(setProjects(projectList));
      return projectList;
    } catch (error) {
      console.error('[Header] Failed to refresh projects:', error);
      return [];
    }
  };

  // 处理创建项目成功
  const handleCreateSuccess = async () => {
    setCreateModalVisible(false);
    await refreshProjects();
  };

  // 处理删除项目
  const handleDeleteProject = () => {
    if (!currentProjectId) {
      message.warning(t('project.selectProjectFirst') || '请先选择一个项目');
      return;
    }
    setDeleteModalVisible(true);
  };

  // 确认删除项目
  const handleConfirmDelete = async () => {
    if (!currentProjectId) return;
    
    try {
      await deleteProject(currentProjectId);
      message.success(t('project.deleteSuccess'));
      setDeleteModalVisible(false);
      
      // 刷新项目列表
      const newProjectList = await refreshProjects();
      
      // 清除当前项目选择
      dispatch(setCurrentProject(null));
      localStorage.removeItem('currentProjectId');
      
      // 如果还有项目，选择第一个
      if (newProjectList.length > 0) {
        dispatch(setCurrentProject(newProjectList[0].id));
        localStorage.setItem('currentProjectId', newProjectList[0].id.toString());
      }
    } catch (error) {
      console.error('[Header] Failed to delete project:', error);
      message.error(t('project.deleteFailed'));
    }
  };

  // 获取当前选中项目的名称
  const getCurrentProjectName = () => {
    const project = projects.find(p => p.id === currentProjectId);
    return project?.name || '';
  };

  // 导航按钮样式
  const getNavButtonStyle = (navKey) => ({
    color: activeNav === navKey ? '#1890ff' : 'rgba(0, 0, 0, 0.65)',
    backgroundColor: activeNav === navKey ? '#e6f7ff' : 'transparent',
    borderRadius: '4px',
    fontWeight: activeNav === navKey ? 500 : 400,
  });

  // 根据角色生成菜单项
  const menuItems = useMemo(() => {
    const role = user?.role;
    const items = [];

    // 系统管理员: 人员管理、提示词管理 (不显示人员分配)
    if (role === 'system_admin') {
      items.push(
        { key: 'users', icon: <UserOutlined />, label: t('menu.users'), action: () => navigate('/users') },
        { key: 'prompt', icon: <BulbOutlined />, label: t('menu.promptManagement'), action: () => navigate('/prompt-management') },
      );
    }
    
    // 项目管理员: 人员管理、提示词管理 (不显示人员分配)
    if (role === 'project_manager') {
      items.push(
        { key: 'users', icon: <UserOutlined />, label: t('menu.users'), action: () => navigate('/users') },
        { key: 'prompt', icon: <BulbOutlined />, label: t('menu.promptManagement'), action: () => navigate('/prompt-management') },
      );
    }
    
    // 项目成员: 提示词管理
    if (role === 'project_member') {
      items.push(
        { key: 'prompt', icon: <BulbOutlined />, label: t('menu.promptManagement'), action: () => navigate('/prompt-management') },
      );
    }
    
    // 所有角色: 个人中心
    items.push({ key: 'profile', icon: <IdcardOutlined />, label: t('menu.personalCenter'), action: () => navigate('/profile') });

    return items;
  }, [user?.role, t, navigate]);

  const handleLogout = () => {
    // 日志埋点：记录用户退出信息
    console.info('[Navigation] User logout:', {
      userId: user?.id,
      username: user?.username,
      timestamp: new Date().toISOString(),
    });

    // Step 1: 清除Redux状态（只清除认证信息，不清除项目选择）
    dispatch(logout());
    
    // Step 1.5: 清除项目列表，但保留项目选择（便于重新登录时恢复）
    dispatch(setProjects([]));
    console.log('[Navigation] Projects cleared, currentProjectId preserved:', currentProjectId);

    // Step 2 & 3: 清除本地存储 (容错处理)
    // 【重要】不清除currentProjectId，登出后项目选择保持不变，便于重新登录时恢复
    try {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user_info');
      // 注意：intentionally NOT removing currentProjectId to preserve user's project selection
      console.log('[Navigation] Logout: currentProjectId preserved in localStorage');
    } catch (error) {
      console.warn('[Navigation] Failed to clear localStorage:', error);
    }

    // Step 4: 跳转到登录页 (替换历史记录)
    navigate('/login', { replace: true });
  };

  // 处理项目管理按钮点击
  const handleProjectsClick = () => {
    // 直接跳转到项目列表页
    navigate('/projects');
  };

  return (
    <AntHeader
      style={{
        background: '#ffffff',
        padding: '0 24px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        borderBottom: '1px solid #f0f0f0',
        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.06)',
      }}
    >
      {/* 左侧：平台名称 */}
      <div style={{ display: 'flex', alignItems: 'baseline', flex: 1, gap: '12px' }}>
        <Title
          level={4}
          style={{
            margin: 0,
            color: '#1890ff',
            fontWeight: 700,
            letterSpacing: '1px',
          }}
        >
          SMART TEST
        </Title>
        <Tag 
          icon={<TagOutlined />} 
          color="blue"
          style={{ 
            fontSize: '11px', 
            padding: '2px 8px',
            fontFamily: 'monospace',
            fontWeight: 500,
          }}
        >
          v{VERSION}
        </Tag>
        <Typography.Text
          style={{
            color: '#8c8c8c',
            fontSize: '12px',
            fontWeight: 400,
          }}
        >
          PEVVD Intelligent Test Platform
        </Typography.Text>
      </div>

      {/* 中间：当前项目选择器 */}
      {(user?.role === 'project_manager' || user?.role === 'project_member') && (
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', flex: 1, gap: '8px' }}>
          <Typography.Text style={{ marginRight: '8px', color: 'rgba(0, 0, 0, 0.85)' }}>
            {t('assign.currentProject')}
          </Typography.Text>
          
          {projects.length > 0 ? (
            <>
              <Select
                value={pendingProjectId !== null ? pendingProjectId : (currentProjectId || undefined)}
                onChange={handleProjectChange}
                style={{ minWidth: '150px' }}
                placeholder={t('assign.selectProjectPlaceholder') || '请选择项目'}
                allowClear
                onClear={() => {
                  // 清除当前项目选择
                  setPendingProjectId(null);
                }}
              >
                {projects.map(project => (
                  <Select.Option key={project.id} value={project.id}>
                    {project.name}
                  </Select.Option>
                ))}
              </Select>
              
              {/* 确认/取消按钮（仅在有未确认的选择时显示） */}
              {pendingProjectId !== null && (
                <Space size={4}>
                  <Button
                    type="primary"
                    size="small"
                    icon={<CheckOutlined />}
                    onClick={handleConfirmProject}
                    loading={isSyncing}
                    title={t('common.confirm') || 'Confirm'}
                  />
                  <Button
                    size="small"
                    icon={<CloseOutlined />}
                    onClick={handleCancelProject}
                    disabled={isSyncing}
                    title={t('common.cancel') || 'Cancel'}
                  />
                </Space>
              )}
            </>
          ) : (
            <Typography.Text type="secondary" style={{ fontSize: '12px' }}>
              {t('project.empty')}
            </Typography.Text>
          )}
          
          {/* 创建项目按钮 (仅项目管理员) */}
          {user?.role === 'project_manager' && (
            <Tooltip title={t('project.createProject')}>
              <Button
                type="text"
                size="small"
                icon={<PlusOutlined />}
                onClick={() => setCreateModalVisible(true)}
                style={{ color: '#1890ff' }}
              />
            </Tooltip>
          )}
          
          {/* 删除项目按钮 (仅项目管理员且有选中项目) */}
          {user?.role === 'project_manager' && currentProjectId && (
            <Tooltip title={t('project.deleteProject')}>
              <Button
                type="text"
                size="small"
                icon={<DeleteOutlined />}
                onClick={handleDeleteProject}
                style={{ color: '#ff4d4f' }}
              />
            </Tooltip>
          )}
        </div>
      )}

      {/* 右侧：操作区域 */}
      <div style={{ display: 'flex', gap: '8px', alignItems: 'center', flex: 1, justifyContent: 'flex-end' }}>
        {/* 项目管理按钮 */}
        {(user?.role === 'project_manager' || user?.role === 'project_member') && (
          <Button
            type="text"
            size="small"
            icon={<ProjectOutlined />}
            onClick={() => {
              setActiveNav('projects');
              handleProjectsClick();
            }}
            style={{ ...getNavButtonStyle('projects'), fontSize: '13px', padding: '4px 8px' }}
          >
            {t('menu.projects')}
          </Button>
        )}
        
        {/* 菜单按钮组 */}
        {menuItems.map(item => (
          <Button
            key={item.key}
            type="text"
            size="small"
            icon={item.icon}
            onClick={() => {
              setActiveNav(item.key);
              item.action();
            }}
            style={{ ...getNavButtonStyle(item.key), fontSize: '13px', padding: '4px 8px' }}
          >
            {item.label}
          </Button>
        ))}

        {/* 语言切换组件 */}
        <LanguageSwitch variant="dropdown" showLabel={false} />

        {/* 用户昵称显示 - 紧靠退出按钮左侧 */}
        {displayName && (
          <Typography.Text
            style={{
              color: 'rgba(0, 0, 0, 0.85)',
              fontSize: '16px',
              fontWeight: 500,
              marginLeft: '8px',
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              minWidth: '64px',
              maxWidth: '120px',
            }}
            className="header-username"
          >
            {displayName}
          </Typography.Text>
        )}
        
        {/* 退出登录按钮 */}
        <Button
          type="text"
          danger
          icon={<LogoutOutlined />}
          onClick={handleLogout}
        >
          {t('common.logout')}
        </Button>
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
        onCancel={() => setDeleteModalVisible(false)}
        okText={t('common.confirm')}
        cancelText={t('common.cancel')}
        okType="danger"
      >
        <p>{t('project.confirmDelete', { name: getCurrentProjectName() })}</p>
      </Modal>
    </AntHeader>
  );
};

export default Header;
