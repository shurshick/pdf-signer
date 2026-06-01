package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type CertInfo struct {
	SubjectCN  string
	IssuerCN   string
	Thumbprint string
	Serial     string
	Container  string
	Provider   string
	Label      string
}

func GetCertificates() ([]CertInfo, error) {
	cmd := exec.Command("/opt/cprocsp/bin/amd64/certmgr", "-list", "-store", "uMy")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %v\n%s", tr(msgCertmgrError), err, string(out))
	}

	certs := parseCertmgrList(string(out))
	if len(certs) == 0 {
		return nil, fmt.Errorf("%s", tr(msgNoCerts))
	}
	return certs, nil
}

func parseCertmgrList(s string) []CertInfo {
	var certs []CertInfo
	var cur *CertInfo

	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		if isCertSeparator(line) {
			if cur != nil {
				cur.Label = buildCertLabel(*cur)
				certs = append(certs, *cur)
			}
			cur = &CertInfo{}
			continue
		}

		if cur == nil {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Субъект"), strings.HasPrefix(line, "Subject"):
			cur.SubjectCN = extractCN(afterColon(line))
		case strings.HasPrefix(line, "Издатель"), strings.HasPrefix(line, "Issuer"):
			cur.IssuerCN = extractCN(afterColon(line))
		case strings.HasPrefix(line, "Серийный номер"), strings.HasPrefix(line, "Serial number"):
			cur.Serial = strings.TrimSpace(afterColon(line))
		case strings.HasPrefix(line, "SHA1 отпечаток"), strings.HasPrefix(line, "SHA1 hash"), strings.HasPrefix(line, "SHA1 thumbprint"):
			cur.Thumbprint = strings.TrimSpace(afterColon(line))
		case strings.HasPrefix(line, "Контейнер"), strings.HasPrefix(line, "Container"):
			cur.Container = strings.TrimSpace(afterColon(line))
		case strings.HasPrefix(line, "Имя провайдера"), strings.HasPrefix(line, "Provider name"):
			cur.Provider = strings.TrimSpace(afterColon(line))
		}
	}

	if cur != nil {
		cur.Label = buildCertLabel(*cur)
		certs = append(certs, *cur)
	}

	return certs
}

func isCertSeparator(s string) bool {
	if len(s) < 8 {
		return false
	}
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i == 0 || i >= len(s) {
		return false
	}
	for ; i < len(s); i++ {
		if s[i] != '-' {
			return false
		}
	}
	return true
}

func afterColon(s string) string {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func extractCN(s string) string {
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "CN=") {
			return strings.TrimPrefix(p, "CN=")
		}
	}
	return s
}

func buildCertLabel(c CertInfo) string {
	var parts []string
	if c.SubjectCN != "" {
		parts = append(parts, c.SubjectCN)
	}
	if c.IssuerCN != "" {
		parts = append(parts, tr(msgIssuer)+": "+c.IssuerCN)
	}
	if c.Serial != "" {
		parts = append(parts, "SN: "+c.Serial)
	}
	return strings.Join(parts, " | ")
}
