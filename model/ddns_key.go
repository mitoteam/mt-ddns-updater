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

func (v KeyAlgoType) Name() string {
	switch v {
	case KeyAlgoSHA1:
		return dns.HmacSHA1
	case KeyAlgoSHA224:
		return dns.HmacSHA224
	case KeyAlgoSHA256:
		return dns.HmacSHA256
	case KeyAlgoSHA384:
		return dns.HmacSHA384
	case KeyAlgoSHA512:
		return dns.HmacSHA512
	case KeyAlgoMD5:
		return dns.HmacMD5
	default:
		return fmt.Sprintf("Unknown key algo(%d)", v)
	}
}

var keyAlgoOptions mttools.Values

// options for select form element
func DdnsKeyAlgoOptions() *mttools.Values {
	return &keyAlgoOptions
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[DdnsKey]())

	//all algos list
	keyAlgoList := []KeyAlgoType{
		KeyAlgoSHA1,
		KeyAlgoSHA224,
		KeyAlgoSHA256,
		KeyAlgoSHA384,
		KeyAlgoSHA512,
		KeyAlgoMD5,
	}

	keyAlgoOptions = mttools.NewValues()

	for _, v := range keyAlgoList {
		keyAlgoOptions.Set(v, v.Name())
	}
}

func (k *DdnsKey) WebhooksCount() int64 {
	goapp.PreQuery[DdnsWebhook]().Where("key_id", k.ID)

	return goapp.CountOL[DdnsWebhook]()
}
