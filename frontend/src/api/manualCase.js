import client from './client';

/**
 * 获取项目手工测试用例元数�?
 * @param {number} projectId - 项目ID
 * @param {string} type - 用例类型 ('overall'|'change')
 * @returns {Promise<{test_version: string, test_env: string, test_date: string, executor: string}>}
 */
export const getMetadata = async (projectId, type = 'overall') => {
  const response = await client.get(`/projects/${projectId}/manual-cases/metadata`, {
    params: { type }
  });
  return response;
};

/**
 * 更新项目手工测试用例元数�?
 * @param {number} projectId - 项目ID
 * @param {string} type - 用例类型 ('overall'|'change')
 * @param {Object} metadata - 元数据对�?
 * @param {string} metadata.test_version - 测试版本
 * @param {string} metadata.test_env - 测试环境
 * @param {string} metadata.test_date - 测试日期 (YYYY-MM-DD)
 * @param {string} metadata.executor - 执行�?
 * @returns {Promise}
 */
export const updateMetadata = async (projectId, type = 'overall', metadata) => {
  const response = await client.put(`/projects/${projectId}/manual-cases/metadata`, metadata, {
    params: { type }
  });
  return response;
};

/**
 * 在指定位置插入用例
 * @param {number} projectId - 项目ID
 * @param {Object} data - 插入参数
 * @param {string} data.caseType - 用例类型
 * @param {string} data.position - 位置 ('before'|'after')
 * @param {string} data.targetCaseId - 目标用例ID
 * @param {string} data.language - 语言
 * @param {string} data.caseGroup - 用例集名称（可选）
 * @returns {Promise}
 */
export const insertCase = async (projectId, data) => {
  const requestData = {
    case_type: data.caseType,
    position: data.position,
    target_case_id: data.targetCaseId,
    language: data.language,
  };
  // 只有当 caseGroup 存在时才添加到请求中
  if (data.caseGroup) {
    requestData.case_group = data.caseGroup;
  }
  const response = await client.post(`/projects/${projectId}/manual-cases/insert`, requestData);
  return response;
};

/**
 * 批量删除用例
 * @param {number} projectId - 项目ID
 * @param {Object} data - 删除参数
 * @param {string} data.caseType - 用例类型
 * @param {Array<string>} data.caseIds - 用例ID数组
 * @returns {Promise}
 */
export const batchDeleteCases = async (projectId, data) => {
  const response = await client.post(`/projects/${projectId}/manual-cases/batch-delete`, {
    case_type: data.caseType,
    case_ids: data.caseIds,
  });
  return response;
};

/**
 * 查询测试用例列表
 * @param {number} projectId - 项目ID
 * @param {Object} params - 查询参数
 * @param {string} params.caseType - 用例类型 ('overall'|'change'|'ai')
 * @param {string} params.language - 语言筛�?('中文'/'English'/'日本�?)
 * @param {string} params.caseGroup - 用例集名称（可选）
 * @param {number} params.page - 页码 (默认1)
 * @param {number} params.size - 每页条数 (默认50)
 * @returns {Promise<{cases: Array, total: number, page: number, size: number, language: string}>}
 */
export const getCasesList = async (projectId, params = {}) => {
  const { caseType = 'overall', language = '中文', caseGroup, page = 1, size = 50 } = params;
  const requestParams = { case_type: caseType, language, page, size };
  // 只有当 caseGroup 不为空时才添加到请求参数中
  if (caseGroup) {
    requestParams.case_group = caseGroup;
  }
  const response = await client.get(`/projects/${projectId}/manual-cases`, {
    params: requestParams
  });
  return response;
};

