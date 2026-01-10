package protocol

import (
	"encoding/json"
	"testing"
)

func TestParseRequest_ValidRequest(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantID   interface{}
		wantMeth string
	}{
		{
			name:     "request with integer id",
			input:    `{"jsonrpc":"2.0","id":1,"method":"test"}`,
			wantID:   float64(1),
			wantMeth: "test",
		},
		{
			name:     "request with string id",
			input:    `{"jsonrpc":"2.0","id":"abc","method":"test/method"}`,
			wantID:   "abc",
			wantMeth: "test/method",
		},
		{
			name:     "request with params",
			input:    `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"test"}}`,
			wantID:   float64(2),
			wantMeth: "tools/call",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := ParseRequest([]byte(tt.input))
			if err != nil {
				t.Fatalf("ParseRequest() error = %v", err)
			}
			if req.ID != tt.wantID {
				t.Errorf("ID = %v, want %v", req.ID, tt.wantID)
			}
			if req.Method != tt.wantMeth {
				t.Errorf("Method = %v, want %v", req.Method, tt.wantMeth)
			}
		})
	}
}

func TestParseRequest_Notification(t *testing.T) {
	input := `{"jsonrpc":"2.0","method":"initialized"}`
	req, err := ParseRequest([]byte(input))
	if err != nil {
		t.Fatalf("ParseRequest() error = %v", err)
	}
	if !req.IsNotification() {
		t.Error("IsNotification() = false, want true")
	}
}

func TestParseRequest_InvalidJSON(t *testing.T) {
	input := `{invalid json}`
	_, err := ParseRequest([]byte(input))
	if err == nil {
		t.Fatal("ParseRequest() expected error for invalid JSON")
	}
	if protoErr, ok := err.(*Error); ok {
		if protoErr.Code != ParseError {
			t.Errorf("error code = %d, want %d", protoErr.Code, ParseError)
		}
	} else {
		t.Error("expected *Error type")
	}
}

func TestParseRequest_MissingMethod(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1}`
	_, err := ParseRequest([]byte(input))
	if err == nil {
		t.Fatal("ParseRequest() expected error for missing method")
	}
	if protoErr, ok := err.(*Error); ok {
		if protoErr.Code != InvalidRequest {
			t.Errorf("error code = %d, want %d", protoErr.Code, InvalidRequest)
		}
	}
}

func TestParseRequest_WrongVersion(t *testing.T) {
	input := `{"jsonrpc":"1.0","id":1,"method":"test"}`
	_, err := ParseRequest([]byte(input))
	if err == nil {
		t.Fatal("ParseRequest() expected error for wrong version")
	}
}

func TestBuildSuccessResponse(t *testing.T) {
	result := map[string]string{"key": "value"}
	resp := BuildSuccessResponse(1, result)

	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %v, want 2.0", resp.JSONRPC)
	}
	if resp.ID != 1 {
		t.Errorf("ID = %v, want 1", resp.ID)
	}
	if resp.Error != nil {
		t.Error("Error should be nil for success response")
	}

	// Verify result can be marshaled
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}
	if len(data) == 0 {
		t.Error("marshaled data is empty")
	}
}

func TestBuildErrorResponse(t *testing.T) {
	resp := BuildErrorResponse(1, InvalidParams, "invalid params", "details")

	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %v, want 2.0", resp.JSONRPC)
	}
	if resp.ID != 1 {
		t.Errorf("ID = %v, want 1", resp.ID)
	}
	if resp.Result != nil {
		t.Error("Result should be nil for error response")
	}
	if resp.Error == nil {
		t.Fatal("Error should not be nil for error response")
	}
	if resp.Error.Code != InvalidParams {
		t.Errorf("Error.Code = %d, want %d", resp.Error.Code, InvalidParams)
	}
}

func TestMapHTTPStatusToError(t *testing.T) {
	tests := []struct {
		status   int
		wantCode int
	}{
		{400, InvalidParams},
		{401, Unauthorized},
		{403, Forbidden},
		{404, NotFound},
		{409, Conflict},
		{500, InternalError},
		{503, InternalError},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.status)), func(t *testing.T) {
			code, _ := MapHTTPStatusToError(tt.status)
			if code != tt.wantCode {
				t.Errorf("MapHTTPStatusToError(%d) = %d, want %d", tt.status, code, tt.wantCode)
			}
		})
	}
}
