import client from './client';

/**
 * 获取任务的执行结果列表
 * @param {string} taskUuid - 任务UUID
 * @returns {Promise<Array>} 执行结果数组
 */
export const getExecutionCaseResults = async (taskUuid) => {
  console.log('📡 [API] getExecutionCaseResults called, taskUuid:', taskUuid);
  try {
    const results = await client.get(`/execution-tasks/${taskUuid}/case-results`);
    console.log('📡 [API] getExecutionCaseResults results:', results);
    // client已经处理了response.data.data的解析
    return Array.isArray(results) ? results : [];
  } catch (error) {
    console.error('📡 [API] getExecutionCaseResults failed:', error);
    throw error;
  }
};

/**
 * 保存或更新执行结果
 * @param {string} taskUuid - 任务UUID
 * @param {Array} resultsArray - 执行结果数组 [{case_id, test_result, bug_id, remark}]
 * @returns {Promise<Object>} 响应数据
 */
export const saveExecutionCaseResults = async (taskUuid, resultsArray) => {
  console.log('📡 [API] saveExecutionCaseResults called');
  console.log('📡 [API] taskUuid:', taskUuid);
  console.log('📡 [API] resultsArray length:', resultsArray?.length);
  
  if (!taskUuid) {
    console.error('📡 [API] ERROR: taskUuid is empty!');
    throw new Error('taskUuid is required');
  }
  if (!resultsArray || resultsArray.length === 0) {
    console.error('📡 [API] ERROR: resultsArray is empty!');
    throw new Error('resultsArray is required');
  }
  
  try {
    const response = await client.patch(`/execution-tasks/${taskUuid}/case-results`, resultsArray);
    console.log('📡 [API] saveExecutionCaseResults success');
    return response;
  } catch (error) {
    console.error('📡 [API] saveExecutionCaseResults failed:', error);
    throw error;
  }
};

/**
 * 获取任务的统计信息
 * @param {string} taskUuid - 任务UUID
 * @returns {Promise<Object>} 统计对象 {total, nr_count, ok_count, ng_count, block_count}
 */
export const getExecutionStatistics = async (taskUuid) => {
  try {
    return await client.get(`/execution-tasks/${taskUuid}/statistics`);
  } catch (error) {
    console.error('[API] getExecutionStatistics failed:', error);
    throw error;
  }
};

/**
 * 初始化任务的执行结果(创建默认NR记录)
 * @param {string} taskUuid - 任务UUID
 * @param {number} projectId - 项目ID
 * @param {string} executionType - 执行类型 (manual/automation/api)
 * @returns {Promise<Object>} 响应数据
 */
export const initExecutionResults = async (taskUuid, projectId, executionType) => {
  try {
    return await client.post(`/execution-tasks/${taskUuid}/case-results/init`, {
      project_id: projectId,
      execution_type: executionType
    });
  } catch (error) {
    console.error('[API] initExecutionResults failed:', error);
    throw error;
  }
};

/**
 * 清空任务的执行结果
 * @param {string} taskUuid - 任务UUID
 * @returns {Promise<Object>} 响应数据
 */
export const clearExecutionResults = async (taskUuid) => {
  try {
    return await client.delete(`/execution-tasks/${taskUuid}/case-results`);
  } catch (error) {
    console.error('[API] clearExecutionResults failed:', error);
    throw error;
  }
};
