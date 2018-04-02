package models

type Estimation struct {
	Duration int64 `json:"duration"`
	Distance float64 `json:"distance"`
	Cost float64 `json:"cost"`
}