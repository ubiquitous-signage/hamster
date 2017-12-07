package buildings

import (
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
type BuildingsHeader struct {
	Version  float64   `json:"version"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
}*/

type Buildings struct {
  BuildingsHeader `bson:",inline"`
	Contents `json:"contents"`
}

type BuildingsHeader struct {
	Version  float64   `json:"version"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
}

type Contents []BuildingsContent

type BuildingsContent struct {
	Room        string `json:"room"`
	Light       string `json:"light"`
	Temperature int    `json:"temperature"`
	Humidity    int    `json:"humidity"`
}

/*
type WordCloud struct {
	Words           `json:"words"`
	WordCloudHeader `bson:",inline"`
}

type Word struct {
	Text  string
	Count int
}

type Words []Word

type WordCloudHeader struct {
	Version  float64   `json:"version"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
}
*/

//http://.../api/buildingsへのpostに対して応答する関数
func PostBuildings(w rest.ResponseWriter, r *rest.Request) {

	//研究室から送られてきたものを整形
	buildingsContents := Contents{}
	err := r.DecodeJsonPayload(&buildingsContents)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(buildingsContents) == 0 {
		rest.Error(w, "No data", 400)
	}
	//log.Println("Update buildings")

	//データをmongoDBへ送る
	storeBuildings(buildingsContents)

	//rabbitに対してデータを送信？
	w.WriteJson(&buildingsContents)
}

//データをmongoDBへ送る
func storeBuildings(buildingsContents Contents) {
	//session initialize
	mgoSession, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	c := mgoSession.DB("ubiquitous-signage").C("buildings")
  log.Println("Upsert buildings")

  fixedHeader := map[string]interface{}{
		"version":  0.0,
		"type":     "buildings",
		"title":    "Buildings",
		"category": "buildings",
	}
	mgoHeader := bson.M(fixedHeader)

  buildings := Buildings{}

	buildings.BuildingsHeader.Version = fixedHeader["version"].(float64)
	buildings.BuildingsHeader.Type = fixedHeader["type"].(string)
	buildings.BuildingsHeader.Title = fixedHeader["title"].(string)
	buildings.BuildingsHeader.Category = fixedHeader["category"].(string)
	buildings.BuildingsHeader.Date = time.Now()
  buildings.Contents = buildingsContents

	c.Upsert(
		mgoHeader,
		buildings,
	)
}
