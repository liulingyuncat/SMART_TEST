import apiClient from './client';

/**
 * 获取项目的测试执行任务列表
 * @param {number} projectId - 项目ID
 * @returns {Promise<Array>} 任务数组
 */
export const getExecutionTasks = (projectId) => {
  return apiClient.get(`/projects/${projectId}/execution-tasks`);
};

/**
 * 创建测试执行任务
 * @param {number} projectId - 项目ID
 * @param {Object} data - 任务数据
 * @param {string} data.task_name - 任务名称
 * @param {string} data.execution_type - 执行内容类型 (manual/automation/api)
 * @param {string} [data.task_status] - 任务状态 (pending/in_progress/completed)
 * @returns {Promise<Object>} 创建的任务对象
 */
export const createExecutionTask = (projectId, data) => {
  return apiClient.post(`/projects/${projectId}/execution-tasks`, data);
};

/**
 * 更新测试执行任务
 * @param {number} projectId - 项目ID
 * @param {string} taskUuid - 任务UUID
 * @param {Object} data - 更新数据
 * @returns {Promise<Object>} 更新后的任务对象
 */
export const updateExecutionTask = (projectId, taskUuid, data) => {
  return apiClient.put(`/projects/${projectId}/execution-tasks/${taskUuid}`, data);
};

/**
 * 删除测试执行任务
 * @param {number} projectId - 项目ID
 * @param {string} taskUuid - 任务UUID
 * @returns {Promise<Object>} 删除结果
 */
export const deleteExecutionTask = (projectId, taskUuid) => {
  return apiClient.delete(`/projects/${projectId}/execution-tasks/${taskUuid}`);
};

/**
 * 执行测试任务
 * @param {number} projectId - 项目ID
 * @param {string} taskUuid - 任务UUID
 * @returns {Promise<Object>} 执行结果统计
 */
export const executeExecutionTask = (projectId, taskUuid) => {
  return apiClient.post(`/projects/${projectId}/execution-tasks/${taskUuid}/execute`, null, {
    timeout: 60000, // 60秒超时，因为执行测试用例可能需要较长时间
  });
};
/**
 * 执行单条测试用例
 * @param {number} projectId - 项目ID
 * @param {string} taskUuid - 任务UUID
 * @param {number} caseResultId - 用例结果ID
 * @returns {Promise<Object>} 执行结果统计
 */
export const executeSingleCase = (projectId, taskUuid, caseResultId) => {
  return apiClient.post(`/projects/${projectId}/execution-tasks/${taskUuid}/cases/${caseResultId}/execute`, null, {
    timeout: 60000, // 60秒超时
  });
};