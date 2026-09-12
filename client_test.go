package streamdeck_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/FlowingSPDG/streamdeck/v2"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type fakeDeck struct {
	url   string
	mu    sync.Mutex
	inbox []map[string]any
	send  chan map[string]any
	conn  *websocket.Conn
}

func newFakeDeck(t *testing.T) *fakeDeck {
	t.Helper()
	f := &fakeDeck{send: make(chan map[string]any, 16)}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			return
		}
		f.mu.Lock()
		f.conn = c
		f.mu.Unlock()
		defer c.Close(websocket.StatusNormalClosure, "")

		ctx := r.Context()
		go func() {
			for {
				var msg map[string]any
				if err := wsjson.Read(ctx, c, &msg); err != nil {
					return
				}
				f.mu.Lock()
				f.inbox = append(f.inbox, msg)
				f.mu.Unlock()
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-f.send:
				if err := wsjson.Write(ctx, c, msg); err != nil {
					return
				}
			}
		}
	}))
	t.Cleanup(srv.Close)
	f.url = "ws" + strings.TrimPrefix(srv.URL, "http")
	return f
}

func (f *fakeDeck) closeConn() {
	f.mu.Lock()
	c := f.conn
	f.mu.Unlock()
	if c != nil {
		_ = c.Close(websocket.StatusGoingAway, "bye")
	}
}

func (f *fakeDeck) waitInbox(t *testing.T, n int) []map[string]any {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		f.mu.Lock()
		if len(f.inbox) >= n {
			out := append([]map[string]any(nil), f.inbox...)
			f.mu.Unlock()
			return out
		}
		f.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	t.Fatalf("inbox size = %d, want >= %d: %+v", len(f.inbox), n, f.inbox)
	return nil
}

func testParams() streamdeck.RegistrationParams {
	return streamdeck.RegistrationParams{
		Port:          0,
		PluginUUID:    "plugin-uuid",
		RegisterEvent: "registerPlugin",
	}
}

