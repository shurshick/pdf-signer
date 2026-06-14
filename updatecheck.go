package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const latestReleaseAPI = "https://api.github.com/repos/shurshick/pdf-signer/releases/latest"

type UpdateInfo struct {
	CurrentVersion string
	LatestVersion  string
	ReleaseName    string
	ReleaseURL     string
	IsNewer        bool
}

func CheckForUpdates() (UpdateInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestReleaseAPI, nil)
	if err != nil {
		return UpdateInfo{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "pdfsigner/"+appVersion)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return UpdateInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return UpdateInfo{}, fmt.Errorf("GitHub release check failed: HTTP %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return UpdateInfo{}, err
	}

	latest := strings.TrimSpace(release.TagName)
	return UpdateInfo{
		CurrentVersion: appVersion,
		LatestVersion:  strings.TrimPrefix(strings.TrimPrefix(latest, "v"), "V"),
		ReleaseName:    strings.TrimSpace(release.Name),
		ReleaseURL:     strings.TrimSpace(release.HTMLURL),
		IsNewer:        compareVersions(latest, appVersion) > 0,
	}, nil
}

func compareVersions(left, right string) int {
	lparts := parseVersion(left)
	rparts := parseVersion(right)

	for i := 0; i < len(lparts) || i < len(rparts); i++ {
		lvalue := versionPart(lparts, i)
		rvalue := versionPart(rparts, i)
		if lvalue > rvalue {
			return 1
		}
		if lvalue < rvalue {
			return -1
		}
	}
	return 0
}

func parseVersion(value string) []int {
	normalized := strings.TrimSpace(value)
	normalized = strings.TrimPrefix(strings.TrimPrefix(normalized, "v"), "V")
	if idx := strings.IndexAny(normalized, "-+"); idx >= 0 {
		normalized = normalized[:idx]
	}

	rawParts := strings.Split(normalized, ".")
	parts := make([]int, 0, len(rawParts))
	for _, raw := range rawParts {
		part, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil {
			parts = append(parts, 0)
			continue
		}
		parts = append(parts, part)
	}
	return parts
}

func versionPart(parts []int, index int) int {
	if index < 0 || index >= len(parts) {
		return 0
	}
	return parts[index]
}
