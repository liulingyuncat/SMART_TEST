import React, { useState, useEffect } from 'react';
import { Card, Radio, Select, Button, Tree, Form, Space, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { getCasesList } from '../../../api/manualCase';

/**
 * 筛选面板组件
 * @param {Object} props
 * @param {string} props.taskExecutionType - 任务执行类型 (manual/automation/api)
 * @param {Function} props.onFilterChange - 筛选条件变更回调
 * @param {number} props.projectId - 项目ID
 * @param {string} props.defaultCaseType - 默认用例类型
 * @param {string} props.defaultLanguage - 默认语言
 */
const FilterPanel = ({ taskExecutionType, onFilterChange, projectId, defaultCaseType, defaultLanguage }) => {
  const { t } = useTranslation();
  const [caseType, setCaseType] = useState(defaultCaseType || 'overall');
  const [roleType, setRoleType] = useState(defaultCaseType || 'role1');
  const [languageType, setLanguageType] = useState(defaultLanguage || 'cn');
  const [selectedFunctions, setSelectedFunctions] = useState([]);
  const [functionTree, setFunctionTree] = useState([]);
  const [loading, setLoading] = useState(false);

  // 当默认值变化时更新状态
  useEffect(() => {
    if (defaultCaseType) {
      if (taskExecutionType === 'manual') {
        setCaseType(defaultCaseType);
      } else {
        setRoleType(defaultCaseType);
      }
    }
    if (defaultLanguage) {
      setLanguageType(defaultLanguage);
    }
  }, [defaultCaseType, defaultLanguage, taskExecutionType]);

  // 加载功能树(manual类型)
  useEffect(() => {
    if (taskExecutionType === 'manual') {
      loadFunctionTree();
    }
  }, [taskExecutionType, caseType, projectId]);

  // 加载大功能/中功能树
  const loadFunctionTree = async () => {
    if (!projectId) return;

    setLoading(true);
    try {
      const response = await getCasesList(projectId, caseType, 1, 999999);
      const cases = response.cases || [];

      // 聚合大功能和中功能
      const functionMap = {};
      cases.forEach(c => {
        const major = c.major_function_cn || '';
        const middle = c.middle_function_cn || '';
        
        if (!major) return;

        if (!functionMap[major]) {
          functionMap[major] = new Set();
        }
        if (middle) {
          functionMap[major].add(middle);
        }
      });

      // 转换为Tree数据结构
      const treeData = Object.keys(functionMap).map((major, index) => ({
        key: `major-${index}`,
        title: major,
        value: major,
        children: Array.from(functionMap[major]).map((middle, midIndex) => ({
          key: `middle-${index}-${midIndex}`,
          title: middle,
          value: `${major}/${middle}`,
        }))
      }));

      setFunctionTree(treeData);
    } catch (error) {
      console.error('Load function tree failed:', error);
      message.error('加载功能分类失败');
    } finally {
      setLoading(false);
    }
  };

  // 处理确认按钮点击
  const handleConfirm = () => {
    const conditions = {
      case_type: taskExecutionType === 'manual' ? caseType : roleType,
      language: taskExecutionType === 'automation' ? languageType : 'cn',
      functions: selectedFunctions
    };

    console.log('[FilterPanel] Confirm conditions:', conditions);
    onFilterChange(conditions);
  };

  // 渲染Manual类型筛选器
  const renderManualFilter = () => (
    <Space direction="vertical" size={12} style={{ width: '100%' }}>
      <Form.Item label="用例类型">
        <Radio.Group value={caseType} onChange={e => setCaseType(e.target.value)}>
          <Radio.Button value="overall">整体用例</Radio.Button>
          <Radio.Button value="acceptance">受入用例</Radio.Button>
          <Radio.Button value="change">变更用例</Radio.Button>
        </Radio.Group>
      </Form.Item>

      <Form.Item label="功能分类">
        <Tree
          checkable
          treeData={functionTree}
          checkedKeys={selectedFunctions}
          onCheck={setSelectedFunctions}
          style={{ maxHeight: 300, overflow: 'auto' }}
        />
      </Form.Item>

      <Button type="primary" onClick={handleConfirm} loading={loading}>
        确认
      </Button>
    </Space>
  );

  // 渲染Automation类型筛选器
  const renderAutomationFilter = () => (
    <Space direction="vertical" size={12} style={{ width: '100%' }}>
      <Form.Item label="ROLE类型">
        <Radio.Group value={roleType} onChange={e => setRoleType(e.target.value)}>
          <Radio.Button value="role1">ROLE1</Radio.Button>
          <Radio.Button value="role2">ROLE2</Radio.Button>
          <Radio.Button value="role3">ROLE3</Radio.Button>
          <Radio.Button value="role4">ROLE4</Radio.Button>
        </Radio.Group>
      </Form.Item>

      <Form.Item label="语言">
        <Radio.Group value={languageType} onChange={e => setLanguageType(e.target.value)}>
          <Radio.Button value="cn">CN</Radio.Button>
          <Radio.Button value="jp">JP</Radio.Button>
          <Radio.Button value="en">EN</Radio.Button>
        </Radio.Group>
      </Form.Item>

      <Button type="primary" onClick={handleConfirm}>
        确认
      </Button>
    </Space>
  );

  // 渲染API类型筛选器
  const renderApiFilter = () => (
    <Space direction="vertical" size={12} style={{ width: '100%' }}>
      <Form.Item label="ROLE类型">
        <Radio.Group value={roleType} onChange={e => setRoleType(e.target.value)}>
          <Radio.Button value="role1">ROLE1</Radio.Button>
          <Radio.Button value="role2">ROLE2</Radio.Button>
          <Radio.Button value="role3">ROLE3</Radio.Button>
          <Radio.Button value="role4">ROLE4</Radio.Button>
        </Radio.Group>
      </Form.Item>

      <Button type="primary" onClick={handleConfirm}>
        确认
      </Button>
    </Space>
  );

  return (
    <Card title="筛选条件" size="small">
      {taskExecutionType === 'manual' && renderManualFilter()}
      {taskExecutionType === 'automation' && renderAutomationFilter()}
      {taskExecutionType === 'api' && renderApiFilter()}
    </Card>
  );
};

export default FilterPanel;
