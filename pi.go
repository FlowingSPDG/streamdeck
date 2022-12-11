package streamdeck

// Settings settings for PI
type Settings interface {
	Initalize()
	IsDefault() bool
}

// PropertyInspector PI settings
type PropertyInspector[T Settings] struct {
	Settings T
}

// PI をactionに紐付けてKeyDownHandlerの際に活用出来るようにする

// NewPropertyInspector Get new PI
func NewPropertyInspector[T Settings](settings T) PropertyInspector[T] {
	return PropertyInspector[T]{
		Settings: settings,
	}
}
