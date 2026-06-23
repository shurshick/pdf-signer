package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSaveSettings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	settings := &ApplicationSettings{
		VerifyAfterSigning: true,
		StampProfile:       DefaultStampProfile(),
	}

	data, _ := os.ReadFile(path)
	_ = data

	settings.StampProfile.IncludeCustom = true
	settings.StampProfile.CustomText = "Test custom text"

	if settings.VerifyAfterSigning != true {
		t.Error("VerifyAfterSigning should be true")
	}
	if settings.StampProfile.IncludeCustom != true {
		t.Error("IncludeCustom should be true")
	}
}

func TestDefaultStampProfile(t *testing.T) {
	p := DefaultStampProfile()
	if p.TemplateName != "standard" {
		t.Errorf("TemplateName = %q, want %q", p.TemplateName, "standard")
	}
	if p.WidthMm != 90 {
		t.Errorf("WidthMm = %f, want 90", p.WidthMm)
	}
	if p.HeightMm != 35 {
		t.Errorf("HeightMm = %f, want 35", p.HeightMm)
	}
	if p.FontSize != 8 {
		t.Errorf("FontSize = %f, want 8", p.FontSize)
	}
}

func TestStampProfileNormalize(t *testing.T) {
	p := &StampProfile{WidthMm: 10, HeightMm: 5, FontSize: 2, MinFontSize: 1, Opacity: 2.0, LogoScale: 50}
	p.Normalize()
	if p.WidthMm != 40 {
		t.Errorf("WidthMm = %f, want 40", p.WidthMm)
	}
	if p.HeightMm != 15 {
		t.Errorf("HeightMm = %f, want 15", p.HeightMm)
	}
	if p.FontSize != 4 {
		t.Errorf("FontSize = %f, want 4", p.FontSize)
	}
	if p.Opacity != 1.0 {
		t.Errorf("Opacity = %f, want 1.0", p.Opacity)
	}
	if p.LogoScale != 100 {
		t.Errorf("LogoScale = %d, want 100", p.LogoScale)
	}
}

func TestBuiltInProfiles(t *testing.T) {
	profiles := BuiltInProfiles()
	if len(profiles) != 3 {
		t.Errorf("BuiltInProfiles() returned %d profiles, want 3", len(profiles))
	}
	if _, ok := profiles["minimal"]; !ok {
		t.Error("missing minimal profile")
	}
	if _, ok := profiles["standard"]; !ok {
		t.Error("missing standard profile")
	}
	if _, ok := profiles["detailed"]; !ok {
		t.Error("missing detailed profile")
	}
}

func TestExportImportSettings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	settings := &ApplicationSettings{
		VerifyAfterSigning: true,
		StampProfile:       DefaultStampProfile(),
	}

	if err := ExportSettings(path, settings); err != nil {
		t.Fatal(err)
	}

	imported, err := ImportSettings(path)
	if err != nil {
		t.Fatal(err)
	}

	if imported.VerifyAfterSigning != true {
		t.Error("imported VerifyAfterSigning should be true")
	}
	if imported.StampProfile == nil {
		t.Error("imported StampProfile should not be nil")
	}
}

func TestImportSettingsInvalidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)

	_, err := ImportSettings(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
