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
				dhtmlbs.NewTextInput("ttl").Label("TTL").Default(wh.Ttl).Require(),
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

		ttl := mttools.AnyToIntOrZero(fd.GetValue("ttl"))

		if ttl < 1 || ttl > 86400 {
			fd.SetError("ttl", "TTL should be between 1 and 86400")
		} else {
			fd.SetControlValue("ttl", ttl) // set it back
		}
	},

	SubmitF: func(fd *dhtmlform.FormData) {
		wh := fd.GetArg("Webhook").(*model.DdnsWebhook)

		wh.DnsRecordName = fd.GetValue("record_name").(string)
		wh.DnsZoneName = fd.GetValue("zone_name").(string)
		wh.ActualIp = fd.GetValue("actual_ip").(string)
		wh.Ttl = uint32(fd.GetValue("ttl").(int))
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

		piece := dhtml.NewHtmlPiece().Append(
			mtweb.RenderInfof("DNS Record: %s", wh.FQDN()),
			dhtmlbs.NewTextInput("ip").
				Label("IP address to update record to").
				Require().
				Note("example: 127.0.0.1"),
			mtweb.NewDefaultSubmitBtn(),
		)

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

		var log *model.DdnsWebhookLog

		goapp.Transaction(func() error {
			log = helper.UpdateIpAddress(ip)

			if log.Error == "" {
				wh.ActualIp = log.Ip
				goapp.SaveObject(wh)
			}

			return nil
		})

		//		fd.SetRedirect(mbr.Url(WebhookRouteController.LogView, "webhook_id", wh.ID, "log_id", log.ID))
	},
}
