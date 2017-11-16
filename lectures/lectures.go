package lectures

import (
	"github.com/ant0ine/go-json-rest/rest"
)

type Lectures struct {
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

func GetLectures(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(
		Lectures{
			Version:  0.0,
			Type:     "table",
			Title:    "授業情報",
			Category: "internal",
			Contents: [][]Content{{
				Content{
					Type:    "Image",
					Payload: "/static/images/lectures/noclass.png",
				},
				Content{
					Type:    "String",
					Payload: "3限",
				},
				Content{
					Type:    "String",
					Payload: "総合情報学特論XX",
				},
			}, {
				Content{
					Type:    "Image",
					Payload: "/static/images/lectures/changed.png",
				},
				Content{
					Type:    "String",
					Payload: "4限",
				},
				Content{
					Type:    "String",
					Payload: "総合情報学基礎XV",
				},
				Content{
					Type:    "String",
					Payload: "301",
				},
				Content{
					Type:    "String",
					Payload: "→",
				},
				Content{
					Type:    "String",
					Payload: "405",
				},
			}},
		},
	)
}
