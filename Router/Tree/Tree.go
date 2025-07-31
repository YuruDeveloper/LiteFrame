// Package Tree는 HTTP 라우팅을 위한 Radix Tree 구현을 제공합니다.
// 효율적인 URL 패턴 매칭과 핸들러 관리를 위한 트리 구조를 구현합니다.
package Tree

import (
	"LiteFrame/Router/Error"
	"LiteFrame/Router/Middleware"
	"LiteFrame/Router/Param"
	"fmt"
	"net/http"
	"sort"
)

// Tree는 HTTP 라우팅을 위한 Radix Tree 구조체입니다.
// 루트 노드, 매개변수 풀, 핸들러 및 미들웨어를 관리합니다.
type Tree struct {
	RootNode          *Node                   // 트리의 루트 노드
	Pool              *Param.ParamsPool       // 매개변수 재사용을 위한 풀
	NotFoundHandler   HandlerFunc             // 404 핸들러
	NotAllowedHandler HandlerFunc             // 405 핸들러
	Middlewares       []Middleware.Middleware // 미들웨어 목록
	CompiledMiddleware func(HandlerFunc) HandlerFunc
}

// NewTree는 새로운 Tree 인스턴스를 생성합니다.
// 루트 노드와 매개변수 풀을 초기화합니다.
func NewTree() Tree {
	return Tree{
		RootNode: NewNode(RootType, "/"),
		Pool:     Param.NewParamsPool(),
	}
}

// IsWildCard는 입력 문자열이 와일드카드 패턴(:param)인지 확인합니다.
//go:inline
func (Instance *Tree) IsWildCard(Input string) bool {
	return len(Input) > 0 && Input[0] == WildCardPrefix
}

// IsCatchAll는 입력 문자열이 캐치올 패턴(*path)인지 확인합니다.
//go:inline
func (Instance *Tree) IsCatchAll(Input string) bool {
	return len(Input) > 0 && Input[0] == CatchAllPrefix
}

// StringToMethodType는 HTTP 메서드 문자열을 MethodType으로 변환합니다.
// 지원되지 않는 메서드의 경우 NotAllowed를 반환합니다.
//go:inline
func (Instance *Tree) StringToMethodType(Method string) MethodType {
	switch Method {
	case "GET":
		return GET
	case "HEAD":
		return HEAD
	case "OPTIONS":
		return OPTIONS
	case "TRACE":
		return TRACE
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	case "CONNECT":
		return CONNECT
	case "PATCH":
		return PATCH
	default:
		return NotAllowed
	}
}

// Match는 PathWithSegment와 문자열의 공통 접두사를 찾아 매칭 결과를 반환합니다.
// 반환값: (완전매칭여부, 매칭된인덱스, 남은PathWithSegment)
// PathWithSegment 기반으로 최적화된 매칭 알고리즘입니다.
//go:inline
func (Instance *Tree) Match(SourcePath PathWithSegment, TargetPath string) (bool, int, PathWithSegment) {
	// PathWithSegment와 문자열 중 짧은 길이를 기준으로 비교 범위 설정

	SourceLength := SourcePath.GetLength()
	Length := min(SourceLength, len(TargetPath))
	// 바이트 단위로 순차 비교하여 공통 접두사 길이 계산
	var Index int
	for Index = 0; Index < Length; Index++ {
		if SourcePath.Body[SourcePath.Start+Index] != TargetPath[Index] {
			break
		}
	}
	// PathWithSegment가 Two의 완전한 접두사인지 확인
	Matched := Index == SourcePath.GetLength()
	if Matched {
		SourcePath.Start = SourcePath.End
	}
	// PathWithSegment에서 매칭되지 않은 나머지 부분으로 업데이트
	if Index < SourcePath.GetLength() {
		SourcePath.Start = SourcePath.Start + Index
	}
	return Matched, Index, SourcePath
}

