package Tree

import (
	"LiteFrame/Router/Middleware"
	"LiteFrame/Router/Param"
	"context"
	"fmt"
	"net/http"
	"strings"
)

var MethodList = map[string]MethodType{
	"GET":     GET,
	"HEAD":    HEAD,
	"OPTIONS": OPTIONS,
	"TRACE":   TRACE,
	"POST":    POST,
	"PUT":     PUT,
	"DELETE":  DELETE,
	"CONNECT": CONNECT,
	"PATCH":   PATCH,
}

type Tree struct {
	RootNode          *Node
	NotFoundHandler   http.HandlerFunc
	NotAllowedHandler http.HandlerFunc
	Middlewares       []Middleware.Middleware
}

func NewTree() Tree {
	return Tree{
		RootNode: NewNode(RootType, "/"),
	}
}

func (Instance *Tree) IsWildCard(Input string) bool {
	return len(Input) > 0 && Input[0] == WildCardPrefix
}

func (Instance *Tree) IsCatchAll(Input string) bool {
	return len(Input) > 0 && Input[0] == CatchAllPrefix
}

func (Instance *Tree) SplitPath(Path string) []string {
	Segments := make([]string, 0, strings.Count(Path, "/")+1)
	for Segment := range strings.SplitSeq(Path, "/") {
		if Segment != "" {
			Segments = append(Segments, Segment)
		}
	}
	return Segments
}

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

func (Instance *Tree) InsertHandler(Node *Node, Method string, Handler http.HandlerFunc) error {
	if TypeMethod, Ok := MethodList[Method]; Ok {
		Node.Handlers[TypeMethod] = Handler
		return nil
	}
	return fmt.Errorf("error : Method %s is not Allowed", Method)
}

func (Instance *Tree)SelectHandler(Node *Node , Method string) http.HandlerFunc{
	if TypeMethod  , Ok := MethodList[Method]; Ok {
		if Handler , Ok2 := Node.Handlers[TypeMethod] ; Ok2{
			return Handler
		}
	}
	return Instance.NotAllowedHandler
} 


func (Instance *Tree) InsertChild(Parent *Node, Path string) (*Node, error) {
	if Instance.IsWildCard(Path) {
		if Parent.WildCard {
			return nil, NewTreeError(Path, "Can not Have Two WildCard Node")
		}
		Parent.WildCard = true
		Child := NewNode(WildCardType, Path)
		Child.Param = Path[1:]
		Child.WildCard = true
		Parent.Indices = Parent.Indices + string(Path[0])
		Parent.Children = append(Parent.Children,Child) 
		return Child, nil
	}
	if Instance.IsCatchAll(Path) {
		if Parent.CatchAll {
			return nil, NewTreeError(Path, "Can not Have Two CatchAll Node")
		}
		Parent.CatchAll = true
		Child := NewNode(CatchAllType, Path)
		Child.CatchAll = true
		Parent.Indices = Parent.Indices + string(Path[0])
		Parent.Children= append(Parent.Children, Child)
		return Child, nil
	}
	Child := NewNode(StaticType, Path)
	Parent.Indices = Parent.Indices + string(Path[0])
	Parent.Children= append(Parent.Children, Child) 
	return Child, nil
}

func (Instance *Tree) SplitNode(Parent *Node, Child *Node, SplitPoint int) (*Node, error) {
	if SplitPoint < 0 || SplitPoint > len(Child.Path) {
		return nil, fmt.Errorf("split point %d is out of range for path %q (length %d)", SplitPoint, Child.Path, len(Child.Path))
	}
	Left := Child.Path[:SplitPoint]
	Right := Child.Path[SplitPoint:]
	//Child.Path = Right
	NewParent := NewNode(StaticType,Left) // Instance.InsertChild(Parent, Left)
	for Index , Indice := range []byte(Parent.Indices) {
		if Indice == Left[0] {
			NewParent.Children[Index] = NewParent
		}
	}
	NewChild  , Err := Instance.InsertChild(NewParent,Right)
	if Err != nil {
		return nil , Err
	}
	NewChild.Indices = Child.Indices
	NewChild.Children = Child.Children
	NewChild.Handlers = Child.Handlers
	
	return NewParent , nil
}

