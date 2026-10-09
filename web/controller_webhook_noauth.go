package web

import (
	"github.com/mitoteam/mbr"
)

type WebhookNoAuthRouteControllerType struct {
	mbr.ControllerBase
}

var WebhookNoAuthRouteController *WebhookNoAuthRouteControllerType

func init() {
	WebhookNoAuthRouteController = &WebhookNoAuthRouteControllerType{}
}

// base route for sub routes
func (c *RootController) WebhookNoAuthSubcontroller() mbr.Route {
	return mbr.Route{PathPattern: "/", ChildController: WebhookNoAuthRouteController}
}

func (c *WebhookNoAuthRouteControllerType) Fire() mbr.Route {
	return mbr.Route{
		PathPattern: "/fire/{secure_token}",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Title("Fire!")

			//mbr.Url(WebhookNoAuthRouteController.Fire, "secure_token", wh.SecureToken)

			return nil
		}),
	}
}
