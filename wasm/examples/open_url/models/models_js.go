//go:build js

package models

import "syscall/js"

func (s *Settings) GetJSObject() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		return js.ValueOf(map[string]any{
			"url": s.URL,
		})
	})
}
