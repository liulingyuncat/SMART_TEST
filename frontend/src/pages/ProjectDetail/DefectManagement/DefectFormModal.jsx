import React, { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Modal, Form, Input, Select, Row, Col } from 'antd';
import {
  DEFECT_STATUS,
  DEFECT_PRIORITY,
  DEFECT_SEVERITY,
  getStatusOptions,
  getPriorityOptions,
  getSeverityOptions,
  getAvailableStatusOptions,
} from '../../../constants/defect';

const { TextArea } = Input;
const { Option } = Select;

/**
 * 缺陷表单模态框 - 用于创建和编辑缺陷
 */
const DefectFormModal = ({
  visible,
  mode, // 'create' | 'edit'
  defect,
  subjects,
  phases,
  onCancel,
  onSubmit,
}) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const isEdit = mode === 'edit';

  // 重置或填充表单
  useEffect(() => {
    if (visible) {
      if (isEdit && defect) {
        form.setFieldsValue({
          summary: defect.summary,
          description: defect.description,
          status: defect.status,
          priority: defect.priority,
          severity: defect.severity,
          subject: defect.subject,
          phase: defect.phase,
          case_id: defect.case_id,
        });
      } else {
        form.resetFields();
        // 设置默认值
        form.setFieldsValue({
          status: DEFECT_STATUS.NEW,
          priority: DEFECT_PRIORITY.B,
          severity: DEFECT_SEVERITY.B,
        });
      }
    }
  }, [visible, isEdit, defect, form]);

  // 获取状态选项
  const statusOptions = isEdit && defect
    ? [
        { value: defect.status, label: t(`defect.status${defect.status}`) },
        ...getAvailableStatusOptions(defect.status, t),
      ]
    : getStatusOptions(t);

  // 提交表单
  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      onSubmit(values);
    } catch (error) {
      console.error('Form validation failed:', error);
    }
  };

  return (
    <Modal
      title={isEdit ? t('defect.edit') : t('defect.create')}
      open={visible}
      onCancel={onCancel}
      onOk={handleOk}
      width={720}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        className="defect-form-container"
      >
        {/* 概述 */}
        <Form.Item
          name="summary"
          label={t('defect.summary')}
          rules={[
            { required: true, message: t('validation.required') },
            { max: 500, message: t('validation.maxLength', { max: 500 }) },
          ]}
        >
          <Input placeholder={t('defect.summary')} maxLength={500} />
        </Form.Item>

        {/* 详细描述 */}
        <Form.Item
          name="description"
          label={t('defect.description')}
          rules={[
            { max: 5000, message: t('validation.maxLength', { max: 5000 }) },
          ]}
        >
          <TextArea
            rows={4}
            placeholder={t('defect.description')}
            maxLength={5000}
            showCount
          />
        </Form.Item>

        <Row gutter={16}>
          {/* 状态 */}
          <Col span={8}>
            <Form.Item
              name="status"
              label={t('defect.status')}
              rules={[{ required: true, message: t('validation.required') }]}
            >
              <Select placeholder={t('defect.status')}>
                {statusOptions.map(opt => (
                  <Option key={opt.value} value={opt.value}>{opt.label}</Option>
                ))}
              </Select>
            </Form.Item>
          </Col>

          {/* 优先级 */}
          <Col span={8}>
            <Form.Item
              name="priority"
              label={t('defect.priority')}
              rules={[{ required: true, message: t('validation.required') }]}
            >
              <Select placeholder={t('defect.priority')}>
                {getPriorityOptions(t).map(opt => (
                  <Option key={opt.value} value={opt.value}>{opt.label}</Option>
                ))}
              </Select>
            </Form.Item>
          </Col>

          {/* 严重程度 */}
          <Col span={8}>
            <Form.Item
              name="severity"
              label={t('defect.severity')}
              rules={[{ required: true, message: t('validation.required') }]}
            >
              <Select placeholder={t('defect.severity')}>
                {getSeverityOptions(t).map(opt => (
                  <Option key={opt.value} value={opt.value}>{opt.label}</Option>
                ))}
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          {/* 主题 */}
          <Col span={8}>
            <Form.Item
              name="subject"
              label={t('defect.subject')}
            >
              <Select placeholder={t('defect.allSubject')} allowClear>
                {subjects.map(s => (
                  <Option key={s.id} value={s.name}>{s.name}</Option>
                ))}
              </Select>
            </Form.Item>
          </Col>

          {/* 发现阶段 */}
          <Col span={8}>
            <Form.Item
              name="phase"
              label={t('defect.phase')}
            >
              <Select placeholder={t('defect.allPhase')} allowClear>
                {phases.map(p => (
                  <Option key={p.id} value={p.name}>{p.name}</Option>
                ))}
              </Select>
            </Form.Item>
          </Col>

          {/* Case ID */}
          <Col span={8}>
            <Form.Item
              name="case_id"
              label={t('defect.caseId')}
            >
              <Input placeholder={t('defect.caseIdPlaceholder')} maxLength={50} />
            </Form.Item>
          </Col>
        </Row>
      </Form>
    </Modal>
  );
};

export default DefectFormModal;
