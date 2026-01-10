import React, { useState, useEffect } from 'react';
import { Modal, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { reorderCasesByDrag, getCasesList } from '../../../../api/manualCase';
import { reorderAutoCases, getAutoCasesList } from '../../../../api/autoCase';

/**
 * ID重排对话框组件
 * 【新需求】按照当前页的位置插入，保持页码对应的编号
 * - 第一页：当前页cases变成No.1-10，其他页顺延
 * - 第二页：前一页保持No.1-10，当前页cases变成No.11-20，其他页顺延
 * 
 * @param {Object} props
 * @param {boolean} props.visible - 对话框是否可见
 * @param {string} props.caseType - 用例类型
 * @param {number} props.projectId - 项目ID
 * @param {string} props.language - 当前语言
 * @param {Array} props.cases - 当前页面显示的cases数组（按显示顺序）
 * @param {number} props.currentPage - 当前页码
 * @param {Function} props.onOk - 确认回调
 * @param {Function} props.onCancel - 取消回调
 */
const ReorderModal = ({ visible, caseType, projectId, language, cases = [], currentPage = 1, onOk, onCancel }) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [totalCount, setTotalCount] = useState(0);
  const [allCases, setAllCases] = useState([]);

  // 当对话框打开时，获取所有用例
  useEffect(() => {
    if (visible && projectId) {
      const isRoleType = caseType && caseType.startsWith('role');
      const apiCall = isRoleType ? getAutoCasesList : getCasesList;
      apiCall(projectId, { caseType, language, page: 1, size: 10000 })
        .then(data => {
          console.log('[ReorderModal] 获取到的用例数据:', data);
          console.log('[ReorderModal] 用例列表:', data.cases);
          setTotalCount(data.total || 0);
          setAllCases(data.cases || []);
        })
        .catch(error => {
          console.error('获取用例数据失败:', error);
        });
    }
  }, [visible, projectId, caseType, language]);

  const handleReorder = async () => {
    try {
      setLoading(true);
      
      // 【修正方案】根据当前页码，将当前页的cases插入到正确位置
      const pageSize = 10; // 固定每页10条
      
      console.log('[ReorderModal] 传入的当前页cases:', cases);
      console.log('[ReorderModal] 从数据库获取的allCases:', allCases);
      
      // 1. 确定当前页的cases：优先使用传入的cases，如果为空则从allCases中提取
      let currentPageCases = cases && cases.length > 0 ? cases : [];
      
      // 如果传入的cases为空，从allCases中按页码提取当前页数据
      if (currentPageCases.length === 0 && allCases.length > 0) {
        const startIndex = (currentPage - 1) * pageSize;
        const endIndex = startIndex + pageSize;
        // 先按id排序
        const sortedAll = [...allCases].sort((a, b) => (a.display_id || a.id) - (b.display_id || b.id));
        currentPageCases = sortedAll.slice(startIndex, endIndex);
        console.log('[ReorderModal] 从allCases提取当前页数据:', currentPageCases);
      }
      
      // 2. 获取当前页的case_id列表
      const currentPageCaseIds = currentPageCases.map(c => c.case_id);
      console.log('[ReorderModal] 当前页case_ids:', currentPageCaseIds);
      
      // 3. 获取所有其他页的cases（排除当前页）
      const otherCases = allCases.filter(c => !currentPageCaseIds.includes(c.case_id));
      console.log('[ReorderModal] 其他页cases数量:', otherCases.length);
      
      // 4. 将其他页的cases按原ID排序
      const sortedOtherCases = otherCases.sort((a, b) => (a.display_id || a.id) - (b.display_id || b.id));
      
      // 5. 计算插入位置：(currentPage - 1) * pageSize
      const insertIndex = (currentPage - 1) * pageSize;
      
      // 6. 构建最终顺序：前面的页 + 当前页 + 后面的页
      const beforeCases = sortedOtherCases.slice(0, insertIndex);
      const afterCases = sortedOtherCases.slice(insertIndex);
      
      const finalOrder = [
        ...beforeCases.map(c => c.case_id),
        ...currentPageCaseIds,
        ...afterCases.map(c => c.case_id)
      ];
      
      console.log('[ReorderModal] finalOrder:', finalOrder);
      
      if (finalOrder.length === 0) {
        message.warning('没有可重排的用例');
        return;
      }
      
      // 6. 调用后端API，按finalOrder重新分配ID
      const isRoleType = caseType && caseType.startsWith('role');
      const reorderAPI = isRoleType ? reorderAutoCases : reorderCasesByDrag;
      await reorderAPI(projectId, caseType, finalOrder);
      
      const startNo = insertIndex + 1;
      const endNo = insertIndex + currentPageCaseIds.length;
      message.success(`成功重排 ${finalOrder.length} 条用例，当前页用例编号为 No.${startNo}-${endNo}`);
      
      // 关闭对话框并刷新表格
      if (onOk) {
        onOk();
      }
    } catch (error) {
      console.error('重排序失败:', error);
      message.error(error.response?.data?.message || 'ID重排失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title="ID重新排序"
      open={visible}
      onOk={handleReorder}
      onCancel={onCancel}
      confirmLoading={loading}
      okText="确定"
      cancelText="取消"
      width={600}
    >
      <div style={{ marginBottom: 16 }}>
        <span style={{ color: '#1890ff', marginRight: 8 }}>ℹ️</span>
        <strong>ID重排说明</strong>
      </div>
      
      <p>将按照当前页位置重新分配ID，当前页的用例会按页码对应的编号范围排列。</p>
      
      <div style={{ marginTop: 16, padding: 12, background: '#e6f7ff', border: '1px solid #91d5ff', borderRadius: 4 }}>
        <strong>📋 重排规则：</strong>
        <ul style={{ marginTop: 8, marginBottom: 0, paddingLeft: 20 }}>
          <li>当前是<strong>第{currentPage}页</strong>，当前页用例将编号为 <strong>No.{(currentPage - 1) * 10 + 1}-{(currentPage - 1) * 10 + cases.length}</strong></li>
          <li>前面的页保持原有顺序（如第1页为No.1-10）</li>
          <li>后面的页顺延排列</li>
          <li>重排后会自动恢复为10条/页</li>
        </ul>
      </div>
      
      <div style={{ marginTop: 16, padding: 12, background: '#fff7e6', border: '1px solid #ffd591', borderRadius: 4 }}>
        <strong>⚠️ 注意：</strong>
        <ul style={{ marginTop: 8, marginBottom: 0, paddingLeft: 20 }}>
          <li>重排操作不可撤销</li>
          <li>{t('common.total')} <strong>{totalCount}</strong> {t('common.items')} {t('manualTest.cases')}, {t('manualTest.currentPageHas')} <strong>{cases.length}</strong> {t('common.items')}</li>
          {(caseType === 'overall' || caseType === 'change' || caseType === 'acceptance') && (
            <li>此操作会同步更新所有语言版本的用例ID</li>
          )}
        </ul>
      </div>
    </Modal>
  );
};

export default ReorderModal;
