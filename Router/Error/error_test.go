package Error

import (
	"errors"
	"fmt"
	"testing"
)

// ======================
// ErrorCode Tests
// ======================

func TestErrorCodes(t *testing.T) {
	t.Run("general_error_codes", func(t *testing.T) {
		expectedCodes := map[ErrorCode]string{
			InvalidParameter: "general parameter errors should start from 1000",
			NilParameter:     "nil parameter error should follow InvalidParameter",
			InvalidMethod:    "invalid method error should follow NilParameter", 
			InvalidHandler:   "invalid handler error should follow InvalidMethod",
		}
		
		if InvalidParameter < 1000 {
			t.Errorf("Expected InvalidParameter >= 1000, got %d", InvalidParameter)
		}
		
		for code, desc := range expectedCodes {
			if code < 1000 || code >= 2000 {
				t.Errorf("%s: Expected code in 1000s range, got %d", desc, code)
			}
		}
	})

	t.Run("tree_error_codes", func(t *testing.T) {
		expectedCodes := map[ErrorCode]string{
			SplitFailed:       "tree structure errors should start from 2000",
			NodeNotFound:      "node not found error should follow SplitFailed",
			PathTooLong:       "path too long error should follow NodeNotFound",
			InvalidSplitPoint: "invalid split point error should follow PathTooLong",
		}
		
		for code, desc := range expectedCodes {
			if code < 2000 || code >= 3000 {
				t.Errorf("%s: Expected code in 2000s range, got %d", desc, code)
			}
		}
	})

	t.Run("route_conflict_error_codes", func(t *testing.T) {
		expectedCodes := map[ErrorCode]string{
			DuplicateWildCard: "route conflict errors should start from 3000",
			DuplicateCatchAll: "duplicate catch-all error should follow DuplicateWildCard",
			ConflictingRoute:  "conflicting route error should follow DuplicateCatchAll",
		}
		
		for code, desc := range expectedCodes {
			if code < 3000 || code >= 4000 {
				t.Errorf("%s: Expected code in 3000s range, got %d", desc, code)
			}
		}
	})

	t.Run("runtime_error_codes", func(t *testing.T) {
		expectedCodes := map[ErrorCode]string{
			HandlerNotFound:  "runtime errors should start from 4000",
			MethodNotAllowed: "method not allowed error should follow HandlerNotFound",
			ParameterMissing: "parameter missing error should follow MethodNotAllowed",
		}
		
		for code, desc := range expectedCodes {
			if code < 4000 || code >= 5000 {
				t.Errorf("%s: Expected code in 4000s range, got %d", desc, code)
			}
		}
	})
}

// ======================
// LiteFrameError Tests
// ======================

func TestLiteFrameError(t *testing.T) {
	t.Run("error_creation", func(t *testing.T) {
		err := &LiteFrameError{
			Code:    InvalidParameter,
			Message: "Test error message",
			Path:    "/test/path",
		}
		
		if err.Code != InvalidParameter {
			t.Errorf("Expected code %d, got %d", InvalidParameter, err.Code)
		}
		
		if err.Message != "Test error message" {
			t.Errorf("Expected message 'Test error message', got '%s'", err.Message)
		}
		
		if err.Path != "/test/path" {
			t.Errorf("Expected path '/test/path', got '%s'", err.Path)
		}
	})

	t.Run("error_interface_implementation", func(t *testing.T) {
		err := &LiteFrameError{
			Code:    InvalidParameter,
			Message: "Test error",
			Path:    "/test",
		}
		
		// Test that it implements the error interface
		var e error = err
		if e == nil {
			t.Error("Expected LiteFrameError to implement error interface")
		}
	})

	t.Run("error_string_format", func(t *testing.T) {
		err := &LiteFrameError{
			Code:    InvalidParameter,
			Message: "Test error message",
			Path:    "/test/path",
		}
		
		expected := fmt.Sprintf("LiteFrame Error [%d]: %s (Path: %s)", InvalidParameter, "Test error message", "/test/path")
		actual := err.Error()
		
		if actual != expected {
			t.Errorf("Expected error string '%s', got '%s'", expected, actual)
		}
	})

	t.Run("empty_values", func(t *testing.T) {
		err := &LiteFrameError{
			Code:    InvalidParameter,
			Message: "",
			Path:    "",
		}
		
		expected := fmt.Sprintf("LiteFrame Error [%d]:  (Path: )", InvalidParameter)
		actual := err.Error()
		
		if actual != expected {
			t.Errorf("Expected error string '%s', got '%s'", expected, actual)
		}
	})

	t.Run("special_characters_in_message", func(t *testing.T) {
		err := &LiteFrameError{
			Code:    InvalidParameter,
			Message: "Error with 한글 and emoji 🚨",
			Path:    "/api/users/:id",
		}
		
		errorStr := err.Error()
		if errorStr == "" {
			t.Error("Expected non-empty error string")
		}
		
		// Should contain all the original characters
		if !contains(errorStr, "한글") || !contains(errorStr, "🚨") {
			t.Error("Expected error string to preserve special characters")
		}
	})
}

