package main

import (
	"testing"
)

func TestApplyPDFStampMissingInput(t *testing.T) {
	err := ApplyPDFStamp(PDFStampOptions{
		InputPDF:   "",
		OutputPDF:  "/tmp/out.pdf",
		StampImage: "/tmp/stamp.png",
		Pages:      "1-",
		Scale:      "0.96",
	})
	if err == nil {
		t.Error("expected error for empty input PDF")
	}
}

func TestApplyPDFStampMissingOutput(t *testing.T) {
	err := ApplyPDFStamp(PDFStampOptions{
		InputPDF:   "/tmp/in.pdf",
		OutputPDF:  "",
		StampImage: "/tmp/stamp.png",
		Pages:      "1-",
		Scale:      "0.96",
	})
	if err == nil {
		t.Error("expected error for empty output PDF")
	}
}

func TestApplyPDFStampMissingStamp(t *testing.T) {
	err := ApplyPDFStamp(PDFStampOptions{
		InputPDF:   "/tmp/in.pdf",
		OutputPDF:  "/tmp/out.pdf",
		StampImage: "",
		Pages:      "1-",
		Scale:      "0.96",
	})
	if err == nil {
		t.Error("expected error for empty stamp image")
	}
}

func TestApplyPDFStampWhitespaceInput(t *testing.T) {
	err := ApplyPDFStamp(PDFStampOptions{
		InputPDF:   "   ",
		OutputPDF:  "/tmp/out.pdf",
		StampImage: "/tmp/stamp.png",
		Pages:      "1-",
		Scale:      "0.96",
	})
	if err == nil {
		t.Error("expected error for whitespace input PDF")
	}
}

func TestPDFStampOptions(t *testing.T) {
	opts := PDFStampOptions{
		InputPDF:   "/tmp/test.pdf",
		OutputPDF:  "/tmp/test_stamped.pdf",
		StampImage: "/tmp/stamp.png",
		Pages:      "1-",
		Scale:      "0.96",
	}
	if opts.InputPDF == "" {
		t.Error("InputPDF should not be empty")
	}
}
