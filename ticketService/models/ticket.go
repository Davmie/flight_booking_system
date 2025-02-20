package models

type Tabler interface {
	TableName() string
}

func (Ticket) TableName() string {
	return "ticket"
}

func TicketToDTO(ticket Ticket) *TicketDTO {
	return &TicketDTO{
		TicketUID:    ticket.TicketUID,
		FlightNumber: ticket.FlightNumber,
		Price:        ticket.Price,
		Status:       ticket.Status,
	}
}

type Ticket struct {
	ID           int    `json:"id" db:"id"`
	TicketUID    string `json:"ticketUid" db:"ticket_uid"`
	Username     string `json:"username" db:"username"`
	FlightNumber string `json:"flightNumber" db:"flight_number"`
	Price        int    `json:"price" db:"price"`
	Status       string `json:"status" db:"status"`
}

type TicketDTO struct {
	TicketUID    string `json:"ticketUid"`
	FlightNumber string `json:"flightNumber"`
	Price        int    `json:"price"`
	Status       string `json:"status"`
}
