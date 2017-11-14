package train

import (
	"github.com/ant0ine/go-json-rest/rest"
)

type Train struct {
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


func GetTrain(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(
		Train{
			Version:  0.0,
			Type:     "table",
			Title:    "運行情報",
			Category: "external",
			Contents: [][]Content{{
				Content{
					Type:    "Image",
					Payload: "/resource/img/marunouchi.png",
				},
				Content{
					Type:    "String",
					Payload: "全線運転見合わせ",
				},
			}, {
				Content{
					Type:    "Image",
					Payload: "/resource/img/chiyoda.png",
				},
				Content{
					Type:    "String",
					Payload: "ダイヤ乱れ",
				},
			}},
		},
	)
}
