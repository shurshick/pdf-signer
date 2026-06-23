package main

import (
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
	w.Resize(fyne.NewSize(950, 580))

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
	outLabel := widget.NewLabel(outputDir)
	certInfoLabel := widget.NewLabel(tr(msgCertNotSelected))

	reasonEntry := widget.NewEntry()
	reasonEntry.SetText(tr(msgDefaultReason))

	scaleEntry := widget.NewEntry()
	scaleEntry.SetText("0.96")

	saveNextToSourceCheck := widget.NewCheck(tr(msgSaveNextToSource), nil)
	saveNextToSourceCheck.SetChecked(true)

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

		certInfoLabel.SetText(
			fmt.Sprintf(
				"%s: %s\n%s: %s\nSN: %s\nSHA1: %s",
				tr(msgOwner),
				c.SubjectCN,
				tr(msgIssuer),
				c.IssuerCN,
				c.Serial,
				c.Thumbprint,
			),
		)
	})

	updatePDFLabel := func() {
		pdfLabel.SetText(selectedPDFsText(pdfPaths))
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
		pdfLabel.SetText(tr(msgPDFNotSelected))
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

		reason := strings.TrimSpace(reasonEntry.Text)
		if reason == "" {
			reason = tr(msgDefaultReason)
		}

		signer := NativeSigner{}
		results := make([]string, 0, len(pdfPaths))

		for _, pdfPath := range pdfPaths {
			outputPath, err := stampedOutputPath(pdfPath, outputDir, saveNextToSourceCheck.Checked)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			sigPath := signatureOutputPath(outputPath)
			stampData := NewStampData(selectedCert, sigPath, reason)

			if selectedMode == SignModeEmbedded {
				stampData.SignatureFN = ""
			}

			stampPNG, cleanupStamp, err := createTempStampPath()
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			err = CreateStampImage(stampPNG, stampData)
			if err != nil {
				cleanupStamp()
				dialog.ShowError(err, w)
				return
			}

			err = ApplyPDFStamp(PDFStampOptions{
				InputPDF:   pdfPath,
				OutputPDF:  outputPath,
				StampImage: stampPNG,
				Pages:      "1-",
				Scale:      scaleEntry.Text,
			})
			cleanupStamp()
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			result := fmt.Sprintf("%s -> %s", filepath.Base(pdfPath), outputPath)

			switch selectedMode {
			case SignModeEmbedded:
				embedRes, err := signer.SignFileEmbedded(outputPath, selectedCert)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				result += "\n" + tr(msgEmbeddedSignature) + ": " + embedRes.SignedPDFPath

			case SignModeDetached:
				detRes, err := signer.SignFileTo(outputPath, selectedCert, sigPath)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				result += "\n" + tr(msgSignature) + ": " + detRes.SignaturePath

			case SignModeBoth:
				detRes, err := signer.SignFileTo(outputPath, selectedCert, sigPath)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				result += "\n" + tr(msgSignature) + ": " + detRes.SignaturePath

				embedRes, err := signer.SignFileEmbedded(outputPath, selectedCert)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				result += "\n" + tr(msgEmbeddedSignature) + ": " + embedRes.SignedPDFPath
			}

			results = append(results, result)
		}

		dialog.ShowInformation(
			tr(msgDone),
			fmt.Sprintf("%s: %d\n%s", tr(msgProcessedFiles), len(results), strings.Join(results, "\n\n")),
			w,
		)
	})

	aboutBtn := widget.NewButton(tr(msgAbout), func() {
		showAboutDialog(w)
	})

	form := container.NewVBox(
		selectPDFBtn,
		clearPDFsBtn,
		pdfLabel,
		widget.NewSeparator(),

		widget.NewLabel(tr(msgSelectOutputFolder)+":"),
		container.NewBorder(nil, nil, nil, selectOutBtn, outLabel),
		saveNextToSourceCheck,
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

		container.NewHBox(runBtn, aboutBtn),
	)

	w.SetContent(form)
	w.ShowAndRun()
}

func selectedPDFsText(paths []string) string {
	if len(paths) == 0 {
		return tr(msgPDFNotSelected)
	}

	lines := make([]string, 0, len(paths)+1)
	lines = append(lines, fmt.Sprintf("%s: %d", tr(msgSelectedPDFs), len(paths)))
	for _, path := range paths {
		lines = append(lines, path)
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
