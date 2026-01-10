import React, { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { Row, Col, Button } from 'antd';
import { LeftOutlined, RightOutlined } from '@ant-design/icons';
import RequirementItemPanel from './RequirementItemPanel';
import ViewpointItemPanel from './ViewpointItemPanel';
import './index.css';

const RequirementManagement = ({ projectId, projectName }) => {
  const { t } = useTranslation();
  
  // 左右栏宽度比例，从localStorage读取，默认各占50%
  const [splitRatio, setSplitRatio] = useState(() => {
    const saved = localStorage.getItem('requirementManagementSplitRatio');
    return saved ? parseFloat(saved) : 50;
  });
  
  const [isDragging, setIsDragging] = useState(false);
  const [requirementPanelCollapsed, setRequirementPanelCollapsed] = useState(false);
  const [viewpointPanelCollapsed, setViewpointPanelCollapsed] = useState(false);
  const containerRef = useRef(null);

  // 保存宽度比例到localStorage
  useEffect(() => {
    localStorage.setItem('requirementManagementSplitRatio', splitRatio.toString());
  }, [splitRatio]);

  // 处理鼠标拖拽
  const handleMouseDown = (e) => {
    e.preventDefault();
    setIsDragging(true);
  };

  useEffect(() => {
    const handleMouseMove = (e) => {
      if (!isDragging || !containerRef.current) return;
      
      const container = containerRef.current;
      const containerRect = container.getBoundingClientRect();
      const newRatio = ((e.clientX - containerRect.left) / containerRect.width) * 100;
      
      // 限制在20%到80%之间
      if (newRatio >= 20 && newRatio <= 80) {
        setSplitRatio(newRatio);
      }
    };

    const handleMouseUp = () => {
      setIsDragging(false);
    };

    if (isDragging) {
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
    }

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isDragging]);

  return (
    <div className="requirement-management-container" ref={containerRef}>
      <Row style={{ height: '100%', position: 'relative' }}>
        {/* 左栏：需求管理 */}
        {!requirementPanelCollapsed && (
          <Col 
            style={{ 
              width: viewpointPanelCollapsed ? '100%' : `${splitRatio}%`, 
              height: '100%',
              display: 'flex',
              flexDirection: 'column',
              borderRight: '1px solid #e8e8e8'
            }}
          >
            <div className="panel-title-bar">
              <span className="panel-title">需求管理</span>
              <Button 
                type="text" 
                size="small" 
                icon={<LeftOutlined />}
                onClick={() => setRequirementPanelCollapsed(true)}
              />
            </div>
            <div style={{ flex: 1, overflow: 'hidden', padding: '0 16px 16px 16px' }}>
              <RequirementItemPanel projectId={projectId} projectName={projectName} />
            </div>
          </Col>
        )}
        {requirementPanelCollapsed && (
          <div 
            className="collapsed-panel-trigger"
            onClick={() => setRequirementPanelCollapsed(false)}
          >
            <RightOutlined />
          </div>
        )}
        
        {/* 拖拽分隔条 */}
        {!requirementPanelCollapsed && !viewpointPanelCollapsed && (
          <div
            className="drag-handle"
            onMouseDown={handleMouseDown}
            style={{
              width: '5px',
              cursor: 'col-resize',
              backgroundColor: isDragging ? '#1890ff' : '#f0f0f0',
              position: 'absolute',
              left: `${splitRatio}%`,
              top: 0,
              bottom: 0,
              zIndex: 10,
              transition: isDragging ? 'none' : 'background-color 0.2s'
            }}
          />
        )}
        
        {/* 右栏：AI观点管理 */}
        {!viewpointPanelCollapsed && (
          <Col 
            style={{ 
              width: requirementPanelCollapsed ? '100%' : `${100 - splitRatio}%`, 
              height: '100%',
              display: 'flex',
              flexDirection: 'column'
            }}
          >
            <div className="panel-title-bar">
              <span className="panel-title">AI观点管理</span>
              <Button 
                type="text" 
                size="small" 
                icon={<RightOutlined />}
                onClick={() => setViewpointPanelCollapsed(true)}
              />
            </div>
            <div style={{ flex: 1, overflow: 'hidden', padding: '0 16px 16px 16px' }}>
              <ViewpointItemPanel projectId={projectId} projectName={projectName} />
            </div>
          </Col>
        )}
        {viewpointPanelCollapsed && (
          <div 
            className="collapsed-panel-trigger right"
            onClick={() => setViewpointPanelCollapsed(false)}
          >
            <LeftOutlined />
          </div>
        )}
      </Row>
    </div>
  );
};

export default RequirementManagement;
