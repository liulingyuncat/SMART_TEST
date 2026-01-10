import apiClient from './client';

/**
 * 获取缺陷列表
 * @param {number|string} projectId - 项目ID
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页条数
 * @param {string} params.status - 状态筛选
 * @param {string} params.priority - 优先级筛选
 * @param {string} params.severity - 严重程度筛选
 * @param {string} params.subject - 主题筛选
 * @param {string} params.phase - 发现阶段筛选
 * @param {string} params.assignee - 负责人筛选
 * @param {string} params.keyword - 关键词搜索
 * @returns {Promise<{items: Array, total: number, page: number, page_size: number}>}
 */
export const fetchDefects = async (projectId, params = {}) => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/defects`, { params });
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 获取单个缺陷详情
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @returns {Promise<Object>}
 */
export const fetchDefect = async (projectId, defectId) => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/defects/${defectId}`);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 创建缺陷
 * @param {number|string} projectId - 项目ID
 * @param {Object} defect - 缺陷数据
 * @returns {Promise<Object>}
 */
export const createDefect = async (projectId, defect) => {
  try {
    const response = await apiClient.post(`/projects/${projectId}/defects`, defect);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 更新缺陷
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @param {Object} defect - 缺陷数据
 * @returns {Promise<Object>}
 */
export const updateDefect = async (projectId, defectId, defect) => {
  try {
    const response = await apiClient.put(`/projects/${projectId}/defects/${defectId}`, defect);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 删除缺陷
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @returns {Promise<Object>}
 */
export const deleteDefect = async (projectId, defectId) => {
  try {
    const response = await apiClient.delete(`/projects/${projectId}/defects/${defectId}`);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 导入缺陷CSV
 * @param {number|string} projectId - 项目ID
 * @param {File} file - CSV文件
 * @returns {Promise<{imported: number, skipped: number, errors: Array}>}
 */
export const importDefects = async (projectId, file) => {
  try {
    console.log('[importDefects] 接收到的file:', file);
    console.log('[importDefects] file.name:', file?.name);
    console.log('[importDefects] file.size:', file?.size);
    
    const formData = new FormData();
    formData.append('file', file);
    
    console.log('[importDefects] FormData内容:');
    for (let pair of formData.entries()) {
      console.log('  ', pair[0], pair[1]);
    }
    
    // 响应拦截器已经提取了data字段，直接返回结果
    const result = await apiClient.post(`/projects/${projectId}/defects/import`, formData);
    console.log('[importDefects] 返回结果:', result);
    return result;
  } catch (error) {
    throw error;
  }
};

/**
 * 导出缺陷（支持CSV和XLSX格式）
 * @param {number|string} projectId - 项目ID
 * @param {string} format - 文件格式（'csv' 或 'xlsx'，默认为 'csv'）
 * @param {Object} params - 筛选参数
 * @param {string} projectName - 项目名称（可选）
 * @returns {Promise<void>} 自动下载文件
 */
export const exportDefects = async (projectId, format = 'csv', params = {}, projectName = '') => {
  try {
    console.log('[exportDefects] 开始导出:', { projectId, format, params, projectName });
    const response = await apiClient.get(`/projects/${projectId}/defects/export`, {
      params: { ...params, format },
      responseType: 'blob'
    });

    console.log('[exportDefects] 响应headers:', response.headers);
    console.log('[exportDefects] Content-Type:', response.headers?.['content-type']);
    console.log('[exportDefects] Content-Disposition:', response.headers?.['content-disposition']);

    // 生成带时间戳和项目名的文件名
    const timestamp = new Date().getTime();
    const projectPrefix = projectName ? `${projectName}_` : '';
    const extension = format === 'xlsx' ? 'xlsx' : 'csv';
    const defaultFilename = `${projectPrefix}defects_export_${timestamp}.${extension}`;

    console.log('[exportDefects] 生成的默认文件名:', defaultFilename);

    // 使用后端返回的文件名或使用默认文件名
    let filename = defaultFilename;
    const contentDisposition = response.headers?.['content-disposition'];
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/);
      if (filenameMatch && filenameMatch[1]) {
        filename = filenameMatch[1].replace(/['"]/g, '');
        console.log('[exportDefects] 使用后端文件名:', filename);
      }
    }

    // 直接使用响应的blob数据，不重新创建blob
    const blob = response.data;
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
    
    console.log('[exportDefects] 下载完成:', filename);
  } catch (error) {
    console.error('[exportDefects] 导出失败:', error);
    throw error;
  }
};

/**
 * 下载模板（支持CSV和XLSX格式）
 * @param {number|string} projectId - 项目ID
 * @param {string} format - 文件格式（'csv' 或 'xlsx'，默认为 'csv'）
 * @returns {Promise<void>} 自动下载模板
 */
export const downloadDefectTemplate = async (projectId, format = 'csv') => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/defects/template`, {
      params: { format },
      responseType: 'blob'
    });

    let mimeType = 'text/csv;charset=utf-8';
    let filename = 'defect_template.csv';

    if (format === 'xlsx') {
      mimeType = 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet';
      filename = 'defect_template.xlsx';
    }

    const blob = new Blob([response.data], { type: mimeType });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    a.click();
    window.URL.revokeObjectURL(url);
  } catch (error) {
    throw error;
  }
};

/**
 * 获取缺陷统计
 * @param {number|string} projectId - 项目ID
 * @returns {Promise<Object>} 统计数据
 */
export const fetchDefectStats = async (projectId) => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/defects/stats`);
    return response;
  } catch (error) {
    throw error;
  }
};

// ==================== 附件管理 API ====================

/**
 * 上传缺陷附件
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @param {File} file - 上传的文件
 * @returns {Promise<Object>}
 */
export const uploadDefectAttachment = async (projectId, defectId, file) => {
  try {
    const formData = new FormData();
    formData.append('file', file);
    const response = await apiClient.post(`/projects/${projectId}/defects/${defectId}/attachments`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 获取缺陷附件列表
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @returns {Promise<Array>}
 */
export const fetchDefectAttachments = async (projectId, defectId) => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/defects/${defectId}/attachments`);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 下载缺陷附件
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @param {number} attachmentId - 附件ID
 * @param {string} filename - 文件名
 * @returns {Promise<void>} 自动下载文件
 */
export const downloadDefectAttachment = async (projectId, defectId, attachmentId, filename) => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/defects/${defectId}/attachments/${attachmentId}`, {
      responseType: 'blob'
    });
    const blob = new Blob([response.data]);
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename || 'attachment';
    a.click();
    window.URL.revokeObjectURL(url);
  } catch (error) {
    throw error;
  }
};

/**
 * 删除缺陷附件
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @param {number} attachmentId - 附件ID
 * @returns {Promise<Object>}
 */
export const deleteDefectAttachment = async (projectId, defectId, attachmentId) => {
  try {
    const response = await apiClient.delete(`/projects/${projectId}/defects/${defectId}/attachments/${attachmentId}`);
    return response;
  } catch (error) {
    throw error;
  }
};

// ==================== 配置管理 API ====================

/**
 * 获取缺陷主题列表
 * @param {number|string} projectId - 项目ID
 * @returns {Promise<Array>}
 */
export const fetchDefectSubjects = async (projectId) => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/defect-subjects`);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 创建缺陷主题
 * @param {number|string} projectId - 项目ID
 * @param {Object} subject - 主题数据 { name, description }
 * @returns {Promise<Object>}
 */
export const createDefectSubject = async (projectId, subject) => {
  try {
    const response = await apiClient.post(`/projects/${projectId}/defect-subjects`, subject);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 更新缺陷主题
 * @param {number|string} projectId - 项目ID
 * @param {number} subjectId - 主题ID
 * @param {Object} subject - 主题数据 { name, description }
 * @returns {Promise<Object>}
 */
export const updateDefectSubject = async (projectId, subjectId, subject) => {
  try {
    const response = await apiClient.put(`/projects/${projectId}/defect-subjects/${subjectId}`, subject);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 删除缺陷主题
 * @param {number|string} projectId - 项目ID
 * @param {number} subjectId - 主题ID
 * @returns {Promise<Object>}
 */
export const deleteDefectSubject = async (projectId, subjectId) => {
  try {
    const response = await apiClient.delete(`/projects/${projectId}/defect-subjects/${subjectId}`);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 获取发现阶段列表
 * @param {number|string} projectId - 项目ID
 * @returns {Promise<Array>}
 */
export const fetchDefectPhases = async (projectId) => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/defect-phases`);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 创建发现阶段
 * @param {number|string} projectId - 项目ID
 * @param {Object} phase - 阶段数据 { name, description }
 * @returns {Promise<Object>}
 */
export const createDefectPhase = async (projectId, phase) => {
  try {
    const response = await apiClient.post(`/projects/${projectId}/defect-phases`, phase);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 更新发现阶段
 * @param {number|string} projectId - 项目ID
 * @param {number} phaseId - 阶段ID
 * @param {Object} phase - 阶段数据 { name, description }
 * @returns {Promise<Object>}
 */
export const updateDefectPhase = async (projectId, phaseId, phase) => {
  try {
    const response = await apiClient.put(`/projects/${projectId}/defect-phases/${phaseId}`, phase);
    return response;
  } catch (error) {
    throw error;
  }
};

/**
 * 删除发现阶段
 * @param {number|string} projectId - 项目ID
 * @param {number} phaseId - 阶段ID
 * @returns {Promise<Object>}
 */
export const deleteDefectPhase = async (projectId, phaseId) => {
  try {
    const response = await apiClient.delete(`/projects/${projectId}/defect-phases/${phaseId}`);
    return response;
  } catch (error) {
    throw error;
  }
};

// ==================== 缺陷说明 API ====================

/**
 * 获取缺陷说明列表
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @returns {Promise<Array>}
 */
export const fetchDefectComments = async (projectId, defectId) => {
  try {
    const response = await apiClient.get(`/projects/${projectId}/defects/${defectId}/comments`);
    return response.comments || [];
  } catch (error) {
    throw error;
  }
};

/**
 * 创建缺陷说明
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @param {string} content - 说明内容
 * @returns {Promise<Object>}
 */
export const createDefectComment = async (projectId, defectId, content) => {
  try {
    const response = await apiClient.post(`/projects/${projectId}/defects/${defectId}/comments`, {
      content,
    });
    return response.comment;
  } catch (error) {
    throw error;
  }
};

/**
 * 更新缺陷说明
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @param {number} commentId - 说明ID
 * @param {string} content - 说明内容
 * @returns {Promise<Object>}
 */
export const updateDefectComment = async (projectId, defectId, commentId, content) => {
  try {
    const response = await apiClient.put(
      `/projects/${projectId}/defects/${defectId}/comments/${commentId}`,
      { content }
    );
    return response.comment;
  } catch (error) {
    throw error;
  }
};

/**
 * 删除缺陷说明
 * @param {number|string} projectId - 项目ID
 * @param {string} defectId - 缺陷UUID
 * @param {number} commentId - 说明ID
 * @returns {Promise<Object>}
 */
export const deleteDefectComment = async (projectId, defectId, commentId) => {
  try {
    const response = await apiClient.delete(
      `/projects/${projectId}/defects/${defectId}/comments/${commentId}`
    );
    return response;
  } catch (error) {
    throw error;
  }
};
