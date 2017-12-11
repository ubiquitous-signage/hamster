package newsletters

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
	var startSecond = viper.GetDuration("newsletter.startDelaySecond")
	var sleepSecond = viper.GetDuration("newsletter.sleepSecond")
	time.Sleep(startSecond * time.Second)

	for {
	session, collection := util.GetPanel()
	defer session.Close()
		log.Println("Upsert newsletters")
		collection.Upsert(
			bson.M{
				"version":  0.0,
				"type":     "table",
				"title.ja": "研究室情報",
				"category": "internal",
			},
			panel.Panel{
				PanelHeader: panel.PanelHeader{
					Version:  0.0,
					Type:     "table",
					Title:    *multiLanguageString.NewMultiLanguageString("研究室情報"),
					Category: "internal",
					Date:     time.Now(),
				},
				Contents: [][]interface{}{{
					*panel.NewStringContent("越塚研究室"),
					*panel.NewStringContent("サイネージ運用開始"),
				}, {
					*panel.NewStringContent("暦本研究室"),
					*panel.NewStringContent("当研究室D１の〇〇くんが〇〇デザインアワード受賞。"),
				}},
			},
		)
		time.Sleep(sleepSecond * time.Second)
	}
}
