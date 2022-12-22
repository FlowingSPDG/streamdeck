package wasm

// WASM: StreamDeck WebSocket Client for Property Inspector.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"syscall/js"

	"github.com/FlowingSPDG/streamdeck"
	"nhooyr.io/websocket"
)

func InitializePropertyInspector[S Settings](ctx context.Context, s S) (*SDClient[S], error) {
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
		fmt.Println("Failed to parse actionInfo:", err)
		return nil, err
	}

	SD, err := connectElgatoStreamDeckSocket(ctx, inPort, inPropertyInspectorUUID, inRegisterEvent, inInfo, inActionInfo, s)
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
func connectElgatoStreamDeckSocket[SettingsT Settings](ctx context.Context, inPort int, inPropertyInspectorUUID string, inRegisterEvent string, inInfo inInfo, inActionInfo inActionInfo[SettingsT], s SettingsT) (*SDClient[SettingsT], error) {
	fmt.Printf("inPort:%#v\n", inPort)
	fmt.Printf("inPropertyInspectorUUID:%#v\n", inPropertyInspectorUUID)
	fmt.Printf("inRegisterEvent:%#v\n", inRegisterEvent)
	fmt.Printf("inInfo:%#v\n", inInfo)
	fmt.Printf("inActionInfo:%#v\n", inActionInfo)

	fmt.Println("Context:", inActionInfo.Context)

	appVersion := js.Global().Get("navigator").Get("appVersion").String()

	p := fmt.Sprintf("ws://127.0.0.1:%d", inPort)
	fmt.Printf("connecting to %s", p)

	c, _, err := websocket.Dial(ctx, p, nil)
	if err != nil {
		// TODO: handle error
		fmt.Println("Failed to connect websocket:", err.Error())
		return nil, err
	}
	// TODO: defer to close websocket

	sdc := &SDClient[SettingsT]{
		c:                                 c,
		uuid:                              inPropertyInspectorUUID,
		registerEventName:                 inRegisterEvent,
		actionInfo:                        inActionInfo,
		inInfo:                            inInfo,
		isQT:                              strings.Contains(appVersion, "QtWebEngine"),
		sendMutex:                         &sync.Mutex{},
		onDidReceiveSettingsHandler:       func(streamdeck.Event) {},
		onDidReceiveGlobalSettingsHandler: func(streamdeck.Event) {},
		onSendToPropertyInspectorHandler:  func(streamdeck.Event) {},
	}
	wrapper := newSdClientJS(sdc, s)

	if err := sdc.Register(ctx); err != nil {
		// TODO: handle error
		fmt.Println("Failed to register Property Inspector:", err.Error())
		return nil, err
	}

	// window.$SD に設定するとJavaScriptからも利用が可能になる
	wrapper.RegisterGlobal("$SD")

	// HTMLに受信したSettingsを反映する
	inActionInfo.Payload.Settings.ApplyHTML()

	go func() {
		for {
			_, message, err := c.Read(context.TODO())
			if err != nil {
				fmt.Println("Failed to read message from websocket:", err)
				return
			}
			event := streamdeck.Event{}
			if err := json.Unmarshal(message, &event); err != nil {
				fmt.Printf("failed to unmarshal received event: %s\n", string(message))
				continue
			}
			fmt.Printf("RCV%#v\n", event)
			switch event.Event {
			case streamdeck.DidReceiveSettings:
				go sdc.onDidReceiveSettingsHandler(event)
			case streamdeck.DidReceiveGlobalSettings:
				go sdc.onDidReceiveGlobalSettingsHandler(event)
			case streamdeck.SendToPropertyInspector:
				go sdc.onSendToPropertyInspectorHandler(event)
			}
		}
	}()

	return sdc, nil
}
