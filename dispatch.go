package streamdeck

import (
	"context"
	"strconv"
)

type queuedEvent struct {
	ctx context.Context
	ev  Event
}

type eventQueue struct {
	ch     chan queuedEvent
	cancel context.CancelFunc
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
	key := ev.Context
	q := c.getOrCreateQueue(key)
	select {
	case q.ch <- queuedEvent{ctx: ctx, ev: ev}:
	case <-ctx.Done():
	}
}

func (c *Client) getOrCreateQueue(key string) *eventQueue {
	c.queuesMu.Lock()
	defer c.queuesMu.Unlock()
	if q, ok := c.queues[key]; ok {
		return q
	}
	qctx, cancel := context.WithCancel(context.Background())
	q := &eventQueue{
		ch:     make(chan queuedEvent, c.queueSize),
		cancel: cancel,
	}
	c.queues[key] = q
	go c.runQueue(qctx, q)
	return q
}

func (c *Client) runQueue(ctx context.Context, q *eventQueue) {
	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-q.ch:
			if !ok {
				return
			}
			c.handleEvent(item.ctx, item.ev)
			if item.ev.Event == WillDisappear && item.ev.Context != "" {
				c.removeQueue(item.ev.Context)
			}
		}
	}
}

func (c *Client) removeQueue(key string) {
	c.queuesMu.Lock()
	q, ok := c.queues[key]
	if ok {
		delete(c.queues, key)
	}
	c.queuesMu.Unlock()
	if ok {
		q.cancel()
	}
}

func (c *Client) stopQueues() {
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
	if ev.Action == "" {
		c.mu.RLock()
		hs := append([]EventHandler(nil), c.handlers[ev.Event]...)
		c.mu.RUnlock()
		for _, h := range hs {
			c.handleError(ctx, h(ctx, c, ev))
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
	c.handleError(ctx, a.dispatch(ctx, ev))
}
