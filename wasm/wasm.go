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

	"nhooyr.io/websocket"
)

func InitializePropertyInspector[S any]() (SDClient[S], error) {
	js.Global().Set("std_connected", false)

	// wasmを読み込む前にconnectElgatoStreamDeckSocketが走ってしまうので、
	// wasmロード前に受け取った値をグローバルに保存して、ロードが終わり次第wasm側から起動する
	inPort := js.Global().Get("port").Int()
	inPropertyInspectorUUID := js.Global().Get("uuid").String()
	inRegisterEvent := js.Global().Get("registerEventName").String() // should be "registerPropertyInspector"
	inInfo := inInfo{}
	if err := json.Unmarshal([]byte(js.Global().Get("Info").String()), &inInfo); err != nil {
		fmt.Println("Failed to parse inInfo:", err)
		return nil, err
	}
	inActionInfo := inActionInfo[S]{}
	if err := json.Unmarshal([]byte(js.Global().Get("actionInfo").String()), &inActionInfo); err != nil {
		fmt.Println("Failed to parse inInfo:", err)
		return nil, err
	}
	SD, err := connectElgatoStreamDeckSocket(inPort, inPropertyInspectorUUID, inRegisterEvent, inInfo, inActionInfo)
	if err != nil {
		fmt.Println("Failed to connect ElgatoStreamDeckSocket:", err)
		return nil, err
	}
	return SD, nil
}

// DidReceiveSettings を受信したり、WebSocketの接続が確立した時にJSに変数を格納したい

// function connectElgatoStreamDeckSocket(inPort, inPropertyInspectorUUID, inRegisterEvent, inInfo, inActionInfo)
// e.g.
// connectElgatoStreamDeckSocket(28196, "F25D3773EA4693AB3C1B4323EA6B00D1", "registerPropertyInspector", '{"application":{"font":".AppleSystemUIFont","language":"en","platform":"mac","platformVersion":"13.1.0","version":"6.0.1.17722"},"colors":{"buttonPressedBackgroundColor":"#303030FF","buttonPressedBorderColor":"#646464FF","buttonPressedTextColor":"#969696FF","disabledColor":"#007AFF7F","highlightColor":"#007AFFFF","mouseDownColor":"#2EA8FFFF"},"devicePixelRatio":2,"devices":[{"id":"7EAEBEB876DC1927A04E7E31610731CF","name":"Stream Deck","size":{"columns":5,"rows":3},"type":0}],"plugin":{"uuid":"dev.flowingspdg.newtek","version":"0.1.4"}}', '{"action":"dev.flowingspdg.newtek.shortcuttcp","context":"52ba9e6590bf53c7ff96b89d61c880b7","device":"7EAEBEB876DC1927A04E7E31610731CF","payload":{"controller":"Keypad","coordinates":{"column":3,"row":2},"settings":{"host":"192.168.100.93","shortcut":"mode","value":"2"}}}')
func connectElgatoStreamDeckSocket[SettingsT any](inPort int, inPropertyInspectorUUID string, inRegisterEvent string, inInfo inInfo, inActionInfo inActionInfo[SettingsT]) (SDClient[SettingsT], error) {
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
		return nil, err
	}
	// TODO: defer to close websocket

	sdc := &sdClient[SettingsT]{
		c:                 c,
		uuid:              inPropertyInspectorUUID,
		registerEventName: inRegisterEvent,
		actionInfo:        inActionInfo,
		inInfo:            inInfo,
		isQT:              strings.Contains(appVersion, "QtWebEngine"),
		sendMutex:         &sync.Mutex{},
	}
	wrapper := newSdClientJS(sdc)

	if err := sdc.Register(ctx); err != nil {
		// TODO: handle error
		fmt.Println("Failed to register Property Inspector:", err.Error())
		return nil, err
	}

	// window.$SD に設定するとJavaScriptからも利用が可能になる
	wrapper.RegisterGlobal("$SD")
	js.Global().Set("std_connected", true)
	return sdc, nil
}
