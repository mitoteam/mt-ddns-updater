package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-ddns-updater/model"
	"github.com/mitoteam/mttools"
	"github.com/mitoteam/mtweb"
)

var formKeyEdit = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		key := fd.GetArg("Key").(*model.DdnsKey)

		formBody.Append(
			dhtml.Div().Class("border bg-light p-3").Append(
				dhtmlbs.NewTextInput("name").Label("Key Name").Default(key.Name).Require(),
				dhtmlbs.NewSelect("algo").
					OptionsMap(model.DdnsKeyAlgoOptions()).
					Label("Algorithm").
					Default(key.Algo).
					Require(),
				dhtmlbs.NewTextInput("secret").
					Label("Secret").
					Default(key.Secret).
					Note("base64").
					Require(),
				dhtmlbs.NewTextarea("description").Label("Description").
					Default(key.Description).Note("Description for key"),
				mtweb.NewDefaultSubmitBtn(),
			),
		)
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		key := fd.GetArg("Key").(*model.DdnsKey)

		key.Name = fd.GetValue("name").(string)

		if algo, ok := mttools.AnyToIntOk(fd.GetValue("algo")); ok {
			key.Algo = algo
		}
		key.Secret = fd.GetValue("secret").(string)
		key.Description = fd.GetValue("description").(string)

		goapp.SaveObject(key)
	},
}
