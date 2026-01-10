import React, { useState, useEffect } from 'react';
import { Table, Button, Popconfirm, message, Space, Row, Col, Card, Empty, Input, Modal } from 'antd';
import { DownloadOutlined, DeleteOutlined, EditOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getVersionList, downloadVersion, deleteVersion, updateVersionRemark } from '../../../../api/manualCase';

/**
 * 版本管理Tab组件(支持三栏垂直布局:可配置左中右栏文档类型)
 * @param {string|number} projectId - 项目ID
 * @param {string} leftDocType - 左栏文档类型(默认'overall')
 * @param {string} middleDocType - 中栏文档类型(可选,默认undefined)
 * @param {string} rightDocType - 右栏文档类型(默认'change')
 * @param {string} leftTitle - 左栏标题(默认'整体用例版本')
 * @param {string} middleTitle - 中栏标题(默认'受入用例版本')
 * @param {string} rightTitle - 右栏标题(默认'变更用例版本')
 * @param {object} apiModule - API模块,包含getVersionList/downloadVersion/deleteVersion方法(默认使用manualCase API)
 */
const VersionManagementTab = ({ 
  projectId, 
  leftDocType = 'overall',
  middleDocType,
  rightDocType,
  leftTitle = '整体用例版本',
  middleTitle = '受入用例版本',
  rightTitle = '变更用例版本',
  leftTitleKey,
  middleTitleKey,
  rightTitleKey,
  apiModule 
}) => {
  const { t } = useTranslation();
  
  // T44: 移除版本管理标题显示
  // const displayLeftTitle = leftTitleKey ? t(leftTitleKey) : leftTitle;
  const displayMiddleTitle = middleTitleKey ? t(middleTitleKey) : middleTitle;
  const displayRightTitle = rightTitleKey ? t(rightTitleKey) : rightTitle;
  
  // 使用传入的apiModule或默认的manualCase API
  const api = apiModule || { getVersionList, downloadVersion, deleteVersion, updateVersionRemark };
  
  const [leftVersions, setLeftVersions] = useState([]);
  const [middleVersions, setMiddleVersions] = useState([]);
  const [rightVersions, setRightVersions] = useState([]);
  const [leftLoading, setLeftLoading] = useState(false);
  const [middleLoading, setMiddleLoading] = useState(false);
  const [rightLoading, setRightLoading] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [editingVersion, setEditingVersion] = useState(null);
  const [editingRemark, setEditingRemark] = useState('');

  useEffect(() => {
    loadLeftVersions();
    loadMiddleVersions();
    loadRightVersions();
  }, [projectId, leftDocType, middleDocType, rightDocType]);

  const loadLeftVersions = async () => {
    if (!projectId) return;
    
    setLeftLoading(true);
    try {
      const data = await api.getVersionList(projectId, leftDocType);
      // 处理可能的数据格式: [] 或 {versions: []} 或 {data: {versions: []}}
      const versionArray = Array.isArray(data) ? data : (data?.versions || data?.data?.versions || []);
      setLeftVersions(versionArray);
    } catch (error) {
      console.error('加载左栏版本失败:', error);
      message.error('加载版本列表失败', 3);
    } finally {
      setLeftLoading(false);
    }
  };

  const loadRightVersions = async () => {
    if (!rightDocType || !projectId) return;
    
    setRightLoading(true);
    try {
      const data = await api.getVersionList(projectId, rightDocType);
      const versionArray = Array.isArray(data) ? data : (data?.versions || data?.data?.versions || []);
      setRightVersions(versionArray);
    } catch (error) {
      console.error('加载右栏版本失败:', error);
      message.error('加载版本列表失败', 3);
    } finally {
      setRightLoading(false);
    }
  };

  const loadMiddleVersions = async () => {
    if (!middleDocType || !projectId) return;
    
    setMiddleLoading(true);
    try {
      const data = await api.getVersionList(projectId, middleDocType);
      const versionArray = Array.isArray(data) ? data : (data?.versions || data?.data?.versions || []);
      setMiddleVersions(versionArray);
    } catch (error) {
      console.error('加载中间栏版本失败:', error);
      message.error('加载版本列表失败', 3);
    } finally {
      setMiddleLoading(false);
    }
  };

  const handleDownload = async (versionID) => {
    try {
      await api.downloadVersion(projectId, versionID);
      message.success('版本文件下载成功', 3);
    } catch (error) {
      console.error('下载版本失败:', error);
      message.error('下载版本失败', 3);
    }
  };

  const handleDelete = async (versionID, side) => {
    try {
      await api.deleteVersion(projectId, versionID);
      message.success('版本删除成功', 3);
      // 根据栏位刷新对应的列表
      if (side === 'left') {
        loadLeftVersions();
      } else if (side === 'middle') {
        loadMiddleVersions();
      } else if (side === 'right') {
        loadRightVersions();
      }
    } catch (error) {
      console.error('删除版本失败:', error);
      message.error('删除版本失败', 3);
    }
  };

  const handleEditRemark = (record) => {
    setEditingVersion(record);
    setEditingRemark(record.remark || '');
    setEditModalVisible(true);
  };

  const handleSaveRemark = async () => {
    if (!editingVersion) return;
    
    try {
      await api.updateVersionRemark(projectId, editingVersion.id, editingRemark);
      message.success('备注更新成功', 3);
      setEditModalVisible(false);
      // 刷新列表
      loadLeftVersions();
      loadMiddleVersions();
      loadRightVersions();
    } catch (error) {
      console.error('更新备注失败:', error);
      message.error('更新备注失败', 3);
    }
  };

  const handleCancelEdit = () => {
    setEditModalVisible(false);
    setEditingVersion(null);
    setEditingRemark('');
  };

  const formatFileSize = (bytes) => {
    if (!bytes) return '-';
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
  };

  const formatDate = (dateStr) => {
    if (!dateStr) return '-';
    try {
      const date = new Date(dateStr);
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      });
    } catch {
      return dateStr;
    }
  };

  const getColumns = (side) => [
    {
      title: t('manualTest.versionId'),
      key: 'index',
      width: 80,
      render: (_, record, index) => {
        // 显示序号而不是UUID
        return index + 1;
      },
    },
    {
      title: t('manualTest.versionFilename'),
      dataIndex: 'filename',
      key: 'filename',
      ellipsis: {
        showTitle: true,
      },
      render: (text, record) => {
        // 如果是接口用例版本,显示4个文件名
        if (record.filename_role1) {
          return (
            <div>
              <div>ROLE1: {record.filename_role1}</div>
              <div>ROLE2: {record.filename_role2}</div>
              <div>ROLE3: {record.filename_role3}</div>
              <div>ROLE4: {record.filename_role4}</div>
            </div>
          );
        }
        // 其他类型显示单个文件名
        return text || '-';
      },
    },
    {
      title: t('manualTest.remark'),
      dataIndex: 'remark',
      key: 'remark',
      ellipsis: {
        showTitle: true,
      },
      width: '20%',
      render: (text) => text || '-',
    },
    {
      title: t('manualTest.operation'),
      key: 'action',
      width: 280,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEditRemark(record)}
          >
            {t('manualTest.editRemark')}
          </Button>
          <Button
            type="link"
            size="small"
            icon={<DownloadOutlined />}
            onClick={() => handleDownload(record.id)}
          >
            {t('manualTest.download')}
          </Button>
          <Popconfirm
            title={t('manualTest.confirmDelete')}
            onConfirm={() => handleDelete(record.id, side)}
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              {t('manualTest.delete')}
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: '16px' }}>
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
      <Row gutter={[16, 16]}>
        <Col span={24}>
          {/* T44: 移除版本管理标题 */}
          <Card bordered={false}>
            <Table
              columns={getColumns('left')}
              dataSource={leftVersions}
              rowKey="id"
              loading={leftLoading}
              scroll={{ x: 'max-content' }}
              pagination={{
                pageSize: 5,
                showSizeChanger: true,
                pageSizeOptions: ['5', '10', '20'],
                showTotal: (total) => `${t('common.total')} ${total} ${t('manualTest.versions')}`,
              }}
              locale={{
                emptyText: <Empty description={t('manualTest.noVersions')} />
              }}
            />
          </Card>
        </Col>
        {middleDocType && (
          <Col span={24}>
            <Card title={displayMiddleTitle} bordered={false}>
              <Table
                columns={getColumns('middle')}
                dataSource={middleVersions}
                rowKey="id"
                loading={middleLoading}
                scroll={{ x: 'max-content' }}
                pagination={{
                  pageSize: 5,
                  showSizeChanger: true,
                  pageSizeOptions: ['5', '10', '20'],
                  showTotal: (total) => `${t('common.total')} ${total} ${t('manualTest.versions')}`,
                }}
                locale={{
                  emptyText: <Empty description={t('manualTest.noVersions')} />
                }}
              />
            </Card>
          </Col>
        )}
        {rightDocType && (
          <Col span={24}>
            <Card title={displayRightTitle} bordered={false}>
              <Table
                columns={getColumns('right')}
                dataSource={rightVersions}
                rowKey="id"
                loading={rightLoading}
                scroll={{ x: 'max-content' }}
                pagination={{
                  pageSize: 5,
                  showSizeChanger: true,
                  pageSizeOptions: ['5', '10', '20'],
                  showTotal: (total) => `${t('common.total')} ${total} ${t('manualTest.versions')}`,
                }}
                locale={{
                  emptyText: <Empty description={t('manualTest.noVersions')} />
                }}
              />
            </Card>
          </Col>
        )}
      </Row>
    </div>
  );
};

export default VersionManagementTab;
