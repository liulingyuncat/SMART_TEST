import React from 'react';
import ViewpointItemPanel from './ViewpointItemPanel';
import './AIViewpointTab.css';

/**
 * AI观点Tab - 独立显示AI观点管理内容
 */
const AIViewpointTab = ({ projectId, projectName }) => {
  return (
    <div className="ai-viewpoint-tab-container">
      <ViewpointItemPanel projectId={projectId} projectName={projectName} />
    </div>
  );
};

export default AIViewpointTab;
