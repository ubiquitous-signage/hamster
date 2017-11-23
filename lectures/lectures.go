package lectures

import (
	"fmt"
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
		fmt.Println("Upsert lectures")
		c.Upsert(
			bson.M{
				"version":  0.0,
				"type":     "table",
				"title":    "授業情報",
				"category": "internal", 
			}, 
			panel.Panel{
				PanelHeader: panel.PanelHeader {
					Version:  0.0,
					Type:     "table",
					Title:    "授業情報",
					Category: "internal",
					Date:     time.Now(),
				},
				Contents: [][]panel.Content{{
					panel.Content{
						Type:    "Image",
						Payload: "/static/images/lectures/noclass.png",
					},
					panel.Content{
						Type:    "String",
						Payload: "3限",
					},
					panel.Content{
						Type:    "String",
						Payload: "総合情報学特論XX",
					},
				}, {
					panel.Content{
						Type:    "Image",
						Payload: "/static/images/lectures/changed.png",
					},
					panel.Content{
						Type:    "String",
						Payload: "4限",
					},
					panel.Content{
						Type:    "String",
						Payload: "総合情報学基礎XV",
					},
					panel.Content{
						Type:    "String",
						Payload: "301",
					},
					panel.Content{
						Type:    "String",
						Payload: "→",
					},
					panel.Content{
						Type:    "String",
						Payload: "405",
					},
				}},
			},
		)
		time.Sleep(2 * time.Second)
	}
}
