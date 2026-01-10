import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Input, message, Popconfirm } from 'antd';
import { DownloadOutlined, DeleteOutlined, EditOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import {
  getWebVersions,
  downloadWebVersion,
  deleteWebVersion,
  updateWebVersionRemark
} from '../../../../api/autoCase';
import {
  getApiVersionList,
  downloadApiVersion,
  deleteApiVersion,
  updateApiVersionRemark
} from '../../../../api/apiCase';
import './WebVersionList.css';

const { TextArea } = Input;

/**
 * 通用版本列表组件
 * @param {string} projectId - 项目ID
 * @param {string} apiModule - API模块类型: 'web-cases'(默认) 或 'api-cases'
 * @param {function} onVersionDeleted - 版本删除后的回调
 */
const WebVersionList = ({ projectId, apiModule = 'web-cases', onVersionDeleted }) => {
  const { t } = useTranslation();
  const [versions, setVersions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [editingVersion, setEditingVersion] = useState(null);
  const [editingRemark, setEditingRemark] = useState('');

  // 根据apiModule选择对应的API方法
  const apiMethods = apiModule === 'api-cases' ? {
    getVersions: getApiVersionList,
    downloadVersion: downloadApiVersion,
    deleteVersion: deleteApiVersion,
    updateRemark: updateApiVersionRemark
  } : {
    getVersions: getWebVersions,
    downloadVersion: downloadWebVersion,
    deleteVersion: deleteWebVersion,
    updateRemark: updateWebVersionRemark
  };

  useEffect(() => {
    loadVersions();
  }, [projectId, apiModule]);

  const loadVersions = async (page = 1, size = 10) => {
    setLoading(true);
    try {
      const result = await apiMethods.getVersions(projectId, page, size);
      setVersions(result.versions || []);
      setPagination({
        current: result.page || page,
        pageSize: result.size || size,
        total: result.total || 0,
      });
    } catch (error) {
      message.error(t('web_version.loadFailed'));
      console.error('Failed to load versions:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = async (versionId) => {
    try {
      await apiMethods.downloadVersion(projectId, versionId);
      message.success(t('web_version.downloadSuccess'));
    } catch (error) {
      message.error(t('web_version.downloadFailed'));
      console.error('Failed to download version:', error);
    }
  };

  const handleDelete = async (versionId) => {
    try {
      await apiMethods.deleteVersion(projectId, versionId);
      message.success(t('web_version.deleteSuccess'));
      loadVersions(pagination.current, pagination.pageSize);
      if (onVersionDeleted) {
        onVersionDeleted();
      }
    } catch (error) {
      message.error(t('web_version.deleteFailed'));
      console.error('Failed to delete version:', error);
    }
  };

  const handleEditRemark = (record) => {
    setEditingVersion(record);
    setEditingRemark(record.remark || '');
    setEditModalVisible(true);
  };

  const handleSaveRemark = async () => {
    try {
      const versionId = apiModule === 'api-cases' ? editingVersion.id : editingVersion.version_id;
      await apiMethods.updateRemark(projectId, versionId, editingRemark);
      message.success(t('web_version.updateRemarkSuccess'));
      setEditModalVisible(false);
      loadVersions(pagination.current, pagination.pageSize);
    } catch (error) {
      message.error(t('web_version.updateRemarkFailed'));
      console.error('Failed to update version remark:', error);
    }
  };

  const handleTableChange = (newPagination) => {
    loadVersions(newPagination.current, newPagination.pageSize);
  };

  // 根据apiModule获取版本ID和文件名字段
  const getVersionId = (record) => {
    return apiModule === 'api-cases' ? record.id : record.version_id;
  };

  const getFilename = (record) => {
    if (apiModule === 'api-cases') {
      // API用例版本：优先使用xlsx_filename，兼容旧版本的CSV文件名
      return record.xlsx_filename || record.filename_role1 || '-';
    }
    return record.zip_filename;
  };

  const columns = [
    {
      title: t('web_version.versionId'),
      dataIndex: apiModule === 'api-cases' ? 'id' : 'version_id',
      key: 'version_id',
      width: 100,
      render: (_, __, index) => (pagination.current - 1) * pagination.pageSize + index + 1,
    },
    {
      title: t('web_version.filename'),
      key: 'filename',
      ellipsis: true,
      render: (_, record) => getFilename(record) || '-',
    },
    {
      title: t('web_version.remark'),
      dataIndex: 'remark',
      key: 'remark',
      width: 200,
      ellipsis: true,
      render: (text) => text || '-',
    },
    {
      title: t('web_version.actions'),
      key: 'actions',
      width: 200,
      render: (_, record) => {
        const versionId = getVersionId(record);
        return (
          <div className="action-buttons">
            <Button
              type="link"
              icon={<EditOutlined />}
              onClick={() => handleEditRemark(record)}
            >
              {t('common.edit')}
            </Button>
            <Button
              type="link"
              icon={<DownloadOutlined />}
              onClick={() => handleDownload(versionId)}
            >
              {t('common.download')}
            </Button>
            <Popconfirm
              title={t('web_version.deleteConfirm')}
              onConfirm={() => handleDelete(versionId)}
              okText={t('common.confirm')}
              cancelText={t('common.cancel')}
            >
              <Button type="link" danger icon={<DeleteOutlined />}>
                {t('common.delete')}
              </Button>
            </Popconfirm>
          </div>
        );
      },
    },
  ];

  return (
    <div className="web-version-list">
      <Table
        columns={columns}
        dataSource={versions}
        loading={loading}
        rowKey={apiModule === 'api-cases' ? 'id' : 'version_id'}
        pagination={{
          ...pagination,
          showSizeChanger: true,
          showQuickJumper: true,
          pageSizeOptions: ['10', '20', '50', '100'],
          showTotal: (total) => t('web_version.totalVersions', { total }),
        }}
        onChange={handleTableChange}
      />

      <Modal
        title={t('web_version.editRemark')}
        open={editModalVisible}
        onOk={handleSaveRemark}
        onCancel={() => setEditModalVisible(false)}
        okText={t('common.save')}
        cancelText={t('common.cancel')}
      >
        <TextArea
          value={editingRemark}
          onChange={(e) => setEditingRemark(e.target.value)}
          placeholder={t('web_version.remarkPlaceholder')}
          maxLength={200}
          rows={4}
          showCount
        />
      </Modal>
    </div>
  );
};

export default WebVersionList;
