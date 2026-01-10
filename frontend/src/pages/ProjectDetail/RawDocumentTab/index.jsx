import React, { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Card,
  Table,
  Button,
  Upload,
  Space,
  message,
  Popconfirm,
  Tag,
  Typography,
  Tooltip,
  Progress,
  Modal,
} from 'antd';
import {
  UploadOutlined,
  DownloadOutlined,
  DeleteOutlined,
  ReloadOutlined,
  SyncOutlined,
  EyeOutlined,
  CopyOutlined,
} from '@ant-design/icons';
import {
  uploadRawDocument,
  fetchRawDocuments,
  convertRawDocument,
  getConvertStatus,
  downloadOriginalDocument,
  downloadConvertedDocument,
  deleteOriginalDocument,
  deleteConvertedDocument,
  previewConvertedDocument,
} from '../../../api/rawDocument';
import './index.css';

const { Text } = Typography;

const RawDocumentTab = ({ projectId }) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [documents, setDocuments] = useState([]);
  const [uploading, setUploading] = useState(false);
  
  // 预览相关状态
  const [previewVisible, setPreviewVisible] = useState(false);
  const [previewContent, setPreviewContent] = useState('');
  const [previewFilename, setPreviewFilename] = useState('');
  const [previewLoading, setPreviewLoading] = useState(false);

  // 加载文档列表
  const loadDocuments = useCallback(async () => {
    setLoading(true);
    try {
      const data = await fetchRawDocuments(projectId);
      // 确保返回的是数组
      const docs = Array.isArray(data) ? data : [];
      // 调试：打印文档的转换状态
      console.log('[RawDocumentTab] documents loaded:', docs.map(d => ({
        id: d.id,
        filename: d.original_filename,
        convert_status: d.convert_status,
        convert_progress: d.convert_progress
      })));
      setDocuments(docs);
    } catch (error) {
      message.error(t('rawDocument.loadFailed'));
      setDocuments([]);
    } finally {
      setLoading(false);
    }
  }, [projectId, t]);

  useEffect(() => {
    if (projectId) {
      loadDocuments();
    }
  }, [projectId, loadDocuments]);
  
  // 文件上传配置 - 支持需求文档规定的15种格式
  const uploadProps = {
    name: 'file',
    showUploadList: false,
    // 文本文档: PDF/DOCX/DOC/TXT/RTF, 图片: PNG/JPG/JPEG/BMP/TIFF, 表格: XLSX/XLS/CSV, 演示: PPTX/PPT
    accept: '.pdf,.docx,.doc,.txt,.rtf,.png,.jpg,.jpeg,.bmp,.tiff,.xlsx,.xls,.csv,.pptx,.ppt',
    beforeUpload: (file) => {
      // 检查文件大小（100MB）
      const isLt100M = file.size / 1024 / 1024 < 100;
      if (!isLt100M) {
        message.error(t('rawDocument.fileTooLarge'));
        return Upload.LIST_IGNORE;
      }
      return true;
    },
    customRequest: async ({ file, onSuccess, onError }) => {
      setUploading(true);
      try {
        const formData = new FormData();
        formData.append('file', file);
        await uploadRawDocument(projectId, formData);
        message.success(t('rawDocument.uploadSuccess'));
        onSuccess();
        loadDocuments();
      } catch (error) {
        message.error(error.message || t('rawDocument.uploadFailed'));
        onError(error);
      } finally {
        setUploading(false);
      }
    },
  };

  /**
   * 转换文档 - 触发后端异步转换任务
   * 
   * 流程：
   * 1. 调用 convertRawDocument API 启动转换
   * 2. 显示转换已启动提示
   * 3. 启动轮询检查转换状态
   * 
   * 错误处理：
   * - 404: 文档不存在
   * - 409: 文档正在转换中（重复请求）
   * - 其他: 显示通用错误信息
   * 
   * @param {Object} record - 文档记录对象，包含 id 和 original_filename
   */
  const handleConvert = async (record) => {
    try {
      console.log(`[Convert Start] documentId=${record.id}, originalFilename=${record.original_filename}`);
      await convertRawDocument(record.id);
      message.success(t('rawDocument.convertStarted'));
      // 轮询检查转换状态
      checkConvertStatus(record.id);
    } catch (error) {
      console.error(`[Convert Failed] documentId=${record.id}, error:`, error);
      if (error.response?.status === 404) {
        message.error(t('rawDocument.documentNotFound'));
      } else if (error.response?.status === 409) {
        message.error(t('rawDocument.convertInProgress'));
      } else {
        message.error(error.message || t('rawDocument.convertFailed'));
      }
    }
  };

  /**
   * 检查转换状态 - 轮询机制
   * 
   * 轮询参数：
   * - 最大尝试次数: 30次
   * - 间隔: 2秒
   * - 总超时: 60秒
   * - 单次失败重试: 最多3次
   * - 进度卡死检测: 10秒后进度仍为0则超时
   * 
   * 状态处理：
   * - completed: 显示成功提示，刷新列表
   * - failed: 显示错误信息，刷新列表
   * - processing: 继续轮询
   * - 超时: 显示超时提示
   * 
   * @param {number} documentId - 文档ID
   */
  const checkConvertStatus = async (documentId) => {
    const maxAttempts = 30; // 最多检查30次（总共60秒）
    const progressTimeout = 5; // 10秒后检查进度（5次轮询 * 2秒）
    let attempts = 0;
    let retryCount = 0;
    const maxRetries = 3;
    const startTime = Date.now();

    const checkStatus = async () => {
      try {
        console.log(`[Polling Attempt] documentId=${documentId}, attempt=${attempts + 1}, status=polling`);
        const status = await getConvertStatus(documentId);
        
        if (status.status === 'completed') {
          console.log(`[Polling Complete] documentId=${documentId}, finalStatus=completed`);
          message.success(t('rawDocument.convertSuccess'));
          loadDocuments();
          return;
        } else if (status.status === 'failed') {
          console.log(`[Polling Complete] documentId=${documentId}, finalStatus=failed, error=${status.error_message}`);
          message.error(status.error_message || t('rawDocument.convertFailed'));
          loadDocuments();
          return;
        } else if (status.status === 'processing') {
          attempts++;
          retryCount = 0; // 重置重试计数
          
          // 检测进度卡死：10秒后进度仍为0则判定超时
          const elapsedSeconds = (Date.now() - startTime) / 1000;
          if (attempts >= progressTimeout && status.progress === 0) {
            console.log(`[Polling Timeout] documentId=${documentId}, progress stuck at 0 after ${elapsedSeconds}s`);
            message.warning(t('rawDocument.convertProgressTimeout') || '转换进度超时，请检查文档格式或稍后重试');
            loadDocuments();
            return;
          }
          
          if (attempts < maxAttempts) {
            console.log(`[Polling Attempt] documentId=${documentId}, attempt=${attempts}/${maxAttempts}, progress=${status.progress}%, continuing...`);
            setTimeout(checkStatus, 2000); // 2秒后再次检查
          } else {
            console.log(`[Polling Complete] documentId=${documentId}, finalStatus=timeout`);
            message.warning(t('rawDocument.convertTimeout'));
            loadDocuments();
          }
        }
      } catch (error) {
        console.error(`[Polling Error] documentId=${documentId}, error:`, error);
        retryCount++;
        if (retryCount < maxRetries) {
          console.log(`[Polling Retry] documentId=${documentId}, retryCount=${retryCount}/${maxRetries}`);
          setTimeout(checkStatus, 2000); // 重试
        } else {
          message.error(t('rawDocument.statusCheckFailed'));
          loadDocuments();
        }
      }
    };

    console.log(`[Polling Start] documentId=${documentId}, maxAttempts=${maxAttempts}`);
    checkStatus();
  };

  // 下载原始文档
  const handleDownloadOriginal = async (record) => {
    try {
      await downloadOriginalDocument(record.id, record.original_filename);
      message.success(t('rawDocument.downloadSuccess'));
    } catch (error) {
      message.error(t('rawDocument.downloadFailed'));
    }
  };

  // 下载转换后的文档
  const handleDownloadConverted = async (record) => {
    try {
      await downloadConvertedDocument(record.id, record.converted_filename);
      message.success(t('rawDocument.downloadSuccess'));
    } catch (error) {
      message.error(t('rawDocument.downloadFailed'));
    }
  };

  // 删除原始文档
  const handleDeleteOriginal = async (record) => {
    try {
      await deleteOriginalDocument(record.id);
      message.success(t('rawDocument.deleteSuccess'));
      loadDocuments();
    } catch (error) {
      message.error(t('rawDocument.deleteFailed'));
    }
  };

  // 删除转换后的文档
  const handleDeleteConverted = async (record) => {
    try {
      await deleteConvertedDocument(record.id);
      message.success(t('rawDocument.deleteConvertedSuccess'));
      loadDocuments();
    } catch (error) {
      message.error(t('rawDocument.deleteFailed'));
    }
  };

  // 预览转换后的文档
  const handlePreviewConverted = async (record) => {
    setPreviewLoading(true);
    setPreviewVisible(true);
    try {
      const response = await previewConvertedDocument(record.id);
      setPreviewFilename(response.filename || record.converted_filename);
      setPreviewContent(response.content || '');
    } catch (error) {
      message.error(t('rawDocument.previewFailed'));
      setPreviewVisible(false);
    } finally {
      setPreviewLoading(false);
    }
  };

  // 关闭预览模态框
  const handlePreviewClose = () => {
    setPreviewVisible(false);
    setPreviewContent('');
    setPreviewFilename('');
  };

  // 格式化时间
  const formatTime = (time) => {
    if (!time) return '-';
    try {
      return new Date(time).toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      });
    } catch (e) {
      return time;
    }
  };

  // 格式化文件大小
  const formatFileSize = (bytes) => {
    if (!bytes || bytes === 0) return '-';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  // 表格列配置 - No./原始文档/文件大小/上传时间/操作/转换文档/文件大小/作成时间/操作
  const columns = [
    {
      title: 'No.',
      key: 'index',
      width: 60,
      align: 'center',
      render: (_, __, index) => index + 1,
    },
    {
      title: t('rawDocument.originalDocument'),
      dataIndex: 'original_filename',
      key: 'original_filename',
      render: (text) => (
        <div style={{ 
          display: 'flex', 
          alignItems: 'center', 
          gap: '8px',
          minWidth: 0,
          width: '100%'
        }}>
          <Tooltip placement="topLeft" title={text}>
            <span style={{
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
              flex: 1,
              minWidth: 0
            }}>{text}</span>
          </Tooltip>
          <Button
            type="text"
            size="small"
            icon={<CopyOutlined />}
            title={t('rawDocument.copyFilename')}
            onClick={(e) => {
              e.stopPropagation();
              navigator.clipboard.writeText(text);
              message.success(t('rawDocument.copySuccess'));
            }}
            style={{ 
              flexShrink: 0, 
              padding: '4px',
              minWidth: '24px',
              height: '24px'
            }}
          />
        </div>
      ),
    },
    {
      title: t('rawDocument.fileSize'),
      dataIndex: 'file_size',
      key: 'file_size',
      width: 100,
      render: (size) => formatFileSize(size),
    },
    {
      title: t('rawDocument.uploadTime'),
      dataIndex: 'created_at',
      key: 'created_at',
      width: 170,
      render: (time) => formatTime(time),
    },
    {
      title: t('rawDocument.originalActions'),
      key: 'original_actions',
      width: 220,
      render: (_, record) => {
        // 判断是否真正在处理中
        const isProcessing = record.convert_status === 'processing';
        return (
        <div>
          <Space size="small" style={{ display: 'flex', flexWrap: 'wrap', gap: '4px', alignItems: 'center' }}>
            {/* 转换按钮 - 所有状态都可以转换/重新转换 */}
            <Tooltip title={record.convert_status === 'completed' ? t('rawDocument.reconvert') : t('rawDocument.convertToMarkdown')}>
              <Button
                type="link"
                size="small"
                icon={<SyncOutlined spin={isProcessing} />}
                onClick={() => handleConvert(record)}
                disabled={isProcessing}
              />
            </Tooltip>
            
              {/* 下载原始文件 */}
            <Tooltip title={t('rawDocument.downloadOriginal')}>
              <Button
                type="link"
                size="small"
                icon={<DownloadOutlined />}
                onClick={() => handleDownloadOriginal(record)}
              />
            </Tooltip>
            
            {/* 删除原始文档 */}
            <Popconfirm
              title={t('rawDocument.deleteConfirm')}
              description={t('rawDocument.deleteWarning')}
              onConfirm={() => handleDeleteOriginal(record)}
              okText={t('common.confirm')}
              cancelText={t('common.cancel')}
            >
            <Button type="link" size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
          
          {/* 转换进度条 - processing状态时显示 */}
          {isProcessing && (
            <Progress 
              percent={record.convert_progress || 0} 
              size="small" 
              status="active"
              style={{ marginLeft: 4, width: 60, marginTop: 0 }}
            />
          )}
          {record.convert_status === 'completed' && (
            <Tag color="success" size="small" style={{ marginLeft: 4, marginTop: 0 }}>
              {t('rawDocument.statusCompleted')}
            </Tag>
          )}
          {record.convert_status === 'failed' && record.convert_error && (
            <Tooltip title={record.convert_error}>
              <Tag color="error" size="small" style={{ marginLeft: 4, marginTop: 0 }}>
                {t('rawDocument.statusFailed')}
              </Tag>
            </Tooltip>
          )}
          </Space>
        </div>
      );
      },
    },
    {
      title: t('rawDocument.convertedDocument'),
      dataIndex: 'converted_filename',
      key: 'converted_filename',
      render: (text, record) => {
        if (!text || record.convert_status !== 'completed') {
          return '-';
        }
        return (
          <div style={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: '8px',
            minWidth: 0,
            width: '100%'
          }}>
            <Tooltip placement="topLeft" title={text}>
              <span style={{
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
                flex: 1,
                minWidth: 0
              }}>{text}</span>
            </Tooltip>
            <Button
              type="text"
              size="small"
              icon={<CopyOutlined />}
              title={t('rawDocument.copyFilename')}
              onClick={(e) => {
                e.stopPropagation();
                navigator.clipboard.writeText(text);
                message.success(t('rawDocument.copySuccess'));
              }}
              style={{ 
                flexShrink: 0, 
                padding: '4px',
                minWidth: '24px',
                height: '24px'
              }}
            />
          </div>
        );
      },
    },
    {
      title: t('rawDocument.convertedFileSize'),
      dataIndex: 'converted_file_size',
      key: 'converted_file_size',
      width: 100,
      render: (size, record) => (record.convert_status === 'completed' && size) ? formatFileSize(size) : '-',
    },
    {
      title: t('rawDocument.convertedTime'),
      dataIndex: 'converted_time',
      key: 'converted_time',
      width: 170,
      render: (time, record) => (record.convert_status === 'completed' && time) ? formatTime(time) : '-',
    },
    {
      title: t('rawDocument.convertedActions'),
      key: 'converted_actions',
      width: 120,
      render: (_, record) => (
        <Space size="small">
          {record.convert_status === 'completed' && record.converted_filename && (
            <>
              {/* 预览转换文件 */}
              <Tooltip title={t('rawDocument.preview')}>
                <Button
                  type="link"
                  size="small"
                  icon={<EyeOutlined />}
                  onClick={() => handlePreviewConverted(record)}
                />
              </Tooltip>
              
              {/* 下载转换文件 */}
              <Tooltip title={t('rawDocument.downloadConverted')}>
                <Button
                  type="link"
                  size="small"
                  icon={<DownloadOutlined />}
                  onClick={() => handleDownloadConverted(record)}
                />
              </Tooltip>
              
              {/* 删除转换文件 */}
              <Popconfirm
                title={t('rawDocument.deleteConvertedConfirm')}
                onConfirm={() => handleDeleteConverted(record)}
                okText={t('common.confirm')}
                cancelText={t('common.cancel')}
              >
                <Button type="link" size="small" danger icon={<DeleteOutlined />} />
              </Popconfirm>
            </>
          )}
          {(!record.converted_filename) && '-'}
        </Space>
      ),
    },
  ];

  return (
    <div className="raw-document-tab">
      <Card
        title={t('rawDocument.title')}
        extra={
          <Space>
            <Upload {...uploadProps}>
              <Button type="primary" icon={<UploadOutlined />} loading={uploading}>
                {t('rawDocument.upload')}
              </Button>
            </Upload>
            <Button icon={<ReloadOutlined />} onClick={loadDocuments}>
              {t('rawDocument.refresh')}
            </Button>
          </Space>
        }
      >
        <div className="raw-document-tip">
          <Text type="secondary">{t('rawDocument.description')}</Text>
        </div>

        <Table
          columns={columns}
          dataSource={documents}
          rowKey="id"
          loading={loading}
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => t('rawDocument.total', { count: total }),
          }}
        />
      </Card>
      
      {/* 预览模态框 */}
      <Modal
        title={previewFilename || t('rawDocument.preview')}
        open={previewVisible}
        onCancel={handlePreviewClose}
        footer={[
          <Button key="close" onClick={handlePreviewClose}>
            {t('common.close')}
          </Button>,
        ]}
        width={800}
        style={{ top: 20 }}
        styles={{ body: { maxHeight: 'calc(100vh - 200px)', overflow: 'auto' } }}
      >
        {previewLoading ? (
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <SyncOutlined spin style={{ fontSize: 24 }} />
            <div style={{ marginTop: 8 }}>{t('common.loading')}</div>
          </div>
        ) : (
          <pre style={{ 
            whiteSpace: 'pre-wrap', 
            wordWrap: 'break-word',
            backgroundColor: '#f5f5f5',
            padding: 16,
            borderRadius: 4,
            maxHeight: 'calc(100vh - 280px)',
            overflow: 'auto',
            fontFamily: 'Monaco, Consolas, "Courier New", monospace',
            fontSize: 13,
            lineHeight: 1.6
          }}>
            {previewContent}
          </pre>
        )}
      </Modal>
    </div>
  );
};

export default RawDocumentTab;
