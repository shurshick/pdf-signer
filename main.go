package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow(tr(msgWindowTitle))
	w.Resize(fyne.NewSize(950, 620))

	settings := LoadSettings()
	var stampProfile *StampProfile
	if settings.StampProfile != nil {
		stampProfile = settings.StampProfile
	} else {
		stampProfile = DefaultStampProfile()
	}

	var pdfPaths []string
	outputDir := defaultOutputDir()
	var selectedCert CertInfo

	certs, err := GetCertificates()
	if err != nil {
		dialog.ShowError(err, w)
		certs = []CertInfo{}
	}

	labels := make([]string, 0, len(certs))
	certMap := make(map[string]CertInfo)

	for _, c := range certs {
		labels = append(labels, c.Label)
		certMap[c.Label] = c
	}

	pdfLabel := widget.NewLabel(tr(msgPDFNotSelected))
	fileSummaryLabel := widget.NewLabel("")
	outLabel := widget.NewLabel(outputDir)
	certInfoLabel := widget.NewLabel(tr(msgCertNotSelected))

	reasonEntry := widget.NewEntry()
	reasonEntry.SetText(tr(msgDefaultReason))

	scaleEntry := widget.NewEntry()
	scaleEntry.SetText(fmt.Sprintf("%.2f", stampProfile.Scale))

	saveNextToSourceCheck := widget.NewCheck(tr(msgSaveNextToSource), nil)
	saveNextToSourceCheck.SetChecked(true)

	verifyAfterCheck := widget.NewCheck(tr(msgVerifyAfterSigning), nil)
	verifyAfterCheck.SetChecked(settings.VerifyAfterSigning)

	modeLabels := []string{
		tr(msgModeEmbedded),
		tr(msgModeDetached),
		tr(msgModeBoth),
	}
	var selectedMode SignMode = SignModeEmbedded
	modeSelect := widget.NewSelect(modeLabels, func(s string) {
		switch s {
		case tr(msgModeEmbedded):
			selectedMode = SignModeEmbedded
		case tr(msgModeDetached):
			selectedMode = SignModeDetached
		case tr(msgModeBoth):
			selectedMode = SignModeBoth
		}
	})
	modeSelect.SetSelected(tr(msgModeEmbedded))

	certSelect := widget.NewSelect(labels, func(s string) {
		c, ok := certMap[s]
		if !ok {
			return
		}
		selectedCert = c

		status := tr(msgReady)

		certInfoLabel.SetText(
			fmt.Sprintf(
				"%s: %s\n%s: %s\nSN: %s\nSHA1: %s\n%s",
				tr(msgOwner),
				c.SubjectCN,
				tr(msgIssuer),
				c.IssuerCN,
				c.Serial,
				c.Thumbprint,
				status,
			),
		)
	})

	updatePDFLabel := func() {
		pdfLabel.SetText(selectedPDFsText(pdfPaths))
		fileSummaryLabel.SetText(fmt.Sprintf(tr(msgFileSummary), len(pdfPaths)))
	}

	selectPDFBtn := widget.NewButton(tr(msgSelectPDF), func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if r == nil {
				return
			}
			defer r.Close()

			pdfPath := r.URI().Path()
			if !containsString(pdfPaths, pdfPath) {
				pdfPaths = append(pdfPaths, pdfPath)
			}
			updatePDFLabel()
		}, w)
	})

	clearPDFsBtn := widget.NewButton(tr(msgClearPDFs), func() {
		pdfPaths = nil
		updatePDFLabel()
	})

	selectOutBtn := widget.NewButton(tr(msgBrowse), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}
			outputDir = uri.Path()
			outLabel.SetText(outputDir)
		}, w)
	})

	diagnosticsBtn := widget.NewButton(tr(msgDiagnostics), func() {
		showDiagnosticsDialog(w)
	})

	verifyBtn := widget.NewButton(tr(msgVerifySignature), func() {
		if len(pdfPaths) == 0 {
			dialog.ShowInformation(tr(msgError), tr(msgChoosePDF), w)
			return
		}
		showVerificationDialog(w, pdfPaths)
	})

	stampEditorBtn := widget.NewButton(tr(msgStampEditor), func() {
		showStampEditor(w, stampProfile, func(p *StampProfile) {
			stampProfile = p
			settings.StampProfile = p
			SaveSettings(settings)
			scaleEntry.Text = fmt.Sprintf("%.2f", p.Scale)
			scaleEntry.Refresh()
		})
	})

	exportSettingsBtn := widget.NewButton(tr(msgExportSettings), func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			settings.StampProfile = stampProfile
			settings.VerifyAfterSigning = verifyAfterCheck.Checked
			if err := ExportSettings(writer.URI().Path(), settings); err != nil {
				dialog.ShowError(err, w)
			} else {
				dialog.ShowInformation(tr(msgDone), tr(msgSettingsExported), w)
			}
		}, w)
	})

	importSettingsBtn := widget.NewButton(tr(msgImportSettings), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			imported, err := ImportSettings(reader.URI().Path())
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			settings = imported
			if imported.StampProfile != nil {
				stampProfile = imported.StampProfile
			}
			verifyAfterCheck.SetChecked(imported.VerifyAfterSigning)
			scaleEntry.SetText(fmt.Sprintf("%.2f", stampProfile.Scale))
			dialog.ShowInformation(tr(msgDone), tr(msgSettingsImported), w)
		}, w)
	})

	runBtn := widget.NewButton(tr(msgSignAndStamp), func() {
		if len(pdfPaths) == 0 {
			dialog.ShowInformation(tr(msgError), tr(msgChoosePDF), w)
			return
		}
		if selectedCert.SubjectCN == "" && selectedCert.Serial == "" {
			dialog.ShowInformation(tr(msgError), tr(msgChooseCert), w)
			return
		}
		if !saveNextToSourceCheck.Checked && strings.TrimSpace(outputDir) == "" {
			dialog.ShowInformation(tr(msgError), tr(msgChooseOutPDF), w)
			return
		}

		LogInfo("Signing started: " + fmt.Sprintf("%d file(s), mode=%d", len(pdfPaths), selectedMode))

		reason := strings.TrimSpace(reasonEntry.Text)
		if reason == "" {
			reason = tr(msgDefaultReason)
		}

		signer := NativeSigner{}
		results := make([]string, 0, len(pdfPaths))

		for _, pdfPath := range pdfPaths {
			outputPath, err := stampedOutputPath(pdfPath, outputDir, saveNextToSourceCheck.Checked)
			if err != nil {
				dialog.ShowError(errors.New(FriendlyErrorMessage(err)), w)
				return
			}

			sigPath := signatureOutputPath(outputPath)
			stampData := NewStampData(selectedCert, sigPath, reason)

			if selectedMode == SignModeEmbedded {
				stampData.SignatureFN = ""
			}

			stampPNG, cleanupStamp, err := createTempStampPath()
			if err != nil {
				dialog.ShowError(errors.New(FriendlyErrorMessage(err)), w)
				return
			}

			err = CreateStampImage(stampPNG, stampData)
			if err != nil {
				cleanupStamp()
				dialog.ShowError(errors.New(FriendlyErrorMessage(err)), w)
				return
			}

			pages := stampProfile.Pages
			if pages == "" {
				pages = "1-"
			}

			err = ApplyPDFStamp(PDFStampOptions{
				InputPDF:   pdfPath,
				OutputPDF:  outputPath,
				StampImage: stampPNG,
				Pages:      pages,
				Scale:      fmt.Sprintf("%.2f", stampProfile.Scale),
			})
			cleanupStamp()
			if err != nil {
				dialog.ShowError(errors.New(FriendlyErrorMessage(err)), w)
				return
			}

			result := fmt.Sprintf("%s -> %s", filepath.Base(pdfPath), outputPath)

			switch selectedMode {
			case SignModeEmbedded:
				embedRes, err := signer.SignFileEmbedded(outputPath, selectedCert)
				if err != nil {
					dialog.ShowError(errors.New(FriendlyErrorMessage(err)), w)
					return
				}
				result += "\n" + tr(msgEmbeddedSignature) + ": " + embedRes.SignedPDFPath

			case SignModeDetached:
				detRes, err := signer.SignFileTo(outputPath, selectedCert, sigPath)
				if err != nil {
					dialog.ShowError(errors.New(FriendlyErrorMessage(err)), w)
					return
				}
				result += "\n" + tr(msgSignature) + ": " + detRes.SignaturePath

			case SignModeBoth:
				detRes, err := signer.SignFileTo(outputPath, selectedCert, sigPath)
				if err != nil {
					dialog.ShowError(errors.New(FriendlyErrorMessage(err)), w)
					return
				}
				result += "\n" + tr(msgSignature) + ": " + detRes.SignaturePath

				embedRes, err := signer.SignFileEmbedded(outputPath, selectedCert)
				if err != nil {
					dialog.ShowError(errors.New(FriendlyErrorMessage(err)), w)
					return
				}
				result += "\n" + tr(msgEmbeddedSignature) + ": " + embedRes.SignedPDFPath
			}

			LogInfo("Signed: " + filepath.Base(pdfPath))
			results = append(results, result)

			if verifyAfterCheck.Checked {
				verifier := &SignatureVerifier{}
				report := verifier.VerifyDetached(sigPath)
				if report.Status == VerifyInvalid {
					dialog.ShowWarning(tr(msgSignatureInvalid)+": "+filepath.Base(pdfPath), w)
				}
			}
		}

		LogInfo("Signing completed: " + fmt.Sprintf("%d file(s)", len(results)))
		dialog.ShowInformation(
			tr(msgDone),
			fmt.Sprintf("%s: %d\n%s", tr(msgProcessedFiles), len(results), strings.Join(results, "\n\n")),
			w,
		)
	})

	aboutBtn := widget.NewButton(tr(msgAbout), func() {
		showAboutDialog(w)
	})

	settingsPanel := container.NewVBox(
		widget.NewLabel(tr(msgSelectOutputFolder)+":"),
		container.NewBorder(nil, nil, nil, selectOutBtn, outLabel),
		saveNextToSourceCheck,
		verifyAfterCheck,
		widget.NewLabel(tr(msgBatchOutputNote)),
		widget.NewSeparator(),
		widget.NewLabel(tr(msgCertificate)+":"),
		certSelect,
		certInfoLabel,
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem(tr(msgReason), reasonEntry),
			widget.NewFormItem(tr(msgScale), scaleEntry),
			widget.NewFormItem(tr(msgSigningMode), modeSelect),
		),
		widget.NewSeparator(),
		container.NewHBox(exportSettingsBtn, importSettingsBtn),
	)

	actionPanel := container.NewHBox(runBtn, verifyBtn, stampEditorBtn, diagnosticsBtn, aboutBtn)

	filesPanel := container.NewVBox(
		container.NewHBox(selectPDFBtn, clearPDFsBtn),
		pdfLabel,
		fileSummaryLabel,
	)

	root := container.NewBorder(
		filesPanel,
		container.NewVBox(widget.NewSeparator(), actionPanel),
		nil,
		nil,
		settingsPanel,
	)

	w.SetContent(root)
	w.ShowAndRun()
}

