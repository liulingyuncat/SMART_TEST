import client from './client';

/**
 * AI报告API模块 (T47)
 * 支持项目级报告管理,包含4种类型:
 * R=用例审阅, A=品质分析, T=测试结果, O=其他
 */

/**
 * 获取项目AI报告列表
 * @param {number} projectId - 项目ID
 * @param {string} reportType - 可选,报告类型(R/A/T/O)
 * @returns {Promise<Array>} AI报告数组
 */
export const fetchAIReports = async (projectId, reportType = '') => {
  const params = reportType ? `?type=${reportType}` : '';
  const response = await client.get(`/projects/${projectId}/ai-reports${params}`);
  return response.items || response || [];
};

/**
 * 创建新AI报告
 * @param {number} projectId - 项目ID
 * @param {string} name - 报告名称
 * @param {string} reportType - 报告类型(R/A/T/O),默认O
 * @returns {Promise<Object>} 创建的AI报告
 */
export const createAIReport = async (projectId, name, reportType = 'O') => {
  const response = await client.post(`/projects/${projectId}/ai-reports`, { name, type: reportType });
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

