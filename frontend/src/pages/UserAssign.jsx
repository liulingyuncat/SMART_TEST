import { useState, useEffect } from 'react';
import { Card, Row, Col, Select, Button, Spin, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { getProjects, getProjectMembers, updateProjectMembers } from '../api/project';
import { getUsers } from '../api/user';
import { getCurrentUser } from '../api/auth';
import MemberTransfer from '../components/MemberTransfer';

const { Option } = Select;

const UserAssign = () => {
  const { t } = useTranslation();
  const [projects, setProjects] = useState([]);
  const [selectedProjectId, setSelectedProjectId] = useState(null);
  const [allUsers, setAllUsers] = useState([]);
  const [currentUserId, setCurrentUserId] = useState(null);
  const [managerTargetKeys, setManagerTargetKeys] = useState([]);
  const [memberTargetKeys, setMemberTargetKeys] = useState([]);
  const [originalManagerKeys, setOriginalManagerKeys] = useState([]);
  const [originalMemberKeys, setOriginalMemberKeys] = useState([]);
  const [loading, setLoading] = useState(false);
  const [loadingMembers, setLoadingMembers] = useState(false);
  const [saving, setSaving] = useState(false);

  // 加载初始数据
  useEffect(() => {
    const fetchInitialData = async () => {
      setLoading(true);
      try {
        // 获取项目列表
        const projectsData = await getProjects();
        const projectList = Array.isArray(projectsData) ? projectsData : [];
        setProjects(projectList);
        
        // 默认选择第一个项目
        if (projectList.length > 0 && !selectedProjectId) {
          setSelectedProjectId(projectList[0].id);
        }

        // 获取所有用户
        const usersData = await getUsers();
        console.log('[UserAssign] 原始用户数据:', usersData);
        // 后端返回格式: {users: [...], total: 10}
        const userList = usersData?.users || usersData;
        console.log('[UserAssign] 提取的用户列表:', userList);
        console.log('[UserAssign] 用户列表是否为数组:', Array.isArray(userList));
        setAllUsers(Array.isArray(userList) ? userList : []);

        // 获取当前用户信息
        console.log('[UserAssign] 🔍 开始获取当前用户信息...');
        const currentUser = await getCurrentUser();
        console.log('[UserAssign] 🔍 getCurrentUser API返回:', currentUser);
        console.log('[UserAssign] 🔍 currentUser 对象的所有键:', Object.keys(currentUser || {}));
        
        // 后端返回的字段可能是 ID、id、user_id 或其他变体
        const userId = currentUser?.ID || currentUser?.id || currentUser?.user_id;
        const userRole = currentUser?.Role || currentUser?.role;
        
        console.log('[UserAssign] 🔍 提取的用户ID:', userId, '角色:', userRole);
        console.log('[UserAssign] 🔍 用户ID类型:', typeof userId);
        
        if (!userId) {
          console.error('[UserAssign] ❌ 错误: 无法获取当前用户ID!', currentUser);
          message.error('无法获取当前用户信息，请重新登录');
        } else {
          setCurrentUserId(userId);
          console.log('[UserAssign] ✅ 成功设置当前用户ID:', userId);
        }
      } catch (error) {
        console.error('[UserAssign] ❌ 加载数据失败:', error);
        message.error(t('assign.loadProjectsError'));
      } finally {
        setLoading(false);
      }
    };
    fetchInitialData();
  }, [t]);

  // 加载项目成员并初始化穿梭框
  useEffect(() => {
    const fetchMembers = async () => {
      if (!selectedProjectId) {
        setManagerTargetKeys([]);
        setMemberTargetKeys([]);
        return;
      }
      setLoadingMembers(true);
      try {
        const data = await getProjectMembers(selectedProjectId);
        console.log('[UserAssign] 项目成员数据:', data);
        console.log('[UserAssign] managers原始数据:', data.managers);
        console.log('[UserAssign] members原始数据:', data.members);
        
        // 分离管理员和成员
        const managers = (data.managers || []).map((m) => m.user_id || m.ID || m.id);
        const members = (data.members || []).map((m) => m.user_id || m.ID || m.id);
        
        console.log('[UserAssign] 提取的manager IDs:', managers);
        console.log('[UserAssign] 提取的member IDs:', members);
        
        setManagerTargetKeys(managers);
        setMemberTargetKeys(members);
        setOriginalManagerKeys([...managers]);
        setOriginalMemberKeys([...members]);
      } catch (error) {
        message.error(t('assign.loadMembersError'));
        setManagerTargetKeys([]);
        setMemberTargetKeys([]);
      } finally {
        setLoadingMembers(false);
      }
    };
    fetchMembers();
  }, [selectedProjectId, t]);

  // 项目选择变化处理
  const handleProjectChange = (value) => {
    setSelectedProjectId(value);
  };

  // 管理员穿梭框变化处理
  const handleManagerChange = (newTargetKeys) => {
    setManagerTargetKeys(newTargetKeys);
  };

  // 成员穿梭框变化处理
  const handleMemberChange = (newTargetKeys) => {
    setMemberTargetKeys(newTargetKeys);
  };

  // 保存按钮点击处理
  const handleSave = async () => {
    if (!selectedProjectId) {
      message.warning(t('assign.selectProjectFirst'));
      return;
    }

    setSaving(true);
    try {
      // 确保当前管理员用户始终在managers列表中（后端要求）
      let finalManagerKeys = [...managerTargetKeys];
      if (isCurrentUserManager && !finalManagerKeys.includes(currentUserId)) {
        finalManagerKeys.push(currentUserId);
        console.log('[UserAssign] 自动添加当前管理员到managers列表:', currentUserId);
      }
      
      const requestData = {
        managers: finalManagerKeys,
        members: memberTargetKeys,
      };

      await updateProjectMembers(selectedProjectId, requestData);
      message.success(t('assign.saveSuccess'));
      
      // 刷新数据以显示最新状态
      const data = await getProjectMembers(selectedProjectId);
      const managers = (data.managers || []).map((m) => m.user_id || m.ID || m.id);
      const members = (data.members || []).map((m) => m.user_id || m.ID || m.id);
      setManagerTargetKeys(managers);
      setMemberTargetKeys(members);
      setOriginalManagerKeys([...managers]);
      setOriginalMemberKeys([...members]);
    } catch (error) {
      const errorMessage = error.message || t('assign.saveFailed');
      message.error(errorMessage);
    } finally {
      setSaving(false);
    }
  };

  // 按角色过滤用户（确保allUsers是数组）
  const safeAllUsers = Array.isArray(allUsers) ? allUsers : [];
  console.log('[UserAssign] 渲染时allUsers:', allUsers);
  console.log('[UserAssign] safeAllUsers长度:', safeAllUsers.length);
  
  // 数据转换：后端返回的字段名是大写开头(ID, Username, Nickname, Role)
  // 需要转换为小写下划线格式(user_id, username, nickname, role)
  const normalizedUsers = safeAllUsers.map((user) => {
    const normalized = {
      user_id: user.ID || user.id || user.user_id,
      username: user.Username || user.username,
      nickname: user.Nickname || user.nickname,
      role: user.Role || user.role,
    };
    console.log('[UserAssign] 用户数据转换:', user, '->', normalized);
    return normalized;
  });

  // 判断当前用户的系统角色是否为项目管理员
  // 注意：只要系统角色是 project_manager，就应该锁定，防止管理员把自己移出项目
  console.log('[UserAssign] 查找当前用户信息, currentUserId:', currentUserId);
  const currentUser = normalizedUsers.find(u => u.user_id === currentUserId);
  console.log('[UserAssign] 找到的当前用户对象:', currentUser);
  
  const currentUserRole = currentUser?.role;
  const isCurrentUserManager = currentUserRole === 'project_manager';
  console.log('[UserAssign] 当前用户是否为管理员:', isCurrentUserManager, 'currentUserId:', currentUserId, 'role:', currentUserRole);
  
  // 锁定用户ID数组（如果当前用户系统角色是管理员则锁定，防止自己把自己移出项目）
  const lockedKeys = isCurrentUserManager ? [currentUserId] : [];
  console.log('[UserAssign] ⚠️ 重要: 锁定的用户IDs:', lockedKeys, '类型:', typeof currentUserId);
  
  console.log('[UserAssign] normalizedUsers:', normalizedUsers);
  console.log('[UserAssign] normalizedUsers中的所有角色:', normalizedUsers.map(u => u.role));
  
  // 项目管理员穿梭框：显示系统角色为 project_manager 的用户
  const projectManagers = normalizedUsers.filter((user) => {
    const match = user.role === 'project_manager';
    console.log(`[UserAssign] 检查用户 ${user.username} (${user.user_id}) 角色=${user.role}, 是否为project_manager: ${match}`);
    return match;
  });
  
  // 项目成员穿梭框：显示系统角色为 project_member 的用户
  const projectMembers = normalizedUsers.filter((user) => {
    const match = user.role === 'project_member';
    console.log(`[UserAssign] 检查用户 ${user.username} (${user.user_id}) 角色=${user.role}, 是否为project_member: ${match}`);
    return match;
  });
  
  console.log('[UserAssign] 过滤后的projectManagers:', projectManagers);
  console.log('[UserAssign] 过滤后的projectMembers:', projectMembers);
  console.log('[UserAssign] managerTargetKeys:', managerTargetKeys);
  console.log('[UserAssign] memberTargetKeys:', memberTargetKeys);

  // 检查是否有变更
  const hasChanges = JSON.stringify(managerTargetKeys.sort()) !== JSON.stringify(originalManagerKeys.sort()) ||
                     JSON.stringify(memberTargetKeys.sort()) !== JSON.stringify(originalMemberKeys.sort());

  return (
    <Card title={t('menu.assign')}>
      <Row gutter={16} align="middle">
        <Col>
          <label style={{ marginRight: 8, fontWeight: 'bold' }}>
            {t('assign.selectProject')}:
          </label>
          <Select
            showSearch
            placeholder={t('assign.selectProjectPlaceholder')}
            style={{ width: 300 }}
            value={selectedProjectId}
            onChange={handleProjectChange}
            loading={loading}
            filterOption={(input, option) =>
              option.children.toLowerCase().includes(input.toLowerCase())
            }
          >
            {Array.isArray(projects) && projects.map((project) => (
              <Option key={project.id} value={project.id}>
                {project.name}
              </Option>
            ))}
          </Select>
        </Col>
        {selectedProjectId && hasChanges && (
          <Col>
            <Button
              type="primary"
              onClick={handleSave}
              loading={saving}
            >
              {t('assign.saveMembers')}
            </Button>
          </Col>
        )}
      </Row>
      
      {selectedProjectId && isCurrentUserManager && (
        <div style={{ marginTop: 8, color: '#ff9800', fontWeight: 'bold' }}>
          ⚠️ {t('assign.lockedUser')}
        </div>
      )}

      {selectedProjectId && (
        <Spin spinning={loadingMembers}>
          <Row gutter={8} style={{ marginTop: 8 }}>
            <Col span={12}>
              <h3>{t('assign.managerTransferTitle')}</h3>
              <MemberTransfer
                dataSource={projectManagers}
                targetKeys={managerTargetKeys}
                lockedKeys={lockedKeys}
                onChange={handleManagerChange}
                title={[t('assign.availableUsers'), t('assign.assignedUsers')]}
              />
            </Col>
            <Col span={12}>
              <h3>{t('assign.memberTransferTitle')}</h3>
              <MemberTransfer
                dataSource={projectMembers}
                targetKeys={memberTargetKeys}
                lockedKeys={[]}
                onChange={handleMemberChange}
                title={[t('assign.availableUsers'), t('assign.assignedUsers')]}
              />
            </Col>
          </Row>
        </Spin>
      )}
    </Card>
  );
};

export default UserAssign;
