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

func TestWithIDs(t *testing.T) {
	ctx := sdcontext.WithIDs(context.Background(), "ctx", "dev", "act")
	if sdcontext.Context(ctx) != "ctx" || sdcontext.Device(ctx) != "dev" || sdcontext.Action(ctx) != "act" {
		t.Fatalf("got %q %q %q", sdcontext.Context(ctx), sdcontext.Device(ctx), sdcontext.Action(ctx))
	}
}

func TestContextWithoutIDs(t *testing.T) {
	if sdcontext.Context(context.Background()) != "" {
		t.Fatal("background context should be empty")
	}
	type foreignKey struct{}
	ctx := context.WithValue(context.Background(), foreignKey{}, 1)
	if sdcontext.Device(ctx) != "" {
		t.Fatal("wrong type should be empty")
	}
}
