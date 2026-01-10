import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import RawDocumentTab from '../index';
import * as rawDocumentApi from '../../../../api/rawDocument';

// Mock i18n
jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key) => {
      const translations = {
        'rawDocument.convert': '转换',
        'rawDocument.converting': '转换中...',
        'rawDocument.convertStarted': '转换已启动',
        'rawDocument.convertSuccess': '转换成功',
        'rawDocument.convertFailed': '转换失败',
        'rawDocument.convertTimeout': '转换超时',
        'rawDocument.statusCheckFailed': '状态检查失败',
        'rawDocument.documentNotFound': '文档不存在',
        'rawDocument.convertInProgress': '文档正在转换中',
        'rawDocument.loadFailed': '加载失败',
        'rawDocument.originalDocument': '原始文档',
        'rawDocument.fileSize': '文件大小',
        'rawDocument.uploadTime': '上传时间',
        'rawDocument.originalActions': '操作',
        'rawDocument.convertedDocument': '转换文档',
        'rawDocument.statusCompleted': '已完成',
        'rawDocument.convertToMarkdown': '转换为Markdown',
        'rawDocument.convertingInProgress': '转换中，请稍候...',
      };
      return translations[key] || key;
    },
  }),
}));

// Mock Ant Design message
const mockMessage = {
  success: jest.fn(),
  error: jest.fn(),
  warning: jest.fn(),
};
jest.mock('antd', () => {
  const antd = jest.requireActual('antd');
  return {
    ...antd,
    message: mockMessage,
  };
});

// Mock API
jest.mock('../../../../api/rawDocument');

