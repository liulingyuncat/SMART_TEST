import React, { useState, useEffect, useCallback, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Space, Typography, message, Popconfirm } from 'antd';
import { SaveOutlined, DownloadOutlined, EditOutlined, CloseOutlined, ImportOutlined, DeleteOutlined } from '@ant-design/icons';
import MarkdownIt from 'markdown-it';
import MdEditor from 'react-markdown-editor-lite';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import 'react-markdown-editor-lite/lib/index.css';
import { debounce } from 'lodash';
import PropTypes from 'prop-types';
import { saveVersion } from '../../../api/requirement';
import './MarkdownEditor.css';

const { Text } = Typography;

const MarkdownEditor = ({ 
  value, 
  onChange, 
  onSave, 
  onSaveVersion,
  onEditCancel,
  onDelete,
  showImport, 
  projectName, 
  projectId,
  docType 
}) => {
  const { t } = useTranslation();
  const [isEditing, setIsEditing] = useState(false);
  const [originalContent, setOriginalContent] = useState('');
  const [saveStatus, setSaveStatus] = useState('idle'); // idle, saving, saved, failed
  const [lastSavedAt, setLastSavedAt] = useState(null);
  const mdParser = useRef(new MarkdownIt());

  // 自动保存(防抖3秒) - 仅在编辑模式下启用
  const autoSave = useCallback(
    debounce(async () => {
      if (!isEditing) return; // 只读模式下禁用自动保存
      
      setSaveStatus('saving');
      const success = await onSave();
      if (success) {
        setSaveStatus('saved');
        setLastSavedAt(new Date());
      } else {
        setSaveStatus('failed');
      }
    }, 3000),
    [onSave, isEditing]
  );

  // 监听内容变化,触发自动保存
  useEffect(() => {
    if (isEditing && value !== undefined && value !== null) {
      autoSave();
    }
    // 清理防抖函数
    return () => {
      autoSave.cancel();
    };
  }, [value, autoSave, isEditing]);

  // 手动保存
  const handleManualSave = async () => {
    autoSave.cancel(); // 取消待执行的自动保存
    setSaveStatus('saving');
    const success = await onSave();
    if (success) {
      setSaveStatus('saved');
      setLastSavedAt(new Date());
      message.success(t('requirement.save') + ' ' + t('message.saveSuccess'));
      setIsEditing(false); // 保存成功后切换到只读模式
    } else {
      setSaveStatus('failed');
      message.error(t('requirement.save') + ' ' + t('message.saveFailed'));
    }
  };

  // 进入编辑模式
  const handleEdit = () => {
    setOriginalContent(value || '');
    setIsEditing(true);
  };

  // 取消编辑
  const handleCancel = () => {
    onChange(originalContent);
    setIsEditing(false);
    setSaveStatus('idle'); // 重置保存状态
    if (onEditCancel) {
      onEditCancel(); // 通知父组件取消编辑
    }
  };

  // 版本保存
  const handleSaveVersion = async () => {
    console.log('[handleSaveVersion] 开始保存版本...');
    console.log('[handleSaveVersion] projectId:', projectId);
    console.log('[handleSaveVersion] docType:', docType);
    console.log('[handleSaveVersion] content length:', value?.length);
    
    if (!value || value.trim() === '') {
      console.log('[handleSaveVersion] 内容为空,取消保存');
      message.warning(t('requirement.emptyContent'));
      return;
    }
    
    try {
      setSaveStatus('saving');
      console.log('[handleSaveVersion] 调用saveVersion API...');
      const result = await saveVersion(projectId, docType, value);
      console.log('[handleSaveVersion] API响应:', result);
      
      const filename = result?.filename || '未知文件名';
      message.success(`版本保存成功: ${filename}`);
      
      if (onSaveVersion) {
        console.log('[handleSaveVersion] 触发版本刷新回调');
        onSaveVersion(); // 通知父组件刷新版本列表
      }
      
      setSaveStatus('saved');
      setIsEditing(false); // 切换到只读模式
    } catch (error) {
      console.error('版本保存失败:', error);
      console.error('错误详情:', error.response || error);
      message.error(`版本保存失败: ${error.message || error}`);
      setSaveStatus('failed');
    }
  };

  // 导入Markdown
  const handleImport = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.md';
    
    input.onchange = async (e) => {
      const file = e.target.files[0];
      if (!file) return;
      
      // 验证文件大小
      const maxSize = 5 * 1024 * 1024; // 5MB
      if (file.size > maxSize) {
        message.error(t('requirement.fileTooLarge'));
        return;
      }
      
      // 读取文件内容
      const reader = new FileReader();
      reader.onload = async (event) => {
        const content = event.target.result;
        // 清空原内容并显示新内容
        onChange(content);
        // 切换到编辑模式
        setIsEditing(true);
        setOriginalContent(value); // 保存当前内容作为原始内容
        message.success(t('requirement.importSuccess'));
      };
      reader.onerror = () => {
        message.error(t('requirement.importFailed'));
      };
      reader.readAsText(file, 'UTF-8');
    };
    
    input.click();
  };

  // 下载功能
  const handleDownload = () => {
    if (!value) {
      message.warning('文档内容为空');
      return;
    }

    const blob = new Blob([value], { type: 'text/markdown;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    
    // 生成文件名: 项目名称-文档类型-YYYYMMDD.md
    const date = new Date();
    const dateStr = date.toISOString().split('T')[0].replace(/-/g, '');
    const docTypeName = t(`requirement.${docType.replace(/-/g, '')}`);
    link.download = `${projectName}-${docTypeName}-${dateStr}.md`;
    
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
    
    message.success(t('requirement.download') + ' ' + t('message.saveSuccess'));
  };

  // 编辑器内容变更
  const handleEditorChange = ({ text }) => {
    onChange(text);
  };

  // 渲染保存状态
  const renderSaveStatus = () => {
    if (saveStatus === 'saving') {
      return <Text type="secondary">{t('requirement.saving')}</Text>;
    }
    if (saveStatus === 'saved' && lastSavedAt) {
      const timeStr = lastSavedAt.toLocaleTimeString('zh-CN', { 
        hour: '2-digit', 
        minute: '2-digit', 
        second: '2-digit' 
      });
      return <Text type="success">{t('requirement.saved')} {timeStr}</Text>;
    }
    if (saveStatus === 'failed') {
      return <Text type="danger">{t('requirement.saveFailed')}</Text>;
    }
    return null;
  };

  return (
    <div className="markdown-editor-container">
      <div className="editor-toolbar">
        <Space>
          {isEditing ? (
            <>
              {/* 编辑模式工具栏：[保存] [取消] */}
              <Button 
                type="primary" 
                icon={<SaveOutlined />} 
                onClick={handleManualSave}
                loading={saveStatus === 'saving'}
              >
                {t('common.save')}
              </Button>
              <Button 
                icon={<CloseOutlined />} 
                onClick={handleCancel}
              >
                {t('common.cancel')}
              </Button>
              {renderSaveStatus()}
            </>
          ) : (
            <>
              {/* 只读模式工具栏：[编辑] [导入] [删除] [下载] */}
              <Button 
                type="primary" 
                icon={<EditOutlined />} 
                onClick={handleEdit}
              >
                编辑
              </Button>
              {showImport && (
                <Button 
                  icon={<ImportOutlined />} 
                  onClick={handleImport}
                >
                  {t('requirement.import')}
                </Button>
              )}
              {onDelete && (
                <Popconfirm
                  title={t('requirement.confirmDelete')}
                  onConfirm={onDelete}
                  okText="确定"
                  cancelText={t('common.cancel')}
                >
                  <Button 
                    danger
                    icon={<DeleteOutlined />}
                  >
                    删除
                  </Button>
                </Popconfirm>
              )}
              <Button 
                icon={<DownloadOutlined />} 
                onClick={handleDownload}
              >
                下载
              </Button>
            </>
          )}
        </Space>
      </div>
      
      {isEditing ? (
        <MdEditor
          value={value || ''}
          style={{ flex: 1, height: 'auto' }}
          renderHTML={(text) => mdParser.current.render(text)}
          onChange={handleEditorChange}
          config={{
            view: {
              menu: true,
              md: true,
              html: true,
            },
            canView: {
              menu: true,
              md: true,
              html: true,
              fullScreen: true,
              hideMenu: true,
            },
          }}
        />
      ) : (
        <div className="markdown-preview" style={{ flex: 1, overflow: 'auto' }}>
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            components={{
              table: ({ node, ...props }) => <table className="markdown-table" {...props} />,
              th: ({ node, ...props }) => <th {...props} />,
              td: ({ node, ...props }) => <td {...props} />,
            }}
          >
            {value || ''}
          </ReactMarkdown>
        </div>
      )}
    </div>
  );
};

MarkdownEditor.propTypes = {
  value: PropTypes.string,
  onChange: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
  onSaveVersion: PropTypes.func,
  onEditCancel: PropTypes.func,
  showImport: PropTypes.bool,
  projectName: PropTypes.string.isRequired,
  projectId: PropTypes.string,
  docType: PropTypes.string.isRequired,
};

MarkdownEditor.defaultProps = {
  value: '',
  onSaveVersion: null,
  onEditCancel: null,
  showImport: false,
  projectId: '',
};

export default MarkdownEditor;
