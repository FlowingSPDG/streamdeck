package wasm

// WASM: StreamDeck WebSocket Client for Property Inspector.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"syscall/js"
	"time"

	"github.com/FlowingSPDG/streamdeck"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	Client SDClient
)

func DeclarePropertyInspectorRegistration[S any]() {
	// wasmを読み込む前にconnectElgatoStreamDeckSocketが走ってしまう
	// wasmロード前に保存して、ロードが終わり次第即座に起動する

	js.Global().Set("std_connected", false)

	// これならパースしてString()とかでも良いのでは？
	vals := make([]js.Value, 0, 5)
	vals = append(vals, js.Global().Get("port"))
	vals = append(vals, js.Global().Get("uuid"))
	vals = append(vals, js.Global().Get("registerEventName"))
	vals = append(vals, js.Global().Get("Info"))
	vals = append(vals, js.Global().Get("actionInfo"))

	js.Global().Set("connectElgatoStreamDeckSocketJS", js.FuncOf(connectElgatoStreamDeckSocketJS[S]))
	connectElgatoStreamDeckSocketJS[S](js.Undefined(), vals)
}

// DidReceiveSettings を受信したり、WebSocketの接続が確立した時にJSに変数を格納したい

// function connectElgatoStreamDeckSocket(inPort, inPropertyInspectorUUID, inRegisterEvent, inInfo, inActionInfo)
// e.g.
// connectElgatoStreamDeckSocket(28196, "F25D3773EA4693AB3C1B4323EA6B00D1", "registerPropertyInspector", '{"application":{"font":".AppleSystemUIFont","language":"en","platform":"mac","platformVersion":"13.1.0","version":"6.0.1.17722"},"colors":{"buttonPressedBackgroundColor":"#303030FF","buttonPressedBorderColor":"#646464FF","buttonPressedTextColor":"#969696FF","disabledColor":"#007AFF7F","highlightColor":"#007AFFFF","mouseDownColor":"#2EA8FFFF"},"devicePixelRatio":2,"devices":[{"id":"7EAEBEB876DC1927A04E7E31610731CF","name":"Stream Deck","size":{"columns":5,"rows":3},"type":0}],"plugin":{"uuid":"dev.flowingspdg.newtek","version":"0.1.4"}}', '{"action":"dev.flowingspdg.newtek.shortcuttcp","context":"52ba9e6590bf53c7ff96b89d61c880b7","device":"7EAEBEB876DC1927A04E7E31610731CF","payload":{"controller":"Keypad","coordinates":{"column":3,"row":2},"settings":{"host":"192.168.100.93","shortcut":"mode","value":"2"}}}')
func connectElgatoStreamDeckSocketJS[SettingsT any](this js.Value, args []js.Value) any {
	fmt.Println("connectElgatoStreamDeckSocket")
	inPort := args[0].Int()
	inPropertyInspectorUUID := args[1].String()
	inRegisterEvent := args[2].String() // should be "registerPropertyInspector"
	inInfo := inInfo{}
	if err := json.Unmarshal([]byte(args[3].String()), &inInfo); err != nil {
		fmt.Println("Failed to parse inInfo:", err)
		return err
	}
	inActionInfo := inActionInfo[SettingsT]{}
	if err := json.Unmarshal([]byte(args[4].String()), &inActionInfo); err != nil {
		fmt.Println("Failed to parse inInfo:", err)
		return err
	}

	connectElgatoStreamDeckSocket(inPort, inPropertyInspectorUUID, inRegisterEvent, inInfo, inActionInfo)

	// 関数を登録する
	setStreamdeckFunctions()

	return nil
}

func setStreamdeckFunctions() {
	fmt.Println("Registering Global Functions")
	js.Global().Set(streamdeck.SetSettings, js.FuncOf(SetSettings))
	js.Global().Set(streamdeck.GetSettings, js.FuncOf(GetSettings))
	js.Global().Set(streamdeck.SetGlobalSettings, js.FuncOf(SetGlobalSettings))
	js.Global().Set(streamdeck.GetGlobalSettings, js.FuncOf(GetGlobalSettings))
	js.Global().Set(streamdeck.OpenURL, js.FuncOf(OpenURL))
	js.Global().Set(streamdeck.LogMessage, js.FuncOf(LogMessage))
	js.Global().Set(streamdeck.SetImage, js.FuncOf(SetImage))
	js.Global().Set(streamdeck.ShowAlert, js.FuncOf(ShowAlert))
	js.Global().Set(streamdeck.ShowOk, js.FuncOf(ShowOk))
	js.Global().Set(streamdeck.SetState, js.FuncOf(SetState))
	js.Global().Set(streamdeck.SwitchToProfile, js.FuncOf(SwitchToProfile))
	js.Global().Set(streamdeck.SendToPropertyInspector, js.FuncOf(SendToPropertyInspector))
	js.Global().Set(streamdeck.SendToPlugin, js.FuncOf(SendToPlugin))
}

func connectElgatoStreamDeckSocket[SettingsT any](inPort int, inPropertyInspectorUUID string, inRegisterEvent string, inInfo inInfo, inActionInfo inActionInfo[SettingsT]) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	fmt.Println("inPort:", inPort)
	fmt.Println("inPropertyInspectorUUID:", inPropertyInspectorUUID)
	fmt.Println("inRegisterEvent:", inRegisterEvent)
	fmt.Println("inInfo:", inInfo)
	fmt.Println("inActionInfo:", inActionInfo)

	appVersion := js.Global().Get("navigator").Get("appVersion").String()

	fmt.Println("Try websocket")
	p := fmt.Sprintf("ws://127.0.0.1:%d", inPort)
	fmt.Printf("connecting to %s", p)

	c, _, err := websocket.Dial(ctx, p, nil)
	if err != nil {
		// TODO: handle error
		fmt.Println("Failed to connect websocket:", err.Error())
		return
	}
	js.Global().Set("std_connected", true)
	// TODO: defer to close websocket

	// イベントを登録
	if err := wsjson.Write(ctx, c, map[string]string{
		"event": inRegisterEvent,
		"uuid":  inPropertyInspectorUUID,
	}); err != nil {
		// TODO: handle error
		fmt.Println("Failed to register Property Inspector:", err.Error())
		return
	}

	Client = &sdClient[SettingsT]{
		c:                 c,
		uuid:              inPropertyInspectorUUID,
		registerEventName: inRegisterEvent,
		actionInfo:        inActionInfo,
		inInfo:            inInfo,
		isQT:              strings.Contains(appVersion, "QtWebEngine"),
		sendMutex:         &sync.Mutex{},
	}
}
