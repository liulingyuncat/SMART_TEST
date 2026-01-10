import React, { useState } from 'react';
import { Modal, Form, Input, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { createUser, checkUnique } from '../api/user';

const CreateUserModal = ({ visible, role, onCancel, onSuccess }) => {
	const { t } = useTranslation();
	const [form] = Form.useForm();
	const [loading, setLoading] = useState(false);

	const handleSubmit = async () => {
		try {
			const values = await form.validateFields();
			setLoading(true);

			await createUser({
				username: values.username,
				nickname: values.nickname,
				role: role,
			});

			const defaultPassword = role === 'project_manager' ? 'admin!123' : 'user!123';
			message.success(t('message.userCreated', { password: defaultPassword }));
			form.resetFields();
			onSuccess();
		} catch (error) {
			if (error.response) {
				message.error(t('message.createFailed') + ': ' + error.response.data.message);
			} else if (!error.errorFields) {
				message.error(t('message.createFailed') + ': ' + error.message);
			}
		} finally {
			setLoading(false);
		}
	};

	const validateUsername = async (_, value) => {
		if (!value) return Promise.resolve();
		try {
			const response = await checkUnique({ username: value });
			if (response.exists) {
				return Promise.reject(new Error(t('validation.usernameExists')));
			}
			return Promise.resolve();
		} catch (error) {
			return Promise.reject(new Error(t('validation.checkFailed')));
		}
	};

	const validateNickname = async (_, value) => {
		if (!value) return Promise.resolve();
		try {
			const response = await checkUnique({ nickname: value });
			if (response.exists) {
				return Promise.reject(new Error(t('validation.nicknameExists')));
			}
			return Promise.resolve();
		} catch (error) {
			return Promise.reject(new Error(t('validation.checkFailed')));
		}
	};

	return (
		<Modal
			title={role === 'project_manager' ? t('user.createProjectManager') : t('user.createProjectMember')}
			open={visible}
			onCancel={onCancel}
			onOk={handleSubmit}
			confirmLoading={loading}
		>
			<Form form={form} layout="vertical">
				<Form.Item
					name="username"
					label={t('user.username')}
					rules={[
						{ required: true, message: t('validation.usernameRequired') },
						{ min: 3, max: 50, message: t('validation.usernameLength') },
						{ pattern: /^[a-zA-Z0-9_]+$/, message: t('validation.usernameFormat') },
						{ validator: validateUsername },
					]}
				>
					<Input placeholder={t('login.usernamePlaceholder')} />
				</Form.Item>

				<Form.Item
					name="nickname"
					label={t('user.nickname')}
					rules={[
						{ required: true, message: t('validation.nicknameRequired') },
						{ min: 2, max: 50, message: t('validation.nicknameLength') },
						{ validator: validateNickname },
					]}
				>
					<Input placeholder={t('validation.nicknamePlaceholder')} />
				</Form.Item>
			</Form>
		</Modal>
	);
};

export default CreateUserModal;
