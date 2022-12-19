package wasm

import "syscall/js"

// Settings Property Inspector setting interface
type Settings interface {
	IsDefault() bool
	Initialize()
	GetJSObject() js.Func
}
