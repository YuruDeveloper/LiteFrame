package tests

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"LiteFrame/Router/Tree"
	"LiteFrame/Router/Tree/Component"
	"github.com/stretchr/testify/assert"
)

// 테스트용 핸들러 함수들
func wildcardHandler1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("WildCard Handler 1"))
}

func wildcardHandler2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("WildCard Handler 2"))
}

func TestNewWildCardNode(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	assert.NotNil(t, node.Identity)
	assert.NotNil(t, node.Container)
	assert.NotNil(t, node.PathHandler)
	assert.Nil(t, node.EndPoint)
	assert.Equal(t, "id", node.GetPath()) // ':'를 제거한 경로
	assert.Equal(t, Component.NodeType(Component.WildCardType), node.GetType())
	assert.Equal(t, Component.PriorityLevel(Component.Middle), node.GetPriority())
	assert.False(t, node.IsLeaf()) // EndPoint가 nil이므로 leaf가 아님
	assert.Equal(t, "", node.Data)
}

// Node 인터페이스 테스트
func TestWildCardNode_GetPriority(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	assert.Equal(t, Component.PriorityLevel(Component.Middle), node.GetPriority())
}

func TestWildCardNode_GetType(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	assert.Equal(t, Component.NodeType(Component.WildCardType), node.GetType())
}

func TestWildCardNode_IsLeaf(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	// EndPoint가 없을 때는 leaf가 아님
	assert.False(t, node.IsLeaf())
	
	// 핸들러 설정 후에는 leaf가 됨
	node.SetHandler("GET", wildcardHandler1)
	assert.True(t, node.IsLeaf())
}

// PathNode 인터페이스 테스트
func TestWildCardNode_GetPath(t *testing.T) {
	node := Tree.NewWildCardNode(":userId")
	assert.Equal(t, "userId", node.GetPath()) // ':' 제거됨
}

func TestWildCardNode_SetPath(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	err := node.SetPath("newId")
	assert.NoError(t, err)
	assert.Equal(t, "newId", node.GetPath())
	
	// 빈 경로 설정 시도
	err = node.SetPath("")
	assert.Error(t, err)
}

func TestWildCardNode_Match(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	// 정상적인 매칭
	matched, matchingChar, leftPath := node.Match("123")
	assert.True(t, matched)
	assert.Equal(t, 3, matchingChar)
	assert.Equal(t, "123", leftPath)
	assert.Equal(t, "123", node.Data) // Data 필드에 매칭된 값 저장
	
	// 다른 값으로 매칭
	matched, matchingChar, leftPath = node.Match("abc")
	assert.True(t, matched)
	assert.Equal(t, 3, matchingChar)
	assert.Equal(t, "abc", leftPath)
	assert.Equal(t, "abc", node.Data)
	
	// 빈 문자열 매칭 (실패해야 함)
	matched, matchingChar, leftPath = node.Match("")
	assert.False(t, matched)
	assert.Equal(t, 0, matchingChar)
	assert.Equal(t, "", leftPath)
}

// NodeContainer 인터페이스 테스트
func TestWildCardNode_AddChild(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	childNode := Tree.NewWildCardNode(":action")
	
	err := node.AddChild("/action", &childNode)
	assert.NoError(t, err)
	assert.Equal(t, 1, node.GetChildrenLength())
	assert.True(t, node.HasChildren())
}

func TestWildCardNode_GetChild(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	childNode := Tree.NewWildCardNode(":action")
	
	node.AddChild("/action", &childNode)
	
	retrievedChild := node.GetChild("/action")
	assert.NotNil(t, retrievedChild)
	
	// 존재하지 않는 자식 조회
	nonExistentChild := node.GetChild("/nonexistent")
	assert.Nil(t, nonExistentChild)
}

func TestWildCardNode_DeleteChild(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	childNode := Tree.NewWildCardNode(":action")
	
	node.AddChild("/action", &childNode)
	assert.Equal(t, 1, node.GetChildrenLength())
	
	err := node.DeleteChild("/action")
	assert.NoError(t, err)
	assert.Equal(t, 0, node.GetChildrenLength())
	assert.False(t, node.HasChildren())
}

