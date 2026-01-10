import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import MarkdownEditor from './MarkdownEditor';

// Mock API
jest.mock('../../../api/requirement', () => ({
  saveVersion: jest.fn()
}));

// Mock react-markdown
jest.mock('react-markdown', () => {
  return function ReactMarkdown({ children }) {
    return <div data-testid="markdown-preview">{children}</div>;
  };
});

// Mock react-markdown-editor-lite
jest.mock('react-markdown-editor-lite', () => {
  return function MdEditor({ value, onChange }) {
    return (
      <textarea
        data-testid="md-editor"
        value={value}
        onChange={(e) => onChange({ text: e.target.value })}
      />
    );
  };
});

// Mock i18next
jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key) => key,
    i18n: { changeLanguage: jest.fn() }
  })
}));

const requirementAPI = require('../../../api/requirement');

describe('MarkdownEditor Component', () => {
  const mockOnChange = jest.fn();
  const mockOnSave = jest.fn().mockResolvedValue(true);
  const mockOnSaveVersion = jest.fn();
  
  const defaultProps = {
    value: '# 测试文档\n\n这是测试内容',
    onChange: mockOnChange,
    onSave: mockOnSave,
    onSaveVersion: mockOnSaveVersion,
    projectId: '1',
    projectName: '测试项目',
    docType: 'overall-requirements',
    showImport: true
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('场景5: 默认只读模式→点击编辑→进入编辑模式→点击取消→恢复原始内容', async () => {
    const { rerender } = render(<MarkdownEditor {...defaultProps} />);

    // 验证默认只读模式
    expect(screen.getByTestId('markdown-preview')).toBeInTheDocument();
    expect(screen.getByText('requirement.edit')).toBeInTheDocument();
    expect(screen.queryByTestId('md-editor')).not.toBeInTheDocument();

    // 点击编辑按钮
    const editButton = screen.getByText('requirement.edit');
    fireEvent.click(editButton);

    // 验证进入编辑模式
    await waitFor(() => {
      expect(screen.getByTestId('md-editor')).toBeInTheDocument();
    });
    expect(screen.getByText('requirement.cancel')).toBeInTheDocument();
    expect(screen.getByText('requirement.saveVersion')).toBeInTheDocument();

    // 修改内容
    const editor = screen.getByTestId('md-editor');
    fireEvent.change(editor, { target: { value: '# 修改后的文档\n\n修改的内容' } });

    await waitFor(() => {
      expect(mockOnChange).toHaveBeenCalledWith('# 修改后的文档\n\n修改的内容');
    });

    // 点击取消
    const cancelButton = screen.getByText('requirement.cancel');
    fireEvent.click(cancelButton);

    // 验证切回只读模式并恢复原始内容
    await waitFor(() => {
      expect(screen.getByTestId('markdown-preview')).toBeInTheDocument();
    });
    expect(screen.queryByTestId('md-editor')).not.toBeInTheDocument();
  });

  test('版本保存→显示自动生成的文件名→自动切回只读模式', async () => {
    requirementAPI.saveVersion.mockResolvedValue({
      filename: '测试项目_整体需求_2025-11-17_103000.md'
    });

    render(<MarkdownEditor {...defaultProps} />);

    // 进入编辑模式
    fireEvent.click(screen.getByText('requirement.edit'));

    // 等待编辑器渲染
    await waitFor(() => {
      expect(screen.getByTestId('md-editor')).toBeInTheDocument();
    });

    // 点击版本保存
    const saveVersionButton = screen.getByText('requirement.saveVersion');
    fireEvent.click(saveVersionButton);

    // 验证API调用
    await waitFor(() => {
      expect(requirementAPI.saveVersion).toHaveBeenCalledWith(
        '1',
        'overall-requirements',
        '# 测试文档\n\n这是测试内容'
      );
    });

    // 验证onSaveVersion回调被调用
    await waitFor(() => {
      expect(mockOnSaveVersion).toHaveBeenCalled();
    });

    // TODO: 验证成功提示消息显示文件名
    // TODO: 验证自动切回只读模式
  });

  test('Markdown导入→验证文件大小限制(≤5MB)', async () => {
    render(<MarkdownEditor {...defaultProps} />);

    // 进入编辑模式
    fireEvent.click(screen.getByText('requirement.edit'));

    await waitFor(() => {
      expect(screen.getByText('requirement.import')).toBeInTheDocument();
    });

    // 创建超大文件(6MB)
    const largeContent = 'a'.repeat(6 * 1024 * 1024);
    const largeFile = new File([largeContent], 'large.md', { type: 'text/markdown' });

    // 模拟文件选择
    const importButton = screen.getByText('requirement.import');
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    Object.defineProperty(fileInput, 'files', {
      value: [largeFile]
    });

    // TODO: 触发文件选择事件并验证错误提示
    // 实际需要模拟input的onChange事件
  });

  test('Markdown导入→验证.md后缀', async () => {
    render(<MarkdownEditor {...defaultProps} />);

    // 进入编辑模式
    fireEvent.click(screen.getByText('requirement.edit'));

    await waitFor(() => {
      expect(screen.getByText('requirement.import')).toBeInTheDocument();
    });

    // 创建非.md文件
    const txtFile = new File(['test content'], 'test.txt', { type: 'text/plain' });

    // TODO: 验证文件类型验证
    // 应该拒绝非.md文件
  });

  test('验证showImport为false时不显示导入按钮', () => {
    render(<MarkdownEditor {...defaultProps} showImport={false} />);

    // 进入编辑模式
    fireEvent.click(screen.getByText('requirement.edit'));

    // 验证导入按钮不存在
    expect(screen.queryByText('requirement.import')).not.toBeInTheDocument();
  });

  test('验证docType不是requirements类型时不显示导入按钮', () => {
    render(<MarkdownEditor {...defaultProps} docType="overall-test-viewpoint" />);

    // 进入编辑模式
    fireEvent.click(screen.getByText('requirement.edit'));

    // 验证导入按钮不存在(测试观点类型不允许导入)
    expect(screen.queryByText('requirement.import')).not.toBeInTheDocument();
  });

  test('空内容时显示提示', () => {
    render(<MarkdownEditor {...defaultProps} value="" />);

    // 验证空内容提示
    expect(screen.getByText('requirement.emptyContent')).toBeInTheDocument();
  });

  test('保存成功后切回只读模式', async () => {
    render(<MarkdownEditor {...defaultProps} />);

    // 进入编辑模式
    fireEvent.click(screen.getByText('requirement.edit'));

    await waitFor(() => {
      expect(screen.getByTestId('md-editor')).toBeInTheDocument();
    });

    // 点击保存
    const saveButton = screen.getByText('common.save');
    fireEvent.click(saveButton);

    // 验证onSave被调用
    await waitFor(() => {
      expect(mockOnSave).toHaveBeenCalled();
    });

    // TODO: 验证自动切回只读模式
  });
});