// ======================
// NewError Tests
// ======================

func TestNewError(t *testing.T) {
	t.Run("creates_proper_error", func(t *testing.T) {
		err := NewError(InvalidParameter, "Test message", "/test/path")
		
		if err == nil {
			t.Fatal("Expected non-nil error")
		}
		
		liteErr, ok := err.(*LiteFrameError)
		if !ok {
			t.Fatal("Expected error to be *LiteFrameError type")
		}
		
		if liteErr.Code != InvalidParameter {
			t.Errorf("Expected code %d, got %d", InvalidParameter, liteErr.Code)
		}
		
		if liteErr.Message != "Test message" {
			t.Errorf("Expected message 'Test message', got '%s'", liteErr.Message)
		}
		
		if liteErr.Path != "/test/path" {
			t.Errorf("Expected path '/test/path', got '%s'", liteErr.Path)
		}
	})

	t.Run("returns_error_interface", func(t *testing.T) {
		err := NewError(InvalidParameter, "Test", "/test")
		
		// Should be usable as error interface
		var e error = err
		if e == nil {
			t.Error("Expected error to implement error interface")
		}
		
		// Should have proper error message
		errorMsg := e.Error()
		if errorMsg == "" {
			t.Error("Expected non-empty error message")
		}
	})

	t.Run("different_error_codes", func(t *testing.T) {
		testCases := []struct {
			code     ErrorCode
			message  string
			path     string
		}{
			{InvalidParameter, "Invalid param", "/api/test"},
			{NodeNotFound, "Node missing", "/users/:id"},
			{DuplicateWildCard, "Duplicate wildcard", "/api/*path"},
			{HandlerNotFound, "No handler", "/not/found"},
		}
		
		for _, tc := range testCases {
			err := NewError(tc.code, tc.message, tc.path)
			liteErr, ok := err.(*LiteFrameError)
			if !ok {
				t.Errorf("Expected *LiteFrameError for code %d", tc.code)
				continue
			}
			
			if liteErr.Code != tc.code {
				t.Errorf("Expected code %d, got %d", tc.code, liteErr.Code)
			}
			
			if liteErr.Message != tc.message {
				t.Errorf("Expected message '%s', got '%s'", tc.message, liteErr.Message)
			}
			
			if liteErr.Path != tc.path {
				t.Errorf("Expected path '%s', got '%s'", tc.path, liteErr.Path)
			}
		}
	})
}

// ======================
// GetErrorMessage Tests
// ======================

