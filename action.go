package streamdeck

import (
	"context"
	"sync"

	sdcontext "github.com/FlowingSPDG/streamdeck/context"
	"golang.org/x/sync/errgroup"
)

// Action action instance
type Action struct {
	uuid     string
	handlers eventHandlers
	contexts contexts
}

type eventHandlers struct {
	m sync.Map // map[string]eventHandlerSlice
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
	m sync.Map
}

func newAction(uuid string) *Action {
	action := &Action{
		uuid: uuid,
		handlers: eventHandlers{
			m: sync.Map{},
		},
		contexts: contexts{m: sync.Map{}},
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
	eh := eventHandlerSlice{
		mutex: &sync.Mutex{},
		eh:    []EventHandler{},
	}

	ehi, loaded := action.handlers.m.LoadOrStore(eventName, eh)
	if loaded {
		eh = ehi.(eventHandlerSlice)
	}

	eh.mutex.Lock()
	defer eh.mutex.Unlock()
	eh.eh = append(eh.eh, handler)

	action.handlers.m.Store(eventName, eh)
}

// Contexts get contexts
func (action *Action) Contexts() []context.Context {
	cs := make([]context.Context, 0) // 0 length/capacity
	action.contexts.m.Range(func(key, value any) bool {
		v := value.(context.Context)
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
