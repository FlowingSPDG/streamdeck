package wasm

import (
	"context"
	"net/url"
	"sync"

	"github.com/FlowingSPDG/streamdeck"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type SDClient[SettingsT any] interface {
	Close() error

	SetSettings(ctx context.Context, settings SettingsT) error
	GetSettings(ctx context.Context) error
	SetGlobalSettings(ctx context.Context, settings SettingsT) error
	GetGlobalSettings(ctx context.Context) error
	OpenURL(ctx context.Context, u *url.URL) error
	LogMessage(ctx context.Context, message string) error
	SendToPlugin(ctx context.Context, payload SettingsT) error

	// TODO: Add handler for didReceiveSettings, didReceiveGlobalSettings, sendToPropertyInspector.
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

// Close close client
func (sd *sdClient[SettingsT]) Close() error {
	return sd.c.Close(websocket.StatusNormalClosure, "")
}

// TODO: WSから受信したメッセージからハンドラを起動する
func (sd *sdClient[SettingsT]) send(ctx context.Context, event streamdeck.Event) error {
	sd.sendMutex.Lock()
	defer sd.sendMutex.Unlock()
	return wsjson.Write(ctx, sd.c, event)
}

// SetSettings Save data persistently for the action's instance.
func (sd *sdClient[SettingsT]) SetSettings(ctx context.Context, settings any) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.SetSettings, settings))
}

// GetSettings Request the persistent data for the action's instance.
func (sd *sdClient[SettingsT]) GetSettings(ctx context.Context) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.GetSettings, nil))
}

// SetGlobalSettings Save data securely and globally for the plugin.
func (sd *sdClient[SettingsT]) SetGlobalSettings(ctx context.Context, settings any) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.SetGlobalSettings, settings))
}

// GetGlobalSettings Request the global persistent data
func (sd *sdClient[SettingsT]) GetGlobalSettings(ctx context.Context) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.GetGlobalSettings, nil))
}

// OpenURL Open an URL in the default browser.
func (sd *sdClient[SettingsT]) OpenURL(ctx context.Context, u *url.URL) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.OpenURL, streamdeck.OpenURLPayload{URL: u.String()}))
}

// LogMessage Write a debug log to the logs file.
func (sd *sdClient[SettingsT]) LogMessage(ctx context.Context, message string) error {
	return sd.send(ctx, streamdeck.NewEvent(nil, streamdeck.LogMessage, streamdeck.LogMessagePayload{Message: message}))
}

// SendToPlugin Send a payload to the plugin.
func (sd *sdClient[SettingsT]) SendToPlugin(ctx context.Context, payload any) error {
	return sd.send(ctx, streamdeck.NewEvent(ctx, streamdeck.SendToPlugin, payload))
}
