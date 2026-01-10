import React, { useState, useEffect } from 'react';
import { Modal, Descriptions, Button, message, Spin, Tag, Table } from 'antd';
import { CopyOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ReactMarkdown from 'react-markdown';
import { fetchPromptById } from '../../api/prompt';

const PromptViewer = ({ visible, prompt, onClose }) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [detailPrompt, setDetailPrompt] = useState(null);

  useEffect(() => {
    if (visible && prompt?.id) {
      loadPromptDetail();
    }
  }, [visible, prompt?.id]);

  const loadPromptDetail = async () => {
    setLoading(true);
    try {
      const data = await fetchPromptById(prompt.id);
      setDetailPrompt(data);
    } catch (error) {
      message.error('加载详情失败');
      onClose();
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = () => {
    const text = detailPrompt?.content || '';
    navigator.clipboard.writeText(text).then(() => {
      message.success(t('prompts.copySuccess'));
    }).catch(() => {
      message.error(t('prompts.copyFailed'));
    });
  };

  const argumentColumns = [
    {
      title: t('prompts.parameterName'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('prompts.parameterDescription'),
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: t('prompts.parameterRequired'),
      dataIndex: 'required',
      key: 'required',
      render: (required) => (
        <Tag color={required ? 'red' : 'default'}>
          {required ? t('prompts.parameterRequired') : '可选'}
        </Tag>
      ),
    },
  ];

  return (
    <Modal
      title={t('prompts.viewTitle')}
      open={visible}
      onCancel={onClose}
      width={800}
      bodyStyle={{ maxHeight: 'calc(100vh - 200px)', overflow: 'hidden', display: 'flex', flexDirection: 'column' }}
      modalRenderToBody={false}
      footer={[
        <Button key="copy" icon={<CopyOutlined />} onClick={handleCopy}>
          {t('prompts.copyFullText') || '复制全文'}
        </Button>,
        <Button key="close" onClick={onClose}>
          {t('prompts.close') || '关闭'}
        </Button>,
      ]}
      onOk={onClose}
      style={{
        pointerEvents: visible ? 'auto' : 'none',
      }}
    >
      {loading ? (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <Spin />
        </div>
      ) : (
        <div style={{ overflow: 'auto', flex: 1 }}>
          <Descriptions column={2} bordered>
            <Descriptions.Item label={t('prompts.name')}>
              {detailPrompt?.name}
            </Descriptions.Item>
            <Descriptions.Item label={t('prompts.version')}>
              {detailPrompt?.version}
            </Descriptions.Item>
            <Descriptions.Item label={t('prompts.scope')}>
              <Tag color={detailPrompt?.scope === 'system' ? 'green' : detailPrompt?.scope === 'project' ? 'blue' : 'orange'}>
                {detailPrompt?.scope === 'system' && t('prompts.systemPrompts')}
                {detailPrompt?.scope === 'project' && t('prompts.projectPrompts')}
                {detailPrompt?.scope === 'user' && t('prompts.userPrompts')}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label={t('prompts.createdAt')}>
              {detailPrompt?.created_at ? new Date(detailPrompt.created_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label={t('prompts.description')} span={2}>
              {detailPrompt?.description || '-'}
            </Descriptions.Item>
          </Descriptions>

          {detailPrompt?.arguments && detailPrompt.arguments.length > 0 && (
            <div style={{ marginTop: 16 }}>
              <h4>{t('prompts.parameters')}</h4>
              <Table
                columns={argumentColumns}
                dataSource={detailPrompt.arguments}
                pagination={false}
                size="small"
                rowKey="name"
              />
            </div>
          )}

          <div style={{ marginTop: 16 }}>
            <h4>{t('prompts.content')}</h4>
            <div
              style={{
                border: '1px solid #d9d9d9',
                borderRadius: 4,
                padding: '20px 24px',
                maxHeight: 400,
                overflow: 'auto',
                backgroundColor: '#fafafa',
                lineHeight: '1.8',
              }}
            >
              {detailPrompt?.content ? (
                <div style={{ 
                  fontSize: '14px',
                  color: '#262626',
                }}>
                  <ReactMarkdown>{detailPrompt.content}</ReactMarkdown>
                </div>
              ) : (
                <span style={{ color: '#999' }}>{t('prompts.noContent') || '无内容'}</span>
              )}
            </div>
          </div>
        </div>
      )}
    </Modal>
  );
};

export default PromptViewer;
