import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { message, Spin, Card, List, Tag, Button, Space, Upload, Dropdown, Menu, Popconfirm, Empty, Pagination, Input, Form, Modal } from 'antd';
import {
  PlusOutlined,
  DownloadOutlined,
  UploadOutlined,
  ExportOutlined,
  SettingOutlined,
  EyeOutlined,
  CopyOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import SubjectEditModal from './components/SubjectEditModal';
import PhaseEditModal from './components/PhaseEditModal';
import {
  fetchDefects,
  importDefects,
  exportDefects,
  downloadDefectTemplate,
} from '../../../api/defect';
import {
  DEFECT_STATUS,
  DEFECT_STATUS_COLORS,
  DEFECT_PRIORITY,
  DEFECT_PRIORITY_COLORS,
  DEFECT_SEVERITY,
  DEFECT_SEVERITY_COLORS,
} from '../../../constants/defect';

// 状态显示顺序：Resolved -> Active -> New -> Closed
const STATUS_ORDER = [
  DEFECT_STATUS.RESOLVED,
  DEFECT_STATUS.ACTIVE,
  DEFECT_STATUS.NEW,
  DEFECT_STATUS.CLOSED,
];

// 分页配置
const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];

/**
 * 缺陷一览页面 - 按状态分组显示
 */
