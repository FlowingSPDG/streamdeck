package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	sdcontext "github.com/FlowingSPDG/streamdeck/v2/context"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type actionRouter interface {
	actionUUID() string
	dispatch(ctx context.Context, ev Event) error
}

// Client communicates with the Stream Deck application over WebSocket.
type Client struct {
	params     RegistrationParams
	c          *websocket.Conn
	actions    map[string]actionRouter
	handlers   map[EventName][]EventHandler
	mu         sync.RWMutex
	sendMu     sync.Mutex
	done       chan struct{}
	closeOnce  sync.Once
	connected  atomic.Bool
	logger     *slog.Logger
	errHandler func(context.Context, error)
	dialURL    string
	queueSize  int

	pending   map[string]chan Event
	pendingMu sync.Mutex
	reqSeq    atomic.Uint64

	queues   map[string]*eventQueue
	queuesMu sync.Mutex
}

// NewClient builds a client from Stream Deck registration parameters.
func NewClient(_ context.Context, params RegistrationParams, opts ...ClientOption) *Client {
	c := &Client{
		params:    params,
		actions:   make(map[string]actionRouter),
		handlers:  make(map[EventName][]EventHandler),
		done:      make(chan struct{}),
		logger:    slog.New(slog.DiscardHandler),
		queueSize: defaultQueueSize,
		pending:   make(map[string]chan Event),
		queues:    make(map[string]*eventQueue),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// UUID returns the plugin UUID supplied by Stream Deck.
func (c *Client) UUID() string {
	return c.params.PluginUUID
}

// Params returns the registration parameters.
func (c *Client) Params() RegistrationParams {
	return c.params
}

// On registers a handler for events that are not bound to an action UUID.
func (c *Client) On(event EventName, handler EventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[event] = append(c.handlers[event], handler)
}

// RegisterNoActionHandler is an alias of On for events such as applicationDidLaunch.
func (c *Client) RegisterNoActionHandler(event EventName, handler EventHandler) {
	c.On(event, handler)
}

// OnApplicationDidLaunch registers a handler for applicationDidLaunch.
func (c *Client) OnApplicationDidLaunch(h func(context.Context, ApplicationDidLaunchPayload) error) {
	c.On(ApplicationDidLaunch, func(ctx context.Context, client *Client, event Event) error {
		p, err := event.Unmarshal[ApplicationDidLaunchPayload]()
		if err != nil {
			return err
		}
		return h(ctx, p)
	})
}

// OnApplicationDidTerminate registers a handler for applicationDidTerminate.
func (c *Client) OnApplicationDidTerminate(h func(context.Context, ApplicationDidTerminatePayload) error) {
	c.On(ApplicationDidTerminate, func(ctx context.Context, client *Client, event Event) error {
		p, err := event.Unmarshal[ApplicationDidTerminatePayload]()
		if err != nil {
			return err
		}
		return h(ctx, p)
	})
}

// OnDeviceDidConnect registers a handler for deviceDidConnect.
func (c *Client) OnDeviceDidConnect(h func(context.Context, Event) error) {
	c.On(DeviceDidConnect, func(ctx context.Context, _ *Client, event Event) error {
		return h(ctx, event)
	})
}

// OnDeviceDidDisconnect registers a handler for deviceDidDisconnect.
func (c *Client) OnDeviceDidDisconnect(h func(context.Context, Event) error) {
	c.On(DeviceDidDisconnect, func(ctx context.Context, _ *Client, event Event) error {
		return h(ctx, event)
	})
}

// OnDeviceDidChange registers a handler for deviceDidChange.
func (c *Client) OnDeviceDidChange(h func(context.Context, Event) error) {
	c.On(DeviceDidChange, func(ctx context.Context, _ *Client, event Event) error {
		return h(ctx, event)
	})
}

// OnSystemDidWakeUp registers a handler for systemDidWakeUp.
func (c *Client) OnSystemDidWakeUp(h func(context.Context) error) {
	c.On(SystemDidWakeUp, func(ctx context.Context, _ *Client, _ Event) error {
		return h(ctx)
	})
}

// OnDidReceiveDeepLink registers a handler for didReceiveDeepLink.
func (c *Client) OnDidReceiveDeepLink(h func(context.Context, DidReceiveDeepLinkPayload) error) {
	c.On(DidReceiveDeepLink, func(ctx context.Context, _ *Client, event Event) error {
		p, err := event.Unmarshal[DidReceiveDeepLinkPayload]()
		if err != nil {
			return err
		}
		return h(ctx, p)
	})
}

// OnDidReceiveGlobalSettings registers a handler for unsolicited global settings updates.
func (c *Client) OnDidReceiveGlobalSettings[G any](h func(context.Context, G) error) {
	c.On(DidReceiveGlobalSettings, func(ctx context.Context, _ *Client, event Event) error {
		p, err := event.Unmarshal[DidReceiveGlobalSettingsPayload[G]]()
		if err != nil {
			return err
		}
		return h(ctx, p.Settings)
	})
}

// OnDidReceiveSecrets registers a handler for didReceiveSecrets.
func (c *Client) OnDidReceiveSecrets[G any](h func(context.Context, G) error) {
	c.On(DidReceiveSecrets, func(ctx context.Context, _ *Client, event Event) error {
		p, err := event.Unmarshal[DidReceiveSecretsPayload[G]]()
		if err != nil {
			return err
		}
		return h(ctx, p.Secrets)
	})
}

// Run connects to Stream Deck, registers, and dispatches events until the
// connection ends or ctx is cancelled. Stream Deck sends os.Interrupt (Ctrl+C)
// when the app shuts down or the plugin is uninstalled; Run listens for that
// in addition to ctx.
func (c *Client) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, shutdownSignals()...)
	defer stop()

	addr := c.dialURL
	if addr == "" {
		addr = fmt.Sprintf("ws://127.0.0.1:%d", c.params.Port)
	}

	conn, _, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}
	c.c = conn
	c.connected.Store(true)

	readErr := make(chan error, 1)
	go func() {
		readErr <- c.readLoop(ctx)
	}()

	if err := c.register(ctx); err != nil {
		_ = c.Close()
		return fmt.Errorf("failed to register with StreamDeck: %w", err)
	}

	select {
	case err := <-readErr:
		c.connected.Store(false)
		if ctx.Err() != nil {
			return nil
		}
		return err
	case <-ctx.Done():
		_ = c.Close()
		return nil
	}
}

