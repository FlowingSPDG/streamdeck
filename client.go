package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	sdcontext "github.com/FlowingSPDG/streamdeck/context"
	"github.com/puzpuzpuz/xsync/v3"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	logger = log.New(io.Discard, "streamdeck", log.LstdFlags)
)

// Log Get logger
func Log() *log.Logger {
	return logger
}

// EventHandler Event handler func
type EventHandler func(ctx context.Context, client *Client, event Event) error

// Client StreamDeck communicating client
type Client struct {
	params    RegistrationParams
	c         *websocket.Conn
	actions   *actions
	handlers  *eventHandlers
	done      chan struct{}
	sendMutex *sync.Mutex
}

type actions struct {
	m *xsync.MapOf[string, *Action]
}

// NewClient Get new client from specified context/params. you can specify "os.Args".
func NewClient(ctx context.Context, params RegistrationParams) *Client {
	return &Client{
		params: params,
		c:      nil,
		actions: &actions{
			m: xsync.NewMapOf[string, *Action](),
		},
		handlers: &eventHandlers{
			m: xsync.NewMapOf[string, *eventHandlerSlice](),
		},
		done:      make(chan struct{}),
		sendMutex: &sync.Mutex{},
	}
}

// UUID get plugin UUID
func (client *Client) UUID() string {
	return client.params.PluginUUID
}

// Action Get action from uuid.
func (client *Client) Action(uuid string) *Action {
	v := newAction(uuid)
	ok := false

	v, ok = client.actions.m.LoadOrStore(uuid, v)
	if !ok {
		v = newAction(uuid)
		client.actions.m.Store(uuid, v)
	}
	return v
}

// RegisterNoActionHandler register event handler with no action such as "applicationDidLaunch".
func (client *Client) RegisterNoActionHandler(eventName string, handler EventHandler) {
	eh, _ := client.handlers.m.LoadOrStore(eventName, &eventHandlerSlice{
		mutex: &sync.Mutex{},
		eh:    []EventHandler{},
	})

	eh.mutex.Lock()
	defer eh.mutex.Unlock()
	eh.eh = append(eh.eh, handler)

	client.handlers.m.Store(eventName, eh)
}

// Run Start communicating with StreamDeck software
func (client *Client) Run(ctx context.Context) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("127.0.0.1:%d", client.params.Port)}
	c, _, err := websocket.Dial(ctx, u.String(), nil)
	if err != nil {
		return err
	}

	client.c = c

	go func() {
		defer close(client.done)
		for {
			_, message, err := client.c.Read(ctx)
			if err != nil {
				logger.Printf("read error: %v\n", err)
				return
			}

			event := Event{}
			if err := json.Unmarshal(message, &event); err != nil {
				logger.Printf("failed to unmarshal received event: %s\n", string(message))
				continue
			}

			logger.Println("recv: ", string(message))

			ctx := sdcontext.WithContext(ctx, event.Context)
			ctx = sdcontext.WithDevice(ctx, event.Device)
			ctx = sdcontext.WithAction(ctx, event.Action)

			if event.Action == "" {
				eh, ok := client.handlers.m.Load(event.Event)
				if ok {
					eh.Execute(ctx, client, event)
				}
				continue
			}

			action, ok := client.actions.m.Load(event.Action)
			if !ok {
				action = client.Action(event.Action)
				action.addContext(ctx)
			}

			eh, ok := action.handlers.m.Load(event.Event)
			if ok {
				eh.Execute(ctx, client, event)
			}
		}
	}()

	if err := client.register(ctx, client.params); err != nil {
		return err
	}

	select {
	case <-client.done:
		return nil
	case <-interrupt:
		logger.Printf("interrupted, closing...\n")
		return client.Close()
	}
}

// Check if WebSocket connection is non-nil.
func (client *Client) IsConnected() bool {
	return client.c != nil
}

func (client *Client) register(ctx context.Context, params RegistrationParams) error {
	if err := client.send(ctx, Event{UUID: params.PluginUUID, Event: params.RegisterEvent}); err != nil {
		client.Close()
		return err
	}
	return nil
}

func (client *Client) send(ctx context.Context, event Event) error {
	client.sendMutex.Lock()
	defer client.sendMutex.Unlock()
	return wsjson.Write(ctx, client.c, event)
}

