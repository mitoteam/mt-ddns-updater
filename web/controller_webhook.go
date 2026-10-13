package web

import (
	"fmt"
	"net/http"

	"github.com/miekg/dns"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-ddns-updater/app"
	"github.com/mitoteam/mt-ddns-updater/model"
	"github.com/mitoteam/mtweb"
)

type WebhookRouteControllerType struct {
	mbr.ControllerBase
}

var WebhookRouteController *WebhookRouteControllerType

func init() {
	WebhookRouteController = &WebhookRouteControllerType{}
	WebhookRouteController.With(AuthMiddleware)
}

// base route for sub routes
func (c *RootController) WebhookSubcontroller() mbr.Route {
	return mbr.Route{PathPattern: "/webhook", ChildController: WebhookRouteController}
}

func baseToolbar(p *PageBuilder, wh *model.DdnsWebhook) *mtweb.BtnPanelElement {
	t := p.Toolbar().
		AddIconBtn(mbr.Url(WebhookRouteController.List), iconWebhook, "All Webhooks")

	if wh != nil {
		t.AddIconBtn(
			mbr.Url(WebhookRouteController.View, "webhook_id", wh.ID), iconWebhook, wh.FQDN(),
		)
	}

	return t
}

func (c *WebhookRouteControllerType) List() mbr.Route {
	return mbr.Route{
		PathPattern: "/",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Title("Webhooks")

			baseToolbar(p, nil).
				AddIconBtn(
					mbr.Url(WebhookRouteController.Edit, "webhook_id", 0), mtweb.FaIconAdd, "New Webhook",
				)

			table := mtweb.NewTable().
				CountTitle("Webhooks added")

			table.
				EmptyLabel("no webhooks created yet").
				Header("DNS Record FQDN").
				Header("Actual IP").
				Header("TTL").
				Header("DNS Server").
				Header("Key").
				Header("Description").
				Header("") // Actions

			for _, wh := range goapp.LoadOL[model.DdnsWebhook]() {
				row := table.NewRow()

				row.Cell(wh.FQDN())
				row.Cell(wh.ActualIp)
				row.Cell(wh.Ttl)
				row.Cell(wh.DnsServerAddress)
				row.Cell(mtweb.Icon(iconKey).Label(wh.GetKey().Name))
				row.Cell(wh.Description).Class("small-muted")

				var actions dhtml.HtmlPiece

				actions.
					Append(
						mtweb.NewSmBtnR(mtweb.FaIconView, WebhookRouteController.View, "webhook_id", wh.ID),
					).
					Append(
						mtweb.NewEditBtn(mbr.Url(WebhookRouteController.Edit, "webhook_id", wh.ID)),
					).
					Append(
						mtweb.NewDeleteBtn(mbr.Url(WebhookRouteController.Delete, "webhook_id", wh.ID), ""),
					)

				row.Cell(actions)
			}

			p.Main(table)

			return nil
		}),
	}
}

func (c *WebhookRouteControllerType) Edit() mbr.Route {
	return mbr.Route{
		PathPattern: "/{webhook_id}/edit",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			wh := goapp.LoadOrCreateO[model.DdnsWebhook](p.Ctx.Request().PathValue("webhook_id"))

			if wh.ID == 0 {
				p.Title("New webhook")

				//pre select key if there is just one
				if goapp.CountOL[model.DdnsKey]() == 1 {
					key := goapp.FirstO[model.DdnsKey]()
					wh.KeyID = key.ID
				}
			} else {
				p.Title("Edit webhook")
			}

			fc := p.NewFormContext().
				SetArg("Webhook", wh).
				SetRedirect(mbr.Url(WebhookRouteController.List))

			p.Main(formWebhookEdit.Render(fc))

			return nil
		}),
	}
}

func (c *WebhookRouteControllerType) Delete() mbr.Route {
	return mbr.Route{
		PathPattern: "/{webhook_id}/delete",
		HandleF: func(ctx *mbr.MbrContext) any {
			goapp.DeleteObject(
				goapp.LoadOMust[model.DdnsWebhook](ctx.Request().PathValue("webhook_id")),
			)
			ctx.RedirectRoute(http.StatusFound, WebhookRouteController.List)
			return nil
		},
	}
}

