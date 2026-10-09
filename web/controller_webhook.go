package web

import (
	"net/http"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
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

func (c *WebhookRouteControllerType) List() mbr.Route {
	return mbr.Route{
		PathPattern: "/",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Title("Webhooks")

			p.Toolbar().AddIconBtn(
				mbr.Url(WebhookRouteController.Edit, "webhook_id", 0), mtweb.FaIconAdd, "New Webhook",
			)

			table := mtweb.NewTable().
				CountTitle("Webhooks added")

			table.
				EmptyLabel("no webhooks created yet").
				Header("DNS Record FQDN").
				Header("DNS Server").
				Header("Key").
				Header("Description").
				Header("") // Actions

			for _, wh := range goapp.LoadOL[model.DdnsWebhook]() {
				row := table.NewRow()

				row.Cell(wh.FQDN())
				row.Cell(wh.DnsServerAddress)
				row.Cell(mtweb.Icon(iconKey).Label(wh.GetKey().Name))
				row.Cell(wh.Description).Class("small-muted")

				var actions dhtml.HtmlPiece

				actions.
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
				SetRedirect(mbr.Url(WebhookRouteController.List)).
				SetArg("Webhook", wh)

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
