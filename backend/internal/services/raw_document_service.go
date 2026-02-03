package services

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
	"webtest/internal/models"
	"webtest/internal/repositories"

	ledongpdf "github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"rsc.io/pdf"

	"gorm.io/gorm"
)

// RawDocumentService 原始文档服务接口
type RawDocumentService interface {
	Upload(projectID, userID uint, file *multipart.FileHeader) (*models.RawDocumentUploadResponse, error)
	List(projectID uint) ([]*models.RawDocumentListItem, error)
	StartConvert(id uint) (*models.ConvertTaskResponse, error)
	GetConvertStatus(id uint) (*models.ConvertStatusResponse, error)
	DownloadOriginal(id uint) (*models.RawDocument, io.ReadCloser, error)
	DownloadConverted(id uint) (*models.RawDocument, io.ReadCloser, error)
	PreviewConverted(id uint) (*models.RawDocument, string, error)
	DeleteOriginal(id uint) error
	DeleteConverted(id uint) error
}

type rawDocumentService struct {
	repo            repositories.RawDocumentRepository
	storageBasePath string
}

// NewRawDocumentService 创建原始文档服务实例
func NewRawDocumentService(repo repositories.RawDocumentRepository, storageBasePath string) RawDocumentService {
	return &rawDocumentService{
		repo:            repo,
		storageBasePath: storageBasePath,
	}
}

// Upload 上传原始文档
func (s *rawDocumentService) Upload(projectID, userID uint, file *multipart.FileHeader) (*models.RawDocumentUploadResponse, error) {
	// 验证文件大小
	if file.Size > models.MaxRawDocumentSize {
		return nil, fmt.Errorf("file size exceeds limit of %d bytes", models.MaxRawDocumentSize)
	}

	// 获取MIME类型
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 验证文件类型（白名单校验）
	if !models.IsRawDocumentAllowedMimeType(mimeType) {
		return nil, errors.New("file type not allowed")
	}

	// 构建存储路径: storage/raw_documents/{projectID}/original/{id}_{timestamp}_{filename}
	timestamp := time.Now().Format("20060102150405")
	fileName := filepath.Base(file.Filename)

	// 先创建临时记录以获取ID
	doc := &models.RawDocument{
		ProjectID:        projectID,
		OriginalFilename: fileName,
		OriginalFilepath: "", // 稍后更新
		FileSize:         file.Size,
		MimeType:         mimeType,
		UploadedBy:       userID,
		ConvertStatus:    "none",
		ConvertProgress:  0,
	}

	if err := s.repo.Create(doc); err != nil {
		return nil, fmt.Errorf("create raw document record: %w", err)
	}

	// 生成完整存储路径
	storageName := fmt.Sprintf("%d_%s_%s", doc.ID, timestamp, fileName)
	relPath := filepath.Join("raw_documents", fmt.Sprintf("%d", projectID), "original")
	fullDir := filepath.Join(s.storageBasePath, relPath)
	fullPath := filepath.Join(fullDir, storageName)

	// 创建目录
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("create destination file: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		// 清理已创建的文件
		os.Remove(fullPath)
		return nil, fmt.Errorf("copy file content: %w", err)
	}

	// 更新文件路径
	doc.OriginalFilepath = filepath.Join(relPath, storageName)
	if err := s.repo.Update(doc); err != nil {
		// 清理已创建的文件
		os.Remove(fullPath)
		return nil, fmt.Errorf("update document filepath: %w", err)
	}

	log.Printf("[RawDocument Upload] project_id=%d, file=%s, size=%d, id=%d", projectID, fileName, file.Size, doc.ID)

	return &models.RawDocumentUploadResponse{
		ID:               doc.ID,
		OriginalFilename: doc.OriginalFilename,
		FileSize:         doc.FileSize,
		MimeType:         doc.MimeType,
		UploadTime:       doc.CreatedAt,
	}, nil
}

// List 获取项目的原始文档列表
func (s *rawDocumentService) List(projectID uint) ([]*models.RawDocumentListItem, error) {
	docs, err := s.repo.ListByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("list raw documents: %w", err)
	}

	items := make([]*models.RawDocumentListItem, 0, len(docs))
	for _, doc := range docs {
		items = append(items, &models.RawDocumentListItem{
			ID:                doc.ID,
			ProjectID:         doc.ProjectID,
			OriginalFilename:  doc.OriginalFilename,
			FileSize:          doc.FileSize,
			MimeType:          doc.MimeType,
			UploadedBy:        doc.UploadedBy,
			ConvertStatus:     doc.ConvertStatus,
			ConvertProgress:   doc.ConvertProgress,
			ConvertedFilename: doc.ConvertedFilename,
			ConvertedFileSize: doc.ConvertedFileSize,
			ConvertedTime:     doc.ConvertedTime,
			ConvertError:      doc.ConvertError,
			CreatedAt:         doc.CreatedAt,
		})
	}

	return items, nil
}

// StartConvert 启动文档转换任务（异步处理）
//
// 流程:
// 1. 从数据库查询文档信息
// 2. 验证文档存在且不在转换中（防止重复转换）
// 3. 生成唯一的转换任务ID
// 4. 更新数据库状态为 processing
// 5. 启动异步 goroutine 执行转换（带60秒超时）
//
// 参数:
//   - id: 文档ID
//
// 返回:
//   - *ConvertTaskResponse: 包含任务ID和状态
//   - error: 文档不存在或已在转换中时返回错误
func (s *rawDocumentService) StartConvert(id uint) (*models.ConvertTaskResponse, error) {
	doc, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// 检查是否已在转换中
	if doc.ConvertStatus == "processing" {
		return nil, errors.New("document conversion already in progress")
	}

	// 生成转换任务ID
	taskID := fmt.Sprintf("convert_%d_%d", time.Now().Unix(), doc.ID)

	// 更新转换状态为 processing
	doc.ConvertStatus = "processing"
	doc.ConvertTaskID = taskID
	doc.ConvertProgress = 0
	doc.ConvertError = ""

	if err := s.repo.Update(doc); err != nil {
		return nil, fmt.Errorf("update convert status: %w", err)
	}

	log.Printf("[Convert Request] documentId=%d, taskId=%s, originalFilename=%s", id, taskID, doc.OriginalFilename)

	// 启动异步转换goroutine（不阻塞HTTP响应）
	// 注意：context 的 cancel 需要在 goroutine 内部处理，不能使用 defer cancel()
	// 因为 defer 会在函数返回时立即执行，导致 context 被过早取消
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		s.convertDocumentAsyncWithContext(ctx, id, taskID)
	}()

	return &models.ConvertTaskResponse{
		TaskID: taskID,
		Status: "processing",
	}, nil
}

// convertDocumentAsyncWithContext 异步执行文档转换（带超时控制）
//
// 使用 select 模式监听转换完成和超时信号:
// - 转换完成: 正常退出
// - 超时(60秒): 更新状态为 failed 并记录错误
//
// 参数:
//   - ctx: 带超时的上下文
//   - documentID: 文档ID
//   - taskID: 转换任务ID（用于日志追踪）
func (s *rawDocumentService) convertDocumentAsyncWithContext(ctx context.Context, documentID uint, taskID string) {
	// 创建一个完成通道
	done := make(chan struct{})

	go func() {
		s.convertDocumentAsync(documentID, taskID)
		close(done)
	}()

	// 等待转换完成或超时
	select {
	case <-done:
		// 转换成功完成
		return
	case <-ctx.Done():
		// 超时，更新状态为失败
		log.Printf("[Convert Timeout] documentId=%d, taskId=%s, exceeded 60 seconds", documentID, taskID)
		s.updateConvertStatus(documentID, "failed", 0, "", "", 0, "conversion timeout exceeded 60 seconds")
		return
	}
}

