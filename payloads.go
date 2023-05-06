package streamdeck

// LogMessagePayload A string to write to the logs file.
type LogMessagePayload struct {
	Message string `json:"message"`
}

// OpenURLPayload An URL to open in the default browser.
type OpenURLPayload struct {
	URL string `json:"url"`
}

// SetTitlePayload The title to display. If there is no title parameter, the title is reset to the title set by the user.
type SetTitlePayload struct {
	Title  string `json:"title"`
	Target Target `json:"target"`
	State  int    `json:"state"`
}

// SetImagePayload The image to display encoded in base64 with the image format declared in the mime type (PNG, JPEG, BMP, ...). svg is also supported. If no image is passed, the image is reset to the default image from the manifest.
type SetImagePayload struct {
	Base64Image string `json:"image"`
	Target      Target `json:"target"`
	State       int    `json:"state"`
}

type SetFeedbackLayoutPayload struct {
	Layout string `json:"layout"`
}

// SetStatePayload A 0-based integer value representing the state requested.
type SetStatePayload struct {
	State int `json:"state"`
}

// SwitchProfilePayload The name of the profile to switch to. The name should be identical to the name provided in the manifest.json file.
type SwitchProfilePayload struct {
	Profile string `json:"profile"`
}

// DidReceiveSettingsPayload This json object contains persistently stored data.
type DidReceiveSettingsPayload[T any] struct {
	Settings        T           `json:"settings,omitempty"`
	Coordinates     Coordinates `json:"coordinates,omitempty"`
	IsInMultiAction bool        `json:"isInMultiAction,omitempty"`
}

// Coordinates The coordinates of the action triggered.
type Coordinates struct {
	Column int `json:"column,omitempty"`
	Row    int `json:"row,omitempty"`
}

// DidReceiveGlobalSettingsPayload This json object contains persistently stored data.
type DidReceiveGlobalSettingsPayload[T any] struct {
	Settings T `json:"settings,omitempty"`
}

// KeyDownPayload A json object
type KeyDownPayload[T any] struct {
	Settings         T           `json:"settings,omitempty"`
	Coordinates      Coordinates `json:"coordinates,omitempty"`
	State            int         `json:"state,omitempty"`
	UserDesiredState int         `json:"userDesiredState,omitempty"`
	IsInMultiAction  bool        `json:"isInMultiAction,omitempty"`
}

// KeyUpPayload A json object
type KeyUpPayload[T any] struct {
	Settings         T           `json:"settings,omitempty"`
	Coordinates      Coordinates `json:"coordinates,omitempty"`
	State            int         `json:"state,omitempty"`
	UserDesiredState int         `json:"userDesiredState,omitempty"`
	IsInMultiAction  bool        `json:"isInMultiAction,omitempty"`
}

// TouchTapPayload A json object
type TouchTapPayload[T any] struct {
	Settings    T           `json:"settings,omitempty"`
	Coordinates Coordinates `json:"coordinates,omitempty"`
	TapPos      [2]int      `json:"tapPos,omitempty"`
	Hold        bool        `json:"hold,omitempty"`
}

type DialDownPayload[T any] struct {
	Settings    T           `json:"settings,omitempty"`
	Coordinates Coordinates `json:"coordinates,omitempty"`
	Controller  string      `json:"controller,omitempty"` // Encoder
}

type DialUpPayload[T any] struct {
	Settings    T           `json:"settings,omitempty"`
	Coordinates Coordinates `json:"coordinates,omitempty"`
	Controller  string      `json:"controller,omitempty"` // Encoder
}

type DialRotatePayload[T any] struct {
	Settings    T           `json:"settings,omitempty"`
	Coordinates Coordinates `json:"coordinates,omitempty"`
	Ticks       int         `json:"ticks,omitempty"`
	Pressed     bool        `json:"pressed,omitempty"`
}

// WillAppearPayload A json object
type WillAppearPayload[T any] struct {
	Settings        T           `json:"settings,omitempty"`
	Coordinates     Coordinates `json:"coordinates,omitempty"`
	State           int         `json:"state,omitempty"`
	IsInMultiAction bool        `json:"isInMultiAction,omitempty"`
}

// WillDisappearPayload A json object
type WillDisappearPayload[T any] struct {
	Settings        T           `json:"settings,omitempty"`
	Coordinates     Coordinates `json:"coordinates,omitempty"`
	State           int         `json:"state,omitempty"`
	IsInMultiAction bool        `json:"isInMultiAction,omitempty"`
}

// TitleParametersDidChangePayload A json object
type TitleParametersDidChangePayload[T any] struct {
	Settings        T               `json:"settings,omitempty"`
	Coordinates     Coordinates     `json:"coordinates,omitempty"`
	State           int             `json:"state,omitempty"`
	Title           string          `json:"title,omitempty"`
	TitleParameters TitleParameters `json:"titleParameters,omitempty"`
}

// TitleParameters A json object
type TitleParameters struct {
	FontFamily     string `json:"fontFamily,omitempty"`
	FontSize       int    `json:"fontSize,omitempty"`
	FontStyle      string `json:"fontStyle,omitempty"`
	FontUnderline  bool   `json:"fontUnderline,omitempty"`
	ShowTitle      bool   `json:"showTitle,omitempty"`
	TitleAlignment string `json:"titleAlignment,omitempty"`
	TitleColor     string `json:"titleColor,omitempty"`
}

// ApplicationDidLaunchPayload A json object
type ApplicationDidLaunchPayload struct {
	Application string `json:"application,omitempty"`
}

// ApplicationDidTerminatePayload A json object
type ApplicationDidTerminatePayload struct {
	Application string `json:"application,omitempty"`
}
