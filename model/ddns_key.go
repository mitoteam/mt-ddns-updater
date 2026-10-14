package model

import (
	"fmt"
	"reflect"

	"github.com/miekg/dns"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mttools"
)

type DdnsKey struct {
	goapp.BaseModel

	Name   string      // key name
	Algo   KeyAlgoType // key algorithm
	Secret string      // secret in base64

	Description string
}

type KeyAlgoType int

const (
	AlgoTypeUnknown KeyAlgoType = iota // 0
	KeyAlgoSHA1                        // 1
	KeyAlgoSHA224                      // 2
	KeyAlgoSHA256                      // 3
	KeyAlgoSHA384                      // 4
	KeyAlgoSHA512                      // 5
	KeyAlgoMD5                         // 6
)

var keyAlgoMap = map[KeyAlgoType]string{
	KeyAlgoSHA1:   dns.HmacSHA1,
	KeyAlgoSHA224: dns.HmacSHA224,
	KeyAlgoSHA256: dns.HmacSHA256,
	KeyAlgoSHA384: dns.HmacSHA384,
	KeyAlgoSHA512: dns.HmacSHA512,
	KeyAlgoMD5:    dns.HmacMD5,
}

func (v KeyAlgoType) String() string {
	if s, ok := keyAlgoMap[v]; ok {
		return s
	} else {
		return fmt.Sprintf("Unknown key algo(%d)", v)
	}
}

// options for select form element
func DdnsKeyAlgoOptions() (r map[string]any) {
	r = make(map[string]any, len(keyAlgoMap))

	for key, value := range keyAlgoMap {
		r[mttools.AnyToString(key)] = value
	}

	return r
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[DdnsKey]())
}

func (k *DdnsKey) WebhooksCount() int64 {
	goapp.PreQuery[DdnsWebhook]().Where("key_id", k.ID)

	return goapp.CountOL[DdnsWebhook]()
}
