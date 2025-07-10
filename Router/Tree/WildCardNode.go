package Tree

import (
	"context"
	"net/http"
	Componet "LiteFrame/Router/Tree/Componet"
)

type treeKey string

type WildCardNode struct {
	Identity *Componet.Identity
	Container *Componet.Container[Componet.Node]
	PathHandler *Componet.PathHandler
	EndPoint *Componet.EndPoint
	Data string
}

func NewWildCardNode(Path string) WildCardNode {
	WildCardNode :=  WildCardNode{
		Identity: Componet.NewIdentity(Componet.Middle, Componet.WildCardType, false),
		PathHandler: Componet.NewPathHandler(Componet.NewError(Componet.WildCardType,"",Path), Path[1:]),
		Container: Componet.NewContainer(Componet.NewError(Componet.WildCardType,"",Path)),
		EndPoint: nil,
		Data: "",
	}
	return  WildCardNode
}

func (Instance *WildCardNode) GetPriority() Componet.PriorityLevel {
	return Instance.Identity.GetPriority()
}

func (Instance *WildCardNode) GetType() Componet.NodeType {
	return Instance.Identity.GetType()
}

func (Instance *WildCardNode) IsLeaf() bool {
	return Instance.EndPoint != nil
} 

func (Instance *WildCardNode)  GetPath() string {
	return  Instance.PathHandler.GetPath()
}

func (Instance *WildCardNode) SetPath(Path string) error {
	return  Instance.PathHandler.SetPath(Path)
}

func (Instance *WildCardNode) Match(Path string) (bool, int, string) {
	if Path == "" {
		return false , 0 , ""
	}
	Instance.Data = Path;
	return true, len(Path), Path
} 

func (Instance *WildCardNode) AddChild(Path string, Child Componet.Node) error {
	return Instance.Container.AddChild(Path, Child)
}

func (Instance *WildCardNode) GetChild(Path string) Componet.Node {
	return Instance.Container.GetChild(Path)
}

func (Instance *WildCardNode) DeleteChild(Key string) error {
	return Instance.Container.DeleteChild(Key)
}

func (Instance *WildCardNode) GetChildrenLength() int {
	return Instance.Container.GetChildrenLength()
}

func (Instance *WildCardNode) GetAllChildren() []Componet.Node {
	return Instance.Container.GetAllChildren()
}

func (Instance *WildCardNode) HasChildren() bool {
	return Instance.Container.HasChildren()
}

func (Instance *WildCardNode) HasMethod(Method string) bool {
	if Instance.EndPoint == nil {
		return false
	}
	return Instance.EndPoint.HasMethod(Method)
}

func (Instance *WildCardNode) GetHandler(Method string) http.HandlerFunc {
	if Instance.EndPoint == nil {
		return nil
	}
	return  Instance.EndPoint.GetHandler(Method)
}

func (Instance *WildCardNode) SetHandler(Method string, Handler http.HandlerFunc) error {
	if Instance.EndPoint == nil {
		Instance.EndPoint = Componet.NewEndPoint(Componet.NewError(Componet.WildCardType,"",Instance.GetPath()))
	}
	return Instance.EndPoint.SetHandler(Method, Handler)
}

func (Instance *WildCardNode) GetAllHandlers() map[string]http.HandlerFunc {
	if Instance.EndPoint == nil {
		return make(map[string]http.HandlerFunc)
	}
	return Instance.EndPoint.GetAllHandlers()
}

func (Instance *WildCardNode) DeleteHandler(Method string) error {
	if Instance.EndPoint == nil {
		return Componet.NewError(Componet.WildCardType, "No handlers to delete", Instance.GetPath())
	}
	return Instance.EndPoint.DeleteHandler(Method)
}

func (Instance *WildCardNode) GetMethodCount() int {
	if Instance.EndPoint == nil {
		return 0
	}
	return Instance.EndPoint.GetMethodCount()
}

func (Instance *WildCardNode) GetAllMethods() []string {
	if Instance.EndPoint == nil {
		return []string{}
	}
	return Instance.EndPoint.GetAllMethods()
}

func (Instance *WildCardNode) SetWildCard(Handler http.Handler) http.Handler {
	return http.HandlerFunc(func(Writer http.ResponseWriter, Request *http.Request) {
		var contextKey treeKey = treeKey(Instance.PathHandler.GetPath())
		ctx := context.WithValue(Request.Context(), contextKey, Instance.Data)
		NewRequest := Request.WithContext(ctx)
		Handler.ServeHTTP(Writer, NewRequest)
	})
}