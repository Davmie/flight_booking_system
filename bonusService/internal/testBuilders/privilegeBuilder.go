package testBuilders

import (
	"flight_booking_system/bonusService/models"
)

type PrivilegeBuilder struct {
	privilege models.Privilege
}

func NewPrivilegeBuilder() *PrivilegeBuilder {
	return &PrivilegeBuilder{}
}

func (b *PrivilegeBuilder) WithID(id int) *PrivilegeBuilder {
	b.privilege.ID = id
	return b
}

func (b *PrivilegeBuilder) WithUsername(name string) *PrivilegeBuilder {
	b.privilege.Username = name
	return b
}

func (b *PrivilegeBuilder) WithStatus(status string) *PrivilegeBuilder {
	b.privilege.Status = status
	return b
}

func (b *PrivilegeBuilder) WithBalance(balance int) *PrivilegeBuilder {
	b.privilege.Balance = balance
	return b
}

func (b *PrivilegeBuilder) Build() models.Privilege {
	return b.privilege
}
