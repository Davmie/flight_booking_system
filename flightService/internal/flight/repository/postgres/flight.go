package postgres

import (
	"flight_booking_system/flightService/internal/flight/repository"
	"flight_booking_system/flightService/models"
	"flight_booking_system/flightService/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pgFlightRepo struct {
	Logger logger.Logger
	DB     *gorm.DB
}

func New(logger logger.Logger, db *gorm.DB) repository.FlightRepositoryI {
	return &pgFlightRepo{
		Logger: logger,
		DB:     db,
	}
}

func (pr *pgFlightRepo) Create(p *models.Flight) error {
	tx := pr.DB.Create(p)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgFlightRepo.Create error while inserting in repo")
	}

	return nil
}

func (pr *pgFlightRepo) Get(id int) (*models.Flight, error) {
	var p models.Flight
	tx := pr.DB.Where("id = ?", id).Take(&p)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgFlightRepo.Get error")
	}

	return &p, nil
}

func (pr *pgFlightRepo) Update(p *models.Flight) error {
	tx := pr.DB.Clauses(clause.Returning{}).Omit("id").Updates(p)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgFlightRepo.Update error while inserting in repo")
	}

	return nil
}

func (pr *pgFlightRepo) Delete(id int) error {
	tx := pr.DB.Delete(&models.Flight{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgFlightRepo.Delete error")
	}

	return nil
}

func (pr *pgFlightRepo) getFlightDTOs(flights []*models.Flight) (ret []*models.FlightDTO, err error) {
	for _, flight := range flights {
		var fromA, toA models.Airport
		tx := pr.DB.Where("id = ?", flight.FromAirportID).Take(&fromA)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "pgFlightRepo.getAirports error")
		}

		tx = pr.DB.Where("id = ?", flight.ToAirportID).Take(&toA)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "pgFlightRepo.getAirports error")
		}

		ret = append(ret, &models.FlightDTO{
			FlightNumber: flight.FlightNumber,
			Date:         flight.DateTime,
			FromAirport:  fromA.City + " " + fromA.Name,
			ToAirport:    toA.City + " " + toA.Name,
			Price:        flight.Price,
		})
	}

	return ret, nil
}

func (pr *pgFlightRepo) GetAll() ([]*models.FlightDTO, error) {
	var flights []*models.Flight

	tx := pr.DB.Find(&flights)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgFlightRepo.GetAll error")
	}

	flightDTOs, err := pr.getFlightDTOs(flights)
	if err != nil {
		return nil, errors.Wrap(tx.Error, "pgFlightRepo.GetAll error")
	}

	return flightDTOs, nil
}

func (pr *pgFlightRepo) GetAllByFlightNumber(flightNumber string) ([]*models.FlightDTO, error) {
	var flights []*models.Flight

	tx := pr.DB.Where("flight_number = ?", flightNumber).Find(&flights)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgFlightRepo.GetAll error")
	}

	flightDTOs, err := pr.getFlightDTOs(flights)
	if err != nil {
		return nil, errors.Wrap(tx.Error, "pgFlightRepo.GetAll error")
	}

	return flightDTOs, nil
}

func (pr *pgFlightRepo) GetAllPaginate(offset, limit int) ([]*models.FlightDTO, error) {
	var flights []*models.Flight

	tx := pr.DB.Offset(offset).Limit(limit).Find(&flights)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgFlightRepo.GetAll error")
	}

	flightDTOs, err := pr.getFlightDTOs(flights)
	if err != nil {
		return nil, errors.Wrap(tx.Error, "pgFlightRepo.GetAll error")
	}

	return flightDTOs, nil
}
