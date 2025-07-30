package tests

import (
	"LiteFrame/Router/Param"
	"LiteFrame/Router/Tree"
	"net/http"
	"net/http/httptest"
)

// ====================
// 공통 Helper 함수들
// ====================

// CreateTestHandler는 기본 테스트 핸들러를 생성합니다
func CreateTestHandler() Tree.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
		w.WriteHeader(http.StatusOK)
	}
}

// CreateHandlerWithResponse는 지정된 응답을 반환하는 테스트 핸들러를 생성합니다
func CreateHandlerWithResponse(response string) Tree.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
		w.Write([]byte(response))
	}
}

// CreateParamCheckHandler는 매개변수를 검증하는 테스트 핸들러를 생성합니다
func CreateParamCheckHandler(expectedParams map[string]string) Tree.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params *Param.Params) {
		for key, expectedValue := range expectedParams {
			actualValue := params.GetByName(key)
			if actualValue != expectedValue {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("param mismatch: " + key))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("params matched"))
	}
}

// ====================
// 테스트 도우미 함수들
// ====================

// SetupTree는 기본 트리를 생성하고 반환합니다
func SetupTree() Tree.Tree {
	return Tree.NewTree()
}

// SetupTreeWithRoutes는 미리 정의된 라우트들로 트리를 설정합니다
func SetupTreeWithRoutes(routes []RouteConfig) (Tree.Tree, error) {
	tree := Tree.NewTree()

	for _, route := range routes {
		err := tree.SetHandler(route.Method, route.Path, route.Handler)
		if err != nil {
			return tree, err
		}
	}

	return tree, nil
}

// ExecuteRequest는 HTTP 요청을 실행하고 결과를 반환합니다
func ExecuteRequest(tree Tree.Tree, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	handlerFunc, params := tree.GetHandler(req, tree.Pool.Get)
	recorder := httptest.NewRecorder()

	if handlerFunc != nil {
		handlerFunc(recorder, req, params)
		// 매개변수 객체를 풀에 반환
		if params != nil {
			tree.Pool.Put(params)
		}
	} else {
		// GetHandler가 nil을 반환하면 404 응답
		recorder.WriteHeader(404)
		recorder.WriteString("handler not found")
	}

	return recorder
}

// ====================
// 테스트 데이터 구조체들
// ====================

// RouteConfig는 라우트 설정 정보를 담습니다
type RouteConfig struct {
	Method  string
	Path    string
	Handler Tree.HandlerFunc
}

// TestCase는 일반적인 테스트 케이스를 정의합니다
type TestCase struct {
	Name     string
	Input    string
	Expected interface{}
}

// HTTPTestCase는 HTTP 요청 테스트 케이스를 정의합니다
type HTTPTestCase struct {
	Name           string
	Method         string
	Path           string
	ExpectedStatus int
	ExpectedBody   string
}

// MatchTestCase는 문자열 매칭 테스트 케이스를 정의합니다
type MatchTestCase struct {
	Name          string
	One           string
	Two           string
	ExpectedMatch bool
	ExpectedIndex int
	ExpectedLeft  string
}

// ====================
// 검증 도우미 함수들
// ====================

// AssertStatusCode는 HTTP 상태 코드를 검증합니다
func AssertStatusCode(t TestingT, recorder *httptest.ResponseRecorder, expected int) {
	if recorder.Code != expected {
		t.Errorf("Expected status %d, got %d", expected, recorder.Code)
	}
}

// AssertResponseBody는 HTTP 응답 본문을 검증합니다
func AssertResponseBody(t TestingT, recorder *httptest.ResponseRecorder, expected string) {
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, actual)
	}
}

// AssertNoError는 에러가 없음을 검증합니다
func AssertNoError(t TestingT, err error, operation string) {
	if err != nil {
		t.Fatalf("%s failed: %v", operation, err)
	}
}

// AssertError는 에러가 있음을 검증합니다
func AssertError(t TestingT, err error, operation string) {
	if err == nil {
		t.Errorf("Expected error for %s, got nil", operation)
	}
}

// TestingT는 테스팅 인터페이스를 정의합니다 (testing.T와 호환)
type TestingT interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}
