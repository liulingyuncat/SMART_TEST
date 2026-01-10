import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Button, message, Space, List, Modal, Input, Popconfirm, Empty, Spin, Tooltip } from 'antd';
import { 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  DownloadOutlined,
  SaveOutlined,
  CloseOutlined,
  FileTextOutlined,
  CopyOutlined
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import {
  getReviewItems,
  createReviewItem,
  getReviewItem,
  updateReviewItem,
  deleteReviewItem,
} from '../../../../api/reviewItem';
import MarkdownIt from 'markdown-it';
import MdEditor from 'react-markdown-editor-lite';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import 'react-markdown-editor-lite/lib/index.css';
import './CaseReviewTab.css';

// 初始化 Markdown 解析器
const mdParser = new MarkdownIt();

/**
 * 用例审阅Tab组件 (T44重构版 - 参考AI质量报告)
 * @param {number} projectId - 项目ID
 */
const CaseReviewTab = ({ projectId }) => {
  const { t } = useTranslation();
  
  // 左栏状态
  const [items, setItems] = useState([]);
  const [selectedItem, setSelectedItem] = useState(null);
  const [loadingItems, setLoadingItems] = useState(false);
  
  // 右栏状态
  const [isEditing, setIsEditing] = useState(false);
  const [editingContent, setEditingContent] = useState('');
  const [saveLoading, setSaveLoading] = useState(false);
  const [tocItems, setTocItems] = useState([]);
  
  // 新建审阅Modal
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [newItemName, setNewItemName] = useState('');
  const [creating, setCreating] = useState(false);
  
  // 编辑名称Modal
  const [editNameModalVisible, setEditNameModalVisible] = useState(false);
  const [editingName, setEditingName] = useState('');
  const [editingItemId, setEditingItemId] = useState(null);

  const editorRef = useRef(null);
  const previewRef = useRef(null);
  
  // 将标题转换为ID
  const titleToId = useCallback((title) => {
    return 'heading-' + title
      .toLowerCase()
      .replace(/[^\u4e00-\u9fa5a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }, []);

  // 解析Markdown生成目录
  const generateTOC = useCallback((markdown) => {
    if (!markdown) return [];
    
    const headings = [];
    const lines = markdown.split('\n');
    
    lines.forEach((line) => {
      const match = line.match(/^(#{1,6})\s+(.+)$/);
      if (match) {
        const level = match[1].length;
        const title = match[2].trim();
        const id = titleToId(title);
        headings.push({ level, title, id });
      }
    });
    
    return headings;
  }, [titleToId]);

  // 加载审阅列表
  useEffect(() => {
    if (projectId) {
      loadItems();
    }
  }, [projectId]);

  const loadItems = async () => {
    setLoadingItems(true);
    try {
      const data = await getReviewItems(projectId);
      setItems(data || []);
    } catch (error) {
      console.error('加载审阅列表失败:', error);
      message.error('加载审阅列表失败');
    } finally {
      setLoadingItems(false);
    }
  };

  const handleCreate = () => {
    setNewItemName('');
    setCreateModalVisible(true);
  };

  const handleCreateConfirm = async () => {
    if (!newItemName.trim()) {
      message.warning('审阅名称不能为空');
      return;
    }

    setCreating(true);
    try {
      await createReviewItem(projectId, newItemName.trim());
      message.success('创建成功');
      setCreateModalVisible(false);
      setNewItemName('');
      await loadItems();
    } catch (error) {
      console.error('创建审阅失败:', error);
      const errorMsg = error.response?.data?.error || '创建审阅失败';
      // 使用红色错误提示
      if (errorMsg.includes('已存在') || errorMsg.includes('重复')) {
        message.error({ content: errorMsg, style: { color: '#ff4d4f' } });
      } else {
        message.error(errorMsg);
      }
    } finally {
      setCreating(false);
    }
  };

  const handleSelectItem = async (item) => {
    setSelectedItem(item);
    setIsEditing(false);
    try {
      const data = await getReviewItem(projectId, item.id);
      const content = data.content || '';
      setEditingContent(content);
      setTocItems(generateTOC(content));
    } catch (error) {
      console.error('加载审阅内容失败:', error);
      message.error('加载审阅内容失败');
    }
  };

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleSave = async () => {
    if (!selectedItem) return;

    if (!selectedItem.name.trim()) {
      message.warning('审阅名称不能为空');
      return;
    }

    setSaveLoading(true);
    try {
      await updateReviewItem(projectId, selectedItem.id, { 
        name: selectedItem.name.trim(),
        content: editingContent 
      });
      message.success('保存成功');
      setIsEditing(false);
      await loadItems();
    } catch (error) {
      console.error('保存审阅失败:', error);
      message.error('保存审阅失败');
    } finally {
      setSaveLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!selectedItem) return;

    try {
      await deleteReviewItem(projectId, selectedItem.id);
      message.success('删除成功');
      setSelectedItem(null);
      setEditingContent('');
      await loadItems();
    } catch (error) {
      console.error('删除审阅失败:', error);
      message.error('删除审阅失败');
    }
  };
  
  const handleEditName = (item) => {
    setEditingItemId(item.id);
    setEditingName(item.name);
    setEditNameModalVisible(true);
  };
  
  const handleEditNameConfirm = async () => {
    if (!editingName.trim()) {
      message.warning('审阅名称不能为空');
      return;
    }
    
    try {
      await updateReviewItem(projectId, editingItemId, { name: editingName.trim() });
      message.success('修改成功');
      setEditNameModalVisible(false);
      await loadItems();
      // 如果当前选中的项目是被编辑的，更新其名称
      if (selectedItem && selectedItem.id === editingItemId) {
        setSelectedItem({ ...selectedItem, name: editingName.trim() });
      }
    } catch (error) {
      console.error('修改名称失败:', error);
      const errorMsg = error.response?.data?.error || '修改名称失败';
      if (errorMsg.includes('已存在') || errorMsg.includes('重复')) {
        message.error({ content: errorMsg, style: { color: '#ff4d4f' } });
      } else {
        message.error(errorMsg);
      }
    }
  };
  
  const handleDeleteItem = async (item) => {
    try {
      await deleteReviewItem(projectId, item.id);
      message.success('删除成功');
      if (selectedItem && selectedItem.id === item.id) {
        setSelectedItem(null);
        setEditingContent('');
      }
      await loadItems();
    } catch (error) {
      console.error('删除审阅失败:', error);
      message.error('删除审阅失败');
    }
  };

  const handleDownload = () => {
    if (!selectedItem) return;

    try {
      const now = new Date();
      const timestamp = now.toISOString().split('T')[0].replace(/-/g, '');
      const filename = `Project${projectId}_${selectedItem.name}_${timestamp}.md`;
      
      const element = document.createElement('a');
      element.setAttribute('href', 'data:text/markdown;charset=utf-8,' + encodeURIComponent(editingContent));
      element.setAttribute('download', filename);
      element.style.display = 'none';
      document.body.appendChild(element);
      element.click();
      document.body.removeChild(element);
      message.success('下载成功');
    } catch (error) {
      console.error('下载审阅失败:', error);
      message.error('下载审阅失败');
    }
  };

  const handleEditorChange = ({ text }) => {
    setEditingContent(text);
  };

  return (
    <div className="case-review-tab-container">
      {/* 左栏 - 审阅列表 */}
      <div className="review-left-panel">
        <div className="left-panel-header">
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={handleCreate}
            size="small"
            block
          >
            {t('manualTest.createReview')}
          </Button>
        </div>

        <div className="left-panel-content">
          {loadingItems ? (
            <div style={{ textAlign: 'center', padding: '20px' }}>
              <Spin />
            </div>
          ) : items.length === 0 ? (
            <Empty description="暂无审阅记录" style={{ marginTop: '20px' }} />
          ) : (
            <div className="review-list">
              {items.map((item) => (
                <div
                  key={item.id}
                  className={`review-list-item ${selectedItem?.id === item.id ? 'selected' : ''}`}
                  onClick={() => handleSelectItem(item)}
                >
                  <div className="review-item-content">
                    <div className="review-item-name">{item.name}</div>
                    <div className="review-item-date">
                      {new Date(item.created_at).toLocaleDateString('zh-CN', {
                        year: '2-digit',
                        month: '2-digit',
                        day: '2-digit'
                      }).replace(/\//g, '/')}
                    </div>
                  </div>
                  <div className="review-item-actions" onClick={(e) => e.stopPropagation()}>
                    <Tooltip title="复制文档名">
                      <CopyOutlined
                        className="action-icon"
                        onClick={() => {
                          navigator.clipboard.writeText(item.name);
                          message.success('复制成功');
                        }}
                      />
                    </Tooltip>
                    <Popconfirm
                      title="确定删除此审阅?"
                      onConfirm={() => handleDeleteItem(item)}
                      okText="确定"
                      cancelText="取消"
                    >
                      <DeleteOutlined className="action-icon" title="删除" />
                    </Popconfirm>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* 右栏 - 内容编辑/预览 */}
      <div className="review-right-panel">
        {selectedItem ? (
          <>
            {/* 操作按钮 */}
            <div className="right-panel-header">
              {isEditing ? (
                <Input
                  value={selectedItem.name}
                  onChange={(e) => setSelectedItem({ ...selectedItem, name: e.target.value })}
                  placeholder="请输入审阅名称"
                  style={{ flex: 1, marginRight: 12 }}
                />
              ) : (
                <div className="review-title">{selectedItem.name}</div>
              )}
              <Space>
                {!isEditing ? (
                  <>
                    <Button icon={<EditOutlined />} onClick={handleEdit}>
                      编辑
                    </Button>
                    <Button icon={<DownloadOutlined />} onClick={handleDownload}>
                      下载
                    </Button>
                  </>
                ) : (
                  <>
                    <Button 
                      type="primary" 
                      icon={<SaveOutlined />}
                      onClick={handleSave} 
                      loading={saveLoading}
                    >
                      保存
                    </Button>
                    <Button icon={<CloseOutlined />} onClick={() => setIsEditing(false)}>
                      取消
                    </Button>
                  </>
                )}
              </Space>
            </div>

            {/* 内容区 */}
            <div className="right-panel-content">
              {isEditing ? (
                <MdEditor
                  ref={editorRef}
                  value={editingContent}
                  onChange={handleEditorChange}
                  renderHTML={(text) => mdParser.render(text)}
                  style={{ height: '100%' }}
                  placeholder="请输入审阅内容（支持Markdown语法）"
                  config={{
                    view: { menu: true, md: true, html: true },
                    canView: { menu: true, md: true, html: true, fullScreen: true, hideMenu: true },
                  }}
                />
              ) : (
                <div className="readonly-container">
                  {/* 左侧目录 */}
                  {tocItems.length > 0 && (
                    <div className="toc-sidebar">
                      <div className="toc-title">目录</div>
                      <div className="toc-list">
                        {tocItems.map((item, index) => (
                          <div
                            key={index}
                            className={`toc-item toc-level-${item.level}`}
                            onClick={() => {
                              const element = document.getElementById(item.id);
                              if (element) {
                                element.scrollIntoView({ behavior: 'smooth', block: 'start' });
                              }
                            }}
                          >
                            {item.title}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                  
                  {/* 右侧内容 */}
                  <div className="markdown-preview" ref={previewRef}>
                    {editingContent ? (
                      <ReactMarkdown
                        children={editingContent}
                        remarkPlugins={[remarkGfm]}
                        components={{
                          h1: ({ node, children, ...props }) => {
                            const extractText = (ch) => {
                              if (typeof ch === 'string') return ch;
                              if (Array.isArray(ch)) return ch.map(extractText).join('');
                              if (ch?.props?.children) return extractText(ch.props.children);
                              return '';
                            };
                            const id = titleToId(extractText(children));
                            return <h1 id={id} {...props}>{children}</h1>;
                          },
                          h2: ({ node, children, ...props }) => {
                            const extractText = (ch) => {
                              if (typeof ch === 'string') return ch;
                              if (Array.isArray(ch)) return ch.map(extractText).join('');
                              if (ch?.props?.children) return extractText(ch.props.children);
                              return '';
                            };
                            const id = titleToId(extractText(children));
                            return <h2 id={id} {...props}>{children}</h2>;
                          },
                          h3: ({ node, children, ...props }) => {
                            const extractText = (ch) => {
                              if (typeof ch === 'string') return ch;
                              if (Array.isArray(ch)) return ch.map(extractText).join('');
                              if (ch?.props?.children) return extractText(ch.props.children);
                              return '';
                            };
                            const id = titleToId(extractText(children));
                            return <h3 id={id} {...props}>{children}</h3>;
                          },
                          h4: ({ node, children, ...props }) => {
                            const extractText = (ch) => {
                              if (typeof ch === 'string') return ch;
                              if (Array.isArray(ch)) return ch.map(extractText).join('');
                              if (ch?.props?.children) return extractText(ch.props.children);
                              return '';
                            };
                            const id = titleToId(extractText(children));
                            return <h4 id={id} {...props}>{children}</h4>;
                          },
                          h5: ({ node, children, ...props }) => {
                            const extractText = (ch) => {
                              if (typeof ch === 'string') return ch;
                              if (Array.isArray(ch)) return ch.map(extractText).join('');
                              if (ch?.props?.children) return extractText(ch.props.children);
                              return '';
                            };
                            const id = titleToId(extractText(children));
                            return <h5 id={id} {...props}>{children}</h5>;
                          },
                          h6: ({ node, children, ...props }) => {
                            const extractText = (ch) => {
                              if (typeof ch === 'string') return ch;
                              if (Array.isArray(ch)) return ch.map(extractText).join('');
                              if (ch?.props?.children) return extractText(ch.props.children);
                              return '';
                            };
                            const id = titleToId(extractText(children));
                            return <h6 id={id} {...props}>{children}</h6>;
                          },
                        }}
                      />
                    ) : (
                      <Empty description="暂无内容，点击'编辑'按钮开始编写" />
                    )}
                  </div>
                </div>
              )}
            </div>
          </>
        ) : (
          <div className="empty-state">
            <FileTextOutlined style={{ fontSize: 64, color: '#ccc', marginBottom: 16 }} />
            <p style={{ color: '#999' }}>{t('manualTest.selectReviewFromLeft')}</p>
          </div>
        )}
      </div>

      {/* 新建审阅Modal */}
      <Modal
        title={t('manualTest.createReview')}
        open={createModalVisible}
        onOk={handleCreateConfirm}
        onCancel={() => setCreateModalVisible(false)}
        confirmLoading={creating}
        okText={t('common.create')}
        cancelText={t('common.cancel')}
      >
        <Input
          placeholder={t('manualTest.enterReviewName')}
          value={newItemName}
          onChange={(e) => setNewItemName(e.target.value)}
          onPressEnter={handleCreateConfirm}
          maxLength={255}
        />
      </Modal>
      
      {/* 编辑名称Modal */}
      <Modal
        title={t('manualTest.editReviewName')}
        open={editNameModalVisible}
        onOk={handleEditNameConfirm}
        onCancel={() => setEditNameModalVisible(false)}
        okText={t('common.ok')}
        cancelText={t('common.cancel')}
      >
        <Input
          placeholder={t('manualTest.enterReviewName')}
          value={editingName}
          onChange={(e) => setEditingName(e.target.value)}
          onPressEnter={handleEditNameConfirm}
          maxLength={255}
        />
      </Modal>
    </div>
  );
};

export default CaseReviewTab;
