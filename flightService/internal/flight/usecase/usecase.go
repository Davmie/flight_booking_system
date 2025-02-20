package usecase

import (
	flightRep "flight_booking_system/flightService/internal/flight/repository"
	"flight_booking_system/flightService/models"
	"github.com/pkg/errors"
)

type FlightUseCaseI interface {
	Create(p *models.Flight) error
	Get(id int) (*models.Flight, error)
	Update(p *models.Flight) error
	Delete(id int) error
	GetAll(flightNumber string) ([]*models.FlightDTO, error)
	GetAllPaginate(offset, limit int) ([]*models.FlightDTO, error)
}

type flightUseCase struct {
	flightRepository flightRep.FlightRepositoryI
}

func New(aRep flightRep.FlightRepositoryI) FlightUseCaseI {
	return &flightUseCase{
		flightRepository: aRep,
	}
}

func (pUC *flightUseCase) Create(p *models.Flight) error {
	err := pUC.flightRepository.Create(p)

	if err != nil {
		return errors.Wrap(err, "flightUseCase.Create error")
	}

	return nil
}

func (pUC *flightUseCase) Get(id int) (*models.Flight, error) {
	resFlight, err := pUC.flightRepository.Get(id)

	if err != nil {
		return nil, errors.Wrap(err, "flightUseCase.Get error")
	}

	return resFlight, nil
}

func (pUC *flightUseCase) Update(p *models.Flight) error {
	_, err := pUC.flightRepository.Get(p.ID)

	if err != nil {
		return errors.Wrap(err, "flightUseCase.Update error: Flight not found")
	}

	err = pUC.flightRepository.Update(p)

	if err != nil {
		return errors.Wrap(err, "flightUseCase.Update error: Can't update in repo")
	}

	return nil
}

func (pUC *flightUseCase) Delete(id int) error {
	_, err := pUC.flightRepository.Get(id)

	if err != nil {
		return errors.Wrap(err, "flightUseCase.Delete error: Flight not found")
	}

	err = pUC.flightRepository.Delete(id)

	if err != nil {
		return errors.Wrap(err, "flightUseCase.Delete error: Can't delete in repo")
	}

	return nil
}

func (pUC *flightUseCase) GetAll(flightNumber string) ([]*models.FlightDTO, error) {
	var flights []*models.FlightDTO
	var err error

	if flightNumber == "" {
		flights, err = pUC.flightRepository.GetAll()
	} else {
		flights, err = pUC.flightRepository.GetAllByFlightNumber(flightNumber)
	}

	if err != nil {
		return nil, errors.Wrap(err, "flightUseCase.GetAll error")
	}

	return flights, nil
}

func (pUC *flightUseCase) GetAllPaginate(offset, limit int) ([]*models.FlightDTO, error) {
	flights, err := pUC.flightRepository.GetAllPaginate(offset, limit)

	if err != nil {
		return nil, errors.Wrap(err, "flightUseCase.GetAllPaginate error")
	}

	return flights, nil
}
