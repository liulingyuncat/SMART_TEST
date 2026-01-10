import React from 'react';
import RequirementItemPanel from './RequirementItemPanel';
import './AIRequirementTab.css';

/**
 * AI需求Tab - 独立显示需求管理内容
 */
const AIRequirementTab = ({ projectId, projectName }) => {
  return (
    <div className="ai-requirement-tab-container">
      <RequirementItemPanel projectId={projectId} projectName={projectName} />
    </div>
  );
};

export default AIRequirementTab;
