package Tree

import (
	"LiteFrame/Router/Param"
	"net/http"
	"net/http/httptest"
)

// ====================
// Common Helper Functions
// ====================

// CreateTestHandler creates a basic test handler
func CreateTestHandler() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
		w.WriteHeader(http.StatusOK)
	}
}

// CreateHandlerWithResponse creates a test handler that returns specified response
func CreateHandlerWithResponse(response string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
		_, _ = w.Write([]byte(response))
	}
}

// CreateParamCheckHandler creates a test handler that validates parameters
func CreateParamCheckHandler(expectedParams map[string]string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
		if params == nil && len(expectedParams) > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("params is nil but expected parameters"))
			return
		}
		
		for key, expectedValue := range expectedParams {
			actualValue := ""
			if params != nil {
				actualValue = params.GetByName(key)
			}
			if actualValue != expectedValue {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("param mismatch: " + key))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("params matched"))
	}
}

// ====================
// Setup Functions
// ====================

// SetupTree creates and returns a basic tree
func SetupTree() Tree {
	return NewTree()
}

// SetupTreeWithRoutes sets up tree with predefined routes
func SetupTreeWithRoutes(routes []RouteConfig) (Tree, error) {
	tree := NewTree()

	for _, route := range routes {
		err := tree.SetHandler(tree.StringToMethodType(route.Method), route.Path, route.Handler)
		if err != nil {
			return tree, err
		}
	}

	return tree, nil
}

// ExecuteRequest executes an HTTP request and returns the result
func ExecuteRequest(tree Tree, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	handlerFunc, params := tree.GetHandler(req, tree.Pool.Get)
	recorder := httptest.NewRecorder()

	if handlerFunc != nil {
		handlerFunc(recorder, req, params)
		if params != nil {
			tree.Pool.Put(params)
		}
	} else {
		recorder.WriteHeader(404)
		_, _ = recorder.WriteString("handler not found")
	}

	return recorder
}

// ====================
// Test Types
// ====================

// RouteConfig holds route configuration information
type RouteConfig struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

// TestCase defines a general test case
type TestCase struct {
	Name     string
	Input    string
	Expected interface{}
}

// HTTPTestCase defines HTTP request test cases
type HTTPTestCase struct {
	Name           string
	Method         string
	Path           string
	ExpectedStatus int
	ExpectedBody   string
}

// MatchTestCase defines string matching test cases
type MatchTestCase struct {
	Name          string
	One           string
	Two           string
	ExpectedMatch bool
	ExpectedIndex int
	ExpectedLeft  string
}

// ====================
// Assertion Functions
// ====================

// AssertStatusCode validates HTTP status code
func AssertStatusCode(t TestingT, recorder *httptest.ResponseRecorder, expected int) {
	if recorder.Code != expected {
		t.Errorf("Expected status %d, got %d", expected, recorder.Code)
	}
}

// AssertResponseBody validates HTTP response body
func AssertResponseBody(t TestingT, recorder *httptest.ResponseRecorder, expected string) {
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, actual)
	}
}

// AssertNoError validates that there is no error
func AssertNoError(t TestingT, err error, operation string) {
	if err != nil {
		t.Fatalf("%s failed: %v", operation, err)
	}
}

// AssertError validates that there is an error
func AssertError(t TestingT, err error, operation string) {
	if err == nil {
		t.Errorf("Expected error for %s, got nil", operation)
	}
}

// TestingT defines testing interface (compatible with testing.T)
type TestingT interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}
