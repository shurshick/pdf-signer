package main

import (
	"testing"
)

func TestParseCertmgrListEmpty(t *testing.T) {
	certs := parseCertmgrList("")
	if len(certs) != 0 {
		t.Fatalf("expected 0 certs, got %d", len(certs))
	}
}

func TestParseCertmgrListSingleCert(t *testing.T) {
	input := `------------------------------------------------------
Субъект: CN=Test User, O=Test Org
Издатель: CN=Test CA, O=Test Org
Серийный номер: 1234567890
SHA1 отпечаток: AB CD EF 01 23 45 67 89
Контейнер: \\.\HDIMAGE\container
Имя провайдера: CryptoPro CSP
`
	certs := parseCertmgrList(input)
	if len(certs) != 1 {
		t.Fatalf("expected 1 cert, got %d", len(certs))
	}
	c := certs[0]
	if c.SubjectCN != "Test User" {
		t.Errorf("SubjectCN = %q, want %q", c.SubjectCN, "Test User")
	}
	if c.IssuerCN != "Test CA" {
		t.Errorf("IssuerCN = %q, want %q", c.IssuerCN, "Test CA")
	}
	if c.Serial != "1234567890" {
		t.Errorf("Serial = %q, want %q", c.Serial, "1234567890")
	}
	if c.Thumbprint != "AB CD EF 01 23 45 67 89" {
		t.Errorf("Thumbprint = %q, want %q", c.Thumbprint, "AB CD EF 01 23 45 67 89")
	}
}

func TestParseCertmgrListEnglish(t *testing.T) {
	input := `------------------------------------------------------
Subject: CN=English User, O=Org
Issuer: CN=CA, O=Org
Serial number: 999
SHA1 thumbprint: AA BB CC
Container: container
Provider name: Provider
`
	certs := parseCertmgrList(input)
	if len(certs) != 1 {
		t.Fatalf("expected 1 cert, got %d", len(certs))
	}
	c := certs[0]
	if c.SubjectCN != "English User" {
		t.Errorf("SubjectCN = %q, want %q", c.SubjectCN, "English User")
	}
	if c.IssuerCN != "CA" {
		t.Errorf("IssuerCN = %q, want %q", c.IssuerCN, "CA")
	}
}

func TestParseCertmgrListMultipleCerts(t *testing.T) {
	input := `------------------------------------------------------
Субъект: CN=First User
Издатель: CN=CA1
Серийный номер: 111
SHA1 отпечаток: AA AA AA
Контейнер: c1
Имя провайдера: p1
------------------------------------------------------
Субъект: CN=Second User
Издатель: CN=CA2
Серийный номер: 222
SHA1 отпечаток: BB BB BB
Контейнер: c2
Имя провайдера: p2
`
	certs := parseCertmgrList(input)
	if len(certs) != 2 {
		t.Fatalf("expected 2 certs, got %d", len(certs))
	}
	if certs[0].SubjectCN != "First User" {
		t.Errorf("first cert SubjectCN = %q, want %q", certs[0].SubjectCN, "First User")
	}
	if certs[1].SubjectCN != "Second User" {
		t.Errorf("second cert SubjectCN = %q, want %q", certs[1].SubjectCN, "Second User")
	}
}

func TestIsCertSeparator(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"------------------------------------------------------", true},
		{"12345678-9012-3456-7890-123456789012", false},
		{"---", false},
		{"", false},
		{"12345678", false},
		{"12345678--------", true},
	}
	for _, tt := range tests {
		got := isCertSeparator(tt.input)
		if got != tt.want {
			t.Errorf("isCertSeparator(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestAfterColon(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Subject: CN=Test", "CN=Test"},
		{"No colon here", ""},
		{"Key: value with: colons", "value with: colons"},
		{"Key:", ""},
	}
	for _, tt := range tests {
		got := afterColon(tt.input)
		if got != tt.want {
			t.Errorf("afterColon(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestExtractCN(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"CN=Test User, O=Org", "Test User"},
		{"O=Org, CN=Test User", "Test User"},
		{"No CN here", "No CN here"},
		{"CN=Only", "Only"},
	}
	for _, tt := range tests {
		got := extractCN(tt.input)
		if got != tt.want {
			t.Errorf("extractCN(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBuildCertLabel(t *testing.T) {
	c := CertInfo{
		SubjectCN:  "Test User",
		IssuerCN:   "Test CA",
		Serial:     "12345",
		Thumbprint: "AABB",
	}
	label := buildCertLabel(c)
	if label == "" {
		t.Fatal("buildCertLabel returned empty string")
	}
	if !containsSubstr(label, "Test User") {
		t.Errorf("label missing SubjectCN, got %q", label)
	}
}

func containsSubstr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsAt(s, sub))
}

func containsAt(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
