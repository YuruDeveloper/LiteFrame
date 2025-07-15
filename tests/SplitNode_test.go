package tests

import (
	"LiteFrame/Router/Tree"
	"LiteFrame/Router/Tree/Component"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPathContainerNode_Split PathContainerNode의 Split 기능 테스트
func TestPathContainerNode_Split(t *testing.T) {

	t.Run("정상적인 Split 테스트", func(t *testing.T) {
		// PathContainerNode 생성
		pathContainer := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", "testpath"),
			"testpath",
		)

		// 자식 노드 추가
		childNode := Tree.NewStaticNode("child")
		err := pathContainer.AddChild("child", &childNode)
		assert.NoError(t, err)

		// 새로운 노드 생성 (Split 결과를 받을 노드)
		newNode := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", ""),
			"",
		)

		// Split 실행 (인덱스 4에서 분할: "test|path")
		result, err := pathContainer.Split(4, newNode)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// 분할 후 경로 확인
		assert.Equal(t, "test", result.GetPath())
		assert.Equal(t, "path", pathContainer.GetPath())

		// 자식 노드가 새로운 노드로 이동했는지 확인
		assert.True(t, result.HasChildren())
		assert.False(t, pathContainer.HasChildren())
	})

	t.Run("잘못된 Split Point - 0", func(t *testing.T) {
		pathContainer := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", "testpath"),
			"testpath",
		)

		newNode := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", ""),
			"",
		)

		_, err := pathContainer.Split(0, newNode)
		assert.Error(t, err)
	})

	t.Run("잘못된 Split Point - 경로 길이와 같음", func(t *testing.T) {
		pathContainer := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", "test"),
			"test",
		)

		newNode := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", ""),
			"",
		)

		_, err := pathContainer.Split(4, newNode)
		assert.Error(t, err)
	})

	t.Run("잘못된 Split Point - 음수", func(t *testing.T) {
		pathContainer := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", "test"),
			"test",
		)

		newNode := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", ""),
			"",
		)

		_, err := pathContainer.Split(-1, newNode)
		assert.Error(t, err)
	})
}

// TestStaticNode_Split StaticNode의 Split 기능 테스트
func TestStaticNode_Split(t *testing.T) {
	handler1 := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("handler1")) }
	handler2 := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("handler2")) }

	t.Run("핸들러가 있는 StaticNode Split", func(t *testing.T) {
		// StaticNode 생성 및 핸들러 설정
		staticNode := Tree.NewStaticNode("userprofile")
		err := staticNode.SetHandler("GET", handler1)
		assert.NoError(t, err)
		err = staticNode.SetHandler("POST", handler2)
		assert.NoError(t, err)

		// 분할용 새로운 PathContainer 생성
		newPathContainer := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", ""),
			"",
		)

		// Split 실행 (인덱스 4에서 분할: "user|profile")
		result, err := staticNode.Split(4, &newPathContainer)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// 분할 후 경로 확인
		assert.Equal(t, "user", result.GetPath())
		assert.Equal(t, "profile", staticNode.GetPath())

		// 핸들러가 결과 노드로 이동했는지 확인
		if handlerNode, ok := result.(Component.HandlerNode); ok {
			assert.True(t, handlerNode.HasMethod("GET"))
			assert.True(t, handlerNode.HasMethod("POST"))
			assert.Equal(t, 2, handlerNode.GetMethodCount())
		}

		// 원본 노드에서 핸들러가 제거되었는지 확인
		assert.False(t, staticNode.HasMethod("GET"))
		assert.False(t, staticNode.HasMethod("POST"))
		assert.Equal(t, 0, staticNode.GetMethodCount())
	})

	t.Run("자식 노드가 있는 StaticNode Split", func(t *testing.T) {
		// StaticNode 생성
		staticNode := Tree.NewStaticNode("userdata")
		
		// 자식 노드 추가
		childNode := Tree.NewStaticNode("settings")
		err := staticNode.AddChild("settings", &childNode)
		assert.NoError(t, err)

		// 분할용 새로운 PathContainer 생성
		newPathContainer := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", ""),
			"",
		)

		// Split 실행 (인덱스 4에서 분할: "user|data")
		result, err := staticNode.Split(4, &newPathContainer)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// 분할 후 경로 확인
		assert.Equal(t, "user", result.GetPath())
		assert.Equal(t, "data", staticNode.GetPath())

		// 자식 노드가 결과 노드로 이동했는지 확인
		assert.True(t, result.HasChildren())
		assert.False(t, staticNode.HasChildren())
	})
}

