package streamdeck

// LogMessagePayload is written to the Stream Deck log file.
type LogMessagePayload struct {
	Message string `json:"message"`
}

// OpenURLPayload opens a URL in the default browser.
type OpenURLPayload struct {
	URL string `json:"url"`
}

// SetTitlePayload updates the title of an action instance.
type SetTitlePayload struct {
	Title  string `json:"title,omitempty"`
	Target Target `json:"target,omitempty"`
	State  *int   `json:"state,omitempty"`
}

// SetImagePayload updates the image of an action instance.
type SetImagePayload struct {
	Base64Image string `json:"image,omitempty"`
	Target      Target `json:"target,omitempty"`
	State       *int   `json:"state,omitempty"`
}

// SetFeedbackLayoutPayload selects a Stream Deck + layout.
type SetFeedbackLayoutPayload struct {
	Layout string `json:"layout"`
}

// SetTriggerDescriptionPayload updates encoder trigger descriptions.
type SetTriggerDescriptionPayload struct {
	LongTouch string `json:"longTouch,omitempty"`
	Push      string `json:"push,omitempty"`
	Rotate    string `json:"rotate,omitempty"`
	Touch     string `json:"touch,omitempty"`
}

// SetStatePayload selects a multi-state action state.
type SetStatePayload struct {
	State int `json:"state"`
}

// SwitchProfilePayload switches to a plugin-distributed profile.
type SwitchProfilePayload struct {
	Profile string `json:"profile,omitempty"`
	Page    *int   `json:"page,omitempty"`
}

// Coordinates locate an action on a device grid.
type Coordinates struct {
	Column int `json:"column"`
	Row    int `json:"row"`
}

// DidReceiveSettingsPayload is returned for getSettings and property inspector edits.
type DidReceiveSettingsPayload[T any] struct {
	Settings        T                 `json:"settings,omitempty"`
	Coordinates     *Coordinates      `json:"coordinates,omitempty"`
	IsInMultiAction bool              `json:"isInMultiAction,omitempty"`
	Controller      Controller        `json:"controller,omitempty"`
	Resources       map[string]string `json:"resources,omitempty"`
	State           *int              `json:"state,omitempty"`
}

// DidReceiveResourcesPayload is returned for getResources and resource edits.
type DidReceiveResourcesPayload[T any] struct {
	Settings        T                 `json:"settings,omitempty"`
	Coordinates     *Coordinates      `json:"coordinates,omitempty"`
	IsInMultiAction bool              `json:"isInMultiAction,omitempty"`
	Controller      Controller        `json:"controller,omitempty"`
	Resources       map[string]string `json:"resources,omitempty"`
	State           *int              `json:"state,omitempty"`
}

// DidReceiveGlobalSettingsPayload contains plugin-wide settings.
type DidReceiveGlobalSettingsPayload[T any] struct {
	Settings T `json:"settings,omitempty"`
}

// DidReceiveSecretsPayload contains plugin secrets.
type DidReceiveSecretsPayload[T any] struct {
	Secrets T `json:"secrets,omitempty"`
}

// KeyDownPayload is sent when a key is pressed.
type KeyDownPayload[T any] struct {
	Settings         T                 `json:"settings,omitempty"`
	Coordinates      *Coordinates      `json:"coordinates,omitempty"`
	State            *int              `json:"state,omitempty"`
	UserDesiredState *int              `json:"userDesiredState,omitempty"`
	IsInMultiAction  bool              `json:"isInMultiAction,omitempty"`
	Controller       Controller        `json:"controller,omitempty"`
	Resources        map[string]string `json:"resources,omitempty"`
}

// KeyUpPayload is sent when a key is released.
type KeyUpPayload[T any] struct {
	Settings         T                 `json:"settings,omitempty"`
	Coordinates      *Coordinates      `json:"coordinates,omitempty"`
	State            *int              `json:"state,omitempty"`
	UserDesiredState *int              `json:"userDesiredState,omitempty"`
	IsInMultiAction  bool              `json:"isInMultiAction,omitempty"`
	Controller       Controller        `json:"controller,omitempty"`
	Resources        map[string]string `json:"resources,omitempty"`
}

// TouchTapPayload is sent when a Stream Deck + touchscreen is tapped.
type TouchTapPayload[T any] struct {
	Settings    T                 `json:"settings,omitempty"`
	Coordinates *Coordinates      `json:"coordinates,omitempty"`
	TapPos      [2]int            `json:"tapPos,omitempty"`
	Hold        bool              `json:"hold,omitempty"`
	Controller  Controller        `json:"controller,omitempty"`
	Resources   map[string]string `json:"resources,omitempty"`
}

