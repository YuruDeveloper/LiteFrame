package Componet

type NodeError struct {
	Type    NodeType
	Message string
	Path    string
}

var NodeTypeNames = map[NodeType]string{
	RootType:       "Root",
	StaticType:     "Static",
	WildCardType:   "WildCard",
	CatchAllType:   "CatchAll",
	MiddlewareType: "Middleware",
}

func (Instance *NodeError) Error() string {
	return Instance.Message
}

func (Instance *NodeError) WithMessage(message string) error {
	return &NodeError{
		Type:    Instance.Type,
		Message: message,
		Path:    Instance.Path,
	}
}

func (Instance *NodeError) WithPath(path string) error {
	return &NodeError{
		Type:    Instance.Type,
		Message: Instance.Message,
		Path:    path,
	}
}

func NewError(nodeType NodeType, message, path string) error {
	return &NodeError{
		Type:    nodeType,
		Message: message,
		Path:    path,
	}
}