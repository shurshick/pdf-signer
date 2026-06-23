package main

import (
	"testing"
)

func TestSignatureVerifierVerifyDetached(t *testing.T) {
	v := &SignatureVerifier{}

	report := v.VerifyDetached("/nonexistent/file.sig")
	if report.Status != VerifyInvalid {
		t.Errorf("expected VerifyInvalid for nonexistent file, got %v", report.Status)
	}

	report = v.VerifyDetached("/tmp/test.pdf")
	if report.Status != VerifyInvalid {
		t.Errorf("expected VerifyInvalid for non-.sig file, got %v", report.Status)
	}
}

func TestSignatureVerifierFormatReport(t *testing.T) {
	v := &SignatureVerifier{}
	reports := []*VerificationReport{
		{
			Status:  VerifyValid,
			Details: []string{"test detail"},
		},
	}
	formatted := v.FormatReport(reports)
	if formatted == "" {
		t.Error("FormatReport returned empty string")
	}
	if !containsSubstr(formatted, "VALID") {
		t.Error("report missing VALID status")
	}
}

func TestVerificationStatus(t *testing.T) {
	if VerifyValid.String() != "VALID" {
		t.Errorf("VerifyValid.String() = %q", VerifyValid.String())
	}
	if VerifyWarning.String() != "WARNING" {
		t.Errorf("VerifyWarning.String() = %q", VerifyWarning.String())
	}
	if VerifyInvalid.String() != "INVALID" {
		t.Errorf("VerifyInvalid.String() = %q", VerifyInvalid.String())
	}
}

func TestFindMatchingPDF(t *testing.T) {
	result := findMatchingPDF("/tmp/nonexistent.sig")
	if result != "" {
		t.Errorf("expected empty for nonexistent file, got %q", result)
	}
}
