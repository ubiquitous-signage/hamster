package schedules

import (
	"log"
	"time"

	"github.com/ubiquitous-signage/hamster/panel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Run() {
	mongoSession, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()

	c := mongoSession.DB("ubiquitous-signage").C("panels")

	for {
		log.Println("Upsert schedules")
		c.Upsert(
			bson.M{
				"version":  0.0,
				"type":     "table",
				"title.ja": "本日の予定",
				"category": "internal",
			},
			panel.Panel{
				PanelHeader: panel.PanelHeader{
					Version:  0.0,
					Type:     "table",
					Title:    *panel.NewMultiLanguageString("本日の予定"),
					Category: "internal",
					Date:     time.Now(),
				},
				Contents: [][]interface{}{{
					*panel.NewStringContent("13:00"),
					*panel.NewStringContent("T-Kernel講習会"),
				}, {
					*panel.NewStringContent("16:00"),
					*panel.NewStringContent("〇〇先生講演会"),
				}},
			},
		)
		time.Sleep(2 * time.Second)
	}
}
