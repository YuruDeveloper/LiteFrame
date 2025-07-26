package bench

import (
	"fmt"
	"testing"
)

// ======================
// 비교 성능 벤치마크
// ======================

// BenchmarkComparative는 다양한 시나리오의 성능을 비교합니다
func BenchmarkComparative(b *testing.B) {
	// 표준 라우트 세트로 트리 설정
	standardTree := SetupBenchTreeWithRoutes(GetStandardRoutes())
	
	// 복잡한 라우트 세트로 트리 설정  
	complexTree := SetupBenchTreeWithRoutes(GetComplexRoutes())

	b.Run("standard_vs_complex", func(b *testing.B) {
		testPath := "/api/v1/users/123"

		b.Run("standard_tree", func(b *testing.B) {
			httpReq := CreateBenchRequest("GET", testPath)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				standardTree.GetHandler(httpReq)
			}
		})

		b.Run("complex_tree", func(b *testing.B) {
			httpReq := CreateBenchRequest("GET", testPath)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				complexTree.GetHandler(httpReq)
			}
		})
	})

	b.Run("route_type_comparison", func(b *testing.B) {
		tree := SetupBenchTreeWithRoutes(GetStandardRoutes())

		routeTests := []struct {
			name string
			path string
		}{
			{"static", "/users"},
			{"single_param", "/users/123"},
			{"double_param", "/users/123/posts/456"},
			{"catch_all", "/files/static/css/main.css"},
		}

		for _, rt := range routeTests {
			b.Run(rt.name, func(b *testing.B) {
				httpReq := CreateBenchRequest("GET", rt.path)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					tree.GetHandler(httpReq)
				}
			})
		}
	})

	b.Run("depth_comparison", func(b *testing.B) {
		tree := SetupBenchTreeWithRoutes(GetComplexRoutes())

		depthTests := []struct {
			name string
			path string
		}{
			{"depth_1", "/health"},
			{"depth_3", "/api/v1/users"},
			{"depth_4", "/api/v1/users/123"},
			{"depth_6", "/api/v1/projects/456/tasks/789"},
		}

		for _, dt := range depthTests {
			b.Run(dt.name, func(b *testing.B) {
				httpReq := CreateBenchRequest("GET", dt.path)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					tree.GetHandler(httpReq)
				}
			})
		}
	})
}

// ======================
// 메모리 효율성 비교 벤치마크
// ======================

// BenchmarkMemoryComparison은 메모리 사용량을 비교합니다
func BenchmarkMemoryComparison(b *testing.B) {

	b.Run("tree_setup_memory", func(b *testing.B) {
		routes := GetStandardRoutes()

		b.Run("small_tree", func(b *testing.B) {
			smallRoutes := routes[:5]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := SetupBenchTree()
				for _, route := range smallRoutes {
					tree.SetHandler(route.Method, route.Path, route.Handler)
				}
			}
		})

		b.Run("medium_tree", func(b *testing.B) {
			mediumRoutes := routes[:10]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := SetupBenchTree()
				for _, route := range mediumRoutes {
					tree.SetHandler(route.Method, route.Path, route.Handler)
				}
			}
		})

		b.Run("large_tree", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := SetupBenchTree()
				for _, route := range routes {
					tree.SetHandler(route.Method, route.Path, route.Handler)
				}
			}
		})
	})

	b.Run("handler_call_memory", func(b *testing.B) {
		tree := SetupBenchTreeWithRoutes(GetStandardRoutes())

		testCases := []struct {
			name string
			path string
		}{
			{"static_no_params", "/users"},
			{"single_param", "/users/123"},
			{"multiple_params", "/users/123/posts/456"},
			{"catch_all", "/files/static/css/main.css"},
		}

		for _, tc := range testCases {
			b.Run(tc.name, func(b *testing.B) {
				httpReq := CreateBenchRequest("GET", tc.path)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					handlerFunc := tree.GetHandler(httpReq)
					if handlerFunc != nil {
						// 실제 핸들러 호출 시뮬레이션 (메모리 할당 측정)
						_ = handlerFunc
					}
				}
			})
		}
	})
}

// ======================
// 확장성 벤치마크
// ======================

// BenchmarkScalability는 확장성을 테스트합니다
func BenchmarkScalability(b *testing.B) {
	handler := CreateBenchHandlerWithParams()

	scaleSizes := []int{10, 50, 100, 500, 1000}

	for _, size := range scaleSizes {
		b.Run(fmt.Sprintf("routes_%d", size), func(b *testing.B) {
			// 트리 설정
			tree := SetupBenchTree()
			for i := 0; i < size; i++ {
				path := generateDynamicPath(i)
				tree.SetHandler("GET", path, handler)
			}

			// 테스트용 경로 추가
			testPath := "/benchmark/test/path/123"
			tree.SetHandler("GET", testPath, handler)

			// 벤치마크 실행
			httpReq := CreateBenchRequest("GET", testPath)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.GetHandler(httpReq)
			}
		})
	}
}

// ======================
// 실제 사용 시나리오 벤치마크  
// ======================

// BenchmarkRealWorldScenarios는 실제 사용 시나리오를 시뮬레이션합니다
func BenchmarkRealWorldScenarios(b *testing.B) {
	tree := SetupBenchTreeWithRoutes(GetComplexRoutes())

	b.Run("api_server_simulation", func(b *testing.B) {
		// API 서버에서 자주 호출되는 엔드포인트들
		apiRequests := []BenchRequest{
			{"health_check", "GET", "/health"},
			{"user_profile", "GET", "/api/v1/users/123/profile"},
			{"user_posts", "GET", "/api/v1/users/123/posts"},
			{"project_tasks", "GET", "/api/v1/projects/456/tasks"},
			{"static_assets", "GET", "/static/css/main.css"},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 각 요청을 순환적으로 실행
			req := apiRequests[i%len(apiRequests)]
			httpReq := CreateBenchRequest(req.Method, req.Path)
			tree.GetHandler(httpReq)
		}
	})

	b.Run("mixed_workload", func(b *testing.B) {
		// 다양한 타입의 요청이 섞인 워크로드
		requests := []string{
			"/api/v1/users",                    // 정적
			"/api/v1/users/123",                // 와일드카드
			"/api/v1/users/123/posts/456",      // 복합 와일드카드
			"/static/js/app.min.js",            // 캐치올
			"/docs/api/reference",              // 정적
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			path := requests[i%len(requests)]
			httpReq := CreateBenchRequest("GET", path)
			tree.GetHandler(httpReq)
		}
	})
}