import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { message } from 'antd';
import LanguageFilter from '../../ManualTestTabs/components/LanguageFilter';
import EditableTable from '../../ManualTestTabs/components/EditableTable';
import ReorderModal from '../../ManualTestTabs/components/ReorderModal';
import './Role1Tab.css';

/**
 * ROLE1 Tab容器组件
 * 管理role1类型的自动化测试用例
 */
const Role1Tab = ({ projectId }) => {
  const [language, setLanguage] = useState('中文');
  const [reorderModalVisible, setReorderModalVisible] = useState(false);
  const [casesForReorder, setCasesForReorder] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [refreshKey, setRefreshKey] = useState(0);

  // 语言筛选变更
  const handleLanguageChange = (newLanguage) => {
    setLanguage(newLanguage);
  };

  // 打开重排对话框
  const handleReorderClick = (currentCases, pageNumber) => {
    setCasesForReorder(currentCases || []);
    setCurrentPage(pageNumber || 1);
    setReorderModalVisible(true);
  };

  // 重排成功回调
  const handleReorderSuccess = () => {
    setReorderModalVisible(false);
    setRefreshKey(prev => prev + 1);
  };

  return (
    <div className="role1-tab">
      <LanguageFilter 
        value={language}
        onChange={handleLanguageChange}
      />

      <EditableTable
        key={refreshKey}
        projectId={projectId}
        caseType="role1"
        language={language}
        onReorderClick={handleReorderClick}
      />

      <ReorderModal
        visible={reorderModalVisible}
        caseType="role1"
        projectId={projectId}
        language={language}
        cases={casesForReorder}
        currentPage={currentPage}
        onOk={handleReorderSuccess}
        onCancel={() => setReorderModalVisible(false)}
      />
    </div>
  );
};

export default Role1Tab;
