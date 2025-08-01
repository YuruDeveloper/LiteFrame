package Tree

import (
	"LiteFrame/Router/Param"
	"net/http"
	"strings"
	"testing"
)

// ======================
// Edge Cases and Boundary Tests
// ======================

func TestEdgeCases_Match(t *testing.T) {
	tree := SetupTree()

	testCases := []MatchTestCase{
		{
			Name:          "both_empty",
			One:           "",
			Two:           "",
			ExpectedMatch: true,
			ExpectedIndex: 0,
			ExpectedLeft:  "",
		},
		{
			Name:          "one_char_vs_empty",
			One:           "a",
			Two:           "",
			ExpectedMatch: false,
			ExpectedIndex: 0,
			ExpectedLeft:  "a",
		},
		{
			Name:          "empty_vs_one_char",
			One:           "",
			Two:           "a",
			ExpectedMatch: true,
			ExpectedIndex: 0,
			ExpectedLeft:  "",
		},
		{
			Name:          "long_vs_short",
			One:           "very_long_string",
			Two:           "v",
			ExpectedMatch: false,
			ExpectedIndex: 1,
			ExpectedLeft:  "ery_long_string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Panic prevention test
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Match(%q, %q) panicked: %v", tc.One, tc.Two, r)
				}
			}()

			// Create PathWithSegment
			var pws *PathWithSegment
			if tc.One == "" {
				pws = NewPathWithSegment("")
			} else {
				pws = NewPathWithSegment("/" + tc.One)
				pws.Next() // Move to first segment
			}

			matched, index, left := tree.Match(*pws, tc.Two)
			remaining := left.Body[left.Start:left.End]

			if matched != tc.ExpectedMatch {
				t.Errorf("Expected match %v, got %v", tc.ExpectedMatch, matched)
			}
			if index != tc.ExpectedIndex {
				t.Errorf("Expected index %d, got %d", tc.ExpectedIndex, index)
			}
			if remaining != tc.ExpectedLeft {
				t.Errorf("Expected remaining '%s', got '%s'", tc.ExpectedLeft, remaining)
			}
		})
	}
}

// ======================
// Path Splitting Edge Case Tests
// ======================

func TestEdgeCases_SplitPath(t *testing.T) {

	t.Run("malformed_paths", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    string
			expected []string
		}{
			{"only_slashes", "///", []string{}},
			{"excessive_slashes", "/////users/////", []string{"users"}},
			{"double_slashes", "//users//123//", []string{"users", "123"}},
			{"dot_segments", "/./users/../admin", []string{".", "users", "..", "admin"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Test by converting to PathWithSegment
				pws := NewPathWithSegment(tc.input)
				var result []string

				for {
					pws.Next()
					if pws.IsSame() {
						break
					}
					segment := pws.Body[pws.Start:pws.End]
					if segment != "" {
						result = append(result, segment)
					}
				}

				if len(result) != len(tc.expected) {
					t.Errorf("Expected length %d, got %d", len(tc.expected), len(result))
					return
				}

				for i, segment := range result {
					if segment != tc.expected[i] {
						t.Errorf("Expected segment[%d] = '%s', got '%s'", i, tc.expected[i], segment)
					}
				}
			})
		}
	})

	t.Run("special_characters", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    string
			expected []string
		}{
			{"url_encoded", "/users/user%20name", []string{"users", "user%20name"}},
			{"dot_in_name", "/files/file.txt", []string{"files", "file.txt"}},
			{"version_with_dot", "/api/v1.0", []string{"api", "v1.0"}},
			{"email_like", "/users/user@domain.com", []string{"users", "user@domain.com"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Test by converting to PathWithSegment
				pws := NewPathWithSegment(tc.input)
				var result []string

				for {
					pws.Next()
					if pws.IsSame() {
						break
					}
					segment := pws.Body[pws.Start:pws.End]
					if segment != "" {
						result = append(result, segment)
					}
				}

				if len(result) != len(tc.expected) {
					t.Errorf("Expected length %d, got %d", len(tc.expected), len(result))
					return
				}

				for i, segment := range result {
					if segment != tc.expected[i] {
						t.Errorf("Expected segment[%d] = '%s', got '%s'", i, tc.expected[i], segment)
					}
				}
			})
		}
	})
}

