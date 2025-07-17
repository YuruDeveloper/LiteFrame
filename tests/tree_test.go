package tests

import (
	"LiteFrame/Router/Tree"
	"net/http"
	"testing"
)

// Test helper function to create a simple handler
func createTestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// TestNewTree tests the Tree constructor
func TestNewTree(t *testing.T) {
	tree := Tree.NewTree()

	if tree.RootNode.Type != Tree.RootType {
		t.Errorf("Expected root node type %d, got %d", Tree.RootType, tree.RootNode.Type)
	}

	if tree.RootNode.Path != "/" {
		t.Errorf("Expected root path '/', got '%s'", tree.RootNode.Path)
	}

	if tree.RootNode.Children == nil {
		t.Error("Expected children map to be initialized")
	}

	if tree.RootNode.Handlers == nil {
		t.Error("Expected handlers map to be initialized")
	}
}

// TestIsWildCard tests wildcard detection
func TestIsWildCard(t *testing.T) {
	tree := Tree.NewTree()

	tests := []struct {
		input    string
		expected bool
		name     string
	}{
		{":id", true, "valid wildcard"},
		{":user", true, "valid wildcard with text"},
		{"", false, "empty string"},
		{"id", false, "regular string"},
		{"*", false, "catch-all character"},
		{"::", true, "double colon"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := tree.IsWildCard(test.input)
			if result != test.expected {
				t.Errorf("IsWildCard(%q) = %v, expected %v", test.input, result, test.expected)
			}
		})
	}
}

// TestIsCatchAll tests catch-all detection
func TestIsCatchAll(t *testing.T) {
	tree := Tree.NewTree()

	tests := []struct {
		input    string
		expected bool
		name     string
	}{
		{"*", true, "valid catch-all"},
		{"*files", true, "catch-all with text"},
		{"", false, "empty string"},
		{"files", false, "regular string"},
		{":", false, "wildcard character"},
		{"**", true, "double asterisk"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := tree.IsCatchAll(test.input)
			if result != test.expected {
				t.Errorf("IsCatchAll(%q) = %v, expected %v", test.input, result, test.expected)
			}
		})
	}
}

// TestSplitPath tests path splitting functionality
func TestSplitPath(t *testing.T) {
	tree := Tree.NewTree()

	tests := []struct {
		input    string
		expected []string
		name     string
	}{
		{"/", []string{}, "root path"},
		{"", []string{}, "empty string"},
		{"/users", []string{"users"}, "single segment"},
		{"/users/123", []string{"users", "123"}, "two segments"},
		{"/users/123/posts", []string{"users", "123", "posts"}, "three segments"},
		{"users/123", []string{"users", "123"}, "no leading slash"},
		{"/users/", []string{"users"}, "trailing slash"},
		{"//users//123//", []string{"users", "123"}, "multiple slashes"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := tree.SplitPath(test.input)
			if len(result) != len(test.expected) {
				t.Errorf("SplitPath(%q) length = %d, expected %d", test.input, len(result), len(test.expected))
				return
			}

			for i, segment := range result {
				if segment != test.expected[i] {
					t.Errorf("SplitPath(%q)[%d] = %q, expected %q", test.input, i, segment, test.expected[i])
				}
			}
		})
	}
}

// TestMatch tests string matching functionality
func TestMatch(t *testing.T) {
	tree := Tree.NewTree()

	tests := []struct {
		one           string
		two           string
		expectedMatch bool
		expectedIndex int
		expectedLeft  string
		name          string
	}{
		{"users", "users", true, 5, "", "exact match"},
		{"user", "users", true, 4, "", "first shorter"},
		{"users", "user", false, 4, "s", "second shorter"},
		{"abc", "def", false, 0, "abc", "no match"},
		{"", "", true, 0, "", "both empty"},
		{"test", "", false, 0, "test", "second empty"},
		{"", "test", true, 0, "", "first empty"},
		{"hello", "help", false, 3, "lo", "partial match"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matched, index, left := tree.Match(test.one, test.two)

			if matched != test.expectedMatch {
				t.Errorf("Match(%q, %q) matched = %v, expected %v", test.one, test.two, matched, test.expectedMatch)
			}

			if index != test.expectedIndex {
				t.Errorf("Match(%q, %q) index = %d, expected %d", test.one, test.two, index, test.expectedIndex)
			}

			if left != test.expectedLeft {
				t.Errorf("Match(%q, %q) left = %q, expected %q", test.one, test.two, left, test.expectedLeft)
			}
		})
	}
}

