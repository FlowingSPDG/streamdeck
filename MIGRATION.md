# v1 から v2 への移行

モジュールパスが変わります。

```
github.com/FlowingSPDG/streamdeck
→ github.com/FlowingSPDG/streamdeck/v2
```

Go 1.27 以降が必要です。

## アクションの作り方

```go
// v1
action := client.Action("com.example.counter")
streamdeck.OnKeyDown(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
    return client.SetSettings(ctx, p.Settings)
})

// v2
action := streamdeck.NewAction[Settings](client, "com.example.counter")
action.OnKeyDown(func(ctx context.Context, e streamdeck.KeyDownEvent[Settings]) error {
    return action.SetSettings(ctx, e.Payload.Settings)
})
```

## シグナル

`Client.Run` は SIGINT を捕捉しません。呼び出し側で `signal.NotifyContext` を使ってください。キャンセルされると `Run` は `ctx.Err()` を返します。

## 名前が変わったもの

| v1 | v2 |
| --- | --- |
| `client.Action(uuid)` | `streamdeck.NewAction[S](client, uuid)` |
| `OnKeyDown[T](action, h)` ほかのパッケージ関数 | `action.OnKeyDown(h)` |
| `client.SetSettings(ctx, s)` | `action.SetSettings(ctx, s)` |
| `event.UnmarshalPayload(&p)` | `event.Unmarshal[T]()` |
| `DeviceInfo.DeviceName` (`json:"deviceName"`) | `DeviceInfo.Name` (`json:"name"`) |
| `OnDidReceivePropertyInspectorMessage` | `OnPropertyInspectorMessage`（wire 名は `sendToPlugin`） |
| `client.SendToPlugin` | 削除。PI へ送るのは `SendToPropertyInspector` |
| `streamdeck.Log()` | `WithLogger(*slog.Logger)` |

## 取得系 API

`GetSettings` / `GetGlobalSettings` / `GetResources` は投げっぱなしではなく、`id` で応答を待ちます。

```go
settings, err := action.GetSettings(ctx)
```

Property Inspector 側の変更通知は、これまでどおり `OnDidReceiveSettings` に届きます。
