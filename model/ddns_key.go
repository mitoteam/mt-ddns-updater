package model

import (
	"reflect"

	"github.com/miekg/dns"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mttools"
)

type DdnsKey struct {
	goapp.BaseModel

	Name   string // key name
	Algo   int    // key algorithm (one of keyAlgoList)
	Secret string // secret in base64

	Description string
}

var keyAlgoList map[int]string

const (
	Unknown       int = iota // 0
	KeyAlgoSHA1              // 1
	KeyAlgoSHA224            // 2
	KeyAlgoSHA256            // 3
	KeyAlgoSHA384            // 4
	KeyAlgoSHA512            // 5
	KeyAlgoMD5               // 6
)

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[DdnsKey]())

	//initialize map
	keyAlgoList = map[int]string{
		KeyAlgoSHA1:   dns.HmacSHA1,
		KeyAlgoSHA224: dns.HmacSHA224,
		KeyAlgoSHA256: dns.HmacSHA256,
		KeyAlgoSHA384: dns.HmacSHA384,
		KeyAlgoSHA512: dns.HmacSHA512,
		KeyAlgoMD5:    dns.HmacMD5,
	}
}

// options for select form element
func DdnsKeyAlgoOptions() (r map[string]any) {
	r = make(map[string]any, len(keyAlgoList))

	for key, value := range keyAlgoList {
		r[mttools.AnyToString(key)] = value
	}

	return r
}

func (k *DdnsKey) GetAlgoName() string {
	return keyAlgoList[k.Algo]
}

func (k *DdnsKey) WebhooksCount() int64 {
	goapp.PreQuery[DdnsWebhook]().Where("key_id", k.ID)

	return goapp.CountOL[DdnsWebhook]()
}
