package entities

import "gorm.io/gorm"

type MangaEntity struct {
	gorm.Model
	Name        string          `json:"name" gorm:"not null;unique;index"`
	Thumbnail   string          `json:"thumbnail" gorm:"not null"`
	Description string          `json:"description" gorm:"not null"`
	Chapters    []ChapterEntity `gorm:"foreignKey:MangaID"`
}
