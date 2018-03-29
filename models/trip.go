package models

type Trip struct {
	Origin string `json:"origin"`
	Destination string `json:"destination"`
	DepartureTime int64 `json:"departure_time"`
}