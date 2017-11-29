package ads

import (
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ubiquitous-signage/hamster/multiLanguageString"
	"gopkg.in/mgo.v2"
)

type InputAd struct {
	Title    string `json:"title"`
	Contents string `json:"contents"`
}

type OutputAd struct {
	Version  float64                                 `bson:"version"`
	Type     string                                  `bson:"type"`
	Title    multiLanguageString.MultiLanguageString `bson:"title"`
	Category string                                  `bson:"category"`
	Date     time.Time                               `bson:"date"`
	Contents multiLanguageString.MultiLanguageString `bson:"contents"`
}

type Context struct {
	Success bool `json:"success"`
}

// POST
func PostAd(w rest.ResponseWriter, r *rest.Request) {
	// Receive JSON
	inputAd := InputAd{}
	err := r.DecodeJsonPayload(&inputAd)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if inputAd.Title == "" {
		rest.Error(w, "title required", 400)
		return
	}
	if inputAd.Contents == "" {
		rest.Error(w, "contents required", 400)
		return
	}

	// Insert Ad into MongoDB
	outputAd := &OutputAd{
		Version:  0.0,
		Type:     "plain",
		Title:    *multiLanguageString.NewMultiLanguageString(inputAd.Title),
		Category: "ad",
		Date:     time.Now(),
		Contents: *multiLanguageString.NewMultiLanguageString(inputAd.Contents),
	}
	err = insert(outputAd)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// JSON response
	context := &Context{
		Success: true,
	}
	w.WriteJson(context)
}

func insert(outputAd *OutputAd) error {
	mongoSession, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()

	c := mongoSession.DB("ubiquitous-signage").C("panels")

	return c.Insert(outputAd)
}
