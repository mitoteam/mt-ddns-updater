package web

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mitoteam/mbr"
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

func (c *RootController) Assets() mbr.Route {
	return mbr.Route{PathPattern: "/assets", StaticFS: webAssetsFS}
}

func (c *RootController) MtWebAssets() mbr.Route {
	return mbr.Route{PathPattern: "/assets/mtweb", StaticFS: mtweb.MtWebAssetsFS}
}

func (c *RootController) FavIcon() mbr.Route {
	return mbr.Route{PathPattern: "/favicon.ico", FileFromFS: "favicon.ico", StaticFS: webAssetsFS}
}

func (c *RootController) Home() mbr.Route {
	route := mbr.Route{
		PathPattern: "/",
		HandleF: func(ctx *mbr.MbrContext) any {
			return "Hello, this is MT DDNS Updater!"
		},
	}

	return route
}
