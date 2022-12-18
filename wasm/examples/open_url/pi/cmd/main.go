package main

import (
	"context"

	"github.com/FlowingSPDG/streamdeck/wasm"
)

// Settings PIの設定に使うJSON形式の構造体
type Settings struct {
	URL string `json:"url"`
}

func main() {
	ctx := context.Background()

	SD, err := wasm.InitializePropertyInspector[Settings]()
	if err != nil {
		panic(err)
	}
	SD.LogMessage(ctx, "PropertyInspector Initialized")
	done := make(chan struct{}, 0)
	<-done
}
