package prompts

import (
	"sort"
	"sync"
	"time"
)

// SystemPrompt 系统提示词定义
type SystemPrompt struct {
	Name        string
	Description string
	Version     string
	Arguments   []PromptArgument
	FilePath    string    // 原始文件路径
	content     string    // 内容（延迟加载，私有字段）
	UpdatedAt   time.Time // 上次更新时间
}

// PromptArgument 提示词参数定义
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// GetContent 获取Prompt内容（支持延迟加载）
func (p *SystemPrompt) GetContent() (string, error) {
	if p.content == "" {
		// 内容未加载，从文件读取
		loader := &PromptLoader{}
		content, err := loader.ReadContent(p.FilePath)
		if err != nil {
			return "", err
		}
		p.content = content
	}
	return p.content, nil
}

// InvalidateCache 使缓存失效，强制重新加载
func (p *SystemPrompt) InvalidateCache() {
	p.content = ""
	p.UpdatedAt = time.Now()
}

// PromptMetadata 提示词元数据（不含content）
type PromptMetadata struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Version     string           `json:"version,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
	UpdatedAt   int64            `json:"updated_at,omitempty"` // Unix时间戳
}

// PromptsRegistry 系统提示词注册表
type PromptsRegistry struct {
	prompts      map[string]*SystemPrompt
	mu           sync.RWMutex
	lastReloadAt time.Time // 上次重新加载时间
}

// NewPromptsRegistry 创建新的PromptsRegistry
func NewPromptsRegistry() *PromptsRegistry {
	return &PromptsRegistry{
		prompts:      make(map[string]*SystemPrompt),
		lastReloadAt: time.Now(),
	}
}

// Register 注册一个系统Prompt
func (r *PromptsRegistry) Register(prompt *SystemPrompt) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if prompt.UpdatedAt.IsZero() {
		prompt.UpdatedAt = time.Now()
	}
	r.prompts[prompt.Name] = prompt
}

// Update 更新一个系统Prompt（用于热更新）
func (r *PromptsRegistry) Update(name string, description, version string, arguments []PromptArgument) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if prompt, exists := r.prompts[name]; exists {
		if description != "" {
			prompt.Description = description
		}
		if version != "" {
			prompt.Version = version
		}
		if arguments != nil {
			prompt.Arguments = arguments
		}
		prompt.InvalidateCache()
		return true
	}
	return false
}

// Unregister 注销一个系统Prompt
func (r *PromptsRegistry) Unregister(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prompts[name]; exists {
		delete(r.prompts, name)
		return true
	}
	return false
}

// Get 根据名称查询系统Prompt
func (r *PromptsRegistry) Get(name string) (*SystemPrompt, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	prompt, ok := r.prompts[name]
	return prompt, ok
}

// List 返回所有系统Prompt的元数据（不含content）
func (r *PromptsRegistry) List() []PromptMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metadata := make([]PromptMetadata, 0, len(r.prompts))
	for _, prompt := range r.prompts {
		metadata = append(metadata, PromptMetadata{
			Name:        prompt.Name,
			Description: prompt.Description,
			Version:     prompt.Version,
			Arguments:   prompt.Arguments,
			UpdatedAt:   prompt.UpdatedAt.Unix(),
		})
	}

	// 按名称排序，确保顺序一致
	sort.Slice(metadata, func(i, j int) bool {
		return metadata[i].Name < metadata[j].Name
	})

	return metadata
}

// GetLastReloadTime 返回上次重新加载的时间
func (r *PromptsRegistry) GetLastReloadTime() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastReloadAt
}

// Count 返回已注册的系统Prompt数量
func (r *PromptsRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.prompts)
}
