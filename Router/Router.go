package Router

import (
	"net/http"
)


// tobo: all change
type Router struct {
	NotFoundHandler http.HandlerFunc
	NotAllowedHandler http.HandlerFunc
}

func NotFoundDefault(Writer http.ResponseWriter, Request *http.Request){
	http.Error(Writer,"404 Not Founded",http.StatusNotFound)
}

func NotAllowedDefault(Writer http.ResponseWriter,Request *http.Request) {
	http.Error(Writer,"Method Not Allowed",http.StatusMethodNotAllowed)
}

func NewRouter() *Router {
	return &Router{
		NotFoundHandler: NotFoundDefault,
		NotAllowedHandler: NotAllowedDefault,
	}
}








