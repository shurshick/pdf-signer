package main

import (
	"testing"
	"time"
)

func TestNewStampData(t *testing.T) {
	cert := CertInfo{
		SubjectCN:  "Test User",
		IssuerCN:   "Test CA",
		Serial:     "12345",
		Thumbprint: "AABBCCDD",
	}
	sigFile := "/path/to/test_stamped.pdf.sig"
	reason := "Test reason"

	data := NewStampData(cert, sigFile, reason)

	if data.Owner != "Test User" {
		t.Errorf("Owner = %q, want %q", data.Owner, "Test User")
	}
	if data.Issuer != "Test CA" {
		t.Errorf("Issuer = %q, want %q", data.Issuer, "Test CA")
	}
	if data.Serial != "12345" {
		t.Errorf("Serial = %q, want %q", data.Serial, "12345")
	}
	if data.Thumbprint != "AABBCCDD" {
		t.Errorf("Thumbprint = %q, want %q", data.Thumbprint, "AABBCCDD")
	}
	if data.Reason != "Test reason" {
		t.Errorf("Reason = %q, want %q", data.Reason, "Test reason")
	}
	if data.SignatureFN != "test_stamped.pdf.sig" {
		t.Errorf("SignatureFN = %q, want %q", data.SignatureFN, "test_stamped.pdf.sig")
	}

	_, err := time.Parse("02.01.2006", data.SignedAt)
	if err != nil {
		t.Errorf("SignedAt is not a valid date: %q, err: %v", data.SignedAt, err)
	}
}

func TestNewStampDataEmptySigFile(t *testing.T) {
	cert := CertInfo{
		SubjectCN: "User",
	}
	data := NewStampData(cert, "", "reason")
	if data.SignatureFN != "." {
		t.Errorf("SignatureFN = %q, want %q", data.SignatureFN, ".")
	}
}
