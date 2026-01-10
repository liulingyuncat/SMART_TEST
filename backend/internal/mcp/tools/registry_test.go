package tools

import (
	"context"
	"testing"
)

// mockHandler is a mock implementation of ToolHandler for testing.
type mockHandler struct {
	name        string
	description string
	schema      map[string]interface{}
}

func (m *mockHandler) Name() string {
	return m.name
}

func (m *mockHandler) Description() string {
	return m.description
}

func (m *mockHandler) InputSchema() map[string]interface{} {
	return m.schema
}

func (m *mockHandler) Annotations() map[string]interface{} {
	return nil
}

func (m *mockHandler) Execute(ctx context.Context, args map[string]interface{}) (ToolResult, error) {
	return NewTextResult("mock result"), nil
}

func TestToolRegistry_Register(t *testing.T) {
	registry := NewToolRegistry()

	handler := &mockHandler{
		name:        "test_tool",
		description: "A test tool",
		schema:      map[string]interface{}{"type": "object"},
	}

	registry.Register(handler)

	if registry.Count() != 1 {
		t.Errorf("Count() = %d, want 1", registry.Count())
	}
}

func TestToolRegistry_Get(t *testing.T) {
	registry := NewToolRegistry()

	handler := &mockHandler{
		name:        "test_tool",
		description: "A test tool",
		schema:      map[string]interface{}{"type": "object"},
	}

	registry.Register(handler)

	// Get existing tool
	got := registry.Get("test_tool")
	if got == nil {
		t.Fatal("Get() returned nil for existing tool")
	}
	if got.Name() != "test_tool" {
		t.Errorf("Get().Name() = %v, want test_tool", got.Name())
	}

	// Get non-existing tool
	got = registry.Get("non_existent")
	if got != nil {
		t.Error("Get() returned non-nil for non-existent tool")
	}
}

func TestToolRegistry_Has(t *testing.T) {
	registry := NewToolRegistry()

	handler := &mockHandler{name: "test_tool"}
	registry.Register(handler)

	if !registry.Has("test_tool") {
		t.Error("Has() = false for existing tool")
	}
	if registry.Has("non_existent") {
		t.Error("Has() = true for non-existent tool")
	}
}

func TestToolRegistry_Names(t *testing.T) {
	registry := NewToolRegistry()

	registry.Register(&mockHandler{name: "tool_a"})
	registry.Register(&mockHandler{name: "tool_b"})
	registry.Register(&mockHandler{name: "tool_c"})

	names := registry.Names()
	if len(names) != 3 {
		t.Errorf("Names() length = %d, want 3", len(names))
	}

	// Check that all names are present (order not guaranteed)
	nameMap := make(map[string]bool)
	for _, n := range names {
		nameMap[n] = true
	}
	for _, expected := range []string{"tool_a", "tool_b", "tool_c"} {
		if !nameMap[expected] {
			t.Errorf("Names() missing %s", expected)
		}
	}
}

func TestToolRegistry_List(t *testing.T) {
	registry := NewToolRegistry()

	registry.Register(&mockHandler{
		name:        "test_tool",
		description: "Test description",
		schema:      map[string]interface{}{"type": "object"},
	})

	list := registry.List()
	if len(list) != 1 {
		t.Fatalf("List() length = %d, want 1", len(list))
	}

	tool := list[0]
	if tool.Name != "test_tool" {
		t.Errorf("tool.Name = %v, want test_tool", tool.Name)
	}
	if tool.Description != "Test description" {
		t.Errorf("tool.Description = %v, want 'Test description'", tool.Description)
	}
}

func TestToolResult_NewTextResult(t *testing.T) {
	result := NewTextResult("hello world")

	if result.IsError {
		t.Error("IsError = true, want false")
	}
	if len(result.Content) != 1 {
		t.Fatalf("Content length = %d, want 1", len(result.Content))
	}
	if result.Content[0].Type != "text" {
		t.Errorf("Content[0].Type = %v, want text", result.Content[0].Type)
	}
	if result.Content[0].Text != "hello world" {
		t.Errorf("Content[0].Text = %v, want 'hello world'", result.Content[0].Text)
	}
}

func TestToolResult_NewErrorResult(t *testing.T) {
	result := NewErrorResult("something went wrong")

	if !result.IsError {
		t.Error("IsError = false, want true")
	}
	if len(result.Content) != 1 {
		t.Fatalf("Content length = %d, want 1", len(result.Content))
	}
	if result.Content[0].Text != "something went wrong" {
		t.Errorf("Content[0].Text = %v, want 'something went wrong'", result.Content[0].Text)
	}
}

func TestToolResult_NewJSONResult(t *testing.T) {
	result := NewJSONResult(`{"key":"value"}`)

	if result.IsError {
		t.Error("IsError = true, want false")
	}
	if len(result.Content) != 1 {
		t.Fatalf("Content length = %d, want 1", len(result.Content))
	}
	if result.Content[0].MimeType != "application/json" {
		t.Errorf("Content[0].MimeType = %v, want application/json", result.Content[0].MimeType)
	}
}
