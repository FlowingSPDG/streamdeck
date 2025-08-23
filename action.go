package streamdeck

import (
	"context"
	"fmt"
	"sync"

	sdcontext "github.com/FlowingSPDG/streamdeck/context"
	"github.com/puzpuzpuz/xsync/v3"
)

// Action action instance
type Action struct {
	uuid     string
	handlers *eventHandlers
	contexts *contexts
}

// TypedEventHandler is a type-safe event handler that automatically unmarshals the payload
type TypedEventHandler[T any] func(ctx context.Context, client *Client, payload T) error

type eventHandlers struct {
	m *xsync.MapOf[string, *eventHandlerSlice]
}

// []EventHandler
type eventHandlerSlice struct {
	mutex *sync.Mutex
	eh    []EventHandler
}

// Execute executes all registered event handlers for this event type.
// Handlers are executed sequentially to avoid race conditions.
func (e *eventHandlerSlice) Execute(ctx context.Context, client *Client, event Event) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	var lastErr error
	for _, handler := range e.eh {
		if err := handler(ctx, client, event); err != nil {
			lastErr = err
			// Log error but continue executing other handlers
			msg := fmt.Sprintf("Error in event handler: %s", err)
			client.LogMessage(ctx, msg)
		}
	}
	return lastErr
}

// map[string]context.Context
type contexts struct {
	m *xsync.MapOf[string, context.Context]
}

func newAction(uuid string) *Action {
	action := &Action{
		uuid: uuid,
		handlers: &eventHandlers{
			m: xsync.NewMapOf[string, *eventHandlerSlice](),
		},
		contexts: &contexts{m: xsync.NewMapOf[string, context.Context]()},
	}

	action.RegisterHandler(WillAppear, func(ctx context.Context, client *Client, event Event) error {
		action.addContext(ctx)
		return nil
	})

	action.RegisterHandler(WillDisappear, func(ctx context.Context, client *Client, event Event) error {
		action.removeContext(ctx)
		return nil
	})

	return action
}

// RegisterHandler Register event handler to specified event. handlers can be multiple(append slice)
func (action *Action) RegisterHandler(eventName string, handler EventHandler) {
	eh, _ := action.handlers.m.LoadOrStore(eventName, &eventHandlerSlice{
		mutex: &sync.Mutex{},
		eh:    []EventHandler{},
	})

	eh.mutex.Lock()
	defer eh.mutex.Unlock()
	eh.eh = append(eh.eh, handler)

	action.handlers.m.Store(eventName, eh)
}

// RegisterTypedHandler registers a type-safe event handler that automatically unmarshals the payload
func RegisterTypedHandler[T any](action *Action, eventName string, handler TypedEventHandler[T]) {
	action.RegisterHandler(eventName, func(ctx context.Context, client *Client, event Event) error {
		var payload T
		if err := event.UnmarshalPayload(&payload); err != nil {
			return fmt.Errorf("failed to unmarshal %s payload: %w", eventName, err)
		}
		return handler(ctx, client, payload)
	})
}

// Event-specific typed handler functions

// OnWillAppear registers a type-safe WillAppear event handler
func OnWillAppear[T any](action *Action, handler TypedEventHandler[WillAppearPayload[T]]) {
	RegisterTypedHandler(action, WillAppear, handler)
}

// OnWillDisappear registers a type-safe WillDisappear event handler
func OnWillDisappear[T any](action *Action, handler TypedEventHandler[WillDisappearPayload[T]]) {
	RegisterTypedHandler(action, WillDisappear, handler)
}

// OnKeyDown registers a type-safe KeyDown event handler
func OnKeyDown[T any](action *Action, handler TypedEventHandler[KeyDownPayload[T]]) {
	RegisterTypedHandler(action, KeyDown, handler)
}

// OnKeyUp registers a type-safe KeyUp event handler
func OnKeyUp[T any](action *Action, handler TypedEventHandler[KeyUpPayload[T]]) {
	RegisterTypedHandler(action, KeyUp, handler)
}

// OnDidReceiveSettings registers a type-safe DidReceiveSettings event handler
func OnDidReceiveSettings[T any](action *Action, handler TypedEventHandler[DidReceiveSettingsPayload[T]]) {
	RegisterTypedHandler(action, DidReceiveSettings, handler)
}

// OnTouchTap registers a type-safe TouchTap event handler
func OnTouchTap[T any](action *Action, handler TypedEventHandler[TouchTapPayload[T]]) {
	RegisterTypedHandler(action, TouchTap, handler)
}

// OnDialDown registers a type-safe DialDown event handler
func OnDialDown[T any](action *Action, handler TypedEventHandler[DialDownPayload[T]]) {
	RegisterTypedHandler(action, DialDown, handler)
}

// OnDialUp registers a type-safe DialUp event handler
func OnDialUp[T any](action *Action, handler TypedEventHandler[DialUpPayload[T]]) {
	RegisterTypedHandler(action, DialUp, handler)
}

// OnDialRotate registers a type-safe DialRotate event handler
func OnDialRotate[T any](action *Action, handler TypedEventHandler[DialRotatePayload[T]]) {
	RegisterTypedHandler(action, DialRotate, handler)
}

// OnPropertyInspectorDidAppear registers a type-safe PropertyInspectorDidAppear event handler
func OnPropertyInspectorDidAppear[T any](action *Action, handler TypedEventHandler[PropertyInspectorDidAppearPayload[T]]) {
	RegisterTypedHandler(action, PropertyInspectorDidAppear, handler)
}

// OnPropertyInspectorDidDisappear registers a type-safe PropertyInspectorDidDisappear event handler
func OnPropertyInspectorDidDisappear[T any](action *Action, handler TypedEventHandler[PropertyInspectorDidDisappearPayload[T]]) {
	RegisterTypedHandler(action, PropertyInspectorDidDisappear, handler)
}

// OnDidReceivePropertyInspectorMessage registers a type-safe DidReceivePropertyInspectorMessage event handler
func OnDidReceivePropertyInspectorMessage[T any](action *Action, handler TypedEventHandler[DidReceivePropertyInspectorMessagePayload[T]]) {
	RegisterTypedHandler(action, DidReceivePropertyInspectorMessage, handler)
}

// OnSendToPlugin registers a type-safe SendToPlugin event handler
func OnSendToPlugin[T any](action *Action, handler TypedEventHandler[T]) {
	RegisterTypedHandler(action, SendToPlugin, handler)
}

// Contexts get contexts
func (action *Action) Contexts() []context.Context {
	cs := make([]context.Context, 0, action.contexts.m.Size())
	action.contexts.m.Range(func(key string, v context.Context) bool {
		cs = append(cs, v)
		return true
	})
	return cs
}

func (action *Action) addContext(ctx context.Context) {
	if sdcontext.Context(ctx) == "" {
		panic("passed non-streamdeck context to addContext")
	}
	action.contexts.m.Store(sdcontext.Context(ctx), ctx)
}

func (action *Action) removeContext(ctx context.Context) {
	if sdcontext.Context(ctx) == "" {
		panic("passed non-streamdeck context to addContext")
	}
	action.contexts.m.Delete(sdcontext.Context(ctx))
}
