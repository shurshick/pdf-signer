package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func CreateStampImage(path string, d StampData) error {
	const w = 1800
	const h = 185

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Полностью прозрачный фон
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	// Синий цвет штампа
	blue := color.RGBA{0, 70, 160, 255}
	lightBlue := color.RGBA{0, 70, 160, 180}

	titleFace, err := loadFont(20)
	if err != nil {
		return err
	}

	textFace, err := loadFont(16)
	if err != nil {
		return err
	}

	smallFace, err := loadFont(14)
	if err != nil {
		return err
	}

	// Рамка
	drawBorder(img, 0, 0, w, h, lightBlue, 2)

	// Заголовок
	drawText(img, 20, 28, tr(msgStampTitle), titleFace, blue)

	// Линия под заголовком
	drawLine(img, 20, 36, w-20, lightBlue)

	// Две колонки
	leftX := 20
	rightX := 900

	y1 := 63
	y2 := 90
	y3 := 117
	y4 := 144

	left1 := tr(msgOwner) + ": " + safeText(d.Owner)
	left2 := tr(msgIssuer) + ": " + safeText(d.Issuer)
	left3 := tr(msgDate) + ": " + safeText(d.SignedAt)
	left4 := tr(msgReason) + ": " + safeText(d.Reason)

	right1 := tr(msgSerialNumber) + ": " + safeText(d.Serial)
	right2 := "SHA1: " + safeText(d.Thumbprint)
	right3 := tr(msgSignatureShort) + ": " + safeText(d.SignatureFN)

	drawWrapped(img, leftX, y1, 840, left1, textFace, blue)
	drawWrapped(img, leftX, y2, 840, left2, textFace, blue)
	drawWrapped(img, leftX, y3, 840, left3, smallFace, blue)
	drawWrapped(img, leftX, y4, 840, left4, smallFace, blue)

	drawWrapped(img, rightX, y1, 870, right1, textFace, blue)
	drawWrapped(img, rightX, y2, 870, right2, smallFace, blue)
	drawWrapped(img, rightX, y3, 870, right3, smallFace, blue)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
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

	if len(lines) > 2 {
		lines = []string{lines[0], ellipsize(lines[1:], face, maxWidth)}
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