// IsConnected reports whether the WebSocket is open.
func (c *Client) IsConnected() bool {
	return c.connected.Load()
}

func (c *Client) register(ctx context.Context) error {
	return c.send(ctx, outgoingEvent{
		Event: EventName(c.params.RegisterEvent),
		UUID:  c.params.PluginUUID,
	})
}

func (c *Client) readLoop(ctx context.Context) error {
	defer close(c.done)
	for {
		_, message, err := c.c.Read(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("%w: %v", ErrReadFailed, err)
		}

		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			c.logger.Warn("failed to unmarshal event", "err", err, "payload", string(message))
			continue
		}

		if c.deliverPending(event) {
			continue
		}

		evCtx := sdcontext.WithContext(ctx, event.Context)
		evCtx = sdcontext.WithDevice(evCtx, event.Device)
		evCtx = sdcontext.WithAction(evCtx, event.Action)
		c.enqueue(evCtx, event)
	}
}

func (c *Client) send(ctx context.Context, event outgoingEvent) error {
	if c.c == nil {
		return ErrNotConnected
	}
	c.sendMu.Lock()
	defer c.sendMu.Unlock()
	if err := wsjson.Write(ctx, c.c, event); err != nil {
		return fmt.Errorf("%w: %v", ErrWriteFailed, err)
	}
	return nil
}

