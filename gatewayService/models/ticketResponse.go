package models

type TicketResponse struct {
	TicketUID    string `json:"ticketUid"`
	FlightNumber string `json:"flightNumber"`
	Price        int    `json:"price"`
	Status       string `json:"status"`
}
