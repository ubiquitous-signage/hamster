package panel

import (
	"time"
)

type Panel struct {
	PanelHeader `bson:",inline"`
	Contents    interface{} `json:"contents"`
}

type PanelHeader struct {
	Version  float64             `json:"version"`
	Type     string              `json:"type"`
	Title    MultiLanguageString `json:"title"`
	Category string              `json:"category"`
	Date     time.Time           `json:"date"`
}

type ImageContent struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func NewImageContent(payload string) *ImageContent {
	return &ImageContent{Type: "Image", Payload: payload}
}

type StringContent struct {
	Type    string              `json:"type"`
	Payload MultiLanguageString `json:"payload"`
}

func NewStringContent(japanese string) *StringContent {
	return &StringContent{Type: "String", Payload: *NewMultiLanguageString(japanese)}
}

func NewMultiLanguageString(japanese string) *MultiLanguageString {
	return &MultiLanguageString{Ja: japanese}
}

type MultiLanguageString struct {
	Ja string `json:"ja"`
	En string `json:"en"`
	Fr string `json:"fr"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	Ko string `json:"ko"`
}
