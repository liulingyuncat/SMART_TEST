import React, { useState } from 'react';
import { Modal, Form, Input, Select, Space, Button, message } from 'antd';
import { useTranslation } from 'react-i18next';
import PropTypes from 'prop-types';
import { createExecutionTask } from '../../../api/executionTask';

const { Option } = Select;

const CreateTaskModal = ({ visible, projectId, onSuccess, onCancel }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      console.log('🚀 [CreateTaskModal] Creating task with values:', values);
      setLoading(true);
      
      const response = await createExecutionTask(projectId, values);
      console.log('✅ [CreateTaskModal] Task created successfully:', response);
      
      message.success(t('testExecution.createModal.createSuccess'));
      form.resetFields();
      console.log('🔄 [CreateTaskModal] Calling onSuccess callback');
      onSuccess();
    } catch (error) {
      console.error('❌ [CreateTaskModal] Create task error:', error);
      if (error.errorFields) {
        // 表单验证失败
        return;
      }
      
      // 检查是否是任务名重复错误 (409 Conflict)
      if (error.response?.status === 409) {
        message.error(t('testExecution.createModal.taskNameExists'));
      } else {
        message.error(t('testExecution.createModal.createFailed'));
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title={t('testExecution.createModal.title')}
      open={visible}
      onCancel={onCancel}
      footer={null}
      width={520}
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{ task_status: 'pending' }}
      >
        <Form.Item
          label={t('testExecution.createModal.taskName')}
          name="task_name"
          rules={[{ required: true, message: t('testExecution.createModal.taskNameRequired') }]}
        >
          <Input
            placeholder={t('testExecution.createModal.taskNamePlaceholder')}
            maxLength={50}
            showCount
          />
        </Form.Item>

        <Form.Item
          label={t('testExecution.createModal.executionType')}
          name="execution_type"
          rules={[{ required: true, message: t('testExecution.createModal.executionTypeRequired') }]}
        >
          <Select placeholder={t('testExecution.createModal.executionTypeRequired')}>
            <Option value="manual">Manual Test</Option>
            <Option value="automation">AI Web</Option>
            <Option value="api">AI API</Option>
          </Select>
        </Form.Item>

        <Form.Item
          label={t('testExecution.createModal.taskStatus')}
          name="task_status"
        >
          <Select>
            <Option value="pending">{t('testExecution.status.pending')}</Option>
            <Option value="in_progress">{t('testExecution.status.inProgress')}</Option>
            <Option value="completed">{t('testExecution.status.completed')}</Option>
          </Select>
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
          <Space>
            <Button onClick={onCancel}>
              {t('testExecution.taskCard.cancel')}
            </Button>
            <Button type="primary" loading={loading} onClick={handleSubmit}>
              {t('testExecution.createModal.submit')}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  );
};

CreateTaskModal.propTypes = {
  visible: PropTypes.bool.isRequired,
  projectId: PropTypes.number.isRequired,
  onSuccess: PropTypes.func.isRequired,
  onCancel: PropTypes.func.isRequired,
};

export default CreateTaskModal;
