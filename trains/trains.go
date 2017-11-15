package trains

import (
	"github.com/ant0ine/go-json-rest/rest"
)

type Trains struct {
	Version  float64     `json:"version"`
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Category string      `json:"category"`
	Contents [][]Content `json:"contents"`
}

type Content struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func GetTrains(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(
		Trains{
			Version:  0.0,
			Type:     "table",
			Title:    "運行情報",
			Category: "external",
			Contents: [][]Content{{
				Content{
					Type:    "Image",
					Payload: "/static/images/metro/M.jpg",
				},
				Content{
					Type:    "String",
					Payload: "全線運転見合わせ",
				},
			}, {
				Content{
					Type:    "Image",
					Payload: "/static/images/metro/C.jpg",
				},
				Content{
					Type:    "String",
					Payload: "ダイヤ乱れ",
				},
			}},
		},
	)
}
