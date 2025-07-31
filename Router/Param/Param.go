// Package Param provides parameter management functionality used in HTTP routing.
// Supports parameter pooling and context-based parameter passing for performance optimization.
package Param

import (
	"context"
	"sync"
)

const (
	DefaultSize       = 2
	MaxSize           = 8
	DefaultBufferSize = 10
)

// NewParams creates a new Params instance.
// Initializes parameter list with default capacity of 2.
// Performance optimization: Most routes have 2 or fewer parameters, minimizing memory reallocation.
func NewParams() *Params {
	return &Params{
		List: make([]Param, 0, DefaultSize), // Default capacity 2: Value optimized for common API patterns
	}
}

// Param is a structure representing a key-value pair of a single parameter.
// Stores parameter name and value extracted from URL path.
type Param struct {
	Key   string // Parameter name (name extracted from :)
	Value string // Parameter value (actual value extracted from URL)
}

// Params is a structure that stores multiple parameters.
// Manages all parameters extracted from a single HTTP request.
type Params struct {
	List []Param // Parameter list
}

// Key is an empty structure for identifying parameters in context.
// Used as a key for parameter storage in context.WithValue.
type Key struct{}

// Add adds a new parameter to the parameter list.
// Stores parameter name and value extracted from URL path.
func (instance *Params) Add(key string, value string) {
	instance.List = append(instance.List, Param{Key: key, Value: value})
}

// GetByName searches for the corresponding value by parameter name.
// Returns empty string if parameter is not found.
func (instance *Params) GetByName(name string) string {
	for _, param := range instance.List {
		if param.Key == name {
			return param.Value
		}
	}
	return ""
}

// GetParamsFromCTX extracts parameters from context.
// Gets parameters from request context in HTTP handlers.
func GetParamsFromCTX(ctx context.Context) (*Params, bool) {
	temp, success := (ctx.Value(Key{})).(*Params)
	return temp, success
}

// NewParamsPool creates a new parameter pool.
// Uses sync.Pool to reuse parameter objects for performance optimization.
// Pre-allocates 10 parameter objects initially and puts them in the pool.
// Warm-up: Prevents memory allocation delays in initial requests.
func NewParamsPool() *ParamsPool {
	instance := &ParamsPool{
		Pool: &sync.Pool{
			// Factory function: Creates new object when pool is empty
			New: func() any {
				return NewParams()
			},
		},
	}
	// Manual warm-up: Pre-create 10 objects to reduce latency of first requests
	for index := 0; index < DefaultBufferSize; index++ {
		instance.Put(NewParams())
	}
	return instance
}

// ParamsPool is a pool structure for reusing Params objects.
// Uses sync.Pool to reduce memory allocation and garbage collection overhead.
type ParamsPool struct {
	Pool *sync.Pool // Parameter object pool
}

// Get retrieves a parameter object from the pool.
// Initializes existing parameter list to prepare for use in new requests.
// Memory efficiency: Maintains capacity and resets only length to 0 for reuse without reallocation
func (instance *ParamsPool) Get() *Params {
	object := instance.Pool.Get().(*Params)
	// Slice reset: Maintains existing memory while setting length to 0 (performance optimization)
	object.List = object.List[:0]
	return object
}

// Put returns a used parameter object to the pool.
// Creates a new slice if list capacity exceeds 8 to prevent memory leaks.
// Important memory management: Prevents excessive capacity growth to maintain memory pool efficiency
func (instance *ParamsPool) Put(object *Params) {
	if object != nil {
		// Prevent memory surge: Create new slice if capacity exceeds threshold (8)
		// This prevents continuous memory growth after requests with large parameters
		if cap(object.List) > MaxSize {
			object.List = make([]Param, 0, DefaultSize) // Reset to default capacity
		}
		instance.Pool.Put(object)
	}
}
