package Tree

import (
	"testing"
)

// ======================
// Tree Constructor Tests
// ======================

func TestNewTree(t *testing.T) {
	tree := SetupTree()

	t.Run("root_node_type", func(t *testing.T) {
		if tree.RootNode.Type != RootType {
			t.Errorf("Expected root node type %d, got %d", RootType, tree.RootNode.Type)
		}
	})

	t.Run("root_path", func(t *testing.T) {
		if tree.RootNode.Path != "/" {
			t.Errorf("Expected root path '/', got '%s'", tree.RootNode.Path)
		}
	})

	t.Run("children_initialized", func(t *testing.T) {
		if tree.RootNode.Children == nil {
			t.Error("Expected children slice to be initialized")
		}
	})

	t.Run("handlers_initialized", func(t *testing.T) {
		if tree.RootNode.Handlers == nil {
			t.Error("Expected handlers map to be initialized")
		}
	})
}

// ======================
// PathWithSegment Tests
// ======================

func TestPathWithSegment(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{"root_path", "/", []string{}},
		{"empty_string", "", []string{}},
		{"single_segment", "/users", []string{"users"}},
		{"two_segments", "/users/123", []string{"users", "123"}},
		{"three_segments", "/users/123/posts", []string{"users", "123", "posts"}},
		{"no_leading_slash", "users/123", []string{"users", "123"}},
		{"trailing_slash", "/users/", []string{"users"}},
		{"multiple_slashes", "//users//123//", []string{"users", "123"}},
		{"wildcard_param", "/users/:id", []string{"users", ":id"}},
		{"catchall_param", "/files/*path", []string{"files", "*path"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pws := NewPathWithSegment(tc.input)
			var result []string

			for {
				pws.Next()
				if pws.IsSame() {
					break
				}
				segment := pws.Get()
				if segment != "" {
					result = append(result, segment)
				}
			}

			if len(result) != len(tc.expected) {
				t.Errorf("Expected length %d, got %d for path '%s'", len(tc.expected), len(result), tc.input)
				t.Errorf("Expected: %v, Got: %v", tc.expected, result)
				return
			}

			for i, segment := range result {
				if segment != tc.expected[i] {
					t.Errorf("Expected segment[%d] = '%s', got '%s'", i, tc.expected[i], segment)
				}
			}
		})
	}
}

// ======================
// 와일드카드 검증 테스트
// ======================

func TestIsWildCard(t *testing.T) {
	tree := SetupTree()

	testCases := []TestCase{
		{"valid_wildcard", ":id", true},
		{"valid_wildcard_with_text", ":user", true},
		{"empty_string", "", false},
		{"regular_string", "id", false},
		{"catch_all_character", "*", false},
		{"double_colon", "::", true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tree.IsWildCard(tc.Input)
			expected := tc.Expected.(bool)

			if result != expected {
				t.Errorf("IsWildCard(%q) = %v, expected %v", tc.Input, result, expected)
			}
		})
	}
}

// ======================
// 캐치올 검증 테스트
// ======================

func TestIsCatchAll(t *testing.T) {
	tree := SetupTree()

	testCases := []TestCase{
		{"valid_catch_all", "*", true},
		{"catch_all_with_text", "*files", true},
		{"empty_string", "", false},
		{"regular_string", "files", false},
		{"wildcard_character", ":", false},
		{"double_asterisk", "**", true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := tree.IsCatchAll(tc.Input)
			expected := tc.Expected.(bool)

			if result != expected {
				t.Errorf("IsCatchAll(%q) = %v, expected %v", tc.Input, result, expected)
			}
		})
	}
}

// ======================
// 문자열 매칭 테스트
// ======================

func TestMatch(t *testing.T) {
	tree := SetupTree()

	testCases := []MatchTestCase{
		{
			Name:          "exact_match",
			One:           "users",
			Two:           "users",
			ExpectedMatch: true,
			ExpectedIndex: 5,
			ExpectedLeft:  "",
		},
		{
			Name:          "first_shorter",
			One:           "user",
			Two:           "users",
			ExpectedMatch: true,
			ExpectedIndex: 4,
			ExpectedLeft:  "",
		},
		{
			Name:          "second_shorter",
			One:           "users",
			Two:           "user",
			ExpectedMatch: false,
			ExpectedIndex: 4,
			ExpectedLeft:  "s",
		},
		{
			Name:          "no_match",
			One:           "abc",
			Two:           "def",
			ExpectedMatch: false,
			ExpectedIndex: 0,
			ExpectedLeft:  "abc",
		},
		{
			Name:          "both_empty",
			One:           "",
			Two:           "",
			ExpectedMatch: true,
			ExpectedIndex: 0,
			ExpectedLeft:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// PathWithSegment 생성
			var pws *PathWithSegment
			if tc.One == "" {
				pws = NewPathWithSegment("")
			} else {
				pws = NewPathWithSegment("/" + tc.One)
				pws.Next() // 첫 번째 세그먼트로 이동
			}

			matched, index, leftPws := tree.Match(*pws, tc.Two)

			// Match 후 남은 부분을 Get()으로 가져오기
			left := leftPws.Get()

			if matched != tc.ExpectedMatch {
				t.Errorf("Expected match %v, got %v", tc.ExpectedMatch, matched)
			}
			if index != tc.ExpectedIndex {
				t.Errorf("Expected index %d, got %d", tc.ExpectedIndex, index)
			}
			if left != tc.ExpectedLeft {
				t.Errorf("Expected left '%s', got '%s'", tc.ExpectedLeft, left)
			}
		})
	}
}

