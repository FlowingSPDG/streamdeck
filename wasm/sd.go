package wasm

import (
	"context"
	"net/url"
	"sync"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/FlowingSPDG/streamdeck"
)

type SDClient[SettingsT Settings] struct {
	c                 *websocket.Conn
	uuid              string
	registerEventName string
	actionInfo        inActionInfo[SettingsT]
	inInfo            inInfo
	// runningApps // ?
	isQT bool

	// Send Mutex lock
	sendMutex *sync.Mutex

	// TODO: WSから受信したメッセージからハンドラを起動する
	onDidReceiveSettingsHandler       func(streamdeck.Event)
	onDidReceiveGlobalSettingsHandler func(streamdeck.Event)
	onSendToPropertyInspectorHandler  func(streamdeck.Event)
}

// Close close client
func (sd *SDClient[SettingsT]) Close() error {
	return sd.c.Close(websocket.StatusNormalClosure, "")
}

func (sd *SDClient[SettingsT]) send(ctx context.Context, event streamdeck.Event) error {
	event.Context = sd.uuid
	sd.sendMutex.Lock()
	defer sd.sendMutex.Unlock()
	return wsjson.Write(ctx, sd.c, event)
}

// SetSettings Save data persistently for the action's instance.
func (sd *SDClient[SettingsT]) SetSettings(ctx context.Context, settings SettingsT) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.SetSettings, settings))
}

// GetSettings Request the persistent data for the action's instance.
func (sd *SDClient[SettingsT]) GetSettings(ctx context.Context) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.GetSettings, nil))
}

// SetGlobalSettings Save data securely and globally for the plugin.
func (sd *SDClient[SettingsT]) SetGlobalSettings(ctx context.Context, settings SettingsT) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.SetGlobalSettings, settings))
}

// GetGlobalSettings Request the global persistent data
func (sd *SDClient[SettingsT]) GetGlobalSettings(ctx context.Context) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.GetGlobalSettings, nil))
}

// OpenURL Open an URL in the default browser.
func (sd *SDClient[SettingsT]) OpenURL(ctx context.Context, u *url.URL) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.OpenURL, streamdeck.OpenURLPayload{URL: u.String()}))
}

// LogMessage Write a debug log to the logs file.
func (sd *SDClient[SettingsT]) LogMessage(ctx context.Context, message string) error {
	return sd.send(ctx, streamdeck.NewEvent(nil, streamdeck.LogMessage, streamdeck.LogMessagePayload{Message: message}))
}

// SendToPlugin Send a payload to the plugin.
func (sd *SDClient[SettingsT]) SendToPlugin(ctx context.Context, payload any) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.SendToPlugin, payload))
}

// Register Register PropertyInspector to StreamDeck
func (sd *SDClient[SettingsT]) Register(ctx context.Context) error {
	return sd.send(ctx, streamdeck.Event{
		Event: sd.registerEventName,
		UUID:  sd.uuid,
	})
}

func (sd *SDClient[SettingsT]) RegisterOnDidReceiveSettingsHandler(ctx context.Context, f func(streamdeck.Event)) {
	sd.onDidReceiveSettingsHandler = f
}

func (sd *SDClient[SettingsT]) RegisterOnDidReceiveGlobalSettingsHandler(ctx context.Context, f func(streamdeck.Event)) {
	sd.onDidReceiveGlobalSettingsHandler = f
}

func (sd *SDClient[SettingsT]) RegisterOnSendToPropertyInspectorHandler(ctx context.Context, f func(streamdeck.Event)) {
	sd.onSendToPropertyInspectorHandler = f
}
