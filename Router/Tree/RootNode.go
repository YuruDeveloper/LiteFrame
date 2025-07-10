package Tree

import Componet "LiteFrame/Router/Tree/Componet"

type RootNode struct {
	Identity *Componet.Identity
	Container *Componet.Container[Componet.Node]
}

func NewRootNode() RootNode {
	return RootNode{
		Identity: Componet.NewIdentity(Componet.High,Componet.RootType,false),
		Container: Componet.NewContainer(Componet.NewError(Componet.RootType,"","/")),
	}
}

func (Instance *RootNode) GetPriority() Componet.PriorityLevel {
	return Instance.Identity.GetPriority()
}

func (Instance *RootNode) GetType() Componet.NodeType {
	return Instance.Identity.GetType()
}

func (Instance *RootNode) IsLeaf() bool {
	return Instance.Identity.IsLeaf()
}

func (Instance *RootNode) AddChild(Path string, Child Componet.Node) error {
	return  Instance.Container.AddChild(Path,Child)
}

func (Instance *RootNode) SetChild(Path string,New string) error {
	return Instance.Container.SetChild(Path,New)
}

func (Instance *RootNode) GetChildrenLength() int {
	return Instance.Container.GetChildrenLength()
}

func (Instance *RootNode) GetChild(Path string) Componet.Node {
	return Instance.Container.GetChild(Path)
}

func (Instance *RootNode) HasChildren() bool {
	return Instance.Container.HasChildren()
}

func (Instance *RootNode) DeleteChild(Key string) error {
	return Instance.Container.DeleteChild(Key)
}

func (Instance *RootNode) GetAllChildren() []Componet.Node {
	return Instance.Container.GetAllChildren()
}