// SelectHandler는 노드에서 메서드에 맞는 핸들러를 선택하고 매개변수를 컨텍스트에 주입합니다.
// 핸들러가 없으면 NotAllowedHandler를 반환합니다.
// 중요: 메모리 풀 관리와 컨텍스트 주입을 동시에 처리하는 핵심 함수입니다.
func (Instance *Tree) SelectHandler(Node *Node, Method MethodType) HandlerFunc {
	if Handler := Node.Handlers[Method]; Handler != nil {
		// 클로저를 통해 매개변수를 컨텍스트에 주입하고 메모리 풀 반환 보장
		return Handler
	}
	return Instance.NotAllowedHandler
}

// InsertUniqueTypeChild는 고유한 타입의 자식 노드(WildCard/CatchAll)를 삽입합니다.
// 중복된 매개변수명이 있으면 에러를 반환합니다.
// 함수형 프로그래밍 패턴: 고차 함수를 활용한 중복 코드 제거

func (Instance *Tree) InsertUniqueTypeChild(Parent *Node, Path string, Target *Node, Type NodeType, ErrorFn func(string) error, SetFn func(*Node, *Node)) (*Node, error) {
	switch {
	// 빈 매개변수명 검증 (":" 또는 "*" 만 있는 경우)
	case Path[1:] == "":
		return nil, Error.NewErrorWithCode(Error.NilParameter, Path)
	// 같은 타입이지만 다른 매개변수명인 경우 충돌 오류
	case Target != nil && Target.Param != Path[1:]:
		return nil, ErrorFn(Path)
	// 동일한 매개변수명의 노드가 이미 존재하는 경우 재사용
	case Target != nil && Target.Param == Path[1:]:
		return Target, nil
	// 새로운 매개변수 노드 생성
	default:
		Child := NewNode(Type, Path)
		Child.Param = Path[1:] // 접두사(: 또는 *) 제거하여 매개변수명만 저장
		SetFn(Parent, Child)   // 함수 포인터를 통한 타입별 노드 설정
		return Child, nil
	}
}

// InsertChild는 부모 노드에 자식 노드를 삽입합니다.
// 경로 타입에 따라 Static, WildCard, CatchAll 노드를 생성합니다.
func (Instance *Tree) InsertChild(Parent *Node, Path string) (*Node, error) {
	switch {
	case Instance.IsWildCard(Path):
		return Instance.InsertUniqueTypeChild(Parent, Path, Parent.WildCard, WildCardType,
			func(Path string) error { return Error.NewErrorWithCode(Error.DuplicateWildCard, Path) },
			func(Parent, Child *Node) { Parent.WildCard = Child })
	case Instance.IsCatchAll(Path):
		return Instance.InsertUniqueTypeChild(Parent, Path, Parent.CatchAll, CatchAllType,
			func(Path string) error { return Error.NewErrorWithCode(Error.DuplicateCatchAll, Path) },
			func(Parent, Child *Node) { Parent.CatchAll = Child })
	default:
		Child := NewNode(StaticType, Path)
		InsertLocation := sort.Search(len(Parent.Indices), func(Index int) bool { return Parent.Indices[Index] >= Path[0] })
		Parent.Indices = append(Parent.Indices[:InsertLocation], append([]byte{Path[0]}, Parent.Indices[InsertLocation:]...)...)
		Parent.Children = append(Parent.Children[:InsertLocation], append([]*Node{Child}, Parent.Children[InsertLocation:]...)...)
		return Child, nil
	}
}

// SplitNode는 기존 노드를 분할점에서 두 개의 노드로 분리합니다.
// 공통 접두사를 가진 새로운 부모 노드를 생성하고 기존 노드를 자식으로 만듭니다.
//go:inline
func (Instance *Tree) SplitNode(Parent *Node, Child *Node, SplitPoint int) (*Node, error) {
	if SplitPoint < 0 || SplitPoint > len(Child.Path) {
		return nil, fmt.Errorf("split point %d is out of range for path %q (length %d)", SplitPoint, Child.Path, len(Child.Path))
	}
	Left := Child.Path[:SplitPoint]
	Right := Child.Path[SplitPoint:]
	if len(Left) == 0 {
		return nil, Error.NewErrorWithCode(Error.SplitFailed, Child.Path)
	}
	NewParent := NewNode(StaticType, Left)
	for Index, Target := range Parent.Children {
		if Target == Child {
			Parent.Children[Index] = NewParent
			if len(Right) > 0 {
				NewParent.Indices = []byte{Right[0]}
				Child.Path = Right
				NewParent.Children = []*Node{Child}
			} else {
				Child.Path = Right
				NewParent.Indices = []byte{}
				NewParent.Children = []*Node{}
			}
		}
	}
	return NewParent, nil
}

