package Param

import (
	"context"
	"sync"
	"testing"
)

// ======================
// Param Constructor Tests
// ======================

func TestNewParams(t *testing.T) {
	params := NewParams()

	t.Run("params_not_nil", func(t *testing.T) {
		if params == nil {
			t.Error("Expected non-nil Params instance")
		}
	})

	t.Run("path_initialized", func(t *testing.T) {
		if params.Path != "" {
			t.Error("Expected empty Path")
		}
	})

	t.Run("count_zero", func(t *testing.T) {
		if params.Count != 0 {
			t.Errorf("Expected Count to be 0, got %d", params.Count)
		}
	})

	t.Run("overflow_empty", func(t *testing.T) {
		if params.Overflow != nil {
			t.Error("Expected Overflow to be nil")
		}
	})
}

// ======================
// Param Add Tests
// ======================

func TestParamsAdd(t *testing.T) {
	t.Run("add_single_param", func(t *testing.T) {
		params := NewParams()
		params.Path = "/users/123/posts/456"
		params.Add("id", 7, 10) // "123"

		if params.Count != 1 {
			t.Errorf("Expected Count 1, got %d", params.Count)
		}

		if params.Fix[0].Key != "id" {
			t.Errorf("Expected key 'id', got '%s'", params.Fix[0].Key)
		}

		if params.Fix[0].Start != 7 || params.Fix[0].End != 10 {
			t.Errorf("Expected Start:7, End:10, got Start:%d, End:%d", params.Fix[0].Start, params.Fix[0].End)
		}
	})

	t.Run("add_multiple_params", func(t *testing.T) {
		params := NewParams()
		params.Path = "/users/123/posts/456/category/books"
		params.Add("id", 7, 10)       // "123"
		params.Add("postId", 17, 20)  // "456"
		params.Add("category", 30, 35) // "books"

		if params.Count != 3 {
			t.Errorf("Expected Count 3, got %d", params.Count)
		}

		// Check fixed parameters (first 2)
		if params.Fix[0].Key != "id" {
			t.Errorf("Expected key 'id', got '%s'", params.Fix[0].Key)
		}
		if params.Fix[1].Key != "postId" {
			t.Errorf("Expected key 'postId', got '%s'", params.Fix[1].Key)
		}

		// Check overflow parameter (3rd)
		if len(params.Overflow) != 1 {
			t.Errorf("Expected 1 overflow parameter, got %d", len(params.Overflow))
		}
		if params.Overflow[0].Key != "category" {
			t.Errorf("Expected overflow key 'category', got '%s'", params.Overflow[0].Key)
		}
	})

	t.Run("add_empty_key", func(t *testing.T) {
		params := NewParams()
		params.Path = "/users/123"
		params.Add("", 7, 10)

		if params.Count != 1 {
			t.Errorf("Expected Count 1, got %d", params.Count)
		}

		if params.Fix[0].Key != "" {
			t.Errorf("Expected empty key, got '%s'", params.Fix[0].Key)
		}
	})
}

// ======================
// GetByName Tests
// ======================

func TestParamsGetByName(t *testing.T) {
	t.Run("existing_param", func(t *testing.T) {
		params := NewParams()
		params.Path = "/users/123/posts/test"
		params.Add("id", 7, 10)     // "123"
		params.Add("name", 17, 21)  // "test"

		value := params.GetByName("id")
		if value != "123" {
			t.Errorf("Expected value '123', got '%s'", value)
		}

		value = params.GetByName("name")
		if value != "test" {
			t.Errorf("Expected value 'test', got '%s'", value)
		}
	})

	t.Run("non_existing_param", func(t *testing.T) {
		params := NewParams()
		params.Path = "/users/123"
		params.Add("id", 7, 10) // "123"

		value := params.GetByName("name")
		if value != "" {
			t.Errorf("Expected empty string, got '%s'", value)
		}
	})

	t.Run("empty_params", func(t *testing.T) {
		params := NewParams()
		params.Path = "/users"

		value := params.GetByName("id")
		if value != "" {
			t.Errorf("Expected empty string, got '%s'", value)
		}
	})

	t.Run("overflow_param", func(t *testing.T) {
		params := NewParams()
		params.Path = "/users/123/posts/456/category/books"
		params.Add("id", 7, 10)       // "123" (Fix[0])
		params.Add("postId", 17, 20)  // "456" (Fix[1])
		params.Add("category", 30, 35) // "books" (Overflow[0])

		value := params.GetByName("category")
		if value != "books" {
			t.Errorf("Expected 'books', got '%s'", value)
		}
	})

	t.Run("case_sensitive", func(t *testing.T) {
		params := NewParams()
		params.Path = "/users/123"
		params.Add("ID", 7, 10) // "123"

		value := params.GetByName("id")
		if value != "" {
			t.Errorf("Expected empty string for case mismatch, got '%s'", value)
		}

		value = params.GetByName("ID")
		if value != "123" {
			t.Errorf("Expected '123', got '%s'", value)
		}
	})
}

// ======================
// Context Tests
// ======================

