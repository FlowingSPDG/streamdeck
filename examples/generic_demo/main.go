package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/FlowingSPDG/streamdeck"
	"golang.org/x/xerrors"
)

// DemoSettings demonstrates the type-safe settings structure
type DemoSettings struct {
	Counter    int    `json:"counter"`
	ButtonText string `json:"buttonText"`
	Color      string `json:"color"`
}

// PropertyInspectorMessage represents messages from Property Inspector
type PropertyInspectorMessage struct {
	Action string `json:"action"`
}

// ResetCompleteResponse represents reset completion response
type ResetCompleteResponse struct {
	Action string `json:"action"`
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
	action := client.Action("com.example.generic_demo")

	// Example 1: WillAppear with automatic unmarshaling
	streamdeck.OnWillAppear(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.WillAppearPayload[DemoSettings]) error {
		// Settings are automatically typed as DemoSettings
		if p.Settings.ButtonText == "" {
			p.Settings.ButtonText = "Click Me!"
		}
		if p.Settings.Color == "" {
			p.Settings.Color = "blue"
		}

		// Set the title
		title := fmt.Sprintf("%s\n%d", p.Settings.ButtonText, p.Settings.Counter)
		return client.SetTitle(ctx, title, streamdeck.HardwareAndSoftware)
	})

	// Example 2: KeyDown with automatic unmarshaling
	streamdeck.OnKeyDown(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[DemoSettings]) error {
		// Increment counter
		p.Settings.Counter++

		// Save settings
		if err := client.SetSettings(ctx, p.Settings); err != nil {
			return xerrors.Errorf("failed to save settings: %w", err)
		}

		// Update title
		title := fmt.Sprintf("%s\n%d", p.Settings.ButtonText, p.Settings.Counter)
		return client.SetTitle(ctx, title, streamdeck.HardwareAndSoftware)
	})

	// Example 3: WillDisappear with automatic unmarshaling
	streamdeck.OnWillDisappear(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.WillDisappearPayload[DemoSettings]) error {
		// Reset counter when button disappears
		p.Settings.Counter = 0
		return client.SetSettings(ctx, p.Settings)
	})

	// Example 4: DidReceiveSettings with automatic unmarshaling
	streamdeck.OnDidReceiveSettings(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.DidReceiveSettingsPayload[DemoSettings]) error {
		// Update title when settings are received
		title := fmt.Sprintf("%s\n%d", p.Settings.ButtonText, p.Settings.Counter)
		return client.SetTitle(ctx, title, streamdeck.HardwareAndSoftware)
	})

	// Example 5: SendToPlugin with automatic unmarshaling
	streamdeck.OnSendToPlugin(action, func(ctx context.Context, client *streamdeck.Client, payload PropertyInspectorMessage) error {
		// Handle messages from Property Inspector
		switch payload.Action {
		case "reset":
			// Send reset confirmation back to Property Inspector
			return client.SendToPropertyInspector(ctx, ResetCompleteResponse{
				Action: "resetComplete",
			})
		}
		return nil
	})

	log.Println("Generic demo plugin started with type-safe event handlers!")
}
