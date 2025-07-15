package tests

import (
	"testing"
	"net/http"
	"LiteFrame/Router/Tree/Component"
	"github.com/stretchr/testify/assert"
)

// 테스트용 핸들러 함수들
func testHandler1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Handler 1"))
}

func testHandler2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Handler 2"))
}

func TestNewEndPoint(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	assert.NotNil(t, endpoint)
	assert.NotNil(t, endpoint.Handlers)
	assert.Equal(t, 0, len(endpoint.Handlers))
}

func TestEndPoint_SetHandler(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	// 성공적인 핸들러 설정
	err = endpoint.SetHandler("GET", testHandler1)
	assert.NoError(t, err)
	assert.True(t, endpoint.HasMethod("GET"))
	
	// 다른 메소드 추가
	err = endpoint.SetHandler("POST", testHandler2)
	assert.NoError(t, err)
	assert.True(t, endpoint.HasMethod("POST"))
	
	// 메소드 개수 확인
	assert.Equal(t, 2, endpoint.GetMethodCount())
}

func TestEndPoint_SetHandler_InvalidInput(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	// 빈 메소드 테스트
	err = endpoint.SetHandler("", testHandler1)
	assert.Error(t, err)
	
	// nil 핸들러 테스트
	err = endpoint.SetHandler("GET", nil)
	assert.Error(t, err)
}

func TestEndPoint_GetHandler(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	// 핸들러 설정
	endpoint.SetHandler("GET", testHandler1)
	
	// 핸들러 조회
	handler := endpoint.GetHandler("GET")
	assert.NotNil(t, handler)
	
	// 존재하지 않는 메소드 조회
	handler = endpoint.GetHandler("POST")
	assert.Nil(t, handler)
}

func TestEndPoint_HasMethod(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	// 핸들러 설정 전
	assert.False(t, endpoint.HasMethod("GET"))
	
	// 핸들러 설정 후
	endpoint.SetHandler("GET", testHandler1)
	assert.True(t, endpoint.HasMethod("GET"))
	assert.False(t, endpoint.HasMethod("POST"))
}

func TestEndPoint_GetAllHandlers(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	// 빈 핸들러 맵 테스트
	handlers := endpoint.GetAllHandlers()
	assert.NotNil(t, handlers)
	assert.Equal(t, 0, len(handlers))
	
	// 핸들러 추가 후 테스트
	endpoint.SetHandler("GET", testHandler1)
	endpoint.SetHandler("POST", testHandler2)
	
	handlers = endpoint.GetAllHandlers()
	assert.Equal(t, 2, len(handlers))
	assert.NotNil(t, handlers["GET"])
	assert.NotNil(t, handlers["POST"])
}

func TestEndPoint_DeleteHandler(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	// 핸들러 설정
	endpoint.SetHandler("GET", testHandler1)
	endpoint.SetHandler("POST", testHandler2)
	
	// 핸들러 삭제
	err = endpoint.DeleteHandler("GET")
	assert.NoError(t, err)
	assert.False(t, endpoint.HasMethod("GET"))
	assert.True(t, endpoint.HasMethod("POST"))
	
	// 존재하지 않는 메소드 삭제 시도
	err = endpoint.DeleteHandler("PUT")
	assert.Error(t, err)
	
	// 빈 메소드 삭제 시도
	err = endpoint.DeleteHandler("")
	assert.Error(t, err)
}

func TestEndPoint_GetMethodCount(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	// 초기 상태
	assert.Equal(t, 0, endpoint.GetMethodCount())
	
	// 핸들러 추가
	endpoint.SetHandler("GET", testHandler1)
	assert.Equal(t, 1, endpoint.GetMethodCount())
	
	endpoint.SetHandler("POST", testHandler2)
	assert.Equal(t, 2, endpoint.GetMethodCount())
	
	// 핸들러 삭제
	endpoint.DeleteHandler("GET")
	assert.Equal(t, 1, endpoint.GetMethodCount())
}

func TestEndPoint_GetAllMethods(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	// 빈 상태
	methods := endpoint.GetAllMethods()
	assert.NotNil(t, methods)
	assert.Equal(t, 0, len(methods))
	
	// 핸들러 추가
	endpoint.SetHandler("GET", testHandler1)
	endpoint.SetHandler("POST", testHandler2)
	
	methods = endpoint.GetAllMethods()
	assert.Equal(t, 2, len(methods))
	assert.Contains(t, methods, "GET")
	assert.Contains(t, methods, "POST")
}

func TestEndPoint_NilHandlers(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	// Handlers를 nil로 설정
	endpoint.Handlers = nil
	
	// nil 상태에서 메소드 호출
	assert.Nil(t, endpoint.GetHandler("GET"))
	assert.False(t, endpoint.HasMethod("GET"))
	assert.Equal(t, 0, endpoint.GetMethodCount())
	assert.Equal(t, 0, len(endpoint.GetAllMethods()))
	
	handlers := endpoint.GetAllHandlers()
	assert.NotNil(t, handlers)
	assert.Equal(t, 0, len(handlers))
	
	// nil 상태에서 핸들러 설정
	err = endpoint.SetHandler("GET", testHandler1)
	assert.NoError(t, err)
	assert.True(t, endpoint.HasMethod("GET"))
}

func TestEndPoint_DeleteHandler_NoHandlers(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	endpoint := Component.NewEndPoint(err)
	
	// Handlers를 nil로 설정
	endpoint.Handlers = nil
	
	// nil 상태에서 삭제 시도
	err = endpoint.DeleteHandler("GET")
	assert.Error(t, err)
}