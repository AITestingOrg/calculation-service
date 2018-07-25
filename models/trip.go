package models

import (
	"regexp"
	"fmt"
	"errors"
)

type Trip struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureTime int64  `json:"departureTime"`
	UserId   	string `json:"userId"`
}

func (trip Trip) ValidateFields(fields ...string) error{
	invalidFields := ""
	for _, field := range fields{
		switch field {
		case "originAddress":
			if !trip.ValidateOrigin(){
				invalidFields += fmt.Sprintf("Invalid originAddress.\n\tGiven: %s\n\tExpected: Non Empty String\n", trip.Origin)
			}
		case "destinationAddress":
			if !trip.ValidateDestination(){
				invalidFields += fmt.Sprintf("Invalid destinationAddress.\n\tGiven: %s\n\tExpected: Non Empty String\n", trip.Destination)
			}
		case "departureTime":
			if !trip.ValidateDepartureTime(){
				invalidFields += fmt.Sprintf("Invalid departureTime.\n\tGiven: %s\n\tExpected: Valid int64 >= 0\n", trip.DepartureTime)
			}
		case "userId":
			if !trip.ValidateUserId(){
				invalidFields += fmt.Sprintf("Invalid userId.\n\tGiven: %s\n\tExpected: Valid UUID in version 4 format\n", trip.UserId)
			}
		}
	}
	if invalidFields != "" {
		return errors.New(invalidFields)
	}
	return nil
}

func (trip Trip) ValidateOrigin() bool {
	if trip.Origin != "" {
		return true
	}
	return false
}

func (trip Trip) ValidateDestination() bool {
	if trip.Destination != "" {
		return true
	}
	return false
}

func (trip Trip) ValidateDepartureTime() bool {
	if trip.DepartureTime >= 0 {
		return true
	}
	return false
}

//This is for Version 4, randomly generated, UUID's. Need to change if anything besides V4 is used in the future
func (trip Trip) ValidateUserId() bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	if r.MatchString(trip.UserId) {
		return true
	}
	return false
}