// SetHandler는 지정된 경로와 메서드에 대한 핸들러를 트리에 등록합니다.
// PathWithSegment를 사용하여 경로를 효율적으로 처리하고 필요한 노드들을 생성하거나 기존 노드를 분할합니다.
//
// 동작 방식:
// 1. PathWithSegment로 경로를 세그먼트 단위로 분석
// 2. 트리를 순회하며 매칭되는 노드 탐색
// 3. 완전 매칭: 다음 세그먼트로 이동
// 4. 부분 매칭: 노드 분할 후 계속 진행
// 5. 매칭 실패: 새로운 자식 노드 생성
func (Instance *Tree) SetHandler(Method MethodType, RawPath string, Handler HandlerFunc) error {
	if RawPath == "" {
		return Error.NewErrorWithCode(Error.InvalidParameter, RawPath)
	}
	if Method == NotAllowed {
		return Error.NewErrorWithCode(Error.MethodNotAllowed, RawPath)
	}
	if Method != CONNECT && Handler == nil {
		return Error.NewErrorWithCode(Error.InvalidParameter, RawPath)
	}
	// PathWithSegment 구조체를 사용한 효율적인 경로 처리
	Path := NewPathWithSegment(RawPath)
	Path.Next()
	if Path.IsSame() {
		Instance.RootNode.Handlers[Method] = Handler
		return nil
	}
	Parent := Instance.RootNode
	SetHelper:
	for {
		if Path.IsSame() {
			Parent.Handlers[Method] = Handler
			break
		}
		Index := sort.Search(len(Parent.Indices), func (Index int) bool {
			return Parent.Indices[Index] >= Path.Body[Path.Start] 
		})
		if Index < len(Parent.Indices) && Parent.Indices[Index] == Path.Body[Path.Start] {
			Matched , MatchingPoint , Left := Instance.Match(*Path,Parent.Children[Index].Path)
			if Matched {
				Path.Next()
				Parent = Parent.Children[Index]
				continue SetHelper
			}
			if MatchingPoint < len(Parent.Children[Index].Path) {
				NewParent, Err := Instance.SplitNode(Parent, Parent.Children[Index], MatchingPoint)
				if Err != nil {
					return Err
				}
				Parent = NewParent
				if Left.GetLength() > 0 {
					Path = &Left
					continue SetHelper
				}
				Path.Next()
				continue SetHelper
			}
			Parent = Parent.Children[Index]
			Path = &Left
			continue SetHelper
		}
		Child, Err := Instance.InsertChild(Parent, Path.Get())
		if Err != nil {
			return Err
		}
		Parent = Child
		Path.Next()
		continue SetHelper
	}
	return nil
}

