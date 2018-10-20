package db

import (
	"os"

	"github.com/AITestingOrg/calculation-service/interfaces/data"
	"gopkg.in/mgo.v2"
)

// MongoDAL is an implementation of DataAccessLayer for MongoDB
type MongoDAL struct {
	session *mgo.Session
	dbName  string
}

// NewMongoDAL creates a MongoDAL
func NewMongoDAL(dbName string) (data.DataAccessInterface, error) {
	dbURI := os.Getenv("MONGO_HOST") + ":" + os.Getenv("MONGO_PORT")
	session, err := mgo.Dial(dbURI)
	mongo := &MongoDAL{
		session: session,
		dbName:  dbName,
	}
	return mongo, err
}

// c is a helper method to get a collection from the session
func (m *MongoDAL) c(collection string) *mgo.Collection {
	return m.session.DB(m.dbName).C(collection)
}

// Implements C public helper of the interface
func (m *MongoDAL) C(collection string) *mgo.Collection {
	return m.c(collection)
}

// Insert stores documents in mongo
func (m *MongoDAL) Insert(collectionName string, docs ...interface{}) error {
	return m.c(collectionName).Insert(docs)
}

// Closes Mongo Connection
func (m *MongoDAL) Close() error {
	return m.Close()
}
