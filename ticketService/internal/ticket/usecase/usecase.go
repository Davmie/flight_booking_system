package usecase

import (
	ticketRep "flight_booking_system/ticketService/internal/ticket/repository"
	"flight_booking_system/ticketService/models"
	"github.com/pkg/errors"
)

type TicketUseCaseI interface {
	Create(p *models.Ticket) error
	Get(id int) (*models.Ticket, error)
	Update(p *models.Ticket) error
	Delete(id int) error
	GetAll(userName string) ([]*models.Ticket, error)
	GetByUID(ticketUid string, userName string) (*models.Ticket, error)
}

type ticketUseCase struct {
	ticketRepository ticketRep.TicketRepositoryI
}

func New(aRep ticketRep.TicketRepositoryI) TicketUseCaseI {
	return &ticketUseCase{
		ticketRepository: aRep,
	}
}

func (pUC *ticketUseCase) Create(p *models.Ticket) error {
	err := pUC.ticketRepository.Create(p)

	if err != nil {
		return errors.Wrap(err, "ticketUseCase.Create error")
	}

	return nil
}

func (pUC *ticketUseCase) Get(id int) (*models.Ticket, error) {
	resTicket, err := pUC.ticketRepository.Get(id)

	if err != nil {
		return nil, errors.Wrap(err, "ticketUseCase.Get error")
	}

	return resTicket, nil
}

func (pUC *ticketUseCase) Update(p *models.Ticket) error {
	//_, err := pUC.ticketRepository.Get(p.ID)
	//
	//if err != nil {
	//	return errors.Wrap(err, "ticketUseCase.Update error: Ticket not found")
	//}

	err := pUC.ticketRepository.Update(p)

	if err != nil {
		return errors.Wrap(err, "ticketUseCase.Update error: Can't update in repo")
	}

	return nil
}

func (pUC *ticketUseCase) Delete(id int) error {
	_, err := pUC.ticketRepository.Get(id)

	if err != nil {
		return errors.Wrap(err, "ticketUseCase.Delete error: Ticket not found")
	}

	err = pUC.ticketRepository.Delete(id)

	if err != nil {
		return errors.Wrap(err, "ticketUseCase.Delete error: Can't delete in repo")
	}

	return nil
}

func (pUC *ticketUseCase) GetAll(userName string) ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	var err error

	if userName == "" {
		tickets, err = pUC.ticketRepository.GetAll()
	} else {
		tickets, err = pUC.ticketRepository.GetAllByUserName(userName)
	}

	if err != nil {
		return nil, errors.Wrap(err, "ticketUseCase.GetAll error")
	}

	return tickets, nil
}

func (pUC *ticketUseCase) GetByUID(ticketUid string, userName string) (*models.Ticket, error) {
	tickets, err := pUC.ticketRepository.GetAllByUserName(userName)
	if err != nil {
		return nil, errors.Wrap(err, "ticketUseCase.GetByUID error")
	}

	for _, ticket := range tickets {
		if ticket.TicketUID == ticketUid {
			return ticket, nil
		}
	}

	return nil, errors.New("ticketUseCase.GetByUID ticket not found")
}
