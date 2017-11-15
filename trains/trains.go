package trains

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/ant0ine/go-json-rest/rest"
)

type Trains struct {
	Version  float64     `json:"version"`
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Category string      `json:"category"`
	Contents [][]Content `json:"contents"`
}

type Content struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

var client = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func GetTrains(w rest.ResponseWriter, r *rest.Request) {
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
	getJson(url, &trainInfomation)
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
	trains := new(Trains)
	trains.Version = 0.0
	trains.Type = "table"
	trains.Title = "東京メトロ運行情報"
	trains.Category = "external"
	for _, line := range trainInfomation {
		symbol := Content{Type: "Image", Payload: "/static/images/metro/" + lineCharacter[line.OdptRailway] + ".jpg"}
		text := Content{Type: "String", Payload: line.OdptTrainInformationText}
		date := Content{Type: "String", Payload: "(" + line.DcDate.Format("15:04") + ")"}
		contentLine := []Content{symbol, text, date}
		trains.Contents = append(trains.Contents, contentLine)
	}
	w.WriteJson(trains)
}
