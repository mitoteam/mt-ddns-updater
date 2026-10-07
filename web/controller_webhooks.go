package web

import (
	"github.com/mitoteam/mbr"
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
			p.Main("Webhook List")

			return nil
		}),
	}
}

func (c *WebhookRouteControllerType) Edit() mbr.Route {
	return mbr.Route{
		PathPattern: "/edit",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Main("Webhook Edit")

			return nil
		}),
	}
}
