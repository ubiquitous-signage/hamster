package trains

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ubiquitous-signage/hamster/multiLanguageString"
	"github.com/ubiquitous-signage/hamster/panel"
	"github.com/ubiquitous-signage/hamster/util"
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

func fetch() (panel.Panel, error) {
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

	lineCharacter := map[string]string{
		"odpt.Railway:TokyoMetro.MarunouchiBranch": "m",
		"odpt.Railway:TokyoMetro.Ginza":            "G",
		"odpt.Railway:TokyoMetro.Marunouchi":       "M",
		"odpt.Railway:TokyoMetro.Hibiya":           "H",
		"odpt.Railway:TokyoMetro.Tozai":            "T",
		"odpt.Railway:TokyoMetro.Chiyoda":          "C",
		"odpt.Railway:TokyoMetro.Yurakucho":        "Y",
		"odpt.Railway:TokyoMetro.Hanzomon":         "Z",
		"odpt.Railway:TokyoMetro.Namboku":          "N",
		"odpt.Railway:TokyoMetro.Fukutoshin":       "F",
	}
	trains := &panel.Panel{
		Contents: []interface{}{},
	}
	trains.Version = 0.0
	trains.Type = "table"
	trains.Title = *multiLanguageString.NewMultiLanguageString("東京メトロ運行情報")
	trains.Category = "external"
	trains.Date = time.Now()
	for _, line := range trainInfomation {
		symbol := *panel.NewImageContent("/static/images/metro/" + lineCharacter[line.OdptRailway] + ".jpg")
		text := *panel.NewStringContent(line.OdptTrainInformationText+" ("+line.DcDate.Format("15:04")+")", true)
		contentLine := []interface{}{symbol, text}
		trains.Contents = append(trains.Contents.([]interface{}), contentLine)
	}
	return *trains, nil
}

func Run() {
	session, collection := util.GetPanel()
	defer session.Close()

	for {
		result, err := fetch()
		if err == nil {
			log.Println("Upsert trains")
			collection.Upsert(
				bson.M{
					"version":  0.0,
					"type":     "table",
					"title.ja": "東京メトロ運行情報",
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
