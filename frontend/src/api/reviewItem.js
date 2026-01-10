import client from './client';

/**
 * 审阅条目API模块 (T44)
 * 支持多文档审阅管理
 */

/**
 * 获取项目所有审阅条目列表
 * @param {number} projectId - 项目ID
 * @returns {Promise<Array>} 审阅条目数组
 */
export const getReviewItems = async (projectId) => {
  const response = await client.get(`/projects/${projectId}/review-items`);
  return response.items || [];
};

/**
 * 创建新审阅条目
 * @param {number} projectId - 项目ID
 * @param {string} name - 审阅名称
 * @returns {Promise<Object>} 创建的审阅条目
 */
export const createReviewItem = async (projectId, name) => {
  const response = await client.post(`/projects/${projectId}/review-items`, { name });
  return response;
};

/**
 * 获取审阅条目详情
 * @param {number} projectId - 项目ID
 * @param {number} itemId - 审阅条目ID
 * @returns {Promise<Object>} 审阅条目详情
 */
export const getReviewItem = async (projectId, itemId) => {
  const response = await client.get(`/projects/${projectId}/review-items/${itemId}`);
  return response;
};

/**
 * 更新审阅条目
 * @param {number} projectId - 项目ID
 * @param {number} itemId - 审阅条目ID
 * @param {Object} data - 更新数据 {name?, content?}
 * @returns {Promise<Object>} 更新后的审阅条目
 */
export const updateReviewItem = async (projectId, itemId, data) => {
  const response = await client.put(`/projects/${projectId}/review-items/${itemId}`, data);
  return response;
};

/**
 * 删除审阅条目
 * @param {number} projectId - 项目ID
 * @param {number} itemId - 审阅条目ID
 * @returns {Promise<Object>} 删除结果
 */
export const deleteReviewItem = async (projectId, itemId) => {
  const response = await client.delete(`/projects/${projectId}/review-items/${itemId}`);
  return response;
};

/**
 * 下载审阅文档为Markdown文件
 * @param {number} projectId - 项目ID
 * @param {number} itemId - 审阅条目ID
 * @returns {Promise<void>}
 */
export const downloadReviewItem = async (projectId, itemId) => {
  try {
    const response = await client.get(`/projects/${projectId}/review-items/${itemId}/download`, {
      responseType: 'blob',
    });

    // 从响应头获取文件名
    const contentDisposition = response.headers['content-disposition'];
    let filename = `review_${itemId}.md`;
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename=(.+)/);
      if (filenameMatch) {
        filename = filenameMatch[1];
      }
    }

    // 创建下载链接
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', filename);
    document.body.appendChild(link);
    link.click();
    link.parentNode.removeChild(link);
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error('Download review item failed:', error);
    throw error;
  }
};
