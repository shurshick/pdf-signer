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
	return s.SignFileTo(pdfPath, cert, pdfPath+".sig")
}

func (s NativeSigner) SignFileTo(pdfPath string, cert CertInfo, sigPath string) (SignResult, error) {
	if strings.TrimSpace(pdfPath) == "" {
		return SignResult{}, fmt.Errorf("%s", tr(msgNoPDF))
	}
	if strings.TrimSpace(sigPath) == "" {
		return SignResult{}, fmt.Errorf("%s", tr(msgOutputSignatureMissing))
	}

	absPDF, err := filepath.Abs(pdfPath)
	if err != nil {
		return SignResult{}, fmt.Errorf("%s: %w", tr(msgAbsPDFError), err)
	}

	absSig, err := filepath.Abs(sigPath)
	if err != nil {
		return SignResult{}, fmt.Errorf("%s: %w", tr(msgAbsSignatureError), err)
	}

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
		"-out", absSig,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return SignResult{}, fmt.Errorf("%s: %v\n%s", tr(msgSignError), err, string(out))
	}

	if _, err := os.Stat(absSig); err != nil {
		return SignResult{}, fmt.Errorf("%s: %s\n%s", tr(msgSignatureMissing), absSig, string(out))
	}

	return SignResult{
		SignaturePath: absSig,
	}, nil
}

func (s NativeSigner) SignFileEmbedded(pdfPath string, cert CertInfo) (SignResult, error) {
	return SignResult{}, fmt.Errorf("%s", tr(msgEmbeddedSignLinuxUnsupported))
}
