package Middleware

import (
	"net/http"
)

type Middleware interface {
	GetHandler() MiddleWareFunc
}

type MiddleWareFunc func(http.Handler) http.Handler
