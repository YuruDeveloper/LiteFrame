// Package Tree provides PathWithSegment structure for efficient path processing.
// A zero-allocation structure for iterating through URL paths segment by segment.
package Tree

// NewPathWithSegment creates a new PathWithSegment instance.
// Path: URL path string to analyze
// Initial state has both Start and End set to 0.
func NewPathWithSegment(path string) *PathWithSegment {
	return &PathWithSegment{
		Body:    path,
		BodyLen: len(path),
		Start:   0,
		End:     0,
	}
}

// PathWithSegment is a structure for processing URL paths segment by segment without memory allocation.
// It traverses path segments using indices without copying strings.
//
// Performance optimizations:
// - Zero allocation: References parts of existing string without creating new strings
// - Iterator pattern: Sequential segment movement through Next()
// - Boundary checking: Provides validation functions for safe index access
type PathWithSegment struct {
	Body    string // Original path string (immutable)
	BodyLen int
	Start   int // Start index of current segment
	End     int // End index of current segment (exclusive)
}

// Next moves to the next path segment.
// Skips path separators ('/') and sets start and end indices for the next segment.
//
// Operation:
// 1. Set current End position as new Start
// 2. Skip consecutive '/' characters
// 3. Set next segment until next '/' or end of string
func (instance *PathWithSegment) Next() {
	instance.Start = instance.End
	if instance.IsEnd() {
		return
	}

	// Skip consecutive path separators ('/')
	for instance.BodyLen > instance.Start && instance.Body[instance.Start] == '/' {
		instance.Start++
	}
	if instance.IsEnd() {
		instance.End = instance.Start
		return
	}
	// Set segment until next path separator or end of string
	instance.End = instance.Start
	for instance.End < instance.BodyLen && instance.Body[instance.End] != PathSeparator {
		instance.End++
	}
}

// IsEnd checks if the end of path has been reached.
// Returns true when Start index equals or exceeds string length.
//
//go:inline
func (instance *PathWithSegment) IsEnd() bool {
	return instance.Start >= instance.BodyLen
}

// IsSame checks if the current segment is empty.
// Returns true when Start equals End (zero-length segment).
//
//go:inline
func (instance *PathWithSegment) IsSame() bool {
	return instance.Start == instance.End
}

// IsNotValid checks if the current index state is invalid.
// Checks for end of path, index out of bounds, and logical errors.
func (instance *PathWithSegment) IsNotValid() bool {
	return instance.IsEnd() || instance.End > len(instance.Body) || instance.Start > instance.End
}

// Get returns the string of the current segment.
// Returns empty string if in invalid state.
// Memory allocation: Creates new string (call only when necessary)
func (instance *PathWithSegment) Get() string {
	if instance.IsNotValid() {
		return ""
	}
	return string(instance.Body[instance.Start:instance.End])
}

// GetToEnd returns the string from current position to end of path.
// Used by CatchAll routes to capture all remaining path.
func (instance *PathWithSegment) GetToEnd() string {
	return string(instance.Body[instance.Start:])
}

// GetLength returns the length of the current segment.
// Memory efficient as it calculates length without creating strings.
func (instance *PathWithSegment) GetLength() int {
	return instance.End - instance.Start
}
