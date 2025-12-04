package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

// SanitizeFilename 过滤文件名中的危险字符,防止路径遍历攻击
func SanitizeFilename(filename string) string {
	// 移除路径遍历字符
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	// 只保留文件名部分
	return filepath.Base(filename)
}

// ValidateFileSize 检查文件大小是否超限
func ValidateFileSize(size int64) error {
	if size > MaxFileSize {
		return fmt.Errorf("file size %d exceeds limit %d", size, MaxFileSize)
	}
	return nil
}

// ValidateFileExtension 验证文件扩展名
func ValidateFileExtension(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".xlsx" && ext != ".xls" {
		return fmt.Errorf("invalid file extension %s, only .xlsx and .xls are supported", ext)
	}
	return nil
}
