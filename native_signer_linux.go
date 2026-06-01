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
		return SignResult{}, fmt.Errorf("%s", tr(msgNoPDF))
	}

	absPDF, err := filepath.Abs(pdfPath)
	if err != nil {
		return SignResult{}, fmt.Errorf("%s: %w", tr(msgAbsPDFError), err)
	}

	sigPath := absPDF + ".sig"

	// Используем CN, как в ручной команде csptest -my
	subject := strings.TrimSpace(cert.SubjectCN)
	if subject == "" {
		return SignResult{}, fmt.Errorf("%s", tr(msgEmptyCertCN))
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
		return SignResult{}, fmt.Errorf("%s: %v\n%s", tr(msgSignError), err, string(out))
	}

	if _, err := os.Stat(sigPath); err != nil {
		return SignResult{}, fmt.Errorf("%s: %s\n%s", tr(msgSignatureMissing), sigPath, string(out))
	}

	return SignResult{
		SignaturePath: sigPath,
	}, nil
}
