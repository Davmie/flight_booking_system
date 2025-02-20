package postgres

import (
	"database/sql"
	"flight_booking_system/ticketService/internal/testBuilders"
	ticketRep "flight_booking_system/ticketService/internal/ticket/repository"
	"flight_booking_system/ticketService/models"
	"flight_booking_system/ticketService/pkg/logger"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type TicketRepoTestSuite struct {
	suite.Suite
	db            *sql.DB
	gormDB        *gorm.DB
	mock          sqlmock.Sqlmock
	repo          ticketRep.TicketRepositoryI
	ticketBuilder *testBuilders.TicketBuilder
}

func TestTicketRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(TicketRepoTestSuite))
}

func (s *TicketRepoTestSuite) BeforeEach(t provider.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("error while creating sql mock")
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatal("error gorm open")
	}

	var logger logger.Logger

	s.db = db
	s.gormDB = gormDB
	s.mock = mock

	s.repo = New(logger, gormDB)
	s.ticketBuilder = testBuilders.NewTicketBuilder()
}

func (s *TicketRepoTestSuite) AfterEach(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *TicketRepoTestSuite) TestCreateTicket(t provider.T) {
	ticket := s.ticketBuilder.
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "ticket" ("uid","username","flight_number","price","status","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(ticket.UID, ticket.Username, ticket.FlightNumber, ticket.Price, ticket.Status, ticket.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectCommit()

	err := s.repo.Create(&ticket)
	t.Assert().NoError(err)
	t.Assert().Equal(1, ticket.ID)
}

func (s *TicketRepoTestSuite) TestGetTicket(t provider.T) {
	ticket := s.ticketBuilder.
		WithID(1).
		WithUID("").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	rows := sqlmock.NewRows([]string{"id", "ticket_uid", "username", "flight_number", "price", "status"}).
		AddRow(
			ticket.ID,
			ticket.UID,
			ticket.Username,
			ticket.FlightNumber,
			ticket.Price,
			ticket.Status,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "ticket" WHERE id = $1 LIMIT $2`)).
		WithArgs(ticket.ID, 1).
		WillReturnRows(rows)

	resTicket, err := s.repo.Get(ticket.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(ticket, *resTicket)
}

func (s *TicketRepoTestSuite) TestUpdateTicket(t provider.T) {
	ticket := s.ticketBuilder.
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	rows := sqlmock.NewRows([]string{"id", "ticket_uid", "username", "flight_number", "price", "status"}).
		AddRow(
			ticket.ID,
			ticket.UID,
			ticket.Username,
			ticket.FlightNumber,
			ticket.Price,
			ticket.Status,
		)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE "ticket" SET "uid"=$1,"username"=$2,"flight_number"=$3,"price"=$4,"status"=$5 WHERE "id" = $6 RETURNING *`)).
		WithArgs(ticket.UID, ticket.Username, ticket.FlightNumber, ticket.Price, ticket.Status, ticket.ID).WillReturnRows(rows)

	s.mock.ExpectCommit()

	err := s.repo.Update(&ticket)
	t.Assert().NoError(err)
}

func (s *TicketRepoTestSuite) TestDeleteTicket(t provider.T) {
	ticket := s.ticketBuilder.
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "ticket" WHERE "ticket"."id" = $1`)).
		WithArgs(ticket.ID).WillReturnResult(sqlmock.NewResult(int64(ticket.ID), 1))

	s.mock.ExpectCommit()

	err := s.repo.Delete(ticket.ID)
	t.Assert().NoError(err)
}

func (s *TicketRepoTestSuite) TestGetAll(t provider.T) {
	tickets := make([]models.Ticket, 10)
	for _, ticket := range tickets {
		err := faker.FakeData(&ticket)
		t.Assert().NoError(err)
	}

	ticketsPtr := make([]*models.Ticket, len(tickets))
	for i, ticket := range tickets {
		ticketsPtr[i] = &ticket
	}

	rowsTickets := sqlmock.NewRows([]string{"id", "ticket_uid", "username", "flight_number", "price", "status"})

	for i := range tickets {
		rowsTickets.AddRow(tickets[i].ID, tickets[i].UID, tickets[i].Username, tickets[i].FlightNumber, tickets[i].Price, tickets[i].Status)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "ticket"`)).
		WillReturnRows(rowsTickets)

	resTickets, err := s.repo.GetAll()
	t.Assert().NoError(err)
	t.Assert().Equal(ticketsPtr, resTickets)
}
