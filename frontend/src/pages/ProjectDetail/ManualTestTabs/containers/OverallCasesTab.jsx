import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { message } from 'antd';
import LanguageFilter from '../components/LanguageFilter';
import EditableTable from '../components/EditableTable';
import ReorderModal from '../components/ReorderModal';
import './OverallCasesTab.css';

/**
 * 整体用例Tab容器组件
 */
const OverallCasesTab = ({ projectId }) => {
  const [language, setLanguage] = useState('中文');
  const [reorderModalVisible, setReorderModalVisible] = useState(false);
  const [casesForReorder, setCasesForReorder] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [refreshKey, setRefreshKey] = useState(0);

  // 语言筛选变更
  const handleLanguageChange = (newLanguage) => {
    setLanguage(newLanguage);
  };

  // 打开重排对话框（接收当前显示的cases数组和页码）
  const handleReorderClick = (currentCases, pageNumber) => {
    setCasesForReorder(currentCases || []);
    setCurrentPage(pageNumber || 1);
    setReorderModalVisible(true);
  };

  // 重排成功回调
  const handleReorderSuccess = () => {
    setReorderModalVisible(false);
    setRefreshKey(prev => prev + 1); // 刷新表格
  };

  return (
    <div className="overall-cases-tab">
      <LanguageFilter 
        value={language}
        onChange={handleLanguageChange}
      />

      <EditableTable
        key={refreshKey}
        projectId={projectId}
        caseType="overall"
        language={language}
        onReorderClick={handleReorderClick}
      />

      <ReorderModal
        visible={reorderModalVisible}
        caseType="overall"
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

export default OverallCasesTab;
