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

	SD, err := wasm.InitializePropertyInspector[*models.Settings](ctx)
	if err != nil {
		panic(err)
	}
	SD.RegisterOnSendToPropertyInspectorHandler(ctx, func(e streamdeck.Event) {
		payload := &models.Settings{}
		if err := json.Unmarshal(e.Payload, payload); err != nil {
			msg := fmt.Sprintf("Failed to parse payload: %v", err)
			SD.LogMessage(ctx, msg)
		}
		settings.URL = payload.URL
	})
	SD.LogMessage(ctx, "PropertyInspector Initialized")
	done := make(chan struct{}, 0)
	<-done
}
