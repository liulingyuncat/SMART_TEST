import { Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';

const RoleGuard = ({ allowedRoles, children }) => {
  const { user, isAuthenticated } = useSelector(state => state.auth);

  // 调试信息
  console.log('RoleGuard - isAuthenticated:', isAuthenticated);
  console.log('RoleGuard - user:', user);
  console.log('RoleGuard - allowedRoles:', allowedRoles);

  // 未登录,跳转登录页
  if (!isAuthenticated || !user) {
    console.log('RoleGuard - Not authenticated, redirecting to login');
    return <Navigate to="/login" replace />;
  }

  // 角色不匹配,跳转403页面
  if (!allowedRoles.includes(user.role)) {
    console.log('RoleGuard - Role not allowed, redirecting to 403');
    console.log('RoleGuard - user.role:', user.role);
    return <Navigate to="/403" replace />;
  }

  // 通过验证,渲染子组件
  console.log('RoleGuard - Access granted');
  return children;
};

export default RoleGuard;
