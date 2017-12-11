package portal

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
					*panel.NewStringContent("研究計画書の提出について", true),
				},
					{
						*panel.NewStringContent("博士コロキウムの実施について", true),
					},
				},
			},
		)
		time.Sleep(2 * time.Second)
	}
}
