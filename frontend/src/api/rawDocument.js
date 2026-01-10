import apiClient from './client';

/**
 * 上传原始需求文档
 * @param {number|string} projectId - 项目ID
 * @param {FormData} formData - 包含文件的表单数据
 * @returns {Promise<Object>} 上传结果 {id, filename, original_filename, file_size, mime_type, upload_time}
 */
export const uploadRawDocument = async (projectId, formData) => {
  try {
    const response = await apiClient.post(`/projects/${projectId}/raw-documents`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 获取项目的原始需求文档列表
 * @param {number|string} projectId - 项目ID
 * @returns {Promise<Array>} 文档列表 [{id, filename, original_filename, file_size, mime_type, upload_time, convert_status, converted_filename}]
 */
export const fetchRawDocuments = async (projectId) => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/raw-documents`);
    // 后端返回 {documents: [...], total: ...}，提取documents数组
    return response.documents || [];
  } catch (error) {
    throw error;
  }
};

/**
 * 转换原始文档为Markdown
 * @param {number|string} documentId - 文档ID
 * @returns {Promise<Object>} 转换结果
 */
export const convertRawDocument = async (documentId) => {
  try {
    const response = await apiClient.post(`/raw-documents/${documentId}/convert`);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 获取文档转换状态
 * @param {number|string} documentId - 文档ID
 * @returns {Promise<Object>} 状态信息 {status, message, converted_filename}
 */
export const getConvertStatus = async (documentId) => {
  try {
    const response = await apiClient.get(`/raw-documents/${documentId}/convert-status`);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 下载原始文档
 * @param {number|string} documentId - 文档ID
 * @param {string} filename - 文件名
 * @returns {Promise<Blob>}
 */
export const downloadOriginalDocument = async (documentId, filename) => {
  try {
    const response = await apiClient.get(`/raw-documents/${documentId}/download`, {
      responseType: 'blob',
    });
    
    // 创建下载链接 - response是完整对象，需要使用response.data
    const blob = response.data || response;
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', filename);
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
    
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 下载转换后的Markdown文档
 * @param {number|string} documentId - 文档ID
 * @param {string} filename - 文件名
 * @returns {Promise<Blob>}
 */
export const downloadConvertedDocument = async (documentId, filename) => {
  try {
    const response = await apiClient.get(`/raw-documents/${documentId}/converted/download`, {
      responseType: 'blob',
    });
    
    // 创建下载链接 - response是完整对象，需要使用response.data
    const blob = response.data || response;
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', filename || 'converted.md');
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
    
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 删除原始文档（同时删除转换后的文档）
 * @param {number|string} documentId - 文档ID
 * @returns {Promise<void>}
 */
export const deleteOriginalDocument = async (documentId) => {
  try {
    await apiClient.delete(`/raw-documents/${documentId}`);
  } catch (error) {
    throw error;
  }
};

/**
 * 仅删除转换后的Markdown文档
 * @param {number|string} documentId - 文档ID
 * @returns {Promise<void>}
 */
export const deleteConvertedDocument = async (documentId) => {
  try {
    await apiClient.delete(`/raw-documents/${documentId}/converted`);
  } catch (error) {
    throw error;
  }
};

/**
 * 预览转换后的Markdown文档内容
 * @param {number|string} documentId - 文档ID
 * @returns {Promise<Object>} {filename, content}
 */
export const previewConvertedDocument = async (documentId) => {
  try {
    const response = await apiClient.get(`/raw-documents/${documentId}/converted/preview`);
    return response;
  } catch (error) {
    throw error;
  }
};