/**
 * 创建测试用例
 * @param {number} projectId - 项目ID
 * @param {Object} caseData - 用例数据
 * @param {string} caseData.case_type - 用例类型 ('overall'|'change'|'ai')
 * @param {string} caseData.language - 语言 ('中文'|'English'|'日本�?)
 * @param {string} caseData.case_number - 用例编号
 * @param {string} caseData.major_function - 一级功�?
 * @param {string} caseData.middle_function - 二级功能
 * @param {string} caseData.minor_function - 三级功能
 * @param {string} caseData.precondition - 前置条件
 * @param {string} caseData.test_steps - 测试步骤
 * @param {string} caseData.expected_result - 期望结果
 * @param {string} caseData.test_result - 测试结果 ('NR'|'Pass'|'Fail')
 * @param {string} caseData.remark - 备注
 * @returns {Promise<Object>}
 */
export const createCase = async (projectId, caseData) => {
  const response = await client.post(`/projects/${projectId}/manual-cases`, caseData);
  return response;
};

/**
 * 更新测试用例(部分更新)
 * @param {number} projectId - 项目ID
 * @param {number} caseId - 用例ID
 * @param {Object} updates - 需要更新的字段
 * @returns {Promise}
 */
export const updateCase = async (projectId, caseId, updates) => {
  const response = await client.patch(`/projects/${projectId}/manual-cases/${caseId}`, updates);
  return response;
};

/**
 * 删除测试用例(支持多语言联动删除)
 * @param {number} projectId - 项目ID
 * @param {number} caseId - 用例ID
 * @returns {Promise}
 */
export const deleteCase = async (projectId, caseId) => {
  const response = await client.delete(`/projects/${projectId}/manual-cases/${caseId}`);
  return response;
};

/**
 * 重新排序测试用例
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change'|'ai')
 * @param {Array<number>} caseIds - 排序后的用例ID数组
 * @returns {Promise<{new_ids: Array<number>}>}
 */
export const reorderCases = async (projectId, caseType, caseIds) => {
  const response = await client.post(`/projects/${projectId}/manual-cases/reorder`, {
    case_type: caseType,
    case_ids: caseIds
  });
  return response;
};

/**
 * 拖拽重新排序测试用例（根据case_id顺序重新分配ID�?
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change'|'ai')
 * @param {Array<string>} caseIDOrder - 排序后的case_id数组（UUID�?
 * @returns {Promise}
 */
export const reorderCasesByDrag = async (projectId, caseType, caseIDOrder) => {
  const response = await client.post(`/projects/${projectId}/manual-cases/reorder-drag`, {
    case_type: caseType,
    case_id_order: caseIDOrder
  });
  return response;
};

/**
 * 按现有ID顺序重新编号所有用例（用于重新排序按钮�?
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change'|'ai')
 * @param {string} language - 语言 ('中文'|'English'|'日本�?)
 * @returns {Promise<{count: number}>}
 */
export const reorderAllCasesByID = async (projectId, caseType, language) => {
  const response = await client.post(`/projects/${projectId}/manual-cases/reorder-all`, {
    case_type: caseType,
    language: language
  });
  return response;
};

/**
 * 清空AI测试用例
 * @param {number} projectId - 项目ID
 * @returns {Promise}
 */
export const clearAICases = async (projectId) => {
  const response = await client.delete(`/projects/${projectId}/manual-cases/clear-ai`);
  return response;
};

// ==================== 导出功能 ====================

/**
 * 导出AI用例(9列单Sheet)
 * @param {number} projectId - 项目ID
 * @returns {Promise<void>} 自动下载文件
 */
export const exportAICases = async (projectId) => {
  const response = await client.get(`/projects/${projectId}/manual-cases/export/ai`, {
    responseType: 'blob'
  });
  // response 拦截器对 blob 类型返回完整�?response 对象
  const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  // �?Content-Disposition header 中提取文件名
  const contentDisposition = response.headers['content-disposition'];
  let filename = 'AI用例.xlsx';
  if (contentDisposition) {
    const filenameMatch = contentDisposition.match(/filename=(.+)/);
    if (filenameMatch && filenameMatch[1]) {
      filename = decodeURIComponent(filenameMatch[1]);
    }
  }
  a.download = filename;
  a.click();
  window.URL.revokeObjectURL(url);
};

