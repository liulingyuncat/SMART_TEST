package tools

import (
	"fmt"
	"strings"
)

// ValidationError represents a schema validation error.
type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Reason)
}

// SchemaValidator validates tool input against JSON Schema.
type SchemaValidator struct{}

// NewSchemaValidator creates a new SchemaValidator.
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{}
}

// Validate validates args against the provided JSON Schema.
func (v *SchemaValidator) Validate(schema map[string]interface{}, args map[string]interface{}) error {
	// Check required fields
	if required, ok := schema["required"].([]interface{}); ok {
		for _, r := range required {
			fieldName, ok := r.(string)
			if !ok {
				continue
			}
			if _, exists := args[fieldName]; !exists {
				return &ValidationError{
					Field:  fieldName,
					Reason: "required field is missing",
				}
			}
		}
	}

	// Validate properties
	if properties, ok := schema["properties"].(map[string]interface{}); ok {
		for fieldName, propSchema := range properties {
			if value, exists := args[fieldName]; exists {
				propSchemaMap, ok := propSchema.(map[string]interface{})
				if !ok {
					continue
				}
				if err := v.validateProperty(fieldName, value, propSchemaMap); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// validateProperty validates a single property value against its schema.
func (v *SchemaValidator) validateProperty(fieldName string, value interface{}, schema map[string]interface{}) error {
	// Check type
	if expectedType, ok := schema["type"].(string); ok {
		if err := v.validateType(fieldName, value, expectedType); err != nil {
			return err
		}
	}

	// Check enum
	if enum, ok := schema["enum"].([]interface{}); ok {
		if err := v.validateEnum(fieldName, value, enum); err != nil {
			return err
		}
	}

	return nil
}

// validateType checks if the value matches the expected JSON Schema type.
func (v *SchemaValidator) validateType(fieldName string, value interface{}, expectedType string) error {
	var valid bool

	switch expectedType {
	case "string":
		switch value.(type) {
		case string:
			valid = true
		case int, int32, int64, float32, float64:
			// 允许数字类型，在Execute中会自动转换为字符串
			valid = true
		default:
			valid = false
		}
	case "integer":
		switch val := value.(type) {
		case int, int32, int64:
			valid = true
		case float64:
			// JSON numbers are decoded as float64, check if it's actually an integer
			valid = val == float64(int64(val))
		case float32:
			valid = val == float32(int32(val))
		case string:
			// 允许字符串形式的整数，在Execute中会自动转换
			valid = true
		default:
			valid = false
		}
	case "number":
		switch value.(type) {
		case int, int32, int64, float32, float64:
			valid = true
		case string:
			// 允许字符串形式的数字
			valid = true
		default:
			valid = false
		}
	case "boolean":
		_, valid = value.(bool)
	case "array":
		_, valid = value.([]interface{})
	case "object":
		_, valid = value.(map[string]interface{})
	case "null":
		valid = value == nil
	default:
		// Unknown type, skip validation
		valid = true
	}

	if !valid {
		return &ValidationError{
			Field:  fieldName,
			Reason: fmt.Sprintf("expected type '%s', got '%T'", expectedType, value),
		}
	}

	return nil
}

// validateEnum checks if the value is one of the allowed enum values.
func (v *SchemaValidator) validateEnum(fieldName string, value interface{}, enum []interface{}) error {
	for _, allowed := range enum {
		if value == allowed {
			return nil
		}
	}

	// Build allowed values string for error message
	var allowedStrs []string
	for _, e := range enum {
		allowedStrs = append(allowedStrs, fmt.Sprintf("%v", e))
	}

	return &ValidationError{
		Field:  fieldName,
		Reason: fmt.Sprintf("value must be one of: %s", strings.Join(allowedStrs, ", ")),
	}
}

// ValidateArgs is a convenience method that creates a schema and validates args.
func ValidateArgs(args map[string]interface{}, required []string, types map[string]string) error {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
	}

	// Add required fields
	if len(required) > 0 {
		reqInterface := make([]interface{}, len(required))
		for i, r := range required {
			reqInterface[i] = r
		}
		schema["required"] = reqInterface
	}

	// Add property types
	props := schema["properties"].(map[string]interface{})
	for field, typeName := range types {
		props[field] = map[string]interface{}{
			"type": typeName,
		}
	}

	validator := NewSchemaValidator()
	return validator.Validate(schema, args)
}
