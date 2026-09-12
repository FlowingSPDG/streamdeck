package streamdeck

import "errors"

// Stream Deck plugin errors.
var (
	ErrMissingPortFlag          = errors.New("missing -port flag")
	ErrMissingPluginUUIDFlag    = errors.New("missing -pluginUUID flag")
	ErrMissingRegisterEventFlag = errors.New("missing -registerEvent flag")
	ErrMissingInfoFlag          = errors.New("missing -info flag")
	ErrConnectionFailed         = errors.New("connection failed")
	ErrWriteFailed              = errors.New("write failed")
	ErrReadFailed               = errors.New("read failed")
	ErrInvalidMessage           = errors.New("invalid message")
	ErrNotConnected             = errors.New("client is not connected")
)
