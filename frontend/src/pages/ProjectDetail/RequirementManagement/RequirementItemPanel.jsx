import React, { useState, useEffect, useCallback, useRef } from 'react';
import { List, Button, Input, Modal, message, Space, Popconfirm, Table, Tooltip, Empty, Spin } from 'antd';
import { PlusOutlined, DeleteOutlined, SaveOutlined, HistoryOutlined, DoubleLeftOutlined, DoubleRightOutlined, DownloadOutlined, EditOutlined, CloseOutlined, ImportOutlined, CopyOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import MdEditor from 'react-markdown-editor-lite';
import MarkdownIt from 'markdown-it';
import 'react-markdown-editor-lite/lib/index.css';
import {
  fetchRequirementItems,
  createRequirementItem,
  updateRequirementItem,
  deleteRequirementItem,
  exportRequirementItemsToZip,
} from '../../../api/requirementItem';
import {
  fetchChunks as fetchRequirementChunks,
  createChunk as createRequirementChunk,
  updateChunk as updateRequirementChunk,
  deleteChunk as deleteRequirementChunk,
} from '../../../api/requirementChunk';
import { getVersionList, downloadVersion, deleteVersion, updateVersionRemark } from '../../../api/requirement';
import ChunkTOC from './ChunkTOC';
import ChunkContent from './ChunkContent';
import './RequirementItemPanel.css';

/**
 * 需求条目管理面板 (T42)
 * 左侧列表 + 右侧Markdown编辑器
 */
const RequirementItemPanel = ({ projectId, projectName }) => {
  const { t } = useTranslation();
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(false);
  const [selectedItem, setSelectedItem] = useState(null);
  const [editingContent, setEditingContent] = useState('');
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [newItemName, setNewItemName] = useState('');
  const [isEditing, setIsEditing] = useState(false);
  const [leftPanelCollapsed, setLeftPanelCollapsed] = useState(false);
  
  // 版本一览相关状态
  const [versionModalVisible, setVersionModalVisible] = useState(false);
  const [versions, setVersions] = useState([]);
  const [versionsLoading, setVersionsLoading] = useState(false);
  const [editRemarkVisible, setEditRemarkVisible] = useState(false);
  const [editingVersion, setEditingVersion] = useState(null);
  const [editingRemark, setEditingRemark] = useState('');
  
  // 编辑名称相关状态
  const [editNameVisible, setEditNameVisible] = useState(false);
  const [editingItemName, setEditingItemName] = useState('');
  
  // TOC导航相关状态
  const [tocItems, setTocItems] = useState([]);
  
  // 编辑模式状态
  const [isEditMode, setIsEditMode] = useState(false);
  const mdParser = React.useRef(new MarkdownIt());
  
  // T54: Chunk相关状态
  const [chunks, setChunks] = useState([]);
  const [chunksLoading, setChunksLoading] = useState(false);
  const [activeChunkId, setActiveChunkId] = useState(null);
  const [editingChunkId, setEditingChunkId] = useState(null);
  const [middlePanelCollapsed, setMiddlePanelCollapsed] = useState(false);
  const contentRef = useRef(null);
  
  // 将标题转换为ID
  const titleToId = useCallback((title) => {
    return 'heading-' + title
      .toLowerCase()
      .replace(/[^\u4e00-\u9fa5a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }, []);
  
  // 生成目录
  const generateTOC = useCallback((markdown) => {
    if (!markdown) return [];
    // 【关键修复】处理 \r\n 和 \n 混合的情况
    const lines = markdown.split(/\r?\n/);
    const toc = [];
    
    lines.forEach(line => {
      // 【关键修复】移除行尾的 \r
      const cleanLine = line.replace(/\r$/, '');
      const match = cleanLine.match(/^(#{1,6})\s+(.+)$/);
      if (match) {
        const level = match[1].length;
        const title = match[2].trim();
        const id = titleToId(title);
        console.log('[generateTOC] 找到标题:', title, '级别:', level);
        toc.push({ level, title, id });
      }
    });
    
    console.log('[generateTOC] 总共生成', toc.length, '个目录项');
    return toc;
  }, [titleToId]);

  // T54: 加载Chunk列表 (必须在loadItems之前定义,避免循环依赖)
  const loadChunks = useCallback(async (itemId) => {
    if (!itemId || !projectId) return;
    setChunksLoading(true);
    try {
      const data = await fetchRequirementChunks(projectId, itemId);
      setChunks(Array.isArray(data) ? data : []);
      setActiveChunkId(null);
      setEditingChunkId(null);
    } catch (error) {
      console.error('加载Chunk失败:', error);
      setChunks([]);
    } finally {
      setChunksLoading(false);
    }
  }, [projectId]);
  
  // 加载需求条目列表
  const loadItems = useCallback(async () => {
    setLoading(true);
    try {
      const data = await fetchRequirementItems(projectId);
      const itemList = Array.isArray(data) ? data : [];
      setItems(itemList);
      
      // 如果当前选中的条目还存在,保持选中
      if (selectedItem) {
        const stillExists = itemList.find(item => item.id === selectedItem.id);
        if (stillExists) {
          setSelectedItem(stillExists);
          setEditingContent(stillExists.content || '');
          // 重新生成目录
          const toc = generateTOC(stillExists.content || '');
          setTocItems(toc);
        } else {
          setSelectedItem(null);
          setEditingContent('');
          setTocItems([]);
        }
      } else if (itemList.length > 0 && !selectedItem) {
        // 如果没有选中任何条目且列表不为空，默认选中第一条
        const firstItem = itemList[0];
        setSelectedItem(firstItem);
        setEditingContent(firstItem.content || '');
        // 为默认选中的条目生成目录
        const toc = generateTOC(firstItem.content || '');
        setTocItems(toc);
        // T54修复: 默认加载第一个条目的Chunks
        loadChunks(firstItem.id);
      }
    } catch (error) {
      message.error(t('requirement.loadFailed'));
    } finally {
      setLoading(false);
    }
  }, [projectId, selectedItem, t, generateTOC, loadChunks]);

  useEffect(() => {
    loadItems();
  }, [projectId]); // 仅在projectId变化时加载

  // 选中条目
  const handleSelectItem = (item) => {
    if (isEditMode) {
      Modal.confirm({
        title: t('requirement.unsavedChanges'),
        content: t('requirement.confirmLeave'),
        onOk: () => {
          setSelectedItem(item);
          setEditingContent(item.content || '');
          setIsEditMode(false);
          const toc = generateTOC(item.content || '');
          setTocItems(toc);
          loadChunks(item.id);
        },
      });
    } else {
      setSelectedItem(item);
      setEditingContent(item.content || '');
      const toc = generateTOC(item.content || '');
      setTocItems(toc);
      loadChunks(item.id);
    }
  };

  // T54: 点击目录项滚动到对应Chunk
  const handleChunkClick = useCallback((chunkId) => {
    setActiveChunkId(chunkId);
    const element = document.querySelector(`[data-chunk-id="${chunkId}"]`);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }, []);

  // T54: 添加新Chunk
  const handleAddChunk = useCallback(async () => {
    if (!selectedItem) return;
    try {
      const newChunk = await createRequirementChunk(projectId, selectedItem.id, '', '');
      setChunks(prev => [...prev, newChunk]);
      setEditingChunkId(newChunk.id);
      message.success(t('chunk.saveSuccess'));
    } catch (error) {
      message.error(error.message || '创建失败');
    }
  }, [projectId, selectedItem, t]);

  // T54: 编辑Chunk
  const handleEditChunk = useCallback((chunkId) => {
    setEditingChunkId(chunkId);
  }, []);

  // T54: 删除Chunk
  const handleDeleteChunk = useCallback(async (chunkId) => {
    try {
      await deleteRequirementChunk(chunkId);
      setChunks(prev => prev.filter(c => c.id !== chunkId));
      message.success(t('chunk.deleteSuccess'));
    } catch (error) {
      message.error(error.message || '删除失败');
    }
  }, [t]);

  // T54: 保存Chunk
  const handleSaveChunk = useCallback(async (chunkId, title, content) => {
    try {
      const updated = await updateRequirementChunk(chunkId, title, content);
      setChunks(prev => prev.map(c => c.id === chunkId ? { ...c, ...updated } : c));
      setEditingChunkId(null);
      message.success(t('chunk.saveSuccess'));
    } catch (error) {
      message.error(error.message || '保存失败');
    }
  }, [t]);

  // T54: 取消编辑Chunk
  const handleCancelChunkEdit = useCallback(() => {
    setEditingChunkId(null);
  }, []);

  // T54: 下载需求的所有Chunks为Markdown
  const handleDownloadChunks = useCallback((item) => {
    if (!chunks || chunks.length === 0) {
      message.warning('暂无段落可下载');
      return;
    }
    const content = chunks.map(c => c.content || '').join('\n\n---\n\n');
    const blob = new Blob([content], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${item?.name || 'requirement'}.md`;
    a.click();
    URL.revokeObjectURL(url);
  }, [chunks]);

  // 显示编辑名称对话框
  const handleShowEditName = (item) => {
    setSelectedItem(item);
    setEditingItemName(item.name);
    setEditNameVisible(true);
  };

  // 保存名称修改
  const handleSaveItemName = async () => {
    if (!editingItemName.trim()) {
      message.warning('名称不能为空');
      return;
    }

    try {
      await updateRequirementItem(selectedItem.id, editingItemName, selectedItem.content);
      message.success('名称修改成功');
      setEditNameVisible(false);
      loadItems();
    } catch (error) {
      message.error('名称修改失败');
    }
  };

  // 创建新条目
  const handleCreateItem = async () => {
    if (!newItemName.trim()) {
      message.warning(t('requirement.nameRequired'));
      return;
    }

    try {
      await createRequirementItem(projectId, newItemName, '');
      message.success(t('requirement.createSuccess'));
      setIsModalVisible(false);
      setNewItemName('');
      await loadItems();
    } catch (error) {
      message.error(t('requirement.createFailed'));
    }
  };

  // 保存内容
  const handleSave = async () => {
    if (!selectedItem) {
      message.warning(t('requirement.selectItemFirst'));
      return false;
    }

    if (!selectedItem.name.trim()) {
      message.warning('需求名称不能为空');
      return false;
    }

    try {
      await updateRequirementItem(selectedItem.id, selectedItem.name, editingContent);
      message.success(t('requirement.saveSuccess'));
      setIsEditMode(false);
      const toc = generateTOC(editingContent);
      setTocItems(toc);
      await loadItems();
      return true;
    } catch (error) {
      message.error(t('requirement.saveFailed'));
      return false;
    }
  };

  // 删除条目
  const handleDelete = async (itemId) => {
    try {
      await deleteRequirementItem(itemId);
      message.success(t('requirement.deleteSuccess'));
      if (selectedItem?.id === itemId) {
        setSelectedItem(null);
        setEditingContent('');
      }
      await loadItems();
    } catch (error) {
      message.error(t('requirement.deleteFailed'));
    }
  };

  // 版本保存（ZIP打包）
  const handleSaveVersion = async () => {
    console.log('=== [handleSaveVersion] 开始保存版本 ===');
    console.log('[handleSaveVersion] projectId:', projectId);
    console.log('[handleSaveVersion] items.length:', items.length);
    
    // 确保items是真正的数组
    const itemsArray = Array.isArray(items) ? items : Array.from(items || []);
    
    // 检查items是否为空或无效
    if (!itemsArray || itemsArray.length === 0) {
      console.warn('[handleSaveVersion] items为空，终止执行');
      message.warning('没有需求条目，无法保存版本');
      return;
    }
    
    // 显示正在保存的loading消息
    const hide = message.loading(`正在保存 ${itemsArray.length} 个需求条目...`, 0);
    console.log('[handleSaveVersion] 显示loading消息');
    
    try {
      const remark = '';
      console.log('[handleSaveVersion] 调用API - projectId:', projectId, 'remark:', remark);
      
      const result = await exportRequirementItemsToZip(projectId, remark);
      
      console.log('[handleSaveVersion] API调用成功:', result);
      
      // 关闭loading
      hide();
      
      const fileCount = result.file_list ? JSON.parse(result.file_list).length : 0;
      message.success(`版本保存成功：${result.filename}，包含 ${fileCount} 个文件`, 3);
      console.log('[handleSaveVersion] 保存成功');
    } catch (error) {
      console.error('[handleSaveVersion] 保存失败:', error);
      console.error('[handleSaveVersion] 错误详情:', error.response?.data);
      
      // 关闭loading
      hide();
      
      message.error(`版本保存失败: ${error.response?.data?.message || error.message || '未知错误'}`, 3);
    }
  };

  // 内容变更
  const handleContentChange = ({ text }) => {
    setEditingContent(text);
  };
  
  // 进入编辑模式
  const handleEdit = () => {
    setIsEditMode(true);
  };

  // 取消编辑
  const handleEditCancel = () => {
    if (selectedItem) {
      setEditingContent(selectedItem.content || '');
      const toc = generateTOC(selectedItem.content || '');
      setTocItems(toc);
    }
    setIsEditMode(false);
  };
  
  // 导入Markdown
  const handleImport = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.md';
    
    input.onchange = async (e) => {
      const file = e.target.files[0];
      if (!file) return;
      
      const maxSize = 5 * 1024 * 1024; // 5MB
      if (file.size > maxSize) {
        message.error('文件大小不能超过5MB');
        return;
      }
      
      const reader = new FileReader();
      reader.onload = async (event) => {
        const content = event.target.result;
        
        if (!selectedItem) {
          message.warning(t('requirement.selectItemFirst'));
          return;
        }
        
        try {
          console.log('[handleImport] 开始导入，文件大小:', file.size);
          console.log('[handleImport] 导入内容长度:', content.length);
          
          // 直接调用API保存
          await updateRequirementItem(selectedItem.id, selectedItem.name, content);
          message.success(t('requirement.importSuccess'));
          
          // 【关键】先更新UI状态
          setEditingContent(content);
          
          // 【关键】立即生成目录，在setTocItems前
          const tocArray = generateTOC(content);
          console.log('[handleImport] 生成的目录项数:', tocArray.length);
          console.log('[handleImport] 目录内容:', tocArray);
          
          // 【关键】直接设置TOC，不通过其他函数
          setTocItems(tocArray);
          setIsEditMode(false);
          
          // 【关键】延迟刷新列表，避免覆盖TOC状态
          setTimeout(() => {
            console.log('[handleImport] 延迟加载条目列表...');
            loadItems();
          }, 100);
          
        } catch (error) {
          console.error('[handleImport] 错误:', error);
          message.error(t('requirement.importFailed') + ': ' + error.message);
        }
      };
      reader.onerror = () => {
        message.error(t('requirement.importFailed'));
      };
      reader.readAsText(file, 'UTF-8');
    };
    
    input.click();
  };
  
  // 下载
  const handleDownload = () => {
    if (!editingContent) {
      message.warning('文档内容为空');
      return;
    }

    const blob = new Blob([editingContent], { type: 'text/markdown;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    
    const date = new Date();
    const dateStr = `${date.getFullYear()}${String(date.getMonth() + 1).padStart(2, '0')}${String(date.getDate()).padStart(2, '0')}`;
    link.download = `Project${projectId}_${selectedItem?.name || 'requirement'}_${dateStr}.md`;
    
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
    
    message.success('下载成功');
  };
  
  // TOC点击导航
  const handleTocClick = (id) => {
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  };
  
  // 提取React元素中的文本内容
  const extractText = (children) => {
    if (typeof children === 'string') return children;
    if (Array.isArray(children)) {
      return children.map(extractText).join('');
    }
    if (children?.props?.children) {
      return extractText(children.props.children);
    }
    return '';
  };

  // 版本一览
  const handleShowVersions = async () => {
    setVersionModalVisible(true);
    // 立即刷新版本列表
    await loadVersions();
  };

  const loadVersions = async () => {
    setVersionsLoading(true);
    try {
      console.log('[loadVersions] 开始加载版本列表, projectId:', projectId);
      // 临时方案：直接查询versions表的所有记录，前端过滤
      // TODO: 后端需要添加专门查询versions表的API
      const baseURL = process.env.REACT_APP_API_BASE_URL || '/api/v1';
      const url = `${baseURL}/versions?project_id=${projectId}&doc_type=`;
      console.log('[loadVersions] 查询URL:', url);
      
      const token = localStorage.getItem('auth_token');
      const response = await fetch(url, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }
      
      const result = await response.json();
      console.log('[loadVersions] 原始响应:', result);
      
      // 尝试多种可能的数据格式
      let data = [];
      if (Array.isArray(result)) {
        data = result;
      } else if (result.data && Array.isArray(result.data)) {
        data = result.data;
      } else if (result.data && result.data.data && Array.isArray(result.data.data)) {
        data = result.data.data;
      }
      
      // 过滤出item_type为requirement-batch的记录
      const filteredData = data.filter(v => v.item_type === 'requirement-batch');
      console.log('[loadVersions] 过滤后的版本列表:', filteredData);
      
      setVersions(filteredData);
    } catch (error) {
      console.error('[loadVersions] 加载版本列表失败:', error);
      message.error('加载版本列表失败');
    } finally {
      setVersionsLoading(false);
    }
  };

  const handleDownloadVersion = async (versionId) => {
    try {
      await downloadVersion(projectId, versionId);
      message.success('下载成功');
    } catch (error) {
      console.error('下载失败:', error);
      message.error('下载失败');
    }
  };

  const handleDeleteVersion = async (versionId) => {
    try {
      await deleteVersion(projectId, versionId);
      message.success('删除成功');
      loadVersions();
    } catch (error) {
      console.error('删除失败:', error);
      message.error('删除失败');
    }
  };

  const handleEditVersionRemark = (version) => {
    setEditingVersion(version);
    setEditingRemark(version.remark || '');
    setEditRemarkVisible(true);
  };

  const handleSaveRemark = async () => {
    try {
      await updateVersionRemark(projectId, editingVersion.id, editingRemark);
      message.success('备注更新成功');
      setEditRemarkVisible(false);
      loadVersions();
    } catch (error) {
      console.error('更新备注失败:', error);
      message.error('更新备注失败');
    }
  };

  const versionColumns = [
    {
      title: t('manualTest.versionListNo'),
      width: 80,
      render: (_, __, index) => index + 1,
    },
    {
      title: t('manualTest.versionListFile'),
      dataIndex: 'filename',
      render: (filename) => (
        <div style={{ wordBreak: 'break-all', whiteSpace: 'normal' }}>
          {filename}
        </div>
      ),
    },
    {
      title: t('manualTest.versionListRemark'),
      dataIndex: 'remark',
      width: 200,
      render: (remark) => (
        <div style={{ wordBreak: 'break-all', whiteSpace: 'normal' }}>
          {remark || '-'}
        </div>
      ),
    },
    {
      title: t('manualTest.versionListCreatedAt'),
      dataIndex: 'created_at',
      width: 120,
      render: (date) => {
        if (!date) return '-';
        const dateObj = new Date(date);
        return dateObj.toLocaleDateString('zh-CN', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit'
        });
      },
    },
    {
      title: t('manualTest.versionListActions'),
      width: 120,
      render: (_, record) => (
        <Space size={4}>
          <Button
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEditVersionRemark(record)}
            title={t('manualTest.editRemark')}
          />
          <Button
            size="small"
            icon={<DownloadOutlined />}
            onClick={() => handleDownloadVersion(record.id)}
            title={t('common.download')}
          />
          <Popconfirm
            title={t('requirement.confirmDelete')}
            onConfirm={() => handleDeleteVersion(record.id)}
            okText="确定"
            cancelText={t('common.cancel')}
          >
            <Button 
              size="small" 
              danger 
              icon={<DeleteOutlined />}
              title={t('common.delete')}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="requirement-item-panel">
      <div className={`panel-left ${leftPanelCollapsed ? 'collapsed' : ''}`}>
        {!leftPanelCollapsed && (
          <>
            <div className="panel-header">
              <Space direction="vertical" style={{ width: '100%' }} size="small">
                <Button
                  type="primary"
                  icon={<PlusOutlined />}
                  onClick={() => setIsModalVisible(true)}
                  block
                  size="small"
                >
                  {t('requirement.createItem')}
                </Button>
                <Button
                  icon={<SaveOutlined />}
                  onClick={handleSaveVersion}
                  disabled={items.length === 0}
                  block
                  size="small"
                >
                  {t('requirement.saveVersion')}
                </Button>
                <Button
                  icon={<HistoryOutlined />}
                  onClick={handleShowVersions}
                  block
                  size="small"
                >
                  {t('requirement.versionList')}
                </Button>
              </Space>
            </div>
            <div className="panel-list-title">
              <span>{t('requirement.itemList')}</span>
              <Button 
                type="text" 
                size="small" 
                icon={<DoubleLeftOutlined />}
                onClick={() => setLeftPanelCollapsed(true)}
                style={{ padding: '0 4px' }}
              />
            </div>
          </>
        )}
        {leftPanelCollapsed && (
          <div className="collapsed-trigger" onClick={() => setLeftPanelCollapsed(false)}>
            <DoubleRightOutlined />
          </div>
        )}
        
        {!leftPanelCollapsed && (
          <div className="requirement-list">
            {loading ? (
              <div style={{ textAlign: 'center', padding: '20px' }}>
                <Spin />
              </div>
            ) : items.length === 0 ? (
              <Empty description="暂无需求" style={{ marginTop: '20px' }} />
            ) : (
              items.map(item => {
                const date = new Date(item.updated_at);
                const dateStr = `${String(date.getFullYear()).slice(2)}/${String(date.getMonth() + 1).padStart(2, '0')}/${String(date.getDate()).padStart(2, '0')}`;
                
                return (
                  <div
                    key={item.id}
                    className={`requirement-list-item ${selectedItem?.id === item.id ? 'selected' : ''}`}
                    onClick={() => handleSelectItem(item)}
                  >
                    <div className="requirement-item-content">
                      <div className="requirement-item-name">{item.name}</div>
                      <div className="requirement-item-date">{dateStr}</div>
                    </div>
                    <div className="requirement-item-actions">
                      <Tooltip title="复制需求名">
                        <CopyOutlined
                          className="action-icon"
                          onClick={(e) => {
                            e.stopPropagation();
                            navigator.clipboard.writeText(item.name);
                            message.success('复制成功');
                          }}
                        />
                      </Tooltip>
                      <Tooltip title={t('common.download')}>
                        <DownloadOutlined
                          className="action-icon"
                          onClick={(e) => {
                            e.stopPropagation();
                            if (selectedItem?.id === item.id) {
                              handleDownloadChunks(item);
                            } else {
                              message.info('请先选中此需求');
                            }
                          }}
                        />
                      </Tooltip>
                      <Popconfirm
                        title={t('requirement.confirmDelete')}
                        onConfirm={(e) => {
                          e.stopPropagation();
                          handleDelete(item.id);
                        }}
                        okText="确定"
                        cancelText={t('common.cancel')}
                      >
                        <DeleteOutlined
                          className="action-icon"
                          onClick={(e) => e.stopPropagation()}
                        />
                      </Popconfirm>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        )}
      </div>

      <div className="panel-right">
        {selectedItem ? (
          <>
            {/* T54: 中栏 - Chunk目录 */}
            <div className={`panel-middle ${middlePanelCollapsed ? 'collapsed' : ''}`}>
              <ChunkTOC
                chunks={chunks}
                activeChunkId={activeChunkId}
                onChunkClick={handleChunkClick}
                onAddChunk={handleAddChunk}
                collapsed={middlePanelCollapsed}
                onCollapse={() => setMiddlePanelCollapsed(!middlePanelCollapsed)}
              />
            </div>
            {/* T54: 右栏 - Chunk内容 */}
            <div className="panel-content" ref={contentRef}>
              {chunksLoading ? (
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%' }}>
                  <Spin />
                </div>
              ) : chunks.length === 0 ? (
                <Empty
                  description={t('chunk.emptyChunks')}
                  style={{ marginTop: '100px' }}
                />
              ) : (
                <ChunkContent
                  chunks={chunks}
                  activeChunkId={activeChunkId}
                  editingChunkId={editingChunkId}
                  onEdit={handleEditChunk}
                  onDelete={handleDeleteChunk}
                  onSave={handleSaveChunk}
                  onCancel={handleCancelChunkEdit}
                />
              )}
            </div>
          </>
        ) : (
          <div className="empty-placeholder">
            {t('requirement.selectItemToEdit')}
          </div>
        )}
      </div>

      {/* 新建需求对话框 */}
      <Modal
        title={t('requirement.createItemTitle')}
        open={isModalVisible}
        onOk={handleCreateItem}
        onCancel={() => {
          setIsModalVisible(false);
          setNewItemName('');
        }}
        okText={t('common.confirm')}
        cancelText={t('common.cancel')}
      >
        <Input
          placeholder="请输入需求名称"
          value={newItemName}
          onChange={(e) => setNewItemName(e.target.value)}
          onPressEnter={handleCreateItem}
        />
      </Modal>

      {/* 版本一览对话框 */}
      <Modal
        title={t('requirement.versionListTitle')}
        open={versionModalVisible}
        onCancel={() => setVersionModalVisible(false)}
        footer={null}
        width={1000}
      >
        <Table
          columns={versionColumns}
          dataSource={versions}
          loading={versionsLoading}
          rowKey="id"
          pagination={{ pageSize: 10 }}
          scroll={{ x: 'max-content' }}
          size="small"
        />
      </Modal>

      {/* 编辑备注对话框 */}
      <Modal
        title={t('requirement.editRemarkTitle')}
        open={editRemarkVisible}
        onOk={handleSaveRemark}
        onCancel={() => setEditRemarkVisible(false)}
        okText={t('requirement.save')}
        cancelText={t('common.cancel')}
      >
        <Input.TextArea
          value={editingRemark}
          onChange={(e) => setEditingRemark(e.target.value)}
          rows={4}
          placeholder="请输入备注信息"
        />
      </Modal>

      {/* 编辑名称对话框 */}
      <Modal
        title={t('requirement.editItemTitle')}
        open={editNameVisible}
        onOk={handleSaveItemName}
        onCancel={() => setEditNameVisible(false)}
        okText={t('requirement.save')}
        cancelText={t('common.cancel')}
      >
        <Input
          value={editingItemName}
          onChange={(e) => setEditingItemName(e.target.value)}
          placeholder="请输入需求名称..."
        />
      </Modal>
    </div>
  );
};

export default RequirementItemPanel;
