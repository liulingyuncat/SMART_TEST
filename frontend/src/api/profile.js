import apiClient from './client';

/**
 * 获取当前用户个人信息
 * @returns {Promise<{user_id: number, username: string, nickname: string, role: string}>}
 */
export async function getProfile() {
  const response = await apiClient.get('/profile');
  // apiClient响应拦截器已经返回response.data.data，直接返回response
  return response;
}

/**
 * 更新当前用户昵称
 * @param {string} nickname - 新昵称 (2-50字符)
 * @returns {Promise<{user_id: number, username: string, nickname: string, role: string}>}
 */
export async function updateNickname(nickname) {
  const response = await apiClient.put('/profile/nickname', { nickname });
  return response;
}

/**
 * 修改当前用户密码 - T23
 * @param {string} currentPassword - 当前密码
 * @param {string} newPassword - 新密码 (6-50字符)
 * @returns {Promise<{message: string}>}
 */
export async function changePassword(currentPassword, newPassword) {
  const response = await apiClient.put('/profile/password', {
    current_password: currentPassword,
    new_password: newPassword,
  });
  return response;
}

/**
 * 生成API Token - T23
 * @returns {Promise<{token: string}>}
 */
export async function generateToken() {
  const response = await apiClient.post('/profile/token');
  return response;
}

/**
 * 获取Token状态 - T23
 * @returns {Promise<{has_token: boolean}>}
 */
export async function getTokenStatus() {
  const response = await apiClient.get('/profile/token');
  return response;
}

/**
 * 设置当前用户选择的当前项目 - T50
 * @param {number} projectId - 项目ID
 * @returns {Promise<{project_id: number}>}
 */
export async function setCurrentProject(projectId) {
  const response = await apiClient.put('/profile/current-project', {
    project_id: projectId,
  });
  return response;
}
