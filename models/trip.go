package models

type Trip struct {
	Origin string `json:"origin"`
	Destination string `json:"destination"`
	DepartureTime int64 `json:"departure_time"`
}

// func (trip Trip) ValidateOrigin(origin string) bool {
// 	return true
// }