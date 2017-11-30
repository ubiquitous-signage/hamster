
package wordCloud

import (
	"log"
	"net/http"	
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/ant0ine/go-json-rest/rest"
)

type WordCloud struct{
	Words []Word `json:"words"`
	WordCloudHeader `bson:",inline"`
}

type Word struct{
	Text  string
	Count int
}

type WordCloudHeader struct {
	Version  float64     `json:"version"`
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Category string      `json:"category"`
	Date	   time.Time	 `json:"date"`
 }

func PostWordCloud(w rest.ResponseWriter, r *rest.Request) {
	wordCloud := WordCloud{}
	err := r.DecodeJsonPayload(&wordCloud)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(wordCloud.Words) == 0 {
		rest.Error(w, "words required", 400)
	}
	log.Println("input wordCloud: ", wordCloud.Words)

	storeWordCloud(wordCloud)
	w.WriteJson(&wordCloud)
}

func storeWordCloud(newWordCloud WordCloud) {
	//session initialize
	mgoSession, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	c := mgoSession.DB("ubiquitous-signage").C("wordCloud")

	//update values
	wordClouds := []WordCloud{}
	c.Find(nil).All(&wordClouds)
	wordCloud := wordClouds[0]

	words := wordCloud.Words
	newWords := newWordCloud.Words

	for ni:= 0; ni < len(newWords); ni++ {
		text := newWords[ni].Text

		isEmerged := false
		for i:= 0; i < len(words); i++ {
			// word := words[i]
			// if text == word.Text {
			// 	word.Count = word.Count + 1
			// 	isEmerged = true
			// 	break
			// }  
			// word:= words[i]のwordの参照はwords[i]にはない？
			// スライスとマップは参照型のはずなのだが
			if text == words[i].Text {
				words[i].Count = words[i].Count + 1
				isEmerged = true
				break
			} 
		}
		if !isEmerged {
			words = append(words, Word{text, 1})
		}
	}

	wordCloud.Words = words

	log.Println(words)

	//upsert preparation
	fixedHeader := map[string]interface{}{
		"version":  0.0,
		"type":     "wordCloud",
		"title":    "Word Cloud",
		"category": "external", 
	}
	mgoHeader := bson.M(fixedHeader)

	wordCloud.Version  = fixedHeader["version"].(float64)
	wordCloud.Type     = fixedHeader["type"].(string)
	wordCloud.Title    = fixedHeader["title"].(string)
	wordCloud.Category = fixedHeader["category"].(string)
	wordCloud.Date     = time.Now()

	c.Upsert(
		mgoHeader,
		wordCloud,
	)
}