// convertDocumentAsync 异步执行文档转换（在goroutine中运行）
func (s *rawDocumentService) convertDocumentAsync(documentID uint, taskID string) {
	log.Printf("[Convert Async Start] documentId=%d, taskId=%s", documentID, taskID)

	doc, err := s.repo.FindByID(documentID)
	if err != nil {
		log.Printf("[Convert Failed] documentId=%d, error: document not found", documentID)
		s.updateConvertStatus(documentID, "failed", 0, "", "", 0, "document not found")
		return
	}

	// 构建原始文件的完整路径
	fullPath := filepath.Join(s.storageBasePath, doc.OriginalFilepath)

	// 读取原始文件
	content, err := os.ReadFile(fullPath)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read file: %v", err)
		log.Printf("[Convert Failed] documentId=%d, error: %s", documentID, errMsg)
		s.updateConvertStatus(documentID, "failed", 0, "", "", 0, errMsg)
		return
	}

	// 执行文本提取和Markdown生成
	log.Printf("[Convert Processing] documentId=%d, fileSize=%d", documentID, len(content))

	// 简单的转换实现：将文件内容转换为Markdown格式
	markdownContent := s.convertToMarkdown(doc, content)

	// 生成转换后的文件名：{原始名}_Trans_{时间戳}.md
	sanitized := s.sanitizeFilename(doc.OriginalFilename)
	timestamp := time.Now().Unix()
	convertedFilename := fmt.Sprintf("%s_Trans_%d.md", sanitized, timestamp)
	log.Printf("[Convert Filename] documentId=%d, original=%s, sanitized=%s, converted=%s", documentID, doc.OriginalFilename, sanitized, convertedFilename)

	// 构建存储路径: storage/raw_documents/{projectID}/converted/{filename}
	convertedDir := filepath.Join(s.storageBasePath, "raw_documents", fmt.Sprintf("%d", doc.ProjectID), "converted")
	if err := os.MkdirAll(convertedDir, 0755); err != nil {
		errMsg := fmt.Sprintf("failed to create converted directory: %v", err)
		log.Printf("[Convert Failed] documentId=%d, error: %s", documentID, errMsg)
		s.updateConvertStatus(documentID, "failed", 0, "", "", 0, errMsg)
		return
	}

	convertedFilepath := filepath.Join(convertedDir, convertedFilename)

	// 保存转换结果
	if err := os.WriteFile(convertedFilepath, []byte(markdownContent), 0644); err != nil {
		errMsg := fmt.Sprintf("failed to save converted file: %v", err)
		log.Printf("[Convert Failed] documentId=%d, error: %s", documentID, errMsg)
		s.updateConvertStatus(documentID, "failed", 0, "", "", 0, errMsg)
		return
	}

	// 获取转换后文件大小
	fileInfo, _ := os.Stat(convertedFilepath)
	convertedFileSize := fileInfo.Size()

	// 计算相对路径用于存储到数据库
	relativeFilepath := filepath.Join("raw_documents", fmt.Sprintf("%d", doc.ProjectID), "converted", convertedFilename)

	// 更新数据库状态为 completed
	log.Printf("[Convert Success] documentId=%d, taskId=%s, convertedFilename=%s, filepath=%s, fileSize=%d", documentID, taskID, convertedFilename, relativeFilepath, convertedFileSize)
	s.updateConvertStatus(documentID, "completed", 100, convertedFilename, relativeFilepath, convertedFileSize, "")
}

// convertToMarkdown 将文件内容转换为Markdown格式
// 根据文件类型智能处理：
// - 图片文件：转换为base64嵌入的Markdown图片
// - PDF文件：提取文本内容
// - Word文件(DOCX)：提取文本内容
// - Excel文件(XLSX/CSV)：转换为Markdown表格
// - 文本文件：智能检测编码并转换为UTF-8
func (s *rawDocumentService) convertToMarkdown(doc *models.RawDocument, content []byte) string {
	mimeType := strings.ToLower(doc.MimeType)
	ext := strings.ToLower(filepath.Ext(doc.OriginalFilename))

	// 检查是否为图片类型
	if s.isImageMimeType(mimeType) {
		return s.convertImageToMarkdown(doc, content)
	}

	// 检查是否为PDF
	if mimeType == "application/pdf" || ext == ".pdf" {
		return s.convertPDFToMarkdown(doc, content)
	}

	// 检查是否为Word文档
	if s.isWordDocument(mimeType, ext) {
		return s.convertWordToMarkdown(doc, content)
	}

	// 检查是否为Excel文档
	if s.isExcelDocument(mimeType, ext) {
		return s.convertExcelToMarkdown(doc, content)
	}

	// 检查是否为PowerPoint文档
	if s.isPowerPointDocument(mimeType, ext) {
		return s.convertPowerPointToMarkdown(doc, content)
	}

	// 文本类文件，进行编码检测和转换
	textContent := s.convertToUTF8(content)

	backtick := "`"
	codeBlock := backtick + backtick + backtick

	markdownContent := fmt.Sprintf("# %s\n\n## 文档信息\n- **文件名**: %s\n- **文件大小**: %d 字节\n- **MIME类型**: %s\n- **转换时间**: %s\n\n## 文档内容\n\n%s\n%s\n%s\n\n---\n*本文档由自动转换工具生成*\n",
		doc.OriginalFilename,
		doc.OriginalFilename,
		doc.FileSize,
		doc.MimeType,
		time.Now().Format("2006-01-02 15:04:05"),
		codeBlock,
		textContent,
		codeBlock,
	)

	return markdownContent
}

// isWordDocument 检查是否为Word文档
func (s *rawDocumentService) isWordDocument(mimeType, ext string) bool {
	wordMimeTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/msword": true,
	}
	wordExts := map[string]bool{
		".docx": true,
		".doc":  true,
	}
	return wordMimeTypes[mimeType] || wordExts[ext]
}

// isExcelDocument 检查是否为Excel文档
func (s *rawDocumentService) isExcelDocument(mimeType, ext string) bool {
	excelMimeTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
		"application/vnd.ms-excel": true,
		"text/csv":                 true,
	}
	excelExts := map[string]bool{
		".xlsx": true,
		".xls":  true,
		".csv":  true,
	}
	return excelMimeTypes[mimeType] || excelExts[ext]
}

// isPowerPointDocument 检查是否为PowerPoint文档
func (s *rawDocumentService) isPowerPointDocument(mimeType, ext string) bool {
	pptMimeTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
		"application/vnd.ms-powerpoint":                                             true,
	}
	pptExts := map[string]bool{
		".pptx": true,
		".ppt":  true,
	}
	return pptMimeTypes[mimeType] || pptExts[ext]
}

// convertPDFToMarkdown 将PDF转换为Markdown
func (s *rawDocumentService) convertPDFToMarkdown(doc *models.RawDocument, content []byte) string {
	// 使用defer recover捕获PDF库可能的panic
	var result string
	var panicErr interface{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = r
				log.Printf("[PDF Convert] Panic recovered: %v", r)
			}
		}()
		result = s.convertPDFToMarkdownInternal(doc, content)
	}()

	if panicErr != nil {
		return s.createErrorMarkdown(doc, fmt.Sprintf("PDF转换失败: 文档格式不兼容 (%v)", panicErr))
	}

	return result
}

