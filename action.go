package streamdeck

import (
	"context"
	"fmt"
	"sync"

	sdcontext "github.com/FlowingSPDG/streamdeck/v2/context"
)

// Action is a typed handle for a plugin action UUID and its settings type S.
type Action[S any] struct {
	uuid     string
	client   *Client
	handlers map[EventName][]func(context.Context, Event) error
	contexts map[string]context.Context
	mu       sync.RWMutex
}

// NewAction registers or returns the action for uuid.
// It panics if uuid is already registered with a different settings type.
func NewAction[S any](c *Client, uuid string) *Action[S] {
	c.mu.Lock()
	defer c.mu.Unlock()
	if r, ok := c.actions[uuid]; ok {
		if a, ok := r.(*Action[S]); ok {
			return a
		}
		panic("streamdeck: action " + uuid + " already registered with a different settings type")
	}
	a := &Action[S]{
		uuid:     uuid,
		client:   c,
		handlers: make(map[EventName][]func(context.Context, Event) error),
		contexts: make(map[string]context.Context),
	}
	a.handlers[WillAppear] = []func(context.Context, Event) error{
		func(ctx context.Context, _ Event) error {
			a.addContext(ctx)
			return nil
		},
	}
	a.handlers[WillDisappear] = []func(context.Context, Event) error{
		func(ctx context.Context, _ Event) error {
			a.removeContext(ctx)
			return nil
		},
	}
	c.actions[uuid] = a
	return a
}

func (a *Action[S]) actionUUID() string { return a.uuid }

