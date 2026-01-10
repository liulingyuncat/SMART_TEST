import React from 'react';
import WebCaseManagementTab from './WebCaseManagementTab';
import './AutoTestTab.css';

/**
 * AIWeb用例库入口组件
 * 直接渲染左右两栏布局的Web用例管理组件
 */
const AutoTestTab = ({ projectId }) => {
  return (
    <div className="auto-test-tab">
      <WebCaseManagementTab projectId={projectId} />
    </div>
  );
};

export default AutoTestTab;
