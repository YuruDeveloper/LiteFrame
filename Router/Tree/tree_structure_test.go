package Tree

import (
	"testing"
)

// ======================
// Tree Structure Validation Tests
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
			err := tree.SetHandler(tree.StringToMethodType(route.Method), route.Path, route.Handler)
			AssertNoError(t, err, "SetHandler for "+route.Path)
		}

		// Validate root structure
		if tree.RootNode.Type != RootType {
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
			err := tree.SetHandler(tree.StringToMethodType(route.Method), route.Path, route.Handler)
			AssertNoError(t, err, "SetHandler for "+route.Path)
		}

		// Validate wildcard structure
		usersNode := findChildNode(tree.RootNode, "users")
		if usersNode == nil {
			t.Fatal("Users node not found")
		}

		if usersNode.WildCard == nil {
			t.Error("Wildcard child not found")
		}

		if usersNode.WildCard.Type != WildCardType {
			t.Error("Wildcard node type mismatch")
		}

		if usersNode.WildCard.Param != "id" {
			t.Errorf("Expected param 'id', got '%s'", usersNode.WildCard.Param)
		}
	})

	t.Run("catch_all_routes", func(t *testing.T) {
		tree := SetupTree()
		handler := CreateTestHandler()

		err := tree.SetHandler(tree.StringToMethodType("GET"), "/files/*path", handler)
		AssertNoError(t, err, "SetHandler")

		// Validate catch-all structure
		filesNode := findChildNode(tree.RootNode, "files")
		if filesNode == nil {
			t.Fatal("Files node not found")
		}

		if filesNode.CatchAll == nil {
			t.Error("CatchAll child not found")
		}

		if filesNode.CatchAll.Type != CatchAllType {
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
			{"GET", "/", handler},                    // Root
			{"GET", "/api/users", handler},           // Static
			{"GET", "/api/users/:id", handler},       // Wildcard
			{"GET", "/api/users/:id/posts", handler}, // Nested wildcard
			{"GET", "/static/*files", handler},       // Catch-all
		}

		for _, route := range routes {
			err := tree.SetHandler(tree.StringToMethodType(route.Method), route.Path, route.Handler)
			AssertNoError(t, err, "SetHandler for "+route.Path)
		}

		// Validate basic structure
		if tree.RootNode.Handlers[GET] == nil {
			t.Error("Root handler not set")
		}

		apiNode := findChildNode(tree.RootNode, "api")
		if apiNode == nil {
			t.Error("API node not found")
		}

		staticNode := findChildNode(tree.RootNode, "static")
		if staticNode == nil {
			t.Fatal("Static node not found")
		}

		if staticNode.CatchAll == nil {
			t.Error("Static catch-all not found")
		}
	})
}

// ======================
// Tree Consistency Validation Tests
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
		err := tree.SetHandler(tree.StringToMethodType(route.Method), route.Path, route.Handler)
		AssertNoError(t, err, "SetHandler")
	}

	t.Run("basic_consistency", func(t *testing.T) {
		if tree.RootNode == nil {
			t.Error("Root node is nil")
		}

		if tree.RootNode.Type != RootType {
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
// Helper Functions
// ======================

// findChildNode finds a node with the specified path among child nodes
func findChildNode(parent *Node, path string) *Node {
	for _, child := range parent.Children {
		if child.Path == path {
			return child
		}
	}
	return nil
}

// checkNodeConsistency recursively validates node consistency
func checkNodeConsistency(t *testing.T, node *Node, path string) {
	if node == nil {
		t.Errorf("Node at path %s is nil", path)
		return
	}

	// Validate child nodes
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

	// Validate wildcard node
	if node.WildCard != nil {
		wildcardPath := path + "/" + node.WildCard.Path
		if path == "/" {
			wildcardPath = "/" + node.WildCard.Path
		}
		checkNodeConsistency(t, node.WildCard, wildcardPath)
	}

	// Validate catch-all node
	if node.CatchAll != nil {
		catchAllPath := path + "/" + node.CatchAll.Path
		if path == "/" {
			catchAllPath = "/" + node.CatchAll.Path
		}
		checkNodeConsistency(t, node.CatchAll, catchAllPath)
	}
}

// ======================
// Tree Integrity Validation Tests
// ======================

func TestTreeIntegrity(t *testing.T) {
	t.Run("duplicate_route_handling", func(t *testing.T) {
		tree := SetupTree()
		handler1 := CreateHandlerWithResponse("handler1")
		handler2 := CreateHandlerWithResponse("handler2")

		// Set handler twice on same path
		err := tree.SetHandler(tree.StringToMethodType("GET"), "/test", handler1)
		AssertNoError(t, err, "First SetHandler")

		err = tree.SetHandler(tree.StringToMethodType("GET"), "/test", handler2)
		AssertNoError(t, err, "Second SetHandler") // Allow overwrite

		// Verify second handler is set
		recorder := ExecuteRequest(tree, "GET", "/test")
		AssertStatusCode(t, recorder, 200)
		AssertResponseBody(t, recorder, "handler2")
	})

	t.Run("path_conflicts", func(t *testing.T) {
		tree := SetupTree()
		handler := CreateTestHandler()

		// Potentially conflicting paths
		routes := []RouteConfig{
			{"GET", "/users/admin", handler},    // Static
			{"GET", "/users/:id", handler},      // Wildcard
			{"GET", "/users/:id/edit", handler}, // Wildcard + Static
		}

		for _, route := range routes {
			err := tree.SetHandler(tree.StringToMethodType(route.Method), route.Path, route.Handler)
			AssertNoError(t, err, "SetHandler for "+route.Path)
		}

		// Verify static route takes precedence
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
