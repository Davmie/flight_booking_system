package models

import "time"

type Tabler interface {
	TableName() string
}

func (Privilege) TableName() string {
	return "privilege"
}

func PrivilegeToDTO(privilege Privilege) PrivilegeDTO {
	return PrivilegeDTO{
		ID:      privilege.ID,
		Balance: privilege.Balance,
		Status:  privilege.Status,
	}
}

type Privilege struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Status   string `json:"status" db:"status"`
	Balance  int    `json:"balance" db:"balance"`
}

type PrivilegeDTO struct {
	ID      int    `json:"id"`
	Balance int    `json:"balance"`
	Status  string `json:"status"`
}

func (PrivilegeHistory) TableName() string {
	return "privilege_history"
}

type PrivilegeHistory struct {
	ID            int       `json:"id" db:"id"`
	PrivilegeID   int       `json:"privilegeId" db:"privilege_id"`
	TicketUID     string    `json:"ticketUid" db:"ticket_uid"`
	DateTime      time.Time `json:"date" gorm:"column:datetime"`
	BalanceDiff   int       `json:"balanceDiff" db:"balance_diff"`
	OperationType string    `json:"operationType" db:"operation_type"`
}

type PrivilegeHistoryDTO struct {
	Date          time.Time `json:"date" db:"datetime"`
	TicketUID     string    `json:"ticketUid" db:"ticket_uid"`
	BalanceDiff   int       `json:"balanceDiff" db:"balance_diff"`
	OperationType string    `json:"operationType" db:"operation_type"`
}

func PrivilegeHistoryToDTO(privilegeHist PrivilegeHistory) PrivilegeHistoryDTO {
	return PrivilegeHistoryDTO{
		Date:          privilegeHist.DateTime,
		TicketUID:     privilegeHist.TicketUID,
		BalanceDiff:   privilegeHist.BalanceDiff,
		OperationType: privilegeHist.OperationType,
	}
}

func PrivilegeHistoryToDTOs(privilegeHists []*PrivilegeHistory) []PrivilegeHistoryDTO {
	res := make([]PrivilegeHistoryDTO, len(privilegeHists))
	for i, ph := range privilegeHists {
		res[i] = PrivilegeHistoryToDTO(*ph)
	}
	return res
}
