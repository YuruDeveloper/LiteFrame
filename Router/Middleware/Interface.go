// Package Middleware defines interfaces for HTTP middleware systems.
// Supports middleware patterns that can execute common logic before and after request processing.
package Middleware

import (
	"LiteFrame/Router/Types"
)

// Middleware is an interface that middleware implementations must follow.
// All middleware must implement the GetHandler method that returns MiddleWareFunc.
//
// Usage example:
//
//	type LoggingMiddleware struct{}
//	func (m LoggingMiddleware) GetHandler() MiddleWareFunc {
//	    return func(next http.Handler) http.Handler {
//	        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	            log.Printf("Request: %s %s", r.Method, r.URL.Path)
//	            next.ServeHTTP(w, r)
//	        })
//	    }
//	}
type Middleware interface {
	GetHandler() MiddleWareFunc // Method that returns middleware function
}

// MiddleWareFunc is a middleware function type.
// A higher-order function that takes the next handler and returns a wrapped handler.
// This pattern allows multiple middleware to be connected in a chain.
type MiddleWareFunc func(Types.HandlerFunc) Types.HandlerFunc