// ======================
// Wildcard/Catch-All Edge Case Tests
// ======================

func TestEdgeCases_WildcardDetection(t *testing.T) {
	tree := SetupTree()

	testCases := []struct {
		name       string
		input      string
		isWildcard bool
		isCatchAll bool
	}{
		{"empty_string", "", false, false},
		{"single_colon", ":", true, false},
		{"single_asterisk", "*", false, true},
		{"colon_asterisk", ":*", true, false},
		{"asterisk_colon", "*:", false, true},
		{"normal_text", "normal", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isWildcard := tree.IsWildCard(tc.input)
			isCatchAll := tree.IsCatchAll(tc.input)

			if isWildcard != tc.isWildcard {
				t.Errorf("IsWildCard(%q) = %v, expected %v", tc.input, isWildcard, tc.isWildcard)
			}

			if isCatchAll != tc.isCatchAll {
				t.Errorf("IsCatchAll(%q) = %v, expected %v", tc.input, isCatchAll, tc.isCatchAll)
			}
		})
	}
}

// ======================
// Child Node Insertion Edge Case Tests
// ======================

func TestEdgeCases_InsertChild(t *testing.T) {
	tree := SetupTree()

	t.Run("duplicate_wildcard_error", func(t *testing.T) {
		parent := NewNode(RootType, "/")

		// First wildcard successful
		child1, err1 := tree.InsertChild(parent, ":id")
		AssertNoError(t, err1, "First wildcard insertion")

		if child1 == nil {
			t.Error("First wildcard child is nil")
		}

		// Second wildcard should fail
		child2, err2 := tree.InsertChild(parent, ":name")
		AssertError(t, err2, "duplicate wildcard")

		if child2 != nil {
			t.Error("Expected nil child for duplicate wildcard")
		}
	})

	t.Run("duplicate_catchall_error", func(t *testing.T) {
		parent := NewNode(RootType, "/")

		// First catch-all successful
		child1, err1 := tree.InsertChild(parent, "*files")
		AssertNoError(t, err1, "First catch-all insertion")

		if child1 == nil {
			t.Error("First catch-all child is nil")
		}

		// Second catch-all should fail
		child2, err2 := tree.InsertChild(parent, "*documents")
		AssertError(t, err2, "duplicate catch-all")

		if child2 != nil {
			t.Error("Expected nil child for duplicate catch-all")
		}
	})

	t.Run("parameter_extraction", func(t *testing.T) {
		testCases := []struct {
			name          string
			path          string
			expectedParam string
		}{
			{"simple_param", ":id", "id"},
			{"underscore_param", ":user_id", "user_id"},
			{"camelCase_param", ":userId", "userId"},
			{"hyphen_param", ":user-id", "user-id"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				parent := NewNode(RootType, "/")
				child, err := tree.InsertChild(parent, tc.path)

				AssertNoError(t, err, "InsertChild")

				if child.Param != tc.expectedParam {
					t.Errorf("Expected param '%s', got '%s'", tc.expectedParam, child.Param)
				}
			})
		}
	})
}

// ======================
// Handler Setting Edge Case Tests
// ======================