// convertPDFToMarkdownInternal 实际的PDF转换逻辑
// 优先使用pdftotext命令行工具（对CJK支持更好），如果不可用则回退到rsc.io/pdf库
func (s *rawDocumentService) convertPDFToMarkdownInternal(doc *models.RawDocument, content []byte) string {
	// 创建临时文件用于PDF解析
	tmpFile, err := os.CreateTemp("", "pdf_*.pdf")
	if err != nil {
		log.Printf("[PDF Convert] Failed to create temp file: %v", err)
		return s.createErrorMarkdown(doc, "PDF转换失败: 无法创建临时文件")
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	if _, err := tmpFile.Write(content); err != nil {
		tmpFile.Close()
		log.Printf("[PDF Convert] Failed to write temp file: %v", err)
		return s.createErrorMarkdown(doc, "PDF转换失败: 无法写入临时文件")
	}
	tmpFile.Close()

	// 优先尝试使用pdftotext（poppler-utils），对CJK支持最好
	extractedText := s.extractPDFWithPdftotext(tmpFileName)
	if extractedText != "" {
		log.Printf("[PDF Convert] Successfully extracted text using pdftotext")
		return s.createDocumentMarkdown(doc, extractedText, "PDF文档")
	}

	// 方案2: 使用ledongthuc/pdf（纯Go，Unicode/CJK支持较好）
	extractedText = s.extractPDFWithLedongthuc(tmpFileName)
	if extractedText != "" {
		log.Printf("[PDF Convert] Successfully extracted text using ledongthuc/pdf")
		return s.createDocumentMarkdown(doc, extractedText, "PDF文档")
	}

	// 方案3: 使用pdfcpu纯Go库
	extractedText = s.extractPDFWithPdfcpu(tmpFileName)
	if extractedText != "" {
		log.Printf("[PDF Convert] Successfully extracted text using pdfcpu")
		return s.createDocumentMarkdown(doc, extractedText, "PDF文档")
	}

	// 方案4: 回退到rsc.io/pdf库
	log.Printf("[PDF Convert] All methods failed, falling back to rsc.io/pdf")
	return s.extractPDFWithRscPdf(doc, tmpFileName)
}

// extractPDFWithPdftotext 使用pdftotext命令行工具提取PDF文本
// pdftotext是poppler-utils的一部分，对CJK语言支持很好
func (s *rawDocumentService) extractPDFWithPdftotext(pdfPath string) string {
	// 检查pdftotext是否可用
	_, err := exec.LookPath("pdftotext")
	if err != nil {
		log.Printf("[PDF Convert] pdftotext not found in PATH")
		return ""
	}

	// 创建临时输出文件
	tmpOutput, err := os.CreateTemp("", "pdftext_*.txt")
	if err != nil {
		log.Printf("[PDF Convert] Failed to create temp output file: %v", err)
		return ""
	}
	tmpOutputPath := tmpOutput.Name()
	tmpOutput.Close()
	defer os.Remove(tmpOutputPath)

	// 执行pdftotext命令
	// -layout: 保持原始布局
	// -enc UTF-8: 输出UTF-8编码
	cmd := exec.Command("pdftotext", "-layout", "-enc", "UTF-8", pdfPath, tmpOutputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[PDF Convert] pdftotext failed: %v, output: %s", err, string(output))
		return ""
	}

	// 读取提取的文本
	textContent, err := os.ReadFile(tmpOutputPath)
	if err != nil {
		log.Printf("[PDF Convert] Failed to read pdftotext output: %v", err)
		return ""
	}

	result := strings.TrimSpace(string(textContent))
	if result == "" {
		return ""
	}

	// 格式化输出，添加页面标识
	return s.formatPdftotextOutput(result)
}

// formatPdftotextOutput 格式化pdftotext输出
func (s *rawDocumentService) formatPdftotextOutput(text string) string {
	// pdftotext使用form feed字符(\f)分隔页面
	pages := strings.Split(text, "\f")

	var result strings.Builder
	pageNum := 0

	for _, page := range pages {
		page = strings.TrimSpace(page)
		if page == "" {
			continue
		}

		pageNum++
		result.WriteString(fmt.Sprintf("### 第 %d 页\n\n", pageNum))
		result.WriteString(page)
		result.WriteString("\n\n")
	}

	return strings.TrimSpace(result.String())
}

// extractPDFWithLedongthuc 使用ledongthuc/pdf库提取PDF文本
// 该库对Unicode和CJK字符有更好的支持
func (s *rawDocumentService) extractPDFWithLedongthuc(pdfPath string) string {
	// 打开PDF文件
	f, r, err := ledongpdf.Open(pdfPath)
	if err != nil {
		log.Printf("[PDF Convert] ledongthuc/pdf: Failed to open file: %v", err)
		return ""
	}
	defer f.Close()

	totalPages := r.NumPage()
	if totalPages == 0 {
		log.Printf("[PDF Convert] ledongthuc/pdf: No pages found")
		return ""
	}

	log.Printf("[PDF Convert] ledongthuc/pdf: Processing %d pages", totalPages)

	var result strings.Builder
	hasContent := false

	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		// 提取页面文本
		pageText, err := page.GetPlainText(nil)
		if err != nil {
			log.Printf("[PDF Convert] ledongthuc/pdf: Failed to get text from page %d: %v", pageNum, err)
			continue
		}

		pageText = strings.TrimSpace(pageText)
		if pageText == "" {
			continue
		}

		// 验证提取的文本质量（检查是否有太多乱码）
		if !s.isValidExtractedText(pageText) {
			log.Printf("[PDF Convert] ledongthuc/pdf: Page %d text quality too low, skipping", pageNum)
			continue
		}

		hasContent = true
		result.WriteString(fmt.Sprintf("### 第 %d 页\n\n", pageNum))
		result.WriteString(pageText)
		result.WriteString("\n\n")
	}

	if !hasContent {
		return ""
	}

	return strings.TrimSpace(result.String())
}

// isValidExtractedText 检查提取的文本是否有效（非乱码）
func (s *rawDocumentService) isValidExtractedText(text string) bool {
	if len(text) == 0 {
		return false
	}

	runes := []rune(text)
	if len(runes) == 0 {
		return false
	}

	validCount := 0
	totalCount := len(runes)

	for _, r := range runes {
		// 有效字符：ASCII可打印字符、CJK字符、常见标点
		if (r >= 0x20 && r <= 0x7E) || // ASCII可打印
			(r >= 0x3000 && r <= 0x303F) || // CJK标点
			(r >= 0x3040 && r <= 0x309F) || // 平假名
			(r >= 0x30A0 && r <= 0x30FF) || // 片假名
			(r >= 0x4E00 && r <= 0x9FFF) || // CJK统一汉字
			(r >= 0xFF00 && r <= 0xFFEF) || // 全角字符
			(r >= 0xAC00 && r <= 0xD7AF) || // 韩文
			r == '\n' || r == '\r' || r == '\t' {
			validCount++
		}
	}

	// 如果有效字符占比超过50%，认为文本有效
	ratio := float64(validCount) / float64(totalCount)
	return ratio > 0.5
}

// extractPDFWithPdfcpu 使用pdfcpu纯Go库提取PDF文本
// pdfcpu对CJK字符支持较好，是纯Go实现，无需外部依赖
func (s *rawDocumentService) extractPDFWithPdfcpu(pdfPath string) string {
	// 打开PDF文件
	file, err := os.Open(pdfPath)
	if err != nil {
		log.Printf("[PDF Convert] pdfcpu: Failed to open file: %v", err)
		return ""
	}
	defer file.Close()

	// 使用pdfcpu配置
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed

	// 读取PDF上下文
	ctx, err := api.ReadContext(file, conf)
	if err != nil {
		log.Printf("[PDF Convert] pdfcpu: Failed to read context: %v", err)
		return ""
	}

	// 优化：确保对象已解析
	if err := ctx.EnsurePageCount(); err != nil {
		log.Printf("[PDF Convert] pdfcpu: Failed to ensure page count: %v", err)
		return ""
	}

	var result strings.Builder
	pageCount := ctx.PageCount
	log.Printf("[PDF Convert] pdfcpu: Processing %d pages", pageCount)

	hasContent := false

	for i := 1; i <= pageCount; i++ {
		// 提取页面文本
		pageText, err := s.extractPdfcpuPageText(ctx, i)
		if err != nil {
			log.Printf("[PDF Convert] pdfcpu: Failed to extract page %d: %v", i, err)
			continue
		}

		pageText = strings.TrimSpace(pageText)
		if pageText == "" {
			continue
		}

		hasContent = true
		result.WriteString(fmt.Sprintf("### 第 %d 页\n\n", i))
		result.WriteString(pageText)
		result.WriteString("\n\n")
	}

	if !hasContent {
		return ""
	}

	return strings.TrimSpace(result.String())
}

// extractPdfcpuPageText 从pdfcpu上下文中提取单页文本
func (s *rawDocumentService) extractPdfcpuPageText(ctx *model.Context, pageNum int) (string, error) {
	// 使用pdfcpu的ExtractPageContent提取页面内容流
	// 然后解析内容流中的文本操作符

	// 获取页面对象
	pageDict, _, _, err := ctx.PageDict(pageNum, false)
	if err != nil {
		return "", err
	}

	if pageDict == nil {
		return "", fmt.Errorf("page %d not found", pageNum)
	}

	// 尝试提取页面内容
	var textBuilder strings.Builder

	// 遍历页面资源，提取文本
	// pdfcpu的文本提取需要解析内容流
	contentStream, err := ctx.PageContent(pageDict, pageNum)
	if err != nil {
		return "", err
	}

	// 解析内容流中的文本
	text := s.extractTextFromContentStream(contentStream)
	textBuilder.WriteString(text)

	return textBuilder.String(), nil
}

