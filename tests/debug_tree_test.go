package tests

import (
	"LiteFrame/Router/Tree"
	"fmt"
	"testing"
)

func TestDebugTreeStructure(t *testing.T) {
	tree := Tree.NewTree()
	handler := createTestHandler()

	// Test the problematic routes one by one
	routes := []string{
		"/api/v1/users",
		"/api/v2/users",
	}

	for i, route := range routes {
		fmt.Printf("\n=== Adding route %d: GET %s ===\n", i+1, route)
		err := tree.SetHandler("GET", route, handler)
		if err != nil {
			t.Fatalf("Failed to set handler: %v", err)
		}

		fmt.Println("Tree structure after adding route:")
		printTreeStructure(tree.RootNode, 0)

		// Debug path splitting
		paths := tree.SplitPath(route)
		fmt.Printf("Split path %s: %v\n", route, paths)
	}
}
