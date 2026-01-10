import apiClient from './client';

/**
 * 观点条目API模块 (T42)
 */

// ==================== 观点条目列表操作 ====================

/**
 * 获取项目的所有观点条目
 * @param {number} projectId - 项目ID
 * @returns {Promise<Array>} 观点条目列表
 */
export const fetchViewpointItems = async (projectId) => {
  const response = await apiClient.get(`/projects/${projectId}/viewpoint-items`);
  return response;
};

/**
 * 创建观点条目
 * @param {number} projectId - 项目ID
 * @param {string} name - 观点名称
 * @param {string} content - 观点内容(Markdown)
 * @returns {Promise<Object>} 创建的观点条目
 */
export const createViewpointItem = async (projectId, name, content) => {
  const response = await apiClient.post(`/projects/${projectId}/viewpoint-items`, {
    name,
    content,
  });
  return response;
};

/**
 * 批量创建观点条目
 * @param {number} projectId - 项目ID
 * @param {Array<{name: string, content: string}>} items - 观点条目列表
 * @returns {Promise<Object>}
 */
export const bulkCreateViewpointItems = async (projectId, items) => {
  const response = await apiClient.post(`/projects/${projectId}/viewpoint-items/bulk`, {
    items,
  });
  return response;
};

// ==================== 单条目CRUD操作 ====================

/**
 * 获取单个观点条目
 * @param {number} id - 观点条目ID
 * @returns {Promise<Object>} 观点条目详情
 */
export const fetchViewpointItem = async (id) => {
  const response = await apiClient.get(`/viewpoint-items/${id}`);
  return response;
};

/**
 * 更新观点条目
 * @param {number} id - 观点条目ID
 * @param {string} name - 观点名称
 * @param {string} content - 观点内容
 * @returns {Promise<Object>}
 */
export const updateViewpointItem = async (id, name, content) => {
  const response = await apiClient.put(`/viewpoint-items/${id}`, {
    name,
    content,
  });
  return response;
};

/**
 * 删除观点条目
 * @param {number} id - 观点条目ID
 * @returns {Promise<Object>}
 */
export const deleteViewpointItem = async (id) => {
  const response = await apiClient.delete(`/viewpoint-items/${id}`);
  return response;
};

/**
 * 批量更新观点条目
 * @param {Array<{id: number, name: string, content: string}>} items - 观点条目列表
 * @returns {Promise<Object>}
 */
export const bulkUpdateViewpointItems = async (items) => {
  const response = await apiClient.put('/viewpoint-items/bulk', {
    items,
  });
  return response;
};

/**
 * 批量删除观点条目
 * @param {Array<number>} ids - 观点条目ID列表
 * @returns {Promise<Object>}
 */
export const bulkDeleteViewpointItems = async (ids) => {
  const response = await apiClient.delete('/viewpoint-items/bulk', {
    data: { ids },
  });
  return response;
};

// ==================== ZIP批量版本操作 ====================

/**
 * 导出观点条目为ZIP批量版本
 * @param {number} projectId - 项目ID
 * @param {string} remark - 备注
 * @returns {Promise<Object>} 版本记录
 */
export const exportViewpointItemsToZip = async (projectId, remark = '') => {
  const response = await apiClient.post(`/projects/${projectId}/viewpoint-items/export`, {
    remark,
  });
  return response;
};

/**
 * 从ZIP批量版本导入观点条目
 * @param {number} projectId - 项目ID
 * @param {File} file - ZIP文件
 * @returns {Promise<Object>}
 */
export const importViewpointItemsFromZip = async (projectId, file) => {
  const formData = new FormData();
  formData.append('file', file);

  const response = await apiClient.post(`/projects/${projectId}/viewpoint-items/import`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response;
};
