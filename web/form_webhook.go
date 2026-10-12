package web

import (
	"net"
	"strings"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-ddns-updater/app"
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
				dhtmlbs.NewTextInput("actual_ip").
					Label("Actual IP").
					Default(wh.ActualIp).
					Note("example: 127.0.0.1, leave empty to use 0.0.0.0"),
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

	ValidateF: func(fd *dhtmlform.FormData) {
		ipString := strings.TrimSpace(fd.GetValue("actual_ip").(string))

		if ipString == "" {
			ipString = "0.0.0.0"
			fd.SetControlValue("actual_ip", ipString) // set it back
		}

		if ip := net.ParseIP(ipString); ip == nil {
			fd.SetError("actual_ip", "Incorrect IP address value")
		}
	},

	SubmitF: func(fd *dhtmlform.FormData) {
		wh := fd.GetArg("Webhook").(*model.DdnsWebhook)

		wh.DnsRecordName = fd.GetValue("record_name").(string)
		wh.DnsZoneName = fd.GetValue("zone_name").(string)
		wh.ActualIp = fd.GetValue("actual_ip").(string)
		wh.DnsServerAddress = fd.GetValue("server_address").(string)

		if key_id, ok := mttools.AnyToInt64Ok(fd.GetValue("key")); ok {
			wh.KeyID = key_id
		}

		wh.Description = fd.GetValue("description").(string)

		goapp.SaveObject(wh)
	},
}

var formWebhookTestUpdate = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		wh := fd.GetArg("Webhook").(*model.DdnsWebhook)

		piece := dhtml.NewHtmlPiece()

		helperI := fd.GetParam("helper")

		var err error
		if errorI := fd.GetParam("error"); errorI != nil {
			err = errorI.(error)
		}

		if helperI != nil {
			// there is helper set, show results
			helper := helperI.(*app.WebhookService)
			var msg *dhtml.HtmlPiece

			if err == nil {
				msg = mtweb.RenderInfof("A-Record for %s updated successfully to %s", wh.FQDN(), fd.GetValue("ip").(string))
			} else {
				msg = mtweb.RenderErrorf("Error during DNS update: %v", err)
			}

			piece.
				Append(dhtml.Div().Class("fw-bold fs-4").Append("Message")).
				Append(dhtml.NewTag("pre").Append(helper.GetLastMessage())).
				Append(dhtml.Div().Class("fw-bold fs-4").Append("Reply")).
				Append(msg).
				Append(dhtml.NewTag("pre").Append(helper.GetLastReply()))
		} else if err != nil {
			// there is no helper but error is set, show it
			piece.Append(mtweb.RenderError("Error sending DNS update message: " + err.Error()))
		} else {
			//first run, render the form
			piece.Append(
				mtweb.RenderInfof("DNS Record: %s", wh.FQDN()),
				dhtmlbs.NewTextInput("ip").
					Label("IP address to update record to").
					Require().
					Note("example: 127.0.0.1"),
				mtweb.NewDefaultSubmitBtn(),
			)
		}

		formBody.Append(
			dhtml.Div().Class("border bg-light p-3").Append(piece),
		)
	},

	ValidateF: func(fd *dhtmlform.FormData) {
		ipString := strings.TrimSpace(fd.GetValue("ip").(string))

		if ip := net.ParseIP(ipString); ip == nil {
			fd.SetError("ip", "Incorrect IP address value")
		}
	},

	SubmitF: func(fd *dhtmlform.FormData) {
		wh := fd.GetArg("Webhook").(*model.DdnsWebhook)

		ipString := fd.GetValue("ip").(string)
		ip := net.ParseIP(ipString)

		helper := app.NewWebhookService(wh)
		fd.SetParam("helper", helper)

		err := helper.UpdateIpAddress(ip)

		if err == nil {
			wh.ActualIp = ip.String()
			goapp.SaveObject(wh)
		} else {
			fd.SetParam("error", err)
		}

		fd.SetRebuild(true)
	},
}
