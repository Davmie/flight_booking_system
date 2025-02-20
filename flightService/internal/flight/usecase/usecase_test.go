package usecase

import (
	flightMocks "flight_booking_system/flightService/internal/flight/repository/mocks"
	"flight_booking_system/flightService/internal/testBuilders"
	"flight_booking_system/flightService/models"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"testing"
)

type FlightTestSuite struct {
	suite.Suite
	uc             FlightUseCaseI
	flightRepoMock *flightMocks.FlightRepositoryI
	flightBuilder  *testBuilders.FlightBuilder
}

func TestFlightTestSuite(t *testing.T) {
	suite.RunSuite(t, new(FlightTestSuite))
}

func (s *FlightTestSuite) BeforeEach(t provider.T) {
	s.flightRepoMock = flightMocks.NewFlightRepositoryI(t)
	s.uc = New(s.flightRepoMock)
	s.flightBuilder = testBuilders.NewFlightBuilder()
}

func (s *FlightTestSuite) TestCreateFlight(t provider.T) {
	flight := s.flightBuilder.WithID(1).
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	s.flightRepoMock.On("Create", &flight).Return(nil)
	err := s.uc.Create(&flight)

	t.Assert().NoError(err)
	t.Assert().Equal(flight.ID, 1)
}

func (s *FlightTestSuite) TestUpdateFlight(t provider.T) {
	var err error
	flight := s.flightBuilder.WithID(1).
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	notFoundFlight := s.flightBuilder.WithID(0).Build()

	s.flightRepoMock.On("Get", flight.ID).Return(&flight, nil)
	s.flightRepoMock.On("Update", &flight).Return(nil)
	s.flightRepoMock.On("Get", notFoundFlight.ID).Return(&notFoundFlight, errors.Wrap(err, "Flight not found"))
	s.flightRepoMock.On("Update", &notFoundFlight).Return(errors.Wrap(err, "Flight not found"))

	cases := map[string]struct {
		ArgData *models.Flight
		Error   error
	}{
		"success": {
			ArgData: &flight,
			Error:   nil,
		},
		"Flight not found": {
			ArgData: &notFoundFlight,
			Error:   errors.Wrap(err, "Flight not found"),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.Update(test.ArgData)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *FlightTestSuite) TestGetFlight(t provider.T) {
	flight := s.flightBuilder.WithID(1).
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	s.flightRepoMock.On("Get", flight.ID).Return(&flight, nil)
	result, err := s.uc.Get(flight.ID)

	t.Assert().NoError(err)
	t.Assert().Equal(&flight, result)
}

func (s *FlightTestSuite) TestDeleteFlight(t provider.T) {
	var err error
	flight := s.flightBuilder.WithID(1).
		WithID(1).
		WithUID("uid").
		WithUsername("username").
		WithFlightNumber("flightNumber").
		WithPrice(20).
		WithStatus("status").
		Build()

	notFoundFlight := s.flightBuilder.WithID(0).Build()

	s.flightRepoMock.On("Get", flight.ID).Return(&flight, nil)
	s.flightRepoMock.On("Delete", flight.ID).Return(nil)
	s.flightRepoMock.On("Get", notFoundFlight.ID).Return(&notFoundFlight, errors.Wrap(err, "Flight not found"))
	s.flightRepoMock.On("Delete", notFoundFlight.ID).Return(errors.Wrap(err, "Flight not found"))

	cases := map[string]struct {
		FlightID int
		Error    error
	}{
		"success": {
			FlightID: flight.ID,
			Error:    nil,
		},
		"Flight not found": {
			FlightID: notFoundFlight.ID,
			Error:    errors.Wrap(err, "Flight not found"),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.Delete(test.FlightID)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *FlightTestSuite) TestGetAll(t provider.T) {
	flights := make([]models.Flight, 0, 10)
	err := faker.FakeData(&flights)
	t.Assert().NoError(err)

	flightsPtr := make([]*models.Flight, len(flights))
	for i, flight := range flights {
		flightsPtr[i] = &flight
	}

	s.flightRepoMock.On("GetAll").Return(flightsPtr, nil)

	cases := map[string]struct {
		Flights []models.Flight
		Error   error
	}{
		"success": {
			Flights: flights,
			Error:   nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			resFlights, err := s.uc.GetAll()
			t.Assert().ErrorIs(err, test.Error)
			t.Assert().Equal(flightsPtr, resFlights)
		})
	}
}
