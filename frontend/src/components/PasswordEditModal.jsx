import React, { useState } from 'react';
import { Modal, Form, Input, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { changePassword } from '../api/profile';

/**
 * 密码修改弹窗组件 - T23
 * @param {Object} props
 * @param {boolean} props.visible - 弹窗是否可见
 * @param {() => void} props.onCancel - 取消回调
 * @param {() => void} props.onSuccess - 成功回调
 */
const PasswordEditModal = ({ visible, onCancel, onSuccess }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      await changePassword(values.current_password, values.new_password);
      
      message.success(t('profile.passwordUpdateSuccess') || '密码修改成功');
      form.resetFields();
      onSuccess?.();
    } catch (error) {
      console.error('密码修改失败:', error);
      const errorMessage = error.response?.data?.message || error.message;
      if (errorMessage === 'current password is incorrect') {
        message.error(t('profile.currentPasswordIncorrect') || '当前密码错误');
      } else if (errorMessage === 'new password cannot be same as current') {
        message.error(t('profile.newPasswordSameAsCurrent') || '新密码不能与当前密码相同');
      } else {
        message.error(t('profile.passwordUpdateFailed') || '密码修改失败');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel?.();
  };

  return (
    <Modal
      title={t('profile.changePassword') || '修改密码'}
      open={visible}
      onOk={handleSubmit}
      onCancel={handleCancel}
      confirmLoading={loading}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        autoComplete="off"
      >
        <Form.Item
          name="current_password"
          label={t('profile.currentPassword') || '当前密码'}
          rules={[
            { required: true, message: t('profile.currentPasswordRequired') || '请输入当前密码' },
          ]}
        >
          <Input.Password placeholder={t('profile.currentPasswordPlaceholder') || '请输入当前密码'} />
        </Form.Item>

        <Form.Item
          name="new_password"
          label={t('profile.newPassword') || '新密码'}
          rules={[
            { required: true, message: t('profile.newPasswordRequired') || '请输入新密码' },
            { min: 6, message: t('profile.passwordMinLength') || '密码长度不能少于6位' },
            { max: 50, message: t('profile.passwordMaxLength') || '密码长度不能超过50位' },
          ]}
        >
          <Input.Password placeholder={t('profile.newPasswordPlaceholder') || '请输入新密码（6-50位）'} />
        </Form.Item>

        <Form.Item
          name="confirm_password"
          label={t('profile.confirmPassword') || '确认密码'}
          dependencies={['new_password']}
          rules={[
            { required: true, message: t('profile.confirmPasswordRequired') || '请确认新密码' },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('new_password') === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error(t('profile.passwordNotMatch') || '两次输入的密码不一致'));
              },
            }),
          ]}
        >
          <Input.Password placeholder={t('profile.confirmPasswordPlaceholder') || '请再次输入新密码'} />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default PasswordEditModal;
