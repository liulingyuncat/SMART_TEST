import React, { useState, useEffect } from 'react';
import { Drawer, Form, Input, Select, Button, Space, message, Divider } from 'antd';
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { createPrompt, updatePrompt } from '../../api/prompt';

const { TextArea } = Input;
const { Option } = Select;

const PromptForm = ({ visible, mode, prompt, onSuccess, onCancel, projectId }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (visible) {
      if (mode === 'edit' && prompt) {
        form.setFieldsValue({
          name: prompt.name,
          version: prompt.version,
          description: prompt.description,
          content: prompt.content,
          scope: prompt.scope,
          arguments: prompt.arguments || [],
        });
      } else {
        form.resetFields();
      }
    }
  }, [visible, mode, prompt, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      if (mode === 'create') {
        await createPrompt({
          project_id: parseInt(projectId),
          name: values.name,
          version: values.version,
          description: values.description || '',
          content: values.content,
          scope: values.scope,
          arguments: values.arguments || [],
        });
        message.success(t('prompts.createSuccess'));
      } else {
        await updatePrompt(prompt.id, {
          version: values.version,
          description: values.description || '',
          content: values.content,
          arguments: values.arguments || [],
        });
        message.success(t('prompts.updateSuccess'));
      }

      onSuccess();
    } catch (error) {
      if (error.errorFields) {
        // 表单验证错误
        return;
      }
      const errorMsg = error.response?.data?.message || error.message;
      message.error(errorMsg || (mode === 'create' ? t('prompts.createFailed') : t('prompts.updateFailed')));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Drawer
      title={mode === 'create' ? t('prompts.createTitle') : t('prompts.editTitle')}
      open={visible}
      onClose={onCancel}
      width={720}
      footer={
        <Space style={{ float: 'right' }}>
          <Button onClick={onCancel}>{t('common.cancel')}</Button>
          <Button type="primary" onClick={handleSubmit} loading={loading}>
            {t('common.save')}
          </Button>
        </Space>
      }
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{
          scope: 'user',
          arguments: [],
        }}
      >
        <Form.Item
          name="name"
          label={t('prompts.name')}
          rules={[
            { required: true, message: t('prompts.nameRequired') },
            { min: 3, max: 50, message: '名称长度3-50字符' },
            { pattern: /^[a-zA-Z_]+$/, message: '仅支持英文字母和下划线' },
          ]}
        >
          <Input
            placeholder={t('prompts.namePlaceholder')}
            disabled={mode === 'edit'}
          />
        </Form.Item>

        <Form.Item
          name="version"
          label={t('prompts.version')}
          rules={[
            { required: true, message: t('prompts.versionRequired') },
            { pattern: /^\d+\.\d+$/, message: '版本格式如: 1.0' },
          ]}
        >
          <Input placeholder={t('prompts.versionPlaceholder')} />
        </Form.Item>

        <Form.Item
          name="description"
          label={t('prompts.description')}
          rules={[{ max: 200, message: '描述最多200字符' }]}
        >
          <TextArea
            rows={2}
            placeholder={t('prompts.descriptionPlaceholder')}
          />
        </Form.Item>

        <Form.Item
          name="scope"
          label={t('prompts.scope')}
          rules={[{ required: true, message: t('prompts.scopeRequired') }]}
        >
          <Select placeholder={t('prompts.scopeRequired')} disabled={mode === 'edit'}>
            <Option value="project">{t('prompts.scopeProject')}</Option>
            <Option value="user">{t('prompts.scopeUser')}</Option>
          </Select>
        </Form.Item>

        <Divider />

        <Form.List name="arguments">
          {(fields, { add, remove }) => (
            <>
              <div style={{ marginBottom: 8 }}>
                <strong>{t('prompts.parameters')}</strong>
                <Button
                  type="dashed"
                  onClick={() => add()}
                  icon={<PlusOutlined />}
                  style={{ marginLeft: 16 }}
                  size="small"
                >
                  {t('prompts.addParameter')}
                </Button>
              </div>
              {fields.map(({ key, name, ...restField }) => (
                <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                  <Form.Item
                    {...restField}
                    name={[name, 'name']}
                    rules={[{ required: true, message: '请输入参数名' }]}
                    style={{ marginBottom: 0 }}
                  >
                    <Input placeholder={t('prompts.parameterName')} style={{ width: 150 }} />
                  </Form.Item>
                  <Form.Item
                    {...restField}
                    name={[name, 'description']}
                    style={{ marginBottom: 0 }}
                  >
                    <Input placeholder={t('prompts.parameterDescription')} style={{ width: 250 }} />
                  </Form.Item>
                  <Form.Item
                    {...restField}
                    name={[name, 'required']}
                    valuePropName="checked"
                    style={{ marginBottom: 0 }}
                  >
                    <Select style={{ width: 100 }} defaultValue={false}>
                      <Option value={true}>必需</Option>
                      <Option value={false}>可选</Option>
                    </Select>
                  </Form.Item>
                  <MinusCircleOutlined onClick={() => remove(name)} />
                </Space>
              ))}
            </>
          )}
        </Form.List>

        <Divider />

        <Form.Item
          name="content"
          label={t('prompts.content')}
          rules={[
            { required: true, message: t('prompts.contentRequired') },
            { max: 10000, message: '内容最多10000字符' },
          ]}
        >
          <TextArea
            rows={12}
            placeholder={t('prompts.contentPlaceholder')}
            showCount
            maxLength={10000}
          />
        </Form.Item>
      </Form>
    </Drawer>
  );
};

export default PromptForm;
