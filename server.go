package main

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ubiquitous-signage/hamster/lectures"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/ok", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson("ok")
		}),
		rest.Get("/lectures", lectures.GetLectures),
	)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
