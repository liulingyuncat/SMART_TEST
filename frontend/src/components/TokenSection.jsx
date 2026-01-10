import React, { useState, useEffect, useCallback } from 'react';
import { Button, message, Alert, Space, Typography, Tooltip } from 'antd';
import { CopyOutlined, CheckOutlined, KeyOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { generateToken, getTokenStatus } from '../api/profile';

const { Text, Paragraph } = Typography;

/**
 * Token管理区域组件 - T23
 * @param {Object} props
 * @param {string} props.role - 当前用户角色
 */
const TokenSection = ({ role }) => {
  const { t } = useTranslation();
  const [hasToken, setHasToken] = useState(false);
  const [token, setToken] = useState(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);

  // 加载Token状态
  const loadTokenStatus = useCallback(async () => {
    try {
      const response = await getTokenStatus();
      // apiClient响应拦截器已经返回response.data.data，所以response就是{ has_token: boolean }
      if (response && response.has_token !== undefined) {
        setHasToken(response.has_token);
      }
    } catch (error) {
      console.error('获取Token状态失败:', error);
    } finally {
      setInitialLoading(false);
    }
  }, []);

  useEffect(() => {
    if (role !== 'system_admin') {
      loadTokenStatus();
    } else {
      setInitialLoading(false);
    }
  }, [role, loadTokenStatus]);

  // 生成Token
  const handleGenerateToken = async () => {
    setLoading(true);
    try {
      const response = await generateToken();
      // apiClient响应拦截器已经返回response.data.data，所以response就是{ token: string }
      const newToken = response?.token;
      if (newToken) {
        setToken(newToken);
        setHasToken(true);
        message.success(t('profile.tokenGenerateSuccess') || 'Token生成成功');
      } else {
        throw new Error('Token generation returned empty response');
      }
    } catch (error) {
      console.error('生成Token失败:', error);
      message.error(t('profile.tokenGenerateFailed') || 'Token生成失败');
    } finally {
      setLoading(false);
    }
  };

  // 复制Token
  const handleCopyToken = async () => {
    if (!token) return;
    try {
      await navigator.clipboard.writeText(token);
      setCopied(true);
      message.success(t('profile.tokenCopySuccess') || 'Token已复制到剪贴板');
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('复制失败:', error);
      message.error(t('profile.tokenCopyFailed') || '复制失败');
    }
  };

  // 系统管理员无权限
  if (role === 'system_admin') {
    return (
      <Alert
        message={t('profile.tokenNoPermission') || '系统管理员无权生成API Token'}
        type="warning"
        showIcon
      />
    );
  }

  if (initialLoading) {
    return <div style={{ padding: '20px', textAlign: 'center' }}>加载中...</div>;
  }

  return (
    <div className="token-section">
      <div className="token-description">
        {t('profile.tokenDescription') || 
          'API Token 用于第三方应用访问系统接口。生成新Token后，旧Token将自动失效。请妥善保管您的Token。'}
      </div>

      {token ? (
        <div style={{ marginBottom: 16 }}>
          <Alert
            message={t('profile.tokenWarning') || '请立即复制并保存此Token，关闭后将无法再次查看'}
            type="warning"
            showIcon={false}
            style={{ marginBottom: 12, fontSize: '12px', background: 'none', border: 'none', padding: 0 }}
          />
          <div style={{ marginBottom: 8 }}>
            <Text strong style={{ fontSize: '13px' }}>{t('profile.yourToken') || '您的Token：'}</Text>
          </div>
          <div className="profile-container .token-display">
            {token}
          </div>
          <Tooltip title={copied ? (t('profile.copied') || '已复制') : (t('profile.copyToken') || '复制Token')}>
            <Button
              icon={copied ? <CheckOutlined /> : <CopyOutlined />}
              onClick={handleCopyToken}
              type={copied ? 'default' : 'primary'}
              size="small"
            >
              {copied ? (t('profile.copied') || '已复制') : (t('profile.copyToken') || '复制Token')}
            </Button>
          </Tooltip>
        </div>
      ) : (
        <div>
          {hasToken && (
            <Alert
              message={t('profile.tokenExists') || '您已生成过Token，重新生成将使旧Token失效'}
              type="info"
              showIcon={false}
              style={{ marginBottom: 16, fontSize: '12px', background: 'none', border: 'none', padding: 0 }}
            />
          )}
          <Button
            type="primary"
            icon={<KeyOutlined />}
            loading={loading}
            onClick={handleGenerateToken}
            size="small"
          >
            {hasToken 
              ? (t('profile.regenerateToken') || '重新生成Token')
              : (t('profile.generateToken') || '生成Token')
            }
          </Button>
        </div>
      )}
    </div>
  );
};

export default TokenSection;
