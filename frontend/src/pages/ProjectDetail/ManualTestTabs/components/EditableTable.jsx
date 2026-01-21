import React, { useState, useEffect, useMemo, useRef, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Table, Button, Space, message, Input, Select, Form, Upload, Alert, Modal, Tooltip } from 'antd';
import { DeleteOutlined, DownloadOutlined, UploadOutlined, SaveOutlined, PlusCircleOutlined, CopyOutlined, ArrowUpOutlined, ArrowDownOutlined, EyeOutlined } from '@ant-design/icons';
import { getCasesList, createCase, updateCase, deleteCase, exportAICases, exportCases, importCases, clearAICases, insertCase, batchDeleteCases, reassignAllIDs } from '../../../../api/manualCase';
import { getAutoCasesList, createAutoCase, updateAutoCase, deleteAutoCase, exportAutoCases, insertAutoCase, batchDeleteAutoCases, reassignAutoIDs } from '../../../../api/autoCase';
import { getApiCasesList, createApiCase, updateApiCase, deleteApiCase, insertApiCase, batchDeleteApiCases } from '../../../../api/apiCase';
import { getExecutionCaseResults, saveExecutionCaseResults } from '../../../../api/executionCaseResult';
import useDebounce from '../hooks/useDebounce';
import MultiLangEditModal from './MultiLangEditModal';
import ApiCaseDetailModal from '../../ApiTestTabs/components/ApiCaseDetailModal';
import WebCaseDetailModal from '../../AutoTestTabs/components/WebCaseDetailModal';
import { maskKnownPasswords } from '../../../../utils/maskPassword';
import './EditableTable.css';

const { TextArea } = Input;

// ==================== 工具函数区域 ====================

/**
 * 复制分类字段工具函数
 * @param {Object} sourceCase - 源用例对象
 * @param {string} caseType - 用例类型 ('ai' | 'overall' | 'change' | 'acceptance')
 * @returns {Object} 包含复制的分类字段的对象
 */
const copyClassificationFields = (sourceCase, caseType) => {
  console.log('[copyClassificationFields] called with caseType:', caseType);
  console.log('[copyClassificationFields] sourceCase:', sourceCase);

  if (!sourceCase) {
    console.log('[copyClassificationFields] sourceCase is null/undefined, returning {}');
    return {};
  }

  let result;
  switch (caseType) {
    case 'ai':
      // AI用例：复制单语言字段
      result = {
        major_function: sourceCase.major_function || '',
        middle_function: sourceCase.middle_function || '',
      };
      break;

    case 'overall':
    case 'change':
    case 'acceptance':
      // 整体/变更/受入用例：复制多语言字段（CN/JP/EN）
      result = {
        major_function_cn: sourceCase.major_function_cn || '',
        middle_function_cn: sourceCase.middle_function_cn || '',
        major_function_jp: sourceCase.major_function_jp || '',
        middle_function_jp: sourceCase.middle_function_jp || '',
        major_function_en: sourceCase.major_function_en || '',
        middle_function_en: sourceCase.middle_function_en || '',
      };
      break;

    default:
      // 未知类型返回空对象
      result = {};
  }

  console.log('[copyClassificationFields] returning fields:', result);
  return result;
};

/**
 * 计算单元格样式工具函数
 * @param {Object} currentCase - 当前行用例对象
 * @param {Object|null} previousCase - 上一行用例对象（第一行为null）
 * @param {string} fieldName - 字段名
 * @returns {Object} React样式对象
 */
const getCellStyle = (currentCase, previousCase, fieldName) => {
  // 第一行无颜色
  if (!previousCase) {
    return {};
  }

  const currentValue = currentCase[fieldName];
  const previousValue = previousCase[fieldName];

  // 空值不应用颜色
  if (!currentValue || !previousValue) {
    return {};
  }

  // 不同值无颜色
  if (currentValue !== previousValue) {
    return {};
  }

  // 相同值，根据字段后缀判断颜色
  if (fieldName.endsWith('_cn')) {
    // 中文字段：极浅的蓝色文字，无背景
    return { color: '#91D5FF' };
  } else if (fieldName.endsWith('_jp') || fieldName.endsWith('_en')) {
    // 英文/日文字段：很浅的灰色文字，无背景
    return { color: '#BFBFBF' };
  } else if (fieldName === 'major_function' || fieldName === 'middle_function') {
    // AI用例字段：浅灰色
    return { backgroundColor: '#F5F5F5' };
  }

  return {};
};

// ==================== 组件定义 ====================

/**
 * 可编辑表格组件
 * @param {Object} props
 * @param {string} props.caseType - 用例类型 ('ai'|'overall'|'change'|'role1-4')
 * @param {string} props.language - 当前语言 ('中文'|'English'|'日本語')
 * @param {Function} props.onRefreshMetadata - 刷新元数据回调
 * @param {string} props.apiModule - API模块标识 ('api-cases' 表示接口测试用例)
 * @param {string} props.projectId - 项目ID (当apiModule='api-cases'时需要传入)
 * @param {boolean} props.executionMode - 是否为执行模式 (true=执行模式，仅编辑test_result/bug_id/remark)
 * @param {boolean} props.selectionMode - 是否为选择模式 (true=选择模式，只读展示)
 * @param {string} props.taskUuid - 执行任务UUID (executionMode=true时必填)
 * @param {Function} props.onResultsChange - 执行结果变更回调 (executionMode=true时可选)
 * @param {Function} props.onCasesLoaded - 用例加载完成回调 (selectionMode=true时可选)
 * @param {Array<string>} props.hiddenButtons - 需要隐藏的按钮标识数组，可选值: ['saveVersion', 'exportTemplate', 'aiSupplement', 'exportCases', 'importCases']
 * @param {string} props.caseGroupFilter - 用例集过滤器，只显示指定用例集的用例
 * @param {Function} props.onBatchDeleteRequest - 批量删除请求回调，返回{selectedCount, executeDelete}对象
 * @param {Array<string>} props.knownPasswords - 已知密码列表（用于脱敏显示）
 * @param {number} props.caseGroupId - 用例集ID（用于脚本测试时获取变量）
 */
