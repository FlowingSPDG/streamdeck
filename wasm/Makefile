GOROOT := $(shell go env GOROOT)

.DEFAULT_GOAL := build
.PHONY: build serve

build: 
	echo $(GOROOT)
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" ./serve/public
	GOOS=js GOARCH=wasm go build -o ./serve/public/main.wasm ./examples

serve: build
	cd ./serve && go run main.go