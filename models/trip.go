package models

import (
	"errors"
	"fmt"
)

type Trip struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureTime int64  `json:"departureTime"`
	UserId        string `json:"userId"`
}

func (trip Trip) ValidateFields(fields ...string) error {
	invalidFields := ""
	for _, field := range fields {
		switch field {
		case "origin":
			if !trip.ValidateOrigin() {
				invalidFields += fmt.Sprintf("Invalid origin.\n\tGiven: %s\n\tExpected: Non Empty String\n", trip.Origin)
			}
		case "destination":
			if !trip.ValidateDestination() {
				invalidFields += fmt.Sprintf("Invalid destination.\n\tGiven: %s\n\tExpected: Non Empty String\n", trip.Destination)
			}
		case "departureTime":
			if !trip.ValidateDepartureTime() {
				invalidFields += fmt.Sprintf("Invalid departureTime.\n\tGiven: %s\n\tExpected: Valid int64 >= 0\n", trip.DepartureTime)
			}
		case "userId":
			if !trip.ValidateUserId() {
				invalidFields += fmt.Sprintf("Invalid userId.\n\tGiven: %s\n\tExpected: Non-empty UUID\n", trip.UserId)
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

func (trip Trip) ValidateUserId() bool {
	if trip.UserId != "" {
		return true
	}
	return false
}
