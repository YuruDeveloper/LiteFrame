package Tree

import (
	"testing"
)

// ======================
// Tree Constructor Tests
// ======================

func TestNewTree(t *testing.T) {
	tree := SetupTree()

	t.Run("RootNodeType", func(t *testing.T) {
		if tree.RootNode.Type != RootType {
			t.Errorf("Expected root node type %d, got %d", RootType, tree.RootNode.Type)
		}
	})

	t.Run("RootPath", func(t *testing.T) {
		if tree.RootNode.Path != "/" {
			t.Errorf("Expected root path '/', got '%s'", tree.RootNode.Path)
		}
	})

	t.Run("ChildrenInitialized", func(t *testing.T) {
		if tree.RootNode.Children == nil {
			t.Error("Expected children slice to be initialized")
		}
	})

	t.Run("HandlersInitialized", func(t *testing.T) {
		if tree.RootNode.Handlers == nil {
			t.Error("Expected handlers map to be initialized")
		}
	})
}

// ======================
// PathWithSegment Tests
// ======================

func TestPathWithSegment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"RootPath", "/", []string{}},
		{"EmptyString", "", []string{}},
		{"SingleSegment", "/users", []string{"users"}},
		{"TwoSegments", "/users/123", []string{"users", "123"}},
		{"ThreeSegments", "/users/123/posts", []string{"users", "123", "posts"}},
		{"NoLeadingSlash", "users/123", []string{"users", "123"}},
		{"TrailingSlash", "/users/", []string{"users"}},
		{"MultipleSlashes", "//users//123//", []string{"users", "123"}},
		{"WildcardParam", "/users/:id", []string{"users", ":id"}},
		{"CatchallParam", "/files/*path", []string{"files", "*path"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pws := NewPathWithSegment(test.input)
			var result []string

			for {
				pws.Next()
				if pws.IsSame() {
					break
				}
				segment := pws.Path[pws.Start:pws.End]
				if segment != "" {
					result = append(result, segment)
				}
			}

			if len(result) != len(test.expected) {
				t.Errorf("Expected length %d, got %d for path '%s'", len(test.expected), len(result), test.input)
				t.Errorf("Expected: %v, Got: %v", test.expected, result)
				return
			}

			for i, segment := range result {
				if segment != test.expected[i] {
					t.Errorf("Expected segment[%d] = '%s', got '%s'", i, test.expected[i], segment)
				}
			}
		})
	}
}

// ======================
// Wildcard Validation Tests
// ======================

func TestIsWildCard(t *testing.T) {
	tree := SetupTree()

	tests := []TestCase{
		{"ValidWildcard", ":id", true},
		{"ValidWildcardWithText", ":user", true},
		{"EmptyString", "", false},
		{"RegularString", "id", false},
		{"CatchAllCharacter", "*", false},
		{"DoubleColon", "::", true},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := tree.IsWildCard(test.Input)
			expected := test.Expected.(bool)

			if result != expected {
				t.Errorf("IsWildCard(%q) = %v, expected %v", test.Input, result, expected)
			}
		})
	}
}

// ======================
// Catch-All Validation Tests
// ======================

func TestIsCatchAll(t *testing.T) {
	tree := SetupTree()

	tests := []TestCase{
		{"ValidCatchAll", "*", true},
		{"CatchAllWithText", "*files", true},
		{"EmptyString", "", false},
		{"RegularString", "files", false},
		{"WildcardCharacter", ":", false},
		{"DoubleAsterisk", "**", true},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := tree.IsCatchAll(test.Input)
			expected := test.Expected.(bool)

			if result != expected {
				t.Errorf("IsCatchAll(%q) = %v, expected %v", test.Input, result, expected)
			}
		})
	}
}

// ======================
// String Matching Tests
// ======================

