package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/FlowingSPDG/streamdeck"
	sdcontext "github.com/FlowingSPDG/streamdeck/context"

	"github.com/FlowingSPDG/streamdeck/wasm/examples/open_url/models"
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
	ac := c.Action("dev.flowingspdg.wasm.openurl")

	// Register WillApeparHandler
	ac.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		client.LogMessage(ctx, "WillAppear on Backend")
		payload := streamdeck.WillAppearPayload[models.Settings]{}
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

	go func() {
		for {
			if !c.IsConnected() {
				continue
			}
			if len(ctxSlice) <= 0 {
				continue
			}
			msg := fmt.Sprintf("Executing %d contexts...", len(ctxSlice))
			c.LogMessage(ctx, msg)
			for _, ctxStr := range ctxSlice {
				ctx := context.Background()
				ctx = sdcontext.WithContext(ctx, ctxStr)

				time.Sleep(time.Second)
				c.SendToPropertyInspector(ctx, &models.Settings{URL: "https://www.elgato.com/"})

				time.Sleep(time.Second)
				c.SendToPropertyInspector(ctx, &models.Settings{URL: "https://go.dev/"})

				time.Sleep(time.Second)
				c.SendToPropertyInspector(ctx, &models.Settings{URL: "https://github.com/FlowingSPDG/streamdeck"})
			}
		}
	}()
	if err := c.Run(ctx); err != nil {
		panic(err)
	}
}
