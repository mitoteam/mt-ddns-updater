package web

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-ddns-updater/app"
	"github.com/mitoteam/mt-ddns-updater/model"
	"github.com/mitoteam/mtweb"
)

type RootController struct {
	mbr.ControllerBase
}

var RootCtl *RootController

func init() {
	RootCtl = &RootController{}

	//using chi middlewares
	RootCtl.With(middleware.Recoverer)
}

// Add standard mtweb assets route
func (c *RootController) MtWebAssets() mbr.Route { return mtweb.AssetsRoute }

// custom assets route as embedded FS
func (c *RootController) Assets() mbr.Route {
	return mbr.Route{PathPattern: "/assets", StaticFS: webAssetsFS}
}

func (c *RootController) FavIcon() mbr.Route {
	return mbr.Route{PathPattern: "/favicon.ico", FileFromFS: "favicon.ico", StaticFS: webAssetsFS}
}

func (c *RootController) Home() mbr.Route {
	r := mbr.Route{
		PathPattern: "/",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Main(dhtml.NewTag("h1").Append(app.App.AppName))

			cards_list := dhtmlbs.NewCardList().Class("row-cols-lg-2")

			// Keys card
			cards_list.Add(
				dhtmlbs.NewCard().
					Header(mtweb.Icon(iconKey).Label("Keys")).
					Body(c.renderKeysCard()),
			)

			// Webhooks card
			cards_list.Add(
				dhtmlbs.NewCard().
					Header(mtweb.Icon(iconWebhook).Label("Webhooks")).
					Body(c.renderWebhooksCard()),
			)

			// render cards
			p.Main(cards_list)

			//test
			p.Main(dhtml.Div().Append(
				dhtml.NewLink(mbr.Url(RecordsRouteController.Test)).Label("test"),
			))

			return nil
		}),
	}

	r.With(AuthMiddleware)

	return r
}

func (c *RootController) Login() mbr.Route {
	return mbr.Route{
		PathPattern: "/login",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			if IsAuthenticated(p.Ctx) { // using raw session check instead of p.IsAuthenticated() because no AuthMiddleware used for this route to set Ctx "IsAuthenticated" field
				p.Main(mtweb.RenderInfo("Already authenticated"))
			} else {
				p.Main(formLogin.Render(p.NewFormContext()))
			}

			return nil
		}),
	}
}

func (c *RootController) Logout() mbr.Route {
	return mbr.Route{
		PathPattern: "/logout",
		HandleF: func(ctx *mbr.MbrContext) any {
			err := Logout(ctx)

			if err != nil {
				return err
			}

			ctx.RedirectRoute(http.StatusFound, RootCtl.Home)
			return nil
		},
	}
}

func (c *RootController) renderKeysCard() (out dhtml.HtmlPiece) {
	actions := dhtml.Div().Class("mt-3")

	if cnt := goapp.CountOL[model.DdnsKey](); cnt > 0 {
		out.Append(
			dhtml.RenderValue("Keys added", cnt),
		)

		actions.Append(
			dhtml.NewLink(mbr.Url(KeyRouteController.List)).
				Class("me-3").
				Label(mtweb.Icon(mtweb.FaIconList).Label("View keys")),
		)
	} else {
		out.Append(
			dhtml.Div().Append(dhtml.EmptyLabel("No keys added yet")),
		)
	}

	actions.Append(
		dhtml.NewLink(mbr.Url(KeyRouteController.Edit, "key_id", 0)).Label(mtweb.Icon(mtweb.FaIconAdd).Label("Add new key")),
	)

	out.Append(actions)

	return out
}

func (c *RootController) renderWebhooksCard() (out dhtml.HtmlPiece) {
	actions := dhtml.Div().Class("mt-3")

	if cnt := goapp.CountOL[model.DdnsWebhook](); cnt > 0 {
		out.Append(
			dhtml.RenderValue("Webhooks added", cnt),
		)

		actions.Append(
			dhtml.NewLink(mbr.Url(WebhookRouteController.List)).
				Class("me-3").
				Label(mtweb.Icon(mtweb.FaIconList).Label("View webhooks")),
		)
	} else {
		out.Append(
			dhtml.Div().Append(dhtml.EmptyLabel("No webhooks added yet")),
		)
	}

	if key_cnt := goapp.CountOL[model.DdnsKey](); key_cnt > 0 {
		actions.Append(
			dhtml.NewLink(mbr.Url(WebhookRouteController.Edit, "webhook_id", 0)).Label(mtweb.Icon(mtweb.FaIconAdd).Label("Add new webhook")),
		)
	} else {
		out.Append(dhtml.Div().Class("mt-3").Append(mtweb.RenderInfo("Webhooks can not be created without keys")))
	}

	if actions.ChildrenCount() > 0 {
		out.Append(actions)
	}

	return out
}
