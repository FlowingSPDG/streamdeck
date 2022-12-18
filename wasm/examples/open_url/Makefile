APPNAME=dev.flowingspdg.wasm.sdPlugin
RELEASEDIR = Release

GOROOT := $(shell go env GOROOT)

.DEFAULT_GOAL := build
.PHONY: build serve distribute

# Replacing "MKDIR" command for Windows PowerShell.
MKDIR = mkdir -p
ifeq ($(OS),Windows_NT)
    MKDIR = New-Item -Force -ItemType Directory
endif

# Replacing Distribute command for Windows PowerShell.
DISTRIBUTION_TOOL = ./DistributionTool.exe
ifeq  ($(shell uname),Darwin)
	DISTRIBUTION_TOOL = ./DistributionTool
endif

build: 
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" ./$(APPNAME)
	GOOS=js GOARCH=wasm go build -o ./$(APPNAME)/main.wasm ./cmd

serve: build
	cd ./serve && go run main.go

distribute: build
	@$(RM) ./$(RELEASEDIR)/*
	@$(MKDIR) $(RELEASEDIR)
	$(DISTRIBUTION_TOOL) -b -i $(APPNAME) -o $(RELEASEDIR)