// extractTextFromContentStream 从PDF内容流中提取文本
func (s *rawDocumentService) extractTextFromContentStream(content []byte) string {
	if len(content) == 0 {
		return ""
	}

	var result strings.Builder
	contentStr := string(content)

	// PDF文本操作符：
	// Tj - 显示字符串
	// TJ - 显示字符串数组
	// ' - 移动到下一行并显示字符串
	// " - 移动到下一行，设置间距并显示字符串

	// 提取括号内的文本 (text) Tj 或 [(text)] TJ
	// 简化的正则匹配
	tjPattern := regexp.MustCompile(`\(([^)]*)\)\s*Tj`)
	matches := tjPattern.FindAllStringSubmatch(contentStr, -1)
	for _, match := range matches {
		if len(match) > 1 {
			text := s.decodePDFString(match[1])
			if text != "" {
				result.WriteString(text)
			}
		}
	}

	// 提取TJ数组中的文本
	tjArrayPattern := regexp.MustCompile(`\[\s*((?:\([^)]*\)|[^]]+)*)\s*\]\s*TJ`)
	arrayMatches := tjArrayPattern.FindAllStringSubmatch(contentStr, -1)
	for _, match := range arrayMatches {
		if len(match) > 1 {
			// 从数组中提取所有字符串
			innerPattern := regexp.MustCompile(`\(([^)]*)\)`)
			innerMatches := innerPattern.FindAllStringSubmatch(match[1], -1)
			for _, inner := range innerMatches {
				if len(inner) > 1 {
					text := s.decodePDFString(inner[1])
					if text != "" {
						result.WriteString(text)
					}
				}
			}
		}
	}

	// 检测换行（BT/ET块之间，或者特定的位置移动）
	// Td, TD, T* 等操作符表示位置移动
	if strings.Contains(contentStr, "Td") || strings.Contains(contentStr, "TD") ||
		strings.Contains(contentStr, "T*") || strings.Contains(contentStr, "Tm") {
		// 在适当位置添加换行
		text := result.String()
		// 简单处理：每隔一定长度添加换行，保持可读性
		text = s.addLineBreaks(text)
		return text
	}

	return result.String()
}

// decodePDFString 解码PDF字符串（处理转义和编码）
func (s *rawDocumentService) decodePDFString(pdfStr string) string {
	// 处理PDF字符串转义
	// \n, \r, \t, \\, \(, \), \ddd (八进制)
	var result strings.Builder
	i := 0
	for i < len(pdfStr) {
		if pdfStr[i] == '\\' && i+1 < len(pdfStr) {
			switch pdfStr[i+1] {
			case 'n':
				result.WriteRune('\n')
				i += 2
			case 'r':
				result.WriteRune('\r')
				i += 2
			case 't':
				result.WriteRune('\t')
				i += 2
			case '\\':
				result.WriteRune('\\')
				i += 2
			case '(':
				result.WriteRune('(')
				i += 2
			case ')':
				result.WriteRune(')')
				i += 2
			default:
				// 检查八进制 \ddd
				if pdfStr[i+1] >= '0' && pdfStr[i+1] <= '7' {
					octal := ""
					j := i + 1
					for j < len(pdfStr) && j < i+4 && pdfStr[j] >= '0' && pdfStr[j] <= '7' {
						octal += string(pdfStr[j])
						j++
					}
					if len(octal) > 0 {
						var val int
						fmt.Sscanf(octal, "%o", &val)
						result.WriteByte(byte(val))
						i = j
					} else {
						result.WriteByte(pdfStr[i])
						i++
					}
				} else {
					result.WriteByte(pdfStr[i])
					i++
				}
			}
		} else {
			result.WriteByte(pdfStr[i])
			i++
		}
	}

	text := result.String()

	// 检查是否是有效的UTF-8
	if !utf8.ValidString(text) {
		// 尝试各种CJK编码
		text = s.convertToUTF8([]byte(text))
	}

	// 清理不可见字符
	return s.cleanPDFText(text)
}

// addLineBreaks 智能添加换行符
func (s *rawDocumentService) addLineBreaks(text string) string {
	if len(text) == 0 {
		return text
	}

	var result strings.Builder
	runes := []rune(text)

	for i, r := range runes {
		result.WriteRune(r)

		// 在句末标点后添加换行
		if r == '。' || r == '！' || r == '？' || r == '.' || r == '!' || r == '?' {
			// 检查下一个字符是否已经是换行
			if i+1 < len(runes) && runes[i+1] != '\n' && runes[i+1] != '\r' {
				result.WriteRune('\n')
			}
		}
	}

	return result.String()
}

// extractPDFWithRscPdf 使用rsc.io/pdf库提取PDF文本（备选方案）
func (s *rawDocumentService) extractPDFWithRscPdf(doc *models.RawDocument, pdfPath string) string {
	var markdownBuilder strings.Builder

	// 使用rsc.io/pdf解析PDF
	pdfReader, err := pdf.Open(pdfPath)
	if err != nil {
		log.Printf("[PDF Convert] Failed to open PDF: %v", err)
		return s.createErrorMarkdown(doc, fmt.Sprintf("PDF转换失败: %v", err))
	}

	numPages := pdfReader.NumPage()
	log.Printf("[PDF Convert] Processing %d pages with rsc.io/pdf", numPages)

	totalTextCount := 0
	skippedPages := 0

	for i := 1; i <= numPages; i++ {
		// 每页单独捕获panic，允许继续处理其他页面
		pageContent := s.extractPDFPageContent(pdfReader, i)
		if pageContent == nil {
			skippedPages++
			continue
		}

		if len(pageContent) == 0 {
			continue
		}

		// 添加页码标题
		markdownBuilder.WriteString(fmt.Sprintf("### 第 %d 页\n\n", i))

		// 智能提取文本：按位置合并相邻文本，检测段落
		var lineTexts []string
		var currentLine strings.Builder
		var lastY float64 = -1
		var lastX float64 = -1

		for _, text := range pageContent {
			// 清理PDF提取的文本，处理CJK字符编码问题
			textStr := s.cleanPDFText(text.S)
			if textStr == "" {
				continue
			}

			// 检测是否为新行（Y坐标变化超过阈值）
			if lastY >= 0 && (lastY-text.Y > 10 || text.Y-lastY > 10) {
				// Y坐标变化较大，可能是新行
				if currentLine.Len() > 0 {
					lineTexts = append(lineTexts, currentLine.String())
					currentLine.Reset()
				}
			} else if lastX >= 0 && text.X-lastX > 20 {
				// X坐标跳跃较大，添加空格分隔
				if currentLine.Len() > 0 {
					currentLine.WriteString(" ")
				}
			}

			currentLine.WriteString(textStr)
			lastY = text.Y
			lastX = text.X + float64(len(textStr)*6) // 估算文本结束位置
			totalTextCount++
		}

		// 添加最后一行
		if currentLine.Len() > 0 {
			lineTexts = append(lineTexts, currentLine.String())
		}

		// 合并连续短行为段落
		var paragraphs []string
		var currentPara strings.Builder

		for _, line := range lineTexts {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 检测是否应该开始新段落
			if currentPara.Len() > 0 {
				paraStr := currentPara.String()
				runes := []rune(paraStr)
				if len(runes) > 0 {
					lastRune := runes[len(runes)-1]
					if lastRune == '。' || lastRune == '！' || lastRune == '？' ||
						lastRune == '.' || lastRune == '!' || lastRune == '?' {
						paragraphs = append(paragraphs, paraStr)
						currentPara.Reset()
					} else {
						currentPara.WriteString(" ")
					}
				}
			}
			currentPara.WriteString(line)
		}

		if currentPara.Len() > 0 {
			paragraphs = append(paragraphs, currentPara.String())
		}

		// 写入段落
		for _, para := range paragraphs {
			para = strings.TrimSpace(para)
			if para != "" {
				markdownBuilder.WriteString(para)
				markdownBuilder.WriteString("\n\n")
			}
		}
	}

	extractedText := strings.TrimSpace(markdownBuilder.String())
	if extractedText == "" || totalTextCount == 0 {
		if skippedPages > 0 {
			extractedText = fmt.Sprintf("(PDF文档部分页面无法解析，跳过了 %d 页。可能为扫描件或图片PDF)", skippedPages)
		} else {
			extractedText = "(PDF文档可能为扫描件或图片PDF，无法提取文本内容)"
		}
	} else if skippedPages > 0 {
		extractedText = fmt.Sprintf("⚠️ 注意：有 %d 页无法解析\n\n%s", skippedPages, extractedText)
	}

	return s.createDocumentMarkdown(doc, extractedText, "PDF文档")
}

