package postgres

import (
	"flight_booking_system/ticketService/internal/ticket/repository"
	"flight_booking_system/ticketService/models"
	"flight_booking_system/ticketService/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type pgTicketRepo struct {
	Logger logger.Logger
	DB     *gorm.DB
}

func New(logger logger.Logger, db *gorm.DB) repository.TicketRepositoryI {
	return &pgTicketRepo{
		Logger: logger,
		DB:     db,
	}
}

func (pr *pgTicketRepo) Create(p *models.Ticket) error {
	tx := pr.DB.Create(p)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgTicketRepo.Create error while inserting in repo")
	}

	return nil
}

func (pr *pgTicketRepo) Get(id int) (*models.Ticket, error) {
	var p models.Ticket
	tx := pr.DB.Where("id = ?", id).Take(&p)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgTicketRepo.Get error")
	}

	return &p, nil
}

func (pr *pgTicketRepo) Update(p *models.Ticket) error {
	tx := pr.DB.Where("ticket_uid = ?", p.TicketUID).Omit("id", "ticket_uid").Updates(p)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgTicketRepo.Update error while inserting in repo")
	}

	return nil
}

func (pr *pgTicketRepo) Delete(id int) error {
	tx := pr.DB.Delete(&models.Ticket{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgTicketRepo.Delete error")
	}

	return nil
}

func (pr *pgTicketRepo) GetAll() ([]*models.Ticket, error) {
	var tickets []*models.Ticket

	tx := pr.DB.Find(&tickets)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgTicketRepo.GetAll error")
	}

	return tickets, nil
}

func (pr *pgTicketRepo) GetAllByUserName(userName string) ([]*models.Ticket, error) {
	var tickets []*models.Ticket

	tx := pr.DB.Find(&tickets, "username = ?", userName)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgTicketRepo.GetAll error")
	}

	return tickets, nil
}
