package bench

import (
	"LiteFrame/Router/Tree"
	"testing"
)

// ======================
// Core Function Benchmarks
// ======================

// BenchmarkPathWithSegment measures path parsing performance
func BenchmarkPathWithSegment(b *testing.B) {
	paths := []string{
		"/users",
		"/users/123/posts",
		"/api/v1/users/123/posts/456",
	}

	for _, path := range paths {
		b.Run(path[1:], func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				pws := Tree.NewPathWithSegment(path)
				for !pws.IsSame() {
					pws.Next()
					_ = pws.GetLength()
				}
			}
		})
	}
}

// BenchmarkMatch measures string matching performance
func BenchmarkMatch(b *testing.B) {
	tree := SetupBenchTree()
	
	tests := []struct{
		one, two string
	}{
		{"users", "users"},
		{"user", "users"},
		{"users", "posts"},
	}

	for _, test := range tests {
		b.Run(test.one+"_vs_"+test.two, func(b *testing.B) {
			pws := Tree.NewPathWithSegment("/" + test.one)
			pws.Next()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.Match(*pws, test.two)
			}
		})
	}
}

// BenchmarkSetHandler measures handler registration performance
func BenchmarkSetHandler(b *testing.B) {
	handler := CreateBenchHandler()
	routes := []string{
		"/users",
		"/users/:id", 
		"/files/*path",
		"/api/v1/users/:id/posts",
	}

	for _, route := range routes {
		b.Run(route[1:], func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := SetupBenchTree()
				methodType := tree.StringToMethodType("GET")
				tree.SetHandler(methodType, route, handler)
			}
		})
	}
}

// ======================
// Utility Function Benchmarks  
// ======================

// BenchmarkRouteTypeValidation measures route type detection performance
func BenchmarkRouteTypeValidation(b *testing.B) {
	tree := SetupBenchTree()
	inputs := []string{":id", "users", "*path", ""}

	b.Run("IsWildCard", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree.IsWildCard(inputs[i%len(inputs)])
		}
	})

	b.Run("IsCatchAll", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree.IsCatchAll(inputs[i%len(inputs)])
		}
	})
}

// BenchmarkInsertChild measures child node insertion performance
func BenchmarkInsertChild(b *testing.B) {
	paths := []string{"users", ":id", "*path"}

	for _, path := range paths {
		b.Run(path, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := SetupBenchTree()
				tree.InsertChild(tree.RootNode, path)
			}
		})
	}
}

// ======================
// Integrated Benchmarks
// ======================

// BenchmarkTreeOperations measures integrated tree performance
func BenchmarkTreeOperations(b *testing.B) {
	handler := CreateBenchHandler()
	routes := []string{"/", "/users", "/users/:id", "/files/*path"}

	b.Run("TreeSetup", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree := SetupBenchTree()
			for _, route := range routes {
				methodType := tree.StringToMethodType("GET")
				tree.SetHandler(methodType, route, handler)
			}
		}
	})

	b.Run("PathProcessing", func(b *testing.B) {
		paths := []string{"/users/123/posts", "/files/main.css"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, pathStr := range paths {
				pws := Tree.NewPathWithSegment(pathStr)
				for !pws.IsSame() {
					pws.Next()
					_ = pws.GetLength()
				}
			}
		}
	})
}

// BenchmarkMemoryUsage measures memory allocation patterns
func BenchmarkMemoryUsage(b *testing.B) {
	handler := CreateBenchHandler()
	
	b.Run("SingleRoute", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree := SetupBenchTree()
			methodType := tree.StringToMethodType("GET")
			tree.SetHandler(methodType, "/users/:id", handler)
		}
	})

	b.Run("MultipleRoutes", func(b *testing.B) {
		routes := GetStandardRoutes()[:5]
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree := SetupBenchTree()
			for _, route := range routes {
				methodType := tree.StringToMethodType(route.Method)
				tree.SetHandler(methodType, route.Path, route.Handler)
			}
		}
	})
}