func TestMatch(t *testing.T) {
	tree := SetupTree()

	tests := []MatchTestCase{
		{
			Name:          "ExactMatch",
			One:           "users",
			Two:           "users",
			ExpectedMatch: true,
			ExpectedIndex: 5,
			ExpectedLeft:  "",
		},
		{
			Name:          "FirstShorter",
			One:           "user",
			Two:           "users",
			ExpectedMatch: true,
			ExpectedIndex: 4,
			ExpectedLeft:  "",
		},
		{
			Name:          "SecondShorter",
			One:           "users",
			Two:           "user",
			ExpectedMatch: false,
			ExpectedIndex: 4,
			ExpectedLeft:  "s",
		},
		{
			Name:          "NoMatch",
			One:           "abc",
			Two:           "def",
			ExpectedMatch: false,
			ExpectedIndex: 0,
			ExpectedLeft:  "abc",
		},
		{
			Name:          "BothEmpty",
			One:           "",
			Two:           "",
			ExpectedMatch: true,
			ExpectedIndex: 0,
			ExpectedLeft:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var pws *PathWithSegment
			if test.One == "" {
				pws = NewPathWithSegment("")
			} else {
				pws = NewPathWithSegment("/" + test.One)
				pws.Next()
			}

			matched, index, leftPws := tree.Match(*pws, test.Two)
			left := leftPws.Path[leftPws.Start:leftPws.End]

			if matched != test.ExpectedMatch {
				t.Errorf("Expected match %v, got %v", test.ExpectedMatch, matched)
			}
			if index != test.ExpectedIndex {
				t.Errorf("Expected index %d, got %d", test.ExpectedIndex, index)
			}
			if left != test.ExpectedLeft {
				t.Errorf("Expected left '%s', got '%s'", test.ExpectedLeft, left)
			}
		})
	}
}

// ======================
// Child Node Insertion Tests
// ======================

func TestInsertChild(t *testing.T) {
	tree := SetupTree()

	t.Run("StaticChild", func(t *testing.T) {
		parent := NewNode(RootType, "/")
		child, err := tree.InsertChild(parent, "users")

		AssertNoError(t, err, "InsertChild")

		if child == nil {
			t.Fatal("Expected child node, got nil")
		}
		if child.Type != StaticType {
			t.Errorf("Expected static type %d, got %d", StaticType, child.Type)
		}
		if child.Path != "users" {
			t.Errorf("Expected path 'users', got '%s'", child.Path)
		}
	})

	t.Run("WildcardChild", func(t *testing.T) {
		parent := NewNode(RootType, "/")
		child, err := tree.InsertChild(parent, ":id")

		AssertNoError(t, err, "InsertChild wildcard")

		if child == nil {
			t.Fatal("Expected child node, got nil")
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

	t.Run("CatchAllChild", func(t *testing.T) {
		parent := NewNode(RootType, "/")
		child, err := tree.InsertChild(parent, "*files")

		AssertNoError(t, err, "InsertChild catch-all")

		if child == nil {
			t.Fatal("Expected child node, got nil")
		}
		if child.Type != CatchAllType {
			t.Errorf("Expected catch-all type %d, got %d", CatchAllType, child.Type)
		}
		if parent.CatchAll == nil {
			t.Error("Expected parent CatchAll to be set")
		}
	})

	t.Run("DuplicateErrors", func(t *testing.T) {
		parent := NewNode(RootType, "/")
		parent.WildCard = NewNode(WildCardType, ":existing")

		_, err := tree.InsertChild(parent, ":id")
		AssertError(t, err, "duplicate wildcard")

		parent2 := NewNode(RootType, "/")
		parent2.CatchAll = NewNode(CatchAllType, "*existing")

		_, err = tree.InsertChild(parent2, "*files")
		AssertError(t, err, "duplicate catch-all")
	})
}

// ======================
// Handler Setting Tests
// ======================

func TestSetHandler(t *testing.T) {
	tree := SetupTree()
	handler := CreateTestHandler()

	tests := []struct {
		name   string
		method string
		path   string
		valid  bool
	}{
		{"RootHandler", "GET", "/", true},
		{"SimplePath", "GET", "/users", true},
		{"NestedPath", "POST", "/users/123/posts", true},
		{"WildcardPath", "GET", "/users/:id", true},
		{"CatchAllPath", "GET", "/files/*path", true},
		{"EmptyMethod", "", "/test", false},
		{"EmptyPath", "GET", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testHandler HandlerFunc
			if test.valid {
				testHandler = handler
			}

			err := tree.SetHandler(tree.StringToMethodType(test.method), test.path, testHandler)

			if test.valid {
				AssertNoError(t, err, "SetHandler")
			} else {
				AssertError(t, err, "SetHandler with invalid params")
			}
		})
	}
}

// ======================
// Node Splitting Tests
// ======================

func TestSplitNode(t *testing.T) {
	tree := SetupTree()

	t.Run("BasicSplit", func(t *testing.T) {
		parent := NewNode(RootType, "/")
		child := NewNode(StaticType, "users")
		parent.Children = append(parent.Children, child)

		newNode, err := tree.SplitNode(parent, child, 4)
		AssertNoError(t, err, "SplitNode")

		if newNode == nil {
			t.Error("Expected new node, got nil")
			return
		}
		if newNode.Path != "user" {
			t.Errorf("Expected new node path 'user', got '%s'", newNode.Path)
		}
		if child.Path != "s" {
			t.Errorf("Expected child path 's', got '%s'", child.Path)
		}
	})
}
