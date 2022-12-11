package streamdeck

// PropertyInspector PI settings
type PropertyInspector[T any] struct {
	Settings T
}

// PI をactionに紐付けてKeyDownHandlerの際に活用出来るようにする

// NewPropertyInspector Get new PI
func NewPropertyInspector[T any](settings T) PropertyInspector[any] {
	return PropertyInspector[any]{
		Settings: settings,
	}
}
