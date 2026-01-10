import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Modal,
  Tag,
  Button,
  Upload,
  message,
  Popconfirm,
  Spin,
  Row,
  Col,
} from 'antd';
import {
  EditOutlined,
  UploadOutlined,
  DownloadOutlined,
  DeleteOutlined,
  PaperClipOutlined,
  CloseOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import {
  DEFECT_STATUS_COLORS,
  DEFECT_STATUS_I18N_KEYS,
  DEFECT_PRIORITY_COLORS,
  DEFECT_PRIORITY_I18N_KEYS,
  DEFECT_SEVERITY_COLORS,
  DEFECT_SEVERITY_I18N_KEYS,
} from '../../../constants/defect';
import {
  fetchDefectAttachments,
  uploadDefectAttachment,
  downloadDefectAttachment,
  deleteDefectAttachment,
} from '../../../api/defect';
import './DefectDetailModal.css';

/**
 * 缺陷详情模态框
 */
const DefectDetailModal = ({
  visible,
  defect,
  projectId,
  onCancel,
  onEdit,
  onRefresh,
}) => {
  const { t } = useTranslation();
  const [attachments, setAttachments] = useState([]);
  const [attachmentLoading, setAttachmentLoading] = useState(false);
  const [uploading, setUploading] = useState(false);

  // 加载附件列表
  const loadAttachments = async () => {
    if (!visible || !defect?.id || !projectId) return;
    setAttachmentLoading(true);
    try {
      const data = await fetchDefectAttachments(projectId, defect.id);
      setAttachments(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Failed to load attachments:', error);
    } finally {
      setAttachmentLoading(false);
    }
  };

  useEffect(() => {
    loadAttachments();
  }, [visible, defect?.id, projectId]);

  // 上传附件
  const handleUpload = async (file) => {
    setUploading(true);
    try {
      await uploadDefectAttachment(projectId, defect.id, file);
      message.success(t('defect.uploadSuccess'));
      loadAttachments();
    } catch (error) {
      console.error('Failed to upload attachment:', error);
      message.error(t('message.uploadFailed'));
    } finally {
      setUploading(false);
    }
    return false; // 阻止默认上传
  };

  // 下载附件
  const handleDownload = async (attachment) => {
    try {
      await downloadDefectAttachment(
        projectId,
        defect.id,
        attachment.id,
        attachment.filename
      );
    } catch (error) {
      console.error('Failed to download attachment:', error);
      message.error(t('message.downloadFailed'));
    }
  };

  // 删除附件
  const handleDeleteAttachment = async (attachment) => {
    try {
      await deleteDefectAttachment(projectId, defect.id, attachment.id);
      message.success(t('message.deleteSuccess'));
      loadAttachments();
    } catch (error) {
      console.error('Failed to delete attachment:', error);
      message.error(t('message.deleteFailed'));
    }
  };

  // 格式化文件大小
  const formatFileSize = (bytes) => {
    if (!bytes) return '-';
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  // 格式化日期
  const formatDate = (dateStr) => {
    if (!dateStr) return '-';
    return dayjs(dateStr).format('YYYY-MM-DD');
  };

  if (!defect) return null;

  // 紧凑型单行字段
  const FieldRow = ({ label, value, children }) => {
    const displayValue = value || '-';
    const isEmpty = !value || value === '-';
    return (
      <div className="dd-field-row">
        <span className="dd-field-label">{label}:</span>
        {children || (
          <span className={`dd-field-value ${isEmpty ? 'empty' : ''}`}>
            {displayValue}
          </span>
        )}
      </div>
    );
  };

  // 多行字段 - 自适应高度
  const MultilineField = ({ label, value }) => {
    const displayValue = value || '-';
    const isEmpty = !value || value === '-';
    if (isEmpty) {
      return (
        <div className="dd-field-row">
          <span className="dd-field-label">{label}:</span>
          <span className="dd-field-value empty">-</span>
        </div>
      );
    }
    return (
      <div className="dd-multiline-field">
        <span className="dd-field-label">{label}:</span>
        <div className="dd-multiline-value">
          {displayValue}
        </div>
      </div>
    );
  };

  return (
    <Modal
      title={
        <div className="dd-modal-title">
          <span>{t('defect.detail')}</span>
          <Tag color="blue" style={{ marginLeft: 8 }}>
            {defect.defect_id}
          </Tag>
          <Tag color={DEFECT_STATUS_COLORS[defect.status]}>
            {t(DEFECT_STATUS_I18N_KEYS[defect.status])}
          </Tag>
        </div>
      }
      open={visible}
      onCancel={onCancel}
      width={750}
      className="dd-modal-compact"
      footer={
        <div className="dd-modal-footer">
          <Button size="small" icon={<CloseOutlined />} onClick={onCancel}>
            {t('common.cancel')}
          </Button>
          <Button type="primary" size="small" icon={<EditOutlined />} onClick={onEdit}>
            {t('common.edit')}
          </Button>
        </div>
      }
      destroyOnClose
    >
      <div className="dd-content">
        {/* 标题 */}
        <div className="dd-summary">
          {defect.summary}
        </div>

        {/* 状态信息行 - 一行显示 */}
        <div className="dd-status-row">
          <Row gutter={16}>
            <Col span={6}>
              <FieldRow label={t('defect.status')}>
                <Tag color={DEFECT_STATUS_COLORS[defect.status]}>
                  {t(DEFECT_STATUS_I18N_KEYS[defect.status])}
                </Tag>
              </FieldRow>
            </Col>
            <Col span={6}>
              <FieldRow label={t('defect.priority')}>
                <Tag color={DEFECT_PRIORITY_COLORS[defect.priority]}>
                  {t(DEFECT_PRIORITY_I18N_KEYS[defect.priority])}
                </Tag>
              </FieldRow>
            </Col>
            <Col span={6}>
              <FieldRow label={t('defect.severity')}>
                <Tag color={DEFECT_SEVERITY_COLORS[defect.severity]}>
                  {t(DEFECT_SEVERITY_I18N_KEYS[defect.severity])}
                </Tag>
              </FieldRow>
            </Col>
            <Col span={6}>
              <FieldRow label={t('defect.recoveryMethod')} value={defect.recovery_method} />
            </Col>
          </Row>
        </div>

        {/* 基本信息行 */}
        <Row gutter={16}>
          <Col span={8}>
            <FieldRow label={t('defect.subject')} value={defect.subject} />
          </Col>
          <Col span={8}>
            <FieldRow label={t('defect.phase')} value={defect.phase} />
          </Col>
          <Col span={8}>
            <FieldRow label={t('defect.caseId')} value={defect.case_id} />
          </Col>
        </Row>

        <Row gutter={16}>
          <Col span={8}>
            <FieldRow label={t('defect.reporter')} value={defect.created_by_user?.username} />
          </Col>
          <Col span={8}>
            <FieldRow label={t('defect.createdAt')} value={formatDate(defect.created_at)} />
          </Col>
          <Col span={8}>
            <FieldRow label={t('defect.updatedAt')} value={formatDate(defect.updated_at)} />
          </Col>
        </Row>

        {/* 描述信息 */}
        <MultilineField label={t('defect.description')} value={defect.description} />

        {/* 附件区域 */}
        <div className="dd-attachments-section">
          <div className="dd-section-header">
            <span className="dd-section-title">
              <PaperClipOutlined /> {t('defect.attachments')} ({attachments.length})
            </span>
            <Upload
              showUploadList={false}
              beforeUpload={handleUpload}
              disabled={uploading}
            >
              <Button
                size="small"
                icon={<UploadOutlined />}
                loading={uploading}
              >
                {t('defect.uploadAttachment')}
              </Button>
            </Upload>
          </div>

          <Spin spinning={attachmentLoading}>
            <div className="dd-attachment-list">
              {attachments.length === 0 ? (
                <div className="dd-no-data">
                  {t('common.noData') || '暂无附件'}
                </div>
              ) : (
                attachments.map((att) => (
                  <div key={att.id} className="dd-attachment-item">
                    <div className="dd-attachment-info">
                      <PaperClipOutlined />
                      <span
                        className="dd-attachment-name"
                        onClick={() => handleDownload(att)}
                      >
                        {att.filename}
                      </span>
                      <span className="dd-attachment-size">
                        ({formatFileSize(att.file_size)})
                      </span>
                    </div>
                    <div className="dd-attachment-actions">
                      <Button
                        type="link"
                        size="small"
                        icon={<DownloadOutlined />}
                        onClick={() => handleDownload(att)}
                      />
                      <Popconfirm
                        title={t('defect.deleteAttachment') + '?'}
                        onConfirm={() => handleDeleteAttachment(att)}
                        okText={t('common.confirm')}
                        cancelText={t('common.cancel')}
                      >
                        <Button
                          type="link"
                          size="small"
                          danger
                          icon={<DeleteOutlined />}
                        />
                      </Popconfirm>
                    </div>
                  </div>
                ))
              )}
            </div>
          </Spin>
        </div>
      </div>
    </Modal>
  );
};

export default DefectDetailModal;
