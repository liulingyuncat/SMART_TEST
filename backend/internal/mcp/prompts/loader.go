package prompts

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// PromptLoader 负责扫描和解析.prompt.md文件
type PromptLoader struct{}

// FrontMatter YAML front-matter结构
type FrontMatter struct {
	Name        string                   `yaml:"name"`
	Description string                   `yaml:"description"`
	Version     string                   `yaml:"version"`
	Arguments   []map[string]interface{} `yaml:"arguments"`
}

// LoadAll 扫描目录并加载所有Prompt
func (l *PromptLoader) LoadAll(dir string, registry *PromptsRegistry) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read directory %s: %v\n", dir, err)
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	fmt.Printf("[DEBUG] Directory %s exists, found %d entries\n", dir, len(entries))

	loadedCount := 0
	failedCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !strings.HasSuffix(fileName, ".prompt.md") {
			continue
		}

		filePath := filepath.Join(dir, fileName)
		prompt, err := l.ParsePromptFile(filePath)
		if err != nil {
			// 解析失败记录警告，不阻塞启动
			fmt.Printf("[WARN] Failed to parse prompt file %s: %v\n", fileName, err)
			failedCount++
			continue
		}

		registry.Register(prompt)
		fmt.Printf("[DEBUG] Registered prompt: %s\n", prompt.Name)
		loadedCount++
	}

	fmt.Printf("[INFO] LoadAll completed: loaded=%d, failed=%d, total=%d\n", loadedCount, failedCount, loadedCount+failedCount)

	return nil
}

// ParsePromptFile 解析单个.prompt.md文件
func (l *PromptLoader) ParsePromptFile(path string) (*SystemPrompt, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	frontMatter, _, err := l.extractFrontMatter(content)
	if err != nil {
		return nil, fmt.Errorf("extract front matter: %w", err)
	}

	// 解析YAML
	var fm FrontMatter
	if err := yaml.Unmarshal([]byte(frontMatter), &fm); err != nil {
		return nil, fmt.Errorf("parse YAML: %w", err)
	}

	// 验证必填字段
	if fm.Name == "" {
		return nil, fmt.Errorf("missing required field: name")
	}

	// 解析arguments
	var args []PromptArgument
	for _, argMap := range fm.Arguments {
		arg := PromptArgument{
			Name:        getString(argMap, "name"),
			Description: getString(argMap, "description"),
			Required:    getBool(argMap, "required"),
		}
		if arg.Name != "" {
			args = append(args, arg)
		}
	}

	prompt := &SystemPrompt{
		Name:        fm.Name,
		Description: fm.Description,
		Version:     fm.Version,
		Arguments:   args,
		FilePath:    path,
		// content字段不加载（延迟加载）
	}

	return prompt, nil
}

// ReadContent 读取Prompt文件的Markdown正文内容
func (l *PromptLoader) ReadContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	_, markdownContent, err := l.extractFrontMatter(content)
	if err != nil {
		return "", fmt.Errorf("extract content: %w", err)
	}

	return markdownContent, nil
}

// extractFrontMatter 分离YAML front-matter和Markdown正文
func (l *PromptLoader) extractFrontMatter(content []byte) (string, string, error) {
	// YAML front-matter格式：以---开头和结尾
	if !bytes.HasPrefix(content, []byte("---\n")) && !bytes.HasPrefix(content, []byte("---\r\n")) {
		return "", "", fmt.Errorf("invalid format: missing YAML front-matter")
	}

	// 去除第一个---
	content = bytes.TrimPrefix(content, []byte("---\n"))
	content = bytes.TrimPrefix(content, []byte("---\r\n"))

	// 查找第二个---
	delimiterIndex := bytes.Index(content, []byte("\n---\n"))
	if delimiterIndex == -1 {
		delimiterIndex = bytes.Index(content, []byte("\r\n---\r\n"))
		if delimiterIndex == -1 {
			return "", "", fmt.Errorf("invalid format: missing closing ---")
		}
		// CRLF情况
		frontMatter := string(content[:delimiterIndex])
		markdown := strings.TrimSpace(string(content[delimiterIndex+6:]))
		return frontMatter, markdown, nil
	}

	// LF情况
	frontMatter := string(content[:delimiterIndex])
	markdown := strings.TrimSpace(string(content[delimiterIndex+5:]))

	return frontMatter, markdown, nil
}

// getString 从map中安全获取string值
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getBool 从map中安全获取bool值
func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}
