package models

type Cost struct {
	Origin      	string	`json:"originAddress" bson:"originAddress"`
	Destination 	string  `json:"destinationAddress" bson:"destinationAddress"`
	DepartureTime 	int64 	`json:"departureTime" bson:"departureTime"`
	Distance    	float64 `json:"distance" bson:"distance"`
	Duration    	int64   `json:"duration" bson:"duration"`
	Cost        	float64 `json:"cost" bson:"cost"`
	UserId 			string  `json:"userId" bson:"userId"`
}
