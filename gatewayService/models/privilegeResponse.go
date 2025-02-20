package models

import "time"

type PrivilegeResponse struct {
	ID      int    `json:"id"`
	Balance int    `json:"balance"`
	Status  string `json:"status"`
}

type PrivilegeHistoryResponse struct {
	Date          time.Time `json:"date" db:"datetime"`
	TicketUID     string    `json:"ticketUid" db:"ticket_uid"`
	BalanceDiff   int       `json:"balanceDiff" db:"balance_diff"`
	OperationType string    `json:"operationType" db:"operation_type"`
}

type PrivilegeFullResponse struct {
	Balance int                        `json:"balance"`
	Status  string                     `json:"status"`
	History []PrivilegeHistoryResponse `json:"history"`
}
