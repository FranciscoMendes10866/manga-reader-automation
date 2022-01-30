package entities

import (
	"github.com/FranciscoMendes10866/queues/helpers"
	"gorm.io/gorm"
)

type CategoriesEntity struct {
	gorm.Model
	ID   string `gorm:"primaryKey;"`
	Name string `json:"name" gorm:"not null;unique;index"`
}

func (category *CategoriesEntity) BeforeCreate(db *gorm.DB) error {
	category.ID = helpers.GenerateID(36)
	return nil
}
