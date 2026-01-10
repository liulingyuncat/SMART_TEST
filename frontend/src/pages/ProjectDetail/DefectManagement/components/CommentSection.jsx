import React, { useState, useEffect } from 'react';
import { List, Button, Typography, Space, Popconfirm, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { fetchDefectComments, createDefectComment, updateDefectComment, deleteDefectComment } from '../../../../api/defect';
import CommentModal from './CommentModal';

const { Title, Text } = Typography;

const CommentSection = ({ projectId, defectId, currentUserId, compact = false }) => {
  const { t } = useTranslation();
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [modalMode, setModalMode] = useState('add');
  const [editingComment, setEditingComment] = useState(null);

  // 调试: 检查i18n是否正常工作
  useEffect(() => {
    console.log('[DEBUG] CommentSection i18n test:', {
      comments: t('defect.comments'),
      addComment: t('defect.addComment'),
      noComments: t('defect.noComments'),
    });
  }, [t]);

  // 定义fallback翻译
  const labels = {
    comments: t('defect.comments', '说明'),
    addComment: t('defect.addComment', '新增说明'),
    noComments: t('defect.noComments', '暂无说明'),
    commentDeleteConfirm: t('defect.commentDeleteConfirm', '确定要删除此说明吗？'),
  };

  useEffect(() => {
    if (projectId && defectId) {
      loadComments();
    }
  }, [projectId, defectId]);

  const loadComments = async () => {
    try {
      setLoading(true);
      const data = await fetchDefectComments(projectId, defectId);
      setComments(data);
    } catch (error) {
      console.error('Failed to load comments:', error);
      message.error(t('message.loadFailed'));
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setModalMode('add');
    setEditingComment(null);
    setModalVisible(true);
  };

  const handleEdit = (comment) => {
    setModalMode('edit');
    setEditingComment(comment);
    setModalVisible(true);
  };

  const handleDelete = async (commentId) => {
    try {
      await deleteDefectComment(projectId, defectId, commentId);
      message.success(t('defect.commentDeleteSuccess'));
      loadComments();
    } catch (error) {
      console.error('Failed to delete comment:', error);
      if (error.response?.status === 403) {
        message.error(t('message.forbidden'));
      } else {
        message.error(t('message.deleteFailed'));
      }
    }
  };

  const handleSave = async (content) => {
    console.log('[DEBUG] CommentSection.handleSave called:', { mode: modalMode, content, projectId, defectId });
    try {
      if (modalMode === 'add') {
        console.log('[DEBUG] Creating new comment...');
        const result = await createDefectComment(projectId, defectId, content);
        console.log('[DEBUG] Comment created:', result);
        message.success(t('defect.commentCreateSuccess', '说明创建成功'));
      } else {
        console.log('[DEBUG] Updating comment:', editingComment.id);
        const result = await updateDefectComment(projectId, defectId, editingComment.id, content);
        console.log('[DEBUG] Comment updated:', result);
        message.success(t('defect.commentUpdateSuccess', '说明更新成功'));
      }
      setModalVisible(false);
      console.log('[DEBUG] Reloading comments...');
      await loadComments();
      console.log('[DEBUG] Comments reloaded');
    } catch (error) {
      console.error('[ERROR] Failed to save comment:', error);
      console.error('[ERROR] Error details:', error.response?.data);
      if (error.response?.status === 403) {
        message.error(t('message.forbidden', '没有权限'));
      } else {
        message.error(t('message.saveFailed', '保存失败'));
      }
      throw error;
    }
  };

  const formatDateTime = (dateStr) => {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div className="comment-section" style={{ marginTop: compact ? 0 : 24 }}>
      <div className="comment-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: compact ? 8 : 16 }}>
        <span style={{ fontSize: compact ? 13 : 16, fontWeight: 500, color: '#303133' }}>{labels.comments}</span>
        <Button size={compact ? 'small' : 'middle'} type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {labels.addComment}
        </Button>
      </div>

      <List
        size={compact ? 'small' : 'default'}
        loading={loading}
        dataSource={comments}
        locale={{ emptyText: labels.noComments }}
        renderItem={(comment) => (
          <List.Item
            style={{ padding: compact ? '6px 0' : '12px 0' }}
            actions={[
              comment.created_by === currentUserId && (
                <Button type="link" size="small" icon={<EditOutlined />} style={{ padding: 0, height: 'auto', fontSize: 12 }} onClick={() => handleEdit(comment)}>
                  {t('common.edit', '编辑')}
                </Button>
              ),
              comment.created_by === currentUserId && (
                <Popconfirm
                  title={labels.commentDeleteConfirm}
                  onConfirm={() => handleDelete(comment.id)}
                  okText={t('common.confirm', '确定')}
                  cancelText={t('common.cancel', '取消')}
                >
                  <Button type="link" size="small" danger icon={<DeleteOutlined />} style={{ padding: 0, height: 'auto', fontSize: 12 }}>
                    {t('common.delete', '删除')}
                  </Button>
                </Popconfirm>
              ),
            ].filter(Boolean)}
          >
            <List.Item.Meta
              title={
                <Space size="small">
                  <Text strong style={{ fontSize: compact ? 12 : 14 }}>{comment.updated_by_user?.username || comment.created_by_user?.username}</Text>
                  <Text type="secondary" style={{ fontSize: compact ? 11 : 12 }}>
                    {formatDateTime(comment.updated_at || comment.created_at)}
                  </Text>
                </Space>
              }
              description={<pre style={{ whiteSpace: 'pre-wrap', fontFamily: 'inherit', margin: 0, fontSize: compact ? 12 : 14, color: '#595959' }}>{comment.content}</pre>}
            />
          </List.Item>
        )}
      />

      <CommentModal
        visible={modalVisible}
        mode={modalMode}
        initialContent={editingComment?.content}
        onOk={handleSave}
        onCancel={() => setModalVisible(false)}
      />
    </div>
  );
};

export default CommentSection;