// TestInsertHandler tests handler insertion
func TestInsertHandler(t *testing.T) {
	tree := Tree.NewTree()
	handler := createTestHandler()

	// Test valid methods
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT"}

	for _, method := range validMethods {
		t.Run("valid_method_"+method, func(t *testing.T) {
			node := Tree.NewNode(Tree.StaticType, "/test")
			err := tree.InsertHandler(node, method, handler)

			if err != nil {
				t.Errorf("InsertHandler() with method %s returned error: %v", method, err)
			}

			if node.Handlers[Tree.MethodList[method]] == nil {
				t.Errorf("Handler for method %s was not set", method)
			}
		})
	}

	// Test invalid method
	t.Run("invalid_method", func(t *testing.T) {
		node := Tree.NewNode(Tree.StaticType, "/test")
		err := tree.InsertHandler(node, "INVALID", handler)

		if err == nil {
			t.Error("Expected error for invalid method, got nil")
		}
	})
}

// TestInsertChild tests child node insertion
func TestInsertChild(t *testing.T) {
	tree := Tree.NewTree()

	t.Run("static_child", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")
		child, err := tree.InsertChild(parent, "users")

		if err != nil {
			t.Errorf("InsertChild() returned error: %v", err)
		}

		if child == nil {
			t.Error("Expected child node, got nil")
		}

		if child.Type != Tree.StaticType {
			t.Errorf("Expected static type %d, got %d", Tree.StaticType, child.Type)
		}

		if child.Path != "users" {
			t.Errorf("Expected path 'users', got '%s'", child.Path)
		}
	})

	t.Run("wildcard_child", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")
		child, err := tree.InsertChild(parent, ":id")

		if err != nil {
			t.Errorf("InsertChild() returned error: %v", err)
		}

		if child == nil {
			t.Error("Expected child node, got nil")
		}

		if child.Type != Tree.WildCardType {
			t.Errorf("Expected wildcard type %d, got %d", Tree.WildCardType, child.Type)
		}

		if child.Param != "id" {
			t.Errorf("Expected param 'id', got '%s'", child.Param)
		}

		if !parent.WildCard {
			t.Error("Expected parent WildCard flag to be true")
		}
	})

	t.Run("catch_all_child", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")
		child, err := tree.InsertChild(parent, "*files")

		if err != nil {
			t.Errorf("InsertChild() returned error: %v", err)
		}

		if child == nil {
			t.Error("Expected child node, got nil")
		}

		if child.Type != Tree.CatchAllType {
			t.Errorf("Expected catch-all type %d, got %d", Tree.CatchAllType, child.Type)
		}

		if !parent.CatchAll {
			t.Error("Expected parent CatchAll flag to be true")
		}
	})

	t.Run("duplicate_wildcard_error", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")
		parent.WildCard = true

		_, err := tree.InsertChild(parent, ":id")

		if err == nil {
			t.Error("Expected error for duplicate wildcard, got nil")
		}
	})

	t.Run("duplicate_catch_all_error", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")
		parent.CatchAll = true

		_, err := tree.InsertChild(parent, "*files")

		if err == nil {
			t.Error("Expected error for duplicate catch-all, got nil")
		}
	})
}

