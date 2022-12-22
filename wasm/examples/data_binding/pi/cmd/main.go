package main

import (
	"context"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/FlowingSPDG/streamdeck"
	"github.com/FlowingSPDG/streamdeck/wasm"

	"github.com/FlowingSPDG/streamdeck/wasm/examples/data_binding/models"
)

func main() {
	ctx := context.Background()

	settings := &models.Settings{}

	js.Global().Set("std_oninput", settings.OnInput())

	// Initialize wasm-based PropertyInspector
	SD, err := wasm.InitializePropertyInspector(ctx, settings)
	if err != nil {
		panic(err)
	}
	js.Global().Set("std_setSettings", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("Saving setting:", settings)
		SD.SetSettings(ctx, settings)
		return js.Undefined()
	}))

	// Register "SendToPropertyInspector" handler
	SD.RegisterOnDidReceiveSettingsHandler(ctx, func(e streamdeck.Event) {
		// SetSettingsされたので、PIに変更を反映する
		fmt.Printf("Received event: %#v\n", e)
		fmt.Printf("raw payload: %s\n", e.Payload)

		// Unmarshal payload
		payload := streamdeck.DidReceiveSettingsPayload[models.Settings]{}
		if err := json.Unmarshal(e.Payload, &payload); err != nil {
			msg := fmt.Sprintf("Failed to parse payload: %v", err)
			fmt.Println(msg)
			SD.LogMessage(ctx, msg)
		}
		fmt.Printf("payload settings: %#v\n", payload.Settings)

		// Apply incoming changes
		if err := payload.Settings.ApplyHTML(); err != nil {
			panic(err)
		}
	})

	msg := "PropertyInspector Initialized"
	SD.LogMessage(ctx, msg)
	fmt.Println(msg)

	// TODO: WebSocket接続時に設定を受信したら構造体に反映する

	// Lock thread
	select {}
}
