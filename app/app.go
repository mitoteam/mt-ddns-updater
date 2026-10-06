package app

import (
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
)

var App *goapp.AppBase

func InitApp() *goapp.AppBase {
	App = goapp.NewAppBase(defaultSettings)

	App.AppName = "MT DDNS Updater"
	App.ExecutableName = "mt-ddns-updater"
	App.LongDescription = `Update DNS records using webhook`

	App.PreRunF = DoPreRun
	App.PostRunF = DoPostRun

	return App
}

func DoPreRun() (err error) {
	// open database and migrate schema
	if err = goapp.DbSchema.Open(App.AppSettings.(*AppSettingsType).LogSql); err != nil {
		return err
	}

	// start forms data expiration handler
	dhtmlform.StartFormDataExpirationHandler(App.BaseContext)

	return nil //no errors
}

func DoPostRun() error {
	goapp.DbSchema.Close()

	return nil //no errors
}
