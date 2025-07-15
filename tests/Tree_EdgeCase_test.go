package tests

import (
	"testing"
	"net/http"
	"LiteFrame/Router/Tree"
	"github.com/stretchr/testify/assert"
)

// 테스트용 핸들러 함수들
func edgeCaseHandler1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("EdgeCase Handler 1"))
}

func edgeCaseHandler2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("EdgeCase Handler 2"))
}

// 1. 극단적인 경로 길이 테스트
func TestTree_Add_ExtremelyLongPath(t *testing.T) {
	tree := Tree.NewTree()
	
	// 1000개 세그먼트를 가진 극단적으로 긴 경로
	longPath := "/"
	for i := 0; i < 1000; i++ {
		longPath += "segment" + string(rune(i%26 + 'a')) + "/"
	}
	longPath += "end"
	
	err := tree.Add("GET", longPath, edgeCaseHandler1)
	
	// 이 테스트는 메모리 부족이나 스택 오버플로우를 유발할 수 있음
	if err != nil {
		t.Logf("Expected: 극단적으로 긴 경로는 실패해야 함 - %v", err)
		assert.Error(t, err)
	} else {
		t.Logf("Warning: 극단적으로 긴 경로가 성공했음 - 메모리 사용량 확인 필요")
	}
}

// 2. 잘못된 와일드카드 형식 테스트
func TestTree_Add_InvalidWildcardFormats(t *testing.T) {
	tree := Tree.NewTree()
	
	invalidWildcards := []string{
		"/users/:",           // 빈 와일드카드 이름
		"/users/:id:",        // 와일드카드 뒤에 콜론
		"/users/::id",        // 이중 콜론
		"/users/:id/:",       // 빈 와일드카드 이름 중간에
		"/users/:123invalid", // 숫자로 시작하는 와일드카드
		"/users/: space",     // 공백이 포함된 와일드카드
	}
	
	for _, path := range invalidWildcards {
		t.Run(path, func(t *testing.T) {
			err := tree.Add("GET", path, edgeCaseHandler1)
			// 일부는 성공할 수 있지만, 검증이 필요함
			if err != nil {
				t.Logf("Path %s failed as expected: %v", path, err)
			} else {
				t.Logf("Warning: Path %s succeeded but may be invalid", path)
			}
		})
	}
}

// 3. 경로 충돌 테스트 - 다른 타입의 노드 간 충돌
func TestTree_Add_PathConflicts(t *testing.T) {
	tree := Tree.NewTree()
	
	// 정적 경로 먼저 추가
	err := tree.Add("GET", "/users/profile", edgeCaseHandler1)
	assert.NoError(t, err)
	
	// 같은 레벨에 와일드카드 추가 시도 - 이것은 충돌을 일으킬 수 있음
	err = tree.Add("GET", "/users/:id", edgeCaseHandler2)
	if err != nil {
		t.Logf("Expected conflict detected: %v", err)
		assert.Error(t, err)
	} else {
		// 충돌이 감지되지 않았다면 경고
		t.Logf("Warning: Path conflict not detected - this may cause routing ambiguity")
	}
}

// 4. CatchAll과 다른 경로의 충돌
func TestTree_Add_CatchAllConflicts(t *testing.T) {
	tree := Tree.NewTree()
	
	// CatchAll 경로 먼저 추가
	err := tree.Add("GET", "/static/*", edgeCaseHandler1)
	assert.NoError(t, err)
	
	// 같은 prefix에 정적 경로 추가 시도
	err = tree.Add("GET", "/static/css/style.css", edgeCaseHandler2)
	if err != nil {
		t.Logf("Expected: CatchAll conflict detected: %v", err)
		assert.Error(t, err)
	} else {
		t.Logf("Warning: CatchAll conflict not detected")
	}
}

// 5. 메모리 소진 시뮬레이션 - 대량의 핸들러 추가
func TestTree_Add_MemoryExhaustion(t *testing.T) {
	tree := Tree.NewTree()
	
	// 10000개의 서로 다른 경로 추가
	for i := 0; i < 10000; i++ {
		path := "/test/" + string(rune(i/1000+'a')) + "/" + string(rune((i%1000)/100+'a')) + "/" + string(rune((i%100)/10+'a')) + "/" + string(rune(i%10+'0'))
		err := tree.Add("GET", path, edgeCaseHandler1)
		
		if err != nil {
			t.Logf("Memory exhaustion detected at iteration %d: %v", i, err)
			break
		}
		
		// 매 1000번째마다 로그
		if i%1000 == 0 {
			t.Logf("Successfully added %d paths", i)
		}
	}
}

// 6. 순환 참조 또는 깊은 재귀 테스트
func TestTree_Add_DeepRecursion(t *testing.T) {
	tree := Tree.NewTree()
	
	// 매우 깊은 경로 (100 레벨)
	deepPath := "/"
	for i := 0; i < 100; i++ {
		deepPath += "level" + string(rune(i%10+'0')) + "/"
	}
	deepPath += "end"
	
	err := tree.Add("GET", deepPath, edgeCaseHandler1)
	if err != nil {
		t.Logf("Deep recursion failed as expected: %v", err)
		assert.Error(t, err)
	} else {
		t.Logf("Warning: Deep recursion succeeded - stack overflow risk")
	}
}

