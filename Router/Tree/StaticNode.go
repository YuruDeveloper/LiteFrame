package Tree

import (
    Component "LiteFrame/Router/Tree/Component"
	"net/http"
)

type StaticNode struct {
	Identity Component.Identity
	PathContainer Component.PathContainer[Component.Node]
	EndPoint *Component.EndPoint
}

func NewStaticNode(Path string) StaticNode {
	return StaticNode{
		Identity: Component.NewIdentity(Component.High, Component.StaticType, false),
		PathContainer: Component.NewPathContainerNode(Component.NewError(Component.StaticType,"",Path),Path),
		EndPoint: nil,
	}
	
}

func (Instance *StaticNode) GetPriority() Component.PriorityLevel {
	return Instance.Identity.GetPriority()
}

func (Instance *StaticNode) GetType() Component.NodeType {
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

func (Instance *StaticNode) Split(SplitPoint int, NewNode *Component.PathContainer[Component.Node]) (Component.PathContainer[Component.Node],error) {
	Result , Err := Instance.PathContainer.Split(SplitPoint, *NewNode)
	if Err != nil {
		return nil , Err
	}
	if HandlerNode , Ok := Result.(Component.HandlerNode) ; Ok {
		for _ , Method := range Instance.GetAllMethods() {
			HandlerNode.SetHandler(Method,Instance.GetHandler(Method))
		}
		for _ , Method := range Instance.GetAllMethods() {
			Instance.DeleteHandler(Method)
		}
	}
	return Result , nil
	//return &Result , nil
}

func (Instance *StaticNode) AddChild(Path string, Child Component.Node) error {
	return  Instance.PathContainer.AddChild(Path,Child)
}

func (Instance *StaticNode) SetChild(Path string,New string) error {
	return Instance.PathContainer.SetChild(Path,New)
}

func (Instance *StaticNode) GetChildrenLength() int {
	return Instance.PathContainer.GetChildrenLength()
}

func (Instance *StaticNode) GetChild(Path string) Component.Node {
	return Instance.PathContainer.GetChild(Path)
}

func (Instance *StaticNode) HasChildren() bool {
	return Instance.PathContainer.HasChildren()
}

func (Instance *StaticNode) DeleteChild(Key string) error {
	return Instance.PathContainer.DeleteChild(Key)
}

func (Instance *StaticNode) GetAllChildren() []Component.Node {
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
		Instance.EndPoint = Component.NewEndPoint(Component.NewError(Component.StaticType,"",Instance.PathContainer.GetPath()))
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
		return Component.NewError(Component.StaticType, "No handlers to delete", Instance.PathContainer.GetPath())
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
