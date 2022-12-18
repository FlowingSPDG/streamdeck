package wasm

import (
	"context"
	"net/url"
	"sync"

	"github.com/FlowingSPDG/streamdeck"
	"github.com/gorilla/websocket"
)

type SDClient interface {
	SetSettings(ctx context.Context, settings any) error
	GetSettings(ctx context.Context) error
	SetGlobalSettings(ctx context.Context, settings any) error
	GetGlobalSettings(ctx context.Context) error
	OpenURL(ctx context.Context, u *url.URL) error
	LogMessage(message string) error
	SetTitle(ctx context.Context, title string, target streamdeck.Target) error
	SetImage(ctx context.Context, base64image string, target streamdeck.Target) error
	ShowAlert(ctx context.Context) error
	ShowOk(ctx context.Context) error
	SetState(ctx context.Context, state int) error
	SwitchToProfile(ctx context.Context, profile string) error
	SendToPropertyInspector(ctx context.Context, payload any) error
	SendToPlugin(ctx context.Context, payload any) error
}

type sdClient[SettingsT any] struct {
	c                 *websocket.Conn
	uuid              string
	registerEventName string
	actionInfo        inActionInfo[SettingsT]
	inInfo            inInfo
	// runningApps // ?
	isQT bool

	// Send Mutex lock
	sendMutex *sync.Mutex
}

// TODO: WSから受信したメッセージからハンドラを起動する
func (sd *sdClient[SettingsT]) send(event streamdeck.Event) error {
	sd.sendMutex.Lock()
	defer sd.sendMutex.Unlock()
	return sd.c.WriteJSON(event)
}

// SetSettings Save data persistently for the action's instance.
func (sd *sdClient[SettingsT]) SetSettings(ctx context.Context, settings any) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.SetSettings, settings))
}

// GetSettings Request the persistent data for the action's instance.
func (sd *sdClient[SettingsT]) GetSettings(ctx context.Context) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.GetSettings, nil))
}

// SetGlobalSettings Save data securely and globally for the plugin.
func (sd *sdClient[SettingsT]) SetGlobalSettings(ctx context.Context, settings any) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.SetGlobalSettings, settings))
}

// GetGlobalSettings Request the global persistent data
func (sd *sdClient[SettingsT]) GetGlobalSettings(ctx context.Context) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.GetGlobalSettings, nil))
}

// OpenURL Open an URL in the default browser.
func (sd *sdClient[SettingsT]) OpenURL(ctx context.Context, u *url.URL) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.OpenURL, streamdeck.OpenURLPayload{URL: u.String()}))
}

// LogMessage Write a debug log to the logs file.
func (sd *sdClient[SettingsT]) LogMessage(message string) error {
	return sd.send(streamdeck.NewEvent(nil, streamdeck.LogMessage, streamdeck.LogMessagePayload{Message: message}))
}

// SetTitle Dynamically change the title of an instance of an action.
func (sd *sdClient[SettingsT]) SetTitle(ctx context.Context, title string, target streamdeck.Target) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.SetTitle, streamdeck.SetTitlePayload{Title: title, Target: target}))
}

// SetImage Dynamically change the image displayed by an instance of an action.
func (sd *sdClient[SettingsT]) SetImage(ctx context.Context, base64image string, target streamdeck.Target) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.SetImage, streamdeck.SetImagePayload{Base64Image: base64image, Target: target}))
}

// ShowAlert Temporarily show an alert icon on the image displayed by an instance of an action.
func (sd *sdClient[SettingsT]) ShowAlert(ctx context.Context) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.ShowAlert, nil))
}

// ShowOk Temporarily show an OK checkmark icon on the image displayed by an instance of an action
func (sd *sdClient[SettingsT]) ShowOk(ctx context.Context) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.ShowOk, nil))
}

// SetState Change the state of the action's instance supporting multiple states.
func (sd *sdClient[SettingsT]) SetState(ctx context.Context, state int) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.SetState, streamdeck.SetStatePayload{State: state}))
}

// SwitchToProfile Switch to one of the preconfigured read-only profiles.
func (sd *sdClient[SettingsT]) SwitchToProfile(ctx context.Context, profile string) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.SwitchToProfile, streamdeck.SwitchProfilePayload{Profile: profile}))
}

// SendToPropertyInspector Send a payload to the Property Inspector.
func (sd *sdClient[SettingsT]) SendToPropertyInspector(ctx context.Context, payload any) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.SendToPropertyInspector, payload))
}

// SendToPlugin Send a payload to the plugin.
func (sd *sdClient[SettingsT]) SendToPlugin(ctx context.Context, payload any) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.SendToPlugin, payload))
}
