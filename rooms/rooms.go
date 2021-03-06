package rooms

import (
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Room struct {
	Version   float64   `json:"version"`
	Date      time.Time `json:"date"`
	RoomState `bson:",inline"`
}

type RoomState struct {
	Name        string  `json:"name"`
	Light       string  `json:"light"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

//http://.../api/roomsへのpostに対して応答する関数
func PostRooms(w rest.ResponseWriter, r *rest.Request) {

	//研究室から送られてきたものを整形
	state := RoomState{}
	err := r.DecodeJsonPayload(&state)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// if len(state) == 0 {
	// 	rest.Error(w, "No data", 400)
	// }
	log.Println("Update buildings")

	//データをmongoDBへ送る
	storeRoom(state)

	//clientへレスポンス
	w.WriteJson(&state)
}

//データをmongoDBへ送る
func storeRoom(state RoomState) {
	//session initialize
	mgoSession, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	c := mgoSession.DB("ubiquitous-signage").C("rooms")
	log.Println("Upsert room")

	room := Room{}

	room.Version = 0.0
	room.Date = time.Now()
	room.RoomState = state

	c.Upsert(bson.M{"name": state.Name}, room)
}

func GetRooms(w rest.ResponseWriter, r *rest.Request) {
	mgoSession, err := mgo.Dial("localhost:27017")
	defer mgoSession.Close()

	if err != nil {
		rest.Error(w, "Failed to connect DB", 500)
		return
	}
	c := mgoSession.DB("ubiquitous-signage").C("rooms")
	result := []bson.M{}
	c.Find(nil).All(&result)
	w.WriteJson(result)
}
