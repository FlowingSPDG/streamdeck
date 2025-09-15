package streamdeck

import (
	"encoding/json"
	"flag"

	"golang.org/x/xerrors"
)

// RegistrationParams Params for registering streamdeck plugin.
type RegistrationParams struct {
	Port          int
	PluginUUID    string
	RegisterEvent string
	Info          Info
}

type Application struct {
	Font            string `json:"font"`
	Language        string `json:"language"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platformVersion"`
	Version         string `json:"version"`
}

type Plugin struct {
	UUID    string `json:"uuid"`
	Version string `json:"version"`
}

type Colors struct {
	ButtonPressedBackgroundColor string `json:"buttonPressedBackgroundColor"`
	ButtonPressedBorderColor     string `json:"buttonPressedBorderColor"`
	ButtonPressedTextColor       string `json:"buttonPressedTextColor"`
	DisabledColor                string `json:"disabledColor"`
	HighlightColor               string `json:"highlightColor"`
	MouseDownColor               string `json:"mouseDownColor"`
}

type Size struct {
	Columns int `json:"columns"`
	Rows    int `json:"rows"`
}

type Device struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size Size   `json:"size"`
	Type int    `json:"type"`
}

type ActionInfoPayload[SettingsT any] struct {
	Coordinates Coordinates `json:"coordinates,omitempty"`
	Settings    SettingsT   `json:"settings,omitempty"`
}

type ActionInfo[SettingsT any] struct {
	Action  string                       `json:"action"`
	Context string                       `json:"context"`
	Device  string                       `json:"device"`
	Payload ActionInfoPayload[SettingsT] `json:"payload"`
}

type Info struct {
	Application      Application `json:"application"`
	Plugin           Plugin      `json:"plugin"`
	DevicePixelRatio int         `json:"devicePixelRatio"`
	Colors           Colors      `json:"colors"`
	Devices          []Device    `json:"devices"`
}

// ParseRegistrationParams Parse parameters. Normally you should use os.Args .
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
		return ret, xerrors.Errorf("%w", ErrMissingPortFlag)
	}
	ret.Port = *port

	if *pluginUUID == "" {
		return ret, xerrors.Errorf("%w", ErrMissingPluginUUIDFlag)
	}
	ret.PluginUUID = *pluginUUID

	if *registerEvent == "" {
		return ret, xerrors.Errorf("%w", ErrMissingRegisterEventFlag)
	}
	ret.RegisterEvent = *registerEvent

	if *info == "" {
		return ret, xerrors.Errorf("%w", ErrMissingInfoFlag)
	}
	infob := []byte(*info)
	if err := json.Unmarshal(infob, &ret.Info); err != nil {
		return ret, err
	}

	return ret, nil
}
