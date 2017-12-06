package wordCloud

import (
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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

func (sl Words) thinOut(f func(x Word) bool) []Word {
	result := make([]Word, 0, len(sl))
	for _, word := range sl {
		if !f(word) {
			word.Count = word.Count - 1
			result = append(result, word)
		}
	}
	return result
}

func countIsOne(word Word) bool {
	return word.Count == 1
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
	wordCloud := WordCloud{}
	c.Find(nil).One(&wordCloud)

	words := wordCloud.Words
	newWords := newWordCloud.Words

	for ni := 0; ni < len(newWords); ni++ {
		text := newWords[ni].Text

		isEmerged := false
		for i := 0; i < len(words); i++ {
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

	if len(words) > 40 {
		log.Println(words)
		words = words.thinOut(countIsOne)
		log.Println("words are thinOuted!")
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

	wordCloud.Version = fixedHeader["version"].(float64)
	wordCloud.Type = fixedHeader["type"].(string)
	wordCloud.Title = fixedHeader["title"].(string)
	wordCloud.Category = fixedHeader["category"].(string)
	wordCloud.Date = time.Now()

	c.Upsert(
		mgoHeader,
		wordCloud,
	)
}
