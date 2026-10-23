package web

import (
	"github.com/miekg/dns"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-ddns-updater/app"
	"github.com/mitoteam/mt-ddns-updater/model"
	"github.com/mitoteam/mttools"
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
				row.Cell(mtweb.Icon(iconIpAddress).Label(wh.ActualIp))
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
						mtweb.NewEditBtn(
							mbr.Url(
								WebhookRouteController.Edit, "webhook_id", wh.ID,
								"destination", p.Ctx.Request().RequestURI,
							),
						),
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
				SetArg("Webhook", wh)

			if !p.Ctx.HasDestination() {
				fc.SetRedirect(mbr.Url(WebhookRouteController.List))
			}

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
			ctx.RedirectRoute(WebhookRouteController.List)
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
					mbr.Url(
						WebhookRouteController.Edit, "webhook_id", wh.ID, "destination", p.Ctx.Request().RequestURI,
					),
					mtweb.FaIconEdit, "Edit",
				).
				AddIconBtn(
					mbr.Url(WebhookRouteController.Question, "webhook_id", wh.ID),
					iconMessage, "Request Current IP",
				).
				AddIconBtn(
					mbr.Url(WebhookRouteController.ManualUpdate, "webhook_id", wh.ID),
					mtweb.FaIconUpdate, "Test Update",
				)

			grid := dhtmlbs.NewGridCards().Class("row-cols-lg-2")

			//values
			lc := dhtmlbs.NewListCard().
				Header(mtweb.Icon(iconWebhook).Label(wh.FQDN())).
				Add(dhtml.RenderValue("DNS Record Name", dns.Fqdn(wh.DnsRecordName))).
				Add(dhtml.RenderValue("DNS Zone Name", dns.Fqdn(wh.DnsZoneName))).
				Add(dhtml.RenderValue("TTL", wh.Ttl)).
				Add(dhtml.RenderValue("DNS Server Address", wh.DnsServerAddress))

			if key := wh.GetKey(); key != nil {
				lc.Add(dhtml.RenderValue(
					mtweb.Icon(iconKey).Label("Key"),
					key.Name+" ("+key.Algo.Name()+")",
				))
			}

			if wh.Description != "" {
				lc.Add(dhtml.Div().Class("small-muted").Append(wh.Description))
			}

			grid.Add(lc)

			//info
			lc = dhtmlbs.NewListCard().
				Header(mtweb.Icon(mtweb.FaIconInfo).Label("Info")).
				Add(dhtml.RenderValue(mtweb.Icon(iconIpAddress).Label("Actual IP"), wh.ActualIp))

			grid.Add(lc)

			// add grid to page
			p.Main(grid)

			// #region History
			log_table := mtweb.NewTable().
				EmptyLabel("no requests made yet").
				Header("ID").
				Header("Time").
				Header("Kind").
				Header("IP").
				Header("Error").
				Header("") // Actions

			for _, log := range wh.LogList() {
				row := log_table.NewRow()

				row.Cell(log.ID)
				row.Cell(mtweb.NewTimestamp(log.CreatedAt))
				row.Cell(log.Kind.String())
				row.Cell(log.Ip)
				if log.Error != "" {
					row.Cell(mtweb.RenderError(log.Error))
				} else {
					row.Cell("")
				}

				var actions dhtml.HtmlPiece

				actions.
					Append(
						mtweb.NewSmBtnR(mtweb.FaIconView, WebhookRouteController.LogView, "webhook_id", wh.ID, "log_id", log.ID),
					)
				row.Cell(actions)
			}

			p.Main(
				dhtmlbs.NewCard().Class("mt-3").
					Header(mtweb.Icon(iconWebhookLog).Label("Request Log")).
					Body(log_table),
			)
			// #endregion History

			return nil
		}),
	}
}

func (c *WebhookRouteControllerType) Question() mbr.Route {
	return mbr.Route{
		PathPattern: "/{webhook_id}/question",
		HandleF: func(ctx *mbr.MbrContext) any {
			wh := goapp.LoadOMust[model.DdnsWebhook](ctx.Request().PathValue("webhook_id"))

			helper := app.NewWebhookService(wh)
			var log *model.DdnsWebhookLog

			goapp.Transaction(func() error {
				log = helper.QuestionIpAddress()

				if log.Error == "" {
					wh.ActualIp = log.Ip
					goapp.SaveObject(wh)
				}
				return nil
			})

			ctx.RedirectRoute(WebhookRouteController.LogView, "webhook_id", wh.ID, "log_id", log.ID)
			return nil
		},
	}
}

func (c *WebhookRouteControllerType) ManualUpdate() mbr.Route {
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

// #region Webhook log entries
func (c *WebhookRouteControllerType) LogView() mbr.Route {
	return mbr.Route{
		PathPattern: "/{webhook_id}/log/{log_id}",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			log := goapp.LoadOMust[model.DdnsWebhookLog](p.Ctx.Request().PathValue("log_id"))
			wh := log.GetWebhook()
			mttools.AssertEqual(mttools.AnyToString(wh.ID), p.Ctx.Request().PathValue("webhook_id"))

			p.Title(wh.FQDN())

			baseToolbar(p, wh)

			lc := dhtmlbs.NewListCard().
				Header(mtweb.Icon(mtweb.FaIconInfo).Label("Log Entry Info")).
				Add(dhtml.RenderValue(mtweb.Icon(iconWebhook).Label("Webhook"), wh.FQDN())).
				Add(mtweb.NewTimestamp(log.CreatedAt)).
				Add(dhtml.RenderValue("Kind", log.Kind.String()))

			if log.Ip != "" {
				lc.Add(dhtml.RenderValue(mtweb.Icon(iconIpAddress).Label("IP Address"), log.Ip))
			}

			if log.Error != "" {
				lc.Add(mtweb.RenderError(log.Error))
			}

			p.Main(lc)

			if log.Message != "" {
				card := dhtmlbs.NewCard().Class("mt-3").
					Header(mtweb.Icon(iconMessage).Label("Message")).
					Body(dhtml.NewTag("pre").Class("mb-0").Append(log.Message))

				p.Main(card)
			}

			if log.Reply != "" {
				card := dhtmlbs.NewCard().Class("mt-3").
					Header(mtweb.Icon(iconReply).Label("Reply")).
					Body(dhtml.NewTag("pre").Class("mb-0").Append(log.Reply))

				p.Main(card)
			}

			return nil
		}),
	}
}

// #endregion
