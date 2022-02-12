package entities

import (
	"github.com/FranciscoMendes10866/queues/helpers"
	"gorm.io/gorm"
)

type MangaEntity struct {
	gorm.Model
	ID          string             `gorm:"primaryKey;"`
	Name        string             `json:"name" gorm:"not null;unique_index;index"`
	Thumbnail   string             `json:"thumbnail" gorm:"not null"`
	Description string             `json:"description" gorm:"not null"`
	Chapters    []ChapterEntity    `json:"chapters" gorm:"foreignKey:MangaID"`
	Categories  []CategoriesEntity `json:"categories" gorm:"many2many:manga_categories;"`
}

func (manga *MangaEntity) BeforeCreate(db *gorm.DB) error {
	manga.ID = helpers.GenerateID(36)
	return nil
}
