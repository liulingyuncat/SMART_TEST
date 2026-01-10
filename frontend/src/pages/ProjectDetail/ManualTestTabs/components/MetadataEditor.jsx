import React, { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { Input, DatePicker, Space, Button, message } from 'antd';
import { SaveOutlined, ClearOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import './MetadataEditor.css';

/**
 * 元数据编辑组件
 * @param {Object} props
 * @param {Function} props.onSave - 保存回调函数
 * @param {Object} props.initialData - 初始数据
 */
const MetadataEditor = ({ onSave, initialData = {} }) => {
  const { t } = useTranslation();
  const [metadata, setMetadata] = useState({
    testVersion: '',
    testEnv: '',
    testDate: '',
    executor: '',
  });
  const [loading, setLoading] = useState(false);
  const hasUserEditedRef = useRef(false); // 标记用户是否编辑过

  // 当initialData变化时同步更新本地state
  // 注意：当initialData的值发生实质性变化时，应该更新组件状态
  useEffect(() => {
    if (initialData && Object.keys(initialData).length > 0) {
      const newMetadata = {
        testVersion: initialData.test_version || '',
        testEnv: initialData.test_env || '',
        testDate: initialData.test_date || '',
        executor: initialData.executor || '',
      };
      
      // 检查是否有实质性变化
      const hasChanged = 
        newMetadata.testVersion !== metadata.testVersion ||
        newMetadata.testEnv !== metadata.testEnv ||
        newMetadata.testDate !== metadata.testDate ||
        newMetadata.executor !== metadata.executor;
      
      if (hasChanged) {
        setMetadata(newMetadata);
        hasUserEditedRef.current = false; // 重置编辑标志
      }
    }
  }, [initialData]); // 移除metadata依赖，避免循环

  // 字段变更处理
  const handleChange = (field, value) => {
    hasUserEditedRef.current = true; // 标记用户已编辑
    setMetadata((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  // 日期变更处理
  const handleDateChange = (date, dateString) => {
    handleChange('testDate', dateString);
  };

  // 保存按钮点击
  const handleSave = async () => {
    setLoading(true);
    try {
      await onSave({
        test_version: metadata.testVersion,
        test_env: metadata.testEnv,
        test_date: metadata.testDate,
        executor: metadata.executor,
      });
      hasUserEditedRef.current = false; // 保存成功后重置编辑标志
      message.success(t('manualTest.saveSuccess'));
    } catch (error) {
      message.error(t('manualTest.saveFailed'));
    } finally {
      setLoading(false);
    }
  };

  // 清空按钮点击
  const handleClear = async () => {
    setLoading(true);
    try {
      // 清空本地状态
      setMetadata({
        testVersion: '',
        testEnv: '',
        testDate: '',
        executor: '',
      });
      // 保存空值到后端
      await onSave({
        test_version: '',
        test_env: '',
        test_date: '',
        executor: '',
      });
      hasUserEditedRef.current = false; // 清空后重置编辑标志
      message.success(t('manualTest.clearSuccess'));
    } catch (error) {
      message.error(t('manualTest.clearFailed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="metadata-editor">
      <Space size="middle" wrap align="center">
        <div className="metadata-field">
          <label>{t('manualTest.testVersion')}</label>
          <Input
            value={metadata.testVersion}
            onChange={(e) => handleChange('testVersion', e.target.value)}
            placeholder={t('manualTest.testVersion')}
            maxLength={50}
            style={{ width: 200 }}
          />
        </div>
        <div className="metadata-field">
          <label>{t('manualTest.testEnv')}</label>
          <Input
            value={metadata.testEnv}
            onChange={(e) => handleChange('testEnv', e.target.value)}
            placeholder={t('manualTest.testEnv')}
            maxLength={100}
            style={{ width: 200 }}
          />
        </div>
        <div className="metadata-field">
          <label>{t('manualTest.testDate')}</label>
          <DatePicker
            value={metadata.testDate ? dayjs(metadata.testDate) : null}
            onChange={handleDateChange}
            format="YYYY-MM-DD"
            placeholder={t('manualTest.testDate')}
            style={{ width: 200 }}
          />
        </div>
        <div className="metadata-field">
          <label>{t('manualTest.executor')}</label>
          <Input
            value={metadata.executor}
            onChange={(e) => handleChange('executor', e.target.value)}
            placeholder={t('manualTest.executor')}
            maxLength={50}
            style={{ width: 200 }}
          />
        </div>
        <Button 
          type="primary" 
          icon={<SaveOutlined />} 
          onClick={handleSave}
          loading={loading}
        >
          {t('common.save')}
        </Button>
        <Button 
          icon={<ClearOutlined />} 
          onClick={handleClear}
          loading={loading}
        >
          {t('common.clear')}
        </Button>
      </Space>
    </div>
  );
};

export default MetadataEditor;
