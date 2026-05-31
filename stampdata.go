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
    SignedAt    string
    SignatureFN string
}

func NewStampData(cert CertInfo, sigFile string) StampData {
    return StampData{
	Owner:       cert.SubjectCN,
	Issuer:      cert.IssuerCN,
	Serial:      cert.Serial,
	Thumbprint:  cert.Thumbprint,
	SignedAt:    time.Now().Format("02.01.2006 15:04:05"),
	SignatureFN: filepath.Base(sigFile), // 👈 только имя файла
    }
}