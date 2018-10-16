package db

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

type Session struct {
	MgoSession *mgo.Session
}

func Connect() (*mgo.Session, error) {
	mongo := os.Getenv("MONGO_HOST") + ":" + os.Getenv("MONGO_PORT")
	session, err := mgo.Dial(mongo)

	if err != nil {
		return nil, err
	}

	if session == nil {
		log.Fatalf("No session")
	}
	return session, nil
}

func (s *Session) Close() {
	if s.MgoSession != nil {
		s.MgoSession.Close()
	}
}
