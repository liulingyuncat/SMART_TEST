import React, { useState, useEffect } from 'react';
import { Modal, Form, Input, message } from 'antd';

const { TextArea } = Input;

/**
 * 多语言编辑对话框组件
 * @param {Object} props
 * @param {boolean} props.visible - 对话框是否可见
 * @param {string} props.title - 对话框标题
 * @param {string} props.fieldName - 字段名称 (major_function/middle_function/minor_function/precondition/test_steps/expected_result)
 * @param {Object} props.data - 数据对象 {cn: string, jp: string, en: string}
 * @param {Function} props.onSave - 保存回调 (data) => void
 * @param {Function} props.onCancel - 取消回调
 */
const MultiLangEditModal = ({ visible, title, fieldName, data, onSave, onCancel }) => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  // 当对话框打开或数据变化时，更新表单
  useEffect(() => {
    if (visible && data) {
      form.setFieldsValue({
        cn: data.cn || '',
        jp: data.jp || '',
        en: data.en || '',
      });
    }
  }, [visible, data, form]);

  // 处理保存
  const handleOk = async () => {
    console.log('🔵 MultiLangEditModal handleOk 开始');
    try {
      setLoading(true);
      console.log('设置 loading = true');
      
      const values = await form.validateFields();
      console.log('表单验证成功，值:', values);
      
      // 调用父组件的保存回调
      console.log('调用 onSave 回调, fieldName:', fieldName);
      await onSave({
        fieldName,
        cn: values.cn,
        jp: values.jp,
        en: values.en,
      });
      console.log('onSave 回调完成');
      
      // 保存成功后重置表单
      form.resetFields();
      console.log('表单已重置');
      console.log('🔵 MultiLangEditModal handleOk 成功结束');
    } catch (error) {
      console.error('❌ MultiLangEditModal handleOk 失败:', error);
      if (error.errorFields) {
        message.error('请检查输入内容');
      } else {
        // 父组件已经显示了错误消息，这里不需要重复显示
        // message.error('保存失败');
      }
    } finally {
      // 确保加载状态被重置
      console.log('设置 loading = false');
      setLoading(false);
    }
  };

  // 处理取消
  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  // 判断是否为文本域字段
  const isTextAreaField = ['precondition', 'test_steps', 'expected_result'].includes(fieldName);

  return (
    <Modal
      title={title}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
      okText="保存"
      cancelText="取消"
      width={600}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        preserve={false}
      >
        <Form.Item
          label="中文 (CN)"
          name="cn"
          rules={[{ max: isTextAreaField ? 2000 : 100, message: `最多${isTextAreaField ? 2000 : 100}个字符` }]}
        >
          {isTextAreaField ? (
            <TextArea
              placeholder="请输入中文内容"
              autoSize={{ minRows: 3, maxRows: 8 }}
            />
          ) : (
            <Input placeholder="请输入中文内容" />
          )}
        </Form.Item>

        <Form.Item
          label="English (EN)"
          name="en"
          rules={[{ max: isTextAreaField ? 2000 : 100, message: `最多${isTextAreaField ? 2000 : 100}个字符` }]}
        >
          {isTextAreaField ? (
            <TextArea
              placeholder="Please enter English content"
              autoSize={{ minRows: 3, maxRows: 8 }}
            />
          ) : (
            <Input placeholder="Please enter English content" />
          )}
        </Form.Item>

        <Form.Item
          label="日本語 (JP)"
          name="jp"
          rules={[{ max: isTextAreaField ? 2000 : 100, message: `最多${isTextAreaField ? 2000 : 100}个字符` }]}
        >
          {isTextAreaField ? (
            <TextArea
              placeholder="日本語の内容を入力してください"
              autoSize={{ minRows: 3, maxRows: 8 }}
            />
          ) : (
            <Input placeholder="日本語の内容を入力してください" />
          )}
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default MultiLangEditModal;
