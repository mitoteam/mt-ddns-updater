package web

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-ddns-updater/app"
	"github.com/mitoteam/mtweb"
)

type RecordsRouteControllerType struct {
	mbr.ControllerBase
}

var RecordsRouteController *RecordsRouteControllerType

func init() {
	RecordsRouteController = &RecordsRouteControllerType{}
	RecordsRouteController.With(AuthMiddleware)
}

// base route for sub routes
func (c *RootController) RecordsSubcontroller() mbr.Route {
	return mbr.Route{PathPattern: "/record", ChildController: RecordsRouteController}
}

func (c *RecordsRouteControllerType) Test() mbr.Route {
	return mbr.Route{
		PathPattern: "/test",
		HandleF: CreatePageBuilderRouteHandler(func(p *PageBuilder) any {
			zone := "mito-team.com"
			name := "ddns-test"
			ip := net.ParseIP("5.5.5." + strconv.Itoa(rand.Intn(253)+1))

			h := app.NewDdnsHelper("***", zone, "***", "***")

			err := h.UpdateARecordIpAddress(name, ip, 55)

			piece := dhtml.NewHtmlPiece().
				Append(dhtml.Div().Class("fw-bold fs-4").Append("Message")).
				Append(dhtml.NewTag("pre").Append(h.GetLastMessage()))

			if err == nil {
				s := fmt.Sprintf("A-Record %s.%s updated successfully to %s", name, zone, ip.String())

				piece.
					Append(dhtml.Div().Class("fw-bold fs-4").Append("Reply")).
					Append(mtweb.RenderInfo(s)).
					Append(dhtml.NewTag("pre").Append(h.GetLastReply()))

				log.Println(s)
			} else {
				s := fmt.Sprintf("Error during DNS update: %v", err)
				piece.
					Append(dhtml.Div().Class("fw-bold fs-4").Append("Error")).
					Append(mtweb.RenderError(s))
				log.Println(s)
			}

			p.Main(piece)

			return nil
		}),
	}
}
