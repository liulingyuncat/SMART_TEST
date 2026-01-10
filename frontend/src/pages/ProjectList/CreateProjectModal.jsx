import React, { useState } from 'react';
import { Modal, Form, Input, Button, message } from 'antd';
import { useTranslation } from 'react-i18next';
import PropTypes from 'prop-types';
import { createProject } from '../../api/project';

const CreateProjectModal = ({ visible, onCancel, onSuccess }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (values) => {
    setSubmitting(true);
    console.log('[CreateProjectModal] Submitting values:', values);
    try {
      const data = await createProject(values);
      console.log('[CreateProjectModal] Success, received data:', data);
      form.resetFields();
      message.success(t('project.createSuccess'));
      
      // 触发全局事件通知侧边栏刷新
      window.dispatchEvent(new CustomEvent('projectCreated', { detail: data }));
      
      onSuccess(data);
    } catch (error) {
      console.error('[CreateProjectModal] Error details:', {
        message: error.message,
        response: error.response,
        status: error.response?.status,
        data: error.response?.data,
        config: error.config
      });
      if (error.response?.status === 400) {
        const errorMsg = error.response?.data?.message || error.message || t('project.nameExists');
        message.error(errorMsg);
      } else if (error.response?.status === 401) {
        message.error('未授权，请重新登录');
      } else {
        const errorMsg = error.response?.data?.message || error.message || t('project.createFailed');
        message.error(errorMsg);
      }
    } finally {
      setSubmitting(false);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={t('project.createTitle')}
      open={visible}
      onCancel={handleCancel}
      footer={[
        <Button key="cancel" onClick={handleCancel}>
          {t('common.cancel')}
        </Button>,
        <Button
          key="submit"
          type="primary"
          loading={submitting}
          onClick={() => form.submit()}
        >
          {t('common.confirm')}
        </Button>,
      ]}
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
      >
        <Form.Item
          name="name"
          label={t('project.name')}
          rules={[
            { required: true, message: t('project.nameRequired') },
            { max: 100, message: t('project.nameTooLong') },
            { whitespace: true, message: t('project.nameNotEmpty') },
          ]}
        >
          <Input
            placeholder={t('project.namePlaceholder')}
            maxLength={100}
          />
        </Form.Item>

        <Form.Item
          name="description"
          label={t('project.description')}
          rules={[
            { max: 500, message: t('project.descTooLong') },
          ]}
        >
          <Input.TextArea
            rows={4}
            placeholder={t('project.descPlaceholder')}
            maxLength={500}
            showCount
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

CreateProjectModal.propTypes = {
  visible: PropTypes.bool.isRequired,
  onCancel: PropTypes.func.isRequired,
  onSuccess: PropTypes.func.isRequired,
};

export default CreateProjectModal;
