package wasm

import (
	"context"
	"net/url"
	"sync"

	"github.com/FlowingSPDG/streamdeck"
	"github.com/gorilla/websocket"
)

type SDClient interface {
	// TODO
	OpenURL(ctx context.Context, u *url.URL) error
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

func (sd *sdClient[SettingsT]) OpenURL(ctx context.Context, u *url.URL) error {
	return sd.send(streamdeck.NewEvent(ctx, streamdeck.OpenURL, streamdeck.OpenURLPayload{URL: u.String()}))
}

func (sd *sdClient[SettingsT]) send(event streamdeck.Event) error {
	sd.sendMutex.Lock()
	defer sd.sendMutex.Unlock()
	return sd.c.WriteJSON(event)
}
