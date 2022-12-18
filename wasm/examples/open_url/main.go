package main

import "github.com/FlowingSPDG/streamdeck/wasm"

// Settings PIの設定に使うJSON形式の構造体
type Settings struct {
	Counter int `json:"counter"`
}

func main() {
	wasm.DeclarePropertyInspectorRegistration[Settings]()
	done := make(chan struct{}, 0)
	<-done
}
