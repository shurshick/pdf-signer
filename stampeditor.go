package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func showStampEditor(w fyne.Window, profile *StampProfile, onSave func(*StampProfile)) {
	editorWindow := fyne.CurrentApp().NewWindow(tr(msgStampEditorTitle))
	editorWindow.Resize(fyne.NewSize(700, 550))

	templateSelect := widget.NewSelect([]string{
		tr(msgTemplateMinimal),
		tr(msgTemplateStandard),
		tr(msgTemplateDetailed),
	}, nil)
	templateSelect.SetSelected(profileToTemplateLabel(profile.TemplateName))

	pagesEntry := widget.NewEntry()
	pagesEntry.SetText(profile.Pages)

	posSelect := widget.NewSelect([]string{
		tr(msgPosBottomRight),
		tr(msgPosBottomLeft),
		tr(msgPosTopRight),
		tr(msgPosTopLeft),
	}, nil)
	posSelect.SetSelected(profileToPosLabel(profile.PositionMode))

	widthEntry := widget.NewEntry()
	widthEntry.SetText(fmt.Sprintf("%.0f", profile.WidthMm))

	heightEntry := widget.NewEntry()
	heightEntry.SetText(fmt.Sprintf("%.0f", profile.HeightMm))

	fontSizeEntry := widget.NewEntry()
	fontSizeEntry.SetText(fmt.Sprintf("%.1f", profile.FontSize))

	opacityEntry := widget.NewEntry()
	opacityEntry.SetText(fmt.Sprintf("%.0f", profile.Opacity*100))

	scaleEntry := widget.NewEntry()
	scaleEntry.SetText(fmt.Sprintf("%.2f", profile.Scale))

	checkOwner := widget.NewCheck(tr(msgOwner), nil)
	checkOwner.SetChecked(profile.IncludeOwner)
	checkIssuer := widget.NewCheck(tr(msgIssuer), nil)
	checkIssuer.SetChecked(profile.IncludeIssuer)
	checkDate := widget.NewCheck(tr(msgDate), nil)
	checkDate.SetChecked(profile.IncludeDate)
	checkReason := widget.NewCheck(tr(msgReason), nil)
	checkReason.SetChecked(profile.IncludeReason)
	checkSerial := widget.NewCheck(tr(msgSerialNumber), nil)
	checkSerial.SetChecked(profile.IncludeSerial)

	autoPlaceCheck := widget.NewCheck(tr(msgAutoPlaceStamp), nil)
	autoPlaceCheck.SetChecked(profile.AutoPlace)

	logoPathEntry := widget.NewEntry()
	logoPathEntry.SetText(profile.LogoPath)
	logoPathEntry.SetPlaceHolder(tr(msgLogoPath))

	logoScaleEntry := widget.NewEntry()
	logoScaleEntry.SetText(fmt.Sprintf("%d", profile.LogoScale))

	chooseLogoBtn := widget.NewButton(tr(msgChooseLogo), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			logoPathEntry.SetText(reader.URI().Path())
		}, editorWindow)
	})

	removeLogoBtn := widget.NewButton(tr(msgRemoveLogo), func() {
		logoPathEntry.SetText("")
	})

	warningLabel := widget.NewLabel("")
	warningLabel.Wrapping = fyne.TextWrapWord

	previewLabel := widget.NewLabel("")

	updatePreview := func() {
		profile := readEditorProfile(templateSelect, pagesEntry, posSelect,
			widthEntry, heightEntry, fontSizeEntry, opacityEntry, scaleEntry,
			checkOwner, checkIssuer, checkDate, checkReason, checkSerial,
			autoPlaceCheck, logoPathEntry, logoScaleEntry)
		profile.Normalize()

		cert := CertInfo{SubjectCN: "Preview User", IssuerCN: "Preview CA", Serial: "12345", Thumbprint: "AABBCCDD"}
		text := BuildStampTextFromProfile(profile, cert, tr(msgDefaultReason))
		previewLabel.SetText(text)

		warnings := validateStampProfile(profile)
		if len(warnings) > 0 {
			warningLabel.SetText(strings.Join(warnings, "\n"))
		} else {
			warningLabel.SetText("")
		}
	}

	templateSelect.OnChanged = func(s string) {
		profiles := BuiltInProfiles()
		key := labelToProfileKey(s)
		if p, ok := profiles[key]; ok {
			pagesEntry.SetText(p.Pages)
			posSelect.SetSelected(profileToPosLabel(p.PositionMode))
			widthEntry.SetText(fmt.Sprintf("%.0f", p.WidthMm))
			heightEntry.SetText(fmt.Sprintf("%.0f", p.HeightMm))
			fontSizeEntry.SetText(fmt.Sprintf("%.1f", p.FontSize))
			opacityEntry.SetText(fmt.Sprintf("%.0f", p.Opacity*100))
			scaleEntry.SetText(fmt.Sprintf("%.2f", p.Scale))
			checkOwner.SetChecked(p.IncludeOwner)
			checkIssuer.SetChecked(p.IncludeIssuer)
			checkDate.SetChecked(p.IncludeDate)
			checkReason.SetChecked(p.IncludeReason)
			checkSerial.SetChecked(p.IncludeSerial)
			autoPlaceCheck.SetChecked(p.AutoPlace)
		}
		updatePreview()
	}

	for _, entry := range []*widget.Entry{pagesEntry, widthEntry, heightEntry, fontSizeEntry, opacityEntry, scaleEntry} {
		entry.OnChanged = func(string) { updatePreview() }
	}
	for _, check := range []*widget.Check{checkOwner, checkIssuer, checkDate, checkReason, checkSerial, autoPlaceCheck} {
		check.OnChanged = func(bool) { updatePreview() }
	}
	posSelect.OnChanged = func(string) { updatePreview() }

	loadProfileBtn := widget.NewButton(tr(msgLoadProfile), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			imported, err := ImportSettings(reader.URI().Path())
			if err != nil {
				dialog.ShowError(err, editorWindow)
				return
			}
			p := imported.StampProfile
			if p != nil {
				templateSelect.SetSelected(profileToTemplateLabel(p.TemplateName))
				pagesEntry.SetText(p.Pages)
				posSelect.SetSelected(profileToPosLabel(p.PositionMode))
				widthEntry.SetText(fmt.Sprintf("%.0f", p.WidthMm))
				heightEntry.SetText(fmt.Sprintf("%.0f", p.HeightMm))
				fontSizeEntry.SetText(fmt.Sprintf("%.1f", p.FontSize))
				opacityEntry.SetText(fmt.Sprintf("%.0f", p.Opacity*100))
				scaleEntry.SetText(fmt.Sprintf("%.2f", p.Scale))
				checkOwner.SetChecked(p.IncludeOwner)
				checkIssuer.SetChecked(p.IncludeIssuer)
				checkDate.SetChecked(p.IncludeDate)
				checkReason.SetChecked(p.IncludeReason)
				checkSerial.SetChecked(p.IncludeSerial)
				autoPlaceCheck.SetChecked(p.AutoPlace)
				logoPathEntry.SetText(p.LogoPath)
				logoScaleEntry.SetText(fmt.Sprintf("%d", p.LogoScale))
				updatePreview()
			}
		}, editorWindow)
	})

	saveProfileBtn := widget.NewButton(tr(msgSaveProfile), func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			profile := readEditorProfile(templateSelect, pagesEntry, posSelect,
				widthEntry, heightEntry, fontSizeEntry, opacityEntry, scaleEntry,
				checkOwner, checkIssuer, checkDate, checkReason, checkSerial,
				autoPlaceCheck, logoPathEntry, logoScaleEntry)
			profile.Normalize()
			settings := &ApplicationSettings{StampProfile: profile}
			if err := ExportSettings(writer.URI().Path(), settings); err != nil {
				dialog.ShowError(err, editorWindow)
			}
		}, editorWindow)
	})

	okBtn := widget.NewButton(tr(msgDone), func() {
		profile := readEditorProfile(templateSelect, pagesEntry, posSelect,
			widthEntry, heightEntry, fontSizeEntry, opacityEntry, scaleEntry,
			checkOwner, checkIssuer, checkDate, checkReason, checkSerial,
			autoPlaceCheck, logoPathEntry, logoScaleEntry)
		profile.Normalize()
		onSave(profile)
		editorWindow.Close()
	})

	cancelBtn := widget.NewButton(tr(msgClose), func() {
		editorWindow.Close()
	})

	settingsPanel := container.NewVBox(
		widget.NewLabel(tr(msgStampTemplate)+":"),
		templateSelect,
		widget.NewSeparator(),
		widget.NewLabel(tr(msgStampPages)+":"),
		pagesEntry,
		widget.NewLabel(tr(msgStampPosition)+":"),
		posSelect,
		widget.NewSeparator(),
		widget.NewLabel(tr(msgStampSize)+":"),
		container.NewHBox(
			widget.NewLabel(tr(msgWidthMm)+":"), widthEntry,
			widget.NewLabel(tr(msgHeightMm)+":"), heightEntry,
		),
		widget.NewLabel(tr(msgFontSize)+":"),
		container.NewHBox(fontSizeEntry,
			widget.NewLabel(tr(msgOpacity)+":"), opacityEntry,
			widget.NewLabel(tr(msgScale)+":"), scaleEntry,
		),
		widget.NewSeparator(),
		checkOwner, checkIssuer, checkDate, checkReason, checkSerial,
		widget.NewSeparator(),
		autoPlaceCheck,
		container.NewHBox(
			widget.NewLabel(tr(msgLogoPath)+":"),
			logoPathEntry,
			chooseLogoBtn,
			removeLogoBtn,
		),
		widget.NewLabel(tr(msgLogoScale)+":"),
		logoScaleEntry,
		widget.NewSeparator(),
		container.NewHBox(loadProfileBtn, saveProfileBtn),
		warningLabel,
		container.NewHBox(okBtn, cancelBtn),
	)

	previewPanel := container.NewVBox(
		widget.NewLabel(tr(msgPreview)+":"),
		previewLabel,
	)

	split := container.NewHSplit(settingsPanel, previewPanel)
	split.SetOffset(0.55)

	editorWindow.SetContent(split)
	editorWindow.Show()
	updatePreview()
}

