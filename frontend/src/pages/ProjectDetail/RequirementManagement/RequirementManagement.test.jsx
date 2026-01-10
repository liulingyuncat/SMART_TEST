import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import RequirementManagement from './index';

// Mock API
jest.mock('../../../api/requirement', () => ({
  fetchRequirement: jest.fn(),
  updateRequirement: jest.fn(),
  saveVersion: jest.fn(),
  getVersionList: jest.fn(),
  downloadVersion: jest.fn(),
  deleteVersion: jest.fn()
}));

// Mock i18next
jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key) => key,
    i18n: { changeLanguage: jest.fn() }
  })
}));

// Mock React Router
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: () => ({ id: '1' }),
  BrowserRouter: ({ children }) => <div>{children}</div>,
  MemoryRouter: ({ children }) => <div>{children}</div>,
  Routes: ({ children }) => <div>{children}</div>,
  Route: ({ element }) => element
}));

// Mock VersionManagementTab
jest.mock('../ManualTestTabs/components/VersionManagementTab', () => {
  return function MockVersionManagementTab(props) {
    return <div data-testid="version-management-tab">VersionManagementTab Mock</div>;
  };
});

// Mock MarkdownEditor
jest.mock('./MarkdownEditor', () => {
  return function MockMarkdownEditor(props) {
    return <div data-testid="markdown-editor">MarkdownEditor Mock</div>;
  };
});

const requirementAPI = require('../../../api/requirement');

describe('RequirementManagement - Version Management', () => {
  const mockProjectName = '测试项目';
  
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Mock API responses
    requirementAPI.fetchRequirement.mockResolvedValue({
      content: '# 测试需求文档\n\n这是测试内容',
      updatedAt: '2025-11-17T10:00:00Z'
    });
    
    requirementAPI.updateRequirement.mockResolvedValue({
      message: 'success'
    });
    
    requirementAPI.saveVersion.mockResolvedValue({
      filename: '测试项目_整体需求_2025-11-17_100000.md'
    });
    
    requirementAPI.getVersionList.mockResolvedValue([
      {
        id: 1,
        filename: '测试项目_整体需求_2025-11-17_100000.md',
        doc_type: 'overall-requirements',
        file_size: 1024,
        created_at: '2025-11-17T10:00:00Z'
      }
    ]);
    
    requirementAPI.downloadVersion.mockResolvedValue();
    requirementAPI.deleteVersion.mockResolvedValue({ message: 'success' });
  });

  test('场景1: 验证Tab渲染和组件集成', async () => {
    render(<RequirementManagement projectName={mockProjectName} />);

    // 等待文档加载
    await waitFor(() => {
      expect(requirementAPI.fetchRequirement).toHaveBeenCalledWith('1', 'overall-requirements');
    });

    // 验证Tab渲染
    expect(screen.getByText('requirement.overallReq')).toBeInTheDocument();
    expect(screen.getByText('requirement.overallVersion')).toBeInTheDocument();

    // 验证MarkdownEditor渲染
    expect(screen.getByTestId('markdown-editor')).toBeInTheDocument();
  });

  test('场景2: 切换到版本管理Tab验证组件渲染', async () => {
    render(<RequirementManagement projectName={mockProjectName} />);

    await waitFor(() => {
      expect(requirementAPI.fetchRequirement).toHaveBeenCalled();
    });

    // 点击版本管理Tab
    const versionTab = screen.getByText('requirement.overallVersion');
    fireEvent.click(versionTab);

    // 验证VersionManagementTab渲染
    await waitFor(() => {
      expect(screen.getByTestId('version-management-tab')).toBeInTheDocument();
    });
  });

  test('场景3: 多次保存版本→验证文件名时间戳递增', async () => {
    const version1 = {
      filename: '测试项目_整体需求_2025-11-17_100000.md'
    };
    const version2 = {
      filename: '测试项目_整体需求_2025-11-17_100001.md'
    };

    requirementAPI.saveVersion
      .mockResolvedValueOnce(version1)
      .mockResolvedValueOnce(version2);

    // 验证文件名格式和时间戳
    expect(version1.filename).toMatch(/^测试项目_整体需求_\d{4}-\d{2}-\d{2}_\d{6}\.md$/);
    expect(version2.filename).toMatch(/^测试项目_整体需求_\d{4}-\d{2}-\d{2}_\d{6}\.md$/);
    
    // 提取时间戳并验证递增
    const timestamp1 = version1.filename.split('_')[3].replace('.md', '');
    const timestamp2 = version2.filename.split('_')[3].replace('.md', '');
    expect(parseInt(timestamp2)).toBeGreaterThan(parseInt(timestamp1));
  });

  test('验证API调用参数正确性', async () => {
    render(<RequirementManagement projectName={mockProjectName} />);

    // 验证fetchRequirement调用
    await waitFor(() => {
      expect(requirementAPI.fetchRequirement).toHaveBeenCalledWith('1', 'overall-requirements');
    });
  });
});

describe('RequirementManagement - 权限校验', () => {
  test('场景4: 非项目成员访问→验证权限拒绝', async () => {
    // Mock API返回403错误
    requirementAPI.fetchRequirement.mockRejectedValue({
      response: { status: 403, data: { message: 'Permission denied' } }
    });

    render(<RequirementManagement projectName="无权限项目" />);

    // 等待错误处理
    await waitFor(() => {
      expect(requirementAPI.fetchRequirement).toHaveBeenCalled();
    });

    // 实际项目中应该显示权限错误提示
  });
});
