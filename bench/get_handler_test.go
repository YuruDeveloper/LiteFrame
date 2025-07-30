package bench

import (
	"fmt"
	"testing"
)

// ======================
// GetHandler 성능 벤치마크
// ======================

// BenchmarkGetHandler는 다양한 라우트 타입에서의 GetHandler 성능을 측정합니다
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

// ======================
// 라우트 타입별 GetHandler 벤치마크
// ======================

// BenchmarkGetHandler_RouteTypes는 라우트 타입별 성능을 자세히 측정합니다
func BenchmarkGetHandler_RouteTypes(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetStandardRoutes())

	b.Run("static_routes", func(b *testing.B) {
		requests := []BenchRequest{
			{"root", "GET", "/"},
			{"simple", "GET", "/users"},
			{"api", "GET", "/api/v1/users"},
		}

		for _, req := range requests {
			b.Run(req.Name, func(b *testing.B) {
				httpReq := CreateBenchRequest(req.Method, req.Path)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					tree.GetHandler(httpReq, tree.Pool.Get)
				}
			})
		}
	})

	b.Run("wildcard_routes", func(b *testing.B) {
		requests := []BenchRequest{
			{"simple_param", "GET", "/users/123"},
			{"nested_params", "GET", "/users/123/posts"},
			{"complex_params", "GET", "/users/123/posts/456"},
			{"api_params", "GET", "/api/v1/users/789"},
		}

		for _, req := range requests {
			b.Run(req.Name, func(b *testing.B) {
				httpReq := CreateBenchRequest(req.Method, req.Path)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					tree.GetHandler(httpReq, tree.Pool.Get)
				}
			})
		}
	})

	b.Run("catch_all_routes", func(b *testing.B) {
		requests := []BenchRequest{
			{"short_file", "GET", "/files/test.txt"},
			{"nested_file", "GET", "/files/dir/test.txt"},
			{"deep_file", "GET", "/files/static/css/components/button.css"},
			{"asset_file", "GET", "/static/js/app.min.js"},
		}

		for _, req := range requests {
			b.Run(req.Name, func(b *testing.B) {
				httpReq := CreateBenchRequest(req.Method, req.Path)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					tree.GetHandler(httpReq, tree.Pool.Get)
				}
			})
		}
	})
}

// ======================
// 경로 복잡도별 GetHandler 벤치마크
// ======================

// BenchmarkGetHandler_PathComplexity는 경로 복잡도에 따른 성능을 측정합니다
func BenchmarkGetHandler_PathComplexity(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetComplexRoutes())

	complexityTests := []struct {
		name string
		path string
	}{
		{"depth_1", "/health"},
		{"depth_3", "/api/v1/users"},
		{"depth_4", "/api/v1/users/123"},
		{"depth_6", "/api/v1/projects/456/tasks/789"},
	}

	for _, test := range complexityTests {
		b.Run(test.name, func(b *testing.B) {
			httpReq := CreateBenchRequest("GET", test.path)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.GetHandler(httpReq, tree.Pool.Get)
			}
		})
	}
}

// ======================
// HTTP 메서드별 GetHandler 벤치마크
// ======================

// BenchmarkGetHandler_HTTPMethods는 HTTP 메서드별 성능을 측정합니다
func BenchmarkGetHandler_HTTPMethods(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetComplexRoutes())
	path := "/api/v1/users/123"

	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		b.Run(method, func(b *testing.B) {
			httpReq := CreateBenchRequest(method, path)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.GetHandler(httpReq, tree.Pool.Get)
			}
		})
	}
}

// ======================
// 트리 크기별 GetHandler 벤치마크
// ======================

// BenchmarkGetHandler_TreeSize는 트리 크기에 따른 성능을 측정합니다
func BenchmarkGetHandler_TreeSize(b *testing.B) {
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
			httpReq := CreateBenchRequest("GET", testPath)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.GetHandler(httpReq, tree.Pool.Get)
			}
		})
	}
}

// ======================
// 동시성 GetHandler 벤치마크
// ======================

// BenchmarkGetHandler_Concurrent는 동시 요청 성능을 측정합니다
func BenchmarkGetHandler_Concurrent(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetStandardRoutes())
	requests := GetStandardRequests()

	b.Run("parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			reqIndex := 0
			for pb.Next() {
				req := requests[reqIndex%len(requests)]
				httpReq := CreateBenchRequest(req.Method, req.Path)
				tree.GetHandler(httpReq, tree.Pool.Get)
				reqIndex++
			}
		})
	})
}

// ======================
// 메모리 할당 벤치마크
// ======================

// BenchmarkGetHandler_Memory는 GetHandler의 메모리 할당을 측정합니다
func BenchmarkGetHandler_Memory(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetStandardRoutes())

	testCases := []struct {
		name string
		path string
	}{
		{"static", "/users"},
		{"wildcard", "/users/123"},
		{"complex", "/users/123/posts/456"},
		{"catch_all", "/files/static/css/main.css"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			httpReq := CreateBenchRequest("GET", tc.path)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.GetHandler(httpReq, tree.Pool.Get)
			}
		})
	}
}

// ======================
// Helper 함수들
// ======================

// generateDynamicPath는 동적으로 경로를 생성합니다
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

