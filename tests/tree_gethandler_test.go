package tests

import (
	"LiteFrame/Router/Param"
	"LiteFrame/Router/Tree"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetHandler tests the GetHandler functionality
func TestGetHandler(t *testing.T) {
	// Helper function to create a test handler with response text
	createHandlerWithResponse := func(response string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(response))
		}
	}

	// Helper function to create a test handler that checks parameters
	createParamCheckHandler := func(expectedParams map[string]string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			params, ok := Param.GetParamsFromCTX(r.Context())
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("no params in context"))
				return
			}

			for key, expectedValue := range expectedParams {
				actualValue := params.GetByName(key)
				if actualValue != expectedValue {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("param mismatch: " + key))
					return
				}
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("params matched"))
		}
	}

	t.Run("root_path", func(t *testing.T) {
		tree := Tree.NewTree()
		handler := createHandlerWithResponse("root")

		err := tree.SetHandler("GET", "/", handler)
		if err != nil {
			t.Fatalf("SetHandler failed: %v", err)
		}

		req := httptest.NewRequest("GET", "/", nil)
		handlerFunc := tree.GetHandler(req)

		recorder := httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}

		if recorder.Body.String() != "root" {
			t.Errorf("Expected 'root', got '%s'", recorder.Body.String())
		}
	})

	t.Run("static_path", func(t *testing.T) {
		tree := Tree.NewTree()
		handler := createHandlerWithResponse("users")

		err := tree.SetHandler("GET", "/users", handler)
		if err != nil {
			t.Fatalf("SetHandler failed: %v", err)
		}

		req := httptest.NewRequest("GET", "/users", nil)
		handlerFunc := tree.GetHandler(req)

		recorder := httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}

		if recorder.Body.String() != "users" {
			t.Errorf("Expected 'users', got '%s'", recorder.Body.String())
		}
	})

	t.Run("nested_static_path", func(t *testing.T) {
		tree := Tree.NewTree()
		handler := createHandlerWithResponse("user profile")

		err := tree.SetHandler("GET", "/users/profile", handler)
		if err != nil {
			t.Fatalf("SetHandler failed: %v", err)
		}

		req := httptest.NewRequest("GET", "/users/profile", nil)
		handlerFunc := tree.GetHandler(req)

		recorder := httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}

		if recorder.Body.String() != "user profile" {
			t.Errorf("Expected 'user profile', got '%s'", recorder.Body.String())
		}
	})

	t.Run("wildcard_path", func(t *testing.T) {
		tree := Tree.NewTree()
		expectedParams := map[string]string{"id": "123"}
		handler := createParamCheckHandler(expectedParams)

		err := tree.SetHandler("GET", "/users/:id", handler)
		if err != nil {
			t.Fatalf("SetHandler failed: %v", err)
		}

		req := httptest.NewRequest("GET", "/users/123", nil)
		handlerFunc := tree.GetHandler(req)

		recorder := httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}

		if recorder.Body.String() != "params matched" {
			t.Errorf("Expected 'params matched', got '%s'", recorder.Body.String())
		}
	})

	t.Run("multiple_wildcards", func(t *testing.T) {
		tree := Tree.NewTree()
		expectedParams := map[string]string{
			"userId": "123",
			"postId": "456",
		}
		handler := createParamCheckHandler(expectedParams)

		err := tree.SetHandler("GET", "/users/:userId/posts/:postId", handler)
		if err != nil {
			t.Fatalf("SetHandler failed: %v", err)
		}

		req := httptest.NewRequest("GET", "/users/123/posts/456", nil)
		handlerFunc := tree.GetHandler(req)

		recorder := httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}

		if recorder.Body.String() != "params matched" {
			t.Errorf("Expected 'params matched', got '%s'", recorder.Body.String())
		}
	})

	t.Run("catch_all_path", func(t *testing.T) {
		tree := Tree.NewTree()
		expectedParams := map[string]string{"path": "static/css/main.css"}
		handler := createParamCheckHandler(expectedParams)

		err := tree.SetHandler("GET", "/files/*path", handler)
		if err != nil {
			t.Fatalf("SetHandler failed: %v", err)
		}

		req := httptest.NewRequest("GET", "/files/static/css/main.css", nil)
		handlerFunc := tree.GetHandler(req)

		recorder := httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}

		if recorder.Body.String() != "params matched" {
			t.Errorf("Expected 'params matched', got '%s'", recorder.Body.String())
		}
	})

	t.Run("not_found", func(t *testing.T) {
		tree := Tree.NewTree()
		tree.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		}

		req := httptest.NewRequest("GET", "/nonexistent", nil)
		handlerFunc := tree.GetHandler(req)

		recorder := httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", recorder.Code)
		}

		if recorder.Body.String() != "not found" {
			t.Errorf("Expected 'not found', got '%s'", recorder.Body.String())
		}
	})

	t.Run("method_not_allowed", func(t *testing.T) {
		tree := Tree.NewTree()
		tree.NotAllowedHandler = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("method not allowed"))
		}

		handler := createHandlerWithResponse("post response")
		err := tree.SetHandler("POST", "/users", handler)
		if err != nil {
			t.Fatalf("SetHandler failed: %v", err)
		}

		req := httptest.NewRequest("GET", "/users", nil)
		handlerFunc := tree.GetHandler(req)

		recorder := httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", recorder.Code)
		}

		if recorder.Body.String() != "method not allowed" {
			t.Errorf("Expected 'method not allowed', got '%s'", recorder.Body.String())
		}
	})

	t.Run("different_methods_same_path", func(t *testing.T) {
		tree := Tree.NewTree()
		
		getHandler := createHandlerWithResponse("GET response")
		postHandler := createHandlerWithResponse("POST response")
		
		err := tree.SetHandler("GET", "/users", getHandler)
		if err != nil {
			t.Fatalf("SetHandler GET failed: %v", err)
		}
		
		err = tree.SetHandler("POST", "/users", postHandler)
		if err != nil {
			t.Fatalf("SetHandler POST failed: %v", err)
		}

		// Test GET request
		req := httptest.NewRequest("GET", "/users", nil)
		handlerFunc := tree.GetHandler(req)
		recorder := httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("GET: Expected status 200, got %d", recorder.Code)
		}
		if recorder.Body.String() != "GET response" {
			t.Errorf("GET: Expected 'GET response', got '%s'", recorder.Body.String())
		}

		// Test POST request
		req = httptest.NewRequest("POST", "/users", nil)
		handlerFunc = tree.GetHandler(req)
		recorder = httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("POST: Expected status 200, got %d", recorder.Code)
		}
		if recorder.Body.String() != "POST response" {
			t.Errorf("POST: Expected 'POST response', got '%s'", recorder.Body.String())
		}
	})

	t.Run("priority_static_over_wildcard", func(t *testing.T) {
		tree := Tree.NewTree()
		
		staticHandler := createHandlerWithResponse("static admin")
		wildcardHandler := createHandlerWithResponse("wildcard user")
		
		err := tree.SetHandler("GET", "/users/:id", wildcardHandler)
		if err != nil {
			t.Fatalf("SetHandler wildcard failed: %v", err)
		}
		
		err = tree.SetHandler("GET", "/users/admin", staticHandler)
		if err != nil {
			t.Fatalf("SetHandler static failed: %v", err)
		}

		// Test static route
		req := httptest.NewRequest("GET", "/users/admin", nil)
		handlerFunc := tree.GetHandler(req)
		recorder := httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}
		if recorder.Body.String() != "static admin" {
			t.Errorf("Expected 'static admin', got '%s'", recorder.Body.String())
		}

		// Test wildcard route
		req = httptest.NewRequest("GET", "/users/123", nil)
		handlerFunc = tree.GetHandler(req)
		recorder = httptest.NewRecorder()
		handlerFunc(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}
		if recorder.Body.String() != "wildcard user" {
			t.Errorf("Expected 'wildcard user', got '%s'", recorder.Body.String())
		}
	})

	t.Run("complex_routing", func(t *testing.T) {
		tree := Tree.NewTree()
		
		// Multiple routes with different patterns
		routes := []struct {
			method   string
			path     string
			response string
		}{
			{"GET", "/", "home"},
			{"GET", "/users", "users list"},
			{"GET", "/users/:id", "user detail"},
			{"POST", "/users", "create user"},
			{"GET", "/users/:id/posts", "user posts"},
			{"GET", "/users/:id/posts/:postId", "user post detail"},
			{"GET", "/files/*path", "file content"},
			{"GET", "/api/v1/users", "api users"},
		}

		for _, route := range routes {
			handler := createHandlerWithResponse(route.response)
			err := tree.SetHandler(route.method, route.path, handler)
			if err != nil {
				t.Fatalf("SetHandler failed for %s %s: %v", route.method, route.path, err)
			}
		}

		// Test requests
		requests := []struct {
			method           string
			path             string
			expectedResponse string
			expectedStatus   int
		}{
			{"GET", "/", "home", http.StatusOK},
			{"GET", "/users", "users list", http.StatusOK},
			{"GET", "/users/123", "user detail", http.StatusOK},
			{"POST", "/users", "create user", http.StatusOK},
			{"GET", "/users/123/posts", "user posts", http.StatusOK},
			{"GET", "/users/123/posts/456", "user post detail", http.StatusOK},
			{"GET", "/files/static/css/main.css", "file content", http.StatusOK},
			{"GET", "/api/v1/users", "api users", http.StatusOK},
		}

		for _, req := range requests {
			httpReq := httptest.NewRequest(req.method, req.path, nil)
			handlerFunc := tree.GetHandler(httpReq)
			recorder := httptest.NewRecorder()
			handlerFunc(recorder, httpReq)

			if recorder.Code != req.expectedStatus {
				t.Errorf("%s %s: Expected status %d, got %d", req.method, req.path, req.expectedStatus, recorder.Code)
			}
			if recorder.Body.String() != req.expectedResponse {
				t.Errorf("%s %s: Expected '%s', got '%s'", req.method, req.path, req.expectedResponse, recorder.Body.String())
			}
		}
	})
}

// BenchmarkGetHandler tests performance of GetHandler
func BenchmarkGetHandler(b *testing.B) {
	tree := Tree.NewTree()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Setup routes
	routes := []string{
		"/",
		"/users",
		"/users/:id",
		"/users/:id/posts",
		"/users/:id/posts/:postId",
		"/files/*path",
		"/api/v1/users",
		"/api/v1/users/:id",
	}

	for _, route := range routes {
		tree.SetHandler("GET", route, handler)
	}

	b.Run("static_route", func(b *testing.B) {
		req := httptest.NewRequest("GET", "/users", nil)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree.GetHandler(req)
		}
	})

	b.Run("wildcard_route", func(b *testing.B) {
		req := httptest.NewRequest("GET", "/users/123", nil)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree.GetHandler(req)
		}
	})

	b.Run("catch_all_route", func(b *testing.B) {
		req := httptest.NewRequest("GET", "/files/static/css/main.css", nil)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree.GetHandler(req)
		}
	})

	b.Run("complex_route", func(b *testing.B) {
		req := httptest.NewRequest("GET", "/users/123/posts/456", nil)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree.GetHandler(req)
		}
	})
}