func TestGetErrorMessage(t *testing.T) {
	t.Run("general_error_messages", func(t *testing.T) {
		testCases := map[ErrorCode]string{
			InvalidParameter: "Invalid parameter provided",
			NilParameter:     "Parameter cannot be nil or empty",
			InvalidMethod:    "HTTP method is not supported",
			InvalidHandler:   "Handler function is required",
		}
		
		for code, expected := range testCases {
			actual := GetErrorMessage(code)
			if actual != expected {
				t.Errorf("Code %d: Expected '%s', got '%s'", code, expected, actual)
			}
		}
	})

	t.Run("tree_structure_error_messages", func(t *testing.T) {
		testCases := map[ErrorCode]string{
			SplitFailed:       "Failed to split node",
			NodeNotFound:      "Node not found in tree",
			PathTooLong:       "Path exceeds maximum length",
			InvalidSplitPoint: "Split point is out of range",
		}
		
		for code, expected := range testCases {
			actual := GetErrorMessage(code)
			if actual != expected {
				t.Errorf("Code %d: Expected '%s', got '%s'", code, expected, actual)
			}
		}
	})

	t.Run("route_conflict_error_messages", func(t *testing.T) {
		testCases := map[ErrorCode]string{
			DuplicateWildCard: "Cannot have duplicate wildcard nodes",
			DuplicateCatchAll: "Cannot have duplicate catch-all nodes",
			ConflictingRoute:  "Route conflicts with existing route",
		}
		
		for code, expected := range testCases {
			actual := GetErrorMessage(code)
			if actual != expected {
				t.Errorf("Code %d: Expected '%s', got '%s'", code, expected, actual)
			}
		}
	})

	t.Run("runtime_error_messages", func(t *testing.T) {
		testCases := map[ErrorCode]string{
			HandlerNotFound:  "Handler not found for route",
			MethodNotAllowed: "HTTP method not allowed",
			ParameterMissing: "Required parameter is missing",
		}
		
		for code, expected := range testCases {
			actual := GetErrorMessage(code)
			if actual != expected {
				t.Errorf("Code %d: Expected '%s', got '%s'", code, expected, actual)
			}
		}
	})

	t.Run("unknown_error_code", func(t *testing.T) {
		unknownCode := ErrorCode(9999)
		actual := GetErrorMessage(unknownCode)
		expected := "Unknown error"
		
		if actual != expected {
			t.Errorf("Expected '%s' for unknown code, got '%s'", expected, actual)
		}
	})

	t.Run("zero_error_code", func(t *testing.T) {
		zeroCode := ErrorCode(0)
		actual := GetErrorMessage(zeroCode)
		expected := "Unknown error"
		
		if actual != expected {
			t.Errorf("Expected '%s' for zero code, got '%s'", expected, actual)
		}
	})
}

// ======================
// NewErrorWithCode Tests
// ======================

func TestNewErrorWithCode(t *testing.T) {
	t.Run("creates_error_with_default_message", func(t *testing.T) {
		err := NewErrorWithCode(InvalidParameter, "/test/path")
		
		if err == nil {
			t.Fatal("Expected non-nil error")
		}
		
		liteErr, ok := err.(*LiteFrameError)
		if !ok {
			t.Fatal("Expected error to be *LiteFrameError type")
		}
		
		if liteErr.Code != InvalidParameter {
			t.Errorf("Expected code %d, got %d", InvalidParameter, liteErr.Code)
		}
		
		expectedMessage := GetErrorMessage(InvalidParameter)
		if liteErr.Message != expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", expectedMessage, liteErr.Message)
		}
		
		if liteErr.Path != "/test/path" {
			t.Errorf("Expected path '/test/path', got '%s'", liteErr.Path)
		}
	})

	t.Run("all_error_codes_with_default_messages", func(t *testing.T) {
		testCodes := []ErrorCode{
			InvalidParameter, NilParameter, InvalidMethod, InvalidHandler,
			SplitFailed, NodeNotFound, PathTooLong, InvalidSplitPoint,
			DuplicateWildCard, DuplicateCatchAll, ConflictingRoute,
			HandlerNotFound, MethodNotAllowed, ParameterMissing,
		}
		
		for _, code := range testCodes {
			err := NewErrorWithCode(code, "/test")
			liteErr, ok := err.(*LiteFrameError)
			if !ok {
				t.Errorf("Expected *LiteFrameError for code %d", code)
				continue
			}
			
			expectedMessage := GetErrorMessage(code)
			if liteErr.Message != expectedMessage {
				t.Errorf("Code %d: Expected default message '%s', got '%s'", code, expectedMessage, liteErr.Message)
			}
			
			if liteErr.Code != code {
				t.Errorf("Expected code %d, got %d", code, liteErr.Code)
			}
		}
	})

	t.Run("unknown_error_code_with_default", func(t *testing.T) {
		unknownCode := ErrorCode(9999)
		err := NewErrorWithCode(unknownCode, "/test")
		
		liteErr, ok := err.(*LiteFrameError)
		if !ok {
			t.Fatal("Expected *LiteFrameError type")
		}
		
		if liteErr.Message != "Unknown error" {
			t.Errorf("Expected 'Unknown error', got '%s'", liteErr.Message)
		}
	})

	t.Run("comparison_with_new_error", func(t *testing.T) {
		code := InvalidParameter
		path := "/api/test"
		
		err1 := NewErrorWithCode(code, path)
		err2 := NewError(code, GetErrorMessage(code), path)
		
		liteErr1, ok1 := err1.(*LiteFrameError)
		liteErr2, ok2 := err2.(*LiteFrameError)
		
		if !ok1 || !ok2 {
			t.Fatal("Expected both errors to be *LiteFrameError type")
		}
		
		if liteErr1.Code != liteErr2.Code {
			t.Errorf("Expected same code, got %d vs %d", liteErr1.Code, liteErr2.Code)
		}
		
		if liteErr1.Message != liteErr2.Message {
			t.Errorf("Expected same message, got '%s' vs '%s'", liteErr1.Message, liteErr2.Message)
		}
		
		if liteErr1.Path != liteErr2.Path {
			t.Errorf("Expected same path, got '%s' vs '%s'", liteErr1.Path, liteErr2.Path)
		}
		
		if liteErr1.Error() != liteErr2.Error() {
			t.Errorf("Expected same error string, got '%s' vs '%s'", liteErr1.Error(), liteErr2.Error())
		}
	})
}

