package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSanitizeLog(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		{"password=secret123", "<redacted>"},
		{"PIN:=12345", "<redacted>"},
		{"normal message", "normal message"},
		{"-----BEGIN CERTIFICATE-----\ndata\n-----END CERTIFICATE-----", "<redacted-certificate>"},
		{"", ""},
	}
	for _, tt := range tests {
		result := SanitizeLog(tt.input)
		if tt.contains == "" {
			if result != tt.input {
				t.Errorf("SanitizeLog(%q) = %q, want %q", tt.input, result, tt.input)
			}
		} else {
			if !containsSubstr(result, tt.contains) {
				t.Errorf("SanitizeLog(%q) = %q, should contain %q", tt.input, result, tt.contains)
			}
		}
	}
}

func TestLogDir(t *testing.T) {
	dir := LogDir()
	if dir == "" {
		t.Error("LogDir() returned empty string")
	}
}

func TestLogInfo(t *testing.T) {
	LogInfo("test log message")
	dir := LogDir()
	logPath := filepath.Join(dir, "app.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("log file was not created")
	}
}
