module github.com/FlowingSPDG/streamdeck/wasm

go 1.19

require (
	github.com/FlowingSPDG/streamdeck v0.0.0-20221216133737-3e9f8326b4bf
	nhooyr.io/websocket v1.8.7
)

require (
	github.com/klauspost/compress v1.10.3 // indirect
	golang.org/x/sync v0.1.0 // indirect
)

// streamdeck package should be relative path
replace github.com/FlowingSPDG/streamdeck => ../
