package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-ddns-updater/model"
	"github.com/mitoteam/mttools"
	"github.com/mitoteam/mtweb"
)

var formWebhookEdit = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		wh := fd.GetArg("Webhook").(*model.DdnsWebhook)

		optionsSelect := dhtmlbs.NewSelect("key").
			Option("0", "[not selected]")

		for _, key := range goapp.LoadOL[model.DdnsKey]() {
			optionsSelect.Option(mttools.AnyToString(key.ID), key.Name)
		}

		formBody.Append(
			dhtml.Div().Class("border bg-light p-3").Append(
				dhtmlbs.NewTextInput("record_name").
					Label("DNS Record Name").
					Default(wh.DnsRecordName).
					Require(),
				dhtmlbs.NewTextInput("zone_name").
					Label("DNS Zone").
					Default(wh.DnsZoneName).
					Require(),
				dhtmlbs.NewTextInput("server_address").
					Label("DNS Server Address and Port").
					Default(wh.DnsServerAddress).
					Note("example: 8.8.8.8:53 or ns2.wildcowdns.net:53").
					Require(),
				optionsSelect.Label("Key").Default(wh.KeyID).Require(),
				dhtmlbs.NewTextarea("description").Label("Description").
					Default(wh.Description).Note("Description for webhook"),
				mtweb.NewDefaultSubmitBtn(),
			),
		)
	},

	SubmitF: func(fd *dhtmlform.FormData) {
		wh := fd.GetArg("Webhook").(*model.DdnsWebhook)

		wh.DnsRecordName = fd.GetValue("record_name").(string)
		wh.DnsZoneName = fd.GetValue("zone_name").(string)
		wh.DnsServerAddress = fd.GetValue("server_address").(string)

		if key_id, ok := mttools.AnyToInt64Ok(fd.GetValue("key")); ok {
			wh.KeyID = key_id
		}

		wh.Description = fd.GetValue("description").(string)

		goapp.SaveObject(wh)
	},
}
