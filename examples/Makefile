GO ?= go
EXAMPLES := counter cpu settings_manager plus

.PHONY: all build test tidy logs $(EXAMPLES)

all: build

build: $(EXAMPLES)

$(EXAMPLES):
	$(GO) build -C $@ -o $@ .

test:
	$(GO) test -C .. ./...

tidy:
	$(GO) mod tidy

logs:
	@echo Stream Deck plugin logs are written by the Stream Deck application.
	@echo On Windows see %APPDATA%\Elgato\StreamDeck\logs
	@echo On macOS see ~/Library/Logs/ElgatoStreamDeck
