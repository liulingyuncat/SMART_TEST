import apiClient from './client';

/**
 * 获取提示词列表
 * @param {Object} params - 查询参数 {project_id, scope, page, page_size}
 * @returns {Promise} 返回提示词列表和总数 {items: [], total: number}
 */
export const fetchPrompts = (params) => {
  return apiClient.get('/prompts', { params });
};

/**
 * 根据ID获取提示词详情
 * @param {number} id - 提示词ID
 * @returns {Promise} 返回提示词对象
 */
export const fetchPromptById = (id) => {
  return apiClient.get(`/prompts/${id}`);
};

/**
 * 创建新提示词
 * @param {Object} promptData - 提示词数据
 * @param {number} promptData.project_id - 项目ID
 * @param {string} promptData.name - 提示词名称
 * @param {string} promptData.description - 描述
 * @param {string} promptData.version - 版本号
 * @param {string} promptData.content - Markdown内容
 * @param {Array} promptData.arguments - 参数列表 [{name, description, required}]
 * @param {string} promptData.scope - 作用域 'project' | 'user'
 * @returns {Promise} 返回新创建的提示词对象
 */
export const createPrompt = (promptData) => {
  return apiClient.post('/prompts', promptData);
};

/**
 * 更新提示词
 * @param {number} id - 提示词ID
 * @param {Object} promptData - 更新数据
 * @param {string} [promptData.description] - 描述
 * @param {string} [promptData.version] - 版本号
 * @param {string} [promptData.content] - Markdown内容
 * @param {Array} [promptData.arguments] - 参数列表
 * @returns {Promise} 返回更新后的提示词对象
 */
export const updatePrompt = (id, promptData) => {
  return apiClient.put(`/prompts/${id}`, promptData);
};

/**
 * 删除提示词
 * @param {number} id - 提示词ID
 * @returns {Promise} 返回删除结果
 */
export const deletePrompt = (id) => {
  return apiClient.delete(`/prompts/${id}`);
};

/**
 * 刷新MCP提示词缓存（热更新）
 * 当提示词发生增删改时调用此接口，通知后端更新MCP层的提示词缓存
 * @param {string} scope - 可选，指定要刷新的作用域: 'system' | 'project' | 'user' | 'all'
 * @returns {Promise} 返回刷新结果
 */
export const refreshPrompts = (scope = 'all') => {
  return apiClient.post('/prompts/refresh', {
    scope
  });
};
