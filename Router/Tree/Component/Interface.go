package Component

import (
	"net/http"
	"LiteFrame/Router/Middleware"
)

// 기본 타입들
type NodeType int
type PriorityLevel int
type TreeKey string

const (
	RootType = iota + 1
	StaticType 
	CatchAllType
	WildCardType
	MiddlewareType
)

const (
	High = 3
	Middle = 2
	Low = 1
)

// 인터페이스들 - Tree 패키지 순환 참조를 피하기 위해 제네릭 타입 대신 interface{} 사용
type Node interface {
	GetPriority() PriorityLevel 
	GetType() NodeType 
	IsLeaf() bool
}

type PathNode interface {
	GetPath() string
	SetPath(Path string) error
	Match(Path string) (Matched bool,MatchingChar int,LeftPath string)
}

type NodeContainer[T Node] interface {
	AddChild(Path string, Child T) error
	SetChild(Path string,New string) error
	GetChild(Path string) T
	DeleteChild(Key string) error
	GetChildrenLength() int
	GetAllChildren() []T
	HasChildren() bool
}

type HandleAccessor interface {
	GetHandler(Method string) http.HandlerFunc
	SetHandler(Method string, Handler http.HandlerFunc) error
	HasMethod(Method string) bool
	GetAllHandlers() map[string]http.HandlerFunc
	DeleteHandler(Method string) error
	GetMethodCount() int
	GetAllMethods() []string
}

type MiddlewareAcessor interface {
	SetMiddleware(Middleware Middleware.Middleware) error
	Apply(Handler http.Handler) http.Handler
}

type PathContainer[T Node] interface {
	PathNode
	NodeContainer[T]
	Split(SplitPoint int,NewNode PathContainer[T]) (PathContainer[T] , error)
}

type RootContainer[T Node] interface{
	NodeContainer[T]
	Node
}

type ContainerNode[T Node] interface {
	NodeContainer[T]
	PathNode
	Node
}

type HandlerNode interface {
	HandleAccessor
	PathNode
	Node
}

type MiddlewareContainerNode[T Node] interface {
	NodeContainer[T]
	MiddlewareAcessor
	PathNode
	Node
}

