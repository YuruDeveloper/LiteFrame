package tests

import (
	"testing"
	"net/http"
	"LiteFrame/Router/Tree"
	"LiteFrame/Router/Tree/Component"
	"github.com/stretchr/testify/assert"
)

// 테스트용 핸들러 함수들
func staticHandler1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Static Handler 1"))
}

func staticHandler2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Static Handler 2"))
}

func TestNewStaticNode(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	assert.NotNil(t, node.Identity)
	assert.NotNil(t, node.PathContainer)
	assert.Nil(t, node.EndPoint)
	assert.Equal(t, "/api/users", node.GetPath())
	assert.Equal(t, Component.NodeType(Component.StaticType), node.GetType())
	assert.Equal(t, Component.PriorityLevel(Component.High), node.GetPriority())
	assert.False(t, node.IsLeaf()) // EndPoint가 nil이므로 leaf가 아님
}

// Node 인터페이스 테스트
func TestStaticNode_GetPriority(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	assert.Equal(t, Component.PriorityLevel(Component.High), node.GetPriority())
}

func TestStaticNode_GetType(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	assert.Equal(t, Component.NodeType(Component.StaticType), node.GetType())
}

func TestStaticNode_IsLeaf(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	// EndPoint가 없을 때는 leaf가 아님
	assert.False(t, node.IsLeaf())
	
	// 핸들러 설정 후에는 leaf가 됨
	node.SetHandler("GET", staticHandler1)
	assert.True(t, node.IsLeaf())
}

// PathNode 인터페이스 테스트
func TestStaticNode_GetPath(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	assert.Equal(t, "/api/users", node.GetPath())
}

func TestStaticNode_SetPath(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	err := node.SetPath("/api/posts")
	assert.NoError(t, err)
	assert.Equal(t, "/api/posts", node.GetPath())
	
	// 빈 경로 설정 시도
	err = node.SetPath("")
	assert.Error(t, err)
}

func TestStaticNode_Match(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	// 정확히 일치
	matched, matchingChar, leftPath := node.Match("/api/users")
	assert.True(t, matched)
	assert.Equal(t, 10, matchingChar)
	assert.Equal(t, "", leftPath)
	
	// 부분 일치
	matched, matchingChar, leftPath = node.Match("/api/users/123")
	assert.True(t, matched)
	assert.Equal(t, 10, matchingChar)
	assert.Equal(t, "/123", leftPath)
	
	// 불일치
	matched, matchingChar, leftPath = node.Match("/api/posts")
	assert.False(t, matched)
	assert.Equal(t, 5, matchingChar)
	assert.Equal(t, "posts", leftPath)
}

// NodeContainer 인터페이스 테스트
func TestStaticNode_AddChild(t *testing.T) {
	node := Tree.NewStaticNode("/api")
	childNode := Tree.NewStaticNode("/api/users")
	
	err := node.AddChild("/users", &childNode)
	assert.NoError(t, err)
	assert.Equal(t, 1, node.GetChildrenLength())
	assert.True(t, node.HasChildren())
}

func TestStaticNode_GetChild(t *testing.T) {
	node := Tree.NewStaticNode("/api")
	childNode := Tree.NewStaticNode("/api/users")
	
	node.AddChild("/users", &childNode)
	
	retrievedChild := node.GetChild("/users")
	assert.NotNil(t, retrievedChild)
	
	// 존재하지 않는 자식 조회
	nonExistentChild := node.GetChild("/posts")
	assert.Nil(t, nonExistentChild)
}

func TestStaticNode_DeleteChild(t *testing.T) {
	node := Tree.NewStaticNode("/api")
	childNode := Tree.NewStaticNode("/api/users")
	
	node.AddChild("/users", &childNode)
	assert.Equal(t, 1, node.GetChildrenLength())
	
	err := node.DeleteChild("/users")
	assert.NoError(t, err)
	assert.Equal(t, 0, node.GetChildrenLength())
	assert.False(t, node.HasChildren())
}

func TestStaticNode_GetAllChildren(t *testing.T) {
	node := Tree.NewStaticNode("/api")
	childNode1 := Tree.NewStaticNode("/api/users")
	childNode2 := Tree.NewStaticNode("/api/posts")
	
	node.AddChild("/users", &childNode1)
	node.AddChild("/posts", &childNode2)
	
	children := node.GetAllChildren()
	assert.Equal(t, 2, len(children))
}

// HandleAccessor 인터페이스 테스트
func TestStaticNode_SetHandler(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	// 처음에는 EndPoint가 nil
	assert.Nil(t, node.EndPoint)
	
	// 핸들러 설정
	err := node.SetHandler("GET", staticHandler1)
	assert.NoError(t, err)
	assert.NotNil(t, node.EndPoint) // EndPoint가 생성됨
	assert.True(t, node.IsLeaf())   // 이제 leaf 노드가 됨
	
	// 다른 메소드 추가
	err = node.SetHandler("POST", staticHandler2)
	assert.NoError(t, err)
	assert.Equal(t, 2, node.GetMethodCount())
}

