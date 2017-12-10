package wordCloud

import (
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
	"github.com/ubiquitous-signage/hamster/util"
)

type WordCloud struct {
	Words           `json:"words" bson:"contents"`
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
	session, collection := util.GetPanel()
	defer session.Close()

	//update values
	wordCloud := WordCloud{}
	collection.Find(bson.M{"type": "wordCloud"}).One(&wordCloud)

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

	if len(words) > viper.GetInt("wordCloud.thinOutThreshold") {
		words = words.thinOut(countIsOne)
		log.Println("[Word-cloud] words are thinOuted!")
	}

	wordCloud.Words = words

	//upsert preparation
	fixedHeader := map[string]interface{}{
		"version":  0.0,
		"type":     "wordCloud",
		"title":    "Word Cloud",
		"category": "wordCloud",
	}
	mgoHeader := bson.M(fixedHeader)

	wordCloud.Version = fixedHeader["version"].(float64)
	wordCloud.Type = fixedHeader["type"].(string)
	wordCloud.Title = fixedHeader["title"].(string)
	wordCloud.Category = fixedHeader["category"].(string)
	wordCloud.Date = time.Now()

	collection.Upsert(
		mgoHeader,
		wordCloud,
	)
}