/**
 * 导出用例模板(23列空模板+示例行)
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change')
 * @returns {Promise<void>} 自动下载文件
 */
export const exportTemplate = async (projectId, caseType) => {
  const response = await client.get(`/projects/${projectId}/manual-cases/export/template`, {
    params: { caseType },
    responseType: 'blob'
  });
  const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  const contentDisposition = response.headers['content-disposition'];
  let filename = '用例模板.xlsx';
  if (contentDisposition) {
    const filenameMatch = contentDisposition.match(/filename=(.+)/);
    if (filenameMatch && filenameMatch[1]) {
      filename = decodeURIComponent(filenameMatch[1]);
    }
  }
  a.download = filename;
  a.click();
  window.URL.revokeObjectURL(url);
};

/**
 * 导出整体/变更用例(双Sheet: 元数据+23列数据)
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change')
 * @param {string} taskUuid - 可选：执行任务UUID，传入后导出25列（增加BugID/ExecutionRemark）
 * @returns {Promise<void>} 自动下载文件
 */
export const exportCases = async (projectId, caseType, taskUuid = null) => {
  const params = { caseType };
  if (taskUuid) {
    params.task_uuid = taskUuid;
  }
  const response = await client.get(`/projects/${projectId}/manual-cases/export/cases`, {
    params,
    responseType: 'blob'
  });
  const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  const contentDisposition = response.headers['content-disposition'];
  let filename = '用例数据.xlsx';
  if (contentDisposition) {
    const filenameMatch = contentDisposition.match(/filename=(.+)/);
    if (filenameMatch && filenameMatch[1]) {
      filename = decodeURIComponent(filenameMatch[1]);
    }
  }
  a.download = filename;
  a.click();
  window.URL.revokeObjectURL(url);
};

// ==================== 导入功能 ====================

/**
 * 导入用例(支持UUID匹配更新)
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change')
 * @param {File} file - 上传的Excel文件
 * @returns {Promise<{updateCount: number, insertCount: number}>}
 */
export const importCases = async (projectId, caseType, file) => {
  const formData = new FormData();
  formData.append('caseType', caseType);
  formData.append('file', file);
  const response = await client.post(`/projects/${projectId}/manual-cases/import`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  });
  // 拦截器已经返回了 response.data，所以这里直接返�?response
  return response;
};

// ==================== 版本管理 ====================

/**
 * 保存版本(导出并存储到服务�?
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change')
 * @returns {Promise<{message: string}>}
 */
export const saveVersion = async (projectId, caseType) => {
  const formData = new FormData();
  formData.append('caseType', caseType);
  const response = await client.post(`/projects/${projectId}/versions/save`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  });
  return response;
};

/**
 * 获取版本列表
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型(可�?'overall'/'change',为空返回所�?
 * @returns {Promise<Array<{id: number, project_id: number, case_type: string, filename: string, file_size: number, created_by: number, created_at: string}>>}
 */
export const getVersionList = async (projectId, caseType = '') => {
  const url = caseType 
    ? `/projects/${projectId}/versions?case_type=${caseType}`
    : `/projects/${projectId}/versions`;
  const response = await client.get(url);
  return response; // 拦截器已经返回了 response.data
};

/**
 * 下载指定版本文件
 * @param {number} projectId - 项目ID
 * @param {number} versionID - 版本ID
 * @returns {Promise<void>} 自动下载文件
 */
export const downloadVersion = async (projectId, versionID) => {
  const response = await client.get(`/projects/${projectId}/versions/${versionID}/download`, {
    responseType: 'blob'
  });
  const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = response.headers['content-disposition']?.split('filename=')[1] || '版本文件.xlsx';
  a.click();
  window.URL.revokeObjectURL(url);
};

/**
 * 删除指定版本(文件+数据库记�?
 * @param {number} projectId - 项目ID
 * @param {number} versionID - 版本ID
 * @returns {Promise<{message: string}>}
 */
