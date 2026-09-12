package streamdeck_test

import (
	"testing"

	"github.com/FlowingSPDG/streamdeck/v2"
)

func TestParseRegistrationParams(t *testing.T) {
	p, err := streamdeck.ParseRegistrationParams([]string{
		"",
		"-port", "12345",
		"-pluginUUID", "213",
		"-registerEvent", "registerPlugin",
		"-info", `{"application":{"language":"en","platform":"mac","version":"4.1.0"},"plugin":{"version":"1.1"},"devicePixelRatio":2,"colors":{"buttonMouseOverBackgroundColor":"#464646FF","buttonPressedBackgroundColor":"#303030FF","buttonPressedBorderColor":"#646464FF","buttonPressedTextColor":"#969696FF","highlightColor":"#0078FFFF"},"devices":[{"id":"55F16B35884A859CCE4FFA1FC8D3DE5B","name":"Device Name","size":{"columns":5,"rows":3},"type":0}]}`,
	})
	if err != nil {
		t.Fatalf("Failed to parse params: %v", err)
	}
	if p.Port != 12345 {
		t.Fatalf("port = %d", p.Port)
	}
	if p.Info.Application.Language != streamdeck.LanguageEN {
		t.Fatalf("language = %q", p.Info.Application.Language)
	}
	if p.Info.Application.Platform != streamdeck.PlatformMac {
		t.Fatalf("platform = %q", p.Info.Application.Platform)
	}
	if p.Info.Colors.ButtonMouseOverBackgroundColor == "" {
		t.Fatal("missing buttonMouseOverBackgroundColor")
	}
	if len(p.Info.Devices) != 1 || p.Info.Devices[0].Name != "Device Name" {
		t.Fatalf("devices = %+v", p.Info.Devices)
	}
	if p.Info.Devices[0].Type != streamdeck.StreamDeck {
		t.Fatalf("device type = %d", p.Info.Devices[0].Type)
	}
}

func TestParseRegistrationParamsMissingFlags(t *testing.T) {
	_, err := streamdeck.ParseRegistrationParams([]string{"", "-pluginUUID", "x", "-registerEvent", "r", "-info", "{}"})
	if err != streamdeck.ErrMissingPortFlag {
		t.Fatalf("err = %v", err)
	}
}
