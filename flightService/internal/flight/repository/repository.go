package repository

import "flight_booking_system/flightService/models"

type FlightRepositoryI interface {
	Create(p *models.Flight) error
	Get(id int) (*models.Flight, error)
	Update(p *models.Flight) error
	Delete(id int) error
	GetAll() ([]*models.FlightDTO, error)
	GetAllByFlightNumber(flightNumber string) ([]*models.FlightDTO, error)
	GetAllPaginate(offset, limit int) ([]*models.FlightDTO, error)
}
