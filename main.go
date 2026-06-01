package main

import (
	"fmt"
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
	w.Resize(fyne.NewSize(950, 380))

	var pdfPaths []string
	var outPDF string
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
	outLabel := widget.NewLabel(tr(msgOutPDFNotSelected))
	certInfoLabel := widget.NewLabel(tr(msgCertNotSelected))

	scaleEntry := widget.NewEntry()
	scaleEntry.SetText("0.96")

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

			if outPDF == "" {
				outPDF = defaultStampedPath(pdfPath)
				outLabel.SetText(outPDF)
			}
		}, w)
	})

	clearPDFsBtn := widget.NewButton(tr(msgClearPDFs), func() {
		pdfPaths = nil
		outPDF = ""
		pdfLabel.SetText(tr(msgPDFNotSelected))
		outLabel.SetText(tr(msgOutPDFNotSelected))
	})

	selectOutBtn := widget.NewButton(tr(msgSavePDFAs), func() {
		dialog.ShowFileSave(func(wc fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if wc == nil {
				return
			}
			outPDF = wc.URI().Path()
			outLabel.SetText(outPDF)
			_ = wc.Close()
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
		if len(pdfPaths) == 1 && outPDF == "" {
			dialog.ShowInformation(tr(msgError), tr(msgChooseOutPDF), w)
			return
		}

		signer := NativeSigner{}
		results := make([]string, 0, len(pdfPaths))

		for _, pdfPath := range pdfPaths {
			outputPath := defaultStampedPath(pdfPath)
			if len(pdfPaths) == 1 {
				outputPath = outPDF
			}

			signRes, err := signer.SignFile(pdfPath, selectedCert)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			stampData := NewStampData(selectedCert, signRes.SignaturePath)
			stampPNG := filepath.Join(filepath.Dir(outputPath), "ep_stamp.png")

			err = CreateStampImage(stampPNG, stampData)
			if err != nil {
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
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			results = append(results, fmt.Sprintf("%s -> %s", filepath.Base(pdfPath), outputPath))
		}

		dialog.ShowInformation(
			tr(msgDone),
			fmt.Sprintf("%s: %d\n%s", tr(msgProcessedFiles), len(results), strings.Join(results, "\n")),
			w,
		)
	})

	form := container.NewVBox(
		selectPDFBtn,
		clearPDFsBtn,
		pdfLabel,
		selectOutBtn,
		outLabel,
		widget.NewLabel(tr(msgBatchOutputNote)),
		widget.NewSeparator(),

		widget.NewLabel(tr(msgCertificate)+":"),
		certSelect,
		certInfoLabel,
		widget.NewSeparator(),

		widget.NewForm(
			widget.NewFormItem(tr(msgScale), scaleEntry),
		),

		runBtn,
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

func defaultStampedPath(pdfPath string) string {
	ext := filepath.Ext(pdfPath)
	base := strings.TrimSuffix(pdfPath, ext)
	return base + "_stamped.pdf"
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
