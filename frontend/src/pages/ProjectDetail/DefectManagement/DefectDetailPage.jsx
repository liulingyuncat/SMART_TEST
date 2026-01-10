import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import dayjs from 'dayjs';
import {
  Form,
  Input,
  Select,
  Button,
  Card,
  Row,
  Col,
  message,
  Space,
  Upload,
  List,
  Popconfirm,
  Spin,
  Tag,
  Typography,
  Divider,
  Image,
  Modal,
} from 'antd';
import {
  ArrowLeftOutlined,
  SaveOutlined,
  UploadOutlined,
  DeleteOutlined,
  DownloadOutlined,
  EditOutlined,
  FileOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import {
  DEFECT_STATUS,
  DEFECT_STATUS_COLORS,
  DEFECT_PRIORITY,
  DEFECT_PRIORITY_COLORS,
  DEFECT_SEVERITY,
  DEFECT_SEVERITY_COLORS,
} from '../../../constants/defect';
import {
  fetchDefect,
  updateDefect,
  deleteDefect,
  fetchDefectAttachments,
  uploadDefectAttachment,
  deleteDefectAttachment,
  downloadDefectAttachment,
} from '../../../api/defect';
import CommentSection from './components/CommentSection';

const { TextArea } = Input;
const { Option } = Select;
const { Text, Title } = Typography;

/**
 * 缺陷详情/编辑页面
 * 支持查看模式和编辑模式切换
 */
const DefectDetailPage = ({
  projectId,
  defectId,
  subjects,
  phases,
  onBack,
  onSuccess,
  onCreate,
}) => {
  const { t, i18n } = useTranslation();
  const currentUser = useSelector((state) => state.auth.user);
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [defect, setDefect] = useState(null);
  const [attachments, setAttachments] = useState([]);
  const [fileList, setFileList] = useState([]);
  const [editing, setEditing] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [previewVisible, setPreviewVisible] = useState(false);
  const [previewImage, setPreviewImage] = useState('');

  // 判断是否为图片类型
  const isImageFile = (fileName) => {
    const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg'];
    const ext = fileName.toLowerCase().substring(fileName.lastIndexOf('.'));
    return imageExtensions.includes(ext);
  };

  // 预览图片
  const handlePreviewImage = async (attachment) => {
    try {
      const token = localStorage.getItem('auth_token');
      const response = await fetch(
        `/api/v1/projects/${projectId}/defects/${defectId}/attachments/${attachment.id}`,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        }
      );
      if (!response.ok) {
        throw new Error('Failed to load image');
      }
      const blob = await response.blob();
      const imageUrl = URL.createObjectURL(blob);
      setPreviewImage(imageUrl);
      setPreviewVisible(true);
    } catch (error) {
      console.error('Failed to preview image:', error);
      message.error('图片预览失败');
    }
  };

  // 使用 useMemo 缓存翻译标签，只在语言变化时重新计算
  const labels = useMemo(() => ({
    detail: t('defect.detail', '缺陷详情'),
    edit: t('defect.edit', '编辑缺陷'),
    save: t('common.save', '保存'),
    cancel: t('common.cancel', '取消'),
    upload: t('common.upload', '上传'),
    download: t('common.download', '下载'),
    delete: t('common.delete', '删除'),
    confirmDelete: t('common.confirmDelete', '确认删除吗？'),
    title: t('defect.title', 'Title'),
    description: t('defect.description', 'Description'),
    recoveryMethod: t('defect.recoveryMethod', 'Recovery Method'),
    frequency: t('defect.frequency', 'Frequency(%)'),
    detectedInRelease: t('defect.detectedInRelease', 'Detected in Release'),
    status: t('defect.status', 'Status'),
    priority: t('defect.priority', 'Priority'),
    severity: t('defect.severity', 'Severity'),
    subject: t('defect.subject', 'Subject'),
    phase: t('defect.phase', 'Phase'),
    attachments: t('defect.attachments', 'Attachments'),
    noAttachments: t('defect.noAttachments', '暂无附件'),
    createdAt: t('defect.createdAt', '创建时间'),
    updatedAt: t('defect.updatedAt', '更新时间'),
    required: t('validation.required', '此字段为必填项'),
    loadFailed: t('message.loadFailed', '加载失败'),
    saveSuccess: t('message.saveSuccess', '保存成功'),
    saveFailed: t('message.saveFailed', '保存失败'),
    deleteSuccess: t('message.deleteSuccess', '删除成功'),
    deleteFailed: t('message.deleteFailed', '删除失败'),
    downloadFailed: t('message.downloadFailed', '下载失败'),
    maxFileSize: t('defect.maxFileSize', '单个文件最大100MB'),
    ok: t('common.ok', '确定'),
  }), [t, i18n.language]);

  // 状态标签映射
  const statusLabels = useMemo(() => ({
    [DEFECT_STATUS.NEW]: t('defect.statusNew', '新建'),
    [DEFECT_STATUS.ACTIVE]: t('defect.statusActive', '处理中'),
    [DEFECT_STATUS.RESOLVED]: t('defect.statusResolved', '已解决'),
    [DEFECT_STATUS.CLOSED]: t('defect.statusClosed', '已关闭'),
  }), [t, i18n.language]);

  // 优先级标签映射（A/B/C/D）
  const priorityLabels = useMemo(() => ({
    [DEFECT_PRIORITY.A]: t('defect.priorityA', 'A'),
    [DEFECT_PRIORITY.B]: t('defect.priorityB', 'B'),
    [DEFECT_PRIORITY.C]: t('defect.priorityC', 'C'),
    [DEFECT_PRIORITY.D]: t('defect.priorityD', 'D'),
  }), [t, i18n.language]);

  // 严重程度标签映射（A/B/C/D）
  const severityLabels = useMemo(() => ({
    [DEFECT_SEVERITY.A]: t('defect.severityA', 'A'),
    [DEFECT_SEVERITY.B]: t('defect.severityB', 'B'),
    [DEFECT_SEVERITY.C]: t('defect.severityC', 'C'),
    [DEFECT_SEVERITY.D]: t('defect.severityD', 'D'),
  }), [t, i18n.language]);

  // 状态选项
  const statusOptions = useMemo(() => Object.values(DEFECT_STATUS).map((value) => ({
    value,
    label: statusLabels[value],
    color: DEFECT_STATUS_COLORS[value],
  })), [statusLabels]);

  // 优先级选项
  const priorityOptions = useMemo(() => Object.values(DEFECT_PRIORITY).map((value) => ({
    value,
    label: priorityLabels[value],
    color: DEFECT_PRIORITY_COLORS[value],
  })), [priorityLabels]);

  // 严重程度选项
  const severityOptions = useMemo(() => Object.values(DEFECT_SEVERITY).map((value) => ({
    value,
    label: severityLabels[value],
    color: DEFECT_SEVERITY_COLORS[value],
  })), [severityLabels]);

  // 加载缺陷详情
  const loadDefect = useCallback(async () => {
    if (!projectId || !defectId) {
      console.log('[DEBUG] loadDefect: missing projectId or defectId', { projectId, defectId });
      return;
    }
    setLoading(true);
    try {
      console.log('[DEBUG] loadDefect: fetching defect', { projectId, defectId });
      const response = await fetchDefect(projectId, defectId);
      console.log('[DEBUG] loadDefect: API response', response);
      // apiClient已经提取了data字段，response直接就是defect对象
      // 但可能有 { defect: {...} } 或直接是 defect对象
      const defectData = response?.defect || response;
      console.log('[DEBUG] loadDefect: defectData to set', defectData);
      if (defectData) {
        setDefect(defectData);
        form.setFieldsValue(defectData);
      } else {
        console.error('[DEBUG] loadDefect: No defect data in response');
        message.error(labels.loadFailed);
      }
    } catch (error) {
      console.error('[DEBUG] loadDefect: Failed to load defect:', error);
      message.error(labels.loadFailed);
    } finally {
      setLoading(false);
    }
  }, [projectId, defectId, form, labels.loadFailed]);

  // 加载附件列表
  const loadAttachments = useCallback(async () => {
    if (!projectId || !defectId) return;
    try {
      const response = await fetchDefectAttachments(projectId, defectId);
      setAttachments(response.attachments || []);
    } catch (error) {
      console.error('Failed to load attachments:', error);
    }
  }, [projectId, defectId]);

  useEffect(() => {
    loadDefect();
    loadAttachments();
  }, [loadDefect, loadAttachments]);

  // 提交更新
  const handleSubmit = async (values) => {
    setSaving(true);
    try {
      await updateDefect(projectId, defectId, values);

      message.success(labels.saveSuccess);
      setEditing(false);
      await loadDefect();
      onSuccess?.();
    } catch (error) {
      console.error('[DEBUG] Failed to update defect:', error);
      message.error(labels.saveFailed);
    } finally {
      setSaving(false);
    }
  };

  // 删除附件
  const handleDeleteAttachment = async (attachmentId) => {
    try {
      await deleteDefectAttachment(projectId, defectId, attachmentId);
      message.success(labels.deleteSuccess);
      loadAttachments();
    } catch (error) {
      console.error('Failed to delete attachment:', error);
      message.error(labels.deleteFailed);
    }
  };

  // 下载附件
  const handleDownloadAttachment = async (attachment) => {
    try {
      await downloadDefectAttachment(projectId, defectId, attachment.id, attachment.file_name);
    } catch (error) {
      console.error('Failed to download attachment:', error);
      message.error(labels.downloadFailed);
    }
  };

  // 文件上传前处理
  const handleBeforeUpload = (file) => {
    setFileList((prev) => [...prev, { uid: file.uid, name: file.name, originFileObj: file }]);
    return false;
  };

  // 移除待上传文件
  const handleRemoveFile = (file) => {
    setFileList((prev) => prev.filter((f) => f.uid !== file.uid));
  };

  // 保存附件（独立操作）
  const handleSaveAttachments = async () => {
    if (fileList.length === 0) return;
    
    setSaving(true);
    try {
      console.log('[DEBUG] Saving attachments:', fileList.length);
      for (const file of fileList) {
        await uploadDefectAttachment(projectId, defectId, file.originFileObj);
      }
      setFileList([]);
      message.success('附件上传成功');
      
      // 刷新附件列表
      await loadAttachments();
      console.log('[DEBUG] Attachments saved and reloaded');
    } catch (error) {
      console.error('[DEBUG] Failed to save attachments:', error);
      message.error('附件上传失败');
    } finally {
      setSaving(false);
    }
  };

  // 删除缺陷
  const handleDeleteDefect = async () => {
    setDeleting(true);
    try {
      await deleteDefect(projectId, defectId);
      message.success(labels.deleteSuccess);
      onBack?.();
    } catch (error) {
      console.error('Failed to delete defect:', error);
      message.error(labels.deleteFailed);
    } finally {
      setDeleting(false);
    }
  };

  // 获取状态标签
  const getStatusTag = (status) => {
    const color = DEFECT_STATUS_COLORS[status];
    const label = statusLabels[status];
    if (!label) return status;
    return <Tag color={color}>{label}</Tag>;
  };

  // 获取优先级标签
  const getPriorityTag = (priority) => {
    const color = DEFECT_PRIORITY_COLORS[priority];
    const label = priorityLabels[priority];
    if (!label) return priority;
    return <Tag color={color}>{label}</Tag>;
  };

  // 获取严重程度标签
  const getSeverityTag = (severity) => {
    const color = DEFECT_SEVERITY_COLORS[severity];
    const label = severityLabels[severity];
    if (!label) return severity;
    return <Tag color={color}>{label}</Tag>;
  };

  // 查看模式下的详情展示 - 紧凑版
  const renderViewMode = () => (
    <div className="defect-view-mode" style={{ fontSize: 13 }}>
      {/* 状态信息行 - 一行显示 */}
      <div style={{ background: '#fafafa', borderRadius: 4, padding: '8px 12px', marginBottom: 12 }}>
        <Row gutter={24}>
          <Col span={6}>
            <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{labels.status}:</span>
            {getStatusTag(defect?.status)}
          </Col>
          <Col span={6}>
            <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{labels.priority}:</span>
            {getPriorityTag(defect?.priority)}
          </Col>
          <Col span={6}>
            <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{labels.severity}:</span>
            {getSeverityTag(defect?.severity)}
          </Col>
          <Col span={6}>
            <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{labels.recoveryMethod}:</span>
            <span style={{ color: '#303133' }}>{defect?.recovery_method || '-'}</span>
          </Col>
        </Row>
      </div>

      {/* 基本信息行 */}
      <Row gutter={24} style={{ marginBottom: 8 }}>
        <Col span={6}>
          <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{labels.frequency}:</span>
          <span style={{ color: '#303133' }}>{defect?.frequency || '-'}</span>
        </Col>
        <Col span={6}>
          <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{labels.detectedInRelease}:</span>
          <span style={{ color: '#303133' }}>{defect?.detected_in_release || '-'}</span>
        </Col>
        <Col span={6}>
          <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{t('defect.caseId', 'Case ID')}:</span>
          <span style={{ color: '#303133' }}>{defect?.case_id || '-'}</span>
        </Col>
        <Col span={6}>
          <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{labels.subject}:</span>
          <span style={{ color: '#303133' }}>{defect?.subject || '-'}</span>
        </Col>
      </Row>

      <Row gutter={24} style={{ marginBottom: 12 }}>
        <Col span={6}>
          <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{labels.phase}:</span>
          <span style={{ color: '#303133' }}>{defect?.phase || '-'}</span>
        </Col>
        <Col span={6}>
          <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{t('defect.detectedBy', '提出人')}:</span>
          <span style={{ color: '#303133' }}>{defect?.created_by_user?.username || '-'}</span>
        </Col>
        <Col span={6}>
          <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{labels.createdAt}:</span>
          <span style={{ color: '#303133' }}>{defect?.created_at ? dayjs(defect.created_at).format('YYYY-MM-DD') : '-'}</span>
        </Col>
        <Col span={6}>
          <span style={{ color: '#8c8c8c', fontSize: 12, marginRight: 8 }}>{labels.updatedAt}:</span>
          <span style={{ color: '#303133' }}>{defect?.updated_at ? dayjs(defect.updated_at).format('YYYY-MM-DD') : '-'}</span>
        </Col>
      </Row>

      {/* 详细描述 - 自适应高度，预设文字高亮 */}
      <div style={{ marginBottom: 8 }}>
        <span style={{ color: '#8c8c8c', fontSize: 12, display: 'block', marginBottom: 4 }}>{labels.description}:</span>
        <div style={{ 
          background: '#fafafa', 
          border: '1px solid #f0f0f0', 
          borderRadius: 4, 
          padding: '10px 12px',
          lineHeight: 1.6,
          fontSize: 14,
        }}>
          {defect?.description ? (
            defect.description.split('\n').map((line, index) => {
              // 检测是否为预设模板行（以 [ 开头的行）
              const isTemplateLine = /^\[.+\]/.test(line.trim());
              return (
                <div key={index} style={{ 
                  color: isTemplateLine ? '#1890ff' : '#303133',
                  fontWeight: isTemplateLine ? 500 : 400,
                  minHeight: line.trim() === '' ? '0.8em' : 'auto'
                }}>
                  {line || '\u00A0'}
                </div>
              );
            })
          ) : '-'}
        </div>
      </div>
    </div>
  );

  // 编辑模式下的表单
  const renderEditMode = () => (
    <Form
      form={form}
      layout="vertical"
      onFinish={handleSubmit}
    >
      <Row gutter={24}>
        <Col span={16}>
          <Form.Item
            name="title"
            label={labels.title}
            rules={[{ required: true, message: labels.required }]}
          >
            <Input maxLength={200} />
          </Form.Item>

          <Form.Item name="subject_id" label={labels.subject}>
            <Select allowClear>
              {subjects?.map((subject) => (
                <Option key={subject.id} value={subject.id}>
                  {subject.name}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item name="description" label={labels.description}>
            <TextArea rows={12} />
          </Form.Item>
        </Col>

        <Col span={8}>
          <Form.Item 
            name="status" 
            label={labels.status}
          >
            <Select>
              {statusOptions.map((opt) => (
                <Option key={opt.value} value={opt.value}>
                  <span style={{ color: opt.color }}>{opt.label}</span>
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item name="recovery_method" label={labels.recoveryMethod}>
            <Input maxLength={500} />
          </Form.Item>

          <Form.Item 
            name="priority" 
            label={labels.priority}
            rules={[{ required: true, message: labels.required }]}
          >
            <Select>
              {priorityOptions.map((opt) => (
                <Option key={opt.value} value={opt.value}>
                  <span style={{ color: opt.color }}>{opt.label}</span>
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item name="severity" label={labels.severity}
            rules={[{ required: true, message: labels.required }]}
          >
            <Select>
              {severityOptions.map((opt) => (
                <Option key={opt.value} value={opt.value}>
                  <span style={{ color: opt.color }}>{opt.label}</span>
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item 
            name="frequency" 
            label={labels.frequency}
            rules={[{ required: true, message: labels.required }]}
          >
            <Input maxLength={10} />
          </Form.Item>

          <Form.Item 
            name="detected_in_release" 
            label={labels.detectedInRelease}
            rules={[{ required: true, message: labels.required }]}
          >
            <Input maxLength={50} />
          </Form.Item>

          <Form.Item name="case_id" label={t('defect.caseId', 'Case ID')}>
            <Input maxLength={50} />
          </Form.Item>

          <Form.Item name="phase_id" label={labels.phase}>
            <Select allowClear>
              {phases?.map((phase) => (
                <Option key={phase.id} value={phase.id}>
                  {phase.name}
                </Option>
              ))}
            </Select>
          </Form.Item>
        </Col>
      </Row>
    </Form>
  );

  return (
    <Spin spinning={loading}>
      <div className="defect-detail-page">
        <Card
          style={{
            borderRadius: '8px',
            boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
            border: '1px solid #f0f0f0'
          }}
          title={
            <Space align="start">
              <Button icon={<ArrowLeftOutlined />} onClick={onBack} style={{ marginTop: 4 }} />
              <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                <div style={{ color: '#8c8c8c', fontSize: 16, fontWeight: 500, flexShrink: 0 }}>{defect?.defect_id || ''}</div>
                {defect?.title && (
                  <div 
                    style={{ 
                      fontSize: 20, 
                      fontWeight: 600, 
                      color: '#262626',
                      whiteSpace: 'nowrap',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      maxWidth: '600px',
                      cursor: 'default'
                    }}
                    title={defect.title}
                  >
                    {defect.title}
                  </div>
                )}
              </div>
            </Space>
          }
          extra={
            editing ? (
              <Space>
                <Button onClick={() => setEditing(false)}>
                  {labels.cancel}
                </Button>
                <Button
                  type="primary"
                  icon={<SaveOutlined />}
                  loading={saving}
                  onClick={() => form.submit()}
                >
                  {labels.save}
                </Button>
              </Space>
            ) : (
              <Space>
                <Button
                  type="primary"
                  icon={<EditOutlined />}
                  onClick={() => setEditing(true)}
                >
                  {labels.edit}
                </Button>
                <Button
                  onClick={() => {
                    if (onCreate) {
                      onCreate();
                    } else if (onBack) {
                      onBack();
                    }
                  }}
                >
                  {t('defect.createDefect', '新建缺陷')}
                </Button>
                <Popconfirm
                  title={t('defect.confirmDeleteDefect', { defectId: defect?.defect_id || '' })}
                  onConfirm={handleDeleteDefect}
                  okText="确定"
                  cancelText="取消"
                  okButtonProps={{ loading: deleting }}
                >
                  <Button danger icon={<DeleteOutlined />}>
                    {labels.delete}
                  </Button>
                </Popconfirm>
              </Space>
            )
          }
        >
          {editing ? renderEditMode() : renderViewMode()}

          <Divider style={{ margin: '12px 0' }} />

          {/* 附件区域 - 紧凑版 */}
          <div className="attachments-section">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
              <span style={{ fontSize: 13, fontWeight: 500, color: '#303133' }}>{labels.attachments}</span>
              <Space size="small">
                <Upload
                  fileList={fileList}
                  beforeUpload={handleBeforeUpload}
                  onRemove={handleRemoveFile}
                  multiple
                  showUploadList={false}
                >
                  <Button size="small" icon={<UploadOutlined />}>
                    {labels.upload}
                  </Button>
                </Upload>
                {fileList.length > 0 && (
                  <Button 
                    size="small"
                    type="primary" 
                    icon={<SaveOutlined />}
                    loading={saving}
                    onClick={handleSaveAttachments}
                  >
                    {labels.save}
                  </Button>
                )}
              </Space>
            </div>

            {/* 待上传附件列表 */}
            {fileList.length > 0 && (
              <div style={{ marginBottom: 8, padding: '6px 10px', background: '#f0f7ff', borderRadius: 4, border: '1px solid #91d5ff' }}>
                <div style={{ marginBottom: 4, fontSize: 12, fontWeight: 500, color: '#1890ff' }}>待上传附件：</div>
                {fileList.map((file) => (
                  <div key={file.uid} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '2px 0', fontSize: 12 }}>
                    <span style={{ color: '#303133' }}>{file.name}</span>
                    <Space size="small">
                      <span style={{ color: '#8c8c8c' }}>{(file.size / 1024).toFixed(1)} KB</span>
                      <Button type="link" size="small" danger style={{ padding: 0, height: 'auto', fontSize: 12 }} onClick={() => handleRemoveFile(file)}>移除</Button>
                    </Space>
                  </div>
                ))}
              </div>
            )}

            {/* 已上传附件列表 */}
            {attachments.length === 0 ? (
              <div style={{ color: '#bfbfbf', fontSize: 12, padding: '8px 0' }}>{labels.noAttachments}</div>
            ) : (
              <div style={{ background: '#fafafa', borderRadius: 4, padding: '4px 0' }}>
                {attachments.map((attachment) => (
                  <div key={attachment.id} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '4px 10px', borderBottom: '1px dashed #f0f0f0' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8, flex: 1, minWidth: 0 }}>
                      <FileOutlined style={{ fontSize: 14, color: '#1890ff' }} />
                      <span style={{ fontSize: 12, color: '#303133', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{attachment.file_name}</span>
                      <span style={{ fontSize: 11, color: '#8c8c8c', flexShrink: 0 }}>{(attachment.file_size / 1024).toFixed(2)} KB</span>
                    </div>
                    <Space size={4}>
                      {isImageFile(attachment.file_name) && (
                        <Button type="link" size="small" icon={<EyeOutlined />} style={{ padding: 0, height: 'auto', fontSize: 12 }} onClick={() => handlePreviewImage(attachment)}>预览</Button>
                      )}
                      <Button type="link" size="small" icon={<DownloadOutlined />} style={{ padding: 0, height: 'auto', fontSize: 12 }} onClick={() => handleDownloadAttachment(attachment)}>{labels.download}</Button>
                      <Popconfirm title={labels.confirmDelete} onConfirm={() => handleDeleteAttachment(attachment.id)}>
                        <Button type="link" size="small" danger icon={<DeleteOutlined />} style={{ padding: 0, height: 'auto', fontSize: 12 }}>{labels.delete}</Button>
                      </Popconfirm>
                    </Space>
                  </div>
                ))}
              </div>
            )}

            {/* 图片预览模态框 */}
            <Modal
              open={previewVisible}
              footer={null}
              onCancel={() => {
                if (previewImage) {
                  URL.revokeObjectURL(previewImage);
                }
                setPreviewImage('');
                setPreviewVisible(false);
              }}
              width="auto"
              centered
              style={{ maxWidth: '90vw' }}
            >
              <div style={{ textAlign: 'center', padding: '20px' }}>
                <Image
                  src={previewImage}
                  alt="预览"
                  style={{ maxWidth: '100%', maxHeight: '80vh' }}
                  preview={false}
                />
              </div>
            </Modal>


          </div>

          <Divider style={{ margin: '12px 0' }} />

          {/* 评论区域 */}
          <CommentSection
            projectId={projectId}
            defectId={defectId}
            currentUserId={currentUser?.id}
            compact={true}
          />
        </Card>
      </div>
    </Spin>
  );
};

export default DefectDetailPage;
