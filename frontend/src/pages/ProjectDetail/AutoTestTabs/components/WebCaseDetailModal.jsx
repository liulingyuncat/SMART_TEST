import React, { useState, useEffect } from 'react';
import { Modal, Form, Input, Row, Col, Button, Tag, message, Spin, Alert } from 'antd';
import { SaveOutlined, CloseOutlined, PlayCircleOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import PropTypes from 'prop-types';
import { testScript } from '../../../../api/scriptTest';
import './WebCaseDetailModal.css';

const { TextArea } = Input;

/**
 * Web用例详细信息弹窗 - 全字段可编辑版
 * 所有字段均可编辑，包括 Script Code
 * 支持脚本测试功能
 * 设计风格与 ApiCaseDetailModal 保持一致
 */
const WebCaseDetailModal = ({
  visible,
  caseData,
  language = 'cn', // 当前语言: cn, jp, en
  projectId,       // 项目ID（用于脚本测试）
  groupId,         // 用例集ID（用于获取变量）
  onSave,
  onCancel,
}) => {
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState(null);
  const [hasScriptCode, setHasScriptCode] = useState(false);

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
      const scriptCode = caseData.script_code || '';
      form.setFieldsValue({
        screen: caseData[`screen${suffix}`] || '',
        function: caseData[`function${suffix}`] || '',
        precondition: caseData[`precondition${suffix}`] || '',
        test_steps: caseData[`test_steps${suffix}`] || '',
        expected_result: caseData[`expected_result${suffix}`] || '',
        script_code: scriptCode,
      });
      // 更新脚本代码状态
      setHasScriptCode(!!scriptCode.trim());
      // 重置测试结果
      setTestResult(null);
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

  // 持久化调试日志
  const saveDebugLog = (key, data) => {
    try {
      const logs = JSON.parse(localStorage.getItem('_script_test_logs') || '[]');
      logs.push({
        timestamp: new Date().toISOString(),
        component: 'WebCaseDetailModal',
        key,
        data: typeof data === 'object' ? JSON.stringify(data, null, 2) : String(data)
      });
      if (logs.length > 30) logs.shift();
      localStorage.setItem('_script_test_logs', JSON.stringify(logs));
    } catch (e) {
      console.error('[saveDebugLog] Failed:', e);
    }
  };

  // 脚本测试
  const handleTestScript = async () => {
    const values = form.getFieldsValue();
    const scriptCode = values.script_code;

    saveDebugLog('START', { projectId, groupId, hasScriptCode: !!scriptCode });

    if (!scriptCode || scriptCode.trim() === '') {
      message.warning('没有脚本代码可执行');
      return;
    }

    if (!projectId) {
      saveDebugLog('ERROR', 'projectId is missing');
      message.warning('项目ID不可用，无法执行测试');
      return;
    }

    setTesting(true);
    setTestResult(null);

    try {
      console.log('🧪 [WebCaseDetailModal] 开始脚本测试...');
      saveDebugLog('CALLING_API', {
        projectId,
        script_code_length: scriptCode.length,
        group_id: groupId || 0,
        group_type: 'web'
      });

      const result = await testScript(projectId, {
        script_code: scriptCode,
        group_id: groupId || 0,
        group_type: 'web',
      });

      console.log('🧪 [WebCaseDetailModal] 测试结果:', result);
      saveDebugLog('SUCCESS', result);
      setTestResult(result);

      if (result.success) {
        message.success(`脚本执行成功 (${result.response_time}ms)`);
      } else {
        message.warning('脚本执行失败');
      }
    } catch (error) {
      console.error('🧪 [WebCaseDetailModal] 脚本测试失败:', error);
      saveDebugLog('CATCH_ERROR', {
        message: error.message,
        response: error.response ? {
          status: error.response.status,
          statusText: error.response.statusText,
          data: error.response.data
        } : 'no response',
        stack: error.stack
      });

      setTestResult({
        success: false,
        error_message: error.message || '执行失败',
        response_time: 0,
      });
      message.error('脚本测试失败: ' + (error.message || '未知错误'));
    } finally {
      setTesting(false);
      saveDebugLog('FINISHED', { testing: false });
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
      width={900}
      className="wcd-modal-compact"
      footer={
        <div className="wcd-modal-footer" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          {/* 左侧：脚本测试按钮 */}
          <div>
            <Button
              icon={<PlayCircleOutlined />}
              onClick={handleTestScript}
              loading={testing}
              disabled={!hasScriptCode}
              style={{
                background: testing ? undefined : 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                borderColor: 'transparent',
                color: testing ? undefined : '#fff',
              }}
            >
              {testing ? '执行中...' : '脚本测试'}
            </Button>
          </div>
          {/* 右侧：取消和保存按钮 */}
          <div>
            <Button icon={<CloseOutlined />} onClick={onCancel} style={{ marginRight: 8 }}>
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
        <Form.Item name="script_code" label="Script Code" style={{ marginBottom: 12 }}>
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
            onChange={(e) => setHasScriptCode(!!e.target.value.trim())}
          />
        </Form.Item>

        {/* 测试结果显示区域 */}
        {testing && (
          <div style={{ textAlign: 'center', padding: '16px 0' }}>
            <Spin tip="正在Docker环境中执行脚本..." />
          </div>
        )}

        {testResult && !testing && (
          <div className="script-test-result" style={{ marginTop: 8 }}>
            <Alert
              type={testResult.success ? 'success' : 'error'}
              icon={testResult.success ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
              message={
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <span>{testResult.success ? '脚本执行成功' : '脚本执行失败'}</span>
                  <Tag color={testResult.success ? 'green' : 'red'}>
                    {testResult.response_time}ms
                  </Tag>
                </div>
              }
              description={
                testResult.error_message ? (
                  <pre style={{
                    margin: '8px 0 0 0',
                    padding: 8,
                    backgroundColor: '#f5f5f5',
                    borderRadius: 4,
                    fontSize: 11,
                    maxHeight: 100,
                    overflow: 'auto',
                    whiteSpace: 'pre-wrap',
                    wordBreak: 'break-word',
                  }}>
                    {testResult.error_message}
                  </pre>
                ) : (
                  testResult.output && (
                    <pre style={{
                      margin: '8px 0 0 0',
                      padding: 8,
                      backgroundColor: '#f5f5f5',
                      borderRadius: 4,
                      fontSize: 11,
                      maxHeight: 100,
                      overflow: 'auto',
                      whiteSpace: 'pre-wrap',
                      wordBreak: 'break-word',
                    }}>
                      {testResult.output}
                    </pre>
                  )
                )
              }
              showIcon
            />
          </div>
        )}
      </Form>
    </Modal>
  );
};

WebCaseDetailModal.propTypes = {
  visible: PropTypes.bool.isRequired,
  caseData: PropTypes.object,
  language: PropTypes.string,
  projectId: PropTypes.oneOfType([PropTypes.number, PropTypes.string]),
  groupId: PropTypes.oneOfType([PropTypes.number, PropTypes.string]),
  onSave: PropTypes.func.isRequired,
  onCancel: PropTypes.func.isRequired,
};

export default WebCaseDetailModal;
