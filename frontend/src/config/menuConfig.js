import {
  ProjectOutlined,
  BarChartOutlined,
  UserOutlined,
  TeamOutlined,
  IdcardOutlined,
} from '@ant-design/icons';

export const menuConfig = [
  {
    key: 'projects',
    label: 'menu.projects',
    icon: ProjectOutlined,
    path: '/projects',
    roles: ['project_manager', 'project_member'],
  },
  {
    key: 'statistics',
    label: 'menu.statistics',
    icon: BarChartOutlined,
    path: '/statistics',
    roles: ['project_manager', 'project_member'],
  },
  {
    key: 'users',
    label: 'menu.users',
    icon: UserOutlined,
    path: '/users',
    roles: ['system_admin', 'project_manager'],
  },
  {
    key: 'assign',
    label: 'menu.assign',
    icon: TeamOutlined,
    path: '/assign',
    roles: ['project_manager'],
  },
  {
    key: 'profile',
    label: 'menu.profile',
    icon: IdcardOutlined,
    path: '/profile',
    roles: ['system_admin', 'project_manager', 'project_member'],
  },
];
