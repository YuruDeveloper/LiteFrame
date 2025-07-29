// Package Tree의 상수와 타입 정의를 포함합니다.
// 노드 타입, HTTP 메서드, 우선순위 등에 대한 열거형과 상수를 정의합니다.
package Tree

import (
	"LiteFrame/Router/Types"
)

// HandlerFunc는 Types 패키지에서 가져온 핸들러 함수 타입입니다.
// 순환 참조를 방지하기 위해 별도 패키지에서 정의된 타입을 재사용합니다.
type HandlerFunc = Types.HandlerFunc

// NodeType은 트리 노드의 타입을 나타내는 열거형입니다.
// 각 노드는 하나의 타입만 가질 수 있으며, 타입에 따라 라우팅 동작이 결정됩니다.
type NodeType uint32

// MethodType은 HTTP 메서드를 나타내는 열거형입니다.
// 배열 인덱스로 사용되어 핸들러 배열에서 O(1) 시간복잡도로 핸들러에 접근합니다.
type MethodType uint32

// NodeType 상수들: 트리 노드의 다양한 타입을 정의합니다.
// 각 노드 타입은 라우팅 성능과 매칭 우선순위에 영향을 미칩니다.
const (
	RootType       = iota // 루트 노드 ("/" 경로, 트리의 시작점)
	StaticType            // 정적 경로 노드 (/users, /api 등, 가장 빠른 매칭)
	CatchAllType          // 캐치올 노드 (*path, 나머지 모든 경로 매칭, 낮은 우선순위)
	WildCardType          // 와일드카드 노드 (:param, 단일 세그먼트 매칭, 중간 우선순위)
	MiddlewareType        // 미들웨어 노드 (요청 처리 파이프라인)
)

// 경로 패턴 상수들: URL 경로에서 사용되는 특수 문자들을 정의합니다.
const (
	WildCardPrefix = ':' // 와일드카드 매개변수 접두사 (:id, :name 등)
	CatchAllPrefix = '*' // 캐치올 매개변수 접두사 (*path, *file 등)
	PathSeparator  = '/' // 경로 구분자
)

// HTTP 메서드 상수들: RFC 7231 표준을 따르는 HTTP 메서드들을 정의합니다.
// 각 메서드는 배열 인덱스로 사용되어 O(1) 핸들러 접근을 제공합니다.
const (
	GET        = iota // GET 메서드 - 리소스 조회 (멱등성, 안전)
	HEAD              // HEAD 메서드 - 헤더 정보만 조회 (GET과 동일하지만 바디 없음)
	OPTIONS           // OPTIONS 메서드 - 지원되는 메서드 조회 (CORS preflight)
	TRACE             // TRACE 메서드 - 요청 경로 추적 (디버깅 목적)
	POST              // POST 메서드 - 리소스 생성 (비멱등성)
	PUT               // PUT 메서드 - 리소스 생성/전체 수정 (멱등성)
	DELETE            // DELETE 메서드 - 리소스 삭제 (멱등성)
	CONNECT           // CONNECT 메서드 - 터널 연결 (프록시 서버용)
	PATCH             // PATCH 메서드 - 리소스 부분 수정 (RFC 5789)
	NotAllowed        // 지원되지 않는 메서드 (405 Method Not Allowed 응답용)
)
