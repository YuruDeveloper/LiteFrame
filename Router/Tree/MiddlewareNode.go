package Tree

import (
	"LiteFrame/Router/Middleware"
	"net/http"
	Component "LiteFrame/Router/Tree/Component"
)

type MiddlewareNode struct {
	Identity Component.Identity
	Container Component.Container[Component.Node]
	MiddlewareHandler *Component.MiddlewareHandler
}

func NewMiddlewareNode() MiddlewareNode {
	return MiddlewareNode{
		Identity: Component.NewIdentity(Component.High, Component.MiddlewareType, false),
		Container: Component.NewContainer(Component.NewError(Component.MiddlewareType,"","")),
		MiddlewareHandler: Component.NewMiddlewareHandler(Component.NewError(Component.MiddlewareType,"","")),
	}
}

func (Instance *MiddlewareNode) GetPriority() Component.PriorityLevel {
	return Instance.Identity.GetPriority()
}

func (Instance *MiddlewareNode) GetType() Component.NodeType {
	return Instance.Identity.GetType()
}

func (Instance *MiddlewareNode) IsLeaf() bool {
	return Instance.Identity.IsLeaf()
}

func (Instance *MiddlewareNode) AddChild(Path string, Child Component.Node) error {
	return  Instance.Container.AddChild(Path,Child)
}

func (Instance *MiddlewareNode) SetChild(Path string,New string) error {
	return Instance.Container.SetChild(Path,New)
}

func (Instance *MiddlewareNode) GetChildrenLength() int {
	return Instance.Container.GetChildrenLength()
}

func (Instance *MiddlewareNode) GetChild(Path string) Component.Node {
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

func (Instance *MiddlewareNode) GetAllChildren() []Component.Node {
	return Instance.Container.GetAllChildren()
}