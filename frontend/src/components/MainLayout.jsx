import { Layout } from 'antd';
import { Outlet } from 'react-router-dom';
import Header from './Header';

const { Content } = Layout;

const MainLayout = () => {
  return (
    <Layout style={{ height: '100%' }}>
      {/* 顶部导航栏 */}
      <Header />

      {/* 内容区域 */}
      <Content
        style={{
          background: '#f0f2f5',
          minHeight: 'calc(100vh - 64px)',
          overflow: 'hidden',
        }}
      >
        <Outlet />
      </Content>
    </Layout>
  );
};

export default MainLayout;
