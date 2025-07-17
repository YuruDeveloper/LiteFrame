package Tree

import "fmt"

type TreeError struct {
	Message string
	Path    string
}

func (Instance *TreeError) Error() string {
	return fmt.Sprintf("error: %s From %s", Instance.Message, Instance.Path)
}

func (Instance *TreeError) WithMessage(message string) error {
	return &TreeError{
		Message: message,
		Path:    Instance.Path,
	}
}

func NewTreeError(Message string, path string) error {
	return &TreeError{
		Message: Message,
		Path:    path,
	}
}
