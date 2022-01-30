package entities

import (
	"github.com/FranciscoMendes10866/queues/helpers"
	"gorm.io/gorm"
)

type ChapterEntity struct {
	gorm.Model
	ID      string        `gorm:"primaryKey;"`
	Name    string        `json:"name" gorm:"not null;unique;index"`
	MangaID string        `json:"manga_id" gorm:"not null"`
	Pages   []PagesEntity `gorm:"foreignKey:ChapterID"`
}

func (chapter *ChapterEntity) BeforeCreate(db *gorm.DB) error {
	chapter.ID = helpers.GenerateID(36)
	return nil
}
