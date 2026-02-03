import React from 'react';
import { PlusOutlined, DoubleRightOutlined } from '@ant-design/icons';
import { Button, Tooltip, Empty } from 'antd';
import { useTranslation } from 'react-i18next';
import './ChunkTOC.css';

/**
 * Chunk目录组件 (T54)
 * 显示需求/观点的Chunk列表，支持点击导航
 */
const ChunkTOC = ({
  chunks = [],
  activeChunkId,
  onChunkClick,
  onAddChunk,
  collapsed = false,
  onCollapse,
}) => {
  const { t } = useTranslation();

  // 折叠状态下只显示展开触发条
  if (collapsed) {
    return (
      <div className="chunk-toc-collapsed" onClick={onCollapse}>
        <Tooltip title={t('chunk.contents')} placement="right">
          <DoubleRightOutlined className="collapse-trigger" />
        </Tooltip>
      </div>
    );
  }

  return (
    <div className="chunk-toc">
      <div className="chunk-toc-header">
        <span className="chunk-toc-title">{t('chunk.contents')}</span>
        <Tooltip title={t('chunk.addChunk')}>
          <Button
            type="text"
            size="small"
            icon={<PlusOutlined />}
            onClick={onAddChunk}
            className="chunk-add-btn"
          />
        </Tooltip>
      </div>
      <div className="chunk-toc-list">
        {chunks.length === 0 ? (
          <Empty
            description={t('chunk.emptyChunks')}
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            className="chunk-empty"
          />
        ) : (
          chunks.map((chunk, index) => (
            <div
              key={chunk.id}
              className={`chunk-toc-item ${activeChunkId === chunk.id ? 'active' : ''}`}
              onClick={() => onChunkClick(chunk.id)}
            >
              <span className="chunk-index">{index + 1}.</span>
              <span className="chunk-title" title={chunk.title || t('chunk.untitled')}>
                {chunk.title || t('chunk.untitled')}
              </span>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default ChunkTOC;
