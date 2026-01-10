import { Menu } from 'antd';
import { Link, useLocation } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import { menuConfig } from '../config/menuConfig';

const Sidebar = () => {
  const location = useLocation();
  const { user } = useSelector(state => state.auth);
  const { t } = useTranslation();

  // 根据角色过滤菜单
  const visibleMenus = menuConfig.filter(menu =>
    menu.roles.includes(user?.role)
  );

  // 生成菜单项(使用t()翻译label)
  const menuItems = visibleMenus.map(menu => ({
    key: menu.key,
    icon: <menu.icon />,
    label: <Link to={menu.path}>{t(menu.label)}</Link>,
  }));

  return (
    <Menu
      mode="inline"
      selectedKeys={[location.pathname.split('/')[1] || 'projects']}
      items={menuItems}
      style={{ height: '100%', borderRight: 0 }}
    />
  );
};

export default Sidebar;
