package streamdeck

import (
	"context"
	"fmt"
	"strconv"
	"sync"
)

type queuedEvent struct {
	ctx context.Context
	ev  Event
}

// eventQueue is an unbounded, FIFO mailbox. The WebSocket reader only takes a
// mutex and never blocks on a handler. A burst such as dialRotate is drained
// in one batch after the current handler returns.
type eventQueue struct {
	mu     sync.Mutex
	items  []queuedEvent
	wait   chan struct{}
	cancel context.CancelFunc
}

func newEventQueue(cancel context.CancelFunc, capHint int) *eventQueue {
	if capHint < 1 {
		capHint = defaultQueueSize
	}
	return &eventQueue{
		items:  make([]queuedEvent, 0, capHint),
		wait:   make(chan struct{}, 1),
		cancel: cancel,
	}
}

func (q *eventQueue) push(item queuedEvent) {
	q.mu.Lock()
	q.items = append(q.items, item)
	q.mu.Unlock()
	select {
	case q.wait <- struct{}{}:
	default:
	}
}

func (q *eventQueue) take() []queuedEvent {
	q.mu.Lock()
	items := q.items
	q.items = nil
	q.mu.Unlock()
	return items
}

func (c *Client) nextRequestID() string {
	return strconv.FormatUint(c.reqSeq.Add(1), 10)
}

func (c *Client) request(ctx context.Context, ev outgoingEvent) (Event, error) {
	id := c.nextRequestID()
	ev.ID = id
	ch := make(chan Event, 1)

	c.pendingMu.Lock()
	c.pending[id] = ch
	c.pendingMu.Unlock()
	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}()

	if err := c.send(ctx, ev); err != nil {
		return Event{}, err
	}

	select {
	case resp := <-ch:
		return resp, nil
	case <-ctx.Done():
		return Event{}, ctx.Err()
	}
}

func (c *Client) deliverPending(ev Event) bool {
	if ev.ID == "" {
		return false
	}
	c.pendingMu.Lock()
	ch, ok := c.pending[ev.ID]
	c.pendingMu.Unlock()
	if !ok {
		return false
	}
	select {
	case ch <- ev:
	default:
	}
	return true
}

func (c *Client) enqueue(ctx context.Context, ev Event) {
	if ctx.Err() != nil {
		return
	}
	c.getOrCreateQueue(ev.Context).push(queuedEvent{ctx: ctx, ev: ev})
}

func (c *Client) getOrCreateQueue(key string) *eventQueue {
	c.queuesMu.Lock()
	defer c.queuesMu.Unlock()
	if q, ok := c.queues[key]; ok {
		return q
	}
	qctx, cancel := context.WithCancel(c.runCtx())
	q := newEventQueue(cancel, c.queueSize)
	c.queues[key] = q
	c.workers.Add(1)
	go func() {
		defer c.workers.Done()
		c.runQueue(qctx, q)
	}()
	return q
}

func (c *Client) runCtx() context.Context {
	c.runCtxMu.RLock()
	defer c.runCtxMu.RUnlock()
	if c.runCancelCtx != nil {
		return c.runCancelCtx
	}
	return context.Background()
}

func (c *Client) setRunCtx(ctx context.Context) {
	c.runCtxMu.Lock()
	c.runCancelCtx = ctx
	c.runCtxMu.Unlock()
}

func (c *Client) runQueue(ctx context.Context, q *eventQueue) {
	for {
		items := q.take()
		if len(items) == 0 {
			select {
			case <-ctx.Done():
				return
			case <-q.wait:
			}
			continue
		}
		for i := range items {
			if ctx.Err() != nil {
				return
			}
			c.handleEvent(items[i].ctx, items[i].ev)
		}
	}
}

func (c *Client) cancelQueues() {
	c.queuesMu.Lock()
	qs := make([]*eventQueue, 0, len(c.queues))
	for key, q := range c.queues {
		qs = append(qs, q)
		delete(c.queues, key)
	}
	c.queuesMu.Unlock()
	for _, q := range qs {
		q.cancel()
	}

	c.pendingMu.Lock()
	c.pending = make(map[string]chan Event)
	c.pendingMu.Unlock()
}

func (c *Client) handleEvent(ctx context.Context, ev Event) {
	defer func() {
		if rec := recover(); rec != nil {
			c.handleError(ctx, fmt.Errorf("panic in event handler: %v", rec))
		}
	}()

	if ev.Action == "" {
		c.mu.RLock()
		hs := c.handlers[ev.Event]
		c.mu.RUnlock()
		for _, h := range hs {
			c.handleError(ctx, callHandler(func() error { return h(ctx, c, ev) }))
		}
		return
	}

	c.mu.RLock()
	a, ok := c.actions[ev.Action]
	c.mu.RUnlock()
	if !ok {
		c.logger.Warn("discarding event for unregistered action", "action", ev.Action, "event", ev.Event)
		return
	}
	_ = a.dispatch(ctx, ev)
}

func callHandler(fn func() error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic in event handler: %v", rec)
		}
	}()
	return fn()
}
