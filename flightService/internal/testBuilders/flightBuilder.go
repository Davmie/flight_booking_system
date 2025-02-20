package testBuilders

import (
	"flight_booking_system/flightService/models"
)

type FlightBuilder struct {
	flight models.Flight
}

func NewFlightBuilder() *FlightBuilder {
	return &FlightBuilder{}
}

func (b *FlightBuilder) WithID(id int) *FlightBuilder {
	b.flight.ID = id
	return b
}

func (b *FlightBuilder) WithFlightNumber(flightNumber string) *FlightBuilder {
	b.flight.FlightNumber = flightNumber
	return b
}

func (b *FlightBuilder) WithDateTime(dateTime string) *FlightBuilder {
	b.flight.DateTime = dateTime
	return b
}

func (b *FlightBuilder) WithFromAirportID(fromAirportID int) *FlightBuilder {
	b.flight.FromAirportID = fromAirportID
	return b
}

func (b *FlightBuilder) WithToAirportID(toAirportID int) *FlightBuilder {
	b.flight.ToAirportID = toAirportID
	return b
}

func (b *FlightBuilder) WithPrice(price int) *FlightBuilder {
	b.flight.Price = price
	return b
}

func (b *FlightBuilder) Build() models.Flight {
	return b.flight
}
