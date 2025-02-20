package postgres

import (
	"database/sql"
	flightRep "flight_booking_system/flightService/internal/flight/repository"
	"flight_booking_system/flightService/internal/testBuilders"
	"flight_booking_system/flightService/models"
	"flight_booking_system/flightService/pkg/logger"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type FlightRepoTestSuite struct {
	suite.Suite
	db            *sql.DB
	gormDB        *gorm.DB
	mock          sqlmock.Sqlmock
	repo          flightRep.FlightRepositoryI
	flightBuilder *testBuilders.FlightBuilder
}

func TestFlightRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(FlightRepoTestSuite))
}

func (s *FlightRepoTestSuite) BeforeEach(t provider.T) {
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
	s.flightBuilder = testBuilders.NewFlightBuilder()
}

func (s *FlightRepoTestSuite) AfterEach(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *FlightRepoTestSuite) TestCreateFlight(t provider.T) {
	flight := s.flightBuilder.
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "flight" ("uid","username","flight_number","price","status","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(flight.UID, flight.Username, flight.FlightNumber, flight.Price, flight.Status, flight.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectCommit()

	err := s.repo.Create(&flight)
	t.Assert().NoError(err)
	t.Assert().Equal(1, flight.ID)
}

func (s *FlightRepoTestSuite) TestGetFlight(t provider.T) {
	flight := s.flightBuilder.
		WithID(1).
		WithUID("").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	rows := sqlmock.NewRows([]string{"id", "flight_uid", "username", "flight_number", "price", "status"}).
		AddRow(
			flight.ID,
			flight.UID,
			flight.Username,
			flight.FlightNumber,
			flight.Price,
			flight.Status,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "flight" WHERE id = $1 LIMIT $2`)).
		WithArgs(flight.ID, 1).
		WillReturnRows(rows)

	resFlight, err := s.repo.Get(flight.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(flight, *resFlight)
}

func (s *FlightRepoTestSuite) TestUpdateFlight(t provider.T) {
	flight := s.flightBuilder.
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	rows := sqlmock.NewRows([]string{"id", "flight_uid", "username", "flight_number", "price", "status"}).
		AddRow(
			flight.ID,
			flight.UID,
			flight.Username,
			flight.FlightNumber,
			flight.Price,
			flight.Status,
		)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE "flight" SET "uid"=$1,"username"=$2,"flight_number"=$3,"price"=$4,"status"=$5 WHERE "id" = $6 RETURNING *`)).
		WithArgs(flight.UID, flight.Username, flight.FlightNumber, flight.Price, flight.Status, flight.ID).WillReturnRows(rows)

	s.mock.ExpectCommit()

	err := s.repo.Update(&flight)
	t.Assert().NoError(err)
}

func (s *FlightRepoTestSuite) TestDeleteFlight(t provider.T) {
	flight := s.flightBuilder.
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "flight" WHERE "flight"."id" = $1`)).
		WithArgs(flight.ID).WillReturnResult(sqlmock.NewResult(int64(flight.ID), 1))

	s.mock.ExpectCommit()

	err := s.repo.Delete(flight.ID)
	t.Assert().NoError(err)
}

func (s *FlightRepoTestSuite) TestGetAll(t provider.T) {
	flights := make([]models.Flight, 10)
	for _, flight := range flights {
		err := faker.FakeData(&flight)
		t.Assert().NoError(err)
	}

	flightsPtr := make([]*models.Flight, len(flights))
	for i, flight := range flights {
		flightsPtr[i] = &flight
	}

	rowsFlights := sqlmock.NewRows([]string{"id", "flight_uid", "username", "flight_number", "price", "status"})

	for i := range flights {
		rowsFlights.AddRow(flights[i].ID, flights[i].UID, flights[i].Username, flights[i].FlightNumber, flights[i].Price, flights[i].Status)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "flight"`)).
		WillReturnRows(rowsFlights)

	resFlights, err := s.repo.GetAll()
	t.Assert().NoError(err)
	t.Assert().Equal(flightsPtr, resFlights)
}
