package web

import (
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
	return mbr.Route{
		PathPattern: "/",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Main(dhtml.NewTag("h1").Append(app.App.AppName))
			p.Main("SOMETHING goes here...")
			return nil
		}),
	}
}
