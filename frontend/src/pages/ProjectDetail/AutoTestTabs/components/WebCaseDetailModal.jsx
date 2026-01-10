import React, { useState, useEffect } from 'react';
import { Modal, Form, Input, Row, Col, Button, Tag, message } from 'antd';
import { SaveOutlined, CloseOutlined } from '@ant-design/icons';
import PropTypes from 'prop-types';
import './WebCaseDetailModal.css';

const { TextArea } = Input;

/**
 * Web用例详细信息弹窗 - 全字段可编辑版
 * 所有字段均可编辑，包括 Script Code
 * 设计风格与 ApiCaseDetailModal 保持一致
 */
const WebCaseDetailModal = ({
  visible,
  caseData,
  language = 'cn', // 当前语言: cn, jp, en
  onSave,
  onCancel,
}) => {
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);

  // 根据语言获取字段后缀
  const getLangSuffix = () => {
    switch (language) {
      case 'jp': return '_jp';
      case 'en': return '_en';
      default: return '_cn';
    }
  };

  // 根据语言获取字段值
  const getFieldValue = (fieldPrefix) => {
    const suffix = getLangSuffix();
    return caseData?.[fieldPrefix + suffix] || '';
  };

  useEffect(() => {
    if (visible && caseData) {
      const suffix = getLangSuffix();
      form.setFieldsValue({
        screen: caseData[`screen${suffix}`] || '',
        function: caseData[`function${suffix}`] || '',
        precondition: caseData[`precondition${suffix}`] || '',
        test_steps: caseData[`test_steps${suffix}`] || '',
        expected_result: caseData[`expected_result${suffix}`] || '',
        script_code: caseData.script_code || '',
      });
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [visible, caseData, language]);

  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      setSaving(true);
      
      const suffix = getLangSuffix();
      const updateData = {
        case_id: caseData.case_id,
        [`screen${suffix}`]: values.screen,
        [`function${suffix}`]: values.function,
        [`precondition${suffix}`]: values.precondition,
        [`test_steps${suffix}`]: values.test_steps,
        [`expected_result${suffix}`]: values.expected_result,
        script_code: values.script_code,
      };
      
      await onSave(updateData);
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

  // 语言标签
  const langLabels = { cn: '中文', jp: '日本語', en: 'English' };
  const langColors = { cn: 'blue', jp: 'purple', en: 'green' };

  return (
    <Modal
      title={
        <div className="wcd-modal-title">
          <span>用例详细信息</span>
          <Tag color="blue" style={{ marginLeft: 8 }}>
            No.{caseData.no || caseData.display_order || caseData.id || '?'}
          </Tag>
          <Tag color="cyan">Web</Tag>
          <Tag color={langColors[language] || 'default'}>
            {langLabels[language] || language}
          </Tag>
        </div>
      }
      open={visible}
      onCancel={onCancel}
      width={850}
      className="wcd-modal-compact"
      footer={
        <div className="wcd-modal-footer">
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
      <Form form={form} layout="vertical" className="wcd-form">
        {/* 第一行: Screen / Function */}
        <Row gutter={16}>
          <Col span={12}>
            <Form.Item name="screen" label="Screen" style={{ marginBottom: 12 }}>
              <Input placeholder="画面名称" />
            </Form.Item>
          </Col>
          <Col span={12}>
            <Form.Item name="function" label="Function" style={{ marginBottom: 12 }}>
              <Input placeholder="功能名称" />
            </Form.Item>
          </Col>
        </Row>

        {/* Precondition */}
        <Form.Item name="precondition" label="Precondition" style={{ marginBottom: 12 }}>
          <TextArea 
            autoSize={{ minRows: 2, maxRows: 4 }} 
            placeholder="前置条件"
            style={{ fontSize: 12 }}
          />
        </Form.Item>

        {/* Test Steps */}
        <Form.Item name="test_steps" label="Test Steps" style={{ marginBottom: 12 }}>
          <TextArea 
            autoSize={{ minRows: 3, maxRows: 6 }} 
            placeholder="测试步骤"
            style={{ fontSize: 12 }}
          />
        </Form.Item>

        {/* Expected Result */}
        <Form.Item name="expected_result" label="Expected Result" style={{ marginBottom: 12 }}>
          <TextArea 
            autoSize={{ minRows: 2, maxRows: 4 }} 
            placeholder="期望结果"
            style={{ fontSize: 12 }}
          />
        </Form.Item>

        {/* Script Code */}
        <Form.Item name="script_code" label="Script Code" style={{ marginBottom: 0 }}>
          <TextArea 
            autoSize={{ minRows: 4, maxRows: 12 }} 
            placeholder={`// Playwright 脚本代码
await page.goto('https://example.com');
await page.click('button[type="submit"]');
await expect(page.locator('.success')).toBeVisible();`}
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

WebCaseDetailModal.propTypes = {
  visible: PropTypes.bool.isRequired,
  caseData: PropTypes.object,
  language: PropTypes.string,
  onSave: PropTypes.func.isRequired,
  onCancel: PropTypes.func.isRequired,
};

export default WebCaseDetailModal;
