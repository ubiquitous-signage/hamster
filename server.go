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
	Contents [][]Message `json:"contents"`
}

type Message struct {
	MessageType    string `json:"message-type"`
	MessagePayload string `json:"message-payload"`
}

func getLectures(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(
		Lectures{
			Version:  0.0,
			Type:     "table",
			Title:    "授業情報",
			Category: "internal",
			Contents: [][]Message{{
				Message{
					MessageType:    "Image",
					MessagePayload: "/resource/img/noclass.png",
				},
				Message{
					MessageType:    "String",
					MessagePayload: "3限",
				},
				Message{
					MessageType:    "String",
					MessagePayload: "総合情報学特論XX",
				},
			},
				{
					Message{
						MessageType:    "Image",
						MessagePayload: "/resource/img/chenged.png",
					},
					Message{
						MessageType:    "String",
						MessagePayload: "4限",
					},
					Message{
						MessageType:    "String",
						MessagePayload: "総合情報学基礎XV",
					},
					Message{
						MessageType:    "String",
						MessagePayload: "301",
					},
					Message{
						MessageType:    "String",
						MessagePayload: "→",
					},
					Message{
						MessageType:    "String",
						MessagePayload: "405",
					},
				}},
		},
	)
}
