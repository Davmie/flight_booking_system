package repository

import "flight_booking_system/bonusService/models"

type PrivilegeRepositoryI interface {
	Create(p *models.Privilege) error
	Get(id int) (*models.Privilege, error)
	GetByUserName(username string) (*models.Privilege, error)
	Update(p *models.Privilege) error
	Delete(id int) error
	GetAll() ([]*models.Privilege, error)
	CreateHistory(p *models.PrivilegeHistory) error
	GetHistory(username string) ([]*models.PrivilegeHistory, error)
}
