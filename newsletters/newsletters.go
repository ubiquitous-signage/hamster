package newsletters

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
		log.Println("Upsert newsletters")
		c.Upsert(
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
					Title:    *panel.NewMultiLanguageString("研究室情報"),
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
		time.Sleep(2 * time.Second)
	}
}
