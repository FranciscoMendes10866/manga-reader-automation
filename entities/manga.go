package entities

import (
	"github.com/FranciscoMendes10866/queues/helpers"
	"gorm.io/gorm"
)

type MangaEntity struct {
	gorm.Model
	ID          string             `gorm:"primaryKey;"`
	Name        string             `json:"name" gorm:"not null;unique;index"`
	Thumbnail   string             `json:"thumbnail" gorm:"not null"`
	Description string             `json:"description" gorm:"not null"`
	Chapters    []ChapterEntity    `gorm:"foreignKey:MangaID"`
	Categories  []CategoriesEntity `gorm:"many2many:manga_categories;"`
}

func (manga *MangaEntity) BeforeCreate(db *gorm.DB) error {
	manga.ID = helpers.GenerateID(36)
	return nil
}
