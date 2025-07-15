package tests

import (
	"LiteFrame/Router/Tree"
	"LiteFrame/Router/Tree/Component"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// ExpectedTreeNode 예상되는 트리 노드 구조를 나타내는 구조체
type ExpectedTreeNode struct {
	Path     string
	Type     Component.NodeType
	Children []ExpectedTreeNode
	Methods  []string
}

// TreeVisualizer 트리를 시각화하는 구조체
type TreeVisualizer struct {
	Builder strings.Builder
}

// NodeTypeToString NodeType을 문자열로 변환하는 헬퍼 함수
func NodeTypeToString(nodeType Component.NodeType) string {
	switch nodeType {
	case Component.RootType:
		return "Root"
	case Component.StaticType:
		return "Static"
	case Component.CatchAllType:
		return "CatchAll"
	case Component.WildCardType:
		return "WildCard"
	case Component.MiddlewareType:
		return "Middleware"
	default:
		return "Unknown"
	}
}

// VisualizeTree 실제 트리를 시각화하는 함수
func (v *TreeVisualizer) VisualizeTree(node Component.Node, prefix string, isLast bool) string {
	if node == nil {
		return ""
	}

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	nodeInfo := fmt.Sprintf("%s%s[%s] %s", prefix, connector, NodeTypeToString(node.GetType()), v.getNodePath(node))
	
	if handlerNode, ok := node.(Component.HandlerNode); ok {
		methods := v.getHandlerMethods(handlerNode)
		if len(methods) > 0 {
			nodeInfo += fmt.Sprintf(" (Methods: %s)", strings.Join(methods, ", "))
		}
	}

	v.Builder.WriteString(nodeInfo + "\n")

	if container, ok := node.(Component.NodeContainer[Component.Node]); ok && container.HasChildren() {
		children := container.GetAllChildren()
		for i, child := range children {
			isChildLast := i == len(children)-1
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			v.VisualizeTree(child, newPrefix, isChildLast)
		}
	}

	return v.Builder.String()
}

// getNodePath 노드의 경로를 가져오는 헬퍼 함수
func (v *TreeVisualizer) getNodePath(node Component.Node) string {
	if pathNode, ok := node.(Component.PathNode); ok {
		return pathNode.GetPath()
	}
	return "/"
}

// getHandlerMethods 핸들러 노드의 메서드들을 가져오는 헬퍼 함수
func (v *TreeVisualizer) getHandlerMethods(handlerNode Component.HandlerNode) []string {
	methods := []string{}
	httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for _, method := range httpMethods {
		if handlerNode.HasMethod(method) {
			methods = append(methods, method)
		}
	}
	return methods
}

// VisualizeExpectedTree 예상되는 트리 구조를 시각화하는 함수
func (v *TreeVisualizer) VisualizeExpectedTree(expected ExpectedTreeNode, prefix string, isLast bool) string {
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	nodeInfo := fmt.Sprintf("%s%s[%s] %s", prefix, connector, NodeTypeToString(expected.Type), expected.Path)
	if len(expected.Methods) > 0 {
		nodeInfo += fmt.Sprintf(" (Methods: %s)", strings.Join(expected.Methods, ", "))
	}

	v.Builder.WriteString(nodeInfo + "\n")

	for i, child := range expected.Children {
		isChildLast := i == len(expected.Children)-1
		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		v.VisualizeExpectedTree(child, newPrefix, isChildLast)
	}

	return v.Builder.String()
}

// CompareTreeStructure 실제 트리와 예상 트리를 비교하는 함수
func CompareTreeStructure(t *testing.T, actual Component.Node, expected ExpectedTreeNode, path string) bool {
	if actual.GetType() != expected.Type {
		t.Errorf("경로 %s에서 노드 타입이 다릅니다. 예상: %s, 실제: %s", path, NodeTypeToString(expected.Type), NodeTypeToString(actual.GetType()))
		return false
	}

	if pathNode, ok := actual.(Component.PathNode); ok {
		if pathNode.GetPath() != expected.Path {
			t.Errorf("경로 %s에서 노드 경로가 다릅니다. 예상: %s, 실제: %s", path, expected.Path, pathNode.GetPath())
			return false
		}
	}

	if handlerNode, ok := actual.(Component.HandlerNode); ok {
		actualMethods := []string{}
		httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
		for _, method := range httpMethods {
			if handlerNode.HasMethod(method) {
				actualMethods = append(actualMethods, method)
			}
		}

		if (expected.Methods == nil || len(expected.Methods) == 0) && len(actualMethods) == 0 {
		} else if !reflect.DeepEqual(actualMethods, expected.Methods) {
			t.Errorf("경로 %s에서 메서드가 다릅니다. 예상: %v, 실제: %v", path, expected.Methods, actualMethods)
			return false
		}
	} else if len(expected.Methods) > 0 {
		t.Errorf("경로 %s에서 HandlerNode가 아닌데 메서드가 예상됨: %v", path, expected.Methods)
		return false
	}

	if container, ok := actual.(Component.NodeContainer[Component.Node]); ok {
		actualChildren := container.GetAllChildren()
		
		if len(actualChildren) != len(expected.Children) {
			t.Errorf("경로 %s에서 자식 노드 개수가 다릅니다. 예상: %d, 실제: %d", path, len(expected.Children), len(actualChildren))
			return false
		}

		for _, expectedChild := range expected.Children {
			found := false
			for _, actualChild := range actualChildren {
				if pathNode, ok := actualChild.(Component.PathNode); ok {
					if pathNode.GetPath() == expectedChild.Path {
						childPath := fmt.Sprintf("%s/%s", path, expectedChild.Path)
						if !CompareTreeStructure(t, actualChild, expectedChild, childPath) {
							return false
						}
						found = true
						break
					}
				}
			}
			if !found {
				t.Errorf("경로 %s에서 예상 자식 노드 '%s'를 찾을 수 없습니다", path, expectedChild.Path)
				return false
			}
		}
	}

	return true
}

// PrintNodeToString Tree_test.go에서 이동된 기존 시각화 함수
func PrintNodeToString(node Component.Node, prefix string, isLast bool) string {
	var result strings.Builder
	var nodeStr string
	var nodeTypeStr string
	var nodePath string

	switch n := node.(type) {
	case *Tree.RootNode:
		nodeTypeStr = "Root"
		nodePath = "/"
	case *Tree.StaticNode:
		nodeTypeStr = "Static"
		nodePath = n.GetPath()
	case *Tree.WildCardNode:
		nodeTypeStr = "Wildcard"
		nodePath = n.GetPath()
	case *Tree.CatchAllNode:
		nodeTypeStr = "CatchAll"
		nodePath = n.GetPath()
	case *Tree.MiddlewareNode:
		nodeTypeStr = "Middleware"
		nodePath = "(Middleware)"
	default:
		nodeTypeStr = "Unknown"
		if pn, ok := node.(Component.PathNode); ok {
			nodePath = pn.GetPath()
		} else {
			nodePath = "N/A"
		}
	}

	nodeStr = fmt.Sprintf("%s (%s)", nodePath, nodeTypeStr)
	result.WriteString(fmt.Sprintf("%s%s%s\n", prefix, GetBranchPrefix(isLast), nodeStr))

	if container, ok := node.(Component.NodeContainer[Component.Node]); ok {
		children := container.GetAllChildren()
		for i, child := range children {
			newPrefix := prefix + GetChildPrefix(isLast)
			result.WriteString(PrintNodeToString(child, newPrefix, i == len(children)-1))
		}
	}

	return result.String()
}

// GetBranchPrefix Tree_test.go에서 이동된 함수
func GetBranchPrefix(isLast bool) string {
	if isLast {
		return "└── "
	}
	return "├── "
}

// GetChildPrefix Tree_test.go에서 이동된 함수
func GetChildPrefix(isLast bool) string {
	if isLast {
		return "    "
	}
	return "│   "
}

// PrintTreeStructure Tree_test.go에서 이동된 함수
func PrintTreeStructure(tree Tree.Tree) string {
	var buffer bytes.Buffer
	buffer.WriteString("LiteFrame Router Tree:\n")
	buffer.WriteString(PrintNodeToString(&tree.Root, "", true))
	return buffer.String()
}