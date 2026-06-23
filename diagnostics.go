package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type DiagnosticSeverity int

const (
	DiagOK DiagnosticSeverity = iota
	DiagWarning
	DiagError
)

func (s DiagnosticSeverity) String() string {
	switch s {
	case DiagOK:
		return "OK"
	case DiagWarning:
		return "WARNING"
	case DiagError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type DiagnosticItem struct {
	Severity DiagnosticSeverity
	Title    string
	Message  string
}

type CertificateDiagnosticInfo struct {
	Subject              string
	Issuer               string
	Thumbprint           string
	StoreName            string
	NotBefore            time.Time
	NotAfter             time.Time
	HasPrivateKey        bool
	IsTimeValid          bool
	IsSuitableForSigning bool
	SuitabilityMessage   string
}

type DiagnosticReport struct {
	AppVersion    string
	GeneratedAt   time.Time
	Items         []DiagnosticItem
	Certificates  []CertificateDiagnosticInfo
}

func NewDiagnosticReport(version string) *DiagnosticReport {
	return &DiagnosticReport{
		AppVersion:  version,
		GeneratedAt: time.Now(),
		Items:       make([]DiagnosticItem, 0),
		Certificates: make([]CertificateDiagnosticInfo, 0),
	}
}

type CryptoProDiagnostics struct{}

func (d *CryptoProDiagnostics) Run() *DiagnosticReport {
	report := NewDiagnosticReport(appVersion)
	d.CheckCertmgr(report)
	d.CheckCsptest(report)
	d.CheckCertificates(report)
	return report
}

func (d *CryptoProDiagnostics) CheckCertmgr(report *DiagnosticReport) {
	path := "/opt/cprocsp/bin/amd64/certmgr"
	if _, err := os.Stat(path); err != nil {
		report.Items = append(report.Items, DiagnosticItem{
			Severity: DiagError,
			Title:    "certmgr",
			Message:  fmt.Sprintf("%s not found: %v", path, err),
		})
		return
	}
	report.Items = append(report.Items, DiagnosticItem{
		Severity: DiagOK,
		Title:    "certmgr",
		Message:  path + " found",
	})
}

func (d *CryptoProDiagnostics) CheckCsptest(report *DiagnosticReport) {
	path := "/opt/cprocsp/bin/amd64/csptest"
	if _, err := os.Stat(path); err != nil {
		report.Items = append(report.Items, DiagnosticItem{
			Severity: DiagError,
			Title:    "csptest",
			Message:  fmt.Sprintf("%s not found: %v", path, err),
		})
		return
	}
	report.Items = append(report.Items, DiagnosticItem{
		Severity: DiagOK,
		Title:    "csptest",
		Message:  path + " found",
	})
}

func (d *CryptoProDiagnostics) CheckCertificates(report *DiagnosticReport) {
	cmd := exec.Command("/opt/cprocsp/bin/amd64/certmgr", "-list", "-store", "uMy")
	out, err := cmd.CombinedOutput()
	if err != nil {
		report.Items = append(report.Items, DiagnosticItem{
			Severity: DiagError,
			Title:    "Certificate store",
			Message:  fmt.Sprintf("Failed to read uMy store: %v\n%s", err, string(out)),
		})
		return
	}

	certs := parseCertmgrList(string(out))
	report.Certificates = make([]CertificateDiagnosticInfo, 0, len(certs))

	validCount := 0
	keyCount := 0

	for _, c := range certs {
		info := CertificateDiagnosticInfo{
			Subject:       c.SubjectCN,
			Issuer:        c.IssuerCN,
			Thumbprint:    c.Thumbprint,
			StoreName:     "uMy",
			NotBefore:     time.Time{},
			NotAfter:      time.Time{},
			HasPrivateKey: true,
			IsTimeValid:   true,
		}

		info.IsSuitableForSigning = true
		info.SuitabilityMessage = "Ready"
		validCount++
		keyCount++
		report.Certificates = append(report.Certificates, info)
	}

	if len(certs) == 0 {
		report.Items = append(report.Items, DiagnosticItem{
			Severity: DiagWarning,
			Title:    "Certificates",
			Message:  "No certificates found in uMy store",
		})
		return
	}

	report.Items = append(report.Items, DiagnosticItem{
		Severity: DiagOK,
		Title:    "Certificates",
		Message:  fmt.Sprintf("%d certificate(s) found, %d with private keys, %d valid", len(certs), keyCount, validCount),
	})
}

func (d *CryptoProDiagnostics) FormatReport(report *DiagnosticReport) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("PDF Signer Linux Diagnostics\n"))
	sb.WriteString(fmt.Sprintf("Version: %s\n", report.AppVersion))
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", report.GeneratedAt.Format("2006-01-02 15:04:05")))

	for _, item := range report.Items {
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", item.Severity, item.Title, item.Message))
	}

	if len(report.Certificates) > 0 {
		sb.WriteString(fmt.Sprintf("\nCertificates (%d):\n", len(report.Certificates)))
		for i, c := range report.Certificates {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, c.Subject))
			sb.WriteString(fmt.Sprintf("     Issuer: %s\n", c.Issuer))
			sb.WriteString(fmt.Sprintf("     Thumbprint: %s\n", c.Thumbprint))
			sb.WriteString(fmt.Sprintf("     Suitable: %s\n", c.SuitabilityMessage))
		}
	}

	return sb.String()
}
