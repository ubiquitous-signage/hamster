package lectures

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/ubiquitous-signage/hamster/multiLanguageString"
	"github.com/ubiquitous-signage/hamster/panel"
	"github.com/ubiquitous-signage/hamster/util"
	"gopkg.in/mgo.v2/bson"
)

func Run() {
	var startSecond = viper.GetDuration("lecture.startDelaySecond")
	var sleepSecond = viper.GetDuration("lecture.sleepSecond")
	time.Sleep(startSecond * time.Second)

	for {
		session, collection := util.GetPanel()
		defer session.Close()

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
					*panel.NewStringContent("総合分析情報学特論V", true),
				}, {
					*panel.NewImageContent("/static/images/lectures/changed.png"),
					*panel.NewStringContent("5限"),
					*panel.NewStringContent("総合分析情報学研究法Ⅱ", true),
					*panel.NewStringContent("306"),
					*panel.NewStringContent("→"),
					*panel.NewStringContent("204"),
				}},
			},
		)
		time.Sleep(sleepSecond * time.Second)
	}
}
