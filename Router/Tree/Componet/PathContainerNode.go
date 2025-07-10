package Componet

func NewPathContainerNode(Error error,Path string) *PathContainerNode[Node] {
	Err := Error.(*NodeError)
	return &PathContainerNode[Node] {
		Error: *Err,
		Path: Path,
		Box: make(map[string]Node),
	}
}

type PathContainerNode[T Node] struct {
	Error NodeError
	Path  string
	Box   map[string]T
}

func (Instance *PathContainerNode[T]) GetPath() string {
	return Instance.Path
}

func (Instance *PathContainerNode[T]) SetPath(Path string) error {
	if Path == "" {
		return Instance.Error.WithMessage("Path is empty")
	}
	Instance.Path = Path
	return nil
}

func (Instance *PathContainerNode[T]) Match(Path string) (Matched bool, MatchingChar int, LeftPath string) {
	MinLen := len(Path)
	if len(Instance.Path) < MinLen {
		MinLen = len(Instance.Path)
	}

	for MatchingChar < MinLen {
		if Path[MatchingChar] != Instance.Path[MatchingChar] {
			break
		}
		MatchingChar++
	}
	Matched = MatchingChar == len(Instance.Path)
	if MatchingChar < len(Path) {
		LeftPath = Path[MatchingChar:]
	}
	return Matched, MatchingChar, LeftPath
}

func (Instance *PathContainerNode[T]) AddChild(Path string, Child T) error {
	if Path == "" {
		return Instance.Error.WithMessage("Path is Empty")
	}
	Instance.Box[Path] = Child
	return nil
}

func (Instance *PathContainerNode[T]) SetChild(Path string, New string) error {
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

func (Instance *PathContainerNode[T]) GetChild(Path string) T {
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

func (Instance *PathContainerNode[T]) DeleteChild(Key string) error {
	if Key == "" {
		return Instance.Error.WithMessage("Path is Empty")
	}
	if _, Exists := Instance.Box[Key]; !Exists {
		return Instance.Error.WithMessage("Object is not exists")
	}
	delete(Instance.Box, Key)
	return nil
}

func (Instance *PathContainerNode[T]) GetAllChildren() []T {
	Children := make([]T, 0, len(Instance.Box))
	for _, Child := range Instance.Box {
		Children = append(Children, Child)
	}
	return Children
}

func (Instance *PathContainerNode[T]) GetChildrenLength() int {
	return len(Instance.Box)
}

func (Instance *PathContainerNode[T]) HasChildren() bool {
	return len(Instance.Box) > 0
}

func (Instance *PathContainerNode[T]) Split(SplitPoint int, NewNode PathContainerNode[T]) (PathContainerNode[T], error) {
	var Zero PathContainerNode[T]
	if SplitPoint <= 0 || SplitPoint >= len(Instance.Path) {
		return Zero , Instance.Error.WithMessage("Invalid Split Point")
	}

	LeftPath := Instance.Path[:SplitPoint]
	RightPath := Instance.Path[SplitPoint:]

	if Err := NewNode.SetPath(LeftPath); Err != nil {
		return Zero, Instance.Error.WithMessage("Failed to set Path on New Node")
	}
	Instance.Path = RightPath
	var Path PathNode
	// 자식 물러주기
	for _ , Child := range Instance.Box {
		Path , _ = any(Child).(PathNode)
		// in path , T
		NewNode.AddChild(Path.GetPath(),Child);
	}
	for k := range Instance.Box {
		delete(Instance.Box, k)
	}
	return NewNode , nil
}