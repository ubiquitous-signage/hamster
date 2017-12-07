package server

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ubiquitous-signage/hamster/ads"
	"github.com/ubiquitous-signage/hamster/panel"
	"github.com/ubiquitous-signage/hamster/wordCloud"
	"github.com/ubiquitous-signage/hamster/buildings"
	"github.com/ubiquitous-signage/hamster/util"
	"gopkg.in/mgo.v2/bson"
)

func Run() {
	//load vars
	var mongoEndpoint string     = viper.GetString("mongo.endpoint")
	var DBName string            = viper.GetString("mongo.rootDBName")
	var chameleonEndpoint string = viper.GetString("chameleon.endpoint")
	var ubiAdEndpoint string     = viper.GetString("ubiAd.endpoint")

	mongoSession := util.Con(mongoEndpoint)
	defer mongoSession.Close()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		OriginValidator: func(origin string, request *rest.Request) bool {
			return origin == chameleonEndpoint || origin == ubiAdEndpoint
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
			}{}
			c.Find(nil).One(&result)
			w.WriteJson(result)
		}),
		rest.Post("/ads", ads.PostAd),
		rest.Post("/word-cloud", wordCloud.PostWordCloud),
		rest.Post("/buildings", buildings.PostBuildings),
	)

	if err != nil {
		panic(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(":9000", nil))
}
