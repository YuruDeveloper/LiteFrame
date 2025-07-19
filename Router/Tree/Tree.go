// Package Tree는 HTTP 라우팅을 위한 Radix Tree 구현을 제공합니다.
// 효율적인 URL 패턴 매칭과 핸들러 관리를 위한 트리 구조를 구현합니다.
package Tree

import (
	"LiteFrame/Router/Middleware"
	"LiteFrame/Router/Param"
	"context"
	"fmt"
	"net/http"
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
func (Instance *Tree) SplitPath(Path string) []string {
	Result := make([]string,0,len(Path) / 2 + 1)
	Index := 0
	Slice := ""
	for Temp := 0; Temp < len(Path); Temp++ {
		if Path[Temp] == '/' {
			if Slice =  Path[Index:Temp];Index != Temp && Slice != ""  {
				Result = append(Result, Slice)
				Slice = ""
			}
			Index = Temp + 1
		}
	}
	if Index < len(Path) {
		Result = append(Result, Path[Index:])
	}
	return Result
}

// StringToMethodType는 HTTP 메서드 문자열을 MethodType으로 변환합니다.
// 지원되지 않는 메서드의 경우 NotAllowed를 반환합니다.
func (Instance *Tree) StringToMethodType(Method string) MethodType {
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
func (Instance *Tree) Match(One string, Two string) (bool, int, string) {
	Length := len(One)
	if Length > len(Two) {
		Length = len(Two)
	}
	var Index int
	for Index = 0; Index < Length; Index++ {
		if One[Index] != Two[Index] {
			break
		}
	}
	Matched := Index == len(One)
	Remain := ""
	if Index < len(One) {
		Remain = One[Index:]
	}
	return Matched, Index, Remain
}

// InsertHandler는 지정된 노드에 HTTP 메서드별 핸들러를 등록합니다.
func (Instance *Tree) InsertHandler(Node *Node, Method MethodType, Handler http.HandlerFunc) {
		Node.Handlers[Method] = Handler
}

// SelectHandler는 노드에서 메서드에 맞는 핸들러를 선택하고 매개변수를 컨텍스트에 주입합니다.
// 핸들러가 없으면 NotAllowedHandler를 반환합니다.
func (Instance *Tree) SelectHandler(Node *Node, Method MethodType,Params *Param.Params) http.HandlerFunc {
	if Handler := Node.Handlers[Method]; Handler != nil {
		return func(Writer http.ResponseWriter, Request *http.Request) {
			Ctx := Request.Context()
			Ctx = context.WithValue(Ctx, Param.Key{}, Params)
			NewRequest := Request.WithContext(Ctx)
			Handler(Writer,NewRequest)
			Instance.Pool.Put(Params)
		} 
	}
	
	Instance.Pool.Put(Params)
	return Instance.NotAllowedHandler
}

// InsertUniqueTypeChild는 고유한 타입의 자식 노드(WildCard/CatchAll)를 삽입합니다.
// 중복된 매개변수명이 있으면 에러를 반환합니다.
func (Instance *Tree) InsertUniqueTypeChild(Parent *Node,Path string , Target *Node,Type NodeType , ErrorFn func(string) error , SetFn func (*Node,*Node)) (*Node, error){
	switch {
		case Path[1:] == "":
			return nil , NewTreeError("Nil Parameter is Not Allowed",Path)
		case Target != nil && Target.Param != Path[1:]:
				return nil , ErrorFn(Path)
		case Target != nil && Target.Param == Path[1:]:
				return Target , nil
		default:
			Child := NewNode(Type,Path)
			Child.Param = Path[1:]
			SetFn(Parent,Child)
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
			Parent.Indices = append(Parent.Indices, Path[0]) 
			Parent.Children = append(Parent.Children, Child)
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
	NewParent := NewNode(StaticType, Left)
	for Index, Target := range Parent.Children {
		if Target == Child {
			Parent.Children[Index] = NewParent
		}
	}
	if len(Left) > 0 {
		NewParent.Indices = make([]byte,Left[0])
	}
	Child.Path = Right
	NewParent.Children = []*Node{Child}
	return NewParent, nil
}

// SetHandler는 지정된 경로와 메서드에 대한 핸들러를 트리에 등록합니다.
// 경로를 분할하고 필요한 노드들을 생성하거나 기존 노드를 분할합니다.
func (Instance *Tree) SetHandler(Method string, Path string, Handler http.HandlerFunc) error {
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
func (Instance *Tree) SetHelper(Parent *Node, Paths []string, Method MethodType, Handler http.HandlerFunc) error {
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
func (Instance *Tree) TryMatch(Parent *Node, Paths []string, Method MethodType, Handler http.HandlerFunc) (bool, error) {
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
func (Instance *Tree) MatchChild(Parent *Node,Child *Node,Paths []string,Method MethodType,Handler http.HandlerFunc) (bool , error){
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
			NewPaths := make([]string, len(Paths))
			copy(NewPaths, Paths)
			NewPaths[0] = LeftPath
			return true, Instance.SetHelper(Child, NewPaths, Method, Handler)
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
func (Instance *Tree) GetHelper(Node *Node, Method MethodType, Paths []string, Params *Param.Params) http.HandlerFunc {
	if len(Paths) == 0 {
		return Instance.SelectHandler(Node,Method,Params)
	} 
	for Index , Indice := range []byte(Node.Indices) {
		switch (Indice) {
			case Paths[0][0]:
				Matched, MatchingPoint, LeftPath := Instance.Match(Paths[0], Node.Children[Index].Path)
				switch {
					case Matched:
						return Instance.GetHelper(Node.Children[Index], Method, Paths[1:], Params)
					case  MatchingPoint > 0 && MatchingPoint < len(Node.Children[Index].Path):
						Paths[0] = LeftPath
						return Instance.GetHelper(Node.Children[Index], Method, Paths, Params)
				}
		}
	}

	if Node.WildCard != nil {
		Params.Add(Node.WildCard.Param, Paths[0])
		return Instance.GetHelper(Node.WildCard, Method, Paths[1:], Params)
	}
	if Node.CatchAll != nil {
		Params.Add(Node.CatchAll.Param, strings.Join(Paths, "/"))
		return Instance.GetHelper(Node.CatchAll, Method, []string{}, Params)
	}

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
