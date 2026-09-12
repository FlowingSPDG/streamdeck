package streamdeck

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEventUnmarshalKeyDown(t *testing.T) {
	raw := mustReadTestdata(t, "events", "keyDown.json")
	var ev Event
	if err := json.Unmarshal(raw, &ev); err != nil {
		t.Fatal(err)
	}
	if ev.Event != KeyDown {
		t.Fatalf("event = %q", ev.Event)
	}
	p, err := ev.Unmarshal[KeyDownPayload[map[string]int]]()
	if err != nil {
		t.Fatal(err)
	}
	if p.Controller != ControllerKeypad {
		t.Fatalf("controller = %q", p.Controller)
	}
	if p.Coordinates == nil || p.Coordinates.Column != 3 {
		t.Fatalf("coordinates = %+v", p.Coordinates)
	}
	if p.Settings["counter"] != 2 {
		t.Fatalf("settings = %+v", p.Settings)
	}
}

func TestEventUnmarshalNilPayload(t *testing.T) {
	ev := Event{}
	p, err := ev.Unmarshal[LogMessagePayload]()
	if err != nil {
		t.Fatal(err)
	}
	if p.Message != "" {
		t.Fatalf("unexpected %+v", p)
	}
}

func TestDeviceDidConnectUsesName(t *testing.T) {
	raw := mustReadTestdata(t, "events", "deviceDidConnect.json")
	var ev Event
	if err := json.Unmarshal(raw, &ev); err != nil {
		t.Fatal(err)
	}
	if ev.DeviceInfo.Name != "Stream Deck" {
		t.Fatalf("name = %q", ev.DeviceInfo.Name)
	}
	if ev.DeviceInfo.Type != StreamDeck {
		t.Fatalf("type = %d", ev.DeviceInfo.Type)
	}
}

func TestGoldenEventsDecode(t *testing.T) {
	dir := filepath.Join("testdata", "events")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		t.Run(e.Name(), func(t *testing.T) {
			raw := mustReadTestdata(t, "events", e.Name())
			var ev Event
			if err := json.Unmarshal(raw, &ev); err != nil {
				t.Fatal(err)
			}
			if ev.Event == "" {
				t.Fatal("empty event name")
			}
		})
	}
}

func TestSetTitleOmitsStateWhenUnset(t *testing.T) {
	data, err := json.Marshal(SetTitlePayload{Title: "Hi", Target: HardwareAndSoftware})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["state"]; ok {
		t.Fatalf("state should be omitted: %s", data)
	}
}

func TestSetTitleIncludesStateZero(t *testing.T) {
	data, err := json.Marshal(SetTitlePayload{Title: "Hi", State: new(0)})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if m["state"] != float64(0) {
		t.Fatalf("state = %v", m["state"])
	}
}

func TestSwitchProfileOmitsPageWhenUnset(t *testing.T) {
	data, err := json.Marshal(SwitchProfilePayload{Profile: "Default"})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["page"]; ok {
		t.Fatalf("page should be omitted: %s", data)
	}
}

func TestCommandGoldens(t *testing.T) {
	dir := filepath.Join("testdata", "commands")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		t.Run(e.Name(), func(t *testing.T) {
			raw := mustReadTestdata(t, "commands", e.Name())
			var ev outgoingEvent
			if err := json.Unmarshal(raw, &ev); err != nil {
				t.Fatal(err)
			}
			if ev.Event == "" {
				t.Fatal("empty command name")
			}
		})
	}
}

func mustReadTestdata(t *testing.T, parts ...string) []byte {
	t.Helper()
	p := filepath.Join(append([]string{"testdata"}, parts...)...)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
