package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-ddns-updater/model"
	"github.com/mitoteam/mtweb"
)

var formWebhookEdit = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		wh := fd.GetArg("Webhook").(*model.DdnsWebhook)

		formBody.Append(
			dhtml.Div().Class("border bg-light p-3").Append(
				dhtmlbs.NewTextInput("record_name").
					Label("DNS Record Name").
					Default(wh.DnsRecordName).
					Require(),
				dhtmlbs.NewTextInput("zone_name").
					Label("DNS Name Name").
					Default(wh.DnsZoneName).
					Require(),
				// dhtmlbs.NewSelect("algo").
				// 	OptionsMap(model.DdnsKeyAlgoOptions()).
				// 	Label("Algorithm").
				// 	Default(key.Algo).
				// 	Require(),
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

		// if algo, ok := mttools.AnyToIntOk(fd.GetValue("algo")); ok {
		// 	key.Algo = algo
		// }
		wh.Description = fd.GetValue("description").(string)

		goapp.SaveObject(wh)
	},
}
