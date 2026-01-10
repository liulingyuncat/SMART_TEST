import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import EditableTable from '../components/EditableTable';
import ReorderModal from '../components/ReorderModal';

/**
 * AI用例Tab容器组件
 * 特点: 
 * - 不显示元数据编辑器(测试版本、测试环境等)
 * - 不显示语言筛选器(AI用例固定为中文)
 * - 只显示可编辑表格
 */
const AICasesTab = ({ projectId }) => {
  const [reorderModalVisible, setReorderModalVisible] = useState(false);
  const [casesForReorder, setCasesForReorder] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [refreshKey, setRefreshKey] = useState(0);

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
    <div className="ai-cases-tab">
      <EditableTable
        key={refreshKey}
        caseType="ai"
        language="中文"
        onReorderClick={handleReorderClick}
        projectId={projectId}
      />

      <ReorderModal
        visible={reorderModalVisible}
        caseType="ai"
        projectId={projectId}
        cases={casesForReorder}
        currentPage={currentPage}
        onOk={handleReorderSuccess}
        onCancel={() => setReorderModalVisible(false)}
      />
    </div>
  );
};

export default AICasesTab;
