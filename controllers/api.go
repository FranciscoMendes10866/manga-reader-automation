package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/go-chi/chi/v5"
)

func GetAllMangas(w http.ResponseWriter, r *http.Request) {
	var mangas []entities.MangaEntity

	config.Database.Table("manga_entities").Find(&mangas)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mangas)
}

func GetMangaDetails(w http.ResponseWriter, r *http.Request) {
	mangaId := chi.URLParam(r, "mangaId")

	var manga entities.MangaEntity
	var chapters []entities.ChapterEntity
	var categories []entities.CategoriesEntity

	config.Database.Table("manga_entities").Where("id = ?", mangaId).First(&manga)
	config.Database.Table("chapter_entities").Where("manga_id = ?", mangaId).Find(&chapters)
	config.Database.Table("manga_categories").Where("manga_entity_id = ?", mangaId).Find(&categories)

	manga.Chapters = chapters
	manga.Categories = categories

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(manga)
}
