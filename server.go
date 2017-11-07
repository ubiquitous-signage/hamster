package main

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/ok", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson("ok")
		}),
		rest.Get("/lectures", getLectures),
	)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

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

func getLectures(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(
		Lectures{
			Version:  0.0,
			Type:     "table",
			Title:    "授業情報",
			Category: "internal",
			Contents: [][]Content{{
				Content{
					Type:    "Image",
					Payload: "/resource/img/noclass.png",
				},
				Content{
					Type:    "String",
					Payload: "3限",
				},
				Content{
					Type:    "String",
					Payload: "総合情報学特論XX",
				},
			},
				{
					Content{
						Type:    "Image",
						Payload: "/resource/img/chenged.png",
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
