package main

import (
	"testing"
	"time"
)

func TestBuildStampTextFromProfile(t *testing.T) {
	profile := DefaultStampProfile()
	cert := CertInfo{SubjectCN: "Test User", IssuerCN: "Test CA", Serial: "12345", Thumbprint: "AABBCCDD", NotBefore: time.Now(), NotAfter: time.Now().AddDate(1, 0, 0)}
	text := BuildStampTextFromProfile(profile, cert, "Test reason")

	if text == "" {
		t.Error("BuildStampTextFromProfile returned empty string")
	}
	if !containsSubstr(text, "Test User") {
		t.Error("stamp text missing owner name")
	}
	if !containsSubstr(text, "Test CA") {
		t.Error("stamp text missing issuer")
	}
	if !containsSubstr(text, tr(msgGostHeader)) {
		t.Error("stamp text missing GOST header")
	}
}

func TestBuildStampTextMinimalProfile(t *testing.T) {
	profile := &StampProfile{
		IncludeOwner:  false,
		IncludeIssuer: false,
		IncludeDate:   false,
		IncludeReason: false,
		IncludeSerial: false,
	}
	cert := CertInfo{SubjectCN: "User", IssuerCN: "CA", Serial: "123", Thumbprint: "XX"}
	text := BuildStampTextFromProfile(profile, cert, "reason")

	if containsSubstr(text, "User") {
		t.Error("minimal profile should not include owner")
	}
	if containsSubstr(text, "CA") {
		t.Error("minimal profile should not include issuer")
	}
}

func TestProfileToPosLabel(t *testing.T) {
	if profileToPosLabel("BottomRight") != tr(msgPosBottomRight) {
		t.Error("wrong label for BottomRight")
	}
	if profileToPosLabel("TopLeft") != tr(msgPosTopLeft) {
		t.Error("wrong label for TopLeft")
	}
}

func TestPosLabelToProfile(t *testing.T) {
	if posLabelToProfile(tr(msgPosBottomRight)) != "BottomRight" {
		t.Error("wrong mode for BottomRight label")
	}
	if posLabelToProfile(tr(msgPosTopLeft)) != "TopLeft" {
		t.Error("wrong mode for TopLeft label")
	}
}

func TestProfileToTemplateLabel(t *testing.T) {
	if profileToTemplateLabel("minimal") != tr(msgTemplateMinimal) {
		t.Error("wrong label for minimal")
	}
	if profileToTemplateLabel("standard") != tr(msgTemplateStandard) {
		t.Error("wrong label for standard")
	}
	if profileToTemplateLabel("detailed") != tr(msgTemplateDetailed) {
		t.Error("wrong label for detailed")
	}
}

func TestLabelToProfileKey(t *testing.T) {
	if labelToProfileKey(tr(msgTemplateMinimal)) != "minimal" {
		t.Error("wrong key for minimal label")
	}
	if labelToProfileKey(tr(msgTemplateStandard)) != "standard" {
		t.Error("wrong key for standard label")
	}
	if labelToProfileKey(tr(msgTemplateDetailed)) != "detailed" {
		t.Error("wrong key for detailed label")
	}
}

func TestValidateStampProfile(t *testing.T) {
	goodProfile := DefaultStampProfile()
	warnings := validateStampProfile(goodProfile)
	if len(warnings) > 0 {
		t.Errorf("default profile should have no warnings, got %d: %v", len(warnings), warnings)
	}

	badProfile := &StampProfile{
		WidthMm:     40,
		HeightMm:    15,
		FontSize:    3,
		MinFontSize: 4,
		Opacity:     0.3,
		Scale:       0.96,
	}
	warnings = validateStampProfile(badProfile)
	if len(warnings) == 0 {
		t.Error("bad profile should have warnings")
	}
}
