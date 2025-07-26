package tests

import (
	"LiteFrame/Router/Tree"
	"testing"
)

// ======================
// 트리 구조 검증 테스트
// ======================

func TestTreeStructure(t *testing.T) {
	t.Run("simple_static_routes", func(t *testing.T) {
		tree := SetupTree()
		handler := CreateTestHandler()

		routes := []RouteConfig{
			{"GET", "/users", handler},
			{"POST", "/users", handler},
			{"GET", "/posts", handler},
		}

		for _, route := range routes {
			err := tree.SetHandler(route.Method, route.Path, route.Handler)
			AssertNoError(t, err, "SetHandler for "+route.Path)
		}

		// 루트 구조 검증
		if tree.RootNode.Type != Tree.RootType {
			t.Error("Root node type mismatch")
		}

		if len(tree.RootNode.Children) != 2 {
			t.Errorf("Expected 2 children, got %d", len(tree.RootNode.Children))
		}
	})

	t.Run("wildcard_routes", func(t *testing.T) {
		tree := SetupTree()
		handler := CreateTestHandler()

		routes := []RouteConfig{
			{"GET", "/users/:id", handler},
			{"PUT", "/users/:id", handler},
			{"GET", "/users/:id/posts", handler},
		}

		for _, route := range routes {
			err := tree.SetHandler(route.Method, route.Path, route.Handler)
			AssertNoError(t, err, "SetHandler for "+route.Path)
		}

		// 와일드카드 구조 검증
		usersNode := findChildNode(tree.RootNode, "users")
		if usersNode == nil {
			t.Fatal("Users node not found")
		}

		if usersNode.WildCard == nil {
			t.Error("Wildcard child not found")
		}

		if usersNode.WildCard.Type != Tree.WildCardType {
			t.Error("Wildcard node type mismatch")
		}

		if usersNode.WildCard.Param != "id" {
			t.Errorf("Expected param 'id', got '%s'", usersNode.WildCard.Param)
		}
	})

	t.Run("catch_all_routes", func(t *testing.T) {
		tree := SetupTree()
		handler := CreateTestHandler()

		err := tree.SetHandler("GET", "/files/*path", handler)
		AssertNoError(t, err, "SetHandler")

		// 캐치올 구조 검증
		filesNode := findChildNode(tree.RootNode, "files")
		if filesNode == nil {
			t.Fatal("Files node not found")
		}

		if filesNode.CatchAll == nil {
			t.Error("CatchAll child not found")
		}

		if filesNode.CatchAll.Type != Tree.CatchAllType {
			t.Error("CatchAll node type mismatch")
		}

		if filesNode.CatchAll.Param != "path" {
			t.Errorf("Expected param 'path', got '%s'", filesNode.CatchAll.Param)
		}
	})

	t.Run("mixed_route_types", func(t *testing.T) {
		tree := SetupTree()
		handler := CreateTestHandler()

		routes := []RouteConfig{
			{"GET", "/", handler},                    // 루트
			{"GET", "/api/users", handler},           // 정적
			{"GET", "/api/users/:id", handler},       // 와일드카드
			{"GET", "/api/users/:id/posts", handler}, // 중첩 와일드카드
			{"GET", "/static/*files", handler},       // 캐치올
		}

		for _, route := range routes {
			err := tree.SetHandler(route.Method, route.Path, route.Handler)
			AssertNoError(t, err, "SetHandler for "+route.Path)
		}

		// 기본 구조 검증
		if tree.RootNode.Handlers[Tree.GET] == nil {
			t.Error("Root handler not set")
		}

		apiNode := findChildNode(tree.RootNode, "api")
		if apiNode == nil {
			t.Error("API node not found")
		}

		staticNode := findChildNode(tree.RootNode, "static")
		if staticNode == nil {
			t.Error("Static node not found")
		}

		if staticNode.CatchAll == nil {
			t.Error("Static catch-all not found")
		}
	})
}

// ======================
// 트리 일관성 검증 테스트
// ======================

