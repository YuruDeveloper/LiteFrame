// Package Middleware는 HTTP 미들웨어 시스템을 위한 인터페이스를 정의합니다.
// 요청 처리 전후에 공통 로직을 실행할 수 있는 미들웨어 패턴을 지원합니다.
package Middleware

import (
	"LiteFrame/Router/Types"
)

// Middleware는 미들웨어 구현체가 따라야 하는 인터페이스입니다.
// 모든 미들웨어는 MiddleWareFunc을 반환하는 GetHandler 메서드를 구현해야 합니다.
//
// 사용 예시:
//   type LoggingMiddleware struct{}
//   func (m LoggingMiddleware) GetHandler() MiddleWareFunc {
//       return func(next http.Handler) http.Handler {
//           return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//               log.Printf("Request: %s %s", r.Method, r.URL.Path)
//               next.ServeHTTP(w, r)
//           })
//       }
//   }
type Middleware interface {
	GetHandler() MiddleWareFunc // 미들웨어 함수를 반환하는 메서드
}

// MiddleWareFunc는 미들웨어 함수 타입입니다.
// 다음 핸들러를 받아 래핑된 핸들러를 반환하는 고차 함수입니다.
// 이 패턴을 통해 여러 미들웨어를 체인형태로 연결할 수 있습니다.
type MiddleWareFunc func(Types.HandlerFunc) Types.HandlerFunc