func readEditorProfile(templateSelect *widget.Select, pagesEntry *widget.Entry,
	posSelect *widget.Select, widthEntry, heightEntry, fontSizeEntry, opacityEntry, scaleEntry *widget.Entry,
	checkOwner, checkIssuer, checkDate, checkReason, checkSerial, autoPlaceCheck *widget.Check,
	logoPathEntry, logoScaleEntry *widget.Entry) *StampProfile {

	profile := DefaultStampProfile()
	profile.TemplateName = labelToProfileKey(templateSelect.Selected)
	profile.Pages = pagesEntry.Text
	profile.PositionMode = posLabelToProfile(posSelect.Selected)
	fmt.Sscanf(widthEntry.Text, "%f", &profile.WidthMm)
	fmt.Sscanf(heightEntry.Text, "%f", &profile.HeightMm)
	fmt.Sscanf(fontSizeEntry.Text, "%f", &profile.FontSize)
	var opacity float64
	fmt.Sscanf(opacityEntry.Text, "%f", &opacity)
	profile.Opacity = opacity / 100.0
	fmt.Sscanf(scaleEntry.Text, "%f", &profile.Scale)
	profile.IncludeOwner = checkOwner.Checked
	profile.IncludeIssuer = checkIssuer.Checked
	profile.IncludeDate = checkDate.Checked
	profile.IncludeReason = checkReason.Checked
	profile.IncludeSerial = checkSerial.Checked
	profile.AutoPlace = autoPlaceCheck.Checked
	profile.LogoPath = logoPathEntry.Text
	fmt.Sscanf(logoScaleEntry.Text, "%d", &profile.LogoScale)
	return profile
}

