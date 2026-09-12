# StreamDeck Go SDK

Elgato Stream Deck プラグインを Go で書くための SDK です。Stream Deck 7.1 時点の WebSocket プラグイン API に合わせています。

必要環境: **Go 1.27 以降**。メソッド固有の型パラメータを使うためです。

## 特徴

- アクションごとに設定型を束ねる `Action[S]`
- Property Inspector やグローバル設定はジェネリックメソッドで型付け
- `id` による `GetSettings` / `GetGlobalSettings` / `GetResources` の応答待ち
- Stream Deck + の dial / touch / setFeedback
- `didReceiveResources` / `didReceiveSecrets` を含む SDK 7.1 の受信イベント

仕様の参照: [Plugin WebSocket API](https://docs.elgato.com/streamdeck/sdk/references/websocket/plugin/)（突き合わせ日: 2026-09-12）。

## 使い方

```go
package main

import (
	"context"
	"os"
	"strconv"

	"github.com/FlowingSPDG/streamdeck/v2"
)

type Settings struct {
	Counter int `json:"counter"`
}

func main() {
	ctx := context.Background()

	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		panic(err)
	}

	client := streamdeck.NewClient(ctx, params)
	action := streamdeck.NewAction[Settings](client, "com.example.counter")

	action.OnKeyDown(func(ctx context.Context, e streamdeck.KeyDownEvent[Settings]) error {
		e.Payload.Settings.Counter++
		if err := action.SetSettings(ctx, e.Payload.Settings); err != nil {
			return err
		}
		return action.SetTitle(ctx, strconv.Itoa(e.Payload.Settings.Counter), streamdeck.HardwareAndSoftware)
	})

	if err := client.Run(ctx); err != nil {
		panic(err)
	}
}
```

Property Inspector のメッセージ型は設定型と別に渡せます。

```go
type PIMessage struct {
	Action string `json:"action"`
}

action.OnPropertyInspectorMessage(func(ctx context.Context, m PIMessage) error {
	return action.SendToPropertyInspector(ctx, map[string]string{"action": "ok"})
})
```

## 例

`examples/` を見てください。

- `counter` — キー押下でカウント
- `cpu` — CPU 使用率のグラフ
- `settings_manager` — 複数インスタンスの設定管理
- `plus` — Stream Deck + のダイヤルと `setFeedback`

```bash
cd examples
make build
```

## 移行

v1 からの破壊的変更は [MIGRATION.md](MIGRATION.md) にまとめています。

## License

MIT License