const EditableTable = ({ caseType, language, onRefreshMetadata, apiModule, projectId: propsProjectId, executionMode = false, selectionMode = false, taskUuid, onResultsChange, onCasesLoaded, hiddenButtons = [], caseGroupFilter = null, onBatchDeleteRequest, knownPasswords = [], caseGroupId }) => {
  const { t } = useTranslation();

  // 调试日志
  console.log('[EditableTable] Props:', { caseType, language, caseGroupFilter, caseGroupId });
  const { id: paramProjectId } = useParams();
  const projectId = propsProjectId || paramProjectId; // 优先使用props传入的projectId
  const [cases, setCases] = useState([]);
  const [loading, setLoading] = useState(false);

  // 执行模式：存储执行结果 (Map: case_id -> {test_result, bug_id, remark})
  const [executionResults, setExecutionResults] = useState(new Map());

  // 从 sessionStorage 恢复页码，实现F5刷新后保持当前页
  const getStorageKey = () => `editable_table_page_${projectId}_${caseType}_${language}`;
  const getSavedPage = () => {
    try {
      const saved = sessionStorage.getItem(getStorageKey());
      return saved ? parseInt(saved, 10) : 1;
    } catch {
      return 1;
    }
  };

  const [pagination, setPagination] = useState({
    current: getSavedPage(),
    pageSize: 10,
    total: 0,
  });
  const [editingKey, setEditingKey] = useState('');
  const [form] = Form.useForm();

  // 存储所有删除按钮的 ref
  const deleteButtonRefs = useRef({});

  // 追踪是否是组件初始挂载
  const isInitialMount = useRef(true);

  // 多语言编辑对话框状态
  const [multiLangModalVisible, setMultiLangModalVisible] = useState(false);
  const [multiLangData, setMultiLangData] = useState({
    record: null,
    fieldName: '',
    title: '',
    cn: '',
    jp: '',
    en: '',
  });

  // T44: 移除批量修改确认对话框状态
  // const [batchConfirmVisible, setBatchConfirmVisible] = useState(false);
  // const [batchConfirmData, setBatchConfirmData] = useState({...});

  // 标记是否有编辑变更
  const [hasEditChanges, setHasEditChanges] = useState(false);

  // API用例详情弹窗状态
  const [apiDetailModalVisible, setApiDetailModalVisible] = useState(false);
  const [apiDetailCaseData, setApiDetailCaseData] = useState(null);

  // Web用例详情弹窗状态
  const [webDetailModalVisible, setWebDetailModalVisible] = useState(false);
  const [webDetailCaseData, setWebDetailCaseData] = useState(null);

  // 批量删除选择状态
  const [selectedRowKeys, setSelectedRowKeys] = useState([]);

  // 将批量删除功能暴露给父组件
  useEffect(() => {
    if (onBatchDeleteRequest) {
      onBatchDeleteRequest({
        selectedCount: selectedRowKeys.length,
        executeDelete: handleBatchDelete
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedRowKeys, onBatchDeleteRequest]);

  // 高亮显示的行ID
  const [highlightedRowId, setHighlightedRowId] = useState(null);

  // 记录本地插入的空行位置 {targetCaseId: string, position: 'before'|'after'}[]
  const [pendingInserts, setPendingInserts] = useState([]);

  // 表格容器ref，用于滚动
  const tableContainerRef = useRef(null);

  // 统一的更新用例API方法(根据caseType选择调用manualCase或autoCase或apiCase)
  const updateCaseAPI = useCallback(async (caseId, updates) => {
    console.log('[updateCaseAPI] caseId:', caseId, 'updates:', updates);
    if (!caseId) {
      console.error('[updateCaseAPI] caseId为空，无法更新！');
      throw new Error('caseId不能为空');
    }
    let apiCall;
    if (apiModule === 'api-cases') {
      apiCall = updateApiCase;
    } else {
      const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');
      apiCall = isAutoType ? updateAutoCase : updateCase;
    }
    return apiCall(projectId, caseId, updates);
  }, [projectId, caseType, apiModule]);

  // 加载用例列表 - 使用 useCallback 避免依赖警告
  const fetchCases = useCallback(async (page = pagination.current, customPageSize = null) => {
    try {
      setLoading(true);

      // 判断使用哪个API: api-cases模块 > role类型 > manual类型
      let apiCall;
      if (apiModule === 'api-cases') {
        apiCall = getApiCasesList;
      } else {
        const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');
        apiCall = isAutoType ? getAutoCasesList : getCasesList;
      }

      // 如果指定了自定义页面大小，使用自定义大小，否则使用当前pageSize
      const requestSize = customPageSize || pagination.pageSize;

      const data = await apiCall(projectId, {
        caseType,
        language,
        page,
        size: requestSize,
        caseGroup: caseGroupFilter, // 传递用例集过滤参数到后端
      });

      console.log('[EditableTable] API返回数据:', { total: data.total, casesCount: data.cases?.length, caseGroupFilter });
      console.log('[EditableTable] API返回的前3条数据:', data.cases?.slice(0, 3));

      // 如果没有数据，显示空表格（不自动创建空行）
      let casesData = data.cases || [];

      // 根据用例集过滤 - 后端已经按 case_group 过滤了，这里只需要处理 null 的情况
      if (caseGroupFilter === null) {
        // 没有选中用例集时，显示空表格
        console.log('[EditableTable] 没有选中用例集，显示空表格');
        casesData = [];
      }

      // 执行模式：加载执行结果并合并到用例数据
      if (executionMode && taskUuid) {
        try {
          const execResults = await getExecutionCaseResults(taskUuid);
          const resultsMap = new Map();
          execResults.forEach(result => {
            resultsMap.set(result.case_id, {
              test_result: result.test_result,
              bug_id: result.bug_id,
              remark: result.remark,
            });
          });
          setExecutionResults(resultsMap);

          // 合并执行结果到用例数据
          casesData = casesData.map(c => {
            const execResult = resultsMap.get(c.case_id);
            if (execResult) {
              return {
                ...c,
                test_result: execResult.test_result,
                bug_id: execResult.bug_id,
                remark: execResult.remark,
              };
            }
            return c;
          });
        } catch (error) {
          console.error('[fetchCases] Failed to load execution results:', error);
          message.error('加载执行结果失败');
        }
      }

      console.log('[EditableTable] 设置cases数据，数量:', casesData.length);
      setCases(casesData);
      setPagination(prev => ({
        ...prev,
        current: page,
        total: data.total || 0,
        // 如果使用了自定义页面大小，更新pageSize
        pageSize: requestSize,
      }));
      console.log('[EditableTable] 设置pagination:', { current: page, total: data.total, pageSize: requestSize });

      // 选择模式：通知父组件用例数据已加载
      if (selectionMode && onCasesLoaded) {
        // 传递所有用例数据（包含total信息）
        onCasesLoaded(casesData, data.total || 0);
      }

      // 保存当前页码到 sessionStorage，实现F5刷新后保持当前页
      try {
        sessionStorage.setItem(getStorageKey(), String(page));
      } catch (e) {
        console.warn('Failed to save page to sessionStorage:', e);
      }

      // 重新加载数据后,清除编辑变更标记
      setHasEditChanges(false);
    } catch (error) {
      console.error('Failed to load cases:', error);
      message.error('加载用例列表失败');
    } finally {
      setLoading(false);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [projectId, caseType, language, pagination.current, pagination.pageSize, executionMode, selectionMode, taskUuid, onCasesLoaded, caseGroupFilter]);

  // 防抖保存 (保留用于未来实现内联编辑功能)
  // eslint-disable-next-line no-unused-vars
  const debouncedSave = useDebounce(async (caseId, field, value) => {
    try {
      await updateCaseAPI(caseId, { [field]: value });
      message.success('保存成功');
    } catch (error) {
      console.error('Failed to update case:', error);
      message.error('保存失败');
    }
  }, 500);

  // 导出AI用例
  const handleExportAICases = async () => {
    try {
      await exportAICases(projectId);
      message.success('导出成功');
    } catch (error) {
      console.error('导出AI用例失败:', error);
      message.error('导出失败');
    }
  };

  // T44: handleExportTemplate 已移除，由ManualCaseManagementTab工具栏提供

  // 导出用例
  const handleExportCases = async () => {
    try {
      const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');
      const apiCall = isAutoType ? exportAutoCases : exportCases;
      // 执行模式：传递taskUuid参数以包含执行结果列
      if (executionMode && taskUuid) {
        await apiCall(projectId, caseType, taskUuid);
      } else {
        await apiCall(projectId, caseType);
      }
      message.success('导出用例成功');
    } catch (error) {
      console.error('导出用例失败:', error);
      message.error('导出用例失败');
    }
  };

  // 导入用例
  const handleImportCases = async (file) => {
    try {
      const result = await importCases(projectId, caseType, file);
      message.success(`导入成功: ${result.insertCount}条新增, ${result.updateCount}条更新`);
      fetchCases(1); // 刷新到第一页
      if (onRefreshMetadata) {
        onRefreshMetadata();
      }
    } catch (error) {
      console.error('导入用例失败:', error);
      message.error('导入用例失败');
    }
    return false; // 阻止Upload组件自动上传
  };

  // T44: handleSaveVersion 已移除，版本管理功能已移至独立模块

  // 创建默认空行(调用后端创建真实记录)
  const createDefaultEmptyRow = useCallback(async () => {
    try {
      let createAPI;
      let createData;

      if (apiModule === 'api-cases') {
        // API测试用例
        createAPI = createApiCase;
        createData = {
          project_id: parseInt(projectId),
          case_type: caseType || 'api',
          case_group: caseGroupFilter || '', // 添加用例集筛选条件
          method: 'GET',
          test_result: 'NR',
        };
      } else {
        // 手工用例或自动化用例
        const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');
        createAPI = isAutoType ? createAutoCase : createCase;
        createData = {
          case_type: caseType,
        };
        // 手工用例需要language参数
        if (!isAutoType) {
          createData.language = language;
        }
        // 如果有 caseGroupFilter，添加到创建数据中
        if (caseGroupFilter) {
          createData.case_group = caseGroupFilter;
        }
      }

      console.log('[createDefaultEmptyRow] Creating with data:', createData);

      // 创建一条空记录
      const newCase = await createAPI(projectId, createData);
      console.log('[createDefaultEmptyRow] Created case:', newCase);
      return newCase;
    } catch (error) {
      console.error('[createDefaultEmptyRow] Failed to create default empty row:', error);
      console.error('[createDefaultEmptyRow] Error response:', error?.response?.data);
      throw error;
    }
  }, [projectId, caseType, language, apiModule, caseGroupFilter]);

  // 用例补足 - 自动填充空的分类字段
  // 用例补足 - 自动填充空的分类字段(一键补充所有用例)
  const handleFillClassification = async () => {
    try {
      setLoading(true);

      const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');
      const listAPI = isAutoType ? getAutoCasesList : getCasesList;
      const updateAPI = isAutoType ? updateAutoCase : updateCase;

      // 获取所有用例(不分页)
      const allCasesData = await listAPI(projectId, {
        caseType,
        language,
        page: 1,
        size: 9999, // 获取所有用例
      });

      let allCases = allCasesData.cases || [];

      if (allCases.length <= 1) {
        message.info('至少需要2条用例才能执行补足操作');
        setLoading(false);
        return;
      }

      let updateCount = 0;

      // 遍历所有用例,从第2条开始(第1条没有上一条可复制)
      for (let i = 1; i < allCases.length; i++) {
        const currentCase = allCases[i];
        const previousCase = allCases[i - 1];

        const updates = {};

        if (caseType === 'ai') {
          // AI用例:检查单语言字段是否为空
          if (!currentCase.major_function && previousCase.major_function) {
            updates.major_function = previousCase.major_function;
          }
          if (!currentCase.middle_function && previousCase.middle_function) {
            updates.middle_function = previousCase.middle_function;
          }
        } else {
          // 其他用例:检查多语言字段是否为空(CN/JP/EN三种语言)
          if (!currentCase.major_function_cn && previousCase.major_function_cn) {
            updates.major_function_cn = previousCase.major_function_cn;
          }
          if (!currentCase.middle_function_cn && previousCase.middle_function_cn) {
            updates.middle_function_cn = previousCase.middle_function_cn;
          }
          if (!currentCase.major_function_jp && previousCase.major_function_jp) {
            updates.major_function_jp = previousCase.major_function_jp;
          }
          if (!currentCase.middle_function_jp && previousCase.middle_function_jp) {
            updates.middle_function_jp = previousCase.middle_function_jp;
          }
          if (!currentCase.major_function_en && previousCase.major_function_en) {
            updates.major_function_en = previousCase.major_function_en;
          }
          if (!currentCase.middle_function_en && previousCase.middle_function_en) {
            updates.middle_function_en = previousCase.middle_function_en;
          }
        }

        // 如果有需要更新的字段,调用API更新
        if (Object.keys(updates).length > 0) {
          console.log(`[用例补足] 更新用例 No.${i + 1}:`, updates);
          // updateCase需要projectId参数,updateAutoCase不需要
          if (isAutoType) {
            await updateAPI(currentCase.case_id, updates);
          } else {
            await updateAPI(projectId, currentCase.case_id, updates);
          }

          // 【关键修复】立即更新内存中的数据,确保后续遍历使用最新值
          allCases[i] = { ...currentCase, ...updates };

          updateCount++;
        }
      }

      if (updateCount > 0) {
        message.success(`已补足 ${updateCount} 条用例的分类字段`);
        // 刷新当前页
        await fetchCases(pagination.current);
      } else {
        message.info('没有需要补足的用例');
      }
    } catch (error) {
      console.error('用例补足失败:', error);
      message.error('用例补足失败: ' + (error?.message || ''));
    } finally {
      setLoading(false);
    }
  };

  // 清空AI用例
  const handleClearAICases = async () => {
    const confirmed = window.confirm('确定要清空所有AI用例吗?此操作不可恢复!');
    if (confirmed) {
      try {
        await clearAICases(projectId);
        message.success('清空成功');

        // 清空后创建一条默认空行
        try {
          const newCase = await createDefaultEmptyRow();
          setCases([newCase]);
          setPagination(prev => ({ ...prev, current: 1, total: 1 }));
        } catch (error) {
          // 如果创建失败,显示空表格
          setCases([]);
          setPagination(prev => ({ ...prev, current: 1, total: 0 }));
        }

        if (onRefreshMetadata) {
          onRefreshMetadata();
        }
      } catch (error) {
        console.error('清空AI用例失败:', error);
        message.error('清空失败');
      }
    }
  };

  // 扩展当前页加载用例(插入后调用，从当前页起始位置多加载几条)
  const fetchExpandedPage = useCallback(async () => {
    try {
      setLoading(true);

      // 判断是否为role类型
      const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');
      const apiCall = isAutoType ? getAutoCasesList : getCasesList;

      // 获取当前页的起始位置
      const currentPage = pagination.current || 1;
      const pageSize = pagination.pageSize || 10;

      // 简单策略：直接请求当前页+扩展的数据
      const expandedSize = pageSize + 2;

      const data = await apiCall(projectId, {
        caseType,
        language,
        page: currentPage,
        size: expandedSize,
      });

      // 如果当前页没有数据，说明插入后数据位置发生了变化
      // 尝试获取前一页的数据
      if (!data.cases || data.cases.length === 0) {
        const prevPage = Math.max(1, currentPage - 1);
        const prevData = await apiCall(projectId, {
          caseType,
          language,
          page: prevPage,
          size: expandedSize,
        });

        setCases(prevData.cases || []);
        setPagination(prev => ({
          ...prev,
          current: prevPage,
          total: prevData.total || prev.total + 1,
        }));
      } else {
        setCases(data.cases || []);
        setPagination(prev => ({
          ...prev,
          current: currentPage,
          total: data.total || prev.total + 1,
        }));
      }

      setHasEditChanges(true);
    } catch (error) {
      console.error('Failed to load expanded page:', error);
      message.error('加载用例失败');
    } finally {
      setLoading(false);
    }
  }, [projectId, caseType, language, pagination]);

  // 防抖版本的fetchExpandedPage，用于合并连续的插入操作
  const debouncedFetchExpandedPage = useDebounce(fetchExpandedPage, 300);

  // 删除用例 - 使用 useCallback 确保引用稳定
  const handleDelete = useCallback(async (record) => {
    console.log('[handleDelete] 开始删除, record:', record);
    console.log('[handleDelete] apiModule:', apiModule);
    console.log('[handleDelete] projectId:', projectId);
    console.log('[handleDelete] case_id:', record.case_id);

    // 使用原生确认对话框替代 Modal.confirm
    const confirmed = window.confirm(
      t('manualTest.deleteCaseConfirm', { caseName: record.case_number || record.case_num || record.case_id })
    );

    if (!confirmed) {
      console.log('[handleDelete] 用户取消删除');
      return;
    }

    console.log('[handleDelete] 用户确认删除');

    try {
      // 判断使用哪个API
      let apiCall;
      if (apiModule === 'api-cases') {
        apiCall = deleteApiCase;
        console.log('[handleDelete] 使用 deleteApiCase API');
      } else {
        const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');
        apiCall = isAutoType ? deleteAutoCase : deleteCase;
        console.log('[handleDelete] 使用其他删除 API');
      }

      console.log('[handleDelete] 调用删除API...');
      await apiCall(projectId, record.case_id);
      console.log('[handleDelete] 删除API调用成功');
      message.success('删除成功');

      // 计算删除后应该停留的页码
      const currentPageSize = pagination.pageSize || 10;
      const totalAfterDelete = pagination.total - 1;
      const maxPage = Math.ceil(totalAfterDelete / currentPageSize) || 1;
      const targetPage = pagination.current > maxPage ? maxPage : pagination.current;

      console.log('[handleDelete] 当前页:', pagination.current, '总数:', pagination.total);
      console.log('[handleDelete] 删除后总数:', totalAfterDelete, '最大页:', maxPage, '目标页:', targetPage);
      console.log('[handleDelete] 调用fetchCases刷新数据...');

      // 立即刷新数据
      await fetchCases(targetPage);
      console.log('[handleDelete] fetchCases完成');
    } catch (error) {
      console.error('[handleDelete] 删除失败:', error);
      console.error('[handleDelete] 错误详情:', error.response);
      const errorMsg = error?.response?.data?.message || error?.message || '删除失败';
      message.error(errorMsg);
    }
  }, [projectId, fetchCases, caseType, apiModule, pagination]);

  // Helper: 根据apiModule创建空行对象
  const createEmptyRowByModule = useCallback(async (targetRow, targetNo) => {
    console.log('[createEmptyRowByModule] targetRow.case_type:', targetRow.case_type);
    console.log('[createEmptyRowByModule] apiModule:', apiModule);

    // 对于api-cases，使用临时ID（与其他用例一致）
    if (apiModule === 'api-cases') {
      return {
        case_id: `temp_${Date.now()}_${Math.random()}`, // 临时ID
        id: targetNo,
        display_id: targetNo,
        project_id: parseInt(projectId),
        case_type: targetRow.case_type || caseType, // 使用目标行的类型，如果没有则使用组件的caseType prop
        case_group: targetRow.case_group || caseGroupFilter || '', // 继承用例集
        case_number: '', // 空字符串，用户可自定义，显示为0
        screen: '',
        url: '',
        header: '',
        method: 'GET',
        body: '',
        response: '',
        test_result: 'NR',
        remark: '',
        created_at: '',
        updated_at: '',
        _isNew: true,
      };
    }

    // 手工/自动化测试用例使用临时ID
    const baseRow = {
      case_id: `temp_${Date.now()}_${Math.random()}`, // 临时ID
      id: targetNo,
      display_id: targetNo,
      project_id: targetRow.project_id || parseInt(projectId),
      case_type: targetRow.case_type || caseType, // 使用目标行的类型，如果没有则使用组件的caseType prop
      case_number: '',
      case_group: targetRow.case_group || caseGroupFilter || '', // 继承用例集
      test_result: '',
      remark: '',
      created_at: '',
      updated_at: '',
      created_by: '',
      updated_by: '',
      _isNew: true,
    };

    // 判断用例类型，创建对应字段
    if (targetRow.case_type === 'ai') {
      // AI用例：使用单语言字段（不带后缀）
      return {
        ...baseRow,
        major_function: '',
        middle_function: '',
        minor_function: '',
        precondition: '',
        test_steps: '',
        expected_result: '',
      };
    } else if (targetRow.case_type === 'web') {
      // Web用例：使用多语言字段
      return {
        ...baseRow,
        screen_cn: '',
        screen_en: '',
        screen_jp: '',
        function_cn: '',
        function_en: '',
        function_jp: '',
        precondition_cn: '',
        precondition_en: '',
        precondition_jp: '',
        test_steps_cn: '',
        test_steps_en: '',
        test_steps_jp: '',
        expected_result_cn: '',
        expected_result_en: '',
        expected_result_jp: '',
      };
    } else {
      // 手动测试用例（overall/change）：使用多语言字段
      return {
        ...baseRow,
        major_function_cn: '',
        middle_function_cn: '',
        minor_function_cn: '',
        major_function_en: '',
        middle_function_en: '',
        minor_function_en: '',
        major_function_jp: '',
        middle_function_jp: '',
        minor_function_jp: '',
        precondition_cn: '',
        test_steps_cn: '',
        expected_result_cn: '',
        precondition_en: '',
        test_steps_en: '',
        expected_result_en: '',
        precondition_jp: '',
        test_steps_jp: '',
        expected_result_jp: '',
      };
    }
  }, [apiModule, projectId, caseGroupFilter]);

  // 在指定行上方插入（本地操作）
  const handleInsertAbove = useCallback(async (targetCaseId) => {
    // 先取消所有正在编辑的行
    if (editingKey) {
      setEditingKey('');
      form.resetFields();
    }

    try {
      setLoading(true);

      // 所有模式：使用本地插入模式
      // 先设置hasEditChanges，防止useEffect触发fetchCases
      setHasEditChanges(true);
      setPendingInserts(prev => [...prev, { targetCaseId, position: 'before' }]);

      // 找到目标行
      const targetIndex = cases.findIndex(c => c.case_id === targetCaseId);
      if (targetIndex === -1) {
        message.error('未找到目标行');
        return;
      }

      const targetRow = cases[targetIndex];
      const targetNo = (pagination.current - 1) * pagination.pageSize + targetIndex + 1;

      console.log('[handleInsertAbove] targetIndex:', targetIndex);
      console.log('[handleInsertAbove] targetRow:', JSON.stringify(targetRow, null, 2));
      console.log('[handleInsertAbove] targetRow.case_type:', targetRow.case_type);
      console.log('[handleInsertAbove] caseType prop:', caseType);

      // 创建一个空行 - 对于api-cases会调用API创建真实UUID
      let emptyRow = await createEmptyRowByModule(targetRow, targetNo);
      console.log('[handleInsertAbove] emptyRow after create:', emptyRow);

      // 【T35需求已修改】不再自动复制分类字段
      // 用户可使用"用例补足"按钮统一填充空的分类字段
      // const sourceIndex = targetIndex - 1;
      // if (sourceIndex >= 0) {
      //   const sourceCase = cases[sourceIndex];
      //   const classificationFields = copyClassificationFields(sourceCase, caseType);
      //   emptyRow = { ...emptyRow, ...classificationFields };
      // }

      // 【修正】使用函数式setState，基于最新状态计算
      setCases(prevCases => {
        const newCases = [...prevCases];
        newCases.splice(targetIndex, 0, emptyRow);

        // api-cases不需要重新计算No.,No.由后端display_order决定
        // 其他类型保持原逻辑
        if (apiModule !== 'api-cases') {
          for (let i = targetIndex + 1; i < newCases.length; i++) {
            newCases[i] = {
              ...newCases[i],
              id: newCases[i].id + 1,
              display_id: newCases[i].id + 1,
            };
          }
        }

        return newCases;
      });
      message.success('已插入空行，点击保存后生效');
    } catch (error) {
      console.error('[handleInsertAbove] 插入失败:', error);
      message.error('插入失败: ' + (error.response?.data?.error || error.message));
    } finally {
      setLoading(false);
    }
  }, [editingKey, form, cases, createEmptyRowByModule, apiModule, projectId, caseType, fetchCases]);

  // 在指定行下方插入（本地操作）
  const handleInsertBelow = useCallback(async (targetCaseId) => {
    // 先取消所有正在编辑的行
    if (editingKey) {
      setEditingKey('');
      form.resetFields();
    }

    try {
      setLoading(true);

      // 所有模式：使用本地插入模式
      // 先设置hasEditChanges，防止useEffect触发fetchCases
      setHasEditChanges(true);
      setPendingInserts(prev => [...prev, { targetCaseId, position: 'after' }]);

      // 找到目标行
      const targetIndex = cases.findIndex(c => c.case_id === targetCaseId);
      if (targetIndex === -1) {
        message.error('未找到目标行');
        return;
      }

      const targetRow = cases[targetIndex];
      const targetNo = (pagination.current - 1) * pagination.pageSize + targetIndex + 1;

      console.log('[handleInsertBelow] targetIndex:', targetIndex);
      console.log('[handleInsertBelow] targetRow.case_type:', targetRow.case_type);

      // 创建一个空行 - 对于api-cases会调用API创建真实UUID
      let emptyRow = await createEmptyRowByModule(targetRow, targetNo + 1);
      console.log('[handleInsertBelow] emptyRow after create:', emptyRow);

      // 【T35需求已修改】不再自动复制分类字段
      // 用户可使用"用例补足"按钮统一填充空的分类字段
      // const sourceCase = cases[targetIndex];
      // const classificationFields = copyClassificationFields(sourceCase, caseType);
      // emptyRow = { ...emptyRow, ...classificationFields };

      // 【修正】使用函数式setState，基于最新状态计算
      setCases(prevCases => {
        const newCases = [...prevCases];
        newCases.splice(targetIndex + 1, 0, emptyRow);

        // 重新计算当前页的No.（插入行之后的所有行+1）
        // 注意：api-cases使用UUID，不需要重新计算id
        if (apiModule !== 'api-cases') {
          for (let i = targetIndex + 2; i < newCases.length; i++) {
            newCases[i] = {
              ...newCases[i],
              id: newCases[i].id + 1,
              display_id: newCases[i].id + 1,
            };
          }
        }

        return newCases;
      });
      message.success('已插入空行，点击保存后生效');
    } catch (error) {
      console.error('[handleInsertBelow] 插入失败:', error);
      message.error('插入失败: ' + (error.response?.data?.error || error.message));
    } finally {
      setLoading(false);
    }
  }, [editingKey, form, cases, createEmptyRowByModule, apiModule, projectId, caseType, fetchCases]);

  // 批量删除用例
  const handleBatchDelete = useCallback(async () => {
    const confirmed = window.confirm(
      t('manualTest.batchDeleteConfirm', { count: selectedRowKeys.length })
    );
    if (!confirmed) return;

    try {
      // 判断使用哪个API
      let apiCall;
      if (apiModule === 'api-cases') {
        apiCall = batchDeleteApiCases;
      } else {
        const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');
        apiCall = isAutoType ? batchDeleteAutoCases : batchDeleteCases;
      }

      await apiCall(projectId, {
        caseType,
        caseIds: selectedRowKeys,
      });

      message.success(`成功删除${selectedRowKeys.length}条用例`);
      setSelectedRowKeys([]); // 清空选择

      // 计算删除后应该停留的页码
      const deletedCount = selectedRowKeys.length;
      const currentPageSize = pagination.pageSize || 10;
      const totalAfterDelete = pagination.total - deletedCount;
      const maxPage = Math.ceil(totalAfterDelete / currentPageSize) || 1;
      const targetPage = pagination.current > maxPage ? maxPage : pagination.current;

      // 立即刷新数据
      await fetchCases(targetPage);
    } catch (error) {
      console.error('Failed to batch delete cases:', error);
      const errorMsg = error?.response?.data?.message || error?.message || '批量删除失败';
      message.error(errorMsg);
    }
  }, [projectId, caseType, selectedRowKeys, fetchCases, apiModule, pagination]);

  // 组件挂载和语言切换时加载数据
  // 【新需求】如果有未保存的编辑变更,不要重新加载数据
  // 【修复】切换用例集时重置页码到第1页
  useEffect(() => {
    // 初始挂载时无条件加载数据，使用保存的页码
    if (isInitialMount.current) {
      isInitialMount.current = false;
      const savedPage = getSavedPage();
      fetchCases(savedPage);
      return;
    }

    // 后续触发时,只有在没有未保存变更时才加载数据
    // 切换语言时保持当前页码，切换用例集时重置到第1页
    if (!hasEditChanges) {
      fetchCases(1); // 重置到第1页
      setPagination(prev => ({ ...prev, current: 1 })); // 同时更新分页状态
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [projectId, language, caseType, caseGroupFilter]);

  // 为所有删除按钮绑定原生事件监听器(解决 React 合成事件失效问题)
  useEffect(() => {
    const listeners = [];

    cases.forEach(record => {
      const buttonEl = deleteButtonRefs.current[`delete-${record.case_id}`];

      if (buttonEl && editingKey !== record.case_id) {
        const handleNativeClick = (e) => {
          e.preventDefault();
          e.stopPropagation();
          handleDelete(record);
        };

        buttonEl.addEventListener('click', handleNativeClick, true);
        listeners.push({ buttonEl, handler: handleNativeClick });
      }
    });

    return () => {
      listeners.forEach(({ buttonEl, handler }) => {
        buttonEl.removeEventListener('click', handler, true);
      });
    };
  }, [cases, editingKey, handleDelete]);

  // 可编辑单元格组件 - 使用 Form.Item 管理输入
  const EditableCell = ({
    editing,
    dataIndex,
    title,
    inputType,
    record,
    index,
    children,
    ...restProps
  }) => {
    let inputNode;

    if (inputType === 'textarea') {
      inputNode = <TextArea autoSize={{ minRows: 2, maxRows: 6 }} />;
    } else if (inputType === 'select') {
      // 根据需求FR-03.1，TestResult选项为OK/NG/Block/NR
      inputNode = (
        <Select style={{ width: '100%' }}>
          <Select.Option value="NR">NR</Select.Option>
          <Select.Option value="OK">OK</Select.Option>
          <Select.Option value="NG">NG</Select.Option>
          <Select.Option value="Block">Block</Select.Option>
        </Select>
      );
    } else {
      inputNode = <Input />;
    }

    return (
      <td {...restProps}>
        {editing ? (
          <Form.Item
            name={dataIndex}
            style={{ margin: 0 }}
          >
            {inputNode}
          </Form.Item>
        ) : (
          children
        )}
      </td>
    );
  };

  // 开始编辑
  const startEdit = useCallback((record) => {
    // 字段值规范化函数：统一处理空值显示
    const getDisplayValue = (value) => {
      if (value === null || value === undefined) return '';
      return String(value);
    };

    // 判断是否为api-cases类型
    const isApiCasesMode = apiModule === 'api-cases';
    // 判断是否为自动化类型(role/web)
    const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');

    // 根据用例类型和语言决定设置哪些字段
    if (isApiCasesMode) {
      // API测试用例：设置api-cases专属字段
      const formValues = {
        screen: getDisplayValue(record.screen),
        url: getDisplayValue(record.url),
        header: getDisplayValue(record.header),
        method: getDisplayValue(record.method),
        body: getDisplayValue(record.body),
        response: getDisplayValue(record.response),
        test_result: getDisplayValue(record.test_result),
        remark: getDisplayValue(record.remark),
      };
      form.setFieldsValue(formValues);
    } else if (caseType === 'ai') {
      // AI用例：设置单语言字段
      const formValues = {
        case_number: getDisplayValue(record.case_number),
        major_function: getDisplayValue(record.major_function),
        middle_function: getDisplayValue(record.middle_function),
        minor_function: getDisplayValue(record.minor_function),
        precondition: getDisplayValue(record.precondition),
        test_steps: getDisplayValue(record.test_steps),
        expected_result: getDisplayValue(record.expected_result),
        remark: getDisplayValue(record.remark),
      };
      form.setFieldsValue(formValues);
    } else if (isAutoType) {
      // 自动化测试用例(role1-4/web)：根据当前语言设置对应的多语言字段
      const langFieldSuffix = language === '中文' ? '_cn' : language === 'English' ? '_en' : '_jp';

      const formValues = {
        case_number: getDisplayValue(record.case_number),
        [`screen${langFieldSuffix}`]: getDisplayValue(record[`screen${langFieldSuffix}`]),
        [`function${langFieldSuffix}`]: getDisplayValue(record[`function${langFieldSuffix}`]),
        [`precondition${langFieldSuffix}`]: getDisplayValue(record[`precondition${langFieldSuffix}`]),
        [`test_steps${langFieldSuffix}`]: getDisplayValue(record[`test_steps${langFieldSuffix}`]),
        [`expected_result${langFieldSuffix}`]: getDisplayValue(record[`expected_result${langFieldSuffix}`]),
        test_result: getDisplayValue(record.test_result),
        remark: getDisplayValue(record.remark),
      };
      form.setFieldsValue(formValues);
    } else {
      // 整体用例/变更用例：根据当前语言设置对应的多语言字段
      const langFieldSuffix = language === '中文' ? '_cn' : language === 'English' ? '_en' : '_jp';

      const formValues = {
        case_number: getDisplayValue(record.case_number),
        [`major_function${langFieldSuffix}`]: getDisplayValue(record[`major_function${langFieldSuffix}`]),
        [`middle_function${langFieldSuffix}`]: getDisplayValue(record[`middle_function${langFieldSuffix}`]),
        [`minor_function${langFieldSuffix}`]: getDisplayValue(record[`minor_function${langFieldSuffix}`]),
        [`precondition${langFieldSuffix}`]: getDisplayValue(record[`precondition${langFieldSuffix}`]),
        [`test_steps${langFieldSuffix}`]: getDisplayValue(record[`test_steps${langFieldSuffix}`]),
        [`expected_result${langFieldSuffix}`]: getDisplayValue(record[`expected_result${langFieldSuffix}`]),
        test_result: getDisplayValue(record.test_result),
        remark: getDisplayValue(record.remark),
      };
      form.setFieldsValue(formValues);
    }
    setEditingKey(record.case_id);
  }, [form, caseType, language, apiModule]);

  // 保存编辑
  const saveEdit = useCallback(async (recordParam) => {
    try {
      // 先验证表单
      const row = await form.validateFields();

      // 从 cases 数组中获取最新的 record 数据
      const currentRecord = cases.find(c => c.case_id === recordParam.case_id);
      if (!currentRecord) {
        console.error('Cannot find current record in cases array');
        message.error('无法找到当前记录');
        return;
      }

      // 字段值规范化函数：统一处理空值
      const normalizeValue = (value) => {
        if (value === null || value === undefined) return '';
        return String(value).trim();
      };

      const updates = {};

      // 判断是否为自动化类型(role/web)
      const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');
      // 判断是否为api-cases类型
      const isApiCasesMode = apiModule === 'api-cases';

      // 根据用例类型处理字段映射
      if (isApiCasesMode) {
        // API测试用例：处理api-cases的字段
        Object.keys(row).forEach(key => {
          const formValue = normalizeValue(row[key]);
          const recordValue = normalizeValue(currentRecord[key]);
          if (formValue !== recordValue) {
            updates[key] = row[key] === undefined || row[key] === null ? '' : String(row[key]);
          }
        });
      } else if (caseType === 'ai') {
        // AI用例：直接使用单语言字段
        Object.keys(row).forEach(key => {
          const formValue = normalizeValue(row[key]);
          const recordValue = normalizeValue(currentRecord[key]);
          if (formValue !== recordValue) {
            // 使用表单的原始值，确保空字符串被正确处理为空字符串而不是null
            updates[key] = row[key] === undefined || row[key] === null ? '' : String(row[key]);
          }
        });
      } else if (isAutoType) {
        // 自动化测试用例(role1-4/web)：处理多语言字段映射
        Object.keys(row).forEach(key => {
          const formValue = normalizeValue(row[key]);
          const recordValue = normalizeValue(currentRecord[key]);
          if (formValue !== recordValue) {
            const updateValue = row[key] === undefined || row[key] === null ? '' : String(row[key]);
            updates[key] = updateValue;
          }
        });
      } else {
        // 整体用例/变更用例：处理多语言字段映射
        Object.keys(row).forEach(key => {
          const formValue = normalizeValue(row[key]);
          const recordValue = normalizeValue(currentRecord[key]);
          if (formValue !== recordValue) {
            // 确保使用表单的实际值，将undefined/null转为空字符串，其他值保持原样
            const updateValue = row[key] === undefined || row[key] === null ? '' : String(row[key]);
            updates[key] = updateValue;
          }
        });
      }

      // 如果有更新内容，发送请求
      if (Object.keys(updates).length > 0) {
        // 执行模式：只保存test_result/bug_id/remark到执行结果表
        if (executionMode && taskUuid) {
          const execUpdates = {
            case_id: currentRecord.case_id,
            test_result: updates.test_result || currentRecord.test_result,
            bug_id: updates.bug_id || currentRecord.bug_id || '',
            remark: updates.remark || currentRecord.remark || '',
          };

          await saveExecutionCaseResults(taskUuid, [execUpdates]);

          // 更新本地executionResults
          setExecutionResults(prev => {
            const newMap = new Map(prev);
            newMap.set(currentRecord.case_id, {
              test_result: execUpdates.test_result,
              bug_id: execUpdates.bug_id,
              remark: execUpdates.remark,
            });
            return newMap;
          });

          // 通知父组件执行结果已变更
          if (onResultsChange) {
            onResultsChange();
          }
        } else {
          // 正常模式：保存到用例表
          await updateCaseAPI(currentRecord.case_id, updates);
        }

        // 立即退出编辑状态，避免状态冲突
        setEditingKey('');

        message.success('保存成功');

        // 更新本地数据而不是重新加载整页
        // 这样可以保留通过插入行添加的新用例（即使它的No.不在当前页范围内）
        setCases(prevCases =>
          prevCases.map(c =>
            c.case_id === currentRecord.case_id
              ? { ...c, ...updates }
              : c
          )
        );
      } else {
        // 没有变更时也要退出编辑状态
        setEditingKey('');
      }

    } catch (error) {
      console.error('Failed to save edit:', error);
      message.error('保存失败');
      // 失败时也要退出编辑状态
      setEditingKey('');
    }
  }, [form, projectId, pagination, fetchCases, caseType, language, cases, apiModule, updateCaseAPI, executionMode, taskUuid, onResultsChange]);

  // 取消编辑
  const cancelEdit = useCallback(() => {
    setEditingKey('');
    form.resetFields();
  }, [form]);

  // 打开多语言编辑对话框
  const openMultiLangModal = useCallback((record, fieldName) => {
    const fieldTitles = {
      screen: '画面',
      function: '功能',
      major_function: '大功能分类',
      middle_function: '中功能分类',
      minor_function: '小功能分类',
      precondition: '前置条件',
      test_steps: '测试步骤',
      expected_result: '期待值',
    };

    setMultiLangData({
      record,
      fieldName,
      title: `编辑${fieldTitles[fieldName]} (${fieldName.split('_').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ')})`,
      cn: record[`${fieldName}_cn`] || '',
      jp: record[`${fieldName}_jp`] || '',
      en: record[`${fieldName}_en`] || '',
    });
    setMultiLangModalVisible(true);
  }, []);

  // 保存多语言编辑
  const handleMultiLangSave = async (data) => {
    console.log('=== handleMultiLangSave 开始 ===');
    console.log('data:', data);
    console.log('multiLangData:', multiLangData);

    try {
      const { record } = multiLangData;
      console.log('record:', record);

      const updates = {
        [`${data.fieldName}_cn`]: data.cn,
        [`${data.fieldName}_jp`]: data.jp,
        [`${data.fieldName}_en`]: data.en,
      };
      console.log('updates:', updates);

      // T44: 已移除批量修改确认逻辑，直接保存当前用例

      // 非大功能/中功能字段，或没有匹配的用例，直接保存
      console.log('执行直接保存逻辑');
      await updateCaseAPI(record.case_id, updates);
      console.log('保存完成');
      message.success('保存成功');

      // 关闭对话框并刷新数据
      console.log('关闭多语言对话框');
      setMultiLangModalVisible(false);
      console.log('开始刷新用例列表');
      await fetchCases(pagination.current);
      console.log('刷新完成');
      console.log('=== handleMultiLangSave 结束 ===');
    } catch (error) {
      console.error('❌ Failed to save multi-lang data:', error);
      message.error('保存失败');
      throw error; // 让Modal显示错误
    }
  };

  // T44: 已移除批量修改确认对话框功能
  /*
  // 处理批量修改确认 - 批量修改所有匹配用例
  const handleBatchConfirmOk = async () => {
    console.log('====== 用户点击了"批量修改" ======');
    const { matchingCases, updates, record } = batchConfirmData;
    
    console.log('📋 批量修改数据:');
    console.log('  - 当前用例 case_id:', record.case_id);
    console.log('  - 匹配用例数量:', matchingCases.length);
    console.log('  - updates 对象:', updates);
    console.log('  - updates 包含的字段:', Object.keys(updates));
    
    try {
      setBatchConfirmVisible(false);
      setLoading(true);
      
      // 批量修改所有匹配的用例
      const updateAPI = apiModule === 'api-cases' ? updateApiCase : (
        (caseType && caseType.startsWith('role')) ? updateAutoCase : updateCase
      );
      console.log('使用的 updateAPI:', updateAPI.name);
      
      console.log('开始更新当前用例:', record.case_id);
      console.log('  应用的 updates:', JSON.stringify(updates, null, 2));
      // 更新当前用例
      await updateAPI(projectId, record.case_id, updates);
      console.log('✅ 当前用例更新完成');
      
      // 批量更新所有匹配的用例
      console.log('开始批量更新其他用例:', matchingCases.length);
      console.log('updates 内容（将应用到所有匹配用例）:', updates);
      for (const matchCase of matchingCases) {
        console.log('更新用例:', matchCase.case_id, '应用 updates:', updates);
        await updateAPI(projectId, matchCase.case_id, updates);
      }
      console.log('批量更新完成，所有用例的 CN/JP/EN 都已更新为新值');
      
      message.success(`成功修改 ${matchingCases.length + 1} 条用例`);
      
      // 刷新数据
      console.log('开始刷新用例列表');
      await fetchCases(pagination.current);
      console.log('刷新完成');
    } catch (error) {
      console.error('❌ 批量修改失败:', error);
      message.error('批量修改失败');
    } finally {
      setLoading(false);
    }
  };

  // 处理批量修改取消 - 仅修改当前用例
  const handleBatchConfirmCancel = async () => {
    console.log('用户点击了"仅修改当前"');
    const { updates, record } = batchConfirmData;
    
    try {
      setBatchConfirmVisible(false);
      setLoading(true);
      
      // 只修改当前用例
      const updateAPI = apiModule === 'api-cases' ? updateApiCase : (
        (caseType && caseType.startsWith('role')) ? updateAutoCase : updateCase
      );
      console.log('使用的 updateAPI:', updateAPI.name);
      console.log('开始更新当前用例:', record.case_id);
      
      await updateAPI(projectId, record.case_id, updates);
      console.log('当前用例更新完成');
      
      message.success('保存成功');
      
      // 刷新数据
      console.log('开始刷新用例列表');
      await fetchCases(pagination.current);
      console.log('刷新完成');
    } catch (error) {
      console.error('❌ 保存失败:', error);
      message.error('保存失败');
    } finally {
      setLoading(false);
    }
  };
  */

  // 基础列定义 - 使用 useMemo 确保引用稳定
  const columns = useMemo(() => {
    // ID列
    const idColumn = {
      title: 'No.',
      dataIndex: 'display_id',
      key: 'display_id',
      width: 80,
      fixed: 'left',
      render: (displayId, record, index) => {
        // 对于api-cases(使用UUID),基于分页计算序号
        if (apiModule === 'api-cases') {
          return (pagination.current - 1) * pagination.pageSize + index + 1;
        }
        // 【用例集模式】每个用例集的No.从1开始，基于当前页的索引
        if (caseGroupFilter) {
          return (pagination.current - 1) * pagination.pageSize + index + 1;
        }
        // 【新需求】拖拽时不改变No.，始终显示原始No.
        // 只有点击"重新排序"后才会更新No.
        // 优先使用display_id，然后id，最后使用索引+1
        return displayId || record.id || (index + 1);
      },
    };

    // 判断是否为自动化用例类型(需要在这里定义，因为caseNumberColumn会用到)
    const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');

    // CaseID列 - 根据需求FR-03.1和FR-03.2，Title栏显示为英文，可直接点击编辑
    // 自动化用例类型(role/web)使用case_num字段，其他类型使用case_number字段
    const caseNumField = isAutoType ? 'case_num' : 'case_number';
    const caseNumberColumn = {
      title: 'CaseID',
      dataIndex: caseNumField,
      key: 'case_number',
      width: 120,
      editable: true,
      render: (text, record) => {
        const isEditing = editingKey === record.case_id;
        const currentValue = isAutoType ? record.case_num : record.case_number;
        if (isEditing) {
          return (
            <Input
              defaultValue={currentValue}
              onBlur={(e) => {
                const newValue = e.target.value;
                if (newValue !== currentValue) {
                  const updateData = isAutoType
                    ? { case_num: newValue }
                    : { case_number: newValue };
                  updateCaseAPI(record.case_id, updateData)
                    .then(() => {
                      message.success('保存成功');
                      // 【新需求】更新本地数据，不调用fetchCases，保持hasEditChanges状态
                      setCases(prevCases =>
                        prevCases.map(c =>
                          c.case_id === record.case_id
                            ? { ...c, [caseNumField]: newValue }
                            : c
                        )
                      );
                    })
                    .catch(error => {
                      console.error('Failed to update case number:', error);
                      message.error('保存失败');
                    });
                }
              }}
              autoFocus
            />
          );
        }
        return (
          <div
            style={{ cursor: 'pointer', minHeight: '22px' }}
            onClick={() => {
              if (editingKey === '') {
                startEdit(record);
              }
            }}
          >
            {currentValue || '-'}
          </div>
        );
      },
    };

    // 根据caseType和language生成字段列
    let dataColumns = [];

    // isAutoType已在上面定义

    if (apiModule === 'api-cases') {
      // 接口测试用例: 使用API专属字段(根据需求FR-02)
      dataColumns = [
        {
          title: 'Screen',
          dataIndex: 'screen',
          key: 'screen',
          width: 150,
          editable: true,
          ellipsis: true,
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'screen',
            title: 'Screen',
            record,
            inputType: 'text',
          }),
        },
        {
          title: 'URL',
          dataIndex: 'url',
          key: 'url',
          width: 250,
          editable: true,
          ellipsis: true,
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'url',
            title: 'URL',
            record,
            inputType: 'text',
          }),
        },
        {
          title: 'Header',
          dataIndex: 'header',
          key: 'header',
          width: 200,
          editable: true,
          ellipsis: true,
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'header',
            title: 'Header',
            record,
            inputType: 'textarea',
          }),
        },
        {
          title: 'Method',
          dataIndex: 'method',
          key: 'method',
          width: 100,
          editable: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing) {
              return (
                <Select
                  defaultValue={text || 'GET'}
                  style={{ width: '100%' }}
                  onChange={(value) => {
                    if (value !== text) {
                      updateCaseAPI(record.case_id, { method: value })
                        .then(() => {
                          message.success('保存成功');
                          setCases(prevCases =>
                            prevCases.map(c =>
                              c.case_id === record.case_id
                                ? { ...c, method: value }
                                : c
                            )
                          );
                        })
                        .catch(error => {
                          console.error('Failed to update method:', error);
                          message.error('保存失败');
                        });
                    }
                  }}
                >
                  <Select.Option value="GET">GET</Select.Option>
                  <Select.Option value="POST">POST</Select.Option>
                  <Select.Option value="PUT">PUT</Select.Option>
                  <Select.Option value="DELETE">DELETE</Select.Option>
                  <Select.Option value="PATCH">PATCH</Select.Option>
                </Select>
              );
            }
            return (
              <div
                style={{ cursor: 'pointer', minHeight: '22px' }}
                onClick={() => {
                  if (editingKey === '') {
                    startEdit(record);
                  }
                }}
              >
                {text || 'GET'}
              </div>
            );
          },
        },
        {
          title: 'Body',
          dataIndex: 'body',
          key: 'body',
          width: 200,
          editable: true,
          ellipsis: true,
          render: (text, record) => {
            // 在显示时对密码进行脱敏
            const maskedText = maskKnownPasswords(text || '', knownPasswords);
            return (
              <Tooltip title={maskedText} placement="topLeft">
                <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{maskedText || '-'}</div>
              </Tooltip>
            );
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'body',
            title: 'Body',
            record,
            inputType: 'textarea',
          }),
        },
        {
          title: 'Response',
          dataIndex: 'response',
          key: 'response',
          width: 200,
          editable: true,
          ellipsis: true,
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'response',
            title: 'Response',
            record,
            inputType: 'textarea',
          }),
        },
      ];
    } else if (caseType === 'ai') {
      // AI用例: 使用单语言字段
      dataColumns = [
        {
          title: 'Maj.Category',
          dataIndex: 'major_function',
          key: 'major_function',
          width: 150,
          editable: true,
          render: (text, record, index) => {
            // 编辑时让EditableCell接管
            if (editingKey === record.case_id) {
              return undefined;
            }

            // 【新增】非编辑状态：应用颜色样式
            const previousCase = index > 0 ? cases[index - 1] : null;
            const style = getCellStyle(record, previousCase, 'major_function');
            return <span style={style}>{text || '-'}</span>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'major_function',
            title: 'Maj.Category',
            record,
            inputType: 'text',
          }),
        },
        {
          title: 'Mid.Category',
          dataIndex: 'middle_function',
          key: 'middle_function',
          width: 150,
          editable: true,
          render: (text, record, index) => {
            // 编辑时让EditableCell接管
            if (editingKey === record.case_id) {
              return undefined;
            }

            // 【新增】非编辑状态：应用颜色样式
            const previousCase = index > 0 ? cases[index - 1] : null;
            const style = getCellStyle(record, previousCase, 'middle_function');
            return <span style={style}>{text || '-'}</span>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'middle_function',
            title: 'Mid.Category',
            record,
            inputType: 'text',
          }),
        },
        {
          title: 'Min.Category',
          dataIndex: 'minor_function',
          key: 'minor_function',
          width: 150,
          editable: true,
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'minor_function',
            title: 'Min.Category',
            record,
            inputType: 'text',
          }),
        },
        {
          title: 'Precondition',
          dataIndex: 'precondition',
          key: 'precondition',
          width: 180,
          editable: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing) {
              return undefined;
            }
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{text || '-'}</div>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'precondition',
            title: 'Precondition',
            record,
            inputType: 'textarea',
          }),
        },
        {
          title: 'Test Step',
          dataIndex: 'test_steps',
          key: 'test_steps',
          width: 200,
          editable: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing) {
              return undefined;
            }
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{text || '-'}</div>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'test_steps',
            title: 'Test Step',
            record,
            inputType: 'textarea',
          }),
        },
        {
          title: 'Expect',
          dataIndex: 'expected_result',
          key: 'expected_result',
          width: 180,
          editable: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing) {
              return undefined;
            }
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{text || '-'}</div>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id,
            dataIndex: 'expected_result',
            title: 'Expect',
            record,
            inputType: 'textarea',
          }),
        },
      ];
    } else if (isAutoType) {
      // 自动化测试用例(role1-4/web): 使用多语言字段,简化为screen+function(无大中小功能)
      const langSuffix = language === '中文' ? 'CN' : language === 'English' ? 'EN' : 'JP';
      const langFieldSuffix = language === '中文' ? '_cn' : language === 'English' ? '_en' : '_jp';
      const isChinese = language === '中文';

      dataColumns = [
        {
          title: `Screen${langSuffix}`,
          dataIndex: `screen${langFieldSuffix}`,
          key: `screen${langFieldSuffix}`,
          width: 150,
          editable: true,
          ellipsis: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            const fieldValue = record[`screen${langFieldSuffix}`];
            if (isChinese) {
              return (
                <div
                  style={{ cursor: 'pointer', color: '#1890ff' }}
                  onClick={() => !isEditing && openMultiLangModal(record, 'screen')}
                >
                  {fieldValue || '-'}
                </div>
              );
            }
            return fieldValue || '-';
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `screen${langFieldSuffix}`,
            title: `Screen${langSuffix}`,
            record,
            inputType: 'text',
          }),
        },
        {
          title: `Function${langSuffix}`,
          dataIndex: `function${langFieldSuffix}`,
          key: `function${langFieldSuffix}`,
          width: 180,
          editable: true,
          ellipsis: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            const fieldValue = record[`function${langFieldSuffix}`];
            if (isChinese) {
              return (
                <div
                  style={{ cursor: 'pointer', color: '#1890ff' }}
                  onClick={() => !isEditing && openMultiLangModal(record, 'function')}
                >
                  {fieldValue || '-'}
                </div>
              );
            }
            return fieldValue || '-';
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `function${langFieldSuffix}`,
            title: `Function${langSuffix}`,
            record,
            inputType: 'text',
          }),
        },
        {
          title: `Precondition${langSuffix}`,
          dataIndex: `precondition${langFieldSuffix}`,
          key: `precondition${langFieldSuffix}`,
          width: 180,
          editable: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            const fieldValue = record[`precondition${langFieldSuffix}`];
            const maskedValue = maskKnownPasswords(fieldValue || '', knownPasswords);
            if (isChinese) {
              return (
                <div
                  style={{ cursor: 'pointer', color: '#1890ff', whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}
                  onClick={() => !isEditing && openMultiLangModal(record, 'precondition')}
                >
                  {maskedValue || '-'}
                </div>
              );
            }
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{maskedValue || '-'}</div>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `precondition${langFieldSuffix}`,
            title: `Precondition${langSuffix}`,
            record,
            inputType: 'textarea',
          }),
        },
        {
          title: `Test Step${langSuffix}`,
          dataIndex: `test_steps${langFieldSuffix}`,
          key: `test_steps${langFieldSuffix}`,
          width: 200,
          editable: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            const fieldValue = record[`test_steps${langFieldSuffix}`];
            const maskedValue = maskKnownPasswords(fieldValue || '', knownPasswords);
            if (isChinese) {
              return (
                <div
                  style={{ cursor: 'pointer', color: '#1890ff', whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}
                  onClick={() => !isEditing && openMultiLangModal(record, 'test_steps')}
                >
                  {maskedValue || '-'}
                </div>
              );
            }
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{maskedValue || '-'}</div>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `test_steps${langFieldSuffix}`,
            title: `Test Step${langSuffix}`,
            record,
            inputType: 'textarea',
          }),
        },
        {
          title: `Expect${langSuffix}`,
          dataIndex: `expected_result${langFieldSuffix}`,
          key: `expected_result${langFieldSuffix}`,
          width: 180,
          editable: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            const fieldValue = record[`expected_result${langFieldSuffix}`];
            const maskedValue = maskKnownPasswords(fieldValue || '', knownPasswords);
            if (isChinese) {
              return (
                <div
                  style={{ cursor: 'pointer', color: '#1890ff', whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}
                  onClick={() => !isEditing && openMultiLangModal(record, 'expected_result')}
                >
                  {maskedValue || '-'}
                </div>
              );
            }
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{maskedValue || '-'}</div>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `expected_result${langFieldSuffix}`,
            title: `Expect${langSuffix}`,
            record,
            inputType: 'textarea',
          }),
        },
      ];
    } else {
      // 整体用例/变更用例: 使用多语言字段
      const langSuffix = language === '中文' ? 'CN' : language === 'English' ? 'EN' : 'JP';
      const langFieldSuffix = language === '中文' ? '_cn' : language === 'English' ? '_en' : '_jp';

      // 中文模式: 点击字段弹出多语言对话框
      const isChinese = language === '中文';

      dataColumns = [
        {
          title: `Maj.Category${langSuffix}`,
          dataIndex: `major_function${langFieldSuffix}`,
          key: `major_function${langFieldSuffix}`,
          width: 150,
          editable: true,
          ellipsis: true,
          render: (text, record, index) => {
            // 如果正在编辑且非中文模式，返回undefined让EditableCell接管
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            // 直接从record中获取值
            const fieldValue = record[`major_function${langFieldSuffix}`];

            // 获取重复字段的颜色样式
            const previousCase = index > 0 ? cases[index - 1] : null;
            const repeatStyle = getCellStyle(record, previousCase, `major_function${langFieldSuffix}`);

            // 中文模式: 可点击打开多语言对话框
            if (isChinese) {
              const baseStyle = { cursor: 'pointer', color: '#1890ff' };
              const combinedStyle = repeatStyle.color ? { ...baseStyle, color: repeatStyle.color } : baseStyle;
              return (
                <div
                  style={combinedStyle}
                  onClick={() => !isEditing && openMultiLangModal(record, 'major_function')}
                >
                  {fieldValue || '-'}
                </div>
              );
            }
            // 非中文模式：应用颜色样式
            return <span style={repeatStyle}>{fieldValue || '-'}</span>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `major_function${langFieldSuffix}`,
            title: `Maj.Category${langSuffix}`,
            record,
            inputType: 'text',
          }),
        },
        {
          title: `Mid.Category${langSuffix}`,
          dataIndex: `middle_function${langFieldSuffix}`,
          key: `middle_function${langFieldSuffix}`,
          width: 150,
          editable: true,
          ellipsis: true,
          render: (text, record, index) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            const fieldValue = record[`middle_function${langFieldSuffix}`];

            // 获取重复字段的颜色样式
            const previousCase = index > 0 ? cases[index - 1] : null;
            const repeatStyle = getCellStyle(record, previousCase, `middle_function${langFieldSuffix}`);

            if (isChinese) {
              const baseStyle = { cursor: 'pointer', color: '#1890ff' };
              const combinedStyle = repeatStyle.color ? { ...baseStyle, color: repeatStyle.color } : baseStyle;
              return (
                <div
                  style={combinedStyle}
                  onClick={() => !isEditing && openMultiLangModal(record, 'middle_function')}
                >
                  {fieldValue || '-'}
                </div>
              );
            }
            // 非中文模式：应用颜色样式
            return <span style={repeatStyle}>{fieldValue || '-'}</span>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `middle_function${langFieldSuffix}`,
            title: `Mid.Category${langSuffix}`,
            record,
            inputType: 'text',
          }),
        },
        {
          title: `Min.Category${langSuffix}`,
          dataIndex: `minor_function${langFieldSuffix}`,
          key: `minor_function${langFieldSuffix}`,
          width: 150,
          editable: true,
          ellipsis: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            const fieldValue = record[`minor_function${langFieldSuffix}`];
            if (isChinese) {
              return (
                <div
                  style={{ cursor: 'pointer', color: '#1890ff' }}
                  onClick={() => !isEditing && openMultiLangModal(record, 'minor_function')}
                >
                  {fieldValue || '-'}
                </div>
              );
            }
            return fieldValue || '-';
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `minor_function${langFieldSuffix}`,
            title: `Min.Category${langSuffix}`,
            record,
            inputType: 'text',
          }),
        },
        {
          title: `Precondition${langSuffix}`,
          dataIndex: `precondition${langFieldSuffix}`,
          key: `precondition${langFieldSuffix}`,
          width: 180,
          editable: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            const fieldValue = record[`precondition${langFieldSuffix}`];
            const maskedValue = maskKnownPasswords(fieldValue || '', knownPasswords);
            if (isChinese) {
              return (
                <div
                  style={{ cursor: 'pointer', color: '#1890ff', whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}
                  onClick={() => !isEditing && openMultiLangModal(record, 'precondition')}
                >
                  {maskedValue || '-'}
                </div>
              );
            }
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{maskedValue || '-'}</div>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `precondition${langFieldSuffix}`,
            title: `Precondition${langSuffix}`,
            record,
            inputType: 'textarea',
          }),
        },
        {
          title: `Test Step${langSuffix}`,
          dataIndex: `test_steps${langFieldSuffix}`,
          key: `test_steps${langFieldSuffix}`,
          width: 200,
          editable: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            const fieldValue = record[`test_steps${langFieldSuffix}`];
            const maskedValue = maskKnownPasswords(fieldValue || '', knownPasswords);
            if (isChinese) {
              return (
                <div
                  style={{ cursor: 'pointer', color: '#1890ff', whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}
                  onClick={() => !isEditing && openMultiLangModal(record, 'test_steps')}
                >
                  {maskedValue || '-'}
                </div>
              );
            }
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{maskedValue || '-'}</div>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `test_steps${langFieldSuffix}`,
            title: `Test Step${langSuffix}`,
            record,
            inputType: 'textarea',
          }),
        },
        {
          title: `Expect${langSuffix}`,
          dataIndex: `expected_result${langFieldSuffix}`,
          key: `expected_result${langFieldSuffix}`,
          width: 180,
          editable: true,
          render: (text, record) => {
            const isEditing = editingKey === record.case_id;
            if (isEditing && !isChinese) {
              return undefined;
            }

            const fieldValue = record[`expected_result${langFieldSuffix}`];
            const maskedValue = maskKnownPasswords(fieldValue || '', knownPasswords);
            if (isChinese) {
              return (
                <div
                  style={{ cursor: 'pointer', color: '#1890ff', whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}
                  onClick={() => !isEditing && openMultiLangModal(record, 'expected_result')}
                >
                  {maskedValue || '-'}
                </div>
              );
            }
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{maskedValue || '-'}</div>;
          },
          onCell: (record) => ({
            editing: editingKey === record.case_id && !isChinese,
            dataIndex: `expected_result${langFieldSuffix}`,
            title: `Expect${langSuffix}`,
            record,
            inputType: 'textarea',
          }),
        },
      ];
    }

    // BugID列 (执行模式专用) - 用于记录执行过程中发现的Bug
    const bugIdColumn = {
      title: 'BugID',
      dataIndex: 'bug_id',
      key: 'bug_id',
      width: 120,
      editable: true,
      render: (text, record) => {
        const isEditing = editingKey === record.case_id;
        if (isEditing) {
          return (
            <Input
              defaultValue={text}
              onBlur={(e) => {
                const newValue = e.target.value;
                if (newValue !== text) {
                  // 执行模式：保存到执行结果表
                  if (executionMode && taskUuid) {
                    const execUpdate = {
                      case_id: record.case_id,
                      test_result: record.test_result || 'NR',
                      bug_id: newValue || '',
                      remark: record.remark || '',
                    };
                    saveExecutionCaseResults(taskUuid, [execUpdate])
                      .then(() => {
                        message.success('保存成功');
                        setCases(prevCases =>
                          prevCases.map(c =>
                            c.case_id === record.case_id
                              ? { ...c, bug_id: newValue }
                              : c
                          )
                        );
                        if (onResultsChange) {
                          onResultsChange();
                        }
                      })
                      .catch(error => {
                        console.error('Failed to update bug_id:', error);
                        message.error('保存失败');
                      });
                  }
                }
              }}
              autoFocus
            />
          );
        }
        return (
          <div
            style={{ cursor: 'pointer', minHeight: '22px' }}
            onClick={() => {
              if (editingKey === '') {
                startEdit(record);
              }
            }}
          >
            {text || '-'}
          </div>
        );
      },
    };

    // Remark列 (所有类型都有) - 根据需求FR-03.1和FR-03.2，Title栏显示为英文，可直接点击编辑
    const remarkColumn = {
      title: 'Remark',
      dataIndex: 'remark',
      key: 'remark',
      width: 150,
      editable: true,
      ellipsis: true,
      render: (text, record) => {
        const isEditing = editingKey === record.case_id;
        if (isEditing) {
          return (
            <TextArea
              defaultValue={text}
              autoSize={{ minRows: 2, maxRows: 6 }}
              onBlur={(e) => {
                const newValue = e.target.value;
                if (newValue !== text) {
                  // 执行模式：保存到执行结果表
                  if (executionMode && taskUuid) {
                    const execUpdate = {
                      case_id: record.case_id,
                      test_result: record.test_result || 'NR',
                      bug_id: record.bug_id || '',
                      remark: newValue || '',
                    };
                    saveExecutionCaseResults(taskUuid, [execUpdate])
                      .then(() => {
                        message.success('保存成功');
                        setCases(prevCases =>
                          prevCases.map(c =>
                            c.case_id === record.case_id
                              ? { ...c, remark: newValue }
                              : c
                          )
                        );
                        if (onResultsChange) {
                          onResultsChange();
                        }
                      })
                      .catch(error => {
                        console.error('Failed to update remark:', error);
                        message.error('保存失败');
                      });
                  } else {
                    // 正常模式：保存到用例表
                    updateCaseAPI(record.case_id, { remark: newValue })
                      .then(() => {
                        message.success('保存成功');
                        setCases(prevCases =>
                          prevCases.map(c =>
                            c.case_id === record.case_id
                              ? { ...c, remark: newValue }
                              : c
                          )
                        );
                      })
                      .catch(error => {
                        console.error('Failed to update remark:', error);
                        message.error('保存失败');
                      });
                  }
                }
              }}
              autoFocus
            />
          );
        }
        return (
          <div
            style={{ cursor: 'pointer', minHeight: '22px' }}
            onClick={() => {
              if (editingKey === '') {
                startEdit(record);
              }
            }}
          >
            {text || '-'}
          </div>
        );
      },
    };

    // 操作列 - 根据需求FR-03.1，Title栏显示为英文
    const actionColumn = {
      title: 'Operation',
      key: 'action',
      width: 160, // 缩短列宽以适应图标按钮
      fixed: 'right',
      render: (_, record) => {
        const isCurrentEditing = editingKey === record.case_id;

        // 编辑状态: 显示保存/取消按钮 - 根据需求FR-03.1使用英文
        if (isCurrentEditing) {
          return (
            <Space>
              <Button
                type="link"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  saveEdit(record);
                }}
                size="small"
              >
                Save
              </Button>
              <Button
                type="link"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  cancelEdit();
                }}
                size="small"
              >
                Cancel
              </Button>
            </Space>
          );
        }

        // 中文模式 + 整体/变更/自动化用例: 显示插入+删除按钮
        // api-cases模式: 显示英文的Above/Below按钮
        // web模式: 与api-cases类似的按钮布局，包含详情按钮
        const isMultiLangType = caseType && (caseType.startsWith('role') || caseType === 'overall' || caseType === 'change' || caseType === 'acceptance');
        const isApiCasesMode = apiModule === 'api-cases';
        const isWebCaseMode = caseType === 'web';

        if (isApiCasesMode) {
          // API测试用例: 使用图标按钮，参考AIWeb用例库
          const isEditing = record.case_id === editingKey;
          return (
            <Space size="small">
              {isEditing ? (
                <>
                  <Tooltip title="保存">
                    <Button
                      type="link"
                      size="small"
                      icon={<SaveOutlined />}
                      onClick={() => saveEdit(record)}
                      disabled={loading}
                      style={{ padding: '4px 4px', minWidth: 'auto' }}
                    />
                  </Tooltip>
                  <Tooltip title="取消">
                    <Button
                      type="link"
                      size="small"
                      icon={<DeleteOutlined />}
                      onClick={cancelEdit}
                      disabled={loading}
                      style={{ padding: '4px 4px', minWidth: 'auto' }}
                    />
                  </Tooltip>
                </>
              ) : (
                <>
                  <Tooltip title="详情">
                    <Button
                      type="link"
                      size="small"
                      icon={<EyeOutlined />}
                      onClick={() => {
                        setApiDetailCaseData({ ...record, no: record.display_order || record.no });
                        setApiDetailModalVisible(true);
                      }}
                      disabled={loading}
                      style={{ padding: '4px 4px', minWidth: 'auto' }}
                    />
                  </Tooltip>
                  <Tooltip title="在上方插入">
                    <Button
                      type="link"
                      size="small"
                      icon={<ArrowUpOutlined />}
                      onClick={() => handleInsertAbove(record.case_id)}
                      disabled={editingKey !== '' || loading}
                      style={{ padding: '4px 4px', minWidth: 'auto' }}
                    />
                  </Tooltip>
                  <Tooltip title="在下方插入">
                    <Button
                      type="link"
                      size="small"
                      icon={<ArrowDownOutlined />}
                      onClick={() => handleInsertBelow(record.case_id)}
                      disabled={editingKey !== '' || loading}
                      style={{ padding: '4px 4px', minWidth: 'auto' }}
                    />
                  </Tooltip>
                  <Tooltip title="删除">
                    <button
                      ref={el => deleteButtonRefs.current[`delete-${record.case_id}`] = el}
                      className="ant-btn ant-btn-link ant-btn-dangerous ant-btn-sm"
                      onClick={() => handleDelete(record)}
                      disabled={editingKey !== '' || loading}
                      style={{ color: '#ff4d4f', padding: '4px 4px', minWidth: 'auto' }}
                    >
                      <DeleteOutlined />
                    </button>
                  </Tooltip>
                </>
              )}
            </Space>
          );
        } else if (isWebCaseMode) {
          // Web用例模式: 与API用例类似的按钮布局，包含详情按钮
          const isEditing = record.case_id === editingKey;
          return (
            <Space size="small">
              {isEditing ? (
                <>
                  <Tooltip title="保存">
                    <Button
                      type="link"
                      size="small"
                      icon={<SaveOutlined />}
                      onClick={() => saveEdit(record)}
                      disabled={loading}
                      style={{ padding: '4px 4px', minWidth: 'auto' }}
                    />
                  </Tooltip>
                  <Tooltip title="取消">
                    <Button
                      type="link"
                      size="small"
                      icon={<DeleteOutlined />}
                      onClick={cancelEdit}
                      disabled={loading}
                      style={{ padding: '4px 4px', minWidth: 'auto' }}
                    />
                  </Tooltip>
                </>
              ) : (
                <>
                  <Tooltip title="详情">
                    <Button
                      type="link"
                      size="small"
                      icon={<EyeOutlined />}
                      onClick={() => {
                        setWebDetailCaseData({ ...record, no: record.display_order || record.no || record.id });
                        setWebDetailModalVisible(true);
                      }}
                      disabled={loading}
                      style={{ padding: '4px 4px', minWidth: 'auto' }}
                    />
                  </Tooltip>
                  <Tooltip title="在上方插入">
                    <Button
                      type="link"
                      size="small"
                      icon={<ArrowUpOutlined />}
                      onClick={() => handleInsertAbove(record.case_id)}
                      disabled={editingKey !== '' || loading}
                      style={{ padding: '4px 4px', minWidth: 'auto' }}
                    />
                  </Tooltip>
                  <Tooltip title="在下方插入">
                    <Button
                      type="link"
                      size="small"
                      icon={<ArrowDownOutlined />}
                      onClick={() => handleInsertBelow(record.case_id)}
                      disabled={editingKey !== '' || loading}
                      style={{ padding: '4px 4px', minWidth: 'auto' }}
                    />
                  </Tooltip>
                  <Tooltip title="删除">
                    <button
                      ref={el => deleteButtonRefs.current[`delete-${record.case_id}`] = el}
                      className="ant-btn ant-btn-link ant-btn-dangerous ant-btn-sm"
                      onClick={() => handleDelete(record)}
                      disabled={editingKey !== '' || loading}
                      style={{ color: '#ff4d4f', padding: '4px 4px', minWidth: 'auto' }}
                    >
                      <DeleteOutlined />
                    </button>
                  </Tooltip>
                </>
              )}
            </Space>
          );
        } else if (isMultiLangType && language === '中文') {
          // CN模式：只显示图标，不显示文字
          return (
            <Space size="small">
              <Tooltip title={t('manualTest.insertAbove')}>
                <Button
                  type="link"
                  size="small"
                  icon={<ArrowUpOutlined />}
                  onClick={() => handleInsertAbove(record.case_id)}
                  disabled={editingKey !== '' || loading}
                  style={{ padding: '4px 4px', minWidth: 'auto' }}
                />
              </Tooltip>
              <Tooltip title={t('manualTest.insertBelow')}>
                <Button
                  type="link"
                  size="small"
                  icon={<ArrowDownOutlined />}
                  onClick={() => handleInsertBelow(record.case_id)}
                  disabled={editingKey !== '' || loading}
                  style={{ padding: '4px 4px', minWidth: 'auto' }}
                />
              </Tooltip>
              <Tooltip title={hasEditChanges ? '有未保存的插入操作，请先保存' : '删除'}>
                <button
                  ref={(el) => {
                    if (el) {
                      deleteButtonRefs.current[`delete-${record.case_id}`] = el;
                    }
                  }}
                  className="ant-btn ant-btn-link ant-btn-dangerous ant-btn-sm"
                  onClick={() => handleDelete(record)}
                  disabled={editingKey !== '' || hasEditChanges}
                  style={{ color: '#ff4d4f', padding: '4px 4px', minWidth: 'auto' }}
                  title={hasEditChanges ? '有未保存的插入操作，请先保存' : ''}
                >
                  <DeleteOutlined />
                </button>
              </Tooltip>
            </Space>
          );
        }

        // T44: 操作列按钮图标化
        return (
          <Space size="small">
            <Tooltip title="在上方插入">
              <Button
                type="link"
                size="small"
                icon={<ArrowUpOutlined />}
                onClick={() => handleInsertAbove(record.case_id)}
                disabled={editingKey !== '' || loading}
                style={{ padding: '4px 4px', minWidth: 'auto' }}
              />
            </Tooltip>
            <Tooltip title="在下方插入">
              <Button
                type="link"
                size="small"
                icon={<ArrowDownOutlined />}
                onClick={() => handleInsertBelow(record.case_id)}
                disabled={editingKey !== '' || loading}
                style={{ padding: '4px 4px', minWidth: 'auto' }}
              />
            </Tooltip>
            <Tooltip title={hasEditChanges ? '有未保存的插入操作，请先保存' : '编辑'}>
              <Button
                type="link"
                onClick={() => startEdit(record)}
                size="small"
                icon={<SaveOutlined />}
                disabled={editingKey !== '' || hasEditChanges}
              />
            </Tooltip>
            <Tooltip title={hasEditChanges ? '有未保存的插入操作，请先保存' : '删除'}>
              <button
                ref={(el) => {
                  if (el) {
                    deleteButtonRefs.current[`delete-${record.case_id}`] = el;
                  }
                }}
                className="ant-btn ant-btn-link ant-btn-dangerous ant-btn-sm"
                onClick={() => handleDelete(record)}
                disabled={editingKey !== '' || hasEditChanges}
                style={{ color: '#ff4d4f', padding: '0 8px' }}
              >
                <DeleteOutlined />
              </button>
            </Tooltip>
          </Space>
        );
      },
    };

    // 组装所有列
    // 整体用例/受入用例/变更用例不显示Remark列
    // role1-4类型(AI自动化和AI接口测试)不显示TestResult和Remark列
    // Web用例也不显示Remark列
    // API用例（api-cases）也不显示TestResult和Remark列
    const shouldHideRemark = caseType === 'overall' || caseType === 'acceptance' || caseType === 'change' ||
      caseType === 'role1' || caseType === 'role2' || caseType === 'role3' || caseType === 'role4' ||
      caseType === 'web' || apiModule === 'api-cases';

    // 执行模式：构建执行专用列（BugID + Remark）
    const executionColumns = executionMode ? [bugIdColumn, remarkColumn] : [];

    // api-cases模式: 不包含CaseID列(case_num/case_number字段)
    let allColumns;
    if (executionMode) {
      // 执行模式：显示数据列 + BugID + Remark + 操作列
      allColumns = apiModule === 'api-cases'
        ? [idColumn, ...dataColumns, ...executionColumns, actionColumn]
        : [idColumn, caseNumberColumn, ...dataColumns, ...executionColumns, actionColumn];
    } else {
      // 正常模式：原有逻辑
      allColumns = apiModule === 'api-cases'
        ? [idColumn, ...dataColumns, ...(shouldHideRemark ? [] : [remarkColumn]), actionColumn]
        : [idColumn, caseNumberColumn, ...dataColumns, ...(shouldHideRemark ? [] : [remarkColumn]), actionColumn];
    }

    // 根据用例类型过滤列 - AI用例和role1-4不显示"测试结果"列
    const shouldHideTestResult = caseType === 'ai' || caseType === 'role1' || caseType === 'role2' ||
      caseType === 'role3' || caseType === 'role4';
    return shouldHideTestResult
      ? allColumns.filter(col => col.key !== 'test_result')
      : allColumns;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [caseType, editingKey, language, pagination, openMultiLangModal, startEdit, saveEdit, cancelEdit, handleDelete, apiModule, executionMode, taskUuid, onResultsChange]);

  // 复选框配置
  const rowSelection = {
    selectedRowKeys,
    onChange: (newSelectedRowKeys) => {
      setSelectedRowKeys(newSelectedRowKeys);
    },
    getCheckboxProps: (record) => ({
      disabled: editingKey !== '' || hasEditChanges, // 有行在编辑或有未保存插入时禁用
    }),
  };

  // 分页变更
  const handleTableChange = (newPagination) => {
    // 如果有未保存的编辑变更，阻止分页切换
    if (hasEditChanges) {
      return;
    }

    // 如果页面大小改变，重置到第一页
    if (newPagination.pageSize !== pagination.pageSize) {
      setPagination({
        ...pagination,
        current: 1,
        pageSize: newPagination.pageSize,
      });
      fetchCases(1, newPagination.pageSize);
    } else {
      fetchCases(newPagination.current);
    }
  };

  return (
    <div className="editable-table-container">
      <div className="table-toolbar">
        <Space wrap>
          {/* AI用例专属按钮 */}
          {caseType === 'ai' && (
            <>
              <Button
                icon={<DownloadOutlined />}
                onClick={handleExportAICases}
                disabled={editingKey !== '' || cases.length === 0}
              >
                导出用例
              </Button>
              <Button
                danger
                icon={<DeleteOutlined />}
                onClick={handleClearAICases}
                disabled={editingKey !== '' || cases.length === 0}
              >
                {t('manualTest.clearAICases')}
              </Button>
            </>
          )}

          {/* 整体用例、变更用例和受入用例共享按钮 */}
          {(caseType === 'overall' || caseType === 'change' || caseType === 'acceptance') && (
            <>
              {!hiddenButtons.includes('aiSupplement') && (
                <Button
                  icon={<CopyOutlined />}
                  onClick={handleFillClassification}
                  disabled={editingKey !== '' || loading || cases.length <= 1}
                  loading={loading}
                >
                  {t('manualTest.fillClassification')}
                </Button>
              )}
              {!hiddenButtons.includes('exportCases') && (
                <Button
                  icon={<DownloadOutlined />}
                  onClick={handleExportCases}
                  disabled={editingKey !== '' || cases.length === 0}
                >
                  {t('manualTest.exportCases')}
                </Button>
              )}
              {!hiddenButtons.includes('importCases') && (
                <Upload
                  accept=".xlsx,.xls"
                  beforeUpload={handleImportCases}
                  showUploadList={false}
                >
                  <Button
                    icon={<UploadOutlined />}
                    disabled={editingKey !== ''}
                  >
                    {t('manualTest.importCases')}
                  </Button>
                </Upload>
              )}
            </>
          )}

          {/* T44: saveVersion按钮已移除，版本管理功能已移至独立模块 */}
          {/* T44: 批量删除按钮已移除，由ManualCaseManagementTab工具栏提供 */}

          {/* 保存按钮 - 仅在有编辑变更时显示 */}
          {hasEditChanges && (
            <Button
              type="primary"
              icon={<SaveOutlined />}
              loading={loading}
              onClick={async () => {
                try {
                  setLoading(true);

                  // 判断使用哪个API
                  const isAutoType = caseType && (caseType.startsWith('role') || caseType === 'web');
                  let insertAPI, reassignAPI, updateAPI;
                  if (apiModule === 'api-cases') {
                    insertAPI = insertApiCase;
                    reassignAPI = null; // API测试用例暂不支持重新分配ID
                    updateAPI = updateApiCase;
                  } else {
                    insertAPI = isAutoType ? insertAutoCase : insertCase;
                    reassignAPI = isAutoType ? reassignAutoIDs : reassignAllIDs;
                    updateAPI = isAutoType ? updateAutoCase : updateCase;
                  }

                  // 依次执行所有插入操作
                  // 策略：找到第一个真实ID，然后按顺序在其之前/之后插入所有临时行
                  const realCases = cases.filter(c => !c.case_id.startsWith('temp_'));

                  if (realCases.length === 0) {
                    // 如果没有真实ID，说明所有行都是新建的，使用createCase而不是insertCase
                    message.error('没有现有用例可作为参考位置，请先创建至少一个用例');
                    setLoading(false);
                    return;
                  }

                  // 使用第一个真实用例作为参考点
                  let referenceId = realCases[0].case_id;

                  // 用于追踪已插入的用例ID，key是临时ID，value是真实ID
                  const insertedIdMap = {};

                  // 按照临时行在cases数组中的顺序插入
                  for (let i = 0; i < cases.length; i++) {
                    const currentCase = cases[i];
                    if (currentCase.case_id.startsWith('temp_')) {
                      // 这是一个临时行，需要插入
                      // 找到这个临时行之前最近的真实ID（包括刚插入的）
                      let beforeRealId = null;
                      for (let j = i - 1; j >= 0; j--) {
                        const prevCase = cases[j];
                        if (!prevCase.case_id.startsWith('temp_')) {
                          // 原本就是真实ID
                          beforeRealId = prevCase.case_id;
                          break;
                        } else if (insertedIdMap[prevCase.case_id]) {
                          // 之前插入的临时行，现在有了真实ID
                          beforeRealId = insertedIdMap[prevCase.case_id];
                          break;
                        }
                      }

                      let insertResult;
                      if (beforeRealId) {
                        // 在找到的真实ID之后插入
                        insertResult = await insertAPI(projectId, {
                          caseType,
                          position: 'after',
                          targetCaseId: beforeRealId,
                          language: !isAutoType ? language : undefined,
                          caseGroup: currentCase.case_group || undefined, // 使用临时行的case_group字段
                        });
                      } else {
                        // 没有找到之前的真实ID，在第一个真实ID之前插入
                        insertResult = await insertAPI(projectId, {
                          caseType,
                          position: 'before',
                          targetCaseId: referenceId,
                          language: !isAutoType ? language : undefined,
                          caseGroup: currentCase.case_group || undefined, // 使用临时行的case_group字段
                        });
                      }

                      // 记录插入结果，用于后续临时行的参考
                      // 注意：auto-cases和api-cases返回case_id(UUID)，manual-cases返回id
                      if (insertResult && insertResult.case_id) {
                        insertedIdMap[currentCase.case_id] = insertResult.case_id;
                      } else if (insertResult && insertResult.id) {
                        insertedIdMap[currentCase.case_id] = insertResult.id;
                      }
                    }
                  }

                  // 调用API重新分配所有ID（api-cases不需要此操作）
                  if (reassignAPI) {
                    await reassignAPI(projectId, caseType);
                  }

                  console.log('[保存] 开始刷新数据，当前页:', pagination.current);

                  // 刷新当前页（先刷新再清空状态，确保能看到新数据）
                  const currentPage = pagination.current || 1;
                  await fetchCases(currentPage);

                  // 清空待插入列表
                  setPendingInserts([]);
                  setHasEditChanges(false);

                  console.log('[保存] 刷新完成，用例数量:', cases.length);
                  message.success('保存成功');
                } catch (error) {
                  console.error('保存失败:', error);
                  message.error('保存失败: ' + (error?.message || ''));
                } finally {
                  setLoading(false);
                }
              }}
            >
              保存
            </Button>
          )}
        </Space>
      </div>

      {/* 插入提示 */}
      {hasEditChanges && (
        <Alert
          message={`当前页已插入 ${pendingInserts.length} 个空行。点击保存按钮后将实际插入到数据库，并重新分配所有用例编号`}
          type="warning"
          showIcon
          closable
          style={{ marginBottom: 16 }}
        />
      )}

      <Form form={form} component={false}>
        <Table
          rowSelection={rowSelection}
          components={{
            body: {
              cell: EditableCell,
            },
          }}
          columns={columns}
          dataSource={cases}
          rowKey="case_id"
          loading={loading}
          rowClassName={(record) => {
            // 如果是高亮行，添加高亮样式类
            return record.case_id === highlightedRowId ? 'highlighted-row' : '';
          }}
          pagination={
            hasEditChanges
              ? {
                ...pagination,
                pageSize: cases.length, // 设置为当前数据长度，确保所有数据都显示
                total: pagination.total + pendingInserts.length, // 总数 = 原总数 + 待插入数
                showTotal: (total) => `${t('common.total')} ${total} ${t('common.items')} (当前页已扩展显示 ${cases.length} 条)`,
                pageSizeOptions: ['10', '20', '50', '100'],
                showSizeChanger: false,  // 禁用页面大小选择器
                simple: false,  // 使用完整分页器显示总数
              }
              : {
                ...pagination,
                showTotal: (total) => `${t('common.total')} ${total} ${t('common.items')}`,
                pageSizeOptions: ['10', '20', '50', '100'],
                showSizeChanger: true,
              }
          }
          onChange={handleTableChange}
          scroll={{
            x: 'max-content',
            // 当数据量<=10时不设置y滚动，让表格自适应高度；超过10条时才启用滚动
            ...(cases.length > 10 ? { y: 'calc(100vh - 400px)' } : {})
          }}
          bordered
          locale={{
            emptyText: (
              <div style={{ padding: '40px 0', textAlign: 'center' }}>
                <p style={{ color: '#999', marginBottom: '16px' }}>暂无数据</p>
                <Button
                  type="primary"
                  icon={<PlusCircleOutlined />}
                  onClick={async () => {
                    try {
                      setLoading(true);
                      const newCase = await createDefaultEmptyRow();
                      message.success('创建成功');
                      await fetchCases(1);
                    } catch (error) {
                      console.error('创建失败:', error);
                      message.error('创建失败: ' + (error.response?.data?.error || error.message));
                    } finally {
                      setLoading(false);
                    }
                  }}
                  disabled={loading}
                >
                  {t('manualTest.addFirstRow')}
                </Button>
              </div>
            ),
          }}
        />
      </Form>

      {/* 多语言编辑对话框 */}
      <MultiLangEditModal
        visible={multiLangModalVisible}
        title={multiLangData.title}
        fieldName={multiLangData.fieldName}
        data={{
          cn: multiLangData.cn,
          jp: multiLangData.jp,
          en: multiLangData.en,
        }}
        onSave={handleMultiLangSave}
        onCancel={() => setMultiLangModalVisible(false)}
      />

      {/* API用例详情弹窗 */}
      {apiModule === 'api-cases' && (
        <ApiCaseDetailModal
          visible={apiDetailModalVisible}
          caseData={apiDetailCaseData}
          projectId={projectId}
          groupId={caseGroupId}
          onSave={async (data) => {
            try {
              await updateApiCase(projectId, data.case_id, {
                screen: data.screen,
                method: data.method,
                url: data.url,
                header: data.header,
                body: data.body,
                response: data.response,
                script_code: data.script_code,
              });
              // 更新本地数据
              setCases(prevCases =>
                prevCases.map(c =>
                  c.case_id === data.case_id
                    ? {
                      ...c,
                      screen: data.screen,
                      method: data.method,
                      url: data.url,
                      header: data.header,
                      body: data.body,
                      response: data.response,
                      script_code: data.script_code,
                    }
                    : c
                )
              );
            } catch (error) {
              console.error('保存失败:', error);
              throw error;
            }
          }}
          onCancel={() => {
            setApiDetailModalVisible(false);
            setApiDetailCaseData(null);
          }}
        />
      )}

      {/* Web用例详情弹窗 */}
      {caseType === 'web' && (
        <WebCaseDetailModal
          visible={webDetailModalVisible}
          caseData={webDetailCaseData}
          language={language === '中文' ? 'cn' : language === '日本語' ? 'jp' : 'en'}
          projectId={projectId}
          groupId={caseGroupId}
          onSave={async (data) => {
            try {
              await updateAutoCase(projectId, data.case_id, data);
              // 更新本地数据
              setCases(prevCases =>
                prevCases.map(c =>
                  c.case_id === data.case_id
                    ? { ...c, ...data }
                    : c
                )
              );
            } catch (error) {
              console.error('保存失败:', error);
              throw error;
            }
          }}
          onCancel={() => {
            setWebDetailModalVisible(false);
            setWebDetailCaseData(null);
          }}
        />
      )}

      {/* T44: 已移除批量修改确认对话框 */}
    </div>
  );
};

export default EditableTable;
