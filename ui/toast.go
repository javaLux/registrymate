package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var BlueTextColor = color.NRGBA{R: 0x00, G: 0x6c, B: 0xff, A: 0xff}

type ToastPopup struct {
	*widget.PopUp
	parent    fyne.Canvas
	textColor color.Color
}

func NewToastPopup(textColor color.Color, parent fyne.Canvas) *ToastPopup {
	return &ToastPopup{
		textColor: textColor,
		parent:    parent,
	}
}

func (t *ToastPopup) ShowToast(message string, timeToShow time.Duration) {
	text := canvas.NewText(message, t.textColor)
	text.Alignment = fyne.TextAlignCenter
	text.TextStyle.Bold = true
	text.TextSize = 16

	content := container.NewStack(text)

	canvasSize := t.parent.Size()
	toastMinSize := content.MinSize()

	toast := widget.NewPopUp(content, t.parent)
	// Position at top center
	position := fyne.NewPos(
		(canvasSize.Width-toastMinSize.Width)/2,
		12,
	)

	// show at top center
	toast.ShowAtPosition(position)

	// Auto-hide after given duration
	go func(popup *widget.PopUp, duration time.Duration) {
		time.Sleep(duration)
		fyne.Do(func() {
			popup.Hide()
		})
	}(toast, timeToShow)
}