// extractPDFPageContent 安全地提取PDF页面内容，捕获可能的panic
func (s *rawDocumentService) extractPDFPageContent(pdfReader *pdf.Reader, pageNum int) (result []pdf.Text) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PDF Convert] Panic on page %d: %v", pageNum, r)
			result = nil
		}
	}()

	page := pdfReader.Page(pageNum)
	if page.V.IsNull() {
		return []pdf.Text{}
	}

	content := page.Content()
	return content.Text
}

// cleanPDFText 清理PDF提取的文本，处理CJK字符编码问题
// PDF内部使用各种字体编码（CID, ToUnicode等），rsc.io/pdf库可能无法正确解析
func (s *rawDocumentService) cleanPDFText(text string) string {
	if text == "" {
		return ""
	}

	// 首先检查是否是有效的UTF-8
	if !utf8.ValidString(text) {
		// 尝试从字节转换
		text = s.convertToUTF8([]byte(text))
	}

	var result strings.Builder
	result.Grow(len(text))

	hasValidChar := false
	consecutiveInvalid := 0

	for _, r := range text {
		// 跳过替换字符（乱码标记）
		if r == '\uFFFD' {
			consecutiveInvalid++
			continue
		}

		// 跳过Private Use Area字符（PDF字体映射失败时常出现）
		if r >= 0xE000 && r <= 0xF8FF {
			consecutiveInvalid++
			continue
		}

		// 跳过无效的控制字符（除了常见的空白字符）
		if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
			continue
		}

		// 检测CID映射失败的字符（通常是高位Unicode区域的无效字符）
		// 这些字符在PDF CJK字体中常见，但实际上是无效的
		if r >= 0x100000 {
			consecutiveInvalid++
			continue
		}

		// 有效字符
		if consecutiveInvalid > 0 && result.Len() > 0 && hasValidChar {
			// 如果之前有跳过的无效字符，可能需要添加空格
			// 但如果连续无效字符太多，可能整段都是乱码
			if consecutiveInvalid < 5 {
				result.WriteRune(' ')
			}
		}
		consecutiveInvalid = 0

		result.WriteRune(r)
		hasValidChar = true
	}

	cleaned := strings.TrimSpace(result.String())

	// 如果清理后内容太少或几乎全是无效字符，返回空
	if len(cleaned) < 2 || float64(len(cleaned))/float64(len(text)) < 0.3 {
		// 可能整个文本都是乱码，尝试其他编码
		return s.tryDecodeAsCJK(text)
	}

	return cleaned
}

// tryDecodeAsCJK 尝试将文本作为CJK编码解码
func (s *rawDocumentService) tryDecodeAsCJK(text string) string {
	// 将字符串转换为字节进行编码检测
	rawBytes := []byte(text)

	// 尝试作为Shift-JIS解码（日文常用）
	result, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), rawBytes)
	if err == nil && utf8.Valid(result) && s.containsValidCJK(string(result)) {
		return strings.TrimSpace(string(result))
	}

	// 尝试作为EUC-JP解码
	result, _, err = transform.Bytes(japanese.EUCJP.NewDecoder(), rawBytes)
	if err == nil && utf8.Valid(result) && s.containsValidCJK(string(result)) {
		return strings.TrimSpace(string(result))
	}

	// 尝试作为GBK解码（简体中文）
	result, _, err = transform.Bytes(simplifiedchinese.GBK.NewDecoder(), rawBytes)
	if err == nil && utf8.Valid(result) && s.containsValidCJK(string(result)) {
		return strings.TrimSpace(string(result))
	}

	return ""
}

// containsValidCJK 检查字符串是否包含有效的CJK字符
func (s *rawDocumentService) containsValidCJK(text string) bool {
	cjkCount := 0
	totalCount := 0

	for _, r := range text {
		totalCount++
		// CJK统一汉字
		if r >= 0x4E00 && r <= 0x9FFF {
			cjkCount++
			continue
		}
		// 日文平假名
		if r >= 0x3040 && r <= 0x309F {
			cjkCount++
			continue
		}
		// 日文片假名
		if r >= 0x30A0 && r <= 0x30FF {
			cjkCount++
			continue
		}
		// CJK扩展A
		if r >= 0x3400 && r <= 0x4DBF {
			cjkCount++
			continue
		}
	}

	// 如果CJK字符占比超过10%，认为是有效的
	return totalCount > 0 && float64(cjkCount)/float64(totalCount) > 0.1
}

// convertWordToMarkdown 将Word文档转换为Markdown
func (s *rawDocumentService) convertWordToMarkdown(doc *models.RawDocument, content []byte) string {
	ext := strings.ToLower(filepath.Ext(doc.OriginalFilename))

	// 只支持DOCX格式的完整解析
	if ext != ".docx" {
		// DOC格式尝试提取文本
		textContent := s.extractTextFromBinary(content)
		if textContent != "" {
			return s.createDocumentMarkdown(doc, textContent, "Word文档(DOC)")
		}
		return s.createErrorMarkdown(doc, "DOC格式暂不支持完整解析，请转换为DOCX格式")
	}

	// 创建临时文件用于DOCX解析
	tmpFile, err := os.CreateTemp("", "docx_*.docx")
	if err != nil {
		log.Printf("[DOCX Convert] Failed to create temp file: %v", err)
		return s.createErrorMarkdown(doc, "DOCX转换失败: 无法创建临时文件")
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(content); err != nil {
		tmpFile.Close()
		log.Printf("[DOCX Convert] Failed to write temp file: %v", err)
		return s.createErrorMarkdown(doc, "DOCX转换失败: 无法写入临时文件")
	}
	tmpFile.Close()

	// 使用docx库解析
	r, err := docx.ReadDocxFile(tmpFile.Name())
	if err != nil {
		log.Printf("[DOCX Convert] Failed to read DOCX: %v", err)
		return s.createErrorMarkdown(doc, fmt.Sprintf("DOCX转换失败: %v", err))
	}
	defer r.Close()

	docxFile := r.Editable()
	textContent := docxFile.GetContent()

	// 清理XML标签，过滤删除线文本
	textContent = s.cleanDocxContentWithStrikethrough(textContent)

	if textContent == "" {
		textContent = "(Word文档内容为空或无法提取)"
	}

	return s.createDocumentMarkdown(doc, textContent, "Word文档(DOCX)")
}

// cleanDocxContentWithStrikethrough 清理DOCX内容中的XML标签，并过滤删除线文本
// Word删除线格式在XML中的表示：
// <w:r><w:rPr><w:strike/></w:rPr><w:t>删除线文本</w:t></w:r>
func (s *rawDocumentService) cleanDocxContentWithStrikethrough(content string) string {
	var result strings.Builder
	decoder := xml.NewDecoder(strings.NewReader(content))

	var inStrikethrough bool
	var inText bool
	var currentRunHasStrike bool

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			// <w:r> 开始一个新的文本运行
			if t.Name.Local == "r" {
				currentRunHasStrike = false
			}
			// <w:strike> 标记删除线
			if t.Name.Local == "strike" {
				currentRunHasStrike = true
				inStrikethrough = true
			}
			// <w:t> 文本节点
			if t.Name.Local == "t" {
				inText = true
			}
		case xml.EndElement:
			// </w:r> 结束文本运行
			if t.Name.Local == "r" {
				inStrikethrough = false
			}
			// </w:t> 结束文本节点
			if t.Name.Local == "t" {
				inText = false
			}
			// </w:p> 段落结束，添加换行
			if t.Name.Local == "p" {
				result.WriteString("\n")
			}
		case xml.CharData:
			// 只保留非删除线的文本
			if inText && !inStrikethrough && !currentRunHasStrike {
				text := strings.TrimSpace(string(t))
				if text != "" {
					result.WriteString(text)
					result.WriteString(" ")
				}
			}
		}
	}

	// 清理多余空白
	text := result.String()
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// 合并多个连续空格
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	// 合并多个连续换行
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}

	return strings.TrimSpace(text)
}

