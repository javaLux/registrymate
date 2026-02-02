package ui

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const projectWebsite = "https://github.com/javaLux/registrymate/"

// =======================
// Singleton infrastructure
// =======================
type aboutSingleton struct {
	metadata fyne.AppMetadata
	parent   fyne.Window
	window   fyne.Window
	isShown  bool
	toast    *ToastPopup
}

var (
	aboutInstance *aboutSingleton
	aboutOnce     sync.Once
)

// ShowAbout shows the singleton About window.
func ShowAbout(metadata fyne.AppMetadata, parent fyne.Window) {
	aboutOnce.Do(func() {
		aboutInstance = &aboutSingleton{
			parent:   parent,
			metadata: metadata,
		}
	})

	aboutInstance.show()
}

// Window lifecycle
func (a *aboutSingleton) show() {
	if a.window == nil {
		a.createWindow()
		a.buildContent()
	}

	if a.isShown {
		a.window.RequestFocus()
		return
	}

	a.isShown = true
	a.window.Show()
}

func (a *aboutSingleton) createWindow() {
	w := fyne.CurrentApp().NewWindow("About")
	w.Resize(fyne.NewSize(420, 285))
	w.SetFixedSize(true)
	w.CenterOnScreen()

	// IMPORTANT:
	// Closing the window destroys it permanently.
	w.SetOnClosed(func() {
		a.window = nil
		a.isShown = false
		a.parent.RequestFocus()
	})

	a.toast = NewToastPopup(BlueTextColor, w.Canvas())
	a.window = w
}

func (a *aboutSingleton) buildContent() {
	header := a.centeredHeader()

	description := widget.NewLabel("Easily create Kubernetes Image-Pull-Secrets")
	description.Wrapping = fyne.TextWrapWord
	description.Alignment = fyne.TextAlignCenter

	version := widget.NewLabel(
		fmt.Sprintf(
			"Version: %s (Build %d)",
			a.metadata.Version,
			a.metadata.Build,
		),
	)
	version.TextStyle = fyne.TextStyle{Italic: true}
	version.Alignment = fyne.TextAlignCenter

	var releaseState string
	if a.metadata.Release {
		releaseState = "Yes"
	} else {
		releaseState = "No"
	}

	release := widget.NewLabel(
		fmt.Sprintf(
			"Fyne Release-Build: %s",
			releaseState,
		),
	)

	release.TextStyle = fyne.TextStyle{Italic: true}
	release.Alignment = fyne.TextAlignCenter

	commit := widget.NewLabel("")
	// only add commit if available
	if a.commit() != "" {
		commit.SetText(fmt.Sprintf("Commit: %s", a.commit()))
		commit.TextStyle = fyne.TextStyle{Italic: true}
		commit.Alignment = fyne.TextAlignCenter
	}

	okBtn := widget.NewButton("OK", func() {
		a.window.Hide()
		a.isShown = false
	})

	// globally key event handler to close the about window
	a.window.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyReturn, fyne.KeyEnter:
			okBtn.OnTapped()
		case fyne.KeyEscape:
			a.window.Hide()
			a.isShown = false
		}
	})

	footer := container.NewVBox(
		widget.NewSeparator(),
		a.footer(),
	)

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		description,
		version,
		release,
		commit,
		footer,
		okBtn,
	)

	a.window.SetContent(container.NewPadded(content))
}

func (a *aboutSingleton) centeredHeader() fyne.CanvasObject {
	icon := widget.NewIcon(fyne.CurrentApp().Icon())

	title := widget.NewLabelWithStyle(
		a.metadata.Name,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	center := container.NewHBox(icon, title)

	return container.NewBorder(
		nil, nil, nil,
		a.copyButton(),
		container.NewHBox(
			layout.NewSpacer(),
			center,
			layout.NewSpacer(),
		),
	)
}

func (a *aboutSingleton) footer() fyne.CanvasObject {
	linkURL, _ := url.Parse(projectWebsite)
	link := widget.NewHyperlink("GitHub", linkURL)
	link.Alignment = fyne.TextAlignCenter

	return container.NewHBox(
		layout.NewSpacer(),
		container.NewVBox(link),
		layout.NewSpacer(),
	)
}

func (a *aboutSingleton) copyButton() *widget.Button {
	btn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		fyne.CurrentApp().Clipboard().SetContent(a.aboutClipboardYAML())
		a.toast.ShowToast("Copied", 2*time.Second)
	})
	return btn
}

func (a *aboutSingleton) aboutClipboardYAML() string {

	if a.commit() == "" {
		return fmt.Sprintf(
			"name: %s\nversion: %s\nbuild: %d\nfyne-release-build: %t\nproject: %s\n",
			a.metadata.Name,
			a.metadata.Version,
			a.metadata.Build,
			a.metadata.Release,
			projectWebsite,
		)
	} else {
		return fmt.Sprintf(
			"name: %s\nversion: %s\nbuild: %d\nfyne-release-build: %t\ncommit: %s\nproject: %s\n",
			a.metadata.Name,
			a.metadata.Version,
			a.metadata.Build,
			a.metadata.Release,
			a.commit(),
			projectWebsite,
		)
	}
}

func (a *aboutSingleton) commit() string {
	commit := a.metadata.Custom["Commit"]

	return strings.TrimSpace(commit)
}

func (a *aboutSingleton) _() bool {
	if a.window == nil {
		return false
	}

	canvas := a.window.Canvas()
	if canvas == nil {
		return false
	}

	return canvas.Focused() != nil
}
