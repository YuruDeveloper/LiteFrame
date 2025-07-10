package Parm

import (
	"context"
)

type Parm struct {
	Key   string
	Value string
}

type Parms []Parm

type ParmKey struct{}

func (Instance *Parms) GetByName(Name string) string{
		for _ , Parm := range *Instance {
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