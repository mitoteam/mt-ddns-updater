package model

import (
	"reflect"

	"github.com/mitoteam/goapp"
)

type DdnsWebhook struct {
	goapp.BaseModel

	SecureToken string

	DnsServerAddress string
	DnsRecordName    string
	DnsZoneName      string

	//fk
	KeyID int64    `gorm:"not null;index"`
	Key   *DdnsKey `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`

	Description string
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[DdnsWebhook]())
}

func (item *DdnsWebhook) GetKey() *DdnsKey {
	if item.Key == nil {
		item.Key = goapp.LoadOMust[DdnsKey](item.KeyID)
	}

	return item.Key
}