const DefectListPage = ({
  projectId,
  projectName,
  subjects,
  phases,
  configLoading,
  onCreate,
  onDetail,
  onConfigUpdate,
}) => {
  const { t, i18n } = useTranslation();
  
  const [loading, setLoading] = useState(false);
  const [defects, setDefects] = useState([]);
  const [statusCounts, setStatusCounts] = useState({});
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  // 排序和检索
  const [sortField, setSortField] = useState('created_at'); // 'status' | 'created_at'
  const [sortOrder, setSortOrder] = useState('desc'); // 'asc' | 'desc'
  const [searchText, setSearchText] = useState('');

  // Subject/Phase 配置弹窗
  const [subjectModalVisible, setSubjectModalVisible] = useState(false);
  const [phaseModalVisible, setPhaseModalVisible] = useState(false);
  const [descriptionModalVisible, setDescriptionModalVisible] = useState(false);

  // 使用 useMemo 缓存翻译标签，只在语言变化时重新计算
  const labels = useMemo(() => ({
    create: t('defect.create', '新建缺陷'),
    import: t('common.import', '导入'),
    export: t('common.export', '导出'),
    downloadTemplate: t('common.downloadTemplate', '下载模板'),
    settings: t('common.settings', '设置'),
    subjects: t('defect.subjects', '主题管理'),
    phases: t('defect.phases', '阶段管理'),
    descriptionSettings: t('defect.descriptionSettings', '详细描述设置'),
    confirmDelete: t('common.confirmDelete', '确认删除吗？'),
    ok: t('common.ok', '确定'),
    cancel: t('common.cancel', '取消'),
    noData: t('common.noData', '暂无数据'),
    loadFailed: t('message.loadFailed', '加载失败'),
    importFailed: t('message.importFailed', '导入失败'),
    exportSuccess: t('message.exportSuccess', '导出成功'),
    exportFailed: t('message.exportFailed', '导出失败'),
    downloadFailed: t('message.downloadFailed', '下载失败'),
    detail: t('defect.detail', '详细'),
    defectId: t('defect.defectId', '缺陷ID'),
    title: t('defect.title', '标题'),
    createTime: t('common.createdAt', '创建时间'),
    operation: t('common.actions', '操作'),
  }), [t, i18n.language]);

  // 缓存状态标签
  const statusLabels = useMemo(() => ({
    [DEFECT_STATUS.NEW]: t('defect.statusNew', '新建'),
    [DEFECT_STATUS.ACTIVE]: t('defect.statusActive', '处理中'),
    [DEFECT_STATUS.RESOLVED]: t('defect.statusResolved', '已解决'),
    [DEFECT_STATUS.CLOSED]: t('defect.statusClosed', '已关闭'),
  }), [t, i18n.language]);

  // 缓存优先级标签（A/B/C/D）
  const priorityLabels = useMemo(() => ({
    [DEFECT_PRIORITY.A]: t('defect.priorityA', 'A'),
    [DEFECT_PRIORITY.B]: t('defect.priorityB', 'B'),
    [DEFECT_PRIORITY.C]: t('defect.priorityC', 'C'),
    [DEFECT_PRIORITY.D]: t('defect.priorityD', 'D'),
  }), [t, i18n.language]);

  // 缓存严重程度标签（A/B/C/D）
  const severityLabels = useMemo(() => ({
    [DEFECT_SEVERITY.A]: t('defect.severityA', 'A'),
    [DEFECT_SEVERITY.B]: t('defect.severityB', 'B'),
    [DEFECT_SEVERITY.C]: t('defect.severityC', 'C'),
    [DEFECT_SEVERITY.D]: t('defect.severityD', 'D'),
  }), [t, i18n.language]);

  // 加载缺陷列表（支持分页、排序、检索）
  const loadDefects = useCallback(async () => {
    if (!projectId) return;
    setLoading(true);
    try {
      const response = await fetchDefects(projectId, { 
        page: currentPage, 
        size: pageSize,
        keyword: searchText,
      });
      let defectList = response.defects || [];
      
      // 前端排序
      if (sortField === 'status') {
        const statusOrder = {
          [DEFECT_STATUS.RESOLVED]: 1,
          [DEFECT_STATUS.ACTIVE]: 2,
          [DEFECT_STATUS.NEW]: 3,
          [DEFECT_STATUS.CLOSED]: 4,
        };
        defectList = defectList.sort((a, b) => {
          const orderA = statusOrder[a.status] || 999;
          const orderB = statusOrder[b.status] || 999;
          return sortOrder === 'asc' ? orderA - orderB : orderB - orderA;
        });
      } else if (sortField === 'created_at') {
        defectList = defectList.sort((a, b) => {
          const timeA = new Date(a.created_at).getTime();
          const timeB = new Date(b.created_at).getTime();
          return sortOrder === 'asc' ? timeA - timeB : timeB - timeA;
        });
      } else if (sortField === 'severity') {
        const severityOrder = {
          [DEFECT_SEVERITY.A]: 1,
          [DEFECT_SEVERITY.B]: 2,
          [DEFECT_SEVERITY.C]: 3,
          [DEFECT_SEVERITY.D]: 4,
        };
        defectList = defectList.sort((a, b) => {
          const orderA = severityOrder[a.severity] || 999;
          const orderB = severityOrder[b.severity] || 999;
          return sortOrder === 'asc' ? orderA - orderB : orderB - orderA;
        });
      } else if (sortField === 'subject') {
        defectList = defectList.sort((a, b) => {
          const subjectA = (a.subject || '').toString();
          const subjectB = (b.subject || '').toString();
          return sortOrder === 'asc' ? subjectA.localeCompare(subjectB) : subjectB.localeCompare(subjectA);
        });
      }
      
      setDefects(defectList);
      setStatusCounts(response.status_counts || {});
      setTotal(response.total || 0);
    } catch (error) {
      console.error('Failed to load defects:', error);
      message.error(labels.loadFailed);
    } finally {
      setLoading(false);
    }
  }, [projectId, currentPage, pageSize, searchText, sortField, sortOrder, labels.loadFailed]);

  // 首次加载
  useEffect(() => {
    loadDefects();
  }, [loadDefects]);

  // 排序切换
  const handleSort = (field) => {
    if (sortField === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortOrder(field === 'created_at' ? 'desc' : 'asc');
    }
  };

  // 检索处理
  const handleSearch = (value) => {
    setSearchText(value);
    setCurrentPage(1);
  };

  // 分页变化处理
  const handlePageChange = (page, size) => {
    setCurrentPage(page);
    if (size !== pageSize) {
      setPageSize(size);
      setCurrentPage(1); // 改变每页大小时回到第一页
    }
  };

  // 导入缺陷
  const handleImport = async (file) => {
    try {
      console.log('[handleImport] 接收到的file对象:', file);
      console.log('[handleImport] file是否为null:', file === null);
      console.log('[handleImport] file类型:', typeof file);
      console.log('[handleImport] file.name:', file?.name);
      console.log('[handleImport] file.size:', file?.size);
      
      const result = await importDefects(projectId, file);
      console.log('[handleImport] 导入结果:', result);
      
      if (result.fail_count > 0) {
        message.warning(t('defect.importPartialSuccess', { 
          success: result.success_count, 
          fail: result.fail_count,
          defaultValue: `导入完成：成功 ${result.success_count} 条，失败 ${result.fail_count} 条`
        }));
      } else {
        message.success(t('defect.importSuccess', { 
          count: result.success_count,
          defaultValue: `成功导入 ${result.success_count} 条缺陷`
        }));
      }
      loadDefects();
    } catch (error) {
      console.error('Failed to import defects:', error);
      message.error(labels.importFailed);
    }
  };

  // 导出缺陷
  const handleExport = async (format) => {
    try {
      await exportDefects(projectId, format, {}, projectName);
      message.success(labels.exportSuccess);
    } catch (error) {
      console.error('Failed to export defects:', error);
      message.error(labels.exportFailed);
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
      message.error(labels.downloadFailed);
    }
  };

  // 确认下载模板
  const handleConfirmDownloadTemplate = async (format) => {
    await handleDownloadTemplate(format);
  };

  // Subject配置更新
  const handleSubjectUpdate = () => {
    setSubjectModalVisible(false);
    onConfigUpdate?.();
  };

  // Phase配置更新
  const handlePhaseUpdate = () => {
    setPhaseModalVisible(false);
    onConfigUpdate?.();
  };

  // 配置菜单
  const configMenu = useMemo(() => (
    <Menu>
      <Menu.Item key="subjects" onClick={() => setSubjectModalVisible(true)}>
        {labels.subjects}
      </Menu.Item>
      <Menu.Item key="phases" onClick={() => setPhaseModalVisible(true)}>
        {labels.phases}
      </Menu.Item>
      <Menu.Item key="description" onClick={() => setDescriptionModalVisible(true)}>
        {labels.descriptionSettings}
      </Menu.Item>
    </Menu>
  ), [labels.subjects, labels.phases, labels.descriptionSettings]);

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

  // 下载菜单
  const downloadMenu = useMemo(() => (
    <Menu>
      <Menu.Item key="csv" onClick={() => handleConfirmDownloadTemplate('csv')}>
        CSV
      </Menu.Item>
      <Menu.Item key="xlsx" onClick={() => handleConfirmDownloadTemplate('xlsx')}>
        XLSX
      </Menu.Item>
    </Menu>
  ), []);

  // 获取状态标签
  const getStatusTag = (status) => {
    return (
      <Tag color={DEFECT_STATUS_COLORS[status]}>
        {statusLabels[status] || status}
      </Tag>
    );
  };

  // 渲染缺陷项
  const renderDefectItem = (defect) => (
    <div 
      key={defect.id}
      style={{ 
        display: 'flex', 
        alignItems: 'center', 
        padding: '12px 16px', 
        borderBottom: '1px solid #f0f0f0',
        cursor: 'pointer'
      }}
      onClick={() => onDetail(defect)}
    >
      {/* 状态列 */}
      <div style={{ width: '100px', flexShrink: 0 }}>
        {getStatusTag(defect.status)}
      </div>
      
      {/* 创建时间 */}
      <div style={{ width: '160px', flexShrink: 0, color: '#666', fontSize: 14 }}>
        {defect.created_at ? dayjs(defect.created_at).format('YYYY-MM-DD') : '-'}
      </div>
      
      {/* 严重程度 */}
      <div style={{ width: '100px', flexShrink: 0 }}>
        {<Tag color={DEFECT_SEVERITY_COLORS[defect.severity]}>
          {severityLabels[defect.severity] || defect.severity}
        </Tag>}
      </div>
      
      {/* Defect ID */}
      <div style={{ width: '140px', flexShrink: 0, display: 'flex', alignItems: 'center', gap: '8px' }}>
        <span style={{ color: '#1890ff', fontWeight: 500 }}>{defect.defect_id}</span>
        <Button
          type="text"
          size="small"
          icon={<CopyOutlined />}
          onClick={(e) => {
            e.stopPropagation();
            navigator.clipboard.writeText(defect.defect_id);
            message.success(t('common.copySuccess', '复制成功'));
          }}
          style={{ padding: '0 4px', minWidth: 'auto' }}
          title={t('common.copy', '复制')}
        />
      </div>
      
      {/* 模块 */}
      <div style={{ width: '84px', flexShrink: 0, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }} title={defect.subject || '-'}>
        <span style={{ color: '#666', fontSize: 14 }}>{defect.subject || '-'}</span>
      </div>
      
      {/* Title */}
      <div style={{ flex: 1, marginRight: 16 }}>
        <span>{defect.title}</span>
      </div>
      
      {/* 提出人列 */}
      <div style={{ width: '100px', flexShrink: 0, color: '#666', fontSize: 14 }}>
        {defect.created_by_user?.username || '-'}
      </div>
      
      {/* Operation */}
      <div style={{ width: '100px', flexShrink: 0, textAlign: 'center' }}>
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            title={labels.detail}
            onClick={(e) => {
              e.stopPropagation();
              onDetail(defect);
            }}
          />
        </Space>
      </div>
    </div>
  );

  return (
    <Spin spinning={loading || configLoading}>
      <div className="defect-list-page">
        {/* 工具栏 */}
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Space>
            <Button type="primary" icon={<PlusOutlined />} onClick={onCreate}>
              {labels.create}
            </Button>

            <Upload
              accept=".csv,.xlsx"
              showUploadList={false}
              beforeUpload={(file) => {
                handleImport(file);
                return false;
              }}
            >
              <Button icon={<UploadOutlined />}>{labels.import}</Button>
            </Upload>

            <Dropdown overlay={exportMenu} placement="bottomRight">
              <Button icon={<ExportOutlined />}>{labels.export}</Button>
            </Dropdown>

            <Dropdown overlay={downloadMenu} placement="bottomRight">
              <Button icon={<DownloadOutlined />}>{labels.downloadTemplate}</Button>
            </Dropdown>

            <Dropdown overlay={configMenu} placement="bottomRight">
              <Button icon={<SettingOutlined />}>{labels.settings}</Button>
            </Dropdown>
          </Space>

          <Input.Search
            placeholder={t('defect.searchPlaceholder')}
            allowClear
            style={{ width: 300 }}
            onSearch={handleSearch}
            onChange={(e) => !e.target.value && handleSearch('')}
          />
        </div>

        {/* 统计信息 */}
        <div style={{ marginBottom: 16, padding: '12px 16px', background: '#fafafa', borderRadius: 4 }}>
          <Space size="large">
            <span>
              <strong>{t('defect.totalCount')}:</strong> <span style={{ color: '#1890ff', fontWeight: 'bold' }}>{total}</span>
            </span>
            <span>
              <strong>{t('defect.resolvedCount')}:</strong> <span style={{ color: '#52c41a', fontWeight: 'bold' }}>{statusCounts[DEFECT_STATUS.RESOLVED] || 0}</span>
            </span>
            <span>
              <strong>{t('defect.activeCount')}:</strong> <span style={{ color: '#faad14', fontWeight: 'bold' }}>{statusCounts[DEFECT_STATUS.ACTIVE] || 0}</span>
            </span>
            <span>
              <strong>{t('defect.newCount')}:</strong> <span style={{ color: '#ff4d4f', fontWeight: 'bold' }}>{statusCounts[DEFECT_STATUS.NEW] || 0}</span>
            </span>
            <span>
              <strong>{t('defect.closedCount')}:</strong> <span style={{ color: '#8c8c8c', fontWeight: 'bold' }}>{statusCounts[DEFECT_STATUS.CLOSED] || 0}</span>
            </span>
          </Space>
        </div>

        {/* 缺陷列表 */}
        <Card bodyStyle={{ padding: 0 }}>
          {/* 表头 */}
          <div style={{ 
            display: 'flex', 
            padding: '12px 16px', 
            background: '#fafafa', 
            borderBottom: '2px solid #e8e8e8',
            fontWeight: 600
          }}>
            <div 
              style={{ width: '100px', flexShrink: 0, cursor: 'pointer', userSelect: 'none' }}
              onClick={() => handleSort('status')}
            >
              {t('defect.tableHeaderStatus')} {sortField === 'status' && (sortOrder === 'asc' ? '↑' : '↓')}
            </div>
            <div 
              style={{ width: '160px', flexShrink: 0, cursor: 'pointer', userSelect: 'none' }}
              onClick={() => handleSort('created_at')}
            >
              {t('defect.tableHeaderCreatedAt')} {sortField === 'created_at' && (sortOrder === 'asc' ? '↑' : '↓')}
            </div>
            <div 
              style={{ width: '100px', flexShrink: 0, cursor: 'pointer', userSelect: 'none' }}
              onClick={() => handleSort('severity')}
            >
              {t('defect.severity')} {sortField === 'severity' && (sortOrder === 'asc' ? '↑' : '↓')}
            </div>
            <div style={{ width: '140px', flexShrink: 0 }}>{t('defect.tableHeaderDefectId')}</div>
            <div 
              style={{ width: '84px', flexShrink: 0, cursor: 'pointer', userSelect: 'none' }}
              onClick={() => handleSort('subject')}
            >
              {t('defect.subject')} {sortField === 'subject' && (sortOrder === 'asc' ? '↑' : '↓')}
            </div>
            <div style={{ flex: 1, marginRight: 16 }}>{t('defect.tableHeaderTitle')}</div>
            <div style={{ width: '100px', flexShrink: 0 }}>{t('defect.reporter')}</div>
            <div style={{ width: '100px', flexShrink: 0, textAlign: 'center' }}>{t('defect.tableHeaderActions')}</div>
          </div>
          
          {/* 数据行 */}
          {defects.length > 0 ? (
            defects.map(renderDefectItem)
          ) : (
            <div style={{ padding: '48px 0', textAlign: 'center' }}>
              <Empty description={labels.noData} />
            </div>
          )}
        </Card>

        {/* 分页器 */}
        {total > 0 && (
          <div style={{ marginTop: 16, display: 'flex', justifyContent: 'flex-end' }}>
            <Pagination
              current={currentPage}
              pageSize={pageSize}
              total={total}
              onChange={handlePageChange}
              onShowSizeChange={handlePageChange}
              showSizeChanger
              showQuickJumper
              pageSizeOptions={PAGE_SIZE_OPTIONS}
              showTotal={(total) => `${t('common.total', '共')} ${total} ${t('common.items', '条')}`}
            />
          </div>
        )}

        {/* Subject配置弹窗 */}
        <SubjectEditModal
          visible={subjectModalVisible}
          projectId={projectId}
          subjects={subjects}
          onClose={() => setSubjectModalVisible(false)}
          onUpdate={handleSubjectUpdate}
        />

        {/* Phase配置弹窗 */}
        <PhaseEditModal
          visible={phaseModalVisible}
          projectId={projectId}
          phases={phases}
          onClose={() => setPhaseModalVisible(false)}
          onUpdate={handlePhaseUpdate}
        />

        {/* 详细描述设置弹窗 */}
        <DescriptionSettingsModal
          visible={descriptionModalVisible}
          onClose={() => setDescriptionModalVisible(false)}
        />
      </div>
    </Spin>
  );
};


