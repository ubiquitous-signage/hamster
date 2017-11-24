package trains

import (
	"log"
	"encoding/json"
	"net/http"	
	"time"
	"github.com/ubiquitous-signage/hamster/panel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var client = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func fetch() (panel.Panel, error){
	var trainInfomation []struct {
		Context                  string    `json:"@context"`
		ID                       string    `json:"@id"`
		DcDate                   time.Time `json:"dc:date"`
		DctValid                 time.Time `json:"dct:valid"`
		OdptOperator             string    `json:"odpt:operator"`
		OdptRailway              string    `json:"odpt:railway"`
		OdptTimeOfOrigin         time.Time `json:"odpt:timeOfOrigin"`
		OdptTrainInformationText string    `json:"odpt:trainInformationText"`
		Type                     string    `json:"@type"`
	}
	
	// accessToken := "東京メトロAPIアクセストークン"
	//     <- 別ファイルで用意
	
	url := "https://api.tokyometroapp.jp/api/v2/datapoints?rdf:type=odpt:TrainInformation&acl:consumerKey=" + accessToken
	r, err := client.Get(url)

	if err != nil {
		return panel.Panel{}, err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&trainInfomation)

	lineCharacter := map[string] string {
		"odpt.Railway:TokyoMetro.MarunouchiBranch": "m", 
		"odpt.Railway:TokyoMetro.Ginza": "G", 
		"odpt.Railway:TokyoMetro.Marunouchi": "M", 
		"odpt.Railway:TokyoMetro.Hibiya": "H", 
		"odpt.Railway:TokyoMetro.Tozai": "T", 
		"odpt.Railway:TokyoMetro.Chiyoda": "C", 
		"odpt.Railway:TokyoMetro.Yurakucho": "Y", 
		"odpt.Railway:TokyoMetro.Hanzomon": "Z", 
		"odpt.Railway:TokyoMetro.Namboku": "N", 
		"odpt.Railway:TokyoMetro.Fukutoshin": "F",
	}
	trains := &panel.Panel{
		Contents: []interface{}{},
	}
	trains.Version = 0.0
	trains.Type = "table"
	trains.Title = "東京メトロ運行情報"
	trains.Category = "external"
	trains.Date = time.Now()
	for _, line := range trainInfomation {
		symbol := panel.Content{Type: "Image", Payload: "/static/images/metro/" + lineCharacter[line.OdptRailway] + ".jpg"}
		text := panel.Content{Type: "String", Payload: line.OdptTrainInformationText}
		date := panel.Content{Type: "String", Payload: "(" + line.DcDate.Format("15:04") + ")"}
		contentLine := []panel.Content{symbol, text, date}
		trains.Contents = append(trains.Contents.([]interface{}), contentLine)
	}
	return *trains, nil
}

func Run() {
	mongoSession, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()

	c := mongoSession.DB("ubiquitous-signage").C("panels")

	for {
		result, err := fetch()
		if err == nil {
			log.Println("Upsert trains")
			c.Upsert(
				bson.M{
					"version":  0.0,
					"type":     "table",
					"title":    "東京メトロ運行情報",
					"category": "external", 
				}, 
				result,
			)
		} else {
			log.Println("Failed to get trains from external server: ", err.Error())
		}
		
		time.Sleep(60 * time.Second)
	}
}