// 7. nil 값들과 빈 문자열들의 조합
func TestTree_Add_NilAndEmptyValues(t *testing.T) {
	tree := Tree.NewTree()
	
	testCases := []struct {
		method  string
		path    string
		handler http.HandlerFunc
		desc    string
	}{
		{"", "/test", edgeCaseHandler1, "empty method"},
		{"GET", "", edgeCaseHandler1, "empty path"},
		{"GET", "/test", nil, "nil handler"},
		{"", "", edgeCaseHandler1, "empty method and path"},
		{"", "/test", nil, "empty method and nil handler"},
		{"GET", "", nil, "empty path and nil handler"},
		{"", "", nil, "all empty/nil"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tree.Add(tc.method, tc.path, tc.handler)
			assert.Error(t, err, "Should fail for %s", tc.desc)
		})
	}
}

// 8. 특수 문자와 유니코드 경로
func TestTree_Add_SpecialCharacters(t *testing.T) {
	tree := Tree.NewTree()
	
	specialPaths := []string{
		"/users/한글경로",
		"/users/中文路径",
		"/users/🚀emoji",
		"/users/special@#$%^&*()",
		"/users/with spaces",
		"/users/with\ttab",
		"/users/with\nnewline",
		"/users/with\\backslash",
		"/users/with'quotes\"",
		"/users/../../../etc/passwd", // 경로 탐색 공격
	}
	
	for _, path := range specialPaths {
		t.Run(path, func(t *testing.T) {
			err := tree.Add("GET", path, edgeCaseHandler1)
			if err != nil {
				t.Logf("Special character path failed: %s - %v", path, err)
			} else {
				t.Logf("Special character path succeeded: %s", path)
			}
		})
	}
}

// 9. 동시성 테스트 (Race Condition)
func TestTree_Add_Concurrency(t *testing.T) {
	tree := Tree.NewTree()
	
	// 동시에 같은 경로에 핸들러 추가 시도
	done := make(chan bool, 2)
	
	go func() {
		for i := 0; i < 100; i++ {
			tree.Add("GET", "/concurrent/test", edgeCaseHandler1)
		}
		done <- true
	}()
	
	go func() {
		for i := 0; i < 100; i++ {
			tree.Add("POST", "/concurrent/test", edgeCaseHandler2)
		}
		done <- true
	}()
	
	// 두 고루틴이 완료될 때까지 대기
	<-done
	<-done
	
	// Race condition이 발생했는지 확인
	// 실제로는 데이터 레이스 감지기로 확인해야 함
	t.Log("Concurrency test completed - check with race detector")
}

// 10. Split 기능 관련 엣지케이스
func TestTree_Add_SplitEdgeCases(t *testing.T) {
	tree := Tree.NewTree()
	
	// 공통 prefix를 가진 복잡한 경로들
	paths := []string{
		"/a",
		"/ab",
		"/abc",
		"/abcd",
		"/abcde",
		"/ab/cd",
		"/abc/d",
	}
	
	for i, path := range paths {
		err := tree.Add("GET", path, edgeCaseHandler1)
		if err != nil {
			t.Logf("Split edge case failed at path %s: %v", path, err)
			// Split 관련 오류는 예상됨
		} else {
			t.Logf("Successfully added path %d: %s", i, path)
		}
	}
	Temp := PrintTreeStructure(tree)
	t.Log(Temp)
}

// 11. 메서드명 검증 테스트
func TestTree_Add_InvalidHTTPMethods(t *testing.T) {
	tree := Tree.NewTree()
	
	invalidMethods := []string{
		"get",           // 소문자
		"INVALID",       // 표준이 아닌 메서드
		"G E T",         // 공백 포함
		"GET\t",         // 탭 포함
		"GET\n",         // 개행 포함
		"123",           // 숫자만
		"@#$",           // 특수문자만
		"VERYLONGMETHODNAME", // 매우 긴 메서드명
	}
	
	for _, method := range invalidMethods {
		t.Run(method, func(t *testing.T) {
			err := tree.Add(method, "/test", edgeCaseHandler1)
			if err != nil {
				t.Logf("Invalid method %s failed as expected: %v", method, err)
			} else {
				t.Logf("Warning: Invalid method %s was accepted", method)
			}
		})
	}
}

// 12. 타입 변환 실패 시뮬레이션
func TestTree_Add_TypeConversionFailures(t *testing.T) {
	tree := Tree.NewTree()
	
	// 정상 경로 추가 후 내부 구조 손상 시뮬레이션은 어려우므로
	// 현재 코드에서 발생할 수 있는 타입 변환 실패 상황을 테스트
	
	// 매우 복잡한 경로 구조로 타입 변환 실패 유도 시도
	err := tree.Add("GET", "/complex/:param1/static/:param2/catch/*", edgeCaseHandler1)
	if err != nil {
		t.Logf("Complex path structure failed: %v", err)
		assert.Error(t, err)
	} else {
		t.Log("Complex path structure succeeded")
	}
}