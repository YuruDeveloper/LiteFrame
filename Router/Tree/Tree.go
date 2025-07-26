// Package Tree는 HTTP 라우팅을 위한 Radix Tree 구현을 제공합니다.
// 효율적인 URL 패턴 매칭과 핸들러 관리를 위한 트리 구조를 구현합니다.
package Tree

import (
	"LiteFrame/Router/Middleware"
	"LiteFrame/Router/Param"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// Tree는 HTTP 라우팅을 위한 Radix Tree 구조체입니다.
// 루트 노드, 매개변수 풀, 핸들러 및 미들웨어를 관리합니다.
type Tree struct {
	RootNode          *Node                    // 트리의 루트 노드
	Pool *Param.ParamsPool                    // 매개변수 재사용을 위한 풀
	NotFoundHandler   http.HandlerFunc        // 404 핸들러
	NotAllowedHandler http.HandlerFunc        // 405 핸들러  
	Middlewares       []Middleware.Middleware // 미들웨어 목록
}

// NewTree는 새로운 Tree 인스턴스를 생성합니다.
// 루트 노드와 매개변수 풀을 초기화합니다.
func NewTree() Tree {
	return Tree{
		RootNode: NewNode(RootType, "/"),
		Pool: Param.NewParamsPool(),
	}
}

// IsWildCard는 입력 문자열이 와일드카드 패턴(:param)인지 확인합니다.
func (Instance *Tree) IsWildCard(Input string) bool {
	return len(Input) > 0 && Input[0] == WildCardPrefix
}

// IsCatchAll는 입력 문자열이 캐치올 패턴(*path)인지 확인합니다.
func (Instance *Tree) IsCatchAll(Input string) bool {
	return len(Input) > 0 && Input[0] == CatchAllPrefix
}

// SplitPath는 URL 경로를 '/' 기준으로 분할하여 문자열 배열로 반환합니다.
// 빈 문자열과 연속된 '/'는 제거됩니다.
// 성능 최적화: 미리 용량을 할당하여 메모리 재할당을 최소화합니다.
func (Instance *Tree) SplitPath(Path string) []string {
	if Path == "" || Path == "/" {
		return []string{}
	}
	SlashCount := 0
	for i := 0 ; i < len(Path) ; i ++ {
		if Path[i] == '/' {
			SlashCount++
		}
	}
	if len(Path) > 0 && Path[0] != '/' {
		SlashCount++
	}
	Result := make([]string,0,SlashCount)
	Index := 0
	Slice := ""
	for Temp := 0; Temp < len(Path); Temp++ {
		if Path[Temp] == '/' {
			// 다중 할당을 통한 효율적인 슬라이스 추출 및 조건 검사
			if Slice =  Path[Index:Temp];Index != Temp && Slice != ""  {
				Result = append(Result, Slice)
				Slice = ""
			}
			// 다음 세그먼트 시작 인덱스 설정 ('/' 다음 문자)
			Index = Temp + 1
		}
	}
	// 마지막 세그먼트 처리 (경로 끝에 '/'가 없는 경우)
	if Index < len(Path) {
		Result = append(Result, Path[Index:])
	}
	return Result
}

// StringToMethodType는 HTTP 메서드 문자열을 MethodType으로 변환합니다.
// 지원되지 않는 메서드의 경우 NotAllowed를 반환합니다.
func (Instance  Tree) StringToMethodType(Method string) MethodType {
	switch (Method) {
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



// Match는 두 문자열의 공통 접두사를 찾아 매칭 결과를 반환합니다.
// 반환값: (완전매칭여부, 매칭된인덱스, 남은문자열)
// Radix Tree의 핵심 알고리즘: 노드 분할과 경로 매칭에 사용됩니다.
func (Instance *Tree) Match(SourcePath string, Two string) (bool, int, string) {
	// 두 문자열 중 짧은 길이를 기준으로 비교 범위 설정
	Length := len(SourcePath)
	if Length > len(Two) {
		Length = len(Two)
	}
	// 바이트 단위로 순차 비교하여 공통 접두사 길이 계산
	var Index int
	for Index = 0; Index < Length; Index++ {
		if SourcePath[Index] != Two[Index] {
			break
		}
	}
	// One이 Two의 완전한 접두사인지 확인
	Matched := Index == len(SourcePath)
	// One에서 매칭되지 않은 나머지 부분 추출
	Remain := ""
	if Index < len(SourcePath) {
		Remain = SourcePath[Index:]
	}
	return Matched, Index, Remain
}

// InsertHandler는 지정된 노드에 HTTP 메서드별 핸들러를 등록합니다.
func (Instance *Tree) InsertHandler(Node *Node, Method MethodType, Handler HandlerFunc) {
		Node.Handlers[Method] = Handler
}

// SelectHandler는 노드에서 메서드에 맞는 핸들러를 선택하고 매개변수를 컨텍스트에 주입합니다.
// 핸들러가 없으면 NotAllowedHandler를 반환합니다.
// 중요: 메모리 풀 관리와 컨텍스트 주입을 동시에 처리하는 핵심 함수입니다.
func (Instance *Tree) SelectHandler(Node *Node, Method MethodType,Params *Param.Params) http.HandlerFunc {
	if Handler := Node.Handlers[Method]; Handler != nil {
		// 클로저를 통해 매개변수를 컨텍스트에 주입하고 메모리 풀 반환 보장
		return func(Writer http.ResponseWriter, Request *http.Request) {
			Handler(Writer,Request,Params)
			// 요청 처리 완료 후 매개변수 객체를 풀에 반환 (메모리 최적화)
			Instance.Pool.Put(Params)
		} 
	}
	
	// 핸들러가 없는 경우에도 매개변수 객체는 반드시 풀에 반환
	Instance.Pool.Put(Params)
	return Instance.NotAllowedHandler
}

// InsertUniqueTypeChild는 고유한 타입의 자식 노드(WildCard/CatchAll)를 삽입합니다.
// 중복된 매개변수명이 있으면 에러를 반환합니다.
// 함수형 프로그래밍 패턴: 고차 함수를 활용한 중복 코드 제거
func (Instance *Tree) InsertUniqueTypeChild(Parent *Node,Path string , Target *Node,Type NodeType , ErrorFn func(string) error , SetFn func (*Node,*Node)) (*Node, error){
	switch {
		// 빈 매개변수명 검증 (":" 또는 "*" 만 있는 경우)
		case Path[1:] == "":
			return nil , NewTreeError("Nil Parameter is Not Allowed",Path)
		// 같은 타입이지만 다른 매개변수명인 경우 충돌 오류
		case Target != nil && Target.Param != Path[1:]:
				return nil , ErrorFn(Path)
		// 동일한 매개변수명의 노드가 이미 존재하는 경우 재사용
		case Target != nil && Target.Param == Path[1:]:
				return Target , nil
		// 새로운 매개변수 노드 생성
		default:
			Child := NewNode(Type,Path)
			Child.Param = Path[1:] // 접두사(: 또는 *) 제거하여 매개변수명만 저장
			SetFn(Parent,Child) // 함수 포인터를 통한 타입별 노드 설정
			return Child , nil
	}
}

// InsertChild는 부모 노드에 자식 노드를 삽입합니다.
// 경로 타입에 따라 Static, WildCard, CatchAll 노드를 생성합니다.
func (Instance *Tree) InsertChild(Parent *Node, Path string) (*Node, error) {
	switch {
		case Instance.IsWildCard(Path):
			return Instance.InsertUniqueTypeChild(Parent,Path,Parent.WildCard, WildCardType,
				func (Path string) error{ return NewTreeError(Path, "Can not Have Two WildCard Node")},
				func(Parent, Child *Node) {Parent.WildCard = Child})
		case Instance.IsCatchAll(Path):
			return Instance.InsertUniqueTypeChild(Parent,Path,Parent.CatchAll, CatchAllType,
				func (Path string) error{ return NewTreeError(Path, "Can not Have Two CatchAll Node")},
				func(Parent, Child *Node) {Parent.CatchAll = Child})
		default:
			Child := NewNode(StaticType, Path)
			InsertLocation := sort.Search(len(Parent.Indices),func(Index int) bool { return Parent.Indices[Index] >= Path[0]})
			Parent.Indices = append(Parent.Indices[:InsertLocation],append([]byte{Path[0]},Parent.Indices[InsertLocation:]...)...) 
			Parent.Children = append(Parent.Children[:InsertLocation], append([]*Node{Child},Parent.Children[InsertLocation:]...)...)
			return Child, nil
	}
}

// SplitNode는 기존 노드를 분할점에서 두 개의 노드로 분리합니다.
// 공통 접두사를 가진 새로운 부모 노드를 생성하고 기존 노드를 자식으로 만듭니다.
func (Instance *Tree) SplitNode(Parent *Node, Child *Node, SplitPoint int) (*Node, error) {
	if SplitPoint < 0 || SplitPoint > len(Child.Path) {
		return nil, fmt.Errorf("split point %d is out of range for path %q (length %d)", SplitPoint, Child.Path, len(Child.Path))
	}
	Left := Child.Path[:SplitPoint]
	Right := Child.Path[SplitPoint:]
	if len(Left) == 0 {
		return nil , NewTreeError("Splite Failed",Child.Path)
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
// 경로를 분할하고 필요한 노드들을 생성하거나 기존 노드를 분할합니다.
func (Instance *Tree) SetHandler(Method string, Path string, Handler HandlerFunc) error {
	if Method == "" || Path == "" {
		return NewTreeError(Path, "Invalid Parameters, Path or Method are Required")
	}
	MethodType := Instance.StringToMethodType(Method)
	if MethodType == NotAllowed {
		return fmt.Errorf("error : Method %s is not Allowed", Method)
	}
	if MethodType != CONNECT && Handler == nil {
		return NewTreeError(Path, "Invalid Parameter, Handler are Required")
	}
	Paths := Instance.SplitPath(Path)
	if len(Paths) == 0 {
		Instance.InsertHandler(Instance.RootNode, MethodType, Handler)
		return nil
	}
	return Instance.SetHelper(Instance.RootNode, Paths, MethodType, Handler)
}

// SetHelper는 SetHandler의 재귀 헬퍼 함수입니다.
// 경로 배열을 순회하며 트리를 구성하고 핸들러를 등록합니다.
func (Instance *Tree) SetHelper(Parent *Node, Paths []string, Method MethodType, Handler HandlerFunc) error {
	if len(Paths) == 0 {
		Instance.InsertHandler(Parent, Method, Handler)
		return nil
	}
	if Ok, Err := Instance.TryMatch(Parent, Paths, Method, Handler); Ok {
		return Err
	}
	Child, Err := Instance.InsertChild(Parent, Paths[0])
	if Err != nil {
		return Err
	}
	return Instance.SetHelper(Child, Paths[1:], Method, Handler)
}

// TryMatch는 부모 노드의 자식들 중에서 현재 경로와 매칭되는 노드를 찾습니다.
// 매칭되는 자식이 있으면 해당 자식으로 라우팅을 계속합니다.
func (Instance *Tree) TryMatch(Parent *Node, Paths []string, Method MethodType, Handler HandlerFunc) (bool, error) {
	if len(Parent.Children) == 0 {
		return false, nil
	}
	for _, Child := range Parent.Children {
		if Ok , Err := Instance.MatchChild(Parent,Child,Paths,Method,Handler) ; Ok{
			return Ok , Err
		}
	}
	return false, nil
}

// MatchChild는 특정 자식 노드와 현재 경로의 매칭을 시도합니다.
// 완전 매칭, 부분 매칭에 따라 노드 분할이나 라우팅을 수행합니다.
func (Instance *Tree) MatchChild(Parent *Node,Child *Node,Paths []string,Method MethodType,Handler HandlerFunc) (bool , error){
		Matched, MatchingPoint, LeftPath := Instance.Match(Paths[0], Child.Path)
		switch {
		case Matched:
			return true, Instance.SetHelper(Child, Paths[1:], Method, Handler)
		case MatchingPoint > 0 && MatchingPoint < len(Child.Path):
			NewParent, Err := Instance.SplitNode(Parent, Child, MatchingPoint)
			if Err != nil {
				return true, Err
			}
			if len(LeftPath) > 0 {
				Paths[0] = LeftPath
				return true, Instance.SetHelper(NewParent, Paths, Method, Handler)
			}
			return true, Instance.SetHelper(NewParent, Paths[1:], Method, Handler)
		case MatchingPoint > 0 && MatchingPoint == len(Child.Path) && len(LeftPath) > 0:
			Paths[0] = LeftPath
			return true, Instance.SetHelper(Child, Paths, Method, Handler)
		}
		return false , nil
}

// GetHandler는 HTTP 요청에 대응하는 핸들러를 트리에서 찾아 반환합니다.
// 경로를 분할하고 매개변수를 추출하여 적절한 핸들러를 선택합니다.
func (Instance *Tree) GetHandler(Request *http.Request) http.HandlerFunc {
	Path := Request.URL.Path
	Method := Instance.StringToMethodType(Request.Method)
	Params := Instance.Pool.Get()
	Paths := Instance.SplitPath(Path)
	if Method == NotAllowed {
		return Instance.NotAllowedHandler
	}
	if len(Paths) == 0 {
		return Instance.SelectHandler(Instance.RootNode, Method,Params)
	}
	return Instance.GetHelper(Instance.RootNode, Method, Paths, Params)
}

// GetHelper는 GetHandler의 재귀 헬퍼 함수입니다.
// 경로를 순회하며 Static, WildCard, CatchAll 노드 순서로 매칭을 시도합니다.
// 라우팅 우선순위: Static > WildCard > CatchAll (성능과 정확성의 균형)
func (Instance *Tree) GetHelper(Node *Node, Method MethodType, Paths []string, Params *Param.Params) http.HandlerFunc {
	// 경로 세그먼트가 모두 소비된 경우 현재 노드에서 핸들러 검색
	if len(Paths) == 0 {
		return Instance.SelectHandler(Node,Method,Params)
	} 
	// 1순위: Static 자식 노드들에서 첫 번째 바이트 기반 빠른 매칭
	for Index , Indice := range []byte(Node.Indices) {
		switch (Indice) {
			case Paths[0][0]:
				// 정확한 문자열 매칭 수행
				Matched, MatchingPoint, LeftPath := Instance.Match(Paths[0], Node.Children[Index].Path)
				switch {
					// 완전 매칭: 다음 세그먼트로 진행
					case Matched:
						return Instance.GetHelper(Node.Children[Index], Method, Paths[1:], Params)
					// 현재 세그먼트가 노드 경로로 시작하고 남은 부분이 있을 때
					case MatchingPoint > 0 && MatchingPoint == len(Node.Children[Index].Path) && len(LeftPath) > 0:
						Paths[0] = LeftPath
						return Instance.GetHelper(Node.Children[Index], Method, Paths, Params)
				}
		}
	}

	// 2순위: WildCard 노드 매칭 (단일 세그먼트 캡처)
	if Node.WildCard != nil {
		// 현재 세그먼트를 매개변수로 저장하고 다음 세그먼트로 진행
		Params.Add(Node.WildCard.Param, Paths[0])
		return Instance.GetHelper(Node.WildCard, Method, Paths[1:], Params)
	}
	// 3순위: CatchAll 노드 매칭 (나머지 모든 경로 캡처)
	if Node.CatchAll != nil {
		// 남은 모든 세그먼트를 '/'로 연결하여 하나의 매개변수로 저장
		Params.Add(Node.CatchAll.Param, strings.Join(Paths, "/"))
		// 빈 경로로 CatchAll 노드에서 핸들러 검색 (경로 소비 완료)
		return Instance.GetHelper(Node.CatchAll, Method, []string{}, Params)
	}

	// 매칭되는 라우트가 없는 경우 매개변수 객체 반환 후 404 핸들러 반환
	Instance.Pool.Put(Params)
	return Instance.NotFoundHandler
}

// SetMiddleware는 트리에 미들웨어를 추가합니다.
// 추가된 미들웨어는 모든 핸들러에 적용됩니다.
func (Instance *Tree) SetMiddleware(Middleware Middleware.Middleware) {
	Instance.Middlewares = append(Instance.Middlewares, Middleware)
}

// ApplyMiddleware는 핸들러에 등록된 미들웨어들을 역순으로 적용합니다.
// 마지막에 등록된 미들웨어가 가장 바깥쪽에 위치하게 됩니다.
func (Instance *Tree) ApplyMiddleware(Handler http.HandlerFunc) http.Handler{
	var Temp http.Handler = Handler
	for Index := len(Instance.Middlewares) -1; Index >= 0; Index-- {
		Temp = Instance.Middlewares[Index].GetHandler()(Temp)
	}
	return Temp
}

// ServeHTTP는 http.Handler 인터페이스를 구현합니다.
// 요청에 대한 핸들러를 찾아 미들웨어를 적용한 후 실행합니다.
func (Instance *Tree) ServeHTTP(Writer http.ResponseWriter, Request *http.Request) {
	Instance.ApplyMiddleware(Instance.GetHandler(Request)).ServeHTTP(Writer,Request)
}