func TestStaticNode_GetHandler(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	// EndPoint가 없을 때
	handler := node.GetHandler("GET")
	assert.Nil(t, handler)
	
	// 핸들러 설정 후
	node.SetHandler("GET", staticHandler1)
	handler = node.GetHandler("GET")
	assert.NotNil(t, handler)
	
	// 존재하지 않는 메소드
	handler = node.GetHandler("PUT")
	assert.Nil(t, handler)
}

func TestStaticNode_HasMethod(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	// EndPoint가 없을 때
	assert.False(t, node.HasMethod("GET"))
	
	// 핸들러 설정 후
	node.SetHandler("GET", staticHandler1)
	assert.True(t, node.HasMethod("GET"))
	assert.False(t, node.HasMethod("POST"))
}

func TestStaticNode_GetAllHandlers(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	// EndPoint가 없을 때
	handlers := node.GetAllHandlers()
	assert.NotNil(t, handlers)
	assert.Equal(t, 0, len(handlers))
	
	// 핸들러 설정 후
	node.SetHandler("GET", staticHandler1)
	node.SetHandler("POST", staticHandler2)
	
	handlers = node.GetAllHandlers()
	assert.Equal(t, 2, len(handlers))
	assert.NotNil(t, handlers["GET"])
	assert.NotNil(t, handlers["POST"])
}

func TestStaticNode_DeleteHandler(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	// EndPoint가 없을 때 삭제 시도
	err := node.DeleteHandler("GET")
	assert.Error(t, err)
	
	// 핸들러 설정 후 삭제
	node.SetHandler("GET", staticHandler1)
	node.SetHandler("POST", staticHandler2)
	
	err = node.DeleteHandler("GET")
	assert.NoError(t, err)
	assert.False(t, node.HasMethod("GET"))
	assert.True(t, node.HasMethod("POST"))
	assert.Equal(t, 1, node.GetMethodCount())
}

func TestStaticNode_GetMethodCount(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	// EndPoint가 없을 때
	assert.Equal(t, 0, node.GetMethodCount())
	
	// 핸들러 추가
	node.SetHandler("GET", staticHandler1)
	assert.Equal(t, 1, node.GetMethodCount())
	
	node.SetHandler("POST", staticHandler2)
	assert.Equal(t, 2, node.GetMethodCount())
}

func TestStaticNode_GetAllMethods(t *testing.T) {
	node := Tree.NewStaticNode("/api/users")
	
	// EndPoint가 없을 때
	methods := node.GetAllMethods()
	assert.NotNil(t, methods)
	assert.Equal(t, 0, len(methods))
	
	// 핸들러 추가
	node.SetHandler("GET", staticHandler1)
	node.SetHandler("POST", staticHandler2)
	
	methods = node.GetAllMethods()
	assert.Equal(t, 2, len(methods))
	assert.Contains(t, methods, "GET")
	assert.Contains(t, methods, "POST")
}

func TestStaticNode_SetChild(t *testing.T) {
	node := Tree.NewStaticNode("/api")
	childNode := Tree.NewStaticNode("/api/users")
	
	node.AddChild("/users", &childNode)
	
	// 경로 변경
	err := node.SetChild("/users", "/newusers")
	assert.NoError(t, err)
	
	// 기존 경로에는 없고 새 경로에 있는지 확인
	assert.Nil(t, node.GetChild("/users"))
	assert.NotNil(t, node.GetChild("/newusers"))
}

func TestStaticNode_Integration(t *testing.T) {
	// 루트 노드 생성
	rootNode := Tree.NewStaticNode("/api")
	
	// 자식 노드들 생성 및 추가
	usersNode := Tree.NewStaticNode("/api/users")
	postsNode := Tree.NewStaticNode("/api/posts")
	
	rootNode.AddChild("/users", &usersNode)
	rootNode.AddChild("/posts", &postsNode)
	
	// 핸들러 설정
	usersNode.SetHandler("GET", staticHandler1)
	usersNode.SetHandler("POST", staticHandler2)
	postsNode.SetHandler("GET", staticHandler1)
	
	// 검증
	assert.Equal(t, 2, rootNode.GetChildrenLength())
	assert.True(t, rootNode.HasChildren())
	assert.False(t, rootNode.IsLeaf()) // 핸들러가 없으므로 leaf가 아님
	
	assert.True(t, usersNode.IsLeaf()) // 핸들러가 있으므로 leaf
	assert.Equal(t, 2, usersNode.GetMethodCount())
	assert.True(t, usersNode.HasMethod("GET"))
	assert.True(t, usersNode.HasMethod("POST"))
	
	assert.True(t, postsNode.IsLeaf())
	assert.Equal(t, 1, postsNode.GetMethodCount())
	assert.True(t, postsNode.HasMethod("GET"))
	assert.False(t, postsNode.HasMethod("POST"))
	
	// 경로 매칭 테스트
	matched, _, _ := usersNode.Match("/api/users/123")
	assert.True(t, matched)
	
	matched, _, _ = postsNode.Match("/api/posts")
	assert.True(t, matched)
}