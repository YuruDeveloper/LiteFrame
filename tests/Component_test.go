package tests

import (
	"testing"
	"LiteFrame/Router/Tree/Component"
	"github.com/stretchr/testify/assert"
)

// NodeError 테스트
func TestNodeError_Error(t *testing.T) {
	err := &Component.NodeError{
		Type:    Component.StaticType,
		Message: "테스트 에러 메시지",
		Path:    "/test/path",
	}
	
	assert.Equal(t, "테스트 에러 메시지", err.Error())
}

func TestNodeError_WithMessage(t *testing.T) {
	originalErr := &Component.NodeError{
		Type:    Component.StaticType,
		Message: "원본 메시지",
		Path:    "/test/path",
	}
	
	newErr := originalErr.WithMessage("새로운 메시지").(*Component.NodeError)
	
	assert.Equal(t, "새로운 메시지", newErr.Message)
	assert.Equal(t, Component.NodeType(Component.StaticType), newErr.Type)
	assert.Equal(t, "/test/path", newErr.Path)
}

func TestNodeError_WithPath(t *testing.T) {
	originalErr := &Component.NodeError{
		Type:    Component.StaticType,
		Message: "테스트 메시지",
		Path:    "/original/path",
	}
	
	newErr := originalErr.WithPath("/new/path").(*Component.NodeError)
	
	assert.Equal(t, "테스트 메시지", newErr.Message)
	assert.Equal(t, Component.NodeType(Component.StaticType), newErr.Type)
	assert.Equal(t, "/new/path", newErr.Path)
}

func TestNewError(t *testing.T) {
	err := Component.NewError(Component.WildCardType, "와일드카드 에러", "/wild/*").(*Component.NodeError)
	
	assert.Equal(t, Component.NodeType(Component.WildCardType), err.Type)
	assert.Equal(t, "와일드카드 에러", err.Message)
	assert.Equal(t, "/wild/*", err.Path)
}

func TestNodeTypeNames(t *testing.T) {
	tests := []struct {
		nodeType     Component.NodeType
		expectedName string
	}{
		{Component.RootType, "Root"},
		{Component.StaticType, "Static"},
		{Component.WildCardType, "WildCard"},
		{Component.CatchAllType, "CatchAll"},
		{Component.MiddlewareType, "Middleware"},
	}
	
	for _, tt := range tests {
		t.Run(tt.expectedName, func(t *testing.T) {
			assert.Equal(t, tt.expectedName, Component.NodeTypeNames[tt.nodeType])
		})
	}
}

// Identity 테스트
func TestNewIdentity(t *testing.T) {
	identity := Component.NewIdentity(Component.High, Component.StaticType, true)
	
	assert.Equal(t, Component.PriorityLevel(Component.High), identity.Priority)
	assert.Equal(t, Component.NodeType(Component.StaticType), identity.Type)
	assert.True(t, identity.Leaf)
}

func TestIdentity_GetPriority(t *testing.T) {
	identity := Component.NewIdentity(Component.Middle, Component.WildCardType, false)
	
	assert.Equal(t, Component.PriorityLevel(Component.Middle), identity.GetPriority())
}

func TestIdentity_GetType(t *testing.T) {
	identity := Component.NewIdentity(Component.Low, Component.CatchAllType, false)
	
	assert.Equal(t, Component.NodeType(Component.CatchAllType), identity.GetType())
}

func TestIdentity_IsLeaf(t *testing.T) {
	leafIdentity := Component.NewIdentity(Component.High, Component.StaticType, true)
	nonLeafIdentity := Component.NewIdentity(Component.High, Component.StaticType, false)
	
	assert.True(t, leafIdentity.IsLeaf())
	assert.False(t, nonLeafIdentity.IsLeaf())
}

// 상수 테스트
func TestNodeTypes(t *testing.T) {
	// NodeType 상수들이 올바른 값을 가지는지 확인
	assert.Equal(t, Component.NodeType(Component.RootType), Component.NodeType(1))
	assert.Equal(t, Component.NodeType(Component.StaticType), Component.NodeType(2))
	assert.Equal(t, Component.NodeType(Component.CatchAllType), Component.NodeType(3))
	assert.Equal(t, Component.NodeType(Component.WildCardType), Component.NodeType(4))
	assert.Equal(t, Component.NodeType(Component.MiddlewareType), Component.NodeType(5))
}

func TestPriorityLevels(t *testing.T) {
	// PriorityLevel 상수들이 올바른 값을 가지는지 확인
	assert.Equal(t, Component.PriorityLevel(Component.High), Component.PriorityLevel(3))
	assert.Equal(t, Component.PriorityLevel(Component.Middle), Component.PriorityLevel(2))
	assert.Equal(t, Component.PriorityLevel(Component.Low), Component.PriorityLevel(1))
}

func TestPriorityOrdering(t *testing.T) {
	// 우선순위 순서가 올바른지 확인
	assert.Greater(t, Component.High, Component.Middle)
	assert.Greater(t, Component.Middle, Component.Low)
	assert.Greater(t, Component.High, Component.Low)
}