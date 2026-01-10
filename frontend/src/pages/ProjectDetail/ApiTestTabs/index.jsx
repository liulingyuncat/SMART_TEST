import React from 'react';
import ApiCaseManagementTab from './containers/ApiCaseManagementTab';

/**
 * 接口测试Tab总容器组件
 * 采用左右分栏布局，左侧用例集管理，右侧用例表格
 */
const ApiTestTabs = ({ projectId }) => {
  return <ApiCaseManagementTab projectId={projectId} />;
};

export default ApiTestTabs;
