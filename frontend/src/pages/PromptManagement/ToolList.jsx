import React, { useState, useEffect } from 'react';
import { Button, message, Spin, Typography, Tooltip, Collapse } from 'antd';
import { CopyOutlined, CaretRightOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import axios from 'axios';

const { Text } = Typography;
const { Panel } = Collapse;

const ToolList = () => {
  const { t } = useTranslation();
  const [toolCategories, setToolCategories] = useState([]);
  const [loading, setLoading] = useState(false);
  const [expandedKeys, setExpandedKeys] = useState([]);

  useEffect(() => {
    loadTools();
  }, []);

  const loadTools = async () => {
    setLoading(true);
    try {
      // 按类别组织的MCP工具列表
      const categorizedTools = [
        {
          key: 'project',
          title: t('prompts.categoryProject'),
          icon: '📁',
          tools: [
            { name: 'get_current_project_name', description: t('prompts.toolDescriptions.get_current_project_name'), params: '无', returns: 'project_id, project_name' },
          ]
        },
        {
          key: 'documents',
          title: t('prompts.categoryDocuments'),
          icon: '📄',
          tools: [
            { name: 'list_raw_documents', description: t('prompts.toolDescriptions.list_raw_documents'), params: 'project_id(required)', returns: '原始文档列表' },
            { name: 'get_converted_document', description: t('prompts.toolDescriptions.get_converted_document'), params: 'project_id(required), document_id(required)', returns: '完整的文档内容' },
          ]
        },
        {
          key: 'requirements',
          title: t('prompts.categoryRequirements'),
          icon: '📋',
          tools: [
            { name: 'list_requirement_items', description: t('prompts.toolDescriptions.list_requirement_items'), params: 'project_id(required)', returns: '需求文档列表' },
            { name: 'get_requirement_item', description: t('prompts.toolDescriptions.get_requirement_item'), params: 'project_id(required), requirement_id(required)', returns: '完整的需求文档内容' },
            { name: 'create_requirement_item', description: t('prompts.toolDescriptions.create_requirement_item'), params: 'project_id(required), name(required), content(required), parent_id(optional)', returns: '新创建的需求文档ID和详细信息' },
            { name: 'update_requirement_item', description: t('prompts.toolDescriptions.update_requirement_item'), params: 'project_id(required), requirement_id(required), name(optional), content(optional)', returns: '更新后的需求文档信息' },
          ]
        },
        {
          key: 'viewpoints',
          title: t('prompts.categoryViewpoints'),
          icon: '👁️',
          tools: [
            { name: 'list_viewpoint_items', description: t('prompts.toolDescriptions.list_viewpoint_items'), params: 'project_id(required)', returns: '观点文档列表' },
            { name: 'get_viewpoint_item', description: t('prompts.toolDescriptions.get_viewpoint_item'), params: 'project_id(required), viewpoint_id(required)', returns: '完整的观点文档内容' },
            { name: 'create_viewpoint_item', description: t('prompts.toolDescriptions.create_viewpoint_item'), params: 'project_id(required), name(required), content(required), requirement_id(optional)', returns: '新创建的观点文档ID和详细信息' },
            { name: 'update_viewpoint_item', description: t('prompts.toolDescriptions.update_viewpoint_item'), params: 'project_id(required), viewpoint_id(required), name(optional), content(optional)', returns: '更新后的观点文档信息' },
          ]
        },
        {
          key: 'manual',
          title: t('prompts.categoryManual'),
          icon: '✋',
          tools: [
            { name: 'list_manual_groups', description: t('prompts.toolDescriptions.list_manual_groups'), params: 'project_id(required)', returns: '用例集列表' },
            { name: 'list_manual_cases', description: t('prompts.toolDescriptions.list_manual_cases'), params: 'project_id(required), case_group_id(required), all_fields(optional)', returns: '用例列表' },
            { name: 'create_case_group', description: t('prompts.toolDescriptions.create_case_group'), params: 'project_id(required), name(required), type(optional), description(optional)', returns: '新创建的用例集ID和详细信息' },
            { name: 'create_manual_cases', description: t('prompts.toolDescriptions.create_manual_cases'), params: 'project_id(required), case_group_id(required), cases(required), continue_on_error(optional)', returns: '创建结果列表' },
            { name: 'update_manual_case', description: t('prompts.toolDescriptions.update_manual_case'), params: 'project_id(required), case_id(required), data(required)', returns: '更新后的用例信息' },
            { name: 'update_manual_cases', description: t('prompts.toolDescriptions.update_manual_cases'), params: 'project_id(required), case_group_id(optional), cases(required), continue_on_error(optional)', returns: '更新结果列表' },
          ]
        },
        {
          key: 'web',
          title: t('prompts.categoryWeb'),
          icon: '🌐',
          tools: [
            { name: 'list_web_groups', description: t('prompts.toolDescriptions.list_web_groups'), params: 'project_id(required)', returns: 'Web用例集列表' },
            { name: 'get_web_group_metadata', description: t('prompts.toolDescriptions.get_web_group_metadata'), params: 'project_id(required), group_id(required)', returns: '用例集元数据' },
            { name: 'list_web_cases', description: t('prompts.toolDescriptions.list_web_cases'), params: 'project_id(required), group_id(required)', returns: 'Web用例列表' },
            { name: 'create_web_cases', description: t('prompts.toolDescriptions.create_web_cases'), params: 'project_id(required), group_id(required), cases(required), continue_on_error(optional)', returns: '创建结果列表' },
            { name: 'update_web_cases', description: t('prompts.toolDescriptions.update_web_cases'), params: 'project_id(required), cases(required), continue_on_error(optional)', returns: '更新结果列表' },
          ]
        },
        {
          key: 'api',
          title: t('prompts.categoryApi'),
          icon: '🔌',
          tools: [
            { name: 'list_api_groups', description: t('prompts.toolDescriptions.list_api_groups'), params: 'project_id(required)', returns: 'API用例集列表' },
            { name: 'get_api_group_metadata', description: t('prompts.toolDescriptions.get_api_group_metadata'), params: 'project_id(required), group_id(required)', returns: '用例集元数据' },
            { name: 'list_api_cases', description: t('prompts.toolDescriptions.list_api_cases'), params: 'project_id(required), group_id(required)', returns: 'API用例列表' },
            { name: 'create_api_case', description: t('prompts.toolDescriptions.create_api_case'), params: 'project_id(required), group_id(required), case_data(required)', returns: '新创建的用例ID和详细信息' },
            { name: 'update_api_case', description: t('prompts.toolDescriptions.update_api_case'), params: 'project_id(required), case_id(required), case_data(required)', returns: '更新后的用例信息' },
          ]
        },
        {
          key: 'review',
          title: t('prompts.categoryReview'),
          icon: '✅',
          tools: [
            { name: 'create_review_item', description: t('prompts.toolDescriptions.create_review_item'), params: 'project_id(required), name(required), content(optional)', returns: '新创建的评审条目ID和详细信息' },
          ]
        },
        {
          key: 'execution',
          title: t('prompts.categoryExecution'),
          icon: '▶️',
          tools: [
            { name: 'list_execution_tasks', description: t('prompts.toolDescriptions.list_execution_tasks'), params: 'project_id(required)', returns: '执行任务列表' },
            { name: 'get_execution_task_metadata', description: t('prompts.toolDescriptions.get_execution_task_metadata'), params: 'project_id(required), task_id(required)', returns: '任务元数据和执行统计' },
            { name: 'get_execution_task_cases', description: t('prompts.toolDescriptions.get_execution_task_cases'), params: 'project_id(required), task_id(required)', returns: '用例列表及其执行结果' },
            { name: 'update_execution_case_result', description: t('prompts.toolDescriptions.update_execution_case_result'), params: 'project_id(required), case_id(required), result(required), comment(optional)', returns: '更新后的结果信息' },
          ]
        },
        {
          key: 'defects',
          title: t('prompts.categoryDefects'),
          icon: '🐛',
          tools: [
            { name: 'list_defects', description: t('prompts.toolDescriptions.list_defects'), params: 'project_id(required), page(optional), page_size(optional)', returns: '缺陷列表和总数' },
            { name: 'update_defect', description: t('prompts.toolDescriptions.update_defect'), params: 'project_id(required), defect_id(required), status(optional), comment(optional)', returns: '更新后的缺陷信息' },
          ]
        },
        {
          key: 'reports',
          title: t('prompts.categoryReports'),
          icon: '📊',
          tools: [
            { name: 'create_ai_report', description: t('prompts.toolDescriptions.create_ai_report'), params: 'project_id(required), title(required), content(required)', returns: '新创建的报告ID和详细信息' },
            { name: 'update_ai_report', description: t('prompts.toolDescriptions.update_ai_report'), params: 'project_id(required), report_id(optional), report_name(optional), content(optional), new_name(optional)', returns: '更新后的报告信息' },
          ]
        },
      ];
      setToolCategories(categorizedTools);
      // 默认展开所有分类
      setExpandedKeys(categorizedTools.map(c => c.key));
    } catch (error) {
      console.error('Failed to load tools:', error);
      message.error(t('prompts.loadToolsFailed'));
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = (toolName) => {
    navigator.clipboard.writeText(toolName).then(() => {
      message.success(t('prompts.copySuccess') || '复制成功');
    }).catch(() => {
      message.error(t('prompts.copyFailed') || '复制失败');
    });
  };

  // 生成工具Tooltip内容
  const renderToolTooltip = (tool) => {
    return (
      <div style={{ textAlign: 'left', maxWidth: '420px' }}>
        <div style={{ marginBottom: '8px', fontWeight: 500, color: '#ffffff' }}>
          {tool.description}
        </div>
        <div style={{ marginBottom: '6px', fontSize: '12px', color: '#f0f0f0' }}>
          <strong>{t('prompts.toolParams')}：</strong> <span style={{ color: '#ffc53d' }}>{tool.params}</span>
        </div>
        <div style={{ fontSize: '12px', color: '#f0f0f0' }}>
          <strong>{t('prompts.toolReturns')}：</strong> <span style={{ color: '#95de64' }}>{tool.returns}</span>
        </div>
      </div>
    );
  };

  // 计算总工具数
  const totalCount = toolCategories.reduce((sum, cat) => sum + cat.tools.length, 0);

  // 渲染单个工具项
  const renderToolItem = (tool) => (
    <div
      key={tool.name}
      className="tool-item"
      style={{
        padding: '8px 16px 8px 24px',
        display: 'flex',
        alignItems: 'center',
        gap: '8px',
        cursor: 'default',
        transition: 'background-color 0.2s ease',
        borderBottom: '1px solid #f5f5f5',
      }}
      onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#f5f5f5'}
      onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
    >
      <Tooltip 
        title={renderToolTooltip(tool)} 
        placement="right"
        overlayStyle={{ maxWidth: '500px' }}
      >
        <Text 
          style={{ 
            fontSize: '12px', 
            color: '#1890ff',
            minWidth: '220px',
            flexShrink: 0,
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
            fontWeight: 500,
            cursor: 'help',
          }}
        >
          {tool.name}
        </Text>
      </Tooltip>
      <Text 
        style={{ 
          fontSize: '11px', 
          color: '#8c8c8c',
          flex: 1,
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          whiteSpace: 'nowrap',
        }}
      >
        {tool.description}
      </Text>
      <Button
        type="text"
        size="small"
        icon={<CopyOutlined />}
        onClick={() => handleCopy(tool.name)}
        style={{ 
          padding: '4px 8px',
          height: '24px',
          minWidth: '24px',
          flexShrink: 0,
        }}
        className="copy-btn"
        title={t('prompts.copyToolName')}
      />
    </div>
  );

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '16px' }}>
        <Spin size="small" />
      </div>
    );
  }

  return (
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column', background: '#fff' }}>
      {/* 头部 */}
      <div style={{ 
        padding: '12px 16px',
        borderBottom: '1px solid #f0f0f0',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        background: '#fff',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <div style={{ fontSize: '14px', fontWeight: 600, color: '#262626' }}>
            {t('prompts.mcpToolList')}
          </div>
          <Text type="secondary" style={{ fontSize: '12px', fontWeight: 400 }}>
            {t('prompts.toolCount', { count: totalCount })}
          </Text>
        </div>
        <div style={{ display: 'flex', gap: '8px' }}>
          <Button 
            type="link" 
            size="small" 
            onClick={() => setExpandedKeys(toolCategories.map(c => c.key))}
            style={{ fontSize: '12px', padding: '0 4px' }}
          >
            {t('prompts.expandAll')}
          </Button>
          <Button 
            type="link" 
            size="small" 
            onClick={() => setExpandedKeys([])}
            style={{ fontSize: '12px', padding: '0 4px' }}
          >
            {t('prompts.collapseAll')}
          </Button>
        </div>
      </div>

      {/* 分类工具列表 */}
      <div style={{ flex: 1, overflow: 'auto' }}>
        <Collapse 
          activeKey={expandedKeys}
          onChange={(keys) => setExpandedKeys(keys)}
          ghost
          expandIcon={({ isActive }) => (
            <CaretRightOutlined rotate={isActive ? 90 : 0} style={{ fontSize: '10px' }} />
          )}
          style={{ background: '#fff' }}
        >
          {toolCategories.map((category) => (
            <Panel 
              key={category.key}
              header={
                <div style={{ 
                  display: 'flex', 
                  alignItems: 'center', 
                  gap: '8px',
                  fontWeight: 500,
                  fontSize: '13px',
                  color: '#262626',
                }}>
                  <span>{category.icon}</span>
                  <span>{category.title}</span>
                  <Text type="secondary" style={{ fontSize: '12px', fontWeight: 400 }}>
                    （{category.tools.length}个）
                  </Text>
                </div>
              }
              style={{ 
                borderBottom: '1px solid #f0f0f0',
                marginBottom: 0,
              }}
            >
              <div style={{ marginLeft: '-12px', marginRight: '-12px' }}>
                {category.tools.map(tool => renderToolItem(tool))}
              </div>
            </Panel>
          ))}
        </Collapse>
      </div>
    </div>
  );
};

export default ToolList;
