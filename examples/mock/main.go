package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/FlowingSPDG/streamdeck"
	"github.com/c-bata/go-prompt"
	"github.com/olahol/melody"
)

var (
	exit = "exit"

	action = ""
)

// Mock用ソフトウェア
// StreamDeckソフトウェア本体の挙動を再現するモックアップソフトウェア

// Mocking software

// mockServer with websocket server
type mockServer struct {
	s *melody.Melody // Websocket sessions
}

func (m *mockServer) SendEvent(event streamdeck.Event) error {
	event.Action = action
	event.Context = "ff158f4368c8b128c7ae45a99c88baa8"
	j, _ := json.Marshal(event)
	log.Printf("Sending event: %s\n", j)
	return m.s.Broadcast(j)
}

func (m *mockServer) KeyDown(ctx context.Context, payload any) error {
	ev := streamdeck.NewEvent(ctx, streamdeck.KeyDown, payload)
	return m.SendEvent(ev)
}

func (m *mockServer) WillAppear(ctx context.Context, payload any) error {
	ev := streamdeck.NewEvent(ctx, streamdeck.WillAppear, payload)
	return m.SendEvent(ev)
}

func (m *mockServer) WillDisappear(ctx context.Context, payload any) error {
	ev := streamdeck.NewEvent(ctx, streamdeck.WillDisappear, payload)
	return m.SendEvent(ev)
}

func (m *mockServer) KeyUp(ctx context.Context, payload any) error {
	ev := streamdeck.NewEvent(ctx, streamdeck.KeyUp, payload)
	return m.SendEvent(ev)
}

func (m *mockServer) SendToPlugin(ctx context.Context, payload any) error {
	ev := streamdeck.NewEvent(ctx, streamdeck.SendToPlugin, payload)
	return m.SendEvent(ev)
}

func newMock() *mockServer {
	m := melody.New()
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		log.Printf("RCV: %s\n", msg)
	})
	return &mockServer{s: m}
}

// TODO: 各種Actionのmock挙動を網羅する
// https://developer.elgato.com/documentation/stream-deck/sdk/events-received/

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: streamdeck.KeyDown, Description: "Send KeyDown event"},
		{Text: streamdeck.KeyUp, Description: "Send KeyUp event"},
		{Text: streamdeck.SendToPlugin, Description: "Send SendToPlugin event"},
		{Text: streamdeck.WillAppear, Description: "Send WillAppear event"},
		{Text: streamdeck.WillDisappear, Description: "Send WillDisappear event"},
		{Text: exit, Description: "Shutdown mock server"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	flag.StringVar(&action, "action", "dev.samwho.streamdeck.cpu", "Action ID")
	flag.Parse()

	ctx := context.Background()
	m := newMock()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m.s.HandleRequest(w, r)
	})

	p := 5000
	log.Println("STARTING on", p)
	line := fmt.Sprintf(`-port %d -pluginUUID C4B10610C736B410125C84369621E33E -registerEvent registerPlugin -info {"application":{"font":"Segoe UI","language":"ja","platform":"windows","platformVersion":"10.0.19043","version":"6.0.1.17722"},"colors":{"buttonMouseOverBackgroundColor":"#464646FF","buttonPressedBackgroundColor":"#303030FF","buttonPressedBorderColor":"#646464FF","buttonPressedTextColor":"#969696FF","highlightColor":"#0078FFFF"},"devicePixelRatio":1,"devices":[{"id":"7EAEBEB876DC1927A04E7E31610731CF","name":"StreamDeck1","size":{"columns":5,"rows":3},"type":0},{"id":"3B1C10164F4A9B3850E88CDA0324656D","name":"Stream Deck XL","size":{"columns":8,"rows":4},"type":2},{"id":"934035044D1CF9F56EBA9607286EFE53","name":"Stream Deck","size":{"columns":5,"rows":3},"type":0},{"id":"D76B6C8774E20D878B50E12759E3CFAA","name":"Stream Deck Mini","size":{"columns":3,"rows":2},"type":1}],"plugin":{"uuid":"dev.samwho.cpu","version":"0.1"}} `, p)
	log.Printf("execute your plugin with following command line: \n%s\n", line)

	// Listen http on localhost:5000
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", p), nil); err != nil {
			panic(err)
		}
	}()

	// Handle prompt
	for {
		t := prompt.Input("StreamDeck Action > ", completer)
		switch t {
		case exit:
			return
		case streamdeck.KeyDown:
			m.KeyDown(ctx, nil)
		case streamdeck.KeyUp:
			m.KeyUp(ctx, nil)
		case streamdeck.SendToPlugin:
			m.SendToPlugin(ctx, nil)
		case streamdeck.WillAppear:
			m.WillAppear(ctx, nil)
		case streamdeck.WillDisappear:
			m.WillDisappear(ctx, nil)
		default:
			m.SendEvent(streamdeck.NewEvent(ctx, t, nil))
		}
	}

}
