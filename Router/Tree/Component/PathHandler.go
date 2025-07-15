package Component

func NewPathHandler(Error error, Path string) PathHandler {
	Err := Error.(*NodeError)
	return PathHandler{
		Error: *Err,
		Path:  Path,
	}
}

type PathHandler struct {
	Error NodeError
	Path  string
}

func (Instance *PathHandler) GetPath() string {
	return Instance.Path
}

func (Instance *PathHandler) SetPath(Path string) error {
	if Path == "" {
		return Instance.Error.WithMessage("Path is empty")
	}
	Instance.Path = Path
	return nil
}

func (Instance *PathHandler) Match(Path string) (Matched bool, MatchingChar int, LeftPath string) {
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