func (a *Action[S]) dispatch(ctx context.Context, ev Event) error {
	a.mu.RLock()
	hs := append([]func(context.Context, Event) error(nil), a.handlers[ev.Event]...)
	a.mu.RUnlock()

	var lastErr error
	for _, handler := range hs {
		if err := handler(ctx, ev); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// UUID returns the action UUID.
func (a *Action[S]) UUID() string { return a.uuid }

// Client returns the parent client.
func (a *Action[S]) Client() *Client { return a.client }

// RegisterHandler appends a raw handler for event. Handlers are snapshotted before execution.
func (a *Action[S]) RegisterHandler(event EventName, handler func(context.Context, Event) error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.handlers[event] = append(a.handlers[event], handler)
}

func (a *Action[S]) onPayload[P any](name EventName, h func(context.Context, ActionEvent[P]) error) {
	a.RegisterHandler(name, func(ctx context.Context, ev Event) error {
		p, err := ev.Unmarshal[P]()
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		return h(ctx, actionEvent(ev, p))
	})
}

// OnWillAppear registers a willAppear handler.
func (a *Action[S]) OnWillAppear(h func(context.Context, WillAppearEvent[S]) error) {
	a.onPayload(WillAppear, h)
}

// OnWillDisappear registers a willDisappear handler.
func (a *Action[S]) OnWillDisappear(h func(context.Context, WillDisappearEvent[S]) error) {
	a.onPayload(WillDisappear, h)
}

// OnKeyDown registers a keyDown handler.
func (a *Action[S]) OnKeyDown(h func(context.Context, KeyDownEvent[S]) error) {
	a.onPayload(KeyDown, h)
}

// OnKeyUp registers a keyUp handler.
func (a *Action[S]) OnKeyUp(h func(context.Context, KeyUpEvent[S]) error) {
	a.onPayload(KeyUp, h)
}

// OnDidReceiveSettings registers a didReceiveSettings handler for unsolicited updates.
func (a *Action[S]) OnDidReceiveSettings(h func(context.Context, DidReceiveSettingsEvent[S]) error) {
	a.onPayload(DidReceiveSettings, h)
}

// OnDidReceiveResources registers a didReceiveResources handler for unsolicited updates.
func (a *Action[S]) OnDidReceiveResources(h func(context.Context, DidReceiveResourcesEvent[S]) error) {
	a.onPayload(DidReceiveResources, h)
}

// OnTouchTap registers a touchTap handler.
func (a *Action[S]) OnTouchTap(h func(context.Context, TouchTapEvent[S]) error) {
	a.onPayload(TouchTap, h)
}

// OnDialDown registers a dialDown handler.
func (a *Action[S]) OnDialDown(h func(context.Context, DialDownEvent[S]) error) {
	a.onPayload(DialDown, h)
}

// OnDialUp registers a dialUp handler.
func (a *Action[S]) OnDialUp(h func(context.Context, DialUpEvent[S]) error) {
	a.onPayload(DialUp, h)
}

// OnDialRotate registers a dialRotate handler.
func (a *Action[S]) OnDialRotate(h func(context.Context, DialRotateEvent[S]) error) {
	a.onPayload(DialRotate, h)
}

// OnTitleParametersDidChange registers a titleParametersDidChange handler.
func (a *Action[S]) OnTitleParametersDidChange(h func(context.Context, TitleParametersDidChangeEvent[S]) error) {
	a.onPayload(TitleParametersDidChange, h)
}

// OnPropertyInspectorDidAppear registers a propertyInspectorDidAppear handler.
func (a *Action[S]) OnPropertyInspectorDidAppear(h func(context.Context, PropertyInspectorDidAppearEvent[S]) error) {
	a.onPayload(PropertyInspectorDidAppear, h)
}

// OnPropertyInspectorDidDisappear registers a propertyInspectorDidDisappear handler.
func (a *Action[S]) OnPropertyInspectorDidDisappear(h func(context.Context, PropertyInspectorDidDisappearEvent[S]) error) {
	a.onPayload(PropertyInspectorDidDisappear, h)
}

// OnPropertyInspectorMessage registers a sendToPlugin handler with a message type independent of S.
func (a *Action[S]) OnPropertyInspectorMessage[M any](h func(context.Context, M) error) {
	a.RegisterHandler(SendToPlugin, func(ctx context.Context, ev Event) error {
		m, err := ev.Unmarshal[M]()
		if err != nil {
			return fmt.Errorf("sendToPlugin: %w", err)
		}
		return h(ctx, m)
	})
}

// SetSettings persists action-instance settings.
func (a *Action[S]) SetSettings(ctx context.Context, settings S) error {
	return a.client.sendCommand(ctx, SetSettings, settings)
}

// GetSettings requests action-instance settings and waits for the matching response.
func (a *Action[S]) GetSettings(ctx context.Context) (S, error) {
	var zero S
	resp, err := a.client.request(ctx, outgoingEvent{
		Event:   GetSettings,
		Context: firstNonEmpty(sdcontext.Context(ctx), ""),
	})
	if err != nil {
		return zero, err
	}
	p, err := resp.Unmarshal[DidReceiveSettingsPayload[S]]()
	if err != nil {
		return zero, err
	}
	return p.Settings, nil
}

// GetResources requests action resources and waits for the matching response.
func (a *Action[S]) GetResources(ctx context.Context) (map[string]string, error) {
	resp, err := a.client.request(ctx, outgoingEvent{
		Event:   GetResources,
		Context: sdcontext.Context(ctx),
	})
	if err != nil {
		return nil, err
	}
	p, err := resp.Unmarshal[DidReceiveResourcesPayload[S]]()
	if err != nil {
		return nil, err
	}
	return p.Resources, nil
}

// SetResources stores files associated with the action instance.
func (a *Action[S]) SetResources(ctx context.Context, resources map[string]string) error {
	return a.client.sendCommand(ctx, SetResources, resources)
}

// SendToPropertyInspector sends a typed message to the property inspector.
func (a *Action[S]) SendToPropertyInspector[M any](ctx context.Context, msg M) error {
	return a.client.send(ctx, outgoingEvent{
		Event:   SendToPropertyInspector,
		Action:  firstNonEmpty(sdcontext.Action(ctx), a.uuid),
		Context: sdcontext.Context(ctx),
		Payload: msg,
	})
}

// SetTitle changes the title of the action instance identified by ctx.
func (a *Action[S]) SetTitle(ctx context.Context, title string, target Target, state ...int) error {
	return a.client.SetTitle(ctx, title, target, state...)
}

// SetImage changes the image of the action instance identified by ctx.
func (a *Action[S]) SetImage(ctx context.Context, base64image string, target Target, state ...int) error {
	return a.client.SetImage(ctx, base64image, target, state...)
}

// SetFeedback updates Stream Deck + layout items for the instance identified by ctx.
func (a *Action[S]) SetFeedback(ctx context.Context, payload Feedback) error {
	return a.client.SetFeedback(ctx, payload)
}

// SetFeedbackLayout sets the layout for the instance identified by ctx.
func (a *Action[S]) SetFeedbackLayout(ctx context.Context, layout string) error {
	return a.client.SetFeedbackLayout(ctx, layout)
}

// ShowAlert shows a temporary alert icon.
func (a *Action[S]) ShowAlert(ctx context.Context) error {
	return a.client.ShowAlert(ctx)
}

// ShowOk shows a temporary OK icon.
func (a *Action[S]) ShowOk(ctx context.Context) error {
	return a.client.ShowOk(ctx)
}

// SetState selects a multi-state action state.
func (a *Action[S]) SetState(ctx context.Context, state int) error {
	return a.client.SetState(ctx, state)
}

// Contexts returns the currently visible action-instance contexts.
func (a *Action[S]) Contexts() []context.Context {
	a.mu.RLock()
	defer a.mu.RUnlock()
	cs := make([]context.Context, 0, len(a.contexts))
	for _, ctx := range a.contexts {
		cs = append(cs, ctx)
	}
	return cs
}

func (a *Action[S]) addContext(ctx context.Context) {
	id := sdcontext.Context(ctx)
	if id == "" {
		return
	}
	a.mu.Lock()
	a.contexts[id] = ctx
	a.mu.Unlock()
}

func (a *Action[S]) removeContext(ctx context.Context) {
	id := sdcontext.Context(ctx)
	if id == "" {
		return
	}
	a.mu.Lock()
	delete(a.contexts, id)
	a.mu.Unlock()
}
