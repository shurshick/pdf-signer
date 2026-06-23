package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type VerificationStatus int

const (
	VerifyValid   VerificationStatus = iota
	VerifyWarning
	VerifyInvalid
)

func (s VerificationStatus) String() string {
	switch s {
	case VerifyValid:
		return "VALID"
	case VerifyWarning:
		return "WARNING"
	case VerifyInvalid:
		return "INVALID"
	default:
		return "UNKNOWN"
	}
}

type VerificationReport struct {
	Status  VerificationStatus
	Details []string
	Errors  []string
	Warnings []string
}

type SignatureVerifier struct{}

func (v *SignatureVerifier) VerifyDetached(sigPath string) *VerificationReport {
	report := &VerificationReport{
		Status:   VerifyValid,
		Details:  make([]string, 0),
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	if !strings.HasSuffix(sigPath, ".sig") {
		report.Status = VerifyInvalid
		report.Errors = append(report.Errors, "Not a .sig file")
		return report
	}

	if _, err := os.Stat(sigPath); os.IsNotExist(err) {
		report.Status = VerifyInvalid
		report.Errors = append(report.Errors, fmt.Sprintf("File not found: %s", sigPath))
		return report
	}

	info, err := os.Stat(sigPath)
	if err != nil {
		report.Status = VerifyInvalid
		report.Errors = append(report.Errors, fmt.Sprintf("Cannot read file: %v", err))
		return report
	}

	report.Details = append(report.Details, fmt.Sprintf("Signature file: %s", filepath.Base(sigPath)))
	report.Details = append(report.Details, fmt.Sprintf("Size: %d bytes", info.Size()))
	report.Details = append(report.Details, fmt.Sprintf("Modified: %s", info.ModTime().Format("2006-01-02 15:04:05")))

	pdfPath := findMatchingPDF(sigPath)
	if pdfPath == "" {
		report.Status = VerifyWarning
		report.Warnings = append(report.Warnings, "Matching PDF file not found alongside .sig")
		report.Details = append(report.Details, "Cannot verify against original PDF")
	} else {
		report.Details = append(report.Details, fmt.Sprintf("Original PDF: %s", filepath.Base(pdfPath)))
	}

	report.Details = append(report.Details, fmt.Sprintf("Verification time: %s", time.Now().Format("2006-01-02 15:04:05")))

	return report
}

func (v *SignatureVerifier) VerifyEmbedded(pdfPath string) *VerificationReport {
	report := &VerificationReport{
		Status:   VerifyValid,
		Details:  make([]string, 0),
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		report.Status = VerifyInvalid
		report.Errors = append(report.Errors, fmt.Sprintf("File not found: %s", pdfPath))
		return report
	}

	info, err := os.Stat(pdfPath)
	if err != nil {
		report.Status = VerifyInvalid
		report.Errors = append(report.Errors, fmt.Sprintf("Cannot read file: %v", err))
		return report
	}

	report.Details = append(report.Details, fmt.Sprintf("PDF file: %s", filepath.Base(pdfPath)))
	report.Details = append(report.Details, fmt.Sprintf("Size: %d bytes", info.Size()))

	cmd := exec.Command("/opt/cprocsp/bin/amd64/csptest", "-sfsign", "-verify", "-in", pdfPath)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		report.Status = VerifyInvalid
		report.Errors = append(report.Errors, fmt.Sprintf("Verification failed: %v", err))
		if output != "" {
			report.Errors = append(report.Errors, output)
		}
	} else {
		report.Details = append(report.Details, "Embedded signature verified successfully")
	}

	report.Details = append(report.Details, fmt.Sprintf("Verification time: %s", time.Now().Format("2006-01-02 15:04:05")))

	return report
}

func (v *SignatureVerifier) FormatReport(reports []*VerificationReport) string {
	var sb strings.Builder
	sb.WriteString("Signature Verification Report\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	for i, r := range reports {
		sb.WriteString(fmt.Sprintf("--- File %d [%s] ---\n", i+1, r.Status))
		for _, d := range r.Details {
			sb.WriteString(fmt.Sprintf("  %s\n", d))
		}
		for _, w := range r.Warnings {
			sb.WriteString(fmt.Sprintf("  WARNING: %s\n", w))
		}
		for _, e := range r.Errors {
			sb.WriteString(fmt.Sprintf("  ERROR: %s\n", e))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func findMatchingPDF(sigPath string) string {
	base := strings.TrimSuffix(sigPath, ".sig")
	candidates := []string{
		base,
		strings.TrimSuffix(base, "-signed"),
		strings.TrimSuffix(base, "_signed"),
		strings.TrimSuffix(base, "-stamped"),
		strings.TrimSuffix(base, "_stamped"),
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}
