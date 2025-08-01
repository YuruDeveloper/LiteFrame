// Package Router provides high-level HTTP router interface.
// It defines the basic router structure and default handlers.
package Router

import (
	"net/http"
)

// Router is the basic HTTP router structure.
// It includes default handlers for 404 (Not Found) and 405 (Method Not Allowed) errors.
// It will be integrated with the Tree router in the future to form a complete routing system.
type Router struct {
	NotFoundHandler   http.HandlerFunc // 404 error handler
	NotAllowedHandler http.HandlerFunc // 405 error handler
}

// NotFoundDefault is the default 404 error handler.
// It returns a standard HTTP 404 response when the requested resource cannot be found.
func NotFoundDefault(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "404 Not Found", http.StatusNotFound)
}

// NotAllowedDefault is the default 405 error handler.
// It returns a standard HTTP 405 response when a request is made with an unsupported HTTP method.
func NotAllowedDefault(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// NewRouter creates a new Router instance.
// It initializes the router with default error handlers.
// Users can replace them with custom error handlers as needed.
func NewRouter() *Router {
	return &Router{
		NotFoundHandler:   NotFoundDefault,   // Set default 404 handler
		NotAllowedHandler: NotAllowedDefault, // Set default 405 handler
	}
}
