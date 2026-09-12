# Settings Manager Example

複数のボタンインスタンスの設定を、context ID をキーにした型付きマップで管理する例です。

## できること

- ボタン文言、色、自動インクリメントの保存
- 2秒ごとの自動カウント
- Property Inspector からの全状態取得と一括リセット

## 使い方

```bash
go build -o settings_manager .
```

ビルドしたバイナリと `manifest.json`、`settings_manager_property_inspector.html` を Stream Deck のプラグインディレクトリへ置きます。
