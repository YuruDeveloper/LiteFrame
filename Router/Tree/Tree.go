package Tree

import (
	"LiteFrame/Router/Middleware"
	"net/http"
	"strings"
	Componet "LiteFrame/Router/Tree/Componet"

)

type NodeFactory struct {}

func (Instance *NodeFactory) IsWildCard(Input string) bool {
	return len(Input) > 0 && Input[:1] == ":"
}

func (Instance *NodeFactory) IsCatchAll(Input string) bool {
	return Input == "*"
}

func (Instance *NodeFactory) CreateNode(Path string) Componet.Node {
	switch {
		case Instance.IsWildCard(Path):
		  	Node := NewWildCardNode(Path) 
		  	return &Node 
		case Instance.IsCatchAll(Path):
			Node := NewCatchAllNode(Path)
			return &Node 
		default:
			Node := NewStaticNode(Path)
			return &Node
	}
}

func (Instance *NodeFactory) CreateHandlerNode(Path string,Method string,Handler http.HandlerFunc) (Componet.HandlerNode ,error) {
	Node := Instance.CreateNode(Path)
	HandlerNode, Ok := Node.(Componet.HandlerNode)
	if !Ok {
		return nil , Componet.NewError(Node.GetType(),"Created Node Does not implement HandlerNode",Path)
	}
	if Err := HandlerNode.SetHandler(Method,Handler); Err != nil {
		return nil , Err
	}
	return HandlerNode , nil
}



type Tree struct {
	Root RootNode
	NodeFactory NodeFactory
	NotFoundHandler http.HandlerFunc
	NotAllowedHandler http.HandlerFunc

}

func NewTree() Tree {
	return Tree{
		Root:NewRootNode(),
	}
}

func (Instance *Tree) IsWildCard(Input string) bool {
	return len(Input) > 0 && Input[:1] == ":"
}

