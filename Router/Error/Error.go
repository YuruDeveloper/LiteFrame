package Error

import "fmt"

type ErrorCode uint32

const (
	// General errors
	InvalidParameter ErrorCode = iota + 1000
	NilParameter
	InvalidMethod
	InvalidHandler
	
	// Tree structure errors
	SplitFailed       ErrorCode = iota + 2000
	NodeNotFound
	PathTooLong
	InvalidSplitPoint
	
	// Route conflicts
	DuplicateWildCard ErrorCode = iota + 3000
	DuplicateCatchAll
	ConflictingRoute
	
	// Runtime errors  
	HandlerNotFound   ErrorCode = iota + 4000
	MethodNotAllowed
	ParameterMissing
)

// LiteFrameError는 LiteFrame에서 발생하는 에러를 나타내는 구조체입니다.
type LiteFrameError struct {
	Code    ErrorCode
	Message string
	Path    string
}

// Error는 error 인터페이스를 구현합니다.
func (e *LiteFrameError) Error() string {
	return fmt.Sprintf("LiteFrame Error [%d]: %s (Path: %s)", e.Code, e.Message, e.Path)
}

// NewError는 새로운 LiteFrameError를 생성합니다.
func NewError(code ErrorCode, message string, path string) error {
	return &LiteFrameError{
		Code:    code,
		Message: message,
		Path:    path,
	}
}

// GetErrorMessage는 에러 코드에 따른 기본 메시지를 반환합니다.
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

// NewErrorWithCode는 에러 코드만으로 새로운 에러를 생성합니다.
func NewErrorWithCode(code ErrorCode, path string) error {
	return NewError(code, GetErrorMessage(code), path)
}