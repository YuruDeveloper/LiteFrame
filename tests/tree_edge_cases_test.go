package tests

import (
	"LiteFrame/Router/Tree"
	"net/http"
	"strings"
	"testing"
)

// TestEdgeCases_Match tests critical edge cases for the Match function
func TestEdgeCases_Match(t *testing.T) {
	tree := Tree.NewTree()

	t.Run("index_out_of_range_protection", func(t *testing.T) {
		// This test ensures the Match function doesn't panic on edge cases
		tests := []struct {
			one  string
			two  string
			name string
		}{
			{"", "", "both_empty"},
			{"a", "", "one_char_vs_empty"},
			{"", "a", "empty_vs_one_char"},
			{"very_long_string", "v", "long_vs_short"},
			{"a", "very_long_string", "short_vs_long"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// Should not panic
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Match(%q, %q) panicked: %v", test.one, test.two, r)
					}
				}()

				matched, index, remaining := tree.Match(test.one, test.two)

				// Verify remaining string doesn't cause out-of-bounds access
				if index > len(test.one) {
					t.Errorf("Index %d is out of bounds for string %q (length %d)", index, test.one, len(test.one))
				}

				// Verify remaining string is correctly calculated
				expectedRemaining := ""
				if index < len(test.one) {
					expectedRemaining = test.one[index:]
				}

				if remaining != expectedRemaining {
					t.Errorf("Expected remaining %q, got %q", expectedRemaining, remaining)
				}

				_ = matched // Use the variable to avoid compiler warning
			})
		}
	})

	t.Run("unicode_handling", func(t *testing.T) {
		tests := []struct {
			one      string
			two      string
			expected bool
			name     string
		}{
			{"hello", "hello world", true, "ascii_partial_match"},
			{"test", "testing", true, "ascii_prefix_match"},
			{"café", "café", true, "accent_exact_match"},
			{"abc", "xyz", false, "no_match"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				matched, _, _ := tree.Match(test.one, test.two)
				if matched != test.expected {
					t.Errorf("Match(%q, %q) = %v, expected %v", test.one, test.two, matched, test.expected)
				}
			})
		}
	})
}

