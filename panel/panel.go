package panel

import (
	"time"

	"github.com/ubiquitous-signage/hamster/multiLanguageString"
)

type Panel struct {
	PanelHeader `bson:",inline"`
	Contents    interface{} `json:"contents"`
}

type PanelHeader struct {
	Version  float64                                 `json:"version"`
	Type     string                                  `json:"type"`
	Title    multiLanguageString.MultiLanguageString `json:"title"`
	Category string                                  `json:"category"`
	Date     time.Time                               `json:"date"`
}

type ImageContent struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func NewImageContent(payload string) *ImageContent {
	return &ImageContent{Type: "Image", Payload: payload}
}

type StringContent struct {
	Type       string                                  `json:"type"`
	Payload    multiLanguageString.MultiLanguageString `json:"payload"`
	Scrollable bool                                    `json:"scrollable"`
}

func NewStringContent(japanese string, options ...bool) *StringContent {
	var scrollable bool
	if len(options) > 0 {
		scrollable = options[0]
	} else {
		scrollable = false
	}
	return &StringContent{
		Type:       "String",
		Payload:    *multiLanguageString.NewMultiLanguageString(japanese),
		Scrollable: scrollable,
	}
}
