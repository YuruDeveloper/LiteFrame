package Tree

import (
	"LiteFrame/Router/Middleware"
	"net/http"
	"strings"
)




type Tree struct {
	RootNode Node
	NotFoundHandler http.HandlerFunc
	NotAllowedHandler http.HandlerFunc
	Middlewares []Middleware.Middleware
}

func NewTree() Tree {
	return Tree{
		RootNode : NewNode(RootType,"/"),
	}
}

func (Instance *Tree) IsWildCard(Input string) bool {
	return len(Input) > 0 && Input[:1] == ":"
}

func (Instance *Tree) IsCatchAll(Input string) bool {
	return Input[0] == '*'
}

func (Instance *Tree) SplitPath(Path string) []string {
	Segments := make([]string,0,strings.Count(Path,"/")+1) 
	for Segment := range strings.SplitSeq(Path, "/") {
		if Segment != "" {
			Segments = append(Segments, Segment)
		}
	}
	return Segments
}

func (Instance *Tree) Match(One string, Two string) (bool , int , string){
	Lenght := len(One)
	var Index int
	for Index = range Lenght {
		if One[Index]  != Two[Index] {
			break
		}
	}
	Matched := Index == len(One)
	return Matched , Index , One[Index:]
} 

func (Instance *Tree) SetHandler(Method MethodType, Path string, Handler http.HandlerFunc) error {
	if Method == "" || Path == "" || (Handler == nil  && Method != CONNECT) {
		return NewTreeError("Invalid Parameters , Path and HJandler are Reqired","/")
	}
	Paths := Instance.SplitPath(Path)
	if len(Paths) == 0 {
		Instance.RootNode.Handlers[Method] = Handler
		return nil
	}
	return Instance.SetHelper(Instance.RootNode, Paths, Method, Handler)
}

func (Instance *Tree) SetHelper(Parent Node,Paths []string,Method MethodType,Handler http.HandlerFunc) error {
	if len(Paths) == 0 {
		return nil
	}
	if len(Parent.Children)  != 0 {

	}
	
	return nil
}






// AddMiddleware 함수 - 전체 트리 또는 특정 경로에 미들웨어를 추가하는 함수
// 의사코드:
// 1. 파라미터 검증
//    - Middleware가 nil인지 확인
//    - nil이면 에러 반환
// 2. 미들웨어 적용 방식 결정
//    - 전역 미들웨어: 모든 경로에 적용
//    - 경로별 미들웨어: 특정 경로에만 적용
// 3. 전역 미들웨어 적용의 경우:
//    - Root 노드에 미들웨어 추가
//    - 모든 자식 노드에 재귀적으로 미들웨어 적용
// 4. 경로별 미들웨어 적용의 경우:
//    - 경로를 분할하여 세그먼트로 변환
//    - 해당 경로의 노드를 찾아 미들웨어 추가
//    - 경로가 존재하지 않으면 미들웨어 노드 생성
// 5. 미들웨어 체인 구성
//    - 기존 미들웨어와 새 미들웨어를 체인으로 연결
//    - 실행 순서: 먼저 추가된 미들웨어가 먼저 실행
// 6. 노드 타입별 처리
//    - HandlerNode: 기존 핸들러를 미들웨어로 래핑
//    - ContainerNode: 자식 노드들에 미들웨어 전파
//    - MiddlewareNode: 미들웨어 체인에 추가
func (Tree *Tree) AddMiddleware(Middleware Middleware.Middleware) {
	// 1. 파라미터 검증
	if Middleware == nil {
		return // 에러 처리 필요
	}
	
	// 2. 전역 미들웨어로 적용
	// Root 노드부터 시작하여 모든 노드에 미들웨어 적용
	Tree.ApplyMiddlewareToAllNodes(Tree.RootNode, Middleware)
}

// ApplyMiddlewareToAllNodes - 모든 노드에 미들웨어를 재귀적으로 적용하는 헬퍼 함수
// 의사코드:
// 1. 현재 노드가 MiddlewareNode인지 확인
//    - MiddlewareNode라면 미들웨어 추가
// 2. 현재 노드가 ContainerNode인지 확인
//    - ContainerNode라면 모든 자식 노드에 재귀 호출
// 3. 현재 노드가 HandlerNode인지 확인
//    - HandlerNode라면 기존 핸들러를 미들웨어로 래핑
func (Tree *Tree) ApplyMiddlewareToAllNodes(Node Node, Middleware Middleware.Middleware) {
	// 구현 예정
}














