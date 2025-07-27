package bench

import (
	"LiteFrame/Router/Tree"
	"testing"
)

// ======================
// Tree 핵심 기능 벤치마크
// ======================

// BenchmarkPathWithSegment는 새로운 경로 분할 성능을 측정합니다
func BenchmarkPathWithSegment(b *testing.B) {
	testCases := []struct {
		name string
		path string
	}{
		{"simple", "/users"},
		{"nested", "/users/123/posts"},
		{"deep", "/users/123/posts/456/comments/789"},
		{"complex", "/api/v1/organizations/123/projects/456/tasks/789/comments"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				path := Tree.NewPathWithSegment(tc.path)
				for !path.IsSame() {
					path.Next()
					_ = path.Get()
				}
			}
		})
	}
}

// BenchmarkMatch는 새로운 문자열 매칭 성능을 측정합니다
func BenchmarkMatch(b *testing.B) {
	tree := SetupBenchTree()
	
	testCases := []struct {
		name string
		one  string
		two  string
	}{
		{"exact_match", "users", "users"},
		{"partial_match", "user", "users"},
		{"no_match", "users", "posts"},
		{"long_strings", "very_long_string_for_testing", "very_long_string_for_comparison"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			sourcePath := Tree.NewPathWithSegment(tc.one)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.Match(*sourcePath, tc.two)
			}
		})
	}
}

// BenchmarkSetHandler는 핸들러 설정 성능을 측정합니다
func BenchmarkSetHandler(b *testing.B) {
	handler := CreateBenchHandler()
	
	testCases := []struct {
		name   string
		method string
		path   string
	}{
		{"static_route", "GET", "/users"},
		{"wildcard_route", "GET", "/users/:id"},
		{"catch_all_route", "GET", "/files/*path"},
		{"nested_route", "GET", "/api/v1/users/:id/posts"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := SetupBenchTree()
				tree.SetHandler(tc.method, tc.path, handler)
			}
		})
	}
}

// ======================
// Tree 유틸리티 함수 벤치마크
// ======================

// BenchmarkIsWildCard는 와일드카드 검증 성능을 측정합니다
func BenchmarkIsWildCard(b *testing.B) {
	tree := SetupBenchTree()
	
	testCases := []struct {
		name  string
		input string
	}{
		{"wildcard", ":id"},
		{"static", "users"},
		{"catch_all", "*path"},
		{"empty", ""},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.IsWildCard(tc.input)
			}
		})
	}
}

// BenchmarkIsCatchAll는 캐치올 검증 성능을 측정합니다
func BenchmarkIsCatchAll(b *testing.B) {
	tree := SetupBenchTree()
	
	testCases := []struct {
		name  string
		input string
	}{
		{"catch_all", "*path"},
		{"wildcard", ":id"},
		{"static", "users"},
		{"empty", ""},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.IsCatchAll(tc.input)
			}
		})
	}
}

// ======================
// 노드 조작 벤치마크
// ======================

// BenchmarkInsertChild는 자식 노드 삽입 성능을 측정합니다
func BenchmarkInsertChild(b *testing.B) {
	testCases := []struct {
		name string
		path string
	}{
		{"static", "users"},
		{"wildcard", ":id"},
		{"catch_all", "*path"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := SetupBenchTree()
				tree.InsertChild(tree.RootNode, tc.path)
			}
		})
	}
}

// ======================
// 통합 성능 벤치마크
// ======================

// BenchmarkTreeOperations는 트리 전체 작업의 통합 성능을 측정합니다
func BenchmarkTreeOperations(b *testing.B) {
	handler := CreateBenchHandler()
	routes := []string{
		"/",
		"/users",
		"/users/:id",
		"/users/:id/posts",
		"/files/*path",
	}

	b.Run("full_setup", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree := SetupBenchTree()
			for _, route := range routes {
				tree.SetHandler("GET", route, handler)
			}
		}
	})

	b.Run("path_processing", func(b *testing.B) {
		paths := []string{
			"/users/123/posts/456",
			"/files/static/css/main.css",
			"/api/v1/users/789",
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, pathStr := range paths {
				path := Tree.NewPathWithSegment(pathStr)
				for !path.IsSame() {
					path.Next()
					_ = path.Get()
				}
			}
		}
	})
}

// ======================
// 메모리 효율성 벤치마크
// ======================

// BenchmarkMemoryUsage는 메모리 사용량을 측정합니다
func BenchmarkMemoryUsage(b *testing.B) {
	handler := CreateBenchHandler()
	
	b.Run("single_tree", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree := SetupBenchTree()
			tree.SetHandler("GET", "/users/:id", handler)
		}
	})

	b.Run("multiple_routes", func(b *testing.B) {
		routes := GetStandardRoutes()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree := SetupBenchTree()
			for _, route := range routes[:5] { // 처음 5개만 사용
				tree.SetHandler(route.Method, route.Path, route.Handler)
			}
		}
	})
}