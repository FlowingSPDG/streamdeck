package streamdeck

import (
	"encoding/json"
	"fmt"
)

// Event is a received Stream Deck WebSocket message.
type Event struct {
	Action     string          `json:"action,omitempty"`
	Event      EventName       `json:"event,omitempty"`
	UUID       string          `json:"uuid,omitempty"`
	Context    string          `json:"context,omitempty"`
	Device     string          `json:"device,omitempty"`
	DeviceInfo DeviceInfo      `json:"deviceInfo,omitempty"`
	ID         string          `json:"id,omitempty"`
	Payload    json.RawMessage `json:"payload,omitempty"`
}

// DeviceInfo describes a Stream Deck device.
type DeviceInfo struct {
	Name string     `json:"name,omitempty"`
	Type DeviceType `json:"type,omitempty"`
	Size DeviceSize `json:"size,omitempty"`
}

// DeviceSize is the key grid size of a device, excluding dials and touchscreens.
type DeviceSize struct {
	Columns int `json:"columns,omitempty"`
	Rows    int `json:"rows,omitempty"`
}

// DeviceType is a Stream Deck hardware type.
// See https://github.com/elgatosf/schemas/blob/main/src/streamdeck/plugins/device-type.ts
type DeviceType int

const (
	StreamDeck        DeviceType = 0
	StreamDeckMini    DeviceType = 1
	StreamDeckXL      DeviceType = 2
	StreamDeckMobile  DeviceType = 3
	CorsairGKeys      DeviceType = 4
	StreamDeckPedal   DeviceType = 5
	CorsairVoyager    DeviceType = 6
	StreamDeckPlus    DeviceType = 7
	SCUFController    DeviceType = 8
	StreamDeckNeo     DeviceType = 9
	StreamDeckStudio  DeviceType = 10
	VirtualStreamDeck DeviceType = 11
	Galleon100SD      DeviceType = 12
	StreamDeckPlusXL  DeviceType = 13
)

// Unmarshal decodes the event payload as T.
func (e Event) Unmarshal[T any]() (T, error) {
	var payload T
	if len(e.Payload) == 0 {
		return payload, nil
	}
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		return payload, fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	return payload, nil
}

// ActionEvent is a typed event bound to an action instance.
type ActionEvent[P any] struct {
	Action  string
	Context string
	Device  string
	ID      string
	Payload P
}

type (
	WillAppearEvent[S any]                    = ActionEvent[WillAppearPayload[S]]
	WillDisappearEvent[S any]                 = ActionEvent[WillDisappearPayload[S]]
	KeyDownEvent[S any]                       = ActionEvent[KeyDownPayload[S]]
	KeyUpEvent[S any]                         = ActionEvent[KeyUpPayload[S]]
	DidReceiveSettingsEvent[S any]            = ActionEvent[DidReceiveSettingsPayload[S]]
	DidReceiveResourcesEvent[S any]           = ActionEvent[DidReceiveResourcesPayload[S]]
	TouchTapEvent[S any]                      = ActionEvent[TouchTapPayload[S]]
	DialDownEvent[S any]                      = ActionEvent[DialDownPayload[S]]
	DialUpEvent[S any]                        = ActionEvent[DialUpPayload[S]]
	DialRotateEvent[S any]                    = ActionEvent[DialRotatePayload[S]]
	TitleParametersDidChangeEvent[S any]      = ActionEvent[TitleParametersDidChangePayload[S]]
	PropertyInspectorDidAppearEvent[S any]    = ActionEvent[PropertyInspectorDidAppearPayload[S]]
	PropertyInspectorDidDisappearEvent[S any] = ActionEvent[PropertyInspectorDidDisappearPayload[S]]
)

type outgoingEvent struct {
	Action  string    `json:"action,omitempty"`
	Event   EventName `json:"event"`
	UUID    string    `json:"uuid,omitempty"`
	Context string    `json:"context,omitempty"`
	Device  string    `json:"device,omitempty"`
	ID      string    `json:"id,omitempty"`
	Payload any       `json:"payload,omitempty"`
}

func actionEvent[P any](ev Event, payload P) ActionEvent[P] {
	return ActionEvent[P]{
		Action:  ev.Action,
		Context: ev.Context,
		Device:  ev.Device,
		ID:      ev.ID,
		Payload: payload,
	}
}
