import apiClient from './client';

/**
 * 获取需求文档
 * @param {number} projectId - 项目ID
 * @param {string} docType - 文档类型 (overall-requirements, overall-test-viewpoint, change-requirements, change-test-viewpoint)
 * @returns {Promise} 返回文档内容和更新时间
 */
export const fetchRequirement = async (projectId, docType) => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/requirements/${docType}`);
    return {
      content: response.content || '',
      updatedAt: response.updated_at,
    };
  } catch (error) {
    // 错误已由 axios interceptor 处理
    throw error;
  }
};

/**
 * 更新需求文档
 * @param {number} projectId - 项目ID
 * @param {string} docType - 文档类型
 * @param {string} content - 文档内容(Markdown格式)
 * @returns {Promise} 返回更新结果
 */
export const updateRequirement = async (projectId, docType, content) => {
  try {
    const response = await apiClient.put(`/projects/${projectId}/requirements/${docType}`, {
      content,
    });
    return response; // 拦截器已经返回了 data
  } catch (error) {
    // 错误已由 axios interceptor 处理
    throw error;
  }
};

/**
 * 保存版本(自动生成文件名)
 * @param {string} projectId - 项目ID
 * @param {string} docType - 文档类型 ('overall-requirements'|'overall-test-viewpoint'|'change-requirements'|'change-test-viewpoint')
 * @param {string} content - 文档内容(Markdown格式)
 * @returns {Promise<{filename: string}>}
 */
export const saveVersion = async (projectId, docType, content) => {
  try {
    const response = await apiClient.post('/versions', {
      project_id: String(projectId),
      doc_type: docType,
      content
    });
    return response; // apiClient拦截器已经返回了response.data
  } catch (error) {
    throw error;
  }
};

/**
 * 获取版本列表
 * @param {string} projectId - 项目ID
 * @param {string} docType - 文档类型
 * @returns {Promise<Array>}
 */
export const getVersionList = async (projectId, docType = '') => {
  try {
    const url = docType 
      ? `/versions?project_id=${projectId}&doc_type=${docType}`
      : `/versions?project_id=${projectId}`;
    console.log('[getVersionList] 请求URL:', url);
    const response = await apiClient.get(url);
    console.log('[getVersionList] 响应:', response);
    return response; // apiClient拦截器已经返回了response.data
  } catch (error) {
    console.error('[getVersionList] 错误:', error);
    throw error;
  }
};

/**
 * 下载指定版本文件
 * @param {string} projectId - 项目ID (保留参数以保持接口一致性)
 * @param {number} versionId - 版本ID
 * @returns {Promise<void>} 自动下载文件
 */
export const downloadVersion = async (projectId, versionId) => {
  try {
    const response = await apiClient.get(`/versions/${versionId}/download`, {
      responseType: 'blob'
    });
    const blob = new Blob([response.data], { type: 'text/markdown;charset=utf-8' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = response.headers['content-disposition']?.split('filename=')[1] || '版本文件.md';
    a.click();
    window.URL.revokeObjectURL(url);
  } catch (error) {
    throw error;
  }
};

/**
 * 删除指定版本(软删除)
 * @param {string} projectId - 项目ID (保留参数以保持接口一致性)
 * @param {number} versionId - 版本ID
 * @returns {Promise<{message: string}>}
 */
export const deleteVersion = async (projectId, versionId) => {
  try {
    const response = await apiClient.delete(`/versions/${versionId}`);
    return response; // apiClient拦截器已经返回了response.data
  } catch (error) {
    throw error;
  }
};

/**
 * 更新版本备注
 * @param {string} projectId - 项目ID (保留参数以保持接口一致性)
 * @param {number} versionId - 版本ID
 * @param {string} remark - 备注内容
 * @returns {Promise<{message: string}>}
 */
export const updateVersionRemark = async (projectId, versionId, remark) => {
  try {
    const response = await apiClient.put(`/versions/${versionId}/remark`, { remark });
    return response; // apiClient拦截器已经返回了response.data
  } catch (error) {
    throw error;
  }
};
