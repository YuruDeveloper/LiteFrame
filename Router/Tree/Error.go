// Package Tree의 에러 처리 기능을 제공합니다.
// 트리 조작 중 발생하는 에러를 처리하기 위한 사용자 정의 에러 타입을 제공합니다.
package Tree

import "fmt"

// TreeError는 트리 조작 중 발생하는 에러를 나타내는 사용자 정의 에러 타입입니다.
// 에러 메시지와 문제가 발생한 경로 정보를 포함합니다.
type TreeError struct {
	Message string // 에러 메시지
	Path    string // 에러가 발생한 경로
}

// Error는 error 인터페이스를 구현하여 에러 메시지를 반환합니다.
// 에러 메시지와 발생 경로를 결합하여 상세한 에러 정보를 제공합니다.
func (Instance *TreeError) Error() string {
	return fmt.Sprintf("error: %s From %s", Instance.Message, Instance.Path)
}

// WithMessage는 동일한 경로를 가지지만 다른 메시지를 가진 새로운 TreeError를 생성합니다.
// 기존 에러의 경로 정보를 유지하며 메시지만 변경할 때 사용합니다.
func (Instance *TreeError) WithMessage(message string) error {
	return &TreeError{
		Message: message,
		Path:    Instance.Path,
	}
}

// NewTreeError는 새로운 TreeError 인스턴스를 생성합니다.
// 에러 메시지와 문제가 발생한 경로를 받아 TreeError를 생성합니다.
func NewTreeError(Message string, path string) error {
	return &TreeError{
		Message: Message,
		Path:    path,
	}
}
