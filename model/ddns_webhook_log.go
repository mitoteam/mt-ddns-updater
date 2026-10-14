package model

import (
	"fmt"
	"reflect"

	"github.com/mitoteam/goapp"
)

type LogKindType int

const (
	LogKindUnknown     LogKindType = iota // 0
	LogKindQuestionGui                    // 1
	LogKindUpdateGui                      // 2
)

var logKindMap = map[LogKindType]string{
	LogKindQuestionGui: "Question [GUI]",
	LogKindUpdateGui:   "Update [GUI]",
}

func (v LogKindType) String() string {
	if s, ok := logKindMap[v]; ok {
		return s
	} else {
		return fmt.Sprintf("Unknown log kind(%d)", v)
	}

}

type DdnsWebhookLog struct {
	goapp.BaseModel

	// using built-in CreatedAt for timestamps

	//fk
	WebhookID int64        `gorm:"not null;index"`
	Webhook   *DdnsWebhook `gorm:"->;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`

	Kind    LogKindType `gorm:"type:integer;not null"`
	Message string
	Reply   string
	Error   string
	Ip      string
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[DdnsWebhookLog]())
}

func NewWebhookLog(webhook *DdnsWebhook, kind LogKindType) *DdnsWebhookLog {
	log := DdnsWebhookLog{
		WebhookID: webhook.ID,
		Kind:      kind,
	}

	return &log
}

func (o *DdnsWebhookLog) GetWebhook() *DdnsWebhook {
	if o.Webhook == nil {
		o.Webhook = goapp.LoadOMust[DdnsWebhook](o.WebhookID)
	}

	return o.Webhook
}
