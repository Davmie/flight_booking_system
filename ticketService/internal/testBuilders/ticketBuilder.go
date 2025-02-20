package testBuilders

import (
	"flight_booking_system/ticketService/models"
)

type TicketBuilder struct {
	ticket models.Ticket
}

func NewTicketBuilder() *TicketBuilder {
	return &TicketBuilder{}
}

func (b *TicketBuilder) WithID(id int) *TicketBuilder {
	b.ticket.ID = id
	return b
}

func (b *TicketBuilder) WithUID(uid string) *TicketBuilder {
	b.ticket.TicketUID = uid
	return b
}

func (b *TicketBuilder) WithUsername(name string) *TicketBuilder {
	b.ticket.Username = name
	return b
}

func (b *TicketBuilder) WithFlightNumber(flightNumber string) *TicketBuilder {
	b.ticket.FlightNumber = flightNumber
	return b
}

func (b *TicketBuilder) WithPrice(price int) *TicketBuilder {
	b.ticket.Price = price
	return b
}

func (b *TicketBuilder) WithStatus(status string) *TicketBuilder {
	b.ticket.Status = status
	return b
}

func (b *TicketBuilder) Build() models.Ticket {
	return b.ticket
}
