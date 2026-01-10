import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { message } from 'antd';
import LanguageFilter from '../components/LanguageFilter';
import EditableTable from '../components/EditableTable';
import ReorderModal from '../components/ReorderModal';

/**
 * 变更用例Tab容器组件
 */
const ChangeCasesTab = ({ projectId }) => {
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
  const handleReorderClick = (cases, pageNumber) => {
    setCasesForReorder(cases);
    setCurrentPage(pageNumber || 1);
    setReorderModalVisible(true);
  };

  // 重排成功回调
  const handleReorderSuccess = () => {
    setReorderModalVisible(false);
    setRefreshKey(prev => prev + 1); // 刷新表格
  };

  return (
    <div className="change-cases-tab">
      <LanguageFilter 
        value={language}
        onChange={handleLanguageChange}
      />

      <EditableTable
        key={refreshKey}
        projectId={projectId}
        caseType="change"
        language={language}
        onReorderClick={handleReorderClick}
      />

      <ReorderModal
        visible={reorderModalVisible}
        caseType="change"
        projectId={projectId}
        cases={casesForReorder}
        currentPage={currentPage}
        onOk={handleReorderSuccess}
        onCancel={() => setReorderModalVisible(false)}
      />
    </div>
  );
};

export default ChangeCasesTab;

