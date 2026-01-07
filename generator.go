package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/javaLux/registrymate/ui"
	"github.com/javaLux/registrymate/utils"
)

const DefaultOutputText = "Nothing to show yet..."

type generator struct {
	appSettings            *utils.AppSettings
	regEntry               *widget.SelectEntry
	userEntry              *widget.Entry
	passEntry              *widget.Entry
	nameSpaceEntry         *widget.SelectEntry
	nameEntry              *widget.SelectEntry
	generateBtn            *widget.Button
	clearRegEntryBtn       *widget.Button
	clearUserEntryBtn      *widget.Button
	clearPassEntryBtn      *widget.Button
	clearNameSpaceEntryBtn *widget.Button
	clearNameEntryBtn      *widget.Button
	clearOutputBtn         *widget.Button
	clearHistoryBtn        *widget.Button
	decodeBtn              *widget.Button
	saveBtn                *widget.Button
	copyBtn                *widget.Button
	themeBtn               *widget.Button
	output                 *widget.Label
	secret                 *Secret
	window                 fyne.Window
	isDecoded              bool
	toast                  *ui.ToastPopup
}

func newGenerator(appSettings *utils.AppSettings) *generator {
	return &generator{
		appSettings: appSettings,
		isDecoded:   false,
	}
}

func (g *generator) loadUI(app fyne.App) {
	// set theme based on saved app settings
	fyne.CurrentApp().Settings().SetTheme(&ui.ForcedVariant{Theme: theme.DefaultTheme(), Variant: g.appSettings.ThemeVariant()})

	w := app.NewWindow(utils.AppName)
	w.Resize(fyne.NewSize(g.appSettings.Width, g.appSettings.Height))

	// save app settings on close event
	w.SetOnClosed(func() {
		g.appSettings.Width = w.Content().Size().Width
		g.appSettings.Height = w.Content().Size().Height
		g.appSettings.SaveAppSettings(app)
	})

	g.window = w

	// --- Labels ---
	g.buildLabels()
	// --- Entries ---
	g.buildEntries()

	// --- Buttons ---
	g.buildButtons()

	// --- Toast Popup---
	g.toast = ui.NewToastPopup(ui.BlueTextColor, g.window.Canvas())

	g.window.SetContent(g.buildLayout())
	g.window.Show()
}

func (g *generator) buildLabels() {
	output := widget.NewLabel(DefaultOutputText)
	output.Wrapping = fyne.TextWrapOff
	output.Importance = widget.HighImportance
	output.TextStyle.Monospace = true
	g.output = output
}

func (g *generator) buildButtons() {
	generateBtn := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), g.buildSecret)
	generateBtn.Disable() // initially disabled until required fields are filled

	decodeBtn := widget.NewButtonWithIcon("", theme.VisibilityOffIcon(), g.decodeOrEncodeSecret)
	decodeBtn.Disable() // initially disabled until a secret is generated

	saveBtn := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), g.saveDialog)
	saveBtn.Disable() // initially disabled until a secret is generated

	copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), g.copyToClipboard)
	copyBtn.Disable() // initially disabled until a secret is generated

	clearOutputBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), g.clearOutput)
	clearOutputBtn.Disable() // initially disabled until a secret is generated

	clearHistoryBtn := widget.NewButtonWithIcon("", theme.HistoryIcon(), g.clearHistory)
	if g.appSettings.History.IsEmpty() {
		clearHistoryBtn.Disable()
	}

	clearRegEntryBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() { g.regEntry.SetText("") })
	clearUserEntryBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() { g.userEntry.SetText("") })
	clearPassEntryBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() { g.passEntry.SetText("") })

	clearNameEntryBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() { g.nameEntry.SetText("") })
	clearNameSpaceEntryBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() { g.nameSpaceEntry.SetText("") })

	var themeBtnIcon fyne.Resource

	if g.appSettings.IsLightTheme() {
		themeBtnIcon = ui.DarkThemeIcon
	} else {
		themeBtnIcon = ui.LightThemeIcon
	}

	var themeBtn *widget.Button
	themeBtn = widget.NewButtonWithIcon("", themeBtnIcon, g.changeTheme)

	g.clearHistoryBtn = clearHistoryBtn
	g.generateBtn = generateBtn
	g.decodeBtn = decodeBtn
	g.saveBtn = saveBtn
	g.copyBtn = copyBtn
	g.clearOutputBtn = clearOutputBtn
	g.clearRegEntryBtn = clearRegEntryBtn
	g.clearUserEntryBtn = clearUserEntryBtn
	g.clearPassEntryBtn = clearPassEntryBtn
	g.clearNameEntryBtn = clearNameEntryBtn
	g.clearNameSpaceEntryBtn = clearNameSpaceEntryBtn
	g.themeBtn = themeBtn
}