func (Instance *Tree) SetHandler(Method string, Path string, Handler http.HandlerFunc) error {
	if Method == "" || Path == "" || (Handler == nil && Method != CONNECT) {
		return NewTreeError(Path,"Invalid Parameters, Path and Handler are Required")
	}
	Paths := Instance.SplitPath(Path)
	if len(Paths) == 0 {
		return Instance.InsertHandler(Instance.RootNode, Method, Handler)
	}
	return Instance.SetHelper(Instance.RootNode, Paths, Method, Handler)
}

func (Instance *Tree) SetHelper(Parent *Node, Paths []string, Method string, Handler http.HandlerFunc) error {
	if len(Paths) == 0 {
		return Instance.InsertHandler(Parent, Method, Handler)
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

func (Instance *Tree) TryMatch(Parent *Node, Paths []string, Method string, Handler http.HandlerFunc) (bool, error) {
	if len(Parent.Children) == 0 {
		return false, nil
	}
	for _, Child := range Parent.Children {
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
					NewPaths := make([]string, len(Paths))
					copy(NewPaths, Paths)
					NewPaths[0] = LeftPath
					return true, Instance.SetHelper(NewParent, NewPaths, Method, Handler)
				}
				return true, Instance.SetHelper(NewParent, Paths[1:], Method, Handler)
			case MatchingPoint > 0 && MatchingPoint == len(Child.Path) && len(LeftPath) > 0: 
				NewPaths := make([]string, len(Paths))
				copy(NewPaths, Paths)
				NewPaths[0] = LeftPath
				return true, Instance.SetHelper(Child, NewPaths, Method, Handler)	
			
		}
	}
	return false, nil
}

func (Instance *Tree) GetHandler(Requeset *http.Request) http.HandlerFunc {
	Path := Requeset.URL.Path
	Method := Requeset.Method
	Params := Param.NewParams()
	Paths := Instance.SplitPath(Path)
	if len(Paths) == 0 {
		return Instance.SelectHandler(Instance.RootNode,Method)
	}
	return Instance.GetHelper(Instance.RootNode,Method,Paths,Params)
}

func (Instance *Tree) GetHelper(Node *Node,Method string,Paths []string,Params Param.Params) http.HandlerFunc {
	if len(Paths) == 0 {
		return  func(Writer http.ResponseWriter, Request *http.Request) { 
			Ctx := Request.Context()
			Ctx = context.WithValue(Ctx,Param.Key{},Params)
			NewRquest  :=Request.WithContext(Ctx)
			Instance.SelectHandler(Node,Method)(Writer,NewRquest)
		} 
	}
	for _ , Child := range Node.Children {
		Matched, MatchingPoint, LeftPath := Instance.Match(Paths[0], Child.Path)
		switch {
			case Matched: 
				return Instance.GetHelper(Child,Method,Paths[1:],Params)
			case MatchingPoint > 0 && MatchingPoint < len(Child.Path):
				NewPaths := make([]string,len(Paths))
				copy(NewPaths,Paths)
				NewPaths[0] = LeftPath
				return Instance.GetHelper(Child,Method,NewPaths,Params)
			case Child.WildCard:
				Params.Add(Child.Param,Paths[0])
				return 	Instance.GetHelper(Child,Method,Paths[1:],Params)
			case Child.CatchAll:
				Params.Add("CactchAll",strings.Join(Paths,"/"))
				return Instance.GetHelper(Child,Method,[]string{},Params)
		}
	}
	return Instance.NotFoundHandler
}



func (Instance *Tree) ServeHTTP(Writer http.ResponseWriter, Request *http.Request) {

}
