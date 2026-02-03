import apiClient from './client';

/**
 * 需求Chunk API模块 (T54)
 */

/**
 * 获取需求的所有Chunk
 * @param {number} projectId - 项目ID
 * @param {number} itemId - 需求条目ID
 * @returns {Promise<Array>} Chunk列表
 */
export const fetchChunks = async (projectId, itemId) => {
  const response = await apiClient.get(`/projects/${projectId}/requirement-items/${itemId}/chunks`);
  return response;
};

/**
 * 创建新Chunk
 * @param {number} projectId - 项目ID
 * @param {number} itemId - 需求条目ID
 * @param {string} title - Chunk标题
 * @param {string} content - Chunk内容
 * @returns {Promise<Object>} 创建的Chunk
 */
export const createChunk = async (projectId, itemId, title, content) => {
  const response = await apiClient.post(`/projects/${projectId}/requirement-items/${itemId}/chunks`, {
    title,
    content,
  });
  return response;
};

/**
 * 更新Chunk
 * @param {number} chunkId - Chunk ID
 * @param {string} title - 新标题
 * @param {string} content - 新内容
 * @returns {Promise<Object>} 更新后的Chunk
 */
export const updateChunk = async (chunkId, title, content) => {
  const response = await apiClient.put(`/requirement-chunks/${chunkId}`, {
    title,
    content,
  });
  return response;
};

/**
 * 删除Chunk
 * @param {number} chunkId - Chunk ID
 * @returns {Promise<Object>}
 */
export const deleteChunk = async (chunkId) => {
  const response = await apiClient.delete(`/requirement-chunks/${chunkId}`);
  return response;
};

/**
 * 重排序Chunk
 * @param {number} projectId - 项目ID
 * @param {number} itemId - 需求条目ID
 * @param {Array<{id: number, sort_order: number}>} chunkOrders - 排序数组
 * @returns {Promise<Object>}
 */
export const reorderChunks = async (projectId, itemId, chunkOrders) => {
  const response = await apiClient.put(`/projects/${projectId}/requirement-items/${itemId}/chunks/reorder`, {
    chunk_orders: chunkOrders,
  });
  return response;
};
