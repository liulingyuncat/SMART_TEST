import React from 'react';
import ManualCaseManagementTab from './containers/ManualCaseManagementTab';
import './index.css';

/**
 * 手工测试用例库入口组件
 * 直接显示手工用例管理内容，不再使用Tab
 */
const ManualTestTabs = ({ projectId }) => {
  return (
    <div className="manual-test-tabs">
      <ManualCaseManagementTab projectId={projectId} />
    </div>
  );
};

export default ManualTestTabs;