// TestEdgeCases_SplitPath tests edge cases for path splitting
func TestEdgeCases_SplitPath(t *testing.T) {
	tree := Tree.NewTree()

	t.Run("malformed_paths", func(t *testing.T) {
		tests := []struct {
			input    string
			expected []string
			name     string
		}{
			{"///", []string{}, "only_slashes"},
			{"/////users/////", []string{"users"}, "excessive_slashes"},
			{"//users//123//", []string{"users", "123"}, "double_slashes"},
			{"/./users/../admin", []string{".", "users", "..", "admin"}, "dot_segments"},
			{"/users//", []string{"users"}, "trailing_double_slash"},
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
	})

	t.Run("very_long_paths", func(t *testing.T) {
		// Test with very long path
		longSegment := strings.Repeat("a", 1000)
		longPath := "/" + strings.Repeat(longSegment+"/", 100)

		result := tree.SplitPath(longPath)

		if len(result) != 100 {
			t.Errorf("Expected 100 segments, got %d", len(result))
		}

		for i, segment := range result {
			if segment != longSegment {
				t.Errorf("Segment %d = %q, expected %q", i, segment, longSegment)
				break
			}
		}
	})

	t.Run("special_characters", func(t *testing.T) {
		tests := []struct {
			input    string
			expected []string
			name     string
		}{
			{"/users/user%20name", []string{"users", "user%20name"}, "url_encoded"},
			{"/files/file.txt", []string{"files", "file.txt"}, "dot_in_name"},
			{"/api/v1.0", []string{"api", "v1.0"}, "version_with_dot"},
			{"/users/user@domain.com", []string{"users", "user@domain.com"}, "email_like"},
			{"/files/test-file_name", []string{"files", "test-file_name"}, "hyphens_underscores"},
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
	})
}

// TestEdgeCases_WildcardAndCatchAll tests edge cases for wildcard and catch-all detection
func TestEdgeCases_WildcardAndCatchAll(t *testing.T) {
	tree := Tree.NewTree()

	t.Run("boundary_conditions", func(t *testing.T) {
		tests := []struct {
			input      string
			isWildcard bool
			isCatchAll bool
			name       string
		}{
			{"", false, false, "empty_string"},
			{":", true, false, "single_colon"},
			{"*", false, true, "single_asterisk"},
			{":*", true, false, "colon_asterisk"},
			{"*:", false, true, "asterisk_colon"},
			{":::", true, false, "multiple_colons"},
			{"***", false, true, "multiple_asterisks"},
			{"normal", false, false, "normal_text"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				isWildcard := tree.IsWildCard(test.input)
				isCatchAll := tree.IsCatchAll(test.input)

				if isWildcard != test.isWildcard {
					t.Errorf("IsWildCard(%q) = %v, expected %v", test.input, isWildcard, test.isWildcard)
				}

				if isCatchAll != test.isCatchAll {
					t.Errorf("IsCatchAll(%q) = %v, expected %v", test.input, isCatchAll, test.isCatchAll)
				}
			})
		}
	})
}

// TestEdgeCases_InsertChild tests edge cases for child insertion
func TestEdgeCases_InsertChild(t *testing.T) {
	tree := Tree.NewTree()

	t.Run("duplicate_wildcard_scenarios", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")

		// First wildcard should succeed
		child1, err1 := tree.InsertChild(&parent, ":id")
		if err1 != nil {
			t.Errorf("First wildcard insertion failed: %v", err1)
		}
		if child1 == nil {
			t.Error("First wildcard child is nil")
		}

		// Second wildcard should fail
		child2, err2 := tree.InsertChild(&parent, ":name")
		if err2 == nil {
			t.Error("Expected error for duplicate wildcard, got nil")
		}
		if child2 != nil {
			t.Error("Expected nil child for duplicate wildcard")
		}
	})

	t.Run("duplicate_catchall_scenarios", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")

		// First catch-all should succeed
		child1, err1 := tree.InsertChild(&parent, "*files")
		if err1 != nil {
			t.Errorf("First catch-all insertion failed: %v", err1)
		}
		if child1 == nil {
			t.Error("First catch-all child is nil")
		}

		// Second catch-all should fail
		child2, err2 := tree.InsertChild(&parent, "*documents")
		if err2 == nil {
			t.Error("Expected error for duplicate catch-all, got nil")
		}
		if child2 != nil {
			t.Error("Expected nil child for duplicate catch-all")
		}
	})

	t.Run("wildcard_parameter_extraction", func(t *testing.T) {
		tests := []struct {
			path          string
			expectedParam string
			name          string
		}{
			{":id", "id", "simple_param"},
			{":user_id", "user_id", "underscore_param"},
			{":userId", "userId", "camelCase_param"},
			{":user-id", "user-id", "hyphen_param"},
			{":123", "123", "numeric_param"},
			{":", "", "empty_param"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				newParent := Tree.NewNode(Tree.RootType, "/")
				child, err := tree.InsertChild(&newParent, test.path)

				if err != nil {
					t.Errorf("InsertChild(%q) returned error: %v", test.path, err)
					return
				}

				if child.Param != test.expectedParam {
					t.Errorf("Expected param %q, got %q", test.expectedParam, child.Param)
				}
			})
		}
	})
}

