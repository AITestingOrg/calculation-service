package models

type Trip struct {
	Origin string `json:"origin"`
	Destination string `json:"destination"`
	DepartureTime int64 `json:"departureTime"`
}

func (trip Trip) ValidateOrigin(origin string) bool {
	if origin != "" {
		return true
	}
	return false
}

func (trip Trip) ValidateDestination(destination string) bool {
	if destination != "" {
		return true
	}
	return false
}