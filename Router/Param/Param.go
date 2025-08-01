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
		Count: 0,
	}
}

// Param is a structure representing a key-value pair of a single parameter.
// Stores parameter name and value extracted from URL path.
type Param struct {
	Key   string // Parameter name (name extracted from :)
	Start int
	End int
}

// Params is a structure that stores multiple parameters.
// Manages all parameters extracted from a single HTTP request.
type Params struct {
	Path string
	Count int
	Fix [DefaultSize]Param // Parameter list
	Overflow []Param
}

// Key is an empty structure for identifying parameters in context.
// Used as a key for parameter storage in context.WithValue.
type Key struct{}

func (instance *Params) Reset() {
	instance.Path = ""
	instance.Count = 0
	instance.Fix[0] = Param{}
	instance.Fix[1] = Param{}
	instance.Overflow = instance.Overflow[:0]
}

// Add adds a new parameter to the parameter list.
// Stores parameter name and value extracted from URL path.
func (instance *Params) Add(key string, start int , end int) {
	if instance.Count < 2 {
		instance.Fix[instance.Count].Key = key
		instance.Fix[instance.Count].Start = start
		instance.Fix[instance.Count].End =end
		instance.Count++
		return
	}
	instance.Overflow = append(instance.Overflow, Param{Key: key,Start: start,End: end})
	instance.Count++
}

// GetByName searches for the corresponding value by parameter name.
// Returns empty string if parameter is not found.
func (instance *Params) GetByName(name string) string {
	if instance.Count > 0 {
		if instance.Fix[0].Key == name {
			return instance.Path[instance.Fix[0].Start:instance.Fix[0].End]
		} 
	}
	if instance.Count > 1 {
		if instance.Fix[1].Key == name {
			return instance.Path[instance.Fix[1].Start:instance.Fix[1].End]
		}
	}
	if instance.Count > 2 {
		for index := 0; index < len(instance.Overflow); index++ {
			if instance.Overflow[index].Key == name {
				return instance.Path[instance.Overflow[index].Start:instance.Overflow[index].End]
			}
		}
	}
	return  ""
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
	object.Reset()
	return object
}

// Put returns a used parameter object to the pool.
// Creates a new slice if list capacity exceeds 8 to prevent memory leaks.
// Important memory management: Prevents excessive capacity growth to maintain memory pool efficiency
func (instance *ParamsPool) Put(object *Params) {
	if object != nil {
		// Prevent memory surge: Create new slice if capacity exceeds threshold (8)
		// This prevents continuous memory growth after requests with large parameters
		if cap(object.Overflow) > MaxSize {
			object.Overflow = make([]Param, 0) // Reset to default capacity
		}
		instance.Pool.Put(object)
	}
}
