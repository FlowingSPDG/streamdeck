package streamdeck

import (
	"context"
	"encoding/json"

	sdcontext "github.com/FlowingSPDG/streamdeck/context"
)

// Event JSON struct. {"action":"com.elgato.example.action1","event":"keyDown","context":"","device":"","payload":{"settings":{},"coordinates":{"column":3,"row":1},"state":0,"userDesiredState":1,"isInMultiAction":false}}
type Event struct {
	Action     string          `json:"action,omitempty"`
	Event      string          `json:"event,omitempty"`
	UUID       string          `json:"uuid,omitempty"`
	Context    string          `json:"context,omitempty"`
	Device     string          `json:"device,omitempty"`
	DeviceInfo DeviceInfo      `json:"deviceInfo,omitempty"`
	Payload    json.RawMessage `json:"payload,omitempty"`
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

// NewEvent Generate new event from specified name and payload. payload will be converted to JSON
func NewEvent(ctx context.Context, name string, payload any) Event {
	p, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	return Event{
		Event:   name,
		Action:  sdcontext.Action(ctx),
		Context: sdcontext.Context(ctx),
		Device:  sdcontext.Device(ctx),
		Payload: p,
	}
}
