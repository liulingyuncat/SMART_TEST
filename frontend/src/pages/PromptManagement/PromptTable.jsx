import React, { useState } from 'react';
import { Table, Button, Space, Input, Popconfirm, message, Tag } from 'antd';
import { PlusOutlined, EyeOutlined, EditOutlined, DeleteOutlined, SearchOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { deletePrompt } from '../../api/prompt';

const PromptTable = ({
  prompts,
  loading,
  total,
  currentPage,
  pageSize,
  onView,
  onCreate,
  onEdit,
  onDelete,
  onPageChange,
  scope,
}) => {
  const { t } = useTranslation();
  const { user } = useSelector((state) => state.auth);
  const [searchText, setSearchText] = useState('');

  const handleDelete = async (record) => {
    try {
      await deletePrompt(record.id);
      message.success(t('prompts.deleteSuccess'));
      onDelete();
    } catch (error) {
      const errorMsg = error.response?.data?.message || error.message;
      message.error(errorMsg || t('prompts.deleteFailed'));
    }
  };

  const canEdit = (record) => {
    // 系统提示词全员只读
    if (record.scope === 'system') return false;
    // 个人提示词仅自己可见可编辑
    if (record.scope === 'user') return record.user_id === user?.id;
    // 全员提示词仅系统管理员可以编辑
    if (record.scope === 'project') return user?.role === 'system_admin';
    return false;
  };

  const filteredPrompts = prompts.filter((prompt) =>
    prompt.name.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: t('prompts.name'),
      dataIndex: 'name',
      key: 'name',
      width: '15%',
      ellipsis: true,
    },
    {
      title: t('prompts.version'),
      dataIndex: 'version',
      key: 'version',
      width: '8%',
    },
    {
      title: t('prompts.description'),
      dataIndex: 'description',
      key: 'description',
      width: '25%',
      ellipsis: true,
    },
    {
      title: t('prompts.parameterCount'),
      key: 'arguments',
      width: '10%',
      render: (_, record) => {
        const count = record.arguments?.length || 0;
        return <Tag color={count > 0 ? 'blue' : 'default'}>{count}</Tag>;
      },
    },
    {
      title: t('prompts.scope'),
      dataIndex: 'scope',
      key: 'scope',
      width: '10%',
      render: (scope) => {
        const colorMap = { system: 'green', project: 'blue', user: 'orange' };
        const labelMap = {
          system: t('prompts.systemPrompts'),
          project: t('prompts.projectPrompts'),
          user: t('prompts.userPrompts'),
        };
        return <Tag color={colorMap[scope]}>{labelMap[scope]}</Tag>;
      },
    },
    {
      title: t('prompts.actions'),
      key: 'actions',
      width: '10%',
      render: (_, record) => (
        <Button
          type="link"
          size="small"
          icon={<EyeOutlined />}
          onClick={() => onView(record)}
        />
      ),
    },
  ];

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <div style={{ 
        marginBottom: 0, 
        display: 'flex', 
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '12px 16px',
        borderBottom: '1px solid #f0f0f0',
        background: '#fff',
      }}>
        <Input
          placeholder={t('prompts.namePlaceholder')}
          prefix={<SearchOutlined />}
          style={{ width: 300 }}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          allowClear
        />
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={onCreate}
          disabled={scope === 'system' || (scope === 'project' && user?.role !== 'system_admin')}
        >
          {t('prompts.create')}
        </Button>
      </div>

      <div style={{ flex: 1, overflow: 'auto' }}>
        <Table
          columns={columns}
          dataSource={filteredPrompts}
          loading={loading}
          rowKey="id"
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            total: total,
            onChange: onPageChange,
            showSizeChanger: false,
            showTotal: (total) => `${t('common.total')} ${total} ${t('common.items')}`,
          }}
          locale={{
            emptyText: t('prompts.noData'),
          }}
          style={{ background: '#fff' }}
        />
      </div>
    </div>
  );
};

export default PromptTable;
