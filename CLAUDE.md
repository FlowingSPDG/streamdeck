# CLAUDE.md

This file provides guidance when working with code in this repository.

## Project Overview

Go library for Elgato Stream Deck plugins. The public module path is `github.com/FlowingSPDG/streamdeck/v2`. It implements the Stream Deck 7.1 WebSocket plugin API.

Requires Go 1.27 or later.

## Development Commands

```bash
go test ./...
go test -race ./...
gofmt -l .

cd examples
make build
```

Example plugins live in `examples/` as a separate module (`examples/go.mod`) wired through `go.work`.

## Architecture

- `Client` (`client.go`) owns the WebSocket, pending request map, and per-context event queues.
- `Action[S]` (`action.go`) is a typed handle for one action UUID. Settings type `S` is fixed at registration.
- Incoming events use `Event.Payload json.RawMessage`. Typed handlers call `Event.Unmarshal[T]()`.
- Outgoing messages use the unexported `outgoingEvent` type.
- `dispatch.go` serializes handlers per Stream Deck context so `GetSettings` can block without stalling the read loop.
- `context/context.go` stores action/context/device IDs on `context.Context`.

## Client lifecycle

1. `ParseRegistrationParams(os.Args)`
2. `streamdeck.NewClient(ctx, params, opts...)`
3. `streamdeck.NewAction[Settings](client, uuid)` and register handlers
4. Cancel `ctx` yourself (`signal.NotifyContext`). `Run` does not install signal handlers.
5. `client.Run(ctx)`

## Event handling

Handlers on `Action[S]` receive `ActionEvent[P]`. Property Inspector messages use a method type parameter:

```go
action.OnPropertyInspectorMessage(func(ctx context.Context, m MyMsg) error {
    return action.SendToPropertyInspector(ctx, Reply{})
})
```

Events without an action UUID (`applicationDidLaunch`, `didReceiveDeepLink`, devices, secrets) are registered on `Client`.

Unknown action UUIDs are discarded.

## Settings

- `action.SetSettings(ctx, settings)` / `action.GetSettings(ctx)`
- `client.SetGlobalSettings[G](ctx, g)` / `client.GetGlobalSettings[G](ctx)`
- Getters wait for the matching `id` on the response.

## Images

Typical key size is 72x72. `streamdeck.Image` encodes `image.Image` as a PNG data URL.

## Dependencies

Library module: `github.com/coder/websocket` only.
Example module may pull `gopsutil` and similar tools.
