import React, { useState, useEffect } from 'react';
import { Modal, Table, Input, Button, Tag, Space, message, Tooltip, Empty, Popconfirm } from 'antd';
import {
  SettingOutlined,
  EditOutlined,
  SaveOutlined,
  CloseOutlined,
  PlusOutlined,
  DeleteOutlined
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import PropTypes from 'prop-types';
import './VariablesModal.css';

/**
 * 用户自定义变量Modal组件
 * 
 * 设计原则（与元数据完全独立）：
 * - 用户自定义变量用于脚本参数化
 * - 所有变量都可以添加、编辑、删除
 * - 变量值在脚本执行时替换 ${VAR_NAME} 占位符
 */
const VariablesModal = ({
  visible,
  onClose,
  groupId,
  groupName,
  groupType,          // web / api
  projectId,          // 项目ID
  variables = [],     // 用户自定义变量列表
  onSave,             // 保存回调
  readOnly = false    // 只读模式（执行任务中查看）
}) => {
  const { t } = useTranslation();
  const [dataSource, setDataSource] = useState([]);
  const [editingKey, setEditingKey] = useState('');
  const [editingRecord, setEditingRecord] = useState(null);

  // 初始化数据
  useEffect(() => {
    if (visible) {
      setDataSource(variables.map((v, idx) => ({ ...v, key: v.id || `var_${idx}` })));
      setEditingKey('');
      setEditingRecord(null);
    }
  }, [visible, variables]);

  // 开始编辑
  const handleEdit = (record) => {
    setEditingKey(record.key);
    setEditingRecord({ ...record });
  };

  // 取消编辑
  const handleCancelEdit = () => {
    // 如果是新增的且没有保存，则删除
    if (editingRecord?.isNew) {
      setDataSource(prev => prev.filter(item => item.key !== editingKey));
    }
    setEditingKey('');
    setEditingRecord(null);
  };

  // 保存单行编辑
  const handleSaveRow = () => {
    if (!editingRecord.var_key) {
      message.warning(t('variables.keyRequired', '请输入变量键名'));
      return;
    }

    // 检查键名是否重复
    const isDuplicate = dataSource.some(
      item => item.key !== editingKey && item.var_key === editingRecord.var_key
    );
    if (isDuplicate) {
      message.warning(t('variables.keyDuplicate', '变量键名已存在'));
      return;
    }

    setDataSource(prev => prev.map(item =>
      item.key === editingKey
        ? {
          ...editingRecord,
          var_name: `\${${editingRecord.var_key.toUpperCase()}}`,
          isNew: false
        }
        : item
    ));
    setEditingKey('');
    setEditingRecord(null);
  };

  // 删除变量
  const handleDelete = (key) => {
    setDataSource(prev => prev.filter(item => item.key !== key));
  };

  // 添加新变量
  const handleAdd = () => {
    const newKey = `new_${Date.now()}`;
    const newVar = {
      key: newKey,
      var_key: '',
      var_name: '',
      var_desc: '',
      var_value: '',
      var_type: 'custom',
      isNew: true
    };
    setDataSource(prev => [...prev, newVar]);
    setEditingKey(newKey);
    setEditingRecord(newVar);
  };

  // 复制变量名
  const handleCopyVarName = (name) => {
    if (!name) return;
    navigator.clipboard.writeText(name).then(() => {
      message.success(t('variables.copySuccess', '已复制到剪贴板'));
    });
  };

  // 保存所有变量到数据库并关闭弹窗
  const handleSaveAll = async () => {
    if (editingKey) {
      message.warning(t('variables.finishEditing', '请先完成当前编辑'));
      return;
    }

    // 过滤掉空的变量
    const validVars = dataSource.filter(v => v.var_key);

    if (onSave) {
      try {
        await onSave(validVars.map(({ key, isNew, ...rest }) => rest));
        // 保存成功后关闭弹窗
        onClose();
      } catch (error) {
        // 错误处理由父组件的onSave函数负责显示
        console.error('[VariablesModal] Save failed:', error);
      }
    }
  };

  // 编辑字段变更
  const handleFieldChange = (field, value) => {
    setEditingRecord(prev => ({ ...prev, [field]: value }));
  };

  // 表格列定义
  const columns = [
    {
      title: 'No.',
      width: 45,
      align: 'center',
      render: (_, __, index) => (
        <span style={{ color: '#999', fontSize: '12px' }}>{index + 1}</span>
      )
    },
    {
      title: t('variables.varName', '变量名'),
      dataIndex: 'var_name',
      width: 140,
      render: (text, record) => {
        if (editingKey === record.key) {
          return (
            <Input
              size="small"
              value={editingRecord?.var_key || ''}
              onChange={(e) => handleFieldChange('var_key', e.target.value.toLowerCase().replace(/[^a-z0-9_]/g, ''))}
              placeholder="var_key"
              addonBefore="${"
              addonAfter="}"
              style={{ fontSize: '12px', width: '130px' }}
            />
          );
        }
        return text ? (
          <Tooltip title={t('variables.clickToCopy', '点击复制')}>
            <code
              style={{
                color: '#1890ff',
                cursor: 'pointer',
                fontSize: '12px',
                padding: '2px 6px',
                background: '#f0f5ff',
                borderRadius: '3px'
              }}
              onClick={() => handleCopyVarName(text)}
            >
              {text}
            </code>
          </Tooltip>
        ) : '-';
      }
    },
    {
      title: t('variables.description', '描述'),
      dataIndex: 'var_desc',
      width: 160,
      ellipsis: true,
      render: (text, record) => {
        if (editingKey === record.key) {
          return (
            <Input
              size="small"
              value={editingRecord?.var_desc || ''}
              onChange={(e) => handleFieldChange('var_desc', e.target.value)}
              placeholder={t('variables.descPlaceholder', '变量描述')}
              style={{ fontSize: '12px' }}
            />
          );
        }
        return <span style={{ fontSize: '12px', color: '#666' }}>{text || '-'}</span>;
      }
    },
    {
      title: t('variables.value', '值'),
      dataIndex: 'var_value',
      width: 200,
      render: (text, record) => {
        if (editingKey === record.key) {
          return (
            <Input
              size="small"
              value={editingRecord?.var_value || ''}
              onChange={(e) => handleFieldChange('var_value', e.target.value)}
              placeholder={t('variables.valuePlaceholder', '变量值')}
              style={{ fontSize: '12px' }}
            />
          );
        }

        // 密码类型显示掩码
        const displayVal = record.var_key === 'password' ? '••••••••' : text;
        return (
          <Tooltip title={record.var_key === 'password' ? t('variables.passwordHidden', '密码已隐藏') : text}>
            <span style={{
              fontSize: '12px',
              color: '#333',
              fontFamily: 'monospace',
              maxWidth: '180px',
              display: 'inline-block',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap'
            }}>
              {displayVal || '-'}
            </span>
          </Tooltip>
        );
      }
    },
    // 类型列已移除 - 根据需求，变量表中不需要显示类型列
  ];

  // 非只读模式添加操作列
  if (!readOnly) {
    columns.push({
      title: t('common.action', '操作'),
      width: 80,
      align: 'center',
      render: (_, record) => {
        if (editingKey === record.key) {
          return (
            <Space size={4}>
              <Button type="link" size="small" onClick={handleSaveRow} style={{ padding: 0, fontSize: '12px' }}>
                <SaveOutlined />
              </Button>
              <Button type="link" size="small" onClick={handleCancelEdit} style={{ padding: 0, fontSize: '12px' }}>
                <CloseOutlined />
              </Button>
            </Space>
          );
        }
        return (
          <Space size={4}>
            <Button
              type="link"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
              style={{ padding: 0 }}
              disabled={editingKey !== ''}
            />
            <Popconfirm
              title={t('variables.deleteConfirm', '确定删除此变量？')}
              onConfirm={() => handleDelete(record.key)}
              okText={t('common.ok', '确定')}
              cancelText={t('common.cancel', '取消')}
            >
              <Button
                type="link"
                size="small"
                danger
                icon={<DeleteOutlined />}
                style={{ padding: 0 }}
                disabled={editingKey !== ''}
              />
            </Popconfirm>
          </Space>
        );
      }
    });
  }

  return (
    <Modal
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <SettingOutlined style={{ color: '#1890ff' }} />
          <span>{t('variables.title', '用户自定义变量')}</span>
          {groupName && (
            <Tag color="processing" style={{ marginLeft: '8px', fontSize: '11px' }}>
              {groupName}
            </Tag>
          )}
        </div>
      }
      open={visible}
      onCancel={onClose}
      width={780}
      footer={
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span style={{ fontSize: '12px', color: '#999' }}>
            {t('variables.hint', '提示：在脚本中使用 ${变量名} 引用这些值，与元数据独立')}
          </span>
          <Space>
            {!readOnly && (
              <Button size="small" type="primary" icon={<SaveOutlined />} onClick={handleSaveAll}>
                {t('common.save', '保存')}
              </Button>
            )}
            <Button size="small" onClick={onClose}>
              {t('common.close', '关闭')}
            </Button>
          </Space>
        </div>
      }
      bodyStyle={{ padding: '12px 16px' }}
      className="variables-modal"
    >
      {/* 添加变量按钮 */}
      {!readOnly && (
        <div style={{ marginBottom: '12px', textAlign: 'right' }}>
          <Button
            type="dashed"
            size="small"
            icon={<PlusOutlined />}
            onClick={handleAdd}
            disabled={editingKey !== ''}
          >
            {t('variables.addVariable', '添加变量')}
          </Button>
        </div>
      )}

      {dataSource.length === 0 ? (
        <Empty
          description={t('variables.noVariables', '暂无变量，点击上方按钮添加')}
          style={{ padding: '40px 0' }}
        />
      ) : (
        <Table
          columns={columns}
          dataSource={dataSource}
          rowKey="key"
          size="small"
          pagination={false}
          scroll={{ y: 350 }}
        />
      )}
    </Modal>
  );
};

VariablesModal.propTypes = {
  visible: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  groupId: PropTypes.number,
  groupName: PropTypes.string,
  groupType: PropTypes.string,
  projectId: PropTypes.number,
  variables: PropTypes.array,
  onSave: PropTypes.func,
  readOnly: PropTypes.bool
};

export default VariablesModal;
