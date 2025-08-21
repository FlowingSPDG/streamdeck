package streamdeck

const (
	// DidReceiveSettings Event received after calling the getSettings API to retrieve the persistent data stored for the action.
	DidReceiveSettings = "didReceiveSettings"
	// DidReceiveGlobalSettings Event received after calling the getGlobalSettings API to retrieve the global persistent data.
	DidReceiveGlobalSettings = "didReceiveGlobalSettings"
	// KeyDown When the user presses a key, the plugin will receive the keyDown event.
	KeyDown = "keyDown"
	// KeyUp When the user releases a key, the plugin will receive the keyUp event.
	KeyUp = "keyUp"
	// TouchTap When the user touches the display, the plugin will receive the touchTap event.
	TouchTap = "touchTap"
	// DialDown When the user presses the encoder down, the plugin will receive the dialDown event (SD+).
	DialDown = "dialDown"
	// DialUp When the user releases a pressed encoder, the plugin will receive the dialUp event (SD+).
	DialUp = "dialUp"
	// DialRotate When the user rotates the encoder, the plugin will receive the dialRotate event.
	DialRotate = "dialRotate"
	// WillAppear When an instance of an action is displayed on the Stream Deck, for example when the hardware is first plugged in, or when a folder containing that action is entered, the plugin will receive a willAppear event.
	WillAppear = "willAppear"
	// WillDisappear When an instance of an action ceases to be displayed on Stream Deck, for example when switching profiles or folders, the plugin will receive a willDisappear event.
	WillDisappear = "willDisappear"
	// TitleParametersDidChange When the user changes the title or title parameters, the plugin will receive a titleParametersDidChange event
	TitleParametersDidChange = "titleParametersDidChange"
	//DeviceDidConnect When a device is plugged to the computer, the plugin will receive a deviceDidConnect event.
	DeviceDidConnect = "deviceDidConnect"
	//DeviceDidDisconnect When a device is unplugged from the computer, the plugin will receive a deviceDidDisconnect event.
	DeviceDidDisconnect = "deviceDidDisconnect"
	//DeviceDidChange When a Stream Deck device changed, for example its name or size. Available from Stream Deck 7.0.
	DeviceDidChange = "deviceDidChange"
	//ApplicationDidLaunch When a monitored application is launched, the plugin will be notified and will receive the applicationDidLaunch event.
	ApplicationDidLaunch = "applicationDidLaunch"
	//ApplicationDidTerminate When a monitored application is terminated, the plugin will be notified and will receive the applicationDidTerminate event.
	ApplicationDidTerminate = "applicationDidTerminate"
	//SystemDidWakeUp 	When the computer is wake up, the plugin will be notified and will receive the systemDidWakeUp event.
	SystemDidWakeUp = "systemDidWakeUp"
	//PropertyInspectorDidAppear Event received when the Property Inspector appears in the Stream Deck software user interface, for example when selecting a new instance.
	PropertyInspectorDidAppear = "propertyInspectorDidAppear"
	//PropertyInspectorDidDisappear Event received when the Property Inspector for an instance is removed from the Stream Deck software user interface, for example when selecting a different instance.
	PropertyInspectorDidDisappear = "propertyInspectorDidDisappear"
	//SendToPlugin 	Event received by the plugin when the Property Inspector uses the sendToPlugin event.
	SendToPlugin = "sendToPlugin"
	// SendToPropertyInspector Event received by the Property Inspector when the plugin uses the sendToPropertyInspector event.
	SendToPropertyInspector = "sendToPropertyInspector"
	// DidReceiveDeepLink Event received when a deep link is activated
	DidReceiveDeepLink = "didReceiveDeepLink"

	// SetSettings Save data persistently for the action's instance.
	SetSettings = "setSettings"
	// GetSettings Request the persistent data for the action's instance.
	GetSettings = "getSettings"
	// SetGlobalSettings Save data securely and globally for the plugin.
	SetGlobalSettings = "setGlobalSettings"
	// GetGlobalSettings Request the global persistent data.
	GetGlobalSettings = "getGlobalSettings"
	// OpenURL Open an URL in the default browser.
	OpenURL = "openUrl"
	// LogMessage Write a debug log to the logs file.
	LogMessage = "logMessage"
	// SetTitle Dynamically change the title of an instance of an action.
	SetTitle = "setTitle"
	// SetImage Dynamically change the image displayed by an instance of an action.
	SetImage = "setImage"
	// SetFeedback
	SetFeedback = "setFeedback"
	// ShowAlert Temporarily show an alert icon on the image displayed by an instance of an action.
	ShowAlert = "showAlert"
	// ShowOk Temporarily show an OK checkmark icon on the image displayed by an instance of an action.
	ShowOk = "showOk"
	// SetState Change the state of the action's instance supporting multiple states.
	SetState = "setState"
	// SwitchToProfile Switch to one of the preconfigured read-only profiles.
	SwitchToProfile = "switchToProfile"
	// SetFeedbackLayout Sets the layout associated with an action instance.
	SetFeedbackLayout = "setFeedbackLayout"
	// SetTriggerDescription Sets the trigger descriptions associated with an encoder action instance.
	SetTriggerDescription = "setTriggerDescription"
)

// Target Specify if you want to display the title on the hardware and software (0), only on the hardware (1) or only on the software (2). Default is 0.
type Target int

const (
	// HardwareAndSoftware hardware and software(0)
	HardwareAndSoftware Target = iota
	// OnlyHardware only on the hardware (1)
	OnlyHardware
	// OnlySoftware only on the software (2)
	OnlySoftware
)
