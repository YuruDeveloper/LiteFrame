// Package Tree provides Radix Tree implementation for HTTP routing.
// Implements tree structure for efficient URL pattern matching and handler management.
package Tree

import (
	"LiteFrame/Router/Error"
	"LiteFrame/Router/Middleware"
	"LiteFrame/Router/Param"
	"net/http"
	"sort"
)

// Tree is a Radix Tree structure for HTTP routing.
// Manages root node, parameter pool, handlers, and middleware.
type Tree struct {
	RootNode           *Node                   // Root node of the tree
	Pool               *Param.ParamsPool       // Pool for parameter reuse
	NotFoundHandler    HandlerFunc             // 404 handler
	NotAllowedHandler  HandlerFunc             // 405 handler
	Middlewares        []Middleware.Middleware // Middleware list
	CompiledMiddleware func(HandlerFunc) HandlerFunc
}

// NewTree creates a new Tree instance.
// Initializes root node and parameter pool.
func NewTree() Tree {
	return Tree{
		RootNode: NewNode(RootType, "/"),
		Pool:     Param.NewParamsPool(),
	}
}

// IsWildCard checks if input string is a wildcard pattern (:param).
//
//go:inline
func (instance *Tree) IsWildCard(input string) bool {
	return len(input) > 0 && input[0] == WildCardPrefix
}

// IsCatchAll checks if input string is a catch-all pattern (*path).
//
//go:inline
func (instance *Tree) IsCatchAll(input string) bool {
	return len(input) > 0 && input[0] == CatchAllPrefix
}

// StringToMethodType converts HTTP method string to MethodType.
// Returns NotAllowed for unsupported methods.
//
//go:inline
func (instance *Tree) StringToMethodType(method string) MethodType {
	switch method {
	case "GET":
		return GET
	case "HEAD":
		return HEAD
	case "OPTIONS":
		return OPTIONS
	case "TRACE":
		return TRACE
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	case "CONNECT":
		return CONNECT
	case "PATCH":
		return PATCH
	default:
		return NotAllowed
	}
}

// Match finds common prefix between PathWithSegment and string, returning match results.
// Returns: (complete match, matched index, remaining PathWithSegment)
// Optimized matching algorithm based on PathWithSegment.
//
//go:inline
func (instance *Tree) Match(sourcePath PathWithSegment, targetPath string) (bool, int, PathWithSegment) {
	// Set comparison range based on shorter length between PathWithSegment and string

	sourceLength := sourcePath.GetLength()
	length := min(sourceLength, len(targetPath))
	// Calculate common prefix length by sequential byte comparison
	var index int
	for index = 0; index < length; index++ {
		if sourcePath.Body[sourcePath.Start+index] != targetPath[index] {
			break
		}
	}
	// Check if PathWithSegment is a complete prefix of target
	matched := index == sourcePath.GetLength()
	if matched {
		sourcePath.Start = sourcePath.End
	}
	// Update PathWithSegment to remaining unmatched portion
	if index < sourcePath.GetLength() {
		sourcePath.Start = sourcePath.Start + index
	}
	return matched, index, sourcePath
}

// SelectHandler selects method-appropriate handler from node and injects parameters into context.
// Returns NotAllowedHandler if handler not found.
// Important: Core function handling both memory pool management and context injection.
func (instance *Tree) SelectHandler(node *Node, method MethodType) HandlerFunc {
	if handler := node.Handlers[method]; handler != nil {
		// Inject parameters into context through closure and ensure memory pool return
		return handler
	}
	return instance.NotAllowedHandler
}

// InsertUniqueTypeChild inserts unique type child nodes (WildCard/CatchAll).
// Returns error if duplicate parameter names exist.
// Functional programming pattern: Eliminates duplicate code using higher-order functions

func (instance *Tree) InsertUniqueTypeChild(parent *Node, path string, target *Node, nodeType NodeType, errorFn func(string) error, setFn func(*Node, *Node)) (*Node, error) {
	switch {
	// Validate empty parameter name (cases with only ":" or "*")
	case path[1:] == "":
		return nil, Error.NewErrorWithCode(Error.NilParameter, path)
	// Conflict error for same type but different parameter name
	case target != nil && target.Param != path[1:]:
		return nil, errorFn(path)
	// Reuse if node with same parameter name already exists
	case target != nil && target.Param == path[1:]:
		return target, nil
	// Create new parameter node
	default:
		child := NewNode(nodeType, path)
		child.Param = path[1:] // Store only parameter name by removing prefix (: or *)
		setFn(parent, child)   // Type-specific node setting through function pointer
		return child, nil
	}
}

