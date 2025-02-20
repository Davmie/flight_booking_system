package models

import "time"

type FlightResponse struct {
	FlightNumber string
	FromAirport  string
	ToAirport    string
	Date         time.Time
	Price        int
}
