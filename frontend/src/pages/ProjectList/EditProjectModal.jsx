import React, { useState, useEffect } from 'react';
import { Modal, Form, Input, message } from 'antd';
import { useTranslation } from 'react-i18next';
import PropTypes from 'prop-types';
import { updateProject } from '../../api/project';

const EditProjectModal = ({ visible, project, onCancel, onSuccess }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  // 当对话框打开且project存在时,设置表单初始值
  useEffect(() => {
    if (visible && project) {
      form.setFieldsValue({
        name: project.name,
      });
    }
  }, [visible, project, form]);

  // 提交表单
  const handleSubmit = async () => {
    try {
      // 验证表单字段
      const values = await form.validateFields();
      setLoading(true);

      // 调用API更新项目
      const response = await updateProject(project.id, { name: values.name });

      // 成功提示
      message.success(t('project.renameSuccess'));

      // 回调父组件刷新数据
      onSuccess(response.data);
    } catch (error) {
      // 错误处理
      if (error.response && error.response.status === 400) {
        message.error(t('project.nameExists'));
      } else {
        message.error(t('project.updateFailed'));
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title={t('project.editProject')}
      open={visible}
      onCancel={onCancel}
      onOk={handleSubmit}
      confirmLoading={loading}
      destroyOnClose
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="name"
          label={t('project.name')}
          rules={[
            { required: true, message: t('project.nameRequired') },
            { max: 100, message: t('project.nameTooLong') },
          ]}
        >
          <Input placeholder={t('project.namePlaceholder')} disabled />
        </Form.Item>
      </Form>
    </Modal>
  );
};

EditProjectModal.propTypes = {
  visible: PropTypes.bool.isRequired,
  project: PropTypes.shape({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
  }),
  onCancel: PropTypes.func.isRequired,
  onSuccess: PropTypes.func.isRequired,
};

export default EditProjectModal;