// convertExcelToMarkdown 将Excel文档转换为Markdown
func (s *rawDocumentService) convertExcelToMarkdown(doc *models.RawDocument, content []byte) string {
	ext := strings.ToLower(filepath.Ext(doc.OriginalFilename))

	// CSV文件直接处理
	if ext == ".csv" {
		textContent := s.convertToUTF8(content)
		return s.convertCSVToMarkdownTable(doc, textContent)
	}

	// XLSX文件使用excelize解析
	if ext != ".xlsx" {
		// XLS格式尝试提取文本
		textContent := s.extractTextFromBinary(content)
		if textContent != "" {
			return s.createDocumentMarkdown(doc, textContent, "Excel文档(XLS)")
		}
		return s.createErrorMarkdown(doc, "XLS格式暂不支持完整解析，请转换为XLSX格式")
	}

	// 创建Reader从内存读取
	reader := bytes.NewReader(content)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		log.Printf("[XLSX Convert] Failed to open XLSX: %v", err)
		return s.createErrorMarkdown(doc, fmt.Sprintf("XLSX转换失败: %v", err))
	}
	defer f.Close()

	var markdownBuilder strings.Builder
	markdownBuilder.WriteString(fmt.Sprintf("# %s\n\n", doc.OriginalFilename))
	markdownBuilder.WriteString("## 文档信息\n")
	markdownBuilder.WriteString(fmt.Sprintf("- **文件名**: %s\n", doc.OriginalFilename))
	markdownBuilder.WriteString(fmt.Sprintf("- **文件大小**: %d 字节\n", doc.FileSize))
	markdownBuilder.WriteString(fmt.Sprintf("- **MIME类型**: %s\n", doc.MimeType))
	markdownBuilder.WriteString(fmt.Sprintf("- **转换时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// 遍历所有工作表
	sheetList := f.GetSheetList()
	for _, sheetName := range sheetList {
		markdownBuilder.WriteString(fmt.Sprintf("## 工作表: %s\n\n", sheetName))

		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Printf("[XLSX Convert] Failed to get rows from sheet %s: %v", sheetName, err)
			continue
		}

		if len(rows) == 0 {
			markdownBuilder.WriteString("(空工作表)\n\n")
			continue
		}

		// 转换为Markdown表格，过滤删除线单元格
		markdownBuilder.WriteString(s.rowsToMarkdownTableWithStrikethrough(f, sheetName, rows))
		markdownBuilder.WriteString("\n")
	}

	markdownBuilder.WriteString("---\n*本文档由自动转换工具生成*\n")

	return markdownBuilder.String()
}