// TestTree_SplitNode Tree의 SplitNode 기능 테스트
func TestTree_SplitNode(t *testing.T) {
	handler1 := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("handler1")) }
	handler2 := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("handler2")) }

	t.Run("Tree SplitNode 정상 동작", func(t *testing.T) {
		tree := Tree.NewTree()

		// StaticNode 생성 및 설정
		targetNode := Tree.NewStaticNode("userprofile")
		err := targetNode.SetHandler("GET", handler1)
		assert.NoError(t, err)

		// 루트에 노드 추가
		err = tree.Root.AddChild("userprofile", &targetNode)
		assert.NoError(t, err)

		// SplitNode 실행
		paths := []string{"userinfo"} // 분할 후 추가할 경로
		newParent, err := tree.SplitNode(targetNode, &tree.Root, 4, paths, "POST", handler2)
		assert.NoError(t, err)
		assert.NotNil(t, newParent)

		// 새로운 부모 노드 확인
		assert.Equal(t, "user", newParent.GetChild("user").(Component.PathNode).GetPath())
		
		// 원본 노드의 경로가 변경되었는지 확인
		assert.Equal(t, "profile", targetNode.GetPath())
	})

	t.Run("실제 경로 분할 시나리오", func(t *testing.T) {
		tree := Tree.NewTree()

		// 첫 번째 경로 추가: "/user"
		err := tree.Add("GET", "/user", handler1)
		assert.NoError(t, err)

		// 두 번째 경로 추가: "/users" (분할이 발생해야 함)
		err = tree.Add("GET", "/users", handler2)
		assert.NoError(t, err)

		// 트리 구조 확인
		userNode := tree.Root.GetChild("user")
		assert.NotNil(t, userNode)

		// 분할이 일어났는지 확인
		if pathNode, ok := userNode.(Component.PathNode); ok {
			// "user"라는 공통 접두사로 분할되어야 함
			assert.Equal(t, "user", pathNode.GetPath())
		}

		// 자식 노드들 확인
		if container, ok := userNode.(Component.NodeContainer[Component.Node]); ok {
			children := container.GetAllChildren()
			assert.True(t, len(children) >= 1) // 적어도 하나의 자식이 있어야 함
		}
	})

	t.Run("복합 경로 분할 시나리오", func(t *testing.T) {
		tree := Tree.NewTree()

		// 여러 경로 추가하여 분할 시나리오 생성
		routes := []struct {
			method string
			path   string
		}{
			{"GET", "/user"},
			{"GET", "/users"},
			{"GET", "/user-profile"},
			{"GET", "/user/settings"},
		}

		for _, route := range routes {
			err := tree.Add(route.method, route.path, handler1)
			assert.NoError(t, err, "경로 추가 실패: %s", route.path)
		}

		// 트리 시각화 (디버깅용)
		visualizer := &TreeVisualizer{}
		treeStructure := visualizer.VisualizeTree(&tree.Root, "", true)
		t.Logf("Split 후 트리 구조:\n%s", treeStructure)

		// 기본 구조 검증
		assert.True(t, tree.Root.HasChildren())
	})
}

// TestSplitNode_EdgeCases Split 기능의 엣지 케이스 테스트
func TestSplitNode_EdgeCases(t *testing.T) {

	t.Run("한 글자 경로 분할", func(t *testing.T) {
		pathContainer := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", "ab"),
			"ab",
		)

		newNode := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", ""),
			"",
		)

		result, err := pathContainer.Split(1, newNode)
		assert.NoError(t, err)
		assert.Equal(t, "a", result.GetPath())
		assert.Equal(t, "b", pathContainer.GetPath())
	})

	t.Run("긴 경로 분할", func(t *testing.T) {
		longPath := "verylongpathnamethatshouldsplitproperly"
		pathContainer := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", longPath),
			longPath,
		)

		newNode := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", ""),
			"",
		)

		splitPoint := 10
		result, err := pathContainer.Split(splitPoint, newNode)
		assert.NoError(t, err)
		assert.Equal(t, longPath[:splitPoint], result.GetPath())
		assert.Equal(t, longPath[splitPoint:], pathContainer.GetPath())
	})

	t.Run("여러 자식이 있는 노드 분할", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {}
		staticNode := Tree.NewStaticNode("parentnode")
		
		// 여러 자식 노드 추가
		for i, childName := range []string{"child1", "child2", "child3"} {
			childNode := Tree.NewStaticNode(childName)
			err := childNode.SetHandler("GET", handler)
			assert.NoError(t, err)
			err = staticNode.AddChild(childName, &childNode)
			assert.NoError(t, err, "자식 노드 %d 추가 실패", i)
		}

		assert.Equal(t, 3, staticNode.GetChildrenLength())

		// Split 실행
		newPathContainer := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", ""),
			"",
		)

		result, err := staticNode.Split(6, &newPathContainer) // "parent|node"
		assert.NoError(t, err)

		// 모든 자식이 결과 노드로 이동했는지 확인
		assert.Equal(t, 3, result.(Component.NodeContainer[Component.Node]).GetChildrenLength())
		assert.Equal(t, 0, staticNode.GetChildrenLength())
	})
}

// TestSplitNode_Performance Split 기능의 성능 테스트
func TestSplitNode_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("성능 테스트는 짧은 테스트에서 건너뜀")
	}

	t.Run("대량 자식 노드가 있는 Split", func(t *testing.T) {
		staticNode := Tree.NewStaticNode("parentnode")

		// 1000개의 자식 노드 추가
		for i := 0; i < 1000; i++ {
			childNode := Tree.NewStaticNode("child" + string(rune(i)))
			err := staticNode.AddChild("child"+string(rune(i)), &childNode)
			assert.NoError(t, err)
		}

		newPathContainer := Component.NewPathContainerNode(
			Component.NewError(Component.StaticType, "", ""),
			"",
		)

		// Split 실행 및 시간 측정
		result, err := staticNode.Split(6, &newPathContainer)
		assert.NoError(t, err)
		assert.Equal(t, 1000, result.(Component.NodeContainer[Component.Node]).GetChildrenLength())
	})
}