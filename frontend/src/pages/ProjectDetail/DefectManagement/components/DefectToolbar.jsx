import React, { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Space, Button, Upload, Dropdown, Menu } from 'antd';
import {
  PlusOutlined,
  DownloadOutlined,
  UploadOutlined,
  ExportOutlined,
  SettingOutlined,
} from '@ant-design/icons';

/**
 * 缺陷列表工具栏
 * 包含：新增、导入、导出、下载模板、配置等操作
 */
const DefectToolbar = ({
  onCreate,
  onDownloadTemplate,
  onImport,
  onExport,
  onEditSubjects,
  onEditPhases,
}) => {
  const { t, i18n } = useTranslation();

  // 使用 useMemo 缓存翻译标签，只在语言变化时重新计算
  const labels = useMemo(() => ({
    create: t('defect.create', '新建缺陷'),
    import: t('common.import', '导入'),
    export: t('common.export', '导出'),
    downloadTemplate: t('defect.downloadTemplate', '下载模板'),
    settings: t('common.settings', '配置'),
    subjects: t('defect.subjectsConfig', '主题配置'),
    phases: t('defect.phasesConfig', '阶段配置'),
  }), [t, i18n.language]);

  // 配置菜单
  const configMenu = (
    <Menu>
      <Menu.Item key="subjects" onClick={onEditSubjects}>
        {labels.subjects}
      </Menu.Item>
      <Menu.Item key="phases" onClick={onEditPhases}>
        {labels.phases}
      </Menu.Item>
    </Menu>
  );

  return (
    <div className="defect-toolbar" style={{ marginBottom: 16 }}>
      <Space>
        <Button type="primary" icon={<PlusOutlined />} onClick={onCreate}>
          {labels.create}
        </Button>

        <Upload
          accept=".csv"
          showUploadList={false}
          beforeUpload={(file) => {
            onImport?.(file);
            return false;
          }}
        >
          <Button icon={<UploadOutlined />}>{labels.import}</Button>
        </Upload>

        <Button icon={<ExportOutlined />} onClick={onExport}>
          {labels.export}
        </Button>

        <Button icon={<DownloadOutlined />} onClick={onDownloadTemplate}>
          {labels.downloadTemplate}
        </Button>

        <Dropdown overlay={configMenu} placement="bottomRight">
          <Button icon={<SettingOutlined />}>{labels.settings}</Button>
        </Dropdown>
      </Space>
    </div>
  );
};

export default DefectToolbar;
