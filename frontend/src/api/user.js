import apiClient from './client';

// 获取用户列表
export const getUsers = () => {
	return apiClient.get('/users');
};

// 创建用户
export const createUser = (data) => {
	return apiClient.post('/users', data);
};

// 更新昵称
export const updateNickname = (userId, nickname) => {
	return apiClient.put(`/users/${userId}`, { nickname });
};

// 删除用户
export const deleteUser = (userId) => {
	return apiClient.delete(`/users/${userId}`);
};

// 重置密码
export const resetPassword = (userId) => {
	return apiClient.post(`/users/${userId}/reset-password`);
};

// 检查唯一性（合并接口�?
export const checkUnique = (params) => {
	return apiClient.get('/users/check', { params });
};

