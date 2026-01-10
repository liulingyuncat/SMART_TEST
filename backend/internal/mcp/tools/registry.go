package tools

import (
	"sort"
	"sync"
)

// ToolAnnotations represents the annotations for a tool.
type ToolAnnotations struct {
	// Title is the human-readable title of the tool (optional).
	Title string `json:"title,omitempty"`
	// ReadOnlyHint indicates if the tool only reads data without side effects.
	ReadOnlyHint bool `json:"readOnlyHint,omitempty"`
	// DestructiveHint indicates if the tool may perform destructive updates.
	DestructiveHint bool `json:"destructiveHint,omitempty"`
	// IdempotentHint indicates if calling the tool multiple times has the same effect.
	IdempotentHint bool `json:"idempotentHint,omitempty"`
	// OpenWorldHint indicates if the tool may interact with external entities.
	OpenWorldHint bool `json:"openWorldHint,omitempty"`
}

// ToolDefinition represents the definition of a tool for the tools/list response.
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
	Annotations *ToolAnnotations       `json:"annotations,omitempty"`
}

// ToolRegistry manages the registration and lookup of tool handlers.
type ToolRegistry struct {
	handlers map[string]ToolHandler
	mu       sync.RWMutex
}

// NewToolRegistry creates a new ToolRegistry.
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		handlers: make(map[string]ToolHandler),
	}
}

// Register registers a tool handler.
// If a handler with the same name already exists, it will be replaced.
func (r *ToolRegistry) Register(handler ToolHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[handler.Name()] = handler
}

// Get returns the handler for the specified tool name, or nil if not found.
func (r *ToolRegistry) Get(name string) ToolHandler {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.handlers[name]
}

// List returns all registered tool definitions sorted by name.
func (r *ToolRegistry) List() []ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	definitions := make([]ToolDefinition, 0, len(r.handlers))
	for _, handler := range r.handlers {
		definitions = append(definitions, ToolDefinition{
			Name:        handler.Name(),
			Description: handler.Description(),
			InputSchema: handler.InputSchema(),
			Annotations: handler.Annotations(),
		})
	}

	// Sort by name for consistent ordering
	sort.Slice(definitions, func(i, j int) bool {
		return definitions[i].Name < definitions[j].Name
	})

	return definitions
}

// Count returns the number of registered tools.
func (r *ToolRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.handlers)
}

// Has returns true if a tool with the specified name is registered.
func (r *ToolRegistry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.handlers[name]
	return exists
}

// Names returns the names of all registered tools.
func (r *ToolRegistry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
