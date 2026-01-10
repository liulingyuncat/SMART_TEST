import React, { useState } from 'react';
import { Modal, Radio, Input, Space, message, Upload } from 'antd';
import { UploadOutlined, PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

/**
 * 导入用例对话框
 * 支持选择已有用例集或创建新用例集
 */
const ImportCaseModal = ({ 
  visible, 
  onCancel, 
  onImport,
  caseGroups = [], // 已有用例集列表
  loading = false 
}) => {
  const { t } = useTranslation();
  const [selectedOption, setSelectedOption] = useState(''); // 选中的用例集或'new'
  const [newCaseGroupName, setNewCaseGroupName] = useState('');
  const [file, setFile] = useState(null);

  // 重置状态
  const handleCancel = () => {
    setSelectedOption('');
    setNewCaseGroupName('');
    setFile(null);
    onCancel();
  };

  // 检查用例集名称是否重复
  const isDuplicate = (name) => {
    return caseGroups.some(group => group === name.trim());
  };

  // 确认导入
  const handleOk = () => {
    // 验证是否选择了用例集
    if (!selectedOption) {
      message.warning('请选择要导入到的用例集');
      return;
    }

    // 如果是新建用例集，验证名称
    if (selectedOption === 'new') {
      const trimmedName = newCaseGroupName.trim();
      if (!trimmedName) {
        message.warning('请输入新用例集名称');
        return;
      }
      if (isDuplicate(trimmedName)) {
        message.error('用例集名称已存在，请使用不同的名称');
        return;
      }
    }

    // 验证是否选择了文件
    if (!file) {
      message.warning('请选择要导入的Excel文件');
      return;
    }

    // 确定最终的用例集名称
    const targetCaseGroup = selectedOption === 'new' ? newCaseGroupName.trim() : selectedOption;

    // 执行导入
    onImport(file, targetCaseGroup);
    
    // 不在这里重置状态，等导入成功后由父组件控制
  };

  // 文件选择配置
  const uploadProps = {
    beforeUpload: (file) => {
      // 验证文件类型
      const isExcel = file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
                      file.type === 'application/vnd.ms-excel' ||
                      file.name.endsWith('.xlsx') ||
                      file.name.endsWith('.xls');
      
      if (!isExcel) {
        message.error('只能上传Excel文件！');
        return false;
      }

      // 验证文件大小 (10MB)
      const isLt10M = file.size / 1024 / 1024 < 10;
      if (!isLt10M) {
        message.error('文件大小不能超过10MB！');
        return false;
      }

      setFile(file);
      return false; // 阻止自动上传
    },
    onRemove: () => {
      setFile(null);
    },
    fileList: file ? [file] : [],
    maxCount: 1,
  };

  return (
    <Modal
      title="导入用例"
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
      okText="确认导入"
      cancelText="取消"
      width={500}
      maskClosable={false}
    >
      <Space direction="vertical" style={{ width: '100%' }} size="middle">
        {/* 文件选择 */}
        <div>
          <div style={{ marginBottom: 8, fontWeight: 500 }}>选择Excel文件:</div>
          <Upload {...uploadProps}>
            <button 
              type="button" 
              style={{
                padding: '8px 16px',
                border: '1px solid #d9d9d9',
                borderRadius: '4px',
                background: '#fff',
                cursor: 'pointer',
                display: 'flex',
                alignItems: 'center',
                gap: '8px'
              }}
            >
              <UploadOutlined /> 选择文件
            </button>
          </Upload>
        </div>

        {/* 用例集选择 */}
        <div>
          <div style={{ marginBottom: 8, fontWeight: 500 }}>选择目标用例集:</div>
          <Radio.Group 
            value={selectedOption} 
            onChange={(e) => setSelectedOption(e.target.value)}
            style={{ width: '100%' }}
          >
            <Space direction="vertical" style={{ width: '100%' }}>
              {/* 已有用例集列表 */}
              {caseGroups.map(group => (
                <Radio key={group} value={group}>
                  {group}
                </Radio>
              ))}
              
              {/* 新建用例集选项 */}
              <Radio value="new">
                <Space>
                  <PlusOutlined />
                  <span>新建用例集</span>
                </Space>
              </Radio>
            </Space>
          </Radio.Group>
        </div>

        {/* 新建用例集输入框 */}
        {selectedOption === 'new' && (
          <div style={{ marginLeft: 24 }}>
            <Input
              placeholder="请输入新用例集名称"
              value={newCaseGroupName}
              onChange={(e) => setNewCaseGroupName(e.target.value)}
              maxLength={100}
              status={newCaseGroupName.trim() && isDuplicate(newCaseGroupName) ? 'error' : ''}
              onPressEnter={handleOk}
            />
            {newCaseGroupName.trim() && isDuplicate(newCaseGroupName) && (
              <div style={{ color: '#ff4d4f', marginTop: 4, fontSize: '12px' }}>
                用例集名称已存在
              </div>
            )}
          </div>
        )}

        {/* 提示信息 */}
        <div style={{ 
          padding: '8px 12px', 
          background: '#e6f7ff', 
          border: '1px solid #91d5ff',
          borderRadius: '4px',
          fontSize: '12px',
          color: '#096dd9'
        }}>
          <div>💡 提示：</div>
          <div>• 导入的用例将添加到选中的用例集中</div>
          <div>• Excel文件会自动读取第一个工作表</div>
          <div>• 如果用例UUID已存在，将更新该用例</div>
        </div>
      </Space>
    </Modal>
  );
};

export default ImportCaseModal;
