package tools

import (
	"testing"
)

func TestSchemaValidator_Required(t *testing.T) {
	v := NewSchemaValidator()

	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{"type": "string"},
			"age":  map[string]interface{}{"type": "integer"},
		},
		"required": []interface{}{"name"},
	}

	// Valid: has required field
	err := v.Validate(schema, map[string]interface{}{"name": "test"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Invalid: missing required field
	err = v.Validate(schema, map[string]interface{}{"age": 10})
	if err == nil {
		t.Error("expected error for missing required field")
	}
}

func TestSchemaValidator_TypeCheck(t *testing.T) {
	v := NewSchemaValidator()

	tests := []struct {
		name      string
		schema    map[string]interface{}
		args      map[string]interface{}
		wantError bool
	}{
		{
			name: "valid string",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{"type": "string"},
				},
			},
			args:      map[string]interface{}{"name": "test"},
			wantError: false,
		},
		{
			name: "invalid string type",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{"type": "string"},
				},
			},
			args:      map[string]interface{}{"name": []string{"not", "a", "string"}},
			wantError: true,
		},
		{
			name: "valid integer",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"count": map[string]interface{}{"type": "integer"},
				},
			},
			args:      map[string]interface{}{"count": float64(10)},
			wantError: false,
		},
		{
			name: "valid boolean",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"active": map[string]interface{}{"type": "boolean"},
				},
			},
			args:      map[string]interface{}{"active": true},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.schema, tt.args)
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestSchemaValidator_Enum(t *testing.T) {
	v := NewSchemaValidator()

	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"status": map[string]interface{}{
				"type": "string",
				"enum": []interface{}{"pending", "active", "done"},
			},
		},
	}

	// Valid enum value
	err := v.Validate(schema, map[string]interface{}{"status": "active"})
	if err != nil {
		t.Errorf("expected no error for valid enum, got %v", err)
	}

	// Invalid enum value
	err = v.Validate(schema, map[string]interface{}{"status": "invalid"})
	if err == nil {
		t.Error("expected error for invalid enum value")
	}
}

func TestSchemaValidator_EmptyArgs(t *testing.T) {
	v := NewSchemaValidator()

	schema := map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}

	err := v.Validate(schema, map[string]interface{}{})
	if err != nil {
		t.Errorf("expected no error for empty args with no required fields, got %v", err)
	}
}

func TestSchemaValidator_NilSchema(t *testing.T) {
	v := NewSchemaValidator()

	err := v.Validate(nil, map[string]interface{}{"key": "value"})
	if err != nil {
		t.Errorf("expected no error for nil schema, got %v", err)
	}
}
