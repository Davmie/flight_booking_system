package postgres

import (
	"flight_booking_system/bonusService/internal/privilege/repository"
	"flight_booking_system/bonusService/models"
	"flight_booking_system/bonusService/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type pgPrivilegeRepo struct {
	Logger logger.Logger
	DB     *gorm.DB
}

func New(logger logger.Logger, db *gorm.DB) repository.PrivilegeRepositoryI {
	return &pgPrivilegeRepo{
		Logger: logger,
		DB:     db,
	}
}

func (pr *pgPrivilegeRepo) Create(p *models.Privilege) error {
	tx := pr.DB.Create(p)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgPrivilegeRepo.Create error while inserting in repo")
	}

	return nil
}

func (pr *pgPrivilegeRepo) Get(id int) (*models.Privilege, error) {
	var p models.Privilege
	tx := pr.DB.Where("id = ?", id).Take(&p)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgPrivilegeRepo.Get error")
	}

	return &p, nil
}

func (pr *pgPrivilegeRepo) GetByUserName(username string) (*models.Privilege, error) {
	var p models.Privilege
	tx := pr.DB.Where("username = ?", username).Take(&p)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgPrivilegeRepo.Get error")
	}

	return &p, nil
}

func (pr *pgPrivilegeRepo) Update(p *models.Privilege) error {
	tx := pr.DB.Where("id = ?", p.ID).Updates(p)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgPrivilegeRepo.Update error while inserting in repo")
	}

	return nil
}

func (pr *pgPrivilegeRepo) Delete(id int) error {
	tx := pr.DB.Delete(&models.Privilege{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgPrivilegeRepo.Delete error")
	}

	return nil
}

func (pr *pgPrivilegeRepo) GetAll() ([]*models.Privilege, error) {
	var privileges []*models.Privilege

	tx := pr.DB.Find(&privileges)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgPrivilegeRepo.GetAll error")
	}

	return privileges, nil
}

func (pr *pgPrivilegeRepo) CreateHistory(p *models.PrivilegeHistory) error {
	tx := pr.DB.Create(p)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgPrivilegeRepo.Create error while inserting in repo")
	}

	return nil
}

func (pr *pgPrivilegeRepo) GetHistory(username string) ([]*models.PrivilegeHistory, error) {
	var privilege models.Privilege
	var privilegeHistory []*models.PrivilegeHistory

	tx := pr.DB.Where("username = ?", username).Take(&privilege)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgPrivilegeRepo.GetHistory error")
	}

	tx = pr.DB.Where("privilege_id = ?", privilege.ID).Find(&privilegeHistory)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgPrivilegeRepo.GetHistory error")
	}

	return privilegeHistory, nil
}