func (g *generator) buildEntries() {
	regEntry := widget.NewSelectEntry(g.appSettings.History.SortedRegistries())
	regEntry.SetPlaceHolder("Registry (e.g. registry.gitlab.com)")
	regEntry.OnChanged = func(s string) { g.canGenerate() }
	regEntry.OnSubmitted = func(string) {
		if g.isRequiredInputFilled() {
			g.buildSecret()
		}
	}

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("Username")
	userEntry.OnChanged = func(string) { g.canGenerate() }
	userEntry.OnSubmitted = func(string) {
		if g.isRequiredInputFilled() {
			g.buildSecret()
		}
	}

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Password")
	passEntry.OnChanged = func(string) { g.canGenerate() }
	passEntry.OnSubmitted = func(string) {
		if g.isRequiredInputFilled() {
			g.buildSecret()
		}
	}

	nameEntry := widget.NewSelectEntry(g.appSettings.History.SortedNames())
	nameEntry.SetPlaceHolder("Secret-Name (optional)")
	nameEntry.OnSubmitted = func(string) {
		if g.isRequiredInputFilled() {
			g.buildSecret()
		}
	}
	nameEntry.AlwaysShowValidationError = true
	nameEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		if !utils.IsK8sNameValid(s) {
			return fmt.Errorf("invalid K8s secret name")
		}
		return nil
	}

	nameSpaceEntry := widget.NewSelectEntry(g.appSettings.History.SortedNamespaces())
	nameSpaceEntry.SetPlaceHolder("Namespace (optional)")
	nameSpaceEntry.OnSubmitted = func(string) {
		if g.isRequiredInputFilled() {
			g.buildSecret()
		}
	}
	nameSpaceEntry.AlwaysShowValidationError = true
	nameSpaceEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		if !utils.IsK8sNameValid(s) {
			return fmt.Errorf("invalid K8s namespace")
		}
		return nil
	}

	g.regEntry = regEntry
	g.userEntry = userEntry
	g.passEntry = passEntry
	g.nameEntry = nameEntry
	g.nameSpaceEntry = nameSpaceEntry
}

func (g *generator) buildLayout() fyne.CanvasObject {
	// Theme toggle button at the top right corner
	topLayout := container.NewHBox(g.clearHistoryBtn, layout.NewSpacer(), g.themeBtn)

	// Registry input with clear buttons
	regEntryContainer := container.NewBorder(nil, nil, nil, g.clearRegEntryBtn, g.regEntry)
	userEntryContainer := container.NewBorder(nil, nil, nil, g.clearUserEntryBtn, g.userEntry)
	passEntryContainer := container.NewBorder(nil, nil, nil, g.clearPassEntryBtn, g.passEntry)
	registryInput := widget.NewCard("Registry", "",
		container.NewVBox(
			regEntryContainer,
			userEntryContainer,
			passEntryContainer,
		))

	// Secret-Metadata input with clear buttons
	nameEntryContainer := container.NewBorder(nil, nil, nil, g.clearNameEntryBtn, g.nameEntry)
	nameSpaceEntryContainer := container.NewBorder(nil, nil, nil, g.clearNameSpaceEntryBtn, g.nameSpaceEntry)
	metadataInput := widget.NewCard("Metadata", "",
		container.NewVBox(
			nameEntryContainer,
			nameSpaceEntryContainer,
		))

	// Combine registry and metadata inputs side by side
	inputContainer := container.NewGridWithColumns(2, registryInput, metadataInput)

	// Buttons at the middle section to generate, save, copy and clear output
	buttonContainer := widget.NewCard("", "",
		container.NewHBox(
			g.generateBtn,
			layout.NewSpacer(),
			g.decodeBtn,
			g.saveBtn,
			g.copyBtn,
			g.clearOutputBtn,
		))

	// YAML header
	yamlHeader := canvas.NewText("Secret", nil)
	yamlHeader.TextSize = theme.Size(theme.SizeNameHeadingText)
	yamlHeader.Alignment = fyne.TextAlignCenter
	yamlHeader.TextStyle.Bold = true

	// Output area with horizontal scroll
	outputScroll := container.NewHScroll(g.output)

	mainContainer := container.NewVBox(
		topLayout,
		inputContainer,
		buttonContainer,
		yamlHeader,
		outputScroll,
	)

	return mainContainer
}

func (g *generator) buildSecret() {

	secretName := strings.TrimSpace(g.nameEntry.Text)
	secretNamespace := strings.TrimSpace(g.nameSpaceEntry.Text)

	if !utils.IsK8sNameValid(secretName) {
		// set a random generated secret name
		secretName = utils.GeneratePullSecretName()
	}

	if !utils.IsK8sNameValid(secretNamespace) {
		// clear invalid namespace
		secretNamespace = ""
	}

	secret, err := NewImagePullSecret(
		strings.TrimSpace(g.regEntry.Text),
		strings.TrimSpace(g.userEntry.Text),
		strings.TrimSpace(g.passEntry.Text),
		secretName,
		secretNamespace,
	)

	if err != nil {
		dialog.ShowError(err, g.window)
		g.saveBtn.Disable()
		return
	}

	if yaml, err := secret.ToYAML(); err != nil {
		dialog.ShowError(err, g.window)
		g.decodeBtn.Disable()
		g.saveBtn.Disable()
		g.copyBtn.Disable()
		g.clearOutputBtn.Disable()
		return
	} else {
		g.secret = secret
		g.output.SetText(yaml)
		g.window.Canvas().Refresh(g.output)
		g.decodeBtn.Enable()
		g.saveBtn.Enable()
		g.copyBtn.Enable()
		g.clearOutputBtn.Enable()
		g.storeHistory()
		g.updateEntries()
	}
}

