package data

import (
	"gopkg.in/mgo.v2"
)

// Interface for MongoDB Data Access
type DataAccessInterface interface {
	C(collection string) *mgo.Collection
	Insert(collectionName string, docs ...interface{}) error
	Close() error
}
