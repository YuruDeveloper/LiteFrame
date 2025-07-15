package Component

func NewIdentity(Priority PriorityLevel, Type NodeType, Leaf bool) Identity {
	return Identity{
		Priority: Priority,
		Type:     Type,
		Leaf:     Leaf,
	}
}

type Identity struct {
	Priority PriorityLevel
	Type     NodeType
	Leaf     bool
}

func (Instance *Identity) GetPriority() PriorityLevel {
	return Instance.Priority
}

func (Instance *Identity) GetType() NodeType {
	return Instance.Type
}

func (Instance *Identity) IsLeaf() bool {
	return Instance.Leaf
}