func showDiagnosticsDialog(w fyne.Window) {
	diagWindow := fyne.CurrentApp().NewWindow(tr(msgDiagnosticsTitle))
	diagWindow.Resize(fyne.NewSize(600, 400))

	reportBox := widget.NewMultiLineEntry()
	reportBox.SetPlaceHolder(tr(msgDiagnosticsRunning))
	reportBox.Disable()

	runBtn := widget.NewButton(tr(msgRunDiagnostics), func() {
		diagnostics := &CryptoProDiagnostics{}
		report := diagnostics.Run()
		reportBox.SetText(diagnostics.FormatReport(report))
	})

	copyBtn := widget.NewButton(tr(msgCopyReport), func() {
		diagWindow.Clipboard().SetContent(reportBox.Text)
	})

	saveBtn := widget.NewButton(tr(msgSaveReport), func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			writer.Write([]byte(reportBox.Text))
		}, diagWindow)
	})

	openLogsBtn := widget.NewButton(tr(msgOpenLogsFolder), func() {
		dir := LogDir()
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}
	})

	buttons := container.NewHBox(runBtn, copyBtn, saveBtn, openLogsBtn)
	content := container.NewBorder(buttons, nil, nil, nil, reportBox)
	diagWindow.SetContent(content)
	diagWindow.Show()

	runBtn.OnTapped()
}

