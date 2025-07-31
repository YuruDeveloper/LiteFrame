// Package Error provides all error types and handling functionality that occur in LiteFrame.
// Supports structured error management based on error codes and detailed error messages.
package Error

import "fmt"

// ErrorCode is an enumeration that distinguishes the types of errors that occur in LiteFrame.
// Categorizes errors by category to facilitate debugging and problem solving.
type ErrorCode uint32

// Error code constants: Managed by category in units of 1000.
const (
	// General parameter and input validation errors (1000s)
	InvalidParameter ErrorCode = iota + 1000 // Invalid parameter
	NilParameter                             // nil or empty parameter
	InvalidMethod                            // Unsupported HTTP method
	InvalidHandler                           // Invalid handler function

	// Tree structure related errors (2000s)
	SplitFailed       ErrorCode = iota + 2000 // Node split failure
	NodeNotFound                              // Node not found
	PathTooLong                               // Path is too long
	InvalidSplitPoint                         // Invalid split point

	// Route conflict errors (3000s)
	DuplicateWildCard ErrorCode = iota + 3000 // Duplicate wildcard node
	DuplicateCatchAll                         // Duplicate catch-all node
	ConflictingRoute                          // Conflicts with existing route

	// Runtime errors (4000s)
	HandlerNotFound  ErrorCode = iota + 4000 // Handler not found
	MethodNotAllowed                         // HTTP method not allowed
	ParameterMissing                         // Required parameter missing
)

// LiteFrameError is a structure representing structured errors that occur in LiteFrame.
// Provides detailed information for debugging including error code, message, and occurrence path.
type LiteFrameError struct {
	Code    ErrorCode // Code for error classification
	Message string    // Detailed description of the error
	Path    string    // Route path where the error occurred
}

// Error implements the error interface.
// Returns formatted error message for use in logging and debugging.
func (e *LiteFrameError) Error() string {
	return fmt.Sprintf("LiteFrame Error [%d]: %s (Path: %s)", e.Code, e.Message, e.Path)
}

// NewError creates a new LiteFrameError.
// Used when creating errors with custom messages.
func NewError(code ErrorCode, message string, path string) error {
	return &LiteFrameError{
		Code:    code,
		Message: message,
		Path:    path,
	}
}

// GetErrorMessage returns the default message according to the error code.
// Provides standardized error messages to support consistent error handling.
func GetErrorMessage(code ErrorCode) string {
	switch code {
	// General errors
	case InvalidParameter:
		return "Invalid parameter provided"
	case NilParameter:
		return "Parameter cannot be nil or empty"
	case InvalidMethod:
		return "HTTP method is not supported"
	case InvalidHandler:
		return "Handler function is required"

	// Tree structure errors
	case SplitFailed:
		return "Failed to split node"
	case NodeNotFound:
		return "Node not found in tree"
	case PathTooLong:
		return "Path exceeds maximum length"
	case InvalidSplitPoint:
		return "Split point is out of range"

	// Route conflicts
	case DuplicateWildCard:
		return "Cannot have duplicate wildcard nodes"
	case DuplicateCatchAll:
		return "Cannot have duplicate catch-all nodes"
	case ConflictingRoute:
		return "Route conflicts with existing route"

	// Runtime errors
	case HandlerNotFound:
		return "Handler not found for route"
	case MethodNotAllowed:
		return "HTTP method not allowed"
	case ParameterMissing:
		return "Required parameter is missing"

	default:
		return "Unknown error"
	}
}

// NewErrorWithCode creates a new error with only the error code.
func NewErrorWithCode(code ErrorCode, path string) error {
	return NewError(code, GetErrorMessage(code), path)
}
