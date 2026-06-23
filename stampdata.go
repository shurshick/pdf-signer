package main

import (
	"path/filepath"
	"time"
)

type StampData struct {
	Owner       string
	Issuer      string
	Serial      string
	Thumbprint  string
	Reason      string
	SignedAt    string
	SignatureFN string
	ValidFrom   string
	ValidTo     string
}

func NewStampData(cert CertInfo, sigFile string, reason string) StampData {
	now := time.Now()
	return StampData{
		Owner:       cert.SubjectCN,
		Issuer:      cert.IssuerCN,
		Serial:      cert.Serial,
		Thumbprint:  cert.Thumbprint,
		Reason:      reason,
		SignedAt:    FormatStampDate(now),
		SignatureFN: filepath.Base(sigFile),
		ValidFrom:   cert.NotBefore.Format("02.01.2006"),
		ValidTo:     cert.NotAfter.Format("02.01.2006"),
	}
}
