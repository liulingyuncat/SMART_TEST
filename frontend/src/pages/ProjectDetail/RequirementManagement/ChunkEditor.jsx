import React, { useState } from 'react';
import { Button, Input, Space } from 'antd';
import { SaveOutlined, CloseOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import MdEditor from 'react-markdown-editor-lite';
import MarkdownIt from 'markdown-it';
import 'react-markdown-editor-lite/lib/index.css';
import './ChunkEditor.css';

/**
 * Chunk编辑器组件 (T54)
 * 封装MdEditor，提供标题和内容编辑
 */
const ChunkEditor = ({
  initialTitle = '',
  initialContent = '',
  onSave,
  onCancel,
}) => {
  const { t } = useTranslation();
  const [title, setTitle] = useState(initialTitle);
  const [content, setContent] = useState(initialContent);
  const mdParser = React.useRef(new MarkdownIt());

  const handleEditorChange = ({ text }) => {
    setContent(text);
  };

  const handleSave = () => {
    onSave(title, content);
  };

  return (
    <div className="chunk-editor">
      <div className="chunk-editor-header">
        <Input
          placeholder={t('chunk.titlePlaceholder') || '段落标题（可选）'}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="chunk-title-input"
        />
        <Space>
          <Button
            type="primary"
            size="small"
            icon={<SaveOutlined />}
            onClick={handleSave}
          >
            {t('common.save')}
          </Button>
          <Button
            size="small"
            icon={<CloseOutlined />}
            onClick={onCancel}
          >
            {t('common.cancel')}
          </Button>
        </Space>
      </div>
      <div className="chunk-editor-body">
        <MdEditor
          value={content}
          style={{ height: '300px' }}
          renderHTML={(text) => mdParser.current.render(text)}
          onChange={handleEditorChange}
          config={{
            view: {
              menu: true,
              md: true,
              html: true,
            },
          }}
        />
      </div>
    </div>
  );
};

export default ChunkEditor;
