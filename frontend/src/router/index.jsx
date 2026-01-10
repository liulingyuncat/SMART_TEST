import React, { useEffect, useState } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { Spin } from 'antd';
import LoginPage from '../pages/Login';
import Forbidden from '../pages/Forbidden';
import ProjectList from '../pages/ProjectList';
import ProjectDetail from '../pages/ProjectDetail';
import ProjectManagement from '../pages/ProjectManagement';
import UserManage from '../pages/UserManage';
import UserAssign from '../pages/UserAssign';
import Profile from '../pages/Profile';
import PromptManagement from '../pages/PromptManagement';
import RoleGuard from '../components/RoleGuard';
import MainLayout from '../components/MainLayout';
import { getProjects } from '../api/project';

// 智能默认路由组件
const DefaultRoute = () => {
  const { user } = useSelector(state => state.auth);
  const [firstProjectId, setFirstProjectId] = useState(null);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const fetchFirstProject = async () => {
      if (user?.role === 'project_manager' || user?.role === 'project_member') {
        try {
          const data = await getProjects();
          const projects = Array.isArray(data) ? data : [];
          if (projects.length > 0) {
            setFirstProjectId(projects[0].id);
          }
        } catch (error) {
          console.error('[DefaultRoute] Failed to fetch projects:', error);
        } finally {
          setLoading(false);
        }
      } else {
        setLoading(false);
      }
    };
    
    fetchFirstProject();
  }, [user?.role]);
  
  // 显示加载状态
  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <Spin size="large" />
      </div>
    );
  }
  
  // 根据用户角色跳转到不同的默认页面
  if (user?.role === 'system_admin') {
    return <Navigate to="/users" replace />;
  } else if (user?.role === 'project_manager' || user?.role === 'project_member') {
    // 跳转到项目管理页面（左侧列表+右侧详情）
    return <Navigate to="/projects" replace />;
  }
  
  // 默认跳转到用户管理页面
  return <Navigate to="/users" replace />;
};

const AppRouter = () => {
  return (
    <BrowserRouter>
      <Routes>
        {/* 公开路由 */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/403" element={<Forbidden />} />

        {/* 受保护路由 - 使用MainLayout包裹 */}
        <Route element={<MainLayout />}>
          {/* 项目管理(PM + PMemb) */}
          <Route 
            path="/projects" 
            element={
              <RoleGuard allowedRoles={['project_manager', 'project_member']}>
                <ProjectManagement />
              </RoleGuard>
            } 
          />

          {/* 项目详情(PM + PMemb) */}
          <Route 
            path="/projects/:id" 
            element={
              <RoleGuard allowedRoles={['project_manager', 'project_member']}>
                <ProjectDetail />
              </RoleGuard>
            } 
          />

          {/* 用户管理(SA + PM) */}
          <Route 
            path="/users" 
            element={
              <RoleGuard allowedRoles={['system_admin', 'project_manager']}>
                <UserManage />
              </RoleGuard>
            } 
          />

          {/* 人员分配(仅PM) */}
          <Route 
            path="/assign" 
            element={
              <RoleGuard allowedRoles={['project_manager']}>
                <UserAssign />
              </RoleGuard>
            } 
          />

          {/* 个人信息(所有角色) */}
          <Route 
            path="/profile" 
            element={
              <RoleGuard allowedRoles={['system_admin', 'project_manager', 'project_member']}>
                <Profile />
              </RoleGuard>
            } 
          />

          {/* 提示词管理(所有角色) */}
          <Route 
            path="/prompt-management" 
            element={
              <RoleGuard allowedRoles={['system_admin', 'project_manager', 'project_member']}>
                <PromptManagement />
              </RoleGuard>
            } 
          />

          {/* 默认路由 */}
          <Route path="/" element={<DefaultRoute />} />
        </Route>

        {/* 404路由 */}
        <Route path="*" element={<Navigate to="/403" replace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default AppRouter;
