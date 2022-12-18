package wasm

// inInfo...
type inInfo struct {
	Application      Application `json:"application"`
	Plugin           Plugin      `json:"plugin"`
	DevicePixelRatio int         `json:"devicePixelRatio"`
	Colors           Colors      `json:"colors"`
	Devices          []Devices   `json:"devices"`
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
type Devices struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size Size   `json:"size"`
	Type int    `json:"type"`
}

type inActionInfo[SettingsT any] struct {
	Action  string             `json:"action"`
	Context string             `json:"context"`
	Device  string             `json:"device"`
	Payload Payload[SettingsT] `json:"payload"`
}
type Coordinates struct {
	Column int `json:"column"`
	Row    int `json:"row"`
}
type Payload[T any] struct {
	Settings    T           `json:"settings"`
	Coordinates Coordinates `json:"coordinates"`
}
