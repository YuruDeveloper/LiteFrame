// Package Types defines common types used in the Router system.
// Separates HandlerFunc and other common types into a separate package to prevent circular references.
package Types

import (
	"LiteFrame/Router/Param"
	"net/http"
)

// HandlerFunc is an HTTP handler function type.
// Parameters are passed through context in GetHandler.
//
// Function signature:
// - http.ResponseWriter: Interface for writing HTTP responses
// - *http.Request: Pointer to structure containing HTTP request information
// - *Param.Params: Parameters extracted from URL path (nil if no parameters)
type HandlerFunc func(http.ResponseWriter, *http.Request, *Param.Params)
