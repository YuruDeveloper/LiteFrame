package Tree

import (
	"LiteFrame/Router/Param"
	"net/http"
	"testing"
)

// ======================
// GetHandler Functionality Tests
// ======================

func TestGetHandler(t *testing.T) {
	// Basic handler path tests
	t.Run("basic_paths", func(t *testing.T) {
		testCases := []HTTPTestCase{
			{
				Name:           "root_path",
				Method:         "GET",
				Path:           "/",
				ExpectedStatus: http.StatusOK,
				ExpectedBody:   "root",
			},
			{
				Name:           "static_path",
				Method:         "GET",
				Path:           "/users",
				ExpectedStatus: http.StatusOK,
				ExpectedBody:   "users",
			},
			{
				Name:           "nested_static_path",
				Method:         "GET",
				Path:           "/users/profile",
				ExpectedStatus: http.StatusOK,
				ExpectedBody:   "user profile",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				tree := SetupTree()
				handler := CreateHandlerWithResponse(tc.ExpectedBody)

				err := tree.SetHandler(tree.StringToMethodType(tc.Method), tc.Path, handler)
				AssertNoError(t, err, "SetHandler")

				recorder := ExecuteRequest(tree, tc.Method, tc.Path)
				AssertStatusCode(t, recorder, tc.ExpectedStatus)
				AssertResponseBody(t, recorder, tc.ExpectedBody)
			})
		}
	})

	// Wildcard parameter tests
	t.Run("wildcard_parameters", func(t *testing.T) {
		t.Run("single_wildcard", func(t *testing.T) {
			tree := SetupTree()
			expectedParams := map[string]string{"id": "123"}
			handler := CreateParamCheckHandler(expectedParams)

			err := tree.SetHandler(tree.StringToMethodType("GET"), "/users/:id", handler)
			AssertNoError(t, err, "SetHandler")

			recorder := ExecuteRequest(tree, "GET", "/users/123")
			AssertStatusCode(t, recorder, http.StatusOK)
			AssertResponseBody(t, recorder, "params matched")
		})

		t.Run("multiple_wildcards", func(t *testing.T) {
			tree := SetupTree()
			expectedParams := map[string]string{
				"userId": "123",
				"postId": "456",
			}
			handler := CreateParamCheckHandler(expectedParams)

			err := tree.SetHandler(tree.StringToMethodType("GET"), "/users/:userId/posts/:postId", handler)
			AssertNoError(t, err, "SetHandler")

			recorder := ExecuteRequest(tree, "GET", "/users/123/posts/456")
			AssertStatusCode(t, recorder, http.StatusOK)
			AssertResponseBody(t, recorder, "params matched")
		})
	})

	// Catch-all path tests
	t.Run("catch_all_paths", func(t *testing.T) {
		tree := SetupTree()
		expectedParams := map[string]string{"path": "static/css/main.css"}
		handler := CreateParamCheckHandler(expectedParams)

		err := tree.SetHandler(tree.StringToMethodType("GET"), "/files/*path", handler)
		AssertNoError(t, err, "SetHandler")

		recorder := ExecuteRequest(tree, "GET", "/files/static/css/main.css")
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseBody(t, recorder, "params matched")
	})

	// Error handling tests
	t.Run("error_handling", func(t *testing.T) {
		t.Run("not_found", func(t *testing.T) {
			tree := SetupTree()
			tree.NotFoundHandler = func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("not found"))
			}

			recorder := ExecuteRequest(tree, "GET", "/nonexistent")
			AssertStatusCode(t, recorder, http.StatusNotFound)
			AssertResponseBody(t, recorder, "not found")
		})

		t.Run("method_not_allowed", func(t *testing.T) {
			tree := SetupTree()
			tree.NotAllowedHandler = func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				_, _ = w.Write([]byte("method not allowed"))
			}

			handler := CreateHandlerWithResponse("post response")
			err := tree.SetHandler(tree.StringToMethodType("POST"), "/users", handler)
			AssertNoError(t, err, "SetHandler")

			recorder := ExecuteRequest(tree, "GET", "/users")
			AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
			AssertResponseBody(t, recorder, "method not allowed")
		})
	})

	// HTTP method tests
	t.Run("multiple_methods", func(t *testing.T) {
		tree := SetupTree()

		routes := []RouteConfig{
			{"GET", "/users", CreateHandlerWithResponse("GET response")},
			{"POST", "/users", CreateHandlerWithResponse("POST response")},
		}

		_, err := SetupTreeWithRoutes(routes)
		AssertNoError(t, err, "SetupTreeWithRoutes")

		// GET request test
		for _, route := range routes {
			err := tree.SetHandler(tree.StringToMethodType(route.Method), route.Path, route.Handler)
			AssertNoError(t, err, "SetHandler")
		}

		// GET test
		recorder := ExecuteRequest(tree, "GET", "/users")
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseBody(t, recorder, "GET response")

		// POST test
		recorder = ExecuteRequest(tree, "POST", "/users")
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseBody(t, recorder, "POST response")
	})

	// Routing priority tests
	t.Run("routing_priority", func(t *testing.T) {
		tree := SetupTree()

		staticHandler := CreateHandlerWithResponse("static admin")
		wildcardHandler := CreateHandlerWithResponse("wildcard user")

		// Register wildcard first
		err := tree.SetHandler(tree.StringToMethodType("GET"), "/users/:id", wildcardHandler)
		AssertNoError(t, err, "SetHandler wildcard")

		// Register static route later
		err = tree.SetHandler(tree.StringToMethodType("GET"), "/users/admin", staticHandler)
		AssertNoError(t, err, "SetHandler static")

		// Static route should take precedence
		recorder := ExecuteRequest(tree, "GET", "/users/admin")
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseBody(t, recorder, "static admin")

		// Wildcard route test
		recorder = ExecuteRequest(tree, "GET", "/users/123")
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseBody(t, recorder, "wildcard user")
	})
}

