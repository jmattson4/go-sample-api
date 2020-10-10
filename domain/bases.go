package domain

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/twinj/uuid"
)

//Base This struct can be used across multiple domains as base.
type Base struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key;type:uuid;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
	DeletedAt time.Time `sql:"index" json:"deleted_at"`
}

//BeforeCreate ...
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid.String())
}
