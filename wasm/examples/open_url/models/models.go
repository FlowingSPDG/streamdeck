package models

import (
	"reflect"
)

// Settings PIの設定に使うJSON形式の構造体
// PI/Plugin 両方で共通したフォーマット
type Settings struct {
	URL string `json:"url"`
}

func (s Settings) IsDefault() bool {
	return reflect.ValueOf(s).IsZero()
}
func (s *Settings) Initialize() {
	s.URL = "https://www.elgato.com/"
}
