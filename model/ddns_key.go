package model

import (
	"reflect"

	"github.com/mitoteam/goapp"
)

type DdnsKey struct {
	goapp.BaseModel

	Name   string // key name
	Algo   string
	Secret string

	Description string
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[DdnsKey]())
}
