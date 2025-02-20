package models

import "time"

type Tabler interface {
	TableName() string
}

func (Airport) TableName() string {
	return "airport"
}

type Airport struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	City    string `json:"city" db:"city"`
	Country string `json:"country" db:"country"`
}

func (Flight) TableName() string {
	return "flight"
}

type Flight struct {
	ID            int       `json:"id" db:"id"`
	FlightNumber  string    `json:"flightNumber" db:"flight_number"`
	DateTime      time.Time `json:"dateTime" db:"datetime"`
	FromAirportID int       `json:"from_airport_id" db:"from_airport_id"`
	ToAirportID   int       `json:"to_airport_id" db:"to_airport_id"`
	Price         int       `json:"price" db:"price"`
}

type FlightDTO struct {
	FlightNumber string    `json:"flightNumber"`
	Date         time.Time `json:"date"`
	FromAirport  string    `json:"fromAirport"`
	ToAirport    string    `json:"toAirport"`
	Price        int       `json:"price"`
}
