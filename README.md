# StreamDeck Go SDK

A Go SDK for building Elgato StreamDeck plugins.

## Features

- Full StreamDeck plugin API support
- Type-safe event handling with generics
- Automatic payload unmarshaling
- Context management
- Image and title manipulation
- Settings persistence

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/FlowingSPDG/streamdeck"
)

type Settings struct {
	Counter    int    `json:"counter"`
	ButtonText string `json:"buttonText"`
}

func main() {
	ctx := context.Background()
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	client := streamdeck.NewClient(ctx, params)
	setup(client)
	
	if err := client.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

func setup(client *streamdeck.Client) {
	action := client.Action("com.example.myaction")

	// Traditional event handling
	action.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		var p streamdeck.WillAppearPayload[Settings]
		if err := event.UnmarshalPayload(&p); err != nil {
			return fmt.Errorf("failed to unmarshal WillAppear payload: %w", err)
		}
		
		// Handle the payload
		return client.SetTitle(ctx, p.Settings.ButtonText, streamdeck.HardwareAndSoftware)
	})

	// Type-safe event handling with automatic unmarshaling
	streamdeck.RegisterTypedHandler(action, streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
		// Payload is automatically unmarshaled and typed
		p.Settings.Counter++
		return client.SetSettings(ctx, p.Settings)
	})
}
```

## Type-Safe Event Handling

The SDK provides intuitive event-specific functions that automatically unmarshal event payloads into the specified type:

```go
// Define your settings type
type MySettings struct {
	Counter int    `json:"counter"`
	Text    string `json:"text"`
}

// Define your message types for Property Inspector communication
type PropertyInspectorMessage struct {
	Action string `json:"action"`
}

type ResetCompleteResponse struct {
	Action string `json:"action"`
}

// Register typed handlers for different events using intuitive method names
streamdeck.OnWillAppear(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.WillAppearPayload[MySettings]) error {
	// p.Settings is automatically typed as MySettings
	return client.SetTitle(ctx, p.Settings.Text, streamdeck.HardwareAndSoftware)
})

streamdeck.OnKeyDown(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[MySettings]) error {
	// p.Settings is automatically typed as MySettings
	p.Settings.Counter++
	return client.SetSettings(ctx, p.Settings)
})

streamdeck.OnDidReceiveSettings(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.DidReceiveSettingsPayload[MySettings]) error {
	// p.Settings is automatically typed as MySettings
	return nil
})

// Fully typed Property Inspector communication
streamdeck.OnSendToPlugin(action, func(ctx context.Context, client *streamdeck.Client, payload PropertyInspectorMessage) error {
	switch payload.Action {
	case "reset":
		return client.SendToPropertyInspector(ctx, ResetCompleteResponse{
			Action: "resetComplete",
		})
	}
	return nil
})
```

### Benefits of Type-Safe Handlers

1. **Automatic Unmarshaling**: No need to manually call `event.UnmarshalPayload()`
2. **Complete Type Safety**: Zero `any` types or type assertions - fully typed from input to output
3. **Intuitive API**: Event-specific function names (e.g., `OnKeyDown`, `OnWillAppear`)
4. **Cleaner Code**: Reduced boilerplate and error handling
5. **Better IDE Support**: Full autocomplete and type inference
6. **Compile-time Safety**: Catch type mismatches at build time, not runtime

### Supported Event Functions

- `OnWillAppear[T]` - Button appears on Stream Deck
- `OnWillDisappear[T]` - Button disappears from Stream Deck
- `OnKeyDown[T]` - Button pressed
- `OnKeyUp[T]` - Button released
- `OnDidReceiveSettings[T]` - Settings received
- `OnTouchTap[T]` - Touch screen tapped (Stream Deck +)
- `OnDialDown[T]` - Dial pressed (Stream Deck +)
- `OnDialUp[T]` - Dial released (Stream Deck +)
- `OnDialRotate[T]` - Dial rotated (Stream Deck +)
- `OnPropertyInspectorDidAppear[T]` - Property Inspector opened
- `OnPropertyInspectorDidDisappear[T]` - Property Inspector closed
- `OnDidReceivePropertyInspectorMessage[T]` - Message from Property Inspector
- `OnSendToPlugin[T]` - Plugin message received

## Traditional Event Handling

For cases where you need more control or want to handle raw events, you can still use the traditional `RegisterHandler` method:

```go
action.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	var p streamdeck.WillAppearPayload[Settings]
	if err := event.UnmarshalPayload(&p); err != nil {
		return fmt.Errorf("failed to unmarshal WillAppear payload: %w", err)
	}
	
	// Handle the payload
	return nil
})
```

## Examples

See the `examples/` directory for complete working examples:

- `counter/` - Simple counter plugin
- `settings_manager/` - Advanced settings management with type-safe handlers

## License

MIT License