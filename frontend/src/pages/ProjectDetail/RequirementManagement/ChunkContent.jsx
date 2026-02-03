import React from 'react';
import { EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { Button, Popconfirm, Empty } from 'antd';
import { useTranslation } from 'react-i18next';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import ChunkEditor from './ChunkEditor';
import './ChunkContent.css';

/**
 * Chunk内容展示组件 (T54)
 * 拼接显示所有Chunk，支持单Chunk编辑
 */
const ChunkContent = ({
  chunks = [],
  editingChunkId,
  activeChunkId,
  onEdit,
  onDelete,
  onSave,
  onCancel,
}) => {
  const { t } = useTranslation();

  if (chunks.length === 0) {
    return (
      <div className="chunk-content-empty">
        <Empty
          description={t('chunk.emptyChunks')}
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        />
      </div>
    );
  }

  return (
    <div className="chunk-content">
      {chunks.map((chunk, index) => (
        <div
          key={chunk.id}
          className={`chunk-item ${activeChunkId === chunk.id ? 'chunk-item-active' : ''}`}
          data-chunk-id={chunk.id}
        >
          {/* 控件栏 - 常驻显示 */}
          <div className="chunk-controls">
            <Button
              type="text"
              size="small"
              icon={<EditOutlined />}
              onClick={() => onEdit(chunk.id)}
              className="chunk-control-btn"
              title={t('chunk.editChunk')}
            />
            <Popconfirm
              title={t('chunk.confirmDelete')}
              onConfirm={() => onDelete(chunk.id)}
              okText={t('common.confirm')}
              cancelText={t('common.cancel')}
            >
              <Button
                type="text"
                size="small"
                icon={<DeleteOutlined />}
                className="chunk-control-btn"
                title={t('chunk.deleteChunk')}
              />
            </Popconfirm>
          </div>

          {/* 内容区域 */}
          {editingChunkId === chunk.id ? (
            <ChunkEditor
              initialTitle={chunk.title}
              initialContent={chunk.content}
              onSave={(title, content) => onSave(chunk.id, title, content)}
              onCancel={onCancel}
            />
          ) : (
            <div className="chunk-markdown">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {chunk.content || ''}
              </ReactMarkdown>
            </div>
          )}

          {/* 分隔线 */}
          {index < chunks.length - 1 && (
            <div className="chunk-divider" />
          )}
        </div>
      ))}
    </div>
  );
};

export default ChunkContent;