func validateStampProfile(profile *StampProfile) []string {
	var warnings []string

	sizeErrors := ValidateStampSize(profile.WidthMm, profile.HeightMm, profile.FontSize)
	if len(sizeErrors) > 0 {
		return sizeErrors
	}

	cert := CertInfo{SubjectCN: "Test User", IssuerCN: "Test CA", Serial: "12345678", Thumbprint: "AABBCCDD", NotBefore: time.Now(), NotAfter: time.Now().AddDate(1, 0, 0)}
	text := BuildStampTextFromProfile(profile, cert, "Test")
	lines := strings.Split(text, "\n")
	neededHeight := float64(len(lines)+1) * (profile.FontSize * 1.5)
	if neededHeight > profile.HeightMm*2.835 {
		warnings = append(warnings, tr(msgStampTooSmall))
	}

	if profile.Opacity < 0.55 {
		warnings = append(warnings, tr(msgOpacityLow))
	}

	if profile.LogoPath != "" {
		if _, err := os.Stat(profile.LogoPath); os.IsNotExist(err) {
			warnings = append(warnings, tr(msgLogoMissing))
		}
		if info, err := os.Stat(profile.LogoPath); err == nil && info.Size() > 1024*1024 {
			warnings = append(warnings, tr(msgLogoTooLarge))
		}
	}

	return warnings
}