// InsertChild inserts child node into parent node.
// Creates Static, WildCard, or CatchAll nodes based on path type.
func (instance *Tree) InsertChild(parent *Node, path string) (*Node, error) {
	switch {
	case instance.IsWildCard(path):
		return instance.InsertUniqueTypeChild(parent, path, parent.WildCard, WildCardType,
			func(path string) error { return Error.NewErrorWithCode(Error.DuplicateWildCard, path) },
			func(parent, child *Node) { parent.WildCard = child })
	case instance.IsCatchAll(path):
		return instance.InsertUniqueTypeChild(parent, path, parent.CatchAll, CatchAllType,
			func(path string) error { return Error.NewErrorWithCode(Error.DuplicateCatchAll, path) },
			func(parent, child *Node) { parent.CatchAll = child })
	default:
		child := NewNode(StaticType, path)
		insertLocation := sort.Search(len(parent.Indices), func(index int) bool { return parent.Indices[index] >= path[0] })

		parent.Indices = append(parent.Indices, 0)
		copy(parent.Indices[insertLocation+1:], parent.Indices[insertLocation:])
		parent.Indices[insertLocation] = path[0]

		parent.Children = append(parent.Children, nil)
		copy(parent.Children[insertLocation+1:], parent.Children[insertLocation:])
		parent.Children[insertLocation] = child
		return child, nil
	}
}

// SplitNode splits existing node into two nodes at split point.
// Creates new parent node with common prefix and makes existing node a child.
//
//go:inline
func (instance *Tree) SplitNode(parent *Node, child *Node, splitPoint int) (*Node, error) {
	if splitPoint < 0 || splitPoint > len(child.Path) {
		return  nil , Error.NewErrorWithCode(Error.InvalidSplitPoint,child.Path) 
	}
	left := child.Path[:splitPoint]
	right := child.Path[splitPoint:]
	if len(left) == 0 {
		return nil, Error.NewErrorWithCode(Error.SplitFailed, child.Path)
	}
	newParent := NewNode(StaticType, left)
	for index, target := range parent.Children {
		if target == child {
			parent.Children[index] = newParent
			if len(right) > 0 {
				newParent.Indices = []byte{right[0]}
				child.Path = right
				newParent.Children = []*Node{child}
			} else {
				child.Path = right
				newParent.Indices = []byte{}
				newParent.Children = []*Node{}
			}
		}
	}
	return newParent, nil
}

// SetHandler registers handler for specified path and method in the tree.
// Uses PathWithSegment for efficient path processing and creates nodes or splits existing nodes as needed.
//
// Operation:
// 1. Analyze path segment by segment using PathWithSegment
// 2. Traverse tree searching for matching nodes
// 3. Complete match: Move to next segment
// 4. Partial match: Split node then continue
// 5. Match failure: Create new child node
func (instance *Tree) SetHandler(method MethodType, rawPath string, handler HandlerFunc) error {
	if rawPath == "" {
		return Error.NewErrorWithCode(Error.InvalidParameter, rawPath)
	}
	if method == NotAllowed {
		return Error.NewErrorWithCode(Error.MethodNotAllowed, rawPath)
	}
	if method != CONNECT && handler == nil {
		return Error.NewErrorWithCode(Error.InvalidParameter, rawPath)
	}
	// Efficient path processing using PathWithSegment structure
	path := NewPathWithSegment(rawPath)
	path.Next()
	if path.IsSame() {
		instance.RootNode.Handlers[method] = handler
		return nil
	}
	parent := instance.RootNode
setHelper:
	for {
		if path.IsSame() {
			parent.Handlers[method] = handler
			break
		}
		index := sort.Search(len(parent.Indices), func(index int) bool {
			return parent.Indices[index] >= path.Body[path.Start]
		})
		if index < len(parent.Indices) && parent.Indices[index] == path.Body[path.Start] {
			matched, matchingPoint, left := instance.Match(*path, parent.Children[index].Path)
			if matched {
				path.Next()
				parent = parent.Children[index]
				continue setHelper
			}
			if matchingPoint < len(parent.Children[index].Path) {
				newParent, err := instance.SplitNode(parent, parent.Children[index], matchingPoint)
				if err != nil {
					return err
				}
				parent = newParent
				if left.GetLength() > 0 {
					path = &left
					continue setHelper
				}
				path.Next()
				continue setHelper
			}
			parent = parent.Children[index]
			path = &left
			continue setHelper
		}
		child, err := instance.InsertChild(parent, path.Get())
		if err != nil {
			return err
		}
		parent = child
		path.Next()
		continue setHelper
	}
	return nil
}

