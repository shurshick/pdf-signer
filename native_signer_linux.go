package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

type NativeSigner struct{}

func (s NativeSigner) SignFile(pdfPath string, cert CertInfo) (SignResult, error) {
    if strings.TrimSpace(pdfPath) == "" {
	return SignResult{}, fmt.Errorf("не указан PDF")
    }

    absPDF, err := filepath.Abs(pdfPath)
    if err != nil {
	return SignResult{}, fmt.Errorf("не удалось получить абсолютный путь к PDF: %w", err)
    }

    sigPath := absPDF + ".sig"

    // Используем CN, как в ручной команде csptest -my
    subject := strings.TrimSpace(cert.SubjectCN)
    if subject == "" {
	return SignResult{}, fmt.Errorf("у сертификата пустой CN")
    }

    cmd := exec.Command(
	"/opt/cprocsp/bin/amd64/csptest",
	"-sfsign",
	"-sign",
	"-detached",
	"-add",
	"-my", subject,
	"-in", absPDF,
	"-out", sigPath,
    )

    out, err := cmd.CombinedOutput()
    if err != nil {
	return SignResult{}, fmt.Errorf("ошибка подписи: %v\n%s", err, string(out))
    }

    if _, err := os.Stat(sigPath); err != nil {
	return SignResult{}, fmt.Errorf("файл подписи не создан: %s\n%s", sigPath, string(out))
    }

    return SignResult{
	SignaturePath: sigPath,
    }, nil
}