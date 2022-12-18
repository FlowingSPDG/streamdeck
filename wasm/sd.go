package wasm

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"syscall/js"

	"github.com/FlowingSPDG/streamdeck"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type SDClient[SettingsT any] interface {
	// JS関連
	RegisterGlobal(name string) // js.Global()に登録する

	// 基礎操作
	Close() error
	Register(ctx context.Context) error

	// StreamDeckとの連携
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

func (sd *sdClient[SettingsT]) RegisterGlobal(name string) {
	fmt.Println("Registering methods")
	window := js.Global()

	window.Set(name, map[string]js.Func{
		streamdeck.SetSettings:       js.FuncOf(SetSettings),
		streamdeck.GetSettings:       js.FuncOf(GetSettings),
		streamdeck.SetGlobalSettings: js.FuncOf(SetGlobalSettings),
		streamdeck.GetGlobalSettings: js.FuncOf(GetGlobalSettings),
		streamdeck.OpenURL:           js.FuncOf(OpenURL),
		streamdeck.LogMessage:        js.FuncOf(LogMessage),
		streamdeck.SendToPlugin:      js.FuncOf(SendToPlugin),
	})
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

// Register Register PropertyInspector to StreamDeck
func (sd *sdClient[SettingsT]) Register(ctx context.Context) error {
	return sd.send(ctx, streamdeck.Event{
		Event: sd.registerEventName,
		UUID:  sd.uuid,
	})
}