func (c *WebhookRouteControllerType) View() mbr.Route {
	return mbr.Route{
		PathPattern: "/{webhook_id}",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			wh := goapp.LoadOMust[model.DdnsWebhook](p.Ctx.Request().PathValue("webhook_id"))

			p.Title(wh.FQDN())

			baseToolbar(p, nil).
				AddIconBtn(
					mbr.Url(WebhookRouteController.Edit, "webhook_id", wh.ID),
					mtweb.FaIconEdit, "Edit",
				).
				AddIconBtn(
					mbr.Url(WebhookRouteController.Question, "webhook_id", wh.ID),
					mtweb.FaIconSearch, "Request Current IP",
				).
				AddIconBtn(
					mbr.Url(WebhookRouteController.TestUpdate, "webhook_id", wh.ID),
					mtweb.FaIconUpdate, "Test Update",
				)

			//values
			lc := dhtmlbs.NewListCard().Header(wh.FQDN()).
				Add(dhtml.RenderValue("DNS Record Name", dns.Fqdn(wh.DnsRecordName))).
				Add(dhtml.RenderValue("DNS Zone Name", dns.Fqdn(wh.DnsZoneName))).
				Add(dhtml.RenderValue("TTL", wh.Ttl)).
				Add(dhtml.RenderValue("DNS Server Address", wh.DnsServerAddress))

			if key := wh.GetKey(); key != nil {
				lc.Add(dhtml.RenderValue(
					mtweb.Icon(iconKey).Label("Key"),
					key.Name+" ("+key.GetAlgoName()+")",
				))
			}

			if wh.Description != "" {
				lc.Add(dhtml.Div().Class("small-muted").Append(wh.Description))
			}

			grid := dhtmlbs.NewGridCards().Class("row-cols-lg-2")

			grid.Add(lc)

			p.Main(grid)

			return nil
		}),
	}
}

func (c *WebhookRouteControllerType) Question() mbr.Route {
	return mbr.Route{
		PathPattern: "/{webhook_id}/question",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			wh := goapp.LoadOMust[model.DdnsWebhook](p.Ctx.Request().PathValue("webhook_id"))

			p.Title(wh.FQDN())

			baseToolbar(p, wh)

			helper := app.NewWebhookService(wh)
			ip, err := helper.QuestionIpAddress()

			piece := dhtml.NewHtmlPiece().
				Append(dhtml.Div().Class("fw-bold fs-4").Append("Message")).
				Append(dhtml.NewTag("pre").Append(helper.GetLastMessage()))

			if err == nil {
				s := fmt.Sprintf("A-Record %s address resolved to %s", wh.FQDN(), ip.String())

				piece.
					Append(dhtml.Div().Class("fw-bold fs-4").Append("Reply")).
					Append(mtweb.RenderInfo(s)).
					Append(dhtml.NewTag("pre").Append(helper.GetLastReply()))

				//log.Println(s)

				wh.ActualIp = ip.String()
				goapp.SaveObject(wh)
			} else {
				s := fmt.Sprintf("Error during DNS query: %v", err)
				piece.
					Append(dhtml.Div().Class("fw-bold fs-4").Append("Error")).
					Append(mtweb.RenderError(s))
				//log.Println(s)
			}

			p.Main(piece)

			return nil
		}),
	}
}

func (c *WebhookRouteControllerType) TestUpdate() mbr.Route {
	return mbr.Route{
		PathPattern: "/{webhook_id}/test-update",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			wh := goapp.LoadOMust[model.DdnsWebhook](p.Ctx.Request().PathValue("webhook_id"))

			p.Title(wh.FQDN())

			baseToolbar(p, wh)

			fc := p.NewFormContext().
				SetArg("Webhook", wh)

			p.Main(formWebhookTestUpdate.Render(fc))

			return nil
		}),
	}
}