// GetHandler는 HTTP 요청에 대응하는 핸들러를 트리에서 찾아 반환합니다.
// PathWithSegment를 사용하여 경로를 효율적으로 처리하고 매개변수를 추출하여 핸들러를 선택합니다.
//
// 매칭 우선순위:
// 1. Static 노드: 정확한 문자열 매칭 (가장 높은 우선순위)
// 2. WildCard 노드: 단일 세그먼트 매개변수 (:param)
// 3. CatchAll 노드: 나머지 모든 경로 (*path, 가장 낮은 우선순위)
//
// 반환값: (핸들러함수, 매개변수객체) - 매개변수가 없으면 nil 반환
func (Instance *Tree) GetHandler(Request *http.Request, GetParams func() *Param.Params) (HandlerFunc, *Param.Params) {
	RawPath := Request.URL.Path
	Method := Instance.StringToMethodType(Request.Method)
	var Params *Param.Params
	// PathWithSegment를 사용한 메모리 효율적인 경로 처리
	Path := NewPathWithSegment(RawPath)
	if Method == NotAllowed {
		return Instance.NotAllowedHandler, nil
	}
	Path.Next()
	Node := Instance.RootNode
GetHelper:
	for {
		if Path.GetLength() == 0 {
			return Instance.SelectHandler(Node, Method), Params
		}

		for Index := 0; Index < len(Node.Indices); Index++ {
			if Path.Body[Path.Start] == Node.Indices[Index] {
				Matched, MatchingPoint, Left := Instance.Match(*Path, Node.Children[Index].Path)
				switch {
				case Matched:
					Path.Next()
					Node = Node.Children[Index]
					continue GetHelper
				case MatchingPoint > 0 && MatchingPoint == len(Node.Children[Index].Path) && Path.GetLength() > 0:
					Node = Node.Children[Index]
					Path = &Left
					continue GetHelper
				}
			}
		}
		// 2순위: WildCard 노드 매칭 (단일 세그먼트 캡처)
		if Node.WildCard != nil {
			if Params == nil {
				Params = GetParams()
			}
			// 현재 세그먼트를 매개변수로 저장하고 다음 세그먼트로 진행
			Params.Add(Node.WildCard.Param, Path.Get())
			Path.Next()
			Node = Node.WildCard
			continue GetHelper
		}
		// 3순위: CatchAll 노드 매칭 (나머지 모든 경로 캡처)
		if Node.CatchAll != nil {
			if Params == nil {
				Params = GetParams()
			}
			// PathWithSegment의 나머지 경로를 하나의 매개변수로 저장
			Params.Add(Node.CatchAll.Param, Path.GetToEnd())
			Path.Start = Path.End
			// 빈 경로로 CatchAll 노드에서 핸들러 검색 (경로 소비 완료)
			Node = Node.CatchAll
			continue GetHelper
		}
		// 매칭되는 라우트가 없는 경우 매개변수 객체 반환 후 404 핸들러 반환
		return Instance.NotFoundHandler, Params
	}
}

// SetMiddleware는 트리에 미들웨어를 추가합니다.
// 추가된 미들웨어는 모든 핸들러에 적용됩니다.
func (Instance *Tree) SetMiddleware(Middleware Middleware.Middleware) {
	Instance.Middlewares = append(Instance.Middlewares, Middleware)
}

func (Instance *Tree) CompileMiddlewares() {
	if len(Instance.Middlewares) == 0 {
		Instance.CompiledMiddleware = nil
		return
	}

	Instance.CompiledMiddleware = func(BaseHandler HandlerFunc) HandlerFunc {
		Temp := BaseHandler
		for Index := len(Instance.Middlewares) -1; Index> -1 ; Index-- {
			Temp = Instance.Middlewares[Index].GetHandler()(Temp)
		}
		return BaseHandler
	}
}


// ApplyMiddleware는 핸들러에 등록된 미들웨어들을 역순으로 적용합니다.
// 마지막에 등록된 미들웨어가 가장 바깥쪽에 위치하게 됩니다.
func (Instance *Tree) ApplyMiddleware(Handler HandlerFunc) HandlerFunc {
	if Instance.CompiledMiddleware == nil {
		return  Handler
	}
	return Instance.CompiledMiddleware(Handler)
}

// ServeHTTP는 http.Handler 인터페이스를 구현합니다.
// 요청에 대한 핸들러를 찾아 미들웨어를 적용한 후 실행합니다.
func (Instance *Tree) ServeHTTP(Writer http.ResponseWriter, Request *http.Request) {
	Handler, Params := Instance.GetHandler(Request, Instance.Pool.Get)
	if Instance.CompiledMiddleware != nil {
		Handler = Instance.CompiledMiddleware(Handler)
	}
	Handler(Writer, Request, Params)

	// 매개변수 객체를 풀에 반환
	if Params != nil {
		Instance.Pool.Put(Params)
	}
}
