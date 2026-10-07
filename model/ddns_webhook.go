package model

import (
	"reflect"

	"github.com/mitoteam/goapp"
)

type DdnsWebhook struct {
	goapp.BaseModel

	Description string
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[DdnsWebhook]())
}
