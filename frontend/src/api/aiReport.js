import client from './client';

/**
 * AI质量报告API模块 (T47)
 * 支持项目级报告管理
 */

/**
 * 获取项目所有AI报告列表
 * @param {number} projectId - 项目ID
 * @returns {Promise<Array>} AI报告数组
 */
export const fetchAIReports = async (projectId) => {
  const response = await client.get(`/projects/${projectId}/ai-reports`);
  return response.items || response || [];
};

/**
 * 创建新AI报告
 * @param {number} projectId - 项目ID
 * @param {string} name - 报告名称
 * @returns {Promise<Object>} 创建的AI报告
 */
export const createAIReport = async (projectId, name) => {
  const response = await client.post(`/projects/${projectId}/ai-reports`, { name });
  return response;
};

/**
 * 获取AI报告详情
 * @param {number} projectId - 项目ID
 * @param {string} reportId - 报告ID
 * @returns {Promise<Object>} AI报告详情
 */
export const fetchAIReportDetail = async (projectId, reportId) => {
  const response = await client.get(`/projects/${projectId}/ai-reports/${reportId}`);
  return response;
};

/**
 * 更新AI报告
 * @param {number} projectId - 项目ID
 * @param {string} reportId - 报告ID
 * @param {Object} data - 更新数据 {name?, content?}
 * @returns {Promise<Object>} 更新后的AI报告
 */
export const updateAIReport = async (projectId, reportId, data) => {
  const response = await client.put(`/projects/${projectId}/ai-reports/${reportId}`, data);
  return response;
};

/**
 * 删除AI报告
 * @param {number} projectId - 项目ID
 * @param {string} reportId - 报告ID
 * @returns {Promise<Object>} 删除结果
 */
export const deleteAIReport = async (projectId, reportId) => {
  const response = await client.delete(`/projects/${projectId}/ai-reports/${reportId}`);
  return response;
};
