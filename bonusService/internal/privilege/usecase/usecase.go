package usecase

import (
	privilegeRep "flight_booking_system/bonusService/internal/privilege/repository"
	"flight_booking_system/bonusService/models"
	"github.com/pkg/errors"
)

type PrivilegeUseCaseI interface {
	Create(p *models.Privilege) error
	Get(id int) (*models.Privilege, error)
	GetByUserName(username string) (*models.Privilege, error)
	Update(p *models.Privilege) error
	Delete(id int) error
	GetAll() ([]*models.Privilege, error)
	GetHistory(username string) ([]*models.PrivilegeHistory, error)
	CreateHistory(p *models.PrivilegeHistory) error
}

type privilegeUseCase struct {
	privilegeRepository privilegeRep.PrivilegeRepositoryI
}

func New(aRep privilegeRep.PrivilegeRepositoryI) PrivilegeUseCaseI {
	return &privilegeUseCase{
		privilegeRepository: aRep,
	}
}

func (pUC *privilegeUseCase) Create(p *models.Privilege) error {
	err := pUC.privilegeRepository.Create(p)

	if err != nil {
		return errors.Wrap(err, "privilegeUseCase.Create error")
	}

	return nil
}

func (pUC *privilegeUseCase) Get(id int) (*models.Privilege, error) {
	resPrivilege, err := pUC.privilegeRepository.Get(id)

	if err != nil {
		return nil, errors.Wrap(err, "privilegeUseCase.Get error")
	}

	return resPrivilege, nil
}

func (pUC *privilegeUseCase) GetByUserName(username string) (*models.Privilege, error) {
	resPrivilege, err := pUC.privilegeRepository.GetByUserName(username)

	if err != nil {
		return nil, errors.Wrap(err, "privilegeUseCase.GetByUserName error")
	}

	return resPrivilege, nil
}

func (pUC *privilegeUseCase) Update(p *models.Privilege) error {
	err := pUC.privilegeRepository.Update(p)

	if err != nil {
		return errors.Wrap(err, "privilegeUseCase.Update error: Can't update in repo")
	}

	return nil
}

func (pUC *privilegeUseCase) Delete(id int) error {
	_, err := pUC.privilegeRepository.Get(id)

	if err != nil {
		return errors.Wrap(err, "privilegeUseCase.Delete error: Privilege not found")
	}

	err = pUC.privilegeRepository.Delete(id)

	if err != nil {
		return errors.Wrap(err, "privilegeUseCase.Delete error: Can't delete in repo")
	}

	return nil
}

func (pUC *privilegeUseCase) GetAll() ([]*models.Privilege, error) {
	privileges, err := pUC.privilegeRepository.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "privilegeUseCase.GetAll error")
	}

	return privileges, nil
}

func (pUC *privilegeUseCase) CreateHistory(p *models.PrivilegeHistory) error {
	err := pUC.privilegeRepository.CreateHistory(p)

	if err != nil {
		return errors.Wrap(err, "privilegeUseCase.CreateHistory error")
	}

	return nil
}

func (pUC *privilegeUseCase) GetHistory(username string) ([]*models.PrivilegeHistory, error) {
	privilegeHists, err := pUC.privilegeRepository.GetHistory(username)
	if err != nil {
		return nil, errors.Wrap(err, "privilegeUseCase.GetHistory error")
	}
	return privilegeHists, nil
}
