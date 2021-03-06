// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mgo "gopkg.in/mgo.v2"
import mock "github.com/stretchr/testify/mock"

// DataAccessInterface is an autogenerated mock type for the DataAccessInterface type
type DataAccessInterface struct {
	mock.Mock
}

// C provides a mock function with given fields: collection
func (_m *DataAccessInterface) C(collection string) *mgo.Collection {
	ret := _m.Called(collection)

	var r0 *mgo.Collection
	if rf, ok := ret.Get(0).(func(string) *mgo.Collection); ok {
		r0 = rf(collection)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mgo.Collection)
		}
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *DataAccessInterface) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: collectionName, docs
func (_m *DataAccessInterface) Insert(collectionName string, docs ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, collectionName)
	_ca = append(_ca, docs...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) error); ok {
		r0 = rf(collectionName, docs...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
