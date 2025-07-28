package Tree

func NewPathWithSegment(Path string) *PathWithSegment {
	return &PathWithSegment{
		Body: Path,
		Start:  0,
		End: 0,
	}
}


type PathWithSegment struct {
	Body string
	Start int 
	End int 
}

func (Instance *PathWithSegment) Next() {
	Instance.Start = Instance.End
	if Instance.Start >= len(Instance.Body) {
		return
	}

	for len(Instance.Body) > Instance.Start && Instance.Body[Instance.Start] == '/' {
		Instance.Start++
	}
	if Instance.Start >= len(Instance.Body) {
		Instance.End = Instance.Start
		return
	}
	Instance.End = Instance.Start
	for Instance.End < len(Instance.Body) && Instance.Body[Instance.End] != '/'{
		Instance.End++
	}
}

func (Instance *PathWithSegment) IsSame() bool {
	return Instance.Start == Instance.End
}

func (Instance *PathWithSegment) Get() string {
	if Instance.Start >= len(Instance.Body) || Instance.End > len(Instance.Body) || Instance.Start > Instance.End {
        return ""
    }
	return string(Instance.Body[Instance.Start:Instance.End])
}

func (Instance *PathWithSegment) GetToEnd() string {
	return string(Instance.Body[Instance.Start:])
}

func (Instance *PathWithSegment) GetLength() int {
	return Instance.End - Instance.Start
}