func BuildStampTextFromProfile(profile *StampProfile, cert CertInfo, reason string) string {
	var parts []string

	parts = append(parts, tr(msgGostHeader))

	if profile.IncludeOwner && cert.SubjectCN != "" {
		parts = append(parts, tr(msgOwner)+": "+cert.SubjectCN)
	}
	if profile.IncludeIssuer && cert.IssuerCN != "" {
		parts = append(parts, tr(msgIssuer)+": "+cert.IssuerCN)
	}
	if profile.IncludeDate {
		parts = append(parts, tr(msgDate)+": "+time.Now().Format("02.01.2006"))
	}
	if profile.IncludeReason && reason != "" {
		parts = append(parts, tr(msgReason)+": "+reason)
	}
	if profile.IncludeSerial && cert.Serial != "" {
		parts = append(parts, tr(msgSerialNumber)+": "+cert.Serial)
	}

	if cert.NotBefore.IsZero() == false && cert.NotAfter.IsZero() == false {
		validFrom := cert.NotBefore.Format("02.01.2006")
		validTo := cert.NotAfter.Format("02.01.2006")
		parts = append(parts, tr(msgGostValidity)+": "+validFrom+" - "+validTo)
	}

	return strings.Join(parts, "\n")
}

func profileToPosLabel(mode string) string {
	switch mode {
	case "BottomRight":
		return tr(msgPosBottomRight)
	case "BottomLeft":
		return tr(msgPosBottomLeft)
	case "TopRight":
		return tr(msgPosTopRight)
	case "TopLeft":
		return tr(msgPosTopLeft)
	default:
		return tr(msgPosBottomRight)
	}
}

func posLabelToProfile(label string) string {
	switch label {
	case tr(msgPosBottomRight):
		return "BottomRight"
	case tr(msgPosBottomLeft):
		return "BottomLeft"
	case tr(msgPosTopRight):
		return "TopRight"
	case tr(msgPosTopLeft):
		return "TopLeft"
	default:
		return "BottomRight"
	}
}

func profileToTemplateLabel(name string) string {
	switch name {
	case "minimal":
		return tr(msgTemplateMinimal)
	case "standard":
		return tr(msgTemplateStandard)
	case "detailed":
		return tr(msgTemplateDetailed)
	default:
		return tr(msgTemplateStandard)
	}
}

func labelToProfileKey(label string) string {
	switch label {
	case tr(msgTemplateMinimal):
		return "minimal"
	case tr(msgTemplateStandard):
		return "standard"
	case tr(msgTemplateDetailed):
		return "detailed"
	default:
		return "standard"
	}
}
