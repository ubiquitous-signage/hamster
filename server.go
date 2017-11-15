package main

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ubiquitous-signage/hamster/lectures"
	"github.com/ubiquitous-signage/hamster/newsletters"
	"github.com/ubiquitous-signage/hamster/portal"
	"github.com/ubiquitous-signage/hamster/schedules"
	"github.com/ubiquitous-signage/hamster/trains"
	"github.com/ubiquitous-signage/hamster/weather"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		OriginValidator: func(origin string, request *rest.Request) bool {
            return origin == "http://localhost:8080"
        },
   RejectNonCorsRequests: false,
	 AllowedMethods: []string{"GET", "POST", "PUT"},
	 AllowedHeaders: []string{"Accept","Authorization", "content-type", "X-Custom-Header", "Origin"},
	 AccessControlAllowCredentials: true,
 })

	router, err := rest.MakeRouter(
		rest.Get("/ok", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson("ok")
		}),
		rest.Get("/lectures", lectures.GetLectures),
		rest.Get("/schedules", schedules.GetSchedules),
		rest.Get("/trains", trains.GetTrains),
		rest.Get("/portal", portal.GetPortal),
		rest.Get("/newsletters", newsletters.GetNewsletters),
		rest.Get("/weather", weather.GetWeather),
	)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":9000", api.MakeHandler()))
}
