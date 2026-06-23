package main

import (
	"fmt"
	"image"
	"os"
	"strings"

	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
)

type PDFStampOptions struct {
	InputPDF   string
	OutputPDF  string
	StampImage string
	Pages      string
	Scale      string
	WidthMm    float64
	HeightMm   float64
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

	scale := calculateStampScale(opts)

	desc := fmt.Sprintf("pos:bc, off:0 8, rot:0, scale:%.4f abs", scale)

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

func calculateStampScale(opts PDFStampOptions) float64 {
	widthMm := opts.WidthMm
	heightMm := opts.HeightMm

	if widthMm <= 0 {
		widthMm = 90
	}
	if heightMm <= 0 {
		heightMm = 35
	}

	stampWidthPx := float64(0)
	stampHeightPx := float64(0)

	if info, err := imageInfo(opts.StampImage); err == nil {
		stampWidthPx = float64(info.Width)
		stampHeightPx = float64(info.Height)
	}

	if stampWidthPx <= 0 {
		stampWidthPx = widthMm * mmToPixels
	}
	if stampHeightPx <= 0 {
		stampHeightPx = heightMm * mmToPixels
	}

	targetWidthPt := widthMm * 72.0 / 25.4

	scale := targetWidthPt / stampWidthPx

	if scale <= 0 || scale > 10 {
		scale = 0.5
	}

	return scale
}

func imageInfo(path string) (struct {
	Width  int
	Height int
}, error) {
	f, err := os.Open(path)
	if err != nil {
		return struct {
			Width  int
			Height int
		}{}, err
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return struct {
			Width  int
			Height int
		}{}, err
	}

	return struct {
		Width  int
		Height int
	}{Width: cfg.Width, Height: cfg.Height}, nil
}
