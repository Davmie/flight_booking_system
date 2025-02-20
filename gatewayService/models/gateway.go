package models

import "time"

type FlightInfo struct {
	FlightNumber string    `json:"flightNumber"`
	FromAirport  string    `json:"fromAirport"`
	ToAirport    string    `json:"toAirport"`
	Date         time.Time `json:"date"`
	Price        int       `json:"price"`
}
type FlightsInfo struct {
	Page    int          `json:"page"`
	Size    int          `json:"pageSize"`
	Total   int          `json:"totalElements"`
	Flights []FlightInfo `json:"items"`
}

type TicketInfoRequest struct {
	UID          string `json:"ticketUid"`
	FlightNumber string `json:"flightNumber"`
	Username     string `json:"username"`
	Price        int    `json:"price"`
	Status       string `json:"status"`
}

type TicketInfo struct {
	UID          string    `json:"ticketUid"`
	FlightNumber string    `json:"flightNumber"`
	FromAirport  string    `json:"fromAirport"`
	ToAirport    string    `json:"toAirport"`
	Date         time.Time `json:"date"`
	Price        int       `json:"price"`
	Status       string    `json:"status"`
}

type PrivilegeInfo struct {
	ID      int    `json:"id"`
	Balance int    `json:"balance"`
	Status  string `json:"status"`
}

type PrivilegeHistoryInfo struct {
	PrivilegeID   int       `json:"privilegeId"`
	TicketUID     string    `json:"ticketUid"`
	Date          time.Time `json:"date" db:"datetime"`
	BalanceDiff   int       `json:"balanceDiff"`
	OperationType string    `json:"operationType"`
}

type UserInfoResponse struct {
	Tickets   []TicketInfo  `json:"tickets"`
	Privilege PrivilegeInfo `json:"privilege"`
}

type BuyTicketInfo struct {
	FlightNumber    string `json:"flightNumber"`
	Price           int    `json:"price"`
	PaidFromBalance bool   `json:"paidFromBalance"`
}

type BuyTicketResponse struct {
	UID           string            `json:"ticketUid"`
	FlightNumber  string            `json:"flightNumber"`
	FromAirport   string            `json:"fromAirport"`
	ToAirport     string            `json:"toAirport"`
	Date          time.Time         `json:"date"`
	Price         int               `json:"price"`
	PaidByMoney   int               `json:"paidByMoney"`
	PaidByBonuses int               `json:"paidByBonuses"`
	Status        string            `json:"status"`
	Privilege     PrivilegeResponse `json:"privilege"`
}
