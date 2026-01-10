import React, { useState, useEffect, useCallback, useRef } from 'react';
import { Form, Input, Select, DatePicker, Button, Space, Empty, message, Row, Col, Modal, Table, Radio, Progress, Tooltip, Tag } from 'antd';
import { FileSearchOutlined, DownloadOutlined, SaveOutlined, EditOutlined, EyeOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import PropTypes from 'prop-types';
import dayjs from 'dayjs';
import * as XLSX from 'xlsx';
import { updateExecutionTask } from '../../../api/executionTask';
import { saveExecutionCaseResults, getExecutionCaseResults } from '../../../api/executionCaseResult';
import CaseSelectionPanel from './CaseSelectionPanel';
import CaseDetailModal from './CaseDetailModal';
import './TaskMetadataPanel.css';

const { Option } = Select;
const { TextArea } = Input;

const TaskMetadataPanel = ({ task, projectId, projectName, onSave }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [caseSelectionVisible, setCaseSelectionVisible] = useState(false);
  const [selectedCasesData, setSelectedCasesData] = useState(null);
  const [caseTableData, setCaseTableData] = useState([]);
  const [displayLanguage, setDisplayLanguage] = useState(null); // 显示语言筛选，初始为null以便使用task.display_language作为后备
  
  // 分页状态
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  
  // 用例详细弹窗状态
  const [caseDetailVisible, setCaseDetailVisible] = useState(false);
  const [selectedCaseForDetail, setSelectedCaseForDetail] = useState(null);
  
  // 用于防抖自动保存的ref
  const saveTimeoutRef = useRef(null);
  const pendingSaveRef = useRef(null);
  
  console.log('🟡 [TaskMetadataPanel] Render with projectId:', projectId, 'task:', task?.task_name);

  // 打开用例详细弹窗
  const handleOpenCaseDetail = (record) => {
    setSelectedCaseForDetail(record);
    setCaseDetailVisible(true);
  };

  // 关闭用例详细弹窗
  const handleCloseCaseDetail = () => {
    setCaseDetailVisible(false);
    setSelectedCaseForDetail(null);
  };

  // 从用例详细弹窗保存数据
  const handleSaveCaseDetail = async (data) => {
    // 更新表格数据
    setCaseTableData(prev => prev.map(c => 
      c.case_id === data.case_id 
        ? { ...c, test_result: data.test_result, bug_id: data.bug_id, remark: data.remark } 
        : c
    ));
    
    // 触发自动保存
    if (data.test_result) {
      await autoSaveCaseResult(data.case_id, 'test_result', data.test_result);
    }
    if (data.bug_id !== undefined) {
      await autoSaveCaseResult(data.case_id, 'bug_id', data.bug_id);
    }
    if (data.remark !== undefined) {
      await autoSaveCaseResult(data.case_id, 'remark', data.remark);
    }
    
    message.success(t('testExecution.caseDetail.saveSuccess'));
  };

  // 立即保存待保存的数据（用于任务切换或组件卸载前）
  const flushPendingSave = useCallback(async (taskUuid) => {
    if (!pendingSaveRef.current || Object.keys(pendingSaveRef.current).length === 0) {
      return;
    }
    
    const dataToSave = Object.values(pendingSaveRef.current);
    pendingSaveRef.current = {};
    
    if (dataToSave.length === 0 || !taskUuid) return;
    
    try {
      console.log('💾 [TaskMetadataPanel] Flushing pending save:', dataToSave.length, 'items');
      await saveExecutionCaseResults(taskUuid, dataToSave);
      console.log('✅ [TaskMetadataPanel] Flush save success');
    } catch (error) {
      console.error('❌ [TaskMetadataPanel] Flush save failed:', error);
    }
  }, []);

  // 任务切换时，先保存当前任务的待保存数据，再加载新任务的执行结果
  const prevTaskUuidRef = useRef(null);
  
  useEffect(() => {
    console.log('🔄 [TaskMetadataPanel] useEffect triggered, task_uuid:', task?.task_uuid, 'task_name:', task?.task_name);
    
    // 如果任务切换了，先保存之前任务的待保存数据
    if (prevTaskUuidRef.current && prevTaskUuidRef.current !== task?.task_uuid) {
      console.log('🔄 [TaskMetadataPanel] Task changed, flushing pending save for:', prevTaskUuidRef.current);
      flushPendingSave(prevTaskUuidRef.current);
      // 清除定时器
      if (saveTimeoutRef.current) {
        clearTimeout(saveTimeoutRef.current);
        saveTimeoutRef.current = null;
      }
    }
    
    prevTaskUuidRef.current = task?.task_uuid;
    
    if (task && task.task_uuid) {
      console.log('🔄 [TaskMetadataPanel] Calling loadSavedCaseResults for task:', task.task_name);
      console.log('🔄 [TaskMetadataPanel] task.display_language:', task.display_language);
      
      // 恢复语言设置：优先 localStorage，其次 task.display_language，最后默认 cn
      const savedFilter = localStorage.getItem(`execution_filter_${task.task_uuid}`);
      if (savedFilter) {
        const filterConditions = JSON.parse(savedFilter);
        const lang = filterConditions.language || task.display_language || 'cn';
        console.log('🔄 [TaskMetadataPanel] Restoring language from localStorage:', lang);
        setDisplayLanguage(lang);
      } else if (task.display_language) {
        // localStorage 没有缓存，使用 task 中保存的语言
        console.log('🔄 [TaskMetadataPanel] Restoring language from task.display_language:', task.display_language);
        setDisplayLanguage(task.display_language);
      } else {
        // 都没有，根据执行类型设置默认值
        const defaultLang = task.execution_type === 'manual' ? 'all' : 
                           task.execution_type === 'api' ? 'en' : 'cn';
        console.log('🔄 [TaskMetadataPanel] Setting default language:', defaultLang);
        setDisplayLanguage(defaultLang);
      }
      
      loadSavedCaseResults();
    } else {
      console.log('🔄 [TaskMetadataPanel] No task, clearing data');
      // 清空数据
      setSelectedCasesData(null);
      setCaseTableData([]);
    }
    
    // 组件卸载时保存待保存的数据
    return () => {
      if (task?.task_uuid) {
        flushPendingSave(task.task_uuid);
      }
    };
  }, [task?.task_uuid, flushPendingSave]);

  // 加载已保存的用例执行结果
  const loadSavedCaseResults = async () => {
    console.log('📥 [TaskMetadataPanel] loadSavedCaseResults called');
    console.log('📥 [TaskMetadataPanel] task:', task?.task_name, 'task_uuid:', task?.task_uuid);
    
    if (!task || !task.task_uuid) {
      console.log('📥 [TaskMetadataPanel] No task or task_uuid, skipping load');
      return;
    }
    
    try {
      console.log('📥 [TaskMetadataPanel] Calling getExecutionCaseResults API...');
      const results = await getExecutionCaseResults(task.task_uuid);
      console.log('📥 [TaskMetadataPanel] API returned results:', results);
      console.log('📥 [TaskMetadataPanel] Results length:', results?.length);
      console.log('📥 [TaskMetadataPanel] Results[0]:', results?.[0]);
      
      if (results && results.length > 0) {
        // 从localStorage恢复筛选条件
        const savedFilter = localStorage.getItem(`execution_filter_${task.task_uuid}`);
        console.log('📥 [TaskMetadataPanel] savedFilter from localStorage:', savedFilter);
        // 注意：默认值不设置language，让它回退到task.display_language
        const parsedFilter = savedFilter ? JSON.parse(savedFilter) : { case_type: 'role1' };
        
        // 语言优先级：1. localStorage中保存的语言 2. 任务中保存的语言(display_language) 3. 默认cn
        const taskLang = task.display_language || '';
        const effectiveLanguage = parsedFilter.language || taskLang || 'cn';
        
        // 确保case_group和language优先使用task中保存的值，防止清除缓存后丢失
        const filterConditions = {
          ...parsedFilter,
          case_group: parsedFilter.case_group || task.case_group_name || '',
          language: effectiveLanguage
        };
        console.log('📥 [TaskMetadataPanel] filterConditions with task fallback:', filterConditions);
        console.log('📥 [TaskMetadataPanel] task.display_language:', task.display_language);
        
        // 同步设置displayLanguage状态
        setDisplayLanguage(effectiveLanguage);
        console.log('📥 [TaskMetadataPanel] setDisplayLanguage:', effectiveLanguage);
        
        // 将结果转换为表格数据
        // 使用后端返回的 display_id 作为 No.（已按 display_id 排序）
        const tableData = results.map((r, index) => ({
          ...r,
          key: r.case_id || index,
          no: r.display_id || (index + 1), // 优先使用保存的 display_id
          test_result: r.test_result || 'Block',
          bug_id: r.bug_id || '',
          remark: r.remark || '',
        }));
        
        console.log('📥 [TaskMetadataPanel] Setting selectedCasesData and caseTableData');
        console.log('📥 [TaskMetadataPanel] tableData[0]:', tableData[0]);
        
        setSelectedCasesData({
          cases: results,
          filterConditions: filterConditions,
          total: results.length
        });
        setCaseTableData(tableData);
        console.log('📥 [TaskMetadataPanel] Data loaded successfully!');
      } else {
        console.log('📥 [TaskMetadataPanel] No results found, clearing data');
        setSelectedCasesData(null);
        setCaseTableData([]);
      }
    } catch (error) {
      console.log('📥 [TaskMetadataPanel] Load failed:', error.message);
      console.log('📥 [TaskMetadataPanel] Error details:', error);
      // 没有保存的数据，保持空状态
      setSelectedCasesData(null);
      setCaseTableData([]);
    }
  };

  // 当选中用例数据变化时，初始化表格数据
  useEffect(() => {
    if (selectedCasesData && selectedCasesData.cases) {
      console.log('🔵 [TaskMetadataPanel] Initializing table data');
      console.log('🔵 [TaskMetadataPanel] selectedCasesData.cases[0]:', selectedCasesData.cases[0]);
      console.log('🔵 [TaskMetadataPanel] execution_type:', selectedCasesData.filterConditions?.execution_type);
      
      const tableData = selectedCasesData.cases.map((c, index) => ({
        ...c,
        key: c.case_id || c.id || index,
        no: index + 1,
        test_result: c.test_result || 'Block',
        bug_id: c.bug_id || '',
        remark: c.remark || '',
      }));
      
      console.log('✅ [TaskMetadataPanel] tableData[0]:', tableData[0]);
      console.log('✅ [TaskMetadataPanel] tableData.length:', tableData.length);
      setCaseTableData(tableData);
    }
  }, [selectedCasesData]);

  // 自动保存单条记录（防抖）
  const autoSaveCaseResult = useCallback(async (caseId, field, value) => {
    if (!task || !task.task_uuid) return;
    
    // 从当前表格数据中获取完整的用例信息
    const caseRecord = caseTableData.find(c => c.case_id === caseId);
    if (!caseRecord) {
      console.log('⚠️ [TaskMetadataPanel] Case not found for auto-save:', caseId);
      return;
    }
    
    // 更新待保存数据
    if (!pendingSaveRef.current) {
      pendingSaveRef.current = {};
    }
    if (!pendingSaveRef.current[caseId]) {
      // 获取当前的用例类型
      const currentCaseType = selectedCasesData?.filterConditions?.case_type || 'overall';
      
      // 初始化时复制完整用例数据（包含手工测试、AI Web和API的所有字段）
      // 注意：display_id 使用 no（当前显示序号），不能使用原始用例的 id
      pendingSaveRef.current[caseId] = {
        case_id: caseRecord.case_id,
        display_id: caseRecord.no || caseRecord.display_id || 0, // 使用当前显示序号
        case_num: caseRecord.case_num || caseRecord.case_number || '', // 用户自定义CaseID
        case_type: caseRecord.case_type || currentCaseType, // 用例类型
        test_result: caseRecord.test_result || 'Block',
        bug_id: caseRecord.bug_id || '',
        remark: caseRecord.remark || '',
        // API 用例字段（无多语言）
        screen: caseRecord.screen || '',
        url: caseRecord.url || '',
        header: caseRecord.header || '',
        method: caseRecord.method || '',
        body: caseRecord.body || '',
        response: caseRecord.response || '',
        response_time: caseRecord.response_time || '',
        // AI Web 用例字段
        screen_cn: caseRecord.screen_cn || '',
        screen_jp: caseRecord.screen_jp || '',
        screen_en: caseRecord.screen_en || '',
        function_cn: caseRecord.function_cn || '',
        function_jp: caseRecord.function_jp || '',
        function_en: caseRecord.function_en || '',
        // 手工测试用例字段
        major_function_cn: caseRecord.major_function_cn || caseRecord.major_function || '',
        major_function_jp: caseRecord.major_function_jp || '',
        major_function_en: caseRecord.major_function_en || '',
        middle_function_cn: caseRecord.middle_function_cn || caseRecord.middle_function || '',
        middle_function_jp: caseRecord.middle_function_jp || '',
        middle_function_en: caseRecord.middle_function_en || '',
        minor_function_cn: caseRecord.minor_function_cn || caseRecord.minor_function || '',
        minor_function_jp: caseRecord.minor_function_jp || '',
        minor_function_en: caseRecord.minor_function_en || '',
        // 公共字段
        precondition_cn: caseRecord.precondition_cn || caseRecord.precondition || '',
        precondition_jp: caseRecord.precondition_jp || '',
        precondition_en: caseRecord.precondition_en || '',
        test_steps_cn: caseRecord.test_steps_cn || caseRecord.test_steps || '',
        test_steps_jp: caseRecord.test_steps_jp || '',
        test_steps_en: caseRecord.test_steps_en || '',
        expected_result_cn: caseRecord.expected_result_cn || caseRecord.expected_result || '',
        expected_result_jp: caseRecord.expected_result_jp || '',
        expected_result_en: caseRecord.expected_result_en || '',
      };
    }
    // 更新变更的字段
    pendingSaveRef.current[caseId][field] = value;
    
    // 清除之前的定时器
    if (saveTimeoutRef.current) {
      clearTimeout(saveTimeoutRef.current);
    }
    
    // 设置新的防抖定时器（500ms）
    saveTimeoutRef.current = setTimeout(async () => {
      const dataToSave = Object.values(pendingSaveRef.current);
      pendingSaveRef.current = {};
      
      if (dataToSave.length === 0) return;
      
      try {
        console.log('💾 [TaskMetadataPanel] Auto-saving:', dataToSave);
        await saveExecutionCaseResults(task.task_uuid, dataToSave);
        console.log('✅ [TaskMetadataPanel] Auto-save success');
      } catch (error) {
        console.error('❌ [TaskMetadataPanel] Auto-save failed:', error);
        message.error('自动保存失败');
      }
    }, 500);
  }, [task, caseTableData]);

  // 更新单条用例的执行结果并自动保存
  const handleCaseFieldChange = useCallback((caseId, field, value) => {
    setCaseTableData(prev => prev.map(c => 
      c.case_id === caseId ? { ...c, [field]: value } : c
    ));
    // 触发自动保存
    autoSaveCaseResult(caseId, field, value);
  }, [autoSaveCaseResult]);

  // 保存所有用例到后端
  const saveAllCasesToBackend = async (cases, filterConditions) => {
    console.log('💾 [TaskMetadataPanel] saveAllCasesToBackend called');
    console.log('💾 [TaskMetadataPanel] task:', task?.task_name, 'task_uuid:', task?.task_uuid);
    console.log('💾 [TaskMetadataPanel] cases count:', cases?.length);
    console.log('💾 [TaskMetadataPanel] filterConditions:', filterConditions);
    
    if (!task || !task.task_uuid) {
      console.error('💾 [TaskMetadataPanel] ERROR: No task or task_uuid!');
      message.error('任务信息缺失');
      return;
    }
    if (!cases || cases.length === 0) {
      console.error('💾 [TaskMetadataPanel] ERROR: No cases to save!');
      message.error('没有用例可保存');
      return;
    }
    
    try {
      // 保存筛选条件到localStorage
      localStorage.setItem(`execution_filter_${task.task_uuid}`, JSON.stringify(filterConditions));
      console.log('💾 [TaskMetadataPanel] Filter saved to localStorage');
      
      // 先从后端加载已有的执行结果，以便合并已保存的 test_result、bug_id、remark
      let existingResults = [];
      try {
        existingResults = await getExecutionCaseResults(task.task_uuid);
        console.log('💾 [TaskMetadataPanel] Loaded existing results:', existingResults?.length || 0);
      } catch (e) {
        console.log('💾 [TaskMetadataPanel] No existing results found');
      }
      
      // 创建已有结果的映射 (case_id -> result)
      const existingMap = new Map();
      if (existingResults && existingResults.length > 0) {
        existingResults.forEach(r => {
          existingMap.set(r.case_id, r);
        });
      }
      
      const isManual = filterConditions?.execution_type === 'manual';
      // 获取用例类型：手工测试用 case_type (overall/acceptance/change)，AI测试用 role
      const caseType = isManual 
        ? (filterConditions?.case_type || 'overall')
        : (filterConditions?.case_type || 'role1');
      
      // 构造保存数据，合并已有的执行结果
      // 按选择顺序重新生成 No.（display_id）
      const dataToSave = cases.map((c, index) => {
        // 查找已有的执行结果
        const existing = existingMap.get(c.case_id);
        
        // 判断是否保留已有的test_result：只有OK/NG才保留（已执行过的结果）
        // NR和Block都视为未执行，重新选择时重置为Block
        const preservedResults = ['OK', 'NG'];
        const shouldPreserveResult = existing?.test_result && preservedResults.includes(existing.test_result);
        
        // 🔍 调试: 打印源用例的 script_code
        if (index === 0) {
          console.log('🔍 [saveAllCasesToBackend] c (source case):', c);
          console.log('🔍 [saveAllCasesToBackend] c.script_code:', c.script_code);
          console.log('🔍 [saveAllCasesToBackend] c keys:', Object.keys(c));
        }
        
        const item = {
          case_id: c.case_id,
          display_id: index + 1, // 按选择顺序重新生成序号（从1开始）
          case_num: c.case_number || c.case_num || '', // 用户自定义CaseID
          case_type: caseType, // 用例类型
          // 只保留OK/NG结果，其他情况默认为Block
          test_result: shouldPreserveResult ? existing.test_result : 'Block',
          bug_id: existing?.bug_id || c.bug_id || '',
          remark: existing?.remark || c.remark || '',
        };
        
        if (isManual) {
          // 手工测试用例的字段
          item.major_function_cn = c.major_function_cn || c.major_function || '';
          item.major_function_jp = c.major_function_jp || '';
          item.major_function_en = c.major_function_en || '';
          item.middle_function_cn = c.middle_function_cn || c.middle_function || '';
          item.middle_function_jp = c.middle_function_jp || '';
          item.middle_function_en = c.middle_function_en || '';
          item.minor_function_cn = c.minor_function_cn || c.minor_function || '';
          item.minor_function_jp = c.minor_function_jp || '';
          item.minor_function_en = c.minor_function_en || '';
          item.precondition_cn = c.precondition_cn || c.precondition || '';
          item.precondition_jp = c.precondition_jp || '';
          item.precondition_en = c.precondition_en || '';
          item.test_steps_cn = c.test_steps_cn || c.test_steps || '';
          item.test_steps_jp = c.test_steps_jp || '';
          item.test_steps_en = c.test_steps_en || '';
          item.expected_result_cn = c.expected_result_cn || c.expected_result || '';
          item.expected_result_jp = c.expected_result_jp || '';
          item.expected_result_en = c.expected_result_en || '';
        } else if (filterConditions?.execution_type === 'api') {
          // API 用例的字段（无多语言）
          item.screen = c.screen || '';
          item.url = c.url || '';
          item.header = c.header || '';
          item.method = c.method || '';
          item.body = c.body || '';
          item.response = c.response || '';
          item.response_time = c.response_time || '';
          item.script_code = c.script_code || ''; // JS脚本代码，用于API测试执行
          // 🔍 调试: 打印构建后的 item.script_code
          if (index === 0) {
            console.log('🔍 [saveAllCasesToBackend] item.script_code:', item.script_code);
            console.log('🔍 [saveAllCasesToBackend] item (built):', item);
          }
        } else {
          // AI Web 用例的字段
          item.screen_cn = c.screen_cn || '';
          item.screen_jp = c.screen_jp || '';
          item.screen_en = c.screen_en || '';
          item.function_cn = c.function_cn || '';
          item.function_jp = c.function_jp || '';
          item.function_en = c.function_en || '';
          item.precondition_cn = c.precondition_cn || '';
          item.precondition_jp = c.precondition_jp || '';
          item.precondition_en = c.precondition_en || '';
          item.test_steps_cn = c.test_steps_cn || '';
          item.test_steps_jp = c.test_steps_jp || '';
          item.test_steps_en = c.test_steps_en || '';
          item.expected_result_cn = c.expected_result_cn || '';
          item.expected_result_jp = c.expected_result_jp || '';
          item.expected_result_en = c.expected_result_en || '';
          item.script_code = c.script_code || ''; // Playwright脚本代码，用于Web自动化执行
        }
        
        return item;
      });
      
      console.log('💾 [TaskMetadataPanel] dataToSave[0]:', dataToSave[0]);
      console.log('💾 [TaskMetadataPanel] dataToSave[0].case_id:', dataToSave[0]?.case_id);
      console.log('💾 [TaskMetadataPanel] Calling saveExecutionCaseResults API...');
      
      await saveExecutionCaseResults(task.task_uuid, dataToSave);
      console.log('✅ [TaskMetadataPanel] All cases saved successfully!');
      message.success(`已保存 ${cases.length} 条用例`);
    } catch (error) {
      console.error('❌ [TaskMetadataPanel] Save cases failed:', error);
      console.error('❌ [TaskMetadataPanel] Error details:', error.response?.data || error.message);
      message.error('保存用例失败: ' + (error.response?.data?.message || error.message));
    }
  };

  // 获取语言后缀 - 使用displayLanguage作为当前显示语言
  // 优先级：displayLanguage状态 > filterConditions > task.display_language > 默认cn
  const getLanguageSuffix = () => {
    const lang = displayLanguage || selectedCasesData?.filterConditions?.language || task?.display_language || 'cn';
    return lang === 'cn' ? '_cn' : lang === 'jp' ? '_jp' : '_en';
  };

  // 获取语言显示名
  const getLanguageDisplay = () => {
    // 优先级：displayLanguage状态 > filterConditions > task.display_language > 默认cn
    const lang = displayLanguage || selectedCasesData?.filterConditions?.language || task?.display_language || 'cn';
    return lang === 'cn' ? 'CN' : lang === 'jp' ? 'JP' : 'EN';
  };

  // 获取执行任务的语言显示值（用于元数据面板显示）
  // - Web: 未选择用例显示"-"，选择后显示选择的语言（EN/JP/CN）
  // - API: 未选择用例显示"-"，选择后显示"EN"
  // - Manual: 未选择用例显示"-"，选择后显示"ALL"
  const getExecutionLanguageDisplay = () => {
    const hasCases = selectedCasesData && selectedCasesData.cases && selectedCasesData.cases.length > 0;
    if (!hasCases) {
      return '-';
    }
    
    const execType = task?.execution_type;
    if (execType === 'automation') {
      // Web类型：显示选择的语言，优先级：filterConditions > displayLanguage状态 > task.display_language > 默认cn
      const lang = selectedCasesData?.filterConditions?.language || displayLanguage || task?.display_language || 'cn';
      return lang === 'cn' ? 'CN' : lang === 'jp' ? 'JP' : 'EN';
    } else if (execType === 'api') {
      // API类型：固定显示EN
      return 'EN';
    } else if (execType === 'manual') {
      // Manual类型：固定显示ALL
      return 'ALL';
    }
    return '-';
  };

  // 判断是否为手工测试类型
  const isManualType = () => {
    return selectedCasesData?.filterConditions?.execution_type === 'manual' || task?.execution_type === 'manual';
  };

  // 处理语言切换
  const handleLanguageChange = (e) => {
    setDisplayLanguage(e.target.value);
  };

  // 下载当前表格内容
  // 手工测试和 AI Web: xlsx 格式
  // AI 接口: csv 格式
  // 文件名格式 (FR-06):
  // - 手工测试: {项目名}_{任务名}_ManualCases_{语言}_{日期}.xlsx
  // - AI Web: {项目名}_{任务名}_AIWebCases_{角色}_{语言}_{日期}.xlsx
  // - AI 接口: {项目名}_{任务名}_AICases_{角色}_{日期}.csv
  const handleDownloadCases = () => {
    if (!caseTableData || caseTableData.length === 0) {
      message.warning('没有用例数据可下载');
      return;
    }

    const langSuffix = getLanguageSuffix();
    const langDisplay = getLanguageDisplay();
    const isManual = isManualType();
    const isAIAPI = task?.execution_type === 'api';
    const isAIWeb = task?.execution_type === 'automation';
    
    // 清理项目名和任务名，去除非法字符
    const safeProjectName = (projectName || 'Project')?.replace(/[\\/:*?"<>|]/g, '_');
    const safeTaskName = (task?.task_name || 'Task')?.replace(/[\\/:*?"<>|]/g, '_');
    const dateStr = dayjs().format('YYYYMMDD');
    const role = (selectedCasesData?.filterConditions?.case_type || 'role1').toUpperCase();

    // 根据执行类型选择不同的表头
    let headers;
    let rows;
    
    if (isManual) {
      // 手工测试用例：No./CaseID/Maj.Category/Mid.Category/Min.Category/Precondition/Test Step/Expect/TestResult/BugID/Remark
      headers = ['No.', 'CaseID', `Maj.Category${langDisplay}`, `Mid.Category${langDisplay}`, `Min.Category${langDisplay}`, `Precondition${langDisplay}`, `Test Step${langDisplay}`, `Expect${langDisplay}`, 'TestResult', 'BugID', 'Remark'];
      
      rows = caseTableData.map((c, index) => [
        index + 1,
        c.case_number || c.case_num || '',
        c[`major_function${langSuffix}`] || c.major_function_cn || c.major_function || '',
        c[`middle_function${langSuffix}`] || c.middle_function_cn || c.middle_function || '',
        c[`minor_function${langSuffix}`] || c.minor_function_cn || c.minor_function || '',
        c[`precondition${langSuffix}`] || c.precondition || '',
        c[`test_steps${langSuffix}`] || c.test_steps || '',
        c[`expected_result${langSuffix}`] || c.expected_result || '',
        c.test_result || 'Block',
        c.bug_id || '',
        c.remark || '',
      ]);
    } else if (isAIAPI) {
      // API 用例：No./CaseID/Screen/URL/Header/Method/Body/Response/ResponseTime/TestResult/BugID/Remark
      headers = ['No.', 'CaseID', 'Screen', 'URL', 'Header', 'Method', 'Body', 'Response', 'ResponseTime', 'TestResult', 'BugID', 'Remark'];
      
      rows = caseTableData.map((c, index) => [
        index + 1,
        c.case_num || c.case_number || '',
        c.screen || '',
        c.url || '',
        c.header || '',
        c.method || '',
        c.body || '',
        c.response || '',
        c.response_time ? `${c.response_time} ms` : '',
        c.test_result || 'Block',
        c.bug_id || '',
        c.remark || '',
      ]);
    } else {
      // AI Web 用例：No./CaseID/Screen/Function/Precondition/Test Step/Expect/TestResult/BugID/Remark
      headers = ['No.', 'CaseID', `Screen${langDisplay}`, `Function${langDisplay}`, `Precondition${langDisplay}`, `Test Step${langDisplay}`, `Expect${langDisplay}`, 'TestResult', 'BugID', 'Remark'];
      
      rows = caseTableData.map((c, index) => [
        index + 1,
        c.case_num || '',
        c[`screen${langSuffix}`] || '',
        c[`function${langSuffix}`] || '',
        c[`precondition${langSuffix}`] || '',
        c[`test_steps${langSuffix}`] || '',
        c[`expected_result${langSuffix}`] || '',
        c.test_result || 'Block',
        c.bug_id || '',
        c.remark || '',
      ]);
    }

    // 构建元数据行
    const metadataRows = [
      ['任务名称', task?.task_name || ''],
      ['执行内容', task?.execution_type === 'manual' ? '手工测试' : task?.execution_type === 'automation' ? 'AI Web' : 'AI接口'],
      ['任务状态', task?.task_status === 'pending' ? '待开始' : task?.task_status === 'in_progress' ? '进行中' : '已完成'],
      ['开始日期', task?.start_date ? dayjs(task.start_date).format('YYYY-MM-DD') : ''],
      ['结束日期', task?.end_date ? dayjs(task.end_date).format('YYYY-MM-DD') : ''],
      ['测试日期', task?.test_date ? dayjs(task.test_date).format('YYYY-MM-DD') : ''],
      ['测试版本', task?.test_version || ''],
      ['测试环境', task?.test_env || ''],
      ['执行人', task?.executor || ''],
      ['任务描述', task?.task_description || ''],
      ['筛选条件', isManual 
        ? `${selectedCasesData?.filterConditions?.case_type_display || '整体'}用例` 
        : role],
      ['语言', langDisplay],
      ['用例数量', caseTableData.length.toString()],
    ];

    // AI 接口测试下载 CSV 格式
    if (isAIAPI) {
      const escapeCSV = (str) => {
        if (str == null) return '';
        const s = String(str);
        if (s.includes(',') || s.includes('"') || s.includes('\n')) {
          return `"${s.replace(/"/g, '""')}"`;
        }
        return s;
      };

      const csvContent = [
        ...metadataRows.map(row => row.map(escapeCSV).join(',')),
        [], // 空行分隔
        headers.map(escapeCSV).join(','),
        ...rows.map(row => row.map(escapeCSV).join(','))
      ].join('\n');

      const BOM = '\uFEFF';
      const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' });
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      // AI 接口: {项目名}_{任务名}_API_TestResult_{时间戳}.csv
      link.download = `${safeProjectName}_${safeTaskName}_API_TestResult_${dateStr}.csv`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    } else {
      // 手工测试和 AI Web 下载 xlsx 格式
      const workbook = XLSX.utils.book_new();
      
      // 构建工作表数据：元数据 + 空行 + 表头 + 数据
      const wsData = [
        ...metadataRows,
        [], // 空行分隔
        headers,
        ...rows
      ];
      
      const worksheet = XLSX.utils.aoa_to_sheet(wsData);
      
      // 设置列宽
      const colWidths = isManual 
        ? [5, 15, 20, 20, 20, 30, 40, 30, 10, 15, 20]
        : [5, 15, 20, 20, 30, 40, 30, 10, 15, 20];
      worksheet['!cols'] = colWidths.map(width => ({ wch: width }));
      
      XLSX.utils.book_append_sheet(workbook, worksheet, 'TestCases');
      
      // 生成文件名
      let filename;
      if (isManual) {
        // 手工测试: {项目名}_{任务名}_Manual_{语言}_TestResult_{时间戳}.xlsx
        filename = `${safeProjectName}_${safeTaskName}_Manual_${langDisplay}_TestResult_${dateStr}.xlsx`;
      } else if (isAIWeb) {
        // AI Web: {项目名}_{任务名}_Web_{语言}_TestResult_{时间戳}.xlsx
        filename = `${safeProjectName}_${safeTaskName}_Web_${langDisplay}_TestResult_${dateStr}.xlsx`;
      }
      
      // 下载文件
      XLSX.writeFile(workbook, filename);
    }
    
    message.success('下载成功');
  };

  // 生成表格列配置
  const getCaseTableColumns = () => {
    const langSuffix = getLanguageSuffix();
    const langDisplay = getLanguageDisplay();
    const isManual = isManualType();
    const isAPI = task?.execution_type === 'api';

    // 展开详情列 - 放在最前面
    const expandColumn = {
      title: '',
      key: 'expand_action',
      width: 50,
      fixed: 'left',
      render: (_, record) => (
        <Button
          type="text"
          size="small"
          icon={<EyeOutlined />}
          onClick={() => handleOpenCaseDetail(record)}
          style={{ color: '#1890ff' }}
          title={t('testExecution.caseDetail.viewDetail')}
        />
      ),
    };

    // 公共列：No. 和 CaseID（API类型不显示CaseID）
    const commonStartColumns = [
      expandColumn,
      {
        title: 'No.',
        dataIndex: 'no',
        key: 'no',
        width: 60,
        fixed: 'left',
      },
    ];
    
    // API类型不显示CaseID列
    if (!isAPI) {
      commonStartColumns.push({
        title: 'CaseID',
        key: 'case_id_display',
        width: 120,
        render: (_, record) => record.case_number || record.case_num || '-',
      });
    }

    // 公共列：TestResult、BugID、Remark
    const commonEndColumns = [
      {
        title: 'TestResult',
        dataIndex: 'test_result',
        key: 'test_result',
        width: 100,
        fixed: 'right',
        render: (value, record) => {
          const getTagColor = (val) => {
            const colorMap = {
              'OK': 'success',
              'NG': 'error',
              'Block': 'warning',
              'NR': 'default',
            };
            return colorMap[val] || 'default';
          };
          
          return (
            <Select
              value={value || 'Block'}
              size="small"
              style={{ width: 90 }}
              onChange={(val) => handleCaseFieldChange(record.case_id, 'test_result', val)}
            >
              <Option value="NR"><Tag color="default" style={{ margin: 0 }}>NR</Tag></Option>
              <Option value="OK"><Tag color="success" style={{ margin: 0 }}>OK</Tag></Option>
              <Option value="NG"><Tag color="error" style={{ margin: 0 }}>NG</Tag></Option>
              <Option value="Block"><Tag color="warning" style={{ margin: 0 }}>Block</Tag></Option>
            </Select>
          );
        },
      },
      {
        title: 'BugID',
        dataIndex: 'bug_id',
        key: 'bug_id',
        width: 120,
        fixed: 'right',
        render: (value, record) => (
          <Input
            defaultValue={value || ''}
            size="small"
            placeholder="Bug ID"
            onBlur={(e) => {
              if (e.target.value !== value) {
                handleCaseFieldChange(record.case_id, 'bug_id', e.target.value);
              }
            }}
            onPressEnter={(e) => {
              e.target.blur();
            }}
          />
        ),
      },
      {
        title: 'Remark',
        dataIndex: 'remark',
        key: 'remark',
        width: isManual ? 200 : 150,
        fixed: 'right',
        render: (value, record) => {
          // Manual类型使用多行TextArea，其他类型使用单行Input
          if (isManual) {
            return (
              <Input.TextArea
                defaultValue={value || ''}
                size="small"
                placeholder="备注"
                autoSize={{ minRows: 2, maxRows: 4 }}
                style={{ resize: 'vertical' }}
                onBlur={(e) => {
                  if (e.target.value !== value) {
                    handleCaseFieldChange(record.case_id, 'remark', e.target.value);
                  }
                }}
              />
            );
          }
          return (
            <Input
              defaultValue={value || ''}
              size="small"
              placeholder="备注"
              onBlur={(e) => {
                if (e.target.value !== value) {
                  handleCaseFieldChange(record.case_id, 'remark', e.target.value);
                }
              }}
              onPressEnter={(e) => {
                e.target.blur();
              }}
            />
          );
        },
      },
    ];

    // 根据执行类型选择中间列
    let middleColumns;
    
    if (isManual) {
      // 手工测试用例的列：大功能/中功能/小功能/前置条件/测试步骤/期望结果
      // 辅助函数：判断是否与上一行相同
      const isSameAsPrevious = (record, field) => {
        const index = caseTableData.findIndex(c => c.key === record.key);
        if (index <= 0) return false;
        const prevRecord = caseTableData[index - 1];
        const currentValue = record[`${field}${langSuffix}`] || record[`${field}_cn`] || record[field] || '';
        const prevValue = prevRecord[`${field}${langSuffix}`] || prevRecord[`${field}_cn`] || prevRecord[field] || '';
        return currentValue === prevValue && currentValue !== '';
      };
      
      // 判断大功能和中功能都相同
      const isSameMajorAndMiddle = (record) => {
        return isSameAsPrevious(record, 'major_function') && isSameAsPrevious(record, 'middle_function');
      };

      middleColumns = [
        {
          title: `Maj.Category${langDisplay}`,
          key: 'major_function',
          width: 120,
          ellipsis: true,
          render: (_, record) => {
            const value = record[`major_function${langSuffix}`] || record.major_function_cn || record.major_function || '-';
            const isSame = isSameAsPrevious(record, 'major_function');
            return <span style={{ color: isSame ? '#d9d9d9' : 'inherit' }}>{value}</span>;
          },
        },
        {
          title: `Mid.Category${langDisplay}`,
          key: 'middle_function',
          width: 120,
          ellipsis: true,
          render: (_, record) => {
            const value = record[`middle_function${langSuffix}`] || record.middle_function_cn || record.middle_function || '-';
            const isSame = isSameMajorAndMiddle(record);
            return <span style={{ color: isSame ? '#d9d9d9' : 'inherit' }}>{value}</span>;
          },
        },
        {
          title: `Min.Category${langDisplay}`,
          key: 'minor_function',
          width: 120,
          ellipsis: true,
          render: (_, record) => record[`minor_function${langSuffix}`] || record.minor_function_cn || record.minor_function || '-',
        },
        {
          title: `Precondition${langDisplay}`,
          key: 'precondition',
          width: 150,
          render: (_, record) => {
            const value = record[`precondition${langSuffix}`] || record.precondition || '-';
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{value}</div>;
          },
        },
        {
          title: `Test Step${langDisplay}`,
          key: 'test_step',
          width: 200,
          render: (_, record) => {
            const value = record[`test_steps${langSuffix}`] || record.test_steps || '-';
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{value}</div>;
          },
        },
        {
          title: `Expect${langDisplay}`,
          key: 'expect',
          width: 150,
          render: (_, record) => {
            const value = record[`expected_result${langSuffix}`] || record.expected_result || '-';
            return <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{value}</div>;
          },
        },
      ];
    } else if (task?.execution_type === 'api') {
      // API 用例的列：URL/Header/Method/Body/Response (Screen列已移除)
      middleColumns = [
        {
          title: 'URL',
          key: 'url',
          dataIndex: 'url',
          width: 200,
          ellipsis: true,
          render: (value) => {
            const displayValue = value || '-';
            return (
              <Tooltip title={displayValue !== '-' ? displayValue : ''} placement="topLeft">
                <div className="single-line-cell">{displayValue}</div>
              </Tooltip>
            );
          },
        },
        {
          title: 'Header',
          key: 'header',
          dataIndex: 'header',
          width: 150,
          ellipsis: true,
          render: (value) => {
            const displayValue = value || '-';
            return (
              <Tooltip title={displayValue !== '-' ? displayValue : ''} placement="topLeft">
                <div className="single-line-cell">{displayValue}</div>
              </Tooltip>
            );
          },
        },
        {
          title: 'Method',
          key: 'method',
          dataIndex: 'method',
          width: 80,
          ellipsis: true,
          render: (value) => {
            const displayValue = value || '-';
            return (
              <Tooltip title={displayValue !== '-' ? displayValue : ''} placement="topLeft">
                <div className="single-line-cell">{displayValue}</div>
              </Tooltip>
            );
          },
        },
        {
          title: 'Body',
          key: 'body',
          dataIndex: 'body',
          width: 200,
          ellipsis: true,
          render: (value) => {
            const displayValue = value || '-';
            return (
              <Tooltip title={displayValue !== '-' ? displayValue : ''} placement="topLeft">
                <div className="single-line-cell">{displayValue}</div>
              </Tooltip>
            );
          },
        },
        {
          title: 'Response',
          key: 'response',
          dataIndex: 'response',
          width: 120,
          ellipsis: true,
          render: (value) => {
            const displayValue = value || '-';
            return (
              <Tooltip title={displayValue !== '-' ? displayValue : ''} placement="topLeft">
                <div className="single-line-cell">{displayValue}</div>
              </Tooltip>
            );
          },
        },
        {
          title: 'ResponseTime',
          key: 'response_time',
          dataIndex: 'response_time',
          width: 120,
          ellipsis: true,
          render: (value) => {
            if (!value || value === '-') {
              return <div className="single-line-cell">-</div>;
            }
            
            const responseTime = Number(value);
            const isSlow = responseTime > 3000; // 超过3秒
            
            return (
              <Tooltip 
                title={isSlow ? t('testExecution.responseTime.slowWarning', { time: responseTime }) : `${responseTime} ms`} 
                placement="topLeft"
              >
                <div 
                  className="single-line-cell" 
                  style={{
                    color: isSlow ? '#ff4d4f' : '#303133',
                    fontWeight: isSlow ? 600 : 400,
                    background: isSlow ? '#fff2f0' : 'transparent',
                    padding: isSlow ? '2px 6px' : '0',
                    borderRadius: isSlow ? '4px' : '0',
                    display: 'inline-block'
                  }}
                >
                  {isSlow && '⚠️ '}{responseTime} ms
                </div>
              </Tooltip>
            );
          },
        },
      ];
    } else {
      // AI Web 用例的列：Screen/Function/Precondition/Test Step/Expect
      middleColumns = [
        {
          title: `Screen${langDisplay}`,
          key: 'screen',
          width: 120,
          ellipsis: true,
          render: (_, record) => {
            const value = record[`screen${langSuffix}`] || '-';
            return (
              <Tooltip title={value !== '-' ? value : ''} placement="topLeft">
                <div className="single-line-cell">{value}</div>
              </Tooltip>
            );
          },
        },
        {
          title: `Function${langDisplay}`,
          key: 'function',
          width: 150,
          ellipsis: true,
          render: (_, record) => {
            const value = record[`function${langSuffix}`] || '-';
            return (
              <Tooltip title={value !== '-' ? value : ''} placement="topLeft">
                <div className="single-line-cell">{value}</div>
              </Tooltip>
            );
          },
        },
        {
          title: `Precondition${langDisplay}`,
          key: 'precondition',
          width: 150,
          ellipsis: true,
          render: (_, record) => {
            const value = record[`precondition${langSuffix}`] || '-';
            return (
              <Tooltip title={value !== '-' ? value : ''} placement="topLeft" overlayStyle={{ maxWidth: 400 }}>
                <div className="single-line-cell">{value}</div>
              </Tooltip>
            );
          },
        },
        {
          title: `Test Step${langDisplay}`,
          key: 'test_step',
          width: 200,
          ellipsis: true,
          render: (_, record) => {
            const value = record[`test_steps${langSuffix}`] || '-';
            return (
              <Tooltip title={value !== '-' ? value : ''} placement="topLeft" overlayStyle={{ maxWidth: 500 }}>
                <div className="single-line-cell">{value}</div>
              </Tooltip>
            );
          },
        },
        {
          title: `Expect${langDisplay}`,
          key: 'expect',
          width: 150,
          ellipsis: true,
          render: (_, record) => {
            const value = record[`expected_result${langSuffix}`] || '-';
            return (
              <Tooltip title={value !== '-' ? value : ''} placement="topLeft" overlayStyle={{ maxWidth: 400 }}>
                <div className="single-line-cell">{value}</div>
              </Tooltip>
            );
          },
        },
      ];
    }

    return [...commonStartColumns, ...middleColumns, ...commonEndColumns];
  };

  useEffect(() => {
    if (task) {
      const formValues = {
        ...task,
        start_date: task.start_date ? dayjs(task.start_date) : null,
        end_date: task.end_date ? dayjs(task.end_date) : null,
        test_date: task.test_date ? dayjs(task.test_date) : null,
      };
      form.setFieldsValue(formValues);
      setIsEditing(false); // 切换任务时重置编辑模式
    }
  }, [task, form]);

  const handleSave = async () => {
    console.log('\ud83d\udcbe [TaskMetadataPanel] handleSave called');
    console.log('\ud83d\udcbe [TaskMetadataPanel] Current task:', task);
    console.log('\ud83d\udcbe [TaskMetadataPanel] isEditing:', isEditing);
    
    try {
      console.log('\ud83d\udcbe [TaskMetadataPanel] Validating form fields...');
      const values = await form.validateFields();
      console.log('\u2705 [TaskMetadataPanel] Form validation passed:', values);
      
      // 验证日期范围
      if (values.start_date && values.end_date && dayjs.isDayjs(values.start_date) && dayjs.isDayjs(values.end_date)) {
        if (values.end_date.isBefore(values.start_date)) {
          console.error('\u274c [TaskMetadataPanel] Invalid date range');
          message.error(t('testExecution.metadata.invalidDateRange'));
          return;
        }
      }

      setSaving(true);
      console.log('\ud83d\udd04 [TaskMetadataPanel] Setting saving to true');

      // 转换日期为ISO 8601格式（RFC3339），符合Go后端期望
      const formattedValues = {};
      
      // 只发送被修改的字段（使用表单当前值）
      formattedValues.task_name = values.task_name;
      formattedValues.execution_type = values.execution_type;
      formattedValues.task_status = values.task_status;
      
      // 日期字段：转换为RFC3339格式（Go的time.Time默认格式）
      // 格式：YYYY-MM-DDTHH:mm:ss+08:00 或 YYYY-MM-DDTHH:mm:ssZ
      if (values.start_date) {
        // 使用午夜时间并转换为ISO格式
        formattedValues.start_date = values.start_date.startOf('day').toISOString();
      } else {
        formattedValues.start_date = null;
      }
      
      if (values.end_date) {
        formattedValues.end_date = values.end_date.startOf('day').toISOString();
      } else {
        formattedValues.end_date = null;
      }
      
      if (values.test_date) {
        formattedValues.test_date = values.test_date.startOf('day').toISOString();
      } else {
        formattedValues.test_date = null;
      }
      
      // 其他可选字段：只在有值时发送
      if (values.test_version) {
        formattedValues.test_version = values.test_version;
      }
      if (values.test_env) {
        formattedValues.test_env = values.test_env;
      }
      if (values.executor) {
        formattedValues.executor = values.executor;
      }
      if (values.task_description) {
        formattedValues.task_description = values.task_description;
      }
      
      console.log('\ud83d\udcbe [TaskMetadataPanel] Original values:', values);
      console.log('\ud83d\udcbe [TaskMetadataPanel] Formatted values:', formattedValues);
      console.log('\ud83d\udcbe [TaskMetadataPanel] Calling API with project_id:', task.project_id, 'task_uuid:', task.task_uuid);

      const response = await updateExecutionTask(task.project_id, task.task_uuid, formattedValues);
      console.log('\u2705 [TaskMetadataPanel] API response:', response);
      
      message.success(t('testExecution.metadata.saveSuccess'));
      setIsEditing(false);
      console.log('\u2705 [TaskMetadataPanel] Exited editing mode');
      
      if (onSave) {
        console.log('\ud83d\udd04 [TaskMetadataPanel] Calling onSave callback');
        onSave({ ...task, ...formattedValues });
      } else {
        console.warn('\u26a0\ufe0f [TaskMetadataPanel] onSave callback is not defined');
      }
    } catch (error) {
      console.error('\u274c [TaskMetadataPanel] handleSave error:', error);
      if (error.errorFields) {
        return;
      }
      
      if (error.response?.status === 409) {
        message.error(t('testExecution.metadata.taskNameExists'));
      } else if (error.response?.status === 400) {
        message.error(t('testExecution.metadata.validationFailed'));
      } else {
        message.error(t('testExecution.metadata.saveFailed'));
      }
    } finally {
      setSaving(false);
    }
  };

  if (!task) {
    return (
      <div className="task-metadata-panel">
        <Empty description={t('testExecution.metadata.selectTask')} />
      </div>
    );
  }

  return (
    <div className="task-metadata-panel">
      <div className="task-metadata-header">
        <Space>
          <Button
            icon={<FileSearchOutlined />}
            onClick={() => setCaseSelectionVisible(true)}
            disabled={selectedCasesData && selectedCasesData.cases && selectedCasesData.cases.length > 0}
          >
            {t('testExecution.metadata.selectCases')}
          </Button>
          <Button
            icon={<DownloadOutlined />}
            onClick={handleDownloadCases}
            disabled={!caseTableData || caseTableData.length === 0}
          >
            {t('testExecution.metadata.download')}
          </Button>
          {!isEditing ? (
            <Button
              type="primary"
              icon={<EditOutlined />}
              onClick={() => setIsEditing(true)}
            >
              {t('testExecution.metadata.edit')}
            </Button>
          ) : (
            <Space>
              <Button
                onClick={() => {
                  setIsEditing(false);
                  // 重置表单到初始值
                  form.setFieldsValue({
                    ...task,
                    start_date: task.start_date ? dayjs(task.start_date) : null,
                    end_date: task.end_date ? dayjs(task.end_date) : null,
                    test_date: task.test_date ? dayjs(task.test_date) : null,
                  });
                }}
              >
                取消
              </Button>
              <Button
                type="primary"
                icon={<SaveOutlined />}
                loading={saving}
                onClick={handleSave}
              >
                {t('testExecution.metadata.save')}
              </Button>
            </Space>
          )}
        </Space>
      </div>

      {/* 元数据区域 - 参考AI接口用例库样式 */}
      <div style={{
        padding: '12px 8px',
        background: '#fafafa',
        borderRadius: '4px',
        marginBottom: '8px'
      }}>
        <Form
          form={form}
          layout="horizontal"
          className="task-metadata-form"
        >
          {/* 第一行：任务名称、任务状态、执行内容、用例集、执行人 */}
          <Row gutter={[8, 8]}>
            <Col style={{ width: '220px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '70px', textAlign: 'right' }}>
                  <span style={{ color: '#ff4d4f' }}>*</span> {t('testExecution.metadata.taskName')}：
                </span>
                <Form.Item name="task_name" rules={[{ required: true, message: t('testExecution.metadata.taskNameRequired') }]} style={{ marginBottom: 0 }}>
                  <Input 
                    size="small" 
                    style={{ width: '140px', fontSize: '12px', backgroundColor: isEditing ? '#fff' : '#f5f5f5' }} 
                    maxLength={50} 
                    disabled={!isEditing} 
                  />
                </Form.Item>
              </div>
            </Col>
            <Col style={{ width: '170px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '60px', textAlign: 'right' }}>{t('testExecution.metadata.taskStatus')}：</span>
                <Form.Item name="task_status" style={{ marginBottom: 0 }}>
                  <Select 
                    size="small" 
                    style={{ width: '100px', fontSize: '12px' }} 
                    disabled={!isEditing}
                    className={!isEditing ? 'metadata-select-readonly' : ''}
                  >
                    <Option value="pending">{t('testExecution.status.pending')}</Option>
                    <Option value="in_progress">{t('testExecution.status.inProgress')}</Option>
                    <Option value="completed">{t('testExecution.status.completed')}</Option>
                  </Select>
                </Form.Item>
              </div>
            </Col>
            <Col style={{ width: '180px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '60px', textAlign: 'right' }}>{t('testExecution.metadata.executionType')}：</span>
                <Form.Item name="execution_type" style={{ marginBottom: 0 }}>
                  <Select 
                    size="small" 
                    style={{ width: '110px', fontSize: '12px' }} 
                    disabled 
                    className="metadata-select-readonly"
                  >
                    <Option value="manual">Manual Test</Option>
                    <Option value="automation">AI Web</Option>
                    <Option value="api">AI API</Option>
                  </Select>
                </Form.Item>
              </div>
            </Col>
            <Col style={{ width: '170px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '60px', textAlign: 'right' }}>{t('testExecution.metadata.caseGroup')}：</span>
                <Form.Item style={{ marginBottom: 0 }}>
                  <Input 
                    size="small"
                    style={{ width: '100px', fontSize: '12px', backgroundColor: '#f5f5f5' }}
                    value={selectedCasesData?.filterConditions?.case_group || task?.case_group_name || '-'} 
                    disabled 
                  />
                </Form.Item>
              </div>
            </Col>
            <Col style={{ width: '170px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '60px', textAlign: 'right' }}>{t('testExecution.metadata.executor')}：</span>
                <Form.Item name="executor" style={{ marginBottom: 0 }}>
                  <Input 
                    size="small" 
                    style={{ width: '100px', fontSize: '12px', backgroundColor: isEditing ? '#fff' : '#f5f5f5' }} 
                    maxLength={50} 
                    disabled={!isEditing} 
                  />
                </Form.Item>
              </div>
            </Col>
            <Col style={{ width: '170px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '60px', textAlign: 'right' }}>{t('testExecution.metadata.language', '语言')}：</span>
                <Form.Item style={{ marginBottom: 0 }}>
                  <Input 
                    size="small"
                    style={{ width: '100px', fontSize: '12px', backgroundColor: '#f5f5f5' }}
                    value={getExecutionLanguageDisplay()} 
                    disabled 
                  />
                </Form.Item>
              </div>
            </Col>
          </Row>

          {/* 第二行：测试环境、测试版本、开始日期、结束日期、测试日期 */}
          <Row gutter={[8, 8]} style={{ marginTop: '8px' }}>
            <Col style={{ width: '220px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '70px', textAlign: 'right' }}>{t('testExecution.metadata.testEnv')}：</span>
                <Form.Item name="test_env" style={{ marginBottom: 0 }}>
                  <Input 
                    size="small" 
                    style={{ width: '140px', fontSize: '12px', backgroundColor: isEditing ? '#fff' : '#f5f5f5' }} 
                    maxLength={100} 
                    disabled={!isEditing} 
                  />
                </Form.Item>
              </div>
            </Col>
            <Col style={{ width: '180px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '70px', textAlign: 'right' }}>{t('testExecution.metadata.testVersion')}：</span>
                <Form.Item name="test_version" style={{ marginBottom: 0 }}>
                  <Input 
                    size="small" 
                    style={{ width: '100px', fontSize: '12px', backgroundColor: isEditing ? '#fff' : '#f5f5f5' }} 
                    maxLength={50} 
                    disabled={!isEditing} 
                  />
                </Form.Item>
              </div>
            </Col>
            <Col style={{ width: '190px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '70px', textAlign: 'right' }}>{t('testExecution.metadata.startDate')}：</span>
                <Form.Item name="start_date" style={{ marginBottom: 0 }}>
                  <DatePicker 
                    size="small" 
                    format="YYYY-MM-DD" 
                    style={{ width: '110px', fontSize: '12px' }} 
                    disabled={!isEditing} 
                    className={!isEditing ? 'metadata-picker-readonly' : ''}
                    placeholder=""
                  />
                </Form.Item>
              </div>
            </Col>
            <Col style={{ width: '180px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '70px', textAlign: 'right' }}>{t('testExecution.metadata.endDate')}：</span>
                <Form.Item name="end_date" style={{ marginBottom: 0 }}>
                  <DatePicker 
                    size="small" 
                    format="YYYY-MM-DD" 
                    style={{ width: '100px', fontSize: '12px' }} 
                    disabled={!isEditing} 
                    className={!isEditing ? 'metadata-picker-readonly' : ''}
                    placeholder=""
                  />
                </Form.Item>
              </div>
            </Col>
            <Col style={{ width: '180px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', whiteSpace: 'nowrap', width: '70px', textAlign: 'right' }}>{t('testExecution.metadata.testDate')}：</span>
                <Form.Item name="test_date" style={{ marginBottom: 0 }}>
                  <DatePicker 
                    size="small" 
                    format="YYYY-MM-DD" 
                    style={{ width: '100px', fontSize: '12px' }} 
                    disabled={!isEditing} 
                    className={!isEditing ? 'metadata-picker-readonly' : ''}
                    placeholder=""
                  />
                </Form.Item>
              </div>
            </Col>
          </Row>
        </Form>
      </div>

      {/* 用例选择弹窗 */}
      <Modal
        title={t('testExecution.metadata.selectCases')}
        open={caseSelectionVisible}
        onCancel={() => {
          console.log('🔴 [TaskMetadataPanel] Modal cancelled');
          setCaseSelectionVisible(false);
        }}
        footer={null}
        width={500}
        destroyOnClose
      >
        <CaseSelectionPanel
          task={task}
          projectId={projectId}
          onConfirm={async (data) => {
            console.log('🟢 [TaskMetadataPanel] onConfirm callback received!');
            console.log('🟢 [TaskMetadataPanel] cases count:', data?.cases?.length);
            console.log('🟢 [TaskMetadataPanel] filterConditions:', data?.filterConditions);
            
            setCaseSelectionVisible(false);
            
            // 设置显示语言为用户选择的语言（AIWeb用例）
            if (data.filterConditions?.language) {
              setDisplayLanguage(data.filterConditions.language);
            }
            
            // 更新任务的用例集信息和显示语言到数据库
            const caseGroupName = data.filterConditions?.case_group || '';
            const selectedLanguage = data.filterConditions?.language || '';
            // 根据执行类型确定保存的语言值
            let displayLangToSave = '';
            if (task?.execution_type === 'automation') {
              displayLangToSave = selectedLanguage || 'cn';
            } else if (task?.execution_type === 'api') {
              displayLangToSave = 'en';
            } else if (task?.execution_type === 'manual') {
              displayLangToSave = 'all';
            }
            
            if (task?.task_uuid) {
              try {
                console.log('💾 [TaskMetadataPanel] Updating task case_group_name:', caseGroupName, 'display_language:', displayLangToSave);
                await updateExecutionTask(projectId, task.task_uuid, {
                  case_group_name: caseGroupName,
                  display_language: displayLangToSave
                });
                console.log('✅ [TaskMetadataPanel] Task case_group_name and display_language updated successfully');
                // 通知父组件更新任务信息
                if (onSave) {
                  onSave({ ...task, case_group_name: caseGroupName, display_language: displayLangToSave });
                }
              } catch (error) {
                console.error('❌ [TaskMetadataPanel] Failed to update task:', error);
                // 不阻止后续流程
              }
            }
            
            // 保存用例到后端（会合并已有的执行结果）
            await saveAllCasesToBackend(data.cases, data.filterConditions);
            
            // 保存后重新加载数据，确保显示最新的执行结果
            await loadSavedCaseResults();
          }}
        />
      </Modal>

      {/* 选中的用例展示区域 */}
      {selectedCasesData && selectedCasesData.cases && selectedCasesData.cases.length > 0 && (
        <div className="selected-cases-section" style={{ marginTop: 16 }}>
          {/* 统计信息区域 */}
          {(() => {
            const total = caseTableData.length;
            const okCount = caseTableData.filter(c => c.test_result === 'OK').length;
            const ngCount = caseTableData.filter(c => c.test_result === 'NG').length;
            const blockCount = caseTableData.filter(c => c.test_result === 'Block').length;
            const nrCount = caseTableData.filter(c => c.test_result === 'NR').length;
            // 实施进度 = (OK + NG + NR) / 总用例数
            const processedCount = okCount + ngCount + nrCount;
            const progressPercent = total > 0 ? Math.round((processedCount / total) * 100) : 0;
            // 通过率 = OK / (总数 - NR)
            const requiredCount = total - nrCount;
            const passRatePercent = requiredCount > 0 ? Math.round((okCount / requiredCount) * 100) : 0;
            
            return (
              <div style={{ marginBottom: 8, display: 'flex', justifyContent: 'flex-start', alignItems: 'center', gap: 16, flexWrap: 'wrap' }}>
                {/* 手工测试显示语言筛选按钮 */}
                {isManualType() && (
                  <Radio.Group 
                    value={displayLanguage} 
                    onChange={handleLanguageChange} 
                    size="small"
                  >
                    <Radio.Button value="cn">CN</Radio.Button>
                    <Radio.Button value="jp">JP</Radio.Button>
                    <Radio.Button value="en">EN</Radio.Button>
                  </Radio.Group>
                )}
                
                {/* 统计数字 */}
                <Space size={12}>
                  <span style={{ color: '#52c41a', fontWeight: 'bold' }}>OK: {okCount}</span>
                  <span style={{ color: '#ff4d4f', fontWeight: 'bold' }}>NG: {ngCount}</span>
                  <span style={{ color: '#faad14', fontWeight: 'bold' }}>Block: {blockCount}</span>
                  <span style={{ color: '#8c8c8c', fontWeight: 'bold' }}>NR: {nrCount}</span>
                </Space>
                
                {/* 实施进度条 */}
                <div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
                  <span style={{ fontSize: 12, color: '#666' }}>{t('testExecution.statistics.progress')}:</span>
                  <div style={{ width: 100, height: 16, backgroundColor: '#f0f0f0', borderRadius: 8, overflow: 'hidden' }}>
                    <div style={{ 
                      width: `${progressPercent}%`, 
                      height: '100%', 
                      backgroundColor: '#1890ff',
                      borderRadius: 8,
                      transition: 'width 0.3s'
                    }} />
                  </div>
                  <span style={{ fontSize: 12, color: '#666', minWidth: 36 }}>{progressPercent}%</span>
                </div>
                
                {/* 通过率条 */}
                <div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
                  <span style={{ fontSize: 12, color: '#666' }}>{t('testExecution.statistics.passRate')}:</span>
                  <div style={{ width: 100, height: 16, backgroundColor: '#f0f0f0', borderRadius: 8, overflow: 'hidden' }}>
                    <div style={{ 
                      width: `${passRatePercent}%`, 
                      height: '100%', 
                      backgroundColor: passRatePercent >= 80 ? '#52c41a' : passRatePercent >= 50 ? '#faad14' : '#ff4d4f',
                      borderRadius: 8,
                      transition: 'width 0.3s'
                    }} />
                  </div>
                  <span style={{ fontSize: 12, color: '#666', minWidth: 36 }}>{passRatePercent}%</span>
                </div>
              </div>
            );
          })()}
          <Table
            columns={getCaseTableColumns()}
            dataSource={caseTableData}
            size="small"
            scroll={{ y: 400 }}
            pagination={{
              current: currentPage,
              pageSize: pageSize,
              showSizeChanger: true,
              showQuickJumper: true,
              pageSizeOptions: ['10', '20', '50', '100'],
              showTotal: (total) => `${t('common.total')} ${total} ${t('common.items')}`,
              onChange: (page, size) => {
                console.log('📄 [Pagination] Page changed:', page, 'Size:', size);
                setCurrentPage(page);
                setPageSize(size);
              },
              onShowSizeChange: (current, size) => {
                console.log('📄 [Pagination] Size changed:', size, 'Current page:', current);
                setCurrentPage(1); // 切换分页大小时重置到第一页
                setPageSize(size);
              },
            }}
            bordered
          />
        </div>
      )}
      
      {/* 用例详细信息弹窗 */}
      <CaseDetailModal
        visible={caseDetailVisible}
        caseData={selectedCaseForDetail}
        executionType={task?.execution_type || 'automation'}
        languageSuffix={getLanguageSuffix()}
        languageDisplay={getLanguageDisplay()}
        onSave={handleSaveCaseDetail}
        onCancel={handleCloseCaseDetail}
      />
    </div>
  );
};

TaskMetadataPanel.propTypes = {
  task: PropTypes.object,
  projectId: PropTypes.number,
  projectName: PropTypes.string,
  onSave: PropTypes.func.isRequired,
};

export default TaskMetadataPanel;
