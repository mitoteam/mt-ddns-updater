package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-ddns-updater/model"
	"github.com/mitoteam/mtweb"
)

var formKeyEdit = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		key := fd.GetArg("Key").(*model.DdnsKey)

		formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlbs.NewTextInput("name").Label("Template Name").Default(key.Name).Require(),
			dhtmlbs.NewTextarea("description").Label("Description").
				Default(key.Description).Note("Description for key"),
			mtweb.NewDefaultSubmitBtn(),
		))
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		key := fd.GetArg("Key").(*model.DdnsKey)

		key.Name = fd.GetValue("name").(string)
		key.Description = fd.GetValue("description").(string)

		goapp.SaveObject(key)
	},
}
