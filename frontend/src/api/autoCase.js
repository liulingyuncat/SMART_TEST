import client from './client';

/**
 * 自动化测试用例API调用�?
 * 支持role1-4类型的CRUD操作
 */

// 获取元数�?
export const getAutoMetadata = (projectId, caseType) => {
  return client.get(`/projects/${projectId}/auto-cases/metadata`, {
    params: { type: caseType }
  });
};

// 更新元数�?
export const updateAutoMetadata = (projectId, caseType, data) => {
  return client.put(`/projects/${projectId}/auto-cases/metadata?type=${caseType}`, {
    screen_cn: data.screen_cn || '',
    screen_jp: data.screen_jp || '',
    screen_en: data.screen_en || '',
  });
};

// 在指定位置插入用例
export const insertAutoCase = (projectId, data) => {
  const requestData = {
    case_type: data.caseType,
    position: data.position,
    target_case_id: data.targetCaseId,
  };
  // 只有当 caseGroup 存在时才添加到请求中（Web用例需要此字段）
  if (data.caseGroup) {
    requestData.caseGroup = data.caseGroup;  // 后端使用驼峰命名
  }
  return client.post(`/projects/${projectId}/auto-cases/insert`, requestData);
};

// 批量删除用例
export const batchDeleteAutoCases = (projectId, data) => {
  return client.post(`/projects/${projectId}/auto-cases/batch-delete`, {
    case_type: data.caseType,
    case_ids: data.caseIds,
  });
};

// 重新分配所有用例的ID
export const reassignAutoIDs = (projectId, caseType) => {
  return client.post(`/projects/${projectId}/auto-cases/reassign-ids`, {
    caseType
  });
};

// 获取用例列表
export const getAutoCasesList = (projectId, params) => {
  const { caseType, language, caseGroup, page = 1, size = 50 } = params;
  return client.get(`/projects/${projectId}/auto-cases`, {
    params: {
      case_type: caseType,
      language,
      case_group: caseGroup,
      page,
      size
    }
  });
};

// 创建用例
export const createAutoCase = (projectId, caseData) => {
  return client.post(`/projects/${projectId}/auto-cases`, caseData);
};

// 更新用例
export const updateAutoCase = (projectId, caseId, updates) => {
  return client.patch(`/projects/${projectId}/auto-cases/${caseId}`, updates);
};

// 删除用例
export const deleteAutoCase = (projectId, caseId) => {
  return client.delete(`/projects/${projectId}/auto-cases/${caseId}`);
};

// ID重新生成（按指定顺序重排�?
export const reorderAutoCases = (projectId, caseType, caseIds) => {
  return client.post(`/projects/${projectId}/auto-cases/reorder`, {
    case_type: caseType,
    case_ids: caseIds || []
  });
};

/**
 * 导出自动化用例(Role1-4)
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('role1'|'role2'|'role3'|'role4')
 * @param {string} taskUuid - 可选：执行任务UUID，传入后导出包含执行结果
 * @returns {Promise<void>} 自动下载文件
 */
export const exportAutoCases = (projectId, caseType, taskUuid = null) => {
  const params = { case_type: caseType };
  if (taskUuid) {
    params.task_uuid = taskUuid;
  }
  return client.get(`/projects/${projectId}/auto-cases/export`, {
    params,
    responseType: 'blob'
  });
};

// 批量保存版本
export const batchSaveAutoVersion = (projectId) => {
  return client.post(`/projects/${projectId}/auto-cases/versions`);
};

// 获取版本列表
export const getAutoVersions = (projectId, page = 1, size = 10) => {
  return client.get(`/projects/${projectId}/auto-cases/versions`, {
    params: { page, size }
  });
};

// 下载版本压缩�?
export const downloadAutoVersion = async (projectId, versionId) => {
  const response = await client.get(`/projects/${projectId}/auto-cases/versions/${versionId}/export`, {
    responseType: 'blob'
  });
  
  // 处理blob响应并触发下�?
  const blob = new Blob([response.data], { type: 'application/zip' });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  
  // 从响应头获取文件�?或使用默认文件名
  const contentDisposition = response.headers['content-disposition'];
  let filename = `${versionId}.zip`;
  if (contentDisposition) {
    const filenameMatch = contentDisposition.match(/filename=(.+)/);
    if (filenameMatch && filenameMatch[1]) {
      filename = filenameMatch[1];
    }
  }
  
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  window.URL.revokeObjectURL(url);
};

