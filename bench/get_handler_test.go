package bench

import (
	"fmt"
	"strings"
	"testing"
)

// ======================
// GetHandler Benchmarks
// ======================

// BenchmarkGetHandler measures GetHandler performance by route type
func BenchmarkGetHandler(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetStandardRoutes())
	requests := GetStandardRequests()

	for _, req := range requests {
		b.Run(req.Name, func(b *testing.B) {
			httpReq := CreateBenchRequest(req.Method, req.Path)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.GetHandler(httpReq, tree.Pool.Get)
			}
		})
	}
}

// BenchmarkGetHandlerByType measures performance by route type
func BenchmarkGetHandlerByType(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetStandardRoutes())

	tests := map[string][]string{
		"Static":   {"/", "/users", "/api/v1/users"},
		"Wildcard": {"/users/123", "/users/123/posts", "/api/v1/users/789"},
		"CatchAll": {"/files/test.txt", "/static/js/app.min.js"},
	}

	for testType, paths := range tests {
		b.Run(testType, func(b *testing.B) {
			for _, path := range paths {
				b.Run(path[1:], func(b *testing.B) {
					req := CreateBenchRequest("GET", path)
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						tree.GetHandler(req, tree.Pool.Get)
					}
				})
			}
		})
	}
}

// BenchmarkGetHandlerDepth measures performance by path depth
func BenchmarkGetHandlerDepth(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetComplexRoutes())

	paths := []string{
		"/health",
		"/api/v1/users", 
		"/api/v1/users/123",
		"/api/v1/projects/456/tasks/789",
	}

	for _, path := range paths {
		b.Run(fmt.Sprintf("depth_%d", len(strings.Split(path, "/"))-1), func(b *testing.B) {
			req := CreateBenchRequest("GET", path)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.GetHandler(req, tree.Pool.Get)
			}
		})
	}
}

// BenchmarkGetHandlerMethods measures performance by HTTP method
func BenchmarkGetHandlerMethods(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetComplexRoutes())
	path := "/api/v1/users/123"
	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		b.Run(method, func(b *testing.B) {
			req := CreateBenchRequest(method, path)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.GetHandler(req, tree.Pool.Get)
			}
		})
	}
}

// BenchmarkGetHandlerTreeSize measures performance by tree size
func BenchmarkGetHandlerTreeSize(b *testing.B) {
	handler := CreateBenchHandlerWithParams()
	treeSizes := []int{50, 100, 200}

	for _, size := range treeSizes {
		b.Run(fmt.Sprintf("routes_%d", size), func(b *testing.B) {
			tree := SetupBenchTree()
			for i := 0; i < size; i++ {
				path := generateDynamicPath(i)
				methodType := tree.StringToMethodType("GET")
				tree.SetHandler(methodType, path, handler)
			}

			testPath := "/users/123/posts/456"
			methodType := tree.StringToMethodType("GET")
			tree.SetHandler(methodType, testPath, handler)
			req := CreateBenchRequest("GET", testPath)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.GetHandler(req, tree.Pool.Get)
			}
		})
	}
}

// BenchmarkGetHandlerConcurrent measures concurrent performance
func BenchmarkGetHandlerConcurrent(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetStandardRoutes())
	requests := GetStandardRequests()

	b.RunParallel(func(pb *testing.PB) {
		reqIndex := 0
		for pb.Next() {
			req := requests[reqIndex%len(requests)]
			httpReq := CreateBenchRequest(req.Method, req.Path)
			tree.GetHandler(httpReq, tree.Pool.Get)
			reqIndex++
		}
	})
}

// BenchmarkGetHandlerMemory measures memory allocation
func BenchmarkGetHandlerMemory(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetStandardRoutes())

	paths := map[string]string{
		"Static":   "/users",
		"Wildcard": "/users/123",
		"Complex":  "/users/123/posts/456",
		"CatchAll": "/files/main.css",
	}

	for name, path := range paths {
		b.Run(name, func(b *testing.B) {
			req := CreateBenchRequest("GET", path)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.GetHandler(req, tree.Pool.Get)
			}
		})
	}
}

// ======================
// Helper Functions
// ======================

// generateDynamicPath generates paths dynamically
func generateDynamicPath(index int) string {
	patterns := []string{
		"/route%d",
		"/api/v1/resource%d",
		"/service%d/endpoint",
		"/module%d/action/:id",
		"/system%d/component/*path",
	}
	
	pattern := patterns[index%len(patterns)]
	return fmt.Sprintf(pattern, index)
}

