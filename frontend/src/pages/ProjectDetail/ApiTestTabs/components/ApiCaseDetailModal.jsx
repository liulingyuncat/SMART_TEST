import React, { useState, useEffect } from 'react';
import { Modal, Form, Input, Select, Row, Col, Button, Tag, message } from 'antd';
import { SaveOutlined, CloseOutlined } from '@ant-design/icons';
import PropTypes from 'prop-types';
import './ApiCaseDetailModal.css';

const { Option } = Select;
const { TextArea } = Input;

/**
 * API用例详细信息弹窗 - 全字段可编辑版
 * 所有字段均可编辑，包括 Script Code
 */
const ApiCaseDetailModal = ({
  visible,
  caseData,
  onSave,
  onCancel,
}) => {
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (visible && caseData) {
      form.setFieldsValue({
        screen: caseData.screen || '',
        method: caseData.method || 'GET',
        url: caseData.url || '',
        header: caseData.header || '',
        body: caseData.body || '',
        response: caseData.response || '',
        script_code: caseData.script_code || '',
      });
    }
  }, [visible, caseData, form]);

  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      setSaving(true);
      await onSave({
        case_id: caseData.case_id,
        screen: values.screen,
        method: values.method,
        url: values.url,
        header: values.header,
        body: values.body,
        response: values.response,
        script_code: values.script_code,
      });
      message.success('保存成功');
      onCancel();
    } catch (error) {
      console.error('保存失败:', error);
      message.error('保存失败');
    } finally {
      setSaving(false);
    }
  };

  if (!caseData) return null;

  return (
    <Modal
      title={
        <div className="acd-modal-title">
          <span>用例详细信息</span>
          <Tag color="blue" style={{ marginLeft: 8 }}>
            No.{caseData.no || caseData.display_order || '?'}
          </Tag>
          <Tag color="orange">API</Tag>
        </div>
      }
      open={visible}
      onCancel={onCancel}
      width={850}
      className="acd-modal-compact"
      footer={
        <div className="acd-modal-footer">
          <Button icon={<CloseOutlined />} onClick={onCancel}>
            取消
          </Button>
          <Button 
            type="primary" 
            icon={<SaveOutlined />} 
            loading={saving} 
            onClick={handleSave}
          >
            保存并关闭
          </Button>
        </div>
      }
      destroyOnClose
    >
      <Form form={form} layout="vertical" className="acd-form">
        {/* 第一行: Screen / Method */}
        <Row gutter={16}>
          <Col span={16}>
            <Form.Item name="screen" label="Screen" style={{ marginBottom: 12 }}>
              <Input placeholder="画面名称，如 [ダッシュボード]" />
            </Form.Item>
          </Col>
          <Col span={8}>
            <Form.Item name="method" label="Method" style={{ marginBottom: 12 }}>
              <Select>
                <Option value="GET"><Tag color="green">GET</Tag></Option>
                <Option value="POST"><Tag color="blue">POST</Tag></Option>
                <Option value="PUT"><Tag color="orange">PUT</Tag></Option>
                <Option value="DELETE"><Tag color="red">DELETE</Tag></Option>
                <Option value="PATCH"><Tag color="purple">PATCH</Tag></Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        {/* URL */}
        <Form.Item name="url" label="URL" style={{ marginBottom: 12 }}>
          <Input placeholder="/api/xxx" style={{ fontFamily: 'monospace' }} />
        </Form.Item>

        {/* Header */}
        <Form.Item name="header" label="Header" style={{ marginBottom: 12 }}>
          <TextArea 
            autoSize={{ minRows: 1, maxRows: 3 }} 
            placeholder='{"Authorization": "Bearer {{token}}"}'
            style={{ fontFamily: 'monospace', fontSize: 12 }}
          />
        </Form.Item>

        {/* Body */}
        <Form.Item name="body" label="Body" style={{ marginBottom: 12 }}>
          <TextArea 
            autoSize={{ minRows: 1, maxRows: 4 }} 
            placeholder='{"key": "value"}'
            style={{ fontFamily: 'monospace', fontSize: 12 }}
          />
        </Form.Item>

        {/* Response */}
        <Form.Item name="response" label="Response" style={{ marginBottom: 12 }}>
          <TextArea 
            autoSize={{ minRows: 1, maxRows: 4 }} 
            placeholder='{"code": 200, "data": {...}}'
            style={{ fontFamily: 'monospace', fontSize: 12 }}
          />
        </Form.Item>

        {/* Script Code */}
        <Form.Item name="script_code" label="Script Code" style={{ marginBottom: 0 }}>
          <TextArea 
            autoSize={{ minRows: 3, maxRows: 10 }} 
            placeholder={`async () => {
  const token = localStorage.getItem('token');
  const res = await fetch('/api/xxx', {
    method: 'GET',
    headers: { 'Authorization': \`Bearer \${token}\` }
  });
  return await res.json();
}`}
            style={{ 
              fontFamily: 'Consolas, Monaco, monospace', 
              fontSize: 12,
              backgroundColor: '#1e1e1e',
              color: '#d4d4d4',
            }}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

ApiCaseDetailModal.propTypes = {
  visible: PropTypes.bool.isRequired,
  caseData: PropTypes.object,
  onSave: PropTypes.func.isRequired,
  onCancel: PropTypes.func.isRequired,
};

export default ApiCaseDetailModal;
