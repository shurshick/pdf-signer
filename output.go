package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func defaultOutputDir() string {
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return "."
	}
	return filepath.Join(home, "Documents", "Signed PDFs")
}

func stampedOutputPath(inputPath, outputDir string, saveNextToSource bool) (string, error) {
	dir := outputDir
	if saveNextToSource {
		dir = filepath.Dir(inputPath)
	}
	if strings.TrimSpace(dir) == "" {
		dir = "."
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), ext)
	return uniquePath(filepath.Join(dir, base+"_stamped"+ext)), nil
}

func signatureOutputPath(stampedPDFPath string) string {
	return uniquePath(stampedPDFPath + ".sig")
}

func uniquePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	for i := 2; ; i++ {
		candidate := base + "_" + strconv.Itoa(i) + ext
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}
