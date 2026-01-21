import React, { useState, useEffect } from 'react';
import { Button, Form, Input, Select, Space, message, Spin, Descriptions, Tag, Table, Checkbox, Popconfirm } from 'antd';
import { ArrowLeftOutlined, SaveOutlined, CopyOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { fetchPromptById, updatePrompt, createPrompt, deletePrompt } from '../../api/prompt';
import './index.css';

const { TextArea } = Input;

const PromptDetail = ({ prompt, mode, onBack, scope, onSuccess }) => {
  const { t } = useTranslation();
  const { user } = useSelector((state) => state.auth);
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [detailPrompt, setDetailPrompt] = useState(null);
  const [isEditing, setIsEditing] = useState(mode === 'create');

  useEffect(() => {
    if (mode === 'create') {
      setIsEditing(true);
      form.resetFields();
      setDetailPrompt(null);
    } else if (prompt?.id) {
      loadPromptDetail();
    }
  }, [prompt?.id, mode, form]);

  const loadPromptDetail = async () => {
    setLoading(true);
    try {
      const data = await fetchPromptById(prompt.id);
      // 解析 arguments JSON 字符串为数组
      const parsedArguments = data.arguments
        ? (typeof data.arguments === 'string' ? JSON.parse(data.arguments) : data.arguments)
        : [];
      setDetailPrompt({ ...data, arguments: parsedArguments });
      form.setFieldsValue({
        name: data.name,
        version: data.version,
        description: data.description,
        content: data.content,
        scope: data.scope,
        arguments: parsedArguments,
      });
    } catch (error) {
      message.error('加载详情失败');
      onBack();
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

  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      // 处理 arguments 数据，确保 required 是布尔值
      const processedArguments = (values.arguments || []).map(arg => ({
        name: arg.name,
        description: arg.description || '',
        required: Boolean(arg.required), // Checkbox 返回布尔值，确保转换
      }));

      if (mode === 'create') {
        // 创建时发送完整数据，不包括project_id
        const createData = {
          name: values.name,
          version: values.version,
          description: values.description || '',
          content: values.content,
          scope: values.scope, // 发送scope字段
          arguments: JSON.stringify(processedArguments), // 转换为JSON字符串
        };

        // 详细调试日志
        console.log('=== [PromptDetail] 创建提示词 ===');
        console.log('[调试] scope:', values.scope);
        console.log('[调试] 最终发送的数据:', JSON.stringify(createData, null, 2));
        console.log('[调试] name:', createData.name, '类型:', typeof createData.name, '长度:', createData.name?.length);
        console.log('[调试] version:', createData.version, '类型:', typeof createData.version);
        console.log('[调试] description:', createData.description, '类型:', typeof createData.description, '长度:', createData.description?.length);
        console.log('[调试] content:', createData.content, '类型:', typeof createData.content, '长度:', createData.content?.length);
        console.log('[调试] arguments:', createData.arguments, '类型:', typeof createData.arguments, '长度:', createData.arguments?.length);
        console.log('[调试] scope:', createData.scope, '类型:', typeof createData.scope);
        console.log('[调试] Object.keys:', Object.keys(createData));

        await createPrompt(createData);
        message.success(t('prompts.createSuccess'));
        onSuccess();
      } else {
        await updatePrompt(prompt.id, {
          version: values.version,
          description: values.description || '',
          content: values.content,
          arguments: JSON.stringify(processedArguments), // 转换为JSON字符串
        });
        message.success(t('prompts.updateSuccess'));
        setIsEditing(false);
        onSuccess();
      }
    } catch (error) {
      if (!error.errorFields) {
        console.error('=== [PromptDetail] API 调用失败 ===');
        console.error('[错误] 完整 error 对象:', error);
        console.error('[错误] error.response:', error.response);
        console.error('[错误] error.response?.status:', error.response?.status);
        console.error('[错误] error.response?.data:', error.response?.data);
        console.error('[错误] error.response?.data?.code:', error.response?.data?.code);
        console.error('[错误] error.response?.data?.message:', error.response?.data?.message);
        console.error('[错误] error.response?.data?.data:', error.response?.data?.data);
        console.error('[错误] error.message:', error.message);

        const errorMsg = error.response?.data?.message || error.message;
        message.error(errorMsg || (mode === 'create' ? t('prompts.createFailed') : t('prompts.updateFailed')));
      }
    } finally {
      setLoading(false);
    }
  };

  const canEdit = () => {
    if (mode === 'create') return true;
    if (detailPrompt?.scope === 'system') return false;
    if (detailPrompt?.scope === 'user') return detailPrompt?.user_id === user?.id;
    if (detailPrompt?.scope === 'project') return user?.role === 'system_admin';
    return false;
  };

  const handleDelete = async () => {
    try {
      setLoading(true);
      await deletePrompt(detailPrompt.id);
      message.success(t('prompts.deleteSuccess'));
      onBack();
      onSuccess();
    } catch (error) {
      const errorMsg = error.response?.data?.message || error.message;
      message.error(errorMsg || t('prompts.deleteFailed'));
    } finally {
      setLoading(false);
    }
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

  if (loading && mode !== 'create') {
    return (
      <div className="prompt-detail">
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', flex: 1 }}>
          <Spin />
        </div>
      </div>
    );
  }

  // 编辑/创建表单
  if (isEditing) {
    return (
      <div className="prompt-detail">
        <div className="prompt-detail-header">
          <Button type="text" icon={<ArrowLeftOutlined />} onClick={onBack} />
          <h2 style={{ margin: 0, flex: 1 }}>
            {mode === 'create' ? t('prompts.createTitle') : detailPrompt?.name || '编辑'}
          </h2>
        </div>

        <Form
          form={form}
          layout="vertical"
          initialValues={{ version: '1.0', scope: scope, arguments: [] }}
          style={{ flex: 1, overflow: 'auto', paddingRight: 8 }}
        >
          <Form.Item
            name="name"
            label={t('prompts.name')}
            rules={[
              { required: true, message: t('prompts.nameRequired') },
              { min: 3, max: 50, message: '名称长度3-50字符' },
              { pattern: /^[a-zA-Z_][a-zA-Z0-9_]*$/, message: '仅支持英文字母、数字和下划线，且首字符不能是数字' },
            ]}
          >
            <Input placeholder="如：my_prompt" disabled={mode !== 'create'} />
          </Form.Item>

          <Form.Item
            name="version"
            label={t('prompts.version')}
            rules={[{ required: true, message: '版本必填' }]}
            initialValue={mode === 'create' ? '1.0' : undefined}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="description"
            label={t('prompts.description')}
          >
            <TextArea rows={1} placeholder="提示词描述" />
          </Form.Item>

          <Form.Item
            name="content"
            label={t('prompts.content')}
            rules={[{ required: true, message: '内容必填' }]}
          >
            <TextArea rows={10} placeholder="提示词内容" />
          </Form.Item>

          <Form.Item
            name="scope"
            label={t('prompts.scope')}
          >
            <Input
              disabled
              value={scope}
              style={{ backgroundColor: '#f5f5f5' }}
            />
          </Form.Item>

          <div style={{ marginBottom: 16 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
              <h4 style={{ margin: 0 }}>{t('prompts.parameters')}</h4>
              {isEditing && (
                <Button type="primary" size="small" onClick={() => {
                  const currentArgs = form.getFieldValue('arguments') || [];
                  const newArg = {
                    name: `param_${currentArgs.length + 1}`,
                    description: '',
                    required: false
                  };
                  form.setFieldValue('arguments', [...currentArgs, newArg]);
                }}>
                  + {t('prompts.addParameter')}
                </Button>
              )}
            </div>
            {isEditing ? (
              <Form.List name="arguments">
                {(fields, { remove }) => {
                  const args = form.getFieldValue('arguments') || [];
                  return args.length > 0 ? (
                    <div style={{ border: '1px solid #d9d9d9', borderRadius: 4, backgroundColor: '#fff' }}>
                      {fields.map((field, index) => (
                        <div
                          key={field.key}
                          style={{
                            padding: '12px',
                            borderBottom: index < fields.length - 1 ? '1px solid #f0f0f0' : 'none',
                            display: 'flex',
                            gap: '12px',
                            alignItems: 'flex-start',
                          }}
                        >
                          <div style={{ flex: 1 }}>
                            <Form.Item
                              {...field}
                              name={[field.name, 'name']}
                              label={t('prompts.parameterName')}
                              rules={[{ required: true, message: t('prompts.parameterNameRequired') }]}
                              style={{ marginBottom: 8 }}
                            >
                              <Input placeholder={t('prompts.parameterName')} />
                            </Form.Item>
                            <Form.Item
                              {...field}
                              name={[field.name, 'description']}
                              label={t('prompts.parameterDescription')}
                              style={{ marginBottom: 0 }}
                            >
                              <Input placeholder={t('prompts.parameterDescription')} />
                            </Form.Item>
                          </div>
                          <div style={{ width: 120, marginTop: 4 }}>
                            <Form.Item
                              {...field}
                              name={[field.name, 'required']}
                              label={t('prompts.parameterRequired')}
                              valuePropName="checked"
                              style={{ marginBottom: 0 }}
                            >
                              <Checkbox />
                            </Form.Item>
                          </div>
                          <Button
                            type="text"
                            danger
                            size="small"
                            onClick={() => remove(field.name)}
                            style={{ marginTop: 4 }}
                          >
                            {t('prompts.removeParameter')}
                          </Button>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div style={{ color: '#999', padding: '16px', textAlign: 'center', border: '1px dashed #d9d9d9', borderRadius: 4 }}>
                      {t('prompts.parameterNameEmpty')}
                    </div>
                  );
                }}
              </Form.List>
            ) : (
              <Table
                columns={argumentColumns}
                dataSource={detailPrompt?.arguments || []}
                pagination={false}
                size="small"
                rowKey="name"
                style={{ backgroundColor: '#fff' }}
              />
            )}
          </div>
        </Form>

        <div style={{ display: 'flex', gap: 8, paddingTop: 16, borderTop: '1px solid #e8e8e8' }}>
          <Button onClick={onBack}>
            {t('common.cancel')}
          </Button>
          <Button type="primary" icon={<SaveOutlined />} onClick={handleSave} loading={loading}>
            {t('common.save')}
          </Button>
        </div>
      </div>
    );
  }

  // 查看模式
  return (
    <div className="prompt-detail">
      <div className="prompt-detail-header">
        <Button type="text" icon={<ArrowLeftOutlined />} onClick={onBack} />
        <h2 style={{ margin: 0, flex: 1 }}>{detailPrompt?.name}</h2>
        <Space>
          {canEdit() && (
            <>
              <Button
                type="primary"
                icon={<EditOutlined />}
                onClick={() => setIsEditing(true)}
              >
                {t('common.edit')}
              </Button>
              <Popconfirm
                title={t('prompts.deleteConfirm')}
                onConfirm={handleDelete}
              >
                <Button
                  type="primary"
                  danger
                  icon={<DeleteOutlined />}
                  loading={loading}
                >
                  {t('common.delete')}
                </Button>
              </Popconfirm>
            </>
          )}
        </Space>
      </div>

      <div className="prompt-detail-content">
        <Descriptions column={1} bordered style={{ marginBottom: 16, backgroundColor: '#fff' }}>
          <Descriptions.Item label={t('prompts.version')}>
            {detailPrompt?.version}
          </Descriptions.Item>
          <Descriptions.Item label={t('prompts.scope')}>
            <Tag
              color={
                detailPrompt?.scope === 'system'
                  ? 'green'
                  : detailPrompt?.scope === 'project'
                    ? 'blue'
                    : 'orange'
              }
            >
              {detailPrompt?.scope === 'system' && t('prompts.systemPrompts')}
              {detailPrompt?.scope === 'project' && t('prompts.projectPrompts')}
              {detailPrompt?.scope === 'user' && t('prompts.userPrompts')}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label={t('prompts.createdAt')}>
            {detailPrompt?.created_at
              ? new Date(detailPrompt.created_at).toLocaleString()
              : '-'}
          </Descriptions.Item>
          <Descriptions.Item label={t('prompts.description')}>
            {detailPrompt?.description || '-'}
          </Descriptions.Item>
        </Descriptions>

        <div style={{ marginBottom: 16 }}>
          <h4>{t('prompts.parameters')}</h4>
          <Table
            columns={argumentColumns}
            dataSource={detailPrompt?.arguments || []}
            pagination={false}
            size="small"
            rowKey="name"
            style={{ backgroundColor: '#fff' }}
          />
        </div>

        <div>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
            <h4 style={{ margin: 0 }}>{t('prompts.content')}</h4>
            <Button
              type="default"
              size="small"
              icon={<CopyOutlined />}
              onClick={handleCopy}
            >
              {t('prompts.copyContent')}
            </Button>
          </div>
          <div
            style={{
              border: '1px solid #d9d9d9',
              borderRadius: 4,
              padding: 12,
              backgroundColor: '#fafafa',
              fontSize: 13,
              lineHeight: 1.6,
            }}
          >
            {detailPrompt?.content ? (
              <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                skipHtml={false}
                components={{
                  // 支持表格样式
                  table: ({ node, ...props }) => (
                    <table style={{
                      borderCollapse: 'collapse',
                      width: '100%',
                      marginTop: '12px',
                      marginBottom: '16px',
                      backgroundColor: '#fff',
                      border: '1px solid #d9d9d9'
                    }} {...props} />
                  ),
                  thead: ({ node, ...props }) => (
                    <thead style={{
                      backgroundColor: '#fafafa',
                      borderBottom: '2px solid #d9d9d9'
                    }} {...props} />
                  ),
                  tbody: ({ node, ...props }) => (
                    <tbody {...props} />
                  ),
                  tr: ({ node, ...props }) => (
                    <tr style={{
                      borderBottom: '1px solid #e8e8e8'
                    }} {...props} />
                  ),
                  th: ({ node, ...props }) => (
                    <th style={{
                      border: '1px solid #d9d9d9',
                      padding: '10px 12px',
                      textAlign: 'left',
                      fontWeight: 600,
                      fontSize: '14px',
                      color: '#333'
                    }} {...props} />
                  ),
                  td: ({ node, ...props }) => (
                    <td style={{
                      border: '1px solid #e8e8e8',
                      padding: '10px 12px',
                      fontSize: '13px'
                    }} {...props} />
                  ),
                  // 只支持 ``` 包裹的代码块
                  code: ({ node, inline, className, children, ...props }) => {
                    // 只有带 className (language-*) 的才是代码块
                    const match = /language-(\w+)/.exec(className || '');

                    if (!inline && match) {
                      const lang = match[1];

                      // 如果是 SVG 代码块，直接渲染 SVG
                      if (lang === 'svg') {
                        const svgCode = String(children).replace(/\n$/, '');
                        return (
                          <div
                            style={{
                              border: '1px solid #e8e8e8',
                              borderRadius: '4px',
                              padding: '16px',
                              marginBottom: '16px',
                              backgroundColor: '#fff',
                              display: 'flex',
                              justifyContent: 'center',
                              alignItems: 'center'
                            }}
                            dangerouslySetInnerHTML={{ __html: svgCode }}
                          />
                        );
                      }

                      // 其他代码块正常显示
                      return (
                        <pre style={{
                          backgroundColor: '#f5f5f5',
                          padding: '12px',
                          borderRadius: '4px',
                          overflow: 'auto',
                          marginBottom: '16px',
                          border: '1px solid #e8e8e8'
                        }}>
                          <code className={className} style={{
                            fontFamily: 'Consolas, Monaco, "Courier New", monospace',
                            fontSize: '13px'
                          }} {...props}>{children}</code>
                        </pre>
                      );
                    }

                    // 行内代码
                    return (
                      <code style={{
                        backgroundColor: '#f5f5f5',
                        padding: '2px 6px',
                        borderRadius: '3px',
                        fontFamily: 'Consolas, Monaco, "Courier New", monospace',
                        fontSize: '0.9em',
                        color: '#c7254e',
                        border: '1px solid #e8e8e8'
                      }} {...props}>{children}</code>
                    );
                  },
                  // 禁用缩进代码块，保留正常段落
                  pre: ({ node, children, ...props }) => {
                    // 检查是否是代码块的 pre 标签
                    const codeChild = React.Children.toArray(children).find(
                      child => child?.props?.node?.tagName === 'code'
                    );

                    // 如果包含 code 标签且有 className，说明是 ``` 代码块，保留
                    if (codeChild?.props?.className) {
                      return <pre {...props}>{children}</pre>;
                    }

                    // 否则作为普通文本处理（禁用缩进代码块）
                    return <div {...props}>{children}</div>;
                  }
                }}
              >
                {detailPrompt.content}
              </ReactMarkdown>
            ) : (
              <span style={{ color: '#999' }}>无内容</span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default PromptDetail;
