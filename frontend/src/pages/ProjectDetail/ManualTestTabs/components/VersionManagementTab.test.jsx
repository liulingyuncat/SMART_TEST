import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import VersionManagementTab from './VersionManagementTab';

// Mock Ant Design message
const mockMessage = {
  success: jest.fn(),
  error: jest.fn()
};

jest.mock('antd', () => {
  const actualAntd = jest.requireActual('antd');
  return {
    ...actualAntd,
    message: mockMessage
  };
});

describe('VersionManagementTab Component', () => {
  const mockAPI = {
    getVersionList: jest.fn(),
    downloadVersion: jest.fn(),
    deleteVersion: jest.fn()
  };

  const defaultProps = {
    projectId: '1',
    leftDocType: 'overall-requirements',
    rightDocType: 'overall-test-viewpoint',
    leftTitle: '整体需求版本管理',
    rightTitle: '整体测试观点版本管理',
    apiModule: mockAPI
  };

  const mockLeftVersions = [
    {
      id: 1,
      filename: '测试项目_整体需求_2025-11-17_100000.md',
      doc_type: 'overall-requirements',
      file_size: 1024,
      created_at: '2025-11-17T10:00:00Z'
    },
    {
      id: 2,
      filename: '测试项目_整体需求_2025-11-17_110000.md',
      doc_type: 'overall-requirements',
      file_size: 2048,
      created_at: '2025-11-17T11:00:00Z'
    }
  ];

  const mockRightVersions = [
    {
      id: 3,
      filename: '测试项目_整体测试观点_2025-11-17_100000.md',
      doc_type: 'overall-test-viewpoint',
      file_size: 1536,
      created_at: '2025-11-17T10:00:00Z'
    }
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    mockAPI.getVersionList.mockImplementation((projectId, docType) => {
      if (docType === 'overall-requirements') {
        return Promise.resolve(mockLeftVersions);
      } else if (docType === 'overall-test-viewpoint') {
        return Promise.resolve(mockRightVersions);
      }
      return Promise.resolve([]);
    });
    mockAPI.downloadVersion.mockResolvedValue();
    mockAPI.deleteVersion.mockResolvedValue({ message: 'success' });
  });

  test('组件挂载后加载左右两栏版本列表', async () => {
    render(<VersionManagementTab {...defaultProps} />);

    // 验证API被调用
    await waitFor(() => {
      expect(mockAPI.getVersionList).toHaveBeenCalledWith('1', 'overall-requirements');
      expect(mockAPI.getVersionList).toHaveBeenCalledWith('1', 'overall-test-viewpoint');
    });

    // 验证双栏标题
    expect(screen.getByText('整体需求版本管理')).toBeInTheDocument();
    expect(screen.getByText('整体测试观点版本管理')).toBeInTheDocument();
  });

  test('版本列表按创建时间倒序显示', async () => {
    render(<VersionManagementTab {...defaultProps} />);

    await waitFor(() => {
      expect(mockAPI.getVersionList).toHaveBeenCalled();
    });

    // TODO: 验证Table中数据顺序
    // 需要检查表格渲染的文件名顺序
    // 最新的版本(id=2)应该显示在最前面
  });

  test('点击下载按钮→调用downloadVersion API', async () => {
    render(<VersionManagementTab {...defaultProps} />);

    await waitFor(() => {
      expect(mockAPI.getVersionList).toHaveBeenCalled();
    });

    // TODO: 模拟点击下载按钮
    // 实际需要找到Table中的下载按钮并点击
    // fireEvent.click(downloadButton);

    // 验证API调用
    // expect(mockAPI.downloadVersion).toHaveBeenCalledWith('1', 1);
  });

  test('点击删除按钮→弹出确认对话框→确认后删除→刷新列表', async () => {
    render(<VersionManagementTab {...defaultProps} />);

    await waitFor(() => {
      expect(mockAPI.getVersionList).toHaveBeenCalled();
    });

    // TODO: 模拟删除操作
    // 1. 点击删除按钮
    // 2. 验证Popconfirm弹出
    // 3. 点击确认
    // 4. 验证deleteVersion被调用
    // 5. 验证getVersionList重新被调用以刷新列表
  });

  test('分页显示,每页10条记录', async () => {
    // 创建12条记录用于测试分页
    const manyVersions = Array.from({ length: 12 }, (_, i) => ({
      id: i + 1,
      filename: `测试项目_整体需求_2025-11-17_${String(100000 + i).padStart(6, '0')}.md`,
      doc_type: 'overall-requirements',
      file_size: 1024,
      created_at: `2025-11-17T${String(10 + Math.floor(i / 6)).padStart(2, '0')}:00:00Z`
    }));

    mockAPI.getVersionList.mockResolvedValue(manyVersions);

    render(<VersionManagementTab {...defaultProps} />);

    await waitFor(() => {
      expect(mockAPI.getVersionList).toHaveBeenCalled();
    });

    // TODO: 验证分页器显示
    // 应该显示"共 12 个版本"
    // 第一页显示10条记录
  });

  test('空版本列表显示提示信息', async () => {
    mockAPI.getVersionList.mockResolvedValue([]);

    render(<VersionManagementTab {...defaultProps} />);

    await waitFor(() => {
      expect(mockAPI.getVersionList).toHaveBeenCalled();
    });

    // 验证空状态提示
    expect(screen.getAllByText('暂无版本记录')).toHaveLength(2); // 左右两栏都显示
  });

  test('API调用失败时显示错误提示', async () => {
    mockAPI.getVersionList.mockRejectedValue(new Error('Network error'));

    render(<VersionManagementTab {...defaultProps} />);

    await waitFor(() => {
      expect(mockMessage.error).toHaveBeenCalledWith('加载版本列表失败', 3);
    });
  });

  test('删除成功后刷新对应栏位的版本列表', async () => {
    render(<VersionManagementTab {...defaultProps} />);

    await waitFor(() => {
      expect(mockAPI.getVersionList).toHaveBeenCalledTimes(2); // 初始加载
    });

    // 模拟删除左栏版本
    // TODO: 触发删除操作
    // await mockAPI.deleteVersion('1', 1);

    // 验证只刷新左栏列表
    // expect(mockAPI.getVersionList).toHaveBeenCalledWith('1', 'overall-requirements');
    // expect(mockAPI.getVersionList).not.toHaveBeenCalledWith('1', 'overall-test-viewpoint');
  });

  test('支持自定义apiModule', async () => {
    const customAPI = {
      getVersionList: jest.fn().mockResolvedValue([]),
      downloadVersion: jest.fn(),
      deleteVersion: jest.fn()
    };

    render(<VersionManagementTab {...defaultProps} apiModule={customAPI} />);

    // 验证使用自定义API
    await waitFor(() => {
      expect(customAPI.getVersionList).toHaveBeenCalled();
      expect(mockAPI.getVersionList).not.toHaveBeenCalled();
    });
  });

  test('变更版本管理Tab配置验证', async () => {
    const changeVersionProps = {
      projectId: '1',
      leftDocType: 'change-requirements',
      rightDocType: 'change-test-viewpoint',
      leftTitle: '变更需求版本管理',
      rightTitle: '变更测试观点版本管理',
      apiModule: mockAPI
    };

    render(<VersionManagementTab {...changeVersionProps} />);

    // 验证API调用参数
    await waitFor(() => {
      expect(mockAPI.getVersionList).toHaveBeenCalledWith('1', 'change-requirements');
      expect(mockAPI.getVersionList).toHaveBeenCalledWith('1', 'change-test-viewpoint');
    });

    // 验证标题
    expect(screen.getByText('变更需求版本管理')).toBeInTheDocument();
    expect(screen.getByText('变更测试观点版本管理')).toBeInTheDocument();
  });
});
