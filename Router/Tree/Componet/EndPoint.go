package Componet

import "net/http"

func NewEndPoint(Error error) *EndPoint {
	Err := Error.(*NodeError)
	return &EndPoint{
		Error:    *Err,
		Handlers: make(map[string]http.HandlerFunc),
	}
}

type EndPoint struct {
	Error    NodeError
	Handlers map[string]http.HandlerFunc
}

func (Instance *EndPoint) GetHandler(Method string) http.HandlerFunc {
	if Instance.Handlers == nil {
		return nil
	}
	return Instance.Handlers[Method]
}

func (Instance *EndPoint) SetHandler(Method string, Handler http.HandlerFunc) error {
	if Handler == nil || Method == "" {
		return Instance.Error.WithMessage("Method and Handler is Required")
	}
	if Instance.Handlers == nil {
		Instance.Handlers = make(map[string]http.HandlerFunc)
	}
	Instance.Handlers[Method] = Handler
	return nil
}

func (Instance *EndPoint) HasMethod(Method string) bool {
	if Instance.Handlers == nil {
		return false
	}
	_, Exists := Instance.Handlers[Method]
	return Exists
}

func (Instance *EndPoint) GetAllHandlers() map[string]http.HandlerFunc {
	if Instance.Handlers == nil {
		return make(map[string]http.HandlerFunc)
	}
	// 복사본을 반환하여 원본 데이터 보호
	HandlersCopy := make(map[string]http.HandlerFunc)
	for Method, Handler := range Instance.Handlers {
		HandlersCopy[Method] = Handler
	}
	return HandlersCopy
}

func (Instance *EndPoint) DeleteHandler(Method string) error {
	if Method == "" {
		return Instance.Error.WithMessage("Method is required")
	}
	if Instance.Handlers == nil {
		return Instance.Error.WithMessage("No handlers to delete")
	}
	if _, Exists := Instance.Handlers[Method]; !Exists {
		return Instance.Error.WithMessage("Method does not exist")
	}
	delete(Instance.Handlers, Method)
	return nil
}

func (Instance *EndPoint) GetMethodCount() int {
	if Instance.Handlers == nil {
		return 0
	}
	return len(Instance.Handlers)
}

func (Instance *EndPoint) GetAllMethods() []string {
	if Instance.Handlers == nil {
		return []string{}
	}
	Methods := make([]string, 0, len(Instance.Handlers))
	for Method := range Instance.Handlers {
		Methods = append(Methods, Method)
	}
	return Methods
}