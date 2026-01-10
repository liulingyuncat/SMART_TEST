import React, { useState, useEffect, useMemo } from 'react';
import { Spin, Alert, Button, message } from 'antd';
import { MenuUnfoldOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import PropTypes from 'prop-types';
import { getExecutionTasks, deleteExecutionTask } from '../../../api/executionTask';
import TaskSidebar from './TaskSidebar';
import TaskMetadataPanel from './TaskMetadataPanel';
import CreateTaskModal from './CreateTaskModal';

const TestExecutionPage = ({ projectId, projectName }) => {
  const { t } = useTranslation();
  const [tasks, setTasks] = useState([]);
  const [selectedTaskUuid, setSelectedTaskUuid] = useState(null);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(() => {
    const saved = localStorage.getItem(`testExecution_sidebarCollapsed_${projectId}`);
    return saved === 'true';
  });
  const [loading, setLoading] = useState(false);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchTasks();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [projectId]);

  const fetchTasks = async () => {
    setLoading(true);
    setError(null);
    try {
      console.log('🔄 [TestExecution] Fetching tasks for projectId:', projectId);
      const response = await getExecutionTasks(projectId);
      console.log('📦 [TestExecution] API response:', response);
      console.log('📦 [TestExecution] API response type:', typeof response);
      console.log('📦 [TestExecution] API response is array:', Array.isArray(response));
      
      // 响应拦截器已经返回了 data 字段，直接使用 response
      const taskList = Array.isArray(response) ? response : [];
      console.log('📋 [TestExecution] Task list length:', taskList.length);
      console.log('📋 [TestExecution] Task list:', JSON.stringify(taskList, null, 2));
      
      setTasks(taskList);
      console.log('✅ [TestExecution] Tasks state updated');

      // 自动选中第一个任务
      if (taskList.length > 0 && !selectedTaskUuid) {
        const firstTask = taskList.find(t => t.task_status === 'in_progress') ||
          taskList.find(t => t.task_status === 'pending') ||
          taskList[0];
        console.log('🎯 [TestExecution] Auto-selecting first task:', firstTask?.task_uuid);
        setSelectedTaskUuid(firstTask?.task_uuid);
      }
    } catch (err) {
      console.error('❌ [TestExecution] Fetch tasks error:', err);
      setError(t('testExecution.messages.loadFailed'));
      message.error(t('testExecution.messages.loadFailed'));
    } finally {
      setLoading(false);
    }
  };

  const handleTaskSelect = (taskUuid) => {
    setSelectedTaskUuid(taskUuid);
  };

  const handleSidebarToggle = () => {
    setSidebarCollapsed(prev => {
      const newValue = !prev;
      localStorage.setItem(`testExecution_sidebarCollapsed_${projectId}`, newValue);
      return newValue;
    });
  };

  const handleTaskCreated = () => {
    setCreateModalVisible(false);
    fetchTasks();
  };

  const handleTaskUpdated = (updatedTask) => {
    setTasks(prevTasks => {
      const newTasks = prevTasks.map(t =>
        t.task_uuid === updatedTask.task_uuid ? { ...t, ...updatedTask } : t
      );
      
      // 如果状态变化,重新排序
      return newTasks.sort((a, b) => {
        const statusOrder = { in_progress: 1, pending: 2, completed: 3 };
        const statusCompare = (statusOrder[a.task_status] || 4) - (statusOrder[b.task_status] || 4);
        if (statusCompare !== 0) return statusCompare;
        return new Date(b.created_at) - new Date(a.created_at);
      });
    });
  };

  const handleTaskDeleted = async (taskUuid) => {
    console.log('\ud83d\uddd1\ufe0f [TestExecution] handleTaskDeleted called with taskUuid:', taskUuid);
    console.log('\ud83d\uddd1\ufe0f [TestExecution] projectId:', projectId);
    
    try {
      console.log('\ud83d\udd04 [TestExecution] Calling deleteExecutionTask API...');
      const response = await deleteExecutionTask(projectId, taskUuid);
      console.log('\u2705 [TestExecution] Delete API response:', response);
      
      setTasks(prevTasks => {
        console.log('\ud83d\udcca [TestExecution] Previous tasks:', prevTasks);
        const filtered = prevTasks.filter(t => t.task_uuid !== taskUuid);
        console.log('\ud83d\udcca [TestExecution] Filtered tasks:', filtered);
        
        // 如果删除的是当前选中任务,选中下一个
        if (taskUuid === selectedTaskUuid) {
          const nextTask = filtered.find(t => t.task_status === 'in_progress') ||
            filtered.find(t => t.task_status === 'pending') ||
            filtered[0];
          console.log('\ud83c\udfaf [TestExecution] Selecting next task:', nextTask?.task_uuid);
          setSelectedTaskUuid(nextTask?.task_uuid);
        }
        
        return filtered;
      });
      
      console.log('\u2705 [TestExecution] Task deleted successfully');
      message.success(t('testExecution.messages.deleteSuccess'));
    } catch (err) {
      console.error('\u274c [TestExecution] Delete task error:', err);
      message.error(t('testExecution.messages.deleteFailed'));
    }
  };

  const selectedTask = useMemo(() => {
    return tasks.find(t => t.task_uuid === selectedTaskUuid);
  }, [tasks, selectedTaskUuid]);

  console.log('🎨 [TestExecution] Rendering with:', {
    loading,
    tasksLength: tasks.length,
    selectedTaskUuid,
    sidebarCollapsed
  });

  if (loading && tasks.length === 0) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: '24px' }}>
        <Alert
          message={t('testExecution.messages.loadFailed')}
          description={error}
          type="error"
          showIcon
          action={
            <Button size="small" onClick={fetchTasks}>
              {t('testExecution.taskCard.cancel')}
            </Button>
          }
        />
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', height: '100%', overflow: 'hidden' }}>
      {!sidebarCollapsed ? (
        <TaskSidebar
          tasks={tasks}
          selectedTaskUuid={selectedTaskUuid}
          onTaskSelect={handleTaskSelect}
          onCollapse={handleSidebarToggle}
          onTaskCreate={() => setCreateModalVisible(true)}
          onTaskDelete={handleTaskDeleted}
        />
      ) : (
        <div style={{ width: 50, padding: '16px 8px', background: '#fafafa', borderRight: '1px solid #f0f0f0', flexShrink: 0 }}>
          <Button
            icon={<MenuUnfoldOutlined />}
            onClick={handleSidebarToggle}
            block
          />
        </div>
      )}

      <div style={{ flex: 1, overflow: 'hidden', minWidth: 0 }}>
        <TaskMetadataPanel
          task={selectedTask}
          projectId={projectId}
          projectName={projectName}
          onSave={handleTaskUpdated}
        />
      </div>

      <CreateTaskModal
        visible={createModalVisible}
        projectId={projectId}
        onSuccess={handleTaskCreated}
        onCancel={() => setCreateModalVisible(false)}
      />
    </div>
  );
};

TestExecutionPage.propTypes = {
  projectId: PropTypes.number.isRequired,
  projectName: PropTypes.string,
};

export default TestExecutionPage;
