package main

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func testFace(t *testing.T) font.Face {
	t.Helper()
	ft, err := opentype.Parse(goregular.TTF)
	if err != nil {
		t.Skip("cannot load Go font for test:", err)
	}
	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		t.Fatal(err)
	}
	return face
}

func TestTextWidth(t *testing.T) {
	face := testFace(t)
	w1 := textWidth(face, "hello")
	w2 := textWidth(face, "hello world")
	if w1 <= 0 || w2 <= 0 {
		t.Errorf("textWidth returned non-positive: %d, %d", w1, w2)
	}
	if w2 <= w1 {
		t.Errorf("wider string should have larger width: %d <= %d", w2, w1)
	}
}

func TestWrapText(t *testing.T) {
	face := testFace(t)
	lines := wrapText("hello world test string", face, 200)
	if len(lines) == 0 {
		t.Fatal("wrapText returned empty")
	}
}

func TestWrapTextEmpty(t *testing.T) {
	face := testFace(t)
	lines := wrapText("", face, 200)
	if len(lines) == 0 {
		t.Fatal("wrapText for empty string returned empty")
	}
}

func TestEllipsize(t *testing.T) {
	face := testFace(t)
	result := ellipsize([]string{"very long text that should be ellipsized"}, face, 100)
	if result == "" {
		t.Fatal("ellipsize returned empty")
	}
}

func TestSafeText(t *testing.T) {
	if got := safeText(""); got != "-" {
		t.Errorf("safeText(\"\") = %q, want %q", got, "-")
	}
	if got := safeText("  hello  "); got != "hello" {
		t.Errorf("safeText(\"  hello  \") = %q, want %q", got, "hello")
	}
	if got := safeText("test"); got != "test" {
		t.Errorf("safeText(\"test\") = %q, want %q", got, "test")
	}
}

func TestDrawText(t *testing.T) {
	face := testFace(t)
	img := image.NewRGBA(image.Rect(0, 0, 200, 50))
	drawText(img, 10, 20, "Hello", face, color.RGBA{0, 0, 0, 255})
}

func TestDrawLine(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 50))
	drawLine(img, 0, 0, 100, color.RGBA{0, 0, 0, 255})
}

func TestDrawBorder(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 50))
	drawBorder(img, 0, 0, 200, 50, color.RGBA{0, 0, 255, 255}, 2)
}

func TestDrawWrapped(t *testing.T) {
	face := testFace(t)
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))
	drawWrapped(img, 10, 20, 180, "Test wrapped text", face, color.RGBA{0, 0, 0, 255})
}

func TestCreateStampImage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stamp.png")

	data := StampData{
		Owner:       "Test User",
		Issuer:      "Test CA",
		Serial:      "12345",
		Thumbprint:  "AABBCCDD",
		Reason:      "Test signing",
		SignedAt:    "01.01.2026",
		SignatureFN: "test.sig",
		ValidFrom:   "01.01.2025",
		ValidTo:     "01.01.2027",
	}

	err := CreateStampImage(path, data)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("stamp image is empty")
	}
}

func TestCreateStampImageWithProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stamp_profile.png")

	profile := DefaultStampProfile()
	data := StampData{
		Owner:       "Test User",
		Issuer:      "Test CA",
		Serial:      "12345",
		Thumbprint:  "AABBCCDD",
		Reason:      "Test signing",
		SignedAt:    "01.01.2026",
		SignatureFN: "test.sig",
		ValidFrom:   "01.01.2025",
		ValidTo:     "01.01.2027",
	}

	err := CreateStampImageWithProfile(path, data, profile)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("stamp image is empty")
	}
}

func TestTruncateHash(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"AABB", "AABB"},
		{"AABBCCDDEE11223344556677889900", "AABBCCDDEE112233...77889900"},
		{"AA BB CC", "AABBCC"},
	}
	for _, tt := range tests {
		got := truncateHash(tt.input)
		if got != tt.want {
			t.Errorf("truncateHash(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestValidateStampSize(t *testing.T) {
	errs := ValidateStampSize(100, 50, 8)
	if len(errs) > 0 {
		t.Errorf("valid stamp should have no errors, got %d", len(errs))
	}

	errs = ValidateStampSize(40, 15, 4)
	if len(errs) == 0 {
		t.Error("undersized stamp should have errors")
	}
}

func TestFormatStampDate(t *testing.T) {
	d := time.Date(2026, 6, 15, 10, 30, 0, 0, time.UTC)
	got := FormatStampDate(d)
	if got != "15.06.2026" {
		t.Errorf("FormatStampDate = %q, want %q", got, "15.06.2026")
	}
}

func TestFormatStampDateTime(t *testing.T) {
	d := time.Date(2026, 6, 15, 10, 30, 45, 0, time.UTC)
	got := FormatStampDateTime(d)
	if got != "15.06.2026 10:30:45" {
		t.Errorf("FormatStampDateTime = %q, want %q", got, "15.06.2026 10:30:45")
	}
}
