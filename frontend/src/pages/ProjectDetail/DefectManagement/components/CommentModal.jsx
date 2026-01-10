import React, { useState, useEffect } from 'react';
import { Modal, Form, Input, message } from 'antd';
import { useTranslation } from 'react-i18next';

const { TextArea } = Input;

const CommentModal = ({ visible, mode, initialContent, onOk, onCancel }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);

  // 定义fallback翻译
  const labels = {
    addComment: t('defect.addComment', '新增说明'),
    editComment: t('defect.editComment', '编辑说明'),
    commentContent: t('defect.commentContent', '说明内容'),
    commentPlaceholder: t('defect.commentPlaceholder', '请输入说明内容（最多2000字符）'),
    save: t('common.save', '保存'),
    cancel: t('common.cancel', '取消'),
    required: t('validation.required', '此字段为必填项'),
  };

  // 调试输出
  useEffect(() => {
    console.log('[DEBUG] CommentModal labels:', labels);
  }, []);

  useEffect(() => {
    if (visible) {
      if (mode === 'edit' && initialContent) {
        form.setFieldsValue({ content: initialContent });
      } else {
        form.resetFields();
      }
    }
  }, [visible, mode, initialContent, form]);

  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      setSaving(true);
      await onOk(values.content);
      form.resetFields();
      setSaving(false);
    } catch (error) {
      console.error('Form validation failed:', error);
      setSaving(false);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={mode === 'add' ? labels.addComment : labels.editComment}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={saving}
      okText={labels.save}
      cancelText={labels.cancel}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="content"
          label={labels.commentContent}
          rules={[
            { required: true, message: labels.required },
            { max: 2000, message: t('validation.maxLength', { max: 2000 }) },
          ]}
        >
          <TextArea
            rows={6}
            placeholder={labels.commentPlaceholder}
            maxLength={2000}
            showCount
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CommentModal;
