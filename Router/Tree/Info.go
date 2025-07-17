package Tree

import ()

type NodeType uint16
type PriorityLevel uint16
type MethodType string

const (
	RootType = iota
	StaticType
	CatchAllType
	WildCardType
	MiddlewareType
)

const (
	High   = 3
	Middle = 2
	Low    = 1
)

const (
	WildCardPrefix = ':'
	CatchAllPrefix = '*'
	PathSeparator  = "/"
)

const (
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	CONNECT = "CONNECT"
	PATCH   = "PATCH"
)
