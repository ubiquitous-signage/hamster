package lectures

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
		log.Println("Upsert lectures")
		collection.Upsert(
			bson.M{
				"version":  0.0,
				"type":     "table",
				"title.ja": "授業情報",
				"category": "internal",
			},
			panel.Panel{
				PanelHeader: panel.PanelHeader{
					Version:  0.0,
					Type:     "table",
					Title:    *multiLanguageString.NewMultiLanguageString("授業情報"),
					Category: "internal",
					Date:     time.Now(),
				},
				Contents: [][]interface{}{{
					*panel.NewImageContent("/static/images/lectures/noclass.png"),
					*panel.NewStringContent("3限"),
					*panel.NewStringContent("総合情報学特論XX"),
				}, {
					*panel.NewImageContent("/static/images/lectures/changed.png"),
					*panel.NewStringContent("4限"),
					*panel.NewStringContent("総合情報学基礎XV"),
					*panel.NewStringContent("301"),
					*panel.NewStringContent("→"),
					*panel.NewStringContent("405"),
				}},
			},
		)
		time.Sleep(2 * time.Second)
	}
}
