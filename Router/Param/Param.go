// Package Param은 HTTP 라우팅에서 사용되는 매개변수 관리 기능을 제공합니다.
// 성능 최적화를 위한 매개변수 풀링과 컨텍스트 기반 매개변수 전달을 지원합니다.
package Param

import (
	"context"
	"sync"
)

// NewParams는 새로운 Params 인스턴스를 생성합니다.
// 기본 용량 2로 매개변수 리스트를 초기화합니다.
// 성능 최적화: 대부분의 라우트는 2개 이하의 매개변수를 가지므로 메모리 재할당을 최소화합니다.
func NewParams() *Params {
	return &Params{
		List: make([]Param, 0, 2), // 기본 용량 2: 일반적인 API 패턴에 최적화된 값
	}
}

// Param은 단일 매개변수의 키-값 쌍을 나타내는 구조체입니다.
// URL 경로에서 추출된 매개변수 이름과 값을 저장합니다.
type Param struct {
	Key   string // 매개변수 이름 (:에서 추출된 이름)
	Value string // 매개변수 값 (URL에서 추출된 실제 값)
}

// Params는 여러 매개변수를 저장하는 구조체입니다.
// 한 번의 HTTP 요청에서 추출된 모든 매개변수들을 관리합니다.
type Params struct {
	List []Param // 매개변수 리스트
}

// Key는 컨텍스트에서 매개변수를 식별하기 위한 빈 구조체입니다.
// context.WithValue에서 매개변수 저장을 위한 키로 사용됩니다.
type Key struct{}

// Add는 매개변수 리스트에 새로운 매개변수를 추가합니다.
// URL 경로에서 추출된 매개변수 이름과 값을 저장합니다.
func (Instance *Params) Add(Key string, Value string) {
	Instance.List = append(Instance.List, Param{Key: Key, Value: Value})
}

// GetByName은 매개변수 이름으로 해당하는 값을 검색합니다.
// 매개변수가 찾지 못하면 빈 문자열을 반환합니다.
func (Instance *Params) GetByName(Name string) string {
	for _, Param := range Instance.List {
		if Param.Key == Name {
			return Param.Value
		}
	}
	return ""
}

// GetParamsFromCTX는 컨텍스트에서 매개변수를 추출합니다.
// HTTP 핸들러에서 요청의 컨텍스트로부터 매개변수를 가져옵니다.
func GetParamsFromCTX(Context context.Context) (*Params, bool) {
	Temp, Success := (Context.Value(Key{})).(*Params)
	return Temp, Success
}

// NewParamsPool은 새로운 매개변수 풀을 생성합니다.
// 성능 최적화를 위해 sync.Pool을 사용하여 매개변수 객체를 재사용합니다.
// 초기에 10개의 매개변수 객체를 미리 할당하여 풀에 넣어둡니다.
// 웜업 예열(Warm-up): 초기 요청에서의 메모리 할당 지연을 방지합니다.
func NewParamsPool() *ParamsPool {
	Instance :=  &ParamsPool{
		Pool : &sync.Pool {
			// Factory 함수: 풀이 비어있을 때 새로운 객체를 생성
			New: func() any {
				return NewParams()
			}, 
		},
	}
	// 수동 웜업: 10개 객체를 미리 생성하여 첫 번째 요청들의 레이턴시 감소
	for Index := 0; Index < 10  ; Index ++ {
		Instance.Put(NewParams())
	}
	return Instance
}


// ParamsPool은 Params 객체를 재사용하기 위한 풀 구조체입니다.
// 메모리 할당과 가비지 컨렉션 오버헤드를 줄이기 위해 sync.Pool을 사용합니다.
type ParamsPool struct {
	Pool *sync.Pool // 매개변수 객체 풀
}

// Get은 풀에서 매개변수 객체를 가져옵니다.
// 기존 매개변수 리스트를 초기화하여 새로운 요청에서 사용할 수 있도록 준비합니다.
// 매모리 효율성: 용량은 유지하고 길이만 0으로 리셋하여 재할당 없이 재사용
func (Instance *ParamsPool) Get() *Params {
	Object := Instance.Pool.Get().(*Params)
	// 슬라이스 리셋: 기존 메모리를 유지하면서 길이만 0으로 설정 (성능 최적화)
	Object.List = Object.List[0:0]
	return Object
}

// Put은 사용이 끝난 매개변수 객체를 풀에 반납합니다.
// 메모리 누수를 방지하기 위해 리스트 용량이 8을 초과하면 새로운 슬라이스를 생성합니다.
// 중요한 메모리 관리: 과도한 용량 증가를 방지하여 메모리 풀의 효율성 유지
func (Instance *ParamsPool) Put(Object *Params) {
	if Object != nil {
		// 메모리 방블론 방지: 용량이 임계값(8)을 초과하면 새 슬라이스 생성
		// 이는 대량의 매개변수를 가진 요청 후 메모리가 계속 증가하는 것을 방지
		if cap(Object.List) > 8 {
			Object.List = make([]Param, 0,2) // 기본 용량으로 리셋
		}
		Instance.Pool.Put(Object)
	}
}