describe('RawDocumentTab', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  const mockDocuments = [
    {
      id: 1,
      project_id: 1,
      original_filename: 'test.pdf',
      file_size: 1024,
      mime_type: 'application/pdf',
      uploaded_by: 1,
      convert_status: 'none',
      convert_progress: 0,
      created_at: '2025-12-17T10:00:00Z',
    },
    {
      id: 2,
      project_id: 1,
      original_filename: 'processing.docx',
      file_size: 2048,
      mime_type: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      uploaded_by: 1,
      convert_status: 'processing',
      convert_progress: 50,
      created_at: '2025-12-17T10:01:00Z',
    },
    {
      id: 3,
      project_id: 1,
      original_filename: 'completed.txt',
      file_size: 512,
      mime_type: 'text/plain',
      uploaded_by: 1,
      convert_status: 'completed',
      convert_progress: 100,
      converted_filename: 'completed_Trans_1702777200.md',
      created_at: '2025-12-17T10:02:00Z',
    },
  ];

  describe('handleConvert', () => {
    it('should call convertRawDocument API when convert button is clicked', async () => {
      rawDocumentApi.fetchRawDocuments.mockResolvedValue(mockDocuments);
      rawDocumentApi.convertRawDocument.mockResolvedValue({ task_id: 'test-task', status: 'processing' });
      rawDocumentApi.getConvertStatus.mockResolvedValue({ status: 'completed', progress: 100 });

      render(<RawDocumentTab projectId={1} />);

      await waitFor(() => {
        expect(screen.getByText('test.pdf')).toBeInTheDocument();
      });

      // Find and click the convert button for the first document
      const convertButtons = screen.getAllByText('转换');
      fireEvent.click(convertButtons[0]);

      await waitFor(() => {
        expect(rawDocumentApi.convertRawDocument).toHaveBeenCalledWith(1);
      });
    });

    it('should show error message when API returns 404', async () => {
      rawDocumentApi.fetchRawDocuments.mockResolvedValue(mockDocuments);
      rawDocumentApi.convertRawDocument.mockRejectedValue({
        response: { status: 404 },
        message: 'Not Found',
      });

      render(<RawDocumentTab projectId={1} />);

      await waitFor(() => {
        expect(screen.getByText('test.pdf')).toBeInTheDocument();
      });

      const convertButtons = screen.getAllByText('转换');
      fireEvent.click(convertButtons[0]);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('文档不存在');
      });
    });

    it('should show error message when API returns 409 (already converting)', async () => {
      rawDocumentApi.fetchRawDocuments.mockResolvedValue(mockDocuments);
      rawDocumentApi.convertRawDocument.mockRejectedValue({
        response: { status: 409 },
        message: 'Conflict',
      });

      render(<RawDocumentTab projectId={1} />);

      await waitFor(() => {
        expect(screen.getByText('test.pdf')).toBeInTheDocument();
      });

      const convertButtons = screen.getAllByText('转换');
      fireEvent.click(convertButtons[0]);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('文档正在转换中');
      });
    });
  });

  describe('checkConvertStatus', () => {
    it('should poll status until completed', async () => {
      rawDocumentApi.fetchRawDocuments.mockResolvedValue([mockDocuments[0]]);
      rawDocumentApi.convertRawDocument.mockResolvedValue({ task_id: 'test-task', status: 'processing' });
      
      // First poll: processing, second poll: completed
      rawDocumentApi.getConvertStatus
        .mockResolvedValueOnce({ status: 'processing', progress: 50 })
        .mockResolvedValueOnce({ status: 'completed', progress: 100 });

      render(<RawDocumentTab projectId={1} />);

      await waitFor(() => {
        expect(screen.getByText('test.pdf')).toBeInTheDocument();
      });

      const convertButtons = screen.getAllByText('转换');
      fireEvent.click(convertButtons[0]);

      // Wait for convert API call
      await waitFor(() => {
        expect(rawDocumentApi.convertRawDocument).toHaveBeenCalled();
      });

      // First poll
      await waitFor(() => {
        expect(rawDocumentApi.getConvertStatus).toHaveBeenCalled();
      });

      // Advance timer for second poll
      jest.advanceTimersByTime(2000);

      await waitFor(() => {
        expect(rawDocumentApi.getConvertStatus).toHaveBeenCalledTimes(2);
      });

      await waitFor(() => {
        expect(mockMessage.success).toHaveBeenCalledWith('转换成功');
      });
    });

    it('should show error when conversion fails', async () => {
      rawDocumentApi.fetchRawDocuments.mockResolvedValue([mockDocuments[0]]);
      rawDocumentApi.convertRawDocument.mockResolvedValue({ task_id: 'test-task', status: 'processing' });
      rawDocumentApi.getConvertStatus.mockResolvedValue({ 
        status: 'failed', 
        progress: 0, 
        error_message: 'Text extraction failed' 
      });

      render(<RawDocumentTab projectId={1} />);

      await waitFor(() => {
        expect(screen.getByText('test.pdf')).toBeInTheDocument();
      });

      const convertButtons = screen.getAllByText('转换');
      fireEvent.click(convertButtons[0]);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('Text extraction failed');
      });
    });

    it('should retry polling on network error', async () => {
      rawDocumentApi.fetchRawDocuments.mockResolvedValue([mockDocuments[0]]);
      rawDocumentApi.convertRawDocument.mockResolvedValue({ task_id: 'test-task', status: 'processing' });
      
      // First poll: error, second retry: completed
      rawDocumentApi.getConvertStatus
        .mockRejectedValueOnce(new Error('Network error'))
        .mockResolvedValueOnce({ status: 'completed', progress: 100 });

      render(<RawDocumentTab projectId={1} />);

      await waitFor(() => {
        expect(screen.getByText('test.pdf')).toBeInTheDocument();
      });

      const convertButtons = screen.getAllByText('转换');
      fireEvent.click(convertButtons[0]);

      await waitFor(() => {
        expect(rawDocumentApi.getConvertStatus).toHaveBeenCalled();
      });

      // Advance timer for retry
      jest.advanceTimersByTime(2000);

      await waitFor(() => {
        expect(rawDocumentApi.getConvertStatus).toHaveBeenCalledTimes(2);
      });

      await waitFor(() => {
        expect(mockMessage.success).toHaveBeenCalledWith('转换成功');
      });
    });
  });

  describe('Button rendering', () => {
    it('should display convert button for documents with status "none"', async () => {
      rawDocumentApi.fetchRawDocuments.mockResolvedValue([mockDocuments[0]]);

      render(<RawDocumentTab projectId={1} />);

      await waitFor(() => {
        expect(screen.getByText('转换')).toBeInTheDocument();
      });
    });

    it('should display "转换中..." for documents with status "processing"', async () => {
      rawDocumentApi.fetchRawDocuments.mockResolvedValue([mockDocuments[1]]);

      render(<RawDocumentTab projectId={1} />);

      await waitFor(() => {
        expect(screen.getByText('转换中...')).toBeInTheDocument();
      });
    });

    it('should display "已完成" tag for documents with status "completed"', async () => {
      rawDocumentApi.fetchRawDocuments.mockResolvedValue([mockDocuments[2]]);

      render(<RawDocumentTab projectId={1} />);

      await waitFor(() => {
        expect(screen.getByText('已完成')).toBeInTheDocument();
      });
    });

    it('should disable button when status is "processing"', async () => {
      rawDocumentApi.fetchRawDocuments.mockResolvedValue([mockDocuments[1]]);

      render(<RawDocumentTab projectId={1} />);

      await waitFor(() => {
        const processingButton = screen.getByText('转换中...').closest('button');
        expect(processingButton).toBeDisabled();
      });
    });
  });
});
