import React, { useState, useEffect } from 'react';
import { useSearchParams, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Tabs } from 'antd';
import ManualCaseManagementTab from './containers/ManualCaseManagementTab';
import CaseReviewTab from './components/CaseReviewTab';
import VersionManagementTab from './components/VersionManagementTab';
import './index.css';

// Tab配置数组 - 精简为3个核心Tab
const MANUAL_TEST_TABS = [
  { 
    key: 'manual-cases', 
    labelKey: 'manualTest.manualCaseManagement', 
    component: ManualCaseManagementTab 
  },
  { 
    key: 'case-review', 
    labelKey: 'manualTest.reviewCases', 
    component: CaseReviewTab, 
    needsProjectId: true, 
    props: { caseType: 'overall' } 
  },
  { 
    key: 'version-management', 
    labelKey: 'manualTest.versionManagement',
    component: VersionManagementTab, 
    needsProjectId: true,
    props: {
      leftDocType: 'manual', // 仅显示manual类型版本
      leftTitleKey: 'manualTest.versionManagement'
    }
  },
];

const ManualTestTabs = ({ projectId }) => {
  const { t } = useTranslation();
  const [searchParams, setSearchParams] = useSearchParams();
  
  // 从sessionStorage恢复上次的Tab，如果没有则默认为'manual-cases'
  const getStorageKey = () => `manual_test_active_tab_${projectId}`;
  const getSavedTab = () => {
    if (!projectId) return 'manual-cases';
    try {
      const saved = sessionStorage.getItem(`manual_test_active_tab_${projectId}`);
      return saved || 'manual-cases';
    } catch {
      return 'manual-cases';
    }
  };
  
  const [activeTab, setActiveTab] = useState(() => {
    // 使用函数形式的初始化，延迟执行
    return getSavedTab();
  });

  // 从URL读取tab参数初始化activeTab
  useEffect(() => {
    const tabParam = searchParams.get('tab');
    if (tabParam && MANUAL_TEST_TABS.find(tab => tab.key === tabParam)) {
      setActiveTab(tabParam);
      // 保存到sessionStorage
      try {
        sessionStorage.setItem(getStorageKey(), tabParam);
      } catch (e) {
        console.warn('Failed to save tab to sessionStorage:', e);
      }
    } else if (!tabParam) {
      // 如果URL没有tab参数，使用保存的tab并更新URL
      const savedTab = getSavedTab();
      setSearchParams({ tab: savedTab });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams]);

  // Tab切换处理
  const handleTabChange = (key) => {
    setActiveTab(key);
    setSearchParams({ tab: key });
    // 保存到sessionStorage
    try {
      sessionStorage.setItem(getStorageKey(), key);
    } catch (e) {
      console.warn('Failed to save tab to sessionStorage:', e);
    }
  };

  // 生成Tab项
  const tabItems = MANUAL_TEST_TABS.map(tab => {
    const componentProps = { ...(tab.props || {}) };
    
    // 所有组件都传递 projectId
    componentProps.projectId = projectId;
    
    // 为版本管理Tab添加key，使其在每次激活时刷新
    if (tab.key === 'version-management' && activeTab === 'version-management') {
      componentProps.key = `version-${Date.now()}`;
    }

    return {
      key: tab.key,
      label: tab.labelKey ? t(tab.labelKey) : tab.label,
      children: React.createElement(tab.component, componentProps),
    };
  });

  return (
    <div className="manual-test-tabs">
      <Tabs
        activeKey={activeTab}
        onChange={handleTabChange}
        items={tabItems}
        destroyInactiveTabPane={false}
      />
    </div>
  );
};

export default ManualTestTabs;
