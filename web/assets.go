package web

import (
	"embed"
	"io/fs"
)

// embedded web assets
//
//go:embed assets/script.min.js
//go:embed assets/vendor/*.css assets/vendor/*.js
//go:embed assets/webfonts/*
//go:embed assets/css/*.css
//go:embed assets/images/*.png
//go:embed assets/favicon.ico
var embedFS embed.FS

var webAssetsFS fs.FS

func init() {
	//prepare fs.FS for embedded subdirectory "/assets"
	webAssetsFS, _ = fs.Sub(embedFS, "assets")
}
