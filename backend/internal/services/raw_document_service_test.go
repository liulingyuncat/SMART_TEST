package services

import (
	"errors"
	"sync"
	"testing"
	"time"
	"webtest/internal/models"

	"gorm.io/gorm"
)

// MockRawDocumentRepository 模拟存储库用于测试
type MockRawDocumentRepository struct {
	mu   sync.RWMutex
	docs map[uint]*models.RawDocument
}

func NewMockRawDocumentRepository() *MockRawDocumentRepository {
	return &MockRawDocumentRepository{
		docs: make(map[uint]*models.RawDocument),
	}
}

func (m *MockRawDocumentRepository) Create(doc *models.RawDocument) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.docs[doc.ID] = doc
	return nil
}

func (m *MockRawDocumentRepository) GetByID(id uint) (*models.RawDocument, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if doc, exists := m.docs[id]; exists {
		// 返回副本以避免数据竞争
		docCopy := *doc
		return &docCopy, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockRawDocumentRepository) FindByID(id uint) (*models.RawDocument, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if doc, exists := m.docs[id]; exists {
		// 返回副本以避免数据竞争
		docCopy := *doc
		return &docCopy, nil
	}
	return nil, errors.New("document not found")
}

func (m *MockRawDocumentRepository) ListByProjectID(projectID uint) ([]*models.RawDocument, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*models.RawDocument
	for _, doc := range m.docs {
		if doc.ProjectID == projectID {
			// 返回副本以避免数据竞争
			docCopy := *doc
			result = append(result, &docCopy)
		}
	}
	return result, nil
}

func (m *MockRawDocumentRepository) Update(doc *models.RawDocument) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.docs[doc.ID] = doc
	return nil
}

func (m *MockRawDocumentRepository) UpdateStatus(id uint, status string, progress int, filename string, filepath string, filesize int64, convertError string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if doc, exists := m.docs[id]; exists {
		doc.ConvertStatus = status
		doc.ConvertProgress = progress
		doc.ConvertedFilename = filename
		doc.ConvertedFileSize = filesize
		if convertError != "" {
			doc.ConvertError = convertError
		}
		m.docs[id] = doc
		return nil
	}
	return errors.New("document not found")
}

func (m *MockRawDocumentRepository) GetConvertStatus(id uint) (*models.ConvertStatusResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if doc, exists := m.docs[id]; exists {
		return &models.ConvertStatusResponse{
			Status:            doc.ConvertStatus,
			Progress:          doc.ConvertProgress,
			ConvertedFilename: doc.ConvertedFilename,
			ErrorMessage:      doc.ConvertError,
		}, nil
	}
	return nil, errors.New("document not found")
}

func (m *MockRawDocumentRepository) Delete(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.docs, id)
	return nil
}

