package models

type Estimation struct {
	Origin      string  `json:"originAddress"`
	Destination string  `json:"destinationAddress"`
	Distance    float64 `json:"distance"`
	Duration    int64   `json:"duration"`
	Cost        float64 `json:"cost"`
	LastUpdated string  `json:"lastUpdated"`
	UserId 		string  `json:"userId"`
}
