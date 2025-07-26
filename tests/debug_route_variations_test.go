package tests

import (
	"fmt"
	"net/http"
	"testing"
)

// TestDebugRouteVariations는 route_variations 실패 원인을 디버깅합니다
func TestDebugRouteVariations(t *testing.T) {
	tree := SetupTree()
	handler := CreateHandlerWithResponse("response")

	// 다양한 패턴의 라우트 등록
	routes := []string{
		"/simple",
		"/with-hyphen", 
		"/with_underscore",
		"/numbers123",
		"/special.file",
	}

	// 라우트 등록 과정을 자세히 확인
	for i, route := range routes {
		fmt.Printf("등록 시도: %d번째 라우트 '%s'\n", i+1, route)
		err := tree.SetHandler("GET", route, handler)
		if err != nil {
			t.Errorf("라우트 등록 실패: %s, 에러: %v", route, err)
		} else {
			fmt.Printf("라우트 등록 성공: %s\n", route)
		}
	}

	// 트리 구조 확인
	fmt.Printf("\n=== 트리 구조 확인 ===\n")
	fmt.Printf("RootNode.Children 개수: %d\n", len(tree.RootNode.Children))
	for i, child := range tree.RootNode.Children {
		fmt.Printf("Child[%d]: Path='%s', Type=%d\n", i, child.Path, child.Type)
		if child.Handlers != nil {
			fmt.Printf("  Handlers: %v\n", child.Handlers)
		}
		// 2차 자식 노드들도 확인
		fmt.Printf("  Children 개수: %d\n", len(child.Children))
		for j, grandchild := range child.Children {
			fmt.Printf("    Grandchild[%d]: Path='%s', Type=%d\n", j, grandchild.Path, grandchild.Type)
			if grandchild.Handlers != nil {
				fmt.Printf("      Handlers: %v\n", grandchild.Handlers)
			}
		}
	}

	// 등록된 라우트들 테스트
	fmt.Printf("\n=== 라우트 테스트 ===\n")
	for i, route := range routes {
		fmt.Printf("테스트 %d: %s\n", i+1, route)
		
		recorder := ExecuteRequest(tree, "GET", route)
		fmt.Printf("  Status: %d\n", recorder.Code)
		fmt.Printf("  Body: '%s'\n", recorder.Body.String())
		
		if recorder.Code != http.StatusOK {
			t.Errorf("라우트 %s: 예상 상태코드 200, 실제 %d", route, recorder.Code)
		}
		if recorder.Body.String() != "response" {
			t.Errorf("라우트 %s: 예상 응답 'response', 실제 '%s'", route, recorder.Body.String())
		}
		fmt.Println()
	}
}