func TestEdgeCases_SetHandler(t *testing.T) {
	t.Run("handler_overwrite", func(t *testing.T) {
		tree := SetupTree()
		handler1 := CreateHandlerWithResponse("response1")
		handler2 := CreateHandlerWithResponse("response2")

		// Set first handler
		err := tree.SetHandler(tree.StringToMethodType("GET"), "/test", handler1)
		AssertNoError(t, err, "First SetHandler")

		// Overwrite with second handler
		err = tree.SetHandler(tree.StringToMethodType("GET"), "/test", handler2)
		AssertNoError(t, err, "Second SetHandler")

		// Verify second handler works
		recorder := ExecuteRequest(tree, "GET", "/test")
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseBody(t, recorder, "response2")
	})

	t.Run("deep_path_handling", func(t *testing.T) {
		tree := SetupTree()
		handler := CreateHandlerWithResponse("deep response")

		// Create deep path
		segments := make([]string, 10)
		for i := range segments {
			segments[i] = "level" + string(rune('0'+(i%10)))
		}
		deepPath := "/" + strings.Join(segments, "/")

		err := tree.SetHandler(tree.StringToMethodType("GET"), deepPath, handler)
		AssertNoError(t, err, "SetHandler with deep path")

		recorder := ExecuteRequest(tree, "GET", deepPath)
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseBody(t, recorder, "deep response")
	})

	t.Run("mixed_route_types", func(t *testing.T) {
		tree := SetupTree()
		handler := func(w http.ResponseWriter, r *http.Request, params *Param.Params) {}

		routes := []struct {
			method string
			path   string
			name   string
		}{
			{"GET", "/users", "static"},
			{"GET", "/users/:id", "wildcard"},
			{"GET", "/files/*path", "catchall"},
			{"POST", "/users/:id/posts", "mixed"},
		}

		for _, route := range routes {
			t.Run(route.name, func(t *testing.T) {
				err := tree.SetHandler(tree.StringToMethodType(route.method), route.path, handler)
				AssertNoError(t, err, "SetHandler for "+route.path)
			})
		}
	})
}

// ======================
// Node Splitting Edge Case Tests
// ======================

func TestEdgeCases_SplitNode(t *testing.T) {
	tree := SetupTree()

	testCases := []struct {
		name          string
		originalPath  string
		splitPoint    int
		expectedLeft  string
		expectedRight string
	}{
		{"split_one_char", "users", 1, "u", "sers"},
		{"split_near_end", "users", 4, "user", "s"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parent := NewNode(RootType, "/")
			child := NewNode(StaticType, tc.originalPath)
			parent.Children = append(parent.Children, child)

			newNode, err := tree.SplitNode(parent, child, tc.splitPoint)
			AssertNoError(t, err, "SplitNode")

			if newNode == nil {
				t.Error("Expected new node, got nil")
				return
			}

			if newNode.Path != tc.expectedLeft {
				t.Errorf("Expected left path '%s', got '%s'", tc.expectedLeft, newNode.Path)
			}

			if child.Path != tc.expectedRight {
				t.Errorf("Expected right path '%s', got '%s'", tc.expectedRight, child.Path)
			}
		})
	}

	// Error case test - split_at_start should cause error as Left becomes empty string
	t.Run("split_at_start_error", func(t *testing.T) {
		parent := NewNode(RootType, "/")
		child := NewNode(StaticType, "users")
		parent.Children = append(parent.Children, child)

		_, err := tree.SplitNode(parent, child, 0)
		if err == nil {
			t.Error("Expected error for split at start (Left would be empty), got nil")
		}
	})
}

// ======================
// Input Validation Tests
// ======================

func TestInputValidation(t *testing.T) {
	tree := SetupTree()
	handler := CreateTestHandler()

	t.Run("invalid_method", func(t *testing.T) {
		err := tree.SetHandler(tree.StringToMethodType(""), "/test", handler)
		AssertError(t, err, "empty method")
	})

	t.Run("invalid_path", func(t *testing.T) {
		err := tree.SetHandler(tree.StringToMethodType("GET"), "", handler)
		AssertError(t, err, "empty path")
	})

	t.Run("nil_handler", func(t *testing.T) {
		err := tree.SetHandler(tree.StringToMethodType("GET"), "/test", nil)
		AssertError(t, err, "nil handler")
	})
}