// ======================
// Complex Routing Scenario Tests
// ======================

func TestComplexRouting(t *testing.T) {
	tree := SetupTree()

	// Complex route configuration
	routes := []struct {
		method   string
		path     string
		response string
	}{
		{"GET", "/", "home"},
		{"GET", "/users", "users list"},
		{"GET", "/users/:id", "user detail"},
		{"POST", "/users", "create user"},
		{"GET", "/users/:id/posts", "user posts"},
		{"GET", "/users/:id/posts/:postId", "user post detail"},
		{"GET", "/files/*path", "file content"},
		{"GET", "/api/v1/users", "api users"},
	}

	// Route setup
	for _, route := range routes {
		handler := CreateHandlerWithResponse(route.response)
		err := tree.SetHandler(tree.StringToMethodType(route.method), route.path, handler)
		AssertNoError(t, err, "SetHandler for "+route.path)
	}

	// Test cases
	testCases := []HTTPTestCase{
		{"home", "GET", "/", http.StatusOK, "home"},
		{"users_list", "GET", "/users", http.StatusOK, "users list"},
		{"user_detail", "GET", "/users/123", http.StatusOK, "user detail"},
		{"create_user", "POST", "/users", http.StatusOK, "create user"},
		{"user_posts", "GET", "/users/123/posts", http.StatusOK, "user posts"},
		{"user_post_detail", "GET", "/users/123/posts/456", http.StatusOK, "user post detail"},
		{"file_content", "GET", "/files/static/css/main.css", http.StatusOK, "file content"},
		{"api_users", "GET", "/api/v1/users", http.StatusOK, "api users"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			recorder := ExecuteRequest(tree, tc.Method, tc.Path)
			AssertStatusCode(t, recorder, tc.ExpectedStatus)
			AssertResponseBody(t, recorder, tc.ExpectedBody)
		})
	}
}

// ======================
// Additional Scenario Tests
// ======================

func TestAdditionalScenarios(t *testing.T) {
	t.Run("deep_path_handling", func(t *testing.T) {
		tree := SetupTree()
		handler := CreateHandlerWithResponse("deep response")

		// Register deep path
		deepPath := "/level1/level2/level3/level4/level5/deep"
		err := tree.SetHandler(tree.StringToMethodType("GET"), deepPath, handler)
		AssertNoError(t, err, "SetHandler for deep path")

		recorder := ExecuteRequest(tree, "GET", deepPath)
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseBody(t, recorder, "deep response")
	})

	t.Run("route_variations", func(t *testing.T) {
		tree := SetupTree()
		handler := CreateHandlerWithResponse("response")

		// Register various route patterns
		routes := []string{
			"/simple",
			"/with-hyphen",
			"/with_underscore",
			"/numbers123",
			"/special.file",
		}

		for _, route := range routes {
			err := tree.SetHandler(tree.StringToMethodType("GET"), route, handler)
			AssertNoError(t, err, "SetHandler for route "+route)
		}

		// Test registered routes
		for _, route := range routes {
			recorder := ExecuteRequest(tree, "GET", route)
			AssertStatusCode(t, recorder, http.StatusOK)
			AssertResponseBody(t, recorder, "response")
		}
	})
}
