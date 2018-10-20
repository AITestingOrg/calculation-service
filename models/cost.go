package models

type Cost struct {
	Origin      string  `json:"origin" bson:"origin"`
	Destination string  `json:"destination" bson:"destination"`
	Cost        float64 `json:"cost" bson:"cost"`
	UserId      string  `json:"userId" bson:"userId"`
}
