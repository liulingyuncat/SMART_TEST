import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import LanguageFilter from '../../ManualTestTabs/components/LanguageFilter';
import EditableTable from '../../ManualTestTabs/components/EditableTable';
import ReorderModal from '../../ManualTestTabs/components/ReorderModal';

const Role3Tab = ({ projectId }) => {
  const [language, setLanguage] = useState('中文');
  const [reorderModalVisible, setReorderModalVisible] = useState(false);
  const [casesForReorder, setCasesForReorder] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleLanguageChange = (newLanguage) => {
    setLanguage(newLanguage);
  };

  const handleReorderClick = (currentCases, pageNumber) => {
    setCasesForReorder(currentCases || []);
    setCurrentPage(pageNumber || 1);
    setReorderModalVisible(true);
  };

  const handleReorderSuccess = () => {
    setReorderModalVisible(false);
    setRefreshKey(prev => prev + 1);
  };

  return (
    <div className="role3-tab">
      <LanguageFilter 
        value={language}
        onChange={handleLanguageChange}
      />

      <EditableTable
        key={refreshKey}
        projectId={projectId}
        caseType="role3"
        language={language}
        onReorderClick={handleReorderClick}
      />

      <ReorderModal
        visible={reorderModalVisible}
        caseType="role3"
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

export default Role3Tab;
