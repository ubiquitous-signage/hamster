package server

import (
	"log"
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ubiquitous-signage/hamster/panel"
	"github.com/ubiquitous-signage/hamster/wordCloud"
)

func Run() {
	mongoSession, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		OriginValidator: func(origin string, request *rest.Request) bool {
            return origin == "http://localhost:8080"
        },
   		RejectNonCorsRequests: false,
	 	AllowedMethods: []string{"GET"},
	 	AllowedHeaders: []string{"Accept","Authorization", "content-type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
	 })

	router, err := rest.MakeRouter(
		rest.Get("/panels", func(w rest.ResponseWriter, r *rest.Request) {
			c := mongoSession.DB("ubiquitous-signage").C("panels")
			result := []panel.Panel{}
			c.Find(nil).All(&result)
			
			w.WriteJson(result)
		}),
		rest.Post("/word-cloud", wordCloud.PostWordCloud),
	)

	if err != nil {
		panic(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(":9000", nil))
}