func TestWildCardNode_GetAllChildren(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	childNode1 := Tree.NewWildCardNode(":action")
	childNode2 := Tree.NewWildCardNode(":type")
	
	node.AddChild("/action", &childNode1)
	node.AddChild("/type", &childNode2)
	
	children := node.GetAllChildren()
	assert.Equal(t, 2, len(children))
}

// HandleAccessor 인터페이스 테스트
func TestWildCardNode_SetHandler(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	// 처음에는 EndPoint가 nil
	assert.Nil(t, node.EndPoint)
	
	// 핸들러 설정
	err := node.SetHandler("GET", wildcardHandler1)
	assert.NoError(t, err)
	assert.NotNil(t, node.EndPoint) // EndPoint가 생성됨
	assert.True(t, node.IsLeaf())   // 이제 leaf 노드가 됨
	
	// 다른 메소드 추가
	err = node.SetHandler("POST", wildcardHandler2)
	assert.NoError(t, err)
	assert.Equal(t, 2, node.GetMethodCount())
}

func TestWildCardNode_GetHandler(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	// EndPoint가 없을 때
	handler := node.GetHandler("GET")
	assert.Nil(t, handler)
	
	// 핸들러 설정 후
	node.SetHandler("GET", wildcardHandler1)
	handler = node.GetHandler("GET")
	assert.NotNil(t, handler)
	
	// 존재하지 않는 메소드
	handler = node.GetHandler("PUT")
	assert.Nil(t, handler)
}

func TestWildCardNode_HasMethod(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	// EndPoint가 없을 때
	assert.False(t, node.HasMethod("GET"))
	
	// 핸들러 설정 후
	node.SetHandler("GET", wildcardHandler1)
	assert.True(t, node.HasMethod("GET"))
	assert.False(t, node.HasMethod("POST"))
}

func TestWildCardNode_GetAllHandlers(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	// EndPoint가 없을 때
	handlers := node.GetAllHandlers()
	assert.NotNil(t, handlers)
	assert.Equal(t, 0, len(handlers))
	
	// 핸들러 설정 후
	node.SetHandler("GET", wildcardHandler1)
	node.SetHandler("POST", wildcardHandler2)
	
	handlers = node.GetAllHandlers()
	assert.Equal(t, 2, len(handlers))
	assert.NotNil(t, handlers["GET"])
	assert.NotNil(t, handlers["POST"])
}

func TestWildCardNode_DeleteHandler(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	// EndPoint가 없을 때 삭제 시도
	err := node.DeleteHandler("GET")
	assert.Error(t, err)
	
	// 핸들러 설정 후 삭제
	node.SetHandler("GET", wildcardHandler1)
	node.SetHandler("POST", wildcardHandler2)
	
	err = node.DeleteHandler("GET")
	assert.NoError(t, err)
	assert.False(t, node.HasMethod("GET"))
	assert.True(t, node.HasMethod("POST"))
	assert.Equal(t, 1, node.GetMethodCount())
}

func TestWildCardNode_GetMethodCount(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	// EndPoint가 없을 때
	assert.Equal(t, 0, node.GetMethodCount())
	
	// 핸들러 추가
	node.SetHandler("GET", wildcardHandler1)
	assert.Equal(t, 1, node.GetMethodCount())
	
	node.SetHandler("POST", wildcardHandler2)
	assert.Equal(t, 2, node.GetMethodCount())
}

func TestWildCardNode_GetAllMethods(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	
	// EndPoint가 없을 때
	methods := node.GetAllMethods()
	assert.NotNil(t, methods)
	assert.Equal(t, 0, len(methods))
	
	// 핸들러 추가
	node.SetHandler("GET", wildcardHandler1)
	node.SetHandler("POST", wildcardHandler2)
	
	methods = node.GetAllMethods()
	assert.Equal(t, 2, len(methods))
	assert.Contains(t, methods, "GET")
	assert.Contains(t, methods, "POST")
}

