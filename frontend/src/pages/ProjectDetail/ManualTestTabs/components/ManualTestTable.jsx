import React from 'react';
import { useTranslation } from 'react-i18next';
import { Table, Empty } from 'antd';

// 用例表格列标题翻译映射（独立于UI语言）
const COLUMN_TITLES = {
  '中文': {
    caseNumber: '用例编号',
    majorFunction: '大功能',
    middleFunction: '中功能',
    minorFunction: '小功能',
    precondition: '前提条件',
    testSteps: '测试步骤',
    expectedResult: '期待值',
    testResult: '测试结果',
    remark: '备考',
    noData: '暂无该语言版本的用例',
    total: '共',
    items: '条',
  },
  'English': {
    caseNumber: 'Case Number',
    majorFunction: 'Major Function',
    middleFunction: 'Middle Function',
    minorFunction: 'Minor Function',
    precondition: 'Precondition',
    testSteps: 'Test Steps',
    expectedResult: 'Expected Result',
    testResult: 'Test Result',
    remark: 'Remark',
    noData: 'No cases in this language',
    total: 'Total',
    items: 'items',
  },
  '日本語': {
    caseNumber: 'ケース番号',
    majorFunction: '大機能',
    middleFunction: '中機能',
    minorFunction: '小機能',
    precondition: '前提条件',
    testSteps: 'テスト手順',
    expectedResult: '期待値',
    testResult: 'テスト結果',
    remark: '備考',
    noData: 'この言語のケースはありません',
    total: '合計',
    items: '件',
  },
};

/**
 * 测试用例表格组件
 * @param {Object} props
 * @param {Array} props.cases - 用例数据数组
 * @param {number} props.total - 总条数
 * @param {number} props.page - 当前页码
 * @param {number} props.pageSize - 每页条数(固定50)
 * @param {string} props.language - 当前语言
 * @param {Function} props.onPageChange - 分页回调
 */
const ManualTestTable = ({ 
  cases = [], 
  total = 0, 
  page = 1, 
  pageSize = 50,
  language = '中文',
  onPageChange 
}) => {
  // 获取当前语言的列标题
  const titles = COLUMN_TITLES[language] || COLUMN_TITLES['中文'];

  // 列配置
  const fullColumns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
      fixed: 'left',
    },
    {
      title: titles.caseNumber,
      dataIndex: 'case_number',
      key: 'case_number',
      width: 200,
    },
    {
      title: titles.majorFunction,
      dataIndex: 'major_function',
      key: 'major_function',
      width: 220,
    },
    {
      title: titles.middleFunction,
      dataIndex: 'middle_function',
      key: 'middle_function',
      width: 220,
    },
    {
      title: titles.minorFunction,
      dataIndex: 'minor_function',
      key: 'minor_function',
      width: 220,
    },
    {
      title: titles.precondition,
      dataIndex: 'precondition',
      key: 'precondition',
      width: 200,
      ellipsis: true,
    },
    {
      title: titles.testSteps,
      dataIndex: 'test_steps',
      key: 'test_steps',
      width: 250,
      ellipsis: true,
    },
    {
      title: titles.expectedResult,
      dataIndex: 'expected_result',
      key: 'expected_result',
      width: 240,
      ellipsis: true,
    },
    {
      title: titles.testResult,
      dataIndex: 'test_result',
      key: 'test_result',
      width: 180,
      render: (text) => text || 'NR',
    },
    {
      title: titles.remark,
      dataIndex: 'remark',
      key: 'remark',
      width: 150,
      ellipsis: true,
    },
  ];

  // 始终使用完整列配置
  const columns = fullColumns;

  // 分页配置
  const pagination = {
    current: page,
    pageSize: pageSize,
    total: total,
    showSizeChanger: false,
    showTotal: (total) => `${titles.total} ${total} ${titles.items}`,
    onChange: onPageChange,
  };

  return (
    <Table
      columns={columns}
      dataSource={cases}
      rowKey="id"
      pagination={pagination}
      scroll={{ x: 'max-content' }}
      locale={{
        emptyText: <Empty description={titles.noData} />,
      }}
    />
  );
};

export default ManualTestTable;
