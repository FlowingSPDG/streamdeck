//go:build js

package models

import (
	"fmt"
	"syscall/js"
)

func (s *Settings) GetJSObject() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		return js.ValueOf(map[string]any{
			"url": s.URL,
		})
	})
}

// OnChange JS's "onchange" event callback.
func (s *Settings) OnInput() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		// 受け取った引数を基にフィールドを上書きする
		fmt.Println("args:", args)

		id := args[0].String()
		value := args[1]

		// 本当はstd タグを使って自動更新したい...
		switch id {
		case "url":
			s.URL = value.String()
		}

		return js.Undefined()
	})
}

func (s *Settings) ApplyHTML() error {
	// 構造体のデータをHTMLに反映する
	doc := js.Global().Get("document")

	elemURL := doc.Call("getElementById", "url")
	elemURL.Set("value", s.URL)

	return nil
}
