package streamdeck

import "errors"

// StreamDeck plugin errors
var (
	ErrMissingPortFlag        = errors.New("missing -port flag")
	ErrMissingPluginUUIDFlag  = errors.New("missing -pluginUUID flag")
	ErrMissingRegisterEventFlag = errors.New("missing -registerEvent flag")
	ErrMissingInfoFlag        = errors.New("missing -info flag")
	ErrJSONMarshal            = errors.New("JSON marshal error")
	ErrSettingsNotFound       = errors.New("couldn't find settings for context")
	ErrConnectionFailed       = errors.New("connection failed")
	ErrWriteFailed            = errors.New("write failed")
	ErrReadFailed             = errors.New("read failed")
	ErrInvalidMessage         = errors.New("invalid message")
)
