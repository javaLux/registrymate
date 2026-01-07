package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/javaLux/registrymate/utils"
)

func main() {
	app := app.New()

	// load app settings from config
	appSettings := utils.LoadAppSettings(app)

	g := newGenerator(appSettings)
	g.loadUI(app)
	app.Run()
}
