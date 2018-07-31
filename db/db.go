package db

import (
	"gopkg.in/mgo.v2"
	"log"
)

type TripDao struct {
	Server   string
	Database string
}

var MgoSession *mgo.Session

func init() {
	session, err := mgo.Dial("mongo:27017")

	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}
	session.SetMode(mgo.Monotonic, true)
	MgoSession = session
}