// GetHandler finds and returns handler corresponding to HTTP request from tree.
// Uses PathWithSegment for efficient path processing and extracts parameters to select handler.
//
// Matching priority:
// 1. Static nodes: Exact string matching (highest priority)
// 2. WildCard nodes: Single segment parameters (:param)
// 3. CatchAll nodes: All remaining paths (*path, lowest priority)
//
// Returns: (handler function, parameter object) - returns nil if no parameters
func (instance *Tree) GetHandler(request *http.Request, getParams func() *Param.Params) (HandlerFunc, *Param.Params) {
	rawPath := request.URL.Path
	method := instance.StringToMethodType(request.Method)
	var params *Param.Params
	// Memory-efficient path processing using PathWithSegment
	path := NewPathWithSegment(rawPath)
	if method == NotAllowed {
		return instance.NotAllowedHandler, nil
	}
	path.Next()
	node := instance.RootNode
getHelper:
	for {
		if path.GetLength() == 0 {
			return instance.SelectHandler(node, method), params
		}

		for index := 0; index < len(node.Indices); index++ {
			if path.Body[path.Start] == node.Indices[index] {
				matched, matchingPoint, left := instance.Match(*path, node.Children[index].Path)
				switch {
				case matched:
					path.Next()
					node = node.Children[index]
					continue getHelper
				case matchingPoint > 0 && matchingPoint == len(node.Children[index].Path) && path.GetLength() > 0:
					node = node.Children[index]
					path = &left
					continue getHelper
				}
			}
		}
		// 2nd priority: WildCard node matching (single segment capture)
		if node.WildCard != nil {
			if params == nil {
				params = getParams()
			}
			// Store current segment as parameter and proceed to next segment
			params.Add(node.WildCard.Param, path.Get())
			path.Next()
			node = node.WildCard
			continue getHelper
		}
		// 3rd priority: CatchAll node matching (capture all remaining paths)
		if node.CatchAll != nil {
			if params == nil {
				params = getParams()
			}
			// Store remaining path of PathWithSegment as single parameter
			params.Add(node.CatchAll.Param, path.GetToEnd())
			path.Start = path.End
			// Search handler in CatchAll node with empty path (path consumption complete)
			node = node.CatchAll
			continue getHelper
		}
		// Return parameter object and 404 handler when no matching route found
		return instance.NotFoundHandler, params
	}
}

// SetMiddleware adds middleware to the tree.
// Added middleware applies to all handlers.
func (instance *Tree) SetMiddleware(middleware Middleware.Middleware) {
	instance.Middlewares = append(instance.Middlewares, middleware)
}

func (instance *Tree) CompileMiddlewares() {
	if len(instance.Middlewares) == 0 {
		instance.CompiledMiddleware = nil
		return
	}

	instance.CompiledMiddleware = func(baseHandler HandlerFunc) HandlerFunc {
		temp := baseHandler
		for index := len(instance.Middlewares) - 1; index > -1; index-- {
			temp = instance.Middlewares[index].GetHandler()(temp)
		}
		return temp
	}
}

// ApplyMiddleware applies registered middleware to handlers in reverse order.
// Last registered middleware becomes outermost layer.
func (instance *Tree) ApplyMiddleware(handler HandlerFunc) HandlerFunc {
	if instance.CompiledMiddleware == nil {
		return handler
	}
	return instance.CompiledMiddleware(handler)
}

// ServeHTTP implements http.Handler interface.
// Finds handler for request, applies middleware, then executes.
func (instance *Tree) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler, params := instance.GetHandler(request, instance.Pool.Get)
	handler = instance.ApplyMiddleware(handler)
	handler(writer, request, params)

	// Return parameter object to pool
	if params != nil {
		instance.Pool.Put(params)
	}
}
