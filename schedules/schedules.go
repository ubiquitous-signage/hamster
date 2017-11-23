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
				"title":    "本日の予定",
				"category": "internal", 
			}, 
			panel.Panel{
				PanelHeader: panel.PanelHeader{
					Version:  0.0,
					Type:     "table",
					Title:    "本日の予定",
					Category: "internal",
					Date:     time.Now(),
				},
				Contents: [][]panel.Content{{
					panel.Content{
						Type:    "String",
						Payload: "13:00",
					},
					panel.Content{
						Type:    "String",
						Payload: "T-Kernel講習会",
					},
				}, {
					panel.Content{
						Type:    "String",
						Payload: "16:00",
					},
					panel.Content{
						Type:    "String",
						Payload: "〇〇先生講演会",
					},
				}},
			},
		)
		time.Sleep(2 * time.Second)
	}
}