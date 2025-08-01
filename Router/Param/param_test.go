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

	t.Run("list_initialized", func(t *testing.T) {
		if params.List == nil {
			t.Error("Expected List to be initialized")
		}
	})

	t.Run("list_empty", func(t *testing.T) {
		if len(params.List) != 0 {
			t.Errorf("Expected empty list, got length %d", len(params.List))
		}
	})

	t.Run("list_capacity", func(t *testing.T) {
		if cap(params.List) != DefaultSize {
			t.Errorf("Expected capacity %d, got %d", DefaultSize, cap(params.List))
		}
	})
}

// ======================
// Param Add Tests
// ======================

func TestParamsAdd(t *testing.T) {
	params := NewParams()

	t.Run("add_single_param", func(t *testing.T) {
		params.Add("id", "123")

		if len(params.List) != 1 {
			t.Errorf("Expected 1 parameter, got %d", len(params.List))
		}

		if params.List[0].Key != "id" {
			t.Errorf("Expected key 'id', got '%s'", params.List[0].Key)
		}

		if params.List[0].Value != "123" {
			t.Errorf("Expected value '123', got '%s'", params.List[0].Value)
		}
	})

	t.Run("add_multiple_params", func(t *testing.T) {
		params := NewParams()
		params.Add("id", "123")
		params.Add("name", "test")
		params.Add("category", "books")

		if len(params.List) != 3 {
			t.Errorf("Expected 3 parameters, got %d", len(params.List))
		}

		expectedParams := []Param{
			{Key: "id", Value: "123"},
			{Key: "name", Value: "test"},
			{Key: "category", Value: "books"},
		}

		for i, expected := range expectedParams {
			if params.List[i].Key != expected.Key {
				t.Errorf("Expected key '%s' at index %d, got '%s'", expected.Key, i, params.List[i].Key)
			}
			if params.List[i].Value != expected.Value {
				t.Errorf("Expected value '%s' at index %d, got '%s'", expected.Value, i, params.List[i].Value)
			}
		}
	})

	t.Run("add_empty_values", func(t *testing.T) {
		params := NewParams()
		params.Add("", "")
		params.Add("key", "")
		params.Add("", "value")

		if len(params.List) != 3 {
			t.Errorf("Expected 3 parameters, got %d", len(params.List))
		}

		if params.List[0].Key != "" || params.List[0].Value != "" {
			t.Error("Expected empty key and value")
		}

		if params.List[1].Key != "key" || params.List[1].Value != "" {
			t.Error("Expected key 'key' and empty value")
		}

		if params.List[2].Key != "" || params.List[2].Value != "value" {
			t.Error("Expected empty key and value 'value'")
		}
	})
}

// ======================
// GetByName Tests
// ======================

func TestParamsGetByName(t *testing.T) {
	t.Run("existing_param", func(t *testing.T) {
		params := NewParams()
		params.Add("id", "123")
		params.Add("name", "test")

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
		params.Add("id", "123")

		value := params.GetByName("name")
		if value != "" {
			t.Errorf("Expected empty string, got '%s'", value)
		}
	})

	t.Run("empty_params", func(t *testing.T) {
		params := NewParams()

		value := params.GetByName("id")
		if value != "" {
			t.Errorf("Expected empty string, got '%s'", value)
		}
	})

	t.Run("duplicate_keys", func(t *testing.T) {
		params := NewParams()
		params.Add("id", "first")
		params.Add("id", "second")

		// Should return the first match
		value := params.GetByName("id")
		if value != "first" {
			t.Errorf("Expected 'first', got '%s'", value)
		}
	})

	t.Run("case_sensitive", func(t *testing.T) {
		params := NewParams()
		params.Add("ID", "123")

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
		params.Add("id", "123")

		ctx := context.WithValue(context.Background(), Key{}, params)

		retrievedParams, ok := GetParamsFromCTX(ctx)
		if !ok {
			t.Error("Expected to retrieve params from context")
		}

		if retrievedParams == nil {
			t.Fatal("Expected non-nil params")
		}

		if len(retrievedParams.List) != 1 {
			t.Errorf("Expected 1 parameter, got %d", len(retrievedParams.List))
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

		if params.List == nil {
			t.Error("Expected List to be initialized")
			return
		}

		if len(params.List) != 0 {
			t.Errorf("Expected empty list, got length %d", len(params.List))
		}
	})

	t.Run("get_resets_existing_params", func(t *testing.T) {
		params := pool.Get()
		params.Add("test", "value")

		// Put it back and get it again
		pool.Put(params)
		resetParams := pool.Get()

		if len(resetParams.List) != 0 {
			t.Errorf("Expected reset list to be empty, got length %d", len(resetParams.List))
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
		params1.Add("key1", "value1")
		params2.Add("key2", "value2")
		params3.Add("key3", "value3")

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
		params.Add("test", "value")

		pool.Put(params)

		retrievedParams := pool.Get()
		if len(retrievedParams.List) != 0 {
			t.Errorf("Expected retrieved params to be reset, got length %d", len(retrievedParams.List))
		}
	})

	t.Run("put_large_capacity_params", func(t *testing.T) {
		params := pool.Get()

		// Add many parameters to exceed MaxSize capacity
		for i := 0; i < MaxSize+5; i++ {
			params.Add("key", "value")
		}

		originalCap := cap(params.List)
		if originalCap <= MaxSize {
			t.Skipf("Capacity %d not large enough to test capacity reset", originalCap)
		}

		pool.Put(params)

		// Get it back and check if capacity was reset
		retrievedParams := pool.Get()
		if cap(retrievedParams.List) != DefaultSize {
			t.Errorf("Expected capacity to be reset to %d, got %d", DefaultSize, cap(retrievedParams.List))
		}
	})

	t.Run("put_normal_capacity_params", func(t *testing.T) {
		params := pool.Get()
		params.Add("key1", "value1")
		params.Add("key2", "value2")

		originalCap := cap(params.List)

		pool.Put(params)
		retrievedParams := pool.Get()

		// Capacity should be preserved if under MaxSize
		if cap(retrievedParams.List) < originalCap {
			t.Errorf("Expected capacity to be preserved, original: %d, retrieved: %d", originalCap, cap(retrievedParams.List))
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
					params.Add("id", "test")
					params.Add("value", "data")

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

		params.Add(longKey, longValue)

		if params.GetByName(longKey) != longValue {
			t.Error("Failed to handle very long key/value pairs")
		}
	})

	t.Run("unicode_characters", func(t *testing.T) {
		params := NewParams()
		params.Add("한글키", "한글값")
		params.Add("🔑", "🎯")

		if params.GetByName("한글키") != "한글값" {
			t.Error("Failed to handle Korean characters")
		}

		if params.GetByName("🔑") != "🎯" {
			t.Error("Failed to handle emoji characters")
		}
	})

	t.Run("special_characters", func(t *testing.T) {
		params := NewParams()
		params.Add("key with spaces", "value with spaces")
		params.Add("key/with/slashes", "value/with/slashes")
		params.Add("key?with&query", "value=with=equals")

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
