package Parm

import (
	"context"
)

type Parm struct {
	Key   string
	Value string
}

type Parms struct {
	List []Parm
}

type ParmKey struct{}

func (Instance *Parms) Add(Key string,Value string) {
	Instance.List = append(Instance.List, Parm{Key: Key,Value: Value})
}

func (Instance *Parms) GetByName(Name string) string{
		for _ , Parm := range Instance.List {
			if Parm.Key == Name {
				return Parm.Value
			}
		}
		return ""
}

func GetParmsFromCTX(Context context.Context) (Parms,bool) {
	Temp , Susscecs:= (Context.Value(ParmKey{})).(Parms)
	return Temp , Susscecs
}