package main

import (
	"fmt"
	"strings"

	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
)

type PDFStampOptions struct {
	InputPDF   string
	OutputPDF  string
	StampImage string
	Pages      string
	Scale      string
}

func ApplyPDFStamp(opts PDFStampOptions) error {
	if strings.TrimSpace(opts.InputPDF) == "" {
		return fmt.Errorf("%s", tr(msgInputPDFMissing))
	}
	if strings.TrimSpace(opts.OutputPDF) == "" {
		return fmt.Errorf("%s", tr(msgOutputPDFMissing))
	}
	if strings.TrimSpace(opts.StampImage) == "" {
		return fmt.Errorf("%s", tr(msgStampFileMissing))
	}

	selectedPages := []string{"1"}
	if p := strings.TrimSpace(opts.Pages); p != "" {
		selectedPages = []string{p}
	}

	scale := "0.96"
	if s := strings.TrimSpace(opts.Scale); s != "" {
		scale = s
	}

	// Горизонтальный штамп по центру внизу страницы,
	// почти на всю ширину, без поворота.
	desc := fmt.Sprintf("pos:bc, off:0 8, rot:0, scale:%s rel", scale)

	// onTop=true => именно stamp, а не watermark под контентом.
	if err := pdfapi.AddImageWatermarksFile(
		opts.InputPDF,
		opts.OutputPDF,
		selectedPages,
		true,
		opts.StampImage,
		desc,
		nil,
	); err != nil {
		return fmt.Errorf("%s: %w", tr(msgStampPDFError), err)
	}

	return nil
}
