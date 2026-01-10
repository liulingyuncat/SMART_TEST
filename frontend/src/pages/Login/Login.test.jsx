import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import { configureStore } from '@reduxjs/toolkit';
import LoginPage from '../index';
import authReducer from '../../../store/authSlice';
import * as authAPI from '../../../api/auth';

// Mock API
jest.mock('../../../api/auth');

// 创建测试 store
const createTestStore = () => {
  return configureStore({
    reducer: {
      auth: authReducer,
    },
  });
};

// 渲染组件辅助函数
const renderLoginPage = () => {
  const store = createTestStore();
  return render(
    <Provider store={store}>
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>
    </Provider>
  );
};

describe('LoginPage Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    localStorage.clear();
  });

  test('renders login form correctly', () => {
    renderLoginPage();
    
    expect(screen.getByPlaceholderText(/请输入用户名/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/请输入密码/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /登录/i })).toBeInTheDocument();
  });

  test('validates required fields', async () => {
    renderLoginPage();
    
    const submitButton = screen.getByRole('button', { name: /登录/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/请输入用户名/i)).toBeInTheDocument();
      expect(screen.getByText(/请输入密码/i)).toBeInTheDocument();
    });
  });

  test('validates username length', async () => {
    renderLoginPage();
    
    const usernameInput = screen.getByPlaceholderText(/请输入用户名/i);
    fireEvent.change(usernameInput, { target: { value: 'ab' } });
    
    const submitButton = screen.getByRole('button', { name: /登录/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/用户名长度为 3-50 个字符/i)).toBeInTheDocument();
    });
  });

  test('submits form successfully', async () => {
    const mockResponse = {
      token: 'test-token-123',
      user: { username: 'admin', role: 'admin' },
    };
    authAPI.login.mockResolvedValue(mockResponse);

    renderLoginPage();
    
    const usernameInput = screen.getByPlaceholderText(/请输入用户名/i);
    const passwordInput = screen.getByPlaceholderText(/请输入密码/i);
    const submitButton = screen.getByRole('button', { name: /登录/i });

    fireEvent.change(usernameInput, { target: { value: 'admin' } });
    fireEvent.change(passwordInput, { target: { value: 'admin123' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(authAPI.login).toHaveBeenCalledWith('admin', 'admin123');
      expect(localStorage.getItem('auth_token')).toBe('test-token-123');
    });
  });

  test('handles login failure', async () => {
    authAPI.login.mockRejectedValue('用户名或密码错误');

    renderLoginPage();
    
    const usernameInput = screen.getByPlaceholderText(/请输入用户名/i);
    const passwordInput = screen.getByPlaceholderText(/请输入密码/i);
    const submitButton = screen.getByRole('button', { name: /登录/i });

    fireEvent.change(usernameInput, { target: { value: 'admin' } });
    fireEvent.change(passwordInput, { target: { value: 'wrongpassword' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(authAPI.login).toHaveBeenCalled();
      expect(localStorage.getItem('auth_token')).toBeNull();
    });
  });

  test('changes language', () => {
    renderLoginPage();
    
    const languageSelect = screen.getByRole('combobox');
    fireEvent.mouseDown(languageSelect);
    
    const englishOption = screen.getByText('English');
    fireEvent.click(englishOption);

    // 验证语言切换后的文本变化
    waitFor(() => {
      expect(screen.getByPlaceholderText(/Please enter username/i)).toBeInTheDocument();
    });
  });
});
