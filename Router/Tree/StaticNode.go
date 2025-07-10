package Tree

import (
	"net/http"
	Componet "LiteFrame/Router/Tree/Componet"
)

type StaticNode struct {
	Identity *Componet.Identity
	PathContainer *Componet.PathContainerNode[Componet.Node]
	EndPoint *Componet.EndPoint
}

func NewStaticNode(Path string) StaticNode {
	return StaticNode{
		Identity: Componet.NewIdentity(Componet.High, Componet.StaticType, false),
		PathContainer: Componet.NewPathContainerNode(Componet.NewError(Componet.StaticType,"",Path),Path),
		EndPoint: nil,
	}
}

func (Instance *StaticNode) GetPriority() Componet.PriorityLevel {
	return Instance.Identity.GetPriority()
}

func (Instance *StaticNode) GetType() Componet.NodeType {
	return Instance.Identity.GetType()
}

func (Instance *StaticNode) IsLeaf() bool {
	return Instance.EndPoint != nil
}

func (Instance *StaticNode)  GetPath() string {
	return  Instance.PathContainer.GetPath()
}

func (Instance *StaticNode) SetPath(Path string) error {
	return  Instance.PathContainer.SetPath(Path)
}

func (Instance *StaticNode) Match(Path string) (bool, int, string) {
	return  Instance.PathContainer.Match(Path)
}

func (Instance *StaticNode) Split(SplitPoint int, NewNode *Componet.PathContainerNode[Componet.Node]) (*Componet.PathContainerNode[Componet.Node],error) {
	Result , Err := Instance.PathContainer.Split(SplitPoint, *NewNode)
	if Err != nil {
		return nil , Err
	}
	if Instance.EndPoint != nil {
		// TODO: EndPoint 핸들러들을 NewNode로 복사하는 로직 구현 필요
	}
	return &Result , nil
}

func (Instance *StaticNode) AddChild(Path string, Child Componet.Node) error {
	return  Instance.PathContainer.AddChild(Path,Child)
}

func (Instance *StaticNode) SetChild(Path string,New string) error {
	return Instance.PathContainer.SetChild(Path,New)
}

func (Instance *StaticNode) GetChildrenLength() int {
	return Instance.PathContainer.GetChildrenLength()
}

func (Instance *StaticNode) GetChild(Path string) Componet.Node {
	return Instance.PathContainer.GetChild(Path)
}

func (Instance *StaticNode) HasChildren() bool {
	return Instance.PathContainer.HasChildren()
}

func (Instance *StaticNode) DeleteChild(Key string) error {
	return Instance.PathContainer.DeleteChild(Key)
}

func (Instance *StaticNode) GetAllChildren() []Componet.Node {
	return Instance.PathContainer.GetAllChildren()
}

func (Instance *StaticNode) GetHandler(Method string) http.HandlerFunc {
	if Instance.EndPoint == nil {
		return nil
	}	
	return Instance.EndPoint.GetHandler(Method)
}

func (Instance *StaticNode) SetHandler(Method string,Handler http.HandlerFunc) error {
	if Instance.EndPoint == nil {
		Instance.EndPoint = Componet.NewEndPoint(Componet.NewError(Componet.StaticType,"",Instance.PathContainer.GetPath()))
	}
	return Instance.EndPoint.SetHandler(Method,Handler)
}	

func (Instance *StaticNode) HasMethod(Method string) bool {
	if Instance.EndPoint == nil {
		return false
	}
	return Instance.EndPoint.HasMethod(Method)
}

func (Instance *StaticNode) GetAllHandlers() map[string]http.HandlerFunc {
	if Instance.EndPoint == nil {
		return make(map[string]http.HandlerFunc)
	}
	return Instance.EndPoint.GetAllHandlers()
}

func (Instance *StaticNode) DeleteHandler(Method string) error {
	if Instance.EndPoint == nil {
		return Componet.NewError(Componet.StaticType, "No handlers to delete", Instance.PathContainer.GetPath())
	}
	return Instance.EndPoint.DeleteHandler(Method)
}

func (Instance *StaticNode) GetMethodCount() int {
	if Instance.EndPoint == nil {
		return 0
	}
	return Instance.EndPoint.GetMethodCount()
}

func (Instance *StaticNode) GetAllMethods() []string {
	if Instance.EndPoint == nil {
		return []string{}
	}
	return Instance.EndPoint.GetAllMethods()
}
