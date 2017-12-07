package weather

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ubiquitous-signage/hamster/multiLanguageString"
	"github.com/ubiquitous-signage/hamster/panel"
	"github.com/ubiquitous-signage/hamster/util"
	"gopkg.in/mgo.v2/bson"
)

//気象庁が公開している情報全体をパースするために必要な構造体
//http://www.data.jma.go.jp/developer/xml/feed/regular_l.xml
type WeatherEntries struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Title   string `xml:"title"`
	Id      string `xml:"id"`
	UpDated string `xml:"updated"`
	Link    struct {
		URL string `xml:"href,attr"`
	} `xml:"link"`
	Content string `xml:"content"`
}

//東京都の天気情報を表現したxmlから情報を取り出すための構造体
type Weather struct {
	Body struct {
		MeteorologicalInfos []MeteorologicalInfos `xml:"MeteorologicalInfos"`
	} `xml:"Body"`
}

type MeteorologicalInfos struct {
	Type           string `xml:"type,attr"`
	TimeSeriesInfo struct {
		TimeDefines struct {
			TimeDefine []struct {
				DateTime string `xml:"DateTime"`
				TimeId   int    `xml:"timeId,attr"`
			} `xml:"TimeDefine"`
		} `xml:"TimeDefines"`
		Items []struct {
			Kind []struct {
				Property struct {
					Type        string `xml:"Type"`
					WeatherPart struct {
						Weather []struct {
							RefID   int    `xml:"refID,attr"`
							Type    string `xml:"type,attr"`
							Weather string `xml:",chardata"`
						} `xml:"http://xml.kishou.go.jp/jmaxml1/elementBasis1/ Weather"`
					} `xml:"WeatherPart"`
					TemperaturePart struct {
						Temperature []struct {
							RefID       int    `xml:"refID,attr"`
							Type        string `xml:"type,attr"`
							Unit        string `xml:"unit,attr"`
							Description string `xml:"desctiption,attr"`
							Temperature string `xml:",chardata"`
						} `xml:"http://xml.kishou.go.jp/jmaxml1/elementBasis1/ Temperature"`
					} `xml:"TemperaturePart"`
				} `xml:"Property"`
			} `xml:"Kind"`
			Area struct {
				Name string `xml:"Name"`
			} `xml:"Area"`
		} `xml:"Item"`
	} `xml:"TimeSeriesInfo"`
}

var client = &http.Client{Timeout: 15 * time.Second}

//URLを受け取って、Bodyの部分を[]byteで返す
func httpGet(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	//responseの内容を読み出し
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, err
}

//日本語の天気表現から英字の画像ファイル名に変換する関数
func filename(jaName string) (enName string) {
	switch jaName {
	case "晴れ":
		return "hare.png"
	case "くもり":
		return "kumori.png"
	case "雨":
		return "ame.png"
	case "雪":
		return "yuki.png"
	default:
		return "noimage.png"
	}
}

func fetch() (panel.Panel, error) {
	//気象庁が公開したxmlの情報がまとまっているxmlを入手してパース
	var entries WeatherEntries

	body, err := httpGet("http://www.data.jma.go.jp/developer/xml/feed/regular_l.xml")
	err = xml.Unmarshal(body, &entries)
	if err != nil {
		return panel.Panel{}, err
	}

	//一番新しい東京都の天気情報を選ぶ
	var neededInfo MeteorologicalInfos
	var neededTemp MeteorologicalInfos
	var latestTime time.Time
	var latestEntry Entry

	for _, entry := range entries.Entries {
		if entry.Content == "【東京都府県天気予報】" {
			temp, _ := time.Parse(time.RFC3339, entry.UpDated)
			if temp.After(latestTime) {
				latestTime = temp
				latestEntry = entry
			}
		}
	}

	//一番新しい東京の天気を取りに行き、パースする
	var weatherInformation Weather

	body, err = httpGet(latestEntry.Link.URL)
	err = xml.Unmarshal(body, &weatherInformation)
	if err != nil {
		return panel.Panel{}, err
	}

	//必要なデータが入っているMeteorologicalInfosタグの中身だけを選ぶ(順番は固定とした)
	neededInfo = weatherInformation.Body.MeteorologicalInfos[2]
	neededTemp = weatherInformation.Body.MeteorologicalInfos[3]

	//panel.Panel形式に整形
	weather := &panel.Panel{
		Contents: []interface{}{},
	}
	weather.Version  = 0.0
	weather.Type     = "table"
	weather.Title    = *multiLanguageString.NewMultiLanguageString("東京の天気")
	weather.Category = "external"
	weather.Date     = time.Now()
	// for i, line := range neededInfo.TimeSeriesInfo.Items[0].Kind[0].Property.WeatherPart.Weather {
	// 	symbol := *panel.NewImageContent("/static/images/weather/" + filename(line.Weather))
	// 	text := *panel.NewStringContent(line.Weather)
	// 	temp, _ := time.Parse(time.RFC3339, neededInfo.TimeSeriesInfo.TimeDefines.TimeDefine[line.RefID-1].DateTime)
	// 	date := *panel.NewStringContent(temp.Format("01/02 15:04  "))
	// 	temperature := *panel.NewStringContent(neededTemp.TimeSeriesInfo.Items[0].Kind[0].Property.TemperaturePart.Temperature[i].Temperature + "℃")
	// 	contentLine := []interface{}{date, temperature, symbol, text}
	// 	weather.Contents = append(weather.Contents.([]interface{}), contentLine)
	// }
	times := []interface{}{}
	symbols := []interface{}{}
	temperatures := []interface{}{}
	for i, line := range neededInfo.TimeSeriesInfo.Items[0].Kind[0].Property.WeatherPart.Weather {
		symbol := *panel.NewImageContent("/static/images/weather/" + filename(line.Weather))
		temp, _ := time.Parse(time.RFC3339, neededInfo.TimeSeriesInfo.TimeDefines.TimeDefine[line.RefID-1].DateTime)
		timeString := ""
		if i%2 == 0 {
			timeString = temp.Format("15:04")
		}
		time := *panel.NewStringContent(timeString)
		temperature := *panel.NewStringContent(neededTemp.TimeSeriesInfo.Items[0].Kind[0].Property.TemperaturePart.Temperature[i].Temperature + "℃")
		times = append(times, time)
		symbols = append(symbols, symbol)
		temperatures = append(temperatures, temperature)
	}
	weather.Contents = []interface{}{times, symbols, temperatures}
	return *weather, nil
}

func Run() {
	session, collection := util.GetPanel()
	defer session.Close()

	for {
		result, err := fetch()
		if err == nil {
			log.Println("Upsert weather")
			collection.Upsert(
				bson.M{
					"version":  0.0,
					"type":     "table",
					"title.ja": "東京の天気",
					"category": "external",
				},
				result,
			)
		} else {
			log.Println("Failed to get weather from external server: ", err.Error())
		}

		time.Sleep(time.Hour)
	}
}
