import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  message,
  Popconfirm,
  Spin,
  Card,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
} from '@ant-design/icons';
import {
  createDefectSubject,
  updateDefectSubject,
  deleteDefectSubject,
  createDefectPhase,
  updateDefectPhase,
  deleteDefectPhase,
} from '../../../api/defect';

const { TextArea } = Input;

/**
 * 缺陷配置管理组件 - 主题和阶段配置
 */
const DefectConfig = ({
  projectId,
  subjects,
  phases,
  loading,
  onUpdate,
}) => {
  const { t } = useTranslation();
  
  // 主题编辑状态
  const [subjectModalVisible, setSubjectModalVisible] = useState(false);
  const [subjectMode, setSubjectMode] = useState('create');
  const [currentSubject, setCurrentSubject] = useState(null);
  const [subjectForm] = Form.useForm();
  const [subjectLoading, setSubjectLoading] = useState(false);

  // 阶段编辑状态
  const [phaseModalVisible, setPhaseModalVisible] = useState(false);
  const [phaseMode, setPhaseMode] = useState('create');
  const [currentPhase, setCurrentPhase] = useState(null);
  const [phaseForm] = Form.useForm();
  const [phaseLoading, setPhaseLoading] = useState(false);

  // ========== 主题操作 ==========
  const handleAddSubject = () => {
    setSubjectMode('create');
    setCurrentSubject(null);
    subjectForm.resetFields();
    setSubjectModalVisible(true);
  };

  const handleEditSubject = (record) => {
    setSubjectMode('edit');
    setCurrentSubject(record);
    subjectForm.setFieldsValue({
      name: record.name,
      description: record.description,
    });
    setSubjectModalVisible(true);
  };

  const handleDeleteSubject = async (record) => {
    try {
      await deleteDefectSubject(projectId, record.id);
      message.success(t('message.deleteSuccess'));
      onUpdate();
    } catch (error) {
      console.error('Failed to delete subject:', error);
      message.error(t('message.deleteFailed'));
    }
  };

  const handleSubjectSubmit = async () => {
    try {
      const values = await subjectForm.validateFields();
      setSubjectLoading(true);
      
      if (subjectMode === 'create') {
        await createDefectSubject(projectId, values);
        message.success(t('message.createSuccess'));
      } else {
        await updateDefectSubject(projectId, currentSubject.id, values);
        message.success(t('message.updateSuccess'));
      }
      
      setSubjectModalVisible(false);
      onUpdate();
    } catch (error) {
      console.error('Failed to save subject:', error);
      if (error.response?.status === 409) {
        message.error(t('defect.duplicateSubject'));
      } else {
        message.error(t('message.saveFailed'));
      }
    } finally {
      setSubjectLoading(false);
    }
  };

  // ========== 阶段操作 ==========
  const handleAddPhase = () => {
    setPhaseMode('create');
    setCurrentPhase(null);
    phaseForm.resetFields();
    setPhaseModalVisible(true);
  };

  const handleEditPhase = (record) => {
    setPhaseMode('edit');
    setCurrentPhase(record);
    phaseForm.setFieldsValue({
      name: record.name,
      description: record.description,
    });
    setPhaseModalVisible(true);
  };

  const handleDeletePhase = async (record) => {
    try {
      await deleteDefectPhase(projectId, record.id);
      message.success(t('message.deleteSuccess'));
      onUpdate();
    } catch (error) {
      console.error('Failed to delete phase:', error);
      message.error(t('message.deleteFailed'));
    }
  };

  const handlePhaseSubmit = async () => {
    try {
      const values = await phaseForm.validateFields();
      setPhaseLoading(true);
      
      if (phaseMode === 'create') {
        await createDefectPhase(projectId, values);
        message.success(t('message.createSuccess'));
      } else {
        await updateDefectPhase(projectId, currentPhase.id, values);
        message.success(t('message.updateSuccess'));
      }
      
      setPhaseModalVisible(false);
      onUpdate();
    } catch (error) {
      console.error('Failed to save phase:', error);
      if (error.response?.status === 409) {
        message.error(t('defect.duplicatePhase'));
      } else {
        message.error(t('message.saveFailed'));
      }
    } finally {
      setPhaseLoading(false);
    }
  };

  // 主题表格列
  const subjectColumns = [
    {
      title: t('defect.subjectName'),
      dataIndex: 'name',
      key: 'name',
      width: 200,
    },
    {
      title: t('defect.subjectDescription'),
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: t('common.actions'),
      key: 'actions',
      width: 120,
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEditSubject(record)}
          />
          <Popconfirm
            title={t('defect.deleteSubject') + '?'}
            onConfirm={() => handleDeleteSubject(record)}
            okText={t('common.confirm')}
            cancelText={t('common.cancel')}
          >
            <Button
              type="link"
              size="small"
              danger
              icon={<DeleteOutlined />}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // 阶段表格列
  const phaseColumns = [
    {
      title: t('defect.phaseName'),
      dataIndex: 'name',
      key: 'name',
      width: 200,
    },
    {
      title: t('defect.phaseDescription'),
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: t('common.actions'),
      key: 'actions',
      width: 120,
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEditPhase(record)}
          />
          <Popconfirm
            title={t('defect.deletePhase') + '?'}
            onConfirm={() => handleDeletePhase(record)}
            okText={t('common.confirm')}
            cancelText={t('common.cancel')}
          >
            <Button
              type="link"
              size="small"
              danger
              icon={<DeleteOutlined />}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Spin spinning={loading}>
      <div className="defect-config-container">
        {/* 主题管理 */}
        <Card
          title={t('defect.subjects')}
          className="defect-config-section"
          extra={
            <Button
              type="primary"
              size="small"
              icon={<PlusOutlined />}
              onClick={handleAddSubject}
            >
              {t('defect.addSubject')}
            </Button>
          }
        >
          <Table
            className="defect-config-table"
            columns={subjectColumns}
            dataSource={subjects}
            rowKey="id"
            pagination={false}
            size="small"
          />
        </Card>

        {/* 阶段管理 */}
        <Card
          title={t('defect.phases')}
          className="defect-config-section"
          extra={
            <Button
              type="primary"
              size="small"
              icon={<PlusOutlined />}
              onClick={handleAddPhase}
            >
              {t('defect.addPhase')}
            </Button>
          }
        >
          <Table
            className="defect-config-table"
            columns={phaseColumns}
            dataSource={phases}
            rowKey="id"
            pagination={false}
            size="small"
          />
        </Card>

        {/* 主题编辑模态框 */}
        <Modal
          title={subjectMode === 'create' ? t('defect.addSubject') : t('defect.editSubject')}
          open={subjectModalVisible}
          onCancel={() => setSubjectModalVisible(false)}
          onOk={handleSubjectSubmit}
          confirmLoading={subjectLoading}
          destroyOnClose
        >
          <Form form={subjectForm} layout="vertical">
            <Form.Item
              name="name"
              label={t('defect.subjectName')}
              rules={[
                { required: true, message: t('validation.required') },
                { max: 100, message: t('validation.maxLength', { max: 100 }) },
              ]}
            >
              <Input placeholder={t('defect.subjectName')} maxLength={100} />
            </Form.Item>
            <Form.Item
              name="description"
              label={t('defect.subjectDescription')}
              rules={[
                { max: 500, message: t('validation.maxLength', { max: 500 }) },
              ]}
            >
              <TextArea
                rows={3}
                placeholder={t('defect.subjectDescription')}
                maxLength={500}
              />
            </Form.Item>
          </Form>
        </Modal>

        {/* 阶段编辑模态框 */}
        <Modal
          title={phaseMode === 'create' ? t('defect.addPhase') : t('defect.editPhase')}
          open={phaseModalVisible}
          onCancel={() => setPhaseModalVisible(false)}
          onOk={handlePhaseSubmit}
          confirmLoading={phaseLoading}
          destroyOnClose
        >
          <Form form={phaseForm} layout="vertical">
            <Form.Item
              name="name"
              label={t('defect.phaseName')}
              rules={[
                { required: true, message: t('validation.required') },
                { max: 100, message: t('validation.maxLength', { max: 100 }) },
              ]}
            >
              <Input placeholder={t('defect.phaseName')} maxLength={100} />
            </Form.Item>
            <Form.Item
              name="description"
              label={t('defect.phaseDescription')}
              rules={[
                { max: 500, message: t('validation.maxLength', { max: 500 }) },
              ]}
            >
              <TextArea
                rows={3}
                placeholder={t('defect.phaseDescription')}
                maxLength={500}
              />
            </Form.Item>
          </Form>
        </Modal>
      </div>
    </Spin>
  );
};

export default DefectConfig;
