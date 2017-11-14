package newsletters

import (
	"github.com/ant0ine/go-json-rest/rest"
)

type Newsletters struct {
	Version  float64 `json:"version"`
	Type     string  `json:"type"`
	Title    string  `json:"title"`
	Category string  `json:"category"`
	Contents [][]Content `json:"contents"`
}

type Content struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}


func GetNewsletters(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(
		Newsletters{
			Version:  0.0,
			Type:     "table",
			Title:    "研究室情報",
			Category: "internal",
			Contents: [][]Content{{
				Content{
					Type:    "String",
					Payload: "越塚研究室",
				},
				Content{
					Type:    "String",
					Payload: "サイネージ運用開始",
				},
			}, {
				Content{
					Type:    "String",
					Payload: "暦本研究室",
				},
				Content{
					Type:    "String",
					Payload: "当研究室D１の〇〇くんが〇〇デザインアワード受賞。",
				},
			}},
		},
	)
}
