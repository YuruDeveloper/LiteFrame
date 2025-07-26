// Package Tree의 Node 구조체와 관련 함수들을 정의합니다.
// Radix Tree의 각 노드를 나타내는 구조체와 생성자 함수를 포함합니다.
package Tree

// NewNode는 새로운 Node 인스턴스를 생성합니다.
// Type: 노드 타입 (Root, Static, WildCard, CatchAll, Middleware)
// Path: 노드가 나타내는 경로 세그먼트
// 성능 최적화: 메모리 할당을 최소화하고 필요한 필드만 초기화합니다.
func NewNode(Type NodeType, Path string) *Node {
	return &Node{
		Type:     Type,
		Path:     Path,
		Children: make([]*Node, 0),           // 동적 확장 가능한 자식 노드 슬라이스
		Handlers: make([]HandlerFunc,int(NotAllowed)), // 모든 HTTP 메서드 지원을 위한 핸들러 배열
		WildCard: nil,                        // 와일드카드 노드는 필요시에만 생성
		CatchAll: nil,                        // 캐치올 노드는 필요시에만 생성
	}
}

// Node는 Radix Tree의 각 노드를 나타내는 구조체입니다.
// 경로 세그먼트, 자식 노드들, HTTP 메서드별 핸들러를 저장합니다.
// 
// 메모리 효율성을 위한 설계:
// - Static 노드: Children + Indices 사용 (O(1) 검색)
// - WildCard/CatchAll: 별도 포인터로 관리 (메모리 절약)
// - Handlers: 배열 인덱스를 통한 직접 접근 (성능 최적화)
type Node struct {
	Type     NodeType           // 노드 타입 (Root, Static, WildCard, CatchAll, Middleware)
	Path     string             // 노드가 나타내는 경로 세그먼트 (압축된 경로)
	Indices  []byte             // 자식 노드들의 첫 번째 바이트 인덱스 (O(1) 검색 최적화)
	Children []*Node            // 정적 자식 노드들 (Indices와 1:1 대응)
	Handlers []HandlerFunc // HTTP 메서드별 핸들러 배열 (MethodType을 인덱스로 사용)
	WildCard *Node              // 와일드카드 자식 노드 (:param, 단일 세그먼트 매칭)
	CatchAll *Node              // 캐치올 자식 노드 (*path, 나머지 모든 경로 매칭)
	Param    string             // 매개변수 이름 (WildCard/CatchAll 노드에서만 사용, ':' '*' 제외)
}
