import React, { useState, useEffect, useCallback } from 'react';
import { Card, Button, Spin, message } from 'antd';
import { EditOutlined, LockOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getProfile } from '../api/profile';
import NicknameEditModal from '../components/NicknameEditModal';
import PasswordEditModal from '../components/PasswordEditModal';
import TokenSection from '../components/TokenSection';
import './Profile.css';

/**
 * 个人信息页面
 * 显示当前登录用户的用户名、昵称、角色信息
 * 支持编辑昵称、修改密码、管理API Token
 */
const Profile = () => {
  const { t } = useTranslation();
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [nicknameModalVisible, setNicknameModalVisible] = useState(false);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);

  // 加载用户信息
  const loadProfile = useCallback(async () => {
    setLoading(true);
    try {
      const response = await getProfile();
      if (response) {
        setProfile(response);
      }
    } catch (error) {
      console.error('获取用户信息失败:', error);
      message.error(t('message.loadProfileFailed') || '获取用户信息失败');
    } finally {
      setLoading(false);
    }
  }, [t]);

  useEffect(() => {
    loadProfile();
  }, [loadProfile]);

  // 角色标签转换
  const getRoleLabel = (role) => {
    switch (role) {
      case 'system_admin':
        return t('user.systemAdmin') || '系统管理员';
      case 'project_manager':
        return t('user.projectManager') || '项目管理员';
      case 'project_member':
        return t('user.projectMember') || '项目成员';
      default:
        return role;
    }
  };

  // 编辑昵称成功后的回调
  const handleNicknameUpdateSuccess = () => {
    setNicknameModalVisible(false);
    loadProfile();
  };

  // 修改密码成功后的回调
  const handlePasswordUpdateSuccess = () => {
    setPasswordModalVisible(false);
  };

  if (loading) {
    return (
      <div className="profile-container">
        <Card className="profile-loading-card">
          <div className="profile-loading">
            <Spin size="large" />
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className="profile-container">
      {/* 个人信息卡片 */}
      <Card
        title={<span className="profile-card-title">{t('profile.title') || '个人信息'}</span>}
      >
        <div className="profile-info-container">
          <div className="profile-info-row">
            <div className="profile-info-label">{t('profile.username') || '用户名'}</div>
            <div className="profile-info-value">{profile?.username}</div>
          </div>
          <div className="profile-info-row">
            <div className="profile-info-label">{t('profile.role') || '角色'}</div>
            <div className="profile-info-value">{getRoleLabel(profile?.role)}</div>
          </div>
          <div className="profile-info-row">
            <div className="profile-info-label">{t('profile.nickname') || '昵称'}</div>
            <div className="profile-info-value-with-action">
              <span>{profile?.nickname || '-'}</span>
              <Button 
                type="text" 
                icon={<EditOutlined />} 
                onClick={() => setNicknameModalVisible(true)} 
                size="small"
                title={t('profile.editNickname') || '编辑昵称'}
                className="profile-action-btn"
              />
            </div>
          </div>
        </div>
      </Card>

      {/* 安全设置卡片 */}
      <Card
        title={<span className="profile-card-title">{t('profile.security') || '安全设置'}</span>}
      >
        <div className="profile-security-item">
          <div className="profile-security-text">
            <div className="profile-security-title-with-action">
              <span>{t('profile.changePassword') || '变更密码'}</span>
              <Button
                type="text"
                icon={<LockOutlined />}
                onClick={() => setPasswordModalVisible(true)}
                size="small"
                title={t('profile.changePassword') || '变更密码'}
                className="profile-action-btn"
              />
            </div>
            <div className="profile-security-hint">
              {t('profile.passwordHint') || '定期更换密码可以保护账户安全'}
            </div>
          </div>
        </div>
      </Card>

      {/* API Token管理卡片 */}
      <Card
        title={<span className="profile-card-title">{t('profile.apiToken') || 'API Token'}</span>}
      >
        <TokenSection role={profile?.role} />
      </Card>

      <NicknameEditModal
        visible={nicknameModalVisible}
        currentNickname={profile?.nickname}
        onCancel={() => setNicknameModalVisible(false)}
        onSuccess={handleNicknameUpdateSuccess}
      />

      <PasswordEditModal
        visible={passwordModalVisible}
        onCancel={() => setPasswordModalVisible(false)}
        onSuccess={handlePasswordUpdateSuccess}
      />
    </div>
  );
};

export default Profile;
