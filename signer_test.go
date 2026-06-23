package main

import (
	"testing"
)

func TestNativeSignerSignFileEmptyPath(t *testing.T) {
	signer := NativeSigner{}
	cert := CertInfo{SubjectCN: "Test"}
	_, err := signer.SignFile("", cert)
	if err == nil {
		t.Error("expected error for empty PDF path")
	}
}

func TestNativeSignerSignFileEmbeddedEmptyPath(t *testing.T) {
	signer := NativeSigner{}
	cert := CertInfo{SubjectCN: "Test"}
	_, err := signer.SignFileEmbedded("", cert)
	if err == nil {
		t.Error("expected error for empty PDF path")
	}
}

func TestNativeSignerSignFileEmptyCert(t *testing.T) {
	signer := NativeSigner{}
	cert := CertInfo{SubjectCN: ""}
	_, err := signer.SignFileTo("/tmp/test.pdf", cert, "/tmp/test.pdf.sig")
	if err == nil {
		t.Error("expected error for empty cert CN")
	}
}

func TestNativeSignerSignFileEmbeddedEmptyCert(t *testing.T) {
	signer := NativeSigner{}
	cert := CertInfo{SubjectCN: ""}
	_, err := signer.SignFileEmbedded("/tmp/test.pdf", cert)
	if err == nil {
		t.Error("expected error for empty cert CN")
	}
}

func TestSignModeConstants(t *testing.T) {
	if SignModeEmbedded != 0 {
		t.Errorf("SignModeEmbedded = %d, want 0", SignModeEmbedded)
	}
	if SignModeDetached != 1 {
		t.Errorf("SignModeDetached = %d, want 1", SignModeDetached)
	}
	if SignModeBoth != 2 {
		t.Errorf("SignModeBoth = %d, want 2", SignModeBoth)
	}
}
