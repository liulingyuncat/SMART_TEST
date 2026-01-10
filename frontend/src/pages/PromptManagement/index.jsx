import React, { useState, useEffect } from 'react';
import { Layout, Tabs, Spin, message, Button, Space } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import ToolList from './ToolList';
import PromptTable from './PromptTable';
import PromptDetail from './PromptDetail';
import { fetchPrompts, refreshPrompts } from '../../api/prompt';
import './index.css';

const { Sider, Content } = Layout;

const PromptManagement = () => {
  const { t } = useTranslation();
  const { user } = useSelector((state) => state.auth);
  const [activeTab, setActiveTab] = useState('system');
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [prompts, setPrompts] = useState([]);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);
  
  // 详情页状态
  const [detailView, setDetailView] = useState(null); // null | 'create' | { mode: 'view' | 'edit', prompt }
  const [selectedPrompt, setSelectedPrompt] = useState(null);

  // 检查用户是否有权限刷新提示词缓存
  const canRefreshCache = user && (user.role === 'SystemAdmin' || user.role === 'ProjectManager');

  // 刷新MCP提示词缓存
  const handleRefreshCache = async () => {
    setRefreshing(true);
    try {
      const scope = activeTab === 'system' ? 'system' : activeTab === 'project' ? 'project' : 'user';
      await refreshPrompts(scope);
      message.success(t('prompts.cacheRefreshed') || '提示词缓存已刷新');
      // 刷新后重新加载列表
      loadPrompts(currentPage);
    } catch (error) {
      message.error(t('prompts.cacheRefreshFailed') || '刷新缓存失败');
    } finally {
      setRefreshing(false);
    }
  };

  // 加载提示词列表
  const loadPrompts = async (page = 1) => {
    setLoading(true);
    try {
      const scopeMap = {
        system: 'system',
        project: 'project',
        user: 'user',
      };
      const params = {
        scope: scopeMap[activeTab],
        page,
        page_size: pageSize,
      };
      console.log('[PromptManagement] 查询参数:', params, '当前Tab:', activeTab);
      const data = await fetchPrompts(params);
      console.log('[PromptManagement] 查询结果:', data);
      setPrompts(data.items || []);
      setTotal(data.total || 0);
      setCurrentPage(page);
    } catch (error) {
      message.error(t('prompts.loadFailed') || '加载失败');
      setPrompts([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // Tab切换时重新加载
  useEffect(() => {
    loadPrompts(1);
  }, [activeTab]);

  // 操作处理函数
  const handleView = (prompt) => {
    setSelectedPrompt(prompt);
    setDetailView('view');
  };

  const handleCreate = () => {
    setSelectedPrompt(null);
    setDetailView('create');
  };

  const handleEdit = (prompt) => {
    setSelectedPrompt(prompt);
    setDetailView('edit');
  };

  const handleBack = () => {
    setDetailView(null);
    setSelectedPrompt(null);
  };

  const handleDelete = () => {
    // 删除成功后刷新列表并更新MCP缓存
    loadPrompts(currentPage);
    if (canRefreshCache) {
      // 异步刷新缓存，不阻塞UI
      refreshPrompts(activeTab).catch(err => {
        console.warn('Cache refresh failed:', err);
      });
    }
  };

  const handleSuccess = () => {
    setDetailView(null);
    setSelectedPrompt(null);
    loadPrompts(currentPage);
    // 创建/更新成功后自动刷新MCP缓存
    if (canRefreshCache) {
      // 异步刷新缓存，不阻塞UI
      refreshPrompts(activeTab).catch(err => {
        console.warn('Cache refresh failed:', err);
      });
    }
  };

  const handlePageChange = (page) => {
    loadPrompts(page);
  };

  const tabItems = [
    {
      key: 'system',
      label: t('prompts.systemPrompts'),
      children: detailView ? null : (
        <PromptTable
          prompts={prompts}
          loading={loading}
          total={total}
          currentPage={currentPage}
          pageSize={pageSize}
          onView={handleView}
          onCreate={handleCreate}
          onEdit={handleEdit}
          onDelete={handleDelete}
          onPageChange={handlePageChange}
          scope="system"
        />
      ),
    },
    {
      key: 'project',
      label: t('prompts.projectPrompts'),
      children: detailView ? null : (
        <PromptTable
          prompts={prompts}
          loading={loading}
          total={total}
          currentPage={currentPage}
          pageSize={pageSize}
          onView={handleView}
          onCreate={handleCreate}
          onEdit={handleEdit}
          onDelete={handleDelete}
          onPageChange={handlePageChange}
          scope="project"
        />
      ),
    },
    {
      key: 'user',
      label: t('prompts.userPrompts'),
      children: detailView ? null : (
        <PromptTable
          prompts={prompts}
          loading={loading}
          total={total}
          currentPage={currentPage}
          pageSize={pageSize}
          onView={handleView}
          onCreate={handleCreate}
          onEdit={handleEdit}
          onDelete={handleDelete}
          onPageChange={handlePageChange}
          scope="user"
        />
      ),
    },
  ];

  return (
    <Layout className="prompt-management">
      <Sider width="30%" className="prompt-sider" style={{ overflow: 'hidden' }}>
        <ToolList />
      </Sider>
      <Content className="prompt-content">
        {detailView ? (
          <PromptDetail
            prompt={selectedPrompt}
            mode={detailView}
            scope={activeTab}
            onBack={handleBack}
            onSuccess={handleSuccess}
          />
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            {canRefreshCache && (
              <div style={{ padding: '12px 0', borderBottom: '1px solid #f0f0f0' }}>
                <Space>
                  <span style={{ fontSize: '12px', color: '#666' }}>
                    {t('prompts.cacheHint') || '提示词缓存会自动更新，你也可以手动刷新'}
                  </span>
                  <Button
                    type="primary"
                    size="small"
                    icon={<ReloadOutlined />}
                    loading={refreshing}
                    onClick={handleRefreshCache}
                  >
                    {t('prompts.refreshCache') || '刷新MCP缓存'}
                  </Button>
                </Space>
              </div>
            )}
            <div style={{ flex: 1, overflow: 'auto' }}>
              <Tabs
                activeKey={activeTab}
                onChange={setActiveTab}
                items={tabItems}
                style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
              />
            </div>
          </div>
        )}
      </Content>
    </Layout>
  );
};

export default PromptManagement;
