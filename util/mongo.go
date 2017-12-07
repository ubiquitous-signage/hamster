package util

import (
	"gopkg.in/mgo.v2"
	"github.com/spf13/viper"
)

func Con (endpoint string) *mgo.Session {
	s, err := mgo.Dial(endpoint)
	if err != nil {
		panic(err)
	}
	return s
}

func GetDB (DBName string, mongoEndpoint string) (s *mgo.Session, db *mgo.Database) {
	s = Con(mongoEndpoint)
	db = s.DB(DBName)
	return
}

func GetCollections (collectionName string, DBName string, mongoEndpoint string) (s *mgo.Session, c *mgo.Collection) {
	s, db := GetDB(DBName, mongoEndpoint)
	c = db.C(collectionName)
	return
}

func GetPanel() (s *mgo.Session, c *mgo.Collection) {
	var mongoEndpoint string  = viper.GetString("mongo.endpoint")
	var DBName string         = viper.GetString("mongo.DBName")
	var collectionName string = "panels"

	s, c = GetCollections(collectionName, DBName, mongoEndpoint)
	return
}