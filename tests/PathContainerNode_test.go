package tests

import (
	"testing"
	"LiteFrame/Router/Tree/Component"
	"github.com/stretchr/testify/assert"
)

func TestNewPathContainerNode(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	assert.NotNil(t, node)
	assert.Equal(t, "/api/users", node.GetPath())
	assert.Equal(t, 0, node.GetChildrenLength())
}

// PathNode 인터페이스 테스트
func TestPathContainerNode_GetPath(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	assert.Equal(t, "/api/users", node.GetPath())
}

func TestPathContainerNode_SetPath(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	// 성공적인 경로 설정
	err = node.SetPath("/api/posts")
	assert.NoError(t, err)
	assert.Equal(t, "/api/posts", node.GetPath())
	
	// 빈 경로 설정 시도
	err = node.SetPath("")
	assert.Error(t, err)
	assert.Equal(t, "/api/posts", node.GetPath()) // 기존 경로 유지
}

func TestPathContainerNode_Match(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	// 정확히 일치하는 경우
	matched, matchingChar, leftPath := node.Match("/api/users")
	assert.True(t, matched)
	assert.Equal(t, 10, matchingChar)
	assert.Equal(t, "", leftPath)
	
	// 부분적으로 일치하는 경우
	matched, matchingChar, leftPath = node.Match("/api/users/123")
	assert.True(t, matched)
	assert.Equal(t, 10, matchingChar)
	assert.Equal(t, "/123", leftPath)
	
	// 일치하지 않는 경우
	matched, matchingChar, leftPath = node.Match("/api/posts")
	assert.False(t, matched)
	assert.Equal(t, 5, matchingChar)
	assert.Equal(t, "posts", leftPath)
}

// NodeContainer 인터페이스 테스트
func TestPathContainerNode_AddChild(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	
	// 성공적인 자식 추가
	err = node.AddChild("/child", testNode)
	assert.NoError(t, err)
	assert.Equal(t, 1, node.GetChildrenLength())
	
	// 빈 경로로 추가 시도
	err = node.AddChild("", testNode)
	assert.Error(t, err)
}

func TestPathContainerNode_GetChild(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	node.AddChild("/child", testNode)
	
	// 존재하는 자식 조회
	child := node.GetChild("/child")
	assert.NotNil(t, child)
	assert.Equal(t, testNode, child)
	
	// 존재하지 않는 자식 조회
	child = node.GetChild("/nonexistent")
	assert.Nil(t, child)
	
	// 빈 경로로 조회
	child = node.GetChild("")
	assert.Nil(t, child)
}

func TestPathContainerNode_SetChild(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	node.AddChild("/old", testNode)
	
	// 성공적인 경로 변경
	err = node.SetChild("/old", "/new")
	assert.NoError(t, err)
	
	// 기존 경로에는 없고 새 경로에 있는지 확인
	child := node.GetChild("/old")
	assert.Nil(t, child)
	
	child = node.GetChild("/new")
	assert.NotNil(t, child)
	assert.Equal(t, testNode, child)
}

func TestPathContainerNode_DeleteChild(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	node.AddChild("/child", testNode)
	
	// 성공적인 삭제
	err = node.DeleteChild("/child")
	assert.NoError(t, err)
	assert.Equal(t, 0, node.GetChildrenLength())
	
	// 존재하지 않는 자식 삭제 시도
	err = node.DeleteChild("/nonexistent")
	assert.Error(t, err)
}

func TestPathContainerNode_GetAllChildren(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	testNode1 := NewTestNode("test1", Component.High, Component.StaticType, true)
	testNode2 := NewTestNode("test2", Component.Middle, Component.WildCardType, false)
	
	// 빈 상태 테스트
	children := node.GetAllChildren()
	assert.NotNil(t, children)
	assert.Equal(t, 0, len(children))
	
	// 자식 추가 후 테스트
	node.AddChild("/child1", testNode1)
	node.AddChild("/child2", testNode2)
	
	children = node.GetAllChildren()
	assert.Equal(t, 2, len(children))
	assert.Contains(t, children, testNode1)
	assert.Contains(t, children, testNode2)
}

func TestPathContainerNode_GetChildrenLength(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	// 초기 상태
	assert.Equal(t, 0, node.GetChildrenLength())
	
	// 자식 추가
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	node.AddChild("/child", testNode)
	assert.Equal(t, 1, node.GetChildrenLength())
}

func TestPathContainerNode_HasChildren(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	// 초기 상태
	assert.False(t, node.HasChildren())
	
	// 자식 추가
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	node.AddChild("/child", testNode)
	assert.True(t, node.HasChildren())
}

// Split 메소드 테스트 - 기본적인 기능만 테스트
func TestPathContainerNode_Split_InvalidSplitPoint(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	newNode := Component.NewPathContainerNode(err, "")
	
	// 유효하지 않은 분할 지점 테스트
	_, err = node.Split(0, newNode)
	assert.Error(t, err)
	
	_, err = node.Split(-1, newNode)
	assert.Error(t, err)
	
	_, err = node.Split(10, newNode) // 경로 길이와 같거나 큰 값
	assert.Error(t, err)
}

func TestPathContainerNode_Integration(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	node := Component.NewPathContainerNode(err, "/api/users")
	
	// 여러 자식 노드 추가
	testNode1 := NewTestNode("test1", Component.High, Component.StaticType, true)
	testNode2 := NewTestNode("test2", Component.Middle, Component.WildCardType, false)
	
	node.AddChild("/child1", testNode1)
	node.AddChild("/child2", testNode2)
	
	assert.Equal(t, 2, node.GetChildrenLength())
	assert.True(t, node.HasChildren())
	
	// 경로 매칭 테스트
	matched, matchingChar, leftPath := node.Match("/api/users/123")
	assert.True(t, matched)
	assert.Equal(t, 10, matchingChar)
	assert.Equal(t, "/123", leftPath)
	
	// 자식 삭제
	node.DeleteChild("/child1")
	assert.Equal(t, 1, node.GetChildrenLength())
	
	// 경로 변경
	node.SetPath("/api/posts")
	assert.Equal(t, "/api/posts", node.GetPath())
}