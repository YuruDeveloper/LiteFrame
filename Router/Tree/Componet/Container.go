package Componet

func NewContainer(Error error) *Container[Node] {
	Err := Error.(*NodeError)
	return &Container[Node]{
		Error: *Err,
		Box:   make(map[string]Node),
	}
}

type Container[T Node] struct {
	Error NodeError
	Box   map[string]T
}

func (Instance *Container[T]) AddChild(Path string, Child T) error {
	if Path == "" {
		return Instance.Error.WithMessage("Path is Empty")
	}
	Instance.Box[Path] = Child
	return nil
}

func (Instance *Container[T]) SetChild(Path string, New string) error {
	if Path == "" || New == "" {
		return Instance.Error.WithMessage("Path is Empty")
	}
	Temp, OK := Instance.Box[Path]

	if !OK {
		return Instance.Error.WithMessage("Wrong Path")
	}
	delete(Instance.Box, Path)
	Instance.Box[New] = Temp
	return nil
}

func (Instance *Container[T]) GetChild(Path string) T {
	if Path == "" {
		var zero T
		return zero
	}
	Temp, OK := Instance.Box[Path]
	if !OK {
		var zero T
		return zero
	}
	return Temp
}

func (Instance *Container[T]) DeleteChild(Key string) error {
	if Key == "" {
		return Instance.Error.WithMessage("Path is Empty")
	}
	if _, Exists := Instance.Box[Key]; !Exists {
		return Instance.Error.WithMessage("Object is not exists")
	}
	delete(Instance.Box, Key)
	return nil
}

func (Instance *Container[T]) GetAllChildren() []T {
	Children := make([]T, 0, len(Instance.Box))
	for _, Child := range Instance.Box {
		Children = append(Children, Child)
	}
	return Children
}

func (Instance *Container[T]) GetChildrenLength() int {
	return len(Instance.Box)
}

func (Instance *Container[T]) HasChildren() bool {
	return len(Instance.Box) > 0
}