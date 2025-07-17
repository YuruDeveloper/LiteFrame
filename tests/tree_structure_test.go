package tests

import (
	"LiteFrame/Router/Tree"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// TreeStructure represents expected tree structure for testing
type TreeStructure struct {
	Path     string
	Type     Tree.NodeType
	Children map[string]*TreeStructure
	Methods  []string
	Param    string
}

// ExpectedRoute represents a route with its expected structure
type ExpectedRoute struct {
	Method      string
	Path        string
	Description string
}

// TestTreeStructureComparison tests actual vs expected tree structures
func TestTreeStructureComparison(t *testing.T) {
	tests := []struct {
		name     string
		routes   []ExpectedRoute
		expected *TreeStructure
	}{
		{
			name: "simple_static_routes",
			routes: []ExpectedRoute{
				{"GET", "/users", "Get all users"},
				{"POST", "/users", "Create user"},
				{"GET", "/posts", "Get all posts"},
			},
			expected: &TreeStructure{
				Path:    "/",
				Type:    Tree.RootType,
				Methods: []string{},
				Children: map[string]*TreeStructure{
					"users": {
						Path:     "users",
						Type:     Tree.StaticType,
						Methods:  []string{"GET", "POST"},
						Children: map[string]*TreeStructure{},
					},
					"posts": {
						Path:     "posts",
						Type:     Tree.StaticType,
						Methods:  []string{"GET"},
						Children: map[string]*TreeStructure{},
					},
				},
			},
		},
		{
			name: "nested_static_routes",
			routes: []ExpectedRoute{
				{"GET", "/api/v1/users", "Get users"},
				{"POST", "/api/v1/users", "Create user"},
				{"GET", "/api/v1/posts", "Get posts"},
				{"GET", "/api/v2/users", "Get users v2"},
			},
			expected: &TreeStructure{
				Path:    "/",
				Type:    Tree.RootType,
				Methods: []string{},
				Children: map[string]*TreeStructure{
					"api": {
						Path:    "api",
						Type:    Tree.StaticType,
						Methods: []string{},
						Children: map[string]*TreeStructure{
							"v": {
								Path:    "v",
								Type:    Tree.StaticType,
								Methods: []string{},
								Children: map[string]*TreeStructure{
									"1": {
										Path:    "1",
										Type:    Tree.StaticType,
										Methods: []string{},
										Children: map[string]*TreeStructure{
											"users": {
												Path:     "users",
												Type:     Tree.StaticType,
												Methods:  []string{"GET", "POST"},
												Children: map[string]*TreeStructure{},
											},
											"posts": {
												Path:     "posts",
												Type:     Tree.StaticType,
												Methods:  []string{"GET"},
												Children: map[string]*TreeStructure{},
											},
										},
									},
									"2": {
										Path:    "2",
										Type:    Tree.StaticType,
										Methods: []string{},
										Children: map[string]*TreeStructure{
											"users": {
												Path:     "users",
												Type:     Tree.StaticType,
												Methods:  []string{"GET"},
												Children: map[string]*TreeStructure{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "wildcard_routes",
			routes: []ExpectedRoute{
				{"GET", "/users/:id", "Get user by ID"},
				{"PUT", "/users/:id", "Update user"},
				{"DELETE", "/users/:id", "Delete user"},
				{"GET", "/users/:id/posts", "Get user posts"},
			},
			expected: &TreeStructure{
				Path:    "/",
				Type:    Tree.RootType,
				Methods: []string{},
				Children: map[string]*TreeStructure{
					"users": {
						Path:    "users",
						Type:    Tree.StaticType,
						Methods: []string{},
						Children: map[string]*TreeStructure{
							":id": {
								Path:    ":id",
								Type:    Tree.WildCardType,
								Param:   "id",
								Methods: []string{"DELETE", "GET", "PUT"},
								Children: map[string]*TreeStructure{
									"posts": {
										Path:     "posts",
										Type:     Tree.StaticType,
										Methods:  []string{"GET"},
										Children: map[string]*TreeStructure{},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "catch_all_routes",
			routes: []ExpectedRoute{
				{"GET", "/files/*path", "Serve static files"},
				{"GET", "/assets/*filepath", "Serve assets"},
			},
			expected: &TreeStructure{
				Path:    "/",
				Type:    Tree.RootType,
				Methods: []string{},
				Children: map[string]*TreeStructure{
					"files": {
						Path:    "files",
						Type:    Tree.StaticType,
						Methods: []string{},
						Children: map[string]*TreeStructure{
							"*path": {
								Path:     "*path",
								Type:     Tree.CatchAllType,
								Methods:  []string{"GET"},
								Children: map[string]*TreeStructure{},
							},
						},
					},
					"assets": {
						Path:    "assets",
						Type:    Tree.StaticType,
						Methods: []string{},
						Children: map[string]*TreeStructure{
							"*filepath": {
								Path:     "*filepath",
								Type:     Tree.CatchAllType,
								Methods:  []string{"GET"},
								Children: map[string]*TreeStructure{},
							},
						},
					},
				},
			},
		},
		{
			name: "complex_mixed_routes",
			routes: []ExpectedRoute{
				{"GET", "/", "Root handler"},
				{"GET", "/api/users", "Get all users"},
				{"POST", "/api/users", "Create user"},
				{"GET", "/api/users/:id", "Get user by ID"},
				{"PUT", "/api/users/:id", "Update user"},
				{"GET", "/api/users/:id/posts", "Get user posts"},
				{"POST", "/api/users/:id/posts", "Create user post"},
				{"GET", "/api/users/:id/posts/:postId", "Get specific post"},
				{"GET", "/static/*files", "Serve static files"},
			},
			expected: &TreeStructure{
				Path:    "/",
				Type:    Tree.RootType,
				Methods: []string{"GET"},
				Children: map[string]*TreeStructure{
					"api": {
						Path:    "api",
						Type:    Tree.StaticType,
						Methods: []string{},
						Children: map[string]*TreeStructure{
							"users": {
								Path:    "users",
								Type:    Tree.StaticType,
								Methods: []string{"GET", "POST"},
								Children: map[string]*TreeStructure{
									":id": {
										Path:    ":id",
										Type:    Tree.WildCardType,
										Param:   "id",
										Methods: []string{"GET", "PUT"},
										Children: map[string]*TreeStructure{
											"posts": {
												Path:    "posts",
												Type:    Tree.StaticType,
												Methods: []string{"GET", "POST"},
												Children: map[string]*TreeStructure{
													":postId": {
														Path:     ":postId",
														Type:     Tree.WildCardType,
														Param:    "postId",
														Methods:  []string{"GET"},
														Children: map[string]*TreeStructure{},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"static": {
						Path:    "static",
						Type:    Tree.StaticType,
						Methods: []string{},
						Children: map[string]*TreeStructure{
							"*files": {
								Path:     "*files",
								Type:     Tree.CatchAllType,
								Methods:  []string{"GET"},
								Children: map[string]*TreeStructure{},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create tree and add routes
			tree := Tree.NewTree()
			handler := createTestHandler()

			for _, route := range test.routes {
				err := tree.SetHandler(route.Method, route.Path, handler)
				if err != nil {
					t.Fatalf("Failed to set handler for %s %s: %v", route.Method, route.Path, err)
				}
			}

			// Print actual tree structure
			fmt.Printf("\n=== Test: %s ===\n", test.name)
			fmt.Println("Added routes:")
			for _, route := range test.routes {
				fmt.Printf("  %s %s - %s\n", route.Method, route.Path, route.Description)
			}

			fmt.Println("\nActual tree structure:")
			printTreeStructure(&tree.RootNode, 0)

			fmt.Println("\nExpected tree structure:")
			printExpectedStructure(test.expected, 0)

			// Compare structures
			actual := buildActualStructure(&tree.RootNode)
			if !compareStructures(actual, test.expected) {
				t.Errorf("Tree structure mismatch for test %s", test.name)
				fmt.Println("\nDetailed comparison:")
				compareAndPrintDifferences(actual, test.expected, "")
			} else {
				fmt.Println("✓ Tree structure matches expected!")
			}
		})
	}
}

// TestTreePathTraversal tests path traversal and outputs all possible paths
func TestTreePathTraversal(t *testing.T) {
	tree := Tree.NewTree()
	handler := createTestHandler()

	routes := []ExpectedRoute{
		{"GET", "/", "Root"},
		{"GET", "/api/users", "Get users"},
		{"POST", "/api/users", "Create user"},
		{"GET", "/api/users/:id", "Get user"},
		{"PUT", "/api/users/:id", "Update user"},
		{"GET", "/api/users/:id/posts", "Get user posts"},
		{"GET", "/api/users/:id/posts/:postId", "Get user post"},
		{"GET", "/static/*files", "Static files"},
		{"GET", "/docs", "Documentation"},
	}

	// Add all routes
	for _, route := range routes {
		err := tree.SetHandler(route.Method, route.Path, handler)
		if err != nil {
			t.Fatalf("Failed to set handler: %v", err)
		}
	}

	fmt.Println("\n=== Tree Path Traversal Test ===")
	fmt.Println("All registered routes:")
	for _, route := range routes {
		fmt.Printf("  %s %s\n", route.Method, route.Path)
	}

	fmt.Println("\nTree structure with all paths:")
	allPaths := extractAllPaths(&tree.RootNode, "")
	//sort.Strings(allPaths)

	fmt.Println("All possible paths in tree:")
	for _, path := range allPaths {
		fmt.Printf("  %s\n", path)
	}

	// Verify expected paths exist
	expectedPaths := []string{
		"/",
		"/api",
		"/api/users",
		"/api/users/:id",
		"/api/users/:id/posts",
		"/api/users/:id/posts/:postId",
		"/static",
		"/static/*files",
		"/docs",
	}

	fmt.Println("\nPath verification:")
	for _, expectedPath := range expectedPaths {
		found := false
		for _, actualPath := range allPaths {
			if actualPath == expectedPath {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("  ✓ %s\n", expectedPath)
		} else {
			fmt.Printf("  ✗ %s (MISSING)\n", expectedPath)
			t.Errorf("Expected path %s not found in tree", expectedPath)
		}
	}
}

// Helper functions

func printTreeStructure(node *Tree.Node, depth int) {
	indent := strings.Repeat("  ", depth)
	nodeType := getNodeTypeName(node.Type)

	methods := []string{}
	for method, handler := range node.Handlers {
		if handler != nil {
			methods = append(methods, string(method))
		}
	}
	//sort.Strings(methods)

	methodsStr := ""
	if len(methods) > 0 {
		methodsStr = fmt.Sprintf(" [%s]", strings.Join(methods, ","))
	}

	param := ""
	if node.Param != "" {
		param = fmt.Sprintf(" (param: %s)", node.Param)
	}

	fmt.Printf("%s%s (%s)%s%s\n", indent, node.Path, nodeType, methodsStr, param)

	// Sort children for consistent output
	childKeys := make([]string, 0, len(node.Children))
	for key := range node.Children {
		childKeys = append(childKeys, key)
	}
	//sort.Strings(childKeys)

	for _, key := range childKeys {
		printTreeStructure(node.Children[key], depth+1)
	}
}

func printExpectedStructure(structure *TreeStructure, depth int) {
	indent := strings.Repeat("  ", depth)
	nodeType := getNodeTypeName(structure.Type)

	methodsStr := ""
	if len(structure.Methods) > 0 {
		methods := make([]string, len(structure.Methods))
		copy(methods, structure.Methods)
		//sort.Strings(methods)
		methodsStr = fmt.Sprintf(" [%s]", strings.Join(methods, ","))
	}

	param := ""
	if structure.Param != "" {
		param = fmt.Sprintf(" (param: %s)", structure.Param)
	}

	fmt.Printf("%s%s (%s)%s%s\n", indent, structure.Path, nodeType, methodsStr, param)

	// Sort children for consistent output
	childKeys := make([]string, 0, len(structure.Children))
	for key := range structure.Children {
		childKeys = append(childKeys, key)
	}
	//sort.Strings(childKeys)

	for _, key := range childKeys {
		printExpectedStructure(structure.Children[key], depth+1)
	}
}

func getNodeTypeName(nodeType Tree.NodeType) string {
	switch nodeType {
	case Tree.RootType:
		return "ROOT"
	case Tree.StaticType:
		return "STATIC"
	case Tree.WildCardType:
		return "WILDCARD"
	case Tree.CatchAllType:
		return "CATCHALL"
	default:
		return "UNKNOWN"
	}
}

func buildActualStructure(node *Tree.Node) *TreeStructure {
	structure := &TreeStructure{
		Path:     node.Path,
		Type:     node.Type,
		Children: make(map[string]*TreeStructure),
		Methods:  []string{},
		Param:    node.Param,
	}

	// Collect methods
	for method, handler := range node.Handlers {
		if handler != nil {
			structure.Methods = append(structure.Methods, string(method))
		}
	}
	sort.Strings(structure.Methods)

	// Build children
	for key, child := range node.Children {
		structure.Children[key] = buildActualStructure(child)
	}

	return structure
}

func compareStructures(actual, expected *TreeStructure) bool {
	if actual.Path != expected.Path {
		return false
	}

	if actual.Type != expected.Type {
		return false
	}

	if actual.Param != expected.Param {
		return false
	}

	if !reflect.DeepEqual(actual.Methods, expected.Methods) {
		return false
	}

	if len(actual.Children) != len(expected.Children) {
		return false
	}

	for key, expectedChild := range expected.Children {
		actualChild, exists := actual.Children[key]
		if !exists {
			return false
		}

		if !compareStructures(actualChild, expectedChild) {
			return false
		}
	}

	return true
}

func compareAndPrintDifferences(actual, expected *TreeStructure, path string) {
	currentPath := path + "/" + actual.Path
	if path == "" {
		currentPath = actual.Path
	}

	if actual.Path != expected.Path {
		fmt.Printf("Path mismatch at %s: actual='%s', expected='%s'\n", currentPath, actual.Path, expected.Path)
	}

	if actual.Type != expected.Type {
		fmt.Printf("Type mismatch at %s: actual=%s, expected=%s\n", currentPath, getNodeTypeName(actual.Type), getNodeTypeName(expected.Type))
	}

	if actual.Param != expected.Param {
		fmt.Printf("Param mismatch at %s: actual='%s', expected='%s'\n", currentPath, actual.Param, expected.Param)
	}

	if !reflect.DeepEqual(actual.Methods, expected.Methods) {
		fmt.Printf("Methods mismatch at %s: actual=%v, expected=%v\n", currentPath, actual.Methods, expected.Methods)
	}

	// Check missing children
	for key := range expected.Children {
		if _, exists := actual.Children[key]; !exists {
			fmt.Printf("Missing child at %s: '%s'\n", currentPath, key)
		}
	}

	// Check extra children
	for key := range actual.Children {
		if _, exists := expected.Children[key]; !exists {
			fmt.Printf("Extra child at %s: '%s'\n", currentPath, key)
		}
	}

	// Recursively check existing children
	for key, expectedChild := range expected.Children {
		if actualChild, exists := actual.Children[key]; exists {
			compareAndPrintDifferences(actualChild, expectedChild, currentPath)
		}
	}
}

func extractAllPaths(node *Tree.Node, currentPath string) []string {
	paths := []string{}

	// Build current path
	fullPath := currentPath
	if node.Type == Tree.RootType {
		fullPath = "/"
	} else {
		if currentPath == "/" {
			fullPath = "/" + node.Path
		} else {
			fullPath = currentPath + "/" + node.Path
		}
	}

	// Add current path if it has handlers or we want to show intermediate paths
	hasHandlers := false
	for _, handler := range node.Handlers {
		if handler != nil {
			hasHandlers = true
			break
		}
	}

	if hasHandlers || node.Type == Tree.RootType || len(node.Children) > 0 {
		paths = append(paths, fullPath)
	}

	// Recursively get paths from children
	for _, child := range node.Children {
		childPaths := extractAllPaths(child, fullPath)
		paths = append(paths, childPaths...)
	}

	return paths
}

// TestTreeConsistency checks tree integrity and consistency
func TestTreeConsistency(t *testing.T) {
	tree := Tree.NewTree()
	handler := createTestHandler()

	routes := []ExpectedRoute{
		{"GET", "/api/users", "Get users"},
		{"POST", "/api/users", "Create user"},
		{"GET", "/api/users/:id", "Get user"},
		{"GET", "/api/users/:id/posts", "Get user posts"},
		{"GET", "/static/*files", "Static files"},
	}

	// Add routes
	for _, route := range routes {
		err := tree.SetHandler(route.Method, route.Path, handler)
		if err != nil {
			t.Fatalf("Failed to set handler: %v", err)
		}
	}

	fmt.Println("\n=== Tree Consistency Test ===")

	// Check tree consistency
	issues := checkTreeConsistency(&tree.RootNode, "/")

	if len(issues) == 0 {
		fmt.Println("✓ Tree structure is consistent!")
	} else {
		fmt.Println("✗ Tree consistency issues found:")
		for _, issue := range issues {
			fmt.Printf("  - %s\n", issue)
			t.Error(issue)
		}
	}
}

func checkTreeConsistency(node *Tree.Node, expectedPath string) []string {
	issues := []string{}

	// Check if node path matches expected
	if node.Type != Tree.RootType && node.Path == "" {
		issues = append(issues, fmt.Sprintf("Node at %s has empty path", expectedPath))
	}

	// Check wildcard constraints
	wildcardCount := 0
	catchAllCount := 0

	for _, child := range node.Children {
		if child.Type == Tree.WildCardType {
			wildcardCount++
		}
		if child.Type == Tree.CatchAllType {
			catchAllCount++
		}
	}

	if wildcardCount > 1 {
		issues = append(issues, fmt.Sprintf("Node at %s has multiple wildcard children", expectedPath))
	}

	if catchAllCount > 1 {
		issues = append(issues, fmt.Sprintf("Node at %s has multiple catch-all children", expectedPath))
	}

	if node.WildCard && wildcardCount == 0 {
		issues = append(issues, fmt.Sprintf("Node at %s has WildCard flag but no wildcard children", expectedPath))
	}

	if node.CatchAll && catchAllCount == 0 {
		issues = append(issues, fmt.Sprintf("Node at %s has CatchAll flag but no catch-all children", expectedPath))
	}

	// Recursively check children
	for key, child := range node.Children {
		childPath := expectedPath
		if expectedPath == "/" {
			childPath = "/" + key
		} else {
			childPath = expectedPath + "/" + key
		}

		childIssues := checkTreeConsistency(child, childPath)
		issues = append(issues, childIssues...)
	}

	return issues
}
