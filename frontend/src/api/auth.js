import apiClient from './client';

/**
 * 用户登录
 * @param {string} username - 用户�?
 * @param {string} password - 密码
 * @returns {Promise<{token: string, user: object}>}
 */
export const login = async (username, password) => {
  const response = await apiClient.post('/auth/login', {
    username,
    password,
  });
  return response;
};

/**
 * 获取当前用户信息
 * @returns {Promise<object>}
 */
export const getCurrentUser = async () => {
  const response = await apiClient.get('/profile');
  return response;
};

