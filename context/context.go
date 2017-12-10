package context

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
	"github.com/ubiquitous-signage/hamster/util"
)

type Context struct {
	Id  int   `json:"id"`
	Lang     string   `json:"lang"`
}

func PostContext(w rest.ResponseWriter, r *rest.Request) {
	context := Context{}
	err := r.DecodeJsonPayload(&context)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	storeContext(context)
	w.WriteJson(&context)
}

func storeContext(context Context) {
	//session initialize
	mongoEndpoint := viper.GetString("mongo.endopint")
	DBName := viper.GetString("mongo.DBName")
	collectionName := "contexts"
	
	session, collection := util.GetCollections(collectionName, DBName, mongoEndpoint)
	defer session.Close()

	//upsert preparation
	fixedHeader := map[string]interface{}{
		"id": context.Id,
	}
	mgoHeader := bson.M(fixedHeader)

	collection.Upsert(
		mgoHeader,
		context,
	)
}
