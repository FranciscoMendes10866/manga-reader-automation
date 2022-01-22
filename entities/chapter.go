package entities

import (
	"github.com/FranciscoMendes10866/queues/helpers"
	"gorm.io/gorm"
)

type ChapterEntity struct {
	gorm.Model
	ID      string `gorm:"primaryKey;"`
	Name    string `json:"name" gorm:"not null;unique;index"`
	Pages   string `json:"pages" gorm:"not null"`
	MangaID string `json:"manga_id" gorm:"not null"`
}

func (chapter *ChapterEntity) BeforeCreate(db *gorm.DB) error {
	chapter.ID = helpers.GenerateID(36)
	return nil
}
