# Stream Deck WebSocket Library for Go

A Go library for creating Stream Deck plugins using the WebSocket API.

## Features

- Full WebSocket API support for Stream Deck plugins
- Type-safe event handling with generics
- Thread-safe event handler registration
- Comprehensive error handling
- Support for all Stream Deck devices (including Stream Deck +)
- Built-in image conversion utilities
- Context-aware event processing

## Quick Start

```go
package main

import (
    "context"
    "os"
    "github.com/FlowingSPDG/streamdeck"
)

type Settings struct {
    Counter int `json:"counter"`
}

func main() {
    ctx := context.Background()
    params, err := streamdeck.ParseRegistrationParams(os.Args)
    if err != nil {
        panic(err)
    }

    client := streamdeck.NewClient(ctx, params)
    setup(client)
    
    if err := client.Run(ctx); err != nil {
        panic(err)
    }
}

func setup(client *streamdeck.Client) {
    action := client.Action("com.example.counter")
    
    action.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
        var payload streamdeck.KeyDownPayload[Settings]
        if err := event.UnmarshalPayload(&payload); err != nil {
            return err
        }
        
        // Handle key press
        return client.SetTitle(ctx, "Pressed!", streamdeck.HardwareAndSoftware)
    })
}
```

## New Features in Latest Version

### Enhanced Event Handling
- Type-safe payload unmarshaling with `event.UnmarshalPayload()`
- Improved error handling and logging
- Sequential event handler execution to prevent race conditions

### New Commands
- `SetTriggerDescription()` - Set encoder action descriptions
- `SetFeedbackLayout()` - Set Stream Deck + touch display layouts
- Enhanced `SwitchToProfile()` with page support
- State-aware `SetTitle()` and `SetImage()` methods

### Device Support
- Full support for Stream Deck + encoder and touch interactions
- Device change event handling
- Improved device type definitions

## Examples

See the `examples/` directory for complete working examples:
- `counter/` - Basic counter plugin
- `cpu/` - CPU usage monitor
- `mock/` - Mock plugin for testing

## API Reference

### Core Types
- `Client` - Main client for Stream Deck communication
- `Event` - WebSocket event structure
- `Action` - Action instance management

### Event Handling
- `RegisterHandler()` - Register event handlers
- `UnmarshalPayload()` - Type-safe payload parsing
- `NewEvent()` - Create new events

### Commands
- `SetTitle()` - Set button title
- `SetImage()` - Set button image
- `SetSettings()` - Save persistent data
- `SetFeedback()` - Stream Deck + touch display
- `SetTriggerDescription()` - Encoder descriptions

## Migration from Previous Versions

The library maintains backward compatibility while adding new features:
- Existing event handlers continue to work
- New optional parameters for enhanced functionality
- Improved error handling without breaking changes

## License

MIT License - see LICENSE file for details.