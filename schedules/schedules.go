package schedules

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
	var startSecond = viper.GetDuration("schedule.startDelaySecond")
	var sleepSecond = viper.GetDuration("schedule.sleepSecond")
	time.Sleep(startSecond * time.Second)

	for {
		session, collection := util.GetPanel()
		defer session.Close()
		log.Println("Upsert schedules")
		collection.Upsert(
			bson.M{
				"version":  0.0,
				"type":     "table",
				"title.ja": "イベント",
				"category": "internal",
			},
			panel.Panel{
				PanelHeader: panel.PanelHeader{
					Version:  0.0,
					Type:     "table",
					Title:    *multiLanguageString.NewMultiLanguageString("イベント"),
					Category: "internal",
					Date:     time.Now(),
				},
				Contents: [][]interface{}{{
					*panel.NewStringContent("12/16"),
					*panel.NewStringContent("第2回　メディアと表現について考えるシンポジウム -「徹底検証　炎上リスク―そのジェンダー表現はアリか」", true),
				}, {
					*panel.NewStringContent("12/15"),
					*panel.NewStringContent("総合分析情報学コース 冬期入試出願締切", true),
				}},
			},
		)
		time.Sleep(sleepSecond * time.Second)
	}
}
