package web

import (
	"time"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-ddns-updater/app"
	"github.com/mitoteam/mtweb"
)

// GUI login form
var formLogin = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlbs.NewFloatingPasswordInput("password").Label("Password").Require(),

			mtweb.NewIconSubmitBtn("arrow-right-to-bracket", "Sign In"),
		))
	},
	ValidateF: func(fd *dhtmlform.FormData) {
		if !fd.HasError() {
			password := fd.GetValue("password").(string)

			ctx := fd.GetParam("MbrContext").(*mbr.MbrContext)

			if password != app.App.AppSettings.(*app.AppSettingsType).GuiPassword {
				Logout(ctx) // clear session if any
				fd.SetError("", "Wrong password given")
			}
		}
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		ctx := fd.GetParam("MbrContext").(*mbr.MbrContext)
		session := GetSession(ctx)

		timeout := time.Now().Add(time.Duration(app.App.AppSettings.(*app.AppSettingsType).GuiSessionTimeout) * 24 * time.Hour)
		session.Values[sessionTimeoutField] = timeout.Format(time.RFC3339)

		session.Save(ctx.Request(), ctx.Writer())
	},
}
