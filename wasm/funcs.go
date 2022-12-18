package wasm

import (
	"context"
	"net/url"
	"syscall/js"

	"github.com/FlowingSPDG/streamdeck"
)

// TODO: スライスの境界外アクセスを止める
// TODO: Client がnilだった場合やめる

func SetSettings(this js.Value, args []js.Value) any {
	return Client.SetSettings(context.TODO(), args[0].String())
}

func GetSettings(this js.Value, args []js.Value) any {
	return Client.GetSettings(context.TODO())
}

func SetGlobalSettings(this js.Value, args []js.Value) any {
	return Client.SetGlobalSettings(context.TODO(), args[0].String())
}

func GetGlobalSettings(this js.Value, args []js.Value) any {
	return Client.GetGlobalSettings(context.TODO())
}

func OpenURL(this js.Value, args []js.Value) any {
	u, err := url.Parse(args[0].String())
	if err != nil {
		return err
	}
	return Client.OpenURL(context.TODO(), u)
}

func LogMessage(this js.Value, args []js.Value) any {
	return Client.LogMessage(args[0].String())
}

func SetImage(this js.Value, args []js.Value) any {
	return Client.SetImage(context.TODO(), args[0].String(), streamdeck.Target(args[1].Int()))
}

func ShowAlert(this js.Value, args []js.Value) any {
	return Client.ShowAlert(context.TODO())
}
func ShowOk(this js.Value, args []js.Value) any {
	return Client.ShowOk(context.TODO())
}

func SetState(this js.Value, args []js.Value) any {
	return Client.SetState(context.TODO(), args[0].Int())
}
func SwitchToProfile(this js.Value, args []js.Value) any {
	return Client.SwitchToProfile(context.TODO(), args[0].String())
}

func SendToPropertyInspector(this js.Value, args []js.Value) any {
	return Client.SendToPropertyInspector(context.TODO(), args[0].String())
}

func SendToPlugin(this js.Value, args []js.Value) any {
	return Client.SendToPlugin(context.TODO(), args[0].String())
}
