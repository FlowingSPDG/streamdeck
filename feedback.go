package streamdeck

import "encoding/json"

// Feedback updates Stream Deck + layout items by id.
type Feedback map[string]FeedbackValue

// FeedbackValue is a setFeedback item: a string, number, or typed layout object.
type FeedbackValue interface {
	isFeedbackValue()
}

// FeedbackString is a shorthand string layout value.
type FeedbackString string

func (FeedbackString) isFeedbackValue() {}

// FeedbackNumber is a shorthand numeric layout value.
type FeedbackNumber float64

func (FeedbackNumber) isFeedbackValue() {}

// FeedbackRange is a min/max pair for bar and gauge items.
type FeedbackRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// BarFeedback is a bar layout item.
type BarFeedback struct {
	Background string         `json:"background,omitempty"`
	BarBgC     string         `json:"bar_bg_c,omitempty"`
	BarBorderC string         `json:"bar_border_c,omitempty"`
	BarFillC   string         `json:"bar_fill_c,omitempty"`
	BorderW    *int           `json:"border_w,omitempty"`
	Enabled    *bool          `json:"enabled,omitempty"`
	Opacity    *float64       `json:"opacity,omitempty"`
	Range      *FeedbackRange `json:"range,omitempty"`
	Subtype    *int           `json:"subtype,omitempty"`
	Value      *float64       `json:"value,omitempty"`
	ZOrder     *int           `json:"zOrder,omitempty"`
}

func (BarFeedback) isFeedbackValue() {}

// GaugeBarFeedback is a gauge-bar layout item.
type GaugeBarFeedback struct {
	Background string         `json:"background,omitempty"`
	BarBgC     string         `json:"bar_bg_c,omitempty"`
	BarBorderC string         `json:"bar_border_c,omitempty"`
	BarFillC   string         `json:"bar_fill_c,omitempty"`
	BarH       *int           `json:"bar_h,omitempty"`
	BorderW    *int           `json:"border_w,omitempty"`
	Enabled    *bool          `json:"enabled,omitempty"`
	Opacity    *float64       `json:"opacity,omitempty"`
	Range      *FeedbackRange `json:"range,omitempty"`
	Subtype    *int           `json:"subtype,omitempty"`
	Value      *float64       `json:"value,omitempty"`
	ZOrder     *int           `json:"zOrder,omitempty"`
}

func (GaugeBarFeedback) isFeedbackValue() {}

// PixmapFeedback is an image layout item.
type PixmapFeedback struct {
	Background string   `json:"background,omitempty"`
	Enabled    *bool    `json:"enabled,omitempty"`
	Opacity    *float64 `json:"opacity,omitempty"`
	Value      string   `json:"value,omitempty"`
	ZOrder     *int     `json:"zOrder,omitempty"`
}

func (PixmapFeedback) isFeedbackValue() {}

// TextAlignment is a setFeedback text alignment.
type TextAlignment string

const (
	TextAlignCenter TextAlignment = "center"
	TextAlignLeft   TextAlignment = "left"
	TextAlignRight  TextAlignment = "right"
)

// TextOverflow is a setFeedback text overflow mode.
type TextOverflow string

const (
	TextOverflowClip     TextOverflow = "clip"
	TextOverflowEllipsis TextOverflow = "ellipsis"
	TextOverflowFade     TextOverflow = "fade"
)

// FeedbackFont is a text item font.
type FeedbackFont struct {
	Size   *int `json:"size,omitempty"`
	Weight *int `json:"weight,omitempty"`
}

// TextFeedback is a text layout item.
type TextFeedback struct {
	Alignment    TextAlignment `json:"alignment,omitempty"`
	Background   string        `json:"background,omitempty"`
	Color        string        `json:"color,omitempty"`
	Enabled      *bool         `json:"enabled,omitempty"`
	Font         *FeedbackFont `json:"font,omitempty"`
	Opacity      *float64      `json:"opacity,omitempty"`
	TextOverflow TextOverflow  `json:"text-overflow,omitempty"`
	Value        string        `json:"value,omitempty"`
	ZOrder       *int          `json:"zOrder,omitempty"`
}

func (TextFeedback) isFeedbackValue() {}

func (f Feedback) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]FeedbackValue(f))
}
