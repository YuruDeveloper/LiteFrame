// Package Types는 Router 시스템에서 사용되는 공통 타입들을 정의합니다.
// 순환 참조를 방지하기 위해 HandlerFunc와 기타 공통 타입들을 별도 패키지로 분리합니다.
package Types

import (
	"LiteFrame/Router/Param"
	"net/http"
)

// HandlerFunc는 HTTP 핸들러 함수 타입입니다.
// 매개변수는 GetHandler에서 컨텍스트를 통해 전달됩니다.
type HandlerFunc func(http.ResponseWriter, *http.Request, *Param.Params)