package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// ForcedVariant is a theme wrapper that forces a specific variant (light or dark)
type ForcedVariant struct {
	fyne.Theme

	Variant fyne.ThemeVariant
}

func (f *ForcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.Variant)
}
