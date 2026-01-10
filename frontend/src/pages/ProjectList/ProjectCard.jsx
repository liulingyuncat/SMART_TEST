import React, { useState } from 'react';
import { Card, Tag, Typography, Dropdown, Button, App, message } from 'antd';
import { CalendarOutlined, EditOutlined, DeleteOutlined, MoreOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import PropTypes from 'prop-types';
import dayjs from 'dayjs';
import { deleteProject } from '../../api/project';
import EditProjectModal from './EditProjectModal';
import './ProjectCard.css';

const { Title, Paragraph } = Typography;

const ProjectCard = ({ project, onUpdate, onDelete }) => {
  const { t } = useTranslation();
  const { modal } = App.useApp();
  const navigate = useNavigate();
  const user = useSelector((state) => state.auth.user);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const isProjectManager = user?.role === 'project_manager';

  console.log('ProjectCard rendered, project:', project.id, 'isProjectManager:', isProjectManager);

  // 卡片点击跳转到项目详情
  const handleCardClick = () => {
    navigate(`/projects/${project.id}`);
  };

  // 打开编辑对话框
  const handleEdit = (e) => {
    e.stopPropagation();
    console.log('handleEdit triggered');
    setEditModalVisible(true);
  };

  // 删除项目
  const handleDelete = (e) => {
    e.stopPropagation();
    console.log('handleDelete triggered for project:', project.id);
    console.log('modal.confirm is:', modal.confirm);
    
    modal.confirm({
      title: t('project.deleteProject'),
      content: t('project.confirmDelete', { name: project.name }),
      okText: t('common.confirm'),
      okType: 'danger',
      cancelText: t('common.cancel'),
      onOk: async () => {
        try {
          console.log('Deleting project:', project.id);
          const response = await deleteProject(project.id);
          console.log('Delete response:', response);
          message.success(t('project.deleteSuccess'));
          if (onDelete) {
            console.log('Calling onDelete callback with id:', project.id);
            onDelete(project.id);
          } else {
            console.warn('onDelete callback is not defined');
          }
        } catch (error) {
          console.error('Delete project error:', error);
          console.error('Error response:', error.response);
          message.error(error.response?.data?.message || t('project.deleteFailed'));
        }
      },
    });
  };

  // 菜单点击处理
  const handleMenuClick = ({ key }) => {
    console.log('Menu item clicked:', key);
    if (key === 'edit') {
      handleEdit();
    } else if (key === 'delete') {
      handleDelete();
    }
  };

  // 操作菜单项
  const menuItems = [
    {
      key: 'delete',
      icon: <DeleteOutlined />,
      label: t('project.delete'),
      danger: true,
    },
  ];

  return (
    <Card
      hoverable
      className="project-card"
      onClick={handleCardClick}
      title={
        <div className="project-card-header">
          <Title level={4} className="project-name" ellipsis={{ tooltip: project.name }}>
            {project.name}
          </Title>
          <div className="project-actions">
            <Tag color="blue">{t('project.active')}</Tag>
            {isProjectManager && (
              <Dropdown 
                menu={{ items: menuItems, onClick: handleMenuClick }} 
                trigger={['click']} 
                placement="bottomRight"
                onOpenChange={(open) => console.log('Dropdown open state:', open)}
              >
                <Button 
                  type="text" 
                  icon={<MoreOutlined />} 
                  onClick={() => console.log('Dropdown button clicked')}
                />
              </Dropdown>
            )}
          </div>
        </div>
      }
    >
      <div className="project-card-content">
        <Paragraph
          className="project-description"
          ellipsis={{ rows: 3, tooltip: project.description }}
        >
          {project.description || t('project.noDescription')}
        </Paragraph>
        <div className="project-meta">
          <CalendarOutlined />
          <span className="project-date">
            {dayjs(project.created_at).format('YYYY-MM-DD')}
          </span>
        </div>
      </div>

      <EditProjectModal
        visible={editModalVisible}
        project={project}
        onCancel={() => setEditModalVisible(false)}
        onSuccess={(updatedProject) => {
          setEditModalVisible(false);
          onUpdate(updatedProject);
        }}
      />
    </Card>
  );
};

ProjectCard.propTypes = {
  project: PropTypes.shape({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    description: PropTypes.string,
    created_at: PropTypes.string.isRequired,
  }).isRequired,
  onUpdate: PropTypes.func,
  onDelete: PropTypes.func,
};

export default ProjectCard;
