package schedules

import (
	"github.com/ant0ine/go-json-rest/rest"
)

type Schedules struct {
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


func GetSchedules(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(
		Schedules{
			Version:  0.0,
			Type:     "table",
			Title:    "本日の予定",
			Category: "internal",
			Contents: [][]Content{{
				Content{
					Type:    "String",
					Payload: "13:00",
				},
				Content{
					Type:    "String",
					Payload: "T-Kernel講習会",
				},
			}, {
				Content{
					Type:    "String",
					Payload: "16:00",
				},
				Content{
					Type:    "String",
					Payload: "〇〇先生講演会",
				},
			}},
		},
	)
}
