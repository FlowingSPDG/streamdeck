package streamdeck_test

import (
	"context"
	"testing"

	"github.com/FlowingSPDG/streamdeck/v2"
)

func TestNewActionSameTypeIsStable(t *testing.T) {
	client := streamdeck.NewClient(context.Background(), testParams())
	a1 := streamdeck.NewAction[struct{ N int }](client, "com.example.a")
	a2 := streamdeck.NewAction[struct{ N int }](client, "com.example.a")
	if a1 != a2 {
		t.Fatal("expected the same action instance")
	}
}

func TestNewActionDifferentTypePanics(t *testing.T) {
	client := streamdeck.NewClient(context.Background(), testParams())
	_ = streamdeck.NewAction[struct{ A int }](client, "com.example.a")
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	_ = streamdeck.NewAction[struct{ B string }](client, "com.example.a")
}
