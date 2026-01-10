import React, { useMemo } from 'react';
import { Button, Typography, Empty } from 'antd';
import { PlusOutlined, MenuFoldOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import PropTypes from 'prop-types';
import TaskCard from './TaskCard';
import './TaskSidebar.css';

const { Text } = Typography;

const TaskSidebar = ({ tasks, selectedTaskUuid, onTaskSelect, onCollapse, onTaskCreate, onTaskDelete }) => {
  const { t } = useTranslation();

  const { inProgressTasks, pendingTasks, completedTasks } = useMemo(() => {
    console.log('📊 [TaskSidebar] Processing tasks:', {
      totalTasks: tasks.length,
      tasks: tasks.map(t => ({
        uuid: t.task_uuid,
        name: t.task_name,
        status: t.task_status,
        statusType: typeof t.task_status
      }))
    });
    
    const grouped = {
      inProgressTasks: tasks.filter(t => t.task_status === 'in_progress'),
      pendingTasks: tasks.filter(t => t.task_status === 'pending'),
      completedTasks: tasks.filter(t => t.task_status === 'completed'),
    };
    
    console.log('📊 [TaskSidebar] Task groups:', {
      inProgress: grouped.inProgressTasks.length,
      pending: grouped.pendingTasks.length,
      completed: grouped.completedTasks.length
    });
    
    return grouped;
  }, [tasks]);

  const getStatusStyle = (status) => {
    const styles = {
      inProgress: {
        color: '#cf1322',
        borderLeft: '3px solid #cf1322'
      },
      pending: {
        color: '#d48806',
        borderLeft: '3px solid #d48806'
      },
      completed: {
        color: '#389e0d',
        borderLeft: '3px solid #389e0d'
      }
    };
    return styles[status] || {};
  };

  const renderTaskGroup = (title, taskList, status) => {
    const statusStyle = getStatusStyle(status);
    return (
      <div className="task-group">
        <div className="task-group-title" style={{ borderLeft: statusStyle.borderLeft }}>
          <Text strong>
            {title}
          </Text>
          <Text style={{ color: statusStyle.color, marginLeft: '8px' }}>({taskList.length})</Text>
        </div>
      <div className="task-group-content">
        {taskList.length === 0 ? (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={t('testExecution.taskList.noTasks')}
            style={{ padding: '16px 0' }}
          />
        ) : (
          taskList.map(task => (
            <TaskCard
              key={task.task_uuid}
              task={task}
              selected={task.task_uuid === selectedTaskUuid}
              onSelect={onTaskSelect}
              onDelete={onTaskDelete}
            />
          ))
        )}
      </div>
    </div>
    );
  };

  return (
    <div className="task-sidebar">
      <div className="task-sidebar-header">
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={onTaskCreate}
        >
          {t('testExecution.createTask')}
        </Button>
        <Button
          type="text"
          icon={<MenuFoldOutlined />}
          size="small"
          onClick={onCollapse}
          style={{ marginTop: '4px' }}
        />
      </div>

      <div className="task-sidebar-body">
        {renderTaskGroup(t('testExecution.taskList.inProgress'), inProgressTasks, 'inProgress')}
        {renderTaskGroup(t('testExecution.taskList.pending'), pendingTasks, 'pending')}
        {renderTaskGroup(t('testExecution.taskList.completed'), completedTasks, 'completed')}
      </div>
    </div>
  );
};

TaskSidebar.propTypes = {
  tasks: PropTypes.array.isRequired,
  selectedTaskUuid: PropTypes.string,
  onTaskSelect: PropTypes.func.isRequired,
  onCollapse: PropTypes.func.isRequired,
  onTaskCreate: PropTypes.func.isRequired,
  onTaskDelete: PropTypes.func.isRequired,
};

export default TaskSidebar;