// convertCSVToMarkdownTable 将CSV内容转换为Markdown表格
func (s *rawDocumentService) convertCSVToMarkdownTable(doc *models.RawDocument, csvContent string) string {
	lines := strings.Split(csvContent, "\n")
	var rows [][]string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 简单的CSV解析（处理逗号分隔）
		cells := strings.Split(line, ",")
		rows = append(rows, cells)
	}

	var markdownBuilder strings.Builder
	markdownBuilder.WriteString(fmt.Sprintf("# %s\n\n", doc.OriginalFilename))
	markdownBuilder.WriteString("## 文档信息\n")
	markdownBuilder.WriteString(fmt.Sprintf("- **文件名**: %s\n", doc.OriginalFilename))
	markdownBuilder.WriteString(fmt.Sprintf("- **文件大小**: %d 字节\n", doc.FileSize))
	markdownBuilder.WriteString(fmt.Sprintf("- **MIME类型**: %s\n", doc.MimeType))
	markdownBuilder.WriteString(fmt.Sprintf("- **转换时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	markdownBuilder.WriteString("## 表格内容\n\n")

	if len(rows) > 0 {
		markdownBuilder.WriteString(s.rowsToMarkdownTable(rows))
	} else {
		markdownBuilder.WriteString("(空文件)\n")
	}

	markdownBuilder.WriteString("\n---\n*本文档由自动转换工具生成*\n")

	return markdownBuilder.String()
}

// rowsToMarkdownTable 将行数据转换为Markdown表格
func (s *rawDocumentService) rowsToMarkdownTable(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	var builder strings.Builder

	// 找出最大列数
	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	// 限制最大列数，避免表格过宽
	if maxCols > 20 {
		maxCols = 20
	}

	// 写入表头
	builder.WriteString("|")
	for i := 0; i < maxCols; i++ {
		if i < len(rows[0]) {
			cell := strings.TrimSpace(rows[0][i])
			cell = strings.ReplaceAll(cell, "|", "\\|")
			builder.WriteString(fmt.Sprintf(" %s |", cell))
		} else {
			builder.WriteString(" |")
		}
	}
	builder.WriteString("\n")

	// 写入分隔行
	builder.WriteString("|")
	for i := 0; i < maxCols; i++ {
		builder.WriteString(" --- |")
	}
	builder.WriteString("\n")

	// 写入数据行（限制最大行数）
	maxRows := 1000
	if len(rows) > maxRows {
		rows = rows[:maxRows]
	}

	for i := 1; i < len(rows); i++ {
		builder.WriteString("|")
		for j := 0; j < maxCols; j++ {
			if j < len(rows[i]) {
				cell := strings.TrimSpace(rows[i][j])
				cell = strings.ReplaceAll(cell, "|", "\\|")
				cell = strings.ReplaceAll(cell, "\n", " ")
				builder.WriteString(fmt.Sprintf(" %s |", cell))
			} else {
				builder.WriteString(" |")
			}
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// rowsToMarkdownTableWithStrikethrough 将行数据转换为Markdown表格，过滤删除线单元格
func (s *rawDocumentService) rowsToMarkdownTableWithStrikethrough(f *excelize.File, sheetName string, rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	var builder strings.Builder

	// 找出最大列数
	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	// 限制最大列数，避免表格过宽
	if maxCols > 20 {
		maxCols = 20
	}

	// 写入表头
	builder.WriteString("|")
	for i := 0; i < maxCols; i++ {
		if i < len(rows[0]) {
			// 检查单元格是否有删除线
			cellName, _ := excelize.CoordinatesToCellName(i+1, 1)
			if s.isCellStrikethrough(f, sheetName, cellName) {
				builder.WriteString(" |")
				continue
			}
			cell := strings.TrimSpace(rows[0][i])
			cell = strings.ReplaceAll(cell, "|", "\\|")
			builder.WriteString(fmt.Sprintf(" %s |", cell))
		} else {
			builder.WriteString(" |")
		}
	}
	builder.WriteString("\n")

	// 写入分隔行
	builder.WriteString("|")
	for i := 0; i < maxCols; i++ {
		builder.WriteString(" --- |")
	}
	builder.WriteString("\n")

	// 写入数据行（限制最大行数）
	maxRows := 1000
	if len(rows) > maxRows {
		rows = rows[:maxRows]
	}

	for i := 1; i < len(rows); i++ {
		builder.WriteString("|")
		for j := 0; j < maxCols; j++ {
			if j < len(rows[i]) {
				// 检查单元格是否有删除线
				cellName, _ := excelize.CoordinatesToCellName(j+1, i+1)
				if s.isCellStrikethrough(f, sheetName, cellName) {
					builder.WriteString(" |")
					continue
				}
				cell := strings.TrimSpace(rows[i][j])
				cell = strings.ReplaceAll(cell, "|", "\\|")
				cell = strings.ReplaceAll(cell, "\n", " ")
				builder.WriteString(fmt.Sprintf(" %s |", cell))
			} else {
				builder.WriteString(" |")
			}
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// isCellStrikethrough 检查Excel单元格是否有删除线格式
func (s *rawDocumentService) isCellStrikethrough(f *excelize.File, sheetName, cellName string) bool {
	styleID, err := f.GetCellStyle(sheetName, cellName)
	if err != nil {
		return false
	}

	style, err := f.GetStyle(styleID)
	if err != nil {
		return false
	}

	// 检查字体是否有删除线
	if style.Font != nil && style.Font.Strike {
		return true
	}

	return false
}

// convertPowerPointToMarkdown 将PowerPoint文档转换为Markdown
func (s *rawDocumentService) convertPowerPointToMarkdown(doc *models.RawDocument, content []byte) string {
	ext := strings.ToLower(filepath.Ext(doc.OriginalFilename))

	// PPTX是ZIP格式，可以提取XML内容
	if ext == ".pptx" {
		textContent := s.extractTextFromPPTX(content)
		if textContent != "" {
			return s.createDocumentMarkdown(doc, textContent, "PowerPoint文档(PPTX)")
		}
	}

	// PPT格式尝试提取文本
	textContent := s.extractTextFromBinary(content)
	if textContent != "" {
		return s.createDocumentMarkdown(doc, textContent, "PowerPoint文档")
	}

	return s.createErrorMarkdown(doc, "PowerPoint文档内容无法提取，请尝试转换为PPTX格式")
}

// extractTextFromPPTX 从PPTX文件中提取文本
// PPTX是Office Open XML格式，实际上是一个ZIP压缩包
// 文本内容在 ppt/slides/slide*.xml 文件中
func (s *rawDocumentService) extractTextFromPPTX(content []byte) string {
	reader := bytes.NewReader(content)

	// 使用archive/zip正确解析PPTX
	zipReader, err := zip.NewReader(reader, int64(len(content)))
	if err != nil {
		log.Printf("[PPTX Extract] Failed to open as ZIP: %v", err)
		return ""
	}

	// 收集所有slide文件并排序
	type slideInfo struct {
		name  string
		index int
		file  *zip.File
	}
	var slides []slideInfo
	slidePattern := regexp.MustCompile(`ppt/slides/slide(\d+)\.xml`)

	for _, file := range zipReader.File {
		if matches := slidePattern.FindStringSubmatch(file.Name); matches != nil {
			var idx int
			fmt.Sscanf(matches[1], "%d", &idx)
			slides = append(slides, slideInfo{name: file.Name, index: idx, file: file})
		}
	}

	// 按slide编号排序
	sort.Slice(slides, func(i, j int) bool {
		return slides[i].index < slides[j].index
	})

	var markdownBuilder strings.Builder

	for _, slide := range slides {
		rc, err := slide.file.Open()
		if err != nil {
			log.Printf("[PPTX Extract] Failed to open slide %s: %v", slide.name, err)
			continue
		}

		slideContent, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			log.Printf("[PPTX Extract] Failed to read slide %s: %v", slide.name, err)
			continue
		}

		// 提取slide中的文本
		slideText := s.extractTextFromSlideXML(slideContent)
		if slideText != "" {
			markdownBuilder.WriteString(fmt.Sprintf("### 第 %d 页\n\n", slide.index))
			markdownBuilder.WriteString(slideText)
			markdownBuilder.WriteString("\n\n")
		}
	}

	return strings.TrimSpace(markdownBuilder.String())
}

// extractTextFromSlideXML 从slide的XML内容中提取文本，过滤删除线文本
// PowerPoint删除线格式：<a:r><a:rPr strike="sngStrike"/><a:t>文本</a:t></a:r>
func (s *rawDocumentService) extractTextFromSlideXML(xmlContent []byte) string {
	var textParts []string

	// 使用XML解码器提取所有<a:t>标签中的文本
	decoder := xml.NewDecoder(bytes.NewReader(xmlContent))
	var currentParagraph []string
	var currentRunHasStrike bool
	var inTextRun bool

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			// <a:p> 是段落开始
			if t.Name.Local == "p" && t.Name.Space == "http://schemas.openxmlformats.org/drawingml/2006/main" {
				currentParagraph = []string{}
			}
			// <a:r> 是文本运行开始
			if t.Name.Local == "r" {
				inTextRun = true
				currentRunHasStrike = false
			}
			// <a:rPr> 是文本属性，检查删除线
			if t.Name.Local == "rPr" && inTextRun {
				for _, attr := range t.Attr {
					if attr.Name.Local == "strike" {
						// strike值可以是 "sngStrike"（单删除线）或 "dblStrike"（双删除线）
						if attr.Value == "sngStrike" || attr.Value == "dblStrike" {
							currentRunHasStrike = true
						}
					}
				}
			}
			// <a:t> 是文本元素
			if t.Name.Local == "t" && !currentRunHasStrike {
				var text string
				if err := decoder.DecodeElement(&text, &t); err == nil && strings.TrimSpace(text) != "" {
					currentParagraph = append(currentParagraph, text)
				}
			} else if t.Name.Local == "t" && currentRunHasStrike {
				// 跳过删除线文本，但仍需消费这个元素
				var text string
				decoder.DecodeElement(&text, &t)
			}
		case xml.EndElement:
			// <a:r> 文本运行结束
			if t.Name.Local == "r" {
				inTextRun = false
				currentRunHasStrike = false
			}
			// <a:p> 段落结束，合并该段落的文本
			if t.Name.Local == "p" && t.Name.Space == "http://schemas.openxmlformats.org/drawingml/2006/main" {
				if len(currentParagraph) > 0 {
					paragraphText := strings.Join(currentParagraph, "")
					if strings.TrimSpace(paragraphText) != "" {
						textParts = append(textParts, paragraphText)
					}
				}
			}
		}
	}

	return strings.Join(textParts, "\n\n")
}

// extractTextFromBinary 从二进制文件中尝试提取可读文本
func (s *rawDocumentService) extractTextFromBinary(content []byte) string {
	var textBuilder strings.Builder
	var currentWord strings.Builder

	for _, b := range content {
		// 检查是否为可打印字符
		if (b >= 0x20 && b <= 0x7E) || b == '\n' || b == '\r' || b == '\t' {
			currentWord.WriteByte(b)
		} else if b >= 0x80 {
			// 可能是多字节UTF-8字符，保留
			currentWord.WriteByte(b)
		} else {
			// 非文本字符，检查当前词是否有意义
			word := currentWord.String()
			if len(word) >= 2 { // 至少2个字符才保留
				textBuilder.WriteString(word)
				textBuilder.WriteString(" ")
			}
			currentWord.Reset()
		}
	}

	// 处理最后一个词
	if currentWord.Len() >= 2 {
		textBuilder.WriteString(currentWord.String())
	}

	result := textBuilder.String()

	// 尝试UTF-8解码
	if !utf8.ValidString(result) {
		result = s.convertToUTF8([]byte(result))
	}

	return strings.TrimSpace(result)
}

// createDocumentMarkdown 创建文档Markdown内容
func (s *rawDocumentService) createDocumentMarkdown(doc *models.RawDocument, textContent, docType string) string {
	return fmt.Sprintf(`# %s

## 文档信息
- **文件名**: %s
- **文件大小**: %d 字节
- **MIME类型**: %s
- **文档类型**: %s
- **转换时间**: %s

## 文档内容

%s

---
*本文档由自动转换工具生成*
`,
		doc.OriginalFilename,
		doc.OriginalFilename,
		doc.FileSize,
		doc.MimeType,
		docType,
		time.Now().Format("2006-01-02 15:04:05"),
		textContent,
	)
}

// createErrorMarkdown 创建错误提示Markdown
func (s *rawDocumentService) createErrorMarkdown(doc *models.RawDocument, errorMsg string) string {
	return fmt.Sprintf(`# %s

## 文档信息
- **文件名**: %s
- **文件大小**: %d 字节
- **MIME类型**: %s
- **转换时间**: %s

## 转换提示

⚠️ %s

---
*本文档由自动转换工具生成*
`,
		doc.OriginalFilename,
		doc.OriginalFilename,
		doc.FileSize,
		doc.MimeType,
		time.Now().Format("2006-01-02 15:04:05"),
		errorMsg,
	)
}

// isImageMimeType 检查MIME类型是否为图片
func (s *rawDocumentService) isImageMimeType(mimeType string) bool {
	imageMimeTypes := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/jpg":  true,
		"image/bmp":  true,
		"image/tiff": true,
		"image/gif":  true,
		"image/webp": true,
	}
	return imageMimeTypes[mimeType]
}

// convertImageToMarkdown 将图片转换为Markdown格式（base64嵌入）
func (s *rawDocumentService) convertImageToMarkdown(doc *models.RawDocument, content []byte) string {
	// 将图片转换为base64
	base64Content := base64.StdEncoding.EncodeToString(content)
	mimeType := doc.MimeType

	// 构建Markdown内容
	markdownContent := fmt.Sprintf(`# %s

## 文档信息
- **文件名**: %s
- **文件大小**: %d 字节
- **MIME类型**: %s
- **转换时间**: %s

## 图片预览

![%s](data:%s;base64,%s)

---
*本文档由自动转换工具生成*
`,
		doc.OriginalFilename,
		doc.OriginalFilename,
		doc.FileSize,
		doc.MimeType,
		time.Now().Format("2006-01-02 15:04:05"),
		doc.OriginalFilename,
		mimeType,
		base64Content,
	)

	return markdownContent
}

// convertToUTF8 将字节内容转换为UTF-8编码的字符串
// 支持自动检测和转换以下编码：UTF-8, UTF-16, GBK, GB2312, Big5, Shift-JIS, EUC-JP
func (s *rawDocumentService) convertToUTF8(content []byte) string {
	// 如果已经是有效的UTF-8，直接返回
	if utf8.Valid(content) {
		// 移除BOM标记（如果存在）
		content = bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF})
		return string(content)
	}

	// 检测UTF-16 BOM
	if len(content) >= 2 {
		// UTF-16 LE BOM
		if content[0] == 0xFF && content[1] == 0xFE {
			decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
			result, _, err := transform.Bytes(decoder, content)
			if err == nil {
				return string(result)
			}
		}
		// UTF-16 BE BOM
		if content[0] == 0xFE && content[1] == 0xFF {
			decoder := unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder()
			result, _, err := transform.Bytes(decoder, content)
			if err == nil {
				return string(result)
			}
		}
	}

	// 尝试各种编码解码
	encodings := []struct {
		name    string
		decoder transform.Transformer
	}{
		{"GBK", simplifiedchinese.GBK.NewDecoder()},
		{"GB18030", simplifiedchinese.GB18030.NewDecoder()},
		{"Big5", traditionalchinese.Big5.NewDecoder()},
		{"Shift-JIS", japanese.ShiftJIS.NewDecoder()},
		{"EUC-JP", japanese.EUCJP.NewDecoder()},
		{"ISO-2022-JP", japanese.ISO2022JP.NewDecoder()},
	}

	for _, enc := range encodings {
		result, _, err := transform.Bytes(enc.decoder, content)
		if err == nil && utf8.Valid(result) {
			log.Printf("[Encoding] Successfully decoded content using %s", enc.name)
			return string(result)
		}
	}

	// 如果所有编码都失败，使用lossy转换（替换无效字符）
	log.Printf("[Encoding] Warning: Could not detect encoding, using lossy UTF-8 conversion")
	return s.toValidUTF8(content)
}

// toValidUTF8 将字节转换为有效的UTF-8，替换无效字符
func (s *rawDocumentService) toValidUTF8(content []byte) string {
	var result strings.Builder
	result.Grow(len(content))

	for len(content) > 0 {
		r, size := utf8.DecodeRune(content)
		if r == utf8.RuneError && size == 1 {
			// 无效字节，用替换字符
			result.WriteRune('�')
			content = content[1:]
		} else {
			result.WriteRune(r)
			content = content[size:]
		}
	}

	return result.String()
}

// sanitizeFilename 清理文件名，移除危险字符
func (s *rawDocumentService) sanitizeFilename(filename string) string {
	// 移除文件扩展名
	name := filename
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		name = filename[:idx]
	}

	// 移除或转义特殊字符
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)

	return replacer.Replace(name)
}

