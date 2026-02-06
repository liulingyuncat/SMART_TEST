import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Table,
  Button,
  Space,
  Tag,
  Input,
  Select,
  Popconfirm,
  message,
  Upload,
  Modal,
  Tooltip,
  Dropdown,
  Menu,
} from 'antd';
import {
  PlusOutlined,
  ImportOutlined,
  ExportOutlined,
  DownloadOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import {
  fetchDefects,
  createDefect,
  updateDefect,
  deleteDefect,
  importDefects,
  exportDefects,
  downloadDefectTemplate,
} from '../../../api/defect';
import {
  DEFECT_STATUS_COLORS,
  DEFECT_STATUS_I18N_KEYS,
  DEFECT_PRIORITY_COLORS,
  DEFECT_PRIORITY_I18N_KEYS,
  DEFECT_SEVERITY_COLORS,
  DEFECT_SEVERITY_I18N_KEYS,
  getStatusOptions,
  getPriorityOptions,
  getSeverityOptions,
} from '../../../constants/defect';
import DefectFormModal from './DefectFormModal';
import DefectDetailModal from './DefectDetailModal';

const { Search } = Input;
const { Option } = Select;

/**
 * 缺陷列表组件
 */
const DefectList = ({ projectId, subjects, phases }) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [defects, setDefects] = useState([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  
  // 筛选状态
  const [filters, setFilters] = useState({
    status: '',
    priority: '',
    severity: '',
    subject: '',
    phase: '',
    keyword: '',
  });
  
  // 模态框状态
  const [formModalVisible, setFormModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [currentDefect, setCurrentDefect] = useState(null);
  const [formMode, setFormMode] = useState('create'); // 'create' | 'edit'
  const [importModalVisible, setImportModalVisible] = useState(false);
  const [importResult, setImportResult] = useState(null);

  // 加载缺陷列表
  const loadDefects = useCallback(async () => {
    if (!projectId) return;
    setLoading(true);
    try {
      const params = {
        page,
        page_size: pageSize,
        ...filters,
      };
      // 移除空值
      Object.keys(params).forEach(key => {
        if (params[key] === '') delete params[key];
      });
      
      const response = await fetchDefects(projectId, params);
      setDefects(response.items || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('Failed to load defects:', error);
      message.error(t('message.loadFailed'));
    } finally {
      setLoading(false);
    }
  }, [projectId, page, pageSize, filters, t]);

  useEffect(() => {
    loadDefects();
  }, [loadDefects]);

  // 筛选变化
  const handleFilterChange = (key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
    setPage(1);
  };

  // 分页变化
  const handleTableChange = (pagination) => {
    setPage(pagination.current);
    setPageSize(pagination.pageSize);
  };

  // 创建缺陷
  const handleCreate = () => {
    setCurrentDefect(null);
    setFormMode('create');
    setFormModalVisible(true);
  };

  // 编辑缺陷
  const handleEdit = (record) => {
    setCurrentDefect(record);
    setFormMode('edit');
    setFormModalVisible(true);
  };

  // 查看详情
  const handleViewDetail = (record) => {
    setCurrentDefect(record);
    setDetailModalVisible(true);
  };

  // 删除缺陷
  const handleDelete = async (record) => {
    try {
      await deleteDefect(projectId, record.id);
      message.success(t('defect.deleteSuccess'));
      loadDefects();
    } catch (error) {
      console.error('Failed to delete defect:', error);
      message.error(t('message.deleteFailed'));
    }
  };

  // 表单提交
  const handleFormSubmit = async (values) => {
    try {
      if (formMode === 'create') {
        await createDefect(projectId, values);
        message.success(t('defect.createSuccess'));
      } else {
        await updateDefect(projectId, currentDefect.id, values);
        message.success(t('defect.updateSuccess'));
      }
      setFormModalVisible(false);
      loadDefects();
    } catch (error) {
      console.error('Failed to save defect:', error);
      message.error(t('message.saveFailed'));
    }
  };

  // 导入CSV
  const handleImport = async (file) => {
    try {
      const result = await importDefects(projectId, file);
      setImportResult(result);
      setImportModalVisible(true);
      loadDefects();
    } catch (error) {
      console.error('Failed to import defects:', error);
      message.error(t('message.importFailed'));
    }
    return false; // 阻止默认上传行为
  };

  // 导出CSV
  const handleExport = async (format) => {
    try {
      await exportDefects(projectId, format, filters, projectName);
      message.success(t('message.exportSuccess'));
    } catch (error) {
      console.error('Failed to export defects:', error);
      message.error(t('message.exportFailed'));
    }
  };

  // 确认导出
  const handleConfirmExport = async (format) => {
    await handleExport(format);
  };

  // 下载模板
  const handleDownloadTemplate = async (format) => {
    try {
      await downloadDefectTemplate(projectId, format);
    } catch (error) {
      console.error('Failed to download template:', error);
      message.error(t('message.downloadFailed'));
    }
  };

  // 确认下载模板
  const handleConfirmDownloadTemplate = async (format) => {
    await handleDownloadTemplate(format);
  };

  // 表格列定义
  const columns = [
    {
      title: t('defect.defectId'),
      dataIndex: 'defect_id',
      key: 'defect_id',
      width: 120,
      fixed: 'left',
      render: (text, record) => (
        <Button type="link" onClick={() => handleViewDetail(record)} className="defect-id-link" style={{ padding: 0 }}>
          {text}
        </Button>
      ),
    },
    {
      title: t('defect.summary'),
      dataIndex: 'summary',
      key: 'summary',
      ellipsis: true,
      render: (text) => (
        <Tooltip title={text}>
          <span className="defect-summary-cell">{text}</span>
        </Tooltip>
      ),
    },
    {
      title: t('defect.status'),
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => (
        <Tag color={DEFECT_STATUS_COLORS[status]}>
          {t(DEFECT_STATUS_I18N_KEYS[status])}
        </Tag>
      ),
    },
    {
      title: t('defect.priority'),
      dataIndex: 'priority',
      key: 'priority',
      width: 80,
      render: (priority) => (
        <Tag color={DEFECT_PRIORITY_COLORS[priority]}>
          {t(DEFECT_PRIORITY_I18N_KEYS[priority])}
        </Tag>
      ),
    },
    {
      title: t('defect.severity'),
      dataIndex: 'severity',
      key: 'severity',
      width: 80,
      render: (severity) => (
        <Tag color={DEFECT_SEVERITY_COLORS[severity]}>
          {t(DEFECT_SEVERITY_I18N_KEYS[severity])}
        </Tag>
      ),
    },
    {
      title: t('defect.subject'),
      dataIndex: 'subject',
      key: 'subject',
      width: 120,
      ellipsis: true,
    },
    {
      title: t('defect.phase'),
      dataIndex: 'phase',
      key: 'phase',
      width: 120,
      ellipsis: true,
    },
    {
      title: t('defect.caseId', 'Case ID'),
      dataIndex: 'case_id',
      key: 'case_id',
      width: 120,
      ellipsis: true,
      render: (text) => text || '-',
    },
    {
      title: t('defect.reporter'),
      dataIndex: ['created_by_user', 'username'],
      key: 'reporter',
      width: 100,
      ellipsis: true,
      render: (text, record) => record.detected_by || text || '-',
    },
    {
      title: t('defect.createdAt'),
      dataIndex: 'created_at',
      key: 'created_at',
      width: 160,
      render: (text) => text ? dayjs(text).format('YYYY-MM-DD') : '-',
    },
    {
      title: t('common.actions'),
      key: 'actions',
      width: 120,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Tooltip title={t('common.view')}>
            <Button
              type="link"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => handleViewDetail(record)}
            />
          </Tooltip>
          <Tooltip title={t('common.edit')}>
            <Button
              type="link"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
        </Space>
      ),
    },
  ];

  // 导出菜单
  const exportMenu = useMemo(() => (
    <Menu>
      <Menu.Item key="csv" onClick={() => handleConfirmExport('csv')}>
        CSV
      </Menu.Item>
      <Menu.Item key="xlsx" onClick={() => handleConfirmExport('xlsx')}>
        XLSX
      </Menu.Item>
    </Menu>
  ), []);

  return (
    <div className="defect-list-container">
      {/* 工具栏 */}
      <div className="defect-list-toolbar">
        <Space className="defect-list-actions">
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            {t('defect.create')}
          </Button>
          <Upload
            accept=".xlsx"
            showUploadList={false}
            beforeUpload={handleImport}
          >
            <Button icon={<ImportOutlined />}>{t('defect.import')}</Button>
          </Upload>
          <Dropdown overlay={exportMenu} placement="bottomRight">
            <Button icon={<ExportOutlined />}>{t('defect.export')}</Button>
          </Dropdown>
          <Button 
            icon={<DownloadOutlined />} 
            onClick={() => handleDownloadTemplate('xlsx')}
          >
            {t('defect.downloadTemplate')}
          </Button>
        </Space>
        <Button icon={<ReloadOutlined />} onClick={loadDefects}>
          {t('common.refresh') || '刷新'}
        </Button>
      </div>

      {/* 筛选行 */}
      <div className="defect-filter-row">
        <Select
          placeholder={t('defect.allStatus')}
          value={filters.status || undefined}
          onChange={(value) => handleFilterChange('status', value || '')}
          allowClear
          style={{ width: 140 }}
        >
          {getStatusOptions(t).map(opt => (
            <Option key={opt.value} value={opt.value}>{opt.label}</Option>
          ))}
        </Select>
        <Select
          placeholder={t('defect.allPriority')}
          value={filters.priority || undefined}
          onChange={(value) => handleFilterChange('priority', value || '')}
          allowClear
          style={{ width: 140 }}
        >
          {getPriorityOptions(t).map(opt => (
            <Option key={opt.value} value={opt.value}>{opt.label}</Option>
          ))}
        </Select>
        <Select
          placeholder={t('defect.allSeverity')}
          value={filters.severity || undefined}
          onChange={(value) => handleFilterChange('severity', value || '')}
          allowClear
          style={{ width: 140 }}
        >
          {getSeverityOptions(t).map(opt => (
            <Option key={opt.value} value={opt.value}>{opt.label}</Option>
          ))}
        </Select>
        <Select
          placeholder={t('defect.allSubject')}
          value={filters.subject || undefined}
          onChange={(value) => handleFilterChange('subject', value || '')}
          allowClear
          style={{ width: 160 }}
        >
          {subjects.map(s => (
            <Option key={s.id} value={s.name}>{s.name}</Option>
          ))}
        </Select>
        <Select
          placeholder={t('defect.allPhase')}
          value={filters.phase || undefined}
          onChange={(value) => handleFilterChange('phase', value || '')}
          allowClear
          style={{ width: 160 }}
        >
          {phases.map(p => (
            <Option key={p.id} value={p.name}>{p.name}</Option>
          ))}
        </Select>
        <Search
          placeholder={t('defect.keyword')}
          value={filters.keyword}
          onChange={(e) => handleFilterChange('keyword', e.target.value)}
          onSearch={() => loadDefects()}
          style={{ width: 200 }}
          allowClear
        />
      </div>

      {/* 缺陷表格 */}
      <Table
        className="defect-table"
        columns={columns}
        dataSource={defects}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `${t('common.total')} ${total} ${t('common.items')}`,
        }}
        onChange={handleTableChange}
        scroll={{ x: 1400 }}
      />

      {/* 创建/编辑模态框 */}
      <DefectFormModal
        visible={formModalVisible}
        mode={formMode}
        defect={currentDefect}
        subjects={subjects}
        phases={phases}
        onCancel={() => setFormModalVisible(false)}
        onSubmit={handleFormSubmit}
      />

      {/* 详情模态框 */}
      <DefectDetailModal
        visible={detailModalVisible}
        defect={currentDefect}
        projectId={projectId}
        onCancel={() => setDetailModalVisible(false)}
        onEdit={() => {
          setDetailModalVisible(false);
          handleEdit(currentDefect);
        }}
        onRefresh={loadDefects}
      />

      {/* 导入结果模态框 */}
      <Modal
        title={t('defect.import')}
        open={importModalVisible}
        onCancel={() => setImportModalVisible(false)}
        footer={[
          <Button key="ok" type="primary" onClick={() => setImportModalVisible(false)}>
            {t('common.ok')}
          </Button>,
        ]}
      >
        {importResult && (
          <div className="import-result">
            <div className="import-result-summary">
              {t('defect.importSuccess', {
                imported: importResult.imported || 0,
                skipped: importResult.skipped || 0,
              })}
            </div>
            {importResult.errors && importResult.errors.length > 0 && (
              <div className="import-result-errors">
                {importResult.errors.map((err, index) => (
                  <div key={index} className="import-result-error-item">
                    {err}
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </Modal>
    </div>
  );
};

export default DefectList;