// 测试用例1: 测试正常转换启动
func TestStartConvert_Success(t *testing.T) {
	// 准备
	mockRepo := NewMockRawDocumentRepository()
	service := NewRawDocumentService(mockRepo, "/tmp/storage")

	// 创建测试文档
	doc := &models.RawDocument{
		ID:               1,
		ProjectID:        1,
		OriginalFilename: "test.pdf",
		OriginalFilepath: "raw_documents/1/original/test.pdf",
		FileSize:         1024,
		MimeType:         "application/pdf",
		UploadedBy:       1,
		ConvertStatus:    "none",
		ConvertProgress:  0,
	}
	mockRepo.Create(doc)

	// 执行
	result, err := service.StartConvert(1)

	// 验证
	if err != nil {
		t.Errorf("StartConvert failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.Status != "processing" {
		t.Errorf("Expected status 'processing', got '%s'", result.Status)
	}

	if result.TaskID == "" {
		t.Error("Expected non-empty task ID")
	}

	// 验证数据库状态
	updatedDoc, _ := mockRepo.GetByID(1)
	if updatedDoc.ConvertStatus != "processing" {
		t.Errorf("Expected database status 'processing', got '%s'", updatedDoc.ConvertStatus)
	}

	// 给异步goroutine一点时间完成以避免race detector警告
	time.Sleep(10 * time.Millisecond)
}

// 测试用例2: 测试重复转换被拒
func TestStartConvert_AlreadyInProgress(t *testing.T) {
	// 准备
	mockRepo := NewMockRawDocumentRepository()
	service := NewRawDocumentService(mockRepo, "/tmp/storage")

	// 创建处于转换中的文档
	doc := &models.RawDocument{
		ID:               2,
		ProjectID:        1,
		OriginalFilename: "test.doc",
		OriginalFilepath: "raw_documents/1/original/test.doc",
		FileSize:         2048,
		MimeType:         "application/msword",
		UploadedBy:       1,
		ConvertStatus:    "processing",
		ConvertProgress:  50,
		ConvertTaskID:    "convert_123_2",
	}
	mockRepo.Create(doc)

	// 执行
	result, err := service.StartConvert(2)

	// 验证
	if err == nil {
		t.Error("Expected error for already in-progress conversion")
	}

	if result != nil {
		t.Error("Expected nil result when conversion already in progress")
	}

	if err.Error() != "document conversion already in progress" {
		t.Errorf("Expected error message 'document conversion already in progress', got '%s'", err.Error())
	}
}

// 测试用例3: 测试转换状态查询准确
func TestGetConvertStatus_Accurate(t *testing.T) {
	// 准备
	mockRepo := NewMockRawDocumentRepository()
	service := NewRawDocumentService(mockRepo, "/tmp/storage")

	// 创建已完成转换的文档
	doc := &models.RawDocument{
		ID:                3,
		ProjectID:         1,
		OriginalFilename:  "test.xlsx",
		OriginalFilepath:  "raw_documents/1/original/test.xlsx",
		FileSize:          5120,
		MimeType:          "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		UploadedBy:        1,
		ConvertStatus:     "completed",
		ConvertProgress:   100,
		ConvertedFilename: "test_Trans_1702777200.md",
		ConvertedFileSize: 8192,
	}
	mockRepo.Create(doc)

	// 执行
	status, err := service.GetConvertStatus(3)

	// 验证
	if err != nil {
		t.Errorf("GetConvertStatus failed: %v", err)
	}

	if status == nil {
		t.Fatal("Expected non-nil status response")
	}

	if status.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", status.Status)
	}

	if status.Progress != 100 {
		t.Errorf("Expected progress 100, got %d", status.Progress)
	}

	if status.ConvertedFilename != "test_Trans_1702777200.md" {
		t.Errorf("Expected filename 'test_Trans_1702777200.md', got '%s'", status.ConvertedFilename)
	}
}

// 测试用例4: 测试转换失败状态
func TestGetConvertStatus_Failed(t *testing.T) {
	// 准备
	mockRepo := NewMockRawDocumentRepository()
	service := NewRawDocumentService(mockRepo, "/tmp/storage")

	// 创建转换失败的文档
	doc := &models.RawDocument{
		ID:               4,
		ProjectID:        1,
		OriginalFilename: "test.pdf",
		OriginalFilepath: "raw_documents/1/original/test.pdf",
		FileSize:         1024,
		MimeType:         "application/pdf",
		UploadedBy:       1,
		ConvertStatus:    "failed",
		ConvertProgress:  0,
		ConvertError:     "failed to extract text from PDF",
	}
	mockRepo.Create(doc)

	// 执行
	status, err := service.GetConvertStatus(4)

	// 验证
	if err != nil {
		t.Errorf("GetConvertStatus failed: %v", err)
	}

	if status.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", status.Status)
	}

	if status.ErrorMessage != "failed to extract text from PDF" {
		t.Errorf("Expected error message, got '%s'", status.ErrorMessage)
	}
}

// 测试用例5: 测试文档不存在
func TestStartConvert_DocumentNotFound(t *testing.T) {
	// 准备
	mockRepo := NewMockRawDocumentRepository()
	service := NewRawDocumentService(mockRepo, "/tmp/storage")

	// 执行 (使用不存在的ID)
	result, err := service.StartConvert(999)

	// 验证
	if err == nil {
		t.Error("Expected error for non-existent document")
	}

	if result != nil {
		t.Error("Expected nil result for non-existent document")
	}
}

// 测试文件名清理功能
func TestSanitizeFilename(t *testing.T) {
	mockRepo := NewMockRawDocumentRepository()
	service := NewRawDocumentService(mockRepo, "/tmp/storage")

	testCases := []struct {
		input    string
		expected string
	}{
		{"document.pdf", "document"},
		{"my/file\\name.docx", "my_file_name"},
		{"test:file*name.doc", "test_file_name"},
		{"normal_file.txt", "normal_file"},
		{"file?.pdf", "file_"},
	}

	for _, tc := range testCases {
		result := service.(*rawDocumentService).sanitizeFilename(tc.input)
		if result != tc.expected {
			t.Errorf("sanitizeFilename(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}
