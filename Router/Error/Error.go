// Package Error는 LiteFrame에서 발생하는 모든 에러 타입과 처리 기능을 제공합니다.
// 에러 코드 기반의 구조화된 에러 관리와 상세한 에러 메시지를 지원합니다.
package Error

import "fmt"

// ErrorCode는 LiteFrame에서 발생하는 에러의 타입을 구분하는 열거형입니다.
// 에러를 카테고리별로 분류하여 디버깅과 문제 해결을 용이하게 합니다.
type ErrorCode uint32

// 에러 코드 상수들: 카테고리별로 1000 단위로 구분하여 관리합니다.
const (
	// 일반적인 매개변수 및 입력 검증 에러 (1000번대)
	InvalidParameter ErrorCode = iota + 1000 // 잘못된 매개변수
	NilParameter                             // nil 또는 빈 매개변수
	InvalidMethod                            // 지원되지 않는 HTTP 메서드
	InvalidHandler                           // 잘못된 핸들러 함수

	// 트리 구조 관련 에러 (2000번대)
	SplitFailed       ErrorCode = iota + 2000 // 노드 분할 실패
	NodeNotFound                              // 노드를 찾을 수 없음
	PathTooLong                               // 경로가 너무 긴 경우
	InvalidSplitPoint                         // 잘못된 분할 지점

	// 라우트 충돌 에러 (3000번대)
	DuplicateWildCard ErrorCode = iota + 3000 // 중복된 와일드카드 노드
	DuplicateCatchAll                         // 중복된 캐치올 노드
	ConflictingRoute                          // 기존 라우트와 충돌

	// 런타임 에러 (4000번대)
	HandlerNotFound  ErrorCode = iota + 4000 // 핸들러를 찾을 수 없음
	MethodNotAllowed                         // 허용되지 않는 HTTP 메서드
	ParameterMissing                         // 필수 매개변수 누락
)

// LiteFrameError는 LiteFrame에서 발생하는 구조화된 에러를 나타내는 구조체입니다.
// 에러 코드, 메시지, 발생 경로를 포함하여 디버깅을 위한 상세 정보를 제공합니다.
type LiteFrameError struct {
	Code    ErrorCode // 에러 분류를 위한 코드
	Message string    // 에러에 대한 상세 설명
	Path    string    // 에러가 발생한 라우트 경로
}

// Error는 error 인터페이스를 구현합니다.
// 형식화된 에러 메시지를 반환하여 로깅과 디버깅에 활용합니다.
func (e *LiteFrameError) Error() string {
	return fmt.Sprintf("LiteFrame Error [%d]: %s (Path: %s)", e.Code, e.Message, e.Path)
}

// NewError는 새로운 LiteFrameError를 생성합니다.
// 사용자 정의 메시지와 함께 에러를 생성할 때 사용합니다.
func NewError(code ErrorCode, message string, path string) error {
	return &LiteFrameError{
		Code:    code,
		Message: message,
		Path:    path,
	}
}

// GetErrorMessage는 에러 코드에 따른 기본 메시지를 반환합니다.
// 표준화된 에러 메시지를 제공하여 일관성 있는 에러 처리를 지원합니다.
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
