// Package Tree contains constant and type definitions.
// Defines enumerations and constants for node types, HTTP methods, priorities, etc.
package Tree

import (
	"LiteFrame/Router/Types"
)

// HandlerFunc is a handler function type imported from the Types package.
// Reuses types defined in separate packages to prevent circular references.
type HandlerFunc = Types.HandlerFunc

// NodeType is an enumeration representing the type of tree node.
// Each node can have only one type, and routing behavior is determined by the type.
type NodeType uint32

// MethodType is an enumeration representing HTTP methods.
// Used as array index to access handlers with O(1) time complexity in handler arrays.
type MethodType uint32

// NodeType constants: Define various types of tree nodes.
// Each node type affects routing performance and matching priority.
const (
	RootType       = iota // Root node ("/" path, starting point of tree)
	StaticType            // Static path node (/users, /api, etc., fastest matching)
	CatchAllType          // Catch-all node (*path, matches all remaining paths, low priority)
	WildCardType          // Wildcard node (:param, single segment matching, medium priority)
)

// Path pattern constants: Define special characters used in URL paths.
const (
	WildCardPrefix = ':' // Wildcard parameter prefix (:id, :name, etc.)
	CatchAllPrefix = '*' // Catch-all parameter prefix (*path, *file, etc.)
	PathSeparator  = '/' // Path separator
)

// HTTP method constants: Define HTTP methods following RFC 7231 standard.
// Each method is used as array index to provide O(1) handler access.
const (
	GET        = iota // GET method - Resource retrieval (idempotent, safe)
	HEAD              // HEAD method - Header information only (same as GET but no body)
	OPTIONS           // OPTIONS method - Supported method inquiry (CORS preflight)
	TRACE             // TRACE method - Request path tracing (for debugging purposes)
	POST              // POST method - Resource creation (non-idempotent)
	PUT               // PUT method - Resource creation/full modification (idempotent)
	DELETE            // DELETE method - Resource deletion (idempotent)
	CONNECT           // CONNECT method - Tunnel connection (for proxy servers)
	PATCH             // PATCH method - Partial resource modification (RFC 5789)
	NotAllowed        // Unsupported method (for 405 Method Not Allowed response)
)