// SetSettings Save data persistently for the action's instance.
func (client *Client) SetSettings(ctx context.Context, settings any) error {
	return client.send(ctx, NewEvent(ctx, SetSettings, SetSettingsPayload[any]{
		Context:  sdcontext.Context(ctx),
		Settings: settings,
	}))
}

// GetSettings Request the persistent data for the action's instance.
func (client *Client) GetSettings(ctx context.Context) error {
	return client.send(ctx, NewEvent(ctx, GetSettings, GetSettingsPayload{
		Context: sdcontext.Context(ctx),
	}))
}

// SetGlobalSettings Save data securely and globally for the plugin.
func (client *Client) SetGlobalSettings(ctx context.Context, settings any) error {
	return client.send(ctx, NewEvent(ctx, SetGlobalSettings, SetGlobalSettingsPayload[any]{
		Settings: settings,
	}))
}

// GetGlobalSettings Request the global persistent data
func (client *Client) GetGlobalSettings(ctx context.Context) error {
	return client.send(ctx, NewEvent(ctx, GetGlobalSettings, GetGlobalSettingsPayload{}))
}

// OpenURL Open an URL in the default browser.
func (client *Client) OpenURL(ctx context.Context, u url.URL) error {
	return client.send(ctx, NewEvent(ctx, OpenURL, OpenURLPayload{URL: u.String()}))
}

// LogMessage Write a debug log to the logs file.
func (client *Client) LogMessage(ctx context.Context, message string) error {
	return client.send(ctx, NewEvent(ctx, LogMessage, LogMessagePayload{Message: message}))
}

// SetTitle Dynamically change the title of an instance of an action.
func (client *Client) SetTitle(ctx context.Context, title string, target Target) error {
	return client.send(ctx, NewEvent(ctx, SetTitle, SetTitlePayload{Title: title, Target: target}))
}

// SetImage Dynamically change the image displayed by an instance of an action.
func (client *Client) SetImage(ctx context.Context, base64image string, target Target) error {
	return client.send(ctx, NewEvent(ctx, SetImage, SetImagePayload{Base64Image: base64image, Target: target}))
}

// SetFeedback The plugin can send a setFeedback event to the Stream Deck application to dynamically change properties of items on the Stream Deck + touch display layout.
func (client *Client) SetFeedback(ctx context.Context, payload any) error {
	return client.send(ctx, NewEvent(ctx, SetFeedback, payload))
}

// SetFeedbackLayout
func (client *Client) SetFeedbackLayout(ctx context.Context, layout string) error {
	return client.send(ctx, NewEvent(ctx, SetImage, SetFeedbackLayoutPayload{Layout: layout}))
}

// ShowAlert Temporarily show an alert icon on the image displayed by an instance of an action.
func (client *Client) ShowAlert(ctx context.Context) error {
	return client.send(ctx, NewEvent(ctx, ShowAlert, nil))
}

// ShowOk Temporarily show an OK checkmark icon on the image displayed by an instance of an action
func (client *Client) ShowOk(ctx context.Context) error {
	return client.send(ctx, NewEvent(ctx, ShowOk, nil))
}

// SetState Change the state of the action's instance supporting multiple states.
func (client *Client) SetState(ctx context.Context, state int) error {
	return client.send(ctx, NewEvent(ctx, SetState, SetStatePayload{State: state}))
}

// SwitchToProfile Switch to one of the preconfigured read-only profiles.
func (client *Client) SwitchToProfile(ctx context.Context, profile string) error {
	return client.send(ctx, NewEvent(ctx, SwitchToProfile, SwitchProfilePayload{Profile: profile}))
}

// SendToPropertyInspector Send a payload to the Property Inspector.
func (client *Client) SendToPropertyInspector(ctx context.Context, payload any) error {
	return client.send(ctx, NewEvent(ctx, SendToPropertyInspector, SendToPropertyInspectorPayload[any]{
		Context: sdcontext.Context(ctx),
		Payload: payload,
	}))
}

// SendToPlugin Send a payload to the plugin.
func (client *Client) SendToPlugin(ctx context.Context, action string, payload any) error {
	return client.send(ctx, NewEvent(ctx, SendToPlugin, SendToPluginPayload[any]{
		Context: sdcontext.Context(ctx),
		Action:  action,
		Payload: payload,
	}))
}

// Close close client
func (client *Client) Close() error {
	err := client.c.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		return err
	}
	select {
	case <-client.done:
	case <-time.After(time.Second):
	}
	return nil
}
