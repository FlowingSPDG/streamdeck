package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	sdcontext "github.com/FlowingSPDG/streamdeck/context"
	"github.com/gorilla/websocket"
)

var (
	logger = log.New(ioutil.Discard, "streamdeck", log.LstdFlags)
)

// Log Get logger
func Log() *log.Logger {
	return logger
}

// EventHandler Event handler func
type EventHandler func(ctx context.Context, client *Client, event Event) error

// Client StreamDeck communicating client
type Client struct {
	ctx       context.Context
	params    RegistrationParams
	c         *websocket.Conn
	actions   actions
	handlers  eventHandlers
	done      chan struct{}
	sendMutex sync.Mutex
}

// map[string]*Action
type actions struct {
	m sync.Map
}

// NewClient Get new client from specified context/params. you can specify "os.Args".
func NewClient(ctx context.Context, params RegistrationParams) *Client {
	return &Client{
		ctx:    ctx,
		params: params,
		c:      nil,
		actions: actions{
			m: sync.Map{},
		},
		handlers: eventHandlers{
			m: sync.Map{},
		},
		done:      make(chan struct{}),
		sendMutex: sync.Mutex{},
	}
}

// Action Get action from uuid.
func (client *Client) Action(uuid string) *Action {
	v := newAction(uuid)
	val, ok := client.actions.m.LoadOrStore(uuid, v)
	if !ok {
		v = newAction(uuid)
		client.actions.m.Store(uuid, v)
	} else {
		v = val.(*Action)
	}
	return v
}

// RegisterNoActionHandler register event handler with no action such as "applicationDidLaunch".
func (client *Client) RegisterNoActionHandler(eventName string, handler EventHandler) {
	eh := eventHandlerSlice{
		mutex: &sync.Mutex{},
		eh:    []EventHandler{},
	}

	ehi, _ := client.handlers.m.LoadOrStore(eventName, eh)
	eh = ehi.(eventHandlerSlice)

	eh.mutex.Lock()
	defer eh.mutex.Unlock()
	eh.eh = append(eh.eh, handler)

	client.handlers.m.Store(eventName, eh)
}

// Run Start communicating with StreamDeck software
func (client *Client) Run() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("127.0.0.1:%d", client.params.Port)}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	client.c = c

	go func() {
		defer close(client.done)
		for {
			messageType, message, err := client.c.ReadMessage()
			if err != nil {
				logger.Printf("read error: %v\n", err)
				return
			}

			if messageType == websocket.PingMessage {
				logger.Printf("received ping message\n")
				if err := client.c.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
					logger.Printf("error while ponging: %v\n", err)
				}
				continue
			}

			event := Event{}
			if err := json.Unmarshal(message, &event); err != nil {
				logger.Printf("failed to unmarshal received event: %s\n", string(message))
				continue
			}

			logger.Println("recv: ", string(message))

			ctx := sdcontext.WithContext(client.ctx, event.Context)
			ctx = sdcontext.WithDevice(ctx, event.Device)
			ctx = sdcontext.WithAction(ctx, event.Action)

			if event.Action == "" {
				v, ok := client.handlers.m.Load(event.Action)
				if ok {
					eh := v.(eventHandlerSlice)
					eh.Execute(ctx, client, event)
				}
				return
			}

			var action *Action
			a, ok := client.actions.m.Load(event.Action)
			if !ok {
				action = client.Action(event.Action)
				action.addContext(ctx)
			} else {
				action = a.(*Action)
			}
			v, ok := action.handlers.m.Load(event.Event)
			if ok {
				eh := v.(eventHandlerSlice)
				eh.Execute(ctx, client, event)
			}
		}
	}()

	if err := client.register(client.params); err != nil {
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

func (client *Client) register(params RegistrationParams) error {
	if err := client.send(Event{UUID: params.PluginUUID, Event: params.RegisterEvent}); err != nil {
		client.Close()
		return err
	}
	return nil
}

func (client *Client) send(event Event) error {
	j, _ := json.Marshal(event)
	client.sendMutex.Lock()
	defer client.sendMutex.Unlock()
	logger.Printf("sending message: %v\n", string(j))
	return client.c.WriteJSON(event)
}

// SetSettings Save data persistently for the action's instance.
func (client *Client) SetSettings(ctx context.Context, settings interface{}) error {
	return client.send(NewEvent(ctx, SetSettings, settings))
}

// GetSettings Request the persistent data for the action's instance.
func (client *Client) GetSettings(ctx context.Context) error {
	return client.send(NewEvent(ctx, GetSettings, nil))
}

// SetGlobalSettings Save data securely and globally for the plugin.
func (client *Client) SetGlobalSettings(ctx context.Context, settings interface{}) error {
	return client.send(NewEvent(ctx, SetGlobalSettings, settings))
}

// GetGlobalSettings Request the global persistent data
func (client *Client) GetGlobalSettings(ctx context.Context) error {
	return client.send(NewEvent(ctx, GetGlobalSettings, nil))
}

// OpenURL Open an URL in the default browser.
func (client *Client) OpenURL(ctx context.Context, u url.URL) error {
	return client.send(NewEvent(ctx, OpenURL, OpenURLPayload{URL: u.String()}))
}

// LogMessage Write a debug log to the logs file.
func (client *Client) LogMessage(message string) error {
	return client.send(NewEvent(nil, LogMessage, LogMessagePayload{Message: message}))
}

// SetTitle Dynamically change the title of an instance of an action.
func (client *Client) SetTitle(ctx context.Context, title string, target Target) error {
	return client.send(NewEvent(ctx, SetTitle, SetTitlePayload{Title: title, Target: target}))
}

// SetImage Dynamically change the image displayed by an instance of an action.
func (client *Client) SetImage(ctx context.Context, base64image string, target Target) error {
	return client.send(NewEvent(ctx, SetImage, SetImagePayload{Base64Image: base64image, Target: target}))
}

// ShowAlert Temporarily show an alert icon on the image displayed by an instance of an action.
func (client *Client) ShowAlert(ctx context.Context) error {
	return client.send(NewEvent(ctx, ShowAlert, nil))
}

// ShowOk Temporarily show an OK checkmark icon on the image displayed by an instance of an action
func (client *Client) ShowOk(ctx context.Context) error {
	return client.send(NewEvent(ctx, ShowOk, nil))
}

// SetState Change the state of the action's instance supporting multiple states.
func (client *Client) SetState(ctx context.Context, state int) error {
	return client.send(NewEvent(ctx, SetState, SetStatePayload{State: state}))
}

// SwitchToProfile Switch to one of the preconfigured read-only profiles.
func (client *Client) SwitchToProfile(ctx context.Context, profile string) error {
	return client.send(NewEvent(ctx, SwitchToProfile, SwitchProfilePayload{Profile: profile}))
}

// SendToPropertyInspector Send a payload to the Property Inspector.
func (client *Client) SendToPropertyInspector(ctx context.Context, payload interface{}) error {
	return client.send(NewEvent(ctx, SendToPropertyInspector, payload))
}

// SendToPlugin Send a payload to the plugin.
func (client *Client) SendToPlugin(ctx context.Context, payload interface{}) error {
	return client.send(NewEvent(ctx, SendToPlugin, payload))
}

// Close close client
func (client *Client) Close() error {
	err := client.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	select {
	case <-client.done:
	case <-time.After(time.Second):
	}
	return client.c.Close()
}
