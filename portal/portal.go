package portal

import (
	"github.com/ant0ine/go-json-rest/rest"
)

type Portal struct {
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

func GetPortal(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(
		Portal{
			Version:  0.0,
			Type:     "table",
			Title:    "学府ポータル",
			Category: "internal",
			Contents: [][]Content{{
				Content{
					Type:    "String",
					Payload: "研究計画書の提出について",
				},
			}, {
				Content{
					Type:    "String",
					Payload: "博士コロキウムの実施について",
				},
			}},
		},
	)
}