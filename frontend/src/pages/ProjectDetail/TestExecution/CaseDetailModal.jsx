import React, { useState, useEffect } from 'react';
import { Modal, Form, Input, Select, Row, Col, Button, Tag } from 'antd';
import { SaveOutlined, CloseOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import PropTypes from 'prop-types';
import './CaseDetailModal.css';

const { Option } = Select;
const { TextArea } = Input;

/**
 * 用例详细信息弹窗 - 简洁紧凑版
 * 直接显示原用例表头，布局紧凑，自适应行数
 */
const CaseDetailModal = ({
  visible,
  caseData,
  executionType,
  languageSuffix,
  languageDisplay,
  onSave,
  onCancel,
}) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (visible && caseData) {
      form.setFieldsValue({
        test_result: caseData.test_result || 'Block',
        bug_id: caseData.bug_id || '',
        remark: caseData.remark || '',
      });
    }
  }, [visible, caseData, form]);

  const getFieldValue = (baseField) => {
    if (!caseData) return '-';
    const valueWithSuffix = caseData[`${baseField}${languageSuffix}`];
    if (valueWithSuffix) return valueWithSuffix;
    const valueWithCn = caseData[`${baseField}_cn`];
    if (valueWithCn) return valueWithCn;
    return caseData[baseField] || '-';
  };

  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      setSaving(true);
      await onSave({
        case_id: caseData.case_id,
        test_result: values.test_result,
        bug_id: values.bug_id,
        remark: values.remark,
      });
      onCancel();
    } catch (error) {
      console.error('保存失败:', error);
    } finally {
      setSaving(false);
    }
  };

  if (!caseData) return null;

  const isManual = executionType === 'manual';
  const isAPI = executionType === 'api';
  const isWeb = executionType === 'automation';

  // 紧凑型单行字段
  const FieldRow = ({ label, value, isCode = false, children }) => {
    const displayValue = value || '-';
    const isEmpty = !value || value === '-';
    return (
      <div className="cd-field-row">
        <span className="cd-field-label">{label}:</span>
        {children || (
          <span className={`cd-field-value ${isCode ? 'code' : ''} ${isEmpty ? 'empty' : ''}`}>
            {displayValue}
          </span>
        )}
      </div>
    );
  };

  // 多行字段 - 自适应高度
  const MultilineField = ({ label, value, isCode = false }) => {
    const displayValue = value || '-';
    const isEmpty = !value || value === '-';
    if (isEmpty) {
      return (
        <div className="cd-field-row">
          <span className="cd-field-label">{label}:</span>
          <span className="cd-field-value empty">-</span>
        </div>
      );
    }
    return (
      <div className="cd-multiline-field">
        <span className="cd-field-label">{label}:</span>
        <div className={`cd-multiline-value ${isCode ? 'code' : ''}`}>
          {displayValue}
        </div>
      </div>
    );
  };

  return (
    <Modal
      title={
        <div className="cd-modal-title">
          <span>{t('testExecution.caseDetail.title')}</span>
          <Tag color="blue" style={{ marginLeft: 8 }}>
            {caseData.case_number || caseData.case_num || `No.${caseData.no}`}
          </Tag>
          <Tag color={isManual ? 'green' : isAPI ? 'orange' : 'blue'}>
            {isManual ? 'Manual' : isAPI ? 'API' : 'Web'}
          </Tag>
          {!isAPI && <span className="cd-lang-tag">{languageDisplay}</span>}
        </div>
      }
      open={visible}
      onCancel={onCancel}
      width={750}
      className="cd-modal-compact"
      footer={
        <div className="cd-modal-footer">
          <Button size="small" icon={<CloseOutlined />} onClick={onCancel}>
            {t('testExecution.caseDetail.cancel')}
          </Button>
          <Button type="primary" size="small" icon={<SaveOutlined />} loading={saving} onClick={handleSave}>
            {t('testExecution.caseDetail.saveAndClose')}
          </Button>
        </div>
      }
      destroyOnClose
    >
      <div className="cd-content">
        {/* Web用例字段 */}
        {isWeb && (
          <>
            <Row gutter={16}>
              <Col span={12}><FieldRow label="Screen" value={getFieldValue('screen')} /></Col>
              <Col span={12}><FieldRow label="Function" value={getFieldValue('function')} /></Col>
            </Row>
            <MultilineField label="Precondition" value={getFieldValue('precondition')} />
            <MultilineField label="Test Steps" value={getFieldValue('test_steps')} />
            <MultilineField label="Expected Result" value={getFieldValue('expected_result')} />
            <MultilineField label="ScriptCode" value={caseData.script_code} isCode />
          </>
        )}

        {/* API用例字段 */}
        {isAPI && (
          <>
            <Row gutter={16}>
              <Col span={8}><FieldRow label="Screen" value={caseData.screen} /></Col>
              <Col span={8}>
                <FieldRow label="Method">
                  <Tag color={caseData.method === 'GET' ? 'green' : caseData.method === 'POST' ? 'blue' : 'orange'}>
                    {caseData.method || '-'}
                  </Tag>
                </FieldRow>
              </Col>
              <Col span={8}>
                <FieldRow label="ResponseTime">
                  <Tag color="purple">
                    {caseData.response_time ? `${caseData.response_time} ms` : '-'}
                  </Tag>
                </FieldRow>
              </Col>
            </Row>
            <FieldRow label="URL" value={caseData.url} isCode />
            <MultilineField label="Header" value={caseData.header} isCode />
            <MultilineField label="Body" value={caseData.body} isCode />
            <MultilineField label="Response" value={caseData.response} isCode />
            <MultilineField label="ScriptCode" value={caseData.script_code} isCode />
          </>
        )}

        {/* Manual用例字段 */}
        {isManual && (
          <>
            <Row gutter={16}>
              <Col span={8}><FieldRow label="Maj.Category" value={getFieldValue('major_function')} /></Col>
              <Col span={8}><FieldRow label="Mid.Category" value={getFieldValue('middle_function')} /></Col>
              <Col span={8}><FieldRow label="Min.Category" value={getFieldValue('minor_function')} /></Col>
            </Row>
            <MultilineField label="Precondition" value={getFieldValue('precondition')} />
            <MultilineField label="Test Steps" value={getFieldValue('test_steps')} />
            <MultilineField label="Expected Result" value={getFieldValue('expected_result')} />
          </>
        )}

        {/* 执行信息 - 可编辑 */}
        <div className="cd-exec-section">
          <div className="cd-exec-title">
            ✏️ {t('testExecution.caseDetail.executionInfo')}
            <span className="cd-edit-hint">({t('testExecution.caseDetail.editable')})</span>
          </div>
          <Form form={form} className="cd-exec-form">
            <Row gutter={12}>
              <Col span={6}>
                <Form.Item
                  name="test_result"
                  label={t('testExecution.caseDetail.testResult')}
                  rules={[{ required: true }]}
                  style={{ marginBottom: 8 }}
                >
                  <Select size="small" style={{ width: '100%' }}>
                    <Option value="NR"><Tag color="default">NR</Tag></Option>
                    <Option value="OK"><Tag color="success">OK</Tag></Option>
                    <Option value="NG"><Tag color="error">NG</Tag></Option>
                    <Option value="Block"><Tag color="warning">Block</Tag></Option>
                  </Select>
                </Form.Item>
              </Col>
              <Col span={18}>
                <Form.Item name="bug_id" label={t('testExecution.caseDetail.bugId')} style={{ marginBottom: 8 }}>
                  <Input size="small" placeholder={t('testExecution.caseDetail.bugIdPlaceholder')} />
                </Form.Item>
              </Col>
            </Row>
            <Form.Item name="remark" label={t('testExecution.caseDetail.remark')} style={{ marginBottom: 0 }}>
              <TextArea
                autoSize={{ minRows: 1, maxRows: 3 }}
                placeholder={t('testExecution.caseDetail.remarkPlaceholder')}
                maxLength={500}
                showCount
              />
            </Form.Item>
          </Form>
        </div>
      </div>
    </Modal>
  );
};

CaseDetailModal.propTypes = {
  visible: PropTypes.bool.isRequired,
  caseData: PropTypes.object,
  executionType: PropTypes.oneOf(['automation', 'api', 'manual']),
  languageSuffix: PropTypes.string,
  languageDisplay: PropTypes.string,
  onSave: PropTypes.func.isRequired,
  onCancel: PropTypes.func.isRequired,
};

CaseDetailModal.defaultProps = {
  languageSuffix: '_cn',
  languageDisplay: 'CN',
};

export default CaseDetailModal;
