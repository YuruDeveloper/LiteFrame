package Tree

import "net/http"

func NewNode(Type NodeType, Path string) Node {
	return Node{
		Type:     Type,
		Path:     Path,
		Children: make(map[string]*Node),
		Handlers: make(map[MethodType]http.HandlerFunc),
		WildCard: false,
		CatchAll: false,
	}
}

type Node struct {
	Type     NodeType
	Path     string
	Children map[string]*Node
	Handlers map[MethodType]http.HandlerFunc
	WildCard bool
	CatchAll bool
	Param    string
}
