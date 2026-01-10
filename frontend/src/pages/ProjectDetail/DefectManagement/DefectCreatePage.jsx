import React, { useState, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Form, Input, Select, Button, Card, Row, Col, message, Space, Upload } from 'antd';
import { SaveOutlined, UploadOutlined } from '@ant-design/icons';
import {
  DEFECT_STATUS,
  DEFECT_STATUS_COLORS,
  DEFECT_PRIORITY,
  DEFECT_PRIORITY_COLORS,
  DEFECT_SEVERITY,
  DEFECT_SEVERITY_COLORS,
} from '../../../constants/defect';
import { createDefect, uploadDefectAttachment } from '../../../api/defect';

const { TextArea } = Input;
const { Option } = Select;

/**
 * 新增缺陷页面
 * 全屏表单页面，填写缺陷详情并保存
 */
const DefectCreatePage = ({
  projectId,
  subjects,
  phases,
  onCancel,
  onSuccess,
}) => {
  const { t, i18n } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [fileList, setFileList] = useState([]);

  // Description默认模板（从localStorage获取，如果没有则使用默认值）
  const getDefaultDescriptionTemplate = () => {
    return localStorage.getItem('defect_description_template') || `[Actual result]

[Relevant validation]

[Test Steps]

→issue occurred

[Expected result]

[Test Environment]

`;
  };

  // 使用 useMemo 缓存翻译标签，只在语言变化时重新计算
  const labels = useMemo(() => ({
    create: t('defect.create', '新建缺陷'),
    save: t('common.save', '保存'),
    cancel: t('common.cancel', '取消'),
    title: t('defect.title', 'Title'),
    titlePlaceholder: t('defect.titlePlaceholder', '请输入缺陷标题'),
    subject: t('defect.subject', 'Subject'),
    description: t('defect.description', 'Description'),
    descriptionPlaceholder: t('defect.descriptionPlaceholder', '请输入详细描述'),
    recoveryMethod: t('defect.recoveryMethod', 'Recovery Method'),
    recoveryMethodPlaceholder: t('defect.recoveryMethodPlaceholder', '请输入恢复方法'),
    priority: t('defect.priority', 'Priority'),
    severity: t('defect.severity', 'Severity'),
    frequency: t('defect.frequency', 'Frequency(%)'),
    frequencyPlaceholder: t('defect.frequencyPlaceholder', '请输入复现频率百分比'),
    detectedInRelease: t('defect.detectedInRelease', 'Detected in Release'),
    detectedInReleasePlaceholder: t('defect.detectedInReleasePlaceholder', '请输入发现版本'),
    caseId: t('defect.caseId', 'Case ID'),
    caseIdPlaceholder: t('defect.caseIdPlaceholder', '请输入关联的Case ID'),
    phase: t('defect.phase', 'Phase'),
    status: t('defect.status', 'Status'),
    attachments: t('defect.attachments', 'Attachments'),
    upload: t('common.upload', '上传'),
    pleaseSelect: t('common.pleaseSelect', '请选择'),
    required: t('validation.required', '此字段为必填项'),
    createSuccess: t('defect.createSuccess', '缺陷创建成功'),
    saveFailed: t('message.saveFailed', '保存失败'),
    maxFileSize: t('defect.maxFileSize', '单个文件最大100MB'),
  }), [t, i18n.language]);

  // 状态选项
  const statusOptions = useMemo(() => [
    { value: DEFECT_STATUS.NEW, label: t('defect.statusNew', '新建'), color: DEFECT_STATUS_COLORS[DEFECT_STATUS.NEW] },
    { value: DEFECT_STATUS.ACTIVE, label: t('defect.statusActive', '处理中'), color: DEFECT_STATUS_COLORS[DEFECT_STATUS.ACTIVE] },
    { value: DEFECT_STATUS.RESOLVED, label: t('defect.statusResolved', '已解决'), color: DEFECT_STATUS_COLORS[DEFECT_STATUS.RESOLVED] },
    { value: DEFECT_STATUS.CLOSED, label: t('defect.statusClosed', '已关闭'), color: DEFECT_STATUS_COLORS[DEFECT_STATUS.CLOSED] },
  ], [t, i18n.language]);

  // 优先级选项（按需求文档FR-03：A/B/C/D）
  const priorityOptions = useMemo(() => [
    { value: DEFECT_PRIORITY.A, label: t('defect.priorityA', 'A'), color: DEFECT_PRIORITY_COLORS[DEFECT_PRIORITY.A] },
    { value: DEFECT_PRIORITY.B, label: t('defect.priorityB', 'B'), color: DEFECT_PRIORITY_COLORS[DEFECT_PRIORITY.B] },
    { value: DEFECT_PRIORITY.C, label: t('defect.priorityC', 'C'), color: DEFECT_PRIORITY_COLORS[DEFECT_PRIORITY.C] },
    { value: DEFECT_PRIORITY.D, label: t('defect.priorityD', 'D'), color: DEFECT_PRIORITY_COLORS[DEFECT_PRIORITY.D] },
  ], [t, i18n.language]);

  // 严重程度选项（按需求文档FR-03：A/B/C/D）
  const severityOptions = useMemo(() => [
    { value: DEFECT_SEVERITY.A, label: t('defect.severityA', 'A'), color: DEFECT_SEVERITY_COLORS[DEFECT_SEVERITY.A] },
    { value: DEFECT_SEVERITY.B, label: t('defect.severityB', 'B'), color: DEFECT_SEVERITY_COLORS[DEFECT_SEVERITY.B] },
    { value: DEFECT_SEVERITY.C, label: t('defect.severityC', 'C'), color: DEFECT_SEVERITY_COLORS[DEFECT_SEVERITY.C] },
    { value: DEFECT_SEVERITY.D, label: t('defect.severityD', 'D'), color: DEFECT_SEVERITY_COLORS[DEFECT_SEVERITY.D] },
  ], [t, i18n.language]);

  // 提交表单
  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      // 过滤掉undefined值，避免后端处理错误
      const filteredValues = Object.fromEntries(
        Object.entries(values).filter(([_, v]) => v !== undefined && v !== null && v !== '')
      );
      const defectData = {
        ...filteredValues,
        project_id: projectId,
      };
      console.log('[DEBUG] createDefect: submitting', defectData);
      const response = await createDefect(projectId, defectData);
      console.log('[DEBUG] createDefect: API response', response);
      
      // apiClient已提取data字段，response可能是 { defect: {...} } 或直接是 defect对象
      const createdDefect = response?.defect || response;
      const defectId = createdDefect?.id;
      console.log('[DEBUG] createDefect: createdDefect', createdDefect, 'defectId', defectId);

      // 上传附件
      if (defectId && fileList.length > 0) {
        for (const file of fileList) {
          await uploadDefectAttachment(projectId, defectId, file.originFileObj);
        }
      }

      message.success(labels.createSuccess);
      // 传递创建的缺陷对象给父组件，使用 defect_id 作为显示ID
      onSuccess?.(createdDefect);
    } catch (error) {
      console.error('[DEBUG] createDefect: Failed to create defect:', error);
      const errorMessage = error.response?.data?.message || labels.saveFailed;
      message.error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // 文件上传前处理
  const handleBeforeUpload = (file) => {
    setFileList((prev) => [...prev, { uid: file.uid, name: file.name, originFileObj: file }]);
    return false;
  };

  // 移除文件
  const handleRemoveFile = (file) => {
    setFileList((prev) => prev.filter((f) => f.uid !== file.uid));
  };

  return (
    <div className="defect-create-page">
      <Card
        title={labels.create}
        extra={
          <Space>
            <Button onClick={onCancel} style={{ border: '1px solid #d9d9d9' }}>
              {labels.cancel}
            </Button>
            <Button
              type="primary"
              icon={<SaveOutlined />}
              loading={loading}
              onClick={() => form.submit()}
            >
              {labels.save}
            </Button>
          </Space>
        }
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{
            status: DEFECT_STATUS.NEW,
            priority: DEFECT_PRIORITY.B,
            severity: DEFECT_SEVERITY.B,
            description: getDefaultDescriptionTemplate(),
          }}
        >
          <Row gutter={24}>
            {/* 左侧: 基本信息 */}
            <Col span={16}>
              {/* Title - 必填 */}
              <Form.Item
                name="title"
                label={labels.title}
                rules={[{ required: true, message: labels.required }]}
              >
                <Input placeholder={labels.titlePlaceholder} maxLength={200} />
              </Form.Item>

              {/* Subject - 可选 */}
              <Form.Item name="subject_id" label={labels.subject}>
                <Select allowClear placeholder={labels.pleaseSelect}>
                  {subjects?.map((subject) => (
                    <Option key={subject.id} value={subject.id}>
                      {subject.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>

              {/* Description - 可选，带默认模板 */}
              <Form.Item name="description" label={labels.description}>
                <TextArea rows={12} placeholder={labels.descriptionPlaceholder} />
              </Form.Item>

              {/* Recovery Method - 可选 */}
              <Form.Item name="recovery_method" label={labels.recoveryMethod}>
                <Input placeholder={labels.recoveryMethodPlaceholder} maxLength={500} />
              </Form.Item>
            </Col>

            {/* 右侧: 属性信息 */}
            <Col span={8}>
              {/* Priority - 必填，默认B */}
              <Form.Item 
                name="priority" 
                label={labels.priority}
                rules={[{ required: true, message: labels.required }]}
              >
                <Select>
                  {priorityOptions.map((opt) => (
                    <Option key={opt.value} value={opt.value}>
                      <span style={{ color: opt.color }}>{opt.label}</span>
                    </Option>
                  ))}
                </Select>
              </Form.Item>

              {/* Severity - 必填，默认B */}
              <Form.Item 
                name="severity" 
                label={labels.severity}
                rules={[{ required: true, message: labels.required }]}
              >
                <Select>
                  {severityOptions.map((opt) => (
                    <Option key={opt.value} value={opt.value}>
                      <span style={{ color: opt.color }}>{opt.label}</span>
                    </Option>
                  ))}
                </Select>
              </Form.Item>

              {/* Frequency(%) - 必填 */}
              <Form.Item 
                name="frequency" 
                label={labels.frequency}
                rules={[{ required: true, message: labels.required }]}
              >
                <Input placeholder={labels.frequencyPlaceholder} maxLength={10} />
              </Form.Item>

              {/* Detected in Release - 必填 */}
              <Form.Item 
                name="detected_in_release" 
                label={labels.detectedInRelease}
                rules={[{ required: true, message: labels.required }]}
              >
                <Input placeholder={labels.detectedInReleasePlaceholder} maxLength={50} />
              </Form.Item>

              {/* Phase - 可选 */}
              <Form.Item name="phase_id" label={labels.phase}>
                <Select allowClear placeholder={labels.pleaseSelect}>
                  {phases?.map((phase) => (
                    <Option key={phase.id} value={phase.id}>
                      {phase.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>

              {/* Case ID - 可选 */}
              <Form.Item name="case_id" label={labels.caseId}>
                <Input placeholder={labels.caseIdPlaceholder} maxLength={50} />
              </Form.Item>

              {/* Attachments - 可选，最大100MB */}
              <Form.Item label={labels.attachments} extra={labels.maxFileSize}>
                <Upload
                  fileList={fileList}
                  beforeUpload={handleBeforeUpload}
                  onRemove={handleRemoveFile}
                  multiple
                  maxCount={10}
                >
                  <Button icon={<UploadOutlined />}>{labels.upload}</Button>
                </Upload>
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Card>
    </div>
  );
};

export default DefectCreatePage;
