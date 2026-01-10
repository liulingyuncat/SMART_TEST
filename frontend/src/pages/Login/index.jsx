import React, { useState, useEffect } from 'react';
import { Form, Input, Button, message, Select } from 'antd';
import { UserOutlined, LockOutlined, SettingOutlined } from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { useTranslation } from 'react-i18next';
import { login as loginAPI } from '../../api/auth';
import { login } from '../../store/authSlice';
import './index.css';

const { Option } = Select;

// 生成随机几何图形
const generateShapes = () => {
  const shapes = [];
  const count = 12; // 图形数量
  for (let i = 0; i < count; i++) {
    const size = Math.random() * 120 + 40; // 40-160px
    const left = Math.random() * 100;
    const top = Math.random() * 100;
    const opacity = Math.random() * 0.15 + 0.05; // 0.05-0.2
    const delay = Math.random() * 5;
    shapes.push({ id: i, size, left, top, opacity, delay });
  }
  return shapes;
};

const LoginPage = () => {
  const [loading, setLoading] = useState(false);
  const [shapes] = useState(generateShapes);
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useDispatch();
  const { t, i18n } = useTranslation();

  const handleLanguageChange = (lang) => {
    i18n.changeLanguage(lang);
    localStorage.setItem('language', lang);
  };

  const onFinish = async (values) => {
    setLoading(true);
    try {
      const response = await loginAPI(values.username, values.password);
      
      // 验证返回数据
      if (!response.user || !response.user.role) {
        throw new Error('Invalid response from server: missing user data');
      }
      
      // 存储 token 和用户信息
      localStorage.setItem('auth_token', response.token);
      localStorage.setItem('user_info', JSON.stringify(response.user));
      
      // 更新 Redux 状态
      dispatch(login({
        token: response.token,
        user: response.user,
      }));

      message.success(t('login.success'));

      // 跳转到目标页面或首页
      const from = new URLSearchParams(location.search).get('redirect') || '/';
      navigate(from);
    } catch (error) {
      console.error('Login error:', error);
      const errorMessage = error?.response?.data?.message || error?.message || t('login.failed');
      message.error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      {/* 背景几何图形 */}
      <div className="login-shapes">
        {shapes.map((shape) => (
          <div
            key={shape.id}
            className="login-shape"
            style={{
              width: shape.size,
              height: shape.size,
              left: `${shape.left}%`,
              top: `${shape.top}%`,
              opacity: shape.opacity,
              animationDelay: `${shape.delay}s`,
            }}
          />
        ))}
      </div>

      {/* 登录卡片 */}
      <div className="login-card">
        {/* Logo */}
        <div className="login-logo">
          <div className="login-logo-icon">
            <span>ST</span>
          </div>
        </div>

        {/* 标题 */}
        <div className="login-title">SMART TEST</div>
        <div className="login-subtitle">PEVVD Intelligent Test Platform</div>

        {/* 表单 */}
        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          size="large"
          className="login-form"
        >
          <div className="login-form-label">
            <span className="required">*</span> {t('login.usernameLabel')}
          </div>
          <Form.Item
            name="username"
            rules={[
              { required: true, message: t('login.usernameRequired') },
              { min: 3, max: 50, message: t('login.usernameLength') },
            ]}
          >
            <Input
              prefix={<UserOutlined style={{ color: '#bfbfbf' }} />}
              placeholder={t('login.usernamePlaceholder')}
            />
          </Form.Item>

          <div className="login-form-label">
            <span className="required">*</span> {t('login.passwordLabel')}
          </div>
          <Form.Item
            name="password"
            rules={[
              { required: true, message: t('login.passwordRequired') },
              { min: 6, max: 50, message: t('login.passwordLength') },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined style={{ color: '#bfbfbf' }} />}
              placeholder={t('login.passwordPlaceholder')}
            />
          </Form.Item>

          <Form.Item style={{ marginTop: 24 }}>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              className="login-submit-btn"
            >
              {t('login.submit')}
            </Button>
          </Form.Item>
        </Form>

        {/* 语言选择 */}
        <div className="login-language">
          <SettingOutlined style={{ marginRight: 8, color: '#8c8c8c' }} />
          <span style={{ color: '#8c8c8c', marginRight: 8 }}>{t('login.language')}</span>
          <Select
            value={i18n.language}
            onChange={handleLanguageChange}
            size="small"
            style={{ width: 100 }}
            bordered={false}
          >
            <Option value="zh">中文</Option>
            <Option value="en">English</Option>
            <Option value="ja">日本語</Option>
          </Select>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
