package entities

import (
	"github.com/FranciscoMendes10866/queues/helpers"
	"gorm.io/gorm"
)

type PagesEntity struct {
	gorm.Model
	ID        string `gorm:"primaryKey;"`
	URL       string `json:"url" gorm:"not null"`
	ChapterID string `json:"chapter_id" gorm:"not null"`
}

func (page *PagesEntity) BeforeCreate(db *gorm.DB) error {
	page.ID = helpers.GenerateID(36)
	return nil
}
