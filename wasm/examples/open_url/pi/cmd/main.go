package main

import (
	"context"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/FlowingSPDG/streamdeck"
	"github.com/FlowingSPDG/streamdeck/wasm"

	"github.com/FlowingSPDG/streamdeck/wasm/examples/open_url/models"
)

func main() {
	settings := &models.Settings{}
	// JS側へ露出する
	js.Global().Set("get_settings", settings.GetJSObject())

	ctx := context.Background()

	// Initialize wasm-based PropertyInspector
	SD, err := wasm.InitializePropertyInspector(ctx, settings)
	if err != nil {
		panic(err)
	}

	// Register "SendToPropertyInspector" handler
	SD.RegisterOnSendToPropertyInspectorHandler(ctx, func(e streamdeck.Event) {
		fmt.Printf("Received event [%#v]\n", e)

		// Unmarshal payload
		payload := &models.Settings{}
		if err := json.Unmarshal(e.Payload, payload); err != nil {
			msg := fmt.Sprintf("Failed to parse payload: %v", err)
			fmt.Println(msg)
			SD.LogMessage(ctx, msg)
		}
		settings.URL = payload.URL
	})
	msg := "PropertyInspector Initialized"
	SD.LogMessage(ctx, msg)
	fmt.Println(msg)

	// Lock thread
	select {}
}
