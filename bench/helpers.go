package bench

import (
	"LiteFrame/Router/Param"
	"LiteFrame/Router/Tree"
	"net/http"
	"net/http/httptest"
)

// ====================
// Benchmark Helpers
// ====================

// CreateBenchHandler creates a basic benchmark handler
func CreateBenchHandler() Tree.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
		w.WriteHeader(http.StatusOK)
	}
}

// CreateBenchHandlerWithParams creates a handler that processes parameters
func CreateBenchHandlerWithParams() Tree.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
		if params != nil {
			// Process fixed parameters
			for i := 0; i < params.Count && i < 2; i++ {
				_ = params.Path[params.Fix[i].Start:params.Fix[i].End]
			}
			// Process overflow parameters
			if params.Count > 2 {
				for i := 0; i < len(params.Overflow); i++ {
					_ = params.Path[params.Overflow[i].Start:params.Overflow[i].End]
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

// SetupBenchTree creates a benchmark tree
func SetupBenchTree() Tree.Tree {
	return Tree.NewTree()
}

// SetupBenchTreeWithRoutes sets up tree with routes
func SetupBenchTreeWithRoutes(routes []BenchRoute) Tree.Tree {
	tree := Tree.NewTree()
	for _, route := range routes {
		methodType := tree.StringToMethodType(route.Method)
		tree.SetHandler(methodType, route.Path, route.Handler)
	}
	return tree
}

// CreateBenchRequest creates a benchmark HTTP request
func CreateBenchRequest(method, path string) *http.Request {
	return httptest.NewRequest(method, path, nil)
}

// ====================
// Benchmark Types
// ====================

// BenchRoute represents a benchmark route
type BenchRoute struct {
	Method  string
	Path    string
	Handler Tree.HandlerFunc
}

// BenchRequest represents a benchmark request
type BenchRequest struct {
	Name   string
	Method string
	Path   string
}

// ====================
// Standard Datasets
// ====================

// GetStandardRoutes returns standard benchmark routes
func GetStandardRoutes() []BenchRoute {
	handler := CreateBenchHandler()
	paramHandler := CreateBenchHandlerWithParams()
	
	return []BenchRoute{
		{"GET", "/", handler},
		{"GET", "/users", handler},
		{"POST", "/users", handler},
		{"GET", "/users/:id", paramHandler},
		{"PUT", "/users/:id", paramHandler},
		{"DELETE", "/users/:id", paramHandler},
		{"GET", "/users/:id/posts", paramHandler},
		{"GET", "/users/:id/posts/:postId", paramHandler},
		{"GET", "/files/*path", paramHandler},
		{"GET", "/static/*filepath", paramHandler},
		{"GET", "/api/v1/users", handler},
		{"GET", "/api/v1/users/:id", paramHandler},
		{"GET", "/api/v2/users", handler},
	}
}

// GetStandardRequests returns standard benchmark requests
func GetStandardRequests() []BenchRequest {
	return []BenchRequest{
		{"Root", "GET", "/"},
		{"SimpleStatic", "GET", "/users"},
		{"SimpleWildcard", "GET", "/users/123"},
		{"NestedWildcard", "GET", "/users/123/posts"},
		{"ComplexWildcard", "GET", "/users/123/posts/456"},
		{"CatchAllShort", "GET", "/files/test.txt"},
		{"CatchAllLong", "GET", "/files/static/css/main.css"},
		{"ApiStatic", "GET", "/api/v1/users"},
		{"ApiWildcard", "GET", "/api/v1/users/123"},
	}
}

// GetComplexRoutes returns complex benchmark routes
func GetComplexRoutes() []BenchRoute {
	handler := CreateBenchHandler()
	paramHandler := CreateBenchHandlerWithParams()
	
	var routes []BenchRoute
	
	// Static routes
	staticPaths := []string{
		"/health", "/metrics", "/docs",
		"/admin/dashboard", "/admin/users",
		"/public/css", "/public/js",
	}
	for _, path := range staticPaths {
		routes = append(routes, BenchRoute{"GET", path, handler})
	}
	
	// API routes with parameters
	apiPaths := []string{
		"/api/v1/users/:id/profile",
		"/api/v1/projects/:projectId/tasks/:taskId",
		"/api/v1/organizations/:orgId/members/:memberId",
	}
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	for _, path := range apiPaths {
		for _, method := range methods {
			routes = append(routes, BenchRoute{method, path, paramHandler})
		}
	}
	
	// Catch-all routes
	catchAllPaths := []string{"/static/*filepath", "/assets/*path"}
	for _, path := range catchAllPaths {
		routes = append(routes, BenchRoute{"GET", path, paramHandler})
	}
	
	return routes
}