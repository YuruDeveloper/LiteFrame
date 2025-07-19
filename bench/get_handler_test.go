package bench

import (
	"LiteFrame/Router/Tree"
	"net/http"
	"net/http/httptest"
	"testing"
)

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