func TestGetParamsFromCTX(t *testing.T) {
	t.Run("existing_params_in_context", func(t *testing.T) {
		params := NewParams()
		params.Path = "/users/123"
		params.Add("id", 7, 10) // "123"

		ctx := context.WithValue(context.Background(), Key{}, params)

		retrievedParams, ok := GetParamsFromCTX(ctx)
		if !ok {
			t.Error("Expected to retrieve params from context")
		}

		if retrievedParams == nil {
			t.Fatal("Expected non-nil params")
		}

		if retrievedParams.Count != 1 {
			t.Errorf("Expected Count 1, got %d", retrievedParams.Count)
		}

		if retrievedParams.GetByName("id") != "123" {
			t.Errorf("Expected value '123', got '%s'", retrievedParams.GetByName("id"))
		}
	})

	t.Run("no_params_in_context", func(t *testing.T) {
		ctx := context.Background()

		retrievedParams, ok := GetParamsFromCTX(ctx)
		if ok {
			t.Error("Expected no params in context")
		}

		if retrievedParams != nil {
			t.Error("Expected nil params")
		}
	})

	t.Run("wrong_type_in_context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Key{}, "not_params")

		retrievedParams, ok := GetParamsFromCTX(ctx)
		if ok {
			t.Error("Expected type assertion to fail")
		}

		if retrievedParams != nil {
			t.Error("Expected nil params")
		}
	})

	t.Run("nil_value_in_context", func(t *testing.T) {
		var nilParams *Params
		ctx := context.WithValue(context.Background(), Key{}, nilParams)

		retrievedParams, ok := GetParamsFromCTX(ctx)
		if !ok {
			t.Error("Expected successful type assertion even for nil value")
		}

		if retrievedParams != nil {
			t.Error("Expected nil params")
		}
	})
}

// ======================
// ParamsPool Tests
// ======================

func TestNewParamsPool(t *testing.T) {
	pool := NewParamsPool()

	t.Run("pool_not_nil", func(t *testing.T) {
		if pool == nil {
			t.Error("Expected non-nil ParamsPool instance")
		}
	})

	t.Run("sync_pool_initialized", func(t *testing.T) {
		if pool.Pool == nil {
			t.Error("Expected Pool to be initialized")
		}
	})

	t.Run("pool_factory_function", func(t *testing.T) {
		// Test that the factory function works
		obj := pool.Pool.Get()
		params, ok := obj.(*Params)
		if !ok {
			t.Error("Expected factory function to return *Params")
		}
		if params == nil {
			t.Error("Expected non-nil Params from factory")
		}
		pool.Pool.Put(obj)
	})

	t.Run("pool_pre_warmed", func(t *testing.T) {
		// Pool should have pre-created objects
		// We can't directly test this, but we can verify that Get() returns objects immediately
		for i := 0; i < DefaultBufferSize; i++ {
			obj := pool.Pool.Get()
			if obj == nil {
				t.Errorf("Expected non-nil object from pool at iteration %d", i)
			}
			pool.Pool.Put(obj)
		}
	})
}

func TestParamsPoolGet(t *testing.T) {
	pool := NewParamsPool()

	t.Run("get_returns_params", func(t *testing.T) {
		params := pool.Get()

		if params == nil {
			t.Error("Expected non-nil Params from pool")
			return
		}

		if params.Count != 0 {
			t.Errorf("Expected Count 0, got %d", params.Count)
		}

		if params.Path != "" {
			t.Errorf("Expected empty Path, got '%s'", params.Path)
		}
	})

	t.Run("get_resets_existing_params", func(t *testing.T) {
		params := pool.Get()
		params.Path = "/test/value"
		params.Add("test", 6, 11) // "value"

		// Put it back and get it again
		pool.Put(params)
		resetParams := pool.Get()

		if resetParams.Count != 0 {
			t.Errorf("Expected reset Count to be 0, got %d", resetParams.Count)
		}
		if resetParams.Path != "" {
			t.Errorf("Expected reset Path to be empty, got '%s'", resetParams.Path)
		}
	})

	t.Run("multiple_gets", func(t *testing.T) {
		params1 := pool.Get()
		params2 := pool.Get()
		params3 := pool.Get()

		if params1 == nil || params2 == nil || params3 == nil {
			t.Error("Expected all Gets to return non-nil Params")
		}

		// They should all be independent
		params1.Path = "/test/value1"
		params1.Add("key1", 6, 12) // "value1"
		
		params2.Path = "/test/value2"
		params2.Add("key2", 6, 12) // "value2"
		
		params3.Path = "/test/value3"
		params3.Add("key3", 6, 12) // "value3"

		if params1.GetByName("key2") != "" || params1.GetByName("key3") != "" {
			t.Error("Expected params1 to only contain its own data")
		}

		if params2.GetByName("key1") != "" || params2.GetByName("key3") != "" {
			t.Error("Expected params2 to only contain its own data")
		}

		if params3.GetByName("key1") != "" || params3.GetByName("key2") != "" {
			t.Error("Expected params3 to only contain its own data")
		}
	})
}

