// Package Tree의 상수와 타입 정의를 포함합니다.
// 노드 타입, HTTP 메서드, 우선순위 등에 대한 열거형과 상수를 정의합니다.
package Tree

import (
	"LiteFrame/Router/Types"
)

// HandlerFunc는 Types 패키지에서 가져온 핸들러 함수 타입입니다.
type HandlerFunc = Types.HandlerFunc

// NodeType은 트리 노드의 타입을 나타내는 열거형입니다.
type NodeType uint32

// MethodType은 HTTP 메서드를 나타내는 열거형입니다.
type MethodType uint32

type ErrorCode string

// NodeType 상수들: 트리 노드의 다양한 타입을 정의합니다.
const (
	RootType       = iota // 루트 노드 ("/" 경로)
	StaticType            // 정적 경로 노드 (/users, /api 등)
	CatchAllType          // 캐치올 노드 (*path, 나머지 모든 경로 매칭)
	WildCardType          // 와일드카드 노드 (:param, 단일 세그먼트 매칭)
	MiddlewareType        // 미들웨어 노드
)

// 경로 패턴 상수들: URL 경로에서 사용되는 특수 문자들을 정의합니다.
const (
	WildCardPrefix = ':' // 와일드카드 매개변수 접두사 (:id, :name 등)
	CatchAllPrefix = '*' // 캐치올 매개변수 접두사 (*path, *file 등)
	PathSeparator  = '/' // 경로 구분자
)

// HTTP 메서드 상수들: 지원되는 HTTP 메서드들을 정의합니다.
const (
	GET        = iota // GET 메서드 - 데이터 조회
	HEAD              // HEAD 메서드 - 헤더 정보만 조회
	OPTIONS           // OPTIONS 메서드 - 지원되는 메서드 조회
	TRACE             // TRACE 메서드 - 요청 경로 추적
	POST              // POST 메서드 - 데이터 생성
	PUT               // PUT 메서드 - 데이터 수정/생성
	DELETE            // DELETE 메서드 - 데이터 삭제
	CONNECT           // CONNECT 메서드 - 터널 연결
	PATCH             // PATCH 메서드 - 데이터 부분 수정
	NotAllowed        // 지원되지 않는 메서드
)

const (
	
)
