package web

import (
	"github.com/mitoteam/mbr"
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
			p.Main("Key List")

			return nil
		}),
	}
}

func (c *KeyRouteControllerType) Edit() mbr.Route {
	return mbr.Route{
		PathPattern: "/add",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Main("Key Edit")

			return nil
		}),
	}
}
