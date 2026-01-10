import React, { useState, useEffect } from 'react';
import { Table, Button, Popconfirm, message, Row, Col, Card, Input, Tabs, Select, Spin } from 'antd';
import { PlusOutlined, DeleteOutlined, KeyOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons';
import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import CreateUserModal from './CreateUserModal';
import { getUsers, updateNickname, deleteUser, resetPassword } from '../api/user';
import { getProjects, getProjectMembers, updateProjectMembers } from '../api/project';
import { getCurrentUser } from '../api/auth';
import MemberTransfer from '../components/MemberTransfer';
import './UserManage.css';

// 可编辑单元格组件
const EditableCell = ({ value, onSave }) => {
	const [editing, setEditing] = useState(false);
	const [inputValue, setInputValue] = useState(value);

	const handleSave = async () => {
		if (inputValue === value) {
			setEditing(false);
			return;
		}
		try {
			await onSave(inputValue);
			setEditing(false);
		} catch (error) {
			setInputValue(value);
			setEditing(false);
		}
	};

	if (editing) {
		return (
			<Input
				value={inputValue}
				onChange={(e) => setInputValue(e.target.value)}
				onPressEnter={handleSave}
				onBlur={handleSave}
				onKeyDown={(e) => e.key === 'Escape' && (setInputValue(value), setEditing(false))}
				autoFocus
			/>
		);
	}

	return (
		<div
			onClick={() => setEditing(true)}
			style={{ cursor: 'pointer', padding: '4px', minHeight: '22px' }}
		>
			{value}
		</div>
	);
};

const UserManage = () => {
	const { t } = useTranslation();
	const [projectManagers, setProjectManagers] = useState([]);
	const [projectMembers, setProjectMembers] = useState([]);
	const [loading, setLoading] = useState(false);
	const [modalVisible, setModalVisible] = useState(false);
	const [createRole, setCreateRole] = useState('');
	
	// 获取当前用户角色
	const { user } = useSelector(state => state.auth);
	const currentUserRole = user?.role;

	// 人员分配页面相关状态
	const [projects, setProjects] = useState([]);
	const [selectedProjectId, setSelectedProjectId] = useState(null);
	const [allUsers, setAllUsers] = useState([]);
	const [currentUserId, setCurrentUserId] = useState(null);
	const [managerTargetKeys, setManagerTargetKeys] = useState([]);
	const [memberTargetKeys, setMemberTargetKeys] = useState([]);
	const [originalManagerKeys, setOriginalManagerKeys] = useState([]);
	const [originalMemberKeys, setOriginalMemberKeys] = useState([]);
	const [loadingMembers, setLoadingMembers] = useState(false);
	const [saving, setSaving] = useState(false);

	const loadUsers = async () => {
		setLoading(true);
		try {
			const response = await getUsers();
			console.log('[UserManage] API Response:', response);
			// apiClient已经提取了data,直接使用response.users
			const users = response.users || [];
			console.log('[UserManage] Users:', users);
			const managers = users.filter(u => u.role === 'project_manager');
			const members = users.filter(u => u.role === 'project_member');
			console.log('[UserManage] Managers:', managers);
			console.log('[UserManage] Members:', members);
			setProjectManagers(managers);
			setProjectMembers(members);
		} catch (error) {
			console.error('[UserManage] Load users error:', error);
			message.error(t('message.loadUsersFailed') + ': ' + (error.response?.data?.message || error.message));
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		loadUsers();
	}, []);

	const handleCreate = (role) => {
		setCreateRole(role);
		setModalVisible(true);
	};

	const handleNicknameChange = async (userId, newNickname) => {
		try {
			await updateNickname(userId, newNickname);
			message.success(t('message.nicknameUpdated'));
			loadUsers();
		} catch (error) {
			message.error(t('message.nicknameUpdateFailed') + ': ' + (error.response?.data?.message || error.message));
			throw error;
		}
	};

	const handleDelete = async (userId, username) => {
		try {
			await deleteUser(userId);
			message.success(t('message.userDeleted', { username }));
			loadUsers();
		} catch (error) {
			message.error(t('message.deleteFailed') + ': ' + (error.response?.data?.message || error.message));
		}
	};

	const handleResetPassword = async (userId, username) => {
		try {
			const response = await resetPassword(userId);
			// apiClient已经提取了data,直接使用response.default_password
			const defaultPassword = response.default_password;
			message.success(t('message.passwordReset', { username, password: defaultPassword }));
		} catch (error) {
			message.error(t('message.resetPasswordFailed') + ': ' + (error.response?.data?.message || error.message));
		}
	};

	// 人员分配相关方法
	useEffect(() => {
		const fetchInitialData = async () => {
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
				const userList = usersData?.users || usersData;
				setAllUsers(Array.isArray(userList) ? userList : []);

				// 获取当前用户信息
				const currentUser = await getCurrentUser();
				const userId = currentUser?.ID || currentUser?.id || currentUser?.user_id;
				
				if (!userId) {
					console.error('[UserManage] 错误: 无法获取当前用户ID!', currentUser);
				} else {
					setCurrentUserId(userId);
				}
			} catch (error) {
				console.error('[UserManage] 加载数据失败:', error);
				message.error(t('assign.loadProjectsError'));
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
				
				// 分离管理员和成员
				const managers = (data.managers || []).map((m) => m.user_id || m.ID || m.id);
				const members = (data.members || []).map((m) => m.user_id || m.ID || m.id);
				
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
	const handleSaveAssign = async () => {
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
	
	// 数据转换：后端返回的字段名是大写开头(ID, Username, Nickname, Role)
	// 需要转换为小写下划线格式(user_id, username, nickname, role)
	const normalizedUsers = safeAllUsers.map((user) => {
		return {
			user_id: user.ID || user.id || user.user_id,
			username: user.Username || user.username,
			nickname: user.Nickname || user.nickname,
			role: user.Role || user.role,
		};
	});

	// 判断当前用户的系统角色是否为项目管理员
	const currentUserInNormalized = normalizedUsers.find(u => u.user_id === currentUserId);
	const currentUserRoleInAssign = currentUserInNormalized?.role;
	const isCurrentUserManager = currentUserRoleInAssign === 'project_manager';
	
	// 锁定用户ID数组（如果当前用户系统角色是管理员则锁定，防止自己把自己移出项目）
	const lockedKeys = isCurrentUserManager ? [currentUserId] : [];
	
	// 项目管理员穿梭框：显示系统角色为 project_manager 的用户
	const projectManagersForAssign = normalizedUsers.filter((user) => user.role === 'project_manager');
	
	// 项目成员穿梭框：显示系统角色为 project_member 的用户
	const projectMembersForAssign = normalizedUsers.filter((user) => user.role === 'project_member');
	
	// 检查是否有变更
	const hasChanges = JSON.stringify(managerTargetKeys.sort()) !== JSON.stringify(originalManagerKeys.sort()) ||
	                    JSON.stringify(memberTargetKeys.sort()) !== JSON.stringify(originalMemberKeys.sort());

	const columns = [
		{ title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
		{ title: t('user.username'), dataIndex: 'username', key: 'username' },
		{
			title: t('user.nickname'),
			dataIndex: 'nickname',
			key: 'nickname',
			render: (text, record) => (
				<EditableCell
					value={text}
					onSave={(newValue) => handleNicknameChange(record.id, newValue)}
				/>
			),
		},
		{
			title: t('user.role'),
			dataIndex: 'role',
			key: 'role',
			render: (role) => role === 'project_manager' ? t('user.projectManager') : t('user.projectMember'),
		},
		{
			title: t('user.actions'),
			key: 'action',
			render: (_, record) => (
				<span>
					<Button
						type="link"
						icon={<KeyOutlined />}
						onClick={() => handleResetPassword(record.id, record.username)}
					>
						{t('user.resetPassword')}
					</Button>
					<Popconfirm
						title={t('user.deleteConfirm', { username: record.username })}
						onConfirm={() => handleDelete(record.id, record.username)}
						okText={t('common.ok')}
						cancelText={t('common.cancel')}
					>
						<Button type="link" danger icon={<DeleteOutlined />}>
							{t('common.delete')}
						</Button>
					</Popconfirm>
				</span>
			),
		},
	];

	return (
		<div style={{ padding: '8px', background: '#fff' }} className="user-manage-container">
			<Tabs
				defaultActiveKey="members"
				items={[
					{
						key: 'members',
						label: t('user.projectMember') || 'Project Members',
						children: (
							<div>
								<Row gutter={16}>
									{/* 仅系统管理员显示项目管理员列表 */}
									{currentUserRole === 'system_admin' && (
										<Col span={12}>
											<Card
												title={t('user.projectManagerList')}
												extra={
													<Button
														type="primary"
														icon={<PlusOutlined />}
														onClick={() => handleCreate('project_manager')}
													>
														{t('user.createProjectManager')}
													</Button>
												}
												className="user-manage-card"
											>
												<Table
													columns={columns}
													dataSource={projectManagers}
													rowKey="id"
													loading={loading}
													pagination={{ pageSize: 10 }}
												/>
											</Card>
										</Col>
									)}
									<Col span={currentUserRole === 'system_admin' ? 12 : 24}>
										<Card
											title={t('user.projectMember')}
											extra={
												<Button
													type="primary"
													icon={<PlusOutlined />}
													onClick={() => handleCreate('project_member')}
												>
													{t('user.createProjectMember')}
												</Button>
											}
											className="user-manage-card"
										>
											<Table
												columns={columns}
												dataSource={projectMembers}
												rowKey="id"
												loading={loading}
												pagination={{ pageSize: 10 }}
											/>
										</Card>
									</Col>
								</Row>

								<CreateUserModal
									visible={modalVisible}
									role={createRole}
									onCancel={() => setModalVisible(false)}
									onSuccess={() => {
										setModalVisible(false);
										loadUsers();
									}}
								/>
							</div>
						),
					},
					{
						key: 'assign',
						label: t('menu.assign') || 'Personnel Assignment',
						children: currentUserRole === 'project_manager' ? (
							<div className="user-assign-section">
								<div className="user-assign-controls">
									<label>
										{t('assign.selectProject')}:
									</label>
									<Select
										showSearch
										placeholder={t('assign.selectProjectPlaceholder')}
										style={{ width: 300 }}
										value={selectedProjectId}
										onChange={handleProjectChange}
										filterOption={(input, option) =>
											option.children.toLowerCase().includes(input.toLowerCase())
										}
									>
										{Array.isArray(projects) && projects.map((project) => (
											<Select.Option key={project.id} value={project.id}>
												{project.name}
											</Select.Option>
										))}
									</Select>
									{selectedProjectId && hasChanges && (
										<Button
											type="primary"
											onClick={handleSaveAssign}
											loading={saving}
										>
											{t('assign.saveMembers')}
										</Button>
									)}
								</div>
								
								{selectedProjectId && isCurrentUserManager && (
									<div style={{ marginBottom: 12, color: '#ff9800', fontWeight: 'bold', fontSize: '13px' }}>
										⚠️ {t('assign.lockedUser')}
									</div>
								)}

								{selectedProjectId && (
									<Spin spinning={loadingMembers}>
										<Row gutter={8}>
											<Col span={12}>
												<h3 style={{ fontSize: '14px', fontWeight: '600', marginBottom: '8px' }}>
													{t('assign.managerTransferTitle')}
												</h3>
												<MemberTransfer
													dataSource={projectManagersForAssign}
													targetKeys={managerTargetKeys}
													lockedKeys={lockedKeys}
													onChange={handleManagerChange}
													title={[t('assign.availableUsers'), t('assign.assignedUsers')]}
												/>
											</Col>
											<Col span={12}>
												<h3 style={{ fontSize: '14px', fontWeight: '600', marginBottom: '8px' }}>
													{t('assign.memberTransferTitle')}
												</h3>
												<MemberTransfer
													dataSource={projectMembersForAssign}
													targetKeys={memberTargetKeys}
													lockedKeys={[]}
													onChange={handleMemberChange}
													title={[t('assign.availableUsers'), t('assign.assignedUsers')]}
												/>
											</Col>
										</Row>
									</Spin>
								)}
							</div>
						) : (
							<div className="user-assign-section">
								<p>{t('user.accessDenied') || 'Only project managers can access personnel assignment'}</p>
							</div>
						),
					},
				]}
			/>
		</div>
	);
};

export default UserManage;
