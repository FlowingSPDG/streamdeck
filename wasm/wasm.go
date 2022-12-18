package wasm

// WASM: StreamDeck WebSocket Client for Property Inspector.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"syscall/js"
)

var (
	Client SDClient
)

func DeclarePropertyInspectorRegistration[S any]() {
	fmt.Println("DeclarePropertyInspectorRegistration")
	js.Global().Set("connectElgatoStreamDeckSocket", js.FuncOf(connectElgatoStreamDeckSocketJS[S]))
}

// function connectElgatoStreamDeckSocket(inPort, inPropertyInspectorUUID, inRegisterEvent, inInfo, inActionInfo)
// e.g.
// connectElgatoStreamDeckSocket(28196, "F25D3773EA4693AB3C1B4323EA6B00D1", "registerPropertyInspector", '{"application":{"font":".AppleSystemUIFont","language":"en","platform":"mac","platformVersion":"13.1.0","version":"6.0.1.17722"},"colors":{"buttonPressedBackgroundColor":"#303030FF","buttonPressedBorderColor":"#646464FF","buttonPressedTextColor":"#969696FF","disabledColor":"#007AFF7F","highlightColor":"#007AFFFF","mouseDownColor":"#2EA8FFFF"},"devicePixelRatio":2,"devices":[{"id":"7EAEBEB876DC1927A04E7E31610731CF","name":"Stream Deck","size":{"columns":5,"rows":3},"type":0}],"plugin":{"uuid":"dev.flowingspdg.newtek","version":"0.1.4"}}', '{"action":"dev.flowingspdg.newtek.shortcuttcp","context":"52ba9e6590bf53c7ff96b89d61c880b7","device":"7EAEBEB876DC1927A04E7E31610731CF","payload":{"controller":"Keypad","coordinates":{"column":3,"row":2},"settings":{"host":"192.168.100.93","shortcut":"mode","value":"2"}}}')
func connectElgatoStreamDeckSocketJS[SettingsT any](this js.Value, args []js.Value) any {
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
	js.Global().Set("std_openURL", js.FuncOf(OpenURL))

	return nil
}

func OpenURL(this js.Value, args []js.Value) any {
	u, err := url.Parse(args[0].String())
	if err != nil {
		return err
	}
	return Client.OpenURL(context.TODO(), u)
}

func connectElgatoStreamDeckSocket[SettingsT any](inPort int, inPropertyInspectorUUID string, inRegisterEvent string, inInfo inInfo, inActionInfo inActionInfo[SettingsT]) {
	fmt.Println("inPort:", inPort)
	fmt.Println("inPropertyInspectorUUID:", inPropertyInspectorUUID)
	fmt.Println("inRegisterEvent:", inRegisterEvent)
	fmt.Println("inInfo:", inInfo)
	fmt.Println("inActionInfo:", inActionInfo)

	appVersion := js.Global().Get("navigator").Get("appVersion").String()

	// TODO: websocketのクライアントを作成する
	// グローバル変数の"Client"を上書きする
	Client = &sdClient[SettingsT]{
		c:                 nil,
		uuid:              inPropertyInspectorUUID,
		registerEventName: inRegisterEvent,
		actionInfo:        inActionInfo,
		inInfo:            inInfo,
		isQT:              strings.Contains(appVersion, "QtWebEngine"),
		sendMutex:         &sync.Mutex{},
	}
}
