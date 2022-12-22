//go:build js

package models

import (
	"fmt"
	"syscall/js"
)

func (s *Settings) GetJSObject() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		return js.ValueOf(s)
	})
}

// OnChange JS's "onchange" event callback.
func (s *Settings) OnInput() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		// TODO: 受信したeventに対してフィールドを上書きする
		fmt.Println("args:", args)

		id := args[0].String()
		value := args[1]

		// 本当はstd タグを使って自動更新したい...
		switch id {
		case "number":
			s.Number = value.Int()
		case "str":
			s.Str = value.String()
		}

		return js.Undefined()
	})
}

// ApplyChange Apply Field change to HTML elements
func (s *Settings) ApplyHTML() error {
	doc := js.Global().Get("document")

	elemNumber := doc.Call("getElementById", "number")
	elemNumber.Set("value", s.Number)

	elemStr := doc.Call("getElementById", "str")
	elemStr.Set("value", s.Str)

	return nil
}
