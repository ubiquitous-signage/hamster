package server

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/spf13/viper"
	"github.com/ubiquitous-signage/hamster/ads"
	"github.com/ubiquitous-signage/hamster/panel"
	"github.com/ubiquitous-signage/hamster/rooms"
	"github.com/ubiquitous-signage/hamster/util"
	"github.com/ubiquitous-signage/hamster/wordCloud"
	"gopkg.in/mgo.v2/bson"
)

func Run() {
	//load vars
	var mongoEndpoint = viper.GetString("mongo.endpoint")
	var DBName = viper.GetString("mongo.DBName")
	// var chameleonEndpoint = viper.GetString("chameleon.endpoint")
	// var ubiAdEndpoint = viper.GetString("ubiAd.endpoint")
	mongoSession := util.Con(mongoEndpoint)
	defer mongoSession.Close()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		// OriginValidator: func(origin string, request *rest.Request) bool {
		// 	return origin == chameleonEndpoint || origin == ubiAdEndpoint
		// },
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		RejectNonCorsRequests:         false,
		AllowedMethods:                []string{"GET", "POST"},
		AllowedHeaders:                []string{"Accept", "Authorization", "content-type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
	})

	// 言語設定初期化
	mongoSession.DB(DBName).C("contexts").Upsert(bson.M{"id": 0}, bson.M{"id": 0, "lang": "ja"})

	router, err := rest.MakeRouter(
		rest.Get("/panels", func(w rest.ResponseWriter, r *rest.Request) {
			c := mongoSession.DB(DBName).C("panels")
			result := []panel.Panel{}
			c.Find(nil).All(&result)
			w.WriteJson(result)
		}),
		rest.Get("/contexts", func(w rest.ResponseWriter, r *rest.Request) {
			c := mongoSession.DB(DBName).C("contexts")
			result := struct {
				Lang string `json:"lang"`
				Id   int    `json:"id"`
			}{}
			c.Find(nil).One(&result)
			w.WriteJson(result)
		}),
		rest.Get("/rooms", rooms.GetRooms),
		rest.Post("/ads", ads.PostAd),
		rest.Post("/word-cloud", wordCloud.PostWordCloud),
		rest.Post("/rooms", rooms.PostRooms),
	)

	if err != nil {
		panic(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(":9000", nil))
}
