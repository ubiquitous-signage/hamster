package server

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ubiquitous-signage/hamster/ads"
	"github.com/ubiquitous-signage/hamster/panel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
		RejectNonCorsRequests:         false,
		AllowedMethods:                []string{"GET"},
		AllowedHeaders:                []string{"Accept", "Authorization", "content-type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
	})

	// 言語設定初期化
	mongoSession.DB("ubiquitous-signage").C("contexts").Upsert(bson.M{"id": 0}, bson.M{"id": 0, "lang": "ja"})

	router, err := rest.MakeRouter(
		rest.Get("/panels", func(w rest.ResponseWriter, r *rest.Request) {
			c := mongoSession.DB("ubiquitous-signage").C("panels")
			result := []panel.Panel{}
			c.Find(nil).All(&result)
			w.WriteJson(result)
		}),
		rest.Get("/contexts", func(w rest.ResponseWriter, r *rest.Request) {
			c := mongoSession.DB("ubiquitous-signage").C("contexts")
			result := struct {
				Lang string `json:"lang"`
			}{}
			c.Find(nil).One(&result)
			w.WriteJson(result)
		}),
		rest.Post("/ads", ads.PostAd),
	)

	if err != nil {
		panic(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(":9000", nil))
}
