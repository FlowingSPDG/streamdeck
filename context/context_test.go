package context_test

import (
	"context"
	"testing"

	sdcontext "github.com/FlowingSPDG/streamdeck/v2/context"
)

func TestContextRoundTrip(t *testing.T) {
	ctx := context.Background()
	ctx = sdcontext.WithContext(ctx, "ctx")
	ctx = sdcontext.WithDevice(ctx, "dev")
	ctx = sdcontext.WithAction(ctx, "act")
	if sdcontext.Context(ctx) != "ctx" || sdcontext.Device(ctx) != "dev" || sdcontext.Action(ctx) != "act" {
		t.Fatalf("got %q %q %q", sdcontext.Context(ctx), sdcontext.Device(ctx), sdcontext.Action(ctx))
	}
}

func TestContextNilAndWrongType(t *testing.T) {
	if sdcontext.Context(nil) != "" {
		t.Fatal("nil context should be empty")
	}
	ctx := context.WithValue(context.Background(), struct{}{}, 1)
	if sdcontext.Device(ctx) != "" {
		t.Fatal("wrong type should be empty")
	}
}
