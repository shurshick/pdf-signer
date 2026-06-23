package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	stampMinWidthMm  = 60
	stampMinHeightMm = 20
	stampMinFontPt   = 6
	stampBlueR       = 0
	stampBlueG       = 74
	stampBlueB       = 173
	mmToPixels       = 3.78
)

func CreateStampImage(path string, d StampData) error {
	profile := DefaultStampProfile()
	return CreateStampImageWithProfile(path, d, profile)
}

func CreateStampImageWithProfile(path string, d StampData, profile *StampProfile) error {
	profile.Normalize()

	widthMm := profile.WidthMm
	heightMm := profile.HeightMm

	if widthMm < stampMinWidthMm {
		widthMm = stampMinWidthMm
	}
	if heightMm < stampMinHeightMm {
		heightMm = stampMinHeightMm
	}

	w := int(float64(widthMm) * mmToPixels)
	h := int(float64(heightMm) * mmToPixels)
	if w < 400 {
		w = 400
	}
	if h < 100 {
		h = 100
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	blue := color.RGBA{stampBlueR, stampBlueG, stampBlueB, 255}
	lightBlue := color.RGBA{stampBlueR, stampBlueG, stampBlueB, 180}

	fontSize := profile.FontSize
	if fontSize < stampMinFontPt {
		fontSize = stampMinFontPt
	}

	titleFace, err := loadFont(fontSize + 2)
	if err != nil {
		return err
	}

	textFace, err := loadFont(fontSize)
	if err != nil {
		return err
	}

	smallFace, err := loadFont(fontSize - 1)
	if err != nil {
		return err
	}

	drawBorder(img, 0, 0, w, h, lightBlue, 2)

	headerText := tr(msgGostHeader)
	drawText(img, int(float64(w)*0.02), int(float64(h)*0.15), headerText, titleFace, blue)

	drawLine(img, int(float64(w)*0.02), int(float64(h)*0.25), w-int(float64(w)*0.02), lightBlue)

	leftX := int(float64(w) * 0.02)
	rightX := w / 2
	colWidth := w/2 - int(float64(w)*0.04)

	lineHeight := int(float64(h) * 0.16)
	y := int(float64(h) * 0.32)

	leftLines := buildLeftColumn(d, profile)
	rightLines := buildRightColumn(d, profile)

	for i, line := range leftLines {
		py := y + i*lineHeight
		if py > h-int(float64(h)*0.05) {
			break
		}
		drawWrapped(img, leftX, py, colWidth, line, textFace, blue)
	}

	for i, line := range rightLines {
		py := y + i*lineHeight
		if py > h-int(float64(h)*0.05) {
			break
		}
		drawWrapped(img, rightX, py, colWidth, line, smallFace, blue)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func buildLeftColumn(d StampData, profile *StampProfile) []string {
	var lines []string

	lines = append(lines, tr(msgOwner)+": "+safeText(d.Owner))

	if profile.IncludeIssuer && d.Issuer != "" {
		lines = append(lines, tr(msgIssuer)+": "+safeText(d.Issuer))
	}

	if profile.IncludeDate {
		lines = append(lines, tr(msgDate)+": "+safeText(d.SignedAt))
	}

	if profile.IncludeReason && d.Reason != "" {
		lines = append(lines, tr(msgReason)+": "+safeText(d.Reason))
	}

	return lines
}

func buildRightColumn(d StampData, profile *StampProfile) []string {
	var lines []string

	lines = append(lines, tr(msgSerialNumber)+": "+safeText(d.Serial))

	if d.ValidFrom != "" && d.ValidTo != "" {
		lines = append(lines, tr(msgGostValidity)+": "+d.ValidFrom+" - "+d.ValidTo)
	}

	if profile.IncludeIssuer && d.Thumbprint != "" {
		lines = append(lines, "SHA1: "+truncateHash(d.Thumbprint))
	}

	if d.SignatureFN != "" {
		lines = append(lines, tr(msgSignatureShort)+": "+safeText(d.SignatureFN))
	}

	return lines
}

func truncateHash(hash string) string {
	h := strings.ReplaceAll(hash, " ", "")
	h = strings.ToUpper(h)
	if len(h) <= 24 {
		return h
	}
	return h[:16] + "..." + h[len(h)-8:]
}

func loadFont(size float64) (font.Face, error) {
	ft, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func drawText(img *image.RGBA, x, y int, s string, face font.Face, c color.Color) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(s)
}

func drawWrapped(img *image.RGBA, x, y, maxWidth int, s string, face font.Face, c color.Color) {
	lines := wrapText(s, face, maxWidth)
	for i, line := range lines {
		drawText(img, x, y+i*20, line, face, c)
	}
}

func wrapText(s string, face font.Face, maxWidth int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	cur := words[0]

	for _, w := range words[1:] {
		test := cur + " " + w
		if textWidth(face, test) <= maxWidth {
			cur = test
		} else {
			lines = append(lines, cur)
			cur = w
		}
	}
	lines = append(lines, cur)

	if len(lines) > 3 {
		lines = []string{lines[0], lines[1], ellipsize(lines[2:], face, maxWidth)}
	}

	return lines
}

func ellipsize(parts []string, face font.Face, maxWidth int) string {
	s := strings.Join(parts, " ")
	r := []rune(s)

	for len(r) > 0 {
		r = r[:len(r)-1]
		t := string(r) + "..."
		if textWidth(face, t) <= maxWidth {
			return t
		}
	}
	return "..."
}

func textWidth(face font.Face, s string) int {
	d := &font.Drawer{Face: face}
	return int(d.MeasureString(s) >> 6)
}

func drawLine(img *image.RGBA, x1, y, x2 int, c color.Color) {
	for x := x1; x <= x2; x++ {
		img.Set(x, y, c)
	}
}

func drawBorder(img *image.RGBA, x, y, w, h int, c color.Color, thickness int) {
	for i := 0; i < thickness; i++ {
		for px := x; px < x+w; px++ {
			img.Set(px, y+i, c)
			img.Set(px, y+h-1-i, c)
		}
		for py := y; py < y+h; py++ {
			img.Set(x+i, py, c)
			img.Set(x+w-1-i, py, c)
		}
	}
}

func safeText(s string) string {
	if s == "" {
		return "-"
	}
	return strings.TrimSpace(s)
}

func ValidateStampSize(widthMm, heightMm float64, fontSize float64) []string {
	var errors []string

	if widthMm < stampMinWidthMm {
		errors = append(errors, fmt.Sprintf(tr(msgStampMinSize), stampMinWidthMm, stampMinHeightMm))
	}
	if heightMm < stampMinHeightMm {
		errors = append(errors, fmt.Sprintf(tr(msgStampMinSize), stampMinWidthMm, stampMinHeightMm))
	}
	if fontSize < stampMinFontPt {
		errors = append(errors, fmt.Sprintf(tr(msgStampMinFont), stampMinFontPt))
	}

	return errors
}

func FormatStampDate(t time.Time) string {
	return t.Format("02.01.2006")
}

func FormatStampDateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04:05")
}
