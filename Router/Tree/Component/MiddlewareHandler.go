package Component

import (
	"net/http"
	"LiteFrame/Router/Middleware"
)

func NewMiddlewareHandler(Error error) *MiddlewareHandler {
	Err := Error.(*NodeError)
	return &MiddlewareHandler{
		Error:       *Err,
		Middlewares: make([]Middleware.Middleware, 0),
	}
}

type MiddlewareHandler struct {
	Error       NodeError
	Middlewares []Middleware.Middleware
}

func (Instance *MiddlewareHandler) SetMiddleware(Middleware Middleware.Middleware) error {
	if Middleware == nil {
		return Instance.Error.WithMessage("MiddleWare is nil")
	}
	Instance.Middlewares = append(Instance.Middlewares, Middleware)
	return nil
}

func (Instance *MiddlewareHandler) Apply(Handler http.Handler) http.Handler {
	var Temp http.Handler = Handler
	for _, Middleware := range Instance.Middlewares {
		Temp = Middleware.GetHandler()(Temp)
	}
	return Temp
}