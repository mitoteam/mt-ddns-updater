package model

import (
	"net"
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
	ActualIp         string
	Ttl              uint32

	//fk
	KeyID int64    `gorm:"not null;index"`
	Key   *DdnsKey `gorm:"->;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`

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

func (wh *DdnsWebhook) BeforeSave(tx *gorm.DB) (err error) {
	if wh.SecureToken == "" {
		wh.SecureToken = mttools.RandomString(64)
	}

	if wh.ActualIp == "" {
		wh.ActualIp = "0.0.0.0"
	}

	return nil //no error
}

func (wh *DdnsWebhook) FQDN() string {
	return dns.Fqdn(wh.DnsRecordName) + dns.Fqdn(wh.DnsZoneName)
}

func (wh *DdnsWebhook) GetParsedActualIp() net.IP {
	ip := net.ParseIP(wh.ActualIp)

	if ip == nil {
		ip = net.ParseIP("0.0.0.0")
	}

	return ip
}
