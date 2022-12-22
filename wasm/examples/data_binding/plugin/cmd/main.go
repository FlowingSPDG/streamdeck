package main

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/FlowingSPDG/streamdeck"

	"github.com/FlowingSPDG/streamdeck/wasm/examples/data_binding/models"
)

func main() {
	ctx := context.Background()
	// Get Parameters
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		panic(err)
	}

	// Initialize StreamDeck Plugin Client
	c := streamdeck.NewClient(ctx, params)

	// Get thread safe slices
	mutex := sync.Mutex{}
	ctxSlice := make([]string, 0, 1)

	// Register action
	ac := c.Action("dev.flowingspdg.binding.bind")

	// Register WillApeparHandler
	ac.RegisterHandler(streamdeck.DidReceiveSettings, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		client.LogMessage(ctx, "DidReceiveSettings on Backend")
		payload := streamdeck.DidReceiveSettingsPayload[models.Settings]{}
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return err
		}

		// Store event.Context
		mutex.Lock()
		defer mutex.Unlock()
		ctxSlice = append(ctxSlice, event.Context)

		if payload.Settings.IsDefault() {
			payload.Settings.Initialize()
			client.SetSettings(ctx, payload.Settings)
		}
		return nil
	})

	if err := c.Run(ctx); err != nil {
		panic(err)
	}
}
