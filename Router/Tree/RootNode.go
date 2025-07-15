package Tree

import Component "LiteFrame/Router/Tree/Component"

type RootNode struct {
	Identity Component.Identity
	Container Component.Container[Component.Node]
}

func NewRootNode() RootNode {
	return RootNode{
		Identity: Component.NewIdentity(Component.High,Component.RootType,false),
		Container: Component.NewContainer(Component.NewError(Component.RootType,"","/")),
	}
}

func (Instance *RootNode) GetPriority() Component.PriorityLevel {
	return Instance.Identity.GetPriority()
}

func (Instance *RootNode) GetType() Component.NodeType {
	return Instance.Identity.GetType()
}

func (Instance *RootNode) IsLeaf() bool {
	return Instance.Identity.IsLeaf()
}

func (Instance *RootNode) AddChild(Path string, Child Component.Node) error {
	return  Instance.Container.AddChild(Path,Child)
}

func (Instance *RootNode) SetChild(Path string,New string) error {
	return Instance.Container.SetChild(Path,New)
}

func (Instance *RootNode) GetChildrenLength() int {
	return Instance.Container.GetChildrenLength()
}

func (Instance *RootNode) GetChild(Path string) Component.Node {
	return Instance.Container.GetChild(Path)
}

func (Instance *RootNode) HasChildren() bool {
	return Instance.Container.HasChildren()
}

func (Instance *RootNode) DeleteChild(Key string) error {
	return Instance.Container.DeleteChild(Key)
}

func (Instance *RootNode) GetAllChildren() []Component.Node {
	return Instance.Container.GetAllChildren()
}