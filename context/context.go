package context

import "context"

type keyType int

const (
	contextKey keyType = iota
	deviceKey
	actionKey
	idsKey
)

type ids struct {
	context string
	device  string
	action  string
}

func Context(ctx context.Context) string {
	if id, ok := idsFrom(ctx); ok {
		return id.context
	}
	return get(ctx, contextKey)
}

func WithContext(ctx context.Context, streamdeckContext string) context.Context {
	return context.WithValue(ctx, contextKey, streamdeckContext)
}

func Device(ctx context.Context) string {
	if id, ok := idsFrom(ctx); ok {
		return id.device
	}
	return get(ctx, deviceKey)
}

func WithDevice(ctx context.Context, streamdeckDevice string) context.Context {
	return context.WithValue(ctx, deviceKey, streamdeckDevice)
}

func Action(ctx context.Context) string {
	if id, ok := idsFrom(ctx); ok {
		return id.action
	}
	return get(ctx, actionKey)
}

func WithAction(ctx context.Context, streamdeckAction string) context.Context {
	return context.WithValue(ctx, actionKey, streamdeckAction)
}

// WithIDs attaches context, device, and action in one value to avoid extra allocations.
func WithIDs(ctx context.Context, contextID, device, action string) context.Context {
	return context.WithValue(ctx, idsKey, ids{
		context: contextID,
		device:  device,
		action:  action,
	})
}

func idsFrom(ctx context.Context) (ids, bool) {
	if ctx == nil {
		return ids{}, false
	}
	val := ctx.Value(idsKey)
	if val == nil {
		return ids{}, false
	}
	id, ok := val.(ids)
	return id, ok
}

func get(ctx context.Context, key keyType) string {
	if ctx == nil {
		return ""
	}

	val := ctx.Value(key)
	if val == nil {
		return ""
	}

	valStr, ok := val.(string)
	if !ok {
		return ""
	}

	return valStr
}