// Checks if required input fields are filled
func (g *generator) isRequiredInputFilled() bool {
	return strings.TrimSpace(g.regEntry.Text) != "" && strings.TrimSpace(g.userEntry.Text) != "" && strings.TrimSpace(g.passEntry.Text) != ""
}

// Enables or disables the generate button based on input fields state
func (g *generator) canGenerate() {
	if g.isRequiredInputFilled() {
		g.generateBtn.Enable()
	} else {
		g.generateBtn.Disable()
	}
}

// Resets the output area and disables related buttons
func (g *generator) clearOutput() {
	g.output.SetText(DefaultOutputText)
	g.decodeBtn.Disable()
	g.saveBtn.Disable()
	g.copyBtn.Disable()
	g.clearOutputBtn.Disable()
}

// Clears the history and updates the entries
func (g *generator) clearHistory() {
	dialog.ShowConfirm("Clear History", "Do you really want to delete the history?", func(confirmed bool) {
		if confirmed {
			g.appSettings.History.Clear()
			g.updateEntries()
			g.clearHistoryBtn.Disable()
			g.toast.ShowToast("History Cleared", 3*time.Second)
		}
	}, g.window)
}

func (g *generator) copyToClipboard() {
	if g.output.Text != DefaultOutputText {
		fyne.CurrentApp().Clipboard().SetContent(g.output.Text)

		g.toast.ShowToast("Copied", 2*time.Second)
	}
}

func (g *generator) decodeOrEncodeSecret() {
	if g.secret != nil {
		if g.isDecoded {
			// back to encoded state
			yaml, err := g.secret.ToYAML()
			if err != nil {
				dialog.ShowError(err, g.window)
				return
			}
			g.output.SetText(yaml)
			g.decodeBtn.SetIcon(theme.VisibilityOffIcon())
		} else {
			// show decoded state
			decoded, err := g.secret.DecodeDockerConfig()
			if err != nil {
				dialog.ShowError(err, g.window)
				return
			}
			g.output.SetText(decoded)
			g.decodeBtn.SetIcon(theme.VisibilityIcon())
		}
		g.isDecoded = !g.isDecoded
	}
}

func (g *generator) saveDialog() {
	if g.output.Text != DefaultOutputText {
		saveDialog := dialog.NewFileSave(
			func(uriWriter fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, g.window)
					return
				}
				if uriWriter == nil {
					// cancelled
					return
				}

				defer uriWriter.Close()

				originalPath := uriWriter.URI().Path()

				// delete the empty file created by the dialog
				_ = os.Remove(originalPath)

				finalPath := utils.EnsureYAMLExt(originalPath)

				if err := utils.WriteFile(finalPath, []byte(g.output.Text)); err != nil {
					dialog.ShowError(err, g.window)
				} else {
					g.toast.ShowToast("Saved", 2*time.Second)
				}
			},
			g.window,
		)

		// Only allow .yaml and .yml files to select
		saveDialog.SetFilter(
			storage.NewExtensionFileFilter([]string{".yaml", ".yml"}),
		)

		// Set default file name
		saveDialog.SetFileName("image-pull-secret.yaml")
		saveDialog.SetTitleText("Save Image-Pull-Secret")
		saveDialog.Resize(fyne.NewSize(600.0, 400.0))
		saveDialog.Show()
	}
}

func (g *generator) changeTheme() {
	if g.appSettings.IsLightTheme() {
		g.appSettings.SetThemeVariant(theme.VariantDark)
		fyne.CurrentApp().Settings().SetTheme(&ui.ForcedVariant{Theme: theme.DefaultTheme(), Variant: g.appSettings.ThemeVariant()})
		g.themeBtn.SetIcon(ui.LightThemeIcon)
	} else {
		g.appSettings.SetThemeVariant(theme.VariantLight)
		fyne.CurrentApp().Settings().SetTheme(&ui.ForcedVariant{Theme: theme.DefaultTheme(), Variant: g.appSettings.ThemeVariant()})
		g.themeBtn.SetIcon(ui.DarkThemeIcon)
	}
}

// store current entries to history
func (g *generator) storeHistory() {
	g.appSettings.History.AddRegistry(g.regEntry.Text)
	g.appSettings.History.AddNamespace(g.nameSpaceEntry.Text)
	g.appSettings.History.AddSecretName(g.nameEntry.Text)
	if g.clearHistoryBtn.Disabled() {
		g.clearHistoryBtn.Enable()
	}
}

// update entries with latest history
func (g *generator) updateEntries() {
	g.regEntry.SetOptions(g.appSettings.History.SortedRegistries())
	g.nameSpaceEntry.SetOptions(g.appSettings.History.SortedNamespaces())
	g.nameEntry.SetOptions(g.appSettings.History.SortedNames())
}
