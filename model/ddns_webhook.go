package model

import (
	"reflect"

	"github.com/miekg/dns"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mttools"
	"gorm.io/gorm"
)

type DdnsWebhook struct {
	goapp.BaseModel

	SecureToken string

	DnsServerAddress string
	DnsRecordName    string
	DnsZoneName      string
	InitialIp        string

	//fk
	KeyID int64    `gorm:"not null;index"`
	Key   *DdnsKey `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`

	Description string
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[DdnsWebhook]())
}

func (wh *DdnsWebhook) GetKey() *DdnsKey {
	if wh.Key == nil {
		wh.Key = goapp.LoadOMust[DdnsKey](wh.KeyID)
	}

	return wh.Key
}

func (wh *DdnsWebhook) FQDN() string {
	return dns.Fqdn(wh.DnsRecordName) + dns.Fqdn(wh.DnsZoneName)
}

func (wh *DdnsWebhook) BeforeSave(tx *gorm.DB) (err error) {
	if wh.SecureToken == "" {
		wh.SecureToken = mttools.RandomString(64)
	}
	return
}
