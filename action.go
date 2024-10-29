package streamdeck

import (
	"context"
	"sync"

	sdcontext "github.com/FlowingSPDG/streamdeck/context"
	"github.com/puzpuzpuz/xsync/v3"
	"golang.org/x/sync/errgroup"
)

// Action action instance
type Action struct {
	uuid     string
	handlers *eventHandlers
	contexts *contexts
}

type eventHandlers struct {
	m *xsync.MapOf[string, *eventHandlerSlice]
}

// []EventHandler
type eventHandlerSlice struct {
	mutex *sync.Mutex
	eh    []EventHandler
}

func (e *eventHandlerSlice) Execute(ctx context.Context, client *Client, event Event) error {
	eg, ectx := errgroup.WithContext(ctx)
	for _, handler := range e.eh {
		h := handler
		eg.Go(func() error {
			return h(ectx, client, event)
		})
	}
	return eg.Wait()
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
