package tests

import (
	"testing"
	"fmt"
	"LiteFrame/Router/Tree/Component"
	"github.com/stretchr/testify/assert"
)

// 테스트용 Node 구현체
type TestNode struct {
	priority Component.PriorityLevel
	nodeType Component.NodeType
	leaf     bool
	name     string
}

func (n *TestNode) GetPriority() Component.PriorityLevel {
	return n.priority
}

func (n *TestNode) GetType() Component.NodeType {
	return n.nodeType
}

func (n *TestNode) IsLeaf() bool {
	return n.leaf
}

func NewTestNode(name string, priority Component.PriorityLevel, nodeType Component.NodeType, leaf bool) *TestNode {
	return &TestNode{
		priority: priority,
		nodeType: nodeType,
		leaf:     leaf,
		name:     name,
	}
}

func TestNewContainer(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	container := Component.NewContainer(err)
	
	assert.NotNil(t, container)
	assert.NotNil(t, container.Box)
	assert.Equal(t, 0, len(container.Box))
}

func TestContainer_AddChild(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	container := Component.NewContainer(err)
	
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	
	// 성공적인 자식 추가
	err = container.AddChild("/test/path", testNode)
	assert.NoError(t, err)
	assert.Equal(t, 1, container.GetChildrenLength())
	
	// 빈 경로로 추가 시도
	err = container.AddChild("", testNode)
	assert.Error(t, err)
}

func TestContainer_GetChild(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	container := Component.NewContainer(err)
	
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	container.AddChild("/test/path", testNode)
	
	// 존재하는 자식 조회
	child := container.GetChild("/test/path")
	assert.NotNil(t, child)
	assert.Equal(t, testNode, child)
	
	// 존재하지 않는 자식 조회
	child = container.GetChild("/nonexistent")
	assert.Nil(t, child)
	
	// 빈 경로로 조회
	child = container.GetChild("")
	assert.Nil(t, child)
}

func TestContainer_SetChild(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	container := Component.NewContainer(err)
	
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	container.AddChild("/old/path", testNode)
	
	// 성공적인 경로 변경
	err = container.SetChild("/old/path", "/new/path")
	assert.NoError(t, err)
	
	// 기존 경로에는 없고 새 경로에 있는지 확인
	child := container.GetChild("/old/path")
	assert.Nil(t, child)
	
	child = container.GetChild("/new/path")
	assert.NotNil(t, child)
	assert.Equal(t, testNode, child)
	
	// 존재하지 않는 경로 변경 시도
	err = container.SetChild("/nonexistent", "/new/path2")
	assert.Error(t, err)
	
	// 빈 경로 테스트
	err = container.SetChild("", "/new/path3")
	assert.Error(t, err)
	
	err = container.SetChild("/new/path", "")
	assert.Error(t, err)
}

func TestContainer_DeleteChild(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	container := Component.NewContainer(err)
	
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	container.AddChild("/test/path", testNode)
	
	// 성공적인 삭제
	err = container.DeleteChild("/test/path")
	assert.NoError(t, err)
	assert.Equal(t, 0, container.GetChildrenLength())
	
	// 존재하지 않는 자식 삭제 시도
	err = container.DeleteChild("/nonexistent")
	assert.Error(t, err)
	
	// 빈 키로 삭제 시도
	err = container.DeleteChild("")
	assert.Error(t, err)
}

func TestContainer_GetAllChildren(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	container := Component.NewContainer(err)
	
	testNode1 := NewTestNode("test1", Component.High, Component.StaticType, true)
	testNode2 := NewTestNode("test2", Component.Middle, Component.WildCardType, false)
	
	// 빈 상태 테스트
	children := container.GetAllChildren()
	assert.NotNil(t, children)
	assert.Equal(t, 0, len(children))
	
	// 자식 추가 후 테스트
	container.AddChild("/path1", testNode1)
	container.AddChild("/path2", testNode2)
	
	children = container.GetAllChildren()
	assert.Equal(t, 2, len(children))
	assert.Contains(t, children, testNode1)
	assert.Contains(t, children, testNode2)
}

func TestContainer_GetChildrenLength(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	container := Component.NewContainer(err)
	
	// 초기 상태
	assert.Equal(t, 0, container.GetChildrenLength())
	
	// 자식 추가
	testNode1 := NewTestNode("test1", Component.High, Component.StaticType, true)
	container.AddChild("/path1", testNode1)
	assert.Equal(t, 1, container.GetChildrenLength())
	
	testNode2 := NewTestNode("test2", Component.Middle, Component.WildCardType, false)
	container.AddChild("/path2", testNode2)
	assert.Equal(t, 2, container.GetChildrenLength())
	
	// 자식 삭제
	container.DeleteChild("/path1")
	assert.Equal(t, 1, container.GetChildrenLength())
}

func TestContainer_HasChildren(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	container := Component.NewContainer(err)
	
	// 초기 상태
	assert.False(t, container.HasChildren())
	
	// 자식 추가
	testNode := NewTestNode("test1", Component.High, Component.StaticType, true)
	container.AddChild("/path1", testNode)
	assert.True(t, container.HasChildren())
	
	// 자식 삭제
	container.DeleteChild("/path1")
	assert.False(t, container.HasChildren())
}

func TestContainer_TypeSafety(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	container := Component.NewContainer(err)
	
	// 타입별 노드 생성
	staticNode := NewTestNode("static", Component.High, Component.StaticType, true)
	wildcardNode := NewTestNode("wildcard", Component.Middle, Component.WildCardType, false)
	
	// 다양한 타입의 노드 추가
	container.AddChild("/static", staticNode)
	container.AddChild("/wildcard", wildcardNode)
	
	// 타입 검증
	retrievedStatic := container.GetChild("/static")
	assert.Equal(t, int(Component.StaticType), int(retrievedStatic.GetType()))
	
	retrievedWildcard := container.GetChild("/wildcard")
	assert.Equal(t, int(Component.WildCardType), int(retrievedWildcard.GetType()))
}

func TestContainer_MultipleOperations(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	container := Component.NewContainer(err)
	
	// 여러 노드 추가
	for i := 0; i < 5; i++ {
		node := NewTestNode(fmt.Sprintf("node%d", i), Component.High, Component.StaticType, true)
		container.AddChild(fmt.Sprintf("/path%d", i), node)
	}
	
	assert.Equal(t, 5, container.GetChildrenLength())
	assert.True(t, container.HasChildren())
	
	// 중간 노드 삭제
	container.DeleteChild("/path2")
	assert.Equal(t, 4, container.GetChildrenLength())
	
	// 경로 변경
	container.SetChild("/path1", "/new_path1")
	assert.Nil(t, container.GetChild("/path1"))
	assert.NotNil(t, container.GetChild("/new_path1"))
	
	// 모든 자식 조회
	children := container.GetAllChildren()
	assert.Equal(t, 4, len(children))
}