func TestRegisterAndKeyDown(t *testing.T) {
	f := newFakeDeck(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := streamdeck.NewClient(ctx, testParams(), streamdeck.WithWebSocketURL(f.url))
	action := streamdeck.NewAction[map[string]int](client, "com.elgato.example.action")

	got := make(chan streamdeck.KeyDownEvent[map[string]int], 1)
	action.OnKeyDown(func(ctx context.Context, e streamdeck.KeyDownEvent[map[string]int]) error {
		got <- e
		return client.SetTitle(ctx, "ok", streamdeck.HardwareAndSoftware)
	})

	errc := make(chan error, 1)
	go func() { errc <- client.Run(ctx) }()

	inbox := f.waitInbox(t, 1)
	if inbox[0]["event"] != "registerPlugin" {
		t.Fatalf("first event = %+v", inbox[0])
	}
	if inbox[0]["uuid"] != "plugin-uuid" {
		t.Fatalf("uuid = %v", inbox[0]["uuid"])
	}

	var ev map[string]any
	if err := json.Unmarshal(mustReadFile(t, "testdata/events/keyDown.json"), &ev); err != nil {
		t.Fatal(err)
	}
	f.send <- ev

	select {
	case e := <-got:
		if e.Payload.Settings["counter"] != 2 {
			t.Fatalf("payload = %+v", e.Payload)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for keyDown")
	}

	inbox = f.waitInbox(t, 2)
	if inbox[1]["event"] != "setTitle" {
		t.Fatalf("second event = %+v", inbox[1])
	}

	cancel()
	select {
	case <-errc:
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return")
	}
}

func TestNoActionHandler(t *testing.T) {
	f := newFakeDeck(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := streamdeck.NewClient(ctx, testParams(), streamdeck.WithWebSocketURL(f.url))
	got := make(chan string, 1)
	client.OnApplicationDidLaunch(func(ctx context.Context, p streamdeck.ApplicationDidLaunchPayload) error {
		got <- p.Application
		return nil
	})

	go func() { _ = client.Run(ctx) }()
	f.waitInbox(t, 1)
	f.send <- map[string]any{
		"event":   "applicationDidLaunch",
		"payload": map[string]any{"application": "notepad.exe"},
	}

	select {
	case app := <-got:
		if app != "notepad.exe" {
			t.Fatalf("app = %q", app)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
}

func TestUnknownActionIsDropped(t *testing.T) {
	f := newFakeDeck(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := streamdeck.NewClient(ctx, testParams(), streamdeck.WithWebSocketURL(f.url))
	go func() { _ = client.Run(ctx) }()
	f.waitInbox(t, 1)
	f.send <- map[string]any{
		"action":  "com.unknown",
		"event":   "keyDown",
		"context": "x",
		"payload": map[string]any{},
	}
	time.Sleep(50 * time.Millisecond)
	if !client.IsConnected() {
		t.Fatal("expected connected")
	}
}

func TestGetSettingsCorrelation(t *testing.T) {
	f := newFakeDeck(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := streamdeck.NewClient(ctx, testParams(), streamdeck.WithWebSocketURL(f.url))
	action := streamdeck.NewAction[map[string]int](client, "com.elgato.example.action")

	go func() { _ = client.Run(ctx) }()
	f.waitInbox(t, 1)

	res := make(chan map[string]int, 1)
	errc := make(chan error, 1)
	go func() {
		s, err := action.GetSettings(ctx)
		if err != nil {
			errc <- err
			return
		}
		res <- s
	}()

	inbox := f.waitInbox(t, 2)
	id, _ := inbox[1]["id"].(string)
	if inbox[1]["event"] != "getSettings" || id == "" {
		t.Fatalf("getSettings = %+v", inbox[1])
	}

	var reply map[string]any
	if err := json.Unmarshal(mustReadFile(t, "testdata/events/didReceiveSettings.json"), &reply); err != nil {
		t.Fatal(err)
	}
	reply["id"] = id
	f.send <- reply

	select {
	case s := <-res:
		if s["counter"] != 4 {
			t.Fatalf("settings = %+v", s)
		}
	case err := <-errc:
		t.Fatal(err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestCloseIsIdempotent(t *testing.T) {
	client := streamdeck.NewClient(context.Background(), testParams())
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
	if client.IsConnected() {
		t.Fatal("expected not connected")
	}
}

func TestRunReturnsOnDisconnect(t *testing.T) {
	f := newFakeDeck(t)
	ctx := context.Background()
	client := streamdeck.NewClient(ctx, testParams(), streamdeck.WithWebSocketURL(f.url))
	errc := make(chan error, 1)
	go func() { errc <- client.Run(ctx) }()
	f.waitInbox(t, 1)
	f.closeConn()
	select {
	case err := <-errc:
		if err == nil {
			t.Fatal("expected read error after disconnect")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not return after disconnect")
	}
}

func TestPropertyInspectorMessage(t *testing.T) {
	f := newFakeDeck(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type msg struct {
		Action string `json:"action"`
	}

	client := streamdeck.NewClient(ctx, testParams(), streamdeck.WithWebSocketURL(f.url))
	action := streamdeck.NewAction[struct{}](client, "com.elgato.example.action")
	got := make(chan string, 1)
	action.OnPropertyInspectorMessage(func(ctx context.Context, m msg) error {
		got <- m.Action
		return action.SendToPropertyInspector(ctx, msg{Action: "resetComplete"})
	})

	go func() { _ = client.Run(ctx) }()
	f.waitInbox(t, 1)

	var ev map[string]any
	if err := json.Unmarshal(mustReadFile(t, "testdata/events/sendToPlugin.json"), &ev); err != nil {
		t.Fatal(err)
	}
	f.send <- ev

	select {
	case actionName := <-got:
		if actionName != "reset" {
			t.Fatalf("action = %q", actionName)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
	inbox := f.waitInbox(t, 2)
	if inbox[1]["event"] != "sendToPropertyInspector" {
		t.Fatalf("reply = %+v", inbox[1])
	}
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
