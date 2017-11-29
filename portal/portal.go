package portal

import (
	"log"
	"time"

	"github.com/ubiquitous-signage/hamster/multiLanguageString"
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
		log.Println("Upsert portal")
		c.Upsert(
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
		time.Sleep(2 * time.Second)
	}
}