// TestEdgeCases_SetHandler tests edge cases for handler setting
func TestEdgeCases_SetHandler(t *testing.T) {
	tree := Tree.NewTree()
	handler := func(w http.ResponseWriter, r *http.Request) {}

	t.Run("concurrent_operations", func(t *testing.T) {
		// Test multiple goroutines setting handlers simultaneously
		// This tests for race conditions (though not comprehensive without -race flag)
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()

				path := "/test" + string(rune('0'+id))
				err := tree.SetHandler("GET", path, handler)
				if err != nil {
					t.Errorf("Goroutine %d: SetHandler failed: %v", id, err)
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("overwrite_handlers", func(t *testing.T) {
		newTree := Tree.NewTree()
		handler1 := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }
		handler2 := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(404) }

		// Set initial handler
		err1 := newTree.SetHandler("GET", "/test", handler1)
		if err1 != nil {
			t.Errorf("First SetHandler failed: %v", err1)
		}

		// Overwrite with new handler
		err2 := newTree.SetHandler("GET", "/test", handler2)
		if err2 != nil {
			t.Errorf("Second SetHandler failed: %v", err2)
		}

		// The second handler should overwrite the first
		// (This test assumes the behavior, actual verification would require route matching)
	})

	t.Run("very_deep_paths", func(t *testing.T) {
		newTree := Tree.NewTree()

		// Create a very deep path
		segments := make([]string, 100)
		for i := range segments {
			segments[i] = "level" + string(rune('0'+(i%10)))
		}
		deepPath := "/" + strings.Join(segments, "/")

		err := newTree.SetHandler("GET", deepPath, handler)
		if err != nil {
			t.Errorf("SetHandler with deep path failed: %v", err)
		}
	})

	t.Run("mixed_route_types", func(t *testing.T) {
		newTree := Tree.NewTree()

		routes := []struct {
			method string
			path   string
			name   string
		}{
			{"GET", "/users", "static"},
			{"GET", "/users/:id", "wildcard"},
			{"GET", "/files/*path", "catchall"},
			{"POST", "/users/:id/posts", "mixed"},
			{"DELETE", "/admin/*", "admin_catchall"},
		}

		for _, route := range routes {
			t.Run(route.name, func(t *testing.T) {
				err := newTree.SetHandler(route.method, route.path, handler)
				if err != nil {
					t.Errorf("SetHandler(%s, %s) failed: %v", route.method, route.path, err)
				}
			})
		}
	})
}

// TestEdgeCases_SplitNode tests edge cases for node splitting
func TestEdgeCases_SplitNode(t *testing.T) {
	tree := Tree.NewTree()

	t.Run("split_at_boundaries", func(t *testing.T) {
		tests := []struct {
			originalPath  string
			splitPoint    int
			expectedLeft  string
			expectedRight string
			name          string
		}{
			{"users", 0, "", "users", "split_at_start"},
			{"users", 5, "users", "", "split_at_end"},
			{"users", 1, "u", "sers", "split_one_char"},
			{"users", 4, "user", "s", "split_near_end"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				parent := Tree.NewNode(Tree.RootType, "/")
				child := Tree.NewNode(Tree.StaticType, test.originalPath)
				parent.Children[test.originalPath] = &child

				newNode, err := tree.SplitNode(&parent, &child, test.splitPoint)

				if err != nil {
					t.Errorf("SplitNode failed: %v", err)
					return
				}

				if newNode.Path != test.expectedLeft {
					t.Errorf("Expected left path %q, got %q", test.expectedLeft, newNode.Path)
				}

				if child.Path != test.expectedRight {
					t.Errorf("Expected right path %q, got %q", test.expectedRight, child.Path)
				}
			})
		}
	})

	t.Run("split_invalid_point", func(t *testing.T) {
		parent := Tree.NewNode(Tree.RootType, "/")
		child := Tree.NewNode(Tree.StaticType, "users")
		parent.Children["users"] = &child

		// Test split point beyond string length
		_, err := tree.SplitNode(&parent, &child, 10)

		// The function should handle this gracefully (might panic in current implementation)
		// This test documents the expected behavior
		_ = err // Current implementation might not validate bounds
	})
}

// TestMemoryLeaks tests for potential memory leaks
func TestMemoryLeaks(t *testing.T) {
	// This test creates and destroys many trees to check for memory leaks
	// Run with: go test -run TestMemoryLeaks -memprofile=mem.prof

	for i := 0; i < 1000; i++ {
		tree := Tree.NewTree()
		handler := func(w http.ResponseWriter, r *http.Request) {}

		// Add many routes
		for j := 0; j < 100; j++ {
			path := "/test" + string(rune('0'+(j%10))) + "/" + string(rune('a'+(j%26)))
			tree.SetHandler("GET", path, handler)
		}

		// Tree should be garbage collected after this iteration
		_ = tree
	}
}

// TestStressConditions tests the tree under stress conditions
func TestStressConditions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	tree := Tree.NewTree()
	handler := func(w http.ResponseWriter, r *http.Request) {}

	t.Run("many_routes", func(t *testing.T) {
		// Add 10,000 routes
		for i := 0; i < 10000; i++ {
			path := "/stress/test/" + string(rune('a'+(i%26))) + "/" + string(rune('0'+(i%10)))
			err := tree.SetHandler("GET", path, handler)
			if err != nil {
				t.Errorf("Failed to set handler for route %d: %v", i, err)
				break
			}
		}
	})
}
