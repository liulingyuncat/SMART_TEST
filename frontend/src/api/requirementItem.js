import apiClient from './client';

/**
 * 需求条目API模块 (T42)
 */

// ==================== 需求条目列表操作 ====================

/**
 * 获取项目的所有需求条目
 * @param {number} projectId - 项目ID
 * @returns {Promise<Array>} 需求条目列表
 */
export const fetchRequirementItems = async (projectId) => {
  const response = await apiClient.get(`/projects/${projectId}/requirement-items`);
  return response;
};

/**
 * 创建需求条目
 * @param {number} projectId - 项目ID
 * @param {string} name - 需求名称
 * @param {string} content - 需求内容(Markdown)
 * @returns {Promise<Object>} 创建的需求条目
 */
export const createRequirementItem = async (projectId, name, content) => {
  const response = await apiClient.post(`/projects/${projectId}/requirement-items`, {
    name,
    content,
  });
  return response;
};

/**
 * 批量创建需求条目
 * @param {number} projectId - 项目ID
 * @param {Array<{name: string, content: string}>} items - 需求条目列表
 * @returns {Promise<Object>}
 */
export const bulkCreateRequirementItems = async (projectId, items) => {
  const response = await apiClient.post(`/projects/${projectId}/requirement-items/bulk`, {
    items,
  });
  return response;
};

// ==================== 单条目CRUD操作 ====================

/**
 * 获取单个需求条目
 * @param {number} id - 需求条目ID
 * @returns {Promise<Object>} 需求条目详情
 */
export const fetchRequirementItem = async (id) => {
  const response = await apiClient.get(`/requirement-items/${id}`);
  return response;
};

/**
 * 更新需求条目
 * @param {number} id - 需求条目ID
 * @param {string} name - 需求名称
 * @param {string} content - 需求内容
 * @returns {Promise<Object>}
 */
export const updateRequirementItem = async (id, name, content) => {
  const response = await apiClient.put(`/requirement-items/${id}`, {
    name,
    content,
  });
  return response;
};

/**
 * 删除需求条目
 * @param {number} id - 需求条目ID
 * @returns {Promise<Object>}
 */
export const deleteRequirementItem = async (id) => {
  const response = await apiClient.delete(`/requirement-items/${id}`);
  return response;
};

/**
 * 批量更新需求条目
 * @param {Array<{id: number, name: string, content: string}>} items - 需求条目列表
 * @returns {Promise<Object>}
 */
export const bulkUpdateRequirementItems = async (items) => {
  const response = await apiClient.put('/requirement-items/bulk', {
    items,
  });
  return response;
};

/**
 * 批量删除需求条目
 * @param {Array<number>} ids - 需求条目ID列表
 * @returns {Promise<Object>}
 */
export const bulkDeleteRequirementItems = async (ids) => {
  const response = await apiClient.delete('/requirement-items/bulk', {
    data: { ids },
  });
  return response;
};

// ==================== ZIP批量版本操作 ====================

/**
 * 导出需求条目为ZIP批量版本
 * @param {number} projectId - 项目ID
 * @param {string} remark - 备注
 * @returns {Promise<Object>} 版本记录
 */
export const exportRequirementItemsToZip = async (projectId, remark = '') => {
  console.log('=== [exportRequirementItemsToZip] API函数被调用 ===');
  console.log('[exportRequirementItemsToZip] projectId:', projectId);
  console.log('[exportRequirementItemsToZip] remark:', remark);
  
  const url = `/projects/${projectId}/requirement-items/export`;
  console.log('[exportRequirementItemsToZip] 请求URL:', url);
  console.log('[exportRequirementItemsToZip] 请求体:', { remark });
  
  try {
    console.log('[exportRequirementItemsToZip] 开始发送POST请求...');
    const response = await apiClient.post(url, { remark });
    console.log('[exportRequirementItemsToZip] API响应成功');
    console.log('[exportRequirementItemsToZip] 响应数据:', response);
    return response;
  } catch (error) {
    console.error('[exportRequirementItemsToZip] API请求失败');
    console.error('[exportRequirementItemsToZip] 错误:', error);
    throw error;
  }
};

/**
 * 从ZIP批量版本导入需求条目
 * @param {number} projectId - 项目ID
 * @param {File} file - ZIP文件
 * @returns {Promise<Object>}
 */
export const importRequirementItemsFromZip = async (projectId, file) => {
  const formData = new FormData();
  formData.append('file', file);

  const response = await apiClient.post(`/projects/${projectId}/requirement-items/import`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response;
};