export const deleteVersion = async (projectId, versionID) => {
  const response = await client.delete(`/projects/${projectId}/versions/${versionID}`);
  return response;
};

/**
 * 更新版本备注
 * @param {number} projectId - 项目ID
 * @param {number} versionID - 版本ID
 * @param {string} remark - 备注内容
 * @returns {Promise<{message: string}>}
 */
export const updateVersionRemark = async (projectId, versionID, remark) => {
  const response = await client.put(`/versions/${versionID}/remark`, { remark });
  return response;
};

/**
 * 重新分配所有用例的ID
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型
 * @returns {Promise<{message: string}>}
 */
export const reassignAllIDs = async (projectId, caseType) => {
  const response = await client.post(`/projects/${projectId}/manual-cases/reassign-ids`, {
    caseType
  });
  return response;
};

// ==================== 评审内容管理 ====================

/**
 * 获取评审内容
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change'|'ai')
 * @returns {Promise<{content: string}>}
 */
export const getCaseReview = async (projectId, caseType) => {
  const response = await client.get(`/projects/${projectId}/review`, {
    params: { caseType }
  });
  return response; // client 拦截器已经返回了 response.data
};

/**
 * 保存评审内容(UPSERT)
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change'|'ai')
 * @param {string} content - 评审内容(Markdown格式)
 * @returns {Promise<{message: string}>}
 */
export const saveCaseReview = async (projectId, caseType, content) => {
  const response = await client.post(`/projects/${projectId}/review`, {
    caseType,
    content
  });
  return response; // client 拦截器已经返回了 response.data
};

/**
 * 保存多语言版本（生成CN/JP/EN三个xlsx打包成zip）
 * @param {number} projectId - 项目ID
 * @returns {Promise<{filename: string, message: string}>}
 */
export const saveMultiLangVersion = async (projectId) => {
  const response = await client.post(`/projects/${projectId}/manual-cases/save-version`);
  return response;
};

/**
 * 导出手工测试用例多语言模版（CN/JP/EN空白xlsx打包成zip）
 * @returns {Promise<Blob>} - 返回zip文件Blob对象
 */
export const exportMultiLangTemplate = async () => {
  const response = await client.get('/manual-cases/template', {
    responseType: 'blob' // 重要：指定响应类型为blob
  });
  // client拦截器对blob类型返回完整response对象，需要取data
  return response.data;
};

// ==================== T44: 按语言导入导出 ====================

/**
 * 按语言导出用例
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change')
 * @param {string} language - 语言 ('中文'|'日本語'|'English')
 * @param {string} caseGroup - 用例集名称
 * @returns {Promise<void>}
 */
export const exportCasesByLanguage = async (projectId, caseType, language, caseGroup) => {
  // 语言映射
  const langMap = {
    '中文': 'CN',
    '日本語': 'JP',
    'English': 'EN'
  };
  
  const params = { 
    caseType,
    language: langMap[language] || 'CN',
    case_group: caseGroup
  };

  const response = await client.get(`/projects/${projectId}/manual-cases/export/cases`, {
    params,
    responseType: 'blob'
  });

  const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  
  const contentDisposition = response.headers['content-disposition'];
  let filename = `${caseGroup}_${langMap[language]}.xlsx`;
  if (contentDisposition) {
    const filenameMatch = contentDisposition.match(/filename=(.+)/);
    if (filenameMatch && filenameMatch[1]) {
      filename = decodeURIComponent(filenameMatch[1]);
    }
  }
  
  a.download = filename;
  a.click();
  window.URL.revokeObjectURL(url);
};

/**
 * 按语言导入用例
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change')
 * @param {File} file - 上传的Excel文件
 * @param {string} language - 语言 ('中文'|'日本語'|'English')
 * @param {string} caseGroup - 用例集名称
 * @returns {Promise<{updateCount: number, insertCount: number}>}
 */
