package models

type Cost struct {
	Origin      	string	`json:"originAddress" bson:"originAddress"`
	Destination 	string  `json:"destinationAddress" bson:"destinationAddress"`
	DepartureTime 	int64 	`json:"departureTime" bson:"departureTime"`
	Cost        	float64 `json:"cost" bson:"cost"`
	UserId 			string  `json:"userId" bson:"userId"`
}
