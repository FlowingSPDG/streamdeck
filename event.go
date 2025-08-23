package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"

	sdcontext "github.com/FlowingSPDG/streamdeck/context"
)

// Event JSON struct. {"action":"com.elgato.example.action1","event":"keyDown","context":"","device":"","payload":{"settings":{},"coordinates":{"column":3,"row":1},"state":0,"userDesiredState":1,"isInMultiAction":false}}
type Event struct {
	Action     string     `json:"action,omitempty"`
	Event      string     `json:"event,omitempty"`
	UUID       string     `json:"uuid,omitempty"`
	Context    string     `json:"context,omitempty"`
	Device     string     `json:"device,omitempty"`
	DeviceInfo DeviceInfo `json:"deviceInfo,omitempty"`
	Payload    any        `json:"payload,omitempty"`
}

// TypedEvent is a type-safe version of Event with a specific payload type
type TypedEvent[T any] struct {
	Action     string     `json:"action,omitempty"`
	Event      string     `json:"event,omitempty"`
	UUID       string     `json:"uuid,omitempty"`
	Context    string     `json:"context,omitempty"`
	Device     string     `json:"device,omitempty"`
	DeviceInfo DeviceInfo `json:"deviceInfo,omitempty"`
	Payload    T          `json:"payload,omitempty"`
}

// DeviceInfo A json object containing information about the device.. {"deviceInfo":{"name":"Device Name","type":0,"size":{"columns":5,"rows":3}}}
type DeviceInfo struct {
	DeviceName string     `json:"deviceName,omitempty"`
	Type       DeviceType `json:"type,omitempty"`
	Size       DeviceSize `json:"size,omitempty"`
}

// DeviceSize The number of columns and rows of keys that the device owns. {"columns":5,"rows":3}
type DeviceSize struct {
	Columns int `json:"columns,omitempty"`
	Rows    int `json:"rows,omitempty"`
}

// DeviceType Type of device. Possible values are kESDSDKDeviceType_StreamDeck (0), kESDSDKDeviceType_StreamDeckMini (1), kESDSDKDeviceType_StreamDeckXL (2), kESDSDKDeviceType_StreamDeckMobile (3), kESDSDKDeviceType_CorsairGKeys (4), kESDSDKDeviceType_StreamDeckPedal (5) and kESDSDKDeviceType_CorsairVoyager (6), kESDSDKDeviceType_StreamDeckPlus (7).
type DeviceType int

const (
	// StreamDeck kESDSDKDeviceType_StreamDeck (0)
	StreamDeck DeviceType = iota
	// StreamDeckMini kESDSDKDeviceType_StreamDeckMini (1)
	StreamDeckMini
	// StreamDeckXL kESDSDKDeviceType_StreamDeckXL (2)
	StreamDeckXL
	// StreamDeckMobile kESDSDKDeviceType_StreamDeckMobile (3)
	StreamDeckMobile
	// CorsairGKeys kESDSDKDeviceType_CorsairGKeys (4)
	CorsairGKeys
	// StreamDeckPedal kESDSDKDeviceType_StreamDeckPedal (5)
	StreamDeckPedal
	// CorsairVoyager kESDSDKDeviceType_CorsairVoyager (6)
	CorsairVoyager
	// StreamDeckPlus kESDSDKDeviceType_StreamDeckPlus (7)
	StreamDeckPlus
	// SCUFController kESDSDKDeviceType_SCUFController (8)
	SCUFController
	// StreamDeckNeo kESDSDKDeviceType_StreamDeckNeo (9)
	StreamDeckNeo
	// StreamDeck Studio...?
	StreamDeckStudio
)

// NewEvent Generate new event from specified name and payload. payload will be stored as raw data
func NewEvent(ctx context.Context, name string, payload any) Event {
	return Event{
		Event:   name,
		Action:  sdcontext.Action(ctx),
		Context: sdcontext.Context(ctx),
		Device:  sdcontext.Device(ctx),
		Payload: payload,
	}
}

// UnmarshalPayload safely unmarshals the event payload into the specified type.
// Returns an error if the payload cannot be unmarshaled into the target type.
func (e Event) UnmarshalPayload(target any) error {
	if e.Payload == nil {
		return nil
	}

	// If Payload is already a []byte, use it directly
	if data, ok := e.Payload.([]byte); ok {
		return json.Unmarshal(data, target)
	}

	// Otherwise, marshal the payload to JSON first
	data, err := json.Marshal(e.Payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// MustUnmarshalPayload unmarshals the event payload into the specified type.
// Panics if the payload cannot be unmarshaled into the target type.
func (e Event) MustUnmarshalPayload(target any) {
	if err := e.UnmarshalPayload(target); err != nil {
		panic(fmt.Sprintf("failed to unmarshal event payload: %v", err))
	}
}

// ToTypedEvent converts an Event to a TypedEvent with the specified payload type
func ToTypedEvent[T any](e Event) (TypedEvent[T], error) {
	var typedEvent TypedEvent[T]

	// Copy common fields
	typedEvent.Action = e.Action
	typedEvent.Event = e.Event
	typedEvent.UUID = e.UUID
	typedEvent.Context = e.Context
	typedEvent.Device = e.Device
	typedEvent.DeviceInfo = e.DeviceInfo

	// Unmarshal payload
	if err := e.UnmarshalPayload(&typedEvent.Payload); err != nil {
		return TypedEvent[T]{}, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return typedEvent, nil
}

// NewTypedEvent creates a new TypedEvent with the specified payload type
func NewTypedEvent[T any](ctx context.Context, name string, payload T) TypedEvent[T] {
	return TypedEvent[T]{
		Event:   name,
		Action:  sdcontext.Action(ctx),
		Context: sdcontext.Context(ctx),
		Device:  sdcontext.Device(ctx),
		Payload: payload,
	}
}
