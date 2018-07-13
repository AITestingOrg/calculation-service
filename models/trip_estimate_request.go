package models

type TripEstimateRequest struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureTime int64  `json:"departureTime"`
	UserId   	string `json:"userId"`
}
