import axios from 'axios';
import { message } from 'antd';
import i18n from '../i18n';
import store from '../store';
import { logout } from '../store/authSlice';

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器 - 添加 JWT Token
apiClient.interceptors.request.use(
  (config) => {
    console.log('=== [apiClient请求拦截器] ===');
    console.log('[请求拦截器] method:', config.method?.toUpperCase());
    console.log('[请求拦截器] url:', config.url);
    console.log('[请求拦截器] baseURL:', config.baseURL);
    console.log('[请求拦截器] 完整URL:', `${config.baseURL}${config.url}`);
    console.log('[请求拦截器] params:', config.params);
    console.log('[请求拦截器] 请求数据:', config.data);
    
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
      console.log('[请求拦截器] Token已添加');
    } else {
      console.warn('[请求拦截器] 未找到Token');
    }
    
    // 如果数据是FormData，不设置Content-Type（让浏览器自动添加boundary）
    if (!(config.data instanceof FormData)) {
      config.headers['Content-Type'] = 'application/json';
    } else {
      // 删除Content-Type header，让axios自动处理
      delete config.headers['Content-Type'];
      console.log('[请求拦截器] FormData检测到，允许浏览器自动处理Content-Type');
    }
    
    console.log('[请求拦截器] 最终headers:', config.headers);
    return config;
  },
  (error) => {
    console.error('[请求拦截器] 错误:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器 - 统一错误处理
apiClient.interceptors.response.use(
  (response) => {
    console.log('=== [apiClient响应拦截器] 收到响应 ===');
    console.log('[响应拦截器] status:', response.status);
    console.log('[响应拦截器] statusText:', response.statusText);
    console.log('[响应拦截器] response.data:', response.data);
    console.log('[响应拦截器] responseType:', response.config.responseType);
    console.log('[响应拦截器] response.headers:', response.headers);
    
    // 对于 blob 类型的响应（文件下载），返回完整的 response 对象以访问 headers
    if (response.config.responseType === 'blob') {
      console.log('[响应拦截器] blob类型，返回完整response');
      return response;
    }
    // 后端统一返回格式: {code: 0, message: "success", data: {...}}
    // 直接返回 data 字段
    const result = response.data?.data !== undefined ? response.data.data : response.data;
    console.log('[响应拦截器] 处理后返回:', result);
    return result;
  },
  (error) => {
    console.error('=== [apiClient响应拦截器] 请求失败 ===');
    console.error('[响应拦截器] error对象:', error);
    console.error('[响应拦截器] error.message:', error.message);
    
    const { response } = error;

    if (response) {
      console.error('[响应拦截器] 服务器响应错误');
      console.error('[响应拦截器] response.status:', response.status);
      console.error('[响应拦截器] response.data:', response.data);
      
      // 服务器返回错误状态码
      switch (response.status) {
        case 401:
          message.error(i18n.t('message.sessionExpired'));
          localStorage.removeItem('auth_token');
          store.dispatch(logout());
          window.location.href = '/login';
          break;

        case 403:
          message.error(i18n.t('message.forbidden'));
          break;

        case 404:
          message.error(i18n.t('message.notFound'));
          break;

        case 409:
          // 冲突错误（如重复数据），不在这里显示消息，让调用方处理
          break;

        case 500:
          message.error(i18n.t('message.serverError'));
          break;

        default:
          message.error(response.data?.message || i18n.t('message.requestFailed'));
      }
    } else {
      console.error('[响应拦截器] 网络错误或请求未发出');
      // 网络错误
      message.error(i18n.t('message.networkError'));
    }

    return Promise.reject(error);
  }
);

export default apiClient;
