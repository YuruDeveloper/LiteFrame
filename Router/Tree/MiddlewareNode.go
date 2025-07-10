package Tree

import (
	"LiteFrame/Router/Middleware"
	"net/http"
	Componet "LiteFrame/Router/Tree/Componet"
)

type MiddlewareNode struct {
	Identity *Componet.Identity
	Container *Componet.Container[Componet.Node]
	MiddlewareHandler *Componet.MiddlewareHandler
}

func NewMiddlewareNode() MiddlewareNode {
	return MiddlewareNode{
		Identity: Componet.NewIdentity(Componet.High, Componet.MiddlewareType, false),
		Container: Componet.NewContainer(Componet.NewError(Componet.MiddlewareType,"","")),
		MiddlewareHandler: Componet.NewMiddlewareHandler(Componet.NewError(Componet.MiddlewareType,"","")),
	}
}

func (Instance *MiddlewareNode) GetPriority() Componet.PriorityLevel {
	return Instance.Identity.GetPriority()
}

func (Instance *MiddlewareNode) GetType() Componet.NodeType {
	return Instance.Identity.GetType()
}

func (Instance *MiddlewareNode) IsLeaf() bool {
	return Instance.Identity.IsLeaf()
}

func (Instance *MiddlewareNode) AddChild(Path string, Child Componet.Node) error {
	return  Instance.Container.AddChild(Path,Child)
}

func (Instance *MiddlewareNode) SetChild(Path string,New string) error {
	return Instance.Container.SetChild(Path,New)
}

func (Instance *MiddlewareNode) GetChildrenLength() int {
	return Instance.Container.GetChildrenLength()
}

func (Instance *MiddlewareNode) GetChild(Path string) Componet.Node {
	return Instance.Container.GetChild(Path)
}

func (Instance *MiddlewareNode) HasChildren() bool {
	return Instance.Container.HasChildren()
}

func (Instance *MiddlewareNode) Apply(Handler http.Handler) http.Handler {
	return Instance.MiddlewareHandler.Apply(Handler)
}


func (Instance *MiddlewareNode) SetMiddleware(Middleware Middleware.Middleware) error {
	return Instance.MiddlewareHandler.SetMiddleware(Middleware)
}

func (Instance *MiddlewareNode) DeleteChild(Key string) error {
	return Instance.Container.DeleteChild(Key)
}

func (Instance *MiddlewareNode) GetAllChildren() []Componet.Node {
	return Instance.Container.GetAllChildren()
}