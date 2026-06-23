package main

import (
	"testing"
)

func TestDiagnosticsReport(t *testing.T) {
	d := &CryptoProDiagnostics{}
	report := d.Run()
	if report == nil {
		t.Fatal("Run() returned nil")
	}
	if report.AppVersion != appVersion {
		t.Errorf("AppVersion = %q, want %q", report.AppVersion, appVersion)
	}
	if len(report.Items) == 0 {
		t.Error("expected at least one diagnostic item")
	}
}

func TestDiagnosticsFormatReport(t *testing.T) {
	d := &CryptoProDiagnostics{}
	report := d.Run()
	formatted := d.FormatReport(report)
	if formatted == "" {
		t.Error("FormatReport returned empty string")
	}
	if !containsSubstr(formatted, "PDF Signer Linux Diagnostics") {
		t.Error("report missing header")
	}
}

func TestDiagnosticSeverity(t *testing.T) {
	if DiagOK.String() != "OK" {
		t.Errorf("DiagOK.String() = %q", DiagOK.String())
	}
	if DiagWarning.String() != "WARNING" {
		t.Errorf("DiagWarning.String() = %q", DiagWarning.String())
	}
	if DiagError.String() != "ERROR" {
		t.Errorf("DiagError.String() = %q", DiagError.String())
	}
}
