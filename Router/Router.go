// Package Router는 고수준 HTTP 라우터 인터페이스를 제공합니다.
// 기본 라우터 구조체와 기본 핸들러들을 정의합니다.
package Router

import (
	"net/http"
)

// Router는 기본적인 HTTP 라우터 구조체입니다.
// 404(Not Found)와 405(Method Not Allowed) 에러를 처리하는 기본 핸들러를 포함합니다.
// 향후 Tree 라우터와 통합하여 완전한 라우팅 시스템을 구성할 예정입니다.
type Router struct {
	NotFoundHandler   http.HandlerFunc // 404 에러 처리 핸들러
	NotAllowedHandler http.HandlerFunc // 405 에러 처리 핸들러
}

// NotFoundDefault는 기본 404 에러 핸들러입니다.
// 요청된 리소스를 찾을 수 없을 때 표준 HTTP 404 응답을 반환합니다.
func NotFoundDefault(Writer http.ResponseWriter, Request *http.Request) {
	http.Error(Writer, "404 Not Found", http.StatusNotFound)
}

// NotAllowedDefault는 기본 405 에러 핸들러입니다.
// 지원되지 않는 HTTP 메서드로 요청이 들어왔을 때 표준 HTTP 405 응답을 반환합니다.
func NotAllowedDefault(Writer http.ResponseWriter, Request *http.Request) {
	http.Error(Writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// NewRouter는 새로운 Router 인스턴스를 생성합니다.
// 기본 에러 핸들러들을 사용하여 라우터를 초기화합니다.
// 사용자는 필요에 따라 커스텀 에러 핸들러로 교체할 수 있습니다.
func NewRouter() *Router {
	return &Router{
		NotFoundHandler:   NotFoundDefault,   // 기본 404 핸들러 설정
		NotAllowedHandler: NotAllowedDefault, // 기본 405 핸들러 설정
	}
}
