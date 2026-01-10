import React, { useState, useEffect, useCallback, useRef } from 'react';
import { Modal, Form, Input, message } from 'antd';
import { useTranslation } from 'react-i18next';
import debounce from 'lodash/debounce';
import { checkUnique } from '../api/user';
import { updateNickname } from '../api/profile';

/**
 * 昵称编辑弹窗组件
 * @param {Object} props
 * @param {boolean} props.visible - 弹窗是否可见
 * @param {string} props.currentNickname - 当前昵称
 * @param {Function} props.onCancel - 取消回调
 * @param {Function} props.onSuccess - 成功回调
 */
const NicknameEditModal = ({ visible, currentNickname, onCancel, onSuccess }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [nickname, setNickname] = useState('');
  const [loading, setLoading] = useState(false);
  const [checking, setChecking] = useState(false);
  const [nicknameError, setNicknameError] = useState('');

  // 使用 ref 保存 debounced 函数以避免重复创建
  const debouncedCheckRef = useRef(null);

  // 弹窗打开时填充当前昵称
  useEffect(() => {
    if (visible) {
      setNickname(currentNickname || '');
      form.setFieldsValue({ nickname: currentNickname || '' });
      setNicknameError('');
      setLoading(false);
      setChecking(false);
    }
  }, [visible, currentNickname, form]);

  // 昵称唯一性检查（防抖）
  const checkNicknameUnique = useCallback(async (value) => {
    if (!value || value === currentNickname) {
      setNicknameError('');
      setChecking(false);
      return;
    }

    setChecking(true);
    try {
      const response = await checkUnique({ nickname: value });
      // apiClient 已返回 data.data，所以直接访问 response.exists
      if (response?.exists) {
        setNicknameError(t('profile.nicknameExists') || '昵称已存在');
      } else {
        setNicknameError('');
      }
    } catch (error) {
      console.error('检查昵称唯一性失败:', error);
      // 检查失败时不阻止提交，后端会再次校验
      setNicknameError('');
    } finally {
      setChecking(false);
    }
  }, [currentNickname, t]);

  // 创建防抖检查函数
  useEffect(() => {
    debouncedCheckRef.current = debounce(checkNicknameUnique, 300);
    return () => {
      if (debouncedCheckRef.current) {
        debouncedCheckRef.current.cancel();
      }
    };
  }, [checkNicknameUnique]);

  // 昵称输入变化处理
  const handleNicknameChange = (e) => {
    const value = e.target.value;
    setNickname(value);
    form.setFieldsValue({ nickname: value });

    // 长度校验
    if (value.length < 2 || value.length > 50) {
      setNicknameError(t('profile.nicknameLengthError') || '昵称长度需要在2-50个字符之间');
      return;
    }

    // 触发防抖的唯一性检查
    if (debouncedCheckRef.current) {
      debouncedCheckRef.current(value);
    }
  };

  // 提交处理
  const handleSubmit = async () => {
    // 验证表单
    try {
      await form.validateFields();
    } catch {
      return;
    }

    // 检查是否有错误
    if (nicknameError) {
      return;
    }

    // 检查昵称是否改变
    if (nickname === currentNickname) {
      message.info(t('profile.nicknameUnchanged') || '昵称未改变');
      onCancel();
      return;
    }

    setLoading(true);
    try {
      await updateNickname(nickname);
      message.success(t('message.nicknameUpdated') || '昵称更新成功');
      onSuccess();
    } catch (error) {
      console.error('更新昵称失败:', error);
      const errorMsg = error.response?.data?.message || error.message;
      if (errorMsg?.includes('already exists')) {
        message.error(t('profile.nicknameExists') || '昵称已存在');
      } else {
        message.error(t('message.nicknameUpdateFailed') || '昵称更新失败');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title={t('profile.changeNickname') || '修改昵称'}
      open={visible}
      onOk={handleSubmit}
      onCancel={onCancel}
      confirmLoading={loading}
      okButtonProps={{ disabled: loading || !!nicknameError || checking }}
      destroyOnClose
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="nickname"
          label={t('profile.nickname') || '昵称'}
          validateStatus={nicknameError ? 'error' : checking ? 'validating' : ''}
          help={nicknameError}
          rules={[
            { required: true, message: t('profile.nicknameRequired') || '请输入昵称' },
            { min: 2, message: t('profile.nicknameTooShort') || '昵称至少2个字符' },
            { max: 50, message: t('profile.nicknameTooLong') || '昵称最多50个字符' },
          ]}
        >
          <Input
            value={nickname}
            onChange={handleNicknameChange}
            placeholder={t('profile.nicknamePlaceholder') || '请输入昵称'}
            maxLength={50}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default NicknameEditModal;
