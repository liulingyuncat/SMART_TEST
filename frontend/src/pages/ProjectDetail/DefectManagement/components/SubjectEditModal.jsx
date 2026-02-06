import React, { useState, useEffect, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Modal, Form, Input, Button, List, Popconfirm, message } from 'antd';
import { PlusOutlined, DeleteOutlined, EditOutlined, SaveOutlined, CloseOutlined } from '@ant-design/icons';
import {
  fetchDefectSubjects,
  createDefectSubject,
  updateDefectSubject,
  deleteDefectSubject,
} from '../../../../api/defect';

/**
 * 缺陷分类(Subject)配置弹窗
 * 支持增删改查操作
 */
const SubjectEditModal = ({
  visible,
  projectId,
  subjects: initialSubjects,
  onClose,
  onUpdate,
}) => {
  const { t, i18n } = useTranslation();
  const [subjects, setSubjects] = useState([]);
  const [loading, setLoading] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [editForm] = Form.useForm();
  const [addForm] = Form.useForm();

  // 使用 useMemo 缓存翻译标签，只在语言变化时重新计算
  const labels = useMemo(() => ({
    editSubjects: t('defect.subjectsConfig', '主题配置'),
    subjectName: t('defect.subjectName', '主题名称'),
    description: t('common.description', '描述'),
    add: t('common.add', '添加'),
    required: t('validation.required', '此字段为必填项'),
    confirmDelete: t('common.confirmDelete', '确认删除吗？'),
    createSuccess: t('message.createSuccess', '创建成功'),
    createFailed: t('message.createFailed', '创建失败'),
    saveSuccess: t('message.saveSuccess', '保存成功'),
    saveFailed: t('message.saveFailed', '保存失败'),
    deleteSuccess: t('message.deleteSuccess', '删除成功'),
    deleteFailed: t('message.deleteFailed', '删除失败'),
  }), [t, i18n.language]);

  // 初始化
  useEffect(() => {
    if (visible) {
      setSubjects(initialSubjects || []);
      setEditingId(null);
    }
  }, [visible, initialSubjects]);

  // 刷新列表
  const refreshSubjects = async () => {
    try {
      const response = await fetchDefectSubjects(projectId);
      // 后端直接返回数组
      setSubjects(Array.isArray(response) ? response : []);
    } catch (error) {
      console.error('Failed to fetch subjects:', error);
    }
  };

  // 添加分类
  const handleAdd = async (values) => {
    if (!values.name?.trim()) {
      message.warning(labels.required);
      return;
    }
    setLoading(true);
    try {
      await createDefectSubject(projectId, values);
      message.success(labels.createSuccess);
      addForm.resetFields();
      await refreshSubjects();
      onUpdate?.();
      // 不再自动关闭对话框，由用户手动关闭
    } catch (error) {
      console.error('Failed to create subject:', error);
      message.error(labels.createFailed);
    } finally {
      setLoading(false);
    }
  };

  // 保存编辑
  const handleSave = async (id) => {
    try {
      const values = await editForm.validateFields();
      setLoading(true);
      await updateDefectSubject(projectId, id, values);
      message.success(labels.saveSuccess);
      setEditingId(null);
      await refreshSubjects();
      onUpdate?.();
      // 不再自动关闭对话框，由用户手动关闭
    } catch (error) {
      console.error('Failed to update subject:', error);
      message.error(labels.saveFailed);
    } finally {
      setLoading(false);
    }
  };

  // 删除分类
  const handleDelete = async (id) => {
    setLoading(true);
    try {
      await deleteDefectSubject(projectId, id);
      message.success(labels.deleteSuccess);
      await refreshSubjects();
      onUpdate?.();
      // 不再自动关闭对话框，由用户手动关闭
    } catch (error) {
      console.error('Failed to delete subject:', error);
      message.error(labels.deleteFailed);
    } finally {
      setLoading(false);
    }
  };

  // 开始编辑
  const startEdit = (subject) => {
    setEditingId(subject.id);
    editForm.setFieldsValue({ name: subject.name });
  };

  // 取消编辑
  const cancelEdit = () => {
    setEditingId(null);
    editForm.resetFields();
  };

  return (
    <Modal
      title={labels.editSubjects}
      open={visible}
      onCancel={onClose}
      footer={null}
      width={600}
    >
      {/* 添加表单 */}
      <Form form={addForm} layout="inline" onFinish={handleAdd} style={{ marginBottom: 16 }}>
        <Form.Item
          name="name"
          rules={[{ required: true, message: labels.required }]}
        >
          <Input placeholder={labels.subjectName} style={{ width: 300 }} />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" icon={<PlusOutlined />} loading={loading}>
            {labels.add}
          </Button>
        </Form.Item>
      </Form>

      {/* 分类列表 */}
      <List
        loading={loading}
        dataSource={subjects}
        renderItem={(subject) => (
          <List.Item
            actions={
              editingId === subject.id
                ? [
                    <Button
                      key="save"
                      type="link"
                      icon={<SaveOutlined />}
                      onClick={() => handleSave(subject.id)}
                    />,
                    <Button
                      key="cancel"
                      type="link"
                      icon={<CloseOutlined />}
                      onClick={cancelEdit}
                    />,
                  ]
                : [
                    <Button
                      key="edit"
                      type="link"
                      icon={<EditOutlined />}
                      onClick={() => startEdit(subject)}
                    />,
                    <Popconfirm
                      key="delete"
                      title={labels.confirmDelete}
                      onConfirm={() => handleDelete(subject.id)}
                    >
                      <Button type="link" danger icon={<DeleteOutlined />} />
                    </Popconfirm>,
                  ]
            }
          >
            {editingId === subject.id ? (
              <Form form={editForm} layout="inline" style={{ width: '100%' }}>
                <Form.Item
                  name="name"
                  rules={[{ required: true, message: labels.required }]}
                  style={{ flex: 1 }}
                >
                  <Input style={{ width: '100%' }} />
                </Form.Item>
              </Form>
            ) : (
              <List.Item.Meta
                title={subject.name}
              />
            )}
          </List.Item>
        )}
      />
    </Modal>
  );
};

export default SubjectEditModal;
