import apiClient from './client';

/**
 * 获取项目列表
 * @returns {Promise} 返回项目列表数组
 */
export const getProjects = () => {
  return apiClient.get('/projects');
};

/**
 * 创建项目
 * @param {Object} projectData - 项目数据 {name, description}
 * @returns {Promise} 返回新创建的项目对象
 */
export const createProject = (projectData) => {
  return apiClient.post('/projects', projectData);
};

/**
 * 更新项目
 * @param {number} projectId - 项目ID
 * @param {Object} projectData - 项目数据 {name}
 * @returns {Promise} 返回更新后的项目对象
 */
export const updateProject = (projectId, projectData) => {
  return apiClient.put(`/projects/${projectId}`, projectData);
};

/**
 * 删除项目
 * @param {number} projectId - 项目ID
 * @returns {Promise} 返回删除结果
 */
export const deleteProject = (projectId) => {
  return apiClient.delete(`/projects/${projectId}`);
};

/**
 * 获取项目详情
 * @param {number} projectId - 项目ID
 * @returns {Promise} 返回项目对象和用户角�?
 */
export const getProjectById = (projectId) => {
  return apiClient.get(`/projects/${projectId}`);
};

/**
 * 获取项目成员列表
 * @param {number} projectId - 项目ID
 * @returns {Promise} 返回包含成员列表的对象 {total, members: [{user_id, username, nickname, role}]}
 */
export const getProjectMembers = (projectId) => {
  return apiClient.get(`/projects/${projectId}/members`);
};

/**
 * 批量更新项目成员
 * @param {number} projectId - 项目ID
 * @param {Object} data - 成员数据 {managers: number[], members: number[]}
 * @returns {Promise} 返回更新后的成员列表对象 {managers: User[], members: User[]}
 */
export const updateProjectMembers = (projectId, data) => {
  return apiClient.put(`/projects/${projectId}/members`, data);
};