// updateConvertStatus 更新转换状态（辅助函数）
func (s *rawDocumentService) updateConvertStatus(documentID uint, status string, progress int, filename string, filepath string, filesize int64, convertError string) {
	if err := s.repo.UpdateStatus(documentID, status, progress, filename, filepath, filesize, convertError); err != nil {
		log.Printf("[Convert Update Failed] documentId=%d, status=%s, error: %v", documentID, status, err)
	}
}

// GetConvertStatus 查询转换状态
func (s *rawDocumentService) GetConvertStatus(id uint) (*models.ConvertStatusResponse, error) {
	// 使用轻量级查询直接获取转换状态
	return s.repo.GetConvertStatus(id)
}

// DownloadOriginal 下载原始文档
func (s *rawDocumentService) DownloadOriginal(id uint) (*models.RawDocument, io.ReadCloser, error) {
	doc, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("document not found")
		}
		return nil, nil, fmt.Errorf("get document: %w", err)
	}

	fullPath := filepath.Join(s.storageBasePath, doc.OriginalFilepath)
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, errors.New("document file not found")
		}
		return nil, nil, fmt.Errorf("open document file: %w", err)
	}

	return doc, file, nil
}

// DownloadConverted 下载转换后的文档
func (s *rawDocumentService) DownloadConverted(id uint) (*models.RawDocument, io.ReadCloser, error) {
	doc, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("document not found")
		}
		return nil, nil, fmt.Errorf("get document: %w", err)
	}

	// 校验转换状态
	if doc.ConvertStatus != "completed" {
		return nil, nil, errors.New("document not converted yet")
	}

	if doc.ConvertedFilepath == "" {
		return nil, nil, errors.New("converted file path not found")
	}

	fullPath := filepath.Join(s.storageBasePath, doc.ConvertedFilepath)
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, errors.New("converted file not found")
		}
		return nil, nil, fmt.Errorf("open converted file: %w", err)
	}

	return doc, file, nil
}

// PreviewConverted 预览转换后的Markdown文档内容
func (s *rawDocumentService) PreviewConverted(id uint) (*models.RawDocument, string, error) {
	doc, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("document not found")
		}
		return nil, "", fmt.Errorf("get document: %w", err)
	}

	// 校验转换状态
	if doc.ConvertStatus != "completed" {
		return nil, "", errors.New("document not converted yet")
	}

	if doc.ConvertedFilepath == "" {
		return nil, "", errors.New("converted file path not found")
	}

	fullPath := filepath.Join(s.storageBasePath, doc.ConvertedFilepath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", errors.New("converted file not found")
		}
		return nil, "", fmt.Errorf("read converted file: %w", err)
	}

	return doc, string(content), nil
}

// DeleteOriginal 删除原始文档（包括转换文件和数据库记录）
func (s *rawDocumentService) DeleteOriginal(id uint) error {
	// 获取文档信息
	doc, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("document not found")
		}
		return fmt.Errorf("get document: %w", err)
	}

	// 删除原始文件
	if doc.OriginalFilepath != "" {
		fullPath := filepath.Join(s.storageBasePath, doc.OriginalFilepath)
		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			log.Printf("[WARN] failed to delete original file %s: %v", fullPath, err)
		}
	}

	// 删除转换文件（如果存在）
	if doc.ConvertedFilepath != "" {
		fullPath := filepath.Join(s.storageBasePath, doc.ConvertedFilepath)
		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			log.Printf("[WARN] failed to delete converted file %s: %v", fullPath, err)
		}
	}

	// 软删除数据库记录
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete document record: %w", err)
	}

	log.Printf("[RawDocument Delete] id=%d", id)
	return nil
}

// DeleteConverted 仅删除转换后的文档（重置转换状态）
func (s *rawDocumentService) DeleteConverted(id uint) error {
	// 获取文档信息
	doc, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("document not found")
		}
		return fmt.Errorf("get document: %w", err)
	}

	// 删除转换文件
	if doc.ConvertedFilepath != "" {
		fullPath := filepath.Join(s.storageBasePath, doc.ConvertedFilepath)
		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			log.Printf("[WARN] failed to delete converted file %s: %v", fullPath, err)
		}
	}

	// 重置转换状态字段
	doc.ConvertStatus = "none"
	doc.ConvertTaskID = ""
	doc.ConvertProgress = 0
	doc.ConvertedFilename = ""
	doc.ConvertedFilepath = ""
	doc.ConvertedTime = nil
	doc.ConvertError = ""

	if err := s.repo.Update(doc); err != nil {
		return fmt.Errorf("reset convert status: %w", err)
	}

	log.Printf("[RawDocument DeleteConverted] id=%d", id)
	return nil
}
