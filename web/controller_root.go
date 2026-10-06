package web

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-ddns-updater/app"
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
			p.Main("SOMETHING goes here...")
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
