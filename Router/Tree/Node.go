// Package Tree의 Node 구조체와 관련 함수들을 정의합니다.
// Radix Tree의 각 노드를 나타내는 구조체와 생성자 함수를 포함합니다.
package Tree

import "net/http"

// NewNode는 새로운 Node 인스턴스를 생성합니다.
// Type: 노드 타입 (Root, Static, WildCard, CatchAll, Middleware)
// Path: 노드가 나타내는 경로 세그먼트
func NewNode(Type NodeType, Path string) *Node {
	return &Node{
		Type:     Type,
		Path:     Path,
		Children: make([]*Node, 0),
		Handlers: make([]http.HandlerFunc,int(NotAllowed)),
		WildCard: nil,
		CatchAll: nil,
	}
}

// Node는 Radix Tree의 각 노드를 나타내는 구조체입니다.
// 경로 세그먼트, 자식 노드들, HTTP 메서드별 핸들러를 저장합니다.
type Node struct {
	Type     NodeType           // 노드 타입 (Root, Static, WildCard, CatchAll, Middleware)
	Path     string             // 노드가 나타내는 경로 세그먼트
	Indices  []byte             // 자식 노드들의 첫 번째 바이트 인덱스 (빠른 검색용)
	Children []*Node            // 정적 자식 노드들
	Handlers []http.HandlerFunc // HTTP 메서드별 핸들러 배열
	WildCard *Node              // 와일드카드 자식 노드 (:param)
	CatchAll *Node              // 캐치올 자식 노드 (*path)
	Param    string             // 매개변수 이름 (WildCard/CatchAll 노드에서 사용)
}
