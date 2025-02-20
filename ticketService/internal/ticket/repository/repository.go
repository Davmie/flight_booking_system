package repository

import "flight_booking_system/ticketService/models"

type TicketRepositoryI interface {
	Create(p *models.Ticket) error
	Get(id int) (*models.Ticket, error)
	Update(p *models.Ticket) error
	Delete(id int) error
	GetAll() ([]*models.Ticket, error)
	GetAllByUserName(userName string) ([]*models.Ticket, error)
}
