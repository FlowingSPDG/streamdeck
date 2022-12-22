package wasm

import (
	"context"
	"fmt"
	"net/url"
	"syscall/js"

	"github.com/FlowingSPDG/streamdeck"
)

// TODO: スライスの境界外アクセスを止める
// TODO: Client がnilだった場合やめる

// sdClientJSのJS wrapper
type sdClientJS[SettingsT Settings] struct {
	c *SDClient[SettingsT]
	s SettingsT
}

func newSdClientJS[SettingsT Settings](c *SDClient[SettingsT], s SettingsT) *sdClientJS[SettingsT] {
	return &sdClientJS[SettingsT]{c: c, s: s}
}

func (sdj *sdClientJS[SettingsT]) RegisterGlobal(name string) {
	fmt.Println("Registering methods")
	window := js.Global()

	window.Set(name, map[string]any{
		streamdeck.SetSettings:       js.FuncOf(sdj.SetSettings),
		streamdeck.GetSettings:       js.FuncOf(sdj.GetSettings),
		streamdeck.SetGlobalSettings: js.FuncOf(sdj.SetGlobalSettings),
		streamdeck.GetGlobalSettings: js.FuncOf(sdj.GetGlobalSettings),
		streamdeck.OpenURL:           js.FuncOf(sdj.OpenURL),
		streamdeck.LogMessage:        js.FuncOf(sdj.LogMessage),
		streamdeck.SendToPlugin:      js.FuncOf(sdj.SendToPlugin),
	})
}

func (sdj *sdClientJS[SettingsT]) SetSettings(this js.Value, args []js.Value) any {
	ctx := context.Background()

	fmt.Println("SetSettings:", sdj.s)

	sdj.c.SetSettings(ctx, sdj.s)
	return nil
}

func (sdj *sdClientJS[SettingsT]) GetSettings(this js.Value, args []js.Value) any {
	return sdj.c.GetSettings(context.TODO())
}

func (sdj *sdClientJS[SettingsT]) SetGlobalSettings(this js.Value, args []js.Value) any {
	// TODO
	return nil
}

func (sdj *sdClientJS[SettingsT]) GetGlobalSettings(this js.Value, args []js.Value) any {
	return sdj.c.GetGlobalSettings(context.TODO())
}

func (sdj *sdClientJS[SettingsT]) OpenURL(this js.Value, args []js.Value) any {
	u, err := url.Parse(args[0].String())
	if err != nil {
		return err
	}
	return sdj.c.OpenURL(context.TODO(), u)
}

func (sdj *sdClientJS[SettingsT]) LogMessage(this js.Value, args []js.Value) any {
	return sdj.c.LogMessage(context.TODO(), args[0].String())
}

func (sdj *sdClientJS[SettingsT]) SendToPlugin(this js.Value, args []js.Value) any {
	return sdj.c.SendToPlugin(context.TODO(), args[0].String())
}
