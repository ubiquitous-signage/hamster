package portal

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
	var startSecond = viper.GetDuration("portal.startDelaySecond")
	var sleepSecond = viper.GetDuration("portal.sleepSecond")
	time.Sleep(startSecond * time.Second)

	for {
	session, collection := util.GetPanel()
	defer session.Close()
		log.Println("Upsert portal")
		collection.Upsert(
			bson.M{
				"version":  0.0,
				"type":     "table",
				"title.ja": "学府ポータル",
				"category": "internal",
			},
			panel.Panel{
				PanelHeader: panel.PanelHeader{
					Version:  0.0,
					Type:     "table",
					Title:    *multiLanguageString.NewMultiLanguageString("学府ポータル"),
					Category: "internal",
					Date:     time.Now(),
				},
				Contents: [][]interface{}{{
					*panel.NewStringContent("研究計画書の提出について"),
				},
					{
						*panel.NewStringContent("博士コロキウムの実施について"),
					},
				},
			},
		)
		time.Sleep(sleepSecond * time.Second)
	}
}
