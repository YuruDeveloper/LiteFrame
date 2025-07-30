package bench

import (
	"LiteFrame/Router/Param"
	"LiteFrame/Router/Tree"
	"net/http"
	"net/http/httptest"
)

// ====================
// 벤치마크 Helper 함수들
// ====================

// CreateBenchHandler는 벤치마크용 기본 핸들러를 생성합니다
func CreateBenchHandler() Tree.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
		w.WriteHeader(http.StatusOK)
	}
}

// CreateBenchHandlerWithParams는 매개변수를 처리하는 벤치마크 핸들러를 생성합니다
func CreateBenchHandlerWithParams() Tree.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
		// 매개변수 처리 시뮬레이션
		if params != nil && len(params.List) > 0 {
			for _, param := range params.List {
				_ = param.Value // 매개변수 값 읽기 시뮬레이션
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

// SetupBenchTree는 벤치마크용 트리를 생성하고 반환합니다
func SetupBenchTree() Tree.Tree {
	return Tree.NewTree()
}

// SetupBenchTreeWithRoutes는 미리 정의된 라우트들로 벤치마크 트리를 설정합니다
func SetupBenchTreeWithRoutes(routes []BenchRoute) Tree.Tree {
	tree := Tree.NewTree()
	
	for _, route := range routes {
		methodType := tree.StringToMethodType(route.Method)
		tree.SetHandler(methodType, route.Path, route.Handler)
	}
	
	return tree
}

// CreateBenchRequest는 벤치마크용 HTTP 요청을 생성합니다
func CreateBenchRequest(method, path string) *http.Request {
	return httptest.NewRequest(method, path, nil)
}

// ====================
// 벤치마크 데이터 구조체들
// ====================

// BenchRoute는 벤치마크용 라우트 설정 정보를 담습니다
type BenchRoute struct {
	Method  string
	Path    string
	Handler Tree.HandlerFunc
}

// BenchmarkConfig는 벤치마크 설정을 정의합니다
type BenchmarkConfig struct {
	Name        string
	Routes      []BenchRoute
	TestRequests []BenchRequest
}

// BenchRequest는 벤치마크 테스트 요청을 정의합니다
type BenchRequest struct {
	Name   string
	Method string
	Path   string
}

// ====================
// 공통 벤치마크 데이터셋
// ====================

// GetStandardRoutes는 표준 벤치마크 라우트 세트를 반환합니다
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
		{"POST", "/users/:id/posts", paramHandler},
		{"GET", "/users/:id/posts/:postId", paramHandler},
		{"GET", "/files/*path", paramHandler},
		{"GET", "/static/*filepath", paramHandler},
		{"GET", "/api/v1/users", handler},
		{"GET", "/api/v1/users/:id", paramHandler},
		{"GET", "/api/v1/users/:id/posts", paramHandler},
		{"GET", "/api/v2/users", handler},
	}
}

// GetStandardRequests는 표준 벤치마크 요청 세트를 반환합니다
func GetStandardRequests() []BenchRequest {
	return []BenchRequest{
		{"root", "GET", "/"},
		{"simple_static", "GET", "/users"},
		{"simple_wildcard", "GET", "/users/123"},
		{"nested_wildcard", "GET", "/users/123/posts"},
		{"complex_wildcard", "GET", "/users/123/posts/456"},
		{"catch_all_short", "GET", "/files/test.txt"},
		{"catch_all_long", "GET", "/files/static/css/main.css"},
		{"api_static", "GET", "/api/v1/users"},
		{"api_wildcard", "GET", "/api/v1/users/123"},
	}
}

// GetComplexRoutes는 복잡한 벤치마크 라우트 세트를 반환합니다
func GetComplexRoutes() []BenchRoute {
	handler := CreateBenchHandler()
	paramHandler := CreateBenchHandlerWithParams()
	
	routes := []BenchRoute{}
	
	// 다양한 정적 라우트
	staticPaths := []string{
		"/health",
		"/metrics",
		"/version",
		"/docs",
		"/swagger",
		"/admin/dashboard",
		"/admin/users",
		"/admin/settings",
		"/public/css",
		"/public/js",
		"/public/images",
	}
	
	for _, path := range staticPaths {
		routes = append(routes, BenchRoute{"GET", path, handler})
	}
	
	// API 라우트
	apiPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/logout",
		"/api/v1/auth/refresh",
		"/api/v1/users/:id/profile",
		"/api/v1/users/:id/settings",
		"/api/v1/projects/:projectId/tasks",
		"/api/v1/projects/:projectId/tasks/:taskId",
		"/api/v1/organizations/:orgId/members",
		"/api/v1/organizations/:orgId/members/:memberId",
	}
	
	for _, path := range apiPaths {
		routes = append(routes, BenchRoute{"GET", path, paramHandler})
		routes = append(routes, BenchRoute{"POST", path, paramHandler})
		routes = append(routes, BenchRoute{"PUT", path, paramHandler})
		routes = append(routes, BenchRoute{"DELETE", path, paramHandler})
	}
	
	// 캐치올 라우트
	catchAllPaths := []string{
		"/static/*filepath",
		"/assets/*path",
		"/uploads/*file",
		"/docs/*page",
	}
	
	for _, path := range catchAllPaths {
		routes = append(routes, BenchRoute{"GET", path, paramHandler})
	}
	
	return routes
}