func showVerificationDialog(w fyne.Window, initialFiles []string) {
	verifyWindow := fyne.CurrentApp().NewWindow(tr(msgVerificationTitle))
	verifyWindow.Resize(fyne.NewSize(600, 400))

	var files []string
	files = append(files, initialFiles...)

	fileList := widget.NewList(
		func() int { return len(files) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(filepath.Base(files[id]))
		},
	)

	addBtn := widget.NewButton(tr(msgSelectPDF), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			files = append(files, reader.URI().Path())
			fileList.Refresh()
		}, verifyWindow)
	})

	clearBtn := widget.NewButton(tr(msgClearPDFs), func() {
		files = nil
		fileList.Refresh()
	})

	reportBox := widget.NewMultiLineEntry()
	reportBox.Disable()

	verifyBtn := widget.NewButton(tr(msgRunVerification), func() {
		if len(files) == 0 {
			return
		}
		verifier := &SignatureVerifier{}
		reports := make([]*VerificationReport, 0, len(files))
		for _, f := range files {
			if strings.HasSuffix(f, ".sig") {
				reports = append(reports, verifier.VerifyDetached(f))
			} else {
				reports = append(reports, verifier.VerifyEmbedded(f))
			}
		}
		reportBox.SetText(verifier.FormatReport(reports))
	})

	exportBtn := widget.NewButton(tr(msgExportTxt), func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			writer.Write([]byte(reportBox.Text))
		}, verifyWindow)
	})

	buttons := container.NewHBox(addBtn, clearBtn, verifyBtn, exportBtn)
	listPanel := container.NewBorder(buttons, nil, nil, nil, fileList)
	split := container.NewHSplit(listPanel, reportBox)
	split.SetOffset(0.35)

	verifyWindow.SetContent(split)
	verifyWindow.Show()
}

func selectedPDFsText(paths []string) string {
	if len(paths) == 0 {
		return tr(msgPDFNotSelected)
	}

	lines := make([]string, 0, len(paths)+1)
	lines = append(lines, fmt.Sprintf("%s: %d", tr(msgSelectedPDFs), len(paths)))
	for _, path := range paths {
		lines = append(lines, filepath.Base(path))
	}
	return strings.Join(lines, "\n")
}

func createTempStampPath() (string, func(), error) {
	f, err := os.CreateTemp("", "pdfsigner-stamp-*.png")
	if err != nil {
		return "", func() {}, err
	}

	path := f.Name()
	if err := f.Close(); err != nil {
		_ = os.Remove(path)
		return "", func() {}, err
	}

	return path, func() { _ = os.Remove(path) }, nil
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
