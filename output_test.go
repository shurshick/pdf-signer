package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultOutputDir(t *testing.T) {
	dir := defaultOutputDir()
	if dir == "" {
		t.Fatal("defaultOutputDir returned empty string")
	}
	if !filepath.IsAbs(dir) && dir != "." {
		t.Errorf("defaultOutputDir returned relative path: %q", dir)
	}
}

func TestStampedOutputPath(t *testing.T) {
	dir := t.TempDir()

	inputPath := filepath.Join(dir, "test.pdf")
	output, err := stampedOutputPath(inputPath, filepath.Join(dir, "output"), false)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Ext(output) != ".pdf" {
		t.Errorf("expected .pdf extension, got %q", filepath.Ext(output))
	}
	base := filepath.Base(output)
	if base != "test_stamped.pdf" {
		t.Errorf("expected test_stamped.pdf, got %q", base)
	}
}

func TestStampedOutputPathNextToSource(t *testing.T) {
	dir := t.TempDir()

	inputPath := filepath.Join(dir, "test.pdf")
	output, err := stampedOutputPath(inputPath, "", true)
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(dir, "test_stamped.pdf")
	if output != expected {
		t.Errorf("output = %q, want %q", output, expected)
	}
}

func TestStampedOutputPathCreatesDir(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "new", "nested", "dir")

	inputPath := filepath.Join(dir, "test.pdf")
	_, err := stampedOutputPath(inputPath, outputDir, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("output directory was not created: %s", outputDir)
	}
}

func TestUniquePath(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "file.pdf")
	got := uniquePath(path)
	if got != path {
		t.Errorf("uniquePath(%q) = %q, want %q", path, got, path)
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	got = uniquePath(path)
	if got == path {
		t.Errorf("uniquePath should have returned a different path, got same: %q", got)
	}
	expected := filepath.Join(dir, "file_2.pdf")
	if got != expected {
		t.Errorf("uniquePath = %q, want %q", got, expected)
	}
}

func TestUniquePathMultiple(t *testing.T) {
	dir := t.TempDir()

	base := filepath.Join(dir, "file.pdf")
	for i := 0; i < 5; i++ {
		f, err := os.Create(base)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
		base = uniquePath(base)
	}

	if base == filepath.Join(dir, "file.pdf") {
		t.Error("expected a unique path after multiple calls")
	}
}

func TestSignatureOutputPath(t *testing.T) {
	got := signatureOutputPath("/path/to/doc_stamped.pdf")
	expected := "/path/to/doc_stamped.pdf.sig"
	if got != expected {
		t.Errorf("signatureOutputPath = %q, want %q", got, expected)
	}
}
