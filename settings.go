package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ApplicationSettings struct {
	VerifyAfterSigning bool         `json:"verifyAfterSigning"`
	StampProfile       *StampProfile `json:"stampProfile,omitempty"`
}

type StampProfile struct {
	TemplateName   string  `json:"templateName"`
	Pages          string  `json:"pages"`
	PositionMode   string  `json:"positionMode"`
	CustomX        float64 `json:"customX"`
	CustomY        float64 `json:"customY"`
	WidthMm        float64 `json:"widthMm"`
	HeightMm       float64 `json:"heightMm"`
	FontSize       float64 `json:"fontSize"`
	MinFontSize    float64 `json:"minFontSize"`
	Opacity        float64 `json:"opacity"`
	Scale          float64 `json:"scale"`
	IncludeOwner   bool    `json:"includeOwner"`
	IncludeIssuer  bool    `json:"includeIssuer"`
	IncludeDate    bool    `json:"includeDate"`
	IncludeReason  bool    `json:"includeReason"`
	IncludeSerial  bool    `json:"includeSerial"`
	IncludeCustom  bool    `json:"includeCustom"`
	CustomText     string  `json:"customText"`
	AutoPlace      bool    `json:"autoPlace"`
	LogoPath       string  `json:"logoPath"`
	LogoScale      int     `json:"logoScale"`
}

func DefaultStampProfile() *StampProfile {
	return &StampProfile{
		TemplateName: "standard",
		Pages:        "1-",
		PositionMode: "BottomRight",
		CustomX:      36,
		CustomY:      36,
		WidthMm:      90,
		HeightMm:     35,
		FontSize:     8,
		MinFontSize:  6,
		Opacity:      1.0,
		Scale:        0.96,
		IncludeOwner: true,
		IncludeIssuer: true,
		IncludeDate:  true,
		IncludeReason: true,
		IncludeSerial: true,
		AutoPlace:    false,
		LogoScale:    100,
	}
}

func (p *StampProfile) Normalize() {
	if p.WidthMm < 40 {
		p.WidthMm = 40
	}
	if p.WidthMm > 200 {
		p.WidthMm = 200
	}
	if p.HeightMm < 15 {
		p.HeightMm = 15
	}
	if p.HeightMm > 100 {
		p.HeightMm = 100
	}
	if p.FontSize < 4 {
		p.FontSize = 4
	}
	if p.FontSize > 16 {
		p.FontSize = 16
	}
	if p.MinFontSize < 4 {
		p.MinFontSize = 4
	}
	if p.MinFontSize > p.FontSize {
		p.MinFontSize = p.FontSize
	}
	if p.Opacity <= 0 {
		p.Opacity = 1.0
	}
	if p.Opacity > 1 {
		p.Opacity = 1
	}
	if p.Scale <= 0 {
		p.Scale = 0.96
	}
	if p.LogoScale < 100 {
		p.LogoScale = 100
	}
	if p.LogoScale > 300 {
		p.LogoScale = 300
	}
}

func settingsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "share", "pdfsigner", "settings.json")
}

func LoadSettings() *ApplicationSettings {
	path := settingsPath()
	if path == "" {
		return &ApplicationSettings{StampProfile: DefaultStampProfile()}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return &ApplicationSettings{StampProfile: DefaultStampProfile()}
	}

	var settings ApplicationSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return &ApplicationSettings{StampProfile: DefaultStampProfile()}
	}

	if settings.StampProfile == nil {
		settings.StampProfile = DefaultStampProfile()
	}

	return &settings
}

func SaveSettings(settings *ApplicationSettings) error {
	path := settingsPath()
	if path == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func ExportSettings(path string, settings *ApplicationSettings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ImportSettings(path string) (*ApplicationSettings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var settings ApplicationSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	if settings.StampProfile == nil {
		settings.StampProfile = DefaultStampProfile()
	}

	settings.StampProfile.Normalize()
	return &settings, nil
}

func BuiltInProfiles() map[string]*StampProfile {
	return map[string]*StampProfile{
		"minimal": {
			TemplateName: "minimal", Pages: "1", PositionMode: "BottomRight",
			WidthMm: 70, HeightMm: 25, FontSize: 7, MinFontSize: 6, Opacity: 1.0, Scale: 0.96,
			IncludeOwner: true, IncludeIssuer: false, IncludeDate: true, IncludeReason: false,
			IncludeSerial: true, AutoPlace: false, LogoScale: 100,
		},
		"standard": DefaultStampProfile(),
		"detailed": {
			TemplateName: "detailed", Pages: "1-", PositionMode: "BottomRight",
			WidthMm: 120, HeightMm: 45, FontSize: 8, MinFontSize: 6, Opacity: 1.0, Scale: 0.96,
			IncludeOwner: true, IncludeIssuer: true, IncludeDate: true, IncludeReason: true,
			IncludeSerial: true, AutoPlace: false, LogoScale: 100,
		},
	}
}
