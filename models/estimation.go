package models

import (
	"errors"
	"fmt"
	"time"
)

type Estimation struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Distance    float64 `json:"distance"`
	Duration    int64   `json:"duration"`
	Cost        float64 `json:"cost"`
	LastUpdated string  `json:"lastUpdated"`
	UserId      string  `json:"userId"`
}

func (estimation Estimation) ValidateFields(fields ...string) error {
	invalidFields := ""
	for _, field := range fields {
		switch field {
		case "origin":
			if !estimation.ValidateOrigin() {
				invalidFields += fmt.Sprintf("Invalid origin.\n\tGiven: %s\n\tExpected: Non Empty String\n", estimation.Origin)
			}
		case "destination":
			if !estimation.ValidateDestination() {
				invalidFields += fmt.Sprintf("Invalid destination.\n\tGiven: %s\n\tExpected: Non Empty String\n", estimation.Destination)
			}
		case "distance":
			if !estimation.ValidateDistance() {
				invalidFields += fmt.Sprintf("Invalid distance.\n\tGiven: %f\n\tExpected: Float64 zero or greater\n", estimation.Distance)
			}
		case "duration":
			if !estimation.ValidateDuration() {
				invalidFields += fmt.Sprintf("Invalid duration.\n\tGiven: %d\n\tExpected: Integer zero or greater\n", estimation.Duration)
			}
		case "cost":
			if !estimation.ValidateCost() {
				invalidFields += fmt.Sprintf("Invalid cost.\n\tGiven: %f\n\tExpected: Float64 zero or greater\n", estimation.Cost)
			}
		case "lastUpdated":
			if !estimation.ValidateLastUpdated() {
				invalidFields += fmt.Sprintf("Invalid lastUpdated.\n\tGiven: %s\n\tExpected: Valid date in the format \"yyyy-mm-dd HH:MM:SS\"\n", estimation.LastUpdated)
			}
		case "userId":
			if !estimation.ValidateUserId() {
				invalidFields += fmt.Sprintf("Invalid userId.\n\tGiven: %s\n\tExpected: Non-empty UUID\n", estimation.UserId)
			}
		}
	}
	if invalidFields != "" {
		return errors.New(invalidFields)
	}
	return nil
}

func (estimation Estimation) ValidateOrigin() bool {
	if estimation.Origin != "" {
		return true
	}
	return false
}

func (estimation Estimation) ValidateDestination() bool {
	if estimation.Destination != "" {
		return true
	}
	return false
}

func (estimation Estimation) ValidateDistance() bool {
	if estimation.Distance >= 0 {
		return true
	}
	return false
}

func (estimation Estimation) ValidateDuration() bool {
	if estimation.Duration >= 0 {
		return true
	}
	return false
}

func (estimation Estimation) ValidateCost() bool {
	if estimation.Cost >= 0 {
		return true
	}
	return false
}

func (estimation Estimation) ValidateLastUpdated() bool {
	if estimation.LastUpdated != "" {
		_, err := time.Parse("2006-01-02 03:04:05", estimation.LastUpdated)
		if err == nil {
			return true
		}
	}
	return false
}

func (estimation Estimation) ValidateUserId() bool {
	if estimation.UserId != "" {
		return true
	}
	return false
}
