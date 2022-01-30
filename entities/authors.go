package entities

import (
	"github.com/FranciscoMendes10866/queues/helpers"
	"gorm.io/gorm"
)

type AuthorsEntity struct {
	gorm.Model
	ID     string        `gorm:"primaryKey;"`
	Name   string        `json:"name" gorm:"not null"`
	Mangas []MangaEntity `gorm:"foreignKey:AuthorID"`
}

func (author *AuthorsEntity) BeforeCreate(db *gorm.DB) error {
	author.ID = helpers.GenerateID(36)
	return nil
}