// DialDownPayload is sent when a Stream Deck + dial is pressed.
type DialDownPayload[T any] struct {
	Settings    T                 `json:"settings,omitempty"`
	Coordinates *Coordinates      `json:"coordinates,omitempty"`
	Controller  Controller        `json:"controller,omitempty"`
	Resources   map[string]string `json:"resources,omitempty"`
}

// DialUpPayload is sent when a Stream Deck + dial is released.
type DialUpPayload[T any] struct {
	Settings    T                 `json:"settings,omitempty"`
	Coordinates *Coordinates      `json:"coordinates,omitempty"`
	Controller  Controller        `json:"controller,omitempty"`
	Resources   map[string]string `json:"resources,omitempty"`
}

// DialRotatePayload is sent when a Stream Deck + dial is rotated.
type DialRotatePayload[T any] struct {
	Settings    T                 `json:"settings,omitempty"`
	Coordinates *Coordinates      `json:"coordinates,omitempty"`
	Ticks       int               `json:"ticks,omitempty"`
	Pressed     bool              `json:"pressed,omitempty"`
	Controller  Controller        `json:"controller,omitempty"`
	Resources   map[string]string `json:"resources,omitempty"`
}

// WillAppearPayload is sent when an action instance becomes visible.
type WillAppearPayload[T any] struct {
	Settings        T                 `json:"settings,omitempty"`
	Coordinates     *Coordinates      `json:"coordinates,omitempty"`
	State           *int              `json:"state,omitempty"`
	IsInMultiAction bool              `json:"isInMultiAction,omitempty"`
	Controller      Controller        `json:"controller,omitempty"`
	Resources       map[string]string `json:"resources,omitempty"`
}

// WillDisappearPayload is sent when an action instance is hidden.
type WillDisappearPayload[T any] struct {
	Settings        T                 `json:"settings,omitempty"`
	Coordinates     *Coordinates      `json:"coordinates,omitempty"`
	State           *int              `json:"state,omitempty"`
	IsInMultiAction bool              `json:"isInMultiAction,omitempty"`
	Controller      Controller        `json:"controller,omitempty"`
	Resources       map[string]string `json:"resources,omitempty"`
}

// TitleParametersDidChangePayload is sent when the user edits title parameters.
type TitleParametersDidChangePayload[T any] struct {
	Settings        T                 `json:"settings,omitempty"`
	Coordinates     *Coordinates      `json:"coordinates,omitempty"`
	State           *int              `json:"state,omitempty"`
	Title           string            `json:"title,omitempty"`
	TitleParameters TitleParameters   `json:"titleParameters,omitempty"`
	Controller      Controller        `json:"controller,omitempty"`
	Resources       map[string]string `json:"resources,omitempty"`
}

// TitleParameters describes how a title is rendered.
type TitleParameters struct {
	FontFamily     string `json:"fontFamily,omitempty"`
	FontSize       int    `json:"fontSize,omitempty"`
	FontStyle      string `json:"fontStyle,omitempty"`
	FontUnderline  bool   `json:"fontUnderline,omitempty"`
	ShowTitle      bool   `json:"showTitle,omitempty"`
	TitleAlignment string `json:"titleAlignment,omitempty"`
	TitleColor     string `json:"titleColor,omitempty"`
}

// ApplicationDidLaunchPayload names a monitored application that launched.
type ApplicationDidLaunchPayload struct {
	Application string `json:"application,omitempty"`
}

// ApplicationDidTerminatePayload names a monitored application that exited.
type ApplicationDidTerminatePayload struct {
	Application string `json:"application,omitempty"`
}

// PropertyInspectorDidAppearPayload is sent when the property inspector opens.
type PropertyInspectorDidAppearPayload[T any] struct {
	Settings        T                 `json:"settings,omitempty"`
	Coordinates     *Coordinates      `json:"coordinates,omitempty"`
	State           *int              `json:"state,omitempty"`
	IsInMultiAction bool              `json:"isInMultiAction,omitempty"`
	Controller      Controller        `json:"controller,omitempty"`
	Resources       map[string]string `json:"resources,omitempty"`
}

// PropertyInspectorDidDisappearPayload is sent when the property inspector closes.
type PropertyInspectorDidDisappearPayload[T any] struct {
	Settings        T                 `json:"settings,omitempty"`
	Coordinates     *Coordinates      `json:"coordinates,omitempty"`
	State           *int              `json:"state,omitempty"`
	IsInMultiAction bool              `json:"isInMultiAction,omitempty"`
	Controller      Controller        `json:"controller,omitempty"`
	Resources       map[string]string `json:"resources,omitempty"`
}

// DidReceiveDeepLinkPayload contains a deep-link path with the prefix omitted.
type DidReceiveDeepLinkPayload struct {
	URL string `json:"url,omitempty"`
}
