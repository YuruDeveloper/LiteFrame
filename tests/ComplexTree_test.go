package tests

import (
	"LiteFrame/Router/Tree"
	"LiteFrame/Router/Tree/Component"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)


// TestComplexRouteTree 복잡한 경로들에 대한 트리 테스트
func TestComplexRouteTree(t *testing.T) {
	// 테스트용 핸들러 함수들
	homeHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("home")) }
	apiHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("api")) }
	userHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("user")) }
	userIDHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("userID")) }
	adminHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("admin")) }
	catchAllHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("catchAll")) }

	// 트리 생성
	tree := Tree.NewTree()

	// 복잡한 경로들 추가
	routes := []struct {
		Method  string
		Path    string
		Handler http.HandlerFunc
	}{
		{"GET", "/", homeHandler},
		{"GET", "/api", apiHandler},
		{"POST", "/api", apiHandler},
		{"GET", "/api/users", userHandler},
		{"POST", "/api/users", userHandler},
		{"GET", "/api/users/:id", userIDHandler},
		{"PUT", "/api/users/:id", userIDHandler},
		{"DELETE", "/api/users/:id", userIDHandler},
		{"GET", "/api/admin", adminHandler},
		{"GET", "/api/admin/users", userHandler},
		{"GET", "/static/*", catchAllHandler},
		{"GET", "/files/*", catchAllHandler},
	}

	// 경로들을 트리에 추가
	for _, route := range routes {
		err := tree.Add(route.Method, route.Path, route.Handler)
		assert.NoError(t, err, "경로 추가 중 오류 발생: %s %s", route.Method, route.Path)
	}

	// 예상되는 트리 구조 정의
	expectedTree := ExpectedTreeNode{
		Path: "/",
		Type: Component.RootType,
		Children: []ExpectedTreeNode{
			{
				Path:    "/",
				Type:    Component.StaticType,
				Methods: []string{"GET"},
			},
			{
				Path: "api",
				Type: Component.StaticType,
				Methods: []string{"GET", "POST"},
				Children: []ExpectedTreeNode{
					{
						Path:    "users",
						Type:    Component.StaticType,
						Methods: []string{"GET", "POST"},
						Children: []ExpectedTreeNode{
							{
								Path:    "id",
								Type:    Component.WildCardType,
								Methods: []string{"GET", "PUT", "DELETE"},
							},
						},
					},
					{
						Path: "admin",
						Type: Component.StaticType,
						Methods: []string{"GET"},
						Children: []ExpectedTreeNode{
							{
								Path:    "users",
								Type:    Component.StaticType,
								Methods: []string{"GET"},
							},
						},
					},
				},
			},
			{
				Path: "static",
				Type: Component.StaticType,
				Children: []ExpectedTreeNode{
					{
						Path:    "*",
						Type:    Component.CatchAllType,
						Methods: []string{"GET"},
					},
				},
			},
			{
				Path: "files",
				Type: Component.StaticType,
				Children: []ExpectedTreeNode{
					{
						Path:    "*",
						Type:    Component.CatchAllType,
						Methods: []string{"GET"},
					},
				},
			},
		},
	}

	// 실제 트리 시각화
	fmt.Println("\n=== 실제 트리 구조 ===")
	actualVisualizer := &TreeVisualizer{}
	actualTreeVisualization := actualVisualizer.VisualizeTree(&tree.Root, "", true)
	fmt.Println(actualTreeVisualization)

	// 예상 트리 시각화
	fmt.Println("\n=== 예상되는 트리 구조 ===")
	expectedVisualizer := &TreeVisualizer{}
	expectedTreeVisualization := expectedVisualizer.VisualizeExpectedTree(expectedTree, "", true)
	fmt.Println(expectedTreeVisualization)

	// 트리 구조 비교
	fmt.Println("\n=== 트리 구조 비교 ===")
	isMatched := CompareTreeStructure(t, &tree.Root, expectedTree, "")
	
	if isMatched {
		fmt.Println("✅ 트리 구조가 예상과 일치합니다!")
	} else {
		fmt.Println("❌ 트리 구조가 예상과 다릅니다!")
		t.Fail()
	}
}

