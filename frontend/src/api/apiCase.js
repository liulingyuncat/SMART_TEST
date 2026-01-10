import client from './client';

/**
 * 接口测试用例API调用�?
 * 支持role1-4类型的CRUD操作
 */

// 获取用例列表
export const getApiCasesList = (projectId, params) => {
  const { caseType, caseGroup, page = 1, size = 50 } = params;
  return client.get(`/projects/${projectId}/api-cases`, {
    params: {
      case_type: caseType,
      case_group: caseGroup, // 添加用例集筛选参数
      page,
      size
    }
  });
};

// 在指定位置插入用例
export const insertApiCase = (projectId, data) => {
  return client.post(`/projects/${projectId}/api-cases/insert`, {
    case_type: data.caseType,
    position: data.position,
    target_case_id: data.targetCaseId, // UUID
    case_group: data.caseGroup, // 用例集名称
    case_data: data.caseData
  });
};

// 批量删除用例
export const batchDeleteApiCases = (projectId, data) => {
  return client.post(`/projects/${projectId}/api-cases/batch-delete`, {
    case_type: data.caseType,
    case_ids: data.caseIds // UUID数组
  });
};

// 创建用例
export const createApiCase = (projectId, caseData) => {
  return client.post(`/projects/${projectId}/api-cases`, caseData);
};

// 更新用例
export const updateApiCase = (projectId, caseId, updates) => {
  return client.patch(`/projects/${projectId}/api-cases/${caseId}`, updates);
};

// 删除用例
export const deleteApiCase = (projectId, caseId) => {
  return client.delete(`/projects/${projectId}/api-cases/${caseId}`);
};

// ========== 版本管理相关API ==========

// 保存版本
export const saveApiVersion = (projectId) => {
  return client.post(`/projects/${projectId}/api-cases/versions`);
};

// 获取版本列表(兼容VersionManagementTab的调用方式)
export const getApiVersionList = (projectId, docTypeOrPage = 1, size = 50) => {
  // 兼容两种调用方式:
  // 1. getApiVersionList(projectId, 'api') - docType传入,用于VersionManagementTab
  // 2. getApiVersionList(projectId, 1, 10) - page/size传入
  const page = typeof docTypeOrPage === 'number' ? docTypeOrPage : 1;
  
  return client.get(`/projects/${projectId}/api-cases/versions`, {
    params: { page, size }
  });
};

// 下载版本压缩�?
export const downloadApiVersion = async (projectId, versionId) => {
  const response = await client.get(`/projects/${projectId}/api-cases/versions/${versionId}/export`, {
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
export const deleteApiVersion = (projectId, versionId) => {
  return client.delete(`/projects/${projectId}/api-cases/versions/${versionId}`);
};

// 更新版本备注
export const updateApiVersionRemark = (projectId, versionId, remark) => {
  return client.put(`/projects/${projectId}/api-cases/versions/${versionId}/remark`, {
    remark
  });
};

// ========== 用例集管理相关API ==========

// 获取用例集列表
export const getApiCaseGroups = (projectId) => {
  return client.get(`/projects/${projectId}/api-case-groups`);
};

// 创建用例集
export const createApiCaseGroup = (projectId, groupName) => {
  return client.post(`/projects/${projectId}/api-case-groups`, {
    group_name: groupName
  });
};

// 更新用例集名称
export const updateApiCaseGroup = (groupId, newGroupName) => {
  return client.put(`/api-case-groups/${groupId}`, {
    group_name: newGroupName
  });
};

// 删除用例集
export const deleteApiCaseGroup = (groupId) => {
  return client.delete(`/api-case-groups/${groupId}`);
};

// 导出API用例模版
export const exportApiTemplate = async (projectId) => {
  const response = await client.get(`/projects/${projectId}/api-cases/template`, {
    responseType: 'blob'
  });
  
  // 处理blob响应并触发下载
  const blob = new Blob([response.data], { 
    type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' 
  });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  
  // 从响应头获取文件名或使用默认文件名
  const contentDisposition = response.headers['content-disposition'];
  let filename = `API_Case_Template_${new Date().toISOString().replace(/[:.]/g, '-')}.xlsx`;
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

// 导入API用例
export const importApiCases = (projectId, file, caseGroup) => {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('case_group', caseGroup);
  
  return client.post(`/projects/${projectId}/api-cases/import`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
};

