package streamdeck

import "context"

// On registers a handler for events that are not bound to an action UUID.
func (c *Client) On(event EventName, handler EventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[event] = appendCopy(c.handlers[event], handler)
}

// RegisterNoActionHandler is an alias of On for events such as applicationDidLaunch.
func (c *Client) RegisterNoActionHandler(event EventName, handler EventHandler) {
	c.On(event, handler)
}

func (c *Client) onPayload[P any](name EventName, h func(context.Context, P) error) {
	c.On(name, func(ctx context.Context, _ *Client, event Event) error {
		p, err := event.Unmarshal[P]()
		if err != nil {
			return err
		}
		return h(ctx, p)
	})
}

// OnApplicationDidLaunch registers a handler for applicationDidLaunch.
func (c *Client) OnApplicationDidLaunch(h func(context.Context, ApplicationDidLaunchPayload) error) {
	c.onPayload(ApplicationDidLaunch, h)
}

// OnApplicationDidTerminate registers a handler for applicationDidTerminate.
func (c *Client) OnApplicationDidTerminate(h func(context.Context, ApplicationDidTerminatePayload) error) {
	c.onPayload(ApplicationDidTerminate, h)
}

func (c *Client) onDevice(name EventName, h func(context.Context, DeviceEvent) error) {
	c.On(name, func(ctx context.Context, _ *Client, event Event) error {
		return h(ctx, DeviceEvent{Device: event.Device, DeviceInfo: event.DeviceInfo})
	})
}

// OnDeviceDidConnect registers a handler for deviceDidConnect.
func (c *Client) OnDeviceDidConnect(h func(context.Context, DeviceEvent) error) {
	c.onDevice(DeviceDidConnect, h)
}

// OnDeviceDidDisconnect registers a handler for deviceDidDisconnect.
func (c *Client) OnDeviceDidDisconnect(h func(context.Context, DeviceEvent) error) {
	c.onDevice(DeviceDidDisconnect, h)
}

// OnDeviceDidChange registers a handler for deviceDidChange.
func (c *Client) OnDeviceDidChange(h func(context.Context, DeviceEvent) error) {
	c.onDevice(DeviceDidChange, h)
}

// OnSystemDidWakeUp registers a handler for systemDidWakeUp.
func (c *Client) OnSystemDidWakeUp(h func(context.Context) error) {
	c.On(SystemDidWakeUp, func(ctx context.Context, _ *Client, _ Event) error {
		return h(ctx)
	})
}

// OnDidReceiveDeepLink registers a handler for didReceiveDeepLink.
func (c *Client) OnDidReceiveDeepLink(h func(context.Context, DidReceiveDeepLinkPayload) error) {
	c.onPayload(DidReceiveDeepLink, h)
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

func appendCopy[T any](old []T, item T) []T {
	next := make([]T, len(old)+1)
	copy(next, old)
	next[len(old)] = item
	return next
}
