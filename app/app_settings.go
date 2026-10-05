package app

import (
	"github.com/mitoteam/goapp"
)

// Settings are stored in .settings.yml and not changeable at runtime
type AppSettingsType struct {
	goapp.AppSettingsBase `yaml:",inline"`
}

var defaultSettings *AppSettingsType

func init() {
	//default settings
	defaultSettings = &AppSettingsType{}

	//default values for goapp.AppSettingsBase options
	defaultSettings.WebserverPort = 15121
}
