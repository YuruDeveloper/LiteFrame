package Param

import (
	"context"
)

func NewParams() Params {
	return Params {
		List: make([]Param, 0),
	}
}

type Param struct {
	Key   string
	Value string
}

type Params struct {
	List []Param
}

type Key struct{ }

func (Instance *Params) Add(Key string, Value string) {
	Instance.List = append(Instance.List, Param{Key: Key, Value: Value})
}

func (Instance *Params) GetByName(Name string) string {
	for _, Param := range Instance.List {
		if Param.Key == Name {
			return Param.Value
		}
	}
	return ""
}

func GetParamsFromCTX(Context context.Context) (Params, bool) {
	Temp, Susscecs := (Context.Value(Key{})).(Params)
	return Temp, Susscecs
}
