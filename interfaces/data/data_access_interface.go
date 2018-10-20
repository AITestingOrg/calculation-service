package data

import mgo "gopkg.in/mgo.v2"

type DataAccessInterface interface {
	C(collection string) *mgo.Collection
	Insert(collectionName string, docs ...interface{}) error
	Close() error
}
