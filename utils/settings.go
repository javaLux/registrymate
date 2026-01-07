package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type AppSettings struct {
	appTheme fyne.ThemeVariant
	Width    float32
	Height   float32
	History  *AppHistory
}

func (a *AppSettings) SetThemeVariant(variant fyne.ThemeVariant) {
	a.appTheme = variant
}

func (a *AppSettings) ThemeVariant() fyne.ThemeVariant {
	return a.appTheme
}

func (a *AppSettings) IsLightTheme() bool {
	return a.appTheme == theme.VariantLight
}

// Load saved App settings or use default
// App settings stored as JSON in the user config directory:  ~/<Config-Dir>/fyne/<appname>/preferences.json
func LoadAppSettings(a fyne.App) *AppSettings {
	isLightTheme := a.Preferences().BoolWithFallback(PrefKeyIsLightTheme, false)
	width := a.Preferences().FloatWithFallback(PrefKeyWindowWidth, 650)
	height := a.Preferences().FloatWithFallback(PrefKeyWindowHeight, 550)
	registries := a.Preferences().StringListWithFallback(PrefKeyRegistries, []string{})
	namespaces := a.Preferences().StringListWithFallback(PrefKeyNamespaces, []string{})
	names := a.Preferences().StringListWithFallback(PrefKeyNames, []string{})

	var appTheme fyne.ThemeVariant
	if isLightTheme {
		appTheme = theme.VariantLight
	} else {
		appTheme = theme.VariantDark
	}

	history := NewAppHistory()
	history.SetRegistries(registries)
	history.SetNamespaces(namespaces)
	history.SetNames(names)

	return &AppSettings{
		appTheme: appTheme,
		Width:    float32(width),
		Height:   float32(height),
		History:  history,
	}
}

// Save App settings before exit
func (a *AppSettings) SaveAppSettings(app fyne.App) {
	app.Preferences().SetBool(PrefKeyIsLightTheme, a.IsLightTheme())
	app.Preferences().SetFloat(PrefKeyWindowWidth, float64(a.Width))
	app.Preferences().SetFloat(PrefKeyWindowHeight, float64(a.Height))
	app.Preferences().SetStringList(PrefKeyRegistries, a.History.Registries)
	app.Preferences().SetStringList(PrefKeyNamespaces, a.History.Namespaces)
	app.Preferences().SetStringList(PrefKeyNames, a.History.Names)
}
