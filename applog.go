package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logDir string

func init() {
	if dir := os.Getenv("PDFSIGNER_LOG_DIR"); dir != "" {
		logDir = dir
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			logDir = "."
			return
		}
		logDir = filepath.Join(home, ".local", "share", "pdfsigner", "logs")
	}
}

func LogInfo(message string) {
	writeLog("INFO", message, nil)
}

func LogError(message string, err error) {
	writeLog("ERROR", message, err)
}

func writeLog(level, message string, err error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return
	}

	logPath := filepath.Join(logDir, "app.log")
	f, openErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return
	}
	defer f.Close()

	line := fmt.Sprintf("[%s] [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), level, SanitizeLog(message))
	if err != nil {
		line += fmt.Sprintf("  error: %s\n", SanitizeLog(err.Error()))
	}

	f.WriteString(line)
}

func SanitizeLog(value string) string {
	if value == "" {
		return value
	}
	result := value

	lower := strings.ToLower(result)
	for _, keyword := range []string{"pin", "password", "пароль", "код", "secret", "key"} {
		idx := strings.Index(lower, keyword)
		if idx >= 0 {
			eqIdx := strings.IndexAny(result[idx:], ":=")
			if eqIdx > 0 {
				start := idx + eqIdx + 1
				if start < len(result) {
					end := strings.IndexAny(result[start:], " \t\n,;")
					if end < 0 {
						end = len(result) - start
					}
					result = result[:start] + "<redacted>" + result[start+end:]
					lower = strings.ToLower(result)
				}
			}
		}
	}

	if strings.Contains(result, "-----BEGIN") {
		start := strings.Index(result, "-----BEGIN")
		end := strings.Index(result[start:], "-----")
		if end > 0 {
			endBlock := strings.Index(result[start+end+5:], "-----END")
			if endBlock > 0 {
				result = result[:start] + "<redacted-certificate>" + result[start+end+5+endBlock+8:]
			}
		}
	}

	return result
}

func LogDir() string {
	return logDir
}
