import apiClient from './client';

/**
 * 用户自定义变量 API
 */

/**
 * 获取用例集的变量列表
 * @param {number} projectId - 项目ID
 * @param {number} groupId - 用例集ID
 * @param {string} groupType - 用例集类型 (web/api)
 * @returns {Promise} 变量列表
 */
export const getVariables = async (projectId, groupId, groupType = 'web') => {
  const response = await apiClient.get(
    `/projects/${projectId}/case-groups/${groupId}/variables`,
    { params: { group_type: groupType } }
  );
  return response;
};

/**
 * 批量保存变量（替换模式）
 * @param {number} projectId - 项目ID
 * @param {number} groupId - 用例集ID
 * @param {string} groupType - 用例集类型 (web/api)
 * @param {Array} variables - 变量数组
 * @returns {Promise} 保存结果
 */
export const saveVariables = async (projectId, groupId, groupType, variables) => {
  const response = await apiClient.put(
    `/projects/${projectId}/case-groups/${groupId}/variables`,
    {
      project_id: projectId,
      group_type: groupType,
      variables: variables,
    }
  );
  return response;
};

/**
 * 添加单个变量
 * @param {number} projectId - 项目ID
 * @param {number} groupId - 用例集ID
 * @param {string} groupType - 用例集类型 (web/api)
 * @param {Object} variable - 变量数据
 * @returns {Promise} 新增的变量
 */
export const addVariable = async (projectId, groupId, groupType, variable) => {
  const response = await apiClient.post(
    `/projects/${projectId}/case-groups/${groupId}/variables`,
    {
      project_id: projectId,
      group_type: groupType,
      var_key: variable.var_key,
      var_desc: variable.var_desc || '',
      var_value: variable.var_value || '',
      var_type: variable.var_type || 'custom',
    }
  );
  return response;
};

/**
 * 更新单个变量
 * @param {number} projectId - 项目ID
 * @param {number} groupId - 用例集ID
 * @param {number} varId - 变量ID
 * @param {Object} variable - 变量数据
 * @returns {Promise} 更新结果
 */
export const updateVariable = async (projectId, groupId, varId, variable) => {
  const response = await apiClient.put(
    `/projects/${projectId}/case-groups/${groupId}/variables/${varId}`,
    {
      var_key: variable.var_key,
      var_desc: variable.var_desc || '',
      var_value: variable.var_value || '',
      var_type: variable.var_type || 'custom',
    }
  );
  return response;
};

/**
 * 删除单个变量
 * @param {number} projectId - 项目ID
 * @param {number} groupId - 用例集ID
 * @param {number} varId - 变量ID
 * @returns {Promise} 删除结果
 */
export const deleteVariable = async (projectId, groupId, varId) => {
  const response = await apiClient.delete(
    `/projects/${projectId}/case-groups/${groupId}/variables/${varId}`
  );
  return response;
};

/**
 * 获取执行任务的变量列表
 * 优先返回任务独立的变量，如果没有则返回用例集的变量
 * @param {number} projectId - 项目ID
 * @param {string} taskUuid - 任务UUID
 * @param {number} groupId - 用例集ID（可选，用于回退到用例集变量）
 * @param {string} groupType - 用例集类型 (web/api)
 * @returns {Promise} 变量列表
 */
export const getTaskVariables = async (projectId, taskUuid, groupId = 0, groupType = 'web') => {
  const response = await apiClient.get(
    `/projects/${projectId}/execution-tasks/${taskUuid}/variables`,
    { params: { group_id: groupId, group_type: groupType } }
  );
  return response;
};

/**
 * 批量保存任务变量（替换模式）
 * @param {number} projectId - 项目ID
 * @param {string} taskUuid - 任务UUID
 * @param {number} groupId - 用例集ID
 * @param {string} groupType - 用例集类型 (web/api)
 * @param {Array} variables - 变量数组
 * @returns {Promise} 保存结果
 */
export const saveTaskVariables = async (projectId, taskUuid, groupId, groupType, variables) => {
  const response = await apiClient.put(
    `/projects/${projectId}/execution-tasks/${taskUuid}/variables`,
    {
      project_id: projectId,
      group_id: groupId,
      group_type: groupType,
      variables: variables,
    }
  );
  return response;
};

export default {
  getVariables,
  saveVariables,
  addVariable,
  updateVariable,
  deleteVariable,
  getTaskVariables,
  saveTaskVariables,
};
