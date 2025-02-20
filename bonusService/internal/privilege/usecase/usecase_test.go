package usecase

import (
	"flight_booking_system/ticketService/internal/testBuilders"
	ticketMocks "flight_booking_system/ticketService/internal/ticket/repository/mocks"
	"flight_booking_system/ticketService/models"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"testing"
)

type TicketTestSuite struct {
	suite.Suite
	uc             TicketUseCaseI
	ticketRepoMock *ticketMocks.TicketRepositoryI
	ticketBuilder  *testBuilders.TicketBuilder
}

func TestTicketTestSuite(t *testing.T) {
	suite.RunSuite(t, new(TicketTestSuite))
}

func (s *TicketTestSuite) BeforeEach(t provider.T) {
	s.ticketRepoMock = ticketMocks.NewTicketRepositoryI(t)
	s.uc = New(s.ticketRepoMock)
	s.ticketBuilder = testBuilders.NewTicketBuilder()
}

func (s *TicketTestSuite) TestCreateTicket(t provider.T) {
	ticket := s.ticketBuilder.WithID(1).
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	s.ticketRepoMock.On("Create", &ticket).Return(nil)
	err := s.uc.Create(&ticket)

	t.Assert().NoError(err)
	t.Assert().Equal(ticket.ID, 1)
}

func (s *TicketTestSuite) TestUpdateTicket(t provider.T) {
	var err error
	ticket := s.ticketBuilder.WithID(1).
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	notFoundTicket := s.ticketBuilder.WithID(0).Build()

	s.ticketRepoMock.On("Get", ticket.ID).Return(&ticket, nil)
	s.ticketRepoMock.On("Update", &ticket).Return(nil)
	s.ticketRepoMock.On("Get", notFoundTicket.ID).Return(&notFoundTicket, errors.Wrap(err, "Ticket not found"))
	s.ticketRepoMock.On("Update", &notFoundTicket).Return(errors.Wrap(err, "Ticket not found"))

	cases := map[string]struct {
		ArgData *models.Ticket
		Error   error
	}{
		"success": {
			ArgData: &ticket,
			Error:   nil,
		},
		"Ticket not found": {
			ArgData: &notFoundTicket,
			Error:   errors.Wrap(err, "Ticket not found"),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.Update(test.ArgData)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *TicketTestSuite) TestGetTicket(t provider.T) {
	ticket := s.ticketBuilder.WithID(1).
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	s.ticketRepoMock.On("Get", ticket.ID).Return(&ticket, nil)
	result, err := s.uc.Get(ticket.ID)

	t.Assert().NoError(err)
	t.Assert().Equal(&ticket, result)
}

func (s *TicketTestSuite) TestDeleteTicket(t provider.T) {
	var err error
	ticket := s.ticketBuilder.WithID(1).
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	notFoundTicket := s.ticketBuilder.WithID(0).Build()

	s.ticketRepoMock.On("Get", ticket.ID).Return(&ticket, nil)
	s.ticketRepoMock.On("Delete", ticket.ID).Return(nil)
	s.ticketRepoMock.On("Get", notFoundTicket.ID).Return(&notFoundTicket, errors.Wrap(err, "Ticket not found"))
	s.ticketRepoMock.On("Delete", notFoundTicket.ID).Return(errors.Wrap(err, "Ticket not found"))

	cases := map[string]struct {
		TicketID int
		Error    error
	}{
		"success": {
			TicketID: ticket.ID,
			Error:    nil,
		},
		"Ticket not found": {
			TicketID: notFoundTicket.ID,
			Error:    errors.Wrap(err, "Ticket not found"),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.Delete(test.TicketID)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *TicketTestSuite) TestGetAll(t provider.T) {
	tickets := make([]models.Ticket, 0, 10)
	err := faker.FakeData(&tickets)
	t.Assert().NoError(err)

	ticketsPtr := make([]*models.Ticket, len(tickets))
	for i, ticket := range tickets {
		ticketsPtr[i] = &ticket
	}

	s.ticketRepoMock.On("GetAll").Return(ticketsPtr, nil)

	cases := map[string]struct {
		Tickets []models.Ticket
		Error   error
	}{
		"success": {
			Tickets: tickets,
			Error:   nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			resTickets, err := s.uc.GetAll()
			t.Assert().ErrorIs(err, test.Error)
			t.Assert().Equal(ticketsPtr, resTickets)
		})
	}
}