func TestTreeConsistency(t *testing.T) {
	tree := SetupTree()
	handler := CreateTestHandler()

	routes := []RouteConfig{
		{"GET", "/api/users", handler},
		{"POST", "/api/users", handler},
		{"GET", "/api/users/:id", handler},
		{"GET", "/static/*files", handler},
	}

	for _, route := range routes {
		err := tree.SetHandler(route.Method, route.Path, route.Handler)
		AssertNoError(t, err, "SetHandler")
	}

	t.Run("basic_consistency", func(t *testing.T) {
		if tree.RootNode == nil {
			t.Error("Root node is nil")
		}

		if tree.RootNode.Type != Tree.RootType {
			t.Error("Root node type mismatch")
		}

		if tree.RootNode.Children == nil {
			t.Error("Root children not initialized")
		}
	})

	t.Run("node_consistency", func(t *testing.T) {
		checkNodeConsistency(t, tree.RootNode, "/")
	})
}

// ======================
// Helper 함수들
// ======================

// findChildNode는 자식 노드들 중에서 지정된 경로를 가진 노드를 찾습니다
func findChildNode(parent *Tree.Node, path string) *Tree.Node {
	for _, child := range parent.Children {
		if child.Path == path {
			return child
		}
	}
	return nil
}

// checkNodeConsistency는 노드의 일관성을 재귀적으로 확인합니다
func checkNodeConsistency(t *testing.T, node *Tree.Node, path string) {
	if node == nil {
		t.Errorf("Node at path %s is nil", path)
		return
	}

	// 자식 노드 검증
	for _, child := range node.Children {
		if child == nil {
			t.Errorf("Child node is nil at path %s", path)
			continue
		}

		childPath := path + "/" + child.Path
		if path == "/" {
			childPath = "/" + child.Path
		}

		checkNodeConsistency(t, child, childPath)
	}

	// 와일드카드 노드 검증
	if node.WildCard != nil {
		wildcardPath := path + "/" + node.WildCard.Path
		if path == "/" {
			wildcardPath = "/" + node.WildCard.Path
		}
		checkNodeConsistency(t, node.WildCard, wildcardPath)
	}

	// 캐치올 노드 검증
	if node.CatchAll != nil {
		catchAllPath := path + "/" + node.CatchAll.Path
		if path == "/" {
			catchAllPath = "/" + node.CatchAll.Path
		}
		checkNodeConsistency(t, node.CatchAll, catchAllPath)
	}
}

// ======================
// 트리 무결성 검증 테스트
// ======================

func TestTreeIntegrity(t *testing.T) {
	t.Run("duplicate_route_handling", func(t *testing.T) {
		tree := SetupTree()
		handler1 := CreateHandlerWithResponse("handler1")
		handler2 := CreateHandlerWithResponse("handler2")

		// 같은 경로에 핸들러 두 번 설정
		err := tree.SetHandler("GET", "/test", handler1)
		AssertNoError(t, err, "First SetHandler")

		err = tree.SetHandler("GET", "/test", handler2)
		AssertNoError(t, err, "Second SetHandler") // 덮어쓰기 허용

		// 두 번째 핸들러가 설정되었는지 확인
		recorder := ExecuteRequest(tree, "GET", "/test")
		AssertStatusCode(t, recorder, 200)
		AssertResponseBody(t, recorder, "handler2")
	})

	t.Run("path_conflicts", func(t *testing.T) {
		tree := SetupTree()
		handler := CreateTestHandler()

		// 충돌 가능한 경로들
		routes := []RouteConfig{
			{"GET", "/users/admin", handler},  // 정적
			{"GET", "/users/:id", handler},    // 와일드카드
			{"GET", "/users/:id/edit", handler}, // 와일드카드 + 정적
		}

		for _, route := range routes {
			err := tree.SetHandler(route.Method, route.Path, route.Handler)
			AssertNoError(t, err, "SetHandler for "+route.Path)
		}

		// 정적 라우트가 우선되는지 확인
		usersNode := findChildNode(tree.RootNode, "users")
		if usersNode == nil {
			t.Fatal("Users node not found")
		}

		adminNode := findChildNode(usersNode, "admin")
		if adminNode == nil {
			t.Error("Admin static route should exist")
		}

		if usersNode.WildCard == nil {
			t.Error("Wildcard route should also exist")
		}
	})
}