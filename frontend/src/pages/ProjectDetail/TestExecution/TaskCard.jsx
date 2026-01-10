import React from 'react';
import { Button, App, Typography, message } from 'antd';
import { DeleteOutlined, ExclamationCircleOutlined, CopyOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import PropTypes from 'prop-types';
import './TaskCard.css';

const { Text } = Typography;

const TaskCard = ({ task, selected, onSelect, onDelete }) => {
  const { t } = useTranslation();
  const { modal } = App.useApp();

  const handleCopy = (e) => {
    e.stopPropagation();
    e.preventDefault();
    navigator.clipboard.writeText(task.task_name);
    message.success(t('common.copySuccess', { defaultValue: '复制成功' }));
  };

  const handleDelete = (e) => {
    console.log('\ud83d\uddd1\ufe0f [TaskCard] handleDelete called');
    console.log('\ud83d\uddd1\ufe0f [TaskCard] Task to delete:', task);
    console.log('\ud83d\uddd1\ufe0f [TaskCard] onDelete callback exists:', !!onDelete);
    
    // 阻止事件冒泡，避免触发卡片的选中事件
    e.stopPropagation();
    e.preventDefault();
    console.log('\ud83d\uddd1\ufe0f [TaskCard] Event propagation stopped');
    
    if (!onDelete) {
      console.error('\u274c [TaskCard] onDelete callback is not defined!');
      return;
    }
    
    console.log('\ud83d\uddd1\ufe0f [TaskCard] Showing confirmation modal...');
    console.log('\ud83d\uddd1\ufe0f [TaskCard] modal instance:', modal);
    
    // 使用App.useApp()提供的modal实例
    modal.confirm({
      title: t('testExecution.taskCard.confirmDelete'),
      icon: <ExclamationCircleOutlined />,
      content: t('testExecution.taskCard.deleteMessage', { taskName: task.task_name }),
      okText: t('testExecution.taskCard.delete'),
      okType: 'danger',
      cancelText: t('testExecution.taskCard.cancel'),
      centered: true,
      onOk() {
        console.log('\u2705 [TaskCard] User confirmed deletion');
        console.log('\ud83d\udd04 [TaskCard] Calling onDelete with task_uuid:', task.task_uuid);
        onDelete(task.task_uuid);
      },
      onCancel() {
        console.log('\u274c [TaskCard] User cancelled deletion');
      }
    });
    
    console.log('\ud83d\uddd1\ufe0f [TaskCard] modal.confirm called');
  };

  return (
    <div
      className={`task-card ${selected ? 'task-card-selected' : ''}`}
      onClick={() => onSelect(task.task_uuid)}
    >
      <div className="task-card-content">
        <Text
          ellipsis={{ tooltip: task.task_name }}
          className="task-card-name"
        >
          {task.task_name}
        </Text>
        <div className="task-card-actions">
          <Button
            type="text"
            icon={<CopyOutlined />}
            size="small"
            className="task-card-copy"
            onClick={handleCopy}
            title={t('common.copy', { defaultValue: '复制' })}
          />
          <Button
            type="text"
            icon={<DeleteOutlined />}
            danger
            size="small"
            className="task-card-delete"
            onClick={handleDelete}
          />
        </div>
      </div>
    </div>
  );
};

TaskCard.propTypes = {
  task: PropTypes.shape({
    task_uuid: PropTypes.string.isRequired,
    task_name: PropTypes.string.isRequired,
  }).isRequired,
  selected: PropTypes.bool.isRequired,
  onSelect: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
};

export default TaskCard;
