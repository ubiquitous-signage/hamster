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
	Type    string                                  `json:"type"`
	Payload multiLanguageString.MultiLanguageString `json:"payload"`
}

func NewStringContent(japanese string) *StringContent {
	return &StringContent{Type: "String", Payload: *multiLanguageString.NewMultiLanguageString(japanese)}
}
