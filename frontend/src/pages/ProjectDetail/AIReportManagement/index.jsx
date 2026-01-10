import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { message, Spin } from 'antd';
import { 
  fetchAIReports, 
  createAIReport, 
  updateAIReport, 
  deleteAIReport,
  fetchAIReportDetail 
} from '../../../api/aiReport';
import AIReportLeftPanel from './AIReportLeftPanel';
import AIReportEditor from './AIReportEditor';
import './index.css';

const AIReportManagementTab = ({ projectId, projectName }) => {
  const { t } = useTranslation();
  const [reports, setReports] = useState([]);
  const [selectedReportId, setSelectedReportId] = useState(null);
  const [selectedReport, setSelectedReport] = useState(null);
  const [loading, setLoading] = useState(true);
  const [listLoading, setListLoading] = useState(false);

  // 加载报告列表
  const loadReports = async () => {
    setListLoading(true);
    try {
      const data = await fetchAIReports(projectId);
      setReports(Array.isArray(data) ? data : []);
    } catch (error) {
      message.error(t('aiReport.loadFailed'));
      console.error('Failed to load reports:', error);
    } finally {
      setListLoading(false);
    }
  };

  // 初始化加载
  useEffect(() => {
    setLoading(true);
    loadReports().finally(() => setLoading(false));
  }, [projectId]);

  // 选择报告
  const handleSelectReport = async (reportId) => {
    setSelectedReportId(reportId);
    try {
      const reportData = await fetchAIReportDetail(projectId, reportId);
      setSelectedReport(reportData);
    } catch (error) {
      message.error(t('aiReport.loadDetailFailed'));
    }
  };

  // 创建报告
  const handleCreateReport = async (name) => {
    try {
      const newReport = await createAIReport(projectId, name);
      await loadReports();
      // 自动选择新创建的报告
      setSelectedReportId(newReport.id);
      setSelectedReport(newReport);
    } catch (error) {
      if (error.response?.status === 409) {
        throw new Error(t('aiReport.reportNameDuplicate'));
      }
      throw error;
    }
  };

  // 编辑报告名称
  const handleEditReport = async (reportId, newName) => {
    try {
      await updateAIReport(projectId, reportId, { name: newName });
      await loadReports();
      if (selectedReportId === reportId) {
        setSelectedReport((prev) => ({ ...prev, name: newName }));
      }
      message.success(t('aiReport.updateSuccess'));
    } catch (error) {
      if (error.response?.status === 409) {
        message.error(t('aiReport.reportNameDuplicate'));
      } else {
        message.error(error.message || t('aiReport.updateFailed'));
      }
      throw error;
    }
  };

  // 删除报告
  const handleDeleteReport = async (reportId) => {
    try {
      await deleteAIReport(projectId, reportId);
      await loadReports();
      
      // 如果删除的是当前选中的报告，清空选择
      if (selectedReportId === reportId) {
        setSelectedReportId(null);
        setSelectedReport(null);
      }
      message.success(t('aiReport.deleteSuccess'));
    } catch (error) {
      message.error(error.message || t('aiReport.deleteFailed'));
      throw error;
    }
  };

  // 保存报告内容
  const handleSaveReport = async (reportId, content) => {
    try {
      await updateAIReport(projectId, reportId, { content });
      setSelectedReport((prev) => ({ ...prev, content }));
      return true;
    } catch (error) {
      throw error;
    }
  };

  // 内容变化
  const handleContentChange = (content) => {
    if (selectedReport) {
      setSelectedReport((prev) => ({ ...prev, content }));
    }
  };

  if (loading) {
    return (
      <div className="ai-report-management-tab">
        <Spin />
      </div>
    );
  }

  return (
    <div className="ai-report-management-tab">
      <div className="left-panel">
        <AIReportLeftPanel
          reports={reports}
          selectedReportId={selectedReportId}
          loading={listLoading}
          onSelect={handleSelectReport}
          onCreate={handleCreateReport}
          onDelete={handleDeleteReport}
        />
      </div>
      <div className="right-panel">
        <AIReportEditor
          report={selectedReport}
          projectName={projectName}
          onSave={handleSaveReport}
          onContentChange={handleContentChange}
          onNameChange={handleEditReport}
          loading={loading}
        />
      </div>
    </div>
  );
};

export default AIReportManagementTab;
