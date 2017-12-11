package schedules

import (
	"log"
	"time"

	"github.com/ubiquitous-signage/hamster/multiLanguageString"
	"github.com/ubiquitous-signage/hamster/panel"
	"github.com/ubiquitous-signage/hamster/util"
	"gopkg.in/mgo.v2/bson"
)

func Run() {
	session, collection := util.GetPanel()
	defer session.Close()

	for {
		log.Println("Upsert schedules")
		collection.Upsert(
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
					Title:    *multiLanguageString.NewMultiLanguageString("本日の予定"),
					Category: "internal",
					Date:     time.Now(),
				},
				Contents: [][]interface{}{{
					*panel.NewStringContent("13:00"),
					*panel.NewStringContent("T-Kernel講習会", true),
				}, {
					*panel.NewStringContent("16:00"),
					*panel.NewStringContent("〇〇先生講演会", true),
				}},
			},
		)
		time.Sleep(2 * time.Second)
	}
}