// TestPathSplittingScenarios 경로 분할 시나리오 테스트
func TestPathSplittingScenarios(t *testing.T) {
	tree := Tree.NewTree()
	handler := func(w http.ResponseWriter, r *http.Request) {}

	// 공통 접두사를 가진 경로들 추가
	routes := []string{
		"/user",
		"/users",
		"/user-profile",
		"/user/settings",
		"/users/:id",
		"/users/:id/profile",
	}

	for _, route := range routes {
		err := tree.Add("GET", route, handler)
		assert.NoError(t, err, "경로 추가 중 오류 발생: %s", route)
	}

	// 트리 시각화
	fmt.Println("\n=== 경로 분할 시나리오 트리 구조 ===")
	visualizer := &TreeVisualizer{}
	treeVisualization := visualizer.VisualizeTree(&tree.Root, "", true)
	fmt.Println(treeVisualization)

	// 예상되는 구조: user라는 공통 접두사로 인한 노드 분할이 발생해야 함
	// user -> [user, users, user-profile] 등으로 분할
}

// TestWildcardAndCatchAllPriority 와일드카드와 Catch-All 우선순위 테스트
func TestWildcardAndCatchAllPriority(t *testing.T) {
	tree := Tree.NewTree()
	
	staticHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("static")) }
	wildcardHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("wildcard")) }
	catchAllHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("catchall")) }

	// 다양한 우선순위의 경로들 추가
	routes := []struct {
		Path    string
		Handler http.HandlerFunc
		Name    string
	}{
		{"/api/users/profile", staticHandler, "static"},
		{"/api/users/:id", wildcardHandler, "wildcard"},
		{"/api/*", catchAllHandler, "catchall"},
		{"/api/admin", staticHandler, "admin-static"},
	}

	for _, route := range routes {
		err := tree.Add("GET", route.Path, route.Handler)
		assert.NoError(t, err, "경로 추가 중 오류 발생: %s", route.Path)
	}

	// 트리 시각화
	fmt.Println("\n=== 우선순위 테스트 트리 구조 ===")
	visualizer := &TreeVisualizer{}
	treeVisualization := visualizer.VisualizeTree(&tree.Root, "", true)
	fmt.Println(treeVisualization)

	// 예상되는 우선순위: Static > Wildcard > CatchAll
	expectedTree := ExpectedTreeNode{
		Path: "/",
		Type: Component.RootType,
		Children: []ExpectedTreeNode{
			{
				Path: "api",
				Type: Component.StaticType,
				Children: []ExpectedTreeNode{
					{
						Path:    "*",
						Type:    Component.CatchAllType,
						Methods: []string{"GET"},
					},
					{
						Path:    "admin",
						Type:    Component.StaticType,
						Methods: []string{"GET"},
					},
					{
						Path: "users",
						Type: Component.StaticType,
						Children: []ExpectedTreeNode{
							{
								Path:    "id",
								Type:    Component.WildCardType,
								Methods: []string{"GET"},
							},
							{
								Path:    "profile",
								Type:    Component.StaticType,
								Methods: []string{"GET"},
							},
						},
					},
				},
			},
		},
	}

	// 구조 비교
	fmt.Println("\n=== 우선순위 트리 구조 비교 ===")
	isMatched := CompareTreeStructure(t, &tree.Root, expectedTree, "")
	
	if isMatched {
		fmt.Println("✅ 우선순위 트리 구조가 예상과 일치합니다!")
	} else {
		fmt.Println("❌ 우선순위 트리 구조가 예상과 다릅니다!")
	}
}

// TestDeepNestedRoutes 깊게 중첩된 경로들 테스트
func TestDeepNestedRoutes(t *testing.T) {
	tree := Tree.NewTree()
	handler := func(w http.ResponseWriter, r *http.Request) {}

	// 깊게 중첩된 경로들
	deepRoutes := []string{
		"/api/v1/users/:userId/posts/:postId/comments/:commentId",
		"/api/v1/users/:userId/posts/:postId/likes",
		"/api/v1/users/:userId/profile/settings/privacy",
		"/api/v1/admin/system/logs/access",
		"/api/v1/admin/system/logs/error",
		"/api/v2/users/:id/data/*",
	}

	for _, route := range deepRoutes {
		err := tree.Add("GET", route, handler)
		assert.NoError(t, err, "경로 추가 중 오류 발생: %s", route)
	}

	// 트리 시각화
	fmt.Println("\n=== 깊게 중첩된 경로 트리 구조 ===")
	visualizer := &TreeVisualizer{}
	treeVisualization := visualizer.VisualizeTree(&tree.Root, "", true)
	fmt.Println(treeVisualization)

	// 깊이 확인 (최소 6레벨 이상)
	assert.True(t, strings.Contains(treeVisualization, "commentId"), "깊게 중첩된 경로가 올바르게 생성되지 않았습니다")
}