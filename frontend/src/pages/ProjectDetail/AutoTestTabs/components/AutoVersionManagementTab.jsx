import React, { useState, useEffect } from 'react';
import { Table, Button, message, Popconfirm, Input, Space, Modal } from 'antd';
import { DownloadOutlined, DeleteOutlined, EditOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getAutoVersions, downloadAutoVersion, deleteAutoVersion, updateAutoVersionRemark } from '../../../../api/autoCase';
import './AutoVersionManagementTab.css';

/**
 * 自动化测试用例版本管理组件
 * 显示版本列表,支持下载、删除、备注编辑
 */
const AutoVersionManagementTab = ({ projectId }) => {
  const { t } = useTranslation();
  const [versions, setVersions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [editingVersion, setEditingVersion] = useState(null);
  const [editingRemark, setEditingRemark] = useState('');

  // 加载版本列表
  const loadVersions = async (page = 1, size = 10) => {
    setLoading(true);
    try {
      const result = await getAutoVersions(projectId, page, size);
      // 拦截器已返回 {versions, total, page, size}
      setVersions(result.versions || []);
      setPagination({
        current: result.page || page,
        pageSize: result.size || size,
        total: result.total || 0
      });
    } catch (error) {
      message.error('加载版本列表失败');
      console.error('[AutoVersionManagementTab] 加载失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 下载版本压缩包
  const handleDownload = async (versionId) => {
    try {
      message.loading({ content: '正在准备下载...', key: 'download' });
      await downloadAutoVersion(projectId, versionId);
      message.success({ content: '下载成功', key: 'download' });
    } catch (error) {
      message.error({ content: '下载失败', key: 'download' });
      console.error('[AutoVersionManagementTab] 下载失败:', error);
    }
  };

  // 删除版本
  const handleDelete = async (versionId) => {
    try {
      await deleteAutoVersion(projectId, versionId);
      message.success('删除成功');
      loadVersions(pagination.current, pagination.pageSize);
    } catch (error) {
      message.error('删除失败');
      console.error('[AutoVersionManagementTab] 删除失败:', error);
    }
  };

  // 打开编辑备注对话框
  const handleEditRemark = (record) => {
    setEditingVersion(record);
    setEditingRemark(record.remark || '');
    setEditModalVisible(true);
  };

  // 保存备注
  const handleSaveRemark = async () => {
    if (!editingVersion) return;
    
    try {
      await updateAutoVersionRemark(projectId, editingVersion.version_id, editingRemark);
      message.success('备注更新成功');
      setEditModalVisible(false);
      loadVersions(pagination.current, pagination.pageSize);
    } catch (error) {
      message.error('备注更新失败');
      console.error('[AutoVersionManagementTab] 更新备注失败:', error);
    }
  };

  // 取消编辑
  const handleCancelEdit = () => {
    setEditModalVisible(false);
    setEditingVersion(null);
    setEditingRemark('');
  };

  // 初始化加载
  useEffect(() => {
    if (projectId) {
      loadVersions();
    }
  }, [projectId]);

  // 分页变化处理
  const handleTableChange = (newPagination) => {
    loadVersions(newPagination.current, newPagination.pageSize);
  };

  // 格式化文件名(移除括号内容)
  const formatFilename = (filename) => {
    if (!filename) return '';
    // 移除括号及其内容,如 "Test_xxx.xlsx (123 KB, 45 条用例)" -> "Test_xxx.xlsx"
    return filename.replace(/\s*\([^)]*\)\s*/g, '').trim();
  };

  // 表格列定义
  const columns = [
    {
      title: t('manualTest.versionId'),
      key: 'index',
      width: 80,
      render: (_, __, index) => {
        // 计算当前页的序号
        return (pagination.current - 1) * pagination.pageSize + index + 1;
      },
    },
    {
      title: t('manualTest.versionFilename'),
      key: 'files',
      ellipsis: true,
      render: (_, record) => (
        <div>
          {record.files && record.files.map((file, index) => (
            <div key={index}>
              {formatFilename(file.filename)}
            </div>
          ))}
        </div>
      ),
    },
    {
      title: t('manualTest.remark'),
      dataIndex: 'remark',
      key: 'remark',
      ellipsis: true,
      width: 200,
      render: (text) => text || '-',
    },
    {
      title: t('manualTest.operation'),
      key: 'actions',
      width: 200,
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEditRemark(record)}
          >
            {t('manualTest.editRemark')}
          </Button>
          <Button
            type="link"
            icon={<DownloadOutlined />}
            onClick={() => handleDownload(record.version_id)}
          >
            {t('manualTest.download')}
          </Button>
          <Popconfirm
            title={t('manualTest.confirmDelete')}
            onConfirm={() => handleDelete(record.version_id)}
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              {t('manualTest.delete')}
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="auto-version-management">
      <Modal
        title={t('manualTest.editRemarkTitle')}
        open={editModalVisible}
        onOk={handleSaveRemark}
        onCancel={handleCancelEdit}
      >
        <Input.TextArea
          rows={4}
          value={editingRemark}
          onChange={(e) => setEditingRemark(e.target.value)}
          placeholder={t('manualTest.remarkPlaceholder')}
        />
      </Modal>
      <Table
        columns={columns}
        dataSource={versions}
        loading={loading}
        rowKey="version_id"
        pagination={{
          ...pagination,
          showSizeChanger: true,
          pageSizeOptions: ['5', '10', '20'],
          showTotal: (total) => `${t('common.total')} ${total} ${t('manualTest.versions')}`,
        }}
        onChange={handleTableChange}
      />
    </div>
  );
};

export default AutoVersionManagementTab;