export const importCasesByLanguage = async (projectId, caseType, file, language, caseGroup) => {
  // 语言映射
  const langMap = {
    '中文': 'CN',
    '日本語': 'JP',
    'English': 'EN'
  };

  console.log('[importCasesByLanguage] 🔍 开始导入:', {
    projectId,
    caseType,
    language,
    mappedLanguage: langMap[language] || 'CN',
    caseGroup,
    fileName: file.name,
    fileSize: file.size
  });

  const formData = new FormData();
  formData.append('caseType', caseType);
  formData.append('language', langMap[language] || 'CN');
  formData.append('case_group', caseGroup);
  formData.append('file', file);

  // 打印FormData内容（用于调试）
  console.log('[importCasesByLanguage] 📦 FormData内容:');
  for (let pair of formData.entries()) {
    if (pair[0] === 'file') {
      console.log(`  ${pair[0]}:`, pair[1].name, `(${pair[1].size} bytes)`);
    } else {
      console.log(`  ${pair[0]}:`, pair[1]);
    }
  }

  const response = await client.post(`/projects/${projectId}/manual-cases/import`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  });

  console.log('[importCasesByLanguage] ✅ 导入响应:', response);
  return response;
};

/**
 * 获取项目的用例集列表
 * @param {number} projectId - 项目ID
 * @param {string} caseType - 用例类型 ('overall'|'change'|'acceptance')
 * @returns {Promise<Array>} 用例集列表
 */
export const getCaseGroups = async (projectId, caseType = 'overall') => {
  const response = await client.get(`/projects/${projectId}/case-groups`, {
    params: { case_type: caseType }
  });
  return response;
};

/**
 * 创建新的用例集
 * @param {number} projectId - 项目ID
 * @param {Object} data - 用例集数据
 * @param {string} data.caseType - 用例类型
 * @param {string} data.groupName - 用例集名称
 * @param {string} data.description - 描述（可选）
 * @param {number} data.displayOrder - 显示顺序（可选）
 * @returns {Promise}
 */
export const createCaseGroup = async (projectId, data) => {
  const requestData = {
    case_type: data.caseType,
    group_name: data.groupName,
    description: data.description || '',
    display_order: data.displayOrder || 0
  };
  const response = await client.post(`/projects/${projectId}/case-groups`, requestData);
  return response;
};

/**
 * 更新用例集
 * @param {number} groupId - 用例集ID
 * @param {Object} data - 更新数据
 * @param {string} data.groupName - 用例集名称（可选）
 * @param {string} data.description - 描述（可选）
 * @param {number} data.displayOrder - 显示顺序（可选）
 * @param {string} data.metaProtocol - 协议（可选）
 * @param {string} data.metaServer - 服务器地址（可选）
 * @param {string} data.metaPort - 端口号（可选）
 * @param {string} data.metaUser - 用户名（可选）
 * @param {string} data.metaPassword - 密码（可选）
 * @returns {Promise}
 */
export const updateCaseGroup = async (groupId, data) => {
  const requestData = {};
  if (data.groupName !== undefined) requestData.group_name = data.groupName;
  if (data.description !== undefined) requestData.description = data.description;
  if (data.displayOrder !== undefined) requestData.display_order = data.displayOrder;
  if (data.metaProtocol !== undefined) requestData.meta_protocol = data.metaProtocol;
  if (data.metaServer !== undefined) requestData.meta_server = data.metaServer;
  if (data.metaPort !== undefined) requestData.meta_port = data.metaPort;
  if (data.metaUser !== undefined) requestData.meta_user = data.metaUser;
  if (data.metaPassword !== undefined) requestData.meta_password = data.metaPassword;
  const response = await client.put(`/case-groups/${groupId}`, requestData);
  return response;
};

/**
 * 删除用例集
 * @param {number} groupId - 用例集ID
 * @returns {Promise}
 */
export const deleteCaseGroup = async (groupId) => {
  const response = await client.delete(`/case-groups/${groupId}`);
  return response;
};