// 删除版本
export const deleteAutoVersion = (projectId, versionId) => {
  return client.delete(`/projects/${projectId}/auto-cases/versions/${versionId}`);
};

// 更新版本备注
export const updateAutoVersionRemark = (projectId, versionId, remark) => {
  return client.put(`/projects/${projectId}/auto-cases/versions/${versionId}/remark`, {
    remark
  });
};

// Web用例版本管理API

// 保存Web用例版本
export const saveWebVersion = (projectId) => {
  return client.post(`/projects/${projectId}/web-cases/versions`);
};

// 获取Web用例版本列表
export const getWebVersions = (projectId, page = 1, size = 10) => {
  return client.get(`/projects/${projectId}/web-cases/versions`, {
    params: { page, size }
  });
};

// 下载Web用例版本
export const downloadWebVersion = async (projectId, versionId) => {
  const response = await client.get(`/projects/${projectId}/web-cases/versions/${versionId}/export`, {
    responseType: 'blob'
  });
  
  // 处理blob响应并触发下载
  const blob = new Blob([response.data], { type: 'application/zip' });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  
  // 从响应头获取文件名或使用默认文件名
  const contentDisposition = response.headers['content-disposition'];
  let filename = `${versionId}.zip`;
  if (contentDisposition) {
    const filenameMatch = contentDisposition.match(/filename=(.+)/);
    if (filenameMatch && filenameMatch[1]) {
      filename = filenameMatch[1];
    }
  }
  
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  window.URL.revokeObjectURL(url);
};

// 删除Web用例版本
export const deleteWebVersion = (projectId, versionId) => {
  return client.delete(`/projects/${projectId}/web-cases/versions/${versionId}`);
};

// 更新Web用例版本备注
export const updateWebVersionRemark = (projectId, versionId, remark) => {
  return client.put(`/projects/${projectId}/web-cases/versions/${versionId}/remark`, {
    remark
  });
};

// Web用例集管理API

// 获取Web用例集列表
export const getWebCaseGroups = (projectId) => {
  return client.get(`/projects/${projectId}/case-groups`, {
    params: { case_type: 'web' }
  });
};

// 创建Web用例集
export const createWebCaseGroup = (projectId, groupData) => {
  return client.post(`/projects/${projectId}/case-groups`, {
    case_type: 'web',
    group_name: groupData.groupName,
    description: groupData.description || '',
    display_order: groupData.displayOrder || 0
  });
};

// 更新Web用例集
export const updateWebCaseGroup = (groupId, updateData) => {
  return client.put(`/case-groups/${groupId}`, updateData);
};

// 删除Web用例集
export const deleteWebCaseGroup = (groupId) => {
  return client.delete(`/case-groups/${groupId}`);
};

// 导出Web用例模版（三语言ZIP包）
export const exportWebTemplate = (projectId) => {
  return client.get(`/projects/${projectId}/web-cases/template`, {
    responseType: 'blob'
  }).then(response => {
    const blob = new Blob([response.data], { type: 'application/zip' });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5);
    link.download = `Web_Case_Template_${timestamp}.zip`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  });
};

// 导入Web用例
export const importWebCases = (projectId, caseGroup, file) => {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('case_group', caseGroup);
  
  return client.post(`/projects/${projectId}/web-cases/import`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
};

// API用例集管理API (使用case_groups表)

// 获取API用例集列表（从case_groups表，case_type='api'）
export const getApiCaseGroupsFromTable = (projectId) => {
  return client.get(`/projects/${projectId}/case-groups`, {
    params: { case_type: 'api' }
  });
};

// 创建API用例集
export const createApiCaseGroupInTable = (projectId, groupData) => {
  return client.post(`/projects/${projectId}/case-groups`, {
    case_type: 'api',
    group_name: groupData.groupName,
    description: groupData.description || '',
    display_order: groupData.displayOrder || 0
  });
};

// 更新API用例集
export const updateApiCaseGroupInTable = (groupId, updateData) => {
  return client.put(`/case-groups/${groupId}`, updateData);
};

// 删除API用例集
export const deleteApiCaseGroupInTable = (groupId) => {
  return client.delete(`/case-groups/${groupId}`);
};