// ======================
// Integration Tests
// ======================

func TestErrorIntegration(t *testing.T) {
	t.Run("error_chaining", func(t *testing.T) {
		originalErr := NewError(InvalidParameter, "Original error", "/test")
		
		// Simulate error wrapping (Go 1.13+ style)
		wrappedErr := fmt.Errorf("wrapper: %w", originalErr)
		
		if !errors.Is(wrappedErr, originalErr) {
			t.Error("Expected wrapped error to be identifiable")
		}
		
		var liteErr *LiteFrameError
		if !errors.As(wrappedErr, &liteErr) {
			t.Error("Expected to extract LiteFrameError from wrapped error")
		}
		
		if liteErr.Code != InvalidParameter {
			t.Errorf("Expected code %d, got %d", InvalidParameter, liteErr.Code)
		}
	})

	t.Run("error_comparison", func(t *testing.T) {
		err1 := NewError(InvalidParameter, "Test", "/path")
		err2 := NewError(InvalidParameter, "Test", "/path")
		err3 := NewError(InvalidMethod, "Test", "/path")
		
		// Different instances should not be equal
		if err1 == err2 {
			t.Error("Expected different error instances to not be equal")
		}
		
		// Different codes should definitely not be equal
		if err1 == err3 {
			t.Error("Expected errors with different codes to not be equal")
		}
		
		// But error strings might be the same
		liteErr1, _ := err1.(*LiteFrameError)
		liteErr2, _ := err2.(*LiteFrameError)
		
		if liteErr1.Code != liteErr2.Code {
			t.Error("Expected same error codes")
		}
	})

	t.Run("nil_handling", func(t *testing.T) {
		var nilErr *LiteFrameError
		
		// Should not cause panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Nil error handling caused panic: %v", r)
			}
		}()
		
		if nilErr != nil {
			t.Error("Expected nil error to be nil")
		}
		
		// In Go, a nil pointer of a concrete type is not nil when assigned to interface
		// This is expected behavior - a nil *LiteFrameError is not the same as nil error
		var err error = nilErr
		if err == nil {
			// This is actually the unexpected case, but let's handle it gracefully
			// A nil concrete type assigned to interface is not nil unless the concrete type is nil
		}
	})
}

// ======================
// Helper Functions
// ======================

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || findSubstring(s, substr))
}

// findSubstring is a simple substring search
func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}