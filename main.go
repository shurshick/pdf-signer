package main

import (
	"fmt"
	"path/filepath"

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

	var pdfPath string
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

			pdfPath = r.URI().Path()
			pdfLabel.SetText(pdfPath)

			if outPDF == "" {
				ext := filepath.Ext(pdfPath)
				base := pdfPath[:len(pdfPath)-len(ext)]
				outPDF = base + "_stamped.pdf"
				outLabel.SetText(outPDF)
			}
		}, w)
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
		if pdfPath == "" {
			dialog.ShowInformation(tr(msgError), tr(msgChoosePDF), w)
			return
		}
		if selectedCert.SubjectCN == "" && selectedCert.Serial == "" {
			dialog.ShowInformation(tr(msgError), tr(msgChooseCert), w)
			return
		}
		if outPDF == "" {
			dialog.ShowInformation(tr(msgError), tr(msgChooseOutPDF), w)
			return
		}

		signer := NativeSigner{}
		signRes, err := signer.SignFile(pdfPath, selectedCert)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		stampData := NewStampData(selectedCert, signRes.SignaturePath)
		stampPNG := filepath.Join(filepath.Dir(outPDF), "ep_stamp.png")

		err = CreateStampImage(stampPNG, stampData)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		err = ApplyPDFStamp(PDFStampOptions{
			InputPDF:   pdfPath,
			OutputPDF:  outPDF,
			StampImage: stampPNG,
			Pages:      "1-",
			Scale:      scaleEntry.Text,
		})
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		dialog.ShowInformation(
			tr(msgDone),
			fmt.Sprintf("%s: %s\n%s: %s",
				tr(msgSignature),
				signRes.SignaturePath,
				tr(msgStampedPDF),
				outPDF,
			),
			w,
		)
	})

	form := container.NewVBox(
		selectPDFBtn,
		pdfLabel,
		selectOutBtn,
		outLabel,
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
