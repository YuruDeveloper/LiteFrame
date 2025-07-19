package bench

import (
	"LiteFrame/Router/Tree"
	"net/http"
	"testing"
)

// Helper function to create a test handler
func createTestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// Benchmark tests for performance
func BenchmarkSplitPath(b *testing.B) {
	tree := Tree.NewTree()
	path := "/users/123/posts/456/comments"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.SplitPath(path)
	}
}

func BenchmarkMatch(b *testing.B) {
	tree := Tree.NewTree()
	one := "users"
	two := "users"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Match(one, two)
	}
}

func BenchmarkSetHandler(b *testing.B) {
	tree := Tree.NewTree()
	handler := createTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.SetHandler("GET", "/users", handler)
	}
}