func (c *Client) sendCommand(ctx context.Context, name EventName, payload any) error {
	return c.send(ctx, outgoingEvent{
		Event:   name,
		Action:  sdcontext.Action(ctx),
		Context: firstNonEmpty(sdcontext.Context(ctx), c.params.PluginUUID),
		Device:  sdcontext.Device(ctx),
		Payload: payload,
	})
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// SetGlobalSettings saves plugin-wide settings.
func (c *Client) SetGlobalSettings[G any](ctx context.Context, settings G) error {
	return c.send(ctx, outgoingEvent{
		Event:   SetGlobalSettings,
		Context: c.params.PluginUUID,
		Payload: settings,
	})
}

// GetGlobalSettings requests plugin-wide settings and waits for the matching response.
func (c *Client) GetGlobalSettings[G any](ctx context.Context) (G, error) {
	var zero G
	resp, err := c.request(ctx, outgoingEvent{
		Event:   GetGlobalSettings,
		Context: c.params.PluginUUID,
	})
	if err != nil {
		return zero, err
	}
	p, err := resp.Unmarshal[DidReceiveGlobalSettingsPayload[G]]()
	if err != nil {
		return zero, err
	}
	return p.Settings, nil
}

// GetSecrets requests plugin secrets and waits for the matching response.
func (c *Client) GetSecrets[G any](ctx context.Context) (G, error) {
	var zero G
	resp, err := c.request(ctx, outgoingEvent{
		Event:   GetSecrets,
		Context: c.params.PluginUUID,
	})
	if err != nil {
		return zero, err
	}
	p, err := resp.Unmarshal[DidReceiveSecretsPayload[G]]()
	if err != nil {
		return zero, err
	}
	return p.Secrets, nil
}

// OpenURL opens a URL in the default browser.
func (c *Client) OpenURL(ctx context.Context, u url.URL) error {
	return c.sendCommand(ctx, OpenURL, OpenURLPayload{URL: u.String()})
}

// LogMessage writes a debug line to the Stream Deck log file.
func (c *Client) LogMessage(ctx context.Context, message string) error {
	return c.sendCommand(ctx, LogMessage, LogMessagePayload{Message: message})
}

// SetTitle changes the title of an action instance identified by ctx.
func (c *Client) SetTitle(ctx context.Context, title string, target Target, state ...int) error {
	payload := SetTitlePayload{Title: title, Target: target}
	if len(state) > 0 {
		payload.State = new(state[0])
	}
	return c.sendCommand(ctx, SetTitle, payload)
}

// SetImage changes the image of an action instance identified by ctx.
func (c *Client) SetImage(ctx context.Context, base64image string, target Target, state ...int) error {
	payload := SetImagePayload{Base64Image: base64image, Target: target}
	if len(state) > 0 {
		payload.State = new(state[0])
	}
	return c.sendCommand(ctx, SetImage, payload)
}

// SetFeedback updates Stream Deck + layout items.
func (c *Client) SetFeedback(ctx context.Context, payload Feedback) error {
	return c.sendCommand(ctx, SetFeedback, payload)
}

// SetFeedbackLayout sets the layout associated with an action instance.
func (c *Client) SetFeedbackLayout(ctx context.Context, layout string) error {
	return c.sendCommand(ctx, SetFeedbackLayout, SetFeedbackLayoutPayload{Layout: layout})
}

// SetTriggerDescription sets encoder trigger descriptions.
func (c *Client) SetTriggerDescription(ctx context.Context, payload SetTriggerDescriptionPayload) error {
	return c.sendCommand(ctx, SetTriggerDescription, payload)
}

// ShowAlert shows a temporary alert icon on an action instance.
func (c *Client) ShowAlert(ctx context.Context) error {
	return c.sendCommand(ctx, ShowAlert, nil)
}

// ShowOk shows a temporary OK icon on an action instance.
func (c *Client) ShowOk(ctx context.Context) error {
	return c.sendCommand(ctx, ShowOk, nil)
}

// SetState selects a multi-state action state.
func (c *Client) SetState(ctx context.Context, state int) error {
	return c.sendCommand(ctx, SetState, SetStatePayload{State: state})
}

// SwitchToProfile switches to a plugin-distributed profile.
func (c *Client) SwitchToProfile(ctx context.Context, profile string, page ...int) error {
	payload := SwitchProfilePayload{Profile: profile}
	if len(page) > 0 {
		payload.Page = new(page[0])
	}
	return c.send(ctx, outgoingEvent{
		Event:   SwitchToProfile,
		Context: firstNonEmpty(sdcontext.Context(ctx), c.params.PluginUUID),
		Device:  sdcontext.Device(ctx),
		Payload: payload,
	})
}

// Close closes the WebSocket. It is safe to call more than once.
func (c *Client) Close() error {
	var err error
	c.closeOnce.Do(func() {
		c.connected.Store(false)
		c.stopQueues()
		if c.c != nil {
			err = c.c.Close(websocket.StatusNormalClosure, "")
		}
		select {
		case <-c.done:
		case <-time.After(time.Second):
		}
	})
	return err
}

func (c *Client) handleError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	c.logger.Warn("event handler error", "err", err)
	if c.errHandler != nil {
		c.errHandler(ctx, err)
	}
}
