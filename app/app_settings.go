package app

import (
	"github.com/mitoteam/goapp"
)

// Settings are stored in .settings.yml and not changeable at runtime
type AppSettingsType struct {
	goapp.AppSettingsBase `yaml:",inline"`

	GuiPassword       string `yaml:"gui_password" yaml_comment:"WebGUI access password"`
	GuiSessionTimeout int    `yaml:"gui_session_timeout" yaml_comment:"WebGUI session timeout in days. Default is 7 days."`
}

var defaultSettings *AppSettingsType

func init() {
	//default settings
	defaultSettings = &AppSettingsType{
		GuiPassword:       "mitoteam",
		GuiSessionTimeout: 7,
	}

	//default values for goapp.AppSettingsBase options
	defaultSettings.WebserverPort = 15121
}