// SetWildCard 메소드 테스트 (컨텍스트 처리)
func TestWildCardNode_SetWildCard(t *testing.T) {
	node := Tree.NewWildCardNode(":userId")
	node.Data = "123" // 매칭된 데이터 설정
	
	// 테스트용 내부 핸들러 - WildCardNode의 실제 키 타입 사용
	innerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// WildCardNode에서 사용하는 키 타입으로 조회
		val := r.Context().Value(Component.TreeKey(node.PathHandler.GetPath()))
		if val != nil {
			w.Write([]byte("User ID: " + val.(string)))
		} else {
			w.Write([]byte("No User ID"))
		}
	})
	
	// SetWildCard로 핸들러 래핑
	wrappedHandler := node.SetWildCard(innerHandler)
	
	// 테스트 요청 생성
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	
	// 핸들러 실행
	wrappedHandler.ServeHTTP(w, req)
	
	// 응답 검증
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "User ID: 123", w.Body.String())
}

func TestWildCardNode_SetWildCard_Context(t *testing.T) {
	node := Tree.NewWildCardNode(":productId")
	node.Data = "abc123"
	
	// 컨텍스트 값을 확인하는 핸들러
	checkHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 컨텍스트에서 값 확인 - 실제 키 사용
		val := r.Context().Value(Component.TreeKey(node.PathHandler.GetPath()))
		assert.Equal(t, "abc123", val)
		w.WriteHeader(http.StatusOK)
	})
	
	wrappedHandler := node.SetWildCard(checkHandler)
	
	req := httptest.NewRequest("GET", "/products/abc123", nil)
	w := httptest.NewRecorder()
	
	wrappedHandler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWildCardNode_Integration(t *testing.T) {
	// 와일드카드 노드 생성
	userNode := Tree.NewWildCardNode(":userId")
	
	// 자식 노드 추가
	profileNode := Tree.NewWildCardNode(":section")
	userNode.AddChild("/profile", &profileNode)
	
	// 핸들러 설정
	userNode.SetHandler("GET", wildcardHandler1)
	profileNode.SetHandler("GET", wildcardHandler2)
	
	// 매칭 테스트
	matched, _, _ := userNode.Match("123")
	assert.True(t, matched)
	assert.Equal(t, "123", userNode.Data)
	
	// 자식 노드 확인
	assert.Equal(t, 1, userNode.GetChildrenLength())
	assert.True(t, userNode.HasChildren())
	
	// 핸들러 확인
	assert.True(t, userNode.IsLeaf())
	assert.True(t, userNode.HasMethod("GET"))
	assert.Equal(t, 1, userNode.GetMethodCount())
	
	// 자식 노드도 확인
	assert.True(t, profileNode.IsLeaf())
	assert.True(t, profileNode.HasMethod("GET"))
}

func TestWildCardNode_MultipleWildcards(t *testing.T) {
	// 다중 와일드카드 시나리오
	userNode := Tree.NewWildCardNode(":userId")
	postNode := Tree.NewWildCardNode(":postId")
	commentNode := Tree.NewWildCardNode(":commentId")
	
	userNode.AddChild("/posts", &postNode)
	postNode.AddChild("/comments", &commentNode)
	
	// 각각 핸들러 설정
	userNode.SetHandler("GET", wildcardHandler1)
	postNode.SetHandler("GET", wildcardHandler1)
	commentNode.SetHandler("GET", wildcardHandler1)
	
	// 계층적 구조 확인
	assert.True(t, userNode.HasChildren())
	assert.True(t, postNode.HasChildren())
	assert.False(t, commentNode.HasChildren())
	
	// 모두 leaf 노드여야 함 (핸들러가 있으므로)
	assert.True(t, userNode.IsLeaf())
	assert.True(t, postNode.IsLeaf())
	assert.True(t, commentNode.IsLeaf())
	
	// 매칭 테스트
	userNode.Match("user123")
	assert.Equal(t, "user123", userNode.Data)
	
	postNode.Match("post456")
	assert.Equal(t, "post456", postNode.Data)
	
	commentNode.Match("comment789")
	assert.Equal(t, "comment789", commentNode.Data)
}

func TestWildCardNode_SetChild(t *testing.T) {
	node := Tree.NewWildCardNode(":id")
	childNode := Tree.NewWildCardNode(":action")
	
	node.AddChild("/action", &childNode)
	
	// 경로 변경
	err := node.SetChild("/action", "/newaction")
	assert.NoError(t, err)
	
	// 기존 경로에는 없고 새 경로에 있는지 확인
	assert.Nil(t, node.GetChild("/action"))
	assert.NotNil(t, node.GetChild("/newaction"))
}