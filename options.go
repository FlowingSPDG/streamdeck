package streamdeck

import (
	"context"
	"log/slog"
)

const defaultQueueSize = 64

// EventHandler handles a received event that is not bound to a registered action.
type EventHandler func(ctx context.Context, client *Client, event Event) error

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithLogger sets the client logger. The default logger discards output.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// WithErrorHandler is called when an event handler returns an error.
func WithErrorHandler(h func(context.Context, error)) ClientOption {
	return func(c *Client) {
		c.errHandler = h
	}
}

// WithWebSocketURL overrides the Stream Deck WebSocket URL. Intended for tests.
func WithWebSocketURL(u string) ClientOption {
	return func(c *Client) {
		c.dialURL = u
	}
}

// WithQueueSize sets the initial per-context event buffer capacity.
// The mailbox grows if a handler is slower than incoming events (for example dialRotate bursts).
// Zero or negative uses the default.
func WithQueueSize(n int) ClientOption {
	return func(c *Client) {
		if n > 0 {
			c.queueSize = n
		}
	}
}
