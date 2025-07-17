package Tree

import "net/http"

func NewNode(Type NodeType,Path string) Node {
	return Node{
		Type: Type,
		Path: Path,
		Children: make(map[string]*Node),
		Handlers: make(map[MethodType]http.HandlerFunc),
		WlidCard: false,
		CactchAll: false,
	}
}


type Node struct {
	Type NodeType
	Indices string
	Path string
	Children map[string]*Node
	Handlers map[MethodType]http.HandlerFunc  
	WlidCard bool
	CactchAll bool
	Parm string
}