// 详细描述设置对话框组件
const DescriptionSettingsModal = ({ visible, onClose }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  // 默认详细描述模板
  const DEFAULT_DESCRIPTION_TEMPLATE = `[Actual result]

[Relevant validation]

[Test Steps]

→issue occurred

[Expected result]

[Test Environment]

`;

  // 组件挂载时加载当前模板
  useEffect(() => {
    if (visible) {
      const savedTemplate = localStorage.getItem('defect_description_template') || DEFAULT_DESCRIPTION_TEMPLATE;
      form.setFieldsValue({ template: savedTemplate });
    }
  }, [visible, form]);

  // 保存模板
  const handleSave = async () => {
    try {
      setLoading(true);
      const values = await form.validateFields();
      localStorage.setItem('defect_description_template', values.template);
      message.success(t('common.saveSuccess', '保存成功'));
      onClose();
    } catch (error) {
      console.error('Save description template failed:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title={t('defect.descriptionSettings', '详细描述设置')}
      open={visible}
      onCancel={onClose}
      onOk={handleSave}
      confirmLoading={loading}
      width={600}
      okText={t('common.save', '保存')}
      cancelText={t('common.cancel', '取消')}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="template"
          label={t('defect.descriptionTemplate', '详细描述模板')}
          rules={[{ required: true, message: t('validation.required', '此字段为必填项') }]}
        >
          <Input.TextArea
            rows={12}
            placeholder={t('defect.descriptionTemplatePlaceholder', '请输入详细描述模板')}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default DefectListPage;
