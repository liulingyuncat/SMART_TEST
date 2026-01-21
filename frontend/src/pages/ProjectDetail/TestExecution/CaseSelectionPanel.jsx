import React, { useState, useEffect } from 'react';
import { Radio, Button, Space, Spin, message, Checkbox, Tree } from 'antd';
import { CheckOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getAutoCasesList, getWebCaseGroups, getApiCaseGroupsFromTable } from '../../../api/autoCase';
import { getCasesList, getCaseGroups } from '../../../api/manualCase';
import { getApiCasesList } from '../../../api/apiCase';

/**
 * 用例选择面板
 * 支持AI Web、AI接口、手工测试三种执行类型
 * @param {Object} props
 * @param {Object} props.task - 任务对象 {task_uuid, execution_type, project_id}
 * @param {number} props.projectId - 项目ID (从父组件传入)
 * @param {Function} props.onConfirm - 确认选择回调
 */
const CaseSelectionPanel = ({ task, projectId, onConfirm }) => {
  const { t } = useTranslation();
  // 通用状态
  const [languageType, setLanguageType] = useState('cn');
  const [selectedCaseGroup, setSelectedCaseGroup] = useState(null); // 选中的用例集
  const [caseGroups, setCaseGroups] = useState([]); // 用例集列表
  
  // 手工测试用例选择状态
  const [manualCaseType, setManualCaseType] = useState('overall'); // overall | acceptance | change
  
  const [loading, setLoading] = useState(false);
  const [loadingCaseGroups, setLoadingCaseGroups] = useState(false);

  // 将language代码转换为API期望的格式
  const mapLanguageCode = (code) => {
    const languageMap = {
      'cn': '中文',
      'jp': '日本語',
      'en': 'English'
    };
    return languageMap[code] || '中文';
  };

  // 手工测试用例类型映射
  const mapManualCaseType = (type) => {
    const typeMap = {
      'overall': 'overall',
      'acceptance': 'acceptance', // 受入用例使用独立的acceptance类型
      'change': 'change'
    };
    return typeMap[type] || 'overall';
  };

  // 加载用例集列表
  useEffect(() => {
    if (projectId && task) {
      loadCaseGroupsList();
    }
  }, [projectId, task?.execution_type]);

  // 加载用例集列表
  const loadCaseGroupsList = async () => {
    setLoadingCaseGroups(true);
    try {
      let groups = [];
      
      if (task.execution_type === 'manual') {
        // 手工测试：加载所有类型的用例集（overall, acceptance, change）
        console.log('🔵 [CaseSelectionPanel] Loading all manual case groups');
        const allGroups = [];
        const types = ['overall', 'acceptance', 'change'];
        
        for (const type of types) {
          try {
            const response = await getCaseGroups(projectId, type);
            if (response && response.length > 0) {
              allGroups.push(...response);
            }
          } catch (error) {
            console.warn(`Failed to load ${type} case groups:`, error);
          }
        }
        
        // 去重（根据group_name）
        const uniqueGroups = Array.from(
          new Map(allGroups.map(g => [g.group_name, g])).values()
        );
        groups = uniqueGroups;
        console.log('✅ [CaseSelectionPanel] Loaded manual case groups:', groups.length);
      } else if (task.execution_type === 'automation') {
        // AI Web：从 case_groups 表获取 web 类型用例集
        console.log('🔵 [CaseSelectionPanel] Loading web case groups');
        const response = await getWebCaseGroups(projectId);
        groups = response || [];
        console.log('✅ [CaseSelectionPanel] Loaded web case groups:', groups.length);
      } else if (task.execution_type === 'api') {
        // AI 接口：从 case_groups 表获取 api 类型用例集（与ApiLeftSidePanel保持一致）
        console.log('🔵 [CaseSelectionPanel] Loading API case groups');
        const response = await getApiCaseGroupsFromTable(projectId);
        // 转换数据格式：提取 group_name
        groups = (response || []).map(group => ({
          group_name: group.group_name,
          id: group.id
        }));
        console.log('✅ [CaseSelectionPanel] Loaded API case groups:', groups.length);
      }

      setCaseGroups(groups);
      
      // 默认选中第一个用例集
      if (groups.length > 0) {
        const firstGroup = groups[0].group_name || groups[0];
        setSelectedCaseGroup(firstGroup);
        console.log('✅ [CaseSelectionPanel] Default selected:', firstGroup);
      }
    } catch (error) {
      console.error('❌ [CaseSelectionPanel] Failed to load case groups:', error);
      message.error('加载用例集列表失败');
    } finally {
      setLoadingCaseGroups(false);
    }
  };

  // 处理AI Web/API确认按钮点击
  const handleConfirm = async () => {
    console.log('🔵 [CaseSelectionPanel] handleConfirm called');
    console.log('🔵 [CaseSelectionPanel] execution_type:', task.execution_type);
    console.log('🔵 [CaseSelectionPanel] selectedCaseGroup:', selectedCaseGroup);
    console.log('🔵 [CaseSelectionPanel] languageType:', languageType);
    
    if (!selectedCaseGroup) {
      message.warning(t('testExecution.messages.selectCaseGroup'));
      return;
    }

    setLoading(true);
    try {
      let cases = [];
      // 查找选中用例集的ID
      const selectedGroup = caseGroups.find(g => 
        (g.group_name || g) === selectedCaseGroup
      );
      const caseGroupId = selectedGroup?.id || 0;
      
      console.log('🔵 [CaseSelectionPanel] Selected group ID:', caseGroupId);
      
      let filterConditions = {
        execution_type: task.execution_type,
        case_group: selectedCaseGroup,
        case_group_id: caseGroupId  // 添加用例集ID
      };

      if (task.execution_type === 'automation') {
        // AI Web：需要语言参数
        const language = mapLanguageCode(languageType);
        console.log('🔵 [CaseSelectionPanel] Loading AI Web cases:', { selectedCaseGroup, language });
        
        const response = await getAutoCasesList(projectId, {
          caseType: 'web',
          language: language,
          caseGroup: selectedCaseGroup,
          page: 1,
          size: 99999
        });
        
        cases = response.cases || [];
        filterConditions.language = languageType;
        filterConditions.languageDisplay = language;
        console.log('✅ [CaseSelectionPanel] Loaded AI Web cases:', cases.length);
      } else if (task.execution_type === 'api') {
        // AI 接口：不需要语言参数
        console.log('🔵 [CaseSelectionPanel] Loading API cases:', { selectedCaseGroup });
        
        const response = await getApiCasesList(projectId, {
          caseType: 'api',
          caseGroup: selectedCaseGroup,
          page: 1,
          size: 99999
        });
        
        cases = response.cases || [];
        console.log('✅ [CaseSelectionPanel] Loaded API cases:', cases.length);
        console.log('🔍 [CaseSelectionPanel] API cases[0]:', cases[0]);
        console.log('🔍 [CaseSelectionPanel] API cases[0].script_code:', cases[0]?.script_code);
        console.log('🔍 [CaseSelectionPanel] API cases[0] keys:', cases[0] ? Object.keys(cases[0]) : 'empty');
      }

      if (cases.length === 0) {
        message.warning('所选用例集中没有用例');
        return;
      }

      if (onConfirm) {
        const resultData = {
          cases: cases,
          total: cases.length,
          filterConditions: filterConditions
        };
        console.log('🔵 [CaseSelectionPanel] Calling onConfirm with:', resultData);
        onConfirm(resultData);
        console.log('✅ [CaseSelectionPanel] onConfirm called successfully');
      }
    } catch (error) {
      console.error('❌ [CaseSelectionPanel] Failed to load cases:', error);
      message.error('加载用例失败: ' + (error.response?.data?.message || error.message));
    } finally {
      setLoading(false);
    }
  };

  // 处理手工测试用例确认按钮点击
  const handleManualConfirm = async () => {
    console.log('🔵 [CaseSelectionPanel] handleManualConfirm called');
    console.log('🔵 [CaseSelectionPanel] selectedCaseGroup:', selectedCaseGroup);

    if (!selectedCaseGroup) {
      message.warning(t('testExecution.messages.selectCaseGroup'));
      return;
    }

    setLoading(true);
    try {
      // 加载选中用例集的所有用例（使用中文，所有类型）
      console.log('🔵 [CaseSelectionPanel] Loading manual cases for case group:', selectedCaseGroup);
      
      // 尝试从所有用例类型中加载该用例集的用例
      let allCases = [];
      const types = ['overall', 'acceptance', 'change'];
      
      for (const type of types) {
        try {
          const response = await getCasesList(projectId, {
            caseType: type,
            language: '中文',
            caseGroup: selectedCaseGroup,
            page: 1,
            size: 99999
          });
          if (response.cases && response.cases.length > 0) {
            allCases.push(...response.cases);
          }
        } catch (error) {
          console.warn(`Failed to load ${type} cases:`, error);
        }
      }

      console.log('✅ [CaseSelectionPanel] Loaded manual cases:', allCases.length);

      if (allCases.length === 0) {
        message.warning('所选用例集中没有用例');
        return;
      }

      // 查找选中用例集的ID
      const selectedGroup = caseGroups.find(g => 
        (g.group_name || g) === selectedCaseGroup
      );
      const caseGroupId = selectedGroup?.id || 0;
      console.log('🔵 [CaseSelectionPanel] Manual selected group ID:', caseGroupId);

      if (onConfirm) {
        const resultData = {
          cases: allCases,
          total: allCases.length,
          filterConditions: {
            language: 'cn',
            languageDisplay: '中文',
            execution_type: 'manual',
            case_group: selectedCaseGroup,
            case_group_id: caseGroupId  // 添加用例集ID
          }
        };
        console.log('🔵 [CaseSelectionPanel] Calling onConfirm with:', resultData);
        onConfirm(resultData);
        console.log('✅ [CaseSelectionPanel] onConfirm called successfully');
      }
    } catch (error) {
      console.error('❌ [CaseSelectionPanel] handleManualConfirm failed:', error);
      message.error('加载用例失败: ' + (error.response?.data?.message || error.message));
    } finally {
      setLoading(false);
    }
  };

  if (!task) {
    return (
      <div style={{ textAlign: 'center', padding: '40px 0' }}>
        请先选择或创建一个测试任务
      </div>
    );
  }

  // 手工测试 (manual) 类型显示：用例集单选列表
  if (task.execution_type === 'manual') {
    return (
      <Spin spinning={loading || loadingCaseGroups}>
        <Space direction="vertical" size={16} style={{ width: '100%' }}>
          <div>
            <div style={{ marginBottom: 8, fontWeight: 'bold' }}>
              {t('testExecution.messages.selectCaseGroupLabel')}
              {caseGroups.length > 0 && (
                <span style={{ fontWeight: 'normal', marginLeft: 8, color: '#666' }}>
                  ({t('testExecution.messages.totalCaseGroups', { count: caseGroups.length })})
                </span>
              )}
            </div>
            {caseGroups.length > 0 ? (
              <Radio.Group 
                value={selectedCaseGroup} 
                onChange={e => setSelectedCaseGroup(e.target.value)}
                style={{ width: '100%' }}
              >
                <Space direction="vertical" style={{ width: '100%' }}>
                  {caseGroups.map((group) => (
                    <Radio 
                      key={group.id || group.group_name} 
                      value={group.group_name}
                      style={{ 
                        width: '100%',
                        padding: '8px 12px',
                        border: '1px solid #d9d9d9',
                        borderRadius: 4,
                        marginLeft: 0
                      }}
                    >
                      {group.group_name}
                    </Radio>
                  ))}
                </Space>
              </Radio.Group>
            ) : (
              <div style={{ textAlign: 'center', padding: '20px 0', color: '#999' }}>
                {loadingCaseGroups ? '正在加载...' : '暂无用例集'}
              </div>
            )}
          </div>

          <div style={{ textAlign: 'right', marginTop: 16 }}>
            <Button
              type="primary"
              icon={<CheckOutlined />}
              onClick={handleManualConfirm}
              loading={loading}
              disabled={!selectedCaseGroup}
            >
              {t('testExecution.messages.confirm')}
            </Button>
          </div>
        </Space>
      </Spin>
    );
  }

  // AI Web (automation) 类型显示：语言选择 + 用例集单选列表
  if (task.execution_type === 'automation') {
    return (
      <Spin spinning={loading || loadingCaseGroups}>
        <Space direction="vertical" size={16} style={{ width: '100%' }}>
          <div>
            <div style={{ marginBottom: 8, fontWeight: 'bold' }}>语言：</div>
            <Radio.Group value={languageType} onChange={e => setLanguageType(e.target.value)}>
              <Radio.Button value="cn">CN</Radio.Button>
              <Radio.Button value="jp">JP</Radio.Button>
              <Radio.Button value="en">EN</Radio.Button>
            </Radio.Group>
          </div>

          <div>
            <div style={{ marginBottom: 8, fontWeight: 'bold' }}>
              {t('testExecution.messages.selectCaseGroupLabel')}
              {caseGroups.length > 0 && (
                <span style={{ fontWeight: 'normal', marginLeft: 8, color: '#666' }}>
                  ({t('testExecution.messages.totalCaseGroups', { count: caseGroups.length })})
                </span>
              )}
            </div>
            {caseGroups.length > 0 ? (
              <Radio.Group 
                value={selectedCaseGroup} 
                onChange={e => setSelectedCaseGroup(e.target.value)}
                style={{ width: '100%' }}
              >
                <Space direction="vertical" style={{ width: '100%' }}>
                  {caseGroups.map((group) => (
                    <Radio 
                      key={group.id} 
                      value={group.group_name}
                      style={{ 
                        width: '100%',
                        padding: '8px 12px',
                        border: '1px solid #d9d9d9',
                        borderRadius: 4,
                        marginLeft: 0
                      }}
                    >
                      {group.group_name}
                    </Radio>
                  ))}
                </Space>
              </Radio.Group>
            ) : (
              <div style={{ textAlign: 'center', padding: '20px 0', color: '#999' }}>
                {loadingCaseGroups ? '正在加载...' : '暂无用例集'}
              </div>
            )}
          </div>

          <div style={{ textAlign: 'right', marginTop: 16 }}>
            <Button
              type="primary"
              icon={<CheckOutlined />}
              onClick={handleConfirm}
              loading={loading}
              disabled={!selectedCaseGroup}
            >
              {t('testExecution.messages.confirm')}
            </Button>
          </div>
        </Space>
      </Spin>
    );
  }

  // API (api) 类型显示：用例集单选列表
  if (task.execution_type === 'api') {
    return (
      <Spin spinning={loading || loadingCaseGroups}>
        <Space direction="vertical" size={16} style={{ width: '100%' }}>
          <div>
            <div style={{ marginBottom: 8, fontWeight: 'bold' }}>
              {t('testExecution.messages.selectCaseGroupLabel')}
              {caseGroups.length > 0 && (
                <span style={{ fontWeight: 'normal', marginLeft: 8, color: '#666' }}>
                  ({t('testExecution.messages.totalCaseGroups', { count: caseGroups.length })})
                </span>
              )}
            </div>
            {caseGroups.length > 0 ? (
              <Radio.Group 
                value={selectedCaseGroup} 
                onChange={e => setSelectedCaseGroup(e.target.value)}
                style={{ width: '100%' }}
              >
                <Space direction="vertical" style={{ width: '100%' }}>
                  {caseGroups.map((group) => (
                    <Radio 
                      key={group.id} 
                      value={group.group_name || group.id}
                      style={{ 
                        width: '100%',
                        padding: '8px 12px',
                        border: '1px solid #d9d9d9',
                        borderRadius: 4,
                        marginLeft: 0
                      }}
                    >
                      {group.group_name || group.id}
                    </Radio>
                  ))}
                </Space>
              </Radio.Group>
            ) : (
              <div style={{ textAlign: 'center', padding: '20px 0', color: '#999' }}>
                {loadingCaseGroups ? '正在加载...' : '暂无用例集'}
              </div>
            )}
          </div>

          <div style={{ textAlign: 'right', marginTop: 16 }}>
            <Button
              type="primary"
              icon={<CheckOutlined />}
              onClick={handleConfirm}
              loading={loading}
              disabled={!selectedCaseGroup}
            >
              {t('testExecution.messages.confirm')}
            </Button>
          </div>
        </Space>
      </Spin>
    );
  }

  // 其他类型暂不支持
  return (
    <div style={{ textAlign: 'center', padding: '40px 0' }}>
      当前执行类型暂不支持用例选择
    </div>
  );
};

export default CaseSelectionPanel;
