package models

import (
	"reflect"
)

// Settings PIの設定に使うJSON形式の構造体
// PI/Plugin 両方で共通したフォーマット
type Settings struct {
	Number int    `json:"number"`
	Str    string `json:"str"`
}

func (s Settings) IsDefault() bool {
	return reflect.ValueOf(s).IsZero()
}
func (s *Settings) Initialize() {
	s.Number = 255
	s.Str = "Hello world"
}