// TestSetHandler tests complete handler setting functionality
func TestSetHandler(t *testing.T) {
	tree := Tree.NewTree()
	handler := createTestHandler()

	t.Run("root_handler", func(t *testing.T) {
		err := tree.SetHandler("GET", "/", handler)
		if err != nil {
			t.Errorf("SetHandler() returned error: %v", err)
		}

		if tree.RootNode.Handlers[Tree.GET] == nil {
			t.Error("Root handler was not set")
		}
	})

	t.Run("simple_path", func(t *testing.T) {
		err := tree.SetHandler("GET", "/users", handler)
		if err != nil {
			t.Errorf("SetHandler() returned error: %v", err)
		}
	})

	t.Run("nested_path", func(t *testing.T) {
		err := tree.SetHandler("POST", "/users/123/posts", handler)
		if err != nil {
			t.Errorf("SetHandler() returned error: %v", err)
		}
	})

	t.Run("wildcard_path", func(t *testing.T) {
		err := tree.SetHandler("GET", "/users/:id", handler)
		if err != nil {
			t.Errorf("SetHandler() returned error: %v", err)
		}
	})

	t.Run("catch_all_path", func(t *testing.T) {
		err := tree.SetHandler("GET", "/files/*path", handler)
		if err != nil {
			t.Errorf("SetHandler() returned error: %v", err)
		}
	})

	t.Run("invalid_parameters", func(t *testing.T) {
		tests := []struct {
			method  string
			path    string
			handler http.HandlerFunc
			name    string
		}{
			{"", "/test", handler, "empty method"},
			{"GET", "", handler, "empty path"},
			{"GET", "/test", nil, "nil handler"},
		}

		for _, test := range tests {
			err := tree.SetHandler(test.method, test.path, test.handler)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", test.name)
			}
		}
	})
}

// TestSplitNode tests node splitting functionality
func TestSplitNode(t *testing.T) {
	tree := Tree.NewTree()

	t.Run("basic_split", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")
		child := Tree.NewNode(Tree.StaticType, "users")
		parent.Children = append(parent.Children, child)

		newNode, err := tree.SplitNode(parent, child, 4)

		if err != nil {
			t.Errorf("SplitNode() returned error: %v", err)
		}

		if newNode == nil {
			t.Error("Expected new node, got nil")
		}

		if newNode.Path != "user" {
			t.Errorf("Expected new node path 'user', got '%s'", newNode.Path)
		}

		if child.Path != "s" {
			t.Errorf("Expected child path 's', got '%s'", child.Path)
		}

		found := false
		for _, c := range parent.Children {
			if c.Path == "users" {
				found = true
			}
		}
		if found {
			t.Error("Original child should be removed from parent")
		}

		found = false
		for _, c := range parent.Children {
			if c.Path == "user" {
				found = true
			}
		}
		if !found {
			t.Error("New node should be added to parent")
		}
	})
}

// TestTryMatch tests the pattern matching logic
func TestTryMatch(t *testing.T) {
	tree := Tree.NewTree()
	handler := createTestHandler()

	t.Run("no_children", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")
		paths := []string{"users"}

		matched, err := tree.TryMatch(parent, paths, "GET", handler)

		if err != nil {
			t.Errorf("TryMatch() returned error: %v", err)
		}

		if matched {
			t.Error("Expected no match, got match")
		}
	})

	t.Run("exact_match", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")
		child := Tree.NewNode(Tree.StaticType, "users")
		parent.Children = append(parent.Children, child)
		paths := []string{"users", "123"}

		matched, err := tree.TryMatch(parent, paths, "GET", handler)

		if err != nil {
			t.Errorf("TryMatch() returned error: %v", err)
		}

		if !matched {
			t.Error("Expected match, got no match")
		}
	})
}

// Benchmark tests for performance
func BenchmarkSplitPath(b *testing.B) {
	tree := Tree.NewTree()
	path := "/users/123/posts/456/comments"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.SplitPath(path)
	}
}

func BenchmarkMatch(b *testing.B) {
	tree := Tree.NewTree()
	one := "users"
	two := "users"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Match(one, two)
	}
}

func BenchmarkSetHandler(b *testing.B) {
	tree := Tree.NewTree()
	handler := createTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.SetHandler("GET", "/users", handler)
	}
}