func (Instance *Tree) IsCatchAll(Input string) bool {
	return Input == "*"
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

func (Instance *Tree) Add(Method string, Path string, Handler http.HandlerFunc) error {
	if Method == "" || Path == "" || Handler == nil {
		return Componet.NewError(Componet.RootType,"Invalid Parameters , Path and HJandler are Reqired","/")
	}
	Paths := Instance.SplitPath(Path)
	if len(Paths) == 0 {
		return Instance.AddRootHandler(Method,Handler)
	}
	return Instance.AddHelper(&Instance.Root,Paths,Method,Handler)  
}

func (Instance *Tree) AddRootHandler(Method string,Handler http.HandlerFunc) error {
	RootKey := "/"
	if Instance.Root.HasChildren() {
		if Child := Instance.Root.GetChild(RootKey); Child != nil {
			if HandlerNode , OK := Child.(Componet.HandlerNode); OK {
				return HandlerNode.SetHandler(Method,Handler)
			}
			return Componet.NewError(Componet.RootType,"Root Node is not a Handler Node","/")
		}
	}

	RootNode := NewStaticNode(RootKey)
	if Err := RootNode.SetHandler(Method,Handler); Err != nil {
		return Err
	}
	return Instance.Root.AddChild(RootKey,&RootNode)
}



func (Instance *Tree) AddHelper(Parent Componet.NodeContainer[Componet.Node],Paths []string,Method string,Handler http.HandlerFunc) error {
	if len(Paths) == 0 {
		return nil
	}
	if !Parent.HasChildren() {
		return Instance.CreateNewNode(Parent,Paths,Method,Handler)
	}
	Matched , Err := Instance.TryMatch(Parent,Paths,Method,Handler)
	if Err != nil {
		return Err
	}
	if !Matched{
		return Instance.CreateNewNode(Parent,Paths,Method,Handler)
	}
	return nil
}

func (Instance *Tree) CreateNewNode(Parent Componet.NodeContainer[Componet.Node],Paths []string,Method string,Handler http.HandlerFunc) error {
	if len(Paths) == 0 {
		return Componet.NewError(Parent.(Componet.Node).GetType(),"Path is empty","")
	}
	if len(Paths) == 1 {
		Node , Err :=  Instance.NodeFactory.CreateHandlerNode(Paths[0],Method,Handler)
		if Err != nil {
			return Err
		}
		if Err := Parent.AddChild(Paths[0],Node); Err != nil {
			return Err
		}
		return nil
	}

	Container := Instance.NodeFactory.CreateNode(Paths[0])
	if Err := Parent.AddChild(Paths[0],Container); Err != nil {
		return  Err
	}
	MewParent , Ok := Container.(Componet.NodeContainer[Componet.Node])
	if !Ok {
		return Componet.NewError(Container.GetType(),"Fail Transform NodeContainer", Paths[0])
	}
	return Instance.AddHelper(MewParent,Paths[1:],Method,Handler)
}


func (Instance *Tree) TryMatch(Parent Componet.NodeContainer[Componet.Node],Paths []string,Method string,Handler http.HandlerFunc) (bool,error) {
	for _ , Child := range Parent.GetAllChildren() {
		Matched , Err := Instance.MatchNode(Parent,Child,Paths,Method,Handler)
		if Err != nil {
			return false , Err
		}
		if Matched {
			return true , nil
		}
	}
	return false , nil
}

func (Instance *Tree) MatchNode(Parent Componet.NodeContainer[Componet.Node],Child Componet.Node,Paths []string,Method string,Handler http.HandlerFunc) (bool, error) {
	PathNode , Ok := Child.(Componet.PathNode)
	if !Ok {
		return false , Componet.NewError(Child.GetType(),"Child Node is not PathNode","")
	}
	Matched, MatchingChar , LeftPath  := PathNode.Match(Paths[0])
	if Matched {
		Err := Instance.HandleMatchedNode(Child,Paths,Method,Handler)
		return true , Err
	}
	if MatchingChar < len(PathNode.GetPath()) {
		return true , Instance.HandlePartialMatchedNode(Parent,Child,Paths,MatchingChar,LeftPath,Method,Handler)
	}
	return false , Componet.NewError(Child.GetType(),"UnExpected Error",PathNode.GetPath())
}

func (Instance *Tree) HandleMatchedNode(Child Componet.Node,Paths []string,Method string,Handler http.HandlerFunc) error {
	switch len(Paths) {
		case 1:
			HandlerNode , Ok := Child.(Componet.HandlerNode)
			if !Ok {
				return Componet.NewError(Child.GetType(),"Leaf Node is not HandlerNode","")
			}
			return  HandlerNode.SetHandler(Method,Handler)
		default: 
			Container , Ok := Child.(Componet.NodeContainer[Componet.Node])
			if !Ok {
				return Componet.NewError(Child.GetType(),"Node is not nodeContainer","")
			}
			return Instance.AddHelper(Container,Paths[1:],Method,Handler)
	}
}
func (Instance *Tree) HandlePartialMatchedNode(Parent Componet.NodeContainer[Componet.Node],Child Componet.Node,Paths []string,MatchingCHar int,LeafPath string,Method string,Handler http.HandlerFunc) error {
	StaticNode , Ok := Child.(*StaticNode)
	if !Ok {
		return Componet.NewError(Child.GetType(),"Error Child is not StaticNode","")
	}
	NewNode , Err:= Instance.SplitNode(*StaticNode,Parent,MatchingCHar,Paths,Method,Handler)
	if Err != nil {
		return Err
	}
	Paths[0] =LeafPath
	return Instance.AddHelper(NewNode,Paths,Method,Handler)
}
func (Instance *Tree) SplitNode(TargetNode StaticNode, Parent Componet.NodeContainer[Componet.Node],SplitPoint int , Paths []string,Method string,Handler http.HandlerFunc) (Componet.NodeContainer[Componet.Node] ,error) {
	// 자식 노드들 백업
	OriginalChildren := TargetNode.GetAllChildren()
	
	Path := TargetNode.GetPath()
	CommonPrefix := Path[:SplitPoint]
	RemainingSuffix := Path[SplitPoint:]

	NewParent := NewStaticNode(CommonPrefix)
	Parent.DeleteChild(TargetNode.GetPath())
	
	if Err := TargetNode.SetPath(RemainingSuffix); Err != nil {
		return nil ,Err
	}

	NewChild , Err:= TargetNode.Split(SplitPoint, NewParent.PathContainer)
	if Err != nil {
		return nil , Componet.NewError(TargetNode.GetType(),"Fail to Split Node",TargetNode.GetPath())
	}
	
	Parent.AddChild(NewChild.GetPath(), &NewParent)

	if TargetNode.GetPath() != "" {
		NewParent.AddChild(TargetNode.GetPath(), &TargetNode)
		// 원래 노드의 모든 자식들을 복원
		for _ , Child := range OriginalChildren {
			if PathChild, ok := Child.(Componet.PathNode); ok {
				TargetNode.AddChild(PathChild.GetPath(), Child)
			}
		}
	}
	Matched , _ , _ := TargetNode.Match(Paths[0])
	if Matched {
		if Err := NewParent.SetHandler(Method,Handler); Err != nil {
			return nil , Err
		}
		return nil , nil
	} 
	return &NewParent , nil
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
	Tree.ApplyMiddlewareToAllNodes(&Tree.Root, Middleware)
}

// ApplyMiddlewareToAllNodes - 모든 노드에 미들웨어를 재귀적으로 적용하는 헬퍼 함수
// 의사코드:
// 1. 현재 노드가 MiddlewareNode인지 확인
//    - MiddlewareNode라면 미들웨어 추가
// 2. 현재 노드가 ContainerNode인지 확인
//    - ContainerNode라면 모든 자식 노드에 재귀 호출
// 3. 현재 노드가 HandlerNode인지 확인
//    - HandlerNode라면 기존 핸들러를 미들웨어로 래핑
func (Tree *Tree) ApplyMiddlewareToAllNodes(Node Componet.Node, Middleware Middleware.Middleware) {
	// 구현 예정
}


// Search 함수 - Radix Tree에서 요청 경로에 맞는 핸들러를 찾는 알고리즘
// 의사코드:
// 1. 요청 경로를 세그먼트로 분할
// 2. 루트 노드부터 시작하여 DFS 탐색
// 3. 각 세그먼트에 대해 우선순위 순으로 노드 매칭:
//    a. StaticNode (High Priority): 정확한 문자열 매칭
//    b. WildcardNode (High Priority): 단일 세그먼트 매칭, 매개변수 추출
//    c. CatchAllNode (Low Priority): 나머지 모든 경로 매칭
// 4. 매칭 프로세스:
//    - 현재 노드의 경로와 요청 경로의 공통 접두사 계산
//    - 완전 매칭이면 다음 세그먼트로 이동
//    - 부분 매칭이면 매칭 실패
// 5. 모든 세그먼트 매칭 완료 시:
//    - 핸들러가 있으면 반환
//    - 핸들러가 없으면 404 처리
// 6. 와일드카드 매개변수를 컨텍스트에 저장
// 7. 미들웨어 체인 적용
func (Instance *Tree) Search(Request *http.Request) http.Handler {
	Paths := Instance.SplitPath(Request.URL.Path)
	if Paths == nil {
		Child := Instance.Root.GetChild("/")
		if Child == nil {
			return Instance.NotFoundHandler
		}
		if !Child.(Componet.HandlerNode).HasMethod(Request.Method) {
			return Instance.NotAllowedHandler
		}
		return Child.(Componet.HandlerNode).GetHandler(Request.Method)
	}
	return Instance.SearchHelper(&Instance.Root,make([]any,0),Paths,Request.Method)
}

func (Instance *Tree) SearchHelper(SearchTarget Componet.Node,ApplyList []any,Paths []string, Method string) http.Handler {
	for _ , Child := range SearchTarget.(Componet.NodeContainer[Componet.Node]).GetAllChildren() {
		switch Child.GetType() {
			case Componet.MiddlewareType:
				Temp := Child.(Componet.MiddlewareAcessor).Apply
				ApplyList = append(ApplyList, Temp)
				Result := Instance.SearchHelper(Child,ApplyList,Paths,Method)
				if Result != nil {
					return Result
				}
			case Componet.StaticType:
				Mathced , count , LeftPath  := Child.(Componet.PathNode).Match(Paths[0])
				if Mathced {
					Result := Instance.SearchHelper(Child,ApplyList,Paths[1:],Method)
					return Result
				}
				if count == 0 {
					continue
				}
				Paths[0] = LeftPath
				Result := Instance.SearchHelper(Child,ApplyList,Paths,Method)
				return Result
			case Componet.WildCardType:

		}
	}
	return nil
}












