package util

import (
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

func Con(endpoint string) *mgo.Session {
	s, err := mgo.Dial(endpoint)
	if err != nil {
		panic(err)
	}
	return s
}

func GetDB(DBName string, mongoEndpoint string) (s *mgo.Session, db *mgo.Database) {
	s = Con(mongoEndpoint)
	db = s.DB(DBName)
	return
}

func GetCollections(collectionName string, DBName string, mongoEndpoint string) (s *mgo.Session, c *mgo.Collection) {
	s, db := GetDB(DBName, mongoEndpoint)
	c = db.C(collectionName)
	return
}

func GetPanel() (s *mgo.Session, c *mgo.Collection) {
	var mongoEndpoint = viper.GetString("mongo.endpoint")
	var DBName = viper.GetString("mongo.DBName")
	var collectionName = "panels"

	s, c = GetCollections(collectionName, DBName, mongoEndpoint)
	return
}
