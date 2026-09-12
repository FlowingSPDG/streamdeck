package streamdeck

import (
	"encoding/json"
	"flag"
	"fmt"
)

// RegistrationParams are the command-line values Stream Deck passes to a plugin.
type RegistrationParams struct {
	Port          int
	PluginUUID    string
	RegisterEvent string
	Info          Info
}

// Application describes the Stream Deck application.
type Application struct {
	Font            string   `json:"font"`
	Language        Language `json:"language"`
	Platform        Platform `json:"platform"`
	PlatformVersion string   `json:"platformVersion"`
	Version         string   `json:"version"`
}

// Plugin describes the running plugin.
type Plugin struct {
	UUID    string `json:"uuid"`
	Version string `json:"version"`
}

// Colors are the Stream Deck application's preferred colors.
type Colors struct {
	ButtonPressedBackgroundColor   string `json:"buttonPressedBackgroundColor"`
	ButtonPressedBorderColor       string `json:"buttonPressedBorderColor"`
	ButtonPressedTextColor         string `json:"buttonPressedTextColor"`
	ButtonMouseOverBackgroundColor string `json:"buttonMouseOverBackgroundColor"`
	HighlightColor                 string `json:"highlightColor"`
}

// Size is a device grid size.
type Size struct {
	Columns int `json:"columns"`
	Rows    int `json:"rows"`
}

// Device is a device listed in registration info.
type Device struct {
	ID   string     `json:"id"`
	Name string     `json:"name"`
	Size Size       `json:"size"`
	Type DeviceType `json:"type"`
}

// ActionInfoPayload is the payload of a property inspector actionInfo object.
type ActionInfoPayload[SettingsT any] struct {
	Coordinates Coordinates `json:"coordinates,omitempty"`
	Settings    SettingsT   `json:"settings,omitempty"`
}

// ActionInfo describes the selected action when a property inspector starts.
type ActionInfo[SettingsT any] struct {
	Action  string                       `json:"action"`
	Context string                       `json:"context"`
	Device  string                       `json:"device"`
	Payload ActionInfoPayload[SettingsT] `json:"payload"`
}

// Info is the registration -info JSON object.
type Info struct {
	Application      Application `json:"application"`
	Plugin           Plugin      `json:"plugin"`
	DevicePixelRatio int         `json:"devicePixelRatio"`
	Colors           Colors      `json:"colors"`
	Devices          []Device    `json:"devices"`
}

// ParseRegistrationParams parses Stream Deck launch arguments. Pass os.Args.
func ParseRegistrationParams(args []string) (RegistrationParams, error) {
	f := flag.NewFlagSet("registration_params", flag.ContinueOnError)

	ret := RegistrationParams{}

	port := f.Int("port", -1, "")
	pluginUUID := f.String("pluginUUID", "", "")
	registerEvent := f.String("registerEvent", "", "")
	info := f.String("info", "", "")

	if err := f.Parse(args[1:]); err != nil {
		return ret, err
	}

	if *port == -1 {
		return ret, ErrMissingPortFlag
	}
	ret.Port = *port

	if *pluginUUID == "" {
		return ret, ErrMissingPluginUUIDFlag
	}
	ret.PluginUUID = *pluginUUID

	if *registerEvent == "" {
		return ret, ErrMissingRegisterEventFlag
	}
	ret.RegisterEvent = *registerEvent

	if *info == "" {
		return ret, ErrMissingInfoFlag
	}
	if err := json.Unmarshal([]byte(*info), &ret.Info); err != nil {
		return ret, fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}

	return ret, nil
}
