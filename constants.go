package streamdeck

// EventName is a Stream Deck WebSocket event or command name.
type EventName string

const (
	// Events received from Stream Deck.

	DidReceiveSettings            EventName = "didReceiveSettings"
	DidReceiveGlobalSettings      EventName = "didReceiveGlobalSettings"
	DidReceiveResources           EventName = "didReceiveResources"
	DidReceiveSecrets             EventName = "didReceiveSecrets"
	DidReceiveDeepLink            EventName = "didReceiveDeepLink"
	KeyDown                       EventName = "keyDown"
	KeyUp                         EventName = "keyUp"
	TouchTap                      EventName = "touchTap"
	DialDown                      EventName = "dialDown"
	DialUp                        EventName = "dialUp"
	DialRotate                    EventName = "dialRotate"
	WillAppear                    EventName = "willAppear"
	WillDisappear                 EventName = "willDisappear"
	TitleParametersDidChange      EventName = "titleParametersDidChange"
	DeviceDidConnect              EventName = "deviceDidConnect"
	DeviceDidDisconnect           EventName = "deviceDidDisconnect"
	DeviceDidChange               EventName = "deviceDidChange"
	ApplicationDidLaunch          EventName = "applicationDidLaunch"
	ApplicationDidTerminate       EventName = "applicationDidTerminate"
	SystemDidWakeUp               EventName = "systemDidWakeUp"
	PropertyInspectorDidAppear    EventName = "propertyInspectorDidAppear"
	PropertyInspectorDidDisappear EventName = "propertyInspectorDidDisappear"
	SendToPlugin                  EventName = "sendToPlugin"
	SendToPropertyInspector       EventName = "sendToPropertyInspector"

	// Commands sent to Stream Deck.

	SetSettings           EventName = "setSettings"
	GetSettings           EventName = "getSettings"
	SetGlobalSettings     EventName = "setGlobalSettings"
	GetGlobalSettings     EventName = "getGlobalSettings"
	GetResources          EventName = "getResources"
	SetResources          EventName = "setResources"
	GetSecrets            EventName = "getSecrets"
	OpenURL               EventName = "openUrl"
	LogMessage            EventName = "logMessage"
	SetTitle              EventName = "setTitle"
	SetImage              EventName = "setImage"
	SetFeedback           EventName = "setFeedback"
	SetFeedbackLayout     EventName = "setFeedbackLayout"
	SetTriggerDescription EventName = "setTriggerDescription"
	ShowAlert             EventName = "showAlert"
	ShowOk                EventName = "showOk"
	SetState              EventName = "setState"
	SwitchToProfile       EventName = "switchToProfile"
)

// Target specifies whether a title or image update applies to hardware, software, or both.
type Target int

const (
	HardwareAndSoftware Target = iota
	OnlyHardware
	OnlySoftware
)

// Controller identifies the input type of an action instance.
type Controller string

const (
	ControllerKeypad  Controller = "Keypad"
	ControllerEncoder Controller = "Encoder"
	ControllerNeo     Controller = "Neo"
)

// Language is a Stream Deck application locale.
type Language string

const (
	LanguageDE   Language = "de"
	LanguageEN   Language = "en"
	LanguageES   Language = "es"
	LanguageFR   Language = "fr"
	LanguageJA   Language = "ja"
	LanguageKO   Language = "ko"
	LanguageZHCN Language = "zh_CN"
	LanguageZHTW Language = "zh_TW"
)

// Platform is the host operating system reported by Stream Deck.
type Platform string

const (
	PlatformMac     Platform = "mac"
	PlatformWindows Platform = "windows"
)
