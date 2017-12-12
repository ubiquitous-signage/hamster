package wordCloud

import (
	"log"
	"net/http"
	"time"
	"math"
	"math/rand"
	"sort"

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
	UpdatedAt time.Time `json:"updated_at"`
	Position Position
}

type Words []Word

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type WordCloudHeader struct {
	Version  float64   `json:"version"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
}

func (w Words) Len() int {
    return len(w)
}

func (w Words) Swap(i, j int) {
    w[i], w[j] = w[j], w[i]
}

// less means early
func (w Words) Less(i, j int) bool {
    return w[i].UpdatedAt.Before(w[j].UpdatedAt)
}

func (sl Words) thinOut(reductionCount int, f func(x Word) bool) []Word {
	result := make([]Word, 0, len(sl))
	for _, word := range sl {
		if reductionCount <= 0 {
			if !f(word){
				word = Reduct(word)
			}
			result = append(result, word)
		} else if !f(word) {
			word = Reduct(word)
			result = append(result, word)
		} else {
			reductionCount = reductionCount - 1
		}
	}
	return result
}

func Reduct(w Word) Word {
	w.Count = int(math.Cbrt(float64(w.Count)))
	return w
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
				words[i].UpdatedAt = time.Now()
				isEmerged = true
				break
			}
		}
		if !isEmerged {
			words = append(words, Word{text, 1, time.Now(), GetPosition()})
		}
	}

	reductionCount := len(words) - viper.GetInt("wordCloud.thinOutThreshold")

	if reductionCount > 0 {
	  sort.Sort(words)
		words = words.thinOut(reductionCount,countIsOne)
		log.Println("[Word-cloud] words are thinOuted!")
	}

	wordCloud.Words = words

	PrintWords(words)

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

func GetPosition() Position {
	x := rand.NormFloat64() * 0.5 + 0.5
	y := rand.NormFloat64() * 0.5 + 0.5
	return Position{x, y}
}

func PrintWords(sl Words){
	for _, i := range sl {
		log.Println(i)
	}
}