package entities

import "gorm.io/gorm"

type ChapterEntity struct {
	gorm.Model
	Name    string `json:"name" gorm:"not null;unique;index"`
	Pages   string `json:"pages" gorm:"not null"`
	MangaID uint   `json:"manga_id" gorm:"not null"`
}
