// Package Tree defines Node structure and related functions.
// Contains structures and constructor functions representing each node of the Radix Tree.
package Tree

// NewNode creates a new Node instance.
// nodeType: Node type (Root, Static, WildCard, CatchAll, Middleware)
// path: Path segment that the node represents
// Performance optimization: Minimizes memory allocation and initializes only necessary fields.
func NewNode(nodeType NodeType, path string) *Node {
	return &Node{
		Type:     nodeType,
		Path:     path,
		Children: make([]*Node, 0),                     // Dynamically expandable child node slice
		Handlers: make([]HandlerFunc, int(NotAllowed)), // Handler array for all HTTP method support
		WildCard: nil,                                  // Wildcard node is created only when needed
		CatchAll: nil,                                  // CatchAll node is created only when needed
	}
}

// Node is a structure representing each node of the Radix Tree.
// Stores path segments, child nodes, and handlers for each HTTP method.
//
// Design for memory efficiency:
// - Static nodes: Use Children + Indices (O(1) search)
// - WildCard/CatchAll: Managed with separate pointers (memory saving)
// - Handlers: Direct access through array index (performance optimization)
type Node struct {
	Type     NodeType      // Node type (Root, Static, WildCard, CatchAll, Middleware)
	Path     string        // Path segment represented by the node (compressed path)
	Indices  []byte        // First byte index of child nodes (O(1) search optimization)
	Children []*Node       // Static child nodes (1:1 correspondence with Indices)
	Handlers []HandlerFunc // Handler array for each HTTP method (using MethodType as index)
	WildCard *Node         // Wildcard child node (:param, single segment matching)
	CatchAll *Node         // CatchAll child node (*path, remaining all path matching)
	Param    string        // Parameter name (used only in WildCard/CatchAll nodes, excluding ':' '*')
}