func TestParamsPoolPut(t *testing.T) {
	pool := NewParamsPool()

	t.Run("put_nil_params", func(t *testing.T) {
		// Should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("The code panicked when putting nil into the pool: %v", r)
			}
		}()
		pool.Put(nil)
	})

	t.Run("put_and_get_cycle", func(t *testing.T) {
		params := pool.Get()
		params.Path = "/test/value"
		params.Add("test", 6, 11) // "value"

		pool.Put(params)

		retrievedParams := pool.Get()
		if retrievedParams.Count != 0 {
			t.Errorf("Expected retrieved params Count to be 0, got %d", retrievedParams.Count)
		}
		if retrievedParams.Path != "" {
			t.Errorf("Expected retrieved params Path to be empty, got '%s'", retrievedParams.Path)
		}
	})

	t.Run("put_large_capacity_params", func(t *testing.T) {
		params := pool.Get()

		// Add many parameters to exceed MaxSize capacity in overflow
		for i := 0; i < MaxSize+5; i++ {
			params.Add("key", 6, 11) // "value"
		}

		originalCap := cap(params.Overflow)
		if originalCap <= MaxSize {
			t.Skipf("Overflow capacity %d not large enough to test capacity reset", originalCap)
		}

		pool.Put(params)

		// Get it back and check if overflow capacity was reset
		retrievedParams := pool.Get()
		if cap(retrievedParams.Overflow) > MaxSize {
			t.Errorf("Expected overflow capacity to be reset to 0, got %d", cap(retrievedParams.Overflow))
		}
	})

	t.Run("put_normal_capacity_params", func(t *testing.T) {
		params := pool.Get()
		params.Path = "/test/value1/value2"
		params.Add("key1", 6, 12)   // "value1"
		params.Add("key2", 13, 19)  // "value2"

		pool.Put(params)
		retrievedParams := pool.Get()

		// Overflow capacity should be preserved if under MaxSize
		if len(params.Overflow) <= MaxSize && cap(retrievedParams.Overflow) > MaxSize {
			t.Errorf("Expected normal capacity to be preserved")
		}
	})
}

// ======================
// Concurrency Tests
// ======================

func TestParamsPoolConcurrency(t *testing.T) {
	pool := NewParamsPool()
	const numGoroutines = 100
	const numOperations = 100

	t.Run("concurrent_get_put", func(t *testing.T) {
		var wg sync.WaitGroup

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				for j := 0; j < numOperations; j++ {
					params := pool.Get()
					params.Path = "/users/test/data"
					params.Add("id", 7, 11)     // "test"
					params.Add("value", 12, 16) // "data"

					// Verify the params work correctly
					if params.GetByName("id") != "test" {
						t.Errorf("Goroutine %d: Expected 'test', got '%s'", id, params.GetByName("id"))
					}

					pool.Put(params)
				}
			}(i)
		}

		wg.Wait()
	})
}

// ======================
// Edge Cases Tests
// ======================

func TestParamsEdgeCases(t *testing.T) {
	t.Run("very_long_key_value", func(t *testing.T) {
		params := NewParams()
		longKey := string(make([]byte, 1000))
		longValue := string(make([]byte, 1000))

		params.Path = "/" + longValue
		params.Add(longKey, 1, 1001) // longValue from position 1 to 1001

		if params.GetByName(longKey) != longValue {
			t.Error("Failed to handle very long key/value pairs")
		}
	})

	t.Run("unicode_characters", func(t *testing.T) {
		params := NewParams()
		params.Path = "/한글값/🎯"
		params.Add("한글키", 1, 10)  // "한글값" (UTF-8 bytes)
		params.Add("🔑", 11, 15)    // "🎯" (UTF-8 bytes)

		if params.GetByName("한글키") != "한글값" {
			t.Error("Failed to handle Korean characters")
		}

		if params.GetByName("🔑") != "🎯" {
			t.Error("Failed to handle emoji characters")
		}
	})

	t.Run("special_characters", func(t *testing.T) {
		params := NewParams()
		params.Path = "/value with spaces/value/with/slashes/value=with=equals"
		params.Add("key with spaces", 1, 18)    // "value with spaces"
		params.Add("key/with/slashes", 19, 37)  // "value/with/slashes"
		params.Add("key?with&query", 38, 55)    // "value=with=equals"

		if params.GetByName("key with spaces") != "value with spaces" {
			t.Error("Failed to handle spaces")
		}

		if params.GetByName("key/with/slashes") != "value/with/slashes" {
			t.Error("Failed to handle slashes")
		}

		if params.GetByName("key?with&query") != "value=with=equals" {
			t.Error("Failed to handle query characters")
		}
	})
}

// ======================
// Constants Tests
// ======================

func TestConstants(t *testing.T) {
	t.Run("default_values", func(t *testing.T) {
		if DefaultSize != 2 {
			t.Errorf("Expected DefaultSize to be 2, got %d", DefaultSize)
		}

		if MaxSize != 8 {
			t.Errorf("Expected MaxSize to be 8, got %d", MaxSize)
		}

		if DefaultBufferSize != 10 {
			t.Errorf("Expected DefaultBufferSize to be 10, got %d", DefaultBufferSize)
		}
	})
}
