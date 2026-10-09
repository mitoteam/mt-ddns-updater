package web

import (
	"net/http"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-ddns-updater/model"
	"github.com/mitoteam/mtweb"
)

type KeyRouteControllerType struct {
	mbr.ControllerBase
}

var KeyRouteController *KeyRouteControllerType

func init() {
	KeyRouteController = &KeyRouteControllerType{}
	KeyRouteController.With(AuthMiddleware)
}

// base route for sub routes
func (c *RootController) KeySubcontroller() mbr.Route {
	return mbr.Route{PathPattern: "/key", ChildController: KeyRouteController}
}

func (c *KeyRouteControllerType) List() mbr.Route {
	return mbr.Route{
		PathPattern: "/",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Title("Keys")

			p.Toolbar().AddIconBtn(
				mbr.Url(KeyRouteController.Edit, "key_id", 0), mtweb.FaIconAdd, "New Key",
			)

			table := mtweb.NewTable().
				CountTitle("Keys added")

			table.
				EmptyLabel("no keys created yet").
				Header("Name").
				Header("Algo").
				Header("Secret").
				Header("Webhooks").
				Header("Description").
				Header("") // Actions

			for _, key := range goapp.LoadOL[model.DdnsKey]() {
				row := table.NewRow()

				row.Cell(key.Name)
				row.Cell(key.GetAlgoName())
				row.Cell("[hidden]").Class("text-muted")
				row.Cell(key.WebhooksCount())
				row.Cell(key.Description).Class("small-muted")

				var actions dhtml.HtmlPiece

				actions.
					Append(mtweb.NewEditBtn(mbr.Url(KeyRouteController.Edit, "key_id", key.ID))).
					Append(mtweb.NewDeleteBtn(mbr.Url(KeyRouteController.Delete, "key_id", key.ID), ""))

				row.Cell(actions)
			}

			p.Main(table)

			return nil
		}),
	}
}

func (c *KeyRouteControllerType) Edit() mbr.Route {
	return mbr.Route{
		PathPattern: "/{key_id}/edit",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			key := goapp.LoadOrCreateO[model.DdnsKey](p.Ctx.Request().PathValue("key_id"))

			if key.ID == 0 {
				p.Title("New key")

				//set default values for key
				key.Algo = model.KeyAlgoSHA256
			} else {
				p.Title("Edit key")
			}

			fc := p.NewFormContext().
				SetRedirect(mbr.Url(KeyRouteController.List)).
				SetArg("Key", key)

			p.Main(formKeyEdit.Render(fc))

			return nil
		}),
	}
}

func (c *KeyRouteControllerType) Delete() mbr.Route {
	return mbr.Route{
		PathPattern: "/{key_id}/delete",
		HandleF: func(ctx *mbr.MbrContext) any {
			goapp.DeleteObject(goapp.LoadOMust[model.DdnsKey](ctx.Request().PathValue("key_id")))
			ctx.RedirectRoute(http.StatusFound, KeyRouteController.List)
			return nil
		},
	}
}
