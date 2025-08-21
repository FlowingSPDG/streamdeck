package streamdeck

import (
	"context"
	"encoding/json"
	"testing"
)

func TestEvent_UnmarshalPayload(t *testing.T) {
	tests := []struct {
		name    string
		event   Event
		target  any
		wantErr bool
	}{
		{
			name: "valid payload",
			event: Event{
				Payload: json.RawMessage(`{"message": "test"}`),
			},
			target:  &LogMessagePayload{},
			wantErr: false,
		},
		{
			name: "nil payload",
			event: Event{
				Payload: nil,
			},
			target:  &LogMessagePayload{},
			wantErr: false,
		},
		{
			name: "invalid json",
			event: Event{
				Payload: json.RawMessage(`{"message": "test"`),
			},
			target:  &LogMessagePayload{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.UnmarshalPayload(tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Event.UnmarshalPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEvent(t *testing.T) {
	ctx := context.Background()
	payload := LogMessagePayload{Message: "test"}

	event := NewEvent(ctx, LogMessage, payload)

	if event.Event != LogMessage {
		t.Errorf("NewEvent() Event = %v, want %v", event.Event, LogMessage)
	}

	if event.Payload == nil {
		t.Error("NewEvent() Payload should not be nil")
	}
}

func TestSetTriggerDescriptionPayload(t *testing.T) {
	payload := SetTriggerDescriptionPayload{
		LongTouch: "Long touch description",
		Push:      "Push description",
		Rotate:    "Rotate description",
		Touch:     "Touch description",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal SetTriggerDescriptionPayload: %v", err)
	}

	var unmarshaled SetTriggerDescriptionPayload
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal SetTriggerDescriptionPayload: %v", err)
	}

	if unmarshaled.LongTouch != payload.LongTouch {
		t.Errorf("LongTouch = %v, want %v", unmarshaled.LongTouch, payload.LongTouch)
	}
}
