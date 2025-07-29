// Package Tree는 효율적인 경로 처리를 위한 PathWithSegment 구조체를 제공합니다.
// 메모리 할당 없이 경로를 세그먼트 단위로 순회할 수 있는 구조체입니다.
package Tree

// NewPathWithSegment는 새로운 PathWithSegment 인스턴스를 생성합니다.
// Path: 분석할 URL 경로 문자열
// 초기 상태에서는 Start와 End가 모두 0으로 설정됩니다.
func NewPathWithSegment(Path string) *PathWithSegment {
	return &PathWithSegment{
		Body: Path,
		Start:  0,
		End: 0,
	}
}

// PathWithSegment는 URL 경로를 메모리 할당 없이 세그먼트 단위로 처리하는 구조체입니다.
// 문자열을 복사하지 않고 인덱스를 이용하여 경로 세그먼트를 순회합니다.
//
// 성능 최적화:
// - 제로 할당: 새로운 문자열을 생성하지 않고 기존 문자열의 부분을 참조
// - 반복자 패턴: Next()를 통해 순차적으로 세그먼트 이동
// - 경계 검사: 안전한 인덱스 접근을 위한 검증 함수들 제공
type PathWithSegment struct {
	Body string    // 원본 경로 문자열 (불변)
	Start int      // 현재 세그먼트의 시작 인덱스
	End int        // 현재 세그먼트의 끝 인덱스 (exclusive)
}

// Next는 다음 경로 세그먼트로 이동합니다.
// 경로 구분자('/')를 건너뛰고 다음 세그먼트의 시작과 끝 인덱스를 설정합니다.
// 
// 동작 방식:
// 1. 현재 End 위치를 새로운 Start로 설정
// 2. 연속된 '/' 문자들을 건너뛰기
// 3. 다음 '/' 또는 문자열 끝까지를 새로운 세그먼트로 설정
func (Instance *PathWithSegment) Next() {
	Instance.Start = Instance.End
	if Instance.IsEnd() {
		return
	}

	// 연속된 경로 구분자('/')들을 건너뛰기
	for len(Instance.Body) > Instance.Start && Instance.Body[Instance.Start] == '/' {
		Instance.Start++
	}
	if Instance.IsEnd() {
		Instance.End = Instance.Start
		return
	}
	// 다음 경로 구분자까지 또는 문자열 끝까지 세그먼트 설정
	Instance.End = Instance.Start
	for Instance.End < len(Instance.Body) && Instance.Body[Instance.End] != PathSeparator {
		Instance.End++
	}
}

// IsEnd는 경로의 끝에 도달했는지 확인합니다.
// Start 인덱스가 문자열 길이와 같거나 클 때 true를 반환합니다.
func (Instance *PathWithSegment) IsEnd() bool {
	return  !(Instance.Start < len((Instance.Body)))
}

// IsSame은 현재 세그먼트가 빈 세그먼트인지 확인합니다.
// Start와 End가 같을 때 true를 반환합니다. (길이가 0인 세그먼트)
func (Instance *PathWithSegment) IsSame() bool {
	return Instance.Start == Instance.End
}

// IsNotVaild는 현재 인덱스 상태가 유효하지 않은지 확인합니다.
// 경로 끝 도달, 인덱스 범위 초과, 논리적 오류를 검사합니다.
func (Instance *PathWithSegment) IsNotVaild() bool {
	return Instance.IsEnd() || Instance.End > len(Instance.Body) || Instance.Start > Instance.End
}

// Get은 현재 세그먼트의 문자열을 반환합니다.
// 유효하지 않은 상태일 경우 빈 문자열을 반환합니다.
// 메모리 할당: 새로운 문자열 생성 (필요한 경우에만 호출)
func (Instance *PathWithSegment) Get() string {
	if Instance.IsNotVaild() {
        return ""
    }
	return string(Instance.Body[Instance.Start:Instance.End])
}

// GetToEnd는 현재 위치부터 경로 끝까지의 문자열을 반환합니다.
// CatchAll 라우트에서 나머지 경로를 모두 캡처할 때 사용됩니다.
func (Instance *PathWithSegment) GetToEnd() string {
	return string(Instance.Body[Instance.Start:])
}

// GetLength는 현재 세그먼트의 길이를 반환합니다.
// 문자열 생성 없이 길이만 계산하므로 메모리 효율적입니다.
func (Instance *PathWithSegment) GetLength() int {
	return Instance.End - Instance.Start
}
