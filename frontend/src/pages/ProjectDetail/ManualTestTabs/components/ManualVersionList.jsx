import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Input, message, Popconfirm } from 'antd';
import { DownloadOutlined, DeleteOutlined, EditOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import {
    getVersionList,
    downloadVersion,
    deleteVersion,
    updateVersionRemark
} from '../../../../api/manualCase';

const { TextArea } = Input;

/**
 * 手工用例版本列表组件
 * @param {string} projectId - 项目ID
 * @param {function} onVersionDeleted - 版本删除后的回调
 */
const ManualVersionList = ({ projectId, onVersionDeleted }) => {
    const { t } = useTranslation();
    const [versions, setVersions] = useState([]);
    const [loading, setLoading] = useState(false);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 10,
        total: 0,
    });
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [editingVersion, setEditingVersion] = useState(null);
    const [editingRemark, setEditingRemark] = useState('');

    useEffect(() => {
        loadVersions();
    }, [projectId]);

    const loadVersions = async (page = 1, size = 10) => {
        setLoading(true);
        try {
            // 获取手工用例版本列表（overall类型）
            const result = await getVersionList(projectId, 'overall');
            // API返回的是数组格式，需要适配分页
            const allVersions = Array.isArray(result) ? result : (result.versions || []);
            setVersions(allVersions);
            setPagination({
                current: page,
                pageSize: size,
                total: allVersions.length,
            });
        } catch (error) {
            message.error(t('manualTest.loadVersionsFailed'));
            console.error('Failed to load versions:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleDownload = async (versionId) => {
        try {
            await downloadVersion(projectId, versionId);
            message.success(t('manualTest.downloadSuccess'));
        } catch (error) {
            message.error(t('manualTest.downloadFailed'));
            console.error('Failed to download version:', error);
        }
    };

    const handleDelete = async (versionId) => {
        try {
            await deleteVersion(projectId, versionId);
            message.success(t('message.deleteSuccess'));
            loadVersions(pagination.current, pagination.pageSize);
            if (onVersionDeleted) {
                onVersionDeleted();
            }
        } catch (error) {
            message.error(t('message.deleteFailed'));
            console.error('Failed to delete version:', error);
        }
    };

    const handleEditRemark = (record) => {
        setEditingVersion(record);
        setEditingRemark(record.remark || '');
        setEditModalVisible(true);
    };

    const handleSaveRemark = async () => {
        try {
            const versionId = editingVersion.id;
            await updateVersionRemark(projectId, versionId, editingRemark);
            message.success(t('manualTest.remarkUpdateSuccess'));
            setEditModalVisible(false);
            loadVersions(pagination.current, pagination.pageSize);
        } catch (error) {
            message.error(t('manualTest.remarkUpdateFailed'));
            console.error('Failed to update version remark:', error);
        }
    };

    const handleTableChange = (newPagination) => {
        setPagination({
            ...pagination,
            current: newPagination.current,
            pageSize: newPagination.pageSize,
        });
    };

    const columns = [
        {
            title: t('manualTest.versionId'),
            dataIndex: 'id',
            key: 'id',
            width: 80,
            render: (_, __, index) => (pagination.current - 1) * pagination.pageSize + index + 1,
        },
        {
            title: t('manualTest.versionFilename'),
            dataIndex: 'filename',
            key: 'filename',
            ellipsis: true,
            render: (text) => text || '-',
        },
        {
            title: t('manualTest.remark'),
            dataIndex: 'remark',
            key: 'remark',
            width: 200,
            ellipsis: true,
            render: (text) => text || '-',
        },
        {
            title: t('manualTest.operation'),
            key: 'actions',
            width: 200,
            render: (_, record) => {
                const versionId = record.id;
                return (
                    <div style={{ display: 'flex', gap: '4px' }}>
                        <Button
                            type="link"
                            size="small"
                            icon={<EditOutlined />}
                            onClick={() => handleEditRemark(record)}
                        >
                            {t('manualTest.edit')}
                        </Button>
                        <Button
                            type="link"
                            size="small"
                            icon={<DownloadOutlined />}
                            onClick={() => handleDownload(versionId)}
                        >
                            {t('manualTest.download')}
                        </Button>
                        <Popconfirm
                            title={t('manualTest.confirmDelete')}
                            onConfirm={() => handleDelete(versionId)}
                            okText={t('common.confirm')}
                            cancelText={t('common.cancel')}
                        >
                            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                                {t('manualTest.delete')}
                            </Button>
                        </Popconfirm>
                    </div>
                );
            },
        },
    ];

    return (
        <div className="manual-version-list">
            <Table
                columns={columns}
                dataSource={versions}
                loading={loading}
                rowKey="id"
                pagination={{
                    ...pagination,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    pageSizeOptions: ['10', '20', '50', '100'],
                    showTotal: (total) => t('web_version.totalVersions', { total }),
                }}
                onChange={handleTableChange}
                size="small"
            />

            <Modal
                title={t('manualTest.editRemarkTitle')}
                open={editModalVisible}
                onOk={handleSaveRemark}
                onCancel={() => setEditModalVisible(false)}
                okText={t('common.save')}
                cancelText={t('common.cancel')}
            >
                <TextArea
                    value={editingRemark}
                    onChange={(e) => setEditingRemark(e.target.value)}
                    placeholder={t('manualTest.remarkPlaceholder')}
                    maxLength={200}
                    rows={4}
                    showCount
                />
            </Modal>
        </div>
    );
};

export default ManualVersionList;