// ======================
// 자식 노드 삽입 테스트
// ======================

func TestInsertChild(t *testing.T) {
	tree := SetupTree()

	t.Run("static_child", func(t *testing.T) {
		parent := NewNode(RootType, "/")
		child, err := tree.InsertChild(parent, "users")

		AssertNoError(t, err, "InsertChild")

		if child == nil {
			t.Error("Expected child node, got nil")
		}
		if child.Type != StaticType {
			t.Errorf("Expected static type %d, got %d", StaticType, child.Type)
		}
		if child.Path != "users" {
			t.Errorf("Expected path 'users', got '%s'", child.Path)
		}
	})

	t.Run("wildcard_child", func(t *testing.T) {
		parent := NewNode(RootType, "/")
		child, err := tree.InsertChild(parent, ":id")

		AssertNoError(t, err, "InsertChild wildcard")

		if child == nil {
			t.Error("Expected child node, got nil")
		}
		if child.Type != WildCardType {
			t.Errorf("Expected wildcard type %d, got %d", WildCardType, child.Type)
		}
		if child.Param != "id" {
			t.Errorf("Expected param 'id', got '%s'", child.Param)
		}
		if parent.WildCard == nil {
			t.Error("Expected parent WildCard to be set")
		}
	})

	t.Run("catch_all_child", func(t *testing.T) {
		parent := NewNode(RootType, "/")
		child, err := tree.InsertChild(parent, "*files")

		AssertNoError(t, err, "InsertChild catch-all")

		if child == nil {
			t.Error("Expected child node, got nil")
		}
		if child.Type != CatchAllType {
			t.Errorf("Expected catch-all type %d, got %d", CatchAllType, child.Type)
		}
		if parent.CatchAll == nil {
			t.Error("Expected parent CatchAll to be set")
		}
	})

	t.Run("duplicate_errors", func(t *testing.T) {
		// 중복 와일드카드 테스트
		parent := NewNode(RootType, "/")
		parent.WildCard = NewNode(WildCardType, ":existing")

		_, err := tree.InsertChild(parent, ":id")
		AssertError(t, err, "duplicate wildcard")

		// 중복 캐치올 테스트
		parent2 := NewNode(RootType, "/")
		parent2.CatchAll = NewNode(CatchAllType, "*existing")

		_, err = tree.InsertChild(parent2, "*files")
		AssertError(t, err, "duplicate catch-all")
	})
}

// ======================
// 핸들러 설정 테스트
// ======================

func TestSetHandler(t *testing.T) {
	tree := SetupTree()
	handler := CreateTestHandler()

	testCases := []struct {
		name   string
		method string
		path   string
		valid  bool
	}{
		{"root_handler", "GET", "/", true},
		{"simple_path", "GET", "/users", true},
		{"nested_path", "POST", "/users/123/posts", true},
		{"wildcard_path", "GET", "/users/:id", true},
		{"catch_all_path", "GET", "/files/*path", true},
		{"empty_method", "", "/test", false},
		{"empty_path", "GET", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var testHandler HandlerFunc
			if tc.valid {
				testHandler = handler
			}

			err := tree.SetHandler(tree.StringToMethodType(tc.method), tc.path, testHandler)

			if tc.valid {
				AssertNoError(t, err, "SetHandler")
			} else {
				AssertError(t, err, "SetHandler with invalid params")
			}
		})
	}
}

// ======================
// 노드 분할 테스트
// ======================

func TestSplitNode(t *testing.T) {
	tree := SetupTree()

	t.Run("basic_split", func(t *testing.T) {
		parent := NewNode(RootType, "/")
		child := NewNode(StaticType, "users")
		parent.Children = append(parent.Children, child)

		newNode, err := tree.SplitNode(parent, child, 4)
		AssertNoError(t, err, "SplitNode")

		if newNode == nil {
			t.Error("Expected new node, got nil")
		}
		if newNode.Path != "user" {
			t.Errorf("Expected new node path 'user', got '%s'", newNode.Path)
		}
		if child.Path != "s" {
			t.Errorf("Expected child path 's', got '%s'", child.Path)
		}
	})
}

// ======================
// 패턴 매칭 테스트 (TryMatch 메서드가 제거되어 주석 처리)
// ======================

// func TestTryMatch(t *testing.T) {
// 	tree := SetupTree()
// 	handler := CreateTestHandler()
//
// 	t.Run("no_children", func(t *testing.T) {
// 		parent := NewNode(RootType, "/")
// 		pws := NewPathWithSegment("/users")
// 		pws.Next() // 첫 번째 세그먼트로 이동
//
// 		matched, err := tree.TryMatch(parent, pws, GET, handler)
// 		AssertNoError(t, err, "TryMatch")
//
// 		if matched {
// 			t.Error("Expected no match, got match")
// 		}
// 	})
//
// 	t.Run("exact_match", func(t *testing.T) {
// 		parent := NewNode(RootType, "/")
// 		child := NewNode(StaticType, "users")
// 		parent.Children = append(parent.Children, child)
// 		pws := NewPathWithSegment("/users/123")
// 		pws.Next() // 첫 번째 세그먼트로 이동
//
// 		matched, err := tree.TryMatch(parent, pws, GET, handler)
// 		AssertNoError(t, err, "TryMatch")
//
// 		if !matched {
// 			t.Error("Expected match, got no match")
// 		}
// 	})
// }
