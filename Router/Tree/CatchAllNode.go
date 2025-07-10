package Tree

import (
	"net/http"
	Componet "LiteFrame/Router/Tree/Componet"
)

type CatchAllNode struct {
	Identity *Componet.Identity
	PathHandler *Componet.PathHandler
	EndPoint *Componet.EndPoint
}
func NewCatchAllNode(Path string) CatchAllNode {
	Node := CatchAllNode{
		Identity: Componet.NewIdentity(Componet.Low, Componet.CatchAllType, true),
		PathHandler: Componet.NewPathHandler(Componet.NewError(Componet.CatchAllType,"",Path), Path),
		EndPoint: Componet.NewEndPoint(Componet.NewError(Componet.CatchAllType,"",Path)),
	}
	return Node 
	
}


func (Instance *CatchAllNode) GetPriority() Componet.PriorityLevel {
	return Instance.Identity.GetPriority()
}

func (Instance *CatchAllNode) GetType() Componet.NodeType {
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
		return Componet.NewError(Componet.CatchAllType, "No handlers to delete", Instance.GetPath())
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
