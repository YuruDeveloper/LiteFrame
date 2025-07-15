package Tree

import (
	"net/http"
	Component "LiteFrame/Router/Tree/Component"
)

type CatchAllNode struct {
	Identity Component.Identity
	PathHandler Component.PathHandler
	EndPoint *Component.EndPoint
}
func NewCatchAllNode(Path string) CatchAllNode {
	Node := CatchAllNode{
		Identity: Component.NewIdentity(Component.Low, Component.CatchAllType, true),
		PathHandler: Component.NewPathHandler(Component.NewError(Component.CatchAllType,"",Path), Path),
		EndPoint: Component.NewEndPoint(Component.NewError(Component.CatchAllType,"",Path)),
	}
	return Node 
	
}


func (Instance *CatchAllNode) GetPriority() Component.PriorityLevel {
	return Instance.Identity.GetPriority()
}

func (Instance *CatchAllNode) GetType() Component.NodeType {
	return Instance.Identity.GetType()
}

func (Instance *CatchAllNode) IsLeaf() bool {
	return Instance.Identity.IsLeaf()
} 

func (Instance *CatchAllNode)  GetPath() string {
	return  Instance.PathHandler.GetPath()
}

func (Instance *CatchAllNode) SetPath(Path string) error {
	return  Instance.PathHandler.SetPath(Path)
}

func (Instance *CatchAllNode) Match(Path string) (bool, int, string) {
	return true, len(Path), ""
} 

func (Instance *CatchAllNode) GetHandler(Method string) http.HandlerFunc {
	return  Instance.EndPoint.GetHandler(Method)
}

func (Instance *CatchAllNode) SetHandler(Method string, Handler http.HandlerFunc) error {
	return Instance.EndPoint.SetHandler(Method, Handler)
}

func (Instance *CatchAllNode) HasMethod(Method string) bool {
	return Instance.EndPoint.HasMethod(Method)
}

func (Instance *CatchAllNode) GetAllHandlers() map[string]http.HandlerFunc {
	if Instance.EndPoint == nil {
		return make(map[string]http.HandlerFunc)
	}
	return Instance.EndPoint.GetAllHandlers()
}

func (Instance *CatchAllNode) DeleteHandler(Method string) error {
	if Instance.EndPoint == nil {
		return Component.NewError(Component.CatchAllType, "No handlers to delete", Instance.GetPath())
	}
	return Instance.EndPoint.DeleteHandler(Method)
}

func (Instance *CatchAllNode) GetMethodCount() int {
	if Instance.EndPoint == nil {
		return 0
	}
	return Instance.EndPoint.GetMethodCount()
}

func (Instance *CatchAllNode) GetAllMethods() []string {
	if Instance.EndPoint == nil {
		return []string{}
	}
	return Instance.EndPoint.GetAllMethods()
}
