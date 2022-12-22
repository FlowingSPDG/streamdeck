//go:build js && wasm

package wasm

import (
	"syscall/js"
)

// Settings Property Inspector setting interface
type Settings interface {
	IsDefault() bool
	Initialize()
	GetJSObject() js.Func
	ApplyHTML() error // 構造体の値をHTMLに適用する
	OnInput() js.Func // oninputに紐づいて構造体の値を更新する
}
